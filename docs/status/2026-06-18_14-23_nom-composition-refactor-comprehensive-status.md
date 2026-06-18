# Nom/ Composition Refactor — Comprehensive Status

**Date:** 2026-06-18 14:23
**Branch:** master (clean, 12 commits this sprint, pushed)
**Sprint Start:** 2026-06-18 13:35
**Plan:** [`docs/planning/2026-06-18_13-40_NOM-COMPOSITION-REFACTOR.md`](../planning/2026-06-18_13-40_NOM-COMPOSITION-REFACTOR.md)

---

## Executive Summary

The nom/ composition refactor is **foundationally complete** — nom/ now imports and embeds root's `output.GraphNode` types, diagram export of live progress state works, and `MultiSubscriber` enables event fanout. The bridge-sync regression was caught and fixed. The migration of `ActivityDisplayState` → `Activity` as primary store is the biggest remaining piece.

**Build:** ✅ All 15 modules compile
**Tests:** ✅ All modules pass
**Race:** ✅ nom + tui clean
**Coverage:** nom 91.4%, tui 90.3%

---

## A. FULLY DONE ✅

| # | Deliverable | Verification |
|---|---|---|
| 1 | **`nom.Activity` type** — embeds `output.GraphNode` (ID, Label, Shape, Style, Metadata) + temporal fields (Status, StartTime, EndTime, EstimatedTime, Err, Symbol, Color, CurrentElapsed, Dependencies) | 7 unit tests, 100% of new methods covered |
| 2 | **`ActivityStatus.NodeShape()` + `.GraphStyle()`** — maps domain status to root's `NodeShape` + `GraphStyle` for diagram export | 2 tests covering all 5 statuses |
| 3 | **`nom.ActivityStore`** — map-backed store with `Nodes() []output.GraphNode` / `Edges() []output.GraphEdge` projection | 10 tests including projection-copy isolation |
| 4 | **`ActivityNode` embeds `Activity`** — tree nodes carry branded GraphNode data; `DisplayState` embedding removed | All tree tests pass with new field accessors |
| 5 | **`nom.MultiSubscriber`** — `io.MultiWriter`-style fanout for `EventSubscriber` | 4 tests: fanout, error isolation, nil-skip, Add() |
| 6 | **`ActivityReader` interface** — `Nodes()/Edges()` contract; both subscriber and standalone store satisfy it | Used by `Store()` adapter |
| 7 | **`subscriberView` on-demand projection** — subscriber projects `Nodes()/Edges()` from existing state; no third copy, no bridge sync | 3 diagram export integration tests |
| 8 | **`NOMStyleSubscriber.Store()`** — returns `ActivityReader` for diagram export | Tested end-to-end with event flow |
| 9 | **ADR 007** — decision documented with trade-offs, implementation status table | [`docs/adr/007-nom-composition-via-root-types.md`](../adr/007-nom-composition-via-root-types.md) |
| 10 | **CHANGELOG updated** — new types, breaking changes documented | [`CHANGELOG.md`](../../CHANGELOG.md) |
| 11 | **Build breakage fixed** — tui/ + integration/ updated for `ActivityNode.ActivityID` removal | All 15 modules build green |
| 12 | **Bridge sync eliminated** — removed 3rd state copy (`syncToStore`); subscriber reads primary state directly | Zero sync calls in event handlers |
| 13 | **Comprehensive plan** — 22 medium tasks, 78 fine tasks, mermaid graph, risk register | [`docs/planning/`](../planning/2026-06-18_13-40_NOM-COMPOSITION-REFACTOR.md) |

---

## B. PARTIALLY DONE 🟡

| # | Item | Status | What's Left |
|---|---|---|---|
| 1 | **`ActivityNode` migration** | Activity embedded (not DisplayState), but `ActivityDisplayState` still exists as subscriber's primary store | Delete `DisplayState` struct; migrate subscriber to store `*Activity` directly |
| 2 | **Diagram export feature** | `ActivityReader` interface + projection works; tests prove it | No real runnable example yet (`examples/nom_progress/diagram_export.go`) |
| 3 | **`nom/` → root dependency** | go.mod updated, builds clean | `go.sum` may need tidy across workspace |
| 4 | **Test coverage of new types** | 26 new tests pass | `ActivityDisplayState` tests still test old type; need consolidation |

---

## C. NOT STARTED 🔲

