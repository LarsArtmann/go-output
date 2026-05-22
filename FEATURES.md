# Features

Complete feature inventory for `go-output` — a Go library providing consistent output formatting across 12 formats for CLI applications.

**Status legend:** FULLY_FUNCTIONAL | PARTIALLY_FUNCTIONAL | DEPRECATED | KNOWN_ISSUE

---

## Output Formats

### Table Data Formats

| Feature | Status | Notes |
|---------|--------|-------|
| **JSON Table** (`FormatJSON`) | FULLY_FUNCTIONAL | Array-of-objects via `JSONTableRenderer`. Headers become keys, rows become values |
| **CSV** (`FormatCSV`) | FULLY_FUNCTIONAL | Streaming writer via `CSVWriter`. Auto-quotes fields with commas/newlines |
| **TSV** (`FormatTSV`) | FULLY_FUNCTIONAL | Streaming writer via `TSVWriter`. Tab-separated with type-switch marshaling (`[][]string`, `[]string`) |
| **Markdown Table** (`FormatMarkdown`) | FULLY_FUNCTIONAL | Column alignment (left/right/center). Auto-calculated column widths via `MarkdownTable` |
| **XML** (`FormatXML`) | FULLY_FUNCTIONAL | Structured `<table><headers>...</headers><rows>...</rows></table>` with XML escaping |
| **YAML Table** (`FormatYAML`) | FULLY_FUNCTIONAL | Sequence of mappings via `YAMLTableRenderer`. Uses `go-faster/yaml` |
| **HTML Table** (`FormatHTML`) | FULLY_FUNCTIONAL | Styled `<table class="data-table">` with XSS escaping. Full-page mode with `RenderFullHTML()` |
| **Terminal Table** (`FormatTable`) | FULLY_FUNCTIONAL | Lipgloss-styled tables in separate `table/` module. Rounded borders, alternating row colors, bold headers |

### Tree Data Formats

| Feature | Status | Notes |
|---------|--------|-------|
| **ASCII Tree** (`FormatTree`) | FULLY_FUNCTIONAL | Box-drawing characters (`├──`, `└──`, `│`). Metadata summary on nodes. `TreeRendererFromTableData()` auto-converts |
| **JSON Tree** | FULLY_FUNCTIONAL | Nested JSON with `id`, `label`, `children`, `metadata` via `JSONTreeRenderer` |
| **YAML Tree** | FULLY_FUNCTIONAL | Nested YAML structure via `YAMLTreeRenderer` |
| **HTML Tree** | FULLY_FUNCTIONAL | Nested `<ul>/<li>` list with CSS styling via `HTMLTreeRenderer`. Full-page mode available |

### Graph/Diagram Formats

| Feature | Status | Notes |
|---------|--------|-------|
| **D2 Diagrams** (`FormatD2`) | FULLY_FUNCTIONAL | Rich domain model: 20 node shapes, 11 arrow types, SQL tables with constraints, grid layouts, nested containers, reusable style classes, icons, links, tooltips |
| **Mermaid** (`FormatMermaid`) | FULLY_FUNCTIONAL | Flowchart diagrams with shape support (diamond, ellipse, circle, hexagon, cylinder, parallelogram, box) |
| **DOT/Graphviz** (`FormatDOT`) | FULLY_FUNCTIONAL | Directed and undirected graphs. Configurable graph ID, node shapes, edge styles/labels |
| **JSON Graph** | FULLY_FUNCTIONAL | `{nodes: [...], edges: [...]}` structure via `JSONGraphRenderer` |
| **YAML Graph** | FULLY_FUNCTIONAL | Same structure as JSON Graph, YAML-serialized via `YAMLGraphRenderer` |

---

## Core Data Model

