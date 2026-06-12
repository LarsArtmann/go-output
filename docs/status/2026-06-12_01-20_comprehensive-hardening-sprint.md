# Comprehensive Hardening Sprint — 2026-06-12

**Session focus:** Execute all 25 prioritized items from the 2026-06-11 audit, plus deep self-review.

---

## a) FULLY DONE

### Bug Fixes (P0) — 7 items

| #   | Change                                                                      | File(s)                       | Impact                       |
| --- | --------------------------------------------------------------------------- | ----------------------------- | ---------------------------- |
| 2   | Add RLock to `GetTimingCache()`                                             | `nom/state_accessors.go:27`   | **Critical** — was lock-free |
| 3   | Remove unused `content` param from `renderHelpOverlay()`                    | `tui/view.go:338`             | Cleanup                      |
| 4   | Replace 4× `//nolint:errcheck` with `mustUpdateActivityStatus` helper       | `integration/nom_tui_test.go` | Correctness                  |
| 5   | Replace `scrollOffset = 9999` with named `scrollToBottomSentinel = 1 << 30` | `tui/model.go:97-101`         | Correctness                  |
| 9   | Modernize `min()` builtin                                                   | `tui/view.go:62`              | Style                        |
| 10  | Use `slices.Contains` / `slices.ContainsFunc`                               | `nom/tree.go:71,79`           | Style                        |
| 22  | Remove unused `Current`/`Total` writes                                      | `tui/messages_test.go:30-31`  | Cleanup                      |

### Architecture (P1) — 5 items

| #   | Change                                                             | File(s)                                                | Impact                                              |
| --- | ------------------------------------------------------------------ | ------------------------------------------------------ | --------------------------------------------------- |
| 1   | Document `GetDependencyTree()` shared-state ownership contract     | `nom/state_accessors.go:19`                            | **Critical** — data race documentation              |
| 7   | Create generic `formatRegistry[T]` consolidating 3 registries      | `registry.go` (NEW), `render_tabledata.go`, `shape.go` | **DRY** — eliminated 3× Register/Get/Mutex patterns |
| 8   | Add `TableData.AddRowChecked()` returning error on column mismatch | `tabledata.go:43-58`                                   | Robustness                                          |
| 6   | Extend `SyncActivityTimingToTree` to also sync Status/Symbol/Color | `nom/activity_management.go:56`                        | Architecture — closes dual-state gap                |
| 21  | Add type-safe event constants (`EventWorkflowStarted`, etc.)       | `nom/event.go:8-15`                                    | Safety — replaces 6 string literals                 |

### Test Coverage (P2) — 5 items

| #   | Module                                                     | Before | After |   Delta    |
| --- | ---------------------------------------------------------- | :----: | :---: | :--------: |
| 11  | plantuml                                                   | 81.4%  | 93.0% | **+11.6%** |
| 12  | serialization                                              | 82.3%  | 89.7% | **+7.4%**  |
| 13  | graph                                                      | 86.8%  | 96.5% | **+9.7%**  |
| 14  | root (RenderAnyData, RegisteredFormats)                    | 92.3%  | 96.5% | **+4.2%**  |
| 15  | tui (Subscriber, SetCancelFunc, SetDisplayMode, lifecycle) | 87.6%  | 89.0% |   +1.4%    |

### Test Dedup (P3) — 5 items

| #   | Change                                                                                | File(s)                         | Clones Eliminated |
| --- | ------------------------------------------------------------------------------------- | ------------------------------- | :---------------: |
| 16  | Add `AssertErrorContains`, `AssertLineCount`, `AssertLastLineContains` to testhelpers | `testhelpers/helpers.go`        |   Cross-module    |
| 17  | Convert `nom/types_test.go` to parameterized tests                                    | `nom/types_test.go`             |         2         |
| 18  | Apply `assertLineCount` / `assertLastLineContains` in delimited                       | `delimited/csv_test.go`         |         2         |
| 19  | Extract `renderTableData`, `parseDelimited`, `assertCell` in integration round-trip   | `integration/roundtrip_test.go` |         4         |
| 20  | Convert Contains pairs to loop                                                        | `tui/view_test.go`              |         1         |

### Additional Work

