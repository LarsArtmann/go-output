# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.17.1] - 2026-06-22

Patch release. Two motivations: (1) ship a real user-facing fix for the NOM
inline-renderer repetition bug, and (2) close out v0.17.0 with the lockstep
dependency realignment that should have shipped inside v0.17.0 itself. The
latter makes the v0.17.x dependency graph self-consistent — consumers pulling
`nom@v0.17.0` previously got a build graph where sibling modules pinned older
roots. No public API changes; drop-in upgrade from v0.17.0.

### Fixed — NOM inline renderer

- **Eliminate the inline-renderer repetition bug.** The renderer repainted the
  full tree every 200ms tick even when nothing changed; over terminals that
  lack synchronized-output support (mode 2026), cursor-up/clear-line codes
  fail to overwrite, causing each tick to append a full copy of the frame.
  Five coordinated fixes:
  - **Frame diffing** — `Draw()` compares the new frame against `lastFrame`
    and emits zero bytes when identical (mirrors bubbletea v2's
    `cursedRenderer.viewEquals()` early-exit).
  - **Writer-aware TTY detection** — `detectNoColor()` / `detectPlainText()`
    now probe the writer's own FD via `x/term.IsTerminal`, not the hardcoded
    `os.Stdout`. Fixes the stdout/stderr mismatch where BuildFlow renders to
    stderr but TTY was detected on stdout.
  - **Sync-output gating** — `\033[?2026h` / `\033[?2026l` codes are emitted
    only when `writerIsTTY` is confirmed; non-TTY writers (pipes, buffers)
    get cursor-up/clear-line but no sync wrapping.
  - **Panic recovery** — `refreshLoop` carries a deferred `recover` that
    restores the cursor on crash.
  - **SIGWINCH handling** — a `listenForResize` goroutine invalidates the
    frame cache on terminal resize, forcing a full redraw.
  All hand-rolled `\033[...]` magic strings replaced with canonical constants
  from `github.com/charmbracelet/x/ansi`. Six new behavioral tests added
  (frame diffing, sync-code gating, SIGWINCH invalidation, panic recovery,
  concurrent safety); golden files regenerated.
- **`writerIsTTY` is now authoritative.** `SetPlainText(false)` on a non-TTY
  writer previously created an impossible state: cursor-up/clear-line codes
  emitted to a pipe. `snapshotConfig()` now computes
  `effectivePlainText = plainText || !writerIsTTY`. A non-TTY writer always
  degrades to plain text; the only valid direction is TTY→plain, never
  pipe→inline. Test helpers set `writerIsTTY=true` directly to simulate a
  terminal for tests that exercise the inline cursor-code path.
- **Lint cleanup.** Replaced deprecated `ansi.SetSynchronizedOutputMode` /
  `ansi.ResetSynchronizedOutputMode` / `ansi.CursorUp1` with current
  equivalents (`SetModeSynchronizedOutput` / `ResetModeSynchronizedOutput` /
  `CUU1`); reworded a comment to avoid a godox trigger; added `os/signal`
  and `syscall` to the depguard allow list for the SIGWINCH handler.

### Changed

- **Lockstep dependency realignment (v0.17.0 follow-up).** All 19 sub-module
  `go.mod` files now reference root `v0.17.0` (and `tui` references
  `nom/v0.17.0`), instead of a mix of stale v0.13.0 / v0.16.0 pins. This is
  the alignment that should have shipped inside v0.17.0; cutting v0.17.1
  makes the v0.17.x graph self-consistent for consumers who pull tagged
  versions rather than source.

### Docs

- Comprehensive README rewrite end-to-end: new headline + tagline; fix broken
  legacy imports throughout (`output.NewMarkdownTable()` →
  `markdown.NewMarkdownTable()`, etc.); promote NOM/TUI to a top-level
  section with event-driven subscription pattern and `DisplayMode` flags; add
  format gallery with real rendered output; reorganize format table by
  module; add `pkg.go.dev` badge, Go 1.26+ requirement, contributing/architecture
  pointers; fix the `cancelFunc` / `tableData` / `ActivityCompleted.Name`
  examples that would not compile or would silently corrupt the timing cache.

## [0.17.0] - 2026-06-20

Breaking type-model release. Three refactors harden the NOM event and activity
types into typed, sealed, decoupled domain forms, removing string-based
dispatch, leaky framework coupling, and lying name conventions. ADR 006 marks
NOM APIs experimental pre-v1.0, so these breaking changes are permitted. The
render hot-path O(n) elimination ships as non-breaking performance work in the
same release.

### Changed — Breaking (NOM type-model hardening)

- **`Event` is now a sealed sum type.** The `Event` interface carries an unexported `isEvent()` marker; concrete structs (`WorkflowStarted`, `ActivityCompleted`, …) replace the open interface + `GetEventType()` string dispatch. All five optional accessor interfaces (`WorkflowEventAccessor`, `DependenciesAccessor`, `HostAccessor`, `DownloadAccessor`, `KindAccessor`) are deleted — handlers read event fields directly via an exhaustive Go type switch. This makes silent event-routing typos a compile error instead of a runtime drop. Callers that implemented local accessor structs (tests, examples, integration) drop the boilerplate (~420 net lines) and construct concrete event types directly.
- **`ActivityKind` replaces the `"phase:"` ID-prefix convention.** Phase detection was `strings.HasPrefix(id, "phase:")` — a lying name that claimed to signal kind but was an opaque string matched at render time. `ActivityKind` (`ActivityKindTask` / `ActivityKindPhase`) is set at construction via `NewPhase()` and threaded through `ActivitySnapshot` (`snap.IsPhase()`). `isPhaseID` and the dead `ActivityNode.IsPhase()` (zero callers) are deleted. Callers that prefixed IDs with `"phase:"` to obtain phase rendering now call `NewPhase(id, name)` instead.
- **`Activity` no longer embeds `output.GraphNode`.** The domain type leaked diagram-export concerns (`Shape`, `Style`, `Metadata`) into the domain model and required `applyVisualStyle()` to mutate graph fields on every lifecycle transition. `Activity` now holds only identity + lifecycle + NOM display state (`ID`/`Label` are direct branded-type fields); `Shape` and `GraphStyle` are projected from `Status` at `subscriberView.Nodes()`, so the framework coupling lives at the export boundary. `applyVisualStyle()` becomes the slim `applyDisplayStyle()` (terminal `Symbol`/`Color` only).
- **`Activity.CurrentElapsed` field removed** (see Performance below). Callers reading elapsed off `*Activity` now read `ActivitySnapshot.CurrentElapsed`.

### Changed — Performance (render hot-path O(n) elimination)

