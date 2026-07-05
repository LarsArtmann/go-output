# Round 4 Comprehensive Polish — Full Status Report

**Date:** 2026-05-28 05:14 CEST
**Branch:** master
**Head commit:** `570c923`
**Session scope:** 4 rounds of footer row feature + polish (commits `98c0ae3..570c923`)
**Files changed:** 58 (+1,333 / -271 lines across all rounds)

---

## Summary

Four rounds of work on the `go-output` library:

- **Round 1** (13 commits): Core footer row feature implementation across 6 modules
- **Round 2** (7 commits): go.mod unification, Validate() wiring, formatter fixes
- **Round 3** (1 commit): WriteFooter tests, SetFooter multi-call bug fix
- **Round 4** (10 commits): Coverage improvement, doc.go, GoDoc, lint fixes, documentation

---

## A) FULLY DONE ✅

### Footer Row Feature (Rounds 1-3)

- **`TableData.Footer []string`** — single source of truth for optional summary row
- **`TableData.Validate()`** — validates footer column count matches headers
- **`Validate()` wired into `RenderTableData()`** — automatic validation before dispatch
- **`TableDataBase`** — `SetFooter()`, `HasFooter()` for sub-module reuse
- **Markdown** — footer after separator, bold-styled, inherits column alignment
- **CSV/TSV** — `WriteFooter()` streaming method + `marshalFromTableData()` shared helper
- **HTML** — `<tfoot>` with `class="footer-cell"` on both regular and streaming renderer
- **XML** — `<footer>` element with cells
- **AsciiDoc** — footer row (uses `HasFooter()` consistently)
- **Terminal Table** — bold footer via `footerRowIndex` tracking (multi-call bug fixed)
- **Data formats** (JSON/YAML/TOML/JSONL) — intentionally skip footer
- **Non-tabular** (Tree/D2/Mermaid/DOT/PlantUML) — intentionally skip footer
- **Tests**: Footer tests in every format module + integration tests + edge cases (Validate mismatch, empty cells, nil footer, footer longer/shorter than headers)
- **Benchmarks**: Markdown, CSV, Table footer benchmarks
- **GoDoc**: `ExampleTable_SetFooter` in table module
- **Examples**: Footer in `examples/basic/main.go`

### Round 2: Architecture Fixes

- All 9 `go.mod` files unified to `go 1.26.3`
- `integration/go.mod` root dep fixed from `v0.5.0` to `v0.0.0`
- `go.work.example` includes `./testhelpers/graphtest`
- Formatter suggestions applied (`errors.New` for static strings, import grouping)
- `delimited/` dedup: `tableDataWriter` interface + `marshalFromTableData()` shared helper

### Round 3: Bug Fix

- **`table.SetFooter` multi-call bug fixed**: Added `footerRowIndex int` field to `Table` struct. Only the last footer row receives bold styling; previous rows become unstyled data (lipgloss limitation). `TestTableSetFooter_MultipleCalls` verifies correct behavior.

### Round 4: Coverage, Documentation, Polish

- **Root coverage: 88.6% → 95.9%** — tested TableDataBase (7 methods), Error types (3), Format.IsValid, NewGraphNode, StreamingRendererFromRenderer, UnmarshalFormat
- **Serialization coverage: 88.3% → 89.0%** — tested JSONL Encode error path
- **`doc.go` added to 8 packages**: output, delimited, d2, graph, markup, plantuml, serialization, testhelpers — pkg.go.dev now shows proper package docs
- **GoDoc on 50+ exported struct fields**: GraphNode, GraphStyle, GraphEdge, EdgeStyle, TreeNode, D2StrokeStyle, D2NodeStyle, D2Node, D2EdgeStyle, D2Edge, D2Column, D2Table
- **6 GoDoc examples**: Format.IsValid, ParseFormat, ColorMode, Shape, TableData.Validate, MustRender
- **`t.Parallel()` added to 13 tests** in `render_tabledata_test.go` + `integration_test.go`
- **`xml_test.go` split** (341→244+101 lines) — proactive, was 9 lines from 350 limit
- **All lint issues fixed**: err113 sentinel errors, modernize `interface{}`→`any`, wsl_v5 whitespace

