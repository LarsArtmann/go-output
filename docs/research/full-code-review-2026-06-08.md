# Full Code Review: go-output

**Date:** 2026-06-08  
**Reviewer:** Senior Staff+ Software Architect  
**Scope:** All 14 modules, ~120 .go files (production + test)  

---

## Executive Summary

**Files Reviewed:** ~120 `.go` files across 14 modules.  
**Overall Assessment:** Well-architected, strongly-typed Go library with excellent test hygiene. Two critical bugs and several API anti-patterns need fixing before production-hardened.

| Severity | Count |
|----------|-------|
| Critical | 2 |
| High | 5 |
| Medium | 12 |
| Low | 15 |

---

## Critical Issues

### 1. D2 `writeClasses` output is non-deterministic
**File:** `d2/d2_render.go`  
**Problem:** `D2Diagram.writeClasses` iterates over `d.classes` map directly. Go map iteration order is randomized. This causes non-deterministic D2 output, breaking snapshot tests and content-addressable caching.  
**Fix:** Sort class names before iterating.

### 2. `D2ArrowType` parse/valid inconsistency
**File:** `d2/d2_enum.go`  
**Problem:** `ParseD2ArrowType("")` returns an error because `d2ArrowTypeValues` does NOT include `D2ArrowNone` (empty string). However, `D2ArrowType.IsValid()` explicitly returns `true` for `D2ArrowNone`. Incoherent behavior.  
**Fix:** Add `D2ArrowNone` to `d2ArrowTypeValues` or remove the special case from `IsValid()`.

---

## High Issues

### 3. `FormatJSON` not registered in `RenderTableData` registry
**File:** `serialization/json.go`  
**Problem:** `serialization/json.go` does not call `output.RegisterTableDataMarshaler` in `init()`. Contrast with `yaml.go`, `toml.go`, `jsonl.go` which all register themselves. `output.RenderTableData(data, output.FormatJSON)` returns `UnsupportedFormatError`.  
**Fix:** Add `init()` calling `output.RegisterTableDataMarshaler(output.FormatJSON, renderJSONTableData)`.

### 4. Variadic `RenderOptions` is misleading
**File:** `render_tabledata.go`  
**Problem:** `RenderTableData(data *TableData, format Format, opts ...RenderOptions)` uses variadic `opts` but only consumes `opts[0]`. Signature implies multiple options are valid.  
**Fix:** Change to single `opts RenderOptions` value, or implement proper merge.

### 5. `StreamingRendererFromRenderer` is a false promise
**File:** `streaming.go`  
**Problem:** The adapter calls `Render()` which builds the entire string in memory, then writes it at once. The name `StreamingRenderer` combined with a non-streaming adapter is misleading.  
**Fix:** Rename to `BufferedRendererAdapter`, or implement chunked writes.

### 6. `GraphRendererMixin.NodesPtr/EdgesPtr` breaks encapsulation
**File:** `graph.go`  
**Problem:** Returns pointers to internal slices, allowing external mutation.  
**Fix:** Remove these methods. Provide controlled mutation methods like `AddNode`/`AddEdge`, or return copies.

---

## Medium Issues

### 7. `TableData` nil safety inconsistent
**File:** `tabledata.go`  
- `Validate()` returns `nil` for nil receiver (dangerous)
- `RowCount()` and `ColCount()` panic on nil receiver
- `ToMapSlice()` and `CreateRowEdges()` correctly check nil

### 8. `RegisterTableDataMarshaler` silently overwrites
**File:** `render_tabledata.go`  
Should return error or panic on duplicate registration.

### 9. `table.buildStyleFunc` allocates `lipgloss.NewStyle()` per row
**File:** `table/table.go`  
Cache base styles at construction time.

### 10. `ColorModeAuto.ShouldColor()` has side effects
**File:** `color.go`  
Reads environment and checks terminal status, making it non-deterministic and hard to test.

