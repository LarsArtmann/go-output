# go-output — Comprehensive Status Report

**Date:** 2026-05-17 05:36
**Session:** Session 7 (testhelpers extraction, Tree/Graph renderers, lint fixes, CI overhaul)
**Branch:** master (uncommitted changes)
**Latest Tag:** v0.4.0 (20+ commits ahead of tag — should be v0.5.0)
**Coverage:** 90.9% (root), 100% (all sub-modules)
**Lint/Race:** 0 issues, 0 races across 9/9 modules

---

## A. FULLY DONE ✅

### Sessions 1-4 (Prior Work)

- [x] MIT license, README rewrite, 27 doc comments, GitHub release v0.4.0
- [x] Shape capability matrix (ADR 002) — `Shape` type, `formatCapabilities` map, `Supports()/Shapes()/FormatsForShape()`
- [x] Structural refactoring — TreeNode→tree.go, TableData→tabledata.go, GraphRendererMixin→graph.go
- [x] Stale PLAN.md deleted, deprecated FormatCategory/IsXxxFormat() methods

### Session 5 (Circular Dep Fix)

- [x] **sort/ circular dependency eliminated** — deleted `Sorter[T]`, modernized `ByField` to return `int`
- [x] **sort/ now zero deps** — stripped go.mod from 21 lines to 3 lines
- [x] **cmdguard/ fixed** — go.mod had proper deps, tests fixed (removed `internal/gentest`)
- [x] **enum/ fixed** — removed `internal/gentest` import, inlined helper
- [x] **table/ fixed** — added missing `go-branded-id` dep via `go mod tidy`
- [x] Root and integration go.mod cleaned of sort/ references
- [x] User journey and integration tests migrated to stdlib `slices.SortStableFunc`

### Session 6 (Documentation + Table Renderers)

- [x] **CHANGELOG.md** — Complete rewrite with version history v0.1.0→v0.4.0 plus [Unreleased]
- [x] **FORMAT_ARCHITECTURE.md** — Replaced old category system with Shape capability matrix
- [x] **AGENTS.md** — Full module table, dependency graph, project structure, build commands
- [x] **CONTRIBUTING.md** — Removed stale `just`/`justfile` references, added go.work setup
- [x] **go.work.example** — Created for contributors
- [x] **JSONTableRenderer** — `NewJSONTableRenderer()`, implements `Renderer` + `TableRenderer`
- [x] **YAMLTableRenderer** — Same pattern, renders TableData as YAML sequence of mappings
- [x] **TableData.ToMapSlice()** — New method converting tabular data to `[]map[string]string`

### Session 7 (This Session)

- [x] **format_test.go split** — 468 lines → 3 files (169+107+256 lines), all under 350 limit
- [x] **testhelpers/ module** — New zero-dep sub-module extracted from `internal/gentest`
  - `AssertStringSliceEqual`, `AssertEqual[T]`, `AssertContains`, `TestEnumIsValid[T]`
  - `TestStructFields`, `StringField`, `IntField`, `FieldCheck`
  - 100% test coverage
- [x] **DRY violation fixed** — `assertStringSliceEqual` no longer duplicated in `enum/` and `cmdguard/`
- [x] **JSONTreeRenderer** — `NewJSONTreeRenderer()`, implements `Renderer` + `TreeOutputRenderer`
- [x] **JSONGraphRenderer** — `NewJSONGraphRenderer()`, implements `Renderer` + `GraphRenderer`
- [x] **YAMLTreeRenderer** — `NewYAMLTreeRenderer()`, same pattern
- [x] **YAMLGraphRenderer** — `NewYAMLGraphRenderer()`, same pattern
- [x] **CI/CD overhaul** — `.github/workflows/ci.yml` now builds/tests/lints all 9 modules
- [x] **README.md updated** — Added JSON/YAML table/tree/graph renderer examples
- [x] **CONTRIBUTING.md updated** — Keep a Changelog link, module count 7→8, testhelpers in go.work
- [x] **All lint issues fixed** — 0 golangci-lint issues across root module
- [x] **depguard config updated** — Added `testhelpers` to allow lists

### Quality Metrics (Current)

| Metric                                  | Value        |
| --------------------------------------- | ------------ |
| Modules building                        | 9/9          |
| Tests passing                           | 9/9          |
| Lint issues                             | 0            |
| Race conditions                         | 0            |
| Root coverage                           | 90.9%        |
| Sub-module coverage                     | 100% (all 6) |
| Files over 350 lines                    | 0            |
| Circular dependencies                   | 0            |
| `internal/gentest` cross-module imports | 0            |
| Duplicated test helpers                 | 0 (was 2)    |

