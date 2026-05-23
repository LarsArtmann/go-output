# ADR: Data Shape × Format Capability Matrix

**Date:** 2026-05-16
**Status:** ACCEPTED & IMPLEMENTED

## Context

The current category system (`FormatCategory` with `CategoryTable`, `CategoryTree`, `CategoryGraph`) treats each format as belonging to a single category. This is fundamentally wrong:

1. **D2 appears in both `tableFormats` and `graphFormats`** — the code already knows D2 is multi-shape, but the `Category()` method can only return one value (it picks "graph").
2. **JSON and YAML serialize `any`** — they're shape-agnostic but have no `Renderer` implementations.
3. **DOT, Mermaid, D2 all support all 3 shapes** via `*FromTableData()` and `*FromTree()` converter functions.
4. **HTML already has two separate renderers** (`HTMLRenderer` for tables, `HTMLTreeRenderer` for trees).

The concept of "category" conflates two independent axes:

- **Data Shape**: What kind of data do you have? (table, tree, graph)
- **Format**: How do you want to render it? (json, csv, d2, ...)

These should be a **capability matrix**, not a classification.

## Decision

Replace `FormatCategory` with a `Shape` type and a capability matrix `map[Format][]Shape`.

### New Types

```go
// Shape represents a data shape that a format can render.
type Shape string

const (
    ShapeTable Shape = "table" // Tabular data (headers + rows)
    ShapeTree  Shape = "tree"  // Hierarchical data (parent-child nodes)
    ShapeGraph Shape = "graph" // Network data (nodes + edges)
)
```

### Capability Matrix

```go
// formatCapabilities maps each format to the data shapes it supports.
var formatCapabilities = map[Format][]Shape{
    FormatTable:    {ShapeTable},
    FormatJSON:     {ShapeTable, ShapeTree, ShapeGraph},
    FormatCSV:      {ShapeTable},
    FormatTSV:      {ShapeTable},
    FormatXML:      {ShapeTable},
    FormatMarkdown: {ShapeTable},
    FormatD2:       {ShapeTable, ShapeGraph},
    FormatYAML:     {ShapeTable, ShapeTree, ShapeGraph},
    FormatHTML:     {ShapeTable, ShapeTree},
    FormatTree:     {ShapeTree},
    FormatMermaid:  {ShapeTable, ShapeGraph},
    FormatDOT:      {ShapeTable, ShapeGraph},
}
```

### New Methods on `Format`

```go
// Supports returns true if the format can render the given data shape.
func (f Format) Supports(s Shape) bool

// Shapes returns all data shapes this format supports.
func (f Format) Shapes() []Shape

// FormatsForShape returns all formats that support the given data shape.
func FormatsForShape(s Shape) []Format
```

### Backward Compatibility

Keep `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, and `Category()` as deprecated wrappers:

```go
// Deprecated: Use f.Supports(ShapeTable) instead.
func (f Format) IsTableFormat() bool { return f.Supports(ShapeTable) }

// Deprecated: Use f.Supports(ShapeTree) instead.
func (f Format) IsTreeFormat() bool { return f.Supports(ShapeTree) }

// Deprecated: Use f.Supports(ShapeGraph) instead.
func (f Format) IsGraphFormat() bool { return f.Supports(ShapeGraph) }

// Deprecated: Use f.Shapes() instead. Returns the primary shape.
func (f Format) Category() FormatCategory
```

`FormatCategory` stays as-is for backward compat, with `FormatShape` as the modern replacement.

## Migration Plan

### Phase 1: Add Shape type + capability matrix (non-breaking)

1. Add `Shape` type, constants, `AllShapes` slice to `format.go`
2. Add `formatCapabilities` matrix
3. Add `Supports()`, `Shapes()`, `FormatsForShape()` methods
4. Deprecate `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` — redirect to new API
5. Update all internal code to use `Supports()` instead of `Is*Format()`
6. Update tests
7. Update README

### Phase 2: Add shape-specific renderer constructors (future, not this PR)

- `NewJSONTableRenderer(data *TableData) Renderer`
- `NewJSONTreeRenderer(root *TreeNode) Renderer`
- `NewYAMLTableRenderer(data *TableData) Renderer`
- etc.

This is out of scope for the initial refactor but the API design enables it.

## Capability Matrix (truth table)

| Format     | Table | Tree | Graph | Notes                          |
| ---------- | :---: | :--: | :---: | ------------------------------ |
| `table`    |  ✅   |  ❌  |  ❌   | Terminal table only            |
| `json`     |  ✅   |  ✅  |  ✅   | Shape-agnostic serializer      |
| `csv`      |  ✅   |  ❌  |  ❌   | Flat rows only                 |
| `tsv`      |  ✅   |  ❌  |  ❌   | Flat rows only                 |
| `xml`      |  ✅   |  ❌  |  ❌   | Currently table-only           |
| `markdown` |  ✅   |  ❌  |  ❌   | Tables only (no tree syntax)   |
| `d2`       |  ✅   |  ❌  |  ✅   | SQL tables + diagrams          |
| `yaml`     |  ✅   |  ✅  |  ✅   | Shape-agnostic serializer      |
| `html`     |  ✅   |  ✅  |  ❌   | Tables + collapsible trees     |
| `tree`     |  ❌   |  ✅  |  ❌   | ASCII tree only                |
| `mermaid`  |  ✅   |  ❌  |  ✅   | Flowcharts (has FromTableData) |
| `dot`      |  ✅   |  ❌  |  ✅   | Graphs (has FromTableData)     |

## Consequences

**Positive:**

- Correctly models reality (formats support multiple shapes)
- `f.Supports(ShapeTable)` reads better than `f.IsTableFormat()`
- Enables `FormatsForShape(ShapeGraph)` — "give me all formats I can use for graph data"
- Backward compatible — deprecated methods still work
- Clean path to future shape-specific renderers

**Negative:**

- `FormatCategory` and `Category()` stay around until next major version
- Marginally more code in format.go (but removes the three separate maps)