### 11. Incomplete AsciiDoc escaping
**File:** `markup/asciidoc.go`  
Only escapes `|` but AsciiDoc has many more special characters (`*`, `_`, `` ` ``, `~`, `^`).

### 12. `MermaidText` escaping incomplete
**File:** `escape/escape.go`  
Replaces `"` with `'` but does not escape `'` itself.

### 13. `StreamingHTMLRenderer` string concatenation allocations
**File:** `markup/streaming.go`  
`[]byte("<th>" + escape.HTML(h) + "</th>\n")` allocates intermediate strings per cell.

### 14. `DOTRenderer.writeEdge` poor slice preallocation
**File:** `graph/dot.go`  
`make([]string, 0)` for 4 attributes — should use `make([]string, 0, 4)`.

### 15. `MermaidRenderer.Render()` uses `fmt.Fprintf` in hot path
**File:** `graph/mermaid.go`  
Parses format string on every node/edge. Use `b.WriteString` + manual concatenation.

### 16. Duplicate slug logic across modules
**Files:** `d2/d2_convert.go`, `graph/dot.go`, `graph/mermaid.go`, `plantuml/convert.go`  
`strings.ReplaceAll(label, " ", "_")` duplicated 4 times.

### 17. `MarkdownTable.writeCell` inconsistent padding approach
**File:** `markdown.go`  
Uses `fmt.Fprintf` for right-alignment but `strings.Repeat` for left/center.

### 18. `Alignment` const shadowing
**File:** `markdown.go`  
`alignmentLeft` appears in two const blocks; second shadows first.

---

## Low Issues

19. `marshal.go` wrappers are shallow error-context passthroughs — YAGNI?
20. `ASCIITreeRenderer` name misleading (uses Unicode, not ASCII)
21. `TreeNode.Depth()` is O(n) and uncached
22. `HandleError` in examples is vague (`PrintAndExit` would be clearer)
23. `testhelpers` mixes deep generic runners with shallow one-liners
24. `ids.go` is a pure re-export of `go-branded-id` — thin seam
25. `updateMaxWidths` pure function extracted just for testability
26. `renderFullHTMLDocument` uses string concatenation instead of `html/template`
27. `XMLWriter.WriteRow` error loses row index context
28. `D2NodeStyle.Opacity` has no bounds validation
29. `sanitizePlantUMLID` incomplete (only spaces/dashes)
30. `escape.D2` chains 4 `strings.ReplaceAll` calls — multiple allocations
31. `MermaidID` vs `MermaidSlug` inconsistent `/` handling
32. `Test*` prefix on exported helpers violates Go convention
33. Missing race tests for concurrent registry access
34. `render_helpers_test.go` ignores errors in many helpers

---

## Praise

1. **Branded IDs** — compile-time type safety for node IDs is exceptional
2. **Multi-module workspace** — clean 14-module split, zero transitive bloat
3. **Test hygiene** — `t.Parallel()`, fuzz tests, benchmarks with `b.Loop()`, table-driven tests
4. **Error wrapping** — consistent `fmt.Errorf("...: %w", err)` with context
5. **XSS prevention** — HTML/XML cell escaping is thorough
6. **Registry extensibility** — `init()`-based registration avoids import cycles
7. **Interface boundaries** — tight, focused interfaces with compile-time checks
8. **Documentation** — thorough GoDoc, package comments, example functions
9. **User journey tests** — high-level workflow testing (`userjourney_test.go`)

---

## Next Steps (Pareto Priority)

1. Fix 2 CRITICAL bugs (D2 determinism, D2ArrowType inconsistency)
2. Register `FormatJSON` in `RenderTableData` registry
3. Fix `RenderOptions` variadic anti-pattern
4. Add nil-receiver safety to `TableData`
5. Cache `lipgloss.NewStyle()` in table renderer
6. Extract shared `SlugifyID` helper
7. Complete AsciiDoc/Mermaid escaping
8. Profile benchmarks after allocation fixes