- **`GetActivityCounts()` is now O(1)** instead of O(n) per render frame. The subscriber maintains an incrementally-updated `ActivityCounts` cache (`applyCountsDelta` on every state transition). Previously, the summary bar scanned all activities every tick (~10×/sec); now it reads a pre-computed value. At 10,000 activities this is ~1000× faster (flat ~10ns regardless of count, verified by benchmark). This ports the core algorithmic advantage of upstream nix-output-monitor's `DependencySummary` monoid. The invariant is verified by `TestActivityCountsCache_LifecycleConsistency` (brute-force recount vs cache after every event).
- **`UpdateRunningActivityElapsed()` deleted.** `CurrentElapsed` is now derived at snapshot time from `StartTime`/`EndTime`/`Status` inside `SnapshotActivities()` (`activity.elapsedAt(now)`), eliminating the per-tick O(n) write scan that ran every 100ms. Renderers read `ActivitySnapshot.CurrentElapsed` unchanged — the snapshot is the single source of truth for display timing. Callers that invoked `UpdateRunningActivityElapsed()` (inline renderer, TUI, examples, integration tests) simply drop the call; snapshots are always current by derivation.

### Changed — Other

- **Lockstep dependency bump.** All 17 sub-module `go.mod` files now reference root `v0.16.0` (tui also `nom/v0.16.0`), reflecting the mono-version tagging policy. Sub-modules build against root unchanged since the v0.17.0 breaking changes are nom-internal.

### Fixed

- **Restored `CODE_OF_CONDUCT.md`.** The file was auto-deleted by the BuildFlow pre-commit hook during the v0.16.0 lockstep dep bump (the same recurring issue documented in AGENTS.md "Gotchas" that struck the v0.14.0 release). Restored from the last known-good revision.

## [0.16.0] - 2026-06-20

Feature release: NOM-style rendering capabilities ported from the original
nix-output-monitor, plus a consolidation of the inline renderer's locking
model and CI/environment hardening. All new features are **backward-compatible**
— host/download/node-class are dormant opt-in annotations that render only
when populated.

### Added — NOM Feature Ports

- **`NodeClass` (root/twig/leaf)** — `ActivityNode.Class()` mirrors NOM's `mapRootsTwigsAndLeaves`. Root nodes now render **bold** so top-level activities stand out from intermediaries and leaves. (Golden files regenerated for the intentional root-styling change.)
- **Host column** — `ActivitySnapshot.Host` + optional `HostAccessor` event interface. Rendered as a dim `@host` tag when an event populates it; dormant otherwise. Mirrors NOM's host column.
- **Per-activity download progress bars** — `ActivitySnapshot.Download` (`DownloadProgress{Downloaded, Total int64}`) + optional `DownloadAccessor`. Rendered as a compact `▕████░░░░▏ 45%` bar while an activity is running; dormant otherwise. `formatBytes` renders human-readable sizes.
- **Height-pressure collapse marker** — when completed children are elided under `maxHeight` pressure, the renderer emits a faint `⋯ N completed` line so the collapsed work is visible rather than silently gone.
- **`SetMaxHeight`** — updates the tree-height cap at runtime, race-free (read via the lock-protected `snapshotConfig`).
- **`SetPlainText`** — runtime override for plain, append-only output (no cursor/sync escape codes). `plainText` is now part of `rendererConfig`, consistent with all other config.

### Added — Tests & Tooling

- **`scripts/pre-tag-check.sh`** (from v0.15.0, now exercised) — builds and tests every module with `-race` on the concurrency-sensitive ones before a tag is cut.
- **`nom_progress` example smoke test** — guards against the blank-render regression class.
- **Inline-renderer failure integration test** — verifies a failed activity + propagated workflow error surface through `Draw` and `Finish`.
- **Render-lock contention benchmarks** — `RenderUnderStepChurn`, `SnapshotActivities_Parallel`, `InlineRenderer_Draw`, `DrawWithChurn`.
- **`NodeClass`, host/download, collapse-marker, `DownloadProgress.Fraction` unit tests.**

### Changed — Consolidations

- **Single config snapshot** — `Draw`, `Finish`, and `renderSummary` each did scattered `tickMu.RLock` field reads (and `renderSummary` re-acquired the lock `Draw` already held). Replaced with one `snapshotConfig()` returning an immutable `rendererConfig`. `maxHeight` is now part of the snapshot. `renderMu` stays separate — merging would reintroduce the documented `Stop()` deadlock.
- **Inline renderer CI degradation** — `detectPlainText()` (evaluated at construction) returns true under `envdetect.IsCI()`; `Draw` then appends frames line-by-line instead of cursor-up/clear-line/sync codes that corrupt captured CI logs.

### Fixed

- **Cyclop lint violations** in `Draw` and `walkSubtree` (both exceeded 13 after feature additions) — extracted `buildRedrawOutput` and `appendCollapseMarker` as focused helpers.

## [0.15.0] - 2026-06-20

Design-smell cleanup following the v0.14.0 type-model refactor. Two breaking
API changes remove footguns that the snapshot refactor exposed: a `*Activity`
parameter that was silently discarded, and convenience renderers that produced
blank output when invoked without snapshots.

### Changed — Breaking

- **Deleted nil-snapshot footgun wrappers**: `RenderString`, `RenderWithWidth`, `VisibleNodes` — they silently passed nil snapshots, producing blank labels. Callers now use `RenderWithSnapshots(snapshots, maxHeight, maxWidth)` / `VisibleNodesWithSnapshots(snapshots, maxHeight)`. The "render with no data" mistake is now unrepresentable.
- **`AddActivity` signature simplified**: removed the dead `*Activity` parameter (it was silently discarded after the snapshot refactor — a "lying" signature). New: `AddActivity(id ActivityID, deps []ActivityID)`.

### Added

- **`scripts/pre-tag-check.sh`** — pre-release verification that builds and tests every module (with `-race` for the concurrency-sensitive `nom`/`tui`/`integration` modules) before a tag is cut. Accepts an optional `vX.Y.Z` argument to also verify the tree is clean and the tag does not already exist.

### Fixed

- **Compile-broken integration test, TUI test, and example** after v0.14.0 refactor — rewritten to use `SnapshotActivities` + `RenderWithSnapshots`.
- Removed an accidentally committed binary (`examples/nom_progress/nom_progress`) from the repository.

## [0.14.1] - 2026-06-20 _(tui module only)_

### Fixed

- **`tui.Subscriber()` re-exported** — the v0.14.0 hardening unexported it, but BuildFlow (the primary external consumer) needs it to dispatch lifecycle events. Re-exported as public API. No other module changed; only `tui/v0.14.1` was tagged.

## [0.14.0] - 2026-06-20

Pre-v1 cleanup, rendering bug fixes, data-race elimination, and dead-code removal.

### Added — API Stability

