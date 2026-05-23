# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across 12 formats (Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT) with type-safe enum-based configuration and a Shape capability matrix.

**Updated:** 2026-05-23

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

| Module                  | go.mod | Deps                                                | Notes                                   |
| ----------------------- | ------ | --------------------------------------------------- | --------------------------------------- |
| Root (`package output`) | ✅     | enum, escape, yaml, x/term, branded-id, testhelpers | Core types + formatters                 |
| `enum/`                 | ✅     | testhelpers (tests only)                            | Generic enum utilities                  |
| `escape/`               | ✅     | None                                                | Format-specific escaping                |
| `testhelpers/`          | ✅     | None                                                | Shared test assertions (non-internal)   |
| `d2/`                   | ✅     | root, escape, testhelpers                           | D2 diagram renderer (rich domain model) |
| `graph/`                | ✅     | root, escape, testhelpers                           | DOT + Mermaid renderers                  |
| `table/`                | ✅     | root, lipgloss                                      | **Lipgloss isolated from root**         |
| `sort/`                 | ✅     | None                                                | **Deprecated** — only `ByField` remains |
| `integration/`          | ✅     | root, table, d2, graph                              | Cross-module tests                      |
| `examples/`             | ✅     | root, table, d2, graph                              | Usage examples                          |

`go.work` is gitignored (local dev only). Each module uses `replace` directives for standalone development.

**Key benefit:** `go get github.com/larsartmann/go-output` pulls ZERO lipgloss deps and ZERO d2/graph deps. Users import only the modules they need.

### Dependency Graph

```
root (output) → enum, escape, yaml, x/term, go-branded-id, testhelpers
enum          → testhelpers (tests only)
escape        → (none)
testhelpers   → (none) — zero deps, shared test assertions
sort          → (none) — zero deps, only ByField helper
d2            → root, escape, testhelpers
graph         → root, escape, testhelpers
table         → root, lipgloss/v2
integration   → root, table, d2, graph
examples      → root, table, d2, graph
```

**No circular dependencies.** Root has zero imports from d2/, graph/, table/, or any sub-module.

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
├── delimited.go, markup.go, marshal.go, render_tabledata.go
├── internal/gentest/         # Generic test helpers (root module only, not importable by sub-modules)
│
├── d2/                       # MODULE: D2 diagram renderer (rich domain model)
├── graph/                    # MODULE: DOT + Mermaid renderers
├── internal/testutils/       # Domain-aware test helpers
│
├── enum/                     # MODULE: Generic enum utilities (zero deps)
├── escape/                   # MODULE: Format-specific escaping (zero deps)
├── testhelpers/              # MODULE: Shared test assertions (non-internal, zero deps)
├── table/                    # MODULE: Lipgloss terminal tables (lipgloss isolated)
├── sort/                     # MODULE: Deprecated — only ByField helper remains (zero deps)
├── integration/              # MODULE: Cross-module integration tests
└── examples/                 # MODULE: Usage examples
```

## Build Commands

```bash
# Nix (recommended)
nix develop                    # Enter dev shell (Go 1.26.2, golangci-lint, gopls)
nix fmt                        # Format .nix files
nix flake check                # Verify formatting + pre-commit hooks

# Inside nix develop (or with Go installed):
go build ./...                  # Build root module
go test ./...                   # Test root module
golangci-lint run ./...         # Lint root module
go mod tidy                     # Tidy root module

# Per-module (required since go.work is gitignored)
for mod in . d2 graph enum escape testhelpers sort table integration examples; do
  (cd $mod && go test ./...)
done
```

**Note:** `go.work` is gitignored. Run per-module commands or create a local `go.work`:

```bash
cat > go.work << 'EOF'
go 1.26.2

use (
  .
  ./d2
  ./enum
  ./escape
  ./graph
  ./testhelpers
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
| d2            | 90%+     | own    |
| graph         | 90%+     | own    |
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
5. **Composition**: GraphRendererMixin in graph.go (root) provides shared state for DOT/Mermaid via accessor methods (Nodes/Edges/NodesPtr/EdgesPtr). tableDataBase in tabledata.go shared by HTML/Streaming.
6. **Registry is opt-in**: Use constructors directly by default. Register/Create for runtime dispatch.

## Common Tasks

### Adding a New Output Format

1. Add format constant to `format.go`
2. Add to `formatCapabilities` map with supported shapes
3. Implement formatter — embed Renderer interface
4. Add tests with >90% coverage
5. Update CHANGELOG.md and README.md

### Adding a New D2 Enum

1. Add type + constants to `d2/d2_enum.go`
2. Add values slice + Parse/IsValid/AllowedValues/String methods
3. Add tests to `d2/d2_enum_test.go`

### Adding a New Graph Renderer

1. Create renderer in `graph/` module
2. Embed `output.GraphRendererMixin` for shared node/edge state
3. Implement `output.GraphRenderer` interface
4. Add tests with >90% coverage

## Architecture Notes

- D2 has richer types than generic graph (shapes, arrows, SQL tables, classes) — lives in `d2/` module
- DOT and Mermaid renderers live in `graph/` module, sharing `GraphRendererMixin` from root via accessor methods
- Tree conversion has renderer-specific addTreeNodes in d2_convert, graph/dot, graph/mermaid — the generic AddTreeNodes in graph.go handles the common case
- Depguard config restricts imports
- escape/ uses `html.EscapeString()` from stdlib for HTML, with `strings.ReplaceAll` for XML `&apos;`
- sort/ is **deprecated** — `Sorter[T]` deleted, only `ByField` helper remains (zero deps). Use `slices.SortStableFunc` + `cmp.Compare` (stdlib)
- Multi-module workspace with 10 independent modules (see ADR 001)
- Shape capability matrix (ADR 002) replaces FormatCategory — deprecated methods redirect to `Supports(Shape)`
- `render_tabledata.go` returns `UnsupportedFormatError` for D2/DOT/Mermaid (these live in separate modules)
- `internal/gentest` and `internal/testutils` are root-only — sub-modules must inline helpers or create their own
- Nix flake uses `flake-parts` + `treefmt-nix` + `git-hooks.nix` — no `gomod2nix` (library, 10 modules, no binary)
- Go checks (build/test/lint) NOT in flake — Nix sandbox blocks `go mod download`; CI handles these reliably
- `.pre-commit-config.yaml` exists for non-Nix users; `git-hooks.nix` auto-installs hooks for Nix users via `nix develop`
