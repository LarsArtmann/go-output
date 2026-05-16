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
fmt.Println(format.IsTableFormat())      // true
fmt.Println(format.Category())           // table
```

## Why go-output?

- **12 formats, one API** — Same data, different renderers. No format-specific code paths.
- **Type-safe enums** — `Format`, `ColorMode`, `SortBy` — all validated at parse time, never raw strings.
- **Zero lipgloss in root module** — `go get go-output` pulls only `go-faster/yaml` and `x/term`. Lipgloss is isolated in the `table/` submodule.
- **Branded IDs** — Phantom types prevent mixing D2NodeID, TreeNodeID, GraphNodeID at compile time.
- **Streaming** — `StreamingHTMLRenderer` for large datasets with minimal memory.
- **Extensible registry** — Register custom renderers for runtime dispatch.

## Supported Formats

| Category  | Formats                                                        | Use Case                           |
| --------- | -------------------------------------------------------------- | ---------------------------------- |
| **Table** | `table`, `json`, `csv`, `tsv`, `xml`, `markdown`, `yaml`, `d2` | Tabular data with rows and columns |
| **Tree**  | `tree`, `html`                                                 | Hierarchical structures            |
| **Graph** | `d2`, `mermaid`, `dot`                                         | Network diagrams and flowcharts    |

All formats implement the `Renderer` interface:

```go
type Renderer interface {
    Render() (string, error)
}
```

### Table Formats

```go
// JSON
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

// YAML
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
tree := output.NewASCIITreeRenderer()

root := output.NewTreeNode("root", "Projects")
root.AddChild(output.NewTreeNode("alpha", "Alpha"))
root.AddChild(output.NewTreeNode("beta", "Beta"))

tree.SetRoot(root)
out, _ := tree.Render()
// Projects
// ├── Alpha
// └── Beta
```

### Graph Formats

```go
// DOT / Graphviz
renderer := output.DOTFromTableData(data)
out, _ := renderer.Render()

// Mermaid flowchart
renderer := output.MermaidFlowchartRenderer(data)
out, _ := renderer.Render()

// D2 diagrams (shapes, SQL tables, grid layouts, nested containers)
d2 := output.NewD2Renderer("Architecture")
d2.AddNode(output.D2Node{
    ID:    output.NewBrandedID[output.D2NodeIDBrand]("api"),
    Label: output.NewBrandedID[output.D2NodeLabelBrand]("API Gateway"),
    Shape: output.D2ShapeHexagon,
})
out, _ := d2.Render()
```

## Format Categories

Formats are classified into three categories for programmatic filtering:

```go
format, _ := output.ParseFormat("d2")
fmt.Println(format.IsTableFormat()) // true (D2 supports SQL tables)
fmt.Println(format.IsGraphFormat()) // true (D2 supports node-edge diagrams)
fmt.Println(format.Category())      // graph (graph takes precedence)

for _, f := range output.AllFormats {
    if f.IsTableFormat() {
        fmt.Println(f) // table, json, csv, tsv, xml, markdown, yaml, d2
    }
}
```

## Installation

```bash
go get github.com/larsartmann/go-output
```

For terminal table styling with lipgloss:

```bash
go get github.com/larsartmann/go-output/table
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
table := output.D2Table{
    Name: "users",
    Columns: []output.D2Column{
        {Name: "id", Type: "INT", Constraint: output.D2ConstraintPrimary},
        {Name: "email", Type: "VARCHAR(255)", Constraint: output.D2ConstraintUnique},
        {Name: "manager_id", Type: "INT", Constraint: output.D2ConstraintForeign},
    },
}

node := output.D2Node{
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

## CLI Flag Integration

The `cmdguard/` subpackage provides types compatible with [cmdguard](https://github.com/larsartmann/cmdguard) for type-safe CLI flags:

```go
import (
    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/cmdguard"
)

type ListFlags struct {
    Format output.OutputFormat `flag:"format" default:"table" help:"Output format (table, json, csv, ...)"`
    SortBy output.SortBy       `flag:"sort-by" default:"name" help:"Sort field"`
    Color  output.ColorMode    `flag:"color" default:"auto" help:"Color mode (auto, always, never)"`
}
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

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a complete example demonstrating all 12 formats:

```bash
go run ./examples/basic/main.go markdown
```

## Development

```bash
go build ./...                  # Build all workspace modules
go test ./...                   # Test all modules
go test -race ./...             # Race detector
go test -cover ./...            # Coverage report
golangci-lint run --fix ./...   # Lint
```

## License

[MIT](LICENSE)