- **`nom.ActivitySnapshot`** — immutable value copy of mutable activity fields, taken under the subscriber's read lock via `SnapshotActivities()`. The dependency tree renderer consumes snapshots instead of reading a shared pointer, eliminating the data race at the type level.
- **`nom.DependencyTree.RenderWithSnapshots(snapshots, maxHeight, maxWidth)`** / **`VisibleNodesWithSnapshots(snapshots, maxHeight)`** — the explicit snapshot-consuming render API.
- **`nom.ActivityCounts.Summary()`** / **`.CompletionPercent()`** — single source of truth for count formatting, shared by the inline renderer and the TUI (eliminates the prior split-brain).

### Changed — Breaking

- **Eliminated the shared `*Activity` pointer from `ActivityNode`.** This was the root cause of every NOM renderer data race — event handlers mutated the same `*Activity` the renderer read, with no synchronization. `ActivityNode` now holds only an `ActivityID` and tree structure. All mutable fields are consumed via `ActivitySnapshot`. The lock-based bandaid (`RenderTree`, `WithSubscriberRLock`) is deleted.
- **Removed the `ActivityStatusPaused` status.** It was fully decorated (symbol, color, shape, sort priority) but had no `EventActivityPaused` constant, no subscriber handler, and `SetPaused()` had zero callers — unreachable through the event system. `ActivityStatusPending` previously _reused_ Paused's symbol/color, making the two indistinguishable; Pending now has its own honest identity.
- **`SymbolPaused` → `SymbolPending`** (glyph `⏸` → `○`), **`SemanticColors.Paused` → `SemanticColors.Pending`**, and the deprecated `ColorPaused` alias is removed.
- **`ActivityStatus.Interest()`** renumbered after Paused removal: `failed=0, running=1, pending=2, completed=3`.
- **Removed deprecated functions:** `nom.EnsureBuild`, `nom.ParseActivityID`/`ParseWorkflowID`, `graph.NewGraphNodeID`/`NewGraphNodeLabel`, `delimited.MarshalTSV(any)` (+ `writeTSVData`, `ErrUnsupportedType`).
- **Removed deprecated color aliases:** `ColorRunning`/`ColorCompleted`/`ColorFailed`/`ColorInfo`/`ColorPhase` — use the `Colors.X` fields.
- **Removed speculative `OperationSymbol` + `OperationTypeDownload`/`OperationTypeUpload`** — a stringly-typed mapping with no production caller. `SymbolDownload`/`SymbolUpload` remain as palette members.
- **`tui/` public surface minimized.** Unexported internals not part of the public contract: Bubble Tea message types, `UpdateType`, `WorkflowState`, `ProgressStep`, `TickCmd`, `MsgNoActivitiesToDisplay`. (`Subscriber()` was later re-exported — see [Unreleased].) Removed dead `CancelMsg` and `SeparatorLineEquals`.

### Fixed

- **Data race in `InlineRenderer.renderSummary`** on `startTime` (TOCTOU) — snapshotted once under `tickMu.RLock()`.
- **Deep-nested trees overflowed `maxWidth`** — the final styled line is now truncated to the terminal width.
- **`FormatDuration` showed `90m` for ≥1h durations** — added hours branch (`1h`, `1h30m`, `24h`).
- **Four inline-renderer races/deadlocks** closed (config setters, background refresh, Finish).
- **`TimingCache.EnsureLoaded` blocked all concurrent `GetMedian` reads during file I/O** — now reads the file lock-free, then re-acquires the write lock to publish (mirrors `saveAsync`).

### Removed — Dead Code

- `nom.FormatTimingInfo` (buggier duplicate of `FormatActivityNodeTiming`), `nom.GetActivitySummaryString` (params named `uploading`/`downloading` but fed `completed`/`failed`), `nom.Activity.Elapsed()` (0 callers; would race), `nom.SummaryWithTotal()` (extracted but never wired).
- `nom.DependencyTree` dead fields `buildOnce sync.Once` and `order []ActivityID`.

### Tests

- End-to-end integration test: full workflow → events → inline render → Finish.
- `FormatDuration` fuzz test (full int64 ns range, 3M+ execs clean).
- `formatActivityLabel` benchmark across all activity statuses.
- Coverage for `Copy()` nil-metadata, `removeChild`, `StripANSI`, `MultiSubscriber.Subscribers()`, and CJK/emoji `VisibleWidth`/`VisibleLineCount`.

## [0.13.0] - 2026-06-18

### Added — New Modules

- **`markdown/` module** — Markdown table renderer extracted from root into an independent sub-module. Self-registers via `init()` for `output.RenderTableData()` dispatch. Root no longer carries Markdown rendering code; importing `markdown/` activates `FormatMarkdown`.
- **`tree/` module** — ASCII tree renderer (`ASCIITreeRenderer`) extracted from root into an independent sub-module. Self-registers via `init()` for `FormatTree` dispatch. Box-drawing characters, depth-based color cycling, metadata summaries.

### Added — nom/ Composition Refactor (ADR 007)

- **`nom.Activity`** — unified activity type embedding `output.GraphNode` (ID, Label, Shape, Style, Metadata) + temporal domain fields (Status, StartTime, EndTime, EstimatedTime, Err, Dependencies). Single source of truth for identity, visual representation, and state.
- **`nom.ActivityStatus.NodeShape()` / `.GraphStyle()`** — mappers from domain status to root's visual types. Same status now drives both terminal lipgloss styling AND diagram export (DOT/Mermaid/D2/PlantUML).
- **`nom.ActivityStore`** — map-backed store with `Nodes() []output.GraphNode` / `Edges() []output.GraphEdge` projections. Any `output.GraphRenderer` can consume live progress state for diagram export.
- **`nom.MultiSubscriber`** — `io.MultiWriter`-style fanout for `EventSubscriber`. Dispatches events to N subscribers; nil subscribers skipped; errors from one don't block others.
- **`nom.NOMStyleSubscriber.Store()`** — exposes the `ActivityStore` for diagram export.
- **`nom.Symbol` typed constants** — `SymbolRunning`, `SymbolCompleted`, `SymbolFailed`, `SymbolPaused`, `SymbolDownload`, `SymbolUpload`, `SymbolTiming`, `SymbolAverage`, `SymbolTotal`, `SymbolPhase`. Replaces bare `string` symbols with a typed `Symbol` type, preventing accidental mixing with arbitrary strings.
- **`nom.AllActivityStatuses`** — complete iterable list of valid `ActivityStatus` values, backing `ParseActivityStatus()`, `IsValid()`, and `AllowedValues()`.

### Added — New Types

