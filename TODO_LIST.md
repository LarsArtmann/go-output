# TODO_LIST.md — go-output

**Last updated:** 2026-06-18
**Open items:** 10
**Blocked:** 1 (needs owner decision)

All items below are **verified open** against the current code (2026-06-18).

---

## P1 — Correctness & Concurrency (non-breaking)

| #   | Task                                                                                                                                                                                                                                                | Effort | Status |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 22  | **Fix TUI reporter data race** — `BubbleTeaProgressReporter` mutates `pr.model.*` fields from the caller goroutine while Bubble Tea's event loop reads/writes the same fields via `Update()`. Route all model mutations through `tea.Program.Send`. | High   | Open   |

## P2 — Type Safety (non-breaking, but exported API)

| #   | Task                                                                                                                                                   | Effort | Status |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 27  | **`EdgeStyle.Style` string → typed LineStyle enum** — free-form string ("solid", "dashed", "dotted") should be a typed enum to prevent invalid values. | Medium | Open   |
| 28  | **`GetActivityCounts()` struct return** — returns 4 unnamed `int` values; callers can swap them. Return an `ActivityCounts` struct.                    | Low    | Open   |

## P3 — Polish & Naming (low value)

| #   | Task                                                          | Effort | Status |
| --- | ------------------------------------------------------------- | ------ | ------ |
| 8   | **`examples/shared`: rename `HandleError` → `Must`**          | Low    | Open   |
| 9   | **Make `ColorModeAuto.ShouldColor()` deterministic/testable** | Low    | Open   |

---

## P4 — Deferred (post-v1, blocked by ADR 006 API freeze)

These have in-code `NOTE(split-brain ...)` markers for traceability.

| #   | Task                                                                                                                                    | Why deferred |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| 12  | **Unify `Marshaler` → `Renderer` terminology** — registry types vs everything-else `Renderer`                                           | Breaking     |
| 13  | **`RenderOptions.GraphID` is dead code** — no marshaler reads it                                                                        | Breaking     |
| M4  | **Rename non-canonical Render methods** — `InlineRenderer.Render()` → `Draw()`, `DependencyTree.Render(h)` → `Format(h)` (in-code NOTE) | Breaking     |
| M5  | **Rename `GraphShape` constants** — `ShapeBox` → `NodeShapeBox` etc. to disambiguate from data-capability `Shape` enum (in-code NOTE)   | Breaking     |
| M6  | **Document canonical branded-ID import path** — D2NodeID reachable via 3 import paths (in-code NOTE)                                    | Breaking     |
| M7  | **Add `output.Direction` enum** — bridge D2Direction ↔ RankDir (in-code NOTE)                                                           | Breaking     |
| M8  | **Align style struct field names** — `GraphStyle.FillColor` ↔ `D2NodeStyle.Fill` (in-code NOTE)                                         | Breaking     |

### Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                 |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only. Option B: unexported + validated setters. Option C: keep both for v0.x. Affects every consumer; v1 stability commitment. |

### Release / Community

| #   | Task                                                                         | Effort | Status |
| --- | ---------------------------------------------------------------------------- | ------ | ------ |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                        | Low    | Open   |
| 16  | **Cut `v1.0.0` tag** — API declared frozen/ready (ADR 006), still at v0.12.x | Low    | Open   |

---

## Resolved This Session (2026-06-17/18) — Do Not Redo

