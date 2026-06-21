# go-output

[![CI](https://github.com/larsartmann/go-output/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/go-output/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/go-output)](https://goreportcard.com/report/github.com/larsartmann/go-output)
[![GoDoc](https://godoc.org/github.com/larsartmann/go-output?status.svg)](https://godoc.org/github.com/larsartmann/go-output)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Write your data once. Render it anywhere.**
>
> One library. Sixteen formats. Three data shapes. Zero lock-in.

`go-output` is a Go library that turns your structured data into **16 output formats** — tables, trees, and diagrams — with type-safe enums, branded IDs, and zero-config color support. It also includes **NOM-style real-time progress visualization** for long-running workflows, inspired by [`nix-output-monitor`](https://github.com/maralorn/nix-output-monitor).

```go
import "github.com/larsartmann/go-output"
```

The root module has **zero heavy dependencies** — only `golang.org/x/term`. YAML, lipgloss, bubbletea, and diagram renderers live in isolated sub-modules you import only when you need them.

---

## Quick Start

Build tabular data once, render it in any format:

```go
import (
    "os"

    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/delimited"
    "github.com/larsartmann/go-output/markdown"
    "github.com/larsartmann/go-output/serialization"
)

// Define your data once
data := output.NewTableData([]string{"Name", "Health", "Complexity"})
data.AddRow([]string{"Alpha", "90%", "7/10"})
data.AddRow([]string{"Beta", "75%", "5/10"})
data.SetFooter([]string{"Total", "2", "-"})

// Render to Markdown
md := markdown.NewMarkdownTable()
md.SetHeaders(data.GetHeaders())
for _, row := range data.GetRows() { md.AddRow(row) }
out, _ := md.Render()

// Or serialize to JSON
jtr := serialization.NewJSONTableRenderer()
jtr.SetData(data)
out, _ = jtr.Render()

// Or stream to CSV
csv := delimited.NewCSVWriter(os.Stdout)
_ = csv.WriteHeader(data.GetHeaders())
for _, row := range data.GetRows() { _ = csv.WriteRow(row) }
csv.Flush()
```

Use the `Format` enum for runtime format selection — perfect for CLI flags:

```go
format, _ := output.ParseFormat("json") // validates input
fmt.Println(format.Supports(output.ShapeTable)) // true
fmt.Println(format.Shapes())                     // [table tree graph]
```

Or dispatch through the unified renderer:

```go
output.RenderTableData(data, output.FormatHTML, output.RenderOptions{
    ColorMode: output.ColorModeAuto,
})
```

---

## Why go-output?

- **16 formats, one API** — Same `TableData`, `TreeNode`, or `GraphNode`. No format-specific code paths.
- **Type-safe everything** — `Format`, `ColorMode`, `ActivityStatus` — all validated at parse time. Branded IDs prevent mixing `D2NodeID` with `TreeNodeID` at compile time.
- **Zero heavy deps in root** — `go get go-output` pulls only `x/term`. YAML, lipgloss, bubbletea, and diagram renderers are opt-in sub-modules.
- **NOM real-time progress** — Dependency trees, activity counts, timing estimates, and inline terminal rendering. O(1) summary bars even at 10,000 activities.
- **Streaming for large data** — `StreamingHTMLRenderer` writes incrementally with minimal memory.
- **Zero-config color** — `ColorModeAuto` detects TTY, respects `NO_COLOR`, `CI`, `FORCE_COLOR`.
- **Pre-v1 API frozen** — ADR 006 locks core interfaces. Breaking changes are documented.

---

## Installation

```bash
go get github.com/larsartmann/go-output
```

Sub-modules for specific formats (import only what you need):

```bash
# Data renderers
go get github.com/larsartmann/go-output/delimited       # CSV + TSV writers
go get github.com/larsartmann/go-output/serialization   # JSON + YAML + TOML + JSONL marshaling
go get github.com/larsartmann/go-output/markup          # XML + HTML + AsciiDoc + Streaming HTML
go get github.com/larsartmann/go-output/table           # Terminal tables with lipgloss
go get github.com/larsartmann/go-output/markdown        # Markdown table renderer
go get github.com/larsartmann/go-output/tree            # ASCII tree renderer

# Diagram renderers
go get github.com/larsartmann/go-output/d2              # D2 diagrams (SQL tables, shapes, grids)
go get github.com/larsartmann/go-output/graph           # DOT + Mermaid flowcharts
go get github.com/larsartmann/go-output/plantuml        # PlantUML component diagrams

# Progress visualization
go get github.com/larsartmann/go-output/nom             # NOM-style real-time progress
go get github.com/larsartmann/go-output/tui             # Bubble Tea interactive TUI (depends on nom)

# Utilities (zero-dep, used internally)
go get github.com/larsartmann/go-output/enum            # Generic enum utilities
go get github.com/larsartmann/go-output/escape          # Format-specific escaping
go get github.com/larsartmann/go-output/envdetect       # CI / NO_COLOR env detection
go get github.com/larsartmann/go-output/testhelpers     # Shared test assertions
```

---

## Supported Formats

| Format     | Table | Tree | Graph | Module        | Notes                                                   |
| ---------- | :---: | :--: | :---: | ------------- | ------------------------------------------------------- |
| `json`     |  ✅   |  ✅  |  ✅   | root          | Shape-agnostic serialization                            |
| `yaml`     |  ✅   |  ✅  |  ✅   | serialization | Shape-agnostic serialization                            |
| `toml`     |  ✅   |  ✅  |  ✅   | serialization | Shape-agnostic serialization                            |
| `csv`      |  ✅   |      |       | delimited     | Streaming writer with auto-quoting                      |
| `tsv`      |  ✅   |      |       | delimited     | Tab-separated with type-switch marshaling               |
| `jsonl`    |  ✅   |      |       | serialization | One JSON object per line                                |
| `xml`      |  ✅   |      |       | markup        | Structured `<table>` with XML escaping                  |
| `html`     |  ✅   |  ✅  |       | markup        | Styled tables + collapsible trees                       |
| `asciidoc` |  ✅   |      |       | markup        | `\|===` borders with pipe escaping                      |
| `markdown` |  ✅   |      |       | root          | Auto column widths, alignment, bold headers             |
| `table`    |  ✅   |      |       | table         | Lipgloss terminal tables with rounded borders           |
| `tree`     |       |  ✅  |       | root          | ASCII box-drawing (`├──`, `└──`) with color cycling     |
| `d2`       |  ✅   |      |  ✅   | d2            | SQL tables, 20 node shapes, grid layouts, style classes |
| `mermaid`  |  ✅   |      |  ✅   | graph         | Flowcharts with 8 node shapes                           |
| `dot`      |  ✅   |      |  ✅   | graph         | Graphviz directed graphs                                |
| `plantuml` |  ✅   |      |  ✅   | plantuml      | Component diagrams with TableData→graph conversion      |

### Data Shape Capabilities

```go
// Check what a format can do
format, _ := output.ParseFormat("d2")
format.Supports(output.ShapeTable) // true (SQL tables)
format.Supports(output.ShapeGraph) // true (node-edge diagrams)
format.Supports(output.ShapeTree)  // false

// Find all formats for a shape
for _, f := range output.FormatsForShape(output.ShapeGraph) {
    fmt.Println(f) // json, yaml, d2, mermaid, dot
}
```

---

## Table Formats

```go
// Terminal table with lipgloss styling (requires go-output/table)
tbl := table.New()
tbl.SetHeaders("Name", "Health")
tbl.AddRow("Alpha", "90%")
out, _ := tbl.Render()

// Markdown table (requires go-output/markdown)
md := markdown.NewMarkdownTable()
md.SetHeaders([]string{"Name", "Health"})
md.AddRow([]string{"Alpha", "90%"})
out, _ = md.Render()

// JSON table (array of objects — requires go-output/serialization)
jtr := serialization.NewJSONTableRenderer()
jtr.SetHeaders([]string{"Name", "Health"})
jtr.AddRow([]string{"Alpha", "90%"})
out, _ = jtr.Render()
// [{"Name": "Alpha", "Health": "90%"}]

// CSV streaming writer (requires go-output/delimited)
csv := delimited.NewCSVWriter(os.Stdout)
_ = csv.WriteHeader([]string{"Name", "Value"})
_ = csv.WriteRow([]string{"Item", "123"})
csv.Flush()

// HTML table with footer (requires go-output/markup)
ht := markup.NewHTMLRenderer()
ht.SetHeaders([]string{"Name", "Value"})
ht.AddRow([]string{"Item", "123"})
ht.SetFooter([]string{"Total", "1"})
out, _ = ht.Render()
```

### Footer Row Support

Set `TableData.Footer` for an optional totals/summary row:

```go
data := output.NewTableData([]string{"Name", "Score"})
data.AddRow([]string{"Alice", "95"})
data.AddRow([]string{"Bob", "87"})
data.SetFooter([]string{"Total", "182"})
```

| Format                                         | Footer | Behavior                                   |
| ---------------------------------------------- | :----: | ------------------------------------------ |
| `table`                                        |   ✅   | Bold-styled footer row                     |
| `markdown`                                     |   ✅   | Second separator + bold footer row         |
| `csv` / `tsv`                                  |   ✅   | Appended as last data row                  |
| `html`                                         |   ✅   | `<tfoot>` section with `footer-cell` class |
| `xml`                                          |   ✅   | `<footer>` element                         |
| `asciidoc`                                     |   ✅   | Footer row cells                           |
| `json` / `yaml` / `toml` / `jsonl`             |   ❌   | Data serialization — footer skipped        |
| `tree` / `d2` / `mermaid` / `dot` / `plantuml` |   ❌   | Non-tabular formats                        |

---

## Tree Formats

```go
root := output.NewTreeNode("root", "Projects")
root.AddChild(output.NewTreeNode("alpha", "Alpha"))
root.AddChild(output.NewTreeNode("beta", "Beta"))

// ASCII tree with box-drawing chars (requires go-output/tree)
tree := tree.NewASCIITreeRenderer()
tree.SetRoot(root)
tree.SetColorMode(output.ColorModeAuto)
out, _ := tree.Render()
// Projects
// ├── Alpha
// └── Beta

// JSON tree (requires go-output/serialization)
jtr := serialization.NewJSONTreeRenderer()
jtr.SetRoot(root)
out, _ = jtr.Render()

// HTML collapsible tree (requires go-output/markup)
ht := markup.NewHTMLTreeRenderer()
ht.SetRoot(root)
out, _ = ht.Render()
```

---

## Graph / Diagram Formats

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
dot := graph.NewDOTRenderer()
dot.SetNodes(nodes)
dot.SetEdges(edges)
out, _ := dot.Render()

// Mermaid flowchart (requires go-output/graph)
mmd := graph.NewMermaidRenderer()
mmd.SetNodes(nodes)
mmd.SetEdges(edges)
out, _ = mmd.Render()

// D2 diagrams with SQL tables (requires go-output/d2)
diagram := d2.NewD2Diagram().
    AddNodeWithShape("api", "API Gateway", d2.D2ShapeHexagon).
    AddEdgeSimple("api", "backend")
out, _ = diagram.Render()
```

### D2 Advanced Features

SQL tables, constraints, grid layouts, and nested containers:

```go
table := d2.D2Table{
    Name: "users",
    Columns: []d2.D2Column{
        {Name: "id", Type: "INT", Constraint: d2.D2ConstraintPrimary},
        {Name: "email", Type: "VARCHAR(255)", Constraint: d2.D2ConstraintUnique},
    },
}

node := d2.D2Node{
    ID:          output.NewBrandedID[output.D2NodeIDBrand]("dashboard"),
    Label:       output.NewBrandedID[output.D2NodeLabelBrand]("Dashboard"),
    GridRows:    3,
    GridColumns: 2,
}
```

### Cross-Shape Conversion

Convert between data shapes without rewriting your data:

```go
// TableData → Graph (auto-generates edges between consecutive rows)
dot := graph.DOTFromTableData(data)
mmd := graph.MermaidFromTableData(data)
plantuml := plantuml.PlantUMLFromTableData(data)
d2Diagram := d2.D2FromTableData(data)

// TableData → Tree (hierarchical from tabular data; requires go-output/tree)
tree := tree.TreeRendererFromTableData(data)

// Tree → Graph
d2Diagram := d2.D2FromTree(root)
dot := graph.DOTFromTree(root)
mmd := graph.MermaidFromTree(root)
plantuml := plantuml.PlantUMLFromTree(root)
```

---

## NOM Real-Time Progress

Track and visualize long-running workflows with dependency trees, inspired by [`nix-output-monitor`](https://github.com/maralorn/nix-output-monitor).

```go
import (
    "context"
    "os"
    "time"

    "github.com/larsartmann/go-output/nom"
)

// Create subscriber and renderer
sub := nom.NewNOMStyleSubscriber()
renderer := nom.NewInlineRenderer(sub, os.Stdout, 20) // maxHeight 20 lines

// Start the render loop
ctx := context.Background()
renderer.Start(ctx, 100*time.Millisecond)
defer renderer.Finish(nil)

// Fire workflow events
sub.OnEvent(ctx, nom.WorkflowStarted{
    ID:   nom.NewWorkflowID("build"),
    Name: nom.NewWorkflowName("Build Project"),
})
sub.OnEvent(ctx, nom.ActivityStarted{
    ID:   nom.NewActivityID("compile"),
    Name: nom.NewActivityName("Compile"),
    Deps: []nom.ActivityID{},
})
sub.OnEvent(ctx, nom.ActivityStarted{
    ID:   nom.NewActivityID("test"),
    Name: nom.NewActivityName("Run Tests"),
    Deps: []nom.ActivityID{nom.NewActivityID("compile")},
})

// Mark activities complete (Name keys the timing cache for future ETAs)
sub.OnEvent(ctx, nom.ActivityCompleted{
    ID:       nom.NewActivityID("compile"),
    Name:     nom.NewActivityName("Compile"),
    Duration: 5 * time.Second,
})
sub.OnEvent(ctx, nom.ActivityCompleted{
    ID:       nom.NewActivityID("test"),
    Name:     nom.NewActivityName("Run Tests"),
    Duration: 12 * time.Second,
})
sub.OnEvent(ctx, nom.WorkflowCompleted{ID: nom.NewWorkflowID("build")})
```

### NOM Features

- **Dependency trees** — Hierarchical parent/child relationships with UTF-8 box-drawing
- **O(1) activity counts** — Summary bar updates in constant time, even with 10,000+ activities
- **Timing cache** — Persists duration history to `~/.cache/nom-timing.csv` for ETA estimates
- **Snapshot-based rendering** — Race-free: renderers read immutable value copies, not shared pointers
- **CI-safe degradation** — Auto-detects CI environments; appends frames line-by-line instead of ANSI cursor codes
- **Height-pressure collapse** — When the tree exceeds `maxHeight`, completed children collapse with a `⋯ N completed` marker
- **Per-activity progress** — Download bars (`▕████░░░░▏ 45%`) and host tags (`@host`)
- **Node classes** — Root/twig/leaf styling with bold roots for top-level visibility

### NOM Diagram Export

Export live NOM state as DOT or Mermaid diagrams:

```go
// Project subscriber state into graph nodes/edges
reader := sub.Store()
nodes := reader.Nodes()
edges := reader.Edges()

// Render as DOT
dot := graph.NewDOTRenderer()
dot.SetNodes(nodes)
dot.SetEdges(edges)
out, _ := dot.Render()
```

---

## Bubble Tea TUI

Full-screen interactive TUI built on Bubble Tea v2, with two display modes:

```go
import (
    "github.com/larsartmann/go-output/nom"
    "github.com/larsartmann/go-output/tui"
)

reporter := tui.NewBubbleTeaProgressReporter()
reporter.SetDisplayMode(tui.DisplayModeNOM) // or DisplayModeUniversal
reporter.SetCancelFunc(cancelFunc)

// Report progress
reporter.ReportStep(1, 5, "Building...")
reporter.ReportProgress(0.35)
reporter.ReportMessage("Compiling module Alpha")

// Or drive via NOM events
reporter.Subscriber().OnEvent(ctx, nom.ActivityStarted{
    ID:   nom.NewActivityID("a1"),
    Name: nom.NewActivityName("Build"),
})

reporter.Start()
defer reporter.Stop()
```

### TUI Controls

| Key          | Action                              |
| ------------ | ----------------------------------- |
| `j` / `↓`    | Scroll down                         |
| `k` / `↑`    | Scroll up                           |
| `pgdown`     | Scroll half page down               |
| `pgup`       | Scroll half page up                 |
| `g` / `Home` | Jump to top                         |
| `G` / `End`  | Jump to bottom                      |
| `?`          | Toggle help overlay                 |
| `q`          | Quit (workflow continues)           |
| `ctrl+c`     | Cancel workflow and quit            |
| **Mouse**    | Wheel scroll, click to select nodes |

---

## Type-Safe Enums

All configuration types provide validation and conversion:

```go
format, err := output.ParseFormat("json")
if format.IsValid() {
    fmt.Println(format.String())      // "json"
    fmt.Println(format.Shapes())     // [table tree graph]
}
allowed := format.AllowedValues() // ["table", "json", "csv", ...]
```

Available enums: `Format` (16 values), `Shape` (3 values), `ColorMode` (auto/always/never), `NodeShape` (8 shapes), `D2NodeShape` (20 shapes), `D2ArrowType` (11 types), `D2Constraint` (3 constraints), `Alignment` (left/right/center).

---

## Branded IDs

Phantom types prevent mixing different ID types at compile time:

```go
nodeID := output.NewBrandedID[output.D2NodeIDBrand]("node-1")
treeID := output.NewBrandedID[output.TreeNodeIDBrand]("root")

// nodeID = treeID  // COMPILE ERROR: different branded types
```

Define your own:

```go
type ProjectIDBrand struct{}
projectID := output.NewBrandedID[ProjectIDBrand]("proj-123")
```

---

## Color Modes

All terminal renderers support `ColorMode` for controlling ANSI output:

```go
// Terminal table: functional option
tbl := table.New(table.WithColorMode(output.ColorModeAlways))

// ASCII tree: setter
tree := tree.NewASCIITreeRenderer()
tree.SetColorMode(output.ColorModeAlways)

// Markdown: chaining setter
md := markdown.NewMarkdownTable().SetColorMode(output.ColorModeAlways)

// Unified dispatch: RenderOptions
output.RenderTableData(data, output.FormatTable,
    output.RenderOptions{ColorMode: output.ColorModeAuto})
```

`ColorModeAuto` detects TTY via `golang.org/x/term`, respects `NO_COLOR`, `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`, `GO_OUTPUT_FORCE_COLOR`, `FORCE_COLOR`.

---

## Streaming Renderer

For large datasets, stream output incrementally:

```go
renderer := markup.NewStreamingHTMLRenderer()
renderer.SetData(tableData)
_ = renderer.Stream(os.Stdout)
```

---

## Escape Functions

The `escape/` subpackage provides safe escaping for each format:

```go
import "github.com/larsartmann/go-output/escape"

safe := escape.HTML("<script>alert('xss')</script>")
safeID := escape.D2("my-node.with.dots")
```

| Function             | Purpose                          |
| -------------------- | -------------------------------- |
| `escape.HTML`        | HTML special characters          |
| `escape.XML`         | XML special characters           |
| `escape.D2`          | D2 diagram identifiers           |
| `escape.DOT`         | DOT graph identifiers            |
| `escape.MermaidID`   | Mermaid node IDs                 |
| `escape.MermaidText` | Mermaid labels                   |
| `escape.PlantUML`    | PlantUML labels                  |
| `escape.SlugifyID`   | Cross-format diagram identifiers |

---

## Examples

Run the basic example to see all 16 formats:

```bash
go run ./examples/basic markdown              # auto color
go run ./examples/basic tree --color always    # force colors
go run ./examples/basic table --color never    # no colors
```

Other examples:

| Example                    | What it demonstrates                                    |
| -------------------------- | ------------------------------------------------------- |
| `examples/basic/`          | All 16 formats with a `Project` dataset                 |
| `examples/nom_progress/`   | NOM workflow events, dependency trees, activity counts  |
| `examples/tui_progress/`   | Bubble Tea TUI with step reporting and NOM display mode |
| `examples/d2/`             | D2 microservice architecture with SQL tables and shapes |
| `examples/diagram_export/` | Export NOM live state as DOT/Mermaid diagrams           |

---

## Development

### Nix (recommended)

```bash
nix develop                    # Enter dev shell (Go 1.26, golangci-lint, gopls)
nix run .#build              # Build all 20 modules
nix run .#test               # Test all 20 modules
nix run .#test-race          # Race-test nom + tui
nix run .#lint               # golangci-lint across all modules
nix run .#tidy               # go mod tidy all modules
nix run .#govulncheck        # Vulnerability scan
nix fmt                      # Format .nix files
nix flake check              # Formatting + pre-commit hooks
```

### Go toolchain

```bash
go build ./...                  # Build all workspace modules
go test ./...                   # Test all modules
go test -race ./...             # Race detector
go test -cover ./...            # Coverage report
golangci-lint run --fix ./...   # Lint
```

---

## API Stability

This library is pre-v1. The following guarantees apply:

- **Root module** (`github.com/larsartmann/go-output`): Public API is stable. Breaking changes documented in CHANGELOG.md.
- **Sub-modules** (`d2`, `graph`, `table`, `nom`, `tui`, etc.): May evolve independently. Import them explicitly to opt in.
- **`Renderer` interface**: Stable — all formats implement `Render() (string, error)`.
- **`internal/` packages**: No stability guarantee. Do not import these.

### Frozen Interfaces (v1 locked)

| Interface            | Methods                                                      | Implementations                                                          |
| -------------------- | ------------------------------------------------------------ | ------------------------------------------------------------------------ |
| `Renderer`           | `Render() (string, error)`                                   | All 16 formats                                                           |
| `TableRenderer`      | `SetHeaders([]string)`, `AddRow([]string)`, `Render()`       | JSON, YAML, TOML, JSONL, HTML, Streaming HTML, AsciiDoc, Markdown, Table |
| `TreeOutputRenderer` | `SetRoot(*TreeNode)`, `Render()`                             | ASCII Tree, JSON Tree, YAML Tree, TOML Tree, HTML Tree                   |
| `GraphRenderer`      | `SetNodes([]GraphNode)`, `SetEdges([]GraphEdge)`, `Render()` | D2, DOT, Mermaid, PlantUML, JSON Graph, YAML Graph, TOML Graph           |
| `StreamingRenderer`  | `Stream(io.Writer) error`, `Render()`                        | Streaming HTML                                                           |

Non-breaking changes until v1: adding new formats, shapes, methods, sub-modules, and renderers.

---

## Architecture

20 modules in a multi-module Go workspace. The root package has **zero imports of any sub-module** — this is the load-bearing architectural guarantee.

- `go get go-output` pulls **no** lipgloss, bubbletea, yaml, d2, graph, table, nom, or tui deps.
- Sub-modules self-register into root's registries via their own `init()`.
- Import a sub-module to activate its renderers automatically.

Read [`docs/FORMAT_ARCHITECTURE.md`](docs/FORMAT_ARCHITECTURE.md) and [`docs/DOMAIN_LANGUAGE.md`](docs/DOMAIN_LANGUAGE.md) for the full design.

---

## License

[MIT](LICENSE)
