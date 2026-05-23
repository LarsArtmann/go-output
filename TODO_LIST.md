# TODO_LIST.md — go-output

**Last updated:** 2026-05-23
**Source audit:** All `docs/`, `.go`, `.github/`, `.golangci.yml`, `CONTRIBUTING.md`, `README.md`, `CHANGELOG.md`, `go.work.example`

---

## P0: Before Merge (Blocking)

### 1. CI workflow missing d2/graph modules
**File:** `.github/workflows/ci.yml`
**Problem:** Build/test/mod-tidy/govulncheck loops iterate `. enum escape testhelpers sort table integration examples` — missing `d2` and `graph`. CI will pass on master but not test 2 of 10 modules.
**Fix:** Add `d2 graph` to all 4 loops in ci.yml.

### 2. README.md: D2/Graph code examples are stale
**File:** `README.md:202-209, 251-280`
**Problem:** D2 examples reference `output.NewD2Renderer("Architecture")` (does not exist), `output.D2Node`, `output.D2ShapeHexagon`, `output.D2Table`, `output.D2Column` — all now in `d2/` module. Graph examples reference `output.DOTFromTableData()` and `output.MermaidFlowchartRenderer()` — now in `graph/` module.
**Fix:** Update all code examples to import and use `d2.NewD2Diagram()`, `d2.D2Node`, `graph.DOTFromTableData()`, `graph.MermaidFlowchartRenderer()`.

### 3. README.md: Installation section missing d2/graph
**File:** `README.md:236-244`
**Problem:** Only shows `go get go-output` and `go get go-output/table`. Missing `go get go-output/d2` and `go get go-output/graph`.
**Fix:** Add installation instructions for d2 and graph sub-modules.

### 4. CONTRIBUTING.md: Outdated module list
**File:** `CONTRIBUTING.md:16`
**Problem:** Says "8 independent modules" and go.work snippet omits `d2` and `graph`.
**Fix:** Update to 10 modules, add `./d2` and `./graph` to go.work snippet.

### 5. go.work.example: Missing d2/graph
**File:** `go.work.example`
**Problem:** Only lists 8 modules, missing `./d2` and `./graph`.
**Fix:** Add both.

### 6. CHANGELOG.md: No entry for d2/graph extraction
**File:** `CHANGELOG.md`
**Problem:** `[Unreleased]` section doesn't mention d2/graph module extraction.
**Fix:** Add entry under `[Unreleased]` noting d2/ and graph/ are now independent modules with import path changes.

---

## P1: Documentation Accuracy

### 7. ADR 001: Stale — references cmdguard, says d2/graph "not done yet"
**File:** `docs/adr/001-multi-module-workspace.md`
**Problem:** Lists 7 modules (now 10), references deleted `cmdguard/`, says "d2/ and graph/ modules not yet extracted (future work)" which is now done.
**Fix:** Update module table to 10, remove cmdguard, update status, add note about completion.

### 8. ADR 002: Status shows PROPOSED but is fully implemented
**File:** `docs/adr/002-shape-capability-matrix.md:4`
**Problem:** `Status: PROPOSED` but Shape capability matrix is fully implemented and in production use.
**Fix:** Change to `Status: ACCEPTED & IMPLEMENTED`.

### 9. ADR 003: Missing — no ADR for d2/graph extraction
**File:** `docs/adr/003-d2-graph-extraction.md` (does not exist)
**Problem:** No ADR documenting the decision to extract d2/ and graph/ as independent modules.
**Fix:** Write ADR 003 with context (root was god-package at 4345 LOC), decision (extract to sub-modules), and consequences (UnsupportedFormatError, accessor methods).

### 10. DEPENDENCY_GRAPH.md: Wrong root LOC
**File:** `docs/modularization/DEPENDENCY_GRAPH.md:19`
**Problem:** Shows "~1,400 production LOC" for root but actual is 5,191. Also shows `internal/testutils/` in root box (deleted).
**Fix:** Update LOC to actual count, remove testutils reference.

### 11. DOMAIN_LANGUAGE.md: Still a template
**File:** `docs/DOMAIN_LANGUAGE.md`
**Problem:** Contains placeholder comments and "Example Term" entries. Should define actual domain vocabulary: TableData, Renderer, Format, Shape, GraphNode, GraphEdge, BrandedID, etc.
**Fix:** Populate with real domain terms from the codebase.

