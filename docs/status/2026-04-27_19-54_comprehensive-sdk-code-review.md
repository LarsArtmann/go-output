# go-output SDK Comprehensive Code Review — Status Report

**Date:** 2026-04-27 19:54 CEST  
**Branch:** master  
**Commit:** 9ec8de8eb1ccd3e97d709b9762459ec072e80a9e  
**Reviewer:** AI Code Review Agent  

---

## a) FULLY DONE

### 1. Architecture & Design Review
- **Package structure** is clean and well-organized with clear separation of concerns
- **Type-safe enums** (`Format`, `SortBy`, `ColorMode`, `GraphShape`, `D2Direction`, `D2NodeShape`, `D2ArrowType`, `D2Constraint`) are consistently implemented with `Parse`, `String`, `AllowedValues`, `IsValid` methods
- **Branded IDs** (`BrandedID[Brand]`) provide compile-time type safety across D2, tree, and graph node types
- **Interface-based design** with `Renderer`, `TableRenderer`, `TreeOutputRenderer`, `GraphRenderer`, `StreamingRenderer` enables extensibility
- **Registry pattern** with thread-safe concurrent access via `sync.RWMutex`
- **Composition over inheritance** throughout (e.g., `GraphRendererMixin`, `tableDataBase`)

### 2. Code Quality Verification
- **All tests pass:** `go test -cover ./...` — 100% pass rate across all packages
- **Linter clean:** `golangci-lint run ./...` — 0 issues
- **Build clean:** `go build ./...` — no errors
- **go vet clean:** no warnings
- **Test coverage:** Main package 91.0%, cmdguard 100%, enum 94.7%, escape 100%, sort 86.7%, table 100%

### 3. Bugs Found & Fixed During Review

#### Bug #1: Unsigned Integer Sorting Panic (CRITICAL)
- **Location:** `sort/sorter.go:110`
- **Issue:** `compareFieldValues` used `a.Int()` for ALL integer types including unsigned (`reflect.Uint`, `Uint8`, etc.). Calling `Int()` on an unsigned `reflect.Value` panics.
- **Fix:** Split into separate cases: signed uses `Int()`, unsigned uses `Uint()`
- **Test added:** `TestSorter_Sort_UnsignedInt` with `uint64` max value validation

#### Bug #2: ColorMode Auto Never Enabled Colors (CRITICAL)
- **Location:** `color.go:77-79`
- **Issue:** `isStdoutTerminal()` only checked `FORCE_COLOR` env vars — it NEVER actually detected if stdout was a terminal. This meant `ColorModeAuto.ShouldColor()` always returned `false` without `FORCE_COLOR` set.
- **Fix:** Added `golang.org/x/term` dependency and `term.IsTerminal(int(os.Stdout.Fd()))` check. Falls back to env var override.
- **Linter config updated:** Added `golang.org/x/term` to depguard allow lists

#### Bug #3: Markdown Alignment Markers Missing (MODERATE)
- **Location:** `markdown.go:131-141`
- **Issue:** The separator row (`|---|---|`) never included GitHub Flavored Markdown alignment markers (`:--`, `--:`, `:--:`). Alignment settings only affected whitespace padding, not actual Markdown semantics.
- **Fix:** Added `getAlignmentMarkers()` method generating proper colon prefixes/suffixes
- **Tests added:** `TestMarkdownAlignmentMarkers` with left/right/center marker validation

#### Bug #4: CSV/TSV Example Formatting Inconsistency (MODERATE)
- **Location:** `examples/basic/main.go:128-139`
- **Issue:** `projectToRow` returned raw integers (`90`, `7`) while all other formats showed formatted values (`90%`, `7/10`). CSV/TSV output was inconsistent with table, markdown, HTML, etc.
- **Fix:** Consolidated `projectToRow` and `projectToTableDataRow` into single function with consistent formatting
- **Dead code removed:** Unused `handleErrorWithContext` function

