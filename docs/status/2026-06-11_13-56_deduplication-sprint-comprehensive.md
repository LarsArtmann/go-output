# Code Deduplication Sprint — 2026-06-11

**Session focus:** Eliminate code duplication at art-dupl thresholds 20 and 15, plus deep architectural audit.

---

## a) FULLY DONE

### Deduplication (22 files, -513/+376 net lines)

| Change | File(s) | Clones Eliminated |
|--------|---------|:-----------------:|
| Extract `ensureStartedAndActive()` | `tui/reporter.go` | 3 |
| Extract `snapshotData()` | `nom/timing_cache_persist.go` | 2 |
| Extract `stepIconAndStyle()` | `tui/view.go` | 3 |
| Extract `finalizeActivityExecution()` | `nom/subscriber_handlers.go` | 2 |
| Extract `setStatusWithElapsed()` | `nom/golden_test.go` | 8 |
| Extract `assertValidJSONLines()` | `serialization/jsonl_test.go` | 7 |
| Extract `assertAllContained()` | `markup/asciidoc_test.go`, `plantuml/plantuml_test.go`, `serialization/toml_renderers_test.go`, `serialization/toml_test.go` | 21 |
| Extract `assertScrollOffset()` | `tui/event_sequence_test.go` | 6 |
| Extract `addTestActivity()` + `setupTestTree()` | `tui/model_test.go`, `tui/view_test.go` | 6 |
| Extract `newProjectTableData()` | `integration/render_helpers_test.go` | 3 |
| Extract `newTestTimingCache()` | `nom/timing_cache_test.go` | 2 |
| Extract `assertErrWrite()` | `testhelpers/writers_test.go` | 3 |
| Convert `TestGetActivitySummaryString` → table-driven | `nom/format_test.go` | 3 |
| Convert tree benchmarks → table-driven sub-benchmarks | `nom/tree_bench_test.go` | 4 |
| Convert `TestBuildActivityCountsSummary` → table-driven | `tui/summary_test.go` | 2 |
| Merge `CanAcceptUpdates` + `CanAcceptTicks` | `tui/state_test.go` | 3 |
| Merge `ReportProgress` + `ReportStep` subtests | `tui/reporter_test.go` | 4 |
| Convert error tests → table-driven | `tabledatastore_test.go` | 3 |
| Reuse `createSummaryStyle()` | `tui/view.go` | 1 |

**Clone group counts:**

| Threshold | Before | After | Eliminated |
|-----------|:------:|:-----:|:----------:|
| t=20 | 29 | 22 | 7 (24%) |
| t=15 | 105 | 94 | 11 (10%) |

**All 15 modules pass tests.** Zero regressions.

### Architectural Audit (Deep Analysis)

Completed full codebase scan across 12 dimensions: TODOs, FIXMEs, nolint directives, registry patterns, TableData struct, Shape/Format matrix, NOM subscriber architecture, TUI state machine, escape overlaps, coverage gaps, untested exports, and examples freshness. Findings documented in sections d–g below.

---

## b) PARTIALLY DONE

### Test Coverage Improvements

Current coverage by module:

| Module | Coverage | Status |
|--------|:--------:|--------|
| enum | 100.0% | Done |
| escape | 100.0% | Done |
| d2 | 97.0% | Done |
| output (root) | 92.3% | Good |
| testhelpers | 91.3% | Good |
| nom | 91.9% | Good |
| delimited | 90.5% | Good |
| markup | 93.8% | Good |
| integration | 95.5% | Good |
| table | 88.4% | Could improve |
| tui | 87.6% | Could improve |
| graph | 86.8% | Should improve |
| serialization | 82.3% | Should improve |
| plantuml | 81.4% | Should improve |

**Lowest 3 modules (plantuml, serialization, graph) need focused test additions.**

---

## c) NOT STARTED

See section f) Top 25 Next Steps below.

---

## d) TOTALLY FUCKED UP

Nothing. All changes compiled, passed tests, and no regressions introduced.

---

## e) WHAT WE SHOULD IMPROVE

### Critical Architectural Issues (Found During Audit)

1. **Data race in `GetDependencyTree()`** — `nom/state_accessors.go:9-14` returns raw `*DependencyTree` pointer under brief RLock. Lock released before caller uses the tree. Concurrent mutation + read = data race.

2. **No lock in `GetTimingCache()`** — `nom/state_accessors.go:17-19` returns `*TimingCache` with zero synchronization. Other accessors use RLock.

3. **Dual-state sync between `ActivityDisplayState` and `TreeNode`** — `nom/activity_display.go` and `nom/tree.go` duplicate `Status`, `Symbol`, `Color`, `StartTime`, `EstimatedTime`, `CurrentElapsed`. Manual `SyncActivityTimingToTree()` keeps them in sync. Fragile — adding a field to one and forgetting the other causes subtle bugs.

4. **Two near-identical registries** — `render_tabledata.go` has `tableDataMarshalers` (lines 28-51) and `anyDataMarshalers` (lines 151-165) with duplicated Register/Get patterns. A generic `Registry[T]` could eliminate this.