- `LineStyle` enum (`LineStyleSolid`, `LineStyleDashed`, `LineStyleDotted`) with `IsValid()`/`String()` — replaces free-form string on `EdgeStyle.Style`
- `Direction` enum (`DirectionDown`/`Up`/`Left`/`Right`) with `ToD2Direction()`/`ToRankDir()` — bridges D2 and DOT layout vocabulary
- `ActivityCounts` struct with `Running`/`Completed`/`Failed`/`Pending` fields and `Total()` method — replaces 4 unnamed int returns from `GetActivityCounts()`

### Added — Documentation

- **ADR 008** — Dedup workflow decision: 5-step checklist for evaluating `art-dupl` clone groups at threshold `t=24`.
- **`RELEASE.md`** — Release process documentation for the 20-module mono-version workspace.
- **`ROADMAP.md`** — Long-term direction and raw ideas.
- Cross-reference comments on `nom.StatusStringUnknown` / `tui.WorkflowStateStringUnknown` documenting both use `"unknown"` for the same purpose.

### Changed — Breaking (0.x — pre-1.0 API cleanup)

- **`nom.ActivityNode`** now embeds `nom.Activity` (which embeds `output.GraphNode`) instead of `DisplayState`. Removed fields: `ActivityNode.ActivityID`, `ActivityNode.ActivityName`. Use `node.ID.Get()` and `node.Label.Get()` from the embedded `GraphNode` instead.
- `GraphShape` → `NodeShape` (type), `ShapeBox` → `NodeShapeBox`, `ShapeEllipse` → `NodeShapeEllipse`, etc. — disambiguated from data-capability `Shape` enum (`ShapeTable`, `ShapeTree`, `ShapeGraph`)
- `GraphStyle.FillColor` → `Fill`, `GraphStyle.StrokeColor` → `Stroke` — aligned with `D2NodeStyle` field names
- `EdgeStyle.Style` field → `EdgeStyle.Line` (type changed from `string` to `LineStyle`)
- `TableDataMarshaler` → `TableDataRenderer`, `AnyDataMarshaler` → `AnyDataRenderer` — unified terminology
- `RegisterTableDataMarshaler` → `RegisterTableDataRenderer`, `RegisterAnyDataMarshaler` → `RegisterAnyDataRenderer`
- `nom.InlineRenderer.Render()` → `nom.InlineRenderer.Draw()` — resolves split-brain M4: `Render()` is now reserved for the `output.Renderer` contract `(string, error)` project-wide. `Draw()` writes a frame to the configured `io.Writer` (void return).
- `nom.DependencyTree.Render(maxHeight)` → `nom.DependencyTree.RenderString(maxHeight)` — same finding: the bare-string return (errors encoded into the string) is now distinguished from the canonical `(string, error)` contract.
- `nom.TreeNode` renamed to `nom.ActivityNode` — eliminates exported name collision with `output.TreeNode` (the generic diagram tree node). Mechanical rename across 12 nom/ files + tui/ references.
- `tui.TimingFormat` renamed to `tui.timingFormatWithIcon` (unexported) — eliminates name collision with `nom.TimingFormat`. The TUI version bakes in the `⏱️` emoji; the NOM version keeps symbol separate from format string.
- `nom.GetWorkflowID()` returns `WorkflowID` instead of `string` — matches the `WorkflowEventAccessor` interface return type.

### Changed — Non-Breaking

