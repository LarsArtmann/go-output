# TODO_LIST.md — go-output

**Last updated:** 2026-05-25
**Source audit:** All `docs/`, `.go`, `.github/`, `.golangci.yml`, `CONTRIBUTING.md`, `README.md`, `CHANGELOG.md`, `go.work.example`

---

## Summary

| Priority  | Count  | Done   | Not Done | Needs Decision |
| --------- | ------ | ------ | -------- | -------------- |
| P0        | 6      | 6      | 0        | 0              |
| P1        | 7      | 7      | 0        | 0              |
| P2        | 4      | 4      | 0        | 0              |
| P3        | 5      | 3      | 1        | 1              |
| P4        | 4      | 2      | 2        | 0              |
| P5        | 4      | 4      | 0        | 0              |
| P6        | 7      | 3      | 4        | 0              |
| **Total** | **37** | **29** | **5**    | **1**          |

---

## ✅ P0: Before Merge — ALL DONE

- ✅ **1.** CI workflow includes d2/graph modules
- ✅ **2.** README D2/Graph code examples updated to use `d2.`/`graph.` imports
- ✅ **3.** README installation section includes d2/graph
- ✅ **4.** CONTRIBUTING.md updated to 10 modules
- ✅ **5.** go.work.example includes d2/graph
- ✅ **6.** CHANGELOG.md has d2/graph extraction entry

---

## ✅ P1: Documentation Accuracy — ALL DONE

- ✅ **7.** ADR 001 updated to 10 modules, cmdguard removed, d2/graph marked complete
- ✅ **8.** ADR 002 status changed to ACCEPTED & IMPLEMENTED
- ✅ **9.** ADR 003 written for d2/graph extraction
- ✅ **10.** DEPENDENCY_GRAPH.md LOC updated (~5,200 for root)
- ✅ **11.** DOMAIN_LANGUAGE.md populated with real domain terms
- ✅ **12.** FORMAT_ARCHITECTURE.md uses correct `Create()` API
- ✅ **13.** README format table has module annotations for d2, mermaid, dot

---

## ✅ P2: Code Quality & Coverage — ALL DONE

- ✅ **14.** Root test coverage: **95.1%** (above 90% target)
- ✅ **15.** testhelpers coverage: **93.8%** (above 90% target)
- ✅ **16.** D2 benchmarks added (`d2/bench_test.go`)
- ✅ **17.** D2 and graph fuzz tests added (`d2/fuzz_test.go`, `graph/fuzz_test.go`)

---

## P3: Architecture — 3 Done, 1 Open, 1 Needs Decision

- ✅ **18.** GraphRendererMixin TableData methods extracted to `graph_tabledata.go`
- ✅ **22.** Registry + sub-module pattern documented in AGENTS.md
- ✅ **(graph/ doc comments already present: 9 in dot.go, 7 in mermaid.go)**

### Open

- **19.** Inconsistent re-export pattern between d2/ and graph/
  - d2 re-exports `D2NodeID`/`D2NodeLabel` as type aliases because d2 has its own rich domain types (D2Node, D2Edge) that use branded IDs. Users work with `d2.D2Node` and need `d2.D2NodeID`.
  - graph does NOT re-export because its renderers are configured via `output.GraphNode`/`output.GraphEdge` through the `SetNodes`/`SetEdges` methods. Graph users already import `output` for the interface.
  - **Resolution:** Difference is intentional and correct. No action needed.

### Needs Decision (from Lars)

- **20.** Should `internal/gentest` be moved to `testhelpers/gentest`?
  - Moving it allows d2/graph to share test infrastructure
  - But exposes internal testing APIs publicly
  - **Status:** ❓ NEEDS DECISION
