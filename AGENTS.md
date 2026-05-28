# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across 16 formats (Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT, JSONL, AsciiDoc, TOML, PlantUML) with type-safe enum-based configuration and a Shape capability matrix.

**Updated:** 2026-05-28

## Location

`/home/lars/projects/go-output/`

## Repository

https://github.com/larsartmann/go-output

## Key Technologies

- Go 1.26+
- charm.land/lipgloss/v2 (terminal styling — **in table/ module only, not root**)
- github.com/go-faster/yaml (YAML support — **in serialization/ module only, not root**)
- github.com/pelletier/go-toml/v2 (TOML support — **in serialization/ module only, not root**)
- golang.org/x/term (terminal detection)
- github.com/larsartmann/go-branded-id (phantom types for type-safe IDs)

## Multi-Module Workspace

This project uses Go workspace modules. Each sub-package with its own `go.mod` is an independent module:

| Module                  | go.mod | Deps                                                               | Notes                                   |
| ----------------------- | ------ | ------------------------------------------------------------------ | --------------------------------------- |
| Root (`package output`) | ✅     | enum, escape, yaml, x/term, branded-id, testhelpers                | Core types + formatters                 |
| `enum/`                 | ✅     | testhelpers (tests only)                                           | Generic enum utilities                  |
| `escape/`               | ✅     | None                                                               | Format-specific escaping                |
| `testhelpers/`          | ✅     | None                                                               | Shared test assertions (non-internal)   |
| `d2/`                   | ✅     | root, escape, testhelpers                                          | D2 diagram renderer (rich domain model) |
| `graph/`                | ✅     | root, escape, testhelpers                                          | DOT + Mermaid renderers                 |
| `plantuml/`             | ✅     | root                                                               | PlantUML diagram renderer               |
| `table/`                | ✅     | root, lipgloss                                                     | **Lipgloss isolated from root**         |
| `integration/`          | ✅     | root, delimited, serialization, markup, table, d2, graph, plantuml | Cross-module tests                      |
| `examples/`             | ✅     | root, delimited, serialization, markup, table, d2, graph, plantuml | Usage examples                          |

`go.work` is gitignored (local dev only). Each module uses `replace` directives for standalone development.

**Key benefit:** `go get github.com/larsartmann/go-output` pulls ZERO lipgloss deps, ZERO yaml deps, ZERO toml deps, and ZERO d2/graph/plantuml deps. Users import only the modules they need.

### Dependency Graph

```
root (output) → enum, x/term, go-branded-id, delimited, serialization
delimited     → root
serialization → root, go-faster/yaml, go-toml/v2
markup        → root, escape
plantuml      → root
enum          → testhelpers (tests only)
escape        → (none)
testhelpers   → (none) — zero deps, shared test assertions
d2            → root, escape, testhelpers
graph         → root, escape, testhelpers
table         → root, lipgloss/v2
integration   → root, delimited, serialization, markup, table, d2, graph, plantuml
examples      → root, delimited, serialization, markup, table, d2, graph, plantuml
```

**No circular dependencies.** Root has zero imports from d2/, graph/, table/, or any sub-module.

## Project Structure

