# go-output — Comprehensive Status Report

**Date:** 2026-04-05 07:23
**Branch:** master
**Last commit:** `9224cc1` docs(status): add comprehensive project audit and status report
**Go version:** 1.26.1 (darwin/arm64)
**Total lines:** ~6,566 (root .go files), ~8,500+ including subpackages

---

## Executive Summary

go-output is a **production-ready Go library** for CLI output formatting across 12 formats. All core functionality is implemented and tested. This session fixed a critical fuzz-discovered index-out-of-range panic in `MarkdownTable.writeRows`. The project has strong coverage (91.1% root, 100% cmdguard/table, 94.6% sort), 185 test functions, 13 benchmarks, and 3 fuzz targets. Remaining work is primarily CI/CD fixes, minor bugs, and polish.

---

## a) FULLY DONE ✅

| Component                | File(s)                  | Coverage  | Status                                              |
| ------------------------ | ------------------------ | --------- | --------------------------------------------------- |
| OutputFormat enum        | `format.go`              | 91.0%     | 12 formats, Parse/Validate/Categories, compat alias |
| SortBy enum              | `sort.go`                | —         | 6 sort fields, full enum pattern                    |
| ColorMode enum           | `color.go`               | —         | Auto/Always/Never with CI env detection             |
| JSON formatter           | `json.go`                | GOOD      | Marshal/Unmarshal + JSONWriter streaming            |
| CSV formatter            | `csv.go`                 | GOOD      | WriteHeader/WriteRow/WriteRows/Flush/Error          |
| YAML formatter           | `yaml.go`                | ADEQUATE  | Marshal/Unmarshal (no streaming writer)             |
| Markdown formatter       | `markdown.go`            | FIXED ✅  | Table builder with alignment — fuzz panic FIXED     |
| D2 diagram builder       | `d2.go`                  | GOOD      | Nodes, edges, shapes, styles, arrows, D2FromTable   |
| DOT graph builder        | `dot.go`                 | GOOD      | Directed/undirected, DOTFromTable, DOTFromTree      |
| Mermaid diagram          | `mermaid.go`             | GOOD      | Flowchart + tree renderers, 8 shapes                |
| HTML renderer            | `html.go`                | GOOD      | Table + tree renderers, full document mode           |
| XML formatter            | `xml.go`                 | GOOD      | Marshal/Unmarshal + XMLWriter streaming              |
| TSV formatter            | `tsv.go`                 | ADEQUATE  | Writer + MarshalTSV                                  |
| ASCII Tree renderer      | `tree.go`                | GOOD      | Unicode box-drawing, TreeRendererFromTableData       |
| Graph types              | `graph.go`               | GOOD      | GraphNode, GraphEdge, 8 shapes, GraphRenderer       |
| Registry                 | `registry.go`            | GOOD      | Thread-safe factory registry with mutex              |
| Streaming framework      | `streaming.go`           | GOOD      | StreamingRenderer interface + HTML impl + adapter    |
| Branded IDs              | `ids.go`                 | —         | Phantom-typed IDs for compile-time safety            |
| Marshal helpers          | `marshal.go`             | —         | Generic marshal/unmarshal with format-aware errors   |
| Shared markup            | `markup.go`              | —         | HTML/XML shared escaping via internal/escape         |
| Table renderer (pkg)     | `table/table.go`         | 100.0%    | lipgloss v2 styled tables with StyleFunc             |
| Sort utilities (pkg)     | `sort/sorter.go`         | 94.6%     | Generic Sorter with string/int/time comparators      |
| Enum package (pkg)       | `enum/enum.go`           | 71.4%     | Generic enum validation helpers                      |
| cmdguard integration     | `cmdguard/`              | 100.0%    | Format, SortBy, ColorMode flag factories             |
| Internal escape          | `internal/escape/`       | —         | Safe string escaping                                 |
| Benchmarks               | `benchmarks_test.go`     | —         | 13 benchmarks across key renderers                   |
| Fuzz tests               | `fuzz_test.go`           | —         | 3 fuzz targets (Format, CSV, Markdown)               |
| Integration tests        | `integration/`           | —         | 4 files: all-formats, workflows, renderers, format   |
| User journey tests       | `userjourney_test.go`    | —         | 5 documented user journeys                           |
| Example                  | `examples/basic/main.go` | —         | All 12 formats demonstrated                          |
| **Total test functions** |                          | **185**   | Across 31 test files                                 |

---

## b) PARTIALLY DONE 🟡

