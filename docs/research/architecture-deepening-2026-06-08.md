# Architecture Deepening Report — go-output

**Date:** 2026-06-08\
**Scope:** Root package, all sub-modules, test infrastructure\
**Method:** Deletion test applied to every module suspected of being shallow

---

## Executive Summary

The codebase is already well-architected: bounded contexts are respected (root imports no sub-modules), `Renderer` is a strong small interface, and the `TableData` / `TreeNode` / `GraphNode` data models provide real leverage. Most findings are minor deepening opportunities around duplication, leaky seams, and shallow wrappers.

| Strength          | Count |
| ----------------- | ----- |
| Strong candidates | 2     |
| Worth exploring   | 4     |
| Speculative       | 3     |

---

## Candidate 1: HTMLRenderer & StreamingHTMLRenderer duplicate table generation

- **Files:** `markup/html.go`, `markup/streaming.go`
- **Problem:** Two independent implementations generate the exact same HTML table structure. `HTMLRenderer` builds a `strings.Builder`; `StreamingHTMLRenderer` writes chunks to `io.Writer`. They share `writeMarkupRow`/`writeMarkupColumns` but replicate the overall sequencing.
- **Solution:** Implement `HTMLRenderer.Render()` by delegating to `StreamingHTMLRenderer.Stream(strings.Builder)`.
- **Benefits:** Single source of truth for HTML table generation. One file owns the sequence; the other merely adapts I/O.
- **Strength:** Strong

---

## Candidate 2: `formatCapabilities` global map breaks bounded-context seams

- **Files:** `shape.go` (lines 66-83)
- **Problem:** Core package owns a hardcoded `map[Format][]Shape` listing all 16 formats. Adding a new format in a sub-module requires editing core. This violates bounded-context separation.
- **Solution:** Invert the dependency. Each format module registers its supported shapes at `init()` time, mirroring the `RegisterTableDataMarshaler` pattern.
- **Benefits:** New formats can be added without touching core. The seam between core and sub-modules becomes clean.
- **Strength:** Strong

---

## Candidate 3: GraphRendererMixin is a partial abstraction with leaky accessors

- **Files:** `graph.go`, `graph/dot.go`, `graph/mermaid.go`, `plantuml/plantuml.go`, `d2/d2_convert.go`
- **Problem:** `GraphRendererMixin` is embedded by DOT, Mermaid, and PlantUML, but D2Diagram bypasses it entirely because D2 types are richer. D2 satisfies `GraphRenderer` via a lossy adapter in `d2_convert.go`.
- **Solution:** Remove `NodesPtr()` and `EdgesPtr()` from the mixin. Consider splitting D2 out of the `GraphRenderer` interface entirely.
- **Strength:** Worth exploring

---

## Candidate 4: `RenderTableData` hardcodes which formats are unsupported

- **Files:** `render_tabledata.go` (lines 59-89)
- **Problem:** Core package code contains knowledge about what sub-modules cannot do. This is a seam leak.
- **Solution:** If a sub-module can't support `RenderTableData` dispatch, it simply doesn't register a `TableDataMarshaler`. Core becomes agnostic to sub-module capabilities.
- **Strength:** Worth exploring

---

## Candidate 5: `enum` package is shallow

- **Files:** `enum/enum.go`, plus every caller
- **Problem:** The package exports `Parse`, `Contains`, `AllowedStrings`, `AllowedValues`. The implementation is trivial generic loops over slices. Every enum type repeats the same boilerplate.
- **Deletion test:** Deleting `enum` would force each enum to write its own 5-line parse loop. Complexity would spread but not increase dramatically.
- **Solution:** Absorb `enum` into core (66 lines, no external deps) or replace repetitive boilerplate with a single generic `Enum[T]` type.
- **Strength:** Worth exploring

---

## Candidate 6: `marshal.go` wrappers are shallow error-context passthroughs

- **Files:** `marshal.go`
- **Problem:** `MarshalFormat`, `UnmarshalFormat`, and `MarshalJSONIndent` are thin wrappers that add `fmt.Errorf` context around stdlib calls.
- **Deletion test:** Deleting them would move ~3 lines of error wrapping into each caller.
- **Solution:** Inline the error wrapping at call sites, or move them into `serialization/` where the callers live.
- **Strength:** Worth exploring

---

## Candidate 7: `testhelpers` mixes deep generic runners with shallow one-liners

- **Files:** `testhelpers/helpers.go`, `testhelpers/renderers.go`, `testhelpers/writers.go`
- **Problem:** `AssertContains` is just `strings.Contains` + `t.Error`. `AssertEqual` is just `!=` + `t.Errorf`. These shallow ones don't justify a separate module.
- **Solution:** Delete shallow helpers and inline them in tests. Keep deep generic test runners.
- **Strength:** Worth exploring

---

## Candidate 8: `updateMaxWidths` pure function extracted just for testability

- **Files:** `markdown.go` (lines 126-147)
- **Problem:** `updateMaxWidths` exists so tests can hit it in isolation, but real bugs live in how widths integrate with alignment and color codes inside `Render()`.
- **Solution:** Inline `updateMaxWidths` back into `calculateColumnWidths`.
- **Strength:** Speculative

---

## Candidate 9: `ids.go` branded ID re-exports are a very thin seam

- **Files:** `ids.go`
- **Problem:** Pure re-export of `go-branded-id`. Interface and implementation are identical in complexity.
- **Deletion test:** Deleting it would force callers to import `go-branded-id` directly.
- **Solution:** Keep if planning to swap implementations later; otherwise delete the indirection.
- **Strength:** Speculative

---

## Top Recommendation

**Merge `HTMLRenderer` and `StreamingHTMLRenderer` table generation into a single implementation.**

Both files encode the same HTML table structure; the only difference is I/O strategy. By implementing the non-streaming renderer as an adapter that streams into a `strings.Builder`, you concentrate the HTML generation logic in one place. This maximizes locality and leverage.

After that, address the **`formatCapabilities` inversion** (Candidate 2). It has broader architectural impact because it cleans the seam between core and all sub-modules, enabling true plug-in formats.