- `ColorModeAuto.ShouldColor()` detection functions are now overridable variables (`stdoutIsTerminal`, `noColorEnv`, `ciEnv`) for deterministic testing
- `BubbleTeaProgressReporter` no longer mutates model fields directly from caller goroutine — all mutations flow through `send()` → `model.Update()`, eliminating the data race (#22)
- **Color detection centralized** — root `isNoColor()` now checks `TERM=dumb`; nom `detectNoColor()` now checks all 5 CI env vars (`CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`). Both delegate to the shared **`envdetect`** module (`envdetect.IsCI()` / `envdetect.IsNoColor()`).
- **`nom.Colors` struct** centralizes the semantic color palette. Deprecated aliases (`ColorRunning`, `ColorCompleted`, etc.) point to `Colors.X` fields for backward compatibility.
- **`delimited.tableDataWriter` interface now includes `WriteFooter`** — footer rows written via dedicated method instead of `WriteRow`, matching the streaming API.
- **`ActivityStatus` documented as intentionally separate from `tui.WorkflowState`** (split-brain m3) — they share `Running`/`Completed` values but model different lifecycles (single activity vs whole workflow).
- **Canonical branded-ID paths documented** (split-brain m6) — `D2NodeID`/`D2NodeLabel` brand types live in root so both root and `d2/` callers see the type alias without a circular import.

### Removed

- `RenderOptions.GraphID` — dead code, no marshaler ever read it. Use `DOTRenderer.SetGraphID()` directly.
- **`nom.ColorWarning`** — was identical to `ColorRunning` (`lipgloss.Color("11")`) with zero callers.
- **`tui.ProgressModel.timingCache` field** — stored a `*nom.TimingCache` that was never read.
- **`nom.ActivityStore` ghost system** (155 LOC) — dead code that was never wired into the subscriber lifecycle. Replaced by the real `ActivityStore` projection.
- **`nom.Activity.Dependencies` and `OperationType` fields** — dead fields with zero readers.

### Fixed

- **Test interface redeclarations eliminated** — `serialization/testhelpers_test.go` now uses `output.GraphRenderer` directly instead of a re-declared `graphRenderer`. `integration/renderer_test.go` now uses `output.Renderer` directly instead of a re-declared `renderer`.
- **All bare event-string literals replaced with `nom.Event*` constants** — 34 literals across `nom/subscriber_test.go`, `tui/`, `integration/`, `examples/` now use `nom.EventWorkflowStarted`, `nom.EventActivityCompleted`, etc.
- **Hardcoded `"No activities to display"` literal** in `tui/view.go` replaced with `MsgNoActivitiesToDisplay` constant.
- **`nom.detectNoColor()` now checks terminal** — aligned with root's `isStdoutTerminal()` gate. Previously nom would emit ANSI color codes even when piped to a file.
- **20/20 split-brain findings resolved** — see `docs/status/` reports for comprehensive elimination details.
- **Hermetic test isolation** — nom/output tests no longer leak global state between test runs.

### Internal

- **CI duplication gate** — `art-dupl -t 24` integrated into CI loops across markdown/ and tree/ modules.
- **Module count: 18 → 20** — markdown/ and tree/ extracted as independent sub-modules.

## [0.12.0] - 2026-06-17

### Added

- **`MermaidRenderer.SetCodeFence(bool)`** — raw flowchart syntax for `.mmd` files, Mermaid CLI, and embedded diagrams.
- **Per-node `GraphStyle` in Mermaid and PlantUML** — `fill`, `stroke`, `font-color`, and `font-size` are now honored. Replaces the hardcoded pink Mermaid classDef and rigid PlantUML skinparam defaults.
- **`GraphRendererState.DedupEdges()`** — removes duplicate edges by `(from, to)` in-place.
- **DOT typed layout enums** — `RankDir` and `SplineStyle` make invalid layout values unrepresentable, with `Parse`/`String`/`IsValid`/`AllowedValues`.
- **`DOTRenderer.SetRankDir(RankDir)`**, **`SetSplines(SplineStyle)`**, **`SetNodeSep`**, **`SetRankSep`** — configurable DOT layout attributes (defaults preserved).

### Changed

- **Root `go.mod` no longer requires `serialization`** — importing `graph`, `plantuml`, or `d2` no longer transitively pulls `go-faster/yaml`, `go-toml/v2`, `go-faster/jx`, or `segmentio/asm` into consumers.
- **Diagram renderer API stability documented** — ADR 006 now defines stable/experimental tiers for `GraphRenderer`, Mermaid, DOT, PlantUML, D2, and NOM/TUI APIs.

### Fixed

- **`escape.SlugifyID`** now handles dots, asterisks, brackets, braces, and parentheses, preventing silent node ID collisions like `foo.bar` ↔ `foo_bar`.

## [0.11.0] - 2026-06-17

### Added

- **`nix run .#test-race`** — new nix app running `go test -race -count=1` for the `nom/` and `tui/` concurrency-sensitive modules.
- **`ActivityStatus.Interest()`** — priority ordering is now a method on the enum type, replacing the standalone `activityInterest()` function.
- **Completed-subtree collapsing** — under height pressure, completed children are elided to prioritize active work (running, failed, pending).
- **Completion percentage in summary bar** — shows `(N%)` where N = (completed + failed) / total.
- **iTerm2 synchronized updates** — inline renderer frames wrapped in `\x1b[?2026h/l` for flicker-free redraws.
- **E2E smoke test** for inline renderer lifecycle (Start → events → Refresh → Stop → Finish).
- **Benchmarks** for `RenderWithWidth` and `childPriority` hot paths.
- **Dedicated tests** for `elideCompletedUnderPressure` and `ActivityStatus.Interest()`.

### Changed

- **TOML renderer wraps rows in `[[row]]`** — TOML cannot encode a bare `[]map[string]string` as document root; rows are now nested under a named key producing valid array-of-tables syntax.
- **ANSI parsing delegated to `charmbracelet/x/ansi`** — `StripANSI`, `VisibleWidth`, `TruncateVisible` now delegate to `ansi.Strip`, `ansi.StringWidth`, `ansi.Truncate`, removing 128 lines of hand-rolled scanner code. Direct `runewidth` dependency eliminated.
- **`Finish()` now calls `Stop()` first** — prevents concurrent tree access between the background render goroutine and the final render.
- **`RenderNode` signature documented** — the second parameter is now named `visibleNodes` (was `_ []*TreeNode`), honestly documenting it as reserved for future width-aware truncation.
- **`log.Printf` → `slog.Error`** in TUI lifecycle per how-to-golang policy.
- **Test refresh rendezvous** — replaced `time.Sleep(50ms)` with deterministic `renderNotify` channel in `TestInlineRenderer_Refresh_TriggersRender`.

### Fixed

- **TOML round-trip test failure** — `nix run .#test` was red on integration/serialization due to `toml: cannot encode a []map[string]string as a document root`.
- **TUI data races** — `go test -race` in tui/ failed on 18 tests; reporter tests now use `newTestReporter()` helper that prevents the real Bubble Tea program from starting.
- **Registry test pollution** — `TestRegisterTableDataMarshaler_ConcurrentAccess` leaked "race-test-\*" formats into the global registry; `TestRegisteredTableDataFormats` now checks known-good formats instead of asserting all are valid.
- **Dead code in `elideCompletedUnderPressure`** — removed unreachable guard `visibleCount+maxHeight <= 0`.

## [0.9.0] - 2026-06-12

### Added

- **`RenderAnyData`** — new registry-based dispatch for rendering arbitrary `any` data (not just `TableData`). JSON, YAML, and TOML register `AnyDataMarshaler` handlers via `init()`.
- **`RegisterAnyDataMarshaler`** / **`RegisteredAnyDataFormats`** — register and introspect any-data marshalers.
- **`RegisteredTableDataFormats`** — returns all formats with registered `TableDataMarshaler`s.
- **`TableData.AddRowChecked(row []string) error`** — fail-fast row addition that returns `ErrColumnMismatch` when column count differs from headers.
- **`nom.Event*` constants** — typed event string constants (`EventWorkflowStarted`, `EventActivityCompleted`, etc.) replacing bare string literals for safer event dispatch.
- **All 16 formats register as `TableDataMarshaler`** — D2, DOT, Mermaid, and PlantUML now register via `init()`, making `RenderTableData()` dispatch all 16 formats when sub-modules are imported (previously only 10 were registered).
- `testhelpers.AssertLineCount`, `AssertLastLineContains`, `AssertErrorContains` — new shared test assertions.
- `tui/colors.go` — extracted all color and layout constants from scattered inline values.
- `escape/fuzz_test.go` — fuzz tests for all escape functions.

### Changed

- **Generic `formatRegistry[T]`** replaces 3 separate mutex+map registry patterns (shape capabilities, table-data marshalers, any-data marshalers). Eliminates duplicate mutex boilerplate.
- `RenderTableData()` doc updated: all 16 formats are now dispatchable with sub-modules imported.
- NOM subscriber handlers simplified — reduced complexity in event routing.
- NOM `DisplayState` embedded directly, removing indirection.
- TUI reporter and view rendering refactored for clarity.
- D2, graph, plantuml `dt.Build()` errors now propagated instead of silently discarded.
- NOM timing cache async save errors now logged instead of swallowed.

### Deprecated

- `delimited.MarshalTSV(data any)` — use `MarshalTSVFromTableData` or `TSVWriter` directly.
- `graph.NewGraphNodeID` / `graph.NewGraphNodeLabel` — use `output.NewBrandedID` directly.
- `nom.ParseActivityID` / `nom.ParseWorkflowID` — use direct type conversion with manual validation.

### Internal

- Coverage improvements across 5 modules (nom tree/render/subscriber, table registry, serialization, testhelpers).
- Test deduplication across 22 files — shared helpers extracted to `testhelpers/`, table-driven patterns.
- Registry benchmarks for generic `formatRegistry[T]`.
- Comprehensive hardening sprint status reports in `docs/status/`.
- Graph DOMAIN_LANGUAGE updated with registry dispatch documentation.

## [0.8.0] - 2026-06-10

### Added

- **TUI keyboard navigation** — arrow keys, mouse scroll, click selection on dependency tree nodes.
- **`tui.SetCancelFunc`** — allows Ctrl+C to cancel workflow context via cancel function.
- **NOM secondary dependencies** — `DependenciesAccessor` and secondary parent labels in tree rendering.
- **NOM `SetDisplayMode`** — switch between NOM and Universal display modes.
- NOM/TUI golden file tests and event sequence tests.

### Changed

- All `panic()` calls eliminated from production code — replaced with error returns.
- NOM `AddActivity` made idempotent — prevents duplicate children.
- NOM dependency tree locking fixed — `GetDependencyTree` now properly synchronized.

### Fixed

- Timing cache race condition in async save.
- TUI help overlay and viewport scroll issues.
- TUI ticks now allowed in Idle state for NOM sync.
- Pre-registered activities preserved correctly.

## [0.7.0] - 2026-06-09

### Added

- **`nom/` module** — NOM-style real-time progress visualization with dependency trees, timing cache (CSV-persisted), and event-driven activity tracking. Zero imports from root module. `NOMStyleSubscriber` implements `EventSubscriber` with string-based event routing and type-assertion accessor interfaces.
- **`tui/` module** — Bubble Tea interactive TUI for workflow progress display. `BubbleTeaProgressReporter` with state machine (`WorkflowState`: Idle→Running→Completed/Errored), step-based progress, and NOM integration. Depends on `nom/`.
- `TableData.Validate()` now rejects nil rows — catches `nil` in `Rows [][]string` before rendering.
- Race test for `RegisterFormatShapes` confirms thread-safety of the shape capability matrix.

### Changed

- **BREAKING**: `RenderTableData()` changed from variadic `opts ...RenderOptions` to single `opts RenderOptions` — only `opts[0]` was ever used, the variadic signature was misleading.
- `GraphRendererMixin` API refactored for cleaner method signatures.
- String escaping optimized across modules — `strings.NewReplacer` for single-pass replacements.
- File renames for naming consistency across the project.
- Transitive dependencies updated across all 15 modules.

### Fixed

- Gosec and wrapcheck lint violations resolved.
- Critical bugs and config gaps from comprehensive quality audit.
- Flaky `TestTimingCache_SaveAndLoad` — async `saveAsync()` goroutine now completes before test cleanup.

### Internal

- `go.work.example` and `AGENTS.md` updated for 15 modules (was 13).
- Coverage table updated: nom 93.1%, tui 84.2%.
- Nix configuration standardized with format and build checks.

## [0.6.3] - 2026-06-03

### Changed

- Updated `charmbracelet/ultraviolet` from `v0.6.0-20260525` to `v0.6.0-20260601` — includes fix for modified Kitty keyboard navigation/function key releases (affects `table/`, `examples/`, `integration/` modules via lipgloss/v2 transitive dependency).
- Updated `golang.org/x/exp` from `v0.0.0-20260527` to `v0.0.0-20260529`.
- Updated nixpkgs flake input to latest revision.

### Internal

- Eliminated all code clones at `art-dupl` threshold 30 across test files in `markup/`, `serialization/`, `table/`, `graph/`, `plantuml/`, `testhelpers/`, and root module.

## [0.6.2] - 2026-06-01

### Fixed

- All internal cross-module `go.mod` references upgraded from `v0.0.0` pseudo-versions to canonical `v0.6.2` tags across all 14 modules. Resolves the chicken-and-egg release issue where v0.6.1 was tagged before dependency versions were bumped.
- `integration/go.mod`: fixed `d2`, `graph`, `markup`, `plantuml`, `table` references from `v0.0.0` to `v0.6.2`.
- `plantuml/go.mod`: fixed `testhelpers/graphtest` reference from `v0.0.0-00010101000000-000000000000` to `v0.6.2`.
- `d2/go.mod`, `graph/go.mod`, `serialization/go.mod`: fixed `testhelpers/graphtest` references from `v0.0.0` to `v0.6.2`.
- `enum/go.mod`, `testhelpers/graphtest/go.mod`: fixed root and testhelpers references from `v0.0.0` to `v0.6.2`.

### Changed

- Mono-version tagging policy documented in `AGENTS.md` — all 14 modules release in lockstep.

### Added

- `docs/research/go-error-family-adoption-report.html` — comprehensive PRO/CONTRA analysis evaluating `github.com/larsartmann/go-error-family v0.3.0` adoption. Verdict: Do Not Adopt (Yet); Strategy B recommended (add to `examples/` module only).

## [0.6.1] - 2026-05-30

### Added

- **Footer row** (`TableData.Footer []string`) — optional totals/summary row on `TableData`. Tabular renderers render it visually: CSV/TSV append as last row, HTML uses `<tfoot>` with `footer-cell` CSS class, XML uses `<footer>`, Markdown adds separator + bold row (with column alignment), AsciiDoc appends footer row, Terminal Table uses bold styling. Data formats (JSON/YAML/TOML/JSONL) skip footer.
- `TableData.GetFooter()`, `TableData.HasFooter()`, `TableDataStore.SetFooter()`, `TableDataStore.HasFooter()` — accessor methods for footer row.
- `TableData.Validate()` — validates footer column count matches headers. Returns `errColumnMismatch` on mismatch. Wired into `RenderTableData()` for automatic validation.
- `MarkdownTable.SetFooter()` — sets footer row on Markdown table renderer (inherits column alignment).
- `table.SetFooter()` — adds bold-styled footer row to lipgloss terminal table. Tracks `footerRowIndex` for correct bold styling on multiple calls.
- `table.FooterProvider` — optional interface checked by `FromTableData()` for automatic footer styling.
- `CSVWriter.WriteFooter()`, `TSVWriter.WriteFooter()` — explicit footer methods for streaming delimited output.
- Package doc.go files for 8 packages (output, delimited, d2, graph, markup, plantuml, serialization, testhelpers) — pkg.go.dev now shows proper package documentation.
- GoDoc examples for `Format.IsValid`, `ParseFormat`, `ColorMode`, `Shape`, `TableData.Validate`, `MustRender`.
- GoDoc comments on all exported struct fields in graph, tree, and d2 types (40+ fields).

### Changed

- `delimited.marshalFromTableData()` — extracted shared helper from `MarshalCSVFromTableData` and `MarshalTSVFromTableData`, eliminating ~70 lines of duplication.
- HTML footer cells now use `class="footer-cell"` for CSS targeting (both `HTMLRenderer` and `StreamingHTMLRenderer`).
- All 14 modules unified to `go 1.26.3`.
- `table.SetFooter()` now correctly tracks `footerRowIndex` — only the last footer row receives bold styling.
- Root module test coverage improved from 88.6% to 96.1%.

### Fixed

- `MarkdownTable.AsTableRenderer()` — adapter wrapping fluent MarkdownTable API as void-returning `TableRenderer` interface.
- `table.Table.AsTableRenderer()` — adapter wrapping lipgloss table builder as `TableRenderer` interface.
- `table.WithFooterStyle(func(lipgloss.Style) lipgloss.Style)` — composable footer styling for lipgloss tables.
- Alignment constants (`AlignmentLeft/Right/Center`) unexported to `alignmentLeft/Right/Center` — were documented as unexported but actually exported.
- `UnsupportedFormatError.Unwrap()` removed — returned nil, semantically identical to not having it.
- ADR 004: Footer row design decision record.
- Coverage: gentest 80.8%→96.2%, integration 82.8%→95.5%, table 85.5%→100%, serialization 89.0%→91.4%.
- `table/table_test.go` split into `table/table_test.go` + `table/color_test.go` (391→274 lines, under 350-line limit).
- GoDoc on all exported testhelpers symbols: `ErrTest`, `ErrorRenderer`, `FixedRenderer`, `ErrWrite`, `ErrorWriter`, `WriteNThenFailWriter`.
- Package doc added to `testhelpers/graphtest`.
- `integration/go.mod` root dep reference fixed from `v0.5.0` to `v0.0.0`.
- AsciiDoc renderer now uses `HasFooter()` consistently with other renderers.
- `TestBrandedIDFormat` updated for go-branded-id v0.3.0 `%#v` output.

## [0.6.0] - 2026-05-25

### Added

- **JSONL** (`jsonl` format) — JSON Lines output, one JSON object per line. `serialization.NewJSONLTableRenderer()`, `serialization.MarshalJSONLFromTableData()`, `serialization.NewJSONLWriter()`. Supports `ShapeTable`.
- **AsciiDoc** (`asciidoc` format) — AsciiDoc table output. `markup.NewAsciiDocTableRenderer()`, `markup.MarshalAsciiDocFromTableData()`. Supports `ShapeTable`.
- **TOML** (`toml` format) — TOML serialization with table and tree support. `serialization.MarshalTOML()`, `serialization.UnmarshalTOML()`, `serialization.NewTOMLTableRenderer()`, `serialization.NewTOMLTreeRenderer()`. Supports `ShapeTable`, `ShapeTree`. Uses `github.com/pelletier/go-toml/v2`.
- **PlantUML** (`plantuml` format) — PlantUML component diagrams. `plantuml.NewPlantUMLDiagram()`, `plantuml.PlantUMLFromTableData()`, `plantuml.PlantUMLFromTree()`. Supports `ShapeTable`, `ShapeGraph`. New independent `plantuml/` module with zero external dependencies.
- 4 new `Format` constants: `FormatJSONL`, `FormatAsciiDoc`, `FormatTOML`, `FormatPlantUML` — 16 formats total.
- Shape capability matrix expanded: JSONL and AsciiDoc support `ShapeTable`; TOML supports `ShapeTable` + `ShapeTree`; PlantUML supports `ShapeTable` + `ShapeGraph`.
- `JSONTableRenderer` — renders `TableData` as a JSON array of objects (implements `Renderer` + `TableRenderer`)
- `YAMLTableRenderer` — renders `TableData` as a YAML sequence of mappings (implements `Renderer` + `TableRenderer`)
- `TableData.ToMapSlice()` — converts tabular data to `[]map[string]string` for serialization
- `UnsupportedFormatError` — renamed from `ErrUnsupportedFormat` (follows Go naming conventions)
- `TableDataMarshaler` registry — sub-modules register via `init()`, root has zero sub-module imports
- `TableDataStore` — exported from root for cross-module embedding
- `RenderTableData` now accepts `RenderOptions.ColorMode` to control ANSI color output for terminal renderers.
- `table.New()` accepts `WithColorMode(ColorMode)` functional option — lipgloss styles conditionally applied based on terminal detection.
- `ASCIITreeRenderer.SetColorMode(ColorMode)` — depth-based ANSI color cycling, bold labels, dim connectors, cyan metadata.
- `MarkdownTable.SetColorMode(ColorMode)` — bold headers, dim separators when terminal detected.
- ColorMode auto-detection respects `NO_COLOR`, `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`, `GO_OUTPUT_FORCE_COLOR`, `FORCE_COLOR` env vars.
- `MarshalFormat()`, `UnmarshalFormat()`, `MarshalJSONIndent()` — exported helpers used by serialization/ and markup/ (BrandedValue removed)
- `flake.nix` — Nix flake with devShell (Go 1.26.2, golangci-lint, gopls), treefmt-nix formatter, git-hooks.nix
- `.envrc` — direnv integration for automatic `nix develop` on cd
- Depguard whitelist for `examples/` module (scoped `examples/**/*.go` rule)
- CI: golangci-lint v2 (`github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

### Changed

- **BREAKING**: CSV and TSV writers moved to `delimited/` sub-module — `output.NewCSVWriter` → `delimited.NewCSVWriter`, `output.NewTSVWriter` → `delimited.NewTSVWriter`. Import `github.com/larsartmann/go-output/delimited`.
- **BREAKING**: JSON and YAML marshalers/renderers moved to `serialization/` sub-module — `output.MarshalJSON` → `serialization.MarshalJSON`, `output.NewJSONTableRenderer` → `serialization.NewJSONTableRenderer`, etc. Import `github.com/larsartmann/go-output/serialization`.
- **BREAKING**: XML, HTML, and StreamingHTML renderers moved to `markup/` sub-module — `output.NewHTMLRenderer` → `markup.NewHTMLRenderer`, `output.MarshalXMLFromTableData` → `markup.MarshalXMLFromTableData`, etc. Import `github.com/larsartmann/go-output/markup`.
- **BREAKING**: D2 diagram types moved to `d2/` sub-module — `output.D2Node` → `d2.D2Node`, `output.NewD2Renderer` → `d2.NewD2Diagram`, etc. Import `github.com/larsartmann/go-output/d2`.
- **BREAKING**: DOT and Mermaid renderers moved to `graph/` sub-module — `output.DOTFromTableData` → `graph.DOTFromTableData`, `output.MermaidFromTableData` → `graph.MermaidFromTableData`. Import `github.com/larsartmann/go-output/graph`.
- `RenderTableData` now returns `UnsupportedFormatError` for D2, Mermaid, and DOT formats (use sub-module constructors directly).
- `RenderTableData` uses registry-based dispatch via `TableDataMarshaler` — sub-modules register via `init()`. Root has zero sub-module imports.
- `RenderTableData` — all writer errors now wrapped with `fmt.Errorf("write X: %w", err)` for pinpoint failure reporting
- `tableDataBase` exported as `TableDataStore` with `Data()` getter — enables cross-module embedding.
- `marshal()`, `unmarshal()`, `brandedValue()` exported as `MarshalFormat()`, `UnmarshalFormat()`, `BrandedValue()` — used by serialization/ and markup/.
- Multi-module workspace: 13 independent modules (see ADR 001, ADR 003).
- Root production code has zero imports from sub-modules (`delimited`, `serialization`, `markup`, `d2`, `graph`, `table`, `plantuml`).
- Root production code has zero `go-faster/yaml`, zero `go-toml/v2`, and zero `escape` imports (isolated in `serialization/` and `markup/`).
- `FilledStrings` — uses `slices.Repeat` (Go 1.26 stdlib) instead of manual make+for loop
- `NewBrandedID` — simplified from `id.NewID[Brand, string](value)` to `id.NewID[Brand](value)` (inferred type arg)
- Added `Nodes()`, `Edges()`, `NodesPtr()`, `EdgesPtr()` accessor methods to `GraphRendererState` for cross-package use
- `enum/enum_test.go` no longer imports `internal/gentest` (inlined helper)
- `.gitignore` — added `result` and `.direnv/` for Nix artifacts
- Deduplication sprints reduced code clones from 44 → 26 (41% total reduction)

### Removed

- **BREAKING**: `format_deprecated.go` removed — `OutputFormat`, `FormatCategory`, `ParseOutputFormat`, `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` all removed. Use `Format`, `Shape`, `ParseFormat()`, `Supports(Shape*)`, `Shapes()` instead.
- **BREAKING**: `sort/` module removed — `ByField` is trivially replaceable with `slices.SortStableFunc` + `cmp.Compare` from stdlib.
- **BREAKING**: `registry.go` removed — `Register()`, `Create()`, `Unregister()`, `RegisteredFormats()`, `IsRegistered()` all removed. Use direct constructors (`d2.NewD2Diagram()`, `serialization.NewJSONTableRenderer()`, etc.) instead.
- **BREAKING**: `SortBy` enum removed from root — `sort.go`, `ParseSortBy()`, `SortBy*` constants all removed. Zero external callers.
- **BREAKING**: `FilledStrings()` removed — use `slices.Repeat` (Go 1.26 stdlib) instead.
- **BREAKING**: `BrandedValue()` removed from `marshal.go` — zero callers after sub-module extraction.
- **BREAKING**: `MermaidFlowchartRenderer` and `MermaidTreeRenderer` removed — use `MermaidFromTableData` and `MermaidFromTree` instead.
- `gci` formatter removed from `.golangci.yml` (conflicted with `goimports` on local-prefix grouping in sub-modules).
- `ErrUnsupportedFormat` — renamed to `UnsupportedFormatError` (breaking change).
- `TestContainsString` in graph/ — tested stdlib `strings.Contains` (zero value).

## [0.4.0] - 2026-05-17

### Added

- MIT license (replaced PROPRIETARY)
- README.md rewritten for general audience and public launch
- 27 doc comments on exported symbols across `d2_enum.go`, `color.go`, `graph.go`, `sort.go`, `enum/enum.go`
- `Shape` type with `ShapeTable`/`ShapeTree`/`ShapeGraph` constants
- `formatCapabilities` map — single source of truth for format-to-shape mapping
- `Format.Supports(Shape)`, `Format.Shapes()`, `FormatsForShape(Shape)` methods
- `ParseShape()`, `Shape.IsValid()`, `Shape.AllowedValues()`, `Shape.String()` enum methods
- `ErrInvalidShape` sentinel error
- `AllShapes` slice for shape iteration
- ADR 002: Shape capability matrix decision record

### Changed

- `FormatCategory`, `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` deprecated — use `Supports(Shape)` instead
- `TreeNode` and `TreeOutputRenderer` extracted from `format.go` to `tree.go`
- `TableData`, `RowEdge`, and `tableDataBase` extracted from `format.go`/`html.go` to `tabledata.go`
- `GraphRendererState` moved from `dot.go` to `graph.go`
- `format.go` reduced from 373 to 291 lines (under 350-line limit)
- `dot.go` reduced from 253 to 199 lines

### Removed

- `PLAN.md` — stale duplicate of AGENTS.md with incorrect info

## [0.3.0] - 2026-05-13

### Added

- Multi-module workspace: `enum/`, `escape/`, `cmdguard/`, `table/`, `sort/`, `integration/`, `examples/`
- `enum` package with generic `Parse`, `Contains`, `AllowedStrings`, `AllowedValues` utilities
- `escape` package with format-specific HTML/XML escaping (using stdlib `html.EscapeString`)
- `cmdguard` package with generic `EnumFlag[T]` CLI flag parser
- `table` package with lipgloss-styled terminal tables (isolated from root)
- Branded ID types via `github.com/larsartmann/go-branded-id` for D2NodeID, TreeNodeID, etc.
- ADR 001: Multi-module workspace decision record
- `go.work` support for local development (gitignored)
- `FormatTSV` constant and TSV formatter implementation
- `FormatXML` constant and XML formatter (`MarshalXML`, `MarshalXMLIndent`, `XMLWriter`, `MarshalXMLFromTableData`)

### Changed

- `ParseFormat`, `ParseSortBy`, `ParseColorMode` refactored to use `enum` helpers
- `AllowedValues()` methods refactored to use `enum` helpers
- `escape/` uses `html.EscapeString()` from stdlib instead of custom implementation
- Go toolchain updated to 1.26+
- `sort/` deprecated with notice pointing to `slices.SortStableFunc` + `cmp.Compare`
- Format classification uses map-based lookup instead of switch statements

### Deprecated

- `OutputFormat` type alias — will be removed in v2.0
- `OutputFormat*` constants — will be removed in v2.0
- `ParseOutputFormat()` function — will be removed in v2.0
- `sort.Sorter[T]` — use stdlib `slices.SortStableFunc` instead

## [0.2.0] - 2026-04-30

### Changed

- Updated dependencies (`charmbracelet/x/exp/golden`)

## [0.1.0] - 2026-01-01

### Added

- Initial release with 12 output formats: Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT
- `Renderer` interface with `Render() (string, error)`
- `TableRenderer` interface with `SetHeaders`/`AddRow`
- `GraphRenderer` interface with `SetNodes`/`SetEdges`
- `TreeOutputRenderer` interface with `SetRoot`
- `StreamingRenderer` interface with `Stream(io.Writer)`
- `ColorMode` enum (Auto/Always/Never) with terminal detection
- `SortBy` enum for sort field selection
- Opt-in renderer registry (`Register`/`Create`)
- D2 diagram renderer with shapes, arrows, SQL tables, classes, user journeys
- DOT/Graphviz renderer with `GraphRendererState`
- Mermaid flowchart renderer
- HTML table renderer with escaping
- Streaming HTML renderer for large datasets
- `MustRender(r)` helper for tests/examples