| Component              | What's Done                         | What's Missing                                        |
| ---------------------- | ----------------------------------- | ----------------------------------------------------- |
| Streaming              | Framework + HTML streaming works    | Only HTML has true streaming; all others use buffered  |
| TTY detection          | `isTerminal()` exists               | `isStderrTerminal()` hardcoded `false`; incomplete     |
| Fuzz tests             | 3 targets (Format, CSV, Markdown)   | Missing fuzz for YAML, XML, HTML, D2, DOT, Mermaid    |
| Benchmarks             | 13 functions across key formats     | Missing benchmarks for D2, DOT, Mermaid, Tree, YAML   |
| Enum coverage          | Core enum package works             | `enum/` only 71.4% coverage                           |
| CI/CD                  | Basic workflows exist               | Release workflow broken, Go version mismatch          |

---

## c) NOT STARTED ⬜

| # | Task                                                                                           | Priority   |
| - | ---------------------------------------------------------------------------------------------- | ---------- |
| 1 | `.goreleaser.yml` — proper release automation                                                  | HIGH       |
| 2 | CONTRIBUTING.md — contributor guidelines                                                       | MEDIUM     |
| 3 | CI matrix — test across multiple Go versions (1.24, 1.25, 1.26)                                | MEDIUM     |
| 4 | Streaming writers for YAML, CSV, TSV, Markdown, XML, D2, DOT, Mermaid                         | MEDIUM     |
| 5 | Expand fuzz targets to all formats                                                             | MEDIUM     |
| 6 | Expand benchmark coverage to all renderers                                                     | MEDIUM     |
| 7 | `pkg/errors` — structured error types (currently placeholder)                                  | LOW        |
| 8 | godoc examples (`Example*` functions) for public APIs                                          | LOW        |

---

## d) TOTALLY FUCKED UP 🔴

| # | Issue                                                                                       | Impact       | Effort   |
| - | ------------------------------------------------------------------------------------------- | ------------ | -------- |
| 1 | **Markdown `writeRows` index-out-of-range** — rows with more cells than headers panicked   | **FIXED** ✅ | 5min     |
| 2 | **Markdown `AlignCenter` bug** — renders identically to `AlignLeft` (uses `%-*s` left-pad) | Formatting   | 15min    |
| 3 | **Release workflow** (`release.yml`) — tries to re-tag existing tag, will fail              | Deployment   | 30min    |
| 4 | **CI Go version mismatch** — CI uses 1.23, `go.mod` says 1.26.1                             | CI broken    | 15min    |
| 5 | **golangci-lint config broken** — `linters-settings` placement wrong, rules ignored         | Linting      | 10min    |
| 6 | **README leaks local paths** — `/Users/larsartmann/projects/...` in documentation           | Privacy      | 5min     |
| 7 | **Example uses deprecated API** — `ParseOutputFormat` instead of `ParseFormat`              | API          | 10min    |

---

## e) WHAT WE SHOULD IMPROVE 📈

1. **Fuzz-first development** — Every public formatter should have a fuzz target. The Markdown panic proves fuzzing catches real bugs that unit tests miss.
2. **Bounds discipline** — Every slice/array access that depends on external input needs bounds checking. The `writeRows` bug is a textbook case.
3. **Test coverage for enum package** — 71.4% is the lowest in the project. Easy wins available.
4. **Dead code removal** — ~11 unused branded ID types in `ids.go` add noise.
5. **Linter configuration** — A broken linter gives false confidence. Fix it, pin the version.
6. **CI/CD alignment** — Go version mismatch between `go.mod`, CI workflows, and toolchain is a ticking time bomb.
7. **Streaming completeness** — Only HTML truly streams. The adapter pattern works but is misleading for consumers expecting streaming behavior.
8. **Documentation consistency** — README, examples, and code should all use the same (current) API surface.
9. **Benchmark coverage** — Missing benchmarks for D2, DOT, Mermaid, Tree, YAML means performance regressions go undetected.
10. **Go version toolchain mismatch** — `go version` is 1.26.1 but `-cover` emits toolchain version mismatch warnings for internal packages. This is cosmetic but messy.

---

## f) TOP 25 THINGS TO DO NEXT