### Coverage Table

| Module        | Coverage | Status            |
| ------------- | -------- | ----------------- |
| output (root) | 95.9%    | ✅ Above 90%      |
| delimited     | 90.2%    | ✅ Above 90%      |
| d2            | 100.0%   | ✅ Perfect        |
| enum          | 100.0%   | ✅ Perfect        |
| escape        | 100.0%   | ✅ Perfect        |
| graph         | 96.0%    | ✅ Above 90%      |
| markup        | 93.9%    | ✅ Above 90%      |
| plantuml      | 97.2%    | ✅ Above 90%      |
| table         | 92.2%    | ✅ Above 90%      |
| testhelpers   | 91.3%    | ✅ Above 90%      |
| integration   | 82.8%    | ⚠️ Below 90%      |
| serialization | 89.0%    | ⚠️ Just below 90% |
| gentest       | 80.8%    | ⚠️ Below 90%      |

### Build/Lint Status

- All 14 modules build clean
- All 14 modules test pass
- All 14 modules lint clean (0 golangci-lint issues)
- No files exceed 350-line limit
- No TODO/FIXME/HACK comments in codebase

---

## B) PARTIALLY DONE ⚠️

### Serialization Coverage (89.0%)

The remaining 1% gap is in marshal error paths for tree/graph renderers (JSON, YAML, TOML). These error paths can't fail in practice because the DTO types are simple structs with string fields. The `renderTable` error path and graph DTO marshal errors are theoretically unreachable.

### Integration Coverage (82.8%)

Integration tests cover the main paths but some error paths and edge cases are not tested. This is expected for cross-module integration tests.

### gentest Coverage (80.8%)

The `internal/gentest` package has some uncovered assertion helpers. This is a pre-existing gap.

---

## C) NOT STARTED ❌

