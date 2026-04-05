# go-output — Comprehensive Status Report

**Date:** 2026-04-05 07:08  
**Branch:** master  
**Last commit:** `8390880` feat: comprehensive test framework and benchmark infrastructure refactor  
**Go version:** 1.26.1 (darwin/arm64)  
**Total lines:** 8,067 (Go source + test files)

---

## Executive Summary

go-output is a **production-ready Go library** for CLI output formatting across 12 formats. All core functionality is implemented, tested, and passing. The project has strong test coverage (91%+ root, 100% cmdguard/table, 94.6% sort), comprehensive integration tests, benchmarks, and fuzz targets. The main risks are CI/CD configuration bugs, minor code defects, and documentation inconsistencies — not missing functionality.

---

## a) FULLY DONE ✅

| Component                | File(s)                  | Coverage | Notes                                                          |
| ------------------------ | ------------------------ | -------- | -------------------------------------------------------------- |
| **OutputFormat enum**    | `format.go`              | 91.0%    | 12 formats, Parse/Validate/Categories, backward compat aliases |
| **SortBy enum**          | `sort.go`                | —        | 6 sort fields, full enum pattern                               |
| **ColorMode enum**       | `color.go`               | —        | Auto/Always/Never with CI env detection                        |
| **JSON formatter**       | `json.go`                | GOOD     | Marshal/Unmarshal + JSONWriter streaming                       |
| **CSV formatter**        | `csv.go`                 | GOOD     | WriteHeader/WriteRow/WriteRows/Flush/Error                     |
| **YAML formatter**       | `yaml.go`                | ADEQUATE | Marshal/Unmarshal (no streaming writer)                        |
| **Markdown formatter**   | `markdown.go`            | ADEQUATE | Table builder with alignment (center has bug)                  |
| **D2 diagram builder**   | `d2.go`                  | GOOD     | Nodes, edges, shapes, styles, arrows, D2FromTableData          |
| **DOT graph builder**    | `dot.go`                 | GOOD     | Directed/undirected, DOTFromTableData, DOTFromTree             |
| **Mermaid diagram**      | `mermaid.go`             | GOOD     | Flowchart + tree renderers, 8 shapes                           |
| **HTML renderer**        | `html.go`                | GOOD     | Table + tree renderers, full document mode                     |
| **XML formatter**        | `xml.go`                 | GOOD     | Marshal/Unmarshal + XMLWriter streaming                        |
| **TSV formatter**        | `tsv.go`                 | ADEQUATE | Writer + MarshalTSV                                            |
| **ASCII Tree renderer**  | `tree.go`                | GOOD     | Unicode box-drawing, TreeRendererFromTableData                 |
| **Graph types**          | `graph.go`               | GOOD     | GraphNode, GraphEdge, 8 shapes, GraphRenderer interface        |
| **Registry**             | `registry.go`            | GOOD     | Thread-safe factory registry with mutex                        |
| **Streaming framework**  | `streaming.go`           | GOOD     | StreamingRenderer interface + HTML impl + adapter              |
| **Branded IDs**          | `ids.go`                 | —        | Phantom-typed IDs for compile-time safety                      |
| **Marshal helpers**      | `marshal.go`             | —        | Generic marshal/unmarshal with format-aware errors             |
| **Shared markup**        | `markup.go`              | —        | HTML/XML shared escaping via internal/escape                   |
| **Table renderer**       | `table/table.go`         | 100.0%   | lipgloss v2 styled tables with StyleFunc                       |
| **Sort utilities**       | `sort/sorter.go`         | 94.6%    | Generic Sorter with string/int/time comparators                |
| **Enum package**         | `enum/enum.go`           | 71.4%    | Generic enum validation helpers                                |
| **cmdguard integration** | `cmdguard/`              | 100.0%   | Format, SortBy, ColorMode flag factories                       |
| **Internal escape**      | `internal/escape/`       | —        | Safe string escaping                                           |
| **Benchmarks**           | `benchmarks_test.go`     | —        | 7 benchmarks across key renderers                              |
| **Fuzz tests**           | `fuzz_test.go`           | —        | 3 fuzz targets (Format, CSV, Markdown)                         |
| **Integration tests**    | `integration/`           | —        | 4 files: all-formats, workflows, renderers, format parse       |
| **User journey tests**   | `userjourney_test.go`    | —        | 5 documented user journeys                                     |
| **Example**              | `examples/basic/main.go` | —        | All 12 formats demonstrated                                    |

### Test Coverage Summary

| Package            | Coverage   |
| ------------------ | ---------- |
| `go-output` (root) | **91.0%**  |
| `cmdguard`         | **100.0%** |
| `table`            | **100.0%** |
| `sort`             | **94.6%**  |
| `enum`             | **71.4%**  |