| Change                                                   | File(s)                                | Impact              |
| -------------------------------------------------------- | -------------------------------------- | ------------------- |
| Add `nom/configuration.go` direct tests (#23)            | `nom/configuration_test.go` (NEW)      | Coverage            |
| Add `tui/lifecycle.go` direct tests (#24)                | `tui/reporter_lifecycle_test.go` (NEW) | Coverage            |
| Refactor shape registry to use `formatRegistry[T]` (#25) | `shape.go`                             | DRY                 |
| Extract `errWriter` for plantuml tests                   | `plantuml/plantuml_test.go`            | Test DRY            |
| Extract `errWriter` + registry tests for graph           | `graph/registry_test.go` (NEW)         | Coverage + Test DRY |

### Session Metrics

- **Files changed:** 29
- **Net lines:** +956/-265
- **Clone groups (t=20):** 22 → 21
- **Clone groups (t=15):** 94 → 89
- **All 15 modules:** Tests passing, `go vet` clean

---

## b) PARTIALLY DONE

### testhelpers Coverage Regression

The 3 new helpers (`AssertLineCount`, `AssertLastLineContains`, `AssertErrorContains`) have **0% coverage within testhelpers itself** because they're only called from other modules. This dragged testhelpers coverage from 91.3% → 73.3%.

**Fix:** Add direct unit tests in `testhelpers/helpers_test.go`.

### nom RenderNode / VisibleNodes at 0%

`nom/tree_render.go` `RenderNode()` and `VisibleNodes()` are exported but uncovered. These are used by the TUI rendering pipeline which runs in integration tests. Direct unit tests are needed.

---

## c) NOT STARTED

See section f) below.

---

## d) TOTALLY FUCKED UP

### Silent Error Swallowing

| File                          | Line   | Code                                                  | Risk                                       |
| ----------------------------- | ------ | ----------------------------------------------------- | ------------------------------------------ |
| `nom/timing_cache_persist.go` | 147    | `_ = writeCacheToFile(...)` in async save             | **Medium** — disk errors silently lost     |
| `nom/tree_render.go`          | 187    | `_ = dt.Build()` in `Render()`                        | **Medium** — empty output on build failure |
| `nom/tree_accessors.go`       | 31, 48 | `_ = dt.Build()` in `GetRootNodes()`, `EnsureBuild()` | **Medium** — stale/empty data              |

### Coverage Gaps

| Module        | Coverage | Target |         Gap          |
| ------------- | :------: | :----: | :------------------: |
| testhelpers   |  73.3%   |  90%+  | -16.7% (regression!) |
| serialization |  89.7%   |  90%+  |        -0.3%         |
| table         |  88.4%   |  90%+  |        -1.6%         |
| tui           |  89.0%   |  90%+  |        -1.0%         |

### Exported Dead Code

These are exported but have zero external callers:

| Package | Symbol                 | Assessment                   |
| ------- | ---------------------- | ---------------------------- |
| `nom`   | `ParseActivityID`      | **Dead** — no callers at all |
| `nom`   | `ParseWorkflowID`      | **Dead** — no callers at all |
| `nom`   | `GetOperationSymbol`   | Test-only                    |
| `nom`   | `SetOperationType`     | Test-only                    |
| `nom`   | `AddDependency`        | Test-only                    |
| `nom`   | `GetHistory`           | Test-only                    |
| `nom`   | `Remove` (TimingCache) | Test-only                    |
| `nom`   | `WaitPendingSaves`     | Test-only                    |
| `nom`   | `RenderNode`           | 0% coverage                  |
| `nom`   | `VisibleNodes`         | 0% coverage                  |
| `graph` | `NewGraphNodeID`       | Test-only                    |
| `graph` | `NewGraphNodeLabel`    | Test-only                    |

---

## e) WHAT WE SHOULD IMPROVE

### 1. Test Helpers Need Tests

The testhelpers module dropped from 91.3% → 73.3% because I added 3 new helpers without adding tests in the module itself. This is a regression I introduced.

### 2. Silent Error Swallowing in nom

Three places silently discard `dt.Build()` errors. If tree construction fails, the user gets no indication — the tree simply appears empty. This should return errors or at minimum log a warning.

### 3. Table Module Still at 88.4%

The table module (lipgloss-based tables) has been below 90% for the entire project lifetime. Needs focused test additions.

### 4. Magic Numbers in TUI View

`tui/view.go` has ~20 numeric literals for lipgloss color indices, progress bar dimensions, help overlay sizing. These should be named constants for clarity.

### 5. Consider `errgroup` for Async Saves

`nom/timing_cache_persist.go` uses raw goroutines for async saves. Using `errgroup` would provide structured error propagation.

### 6. `MarshalTSV` Uses `any` Where Union Type Would Be Better

`delimited/tsv.go:67` `MarshalTSV(data any)` does a type switch on `[][]string`/`[]string`/default. Could use a `TSVData` interface or explicit overloads.

---

## f) Top 25 Next Steps

Sorted by impact / effort (highest ROI first):

### P0 — Regressions I Caused (Must Fix)

| #   | Task                                                                                                 | Effort | Impact                      |
| --- | ---------------------------------------------------------------------------------------------------- | ------ | --------------------------- |
| 1   | Add unit tests for `AssertLineCount`, `AssertLastLineContains`, `AssertErrorContains` in testhelpers | 15min  | Fix 73.3% → 90%+ regression |
| 2   | Add tests for `nom.RenderNode()` and `nom.VisibleNodes()` (currently 0%)                             | 20min  | Fix coverage gap            |
| 3   | Add tests for `table` module (88.4% → 90%+)                                                          | 30min  | Close last sub-90% module   |

### P1 — Silent Error Handling (Safety)

| #   | Task                                                                                | Effort | Impact                  |
| --- | ----------------------------------------------------------------------------------- | ------ | ----------------------- |
| 4   | Make `dt.Build()` errors propagate in `Render()`, `GetRootNodes()`, `EnsureBuild()` | 30min  | Fix silent empty output |
| 5   | Replace `_ = writeCacheToFile()` with structured error logging                      | 15min  | Fix silent data loss    |
| 6   | Use `sync.WaitGroup` or `errgroup` for `saveAsync()` goroutine lifecycle            | 30min  | Structured concurrency  |

### P2 — Code Quality

| #   | Task                                                                                                         | Effort | Impact                |
| --- | ------------------------------------------------------------------------------------------------------------ | ------ | --------------------- |
| 7   | Extract lipgloss color constants in `tui/view.go`, `tui/summary.go`, `nom/symbols.go`                        | 30min  | Clarity               |
| 8   | Extract layout dimension constants (progress bar width, help overlay size) in `tui/view.go`                  | 15min  | Clarity               |
| 9   | Unexport or remove dead symbols: `ParseActivityID`, `ParseWorkflowID`                                        | 10min  | API surface reduction |
| 10  | Review `nom` test-only exports: should `SetOperationType`, `AddDependency`, etc. be unexported?              | 20min  | API surface reduction |
| 11  | Replace `any` in `MarshalTSV` with typed overloads or interface                                              | 20min  | Type safety           |
| 12  | Add `// Deprecated` or remove `graph.NewGraphNodeID`/`NewGraphNodeLabel` (use `output.NewBrandedID` instead) | 10min  | API cleanup           |

### P3 — Test Coverage

| #   | Task                                                                         | Effort | Impact       |
| --- | ---------------------------------------------------------------------------- | ------ | ------------ |
| 13  | Add `serialization` tests for writer error paths in JSON/YAML/TOML renderers | 20min  | 89.7% → 92%+ |
| 14  | Add `tui` tests for help overlay rendering, universal progress rendering     | 20min  | 89.0% → 92%+ |
| 15  | Add `nom` tests for `GetActivityCounts` Paused branch (currently 70%)        | 10min  | Coverage     |
| 16  | Add `nom` tests for `TimingCache` async save path                            | 15min  | Coverage     |
| 17  | Add `delimited` tests for streaming writer error paths                       | 15min  | Coverage     |
| 18  | Add `markup` tests for streaming HTML error paths                            | 15min  | Coverage     |

### P4 — Architecture (Higher Effort)

| #   | Task                                                                                    | Effort | Impact                             |
| --- | --------------------------------------------------------------------------------------- | ------ | ---------------------------------- |
| 19  | Unify `ActivityDisplayState` + `TreeNode` into single source of truth with embedding    | 2h     | Eliminate dual-state sync entirely |
| 20  | Replace raw goroutines in nom with `errgroup.Group` for structured lifecycle            | 1h     | Safety                             |
| 21  | Add integration test for `RenderNode`/`VisibleNodes` through TUI pipeline               | 30min  | Integration coverage               |
| 22  | Extract color theme into configurable struct (replace 20+ scattered lipgloss calls)     | 1h     | Customizability                    |
| 23  | Add fuzz targets for escape package (D2, DOT, Mermaid, AsciiDoc)                        | 1h     | Security hardening                 |
| 24  | Benchmark `formatRegistry.get()` vs hand-written map access                             | 15min  | Performance verification           |
| 25  | Add `go:generate` for enum boilerplate in `d2/d2_enum.go` (4 near-identical enum types) | 1h     | DRY                                |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should `nom.RenderNode()` and `nom.VisibleNodes()` remain exported?**

These are at 0% coverage and have zero external callers. They appear to be convenience wrappers around the tree's internal rendering logic. However:

- Removing them would be a breaking API change
- They might be used by downstream consumers we can't see
- They could be replaced by the TUI calling `tree.Render()` directly

**Options:**

- **(A)** Keep exported, add tests → safe but adds test maintenance
- **(B)** Unexport (`renderNode`, `visibleNodes`) → cleaner API, but breaking
- **(C)** Mark `// Deprecated: Use tree.Render() instead` → soft transition

I recommend **(C)** — mark deprecated, add minimal tests, plan removal in next major version.

---

## Coverage Summary

| Module        | Coverage |   Trend    |
| ------------- | :------: | :--------: |
| enum          |  100.0%  |     —      |
| escape        |  100.0%  |     —      |
| d2            |  97.0%   |     —      |
| output (root) |  96.5%   | ↑ (+4.2%)  |
| graph         |  96.5%   | ↑ (+9.7%)  |
| integration   |  95.5%   |     —      |
| markup        |  93.8%   |     —      |
| plantuml      |  93.0%   | ↑ (+11.6%) |
| nom           |  91.9%   |     —      |
| delimited     |  90.5%   |     —      |
| serialization |  89.7%   | ↑ (+7.4%)  |
| tui           |  89.0%   | ↑ (+1.4%)  |
| table         |  88.4%   |     —      |
| testhelpers   |  73.3%   | ↓ (-18.0%) |

## Clone Group Summary

| Threshold | Before Session | After Session | Delta |
| --------- | :------------: | :-----------: | :---: |
| t=20      |       22       |      21       |  -1   |
| t=15      |       94       |      89       |  -5   |