| # | Item | Impact | Effort | Notes |
|---|---|---|---|---|
| 1 | **Delete `ActivityDisplayState`** entirely — subscriber stores `*Activity` | 🔥🔥🔥 | 3h | The big one. Eliminates `SyncActivityTimingToTree`, `DisplayState`, the bridge sync function, ~150 LOC of duplicated state |
| 2 | **Fix `InlineRenderer.Render()` contract** (split-brain M4) | 🔥🔥 | 30m | Rename to `Draw()`, add `Render() (string, error)`. Already documented as TODO |
| 3 | **Real diagram export example** | 🔥🔥 | 45m | `examples/nom_progress/diagram_export.go` showing DOT/Mermaid export |
| 4 | **Update `FORMAT_ARCHITECTURE.md`** nom section | 🔥 | 30m | Document Activity embedding + diagram export |
| 5 | **Update `AGENTS.md`** nom patterns | 🔥 | 30m | New types, new file structure |
| 6 | **Update `FEATURES.md`** — diagram export + MultiSubscriber | 🔥 | 15m | New feature entries |
| 7 | **Update `DOMAIN_LANGUAGE.md`** — Activity, ActivityStore terms | 🟡 | 15m | Vocabulary update |
| 8 | **`ActivityStore` integration into DependencyTree** | 🟡 | 2h | Tree could derive rendering from store instead of parallel nodes map |
| 9 | **Performance benchmark** — projection overhead | 🟡 | 1h | Verify on-demand projection is acceptable for render-tick frequency |
| 10 | **Theme injection** — custom symbols/colors | 🟢 | 2h | Low priority, speculative |

---

## D. TOTALLY FUCKED UP 💥 (and Fixed)