5. **Unused parameter** — `tui/view.go:338` `renderHelpOverlay(content string)` never uses `content`.

6. **Magic number** — `tui/model.go:99` uses `scrollOffset = 9999` for scroll-to-bottom instead of computing from content height.

### Type Model Improvements

- **`TableData.AddRow` doesn't validate column count** — silently accepts mismatched row lengths. Should validate against `len(Headers)`.
- **`Shape` registration decoupled from marshaler registration** — shape.go registers all 16 formats in root init(), but marshalers are in sub-modules. Can query shape support for a format with no marshaler.
- **String-based event routing** — nom subscriber uses string comparison for event types. Typos silently go to default case returning nil.

### Library Opportunities

- **`slices.Contains`** — `nom/tree.go:82-89` has a manual loop that could use stdlib `slices.Contains`.
- **`min()` builtin** — `tui/view.go:62` can use Go 1.21+ `min()` instead of if-statement.
- **`slices.ContainsFunc`** — Several manual search loops across the codebase could use this.

---

## f) Top 25 Things We Should Get Done Next

Sorted by impact / effort (highest ROI first):

### P0 — Bugs / Safety (Must Fix)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 1 | Fix data race: `GetDependencyTree()` should return copy or document unsafe | 30min | Critical |
| 2 | Add lock to `GetTimingCache()` | 5min | Critical |
| 3 | Fix unused param `renderHelpOverlay(content string)` → `renderHelpOverlay()` | 10min | Cleanup |
| 4 | Replace 4× `//nolint:errcheck` in `integration/nom_tui_test.go` with `require.NoError` | 10min | Correctness |
| 5 | Replace `scrollOffset = 9999` with computed `maxOffset` | 15min | Correctness |

### P1 — Architecture (High Value)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 6 | Consolidate `ActivityDisplayState` + `TreeNode` timing fields into single source of truth | 2h | Architecture |
| 7 | Generic `Registry[T]` to deduplicate two marshaler registries | 1h | DRY |
| 8 | Add `TableData.AddRow` column count validation | 15min | Robustness |
| 9 | Modernize `min()` in `tui/view.go:62` | 5min | Style |
| 10 | Use `slices.Contains` in `nom/tree.go:82-89` | 5min | Style |

### P2 — Test Coverage (Closing Gaps)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 11 | Add tests for `plantuml` (81.4% → 90%+) | 30min | Coverage |
| 12 | Add tests for `serialization` (82.3% → 90%+) | 30min | Coverage |
| 13 | Add tests for `graph` (86.8% → 90%+) | 20min | Coverage |
| 14 | Add tests for `RenderAnyData()`, `RegisteredTableDataFormats()`, `RegisteredAnyDataFormats()` | 20min | Coverage |
| 15 | Add direct tests for `tui/Subscriber()`, `SetCancelFunc()`, `SetDisplayMode()` | 15min | Coverage |

### P3 — Further Deduplication (Diminishing Returns)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 16 | Extract `assertErrorContains` helper across integration/serialization/root error tests | 20min | Test DRY |
| 17 | Convert remaining `nom/types_test.go` ActivityID/WorkflowID to parameterized tests | 20min | Test DRY |
| 18 | Extract `assertFooterLineCountContains` across delimited/integration/root footer tests | 30min | Test DRY |
| 19 | Extract round-trip parse helpers in `integration/roundtrip_test.go` | 30min | Test DRY |
| 20 | Convert `tui/view_test.go` Contains pairs to table-driven | 10min | Test DRY |

### P4 — Nice-to-Have

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 21 | Add type-safe event constants (replace string-based routing in nom) | 1h | Safety |
| 22 | Fix unused struct field writes in `tui/messages_test.go:30-31` | 5min | Cleanup |
| 23 | Add `nom/configuration.go` direct tests | 15min | Coverage |
| 24 | Add `tui/lifecycle.go` direct tests | 15min | Coverage |
| 25 | Review shape registration coupling with marshaler registration | 30min | Architecture |

---

## g) Top #1 Question I Cannot Figure Out Myself

**The `GetDependencyTree()` data race:** Should we:
- **(A)** Return a deep copy (expensive for large trees, but safe)?
- **(B)** Return the raw pointer and document "caller must not mutate" (current behavior, but the TUI *does* mutate it — e.g., `model.nomSubscriber.GetDependencyTree().AddActivity(...)`)?
- **(C)** Make `DependencyTree` internally synchronized (add its own mutex)?

The TUI module actively mutates the tree via the returned pointer (e.g., `event_sequence_test.go:100-104`). Option C seems right but is a significant API change. **What's the intended ownership model?**

---

## Session Metrics

- **Files changed:** 22
- **Net lines removed:** 137 (513 removed, 376 added)
- **Clone groups eliminated (t=20):** 7 (29→22)
- **Clone groups eliminated (t=15):** 11 (105→94)
- **All 15 modules:** Tests passing, zero regressions
- **Production code changes:** 4 files (reporter.go, timing_cache_persist.go, subscriber_handlers.go, view.go)
- **Test code changes:** 18 files