| Priority | # | Task                                                                    | Effort    | Status          |
| -------- | - | ----------------------------------------------------------------------- | --------- | --------------- |
| 🔴 P0    | 1 | ~~Fix Markdown `writeRows` panic~~                                      | ~~5min~~  | **DONE** ✅      |
| 🔴 P0    | 2 | Fix Markdown `AlignCenter` rendering bug                                | 15min     | NOT STARTED     |
| 🔴 P0    | 3 | Fix release workflow (`release.yml`)                                    | 30min     | NOT STARTED     |
| 🔴 P0    | 4 | Align Go versions across CI, go.mod, toolchain                         | 15min     | NOT STARTED     |
| 🔴 P1    | 5 | Fix golangci-lint config                                                | 10min     | NOT STARTED     |
| 🟡 P1    | 6 | Scrub README local filesystem paths                                     | 5min      | NOT STARTED     |
| 🟡 P1    | 7 | Fix example deprecated API usage                                        | 10min     | NOT STARTED     |
| 🟡 P1    | 8 | Add `.goreleaser.yml`                                                   | 30min     | NOT STARTED     |
| 🟡 P1    | 9 | Complete TTY detection (`isStderrTerminal`, `isTerminal`)              | 30min     | PARTIAL         |
| 🟡 P2    | 10 | Clean up dead branded IDs in `ids.go`                                   | 20min     | NOT STARTED     |
| 🟡 P2    | 11 | Pin golangci-lint version in CI                                         | 10min     | NOT STARTED     |
| 🟡 P2    | 12 | Add CONTRIBUTING.md                                                     | 30min     | NOT STARTED     |
| 🟡 P2    | 13 | Increase enum package test coverage to 90%+                            | 30min     | NOT STARTED     |
| 🟢 P2    | 14 | Add fuzz targets for YAML, XML, HTML                                    | 30min     | NOT STARTED     |
| 🟢 P2    | 15 | Add fuzz targets for D2, DOT, Mermaid                                   | 30min     | NOT STARTED     |
| 🟢 P2    | 16 | Add benchmarks for D2, DOT, Mermaid, Tree, YAML                        | 30min     | NOT STARTED     |
| 🟢 P3    | 17 | Add CI matrix for Go 1.24, 1.25, 1.26                                  | 20min     | NOT STARTED     |
| 🟢 P3    | 18 | Add godoc `Example*` functions for key public APIs                     | 60min     | NOT STARTED     |
| 🟢 P3    | 19 | Streaming writers for remaining formats                                 | 2-4h      | NOT STARTED     |
| 🟢 P3    | 20 | Update CHANGELOG.md for next release                                    | 20min     | NOT STARTED     |
| 🟢 P3    | 21 | Review and update PLAN.md                                               | 20min     | NOT STARTED     |
| ⚪ P4    | 22 | Flesh out `pkg/errors` with structured error types                      | 1-2h      | NOT STARTED     |
| ⚪ P4    | 23 | Add YAML assertions in tests (structured comparison)                   | 30min     | NOT STARTED     |
| ⚪ P4    | 24 | Add integration test for cmdguard flag parsing with real cobra/command  | 30min     | NOT STARTED     |
| ⚪ P4    | 25 | Explore GoReleaser vs pure library versioning strategy                 | Research  | BLOCKED         |

---

## g) TOP #1 QUESTION

**What is the intended release and versioning strategy?**

This project is a Go library (no `main` package, no binary). Yet it has a `release.yml` GitHub Actions workflow suggesting binary releases. This affects:

- Whether we need `.goreleaser.yml` at all (libraries just need tagged commits)
- Whether the release workflow should build binaries or just run tests + tag
- Whether `pkg/errors`, `internal/gentest`, and `examples/basic` should be separate modules
- Whether we should publish to a Go proxy as-is or restructure

The entire CI/CD architecture depends on this answer. Current state suggests confusion between library and application packaging.

---

## Session Work (2026-04-05 07:20–07:23)

### Fixed: Markdown Table Fuzz Panic

**Bug:** `FuzzMarkdownTable` discovered a `panic: runtime error: index out of range [1] with length 1` in `markdown.go:126` (`writeCell`).

**Root cause:** `writeRows` iterated over row cells (`for i, cell := range row`), but `colWidths` is sized by headers. When a row had more cells than headers, `colWidths[i]` panicked.

**Fix:** Changed `writeRows` to iterate over `m.headers` (the authoritative column count), using empty string for missing cells and silently dropping extras beyond the header count.

**Verification:**
- Fuzz tested for 30 seconds: 613,412 executions, 0 panics
- All 185 unit tests pass
- Full test suite: all packages green
- Coverage: root 91.1%, cmdguard 100%, table 100%, sort 94.6%

---

## Coverage Summary

| Package                     | Coverage |
| --------------------------- | -------- |
| `github.com/larsartmann/go-output` (root) | 91.1%    |
| `github.com/larsartmann/go-output/cmdguard` | 100.0%   |
| `github.com/larsartmann/go-output/enum`    | 71.4%    |
| `github.com/larsartmann/go-output/integration` | — (integration) |
| `github.com/larsartmann/go-output/sort`    | 94.6%    |
| `github.com/larsartmann/go-output/table`   | 100.0%   |

## Test Inventory

- **185 test functions** across 31 test files
- **13 benchmark functions** (7 root + 3 JSON + 2 YAML + 1 TSV)
- **3 fuzz targets** (OutputFormat, CSVWriter, MarkdownTable)
- **5 user journey tests**
- **4 integration test files**