```
go-output/                    # Root module (package output) — core types, Markdown, Tree
├── format.go                 # Format enum, ParseFormat, InvalidFormatError, AllFormats
├── shape.go                  # Shape enum, capability matrix, Supports/Shapes/FormatsForShape
├── renderer.go               # Renderer, MustRender, TableRenderer interfaces
├── color.go                  # ColorMode enum + terminal detection (wired into table, tree, markdown)
├── ids.go                    # BrandedID phantom types
├── tabledata.go              # TableData, RowEdge, TableDataBase (exported for sub-modules)
├── tree.go                   # TreeNode, TreeOutputRenderer, colored tree rendering
├── graph.go                  # GraphNode, GraphEdge, GraphRenderer, GraphRendererMixin, AddTreeNodes
├── markdown.go               # Markdown table renderer with ColorMode (bold headers, dim separators)
├── marshal.go                # MarshalFormat, UnmarshalFormat, MarshalJSONIndent
├── render_tabledata.go       # Registry-based TableData dispatch with ColorMode in RenderOptions
├── streaming.go              # StreamingRenderer interface + adapter
├── internal/gentest/         # Generic test helpers (root module only, not importable by sub-modules)
│
├── delimited/                # MODULE: CSV + TSV writers and formatters
├── serialization/            # MODULE: JSON + YAML + TOML + JSONL marshaling and renderers
├── markup/                   # MODULE: XML + HTML + AsciiDoc + Streaming HTML renderers
├── d2/                       # MODULE: D2 diagram renderer (rich domain model)
├── graph/                    # MODULE: DOT + Mermaid renderers
├── plantuml/                 # MODULE: PlantUML diagram renderer
│
├── enum/                     # MODULE: Generic enum utilities (zero deps)
├── escape/                   # MODULE: Format-specific escaping (zero deps)
├── testhelpers/              # MODULE: Shared test assertions (non-internal, zero deps)
├── table/                    # MODULE: Lipgloss terminal tables (lipgloss isolated)
├── integration/              # MODULE: Cross-module integration tests
└── examples/                 # MODULE: Usage examples
```

## Build Commands

```bash
# Nix (recommended)
nix develop                    # Enter dev shell (Go 1.26, golangci-lint, gopls)
nix develop .#ci               # CI dev shell (no gopls, smaller closure)
nix fmt                        # Format .nix files
nix flake check                # Verify formatting + pre-commit hooks

# Nix apps (iterate all 13 modules automatically)
nix run .#build                # Build all modules
nix run .#test                 # Test all modules
nix run .#lint                 # Lint all modules
nix run .#tidy                 # go mod tidy all modules
nix run .#setup-workspace      # Generate go.work from go.work.example

# Inside nix develop (or with Go installed) — single module:
go build ./...                  # Build root module
go test ./...                   # Test root module
golangci-lint run ./...         # Lint root module
go mod tidy                     # Tidy root module
```

**Note:** `go.work` is gitignored. Run per-module commands or create a local `go.work`:

```bash
cat > go.work << 'EOF'
go 1.26.3

use (
  .
  ./delimited
  ./d2
  ./enum
  ./escape
  ./examples
  ./graph
  ./integration
  ./markup
  ./plantuml
  ./serialization
  ./table
  ./testhelpers
  ./testhelpers/graphtest
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
| output (root) | 96.1%    | root   |
| delimited     | 90.2%    | own    |
| serialization | 91.4%    | own    |
| markup        | 93.9%    | own    |
| d2            | 100%     | own    |
| graph         | 96.0%    | own    |
| enum          | 100%     | own    |
| escape        | 100%     | own    |
| table         | 100%     | own    |
| testhelpers   | 91.3%    | own    |
| integration   | 95.5%    | own    |
| gentest       | 96.2%    | root   |

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
5. **Composition**: GraphRendererMixin in graph.go (root) provides shared state for DOT/Mermaid via accessor methods (Nodes/Edges/NodesPtr/EdgesPtr). TableDataBase in tabledata.go shared by delimited, markup, serialization via Data() getter.
6. **Registry dispatch**: `RenderTableData()` dispatches to registered `TableDataMarshaler` functions. Sub-modules register via `init()`. Root has ZERO sub-module imports.
7. **ColorMode wired into renderers**: `ColorMode` (auto/always/never) controls ANSI color output. `table.New(WithColorMode(mode))` for lipgloss tables, `ASCIITreeRenderer.SetColorMode(mode)` for trees, `MarkdownTable.SetColorMode(mode)` for markdown, `RenderOptions.ColorMode` for `RenderTableData()` dispatch.
8. **Footer row on TableData**: `TableData.Footer []string` provides an optional totals/summary row. Tabular renderers (CSV, TSV, Markdown, HTML, XML, AsciiDoc, Table) render it visually. HTML uses `<tfoot>` with `footer-cell` CSS class, XML uses `<footer>`, Markdown adds a second separator + bold footer (inherits column alignment). The `table.FooterProvider` optional interface enables `table.FromTableData()` to bold-style the footer. `CSVWriter.WriteFooter()` and `TSVWriter.WriteFooter()` provide explicit streaming footer methods. `TableData.Validate()` checks footer column count matches headers. Data formats (JSON, YAML, TOML, JSONL) and non-tabular formats (Tree, Graph) skip the footer.
9. **Delimited dedup**: `delimited.marshalFromTableData()` is a shared helper used by both `MarshalCSVFromTableData` and `MarshalTSVFromTableData`, eliminating structural duplication via a `tableDataWriter` interface.
10. **TableRenderer adapter pattern**: Both `MarkdownTable` and `table.Table` use fluent/builder APIs (returning `*Self`) incompatible with the void-returning `TableRenderer` interface. `AsTableRenderer()` adapter methods wrap the builder with void-returning delegates. Pattern: `md.AsTableRenderer()` for MarkdownTable, `tbl.AsTableRenderer()` for table.Table.
11. **WithFooterStyle option**: `table.WithFooterStyle(func(lipgloss.Style) lipgloss.Style)` provides composable footer styling for lipgloss tables. Receives the base style (with padding) and returns the styled result.

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

### Using ColorMode

All terminal renderers support `ColorMode` for controlling ANSI color output:

```go
// Table: functional options pattern
tbl := table.New(table.WithColorMode(output.ColorModeAlways))

