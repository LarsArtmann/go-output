# Root Cleanup Plan — Dead Code Elimination & Architecture Polish

**Date:** 2026-05-25
**Branch:** `cleanup/root-dead-code`
**Scope:** Remove dead code, unused utilities, and deprecated APIs from root module
**Constraint:** Each task ≤ 12 min

---

## Executive Summary

Root has 16 production files. **4 are completely dead** (zero external consumers) and 1 exported function has zero callers. Removing them:

- Deletes **743 lines** of production + test code
- Removes `golang.org/x/term` external dependency from `go.mod`
- Drops root from 16 → 12 production files (~1,324 LOC)
- Leaves a clean core module: types, interfaces, data models, two lightweight formatters, dispatch

---

## Dependency Analysis — What's Dead

| File/Symbol                                        | External Prod Callers | External Test Callers                    | Verdict                      |
| -------------------------------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| `registry.go` (Register, Unregister, Create, etc.) | **ZERO**              | `integration/format_test.go` (1 test)    | DELETE                       |
| `sort.go` (SortBy enum)                            | **ZERO**              | **ZERO**                                 | DELETE                       |
| `color.go` (ColorMode, ShouldColor)                | **ZERO**              | **ZERO**                                 | DELETE                       |
| `slices.go` (FilledStrings)                        | **ZERO**              | `integration/workflow_test.go` (2 calls) | DELETE — inline at call site |
| `BrandedValue()` in `marshal.go`                   | **ZERO**              | **ZERO**                                 | DELETE from marshal.go       |

---

## What Survives (12 production files)

| File                  | Lines | Concern                                              |
| --------------------- | ----- | ---------------------------------------------------- |
| `format.go`           | 99    | Format enum                                          |
| `shape.go`            | 107   | Shape enum + capability matrix                       |
| `renderer.go`         | 30    | Renderer/TableRenderer interfaces                    |
| `ids.go`              | 51    | Branded ID phantom types                             |
| `tabledata.go`        | 127   | TableData, TableDataBase, RowEdge                    |
| `graph.go`            | 205   | GraphNode, GraphEdge, GraphRendererMixin, GraphShape |
| `graph_tabledata.go`  | 76    | TableData→Graph conversion                           |
| `tree.go`             | 182   | TreeNode + ASCIITreeRenderer                         |
| `markdown.go`         | 205   | MarkdownTable renderer                               |
| `marshal.go`          | ~44   | MarshalFormat, UnmarshalFormat, MarshalJSONIndent    |
| `streaming.go`        | 54    | StreamingRenderer interface + adapter                |
| `render_tabledata.go` | 149   | RenderTableData dispatcher + marshaler registry      |

---

## Overall Progress

| Phase     | Description                     | Status    | Tasks Done |
| --------- | ------------------------------- | --------- | ---------- |
| Phase 1   | Delete dead production files    | ⬜ TODO   | 0/4        |
| Phase 2   | Delete dead test files          | ⬜ TODO   | 0/5        |
| Phase 3   | Delete dead code from survivors | ⬜ TODO   | 0/3        |
| Phase 4   | Fix external references         | ⬜ TODO   | 0/3        |
| Phase 5   | Clean go.mod + golangci.yml     | ⬜ TODO   | 0/4        |
| Phase 6   | Update documentation            | ⬜ TODO   | 0/4        |
| Phase 7   | Verify + test all modules       | ⬜ TODO   | 0/3        |
| **Total** |                                 | **⬜ 0%** | **0/26**   |

---

## Execution Plan

Tasks sorted by: Impact → Customer Value → Effort (ascending). Each ≤ 12 min.

### Phase 1: Delete Dead Production Files (1% → 51%)

Highest impact — removes dead code immediately visible to users.