### Dependency Graph (Current)

```
root (output) → enum, escape, yaml, x/term, go-branded-id, testhelpers
enum          → testhelpers (tests only)
escape        → (none)
testhelpers   → (none) — zero deps, shared test assertions
sort          → (none) — zero deps, only ByField helper
cmdguard      → root (tests only), testhelpers (tests); prod code standalone
table         → root, lipgloss/v2
integration   → root, table
examples      → root, table
```

**No circular dependencies.**

---

## B. PARTIALLY DONE 🟡

### None — all started work items are complete.

---

## C. NOT STARTED ⬜

### Formats

1. TOML format support
2. JSONL format support
3. PlantUML format support
4. AsciiDoc format support

### Infrastructure

5. LSP workspace errors — 106 phantom errors from missing go.work (by design)
6. `go.work.example` → actual `go.work` auto-creation script or make target

### Documentation & Visibility

7. Post to r/golang
8. Submit to Awesome Go
9. Write blog post about Shape capability matrix

### Code Cleanup

10. Remove `FormatCategory`/`IsTableFormat`/`IsTreeFormat`/`IsGraphFormat`/`Category()` deprecated code
11. Remove `OutputFormat` type alias + constants + `ParseOutputFormat()`
12. Evaluate `cmdguard/` extraction to own repo
13. `internal/gentest` still exists — should be deleted or merged into testhelpers
14. `internal/testutils` still references `internal/gentest` — should use `testhelpers`

### API & Design

15. `FormatsForShape` integration with registry — runtime dispatch by shape
16. D2 `ShapeTable` renderer — `D2FromTableData` exists but isn't a `Renderer`
17. API stability audit — pre-v1.0 review of all exported symbols
18. Register HTMLTreeRenderer in registry

---

## D. TOTALLY FUCKED UP 💥

### 1. v0.4.0 Tag is Severely Stale

The commit tagged `v0.4.0` was `01c7e21` (MIT license change). Since then, **20+ commits** landed with major features (Shape matrix, JSON/YAML renderers, circular dep fix, structural refactoring, testhelpers module, Tree/Graph renderers). Anyone installing `@v0.4.0` gets code WITHOUT Shape, WITHOUT Tree/Graph renderers, WITH the circular dependency.

**Action needed:** Tag current HEAD as v0.5.0.

### 2. `internal/gentest` Still Exists

The `internal/gentest` package is still present and still used by 7 root test files. It was NOT deleted when `testhelpers/` was created — instead, `testhelpers/` was created as a parallel module that sub-modules can import. Root still uses `internal/gentest` for domain-specific helpers (HTMLEscapeTestRenderer, MarshalError, ExpectedOutput).

This is fine architecturally (root uses gentest, sub-modules use testhelpers) but creates confusion. The shared helpers (AssertStringSliceEqual, TestEnumIsValid, struct helpers) should ideally live only in `testhelpers/`.

### 3. `internal/testutils` Still References `internal/gentest`

The `internal/testutils/test_helpers.go` file re-exports from `internal/gentest`. This is a domain-specific wrapper package that should be updated to use `testhelpers/` for shared assertions.

---

## E. WHAT WE SHOULD IMPROVE 🔧

### Critical

1. **Tag v0.5.0** — current tag is 20+ commits behind reality, users get stale code
2. **Merge `internal/gentest` into `testhelpers/`** — eliminate the parallel test-helper packages

### High Impact

