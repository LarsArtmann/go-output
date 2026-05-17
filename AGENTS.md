# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across 12 formats (Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT) with type-safe enum-based configuration and a Shape capability matrix.

**Updated:** 2026-05-17

## Location

`/home/lars/projects/go-output/`

## Repository

https://github.com/larsartmann/go-output

## Key Technologies

- Go 1.26+
- charm.land/lipgloss/v2 (terminal styling — **in table/ module only, not root**)
- github.com/go-faster/yaml (YAML support)
- golang.org/x/term (terminal detection)
- github.com/larsartmann/go-branded-id (phantom types for type-safe IDs)

## Multi-Module Workspace

This project uses Go workspace modules. Each sub-package with its own `go.mod` is an independent module:

| Module                  | go.mod | Deps                                   | Notes                                   |
| ----------------------- | ------ | -------------------------------------- | --------------------------------------- |
| Root (`package output`) | ✅     | enum, escape, yaml, x/term, branded-id, testhelpers | Core types + formatters                 |
| `enum/`                 | ✅     | testhelpers (tests only)               | Generic enum utilities                  |
| `escape/`               | ✅     | None                                   | Format-specific escaping                |
| `testhelpers/`          | ✅     | None                                   | Shared test assertions (non-internal)   |
| `cmdguard/`             | ✅     | root (tests only), testhelpers (tests) | CLI flag parsing (prod standalone)      |
| `table/`                | ✅     | root, lipgloss                         | **Lipgloss isolated from root**         |
| `sort/`                 | ✅     | None                                   | **Deprecated** — only `ByField` remains |
| `integration/`          | ✅     | root, table                            | Cross-module tests                      |
| `examples/`             | ✅     | root, table                            | Usage examples                          |

`go.work` is gitignored (local dev only). Each module uses `replace` directives for standalone development.

**Key benefit:** `go get github.com/larsartmann/go-output` pulls ZERO lipgloss deps. Users who need terminal tables import `go-output/table` explicitly.

### Dependency Graph

```
root (output) → enum, escape, yaml, x/term, go-branded-id, testhelpers
enum          → testhelpers (tests only)
escape        → (none)
testhelpers   → (none) — zero deps, shared test assertions
sort          → (none) — zero deps, only ByField helper
cmdguard      → root (tests only), testhelpers (tests); prod code standalone
table         → root, lipgloss/v2
integration   → root, table
examples      → root, table
```

**No circular dependencies.** The sort/ → root cycle was eliminated in v0.4.0 by removing `Sorter[T]`.

## Project Structure

```
go-output/                    # Root module (package output) — types, interfaces, formatters
├── format.go                 # Format enum, Shape enum, capability matrix, Renderer/TableRenderer interfaces
├── format_deprecated.go      # FormatCategory/OutputFormat backward compat (deprecated)
├── sort.go                   # SortBy enum
├── color.go                  # ColorMode enum + terminal detection
├── ids.go                    # BrandedID phantom types
├── registry.go               # Opt-in renderer registry
├── slices.go                 # FilledStrings utility
├── tabledata.go              # TableData, RowEdge, tableDataBase (extracted from format.go/html.go)
├── tree.go                   # TreeNode, TreeOutputRenderer (extracted from format.go)
├── json.go, json_renderers.go, csv.go, tsv.go, yaml.go, yaml_renderers.go, xml.go, markdown.go
├── html.go, streaming.go
├── graph.go                  # GraphNode, GraphEdge, GraphRenderer, GraphRendererMixin, AddTreeNodes
├── dot.go                    # DOT/Graphviz renderer
├── mermaid.go                # Mermaid diagram renderer
├── delimited.go, markup.go, marshal.go
├── d2.go, d2_enum.go, d2_render.go, d2_write.go, d2_convert.go
├── internal/gentest/         # Generic test helpers (root module only, not importable by sub-modules)
├── internal/testutils/       # Domain-aware test helpers
│
├── enum/                     # MODULE: Generic enum utilities (zero deps)
├── escape/                   # MODULE: Format-specific escaping (zero deps)
├── testhelpers/              # MODULE: Shared test assertions (non-internal, zero deps)
├── cmdguard/                 # MODULE: CLI flag parsing (prod standalone, tests import root)
├── table/                    # MODULE: Lipgloss terminal tables (lipgloss isolated)
├── sort/                     # MODULE: Deprecated — only ByField helper remains (zero deps)
├── integration/              # MODULE: Cross-module integration tests
└── examples/                 # MODULE: Usage examples
```