1. **Serialization footer option** — JSON/YAML/TOML could include footer with `_footer: true` marker (controversial — data formats shouldn't mix visual concerns)
2. **Multiple footer rows** — Only single footer supported (`[]string`, not `[][]string`)
3. **`WithFooterStyle` option** — Composable footer styling for users who override `table.StyleFunc`
4. **MarkdownTable satisfies TableRenderer** — Method signatures differ (chaining vs void return). Would require API redesign.
5. **RenderOptions functional options** — Current struct has public fields; functional options would be cleaner but breaking change
6. **Move `internal/gentest` to `testhelpers/gentest`** — Needs user decision (TODO_LIST #20)
7. **Pre-v1 API stability audit** — TODO_LIST #39
8. **Community outreach** — Post to r/golang, submit to Awesome Go (TODO_LIST #40)
9. **ADR for footer design decision** — Document why footer is on TableData vs RenderOptions
10. **Graph FromTableData/FromTree dedup** — dot.go and mermaid.go have similar helper skeletons (~30 lines each)

---

## D) TOTALLY FUCKED UP 💥

Nothing catastrophic. Minor issues resolved during the session:

- **`table.SetFooter` multi-call bug**: Was incrementing `rowCount` and only bolding the last row. **Fixed** by tracking `footerRowIndex` separately.
- **Stale LSP diagnostics**: `markdown_footer_test.go` showed false-positive "duplicate decl" in gopls. File is valid Go, builds/tests clean. LSP needed restart.
- **`go.work` missing graphtest**: The workspace file didn't include `./testhelpers/graphtest`, causing build failures. **Fixed** by adding it to `go.work` (already in `go.work.example`).
- **TSV WriteFooter test missing `bytes` import**: Leftover from previous session. **Fixed** immediately.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **InvalidError type consistency**: `InvalidFormatError` has `Allowed []Format` field, while `InvalidShapeError`, `InvalidColorModeError`, `InvalidGraphShapeError` only have `Value string`. Design choice for now (format is user-facing, others are internal), but worth standardizing.
2. **`TableRenderer` interface gap**: `MarkdownTable` and `table.Table` don't satisfy `TableRenderer` due to chaining method signatures (`*Table` returns vs void). This means code working against the interface can't accept them. Architectural inconsistency.
3. **Serialization tree/graph renderer duplication**: 3 TreeRenderers and 3 GraphRenderers in `serialization/` with identical structure — only marshal function differs. Could be a generic `SerializationRenderer` parameterized by marshal func, but would hurt GoDoc discoverability.

### Testing

4. **Integration coverage at 82.8%** — below 90% target
5. **Serialization coverage at 89.0%** — 1% below target, marshal error paths untested
6. **gentest at 80.8%** — pre-existing gap in assertion helpers

### Documentation

7. **17 "dead" exported symbols** — Never referenced from other packages. Need audit: document as intended public API or internalize.
8. **testhelpers GoDoc** — 8 exported symbols in testhelpers lack doc comments

### Type Model

9. **`RenderOptions` struct** has public fields (`Title`, `GraphID`, `Writer`, `ColorMode`) that can be set inconsistently. Functional options pattern would be safer but is a breaking change.
10. **Graph DTO types** (`tree_dto.go`, `graph_dto.go`) in serialization/ are unexported — fine, but there's no way for users to customize serialization output.

---

## F) Top #25 Things to Do Next

### High Impact, Low Effort

1. 🔍 Audit 17 "dead" exported symbols — document or internalize
2. 🧪 Add integration test coverage for streaming HTML footer (already partially done — verify)
3. 📝 Add GoDoc to 8 testhelpers exports (ErrWrite, ErrorWriter, FixedRenderer, etc.)
4. 🏷️ Add `//nolint:testableexamples` explanation comment to example files
5. 🧹 Remove stale `docs/planning/2026-05-28_round4-polish.md` or keep as reference

### High Impact, Medium Effort

6. 🏗️ Fix `MarkdownTable` to satisfy `TableRenderer` — dual-method or adapter pattern
7. 📊 Boost integration coverage to 90%+ — add missing error path tests
8. 📦 Boost serialization coverage to 90%+ — find testable error paths
9. 🎯 Add `WithFooterStyle` option to `table.New()` for composable footer styling
10. 📖 Write ADR for footer design decision (data on TableData vs RenderOptions)

### Medium Impact, Low Effort

11. 📐 Add `BenchmarkTableWithFooter` for terminal table
12. 🌳 Document tree format footer exclusion rationale in FEATURES.md
13. 📋 Add `// Deprecated` comment on direct `.Footer =` usage, pointing to `SetFooter()`
14. 🎨 Consider dark-mode footer styling for terminal table
15. 🔗 Cross-reference footer in all format entries in CHANGELOG.md

### Medium Impact, Medium Effort

16. 📊 Consider `_footer: true` marker in JSON/YAML serialization for round-trip fidelity
17. 🧪 Add fuzz tests for footer with edge-case characters
18. 📖 Extract footer rendering into `FooterRenderer` interface for custom footer logic
19. 🔧 Make error types consistent — all InvalidXxxError include Allowed values
20. 🏗️ Graph module FromTableData/FromTree dedup — shared helpers on GraphRendererMixin

### Lower Priority

21. 🤔 Multiple footer rows (`Footer [][]string` instead of `Footer []string`)
22. 📐 Footer row spanning (merged cells)
23. 📦 Version bump for v0.7.0 (footer feature is significant enough)
24. 🌙 Pre-v1 API stability audit (TODO_LIST #39)
25. 🧩 Community: Post to r/golang, submit to Awesome Go (TODO_LIST #40)

---

## G) Top #1 Question I Cannot Figure Out Myself

**Should JSON/YAML/TOML/JSONL include the footer row in their serialized output?**

The current behavior (skip footer) means `TableData` → JSON → `TableData` round-trips lose footer data. This is philosophically consistent ("footer is visual") but practically surprising.

Arguments for inclusion:

- Round-trip fidelity — serialize/deserialize should preserve all data
- `TableData` has the field, silently dropping it feels wrong
- Users exporting data might want totals included

Arguments against:

- Footer is a presentation concept, not data
- JSON `[{...}, {_footer: true, ...}]` is ugly and non-standard
- Data formats should remain pure data — use CSV for footer preservation

**My recommendation:** Keep current behavior. Add a `RenderOptions.IncludeFooter` flag (default false) that lets users opt-in to including footer in data formats. This preserves backward compatibility while enabling the use case.