| Feature | Status | Notes |
|---------|--------|-------|
| **TableData** | FULLY_FUNCTIONAL | `Headers []string` + `Rows [][]string`. Central data type shared across all table renderers |
| **tableDataBase** | FULLY_FUNCTIONAL | Embedded struct providing `SetHeaders()`, `AddRow()`, `SetData()`. Shared by JSON, YAML, HTML, Streaming renderers |
| **ToMapSlice()** | FULLY_FUNCTIONAL | Converts `TableData` to `[]map[string]string` (header→cell). Used by JSON/YAML table renderers |
| **CreateRowEdges()** | FULLY_FUNCTIONAL | Generates directed edges between consecutive rows. Used by graph renderers for `TableData`→graph conversion |
| **TreeNode** | FULLY_FUNCTIONAL | Hierarchical node with `ID`, `Label`, `Children`, `Metadata`, `Parent()`, `Depth()` |
| **GraphNode / GraphEdge** | FULLY_FUNCTIONAL | Generic graph model with `ID`, `Label`, `Shape`, `Style`, `Metadata`. Shared by DOT/Mermaid/JSON/YAML |
| **GraphRendererMixin** | FULLY_FUNCTIONAL | Shared composition for DOT and Mermaid. Provides `SetNodes()`, `SetEdges()`, `SetNodesFromTableData()`, `AddRowEdges()` |

---

## Data Shape System

| Feature | Status | Notes |
|---------|--------|-------|
| **Shape enum** (`ShapeTable`, `ShapeTree`, `ShapeGraph`) | FULLY_FUNCTIONAL | Classifies data shapes that formats can render |
| **Capability Matrix** (`formatCapabilities`) | FULLY_FUNCTIONAL | Single source of truth: maps `Format` → `[]Shape`. Query with `f.Supports(shape)` |
| **FormatsForShape()** | FULLY_FUNCTIONAL | Reverse lookup: given a shape, returns all formats that support it |
| **Format.Shapes()** | FULLY_FUNCTIONAL | Returns all shapes a format supports |

---

## Type-Safe Enums

| Feature | Status | Notes |
|---------|--------|-------|
| **Format enum** | FULLY_FUNCTIONAL | 12 format constants. `ParseFormat()`, `String()`, `IsValid()`, `AllowedValues()` |
| **Shape enum** | FULLY_FUNCTIONAL | 3 shape constants. `ParseShape()`, `String()`, `IsValid()`, `AllowedValues()` |
| **ColorMode enum** | FULLY_FUNCTIONAL | `auto`, `always`, `never`. `ParseColorMode()`, `ShouldColor()` |
| **SortBy enum** | FULLY_FUNCTIONAL | 6 sort fields (name, importance, created_at, updated_at, health, complexity) |
| **GraphShape enum** | FULLY_FUNCTIONAL | 8 node shapes (box, ellipse, diamond, circle, cylinder, hexagon, parallelogram, rect) |
| **D2Direction enum** | FULLY_FUNCTIONAL | 4 directions (down, right, left, up). Default is down |
| **D2NodeShape enum** | FULLY_FUNCTIONAL | 20 shapes (rectangle, circle, diamond, hexagon, cloud, person, queue, sql_table, class, code, etc.) |
| **D2ArrowType enum** | FULLY_FUNCTIONAL | 11 arrow types (arrow, triangle, diamond, circle, filled, box, cross, CF variants) |
| **D2Constraint enum** | FULLY_FUNCTIONAL | 3 SQL constraints (primary_key, foreign_key, unique) |
| **Alignment enum** | FULLY_FUNCTIONAL | Markdown column alignment: left, right, center |
| **enum package** | FULLY_FUNCTIONAL | Generic `Parse[T]()`, `Contains[T]()`, `AllowedValues[T]()`, `AllowedStrings[T]()`. Zero-dependency sub-module |
| **FormatCategory** | DEPRECATED | Replaced by `Shape`. `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` redirect to `Supports()` |
| **OutputFormat** | DEPRECATED | Type alias for `Format`. All `OutputFormat*` constants redirect to `Format*` |

---

## Branded IDs (Phantom Types)

