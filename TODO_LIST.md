# TODO_LIST.md — go-output

**Last updated:** 2026-06-18
**Open items:** 3
**Blocked:** 1 (needs owner decision)

---

## Open Items

| #   | Task                                                                                                                                            | Effort | Status                      |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------ | --------------------------- |
| M4  | **Rename non-canonical Render methods** — `InlineRenderer.Render()` → `Draw()`, `DependencyTree.Render(h)` → `Format(h)`. Pure naming refactor. | Medium | Open                        |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                                                                                           | Low    | Open (needs owner account)  |
| 16  | **Cut `v1.0.0` tag** — API declared frozen/ready (ADR 006), still at v0.12.x                                                                    | Low    | Open (needs owner decision) |

### Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                 |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only. Option B: unexported + validated setters. Option C: keep both for v0.x. Affects every consumer; v1 stability commitment. |

---

## Resolved This Session (2026-06-17/18) — Do Not Redo

| Task                                                              | Resolution                                                                   |
| ----------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `FormatDuration` displaying "60.0s" at 59.95s                     | **Fixed** — integer math (tenths) prevents float rounding                    |
| `GetActivitySummaryString` returned bracketed Go slice repr       | **Fixed** — `strings.Join` instead of `fmt.Sprintf("%s", parts)`             |
| DOT/Mermaid double-escaping in TableData conversion               | **Fixed** — removed pre-escape; render-time escape is single source          |
| D2 Icon/Link fields not escaped                                   | **Fixed** — `escape.D2()` applied                                            |
| PlantUML labels not escaped                                       | **Fixed** — added `escape.PlantUML()`                                        |
| Mermaid node/edge IDs not sanitized in graph Render()             | **Fixed** — `escape.MermaidID()` applied                                     |
| GraphEdge.Style and GraphStyle.FontColor dropped in D2 conversion | **Fixed** — edge color and font-color now mapped                             |
| Timing cache goroutine race                                       | **Fixed** — saveMu serializes file writes                                    |
| Timing cache entries not capped on load                           | **Fixed** — maxCachedEntries applied during load                             |
| InlineRenderer config setter data race                            | **Fixed** — RWMutex guards all setters and reads                             |
| Tree re-parenting phantom edges                                   | **Fixed** — removeChild from old parent                                      |
| DOT nodesep/ranksep injection                                     | **Fixed** — numeric validation                                               |
| DOT edge style attrs not escaped                                  | **Fixed** — escape.DOT applied                                               |
| `GetOperationSymbol` redundant Get prefix                         | **Fixed** — renamed to `OperationSymbol`                                     |
| `TreeNode.Depth()` O(n) parent-chain walk                         | **Fixed** — cached field with subtree propagation                            |
| `D2NodeStyle.Opacity` no bounds validation                        | **Fixed** — clamped to [0.0, 1.0]                                            |
| `ProgressStep.IsActive` allowed impossible state                  | **Fixed** — derived method from CompletedAt                                  |
| `DisplayMode` stringly-typed                                      | **Fixed** — int enum with iota                                               |
| `tui/view.go` over 350-line limit                                 | **Fixed** — split to render_nom.go                                           |
| `nom/inline_renderer.go` over 350-line limit                      | **Fixed** — split to inline_renderer_summary.go                              |
| Split-brain C1-C5, M1-M3, M9, m2                                  | **Fixed** — all resolved                                                     |
| TUI reporter data race                                            | **Fixed** — reporter owns workflowState; all mutations via send() (TODO #22) |
| `RenderOptions.GraphID` dead code                                 | **Removed** (TODO #13)                                                       |
| `EdgeStyle.Style` free-form string                                | **Fixed** — typed `LineStyle` enum; field renamed to `Line` (TODO #27)       |
| `GetActivityCounts()` 4 unnamed int returns                       | **Fixed** — returns `ActivityCounts` struct (TODO #28)                       |
| `ColorModeAuto.ShouldColor()` not testable                        | **Fixed** — overridable detection vars; deterministic test (TODO #9)         |
| `GraphShape` constants collide with `Shape` enum                  | **Fixed** — renamed to `NodeShape`/`NodeShapeBox` etc. (M5)                  |
| `GraphStyle.FillColor`/`StrokeColor` diverge from D2NodeStyle     | **Fixed** — renamed to `Fill`/`Stroke` (M8)                                  |
| No bridge between D2Direction and RankDir                         | **Fixed** — added `output.Direction` (M7)                                    |
| D2NodeID canonical import path undocumented                       | **Fixed** — doc comment added (M6)                                           |
| `Marshaler` terminology inconsistent with `Renderer`              | **Fixed** — renamed to `TableDataRenderer`/`AnyDataRenderer` (#12)           |
| #8: `HandleError` → `Must` suggestion                             | **Won't-fix** — `HandleError` is honest; `Must` implies panic                |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-18 v1 Preparation Sprint:** 11 TODO items resolved.
- **2026-06-17 Split-Brain Sprint + Full Code Review:** 20 findings + 13 review items fixed.
- **2026-06-15 Full Code Review:** BDD module added, zero lint achieved.
- **2026-06-08 Architecture & Naming Sprint:** SlugifyID, escape optimization.
- **2026-05-28 Round 6:** Footer row feature, API stability audit (ADR 006).
- **2026-05-25 Modularization:** 16 formats, Shape matrix, zero transitive deps.
