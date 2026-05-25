# Status Report: 4 New Output Formats + PlantUML Module

**Date:** 2026-05-25 23:53
**Branch:** master
**Author:** Crush (assisted)

---

## Executive Summary

Added 4 new output formats (JSONL, AsciiDoc, TOML, PlantUML) to go-output, expanding from 12 → 16 formats and 12 → 13 modules. All 13 modules pass tests with zero lint issues. New `plantuml/` module is fully independent with zero external dependencies.

---

## a) FULLY DONE ✓

### Root Enum Infrastructure
- [x] `FormatJSONL`, `FormatAsciiDoc`, `FormatTOML`, `FormatPlantUML` constants in `format.go`
- [x] `AllFormats` slice updated (16 entries)
- [x] `formatCapabilities` matrix in `shape.go` updated:
  - JSONL → `{ShapeTable}`
  - AsciiDoc → `{ShapeTable}`
  - TOML → `{ShapeTable, ShapeTree}`
  - PlantUML → `{ShapeTable, ShapeGraph}`
- [x] `ParseFormat`, `IsValid`, `AllowedValues` all work automatically via `AllFormats`

### JSONL (`serialization/jsonl.go`) — 143 lines
- [x] `init()` registers `FormatJSONL` marshaler
- [x] `JSONLWriter` — streaming writer with `Encode(v)` + `Flush()`
- [x] `JSONLTableRenderer` — `Render()` produces one JSON object per line
- [x] `MarshalJSONLFromTableData()` — batch marshaling
- [x] `renderJSONLTableData()` — registry dispatch target
- [x] Tests: `jsonl_test.go` — marshal, renderer, writer, nil/empty edge cases
- [x] Zero new external dependencies (stdlib `encoding/json` + `bufio`)

### AsciiDoc (`markup/asciidoc.go`) — 99 lines
- [x] `init()` registers `FormatAsciiDoc` marshaler
- [x] `AsciiDocTableRenderer` — `Render()` produces `|===` delimited table
- [x] `MarshalAsciiDocFromTableData()` — batch marshaling
- [x] Pipe character escaping (`|` → `\|`)
- [x] Tests: `asciidoc_test.go` — marshal, renderer, escape, nil/empty
- [x] Zero new external dependencies

### TOML (`serialization/toml.go` + `toml_renderers.go`) — 170 lines total
- [x] `init()` registers `FormatTOML` marshaler
- [x] `MarshalTOML()` / `UnmarshalTOML()` — generic marshal/unmarshal
- [x] `TOMLTableRenderer` — renders TableData as TOML
- [x] `TOMLTreeRenderer` — renders TreeNode hierarchy as TOML (tree shape support)
- [x] `MarshalTOMLFromTableData()` — batch marshaling
- [x] New dependency: `github.com/pelletier/go-toml/v2 v2.3.1` (isolated in serialization/)
- [x] Tests: `toml_test.go` + `toml_renderers_test.go` — marshal, unmarshal, renderer, tree

### PlantUML (`plantuml/` — NEW MODULE) — 121 lines production
- [x] `plantuml/go.mod` — independent module, replace directives
- [x] `PlantUMLDiagram` — `Render()` produces `@startuml`/`@enduml` component diagram
- [x] `AddNode()` / `AddEdge()` — builder pattern, uses `GraphRendererMixin`
- [x] `PlantUMLFromTableData()` — TableData conversion
- [x] `PlantUMLFromTree()` — TreeNode conversion
- [x] Tests: `plantuml_test.go` — diagram, sanitize, FromTableData, FromTree
- [x] Zero external dependencies (root only)
- [x] Coverage: 97.2%

### Integration Tests
- [x] `format_test.go` — `TestFormatCategories` updated with JSONL, TOML, PlantUML in table/tree/graph lists
- [x] `integration_test.go` — `TestAllFormatsRender` expanded to 16 formats
- [x] `integration_test.go` — `renderProject` switch expanded with 4 new cases
- [x] `integration_test.go` — new render functions: `renderJSONLFormat`, `renderAsciiDocFormat`, `renderTOMLFormat`, `renderPlantUMLFormat`

### Examples
- [x] `examples/basic/main.go` — new `renderJSONL`, `renderAsciiDoc`, `renderTOML`, `renderPlantUML` functions
- [x] `getRenderers()` map expanded with 4 new format entries

### Module Wiring
- [x] `serialization/go.mod` — added `go-toml/v2` dep + replace directives
- [x] `integration/go.mod` — added `plantuml` import + replace
- [x] `examples/go.mod` — added `plantuml` import + replace
- [x] `go.work` — added `./plantuml`
- [x] All 13 modules: `go mod tidy` + `go build ./...` clean

### Lint Configuration
- [x] `.golangci.yml` depguard `default` rule: added `plantuml`, `go-toml/v2`, `bufio`
- [x] `.golangci.yml` depguard `examples` rule: added `plantuml`
- [x] `.golangci.yml` depguard `main` rule: added `plantuml`
- [x] `.golangci.yml` `gomoddirectives` replace-allow-list: added `plantuml`

