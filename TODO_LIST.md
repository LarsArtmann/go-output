# TODO_LIST.md — go-output

**Last updated:** 2026-06-20
**Open items:** 2
**Blocked:** 0

---

## Open Items

| #   | Task                                                                         | Effort | Status                      |
| --- | ---------------------------------------------------------------------------- | ------ | --------------------------- |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                        | Low    | Open (needs owner account)  |
| 16  | **Cut `v1.0.0` tag** — API frozen (ADR 006); CHANGELOG + full checklist done | Low    | Prepared — awaiting owner tag |

---

## Resolved This Session (2026-06-20) — Do Not Redo

| Task                                                              | Resolution                                                                                                                                         |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| #15 `TableData` fields vs getters for v1                          | **Resolved — keep both.** Getters (`GetHeaders`/`GetRows`/`GetFooter`) satisfy the `TableDataProvider` interface in `table/` + test doubles. Load-bearing, not duplication. |
| `ActivityStatusPaused` / `SetPaused` (dead event path)           | **Removed.** No `EventActivityPaused`, zero callers, unreachable. `SymbolPaused`→`SymbolPending` (glyph ⏸→○); Pending given its own honest identity. |
| Deprecated APIs (7 markers)                                       | **Removed.** `EnsureBuild`, `ParseActivityID`/`ParseWorkflowID`, `NewGraphNodeID`/`NewGraphNodeLabel`, `MarshalTSV(any)`, 6 color aliases — all had zero prod callers. Zero `Deprecated:` markers remain. |
| `OperationSymbol` + `OperationType*` (speculative)               | **Removed.** Zero production callers; stringly-typed mapping to nowhere. `SymbolDownload`/`Upload` kept as palette. |
| Over-exposed tui public API (~15 symbols)                         | **Unexported.** Msg types, `WorkflowState`, `ProgressStep`, `UpdateType`, `TickCmd`, etc. Public surface is now `NewBubbleTeaProgressReporter` + `Report*` + `DisplayMode` + `ProgressModel`. |
| Pre-release checklist (build/test/race/lint/vuln/tidy)           | **All green.** 20/20 modules pass, 0 lint issues, `-race` clean, govulncheck 0 vulnerabilities, go mod tidy clean. |

---

## Resolved Earlier (Summary — details in git history)