#### Bug #5: Missing Integration Test Coverage (MINOR)
- **Location:** `integration/integration_test.go:46-57`
- **Issue:** `TestAllFormatsRender` omitted `FormatTSV` and `FormatXML` from the formats slice, leaving two formats untested in the comprehensive integration test.
- **Fix:** Added both formats to the test list

#### Bug #6: Sort Field Validation Gap (MINOR)
- **Location:** `sort/sorter.go:63-70`
- **Issue:** `defaultLess` validated that field existed in `a` but not in `b`. If sorting heterogeneous interface slices where `a` has the field but `b` doesn't, `compareFieldValues` receives an invalid `reflect.Value` for `b`.
- **Fix:** Added `!fieldB.IsValid()` check alongside existing `!fieldA.IsValid()` guard

---

## b) PARTIALLY DONE

### 1. Streaming Error Handling
- **Status:** Acknowledged limitation
- **Issue:** `StreamingHTMLRenderer.Render()` ignores `Stream()` error (`_ = r.Stream(&b)`)
- **Constraint:** The `Renderer` interface returns `string` (not `(string, error)`), so error propagation is structurally impossible without breaking API compatibility
- **Mitigation:** `strings.Builder.Write()` never returns errors in practice, making this safe for the `Render()` use case. True streaming via `Stream(io.Writer)` properly returns errors.

### 2. Example Code Completeness
- **Status:** Functional but could be richer
- **Issue:** `examples/d2/main.go` is excellent, but `examples/basic/main.go` could demonstrate more advanced features (sorting, registry, branded IDs)

### 3. Documentation Accuracy
- **Status:** Mostly accurate
- **Issue:** `README.md` still shows `charm.land/lipgloss/v2 v2.0.2` in Dependencies section but `go.mod` has `v2.0.3`

---

## c) NOT STARTED

### 1. Performance Optimizations
- No performance profiling or benchmarking improvements attempted
- `enum.go` `joinStrings()` uses recursion instead of `strings.Join` — not addressed
- Markdown table column width calculation rebuilds on every `Render()` call — could be optimized for repeated renders

### 2. API Completeness Gaps
- `MarshalYAMLIndent` function does not exist (inconsistent with JSON/XML)
- `UnmarshalXML` function does not exist
- `MarshalTSV` doesn't support `*TableData` with headers (only `[][]string`)

### 3. Advanced Features
- No validation that table rows match header column count
- No support for Markdown table captions
- No support for D2 diagram themes or custom D2 configuration blocks

### 4. Cross-Platform Testing
- No Windows-specific terminal detection validation
- No testing in CI environment simulation

---

## d) TOTALLY FUCKED UP!

**Nothing.** After fixes, the codebase is in excellent shape. All tests pass, linter is clean, build succeeds, and the architecture is sound. The issues found were real bugs but all have been resolved.

---

## e) WHAT WE SHOULD IMPROVE

### 1. API Consistency
- Add `MarshalYAMLIndent` for parity with `MarshalJSONIndent` and `MarshalXMLIndent`
- Consider adding `UnmarshalXML` for completeness
- `MarshalTSV` should accept `*TableData` to write headers (or add `MarshalTSVFromTableData`)

### 2. Error Handling in Examples
- Integration test `renderCSVFormat` and `renderTSVFormat` ignore write errors with `_ =`
- Examples should demonstrate proper error handling patterns for production use

### 3. Documentation
- Update `README.md` dependency versions to match `go.mod`
- Add godoc examples (`ExampleXXX` functions) for key APIs — currently missing
- Document the `StreamingRenderer` adapter limitation more prominently

### 4. Performance
- Cache column widths in `MarkdownTable` if headers/rows haven't changed
- Replace recursive `joinStrings` in `enum.go` with `strings.Join`
- Consider pooling `strings.Builder` instances in high-throughput scenarios

### 5. Test Coverage Gaps
- `sort` package at 86.7% — the `defaultLess` struct-kind check and reflect edge cases could use more coverage
- Examples packages at 0% — not critical but would demonstrate usage
- No tests for `MarshalXML` with `xml.Marshaler` interfaces

