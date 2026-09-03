# Deduplication Sprint #2 — Status Report

**Date:** 2026-05-25 23:02
**Branch:** master
**Commit:** 5f96c23 (clean, pre-session)

---

## Session Summary

Continued the semantic deduplication effort using `art-dupl --semantic --sort total-tokens -t 15`. Reduced clone groups from **44 → 29** (34% reduction, net -55 lines). All 12 modules build, test, and lint cleanly.

---

## a) FULLY DONE ✅

### Production Code Deduplication

| What                                | Files                                                                             | Change                                                                 |
| ----------------------------------- | --------------------------------------------------------------------------------- | ---------------------------------------------------------------------- |
| Shared `renderDelimitedTableData()` | `delimited/csv.go`, `delimited/tsv.go`                                            | Extracted shared helper; CSV/TSV init() closures now call one function |
| Shared `renderTable()`              | `serialization/render.go` (NEW), `serialization/json.go`, `serialization/yaml.go` | Extracted table rendering logic shared between JSON and YAML renderers |
| `writeRowBoundary()`                | `markup/streaming.go`                                                             | Extracted row start/end chunk writing into single helper               |

### Test Code Deduplication

| What                                                             | Files                                                                                   | Change                                                                                    |
| ---------------------------------------------------------------- | --------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `buildBenchmarkTree()`                                           | `benchmarks_test.go`                                                                    | Eliminated duplicate tree construction in ASCIITreeRenderer + ASCIITreeColored benchmarks |
| `assertSliceLen[T]()`                                            | `graph_mixin_test.go`                                                                   | Generic helper replacing 6 identical len-assertion patterns                               |
| `testRenderTableDataWriterError()` + table-driven partial writes | `render_tabledata_test.go`                                                              | Consolidated 4 writer-error tests into helper + table-driven pattern                      |
| `newTestGraphStyle()`                                            | `graph/graph_test.go`                                                                   | Shared GraphStyle test constructor                                                        |
| `unmarshalTestCase` + `testUnmarshalCases()`                     | `serialization/testhelpers_test.go`, `json_test.go`, `yaml_test.go`                     | Shared unmarshal test infrastructure                                                      |
| `graphRenderer` interface + `testGraphRendererNodeWithShape()`   | `serialization/testhelpers_test.go`, `json_renderers_test.go`, `yaml_renderers_test.go` | Shared graph renderer shape test                                                          |
| `assertColorMode()`                                              | `integration/integration_test.go`                                                       | Consolidated 4 ColorMode subtests into table-driven helper                                |
| Table-driven mismatch tests                                      | `integration/test_helpers_test.go`                                                      | Merged 2 identical assertTableData mismatch patterns                                      |
| `htmlEscapeTestRenderer` type + table-driven HTML escape         | `markup/testhelpers_test.go`, `html_test.go`, `streaming_test.go`                       | Named interface replacing anonymous inline type; table-driven test cases                  |
| `assertEscaped()`                                                | `graph/fuzz_test.go`                                                                    | Extracted MermaidText escape assertion                                                    |
| `addBenchmarkD2Tables()`                                         | `d2/bench_test.go`                                                                      | Shared benchmark table-addition helper                                                    |
| `newMarkdownTableWithData()` + `newMarkdownTableWithSingleRow()` | `markdown_test.go`                                                                      | Extracted markdown test table constructors                                                |

### Quality Gates

- ✅ All 12 modules: `go build ./...` — clean
- ✅ All 12 modules: `go test ./...` — all pass
- ✅ All 12 modules: `golangci-lint run ./...` — 0 issues
- ✅ Root coverage: 89.6%
- ✅ 4 modules at 100% coverage: d2, enum, escape, table

---

## b) PARTIALLY DONE 🔧

### Clone Reduction (44 → 29 remaining)

The 29 remaining clone groups are **structural/architectural** and fall into three categories:

**Cross-module test helpers (12 groups):** Graph node constructors (`newTestNode`, `testNodesAB`, `testEdgeAB`, `testNodesABC`, `testEdgesABC`) are duplicated across `graph/helpers_test.go`, `serialization/testhelpers_test.go`, and `output_test_helpers_test.go`. These CANNOT be shared because:

