# Features

Complete feature inventory for `go-output` — a Go library providing consistent output formatting across 16 formats for CLI applications.

**Status legend:** FULLY_FUNCTIONAL | PARTIALLY_FUNCTIONAL | DEPRECATED | KNOWN_ISSUE

---

## Output Formats

### Table Data Formats

| Feature                               | Status           | Notes                                                                                                     |
| ------------------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------- |
| **JSON Table** (`FormatJSON`)         | FULLY_FUNCTIONAL | Array-of-objects via `JSONTableRenderer`. Headers become keys, rows become values                         |
| **CSV** (`FormatCSV`)                 | FULLY_FUNCTIONAL | Streaming writer via `CSVWriter`. Auto-quotes fields with commas/newlines                                 |
| **TSV** (`FormatTSV`)                 | FULLY_FUNCTIONAL | Streaming writer via `TSVWriter`. Tab-separated with type-switch marshaling (`[][]string`, `[]string`)    |
| **Markdown Table** (`FormatMarkdown`) | FULLY_FUNCTIONAL | Column alignment (left/right/center). Auto-calculated column widths via `MarkdownTable`                   |
| **XML** (`FormatXML`)                 | FULLY_FUNCTIONAL | Structured `<table><headers>...</headers><rows>...</rows></table>` with XML escaping                      |
| **YAML Table** (`FormatYAML`)         | FULLY_FUNCTIONAL | Sequence of mappings via `YAMLTableRenderer`. Uses `go-faster/yaml`                                       |
| **HTML Table** (`FormatHTML`)         | FULLY_FUNCTIONAL | Styled `<table class="data-table">` with XSS escaping. Full-page mode with `RenderFullHTML()`             |
| **JSONL** (`FormatJSONL`)             | FULLY_FUNCTIONAL | JSON Lines — one JSON object per line via `JSONLTableRenderer`. Streaming via `JSONLWriter`               |
| **AsciiDoc** (`FormatAsciiDoc`)       | FULLY_FUNCTIONAL | AsciiDoc tables with `\|===` borders via `AsciiDocTableRenderer`. Pipe escaping for cell content          |
| **TOML** (`FormatTOML`)               | FULLY_FUNCTIONAL | TOML serialization via `TOMLTableRenderer`, `TOMLTreeRenderer`, `TOMLGraphRenderer`. Uses `go-toml/v2`    |
| **Terminal Table** (`FormatTable`)    | FULLY_FUNCTIONAL | Lipgloss-styled tables in separate `table/` module. Rounded borders, alternating row colors, bold headers |

### Tree Data Formats

| Feature                       | Status           | Notes                                                                                                              |
| ----------------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------ |
| **ASCII Tree** (`FormatTree`) | FULLY_FUNCTIONAL | Box-drawing characters (`├──`, `└──`, `│`). Metadata summary on nodes. `TreeRendererFromTableData()` auto-converts |
| **JSON Tree**                 | FULLY_FUNCTIONAL | Nested JSON with `id`, `label`, `children`, `metadata` via `JSONTreeRenderer`                                      |
| **TOML Tree**                 | FULLY_FUNCTIONAL | Nested TOML structure via `TOMLTreeRenderer`. Uses shared treeNodeDTO with json+yaml+toml tags                     |
| **YAML Tree**                 | FULLY_FUNCTIONAL | Nested YAML structure via `YAMLTreeRenderer`                                                                       |
| **HTML Tree**                 | FULLY_FUNCTIONAL | Nested `<ul>/<li>` list with CSS styling via `HTMLTreeRenderer`. Full-page mode available                          |

### Graph/Diagram Formats