### 12. FORMAT_ARCHITECTURE.md: Stale registry API reference
**File:** `docs/FORMAT_ARCHITECTURE.md:170`
**Problem:** References `GetRenderer(format Format) (Renderer, error)` but actual API is `Create(format Format) (Renderer, error)`.
**Fix:** Update method name.

### 13. README.md: Supported formats table missing module annotations
**File:** `README.md:52-66`
**Problem:** `d2`, `mermaid`, `dot` rows don't note they're in separate modules (like `table` does).
**Fix:** Add "(separate `d2/` module)" and "(separate `graph/` module)" notes.

---

## P2: Code Quality & Coverage

### 14. Root test coverage: 82.2% (target: 90%+)
**File:** `internal/gentest/` (0% coverage)
**Problem:** Root module coverage is 82.2%, below the 90% target in AGENTS.md. Main gap: `internal/gentest` at 0%.
**Fix:** Either add tests for gentest or exclude from coverage reporting.

### 15. testhelpers coverage: 75% (target: 90%+)
**File:** `testhelpers/`
**Problem:** Only module below 90% target.
**Fix:** Add tests for uncovered assertions.

### 16. Add D2 benchmark tests
**File:** `d2/` (no bench_test.go)
**Problem:** `graph/` has benchmarks but `d2/` doesn't. No performance regression detection for D2 rendering.
**Fix:** Add `d2/bench_test.go` matching `graph/bench_test.go` pattern.

### 17. Add fuzz tests for d2/ and graph/
**File:** `d2/`, `graph/`
**Problem:** Root has `fuzz_test.go` but sub-modules don't.
**Fix:** Add fuzz targets for renderers in both modules.

---

## P3: Architecture

### 18. GraphRendererMixin couples Graph types with TableData
**File:** `graph.go:169-270`
**Problem:** `SetNodesFromTableData`, `AddRowEdges`, `NodesFromTableData` mix graph and table concerns in one file.
**Fix:** Extract `GraphTableConverter` or move TableData methods to a separate file.

### 19. Inconsistent re-export pattern between d2/ and graph/
**File:** `d2/d2.go:5-6` vs `graph/`
**Problem:** d2 re-exports `D2NodeID`/`D2NodeLabel` as type aliases. graph does not re-export `GraphNodeID`/`GraphNodeLabel`.
**Fix:** Pick one pattern. Either add graph aliases or remove d2 aliases.

### 20. Decide: migrate internal/gentest to testhelpers/gentest?
**File:** `internal/gentest/`
**Problem:** Sub-modules (d2, graph) cannot import `internal/` packages from root. They had to duplicate test helpers. Trade-off: public test API vs sub-module reuse.
**Status:** ❓ NEEDS DECISION (from Lars)