| Task                                                                | Resolution                                                            |
| ------------------------------------------------------------------- | --------------------------------------------------------------------- |
| `FormatDuration` displaying "60.0s" at 59.95s                       | **Fixed** — integer math (tenths) prevents float rounding             |
| `GetActivitySummaryString` returned bracketed Go slice repr         | **Fixed** — `strings.Join` instead of `fmt.Sprintf("%s", parts)`      |
| DOT/Mermaid double-escaping in TableData conversion                 | **Fixed** — removed pre-escape; render-time escape is single source   |
| D2 Icon/Link fields not escaped                                     | **Fixed** — `escape.D2()` applied                                     |
| PlantUML labels not escaped                                         | **Fixed** — added `escape.PlantUML()`, applied to all label emissions |
| Mermaid node/edge IDs not sanitized in graph Render()               | **Fixed** — `escape.MermaidID()` applied                              |
| GraphEdge.Style and GraphStyle.FontColor dropped in D2 conversion   | **Fixed** — edge color and font-color now mapped                      |
| d2/doc.go example referenced nonexistent constructors               | **Fixed** — updated to AddNodeSimple/AddEdgeSimple                    |
| serialization/doc.go referenced nonexistent Marshal functions       | **Fixed** — updated to actual API entry points                        |
| TUI `ProgressModel.messages` unbounded memory leak (never rendered) | **Fixed** — field and append removed                                  |
| Timing cache goroutine race (unbounded saveAsync on same file)      | **Fixed** — saveMu serializes file writes (TODO #20)                  |
| Timing cache entries not capped on load                             | **Fixed** — maxCachedEntries applied during load (TODO #23)           |
| Timing cache parse asymmetry (Sscanf vs FormatInt)                  | **Fixed** — strconv.ParseInt for symmetry                             |
| Timing cache non-deterministic write order                          | **Fixed** — sorted keys before writing                                |
| InlineRenderer config setter data race                              | **Fixed** — tickMu (RWMutex) guards all setters and reads (TODO #21)  |
| Tree re-parenting phantom edges                                     | **Fixed** — removeChild from old parent before reassigning (TODO #24) |
| DOT nodesep/ranksep injection                                       | **Fixed** — numeric validation (TODO #31)                             |
| DOT edge style attrs not escaped                                    | **Fixed** — escape.DOT applied (TODO #32)                             |
| `GetOperationSymbol` redundant Get prefix                           | **Fixed** — renamed to `OperationSymbol` (TODO #7)                    |
| `TreeNode.Depth()` O(n) parent-chain walk                           | **Fixed** — cached field with subtree propagation (TODO #10)          |
| `D2NodeStyle.Opacity` no bounds validation                          | **Fixed** — clamped to [0.0, 1.0] (TODO #11)                          |
| `ProgressStep.IsActive` allowed impossible state                    | **Fixed** — derived method from CompletedAt (TODO #25)                |
| `DisplayMode` stringly-typed                                        | **Fixed** — int enum with iota (TODO #26)                             |
| `tui/view.go` over 350-line limit                                   | **Fixed** — split to render_nom.go (TODO #29)                         |
| `nom/inline_renderer.go` over 350-line limit                        | **Fixed** — split to inline_renderer_summary.go (TODO #30)            |
| Split-brain C1: duplicate `TreeNode` type                           | **Fixed** — renamed to `nom.ActivityNode`                             |
| Split-brain C2: ProgressModel duplicating subscriber state          | **Fixed** — activities field removed                                  |
| Split-brain C3: `TimingFormat` constant divergence                  | **Fixed**                                                             |
| Split-brain C4/C5: test-only interface redeclarations               | **Fixed** — use canonical types                                       |
| Split-brain M1: `ColorWarning` duplicate of `ColorRunning`          | **Fixed** — deleted                                                   |
| Split-brain M2: divergent color detection logic                     | **Fixed** — aligned to union of env checks                            |
| Split-brain M3: hardcoded "No activities to display"                | **Fixed** — uses constant                                             |
| Split-brain M9: example `delimitedWriter` interface drift           | **Fixed**                                                             |
| Split-brain m2: bare event string literals                          | **Fixed** — replaced with `Event*` constants                          |
| Split-brain m4: stale GraphEdge in FORMAT_ARCHITECTURE.md           | **Fixed** — added `Style EdgeStyle` field                             |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-17 Split-Brain Sprint + Full Code Review:** 20 findings audited + 13 code review items fixed. 35 commits total since v0.12.0.
- **2026-06-15 Full Code Review:** BDD module added (19 Ginkgo specs), zero lint achieved.
- **2026-06-08 Architecture & Naming Sprint:** SlugifyID, GraphRendererMixin→State, escape optimization.
- **2026-05-28 Round 6:** Footer row feature, API stability audit (ADR 006).
- **2026-05-25 Modularization:** 16 formats, Shape matrix, zero transitive deps.
