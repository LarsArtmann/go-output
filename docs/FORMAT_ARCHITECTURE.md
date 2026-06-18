# go-output Format Architecture

## Overview

This document describes the extensible format architecture for go-output: **16 output formats** across **3 data shapes**, dispatched through a registry-based core with zero sub-module imports in the root package.

The 16 formats are: `table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml`.

> **Scope note.** This document covers the _data-rendering_ axis only. The `nom/` and `tui/` modules are a **separate, orthogonal subsystem** — real-time progress visualization — and are described in [Relationship to nom/ and tui/](#relationship-to-nom-and-tui). They are not formats, do not appear in the shape matrix, share no types with root's `tree` renderer, and do not depend on the root package.

## Data Shapes

Formats are classified by the data shapes they support. Each format may support multiple shapes.

### Shape Capability Matrix

The single source of truth is the `init()` block in [`shape.go`](../shape.go) (root defaults); sub-modules may override via their own `init()` calling `RegisterFormatShapes`. The matrix below reflects the registered defaults.

| Format   | ShapeTable | ShapeTree | ShapeGraph |
| -------- | :--------: | :-------: | :--------: |
| table    |     Y      |           |            |
| json     |     Y      |     Y     |     Y      |
| csv      |     Y      |           |            |
| tsv      |     Y      |           |            |
| xml      |     Y      |           |            |
| markdown |     Y      |           |            |
| yaml     |     Y      |     Y     |     Y      |
| html     |     Y      |     Y     |            |
| tree     |            |     Y     |            |
| d2       |     Y      |     Y     |     Y      |
| mermaid  |     Y      |     Y     |     Y      |
| dot      |     Y      |     Y     |     Y      |
| jsonl    |     Y      |           |            |
| asciidoc |     Y      |           |            |
| toml     |     Y      |     Y     |     Y      |
| plantuml |     Y      |     Y     |     Y      |

### Querying Capabilities

```go
// Check if a format supports a specific shape
output.FormatJSON.Supports(output.ShapeTable) // true

// Get all shapes a format supports
output.FormatD2.Shapes() // [ShapeTable, ShapeTree, ShapeGraph]

// Get all formats that support a shape (iterates AllFormats order)
output.FormatsForShape(output.ShapeGraph) // [json, d2, yaml, mermaid, dot, toml, plantuml]
```

### Deprecated Methods

The following are deprecated and redirect to the Shape API:

- `f.IsTableFormat()` → `f.Supports(ShapeTable)`
- `f.IsTreeFormat()` → `f.Supports(ShapeTree)`
- `f.IsGraphFormat()` → `f.Supports(ShapeGraph)`
- `f.Category()` → `f.Shapes()`

## Data Structures

### TableData

Unified data structure for all tabular outputs (defined in [`tabledata.go`](../tabledata.go)):

```go
type TableData struct {
    Headers []string
    Rows    [][]string
    Footer  []string // optional totals/summary row (ADR 004)
}
```

The optional `Footer` is rendered by **tabular** formats (CSV, TSV, Markdown, HTML, XML, AsciiDoc, Table) and skipped by **data** formats (JSON, YAML, TOML, JSONL) and **non-tabular** formats (Tree, Graph). HTML uses `<tfoot>`, XML uses `<footer>`, Markdown adds a second separator + bold footer inheriting column alignment. `TableData.Validate()` checks the footer column count matches the headers.

Helpers: `NewTableData`, `AddRow`, `AddRowChecked`, `SetFooter`, `HasFooter`, `Validate`, `ToMapSlice`, `CreateRowEdges`.

### TableDataProvider

Interface for types that provide tabular data. **Lives in the `table/` sub-module** ([`table/table.go`](../table/table.go)); the root `TableData` satisfies it via its `GetHeaders()` / `GetRows()` methods:

```go
// table.TableDataProvider
type TableDataProvider interface {
    GetHeaders() []string
    GetRows() [][]string
}
```

`table.FromTableData(data TableDataProvider, opts ...Option)` consumes this. `TableData` additionally satisfies the optional `table.FooterProvider` interface via `GetFooter()`.

### TreeNode

Hierarchical data structure for tree outputs (defined in [`tree.go`](../tree.go)):

```go
type TreeNode struct {
    ID       TreeNodeID
    Label    TreeNodeLabel
    Children []*TreeNode
    Metadata map[string]string
    parent   *TreeNode // unexported; use Parent()
    depth    int       // unexported; use Depth() (root = 0)
}
```

`NewTreeNode`, `AddChild` (maintains cached depth across subtrees), `Depth()`, and `Parent()` are the construction/access API. IDs and Labels are branded phantom types (see [`ids.go`](../ids.go)).

### GraphNode and GraphEdge

Data structures for graph/diagram outputs (defined in [`graph.go`](../graph.go)):

```go
type GraphNode struct {
    ID       GraphNodeID
    Label    GraphNodeLabel
    Shape    NodeShape
    Style    GraphStyle
    Metadata map[string]string
}

type GraphEdge struct {
    From  GraphNodeID
    To    GraphNodeID
    Label GraphNodeLabel
    Style EdgeStyle
}
```

## Interfaces

All defined in root ([`renderer.go`](../renderer.go), [`tree.go`](../tree.go), [`graph.go`](../graph.go), [`streaming.go`](../streaming.go)).

### Renderer

Base interface for all renderers:

```go
type Renderer interface {
    Render() (string, error)
}
```

### TableRenderer

For flat tabular data with SetHeaders/AddRow pattern:

```go
type TableRenderer interface {
    Renderer
    SetHeaders(headers []string)
    AddRow(row []string)
}
```

Note: `MarkdownTable` uses variadic `AddRow(row ...string) *MarkdownTable` and returns self for chaining, so it does not implement `TableRenderer`. Builders (`MarkdownTable`, `table.Table`) expose `AsTableRenderer()` adapters that wrap their fluent APIs behind the void-returning interface.

### TreeOutputRenderer

For hierarchical tree data:

```go
type TreeOutputRenderer interface {
    Renderer
    SetRoot(node *TreeNode)
}
```

### GraphRenderer

For diagram/graph data:

```go
type GraphRenderer interface {
    Renderer
    SetNodes(nodes []GraphNode)
    SetEdges(edges []GraphEdge)
}
```

### StreamingRenderer

For renderers that support streaming output:

```go
type StreamingRenderer interface {
    Renderer
    Stream(w io.Writer) error
}
```

`StreamingRendererFromRenderer(r)` adapts any `Renderer`, but **collects output before writing** — it does not provide true streaming. Only `markup.StreamingHTMLRenderer` provides genuine streaming behavior.

## Dispatch and Registry Architecture

The root package exposes **function-based dispatch** (not a `Create` factory). Two entry points live in [`render_tabledata.go`](../render_tabledata.go):

```go
func RenderTableData(data *TableData, format Format, opts RenderOptions) error
func RenderAnyData(data any, format Format, opts RenderOptions) error
```

Both validate, pick a writer (default `os.Stdout`), look up a registered marshaler, and return `*UnsupportedFormatError` if none is registered for the format. Sub-modules register their marshalers in `init()`, so a format is only available if the user imports the corresponding sub-module.

### RenderOptions

```go
type RenderOptions struct {
    Title     string    // HTML document title / Markdown header
    GraphID   string    // reserved for DOT via dispatch (use DOTRenderer.SetGraphID directly today)
    Writer    io.Writer // defaults to os.Stdout when nil
    ColorMode ColorMode // terminal color control (see below)
}
```

### Three registries, one generic type

All runtime dispatch is backed by a single generic, thread-safe container — `formatRegistry[T]` in [`registry.go`](../registry.go) — replacing what was previously three separate mutex+map boilerplates:

| Registry             | Value type           | Populated by                                                                        | Queried by               |
| -------------------- | -------------------- | ----------------------------------------------------------------------------------- | ------------------------ |
| `formatCapabilities` | `[]Shape`            | root `init()` + sub-module overrides                                                | `Format.Supports/Shapes` |
| `tableDataRegistry`  | `TableDataRenderer` | delimited, serialization, markup, d2, graph, plantuml, table, root (markdown, tree) | `RenderTableData`        |
| `anyDataRegistry`    | `AnyDataRenderer`   | serialization (JSON, YAML, TOML)                                                    | `RenderAnyData`          |

`RegisterTableDataRenderer`, `RegisterAnyDataRenderer`, `RegisterFormatShapes`, and the read-only `RegisteredTableDataFormats` / `RegisteredAnyDataFormats` helpers form the registration API.

### ColorMode

`ColorMode` (auto / always / never) controls ANSI color output and is wired through every terminal renderer: `table.WithColorMode(...)`, `ASCIITreeRenderer.SetColorMode(...)`, `MarkdownTable.SetColorMode(...)`, and `RenderOptions.ColorMode` for `RenderTableData` dispatch. `ColorModeAuto` (default) detects the terminal via `golang.org/x/term` and honors `NO_COLOR`, `CI`, `FORCE_COLOR`.

## Composition Shared Across Formats

Two exported composites in root eliminate per-renderer boilerplate:

- **`GraphRendererState`** ([`graph.go`](../graph.go)) — shared node/edge state with `SetNodes`, `SetEdges`, `AddNode`, `AddEdge`, `DedupEdges`, `SetNodesFromTableData`, `AddRowEdges`. Embedded by DOT, Mermaid, PlantUML, and the serialization graph renderers (JSON/YAML/TOML graph views). Tree-to-graph conversion flows through the `NodeEdgeAppender` interface (`AddTreeNodes`). _(Renamed from `GraphRendererMixin` on 2026-06-08.)_
- **`TableDataStore`** ([`tabledata.go`](../tabledata.go)) — shared table storage with `SetHeaders`, `AddRow`, `SetData`, `Data`, `SetFooter`, `HasFooter`. Embedded by delimited (CSV/TSV), serialization (JSON/YAML/TOML/JSONL), and markup (HTML/AsciiDoc/Streaming) renderers. _(Renamed/exported from `tableDataBase` on 2026-06-08.)_

## Format-Specific Notes

### JSON and YAML

JSON and YAML declare support for all three shapes (Table, Tree, Graph). They work as generic serialization functions — pass any data structure. For typed table rendering, use `serialization.MarshalJSONFromTableData` / `serialization.MarshalYAMLFromTableData`. TOML mirrors this (Table/Tree/Graph).

### D2

D2 supports all three shapes: table data via `d2.D2FromTableData`, and graph data via `d2.D2FromTree` or `GraphNode`/`GraphEdge`. D2 has richer domain-specific types than generic graph (shapes, arrows, SQL tables, classes, user journeys) and lives in the `d2/` module with its own `D2Node`/`D2Edge` types. It re-exports `D2NodeID`/`D2NodeLabel` from root.

### HTML

HTML supports table data (`markup.HTMLRenderer`) and tree data (`markup.HTMLTreeRenderer`). `markup.StreamingHTMLRenderer` provides true streaming for large datasets. Uses `html/template` for auto-escaping and `<tfoot>`/`footer-cell` CSS for the footer row.

### XMLWriter

`markup.XMLWriter` requires an `io.Writer` in its constructor. For string output, use with a `strings.Builder`:

```go
var buf strings.Builder
w := markup.NewXMLWriter(&buf)
_ = w.WriteHeader([]string{"Name", "Value"})
_ = w.WriteRow([]string{"Alice", "30"})
_ = w.WriteFooter()
fmt.Print(buf.String())
```

### Diagram renderers (DOT, Mermaid, PlantUML)

All three embed `GraphRendererState`, honor `GraphStyle` for per-node colors, and support Tree→graph conversion via `AddTreeNodes`. `DOTRenderer` exposes typed `RankDir`/`SplineStyle` enums plus `SetNodeSep`/`SetRankSep`/`SetGraphID`. `MermaidRenderer.SetCodeFence(bool)` toggles the markdown code fence (default on). IDs are sanitized with `escape.SlugifyID`.

## Relationship to nom/ and tui/

`nom/` and `tui/` are **not output formats**. They are a separate, orthogonal subsystem for **real-time progress visualization** driven by an external workflow engine that emits events. They do not register in any of the three registries above, are not reachable via `RenderTableData`/`RenderAnyData`, and do not appear in the shape matrix. Root has zero imports from either module (verified: `nom/go.mod` requires only `lipgloss/v2` + stdlib-adjacent helpers — no `output` dependency).

### "Isn't nom/ just a fancy real-time tree?"

The **rendered output** is a tree, yes. But "just a fancy tree" undersells what makes it NOM-style rather than a live-re-rendering of `output.TreeRenderer`. nom/ and root's `tree` renderer share the _word_ "tree" but **share no types and no code path**:

| Concern             | Root `tree` format                  | `nom/`                                                                                                            |
| ------------------- | ----------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| Data source         | Static `*output.TreeNode` graph     | Live `nom.Event` stream (`EventSubscriber.OnEvent`)                                                               |
| Node type           | `output.TreeNode`                   | `nom.ActivityNode` (separate type, no embedding, no conversion path)                                              |
| Ordering            | Insertion order                     | **Self-resorting**: Failed > Running > Paused > Pending > Completed; ties by elapsed-then-ID (`tree_priority.go`) |
| Time-awareness      | None                                | Persistent cross-run CSV at `~/.cache/nom-timing.csv`; **median** of last ≤10 runs predicts pending durations     |
| State model         | Single (`TreeNode`)                 | **Two**, synced: flat `ActivityDisplayState` + hierarchical `ActivityNode` (`SyncActivityTimingToTree`)           |
| Lifecycle           | One-shot `Render() (string, error)` | `InlineRenderer` with Start/Stop/Refresh, ANSI cursor-up redraw, 1s max-frame timer                               |
| Package dep on root | —                                   | **None** (cannot import `output` — would add transitive deps to nom's closure)                                    |

Strip the timing cache and the priority resorting and yes, you'd be left with a fancy real-time tree. Those two features — **time prediction** and **adaptive ordering** — are what make it NOM.

### nom/ architecture

Event-driven; no Bubble Tea dependency. Three layers:

- **Ingestion** — `nom.Event` / `nom.EventSubscriber` (`OnEvent`) is the integration contract. Event types are the `nom.Event*` constants (`EventWorkflowStarted`, `EventActivityStarted`, `EventActivityCompleted`, `EventActivityFailed`, …). `NOMStyleSubscriber` reads payloads via type-assertion accessor interfaces (`WorkflowEventAccessor`, `ActivityEventAccessor`, `DurationAccessor`, `ErrorAccessor`) so it never couples to a concrete event struct.
- **State** — `NOMStyleSubscriber` owns the activity map, the `DependencyTree`, and the `TimingCache` under one `sync.RWMutex`. `ActivityDisplayState` is the flat record of record; `ActivityNode` is the derived tree view; both carry the shared `DisplayState` (status, symbol, color, timing).
- **Rendering** — `DependencyTree` renders the priority-sorted hierarchy; `InlineRenderer` does the in-place ANSI redraw **without** alt-screen takeover (`Refresh()` for on-demand updates, terminal-width-aware truncation, cursor hide/show lifecycle).

### tui/ — Bubble Tea interactive TUI (depends on nom/)

Wraps `nom/` with a full-screen Bubble Tea program.

- `tui.BubbleTeaProgressReporter` — the primary entry point. Implements a `ProgressReporter` contract (idle → running → completed/errored state machine via `WorkflowState`). Owns a `ProgressModel` and lazily starts the TUI on first report (double-checked locking).
- `tui.ProgressModel` — the Bubble Tea model; holds its own `workflowState` on the TUI goroutine. All model mutations flow through a `send()` dispatcher to eliminate cross-goroutine races.
- `tui.DisplayMode` — `DisplayModeNOM` (renders the dependency tree) vs `DisplayModeUniversal` (step-based progress, `nh darwin switch` style).
- `Subscriber()` exposes the underlying `*nom.NOMStyleSubscriber` so callers can feed events.

### Where to look next

- Module dependency graph and per-module `go.mod` details: [`AGENTS.md`](../AGENTS.md#multi-module-workspace).
- Domain vocabulary (Activity, Workflow, DependencyTree, etc.): [`docs/DOMAIN_LANGUAGE.md`](DOMAIN_LANGUAGE.md).
- Design rationale for the shape matrix and the registry pattern: [`docs/adr/002-shape-capability-matrix.md`](adr/002-shape-capability-matrix.md).
- Footer row behavior: [`docs/adr/004-footer-row.md`](adr/004-footer-row.md).
