# Deduplication Sprint 3 — Status Report

**Date:** 2026-05-28 11:05  
**Scope:** Full codebase deduplication at `art-dupl -t 15`  
**Result:** 60 → 53 clone groups (net -7, ~12% reduction)

---

## a) FULLY DONE

### Production Code (4 changes)

| File | Change | Lines Saved |
|------|--------|-------------|
| `markup/asciidoc.go` | Extracted `writeAsciiDocCells()` — row/footer cell loop | -8 |
| `markdown.go` | Extracted `updateMaxWidths()` — rows/footer width calc | -5 |
| `d2/d2.go` | `D2Edge.hasBlockAttrs()` reuses `D2StrokeStyle.isSet()` | -3 |
| `serialization/render.go` | Extracted `renderViaRenderer()` — YAML/TOML identical pattern | -25 |

### Test Code (6 changes)

| File | Change | Lines Saved |
|------|--------|-------------|
| `delimited/testhelpers_test.go` | Replaced local `assertContains` with `testhelpers.AssertContains` alias | -6 |
| `testing_test.go` + `graph/helpers_test.go` | Removed `testParseEnum`/`testEnumString`/`testAllowedValues` wrappers | -50 |
| `serialization/toml_renderers_test.go` | Uses `testNodesAB()`/`testEdgesAB()`/`newTestNode()` | -20 |
| `serialization/error_test.go` | Removed 3 duplicate NilData tests | -44 |
| `d2/fuzz_test.go` + `graph/fuzz_test.go` | Uses `graphtest.AssertEscape()` | -8 |
| `serialization/registry_test.go` + `markup/registry_test.go` | Table-driven NilData/WriterError (9+9 → 2+2 tests) | -112 |

### Total: 4 commits, 20 files, **-281 net lines removed**, all 13 modules build+test pass.

---

## b) PARTIALLY DONE