| #  | Task                                                                                                                               | File(s)       | Impact                                      | Effort |
| -- | ---------------------------------------------------------------------------------------------------------------------------------- | ------------- | ------------------------------------------- | ------ |
| 01 | Delete `registry.go` — deprecated renderer registry (106 LOC). All 5 functions deprecated, zero external prod callers.             | `registry.go` | 🔴 Critical — removes deprecated public API | 1min   |
| 02 | Delete `sort.go` — deprecated SortBy enum (67 LOC). Zero external callers. Was supposed to be removed with sort/ module.           | `sort.go`     | 🔴 Critical — removes deprecated public API | 1min   |
| 03 | Delete `color.go` — ColorMode enum + terminal detection (105 LOC). Zero external callers. Only file importing `golang.org/x/term`. | `color.go`    | 🟠 High — removes unused external dep       | 1min   |
| 04 | Delete `slices.go` — FilledStrings utility (9 LOC). Zero external prod callers. Trivial `slices.Repeat` wrapper.                   | `slices.go`   | 🟢 Low — removes trivial dead code          | 1min   |

**Commit:** `refactor: delete dead code — registry, sort, color, slices`

---

### Phase 2: Delete Dead Test Files (4% → 64%)

| #  | Task                                                                                                                | File(s)                      | Impact                           | Effort |
| -- | ------------------------------------------------------------------------------------------------------------------- | ---------------------------- | -------------------------------- | ------ |
| 05 | Delete `registry_test.go` — tests for deleted registry (181 LOC).                                                   | `registry_test.go`           | 🟠 High — removes orphaned tests | 1min   |
| 06 | Delete `sort_test.go` — tests for deleted SortBy enum (82 LOC). Includes FuzzParseSortBy.                           | `sort_test.go`               | 🟠 High — removes orphaned tests | 1min   |
| 07 | Delete `color_test.go` — tests for deleted ColorMode (152 LOC).                                                     | `color_test.go`              | 🟠 High — removes orphaned tests | 1min   |
| 08 | Delete `slices_test.go` — tests for deleted FilledStrings (41 LOC).                                                 | `slices_test.go`             | 🟢 Low — removes orphaned tests  | 1min   |
| 09 | Delete `TestFormatRegistry` from `integration/format_test.go` (lines 112-136, ~25 LOC). Tests deleted registry API. | `integration/format_test.go` | 🟠 High — removes broken test    | 3min   |

**Commit:** `test: delete orphaned tests for removed code`

---

### Phase 3: Delete Dead Code From Surviving Files (20% → 80%)

| #  | Task                                                                                                                                                                                       | File(s)           | Impact                                | Effort |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------- | ------------------------------------- | ------ |
| 10 | Delete `BrandedValue()` from `marshal.go` (~7 LOC). Zero production callers. Only referenced in docs.                                                                                      | `marshal.go`      | 🟡 Medium — removes dead exported API | 2min   |
| 11 | Remove `errorRenderer` type from `format_test.go` (lines 163-167). Only used for MustRender panic test — check if still needed.                                                            | `format_test.go`  | 🟢 Low — cleanup                      | 2min   |
| 12 | Verify `testing_test.go` helpers still used: `testParseEnum`, `testEnumString`, `testAllowedValues`, `runSubtest`. Check if `sort_test.go` and `color_test.go` deletion leaves any unused. | `testing_test.go` | 🟡 Medium — prevents dead helpers     | 5min   |

**Commit:** `refactor: remove BrandedValue dead code, clean up test helpers`

---

### Phase 4: Fix External References (80% → 90%)

| #  | Task                                                                                                                                                                                                                             | File(s)                        | Impact                            | Effort |
| -- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------ | --------------------------------- | ------ |
| 13 | Replace `output.FilledStrings(10, "Col")` in `integration/workflow_test.go` lines 152, 156 with inline `slices.Repeat([]string{"Col"}, 10)`. Add `"slices"` import.                                                              | `integration/workflow_test.go` | 🔴 Critical — fixes broken import | 3min   |
| 14 | Search entire codebase for any remaining references to deleted symbols: `Register(`, `Unregister(`, `Create(format`, `IsRegistered`, `RegisteredFormats`, `SortBy`, `ColorMode`, `FilledStrings`, `BrandedValue`. Fix any found. | all files                      | 🔴 Critical — catches stale refs  | 8min   |
| 15 | Run `go build ./...` in root to verify compilation. Fix any compile errors.                                                                                                                                                      | root module                    | 🔴 Critical — verifies build      | 3min   |

**Commit:** `fix: update external references for deleted root code`

---

### Phase 5: Clean go.mod + golangci.yml (90% → 95%)