| Feature | Status | Notes |
|---------|--------|-------|
| **D2NodeID / D2NodeLabel** | FULLY_FUNCTIONAL | Compile-time type safety for D2 node identifiers |
| **TreeNodeID / TreeNodeLabel** | FULLY_FUNCTIONAL | Compile-time type safety for tree node identifiers |
| **GraphNodeID / GraphNodeLabel** | FULLY_FUNCTIONAL | Compile-time type safety for graph node identifiers |
| **Generic BrandedID[T]** | FULLY_FUNCTIONAL | `NewBrandedID[Brand](value)`. Users can define custom brand types |

---

## Cross-Shape Conversion

| Feature | Status | Notes |
|---------|--------|-------|
| **TableData → Graph** | FULLY_FUNCTIONAL | `D2FromTableData()`, `DOTFromTableData()`, `MermaidFlowchartRenderer()`, `NodesFromTableData()` |
| **TableData → Tree** | FULLY_FUNCTIONAL | `TreeRendererFromTableData()` creates hierarchical tree from tabular data |
| **Tree → Graph** | FULLY_FUNCTIONAL | `D2FromTree()`, `DOTFromTree()`, `MermaidTreeRenderer()`. Generic `AddTreeNodes()` for custom renderers |
| **GraphNode → D2Node** | FULLY_FUNCTIONAL | `graphNodeToD2()`, `graphEdgeToD2()`, `graphShapeToD2()` — automatic type mapping for `SetNodes()`/`SetEdges()` |

---

## Rendering Infrastructure

| Feature | Status | Notes |
|---------|--------|-------|
| **Renderer interface** | FULLY_FUNCTIONAL | `Render() (string, error)`. Universal contract for all formats |
| **TableRenderer interface** | FULLY_FUNCTIONAL | Extends `Renderer` with `SetHeaders([]string)` and `AddRow([]string)` |
| **TreeOutputRenderer interface** | FULLY_FUNCTIONAL | Extends `Renderer` with `SetRoot(*TreeNode)` |
| **GraphRenderer interface** | FULLY_FUNCTIONAL | Extends `Renderer` with `SetNodes([]GraphNode)` and `SetEdges([]GraphEdge)` |
| **StreamingRenderer interface** | FULLY_FUNCTIONAL | `Stream(io.Writer) error`. Incremental output for large datasets |
| **StreamingHTMLRenderer** | FULLY_FUNCTIONAL | True streaming HTML table output. Writes chunks incrementally |
| **StreamingRendererFromRenderer()** | FULLY_FUNCTIONAL | Adapter wrapping standard `Renderer` as `StreamingRenderer` (collects then writes) |
| **MustRender()** | FULLY_FUNCTIONAL | `MustRender(r Renderer) string` — panics on error. For tests/examples |
| **RenderTableData()** | FULLY_FUNCTIONAL | Unified dispatcher: renders `TableData` in any format to `io.Writer` (defaults to stdout). Accepts `RenderOptions` (Title, GraphID, Writer) |

---

## Writer APIs (Streaming I/O)

| Feature | Status | Notes |
|---------|--------|-------|
| **CSVWriter** | FULLY_FUNCTIONAL | `WriteHeader()`, `WriteRow()`, `WriteRows()`, `Flush()`, `Error()` |
| **TSVWriter** | FULLY_FUNCTIONAL | Same API as CSVWriter with tab delimiter |
| **JSONWriter** | FULLY_FUNCTIONAL | `Encode(v any) error` — streaming JSON encoder with indentation |
| **XMLWriter** | FULLY_FUNCTIONAL | `WriteHeader()`, `WriteRow()`, `WriteRows()`, `WriteFooter()` |
| **DelimitedWriter** | FULLY_FUNCTIONAL | Shared base for CSV/TSV writers. Configurable delimiter |

---

## Serialization Helpers