- `testhelpers/` has **zero dependencies** (architectural constraint)
- Each package's test files are isolated (Go restriction)
- `internal/gentest` is root-only (can't be imported by sub-modules)

**Inherent Go patterns (9 groups):**

- Function signatures: `render_tabledata.go` dispatch functions, `markup/streaming.go` method signatures
- `init()` closures: `delimited/csv.go`, `delimited/tsv.go`, `markup/xml.go`, `serialization/yaml.go` — RegisterTableDataMarshaler pattern
- `d2/d2.go` `isSet()` methods on different struct types
- `testhelpers/helpers.go` function signatures

**Small assertion patterns across modules (8 groups):**

- 3-line `strings.Contains` assertions in test files across delimited, graph, serialization, table, tree
- Table-driven test struct literal similarity
- Single-line references (e.g., `markup/markup_test.go`, `markup/xml_test.go`)

---

## c) NOT STARTED ⏳

- No further deduplication work planned — remaining 29 groups are architectural limits
- No new features or formats in this session

---

## d) TOTALLY FUCKED UP 💥

Nothing. All changes compile, test, and lint cleanly. Zero regressions.

---

## e) WHAT WE SHOULD IMPROVE 📈

1. **testhelpers coverage at 61.2%** — lowest in the project. `TestAllowedValues`, `TestStructFields`, `IntField` are only tested via other modules.
2. **gentest at 80.8%** — the `AssertHTMLEscape` helper is only tested indirectly through markup tests.
3. **delimited at 86.2%** — some error paths in `DelimitedWriter` could use more coverage.
4. **serialization at 83.3%** — JSON/YAML writer error paths need direct tests.
5. **markup at 86.9%** — streaming error paths partially covered.
6. **integration at 82.8%** — cross-module edge cases could be expanded.
7. **No benchmarks for delimited/serialization/markup modules** — only root, d2, and table have benchmarks.
8. **No fuzz tests for markup (HTML/XML escaping)** — only d2 and graph have fuzz coverage.
9. **`go.work` is gitignored** — contributors must create it manually; could provide a `go.work.example`.
10. **examples/shared has 0% coverage** — example helpers are untested.
11. **No `CHANGELOG.md` entry for the deduplication work** — should document the refactoring.
12. **`docs/adr/` doesn't have an ADR for the zero-dep testhelpers decision** — this is the key reason 29 clones remain.

---

## f) Top 25 Things to Do Next

| #  | Priority | Task                                                                    | Impact               |
| -- | -------- | ----------------------------------------------------------------------- | -------------------- |
| 1  | 🔴 HIGH  | Add direct tests for `testhelpers` to push coverage to 90%+             | Quality gate         |
| 2  | 🔴 HIGH  | Add fuzz tests for `escape.HTML()` and `escape.XML()`                   | Security             |
| 3  | 🔴 HIGH  | Add `CHANGELOG.md` entry for v0.x deduplication work                    | Documentation        |
| 4  | 🟡 MED   | Write ADR 003: Zero-dep testhelpers boundary                            | Architecture clarity |
| 5  | 🟡 MED   | Add benchmarks for `delimited` module (CSV/TSV throughput)              | Performance baseline |
| 6  | 🟡 MED   | Add benchmarks for `serialization` module (JSON/YAML throughput)        | Performance baseline |
| 7  | 🟡 MED   | Add benchmarks for `markup` module (HTML/XML throughput)                | Performance baseline |
| 8  | 🟡 MED   | Add error-path tests for `delimited` to push coverage to 90%+           | Coverage             |
| 9  | 🟡 MED   | Add error-path tests for `serialization` to push coverage to 90%+       | Coverage             |
| 10 | 🟡 MED   | Add streaming error-path tests for `markup` to push to 90%+             | Coverage             |
| 11 | 🟡 MED   | Create `go.work.example` for contributors                               | DX                   |
| 12 | 🟡 MED   | Add integration tests for edge cases (empty rows, single column)        | Coverage             |
| 13 | 🟡 MED   | Test `examples/shared` package directly                                 | Coverage             |
| 14 | 🟡 MED   | Add `// Example` functions to key public APIs for godoc                 | Documentation        |
| 15 | 🟢 LOW   | Investigate `gob` format support (binary serialization)                 | Feature              |
| 16 | 🟢 LOW   | Add `FormatString()` method to Format enum for human-readable names     | DX                   |
| 17 | 🟢 LOW   | Add `RenderToWriter()` convenience function                             | DX                   |
| 18 | 🟢 LOW   | Profile memory allocations in hot paths (table rendering)               | Performance          |
| 19 | 🟢 LOW   | Add `cio` (Colored IO) streaming renderer with ANSI support             | Feature              |
| 20 | 🟢 LOW   | Consider `text/tabwriter` integration for aligned terminal output       | Feature              |
| 21 | 🟢 LOW   | Add `README.md` badges (coverage, godoc, Go version)                    | Presentation         |
| 22 | 🟢 LOW   | Add `CONTRIBUTING.md` with dev setup instructions                       | DX                   |
| 23 | 🟢 LOW   | Explore `golang.org/x/term` terminal width detection for table wrapping | Feature              |
| 24 | 🟢 LOW   | Add property-based testing with `rapid` for format round-trips          | Quality              |
| 25 | 🟢 LOW   | Evaluate `slog` handler for structured logging output                   | Feature              |

---

## g) Top #1 Question I Cannot Answer Myself 🤔

**Should the zero-dep `testhelpers/` constraint be relaxed to allow importing the root `output` package?**

This would let us extract shared graph-node test constructors (`newTestNode`, `testNodesAB`, `testEdgeAB`, etc.) into `testhelpers/`, eliminating 12 of the 29 remaining clone groups. But it would break the "zero transitive deps" promise documented in AGENTS.md.

The tradeoff:

- **Pro:** -12 clone groups, cleaner test code, single source of truth for test data
- **Con:** `testhelpers` would depend on `output`, meaning `go get testhelpers` pulls in the root module's deps (x/term, branded-id). But `testhelpers` is a test-only package — no production user should import it.

**My recommendation:** Keep the constraint. The 29 remaining clones are all structural and don't indicate real duplication. They're the cost of clean module boundaries. But I'd like explicit confirmation.

---

## Metrics Snapshot

| Metric                      | Value                              |
| --------------------------- | ---------------------------------- |
| Total Go LOC                | 14,058                             |
| Clone groups (pre-session)  | 44                                 |
| Clone groups (post-session) | 29                                 |
| Clone reduction             | 34% (-15 groups)                   |
| Net line change             | -55 lines (264 added, 319 removed) |
| Files changed               | 22 (+1 new)                        |
| Modules                     | 12                                 |
| Lint issues                 | 0                                  |
| Test failures               | 0                                  |

### Coverage by Module

| Module        | Coverage |
| ------------- | -------- |
| d2            | 100.0%   |
| enum          | 100.0%   |
| escape        | 100.0%   |
| table         | 100.0%   |
| graph         | 96.0%    |
| root (output) | 89.6%    |
| markup        | 86.9%    |
| delimited     | 86.2%    |
| serialization | 83.3%    |
| integration   | 82.8%    |
| gentest       | 80.8%    |
| testhelpers   | 61.2%    |
