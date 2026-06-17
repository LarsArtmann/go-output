# go-output Format Architecture

## Overview

This document describes the extensible format architecture for go-output, supporting 16 output formats across 3 data shapes in a unified way.

## Data Shapes

Formats are classified by the data shapes they support. Each format may support multiple shapes.

### Shape Capability Matrix

| Format   | ShapeTable | ShapeTree | ShapeGraph |
| -------- | :--------: | :-------: | :--------: |
| table    |     Y      |           |            |
| json     |     Y      |     Y     |     Y      |
| csv      |     Y      |           |            |
| tsv      |     Y      |           |            |
| xml      |     Y      |           |            |
| markdown |     Y      |           |            |
| yaml     |     Y      |     Y     |     Y      |
| html     |     Y      |     Y     |            |
| tree     |            |     Y     |            |
| d2       |     Y      |     Y     |     Y      |
| mermaid  |     Y      |     Y     |     Y      |
| dot      |     Y      |     Y     |     Y      |
| jsonl    |     Y      |           |            |
| asciidoc |     Y      |           |            |
| toml     |     Y      |     Y     |     Y      |
| plantuml |     Y      |     Y     |     Y      |

### Querying Capabilities

```go
// Check if a format supports a specific shape
output.FormatJSON.Supports(output.ShapeTable) // true

// Get all shapes a format supports
output.FormatD2.Shapes() // [ShapeTable, ShapeGraph]

// Get all formats that support a shape
output.FormatsForShape(output.ShapeGraph) // [json, yaml, d2, mermaid, dot, toml, plantuml]
```

### Deprecated Methods

The following are deprecated and redirect to the Shape API:

- `f.IsTableFormat()` → `f.Supports(ShapeTable)`
- `f.IsTreeFormat()` → `f.Supports(ShapeTree)`
- `f.IsGraphFormat()` → `f.Supports(ShapeGraph)`
- `f.Category()` → `f.Shapes()`

## Data Structures

### TableData

Unified data structure for all tabular outputs (defined in `tabledata.go`):

```go
type TableData struct {
    Headers []string
    Rows    [][]string
}
```

### TableDataProvider

Interface for types that provide tabular data:

```go
type TableDataProvider interface {
    GetHeaders() []string
    GetRows() [][]string
}
```

### TreeNode

Hierarchical data structure for tree outputs (defined in `tree.go`):

```go
type TreeNode struct {
    ID       TreeNodeID
    Label    TreeNodeLabel
    Children []*TreeNode
    Metadata map[string]string
}
```

### GraphNode and GraphEdge

Data structures for graph/diagram outputs (defined in `graph.go`):

```go
type GraphNode struct {
    ID       GraphNodeID
    Label    GraphNodeLabel
    Shape    GraphShape
    Style    GraphStyle
    Metadata map[string]string
}

type GraphEdge struct {
    From  GraphNodeID
    To    GraphNodeID
    Label GraphNodeLabel
    Style EdgeStyle
}
```

## Interfaces

### Renderer

Base interface for all renderers:

```go
type Renderer interface {
    Render() (string, error)
}
```

### TableRenderer

For flat tabular data with SetHeaders/AddRow pattern:

```go
type TableRenderer interface {
    Renderer
    SetHeaders(headers []string)
    AddRow(row []string)
}
```

Note: MarkdownTable uses variadic `AddRow(row ...string) *MarkdownTable` and returns self for chaining, so it does not implement TableRenderer.

### TreeOutputRenderer

For hierarchical tree data:

```go
type TreeOutputRenderer interface {
    Renderer
    SetRoot(node *TreeNode)
}
```

### GraphRenderer

For diagram/graph data:

```go
type GraphRenderer interface {
    Renderer
    SetNodes(nodes []GraphNode)
    SetEdges(edges []GraphEdge)
}
```

### StreamingRenderer

For renderers that support streaming output:

```go
type StreamingRenderer interface {
    Renderer
    Stream(w io.Writer) error
}
```

## Implementation Strategy

1. **Shape capability matrix**: `formatCapabilities` map in `format.go` is the single source of truth
2. **Format Registry**: Opt-in runtime dispatch via `Create(format Format) (Renderer, error)`
3. **Adapter Pattern**: Each format implements its specific rendering
4. **Unified Data Model**: TableData works across all table-capable formats
5. **Tree-specific Model**: TreeNode for tree-capable formats
6. **Graph-specific Model**: GraphNode/GraphEdge for graph-capable formats
7. **Composition**: `GraphRendererMixin` (in `graph.go`) shared by DOT/Mermaid, `tableDataBase` (in `tabledata.go`) shared by HTML/Streaming

## Format-Specific Notes

### JSON and YAML

JSON and YAML declare support for all three shapes (Table, Tree, Graph) via `MarshalJSON`/`MarshalYAML`. These work as generic serialization functions — pass any data structure. For typed table rendering, use `MarshalJSONFromTableData`.

### D2

D2 supports both table data (via `D2FromTableData`) and graph data (via `D2FromTree` or `GraphNode`/`GraphEdge`). D2 has richer types than generic graph (shapes, arrows, SQL tables, classes, user journeys).

### HTML

HTML supports both table data (via `HTMLRenderer`) and tree data (via `HTMLTreeRenderer`). The `StreamingHTMLRenderer` provides true streaming for large datasets.

### XMLWriter

XMLWriter requires an `io.Writer` in its constructor. For string output, use with a `strings.Builder`:

```go
var buf strings.Builder
w := NewXMLWriter(&buf)
_ = w.WriteHeader([]string{"Name", "Value"})
_ = w.WriteRow([]string{"Alice", "30"})
_ = w.WriteFooter()
fmt.Print(buf.String())
```