| Feature | Status | Notes |
|---------|--------|-------|
| **MarshalJSON / UnmarshalJSON** | FULLY_FUNCTIONAL | Wrapper over `encoding/json` with type-context errors |
| **MarshalJSONIndent** | FULLY_FUNCTIONAL | Indented JSON with configurable prefix/indent |
| **MarshalYAML / UnmarshalYAML** | FULLY_FUNCTIONAL | Wrapper over `go-faster/yaml` with type-context errors |
| **MarshalXML / MarshalXMLIndent** | FULLY_FUNCTIONAL | Wrapper over `encoding/xml` with type-context errors |
| **MarshalCSVFromTableData** | FULLY_FUNCTIONAL | One-shot CSV from `TableData` |
| **MarshalTSV / MarshalTSVFromTableData** | FULLY_FUNCTIONAL | One-shot TSV from `TableData` or raw data |
| **MarshalXMLFromTableData** | FULLY_FUNCTIONAL | One-shot XML from `TableData` |

---

## Registry System

| Feature | Status | Notes |
|---------|--------|-------|
| **Register()** | FULLY_FUNCTIONAL | Register a `RendererFactory` for a `Format`. Thread-safe. Returns error on duplicate |
| **Create()** | FULLY_FUNCTIONAL | Create a `Renderer` instance from registered factory |
| **Unregister()** | FULLY_FUNCTIONAL | Remove a format from the registry |
| **RegisteredFormats()** | FULLY_FUNCTIONAL | Returns sorted list of registered formats |
| **IsRegistered()** | FULLY_FUNCTIONAL | Check if a format has a registered renderer |

---

## Color & Terminal Detection

| Feature | Status | Notes |
|---------|--------|-------|
| **ColorMode.Auto** | FULLY_FUNCTIONAL | Checks TTY, `NO_COLOR`, `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`, `GO_OUTPUT_FORCE_COLOR`, `FORCE_COLOR` |
| **ColorMode.Always / Never** | FULLY_FUNCTIONAL | Force enable/disable |

---

## Escape/Sanitization

| Feature | Status | Notes |
|---------|--------|-------|
| **escape.HTML()** | FULLY_FUNCTIONAL | Uses `html.EscapeString` from stdlib |
| **escape.XML()** | FULLY_FUNCTIONAL | Uses `html.EscapeString` + `&apos;` for apostrophe |
| **escape.D2()** | FULLY_FUNCTIONAL | Escapes `\`, `"`, `\n`, `\t` for D2 diagram strings |
| **escape.DOT()** | FULLY_FUNCTIONAL | Escapes `\`, `"`, `\n` for DOT/Graphviz strings |
| **escape.MermaidID()** | FULLY_FUNCTIONAL | Sanitizes for Mermaid node identifiers (alphanumeric + underscore) |
| **escape.MermaidSlug()** | FULLY_FUNCTIONAL | Fallback slug sanitization (spaces/hyphens/slashes → underscores) |
| **escape.MermaidText()** | FULLY_FUNCTIONAL | Escapes brackets, braces, quotes, newlines for Mermaid labels |

---

## Testing Infrastructure

| Feature | Status | Notes |
|---------|--------|-------|
| **testhelpers package** | FULLY_FUNCTIONAL | Zero-dep, publicly importable. `AssertStringSliceEqual()`, `AssertContains()`, `AssertEqual[T]()`, `TestEnumIsValid[T]()`, `TestStructFields()`, `StringField()`, `IntField()` |
| **internal/gentest** | FULLY_FUNCTIONAL | Root-only helpers: `AssertOutputContains()`, `AssertValidJSON()`, `AssertValidYAML()`, `AssertHTMLEscape()`, `AssertMarshalError()`, `ExpectedOutput` |
| **Fuzz tests** | FULLY_FUNCTIONAL | `FuzzCSVWriter`, `FuzzMarkdownTable` — seed corpus + coverage-guided fuzzing |
| **Benchmarks** | FULLY_FUNCTIONAL | `BenchmarkASCIITreeRenderer`, `BenchmarkHTMLRenderer`, `BenchmarkMermaidRenderer`, `BenchmarkDOTRenderer`, `BenchmarkCSVWriter`, `BenchmarkMarkdownTable`, `BenchmarkTableDataCreateRowEdges` |
| **Integration tests** | FULLY_FUNCTIONAL | Cross-module tests in `integration/` package. Tests all 12 formats, streaming, tree depth, edge creation, large datasets |
| **User journey tests** | FULLY_FUNCTIONAL | End-to-end tests simulating CLI developer workflows in `userjourney_test.go` |