- **PlantUML graphtest adoption**: Serialization tests now use graphtest helpers, but plantuml tests still use inline node construction (different test data "svc-a"/"svc-b" vs "A"/"B" — changing would alter test semantics)
- **Registry test dedup**: serialization and markup are table-driven, but delimited still uses its internal `renderDelimitedTableData` function (different pattern, can't unify)

---

## c) NOT STARTED

These items from the plan were analyzed but not executed:

| # | Item | Reason |
|---|------|--------|
| 1 | `AssertLineCount` helper in testhelpers | testhelpers is zero-dep; can't add Format param. Local helpers sufficient. |
| 2 | PlantUML bench/example graphtest helpers | Would require adding graphtest dep for 2 nodes in bench. Overkill. |
| 3 | `testEmptyRendererOutput` extraction to testhelpers | In different modules (graph vs markup); can't share without cross-dep. |
| 4 | `render*TableData` init() registration unification | Each is 1-line unique binding. Category E (single-line). |
| 5 | Streaming HTML `writeHeaders`/`writeRow`/`writeFooter` dedup | Different tag names (`<th>` vs `<td>`), different error messages. Not worth abstracting. |

---

## d) TOTALLY FUCKED UP

Nothing. All changes compile, pass tests, and pass lint.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture Observations

1. **`testhelpers` zero-dep constraint limits cross-module sharing**: The `testhelpers` module intentionally has zero deps (not even `output`). This prevents adding `AssertNilDataRendersEmpty(format)` or `AssertLineCount(output)`. Each module must keep local table-driven wrappers. **This is the right tradeoff** — zero deps means users who `go get` testhelpers get nothing extra.

2. **`render*TableData` registration pattern is boilerplate**: Every format module has identical `init() + render*TableData` glue. A generic `output.RegisterTableDataMarshalerFunc(format, rendererFactory)` could eliminate this, but would couple the registry to the Renderer interface. Current approach is more flexible.

3. **Multi-module workspace has inherent duplication**: Go's module boundaries force interface re-declaration, type aliasing, and helper re-exporting across modules. This is **by design** — each module is independently versionable. The ADR policy should explicitly accept this category.

### Process Improvements

4. **Threshold 15 is too aggressive for production use**: The SKILL.md says 15 is "aggressive — catches incidental patterns that are Go idioms, not duplication." We should use **threshold 30-40** for actionable dedup work. At t=15, 80%+ of clones are Go test idioms or module boundary noise.

5. **Should have read ADR policy before starting**: The project already has nuanced duplicate-code guidance. I should have checked `docs/adr/` first.

---

## f) Top 25 Things We Should Get Done Next

### High Impact, Low Effort (do now)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 1 | Update AGENTS.md with dedup learnings (zero-dep testhelpers, threshold guidance) | 5min | High |
| 2 | Add ADR for acceptable clone categories at t=15 | 15min | High |
| 3 | Run art-dupl at t=30 to find genuinely actionable clones | 2min | High |
| 4 | Fix pre-existing lint issues in serialization/error_test.go (errchkjson, wsl) | 5min | Medium |
| 5 | Fix pre-existing lint issues in table/color_test.go (staticcheck, wsl) | 5min | Medium |

### Medium Impact, Medium Effort (plan for next sprint)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 6 | Extract `renderMarshalAndWrite(w, marshalFn, formatName)` for AsciiDoc/XML pattern | 30min | Medium |
| 7 | Table-drive delimited NoHeaders tests (CSV + TSV identical) | 15min | Medium |
| 8 | Extract shared `testGraphRendererEmpty` helper to graphtest (graph+plantuml+serialization) | 20min | Medium |
| 9 | Consolidate `emptyYAML`/`emptyTOML` sentinel constants in serialization | 10min | Low |
| 10 | Add `generateBenchmarkNodes(n)` to graphtest (graph+plantuml benches) | 20min | Low |

### Lower Priority (backlog)

| # | Task | Effort | Impact |
|---|------|--------|--------|
| 11 | Investigate generic `RegisterSimpleMarshaler(format, func(data) []byte)` | 1hr | Medium |
| 12 | Unify streaming.go HTML cell writing (`<th>`/`<td>` with templates) | 1hr | Low |
| 13 | Consider `go-structure-linter` suppressions for root package files | 15min | Low |
| 14 | Add `go-error-family` dependency for structured errors | 1hr | Medium |
| 15 | Migrate `go.work.example` to auto-generated from go.mod scan | 30min | Low |
| 16 | Review D2 `D2NodeStyle.isSet()` vs `D2StrokeStyle.isSet()` overlap | 20min | Low |
| 17 | Consider `cmp.Diff` for test assertions instead of manual `strings.Contains` | 2hr | Medium |
| 18 | Add `testify` or stay with stdlib — document decision in ADR | 1hr | High |
| 19 | Review examples/ for consistency (some use graphtest, some don't) | 30min | Low |
| 20 | Run coverage report across all modules | 5min | Medium |
| 21 | Check if `go-faster/yaml` or `go-toml/v2` have newer versions | 10min | Low |
| 22 | Review `graph/fuzz_test.go` escape predicate vs `escape/escape.go` duplication | 15min | Low |
| 23 | Consider `go:generate stringer` for enum types | 1hr | Low |
| 24 | Document the `dataSetter` interface pattern in serialization/render.go | 5min | Low |
| 25 | Clean up `coverage.out` in root (go-structure-linter complaint) | 2min | Low |

---

## g) Top Question

**Should we adopt `testify/assert` (or similar) for test assertions?**

Currently every test uses manual `if !strings.Contains(got, want) { t.Errorf(...) }` which is the #1 source of "clones" at threshold 15 (36+ occurrences). Using `assert.Contains(t, got, want)` would:
- Eliminate ~20 clone groups at t=15
- Improve test readability significantly
- But add a dependency to every test package

The project already has `testhelpers.AssertContains` which is a subset. Should we:
1. Keep stdlib-only assertions (current approach)
2. Adopt `testify/assert` project-wide
3. Expand `testhelpers` with more assertion helpers (keeping zero-dep)

---

## Metrics

| Metric | Before | After | Delta |
|--------|--------|-------|-------|
| Clone groups (t=15) | 60 | 53 | -7 (-12%) |
| Total lines | ~8,500 | ~8,219 | -281 |
| Production clones fixed | 4 | 4 | Done |
| Test clones fixed | 6 | 6 | Done |
| Modules passing | 13/13 | 13/13 | 100% |