### 21. Duplicated test helpers in graph/helpers_test.go
**File:** `graph/helpers_test.go`
**Problem:** Copies of helpers from root's `output_test_helpers.go`. If root's versions change, graph's won't.
**Fix:** Move shared helpers to `testhelpers/` package (depends on #20 decision).

### 22. Document registry + sub-module pattern
**File:** `AGENTS.md`
**Problem:** Registry in root can't register d2/graph renderers. Users must import sub-modules directly. No guidance in docs.
**Fix:** Add "Using Registry with Sub-Modules" to AGENTS.md Common Tasks.

---

## P4: Build & Config Hygiene

### 23. .golangci.yml depguard: missing d2/graph in allow-lists
**File:** `.golangci.yml:121-192`
**Problem:** depguard `default` and `examples` rules don't mention `go-output/d2` or `go-output/graph`. Works per-module but should be explicit.
**Fix:** Add `github.com/larsartmann/go-output/d2` and `github.com/larsartmann/go-output/graph` to depguard allow-lists.

### 24. Fix or disable broken pre-commit hooks
**Problem:** `go-structure-linter` and `todo-check` hooks fail on pre-existing issues, forcing `--no-verify` for every commit.
**Fix:** Either fix the 29 root-package-files issues, add exceptions, or disable these hooks.

### 25. Verify go mod tidy is idempotent across all 10 modules
**Problem:** Not verified that `go mod tidy` produces no changes on already-tidy modules.
**Fix:** Run `go mod tidy` in all 10 modules and verify zero diff.

### 26. flake.nix: verify d2/graph included in devShell
**File:** `flake.nix`
**Problem:** If flake builds/tests modules, verify d2 and graph are included.
**Fix:** Check flake.nix and add if missing.

---

## P5: Polish & DX

### 27. Add doc comments to graph/ public API
**File:** `graph/dot.go`, `graph/mermaid.go`
**Problem:** `NewDOTRenderer`, `NewMermaidRenderer`, `DOTFromTableData`, `MermaidFlowchartRenderer` lack doc comments.
**Fix:** Add godoc-compatible comments.

### 28. Add API stability section to README
**Problem:** Pre-v1 library with no documented stability promises.
**Fix:** Add section about pre-v1 API guarantees and import path stability.

### 29. Add Example test functions for godoc
**File:** `graph/`, `d2/`
**Problem:** No `func Example*` test functions for godoc discoverability.
**Fix:** Add `ExampleDOTFromTableData`, `ExampleMermaidFlowchartRenderer`, `ExampleNewD2Diagram`.

### 30. Delete stale docs/status/ reports
**Problem:** Multiple older status reports reference broken/incomplete states that are now resolved.
**Fix:** Keep the latest, prune older reports that are fully superseded.

---

## P6: Future (Not Blocking)

### 31. Tag next release (v0.5.0?)
**Problem:** Multiple major features since v0.4.0: d2/graph extraction, Shape matrix, JSON/YAML renderers, deduplication.
**Fix:** Update CHANGELOG, tag release.

### 32. Remove deprecated FormatCategory code
**File:** `format_deprecated.go`, `format.go:171-235`
**Problem:** `FormatCategory`, `IsTableFormat()`, `Category()` all deprecated. Redirect to `Supports()`.
**Fix:** Remove in next major version.

### 33. Remove deprecated OutputFormat aliases
**File:** `format_deprecated.go`
**Problem:** `OutputFormat` type alias and all `OutputFormat*` constants.
**Fix:** Remove in v2.0.

### 34. ADR 002 Phase 2: Shape-specific renderer constructors
**File:** `docs/adr/002-shape-capability-matrix.md:104-111`
**Problem:** `NewJSONTableRenderer(data)`, `NewYAMLTreeRenderer(root)` etc. were out of scope for initial Shape refactor.
**Fix:** Implement shape-specific constructors in future iteration.

### 35. Add TOML format
**Problem:** Listed as P2 item in earlier planning.
**Fix:** New module `toml/` with table rendering.

### 36. Add JSONL format
**Problem:** Listed as P2 item in earlier planning.
**Fix:** New renderer for JSON Lines output.

### 37. Add PlantUML format
**Problem:** Listed as P3 item in earlier planning.
**Fix:** New module `plantuml/` with UML rendering.

### 38. Add AsciiDoc format
**Problem:** Listed as P3 item in earlier planning.
**Fix:** New renderer for AsciiDoc tables.

### 39. Pre-v1 API stability audit
**Problem:** No comprehensive review of all public APIs for breaking changes.
**Fix:** Audit all exported types, functions, constants. Document stability guarantees.

### 40. Community: Post to r/golang, submit to Awesome Go
**Fix:** Marketing/community items from earlier planning.

---

## Summary

| Priority | Count | Not Done | Needs Decision |
|----------|-------|----------|----------------|
| P0: Before Merge | 6 | 6 | 0 |
| P1: Documentation | 7 | 7 | 0 |
| P2: Code Quality | 4 | 4 | 0 |
| P3: Architecture | 5 | 4 | 1 |
| P4: Build/Config | 4 | 4 | 0 |
| P5: Polish | 4 | 4 | 0 |
| P6: Future | 10 | 10 | 0 |
| **Total** | **40** | **39** | **1** |

## Completed (do not re-do)

- ✅ D2 module extraction from root
- ✅ Graph module extraction from root
- ✅ Root module zero sub-module imports
- ✅ All 10 modules build/test/vet/lint clean
- ✅ AGENTS.md updated to 10-module table
- ✅ DEPENDENCY_GRAPH.md rewritten for current state
- ✅ PROPOSAL.md → ACCEPTED & IMPLEMENTED
- ✅ EXECUTION_PLAN.md → COMPLETED
- ✅ Dead code removed from root (test helpers, benchmarks)
- ✅ gci/goimports formatter conflict resolved
- ✅ Code deduplication (0 clone groups)
- ✅ Sort dependency cleanup from root go.mod
- ✅ Shape capability matrix implemented
- ✅ JSON/YAML table renderers implemented
- ✅ Graph benchmarks added to graph/ module
- ✅ UnsupportedFormatError message cleaned