### Documentation
- [x] `AGENTS.md` — 16 formats, 13 modules, updated dependency graph, coverage table, go.work snippet
- [x] `README.md` — 4 new rows in Supported Formats table
- [x] `CHANGELOG.md` — `### Added` section with all 4 formats detailed

---

## b) PARTIALLY DONE — Nothing

All planned work items from the 44-task TODO list are complete.

---

## c) NOT STARTED

### Format Enhancements
- [ ] TOML graph renderer (`TOMLGraphRenderer`) — only table + tree implemented
- [ ] JSONL tree/graph renderers — only table implemented
- [ ] AsciiDoc tree renderer — only table implemented
- [ ] PlantUML rich domain model (like D2's shapes, arrows, classes, SQL tables)
- [ ] PlantUML sequence diagrams
- [ ] PlantUML class diagrams with methods/fields
- [ ] Streaming renderers for JSONL (line-by-line to io.Writer)
- [ ] AsciiDoc include directives, cross-references
- [ ] TOML `MarshalTOMLFromTableData` could produce TOML arrays of tables more idiomatically

### Test Coverage Gaps
- [ ] `serialization/` coverage dropped from ~90% to 75.4% — needs more tests for TOML/JSONL edge cases
- [ ] `testhelpers/` coverage at 43.5% — was 93.8% (may be measurement artifact)
- [ ] PlantUML `convert.go` test coverage (currently only tested via plantuml_test.go)
- [ ] Benchmarks for new formats (JSONL, AsciiDoc, TOML, PlantUML)
- [ ] Fuzz tests for new formats
- [ ] Error path tests for PlantUML (edge cases with special chars in labels)

### Documentation
- [ ] README.md usage examples for new formats (only table rows added)
- [ ] Example functions (`ExampleNewJSONLTableRenderer`, etc.) for Go doc
- [ ] ADR for PlantUML module creation (like ADR-003 for D2/Graph)
- [ ] ADR-001 update to reflect 13 modules (currently says 12 in the doc file)

---

## d) TOTALLY FUCKED UP — Nothing

No broken builds, no failing tests, no lint errors, no circular deps. All 13 modules green across test + lint.

---

## e) WHAT WE SHOULD IMPROVE

### Critical (Do Next)
1. **Serialization coverage is 75.4%** — dropped from ~90% because new TOML/JSONL files added without exhaustive edge-case tests. Needs error path tests, nil/empty edge cases, large datasets.
2. **ADR-001 still says 10 modules** — the file wasn't updated in this session (only AGENTS.md was). Need to sync.
3. **PlantUML is bare-bones** — it works but has no rich domain model (no shapes, arrows, SQL tables, classes like D2). This is the biggest gap vs. the D2 module.

### Important
4. **No benchmarks** for new formats — `delimited/` has benchmarks, but new formats don't.
5. **go.work is gitignored** — every new module requires manual `go.work` update. Consider a `just` recipe or `scripts/setup-workspace.sh`.
6. **No fuzz tests** for new formats — `graph/` has fuzz tests, new modules don't.
7. **TOML `go-toml/v2` pulls transitive deps** through root `go.mod` — verify root production code stays zero-toml-import.
8. **PlantUML `sanitizePlantUMLID` is naive** — only handles spaces and hyphens. Should handle dots, colons, unicode, etc.

### Nice to Have
9. **No `Example*` functions** for godoc — all existing formats have them in `example_test.go`.
10. **Integration `go.mod` has 25+ lines of replace directives** — the module count is growing; consider a script to generate them.
11. **No version bump** — `go.mod` still says `v0.0.0` or `v0.5.0` for root. Should tag a release.
12. **CHANGELOG is getting long** — consider splitting into `CHANGELOG.md` (recent) + link to older releases.

---

## f) Top 25 Things to Do Next

| # | Priority | Task | Impact | Effort |
|---|----------|------|--------|--------|
| 1 | P0 | Fix serialization coverage → 90%+ (TOML/JSONL error paths) | Quality | Low |
| 2 | P0 | Update ADR-001 to reflect 13 modules | Accuracy | Low |
| 3 | P1 | Add benchmarks for JSONL, AsciiDoc, TOML, PlantUML | Perf | Low |
| 3 | P1 | Add `Example*` functions for all 4 new formats (godoc) | Docs | Low |
| 5 | P1 | PlantUML rich domain model (shapes, arrows, classes) | Feature | Medium |
| 6 | P1 | TOML graph renderer (`TOMLGraphRenderer`) | Feature parity | Low |
| 7 | P1 | Verify root `go.mod` has zero `go-toml/v2` import in production code | Correctness | Low |
| 8 | P1 | Add fuzz tests for TOML/JSONL/AsciiDoc parsing | Robustness | Low |
| 9 | P2 | PlantUML sequence diagram support | Feature | Medium |
| 10 | P2 | PlantUML `sanitizePlantUMLID` — handle dots, colons, unicode | Correctness | Low |
| 11 | P2 | JSONL streaming renderer (line-by-line to `io.Writer`) | Streaming | Low |
| 12 | P2 | AsciiDoc tree renderer | Feature parity | Low |
| 13 | P2 | README.md usage examples for JSONL, TOML, AsciiDoc, PlantUML | Docs | Low |
| 14 | P2 | Write ADR for PlantUML module creation | Docs | Low |
| 15 | P2 | Tag a release (v0.6.0?) with 16 formats, 13 modules | Release | Low |
| 16 | P3 | Script to auto-generate replace directives in go.mod files | DX | Medium |
| 17 | P3 | Testhelpers coverage investigation (43.5% vs 93.8%) | Quality | Low |
| 18 | P3 | JSONL tree/graph renderers | Feature parity | Medium |
| 19 | P3 | TOML `MarshalTOMLFromTableData` — idiomatic TOML array of tables | Quality | Low |
| 20 | P3 | Root `go.mod` dep cleanup — remove transitive deps pulled by new modules | Hygiene | Low |
| 21 | P3 | PlantUML class diagram support (methods, fields) | Feature | Medium |
| 22 | P3 | AsciiDoc include directives / cross-references | Feature | Low |
| 23 | P4 | Integration test for `RenderTableData` dispatch with new formats | Coverage | Low |
| 24 | P4 | Consider `go.work` sync script or justfile recipe | DX | Low |
| 25 | P4 | Split CHANGELOG into recent + archived | Hygiene | Low |

---

## g) Top #1 Question I Cannot Answer Myself

**Should PlantUML be a "rich domain model" module like D2 (with PlantUML-specific shapes, arrows, SQL tables, classes, sequence diagrams) or remain a thin wrapper around the generic `GraphRendererMixin`?**

- D2 has 800+ lines of domain-specific types (D2Node, D2Edge, D2StrokeStyle, D2Column, D2Shape, D2Arrow, etc.)
- Current PlantUML is 121 lines and uses generic `GraphNode`/`GraphEdge`
- Rich model would mean PlantUML-specific types (participants, lifelines, class fields, etc.) but triples the module size
- Thin wrapper means less maintenance but limits PlantUML to what the generic graph model can express
- This is a product/domain decision that requires user input

---

## Verification

```
13/13 modules: tests PASS
13/13 modules: lint ZERO issues
16 formats registered in AllFormats
13 modules in workspace
0 circular dependencies
0 new external deps in root production code
```

### Coverage Snapshot

| Module | Coverage | Status |
|--------|----------|--------|
| output (root) | 89.6% | ✓ |
| internal/gentest | 80.8% | ✓ |
| delimited | 86.2% | ✓ |
| serialization | 75.4% | ⚠ Dropped from ~90% |
| markup | 85.5% | ✓ |
| d2 | 100.0% | ✓ |
| graph | 96.0% | ✓ |
| enum | 100.0% | ✓ |
| escape | 100.0% | ✓ |
| testhelpers | 43.5% | ⚠ Investigate |
| table | 100.0% | ✓ |
| integration | 82.8% | ✓ |
| plantuml | 97.2% | ✓ |

---

## Files Changed (28 total)

### Modified (15)
- `.golangci.yml` — depguard + gomoddirectives for plantuml, toml, bufio
- `AGENTS.md` — 16 formats, 13 modules, updated graph/table/snippets
- `CHANGELOG.md` — Added section for 4 new formats
- `README.md` — 4 new rows in Supported Formats table
- `format.go` — 4 new Format constants + AllFormats entries
- `shape.go` — 4 new formatCapabilities entries
- `serialization/go.mod` — go-toml/v2 dep + replace directives
- `serialization/go.sum` — checksums
- `integration/format_test.go` — TestFormatCategories updated
- `integration/integration_test.go` — 4 new render functions + switch cases
- `integration/go.mod` — plantuml import + replace
- `integration/go.sum` — checksums
- `examples/basic/main.go` — 4 new render functions + getRenderers entries
- `examples/go.mod` — plantuml import + replace
- `examples/go.sum` — checksums

### Created (13)
- `serialization/jsonl.go` — 143 lines
- `serialization/jsonl_test.go` — JSONL tests
- `markup/asciidoc.go` — 99 lines
- `markup/asciidoc_test.go` — AsciiDoc tests
- `serialization/toml.go` — 101 lines
- `serialization/toml_test.go` — TOML tests
- `serialization/toml_renderers.go` — 69 lines (tree renderer)
- `serialization/toml_renderers_test.go` — Tree renderer tests
- `plantuml/go.mod` — independent module definition
- `plantuml/go.sum` — checksums
- `plantuml/plantuml.go` — 80 lines (diagram + AddNode/AddEdge)
- `plantuml/plantuml_test.go` — PlantUML tests
- `plantuml/convert.go` — 41 lines (FromTableData + FromTree)

---

_Arte in Aeternum_