---

## Multi-Module Architecture

| Feature | Status | Notes |
|---------|--------|-------|
| **Root module** (`package output`) | FULLY_FUNCTIONAL | Core types + formatters. Zero lipgloss dependency |
| **enum/** | FULLY_FUNCTIONAL | Generic enum utilities. Zero dependencies |
| **escape/** | FULLY_FUNCTIONAL | Format-specific escaping. Zero dependencies |
| **testhelpers/** | FULLY_FUNCTIONAL | Shared test assertions. Zero dependencies, publicly importable |
| **sort/** | DEPRECATED | Only `ByField[T]()` helper remains. Use `slices.SortStableFunc` + `cmp.Compare` instead |
| **table/** | FULLY_FUNCTIONAL | Lipgloss terminal tables. Isolated from root module |
| **integration/** | FULLY_FUNCTIONAL | Cross-module integration tests |
| **examples/** | FULLY_FUNCTIONAL | Working examples demonstrating all 12 formats |
| **go.work** | FULLY_FUNCTIONAL | Gitignored. `go.work.example` provided for local development |

---

## CI/CD

| Feature | Status | Notes |
|---------|--------|-------|
| **Build & Test** | FULLY_FUNCTIONAL | Per-module build and test with `-race` flag |
| **Lint** | FULLY_FUNCTIONAL | `golangci-lint` v2.12 across all modules |
| **govulncheck** | FULLY_FUNCTIONAL | Vulnerability scanning across all modules |
| **go mod tidy check** | FULLY_FUNCTIONAL | Verifies all module `go.mod` files are tidy |
| **Nix flake** | FULLY_FUNCTIONAL | Dev shell with Go 1.26.2, golangci-lint, gopls. Uses `flake-parts` + `treefmt-nix` + `git-hooks.nix` |
| **Pre-commit hooks** | FULLY_FUNCTIONAL | Auto-installed via `nix develop`. Also `.pre-commit-config.yaml` for non-Nix users |

---

## Documentation

| Feature | Status | Notes |
|---------|--------|-------|
| **README.md** | KNOWN_ISSUE | References `NewD2Renderer("Architecture")` which does not exist. Only `NewD2Diagram()` exists |
| **CHANGELOG.md** | FULLY_FUNCTIONAL | Version history |
| **CONTRIBUTING.md** | FULLY_FUNCTIONAL | Contribution guidelines |
| **ADR 001** | FULLY_FUNCTIONAL | Multi-module workspace decision |
| **ADR 002** | FULLY_FUNCTIONAL | Shape capability matrix decision |
| **DOMAIN_LANGUAGE.md** | FULLY_FUNCTIONAL | Domain vocabulary |
| **FORMAT_ARCHITECTURE.md** | FULLY_FUNCTIONAL | Format architecture documentation |

---

## Utility

| Feature | Status | Notes |
|---------|--------|-------|
| **FilledStrings()** | FULLY_FUNCTIONAL | Creates a slice of n identical strings. Used in benchmarks and test setup |
| **InvalidFormatError** | FULLY_FUNCTIONAL | Descriptive error with allowed values list |
| **UnsupportedFormatError** | FULLY_FUNCTIONAL | Returned by `RenderTableData()` for table/json formats |
| **D2 utility methods** | FULLY_FUNCTIONAL | `AddNodeSimple()`, `AddNodeWithShape()`, `AddEdgeSimple()`, `AddLabeledEdge()` — builder pattern shortcuts |

---

**Last audited:** 2026-05-23
**Total features:** 112
**Fully functional:** 108
**Deprecated:** 3 (FormatCategory, OutputFormat, sort package)
**Known issues:** 1 (README references non-existent `NewD2Renderer`)