| Feature                         | Status           | Notes                                                                                                                                                           |
| ------------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **D2 Diagrams** (`FormatD2`)    | FULLY_FUNCTIONAL | Rich domain model: 20 node shapes, 11 arrow types, SQL tables with constraints, grid layouts, nested containers, reusable style classes, icons, links, tooltips |
| **Mermaid** (`FormatMermaid`)   | FULLY_FUNCTIONAL | Flowchart diagrams with shape support (diamond, ellipse, circle, hexagon, cylinder, parallelogram, box)                                                         |
| **DOT/Graphviz** (`FormatDOT`)  | FULLY_FUNCTIONAL | Directed and undirected graphs. Configurable graph ID, node shapes, edge styles/labels                                                                          |
| **JSON Graph**                  | FULLY_FUNCTIONAL | `{nodes: [...], edges: [...]}` structure via `JSONGraphRenderer`                                                                                                |
| **YAML Graph**                  | FULLY_FUNCTIONAL | Same structure as JSON Graph, YAML-serialized via `YAMLGraphRenderer`                                                                                           |
| **TOML Graph**                  | FULLY_FUNCTIONAL | Same structure as JSON Graph, TOML-serialized via `TOMLGraphRenderer`                                                                                           |
| **PlantUML** (`FormatPlantUML`) | FULLY_FUNCTIONAL | Component diagrams via `PlantUMLDiagram`. Uses `GraphRendererState`. Supports TableData→graph and Tree→graph conversion                                         |

---

## Core Data Model

| Feature                   | Status           | Notes                                                                                                                                                                    |
| ------------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **TableData**             | FULLY_FUNCTIONAL | `Headers []string` + `Rows [][]string` + `Footer []string`. Central data type shared across all table renderers. Footer renders as totals/summary row in tabular formats |
| **TableDataStore**        | FULLY_FUNCTIONAL | Exported embedded struct providing `SetHeaders()`, `AddRow()`, `SetData()`, `Data()`, `SetFooter()`. Shared by JSON, YAML, TOML, HTML, AsciiDoc, Streaming renderers     |
| **ToMapSlice()**          | FULLY_FUNCTIONAL | Converts `TableData` to `[]map[string]string` (header→cell). Used by JSON/YAML table renderers                                                                           |
| **CreateRowEdges()**      | FULLY_FUNCTIONAL | Generates directed edges between consecutive rows. Used by graph renderers for `TableData`→graph conversion                                                              |
| **TreeNode**              | FULLY_FUNCTIONAL | Hierarchical node with `ID`, `Label`, `Children`, `Metadata`, `Parent()`, `Depth()`                                                                                      |
| **GraphNode / GraphEdge** | FULLY_FUNCTIONAL | Generic graph model with `ID`, `Label`, `Shape`, `Style`, `Metadata`. Shared by DOT/Mermaid/JSON/YAML                                                                    |
| **GraphRendererState**    | FULLY_FUNCTIONAL | Shared composition for all graph renderers (DOT, Mermaid, JSON, YAML, TOML, PlantUML). Provides `SetNodes()`, `SetEdges()`, `SetNodesFromTableData()`, `AddRowEdges()`   |

---

## Data Shape System

| Feature                                                  | Status           | Notes                                                                             |
| -------------------------------------------------------- | ---------------- | --------------------------------------------------------------------------------- |
| **Shape enum** (`ShapeTable`, `ShapeTree`, `ShapeGraph`) | FULLY_FUNCTIONAL | Classifies data shapes that formats can render                                    |
| **Capability Matrix** (`formatCapabilities`)             | FULLY_FUNCTIONAL | Single source of truth: maps `Format` → `[]Shape`. Query with `f.Supports(shape)` |
| **FormatsForShape()**                                    | FULLY_FUNCTIONAL | Reverse lookup: given a shape, returns all formats that support it                |
| **Format.Shapes()**                                      | FULLY_FUNCTIONAL | Returns all shapes a format supports                                              |

---

## Type-Safe Enums