| # | What Happened | Severity | How Fixed |
|---|---|---|---|
| 1 | **Broke build for tui/ + integration/** — removed `ActivityNode.ActivityID` field without updating 4 downstream call sites | 🔴 Critical | Fixed all references: `node.ActivityID` → `nom.ActivityID(node.ID.Get())` |
| 2 | **VERSCHLIMMBESSER: Created 3rd state copy** — added `ActivityStore` as bridge-synced parallel state, making split-brain WORSE (2→3 copies, 1→2 syncs) | 🔴 Critical | Eliminated entirely: subscriber now implements `ActivityReader` directly via `subscriberView` adapter that projects on-demand from existing state |
| 3 | **Didn't check existing code for reuse** — invented `ActivityStore` as a new abstraction when the subscriber could satisfy `Nodes()/Edges()` directly | 🟡 Design flaw | `ActivityStore` kept as standalone utility (useful for testing), but subscriber doesn't use it internally |
| 4 | **Allocated new `Activity` per event** in bridge sync — performance regression on every `handleActivityStarted` | 🟡 Perf | Eliminated with bridge sync removal; projection allocates only on diagram export (infrequent) |

---

## E. WHAT WE SHOULD IMPROVE

### Architecture

1. **Kill `ActivityDisplayState` + `DisplayState` entirely.** This is the root cause of all remaining split-brain. The subscriber should store `map[ActivityID]*Activity` directly. When this is done, `SyncActivityTimingToTree` vanishes, `syncActivityToNode` vanishes, and there is exactly ONE source of truth.

2. **Consider whether `ActivityStore` is even needed.** If the subscriber stores `*Activity` directly and implements `ActivityReader`, the standalone `ActivityStore` may be unnecessary abstraction. Keep it only if there's a use case for a store independent of event subscription (e.g., testing, replay).

3. **Fix `InlineRenderer.Render()` contract (M4).** This is a documented TODO from before this sprint. Renaming to `Draw()` and adding `Render() (string, error)` makes it conform to `output.Renderer`, enabling composition with other renderers.

### Process

4. **Run ALL module builds before committing.** The build breakage happened because I only ran `cd nom && go test`. I should run the full workspace build (`go work build` or iterate all modules) after ANY change to nom/ types that are consumed by tui/.

5. **Check for verschlimmbessern before adding abstractions.** The bridge sync looked reasonable in isolation but made the system worse. The test: "Does this reduce the total number of state copies + sync calls?" If no, don't add it.

6. **Reuse existing types before inventing new ones.** The initial proposal invented `ActivityStore`/`Hierarchy`/`DurationHistory` interfaces. The simpler answer was making the subscriber satisfy `ActivityReader` directly. Always ask: "Can existing code satisfy this contract?"

### Code Quality

7. **`subscriberView` accesses `dependencyTree.mu` directly** — this is a lock-ordering concern. If the tree's mutex and subscriber's mutex are ever held in different orders by different callers, it's a deadlock risk. Should be refactored to call tree methods that acquire their own lock.

8. **`Activity` struct is getting large** (12+ fields). Consider grouping temporal fields into a `Timing` sub-struct if more fields are added.

---

## F. TOP 25 THINGS TO GET DONE NEXT

**Sorted by impact/effort ratio (highest first).**

| # | Task | Impact | Effort | Category |
|---|---|---|---|---|
| 1 | Delete `ActivityDisplayState` → subscriber stores `*Activity` directly | 🔥🔥🔥 | 3h | Split-brain elimination |
| 2 | Eliminate `SyncActivityTimingToTree` (auto-follows from #1) | 🔥🔥🔥 | 0m | Split-brain elimination |
| 3 | Delete `DisplayState` struct (auto-follows from #1) | 🔥🔥🔥 | 0m | Split-brain elimination |
| 4 | Delete `syncActivityToNode` bridge function | 🔥🔥 | 0m | Split-brain elimination |
| 5 | Fix `InlineRenderer.Render()` → `Draw()` + `Render() (string, error)` | 🔥🔥 | 30m | Contract conformance |
| 6 | Add `examples/nom_progress/diagram_export.go` — real DOT export demo | 🔥🔥 | 45m | Feature showcase |
| 7 | Update `FORMAT_ARCHITECTURE.md` — Activity embedding + diagram export | 🔥 | 30m | Documentation |
| 8 | Update `AGENTS.md` — nom patterns, file structure | 🔥 | 30m | Documentation |
| 9 | Update `FEATURES.md` — diagram export + MultiSubscriber | 🔥 | 15m | Documentation |
| 10 | Migrate `activity_display_test.go` → test `Activity` directly | 🔥 | 1h | Test consolidation |
| 11 | Migrate `subscriber_test.go` — test with `*Activity` not `*ActivityDisplayState` | 🔥 | 1.5h | Test consolidation |
| 12 | Migrate `format_test.go` — `FormatTimingInfo` takes `*Activity` | 🟡 | 30m | Test consolidation |
| 13 | Migrate `configuration_test.go` — use new types | 🟡 | 15m | Test consolidation |
| 14 | Consolidate `activity_management.go` — remove dead code after #1 | 🟡 | 30m | Code cleanup |
| 15 | Refactor `subscriberView.Edges()` to avoid direct tree.mu access | 🟡 | 30m | Lock safety |
| 16 | Add `DependencyTree.SetStore(*ActivityStore)` option for store-backed mode | 🟡 | 1h | Architecture option |
| 17 | Add golden test for DOT diagram export output | 🟡 | 30m | Test quality |
| 18 | Benchmark: projection overhead vs bridge sync | 🟢 | 1h | Performance verification |
| 19 | Add `RenderOptions` for diagram export (title, color theme) | 🟢 | 30m | Feature polish |
| 20 | Consider `Theme` struct for custom symbols/colors | 🟢 | 2h | Extensibility |
| 21 | Add `examples/tui_progress/diagram_export.go` | 🟢 | 30m | Feature showcase |
| 22 | Integration test: full workflow → DOT diagram export | 🟡 | 45m | Integration coverage |
| 23 | Update `TODO_LIST.md` with remaining nom composition tasks | 🟢 | 15m | Planning |
| 24 | Run `go mod tidy` across all 16 modules | 🟢 | 10m | Hygiene |
| 25 | Final lint sweep across all modules | 🟢 | 10m | Hygiene |

---

## G. TOP #1 QUESTION

**Should `ActivityStore` exist as a standalone type, or should it be removed in favor of the subscriber implementing `ActivityReader` directly?**

**Context:** `ActivityStore` was originally created as the "projection layer for diagram export." But after eliminating the bridge sync, the subscriber now implements `ActivityReader` directly via `subscriberView`. The standalone `ActivityStore` is:
- Used in 3 tests (`TestDiagramExport_StatusShapes`, `TestDiagramExport_EdgeStructure`, and its own test suite)
- NOT used by the subscriber internally
- NOT used by tui/
- A candidate for YAGNI unless there's a real use case

**Arguments for keeping it:**
- Useful for testing (build a store without a full subscriber)
- Useful for replay scenarios (load activities from a log file without events)
- Clean separation of storage vs event handling

**Arguments for removing it:**
- YAGNI — no production code uses it outside tests
- The subscriber already provides the same projection
- Removing it simplifies the API surface

**I cannot decide this myself** because it depends on whether you envision use cases like "load progress state from a file" or "replay a recorded event stream into a standalone store." If those are real planned features, keep it. If not, delete it.
