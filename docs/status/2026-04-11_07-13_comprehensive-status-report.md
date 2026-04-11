# go-output — Comprehensive Status Report

**Date:** 2026-04-11 07:13 CEST
**Branch:** master
**Last Commit:** `1fef47d` — refactor(tests): extract test assertion helpers to reduce code duplication
**Previous Status:** 2026-04-05 (6 days ago)

---

## Executive Summary

go-output is a **mature Go library** for CLI output formatting supporting 12 formats across 3 data models (table, tree, graph). The library itself is **solid** — 91.5% coverage, clean API, branded IDs, streaming support. However, the surrounding infrastructure has **critical broken states**: integration tests won't compile, CI uses wrong Go version, and release workflow references missing goreleaser config. There are also **9 uncommitted files** with half-finished refactoring work in the working tree.

**Overall Health: 🟡 Functional but needs infrastructure repair before next release.**

---

## A) FULLY DONE ✅

### Core Library (Production Quality)

| Component                | Status      | Details                                                                    |
| ------------------------ | ----------- | -------------------------------------------------------------------------- |
| **12 Output Formats**    | ✅ Complete | Table, JSON, CSV, TSV, Markdown, XML, YAML, D2, HTML, Tree, Mermaid, DOT   |
| **3 Data Models**        | ✅ Complete | `TableData`, `TreeNode`, `Graph` with cross-format bridges                 |
| **Type-Safe Enums**      | ✅ Complete | `Format`, `SortBy`, `ColorMode` with `Parse`/`Validate`/`IsValid`          |
| **Branded IDs**          | ✅ Complete | `TreeNodeID`, `GraphNodeID`, `GraphNodeLabel` — prevents type confusion    |
| **Format Registry**      | ✅ Complete | Thread-safe factory pattern with `Register`/`Create`/`IsRegistered`        |
| **Streaming**            | ✅ Complete | `StreamingRenderer` interface + `StreamingHTMLRenderer` (true incremental) |
| **Error Handling**       | ✅ Complete | Typed errors (`InvalidFormatError`, sentinel errors in `pkg/errors`)       |
| **cmdguard Integration** | ✅ Complete | `EnumFlag[T]` generic flag parser for CLI frameworks                       |
| **Sorting**              | ✅ Complete | Generic `Sorter[T]` with multiple sort fields and custom `LessFunc`        |
| **Internal Helpers**     | ✅ Complete | `escape`, `gentest`, `testutils` packages                                  |

### Test Coverage

| Package            | Coverage  | Tests         |
| ------------------ | --------- | ------------- |
| `go-output` (root) | **91.5%** | Comprehensive |
| `cmdguard`         | **100%**  | Full          |
| `sort`             | **94.6%** | Strong        |
| `table`            | **100%**  | Full          |
| `enum`             | **71.4%** | Adequate      |
| `integration`      | 🔴 BROKEN | Won't compile |

### Test Infrastructure

| Type              | Count   | Status                                                                       |
| ----------------- | ------- | ---------------------------------------------------------------------------- |
| Unit Tests        | ~185+   | ✅ Passing                                                                   |
| Benchmarks        | 13      | ✅ Present                                                                   |
| Fuzz Targets      | 3       | ✅ Present (`FuzzParseOutputFormat`, `FuzzParseSortBy`, `FuzzMarkdownTable`) |
| Integration Tests | 4 files | 🔴 BROKEN                                                                    |

### Code Quality

- **Zero TODO/FIXME comments** in codebase
- No code duplication (jscpd reports clean)
- Consistent error wrapping patterns
- File size limit (350 lines) mostly respected — `sort/sort_test.go` at 400 lines is the only violation

---

## B) PARTIALLY DONE 🟡

### Uncommitted Refactoring Work (9 files modified)

These changes are **in the working tree but not committed** — appears to be an in-progress refactoring session:

