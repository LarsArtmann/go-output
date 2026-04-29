# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across 12 formats (Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT) with type-safe enum-based configuration.

**Updated:** 2026-04-29

## Location

`/home/lars/projects/go-output/`

## Repository

https://github.com/larsartmann/go-output

## Key Technologies

- Go 1.26+
- charm.land/lipgloss/v2 (terminal styling)
- github.com/go-faster/yaml (YAML support)
- golang.org/x/term (terminal detection)
- github.com/larsartmann/cmdguard/v2 (optional CLI flag integration - add separately)

## Project Structure

```
go-output/
├── format.go              # Format enum + Renderer/TableData/TreeNode types
├── format_deprecated.go   # OutputFormat backward compat aliases
├── sort.go                # SortBy enum
├── color.go               # ColorMode enum + terminal detection
├── ids.go                 # BrandedID phantom types
├── registry.go            # Opt-in renderer registry (plugin system)
│
├── json.go                # JSON marshal/unmarshal + JSONWriter
├── csv.go                 # CSV writer
├── tsv.go                 # TSV writer + MarshalTSV
├── yaml.go                # YAML marshal/unmarshal
├── xml.go                 # XML writer + MarshalXMLFromTableData
├── markdown.go            # Markdown table builder with alignment
├── html.go                # HTML table + tree renderers + tableDataBase
├── tree.go                # ASCII tree renderer
├── delimited.go           # Shared CSV/TSV DelimitedWriter
├── markup.go              # Shared XML/HTML row writing helpers
├── marshal.go             # Shared marshal/unmarshal error wrapping
├── streaming.go           # Streaming HTML renderer + adapter
├── slices.go              # FilledStrings utility
│
├── d2.go                  # D2 domain types (D2Node, D2Edge, D2Table)
├── d2_enum.go             # D2 enums (Direction, NodeShape, ArrowType, Constraint)
├── d2_render.go           # D2Diagram builder + Render()
├── d2_write.go            # D2 style/edge writing helpers
├── d2_convert.go          # TableData/Tree → D2 conversion
│
├── graph.go               # Generic graph types (GraphNode, GraphEdge, GraphShape)
├── dot.go                 # DOT/Graphviz renderer + GraphRendererMixin
├── mermaid.go             # Mermaid diagram renderer
│
├── enum/                  # Generic enum utilities (Parse, Contains, AllowedValues)
├── table/                 # Lipgloss-based terminal table renderer
├── sort/                  # Generic Sorter[T] with reflect-based field comparison
├── cmdguard/              # Generic EnumFlag[T] for cmdguard integration
├── internal/escape/       # Format-specific escaping (HTML, XML, D2, DOT, Mermaid)
└── examples/              # Usage examples
```

## Build Commands

```bash
just build     # go build ./...
just test      # go test ./...
just lint      # golangci-lint run --fix ./...
just verify    # build + test + lint
```

## Code Quality Standards

- All code must pass `golangci-lint` with `.golangci.yml` configuration
- Tests required for all public APIs
- 90%+ test coverage target
- File size limit: 350 lines per file
- No code duplication (threshold: 30 tokens)

## Current Coverage

| Package | Coverage |
|---------|----------|
| output (root) | 91.0% |
| cmdguard | 100% |
| enum | 100% |
| internal/escape | 100% |
| sort | 95.5% |
| table | 100% |

## Testing

```bash
go test ./...              # Unit tests
go test -race ./...        # Race detector
go test -cover ./...       # Coverage
go test -bench=. -benchmem ./...  # Benchmarks
```

## Key Design Patterns

1. **Type-safe enums**: String constants with Parse/Validate via `enum` package
2. **Branded IDs**: Phantom types prevent mixing D2NodeID/TreeNodeID/etc
3. **Interface-based design**: Renderer, GraphRenderer, TableRenderer interfaces
4. **Composition**: GraphRendererMixin shared by DOT/Mermaid, tableDataBase shared by HTML/Streaming
5. **Registry is opt-in**: Use constructors directly by default. Register/Create for runtime dispatch.

## Common Tasks

### Adding a New Output Format

1. Add format constant to `format.go`
2. Implement formatter — embed Renderer interface
3. Add to format category maps if table/tree/graph
4. Add tests with >90% coverage
5. Update cmdguard if needed (EnumFlag already generic)

### Adding a New D2 Enum

1. Add type + constants to `d2_enum.go`
2. Add values slice + Parse/IsValid/AllowedValues/String methods
3. Add tests to `d2_enum_test.go`

## Architecture Notes

- D2 has richer types than generic graph (shapes, arrows, SQL tables, classes) — intentional split
- Tree conversion has renderer-specific addTreeNodes in d2_convert, dot, mermaid — the generic AddTreeNodes in graph.go handles the common case
- Depguard config restricts imports — `cmp` is blocked, use manual comparison
- CI uses Go 1.26 (must match go.mod)
