# go-output

[![CI](https://github.com/larsartmann/go-output/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/go-output/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/go-output)](https://goreportcard.com/report/github.com/larsartmann/go-output)
[![GoDoc](https://godoc.org/github.com/larsartmann/go-output?status.svg)](https://godoc.org/github.com/larsartmann/go-output)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library that formats structured data into **12 output formats** — tables, trees, and diagrams — with type-safe enums, branded IDs, and zero-config color support. Write your data once, render it anywhere.

```go
import "github.com/larsartmann/go-output"
```

## Quick Start

Build tabular data once, render it in any format:

```go
data := output.NewTableData([]string{"Name", "Health", "Complexity"})
data.AddRow([]string{"Alpha", "90%", "7/10"})
data.AddRow([]string{"Beta", "75%", "5/10"})

// Markdown table
md := output.NewMarkdownTable()
md.SetHeaders(data.GetHeaders())
for _, row := range data.GetRows() {
    md.AddRow(row)
}
out, _ := md.Render()

// JSON
json, _ := output.MarshalJSONIndent(projects, "", "  ")

// CSV
w := output.NewCSVWriter(os.Stdout)
w.WriteHeader(data.GetHeaders())
for _, row := range data.GetRows() {
    w.WriteRow(row)
}
w.Flush()
```

Use the `Format` enum for runtime format selection — perfect for CLI flags:

```go
format, _ := output.ParseFormat("json") // validates input
fmt.Println(format.Supports(output.ShapeTable)) // true
fmt.Println(format.Shapes())                     // [table tree graph]
```

## Why go-output?

- **12 formats, one API** — Same data, different renderers. No format-specific code paths.
- **Type-safe enums** — `Format`, `ColorMode`, `SortBy` — all validated at parse time, never raw strings.
- **Zero heavy deps in root module** — `go get go-output` pulls only `go-faster/yaml` and `x/term`. Lipgloss is isolated in `table/`, D2 and graph renderers in their own modules.
- **Branded IDs** — Phantom types prevent mixing D2NodeID, TreeNodeID, GraphNodeID at compile time.
- **Streaming** — `StreamingHTMLRenderer` for large datasets with minimal memory.
- **Extensible registry** — Register custom renderers for runtime dispatch.

## Supported Formats

| Format     | Table | Tree | Graph | Notes                                                            |
| ---------- | :---: | :--: | :---: | ---------------------------------------------------------------- |
| `table`    |  ✅   |      |       | Terminal tables with lipgloss styling (separate `table/` module) |
| `json`     |  ✅   |  ✅  |  ✅   | Shape-agnostic serialization                                     |
| `csv`      |  ✅   |      |       | Comma-separated export                                           |
| `tsv`      |  ✅   |      |       | Tab-separated export                                             |
| `xml`      |  ✅   |      |       | XML with table structure                                         |
| `markdown` |  ✅   |      |       | Markdown tables                                                  |
| `yaml`     |  ✅   |  ✅  |  ✅   | Shape-agnostic serialization                                     |
| `d2`       |  ✅   |      |  ✅   | SQL tables + node-edge diagrams (separate `d2/` module)          |
| `html`     |  ✅   |  ✅  |       | HTML tables + collapsible tree                                   |
| `tree`     |       |  ✅  |       | ASCII tree with box-drawing chars                                |
| `mermaid`  |  ✅   |      |  ✅   | Mermaid flowchart diagrams (separate `graph/` module)            |
| `dot`      |  ✅   |      |  ✅   | DOT/Graphviz directed graphs (separate `graph/` module)          |

All formats implement the `Renderer` interface:

```go
type Renderer interface {
    Render() (string, error)
}
```

### Table Formats

```go
// JSON table (array of objects)
jt := output.NewJSONTableRenderer()
jt.SetHeaders([]string{"Name", "Health"})
jt.AddRow([]string{"Alpha", "90%"})
out, _ := jt.Render()
// [{"Name": "Alpha", "Health": "90%"}]

// YAML table (sequence of mappings)
yt := output.NewYAMLTableRenderer()
yt.SetHeaders([]string{"Name", "Health"})
yt.AddRow([]string{"Alpha", "90%"})
out, _ := yt.Render()

// JSON (any data)
data, _ := output.MarshalJSONIndent(projects, "", "  ")

// CSV
w := output.NewCSVWriter(os.Stdout)
w.WriteHeader([]string{"Name", "Value"})
w.WriteRow([]string{"Item", "123"})
w.Flush()

// TSV
tw := output.NewTSVWriter(os.Stdout)
tw.WriteHeader([]string{"Name", "Value"})
tw.WriteRow([]string{"Item", "123"})
tw.Flush()

// XML
data, _ := output.MarshalXMLFromTableData(tableData)

// YAML (any data)
data, _ := output.MarshalYAML(projects)

// Markdown table
md := output.NewMarkdownTable()
md.SetHeaders([]string{"Name", "Health"})
md.AddRow([]string{"Alpha", "90%"})
out, _ := md.Render()

// Terminal table with lipgloss styling (requires go-output/table)
tbl := table.New()
tbl.SetHeaders("Name", "Health")
tbl.AddRow("Alpha", "90%")
out, _ := tbl.Render()
```

### Tree Formats

```go
root := output.NewTreeNode("root", "Projects")
root.AddChild(output.NewTreeNode("alpha", "Alpha"))
root.AddChild(output.NewTreeNode("beta", "Beta"))

// ASCII tree
tree := output.NewASCIITreeRenderer()
tree.SetRoot(root)
out, _ := tree.Render()
// Projects
// ├── Alpha
// └── Beta

// JSON tree
jt := output.NewJSONTreeRenderer()
jt.SetRoot(root)
out, _ := jt.Render()
// {"id": "root", "label": "Projects", "children": [...]}

// YAML tree
yt := output.NewYAMLTreeRenderer()
yt.SetRoot(root)
out, _ := yt.Render()
// id: root
// label: Projects
// children: ...

// HTML tree (collapsible)
ht := output.NewHTMLTreeRenderer()
ht.SetRoot(root)
out, _ := ht.Render()
```

### Graph Formats

```go
import (
    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/d2"
    "github.com/larsartmann/go-output/graph"
)

nodes := []output.GraphNode{
    output.NewGraphNode("a", "API Gateway"),
    output.NewGraphNode("b", "Backend"),
}
edges := []output.GraphEdge{
    output.NewGraphEdge("a", "b"),
}

// DOT / Graphviz (requires go-output/graph)
renderer := graph.DOTFromTableData(data)
out, _ := renderer.Render()

// Mermaid flowchart (requires go-output/graph)
renderer := graph.MermaidFromTableData(data)
out, _ := renderer.Render()

// JSON graph
jg := output.NewJSONGraphRenderer()
jg.SetNodes(nodes)
jg.SetEdges(edges)
out, _ := jg.Render()
// {"nodes": [...], "edges": [...]}

// YAML graph
yg := output.NewYAMLGraphRenderer()
yg.SetNodes(nodes)
yg.SetEdges(edges)
out, _ := yg.Render()

// D2 diagrams (requires go-output/d2)
diagram := d2.NewD2Diagram().
    AddNodeWithShape("api", "API Gateway", d2.D2ShapeHexagon).
    AddEdgeSimple("api", "backend")
out, _ := diagram.Render()
```

## Data Shapes

Every format declares which data shapes it supports via the capability matrix:

```go
// Check if a format supports a specific data shape
format, _ := output.ParseFormat("d2")
fmt.Println(format.Supports(output.ShapeTable)) // true (D2 supports SQL tables)
fmt.Println(format.Supports(output.ShapeGraph)) // true (D2 supports node-edge diagrams)
fmt.Println(format.Supports(output.ShapeTree))  // false

// Get all shapes a format supports
for _, shape := range format.Shapes() {
    fmt.Println(shape) // "table", "graph"
}

// Find all formats that can render graph data
for _, f := range output.FormatsForShape(output.ShapeGraph) {
    fmt.Println(f) // json, yaml, d2, mermaid, dot
}
```

## Installation

```bash
go get github.com/larsartmann/go-output
```

Sub-modules for specific formats:

```bash
go get github.com/larsartmann/go-output/table   # Terminal tables with lipgloss
go get github.com/larsartmann/go-output/d2     # D2 diagrams
go get github.com/larsartmann/go-output/graph  # DOT + Mermaid renderers
```

## Branded IDs

Type-safe identifiers prevent mixing different ID types at compile time:

```go
nodeID := output.NewBrandedID[output.D2NodeIDBrand]("node-1")
treeID := output.NewBrandedID[output.TreeNodeIDBrand]("root")

// nodeID = treeID  // COMPILE ERROR: different branded types
```

Define your own branded types:

```go
type ProjectIDBrand struct{}
projectID := output.NewBrandedID[ProjectIDBrand]("proj-123")
```

## D2 Advanced Features

D2 diagrams support SQL tables, constraints, grid layouts, and nested containers:

```go
table := d2.D2Table{
    Name: "users",
    Columns: []d2.D2Column{
        {Name: "id", Type: "INT", Constraint: d2.D2ConstraintPrimary},
        {Name: "email", Type: "VARCHAR(255)", Constraint: d2.D2ConstraintUnique},
        {Name: "manager_id", Type: "INT", Constraint: d2.D2ConstraintForeign},
    },
}

node := d2.D2Node{
    ID:          output.NewBrandedID[output.D2NodeIDBrand]("dashboard"),
    Label:       output.NewBrandedID[output.D2NodeLabelBrand]("Dashboard"),
    GridRows:    3,
    GridColumns: 2,
    GridGap:     8,
}
```

## Registry System

Register custom renderers for runtime dispatch:

```go
output.Register(output.Format("custom"), func() output.Renderer {
    return &myCustomRenderer{}
})

renderer, _ := output.Create(output.FormatTable)
formats := output.RegisteredFormats()
```

## Streaming Renderer

For large datasets, stream output incrementally:

```go
renderer := output.NewStreamingHTMLRenderer()
renderer.SetData(tableData)
_ = renderer.Stream(os.Stdout)
```

## Color Modes

| Mode     | Description                                    |
| -------- | ---------------------------------------------- |
| `auto`   | Respect `NO_COLOR`, CI env vars, TTY detection |
| `always` | Force ANSI colors                              |
| `never`  | Disable colors                                 |

## Type-Safe Enums

All configuration types provide validation and string conversion:

```go
format, err := output.ParseFormat("json")
if format.IsValid() {
    fmt.Println(format.String()) // "json"
}
allowed := format.AllowedValues() // []string{"table", "json", "csv", ...}
```

## Escape Functions

The `escape/` subpackage provides safe escaping for each format:

```go
import "github.com/larsartmann/go-output/escape"

safe := escape.HTML("<script>alert('xss')</script>")
safeID := escape.D2("my-node.with.dots")
```

| Function           | Purpose                 |
| ------------------ | ----------------------- |
| `escape.HTML`      | HTML special characters |
| `escape.XML`       | XML special characters  |
| `escape.D2`        | D2 diagram identifiers  |
| `escape.DOT`       | DOT graph identifiers   |
| `escape.MermaidID` | Mermaid node IDs        |

## Dependencies

Root module — zero lipgloss dependencies:

```go
require (
    github.com/go-faster/yaml v0.4.6
    golang.org/x/term v0.42.0
)
```

Terminal table module (install separately):

```go
require (
    charm.land/lipgloss/v2 v2.0.3
)
```

D2 diagram module (install separately):

```go
require github.com/larsartmann/go-output/d2 v0.0.0
```

DOT + Mermaid graph module (install separately):

```go
require github.com/larsartmann/go-output/graph v0.0.0
```

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a complete example demonstrating all 12 formats:

```bash
go run ./examples/basic/main.go markdown
```

## Development

### Nix (recommended)

```bash
nix develop                    # Enter dev shell (Go 1.26, golangci-lint, gopls)
nix fmt                        # Format .nix files
nix flake check                # Verify formatting + pre-commit hooks
```

### Go toolchain (manual)

```bash
go build ./...                  # Build all workspace modules
go test ./...                   # Test all modules
go test -race ./...             # Race detector
go test -cover ./...            # Coverage report
golangci-lint run --fix ./...   # Lint
```

## API Stability

This library is pre-v1. The following guarantees apply:

- **Root module** (`github.com/larsartmann/go-output`): Public API is stable. Breaking changes will be documented in CHANGELOG.md.
- **Sub-modules** (`d2`, `graph`, `table`): May evolve independently. Import them explicitly to opt in.
- **`Renderer` interface**: Stable — all formats implement `Render() (string, error)`.
- **`internal/` packages**: No stability guarantee. Do not import these.

## License

[MIT](LICENSE)
