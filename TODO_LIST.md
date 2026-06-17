# TODO_LIST.md — go-output

**Last updated:** 2026-06-17
**Open items:** 14
**Blocked:** 1 (needs owner decision)

All items below are **verified open** against the current code (2026-06-17).

---

## P1 — Correctness & Data Integrity (non-breaking)

| #   | Task                                                                                                                                                                                                                                                                                                                     | Effort | Status |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 20  | **Fix timing cache goroutine race** — `nom.TimingCache.Record()` spawns a new goroutine per call (`go tc.saveAsync()`). Burst of N completions spawns N goroutines racing to `os.Create`+write the same CSV path → file corruption. Debounce via single goroutine + channel, or serialize with mutex inside `saveAsync`. | Medium | Open   |
| 21  | **Fix inline renderer setter data race** — `nom.InlineRenderer.SetHideCursor/SetNoColor/SetAppName/SetStartTime` mutate fields with no synchronization while `refreshLoop` goroutine reads them. Guard with mutex or atomic types.                                                                                       | Medium | Open   |
| 22  | **Fix TUI reporter data race** — `BubbleTeaProgressReporter` mutates `pr.model.*` fields (workflowState, currentProgress, steps) from the caller goroutine while Bubble Tea's event loop reads/writes the same fields via `Update()`. Route all model mutations through `tea.Program.Send`.                              | High   | Open   |
| 23  | **Cap timing cache entries on load** — `loadLocked` does not apply `maxCachedEntries` cap during load. A hand-edited file with 10,000 entries for one activity loads in full. Apply cap during load.                                                                                                                     | Low    | Open   |
| 24  | **Fix tree re-parenting phantom edges** — `nom.DependencyTree.AddActivity` re-adding a node with different deps overwrites `Parent` but never removes the node from the OLD parent's `Children` slice.                                                                                                                   | Low    | Open   |

## P2 — Type Safety Improvements (non-breaking)

| #   | Task                                                                                                                                                                                                                                      | Effort | Status |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 25  | **`ProgressStep.IsActive` → derived method** — currently a mutable bool allowing impossible states (`CompletedAt != nil && IsActive == true`). Make it a method: `func (s ProgressStep) IsActive() bool { return s.CompletedAt == nil }`. | Low    | Open   |
| 26  | **`DisplayMode` string → int enum** — `type DisplayMode string` is stringly-typed. Change to `type DisplayMode int` with iota constants, matching `WorkflowState` pattern.                                                                | Low    | Open   |
| 27  | **`EdgeStyle.Style` string → typed enum** — free-form string ("solid", "dashed", "dotted") should be a typed enum to prevent invalid values.                                                                                              | Medium | Open   |
| 28  | **`GetActivityCounts()` struct return** — returns 4 unnamed `int` values; callers can swap them. Return an `ActivityCounts` struct.                                                                                                       | Low    | Open   |

## P3 — File Size & Structure (non-breaking)

| #   | Task                                                                                                                                                                    | Effort | Status |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 29  | **Split `tui/view.go`** — 381 lines (>350 limit). Split into `view.go` (View + viewport), `render_universal.go` (step rendering), `render_nom.go` (NOM tree rendering). | Low    | Open   |
| 30  | **Split `nom/inline_renderer.go`** — 393 lines (>350 limit). Extract summary rendering and lifecycle loop into separate files.                                          | Low    | Open   |

## P4 — Security Hardening (non-breaking)

| #   | Task                                                                                                                                               | Effort | Status |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 31  | **DOT nodesep/ranksep injection** — `SetNodeSep(sep string)` / `SetRankSep(sep string)` write raw user input into DOT output. Validate as numeric. | Low    | Open   |
| 32  | **DOT edge style attrs not escaped** — `edge.Style.Color` and `edge.Style.Style` written raw in `writeEdge()`. Escape or validate against enum.    | Low    | Open   |

---

## P5 — Naming & Polish (low value, carried from prior audit)

| #   | Task                                                            | Effort | Status |
| --- | --------------------------------------------------------------- | ------ | ------ |
| 7   | **`nom`: rename `GetOperationSymbol` → `OperationSymbol`**      | Low    | Open   |
| 8   | **`examples/shared`: rename `HandleError` → `Must`**            | Low    | Open   |
| 9   | **Make `ColorModeAuto.ShouldColor()` deterministic/testable**   | Low    | Open   |
| 10  | **Cache `TreeNode.Depth()`** — currently O(n) parent-chain walk | Low    | Open   |
| 11  | **Add bounds validation for `D2NodeStyle.Opacity`**             | Low    | Open   |

---

## P6 — Deferred (post-v1, blocked by ADR 006 API freeze)

These have in-code `TODO(split-brain ...)` markers for traceability.

