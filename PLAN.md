# go-output Library

A reusable Go library for CLI applications providing consistent output formatting across multiple formats with type-safe enum-based configuration.

## Package Structure

```
go-output/
├── format.go              # Format enum + type definitions (Renderer, TableData, TreeNode, etc.)
├── format_deprecated.go   # Backward compatibility aliases (OutputFormat)
├── sort.go                # SortBy enum
├── color.go               # ColorMode enum + terminal detection
├── ids.go                 # BrandedID phantom types for type-safe IDs
│
├── json.go                # JSON marshal/unmarshal + JSONWriter
├── csv.go                 # CSV writer
├── tsv.go                 # TSV writer + MarshalTSV
├── yaml.go                # YAML marshal/unmarshal
├── xml.go                 # XML writer + MarshalXMLFromTableData
├── markdown.go            # Markdown table builder with alignment
├── html.go                # HTML table + tree renderers
├── tree.go                # ASCII tree renderer
│
├── d2.go                  # D2 domain types (D2Node, D2Edge, D2Table, etc.)
├── d2_render.go           # D2Diagram builder + Render()
├── d2_write.go            # D2 style/edge writing helpers
├── d2_convert.go          # TableData/Tree → D2 conversion
│
├── graph.go               # Generic graph types (GraphNode, GraphEdge, GraphShape)
├── dot.go                 # DOT/Graphviz renderer + GraphRendererMixin
├── mermaid.go             # Mermaid diagram renderer
│
├── delimited.go           # Shared CSV/TSV DelimitedWriter
├── markup.go              # Shared XML/HTML row writing helpers
├── marshal.go             # Shared marshal/unmarshal error wrapping
├── streaming.go           # Streaming HTML renderer + adapter
├── slices.go              # Slice utilities
├── registry.go            # Opt-in renderer registry (plugin system)
│
├── enum/                  # Generic enum utilities (Parse, Contains, AllowedValues)
├── table/                 # Lipgloss-based terminal table renderer
├── sort/                  # Generic Sorter[T] with reflect-based field comparison
├── cmdguard/              # Generic EnumFlag[T] for cmdguard integration
├── internal/escape/       # Format-specific escaping (HTML, XML, D2, DOT, Mermaid)
├── internal/gentest/      # Test assertion helpers
└── internal/testutils/    # Test helper utilities
```

## Supported Formats

| Format | Constant | Category | Renderer |
|--------|----------|----------|----------|
| Table | `FormatTable` | table | `table.Table` (lipgloss) |
| JSON | `FormatJSON` | table | `MarshalJSON` / `JSONWriter` |
| CSV | `FormatCSV` | table | `CSVWriter` |
| TSV | `FormatTSV` | table | `TSVWriter` |
| Markdown | `FormatMarkdown` | table | `MarkdownTable` |
| XML | `FormatXML` | table | `XMLWriter` |
| YAML | `FormatYAML` | table | `MarshalYAML` |
| HTML | `FormatHTML` | tree | `HTMLRenderer` / `StreamingHTMLRenderer` |
| Tree | `FormatTree` | tree | `ASCIITreeRenderer` |
| D2 | `FormatD2` | graph | `D2Diagram` |
| Mermaid | `FormatMermaid` | graph | `MermaidRenderer` |
| DOT | `FormatDOT` | graph | `DOTRenderer` |

## Core Enums

| Enum | Values | Purpose |
|------|--------|---------|
| `Format` | 12 formats | Output format selection |
| `SortBy` | name, importance, created_at, updated_at, health, complexity | Sort field selection |
| `ColorMode` | auto, always, never | Color output control |

All enums provide `ParseX()`, `String()`, `AllowedValues()`, `IsValid()` methods.

## Quick Start

```go
import (
    "fmt"
    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/table"
)

// JSON
jsonBytes, _ := output.MarshalJSONIndent(data, "", "  ")

// CSV
w := output.NewCSVWriter(os.Stdout)
w.WriteHeader([]string{"Name", "Age"})
w.WriteRow([]string{"Alice", "30"})
w.Flush()

// Markdown
md := output.NewMarkdownTable()
md.SetHeaders([]string{"Name", "Age"})
md.AddRow([]string{"Alice", "30"})
fmt.Println(md.Render())

// Terminal table
t := table.FromTableData(output.NewTableData([]string{"Name"}))
t.AddRow("Alice")
fmt.Println(t.Render())

// D2 diagram
d := output.NewD2Diagram()
d.AddNodeSimple("server", "Server")
d.AddNodeSimple("db", "Database")
d.AddEdgeSimple("server", "db")
fmt.Println(d.Render())
```

## Dependencies

- `charm.land/lipgloss/v2` — Terminal styling and table rendering
- `github.com/go-faster/yaml` — YAML marshaling
- `golang.org/x/term` — Terminal detection

## Build Commands

```bash
just build    # go build ./...
just test     # go test ./...
just lint     # golangci-lint run --fix ./...
just verify   # build + test + lint
```

## Status

All phases complete. Library is production-ready.
