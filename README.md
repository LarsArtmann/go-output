# go-output

[![CI](https://github.com/larsartmann/go-output/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/go-output/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/go-output)](https://goreportcard.com/report/github.com/larsartmann/go-output)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library that formats structured data (tables, trees, graphs) into 12 different output formats with type-safe enums and zero-config color support.

## Purpose

A unified output formatting library for Go CLI applications — write your data once, render it in any of 12 formats.

## Quick Start

```go
import "github.com/larsartmann/go-output"

// JSON output
data, _ := output.MarshalJSONIndent(projects, "", "  ")
fmt.Println(string(data))

// Markdown table
md := output.NewMarkdownTable()
md.SetHeaders([]string{"Name", "Health", "Complexity"})
md.AddRow([]string{"Alpha", "90%", "7/10"})
out, err := md.Render()
if err != nil {
    log.Fatal(err)
}
fmt.Println(out)

// CSV output
w := output.NewCSVWriter(os.Stdout)
w.WriteHeader([]string{"Name", "Value"})
w.WriteRow([]string{"Item", "123"})
w.Flush()

// TSV output
tw := output.NewTSVWriter(os.Stdout)
tw.WriteHeader([]string{"Name", "Value"})
tw.WriteRow([]string{"Item", "123"})
tw.Flush()

// XML output
data, _ := output.MarshalXMLFromTableData(tableData)
fmt.Println(string(data))
```

## Format Categories

Formats are classified into three categories for programmatic filtering:

| Category  | Formats                                                        | Use Case                           |
| --------- | -------------------------------------------------------------- | ---------------------------------- |
| **Table** | `table`, `json`, `csv`, `tsv`, `xml`, `markdown`, `yaml`, `d2` | Tabular data with rows and columns |
| **Tree**  | `tree`, `html`                                                 | Hierarchical structures            |
| **Graph** | `d2`, `mermaid`, `dot`                                         | Network diagrams and flowcharts    |

```go
// Check format category
format, _ := output.ParseFormat("d2")
fmt.Println(format.IsTableFormat()) // true (D2 supports SQL tables)
fmt.Println(format.IsGraphFormat()) // true (D2 supports node-edge diagrams)
fmt.Println(format.Category())      // graph (graph takes precedence)

// Filter formats by category
for _, f := range output.AllFormats {
    if f.IsTableFormat() {
        fmt.Println(f) // table, json, csv, tsv, xml, markdown, yaml, d2
    }
}
```

## Streaming Renderer

For large datasets, use streaming output to minimize memory usage:

```go
// Streaming HTML output - writes incrementally
renderer := output.NewStreamingHTMLRenderer()
renderer.SetData(tableData)

// Stream directly to stdout
_ = renderer.Stream(os.Stdout)

// Wrap any renderer for streaming interface compliance
streamable := output.StreamingRendererFromRenderer(renderer)
_ = streamable.Stream(writer)
```

## Registry System

Register custom renderers for extensibility:

```go
// Register a custom format
err := output.Register(output.Format("custom"), func() output.Renderer {
    return &myCustomRenderer{}
})

// Create renderer by format
renderer, err := output.Create(output.FormatTable)

// Check what's registered
formats := output.RegisteredFormats()
isRegistered := output.IsRegistered(output.FormatJSON)
```

## Branded IDs

Type-safe identifiers prevent mixing different ID types:

```go
// D2 diagram nodes
nodeID := output.NewBrandedID[output.D2NodeIDBrand]("node-1")
nodeLabel := output.NewBrandedID[output.D2NodeLabelBrand]("My Node")

// Tree nodes
treeID := output.NewBrandedID[output.TreeNodeIDBrand]("root")
treeLabel := output.NewBrandedID[output.TreeNodeLabelBrand]("Root Node")

// Graph nodes
graphID := output.NewBrandedID[output.GraphNodeIDBrand]("vertex-a")
graphLabel := output.NewBrandedID[output.GraphNodeLabelBrand]("Vertex A")

// Generic branded ID
type ProjectIDBrand struct{}
projectID := output.NewBrandedID[ProjectIDBrand]("proj-123")

// Type safety - these won't compile:
// nodeID = treeID  // ERROR: cannot use treeID (type BrandedID[TreeNodeIDBrand]) as type BrandedID[D2NodeIDBrand]
```

## D2 Advanced Features

D2 diagrams support SQL tables, constraints, grid layouts, and nested containers:

```go
// SQL table with constraints
table := output.D2Table{
    Name: "users",
    Columns: []output.D2Column{
        {Name: "id", Type: "INT", Constraint: output.D2ConstraintPrimary},
        {Name: "email", Type: "VARCHAR(255)", Constraint: output.D2ConstraintUnique},
        {Name: "manager_id", Type: "INT", Constraint: output.D2ConstraintForeign},
    },
}

// Node with grid layout
node := output.D2Node{
    ID:          output.NewBrandedID[output.D2NodeIDBrand]("dashboard"),
    Label:       output.NewBrandedID[output.D2NodeLabelBrand]("Dashboard"),
    GridRows:    3,
    GridColumns: 2,
    GridGap:     8,
}

// Nested container
nestedNode := output.D2Node{
    ID:      output.NewBrandedID[output.D2NodeIDBrand]("container"),
    Nested:  "inner.nested.node",
}
```

## Supported Formats

### Table Formats

| Format     | Description                           | Package                            |
| ---------- | ------------------------------------- | ---------------------------------- |
| `table`    | Terminal tables with lipgloss styling | `github.com/larsartmann/go-output` |
| `json`     | JSON output with indentation          | `github.com/larsartmann/go-output` |
| `csv`      | CSV export with headers               | `github.com/larsartmann/go-output` |
| `tsv`      | TSV (Tab-Separated Values) export     | `github.com/larsartmann/go-output` |
| `xml`      | XML export with table structure       | `github.com/larsartmann/go-output` |
| `markdown` | Markdown tables                       | `github.com/larsartmann/go-output` |
| `yaml`     | YAML serialization                    | `github.com/larsartmann/go-output` |
| `d2`       | D2 diagram shapes                     | `github.com/larsartmann/go-output` |

### Tree Formats

| Format | Description                         | Package                            |
| ------ | ----------------------------------- | ---------------------------------- |
| `tree` | ASCII tree with box-drawing chars   | `github.com/larsartmann/go-output` |
| `html` | HTML tree with collapsible sections | `github.com/larsartmann/go-output` |

### Graph Formats

| Format    | Description                  | Package                            |
| --------- | ---------------------------- | ---------------------------------- |
| `d2`      | D2 diagram shapes            | `github.com/larsartmann/go-output` |
| `mermaid` | Mermaid flowchart diagrams   | `github.com/larsartmann/go-output` |
| `dot`     | DOT/Graphviz directed graphs | `github.com/larsartmann/go-output` |

## Supported Sort Options

| Option       | Description               |
| ------------ | ------------------------- |
| `name`       | Sort by name              |
| `importance` | Sort by importance level  |
| `created_at` | Sort by creation date     |
| `updated_at` | Sort by last update       |
| `health`     | Sort by health score      |
| `complexity` | Sort by complexity metric |

```go
import "github.com/larsartmann/go-output/sort"

type Project struct {
    Name       string
    Complexity int
}

items := []Project{
    {Name: "zebra", Complexity: 8},
    {Name: "apple", Complexity: 3},
}

// Sort by name using ByField for type-safe field comparison
sort.New(items, output.SortByName, false).
    WithLessFunc(sort.ByField(func(p Project) string { return p.Name })).
    Sort()
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
// Parse with validation
format, err := output.ParseFormat("json")
if err != nil {
    // handle error
}

// Check validity
if format.IsValid() {
    fmt.Println(format.String()) // "json"
}

// Get allowed values for CLI help
allowed := format.AllowedValues() // []string{"table", "json", "csv", ...}
```

## CLI Flag Integration

The `cmdguard/` subpackage provides helper types compatible with [cmdguard](https://github.com/larsartmann/cmdguard) for type-safe flags. Add cmdguard separately to your project:

```bash
go get github.com/larsartmann/cmdguard/v2
```

Example with cmdguard:

```go
import (
    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/cmdguard"
)

type ListFlags struct {
    Format output.OutputFormat `flag:"format" default:"table" help:"Output format (table, json, csv, tsv, xml, markdown, yaml, tree, html, d2, mermaid, dot)"`
    SortBy output.SortBy       `flag:"sort-by" default:"name" help:"Sort by (name, importance, created_at, updated_at, health, complexity)"`
    Color  output.ColorMode    `flag:"color" default:"auto" help:"Color mode (auto, always, never)"`
}
```

Flags validate against allowed values and provide bash/zsh completion.

## Installation

```bash
go get github.com/larsartmann/go-output
```

## Dependencies

Root module (zero lipgloss dependencies):

```go
require (
    github.com/go-faster/yaml v0.4.6
    golang.org/x/term v0.42.0
)
```

Terminal table module (install separately: `go get github.com/larsartmann/go-output/table`):

```go
require (
    charm.land/lipgloss/v2 v2.0.3
)
```

## Escape Functions

Safe escaping for various output formats:

| Function             | Purpose                 | Used By             |
| -------------------- | ----------------------- | ------------------- |
| `escape.HTML`        | HTML special characters | HTML, StreamingHTML |
| `escape.XML`         | XML special characters  | XML                 |
| `escape.D2`          | D2 diagram identifiers  | D2                  |
| `escape.DOT`         | DOT graph identifiers   | DOT                 |
| `escape.MermaidID`   | Mermaid node IDs        | Mermaid             |
| `escape.MermaidSlug` | Mermaid text (URL-safe) | Mermaid             |
| `escape.MermaidText` | Mermaid labels          | Mermaid             |

```go
import "github.com/larsartmann/go-output/escape"

// Escape HTML content
safe := escape.HTML("<script>alert('xss')</script>")
// Result: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"

// Escape D2 node ID
safeID := escape.D2("my-node.with.dots")
// Result: "my-node.with.dots" (dots preserved for D2 nesting)
```

## Development

```bash
# Build
just build

# Test (includes benchmarks and fuzz tests)
just test

# Lint
just lint

# Full verification
just verify

# Run example
go run ./examples/basic/main.go markdown

# Pre-commit hooks (install once)
pre-commit install
```

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a complete example demonstrating all formats.

## License

See [LICENSE](LICENSE) file for details.