| #   | Task                                                                                                                                    | Why deferred |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| 12  | **Unify `Marshaler` → `Renderer` terminology** — registry types vs everything-else `Renderer`                                           | Breaking     |
| 13  | **`RenderOptions.GraphID` is dead code** — no marshaler reads it                                                                        | Breaking     |
| M4  | **Rename non-canonical Render methods** — `InlineRenderer.Render()` → `Draw()`, `DependencyTree.Render(h)` → `Format(h)` (in-code TODO) | Breaking     |
| M5  | **Rename `GraphShape` constants** — `ShapeBox` → `NodeShapeBox` etc. to disambiguate from data-capability `Shape` enum (in-code TODO)   | Breaking     |
| M6  | **Document canonical branded-ID import path** — D2NodeID reachable via 3 import paths (in-code TODO)                                    | Breaking     |
| M7  | **Add `output.Direction` enum** — bridge D2Direction ↔ RankDir (in-code TODO)                                                           | Breaking     |
| M8  | **Align style struct field names** — `GraphStyle.FillColor` ↔ `D2NodeStyle.Fill` (in-code TODO)                                         | Breaking     |

### Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                                          |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only. Option B: unexported + validated setters (impossible-state-proof). Option C: keep both for v0.x. Affects every consumer; v1 stability commitment. |

### Release / Community

| #   | Task                                                                         | Effort | Status |
| --- | ---------------------------------------------------------------------------- | ------ | ------ |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                        | Low    | Open   |
| 16  | **Cut `v1.0.0` tag** — API declared frozen/ready (ADR 006), still at v0.12.x | Low    | Open   |

---

## Resolved This Session (2026-06-17) — Do Not Redo

| Task                                                                | Resolution                                                               |
| ------------------------------------------------------------------- | ------------------------------------------------------------------------ |
| `FormatDuration` displaying "60.0s" at 59.95s                       | **Fixed** — integer math (tenths) prevents float rounding                |
| DOT/Mermaid double-escaping in TableData conversion                 | **Fixed** — removed pre-escape; render-time escape is single source      |
| D2 Icon/Link fields not escaped                                     | **Fixed** — `escape.D2()` applied                                        |
| PlantUML labels not escaped                                         | **Fixed** — added `escape.PlantUML()`, applied to all label emissions    |
| Mermaid node/edge IDs not sanitized in graph Render()               | **Fixed** — `escape.MermaidID()` applied (was already done in tree path) |
| GraphEdge.Style and GraphStyle.FontColor dropped in D2 conversion   | **Fixed** — edge color and font-color now mapped                         |
| d2/doc.go example referenced nonexistent constructors               | **Fixed** — updated to AddNodeSimple/AddEdgeSimple                       |
| serialization/doc.go referenced nonexistent Marshal functions       | **Fixed** — updated to actual API entry points                           |
| `GetActivitySummaryString` returned bracketed slice repr            | **Fixed** — `strings.Join` instead of `fmt.Sprintf("%s", parts)`         |
| TUI `ProgressModel.messages` unbounded memory leak (never rendered) | **Fixed** — field and append removed                                     |
| Split-brain C1: duplicate `TreeNode` type                           | **Fixed** — renamed to `nom.ActivityNode`                                |
| Split-brain C2: ProgressModel duplicating subscriber state          | **In progress** — concurrent session (TUI module mid-refactor)           |
| Split-brain C3: `TimingFormat` constant divergence                  | **Fixed**                                                                |
| Split-brain C4/C5: test-only interface redeclarations               | **Fixed** — use canonical `output.GraphRenderer`/`output.Renderer`       |
| Split-brain M1: `ColorWarning` duplicate of `ColorRunning`          | **Fixed** — unused duplicate deleted                                     |
| Split-brain M2: divergent color detection logic                     | **Fixed** — aligned to union of env checks                               |
| Split-brain M3: hardcoded "No activities to display"                | **Fixed** — uses `MsgNoActivitiesToDisplay` constant                     |
| Split-brain M9: example `delimitedWriter` interface drift           | **Fixed**                                                                |
| Split-brain m2: bare event string literals                          | **Fixed** — replaced with `Event*` constants                             |
| Split-brain m4: stale GraphEdge in FORMAT_ARCHITECTURE.md           | **Fixed** — added `Style EdgeStyle` field                                |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-17 Split-Brain Sprint:** 20 findings audited. 14 fixed (6 critical/major + 8 minor). 7 deferred as API-breaking TODOs for v1.
- **2026-06-15 Full Code Review:** nom internal decomposition planned, BDD module added (19 Ginkgo specs), zero lint achieved, all wsl_v5 whitespace fixed.
- **2026-06-08 Architecture & Naming Sprint:** SlugifyID extracted, GraphRendererMixin→State, TableDataBase→Store, DTO suffix removed, formatCapabilities inverted, HTML table unified, html/template adopted, escape.D2/MermaidText optimized, AsciiDoc escaping completed, lipgloss style cached.
- **2026-05-28 Round 6:** Footer row feature, pre-v1 API stability audit (ADR 006, 228 symbols frozen), round-trip integration tests, root coverage 82%→96%.
- **2026-05-25 Modularization:** D2/graph/table/delimited/serialization/markup/plantuml extraction, JSONL+AsciiDoc+TOML+PlantUML added (12→16 formats), Shape capability matrix (ADR 002), zero transitive deps, deduplication to 0 actionable clones.
