# Footer Row Polish — Status Report (Round 2 + Round 3)

**Date:** 2026-05-28 02:50
**Branch:** master
**Base commit:** a564b67 (Round 1 end)
**Head commit:** (pending)
**Total commits:** 12+ (across Round 2 + Round 3)
**Scope:** Architecture fixes, go.mod consistency, Validate() wiring, test gaps, SetFooter bug fix

---

## Summary

Three rounds of polishing the footer row feature:

- **Round 1** (13 commits): Core footer implementation across 6 modules
- **Round 2** (7 commits): go.mod unification, Validate() wiring, formatter fixes
- **Round 3** (ongoing): WriteFooter tests, SetFooter multi-call bug fix

---

## A) FULLY DONE ✅

### Round 1 (commits 98c0ae3..a564b67)

- `TableData.Validate()` method with `errColumnMismatch` error type
- `delimited/` dedup: `tableDataWriter` interface + `marshalFromTableData()` shared helper
- `WriteFooter()` method on `CSVWriter` and `TSVWriter`
- `class="footer-cell"` CSS on HTML `<tfoot>` and streaming HTML
- AsciiDoc switched from raw `len()` to `HasFooter()`
- `ExampleTable_SetFooter` GoDoc example
- Footer benchmarks for Markdown and CSV
- Streaming HTML footer integration test
- Markdown footer alignment test (verified working)
- Footer format matrix in README
- Footer in examples/basic
- Updated FEATURES.md, CHANGELOG.md, TODO_LIST.md, AGENTS.md
- `TestBrandedIDFormat` fix for go-branded-id v0.3.0

### Round 2 (commits 3894ad2..d635279)

- All 9 `go.mod` files unified to `go 1.26.3` (from mixed 1.26.2/1.26.3)
- `integration/go.mod` root dep fixed: `v0.5.0` → `v0.0.0`
- `go.work.example` updated with `./testhelpers/graphtest`
- All `go.sum` files regenerated
- `data.Validate()` wired into `RenderTableData()` — catches footer/header column mismatch before dispatch
- `TestRenderTableData_ValidateRejectsFooterMismatch` added
- Formatter suggestions applied: `errors.New` for static strings, import grouping

### Round 3 (in progress)

- TSV WriteFooter test fixed (missing `bytes` import)
- `table.SetFooter` multi-call bug fixed: `footerRowIndex` field tracks the correct bold row
- `TestTableSetFooter_MultipleCalls` added to verify fix

### Code Quality

- All 14 modules build clean
- All tests pass
- `golangci-lint` clean on all modules
- All files under 350-line limit

---

## B) PARTIALLY DONE ⚠️

Nothing partially done — all items from the previous status report's "Partially Done" section have been completed:

- ✅ `delimited/` dedup — extracted `marshalFromTableData()`
- ✅ Example coverage — footer added to `projectsToTableData()`
- ✅ Footer format matrix in README
- ✅ GoDoc example for `SetFooter`
- ✅ Benchmarks for footer rendering
- ✅ `class="footer-cell"` CSS
- ✅ `WriteFooter()` streaming methods

---

## C) NOT STARTED ❌

1. **Serialization footer option** — JSON/YAML/TOML could include footer with `_footer: true` marker (controversial)
2. **Multiple footer rows** — Only single footer supported (slice of strings, not slice of slices)
3. **Footer-specific escaping** — No format-specific footer escaping beyond cell escaping
4. **`WithFooterStyle` option** — Composable footer styling for users who override StyleFunc
5. **ADR for footer design** — Document the decision to put footer on TableData vs RenderOptions

---

## D) TOTALLY FUCKED UP 💥

Nothing catastrophic. Minor issues resolved:

- **SetFooter multi-call bug**: Fixed in Round 3. `footerRowIndex` now correctly tracks the bold-styled row. Previous footer rows become unstyled data rows (lipgloss doesn't support row removal — documented limitation).
- **Stale LSP diagnostics**: `markdown_footer_test.go` shows false-positive "duplicate decl" in gopls — the file is valid Go, builds and tests clean. LSP needs restart.

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **`table.SetFooter` doc could warn more strongly** — Now tracked via `footerRowIndex`, but previous footer rows still appear as data. Consider panicking on 2nd call or making it truly idempotent.
2. **`WithFooterStyle` option** — Users overriding `StyleFunc` lose footer bold. A composable option would help.
3. **Serialization footer** — JSON/YAML round-trip fidelity loss (footer silently dropped).

### Testing

4. **Edge case: footer longer than headers** — No test
5. **Edge case: empty footer cells** — `Footer: []string{"", ""}` untested
6. **Coverage**: Root module at ~88.6% (below 90% target)

### Documentation

7. **Streaming `WriteFooter` docs** — Undocumented that CSV/TSV writers support `WriteFooter()`
8. **Footer design rationale** — No ADR

---

## F) Top #25 Things to Do Next

### High Impact, Low Effort

1. 🧪 Add edge case tests (footer longer than headers, empty cells, nil footer)
2. 📊 Check root module coverage — get back above 90%
3. 🔍 Audit all renderers for `HasFooter()` consistency (all should use it)
4. 📝 Document `WriteFooter()` on CSVWriter/TSVWriter in GoDoc

### High Impact, Medium Effort

5. 🏗️ Add `WithFooterStyle` option to `table.New()` for composable footer styling
6. 📖 Write ADR for footer on TableData vs RenderOptions
7. 🎯 Make `SetFooter` idempotent or panic on 2nd call
8. 🧹 Clean up pre-existing gopls false positives

### Medium Impact, Low Effort

9. 🌳 Document tree format footer exclusion rationale
10. 📋 Add `// Deprecated` on `.Footer =` pointing to `SetFooter()`
11. 🔗 Cross-reference all format entries in CHANGELOG.md
12. 📐 Add `BenchmarkTableWithFooter` for terminal table

### Medium Impact, Medium Effort

13. 📊 Consider `_footer: true` marker in JSON/YAML for round-trip fidelity
14. 🧪 Add integration test for streaming HTML with footer (DONE in Round 1 — verify)
15. 🎨 Add dark-mode footer styling for terminal table
16. 🧩 Extract footer rendering into `FooterRenderer` interface

### Lower Priority

17. 🤔 Multiple footer rows (`Footer [][]string`)
18. 📐 Footer row spanning (merged cells)
19. 🌙 Custom footer CSS class names
20. 🧪 Fuzz testing for footer with edge-case characters
21. 📖 Blog post on footer design patterns for CLI libraries
22. 🔧 `go-error-family` integration for typed error hierarchy
23. 📦 Version bump for v0.6.0 (footer is a minor version bump)
24. 🎯 Performance: lazy footer rendering (skip if no footer set)
25. 🧩 Plugin system for custom footer formatters

---

## G) Top #1 Question I Cannot Figure Out Myself

**Same as previous report — should JSON/YAML/TOML/JSONL include the footer?**

The question remains unresolved. The current behavior (skip footer in data formats) is consistent with the "footer is a visual concept" philosophy. But it means `TableData` → JSON → `TableData` round-trips lose footer data.

The pragmatic answer: **keep it as-is**. Users who need round-trip fidelity can use CSV (which includes footer). Data formats should remain pure data.