### All Tests: PASSING ✅

---

## b) PARTIALLY DONE ⚠️

| Component                     | Status                           | What's Missing                                                                                                                                                                      |
| ----------------------------- | -------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Streaming support**         | HTML-only true streaming         | JSON, CSV, YAML, XML, etc. use buffered adapter — no incremental output for large datasets                                                                                          |
| **Markdown center alignment** | Implemented but **buggy**        | `AlignCenter` uses `%-*s` (left-align format) — doesn't actually center text                                                                                                        |
| **ColorMode auto-detection**  | Partial implementation           | `isStderrTerminal()` hardcoded to `false`; `isTerminal()` only checks env vars, not actual TTY via `os.Stat()` or `golang.org/x/term`                                               |
| **Branded IDs (ids.go)**      | 13 types defined, ~11 never used | `D2EdgeFromBrand`, `DOTNodeIDBrand`, `DOTEdgeFromBrand`, `DOTEdgeToBrand`, `DOTEdgeLabelBrand`, `TreeParentID`, `MermaidNodeID`, `MermaidParentID`, `HTMLTitle`, etc. are dead code |
| **Fuzz testing**              | 3 targets only                   | Missing: JSON, YAML, XML, HTML, DOT, Mermaid, D2, TSV                                                                                                                               |
| **Benchmarks**                | 7 renderers covered              | Missing: JSON, YAML, XML, TSV, D2                                                                                                                                                   |
| **YAML tests**                | No exact output assertions       | Tests only verify "not empty" — no structural validation                                                                                                                            |
| **Markdown tests**            | No exact output verification     | Tests check rendering runs, not output correctness                                                                                                                                  |
| **enum package coverage**     | 71.4%                            | Missing fuzz test, missing `IsValid` direct test                                                                                                                                    |

---

## c) NOT STARTED ❌

| Component                                 | Priority | Notes                                                                  |
| ----------------------------------------- | -------- | ---------------------------------------------------------------------- |
| **CONTRIBUTING.md**                       | Medium   | No contribution guidelines exist                                       |
| **Go module caching in CI**               | Medium   | No `setup-go` cache or `actions/cache`                                 |
| **CI matrix testing**                     | Low      | Only ubuntu-latest; no macOS/Windows, no Go version matrix             |
| **GoReleaser config**                     | High     | Release workflow references goreleaser but no `.goreleaser.yml` exists |
| **True streaming for non-HTML formats**   | Low      | All formats except HTML use buffered adapter                           |
| **TTY detection via `golang.org/x/term`** | Medium   | `isTerminal()` and `isStderrTerminal()` are incomplete stubs           |
| **TSV error writer test**                 | Low      | CSV has it, TSV doesn't                                                |
| **Table column mismatch test**            | Low      | No test for rows with more/fewer values than headers                   |

---

## d) TOTALLY FUCKED UP 💥

| Issue                                    | File                                  | Severity    | Details                                                                                                                                  |
| ---------------------------------------- | ------------------------------------- | ----------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| **Release workflow will FAIL**           | `.github/workflows/release.yml:57-63` | 🔴 CRITICAL | Tries to `git tag` + `git push` for a tag that already exists (workflow is triggered by that tag). Will error with "tag already exists." |
| **CI Go version mismatch**               | `ci.yml`, `release.yml`               | 🔴 HIGH     | CI uses Go 1.23 but `go.mod` says 1.26.1, `.golangci.yml` says 1.26.1. Builds may fail or produce inconsistent results.                  |
| **golangci-lint config silently broken** | `.golangci.yml:126`                   | ⚠️ MEDIUM   | `settings` key at `linters:` level instead of `linters-settings:` — revive rule disabling `var-naming` is **silently ignored**.          |
| **README leaks local filesystem paths**  | `README.md:12-15`                     | ⚠️ MEDIUM   | Exposes `/Users/larsartmann/projects/...` — personal/private detail in public README.                                                    |
| **Example uses deprecated API**          | `examples/basic/main.go:52`           | ⚠️ MEDIUM   | Uses `ParseOutputFormat` which is deprecated. Should use `ParseFormat`.                                                                  |
| **Markdown center alignment bug**        | `markdown.go:133`                     | 🐛 BUG      | `AlignCenter` renders identically to `AlignLeft` — format string `%-*s` is left-align.                                                   |
| **Go toolchain version mismatch**        | `go.mod` vs installed                 | ⚠️ MEDIUM   | `go.mod` says `go 1.26.0` but installed Go is `1.26.1` — causes spurious `compile: version does not match` warnings in coverage builds.  |

---

## e) WHAT WE SHOULD IMPROVE 📈

### High Priority

