# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across 12 formats (Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT) with type-safe enum-based configuration.

**Updated:** 2026-05-07

## Location

`/home/lars/projects/go-output/`

## Repository

https://github.com/larsartmann/go-output

## Key Technologies

- Go 1.26+
- charm.land/lipgloss/v2 (terminal styling — **in table/ module only, not root**)
- github.com/go-faster/yaml (YAML support)
- golang.org/x/term (terminal detection)

## Multi-Module Workspace

This project uses Go workspace modules. Each sub-package with its own `go.mod` is an independent module:

| Module                  | go.mod | Deps                       | Notes                           |
| ----------------------- | ------ | -------------------------- | ------------------------------- |
| Root (`package output`) | ✅     | enum, escape, yaml, x/term | Core types + formatters         |
| `enum/`                 | ✅     | None                       | Generic enum utilities          |
| `escape/`               | ✅     | None                       | Format-specific escaping        |
| `cmdguard/`             | ✅     | None                       | CLI flag parsing                |
| `table/`                | ✅     | root, lipgloss             | **Lipgloss isolated from root** |
| `sort/`                 | ✅     | root                       | **Deprecated** — use stdlib     |
| `integration/`          | ✅     | root, sort, table          | Cross-module tests              |
| `examples/`             | ✅     | root, table                | Usage examples                  |

`go.work` is gitignored (local dev only). Each module uses `replace` directives for standalone development.

**Key benefit:** `go get github.com/larsartmann/go-output` pulls ZERO lipgloss deps. Users who need terminal tables import `go-output/table` explicitly.

## Project Structure

```
go-output/                    # Root module (package output) — types, interfaces, formatters
├── format.go                 # Format enum + Renderer/TableData/TreeNode types
├── format_deprecated.go      # OutputFormat backward compat aliases
├── sort.go                   # SortBy enum
├── color.go                  # ColorMode enum + terminal detection
├── ids.go                    # BrandedID phantom types
├── registry.go               # Opt-in renderer registry
├── slices.go                 # FilledStrings utility
├── json.go, csv.go, tsv.go, yaml.go, xml.go, markdown.go
├── html.go, tree.go, streaming.go
├── graph.go                  # GraphNode, GraphEdge, GraphRenderer, AddTreeNodes
├── dot.go                    # DOT/Graphviz renderer + GraphRendererMixin
├── mermaid.go                # Mermaid diagram renderer
├── delimited.go, markup.go, marshal.go
├── d2.go, d2_enum.go, d2_render.go, d2_write.go, d2_convert.go
├── internal/gentest/         # Generic test helpers
├── internal/testutils/       # Domain-aware test helpers
│
├── enum/                     # MODULE: Generic enum utilities (zero deps)
├── escape/                   # MODULE: Format-specific escaping (zero deps)
├── cmdguard/                 # MODULE: CLI flag parsing (zero deps)
├── table/                    # MODULE: Lipgloss terminal tables (lipgloss isolated)
├── sort/                     # MODULE: Generic sorting (DEPRECATED — use stdlib)
├── integration/              # MODULE: Cross-module integration tests
└── examples/                 # MODULE: Usage examples
```

## Build Commands

```bash
go build ./...                  # Build all workspace modules
go test ./...                   # Test all workspace modules
golangci-lint run --fix ./...   # Lint all modules
go mod tidy                     # Tidy root module (run in each submodule too)
```

**Note:** `go.work` is gitignored. Run from project root to use workspace mode.

## Code Quality Standards

- All code must pass `golangci-lint` with `.golangci.yml` configuration
- Tests required for all public APIs
- 90%+ test coverage target
- File size limit: 350 lines per file
- No code duplication (threshold: 30 tokens)
- Each module's `go.mod` must have `replace` directives for sibling deps

## Current Coverage

| Package       | Coverage | Module |
| ------------- | -------- | ------ |
| output (root) | 90.3%    | root   |
| cmdguard      | 100%     | own    |
| enum          | 100%     | own    |
| escape        | 100%     | own    |
| sort          | 100%     | own    |
| table         | 100%     | own    |

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
3. **Interface-based design**: Renderer, GraphRenderer, TableRenderer interfaces — `Render() (string, error)`. Use `MustRender(r)` helper for tests/examples.
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
- Depguard config restricts imports — `cmp` is allowed for sort.ByField
- escape/ uses `html.EscapeString()` from stdlib for HTML, with `strings.ReplaceAll` for XML `&apos;`
- sort/ is **deprecated** — use `slices.SortStableFunc` + `cmp.Compare` (stdlib, Go 1.21+)
- SortBy enum kept in root — used by cmdguard tests as example enum type
- Multi-module workspace with 7 independent modules (see ADR 001)
- GraphRendererMixin defined in dot.go — should move to graph.go or own file when graph/ is extracted as module