| Feature               | Status           | Notes                                                                                                          |
| --------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------- |
| **Format enum**       | FULLY_FUNCTIONAL | 16 format constants. `ParseFormat()`, `String()`, `IsValid()`, `AllowedValues()`                               |
| **Shape enum**        | FULLY_FUNCTIONAL | 3 shape constants. `ParseShape()`, `String()`, `IsValid()`, `AllowedValues()`                                  |
| **ColorMode enum**    | FULLY_FUNCTIONAL | `auto`, `always`, `never`. `ParseColorMode()`, `ShouldColor()`. Wired into table, tree, markdown renderers     |
| **SortBy enum**       | REMOVED          | Deleted — zero external callers. Use `slices.SortStableFunc` + `cmp.Compare` (stdlib)                          |
| **NodeShape enum**    | FULLY_FUNCTIONAL | 8 node shapes (box, ellipse, diamond, circle, cylinder, hexagon, parallelogram, rect)                          |
| **D2Direction enum**  | FULLY_FUNCTIONAL | 4 directions (down, right, left, up). Default is down                                                          |
| **D2NodeShape enum**  | FULLY_FUNCTIONAL | 20 shapes (rectangle, circle, diamond, hexagon, cloud, person, queue, sql_table, class, code, etc.)            |
| **D2ArrowType enum**  | FULLY_FUNCTIONAL | 11 arrow types (arrow, triangle, diamond, circle, filled, box, cross, CF variants)                             |
| **D2Constraint enum** | FULLY_FUNCTIONAL | 3 SQL constraints (primary_key, foreign_key, unique)                                                           |
| **Alignment enum**    | FULLY_FUNCTIONAL | Markdown column alignment: left, right, center                                                                 |
| **enum package**      | FULLY_FUNCTIONAL | Generic `Parse[T]()`, `Contains[T]()`, `AllowedValues[T]()`, `AllowedStrings[T]()`. Zero-dependency sub-module |
| **FormatCategory**    | REMOVED          | Replaced by `Shape`. `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` removed              |
| **OutputFormat**      | REMOVED          | Type alias for `Format`. All `OutputFormat*` constants removed                                                 |

---

## Branded IDs (Phantom Types)

| Feature                          | Status           | Notes                                                             |
| -------------------------------- | ---------------- | ----------------------------------------------------------------- |
| **D2NodeID / D2NodeLabel**       | FULLY_FUNCTIONAL | Compile-time type safety for D2 node identifiers                  |
| **TreeNodeID / TreeNodeLabel**   | FULLY_FUNCTIONAL | Compile-time type safety for tree node identifiers                |
| **GraphNodeID / GraphNodeLabel** | FULLY_FUNCTIONAL | Compile-time type safety for graph node identifiers               |
| **Generic BrandedID[T]**         | FULLY_FUNCTIONAL | `NewBrandedID[Brand](value)`. Users can define custom brand types |

---

## Cross-Shape Conversion

| Feature                | Status           | Notes                                                                                                                     |
| ---------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------- |
| **TableData → Graph**  | FULLY_FUNCTIONAL | `D2FromTableData()`, `DOTFromTableData()`, `MermaidFromTableData()`, `PlantUMLFromTableData()`, `NodesFromTableData()`    |
| **TableData → Tree**   | FULLY_FUNCTIONAL | `TreeRendererFromTableData()` creates hierarchical tree from tabular data                                                 |
| **Tree → Graph**       | FULLY_FUNCTIONAL | `D2FromTree()`, `DOTFromTree()`, `MermaidFromTree()`, `PlantUMLFromTree()`. Generic `AddTreeNodes()` for custom renderers |
| **GraphNode → D2Node** | FULLY_FUNCTIONAL | `graphNodeToD2()`, `graphEdgeToD2()`, `graphShapeToD2()` — automatic type mapping for `SetNodes()`/`SetEdges()`           |

---

## Rendering Infrastructure