| File                  | Change                                                                   | Risk                    |
| --------------------- | ------------------------------------------------------------------------ | ----------------------- |
| `cmdguard/flag.go`    | `NewEnumFlag` signature changed to use `enumFlagConfig[T]` struct        | **Breaking API change** |
| `cmdguard/color.go`   | Updated to new `NewEnumFlag` signature                                   | Depends on flag.go      |
| `cmdguard/format.go`  | Updated to new `NewEnumFlag` signature                                   | Depends on flag.go      |
| `cmdguard/sort.go`    | Updated to new `NewEnumFlag` signature                                   | Depends on flag.go      |
| `json.go`             | Inlined `marshalJSONIndent` → direct `marshalIndent` call                | Cleanup                 |
| `marshal.go`          | Removed `marshalJSONIndent` and `marshalXMLIndent` helpers               | Cleanup                 |
| `xml.go`              | Inlined `marshalXMLIndent` → direct `marshalIndent` call + named returns | Cleanup                 |
| `sort/sort_test.go`   | Extracted `countLessFunc` helper                                         | Test improvement        |
| `userjourney_test.go` | Extracted `makeProjects` helper, used named test data                    | Test improvement        |

**Assessment:** The cmdguard refactoring is a **breaking API change** — `NewEnumFlag` now requires a config struct instead of positional parameters. The marshal cleanup is pure improvement.

### Streaming Support

- ✅ `StreamingHTMLRenderer` — true incremental streaming
- 🟡 Other formats use `StreamingRendererFromRenderer` adapter (not true streaming, just `Render()` + write)
- ⬜ No streaming for JSON, CSV, YAML, etc.

### CI/CD Pipeline

- ✅ Three workflow jobs (test, verify, lint)
- 🔴 Go version mismatch: CI uses `1.23`, `go.mod` requires `1.26.0`
- 🔴 No goreleaser config exists but release workflow runs `goreleaser`
- 🟡 No caching of Go modules
- ⬜ No matrix testing (multiple Go versions)

### Documentation

- ✅ README with usage examples
- ✅ CHANGELOG started
- ✅ AGENTS.md for AI agents
- 🟡 Example uses deprecated `ParseOutputFormat` API
- ⬜ No CONTRIBUTING.md
- ⬜ No API reference documentation

---

## C) NOT STARTED ⬜

| Item                                | Priority     | Notes                                                      |
| ----------------------------------- | ------------ | ---------------------------------------------------------- |
| `.goreleaser.yml`                   | High         | Release workflow calls goreleaser but config doesn't exist |
| CONTRIBUTING.md                     | Medium       | Standard for open-source projects                          |
| CI Go version fix                   | **Critical** | CI is on 1.23, project is 1.26 — will fail                 |
| Streaming writers for JSON/CSV/YAML | Low          | Adapter exists but not true streaming                      |
| API reference (godoc/pkgsite)       | Medium       | No published docs                                          |
| Multi-version CI matrix             | Low          | Only tested on single Go version                           |
| Go module caching in CI             | Low          | Performance optimization                                   |
| SARIF output format                 | N/A          | Investigated — **not appropriate** for this library        |

---

## D) TOTALLY FUCKED UP 🔴

### 1. Integration Tests Won't Compile (BLOCKER)

**File:** `integration/workflow_test.go:28`
**Error:** `undefined: testutils.AssertTableData`

`testutils.AssertTableData(t, data, 2, 3)` is called but the function **does not exist** in `internal/testutils/test_helpers.go`. This appears to be a leftover from the test helper extraction refactoring (commit `1fef47d` and friends). The function was likely moved or renamed but this call site was missed.

**Impact:** `go test ./...` **FAILS**. Integration tests are completely dead.

### 2. CI Will Fail on Every Push

**Both `ci.yml` and `release.yml` use Go `1.23`.** The project requires Go `1.26.0`. This means:

- `go build` will fail (unknown Go version directive)
- Or if 1.23 is treated as compatible, features like `range` over integers will cause compile errors

**Impact:** No CI protection. Broken code can merge without detection.

### 3. Release Workflow is Broken