| #  | Task                                                                                                                                 | File(s)         | Impact                           | Effort |
| -- | ------------------------------------------------------------------------------------------------------------------------------------ | --------------- | -------------------------------- | ------ |
| 16 | Remove `golang.org/x/term` from `go.mod` require block. Was only imported by deleted `color.go`.                                     | `go.mod`        | 🟠 High — removes unused dep     | 2min   |
| 17 | Remove `golang.org/x/term` from depguard allow-lists in `.golangci.yml` (likely lines ~163 and ~205).                                | `.golangci.yml` | 🟡 Medium — lint config accuracy | 3min   |
| 18 | Run `go mod tidy` in root. Verify `golang.org/x/sys` indirect also removed (was only needed by x/term).                              | root module     | 🟠 High — dependency hygiene     | 2min   |
| 19 | Run `go mod tidy` in all sub-modules that might reference deleted symbols (delimited, markup, serialization, integration, examples). | all modules     | 🟡 Medium — hygiene              | 5min   |

**Commit:** `chore: remove golang.org/x/term dep, tidy all go.mod files`

---

### Phase 6: Update Documentation (95% → 99%)

| #  | Task                                                                                                                                                                                                                     | File(s)                                   | Impact                            | Effort |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------- | --------------------------------- | ------ |
| 20 | Update `AGENTS.md`: remove references to `sort.go`, `color.go`, `registry.go`, `slices.go`, `BrandedValue`, `FilledStrings`, `SortBy`, `ColorMode`. Update module table, file listing, coverage table, dependency graph. | `AGENTS.md`                               | 🟠 High — developer docs accuracy | 10min  |
| 21 | Update `CHANGELOG.md`: add `[Unreleased]` entry listing removed symbols with migration notes.                                                                                                                            | `CHANGELOG.md`                            | 🟠 High — user-facing changelog   | 5min   |
| 22 | Update `README.md`: remove any references to deleted symbols. Verify code examples don't reference `SortBy`, `ColorMode`, `Register`, `BrandedValue`, `FilledStrings`.                                                   | `README.md`                               | 🟡 Medium — user docs             | 5min   |
| 23 | Update `docs/modularization/DEPENDENCY_GRAPH.md`: update root LOC, file count, dependency list. Remove `x/term` reference.                                                                                               | `docs/modularization/DEPENDENCY_GRAPH.md` | 🟢 Low — planning docs accuracy   | 5min   |

**Commit:** `docs: update AGENTS.md, CHANGELOG, README for dead code removal`

---

### Phase 7: Verify + Test All Modules (99% → 100%)

| #  | Task                                                                                                                                     | File(s)     | Impact                       | Effort |
| -- | ---------------------------------------------------------------------------------------------------------------------------------------- | ----------- | ---------------------------- | ------ |
| 24 | Run full root test suite: `go test ./...` in root. Verify all pass.                                                                      | root        | 🔴 Critical — final gate     | 3min   |
| 25 | Run all sub-module tests: `go test ./...` in delimited, markup, serialization, d2, graph, table, integration, examples. Verify all pass. | all modules | 🔴 Critical — no regressions | 8min   |
| 26 | Run lint: `golangci-lint run ./...` in root. Fix any issues (unused imports, missing nolint directives).                                 | root        | 🟠 High — lint hygiene       | 5min   |

**Commit:** `chore: verify all modules build, test, and lint clean`

---

## Summary Table

