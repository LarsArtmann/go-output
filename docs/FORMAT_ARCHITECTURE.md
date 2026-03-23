# go-output Format Architecture

## Overview

This document describes the extensible format architecture for go-output, supporting multiple output formats in a unified way.

## Format Categories

### 1. Table Formats (Flat Data)
- `table` - Terminal tables with lipgloss styling
- `json` - Formatted JSON
- `csv` - CSV with headers
- `markdown` - Markdown tables
- `yaml` - YAML output

### 2. Tree Formats (Hierarchical Data)
- `tree` - ASCII tree representation
- `html` - HTML with collapsible tree

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

### TreeNode
Hierarchical data structure for tree outputs:
```go
type TreeNode struct {
    ID       string
    Label    string
    Children []*TreeNode
    Metadata map[string]string
}
```

## Interfaces

### Renderer
Base interface for all renderers:
```go
type Renderer interface {
    Render() string
}
```

### TableRenderer
For flat tabular data:
```go
type TableRenderer interface {
    Renderer
    SetData(data TableData)
}
```

### TreeRenderer
For hierarchical tree data:
```go
type TreeRenderer interface {
    Renderer
    SetTree(root *TreeNode)
}
```

### GraphRenderer
For diagram/graph data:
```go
type GraphRenderer interface {
    Renderer
    SetGraph(nodes []GraphNode, edges []GraphEdge)
}
```

## Implementation Strategy

1. **Format Registry**: Central place to get format renderers
2. **Adapter Pattern**: Each format implements its specific rendering
3. **Unified Data Model**: TableData works across all table formats
4. **Tree-specific Model**: TreeNode for tree/graph formats