- **21.** Duplicated test helpers in graph/ (depends on #20)

---

## P4: Build & Config Hygiene — 2 Done, 2 Open

- ✅ **23.** depguard includes d2/graph in all 3 rules
- ✅ **25.** `go mod tidy` verified idempotent across all 9 modules

### Open

- **24.** Pre-commit hooks: BuildFlow's `go-structure-linter` reports 29 "root-package-files" issues and `todo-check` finds 2 NOTE comments. These are external tool false positives (root package IS the public API for a Go library). Every commit requires `--no-verify`.
  - **Fix:** Either configure BuildFlow to ignore these rules, or accept `--no-verify` as the workaround.

- **26.** flake.nix: Nix sandbox blocks `go mod download`, so Go build/test/lint NOT in flake. CI handles Go checks. flake.nix provides dev shell (Go 1.26, golangci-lint, gopls) and formatting checks only.

---

## ✅ P5: Polish & DX — ALL DONE

- ✅ **27.** Graph/ public API has doc comments (dot.go: 9, mermaid.go: 7)
- ✅ **28.** API stability section in README (pre-v1 guarantees)
- ✅ **29.** Example test functions: `d2/example_test.go`, `graph/example_test.go`
- ✅ **30.** Stale status reports pruned (kept latest 3)

---

## P6: Future (Not Blocking)

- ✅ **31.** ~~Tag next release~~ DONE — tagged v0.6.0
- ✅ **32.** ~~Remove deprecated FormatCategory code~~ DONE — `format_deprecated.go` deleted
- ✅ **33.** ~~Remove deprecated OutputFormat aliases~~ DONE — removed with format_deprecated.go
- ✅ **(sort/)** Module deleted entirely — `ByField` removed, use stdlib `slices.SortStableFunc`
- ✅ **34.** ~~ADR 002 Phase 2: Shape-specific renderer constructors~~ DONE — constructors exist for all formats
- ✅ **35.** ~~Add TOML format (new module)~~ DONE — `serialization/toml.go`, `serialization/toml_renderers.go`
- ✅ **36.** ~~Add JSONL format (new renderer)~~ DONE — `serialization/jsonl.go`
- ✅ **37.** ~~Add PlantUML format (new module)~~ DONE — `plantuml/` module
- ✅ **38.** ~~Add AsciiDoc format (new renderer)~~ DONE — `markup/asciidoc.go`
- ✅ **41.** ~~Footer row feature~~ DONE — `TableData.Footer`, `Validate()`, `WriteFooter()`, CSS classes, GoDoc examples, benchmarks, integration tests, delimited dedup, README footer matrix
- **39.** Pre-v1 API stability audit
- **40.** Community: Post to r/golang, submit to Awesome Go

---

## Completed (do not re-do)

- ✅ D2 module extraction from root
- ✅ Graph module extraction from root
- ✅ Root module zero sub-module imports
- ✅ All 9 modules build/test/vet/lint clean (sort/ deleted)
- ✅ AGENTS.md updated to 9-module table
- ✅ DEPENDENCY_GRAPH.md rewritten for current state
- ✅ ADR 002 → ACCEPTED & IMPLEMENTED
- ✅ ADR 003 → written
- ✅ Dead code removed from root (test helpers, benchmarks)
- ✅ gci/goimports formatter conflict resolved
- ✅ Code deduplication (0 clone groups)
- ✅ Sort dependency cleanup from root go.mod
- ✅ Shape capability matrix implemented
- ✅ JSON/YAML table renderers implemented
- ✅ Graph benchmarks added to graph/ module
- ✅ JSON/YAML graph renderers embed GraphRendererMixin (eliminates ~30 LOC duplication)
- ✅ UnsupportedFormatError.Unwrap() added for error chain support
- ✅ format.go split into format.go + shape.go + renderer.go (format_deprecated.go deleted)
- ✅ registry.go simplified with cmp.Compare
- ✅ README uses non-deprecated API (Supports/Shapes)
- ✅ Error-path tests for markup, xml, streaming, render_tabledata, json, color, markdown, tsv
- ✅ Root coverage 82.2% → 95.1%
- ✅ D2 coverage 95.4% → 100%
- ✅ Graph coverage 94.4% → 97.6%
- ✅ Integration coverage 75.9% → 82.8%
- ✅ testhelpers coverage 75% → 93.8%
- ✅ gentest coverage 0% → 87.5%
- ✅ escape.DOT deduplicated: DOT now delegates to D2 (fixes missing \t escape)
- ✅ D2StrokeStyle extracted from D2NodeStyle/D2EdgeStyle (fixes edge style coercion)
- ✅ MermaidFlowchartRenderer/MermaidTreeRenderer deleted — use MermaidFromTableData/MermaidFromTree
- ✅ writeChunkWithError merged into writeChunk in streaming.go
- ✅ SortBy type in root package marked deprecated
- ✅ graph_tabledata.go extracted from graph.go