| Feature                             | Status           | Notes                                                                                                                                                                                                                                                                  |
| ----------------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Renderer interface**              | FULLY_FUNCTIONAL | `Render() (string, error)`. Universal contract for all formats                                                                                                                                                                                                         |
| **TableRenderer interface**         | FULLY_FUNCTIONAL | Extends `Renderer` with `SetHeaders([]string)` and `AddRow([]string)`                                                                                                                                                                                                  |
| **TreeOutputRenderer interface**    | FULLY_FUNCTIONAL | Extends `Renderer` with `SetRoot(*TreeNode)`                                                                                                                                                                                                                           |
| **GraphRenderer interface**         | FULLY_FUNCTIONAL | Extends `Renderer` with `SetNodes([]GraphNode)` and `SetEdges([]GraphEdge)`                                                                                                                                                                                            |
| **StreamingRenderer interface**     | FULLY_FUNCTIONAL | `Stream(io.Writer) error`. Incremental output for large datasets                                                                                                                                                                                                       |
| **StreamingHTMLRenderer**           | FULLY_FUNCTIONAL | True streaming HTML table output. Writes chunks incrementally                                                                                                                                                                                                          |
| **StreamingRendererFromRenderer()** | FULLY_FUNCTIONAL | Adapter wrapping standard `Renderer` as `StreamingRenderer` (collects then writes)                                                                                                                                                                                     |
| **MustRender()**                    | FULLY_FUNCTIONAL | `MustRender(r Renderer) string` — panics on error. For tests/examples                                                                                                                                                                                                  |
| **RenderTableData()**               | FULLY_FUNCTIONAL | Unified dispatcher: renders `TableData` in any format to `io.Writer` (defaults to stdout). Accepts `RenderOptions` (Title, GraphID, Writer)                                                                                                                            |
| **Footer row**                      | FULLY_FUNCTIONAL | `TableData.Footer []string` — optional totals/summary row. CSV/TSV append as last row. HTML uses `<tfoot>` with `footer-cell` class. XML uses `<footer>`. Markdown adds separator + bold row. Table uses bold styling. Data formats (JSON/YAML/TOML/JSONL) skip footer |
| **TableData.Validate()**            | FULLY_FUNCTIONAL | Validates footer column count matches headers. Returns `errColumnMismatch` when counts differ                                                                                                                                                                          |

---

## Writer APIs (Streaming I/O)

| Feature             | Status           | Notes                                                                               |
| ------------------- | ---------------- | ----------------------------------------------------------------------------------- |
| **CSVWriter**       | FULLY_FUNCTIONAL | `WriteHeader()`, `WriteRow()`, `WriteRows()`, `WriteFooter()`, `Flush()`, `Error()` |
| **TSVWriter**       | FULLY_FUNCTIONAL | Same API as CSVWriter with tab delimiter                                            |
| **JSONWriter**      | FULLY_FUNCTIONAL | `Encode(v any) error` — streaming JSON encoder with indentation                     |
| **JSONLWriter**     | FULLY_FUNCTIONAL | `Encode(v any) error`, `Flush() error` — streaming JSON Lines encoder               |
| **XMLWriter**       | FULLY_FUNCTIONAL | `WriteHeader()`, `WriteRow()`, `WriteRows()`, `WriteFooter()`                       |
| **DelimitedWriter** | FULLY_FUNCTIONAL | Shared base for CSV/TSV writers. Configurable delimiter                             |

---

## Serialization Helpers

| Feature                           | Status           | Notes                                                             |
| --------------------------------- | ---------------- | ----------------------------------------------------------------- |
| **MarshalJSON / UnmarshalJSON**   | FULLY_FUNCTIONAL | Wrapper over `encoding/json` with type-context errors             |
| **MarshalJSONIndent**             | FULLY_FUNCTIONAL | Indented JSON with configurable prefix/indent                     |
| **MarshalYAML / UnmarshalYAML**   | FULLY_FUNCTIONAL | Wrapper over `go-faster/yaml` with type-context errors            |
| **MarshalXML / MarshalXMLIndent** | FULLY_FUNCTIONAL | Wrapper over `encoding/xml` with type-context errors              |
| **MarshalCSVFromTableData**       | FULLY_FUNCTIONAL | One-shot CSV from `TableData`                                     |
| **MarshalTSVFromTableData**       | FULLY_FUNCTIONAL | One-shot TSV from `TableData`; `TSVWriter` for streaming raw rows |
| **MarshalXMLFromTableData**       | FULLY_FUNCTIONAL | One-shot XML from `TableData`                                     |
| **MarshalTOML / UnmarshalTOML**   | FULLY_FUNCTIONAL | Wrapper over `go-toml/v2` with type-context errors                |
| **MarshalJSONLFromTableData**     | FULLY_FUNCTIONAL | One-shot JSON Lines from `TableData`                              |
| **MarshalAsciiDocFromTableData**  | FULLY_FUNCTIONAL | One-shot AsciiDoc table from `TableData`                          |
| **MarshalTOMLFromTableData**      | FULLY_FUNCTIONAL | One-shot TOML from `TableData`                                    |