| #  | Phase | Task                                                       | Impact      | Effort  | Est      |
| -- | ----- | ---------------------------------------------------------- | ----------- | ------- | -------- |
| 01 | 1     | Delete `registry.go`                                       | 🔴 Critical | Trivial | 1m       |
| 02 | 1     | Delete `sort.go`                                           | 🔴 Critical | Trivial | 1m       |
| 03 | 1     | Delete `color.go`                                          | 🟠 High     | Trivial | 1m       |
| 04 | 1     | Delete `slices.go`                                         | 🟢 Low      | Trivial | 1m       |
| 05 | 2     | Delete `registry_test.go`                                  | 🟠 High     | Trivial | 1m       |
| 06 | 2     | Delete `sort_test.go`                                      | 🟠 High     | Trivial | 1m       |
| 07 | 2     | Delete `color_test.go`                                     | 🟠 High     | Trivial | 1m       |
| 08 | 2     | Delete `slices_test.go`                                    | 🟢 Low      | Trivial | 1m       |
| 09 | 2     | Delete `TestFormatRegistry` from integration               | 🟠 High     | Trivial | 3m       |
| 10 | 3     | Delete `BrandedValue()` from `marshal.go`                  | 🟡 Medium   | Trivial | 2m       |
| 11 | 3     | Clean `errorRenderer` from `format_test.go`                | 🟢 Low      | Trivial | 2m       |
| 12 | 3     | Verify `testing_test.go` helpers still used                | 🟡 Medium   | Small   | 5m       |
| 13 | 4     | Fix `FilledStrings` refs in `integration/workflow_test.go` | 🔴 Critical | Trivial | 3m       |
| 14 | 4     | Search & fix all stale references                          | 🔴 Critical | Small   | 8m       |
| 15 | 4     | Verify root `go build ./...`                               | 🔴 Critical | Trivial | 3m       |
| 16 | 5     | Remove `golang.org/x/term` from `go.mod`                   | 🟠 High     | Trivial | 2m       |
| 17 | 5     | Remove `golang.org/x/term` from depguard                   | 🟡 Medium   | Trivial | 3m       |
| 18 | 5     | `go mod tidy` root — verify `x/sys` removed                | 🟠 High     | Trivial | 2m       |
| 19 | 5     | `go mod tidy` all sub-modules                              | 🟡 Medium   | Small   | 5m       |
| 20 | 6     | Update `AGENTS.md`                                         | 🟠 High     | Small   | 10m      |
| 21 | 6     | Update `CHANGELOG.md`                                      | 🟠 High     | Trivial | 5m       |
| 22 | 6     | Update `README.md`                                         | 🟡 Medium   | Trivial | 5m       |
| 23 | 6     | Update `DEPENDENCY_GRAPH.md`                               | 🟢 Low      | Trivial | 5m       |
| 24 | 7     | Run full root test suite                                   | 🔴 Critical | Trivial | 3m       |
| 25 | 7     | Run all sub-module tests                                   | 🔴 Critical | Small   | 8m       |
| 26 | 7     | Run lint on root                                           | 🟠 High     | Small   | 5m       |
|    |       | **Total**                                                  |             |         | **~85m** |

---

## Risks

| Risk                                                     | Mitigation                                                                          |
| -------------------------------------------------------- | ----------------------------------------------------------------------------------- |
| External consumers importing deleted symbols             | These are all deprecated/unreleased symbols — check `go mod graph` for any consumer |
| `golang.org/x/sys` might be needed by other deps         | `go mod tidy` will keep it if needed; verify after                                  |
| `testing_test.go` helpers might be unused after deletion | Task 12 explicitly checks this                                                      |
| Depguard false positives after removing x/term           | Task 17 explicitly cleans depguard config                                           |

---

## Commit Strategy

| Phase | Commits | Message Pattern                                                   |
| ----- | ------- | ----------------------------------------------------------------- |
| 1     | 1       | `refactor: delete dead code — registry, sort, color, slices`      |
| 2     | 1       | `test: delete orphaned tests for removed code`                    |
| 3     | 1       | `refactor: remove BrandedValue dead code, clean up test helpers`  |
| 4     | 1       | `fix: update external references for deleted root code`           |
| 5     | 1       | `chore: remove golang.org/x/term dep, tidy all go.mod files`      |
| 6     | 1       | `docs: update AGENTS.md, CHANGELOG, README for dead code removal` |
| 7     | 0       | Verification only — fix issues in previous commits if needed      |

---

## Future Considerations (NOT in this plan)

These are architectural improvements for a later iteration:

1. **Move `markdown.go` to registry dispatch** — makes `RenderTableData` a pure dispatcher
2. **Extract `graph.go` + `graph_tabledata.go` into their own concern** — but they're shared by d2/, graph/, serialization/
3. **Consider if `RenderTableData` is worth keeping** — zero production callers outside tests
4. **Move `Alignment` type** — currently in `markdown.go`, could be shared if more table renderers need it