// Tree: setter method
tree := output.NewASCIITreeRenderer()
tree.SetColorMode(output.ColorModeAlways)

// Markdown: setter method (chains)
md := output.NewMarkdownTable().SetColorMode(output.ColorModeAlways)

// RenderTableData dispatch: pass via RenderOptions
output.RenderTableData(data, output.FormatTree, output.RenderOptions{ColorMode: output.ColorModeAlways})
```

`ColorModeAuto` (default) detects terminal via `golang.org/x/term`, respects `NO_COLOR`, `CI`, `FORCE_COLOR`.

### Sub-Module Usage Pattern

Each sub-module is independently versioned. Users import only what they need:

```go
import "github.com/larsartmann/go-output"                  // core types + Markdown/Tree formatters
import "github.com/larsartmann/go-output/delimited"          // CSV + TSV (optional)
import "github.com/larsartmann/go-output/serialization"      // JSON + YAML + TOML + JSONL (optional)
import "github.com/larsartmann/go-output/markup"             // XML + HTML + AsciiDoc + Streaming (optional)
import "github.com/larsartmann/go-output/d2"                 // D2 diagrams (optional)
import "github.com/larsartmann/go-output/graph"              // DOT + Mermaid (optional)
import "github.com/larsartmann/go-output/table"              // Lipgloss tables (optional)
import "github.com/larsartmann/go-output/plantuml"            // PlantUML diagrams (optional)
```

## Architecture Notes

- **Root has ZERO sub-module imports** — verified via `go mod graph`. Users get zero transitive deps from sub-modules they don't import.
- Root no longer imports `go-faster/yaml` or `escape` in production code — these deps are isolated in serialization/ and markup/ respectively
- D2 has richer types than generic graph (shapes, arrows, SQL tables, classes) — lives in `d2/` module with its own `D2Node`/`D2Edge` types
- D2 re-exports `D2NodeID`/`D2NodeLabel` from root so users don't need to import both `d2` and `output` for ID construction
- DOT and Mermaid renderers live in `graph/` module, sharing `GraphRendererMixin` from root via accessor methods
- Tree conversion has renderer-specific addTreeNodes in d2_convert, graph/dot, graph/mermaid — the generic AddTreeNodes in graph.go handles the common case
- Depguard config restricts imports
- escape/ uses `html.EscapeString()` from stdlib for HTML, with `strings.ReplaceAll` for XML `&apos;`
- sort/ is **removed** — was deprecated, only `ByField` helper remained (zero deps). Use `slices.SortStableFunc` + `cmp.Compare` (stdlib)
- registry.go **removed** — deprecated renderer registry had zero external callers. Use direct constructors instead.
- slices.go **removed** — `FilledStrings` was trivially replaced by `slices.Repeat` (stdlib)
- `BrandedValue()` removed from `marshal.go` — zero callers after sub-module extraction
- `ColorMode` wired into table (via `WithColorMode` functional option), tree (via `SetColorMode()`), markdown (via `SetColorMode()`), and `RenderTableData()` (via `RenderOptions.ColorMode`)
- Tree rendering uses depth-based ANSI color cycling when color is enabled
- Markdown rendering uses bold headers and dim separators when color is enabled
- Multi-module workspace with 13 independent modules (see ADR 001)
- Shape capability matrix (ADR 002) replaces FormatCategory — deprecated methods redirect to `Supports(Shape)`
- API stability audit (ADR 006) completed — all exported symbols frozen, capability matrix fixed (D2/Mermaid/DOT/PlantUML now declare ShapeTree, TOML now declares ShapeGraph)
- Round-trip integration tests in `integration/roundtrip_test.go` verify 16 formats: 8 parseable round-trips (JSON, CSV, TSV, YAML, TOML, JSONL, XML, HTML) + 8 structural verifications (Markdown, Table, Tree, AsciiDoc, D2, Mermaid, DOT, PlantUML)
- `format.go` split into focused files: `format.go` (Format enum), `shape.go` (Shape + capability matrix), `renderer.go` (Renderer/TableRenderer interfaces) — `format_deprecated.go` deleted
- `render_tabledata.go` uses registry-based dispatch via `TableDataMarshaler` — sub-modules register in `init()`. Returns `UnsupportedFormatError` for D2/DOT/Mermaid (these live in separate modules)
- `tableDataBase` exported as `TableDataBase` with `Data()` getter — allows sub-modules to access unexported `data` field
- `marshal.go` exports `MarshalFormat()`, `UnmarshalFormat()`, `MarshalJSONIndent()` — used by both serialization/ and markup/ sub-modules
- Root format files (CSV, TSV, JSON, YAML, XML, HTML, Streaming) extracted to delimited/, serialization/, markup/ modules respectively
- `internal/gentest` and `internal/testutils` are root-only — sub-modules must inline helpers or create their own. Decision rationale: exposing test helpers publicly via `testhelpers/gentest` would freeze internal testing APIs; each module having its own test helpers allows independent evolution
- Nix flake uses `flake-parts` + `treefmt-nix` + `git-hooks.nix` — no `gomod2nix` (library, 13 modules, no binary)
- Go checks (build/test/lint) run via `nix run .#test` / `nix run .#lint` / `nix run .#build` — these iterate all 13 modules. Not in `nix flake check` because the Nix sandbox blocks `go mod download`; CI handles these reliably
- `.pre-commit-config.yaml` exists for non-Nix users; `git-hooks.nix` auto-installs hooks for Nix users via `nix develop`

## Code Duplication Policy

**Updated:** 2026-05-28

At `art-dupl -t 50` (industry standard), this codebase has **zero actionable clones**.

At `art-dupl -t 15` (aggressive), ~50 clone groups appear. These are categorized as:

| Category | Description | Action |
|----------|-------------|--------|
| **Go test idioms** | `strings.Contains` assertions, `t.Errorf` patterns, `t.Parallel()` | Accept — language patterns |
| **Module boundary** | Interface re-declarations, type aliases across modules | Accept — Go design constraint |
| **Example/docs** | Full API usage in `example_test.go` and `examples/` | Accept — intentional for documentation |
| **Single-line** | `render*TableData` signatures, `init()` registrations | Accept — interface compliance |

### Key Decisions

- `testhelpers` is **zero-dep by design** — cannot import `output`. Cross-module test helpers must stay local or use table-driven patterns within each module
- `serialization/render.go` has `renderViaRenderer()` shared helper for YAML/TOML (identical `render*TableData` bodies)
- `graphtest.NewTestNode`/`TestEdgeAB` used in serialization, graph, d2, plantuml benches — not in examples (must show full API)
- **Threshold 15 is too aggressive for action** — use t=30-40 for meaningful dedup work