The `goreleaser` job references `goreleaser/goreleaser-action@v5` but **no `.goreleaser.yml` exists**. Additionally, the `release` job creates a tag and pushes it — then the `goreleaser` job also runs. This will fail on every tag push.

**Impact:** Cannot ship releases through CI.

### 4. golangci-lint Config Issues

Per previous status report (2026-04-05), the `.golangci.yml` has:

- `linters.settings` placed under wrong key — most settings silently ignored
- Blanket exclusion rules that disable most linters for all `*.go` files

**Impact:** Linter gives false confidence — many checks are effectively disabled.

### 5. Markdown AlignCenter Bug (Known, Unfixed)

`MarkdownTable` `AlignCenter` renders identically to `AlignLeft`. Identified in previous status reports, still not fixed.

---

## E) WHAT WE SHOULD IMPROVE

### Immediate (Before Next Commit)

1. **Fix `AssertTableData`** — Add the missing function to `testutils` or rewrite the test
2. **Commit or revert the 9 uncommitted files** — Don't leave half-finished work in the tree
3. **Fix the cmdguard breaking change** — Either commit it as a major change or revert; document the migration

### Short Term (This Week)

4. **Fix CI Go version** — Change `1.23` → `1.26` in both workflow files
5. **Remove or fix goreleaser** — Either add `.goreleaser.yml` or remove the goreleaser job
6. **Fix `.golangci.yml` config** — Move settings to correct key, review blanket exclusions
7. **Update example** — Replace deprecated `ParseOutputFormat` with `ParseFormat`

### Medium Term (Next Sprint)

8. **Add CONTRIBUTING.md** — Standardize contribution process
9. **Trim status docs** — 14 status documents is excessive; archive older ones
10. **Fix Markdown AlignCenter** — Known rendering bug
11. **Add `AssertTableData` to gentest** — Proper home for shared test assertions

### Architecture Level

12. **Reconsider file sizes** — `sort/sort_test.go` (400 lines) and `d2.go` (305 lines) approach limits
13. **Streaming strategy** — Decide if true streaming is needed beyond HTML
14. **Version strategy** — `v0.1.0` is ancient; decide on v1.0 roadmap or keep experimental
15. **Deprecation timeline** — `OutputFormat` aliases are deprecated but no removal date set

---

## F) Top #25 Things We Should Get Done Next

### Priority 1: STOP THE BLEEDING (Critical Fixes)

| #   | Task                                                             | Effort | Impact     |
| --- | ---------------------------------------------------------------- | ------ | ---------- |
| 1   | Fix `testutils.AssertTableData` — make integration tests compile | 15min  | 🔴 BLOCKER |
| 2   | Run `go test ./...` and verify ALL tests pass                    | 5min   | 🔴 BLOCKER |
| 3   | Fix CI Go version: `1.23` → `1.26` in `ci.yml` and `release.yml` | 5min   | 🔴 BLOCKER |
| 4   | Commit or revert the 9 uncommitted files with proper messages    | 10min  | 🔴 BLOCKER |
| 5   | Fix or remove goreleaser job from `release.yml`                  | 15min  | 🔴 Broken  |

### Priority 2: FIX KNOWN BUGS

| #   | Task                                                                         | Effort | Impact        |
| --- | ---------------------------------------------------------------------------- | ------ | ------------- |
| 6   | Fix Markdown `AlignCenter` rendering bug                                     | 30min  | 🟡 Bug        |
| 7   | Update `examples/basic/main.go` to use `ParseFormat` not `ParseOutputFormat` | 5min   | 🟡 Deprecated |
| 8   | Fix `.golangci.yml` settings placement                                       | 30min  | 🟡 Config     |

### Priority 3: CLEAN UP

