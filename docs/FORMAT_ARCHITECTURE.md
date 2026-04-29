# go-output Format Architecture

## Overview

This document describes the extensible format architecture for go-output, supporting multiple output formats in a unified way.

## Format Categories

### 1. Table Formats (Flat Data)

- `table` - Terminal tables with lipgloss styling
- `json` - Formatted JSON
- `csv` - CSV with headers
- `tsv` - TSV (Tab-Separated Values) with headers
- `xml` - XML with headers and rows
- `markdown` - Markdown tables
- `yaml` - YAML output

### 2. Tree Formats (Hierarchical Data)

- `tree` - ASCII tree representation
- `html` - HTML with nested lists (see HTMLTreeRenderer)

### 3. Graph Formats (Network/Diagram Data)

- `d2` - D2 diagram shapes
- `dot` - DOT/Graphviz format
- `mermaid` - Mermaid flowchart syntax

## Data Structures

### TableData

Unified data structure for all tabular outputs:

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

Hierarchical data structure for tree outputs:

```go
type TreeNode struct {
    ID       TreeNodeID
    Label    TreeNodeLabel
    Children []*TreeNode
    Metadata map[string]string
}
```

### GraphNode and GraphEdge

Data structures for graph/diagram outputs:

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

1. **Format Registry**: Central place to get format renderers via `GetRenderer(format Format) (Renderer, error)`
2. **Adapter Pattern**: Each format implements its specific rendering
3. **Unified Data Model**: TableData works across all table formats
4. **Tree-specific Model**: TreeNode for tree/graph formats
5. **Graph-specific Model**: GraphNode/GraphEdge for diagram formats

## Format-Specific Notes

### D2

D2 diagrams support both flat data (via `D2FromTableData`) and hierarchical data (via `D2FromTree`). D2 is categorized as a Graph format but can render table-like data.

### XMLWriter

XMLWriter requires an `io.Writer` in its constructor (v2.0.0+). For string output, use with a `strings.Builder`:

```go
var buf strings.Builder
w := NewXMLWriter(&buf)
_ = w.WriteHeader([]string{"Name", "Value"})
_ = w.WriteRow([]string{"Alice", "30"})
_ = w.WriteFooter()
fmt.Print(buf.String())
```