---

## Registry System

| Feature                 | Status  | Notes                                                                               |
| ----------------------- | ------- | ----------------------------------------------------------------------------------- |
| **Register()**          | REMOVED | Deprecated renderer registry deleted. Use direct constructors (`d2.NewD2Diagram()`) |
| **Create()**            | REMOVED | Removed with registry. Use format-specific constructors                             |
| **Unregister()**        | REMOVED | Removed with registry                                                               |
| **RegisteredFormats()** | REMOVED | Removed with registry                                                               |
| **IsRegistered()**      | REMOVED | Removed with registry                                                               |

---

## Color & Terminal Detection

| Feature                      | Status           | Notes                                                                                                                           |
| ---------------------------- | ---------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| **ColorMode.Auto**           | FULLY_FUNCTIONAL | Checks TTY, `NO_COLOR`, `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`, `GO_OUTPUT_FORCE_COLOR`, `FORCE_COLOR` |
| **ColorMode.Always / Never** | FULLY_FUNCTIONAL | Force enable/disable                                                                                                            |

### ColorMode Support by Format

| Format                        | ColorMode | Behavior                                                                         |
| ----------------------------- | --------- | -------------------------------------------------------------------------------- |
| Table (lipgloss)              | ✅        | Bold headers, alternating row colors, border styling via `table.WithColorMode()` |
| Tree                          | ✅        | Depth-based ANSI color cycling, bold labels, dim connectors via `SetColorMode()` |
| Markdown                      | ✅        | Bold headers, dim separators via `SetColorMode()`                                |
| CSV / TSV                     | ❌        | Plain text — no ANSI support                                                     |
| JSON / YAML / TOML / JSONL    | ❌        | Structured data — no ANSI support                                                |
| XML / HTML / AsciiDoc         | ❌        | Markup formats — styling via markup syntax, not ANSI                             |
| D2 / Mermaid / DOT / PlantUML | ❌        | Diagram formats — styling via diagram syntax                                     |

---

## Escape/Sanitization

| Feature                  | Status           | Notes                                                                         |
| ------------------------ | ---------------- | ----------------------------------------------------------------------------- |
| **escape.HTML()**        | FULLY_FUNCTIONAL | Uses `html.EscapeString` from stdlib                                          |
| **escape.XML()**         | FULLY_FUNCTIONAL | Uses `html.EscapeString` + `&apos;` for apostrophe                            |
| **escape.D2()**          | FULLY_FUNCTIONAL | Escapes `\`, `"`, `\n`, `\t` for D2 diagram strings                           |
| **escape.DOT()**         | FULLY_FUNCTIONAL | Escapes `\`, `"`, `\n` for DOT/Graphviz strings                               |
| **escape.MermaidID()**   | FULLY_FUNCTIONAL | Sanitizes for Mermaid node identifiers (alphanumeric + underscore)            |
| **escape.MermaidSlug()** | FULLY_FUNCTIONAL | Fallback slug sanitization (spaces/hyphens/slashes → underscores)             |
| **escape.MermaidText()** | FULLY_FUNCTIONAL | Escapes brackets, braces, quotes, newlines for Mermaid labels                 |
| **escape.SlugifyID()**   | FULLY_FUNCTIONAL | Sanitizes strings for diagram node identifiers across D2/DOT/Mermaid/PlantUML |
| **escape.PlantUML()**    | FULLY_FUNCTIONAL | Escapes `]`, newline, backslash, quote for PlantUML labels                    |