## Build Commands

```bash
go build ./...                  # Build root module
go test ./...                   # Test root module
golangci-lint run ./...         # Lint root module
go mod tidy                     # Tidy root module

# Per-module (required since go.work is gitignored)
for mod in . enum escape testhelpers sort cmdguard table integration; do
  (cd $mod && go test ./...)
done
```

**Note:** `go.work` is gitignored. Run per-module commands or create a local `go.work`:

```bash
cat > go.work << 'EOF'
go 1.26.2

use (
  .
  ./enum
  ./escape
  ./cmdguard
  ./sort
  ./table
  ./integration
  ./examples
)
EOF
```

## Code Quality Standards

- All code must pass `golangci-lint` with `.golangci.yml` configuration
- Tests required for all public APIs
- 90%+ test coverage target
- File size limit: 350 lines per file
- No code duplication (threshold: 30 tokens)
- Each module's `go.mod` must have `replace` directives for sibling deps
- Sub-modules must NOT import `internal/` packages from root (Go restriction)

## Current Coverage

| Package       | Coverage | Module |
| ------------- | -------- | ------ |
| output (root) | 90%+     | root   |
| cmdguard      | 100%     | own    |
| enum          | 100%     | own    |
| escape        | 100%     | own    |
| sort          | 100%     | own    |
| table         | 100%     | own    |
| testhelpers   | 100%     | own    |

## Testing

```bash
go test ./...              # Unit tests (root module)
go test -race ./...        # Race detector
go test -cover ./...       # Coverage
go test -bench=. -benchmem ./...  # Benchmarks
```

## Key Design Patterns

1. **Type-safe enums**: String constants with Parse/Validate via `enum` package. Every enum has `Parse()`, `String()`, `IsValid()`, `AllowedValues()`.
2. **Shape capability matrix**: Each format declares supported data shapes via `formatCapabilities` map[Format][]Shape. Use `f.Supports(shape)` to query.
3. **Branded IDs**: Phantom types prevent mixing D2NodeID/TreeNodeID/etc via `go-branded-id`.
4. **Interface-based design**: Renderer, GraphRenderer, TableRenderer, TreeOutputRenderer interfaces — all have `Render() (string, error)`. Use `MustRender(r)` for tests/examples.
5. **Composition**: GraphRendererMixin in graph.go shared by DOT/Mermaid, tableDataBase in tabledata.go shared by HTML/Streaming.
6. **Registry is opt-in**: Use constructors directly by default. Register/Create for runtime dispatch.

## Common Tasks

### Adding a New Output Format

1. Add format constant to `format.go`
2. Add to `formatCapabilities` map with supported shapes
3. Implement formatter — embed Renderer interface
4. Add tests with >90% coverage
5. Update cmdguard if needed (EnumFlag already generic)
6. Update CHANGELOG.md and README.md

### Adding a New D2 Enum

1. Add type + constants to `d2_enum.go`
2. Add values slice + Parse/IsValid/AllowedValues/String methods
3. Add tests to `d2_enum_test.go`

## Architecture Notes

- D2 has richer types than generic graph (shapes, arrows, SQL tables, classes) — intentional split
- Tree conversion has renderer-specific addTreeNodes in d2_convert, dot, mermaid — the generic AddTreeNodes in graph.go handles the common case
- Depguard config restricts imports
- escape/ uses `html.EscapeString()` from stdlib for HTML, with `strings.ReplaceAll` for XML `&apos;`
- sort/ is **deprecated** — `Sorter[T]` deleted, only `ByField` helper remains (zero deps). Use `slices.SortStableFunc` + `cmp.Compare` (stdlib)
- SortBy enum kept in root — used by cmdguard tests as example enum type
- Multi-module workspace with 7 independent modules (see ADR 001)
- GraphRendererMixin in `graph.go` — shared by DOT and Mermaid renderers
- Shape capability matrix (ADR 002) replaces FormatCategory — deprecated methods redirect to `Supports(Shape)`
- `internal/gentest` and `internal/testutils` are root-only — sub-modules must inline helpers or create their own
- cmdguard prod code (`flag.go`) has zero external deps — only tests import root