### 6. Edge Cases
- `DOTRenderer.writeNodeAttr` doesn't quote attribute values — works for typical values but could break with spaces/special chars
- `MermaidRenderer` doesn't sanitize node IDs in `MermaidFlowchartRenderer` (IDs are auto-generated as `rowN` so safe, but custom IDs could break)
- `D2Diagram.writeClasses` iterates over map in non-deterministic order — D2 output could vary between runs

---

## f) Top #25 Things We Should Get Done Next

### Critical (P0)
1. **Add godoc examples** (`ExampleMarshalJSON`, `ExampleNewMarkdownTable`, etc.) — improves discoverability
2. **Update README dependency versions** to match actual `go.mod`
3. **Add `MarshalYAMLIndent`** for API consistency

### High Priority (P1)
4. Cache `MarkdownTable` column widths to avoid recalculation on every `Render()`
5. Replace recursive `joinStrings` with `strings.Join` in `enum.go`
6. Add `MarshalTSVFromTableData` that writes headers from `*TableData`
7. Fix `D2Diagram.writeClasses` map iteration order for deterministic output
8. Add `UnmarshalXML` wrapper for API completeness
9. Improve example error handling (don't ignore `WriteHeader`/`WriteRow` errors)

### Medium Priority (P2)
10. Add row length validation to `TableData.AddRow()` (warn if row doesn't match header count)
11. Escape DOT node attributes that could contain spaces (fillcolor, color)
12. Sanitize Mermaid node IDs in `MermaidFlowchartRenderer` (not just labels)
13. Add `MarshalXMLIndent` tests with actual struct types
14. Add benchmark for `sort.Sorter` with large datasets
15. Add test for `ColorModeAuto.ShouldColor()` with mocked terminal detection
16. Document `StreamingRenderer` adapter behavior in FORMAT_ARCHITECTURE.md

### Lower Priority (P3)
17. Add `TableData.RemoveRow()` method
18. Support Markdown table captions (`<caption>` tag in HTML mode)
19. Add `D2Diagram.SetTheme()` for D2 theme configuration
20. Consider adding `Format.AutoDetect()` based on file extension or MIME type
21. Add `TreeNode.RemoveChild()` method
22. Support CSV/TSV reading/parsing (bidirectional conversion)
23. Add visual regression tests for terminal table output (snapshot testing)
24. Consider adding `lipgloss` adaptive color support (light/dark mode)
25. Add fuzz tests for `MarshalJSON`/`UnmarshalJSON` roundtrip

---

## g) Top #1 Question I Cannot Figure Out Myself

### Why does the `Renderer` interface return `string` instead of `(string, error)`?

The `Renderer` interface is defined as:
```go
type Renderer interface {
    Render() string
}
```

This design decision forces implementations like `StreamingHTMLRenderer.Render()` to silently swallow errors from `Stream()`:
```go
func (r *StreamingHTMLRenderer) Render() string {
    var b strings.Builder
    _ = r.Stream(&b)  // Error intentionally discarded
    return b.String()
}
```

While `strings.Builder` never returns write errors in practice, this is a structural limitation. If a custom `io.Writer` that CAN fail is passed to `Stream()`, the error is lost when using `Render()`.

**Was this intentional** to keep the simple API surface minimal, or **should it be reconsidered** for v3 to return `(string, error)`? The tradeoff is:
- **Current:** Simple API, no error handling boilerplate for 99% of use cases
- **Alternative:** Safer API but every renderer call requires error checking

I suspect this was intentional for simplicity, but I'd like confirmation before proposing a v3 API change.

---

## Commit Summary

All review findings have been committed as:

```
commit 9ec8de8eb1ccd3e97d709b9762459ec072e80a9e
Author: Lars Artmann <git@lars.software>
Date:   Mon Apr 27 13:30:03 2026 +0200

    feat: add terminal detection, GFM alignment markers, and unsigned int sorting
```

**Files changed:** 10 files, +111 lines, -25 lines  
**Key areas:** sort/sorter.go, color.go, markdown.go, examples/basic/main.go, integration/integration_test.go, sort/sort_test.go, markdown_test.go, .golangci.yml, go.mod, go.sum