---

## Progress Visualization (NOM & TUI)

### NOM-style Real-time Progress (`nom/` module)

| Feature                              | Status           | Notes                                                                                                                                                                |
| ------------------------------------ | ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NOMStyleSubscriber**               | FULLY_FUNCTIONAL | Event-driven subscriber implementing `EventSubscriber`. Routes string-based events without sharing concrete types                                                    |
| **DependencyTree**                   | FULLY_FUNCTIONAL | Hierarchical activity visualization. Priority filtering (Running > Failed > Pending > Completed), depth-aware prefixes                                               |
| **InlineRenderer**                   | FULLY_FUNCTIONAL | Real-time inline terminal renderer. Snapshot-based rendering (race-free); Start/Stop/Finish lifecycle; config setters (`SetStartTime`/`SetAppName`/`SetNoColor`/`SetHideCursor`/`SetMaxHeight`) safe during the render loop; CI plain-text degradation (skips cursor codes under CI); height-pressure collapse marker (`⋯ N completed`); no-color mode; ANSI redraw                                                           |
| **TimingCache**                      | FULLY_FUNCTIONAL | Persists activity durations as CSV at `~/.cache/nom-timing.csv`. Serialized saves (saveMu), caps 10 entries/activity, applies cap on load                            |
| **ActivityStatus enum**              | FULLY_FUNCTIONAL | 4 states: Pending, Running, Completed, Failed (with symbol/color/shape mapping)                                                                                      |
| **Branded IDs**                      | FULLY_FUNCTIONAL | `ActivityID`, `ActivityName`, `WorkflowID`, `WorkflowName` — named types over `string` (phantom-branding upgrade tracked in TODO §E)                                 |
| **Event accessor interfaces**        | FULLY_FUNCTIONAL | `WorkflowEventAccessor`, `ActivityEventAccessor`, `DurationAccessor`, `ErrorAccessor` — type-assertion routing                                                       |
| **Activity symbols**                 | FULLY_FUNCTIONAL | `SymbolRunning`, `SymbolCompleted`, `SymbolFailed`, `SymbolPending`, `SymbolDownload`, `SymbolUpload`, `SymbolTiming`, `SymbolAverage`, `SymbolTotal`, `SymbolPhase` |
| **Lazy build (double-checked lock)** | FULLY_FUNCTIONAL | `DependencyTree.Build()` uses the `loaded` flag with double-checked locking to prevent rebuild under read lock                                                       |

### Bubble Tea Interactive TUI (`tui/` module)

| Feature                       | Status           | Notes                                                                                                           |
| ----------------------------- | ---------------- | --------------------------------------------------------------------------------------------------------------- |
| **BubbleTeaProgressReporter** | FULLY_FUNCTIONAL | Reports progress to a Bubble Tea program. Lazy TUI start via double-checked locking                             |
| **ProgressModel**             | FULLY_FUNCTIONAL | Bubble Tea `Model`: `Init()`/`Update()`/`View()`. Handles `TickMsg`, `ProgressUpdateMsg`, `CancelMsg`           |
| **WorkflowState machine**     | FULLY_FUNCTIONAL | Enforces `Idle → Running → Completed/Errored`. Terminal states reject updates/ticks                             |
| **DisplayModeNOM**            | FULLY_FUNCTIONAL | Renders NOM dependency tree + activity status                                                                   |
| **DisplayModeUniversal**      | FULLY_FUNCTIONAL | Renders step-based progress (like `nh darwin switch`)                                                           |
| **Reporter API**              | FULLY_FUNCTIONAL | `Start()`, `Stop()`, `ReportStep()`, `ReportProgress()`, `ReportError()`, `ReportMessage()`, `SetDisplayMode()` |

---

## Testing Infrastructure