1. **Fix release workflow** — The `git tag` step will fail; re-architect release process
2. **Align Go versions** — `go.mod`, `.golangci.yml`, `ci.yml`, `release.yml` all differ
3. **Fix golangci-lint config** — `linters-settings` placement is wrong; revive config is ignored
4. **Fix markdown center alignment bug** — `%-*s` → proper centering calculation
5. **Scrub README** — Remove local filesystem paths from public documentation
6. **Fix example to use non-deprecated API** — `ParseOutputFormat` → `ParseFormat`

### Medium Priority

7. **Complete color auto-detection** — Implement real TTY check for stdout/stderr
8. **Add `.goreleaser.yml`** — Release workflow depends on it but it doesn't exist
9. **Pin `golangci-lint` version in CI** — `@latest` can break CI unexpectedly
10. **Clean up dead branded ID types** — Remove ~11 unused phantom types from `ids.go`
11. **Add CONTRIBUTING.md** — Document contribution process
12. **Fix `isStderrTerminal()`** — Currently hardcoded to `false`

### Low Priority

13. **Add YAML exact output assertions** — Test structural correctness
14. **Expand fuzz targets** — Add for JSON, YAML, XML, HTML, DOT, Mermaid, D2, TSV
15. **Expand benchmarks** — Add for JSON, YAML, XML, TSV, D2
16. **Add CI matrix** — Test on macOS, Windows, multiple Go versions
17. **Add Go module caching to CI** — Speed up CI runs
18. **Fix enum package coverage** — Bring from 71.4% to 90%+
19. **Move PLAN.md to docs/** — All phases complete; historical artifact
20. **Fix pre-commit `--fix` auto-modification** — Can cause confusing commit behavior

---

## f) Top #25 Things We Should Get Done Next

| #   | Task                                                                       | Effort | Impact      | Type        |
| --- | -------------------------------------------------------------------------- | ------ | ----------- | ----------- |
| 1   | Fix release workflow (`release.yml` re-tagging bug)                        | 30min  | 🔴 Critical | Bug fix     |
| 2   | Align Go versions across all config files (1.26.1)                         | 15min  | 🔴 High     | Consistency |
| 3   | Fix golangci-lint config (`linters-settings` placement)                    | 10min  | 🔴 High     | Bug fix     |
| 4   | Fix markdown `AlignCenter` bug                                             | 15min  | 🐛 Bug      | Bug fix     |
| 5   | Scrub README (remove local paths)                                          | 5min   | ⚠️ Medium   | Security    |
| 6   | Fix example to use `ParseFormat` instead of deprecated `ParseOutputFormat` | 10min  | ⚠️ Medium   | Cleanup     |
| 7   | Add `.goreleaser.yml` configuration                                        | 30min  | ⚠️ Medium   | Infra       |
| 8   | Complete `isTerminal()` with real TTY detection                            | 30min  | ⚠️ Medium   | Feature     |
| 9   | Fix `isStderrTerminal()` (currently hardcoded `false`)                     | 15min  | ⚠️ Medium   | Bug fix     |
| 10  | Clean up unused branded ID types in `ids.go`                               | 20min  | Low         | Dead code   |
| 11  | Pin `golangci-lint` version in CI workflows                                | 5min   | ⚠️ Medium   | Infra       |
| 12  | Add CONTRIBUTING.md                                                        | 20min  | Low         | Docs        |
| 13  | Add YAML exact output assertions in tests                                  | 20min  | Low         | Quality     |
| 14  | Add fuzz targets for JSON, YAML, XML, HTML, D2, DOT, Mermaid, TSV          | 1h     | Low         | Quality     |
| 15  | Add benchmarks for JSON, YAML, XML, TSV, D2                                | 30min  | Low         | Perf        |
| 16  | Add Go module caching to CI (`setup-go` cache mode)                        | 10min  | Low         | Infra       |
| 17  | Improve enum package coverage (71.4% → 90%+)                               | 20min  | Low         | Quality     |
| 18  | Update CHANGELOG — cut a new version release                               | 15min  | Low         | Release     |
| 19  | Add TSV error writer test                                                  | 10min  | Low         | Quality     |
| 20  | Add table column mismatch test                                             | 10min  | Low         | Quality     |
| 21  | Move PLAN.md to `docs/` (all phases complete)                              | 5min   | Low         | Cleanup     |
| 22  | Fix `go.mod` toolchain directive (1.26.0 → 1.26.1)                         | 5min   | Low         | Consistency |
| 23  | Add CI matrix testing (macOS, Windows)                                     | 20min  | Low         | Infra       |
| 24  | Remove pre-commit `--fix` flag from golangci-lint                          | 5min   | Low         | DX          |
| 25  | Add true streaming support for JSON/CSV/XML formats                        | 2h     | Low         | Feature     |

---

## g) Top #1 Question I Cannot Figure Out Myself 🤔

**What is the intended release/versioning strategy for this library?**

The project has:

- A `CHANGELOG.md` with `[Unreleased]` accumulating entries + a `[0.1.0]` dated 2026-01-01
- A `release.yml` workflow triggered by `v*` tags with GoReleaser
- But **no `.goreleaser.yml`** configuration file
- No semver version in `go.mod` (just the module path)
- The release workflow has a broken re-tagging step

**I need to know:** Is the intent to use GoReleaser for proper binary releases? Or is this a pure Go library (no binaries) where GoReleaser doesn't apply? Should we use GoReleaser's `gomod` mode for module proxying, or simplify the release workflow to just create GitHub releases with the changelog? This determines whether we fix the GoReleaser setup or replace it with something simpler.

---

## File Inventory

### Source Files (Production Code)

| File             | Lines | Purpose                                     |
| ---------------- | ----- | ------------------------------------------- |
| `format.go`      | 302   | OutputFormat enum + data types + interfaces |
| `d2.go`          | 302   | D2 diagram builder                          |
| `dot.go`         | 258   | DOT/Graphviz renderer                       |
| `sort/sorter.go` | 238   | Generic sorting utilities                   |
| `html.go`        | 205   | HTML table + tree renderers                 |
| `mermaid.go`     | 186   | Mermaid diagram renderer                    |
| `streaming.go`   | 183   | Streaming renderer framework                |
| `tree.go`        | 123   | ASCII tree renderer                         |
| `ids.go`         | 156   | Branded phantom-typed IDs                   |
| `graph.go`       | 129   | Core graph types                            |
| `markdown.go`    | 153   | Markdown table builder                      |
| `tsv.go`         | 109   | TSV formatter                               |
| `color.go`       | 107   | Color mode enum + auto-detection            |
| `xml.go`         | 103   | XML formatter                               |
| `registry.go`    | 85    | Thread-safe renderer registry               |
| `yaml.go`        | ~70   | YAML formatter                              |
| `json.go`        | ~60   | JSON formatter                              |
| `csv.go`         | ~60   | CSV formatter                               |
| `marshal.go`     | ~40   | Generic marshal helpers                     |
| `markup.go`      | ~30   | Shared HTML/XML helpers                     |
| `slices.go`      | ~15   | Slice utilities                             |
| `sort.go`        | ~40   | SortBy enum                                 |
| `table/table.go` | ~80   | lipgloss table wrapper                      |
| `enum/enum.go`   | ~40   | Generic enum helpers                        |
| `cmdguard/*.go`  | ~100  | CLI flag factories                          |

### Test Files

| File                        | Lines | Type               |
| --------------------------- | ----- | ------------------ |
| `sort/sort_test.go`         | 363   | Unit               |
| `d2_test.go`                | 321   | Unit               |
| `format_test.go`            | 309   | Unit + Fuzz        |
| `cmdguard/cmdguard_test.go` | 296   | Unit               |
| `mermaid_test.go`           | 283   | Unit               |
| `userjourney_test.go`       | 263   | User Journey       |
| `html_test.go`              | 199   | Unit               |
| `streaming_test.go`         | 198   | Unit               |
| `benchmarks_test.go`        | 197   | Benchmarks         |
| `json_test.go`              | 190   | Unit + Benchmarks  |
| `xml_test.go`               | 181   | Unit               |
| `table/table_test.go`       | 180   | Unit               |
| `dot_test.go`               | 179   | Unit               |
| `registry_test.go`          | 174   | Unit + Concurrency |
| `testing_test.go`           | 167   | Test Helpers       |
| `graph_test.go`             | 167   | Unit               |
| `tree_test.go`              | 156   | Unit               |
| `color_test.go`             | 136   | Unit               |
| `csv_test.go`               | 130   | Unit               |
| `markdown_test.go`          | 82    | Unit               |
| `yaml_test.go`              | 117   | Unit               |
| `tsv_test.go`               | 107   | Unit               |
| `fuzz_test.go`              | 104   | Fuzz               |
| `sort_test.go`              | 81    | Unit               |
| `integration/` (4 files)    | ~400  | Integration        |

---

## Uncommitted Changes (Working Tree)

| File             | Status    | Description                                                       |
| ---------------- | --------- | ----------------------------------------------------------------- |
| `.gitignore`     | Modified  | Added `coverage/` dir and `*.db` to ignores                       |
| `AGENTS.md`      | Modified  | Fixed lipgloss import paths (`github.com/...` → `charm.land/...`) |
| `go.mod`         | Modified  | Updated `go 1.26.0` → `go 1.26.1`, bumped indirect deps           |
| `go.sum`         | Modified  | Updated checksums for bumped deps                                 |
| `.gitattributes` | Untracked | New file: excludes `docs/status/*.md` from linguist               |

---

_Generated by Crush — 2026-04-05_