| Task                                                              | Resolution                                                                                                                                         |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| `FormatDuration` displaying "60.0s" at 59.95s                     | **Fixed** — integer math (tenths) prevents float rounding                                                                                          |
| `GetActivitySummaryString` returned bracketed Go slice repr       | **Fixed** — `strings.Join` instead of `fmt.Sprintf("%s", parts)`                                                                                   |
| DOT/Mermaid double-escaping in TableData conversion               | **Fixed** — removed pre-escape; render-time escape is single source                                                                                |
| D2 Icon/Link fields not escaped                                   | **Fixed** — `escape.D2()` applied                                                                                                                  |
| PlantUML labels not escaped                                       | **Fixed** — added `escape.PlantUML()`                                                                                                              |
| Mermaid node/edge IDs not sanitized in graph Render()             | **Fixed** — `escape.MermaidID()` applied                                                                                                           |
| GraphEdge.Style and GraphStyle.FontColor dropped in D2 conversion | **Fixed** — edge color and font-color now mapped                                                                                                   |
| Timing cache goroutine race                                       | **Fixed** — saveMu serializes file writes                                                                                                          |
| Timing cache entries not capped on load                           | **Fixed** — maxCachedEntries applied during load                                                                                                   |
| InlineRenderer config setter data race                            | **Fixed** — RWMutex guards all setters and reads                                                                                                   |
| Tree re-parenting phantom edges                                   | **Fixed** — removeChild from old parent                                                                                                            |
| DOT nodesep/ranksep injection                                     | **Fixed** — numeric validation                                                                                                                     |
| DOT edge style attrs not escaped                                  | **Fixed** — escape.DOT applied                                                                                                                     |
| `GetOperationSymbol` redundant Get prefix                         | **Fixed** — renamed to `OperationSymbol`                                                                                                           |
| `TreeNode.Depth()` O(n) parent-chain walk                         | **Fixed** — cached field with subtree propagation                                                                                                  |
| `D2NodeStyle.Opacity` no bounds validation                        | **Fixed** — clamped to [0.0, 1.0]                                                                                                                  |
| `ProgressStep.IsActive` allowed impossible state                  | **Fixed** — derived method from CompletedAt                                                                                                        |
| `DisplayMode` stringly-typed                                      | **Fixed** — int enum with iota                                                                                                                     |
| `tui/view.go` over 350-line limit                                 | **Fixed** — split to render_nom.go                                                                                                                 |
| `nom/inline_renderer.go` over 350-line limit                      | **Fixed** — split to inline_renderer_summary.go                                                                                                    |
| Split-brain C1-C5, M1-M3, M9, m2                                  | **Fixed** — all resolved                                                                                                                           |
| TUI reporter data race                                            | **Fixed** — reporter owns workflowState; all mutations via send() (TODO #22)                                                                       |
| `RenderOptions.GraphID` dead code                                 | **Removed** (TODO #13)                                                                                                                             |
| `EdgeStyle.Style` free-form string                                | **Fixed** — typed `LineStyle` enum; field renamed to `Line` (TODO #27)                                                                             |
| `GetActivityCounts()` 4 unnamed int returns                       | **Fixed** — returns `ActivityCounts` struct (TODO #28)                                                                                             |
| `ColorModeAuto.ShouldColor()` not testable                        | **Fixed** — overridable detection vars; deterministic test (TODO #9)                                                                               |
| `GraphShape` constants collide with `Shape` enum                  | **Fixed** — renamed to `NodeShape`/`NodeShapeBox` etc. (M5)                                                                                        |
| `GraphStyle.FillColor`/`StrokeColor` diverge from D2NodeStyle     | **Fixed** — renamed to `Fill`/`Stroke` (M8)                                                                                                        |
| No bridge between D2Direction and RankDir                         | **Fixed** — added `output.Direction` (M7)                                                                                                          |
| D2NodeID canonical import path undocumented                       | **Fixed** — doc comment added (M6)                                                                                                                 |
| `Marshaler` terminology inconsistent with `Renderer`              | **Fixed** — renamed to `TableDataRenderer`/`AnyDataRenderer` (#12)                                                                                 |
| M4: `Render()` name collision (InlineRenderer/DependencyTree)     | **Fixed** — `Draw()` / `RenderString`/`RenderWithWidth`; zero NOTE markers                                                                         |
| `reflect` depguard violation in integration tests                 | **Fixed** — replaced with `fmt.Sprintf("%T")` (no reflect import)                                                                                  |
| Stale `getTableDataMarshaler` naming                              | **Fixed** — renamed to `getTableDataRenderer`/`getAnyDataRenderer`                                                                                 |
| Nom timing cache tests wrote to real `~/.cache`                   | **Fixed** — `newTempTimingCache(t)` isolates direct cache tests; `WithCachePath` + `newTestSubscriber(t)` isolate all subscriber/integration tests |
| #8: `HandleError` → `Must` suggestion                             | **Won't-fix** — `HandleError` is honest; `Must` implies panic                                                                                      |
| O8: `ActivityStore` YAGNI                                         | **Resolved** — removed as ghost system (155 LOC dead prod code)                                                                                    |
| nom/tui/bdd/envdetect missing from CI                             | **Fixed** — added to all ci.yml + release.yml loops                                                                                                |
| Dead fields on `Activity` (Dependencies, OperationType)           | **Fixed** — removed unused fields + methods                                                                                                        |
| Stale `mustUpdateActivityStatus` comment + unused params          | **Fixed** — updated comment, dropped 3 unused params                                                                                               |
| Lock-ordering protocol undocumented                               | **Fixed** — documented `ns.mu → tree.mu` on subscriberView.Edges()                                                                                 |
| ADR 007 `tui/` migration status stale                             | **Fixed** — marked done (tests already use new types)                                                                                              |
| Test files over 350-line limit                                    | **Fixed** — all split under 350 lines (incl. `render_tabledata_test.go` → `render_registry_test.go`)                                               |
| No local govulncheck                                              | **Fixed** — added `nix run .#govulncheck` app                                                                                                      |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-18 v1 Preparation Sprint:** 11 TODO items resolved.
- **2026-06-17 Split-Brain Sprint + Full Code Review:** 20 findings + 13 review items fixed.
- **2026-06-15 Full Code Review:** BDD module added, zero lint achieved.
- **2026-06-08 Architecture & Naming Sprint:** SlugifyID, escape optimization.
- **2026-05-28 Round 6:** Footer row feature, API stability audit (ADR 006).
- **2026-05-25 Modularization:** 16 formats, Shape matrix, zero transitive deps.