| Feature                 | Status           | Notes                                                                                                                                                                                                                                    |
| ----------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **testhelpers package** | FULLY_FUNCTIONAL | Zero-dep, publicly importable. `AssertStringSliceEqual()`, `AssertContains()`, `AssertEqual[T]()`, `TestEnumIsValid[T]()`, `TestStructFields()`, `StringField()`, `IntField()`                                                           |
| **internal/gentest**    | FULLY_FUNCTIONAL | Root-only helpers: `AssertOutputContains()`, `AssertValidJSON()`, `AssertValidYAML()`, `AssertHTMLEscape()`, `AssertMarshalError()`, `ExpectedOutput`                                                                                    |
| **Fuzz tests**          | FULLY_FUNCTIONAL | `FuzzMarkdownTable` — seed corpus + coverage-guided fuzzing                                                                                                                                                                              |
| **Benchmarks**          | FULLY_FUNCTIONAL | `BenchmarkASCIITreeRenderer`, `BenchmarkHTMLRenderer`, `BenchmarkMermaidRenderer`, `BenchmarkDOTRenderer`, `BenchmarkCSVWriter`, `BenchmarkMarkdownTableColored`, `BenchmarkMarkdownTableWithFooter`, `BenchmarkTableDataCreateRowEdges` |
| **Integration tests**   | FULLY_FUNCTIONAL | Cross-module tests in `integration/` package. Tests all 16 formats, streaming, tree depth, edge creation, large datasets                                                                                                                 |
| **User journey tests**  | FULLY_FUNCTIONAL | End-to-end tests simulating CLI developer workflows in `userjourney_test.go`                                                                                                                                                             |

---

## Multi-Module Architecture