| #   | Task                                                                                           | Effort | Impact           |
| --- | ---------------------------------------------------------------------------------------------- | ------ | ---------------- |
| 9   | Archive old status docs (keep only last 2-3)                                                   | 10min  | 🟢 Hygiene       |
| 10  | Update CHANGELOG.md with recent changes                                                        | 30min  | 🟢 Documentation |
| 11  | Add `AssertTableData` to `internal/gentest` properly                                           | 15min  | 🟢 Test infra    |
| 12  | Fix 4 files with unnecessary type arguments (gopls warnings)                                   | 10min  | 🟢 Cleanup       |
| 13  | Remove `marshalJSONIndent`/`marshalXMLIndent` from marshal.go (already in uncommitted changes) | Done   | ✅               |

### Priority 4: IMPROVE

| #   | Task                                               | Effort | Impact         |
| --- | -------------------------------------------------- | ------ | -------------- |
| 14  | Add `CONTRIBUTING.md`                              | 30min  | 🟢 Process     |
| 15  | Review and tighten `.golangci.yml` exclusion rules | 30min  | 🟢 Quality     |
| 16  | Add CI Go module caching                           | 10min  | 🟢 Performance |
| 17  | Split `sort/sort_test.go` (400 lines → under 350)  | 15min  | 🟢 Standards   |
| 18  | Add CI matrix testing (Go 1.26 + 1.27)             | 15min  | 🟢 Robustness  |

### Priority 5: STRATEGIC

| #   | Task                                                                   | Effort     | Impact       |
| --- | ---------------------------------------------------------------------- | ---------- | ------------ |
| 19  | Decide on v1.0 release roadmap                                         | Discussion | 🔵 Strategic |
| 20  | Set deprecation timeline for `OutputFormat` aliases                    | Discussion | 🔵 Strategic |
| 21  | Evaluate if true streaming needed for JSON/CSV                         | Discussion | 🔵 Strategic |
| 22  | Add `.goreleaser.yml` for proper release automation                    | 1hr        | 🔵 Release   |
| 23  | Publish godoc/pkgsite documentation                                    | 30min      | 🔵 Adoption  |
| 24  | Add `go test -race` to CI (already in `test` job but not `verify`)     | 5min       | 🔵 Quality   |
| 25  | Review cmdguard API stability — is `enumFlagConfig` the right pattern? | Discussion | 🔵 API       |

---

## G) Top #1 Question I Cannot Figure Out Myself

### What is the intended API for `cmdguard.NewEnumFlag`?

The working tree contains a **breaking change** to `NewEnumFlag`:

**Before (committed):**

```go
func NewEnumFlag[T EnumValue](
    value *T,
    name string,
    parseFunc func(string) (T, error),
) *EnumFlag[T]
```

**After (uncommitted):**

```go
type enumFlagConfig[T EnumValue] struct {
    value     *T
    name      string
    parseFunc func(string) (T, error)
}

func NewEnumFlag[T EnumValue](cfg enumFlagConfig[T]) *EnumFlag[T]
```

This is a **breaking API change** for any external caller of `NewEnumFlag`. The deprecated wrapper functions (`NewOutputFormatFlag`, `NewColorModeFlag`, `NewSortByFlag`) still work because they construct the config struct internally — but direct callers of `NewEnumFlag` will break.

**Questions only the author can answer:**

1. Is this the intended direction, or was this experimentation?
2. Should this be a v2 breaking change, or is the API considered unstable?
3. Is there a plan for functional options instead of a config struct?

---

## Project Metrics Snapshot

| Metric            | Value                                  |
| ----------------- | -------------------------------------- |
| Total Go Lines    | ~9,119                                 |
| Go Files          | ~50                                    |
| Output Formats    | 12                                     |
| Data Models       | 3 (Table, Tree, Graph)                 |
| Root Coverage     | 91.5%                                  |
| Test Files        | ~20                                    |
| Uncommitted Files | 9                                      |
| Known Bugs        | 2 (AlignCenter, AssertTableData)       |
| CI Status         | 🔴 Will fail (Go version mismatch)     |
| Last Release      | v0.1.0 (2026-01-01)                    |
| Dependencies      | 2 direct (lipgloss/v2, go-faster/yaml) |
| License           | Present                                |
| Code Owners       | AUTHORS file present                   |