3. **CI/CD pipeline** — ✅ DONE (this session)
4. **JSON/YAML Tree+Graph renderers** — ✅ DONE (this session)
5. **Delete or keep sort/** — ✅ DECIDED: keep (zero-dep convenience wrapper)
6. **Remove deprecated FormatCategory/OutputFormat** — still present, still has tests

### Medium Impact

7. **README examples** — ✅ DONE (this session)
8. **r/golang post** — public visibility for the library
9. **Awesome Go submission** — discoverability
10. **Register HTMLTreeRenderer** — exists but not in registry

### Low Impact

11. **D2 ShapeTable Renderer** — wrap D2FromTableData in Renderer interface
12. **Registry + Shapes integration** — `GetRenderersForShape(shape) []Renderer`
13. **Blog post** about Shape capability matrix design
14. **Merge `internal/testutils` into `testhelpers/`** — unify test helper strategy

---

## F. TOP 25 THINGS TO DO NEXT 🎯

| #   | Priority | Task                                                     | Impact | Effort  |
| --- | -------- | -------------------------------------------------------- | ------ | ------- |
| 1   | **P0**   | **Tag v0.5.0** — 20+ commits past v0.4.0                 | High   | Trivial |
| 2   | **P0**   | **Merge `internal/gentest` into `testhelpers/`**         | Medium | Medium  |
| 3   | **P0**   | **Update `internal/testutils` to use `testhelpers`**     | Medium | Low     |
| 4   | P1       | **Remove deprecated FormatCategory code**                | Low    | Low     |
| 5   | P1       | **Remove deprecated OutputFormat aliases**               | Low    | Low     |
| 6   | P1       | **Post to r/golang**                                     | High   | Low     |
| 7   | P1       | **Submit to Awesome Go**                                 | High   | Low     |
| 8   | P1       | **Register HTMLTreeRenderer in registry**                | Low    | Trivial |
| 9   | P2       | **Add TOML format**                                      | Medium | Medium  |
| 10  | P2       | **Add JSONL format**                                     | Medium | Medium  |
| 11  | P2       | **Write blog post** about Shape matrix                   | Medium | High    |
| 12  | P2       | **D2 ShapeTable Renderer** — wrap D2FromTableData        | Low    | Low     |
| 13  | P2       | **Registry + Shapes integration**                        | Medium | Medium  |
| 14  | P2       | **Evaluate cmdguard/ extraction**                        | Low    | Low     |
| 15  | P3       | **Add PlantUML format**                                  | Low    | Medium  |
| 16  | P3       | **Add AsciiDoc format**                                  | Low    | Medium  |
| 17  | P3       | **API stability audit** — pre-v1.0 review                | Medium | Medium  |
| 18  | P3       | **Fix LSP workspace** — go.work.example → go.work script | Low    | Trivial |
| 19  | P3       | **Benchmark Tree/Graph renderers**                       | Low    | Trivial |
| 20  | P4       | **integration_test.go approaching limit** (340 lines)    | Low    | Low     |
| 21  | P4       | **d2_node_test.go** (306 lines) — approaching limit      | Low    | Low     |
| 22  | P4       | **Add Examples tests** (GoDoc-style) for key types       | Medium | Medium  |
| 23  | P4       | **Fuzz tests for JSON/YAML renderers**                   | Low    | Low     |
| 24  | P4       | **Performance benchmarks** for all Shape renderers       | Low    | Medium  |
| 25  | P4       | **Codecov/Istanbul coverage reporting** in CI            | Low    | Low     |

---

## G. TOP #1 QUESTION 🤔

**Should we tag v0.5.0 now with the Tree/Graph renderer story complete, or wait for the deprecated code removal (FormatCategory/OutputFormat purge) to make v0.5.0 a clean "breaking changes" release?**

Arguments for tagging now:

- 20+ commits with massive features since v0.4.0 tag
- Shape capability matrix, all 3 shapes now have typed renderers for JSON/YAML
- testhelpers module extraction (DRY fix)
- CI covering all 9 modules
- Anyone installing `@v0.4.0` gets broken code with circular deps

Arguments for waiting:

- Could remove deprecated FormatCategory/OutputFormat first — that's a breaking change that belongs in v0.5.0
- Could merge `internal/gentest` into `testhelpers/` for a cleaner story

**My recommendation:** Tag v0.5.0 NOW with a `[Unreleased]` section in CHANGELOG.md noting deprecated APIs will be removed in v0.6.0. The current state is significantly better than v0.4.0 and users shouldn't have to wait for a cleanup pass.

---

## Session 7 Changes Summary

### Net Delta Since v0.4.0 Tag

- **New module:** `testhelpers/` (9th module, zero deps, 100% coverage)
- **New files:** `json_renderers.go`, `json_renderers_test.go`, `yaml_renderers.go`, `yaml_renderers_test.go`, `format_shape_test.go`, `format_deprecated_test.go`, `testhelpers/helpers.go`, `testhelpers/helpers_test.go`
- **New public API:** `JSONTreeRenderer`, `JSONGraphRenderer`, `YAMLTreeRenderer`, `YAMLGraphRenderer`
- **New shared API:** `testhelpers.AssertStringSliceEqual`, `testhelpers.TestEnumIsValid`, `testhelpers.AssertEqual`, `testhelpers.AssertContains`, `testhelpers.TestStructFields`
- **DRY violations fixed:** 0 (was 2: assertStringSliceEqual duplicated in enum/ and cmdguard/)
- **Files over 350 lines:** 0 (was 2: format_test.go at 468, format_category_test.go would have been created)
- **Root coverage:** 90.3% → 90.9%
- **Lint issues:** 0 (was 0, but new code introduced and fixed 11 issues)
- **CI coverage:** 9/9 modules (was 1/9 — root only)
- **Module count:** 9 (root, enum, escape, testhelpers, sort, cmdguard, table, integration, examples)