| Feature                            | Status           | Notes                                                                         |
| ---------------------------------- | ---------------- | ----------------------------------------------------------------------------- |
| **Root module** (`package output`) | FULLY_FUNCTIONAL | Core types + formatters. Zero lipgloss dependency                             |
| **enum/**                          | FULLY_FUNCTIONAL | Generic enum utilities. Zero dependencies                                     |
| **escape/**                        | FULLY_FUNCTIONAL | Format-specific escaping. Zero dependencies                                   |
| **envdetect/**                     | FULLY_FUNCTIONAL | Shared CI/NO_COLOR env detection. Zero dependencies                           |
| **testhelpers/**                   | FULLY_FUNCTIONAL | Shared test assertions. Zero dependencies, publicly importable                |
| **testhelpers/graphtest/**         | FULLY_FUNCTIONAL | Shared graph test fixtures. Zero dependencies                                 |
| **table/**                         | FULLY_FUNCTIONAL | Lipgloss terminal tables. Isolated from root module                           |
| **markdown/**                      | FULLY_FUNCTIONAL | Markdown table renderer. Self-registers via `init()`. Zero deps beyond root   |
| **tree/**                          | FULLY_FUNCTIONAL | ASCII tree renderer (box-drawing, color cycling). Self-registers via `init()` |
| **integration/**                   | FULLY_FUNCTIONAL | Cross-module integration tests                                                |
| **examples/**                      | FULLY_FUNCTIONAL | Working examples demonstrating all 16 formats                                 |
| **bdd/**                           | FULLY_FUNCTIONAL | BDD test suite (Ginkgo/Gomega). Test-only module                              |
| **delimited/**                     | FULLY_FUNCTIONAL | CSV/TSV writers and marshalers. Isolated from root module                     |
| **d2/**                            | FULLY_FUNCTIONAL | D2 diagram builder. Isolated from root module                                 |
| **graph/**                         | FULLY_FUNCTIONAL | DOT and Mermaid renderers. Isolated from root module                          |
| **markup/**                        | FULLY_FUNCTIONAL | HTML, XML, AsciiDoc renderers. Isolated from root module                      |
| **plantuml/**                      | FULLY_FUNCTIONAL | PlantUML diagram renderer. Isolated from root module                          |
| **serialization/**                 | FULLY_FUNCTIONAL | JSON, YAML, TOML, JSONL renderers. Isolated from root module                  |
| **nom/**                           | FULLY_FUNCTIONAL | NOM-style real-time progress (dependency trees, timing cache). Lipgloss-only  |
| **tui/**                           | FULLY_FUNCTIONAL | Bubble Tea interactive TUI (depends on nom + bubbletea + lipgloss)            |
| **go.work**                        | FULLY_FUNCTIONAL | Gitignored. `go.work.example` provided for local development                  |

---

## CI/CD

| Feature               | Status           | Notes                                                                                                |
| --------------------- | ---------------- | ---------------------------------------------------------------------------------------------------- |
| **Build & Test**      | FULLY_FUNCTIONAL | Per-module build and test with `-race` flag                                                          |
| **Lint**              | FULLY_FUNCTIONAL | `golangci-lint` v2.12 across all modules                                                             |
| **govulncheck**       | FULLY_FUNCTIONAL | Vulnerability scanning across all modules                                                            |
| **go mod tidy check** | FULLY_FUNCTIONAL | Verifies all module `go.mod` files are tidy                                                          |
| **Nix flake**         | FULLY_FUNCTIONAL | Dev shell with Go 1.26.2, golangci-lint, gopls. Uses `flake-parts` + `treefmt-nix` + `git-hooks.nix` |
| **Pre-commit hooks**  | FULLY_FUNCTIONAL | Auto-installed via `nix develop`. Also `.pre-commit-config.yaml` for non-Nix users                   |

---

## Documentation

| Feature                    | Status           | Notes                                                      |
| -------------------------- | ---------------- | ---------------------------------------------------------- |
| **README.md**              | FULLY_FUNCTIONAL | All examples verified correct after deprecated API removal |
| **CHANGELOG.md**           | FULLY_FUNCTIONAL | Version history                                            |
| **CONTRIBUTING.md**        | FULLY_FUNCTIONAL | Contribution guidelines                                    |
| **ADR 001**                | FULLY_FUNCTIONAL | Multi-module workspace decision                            |
| **ADR 002**                | FULLY_FUNCTIONAL | Shape capability matrix decision                           |
| **ADR 003**                | FULLY_FUNCTIONAL | D2/graph module extraction decision                        |
| **ADR 004**                | FULLY_FUNCTIONAL | Footer row design decision                                 |
| **ADR 005**                | FULLY_FUNCTIONAL | Code duplication thresholds decision                       |
| **ADR 006**                | FULLY_FUNCTIONAL | Pre-v1 API stability guarantees                            |
| **ADR 007**                | FULLY_FUNCTIONAL | nom composition via root types                             |
| **ADR 008**                | FULLY_FUNCTIONAL | Dedup workflow decision (art-dupl threshold + checklist)   |
| **RELEASE.md**             | FULLY_FUNCTIONAL | Release process for 20-module mono-version workspace       |
| **ROADMAP.md**             | FULLY_FUNCTIONAL | Long-term direction and raw ideas                          |
| **DOMAIN_LANGUAGE.md**     | FULLY_FUNCTIONAL | Domain vocabulary                                          |
| **FORMAT_ARCHITECTURE.md** | FULLY_FUNCTIONAL | Format architecture documentation                          |

---

## Utility

| Feature                    | Status           | Notes                                                                                                      |
| -------------------------- | ---------------- | ---------------------------------------------------------------------------------------------------------- |
| **FilledStrings()**        | REMOVED          | Replaced by `slices.Repeat` (stdlib). Zero external callers.                                               |
| **InvalidFormatError**     | FULLY_FUNCTIONAL | Descriptive error with allowed values list                                                                 |
| **UnsupportedFormatError** | FULLY_FUNCTIONAL | Returned by `RenderTableData()` for table/json formats                                                     |
| **D2 utility methods**     | FULLY_FUNCTIONAL | `AddNodeSimple()`, `AddNodeWithShape()`, `AddEdgeSimple()`, `AddLabeledEdge()` — builder pattern shortcuts |

---

**Last audited:** 2026-06-18
**Total features:** 172
**Fully functional:** 161
**Partially functional:** 0
**Removed:** 9 (FormatCategory, OutputFormat, SortBy, FilledStrings, Register, Create, Unregister, RegisteredFormats, IsRegistered)
**Known issues:** 0
