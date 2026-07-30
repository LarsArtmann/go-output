# Domain Language

A **Unified Language** for `go-output` — shared across Customer, Product Owner, Developer, and AI.
Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

Every term below should mean the **same thing** to everyone who reads it.

## Glossary

| Term             | Definition                                                                            | Context                                                                      |
| ---------------- | ------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| **Format**       | A string enum identifying an output format (e.g., "json", "csv", "d2")                | Used for CLI flags, runtime dispatch, capability queries                     |
| **Shape**        | A data shape an output format can render: table, tree, or graph                       | `ShapeTable`, `ShapeTree`, `ShapeGraph` — formats declare which they support |
| **Renderer**     | The core interface: `Render() (string, error)`                                        | Every format implements this — JSON, CSV, HTML, D2, DOT, Mermaid, etc.       |
| **Table**        | Tabular data structure with headers, rows, and optional footer                        | Central data type shared across all table-capable formats                    |
| **TreeNode**     | A node in a tree hierarchy with id, label, children, metadata                         | Used by Tree, JSON Tree, YAML Tree, HTML Tree renderers                      |
| **GraphNode**    | A node in a graph with branded ID, label, and optional shape                          | Used by DOT, Mermaid, JSON Graph, YAML Graph renderers                       |
| **GraphEdge**    | A directed edge between two GraphNodes with optional label and style                  | Used alongside GraphNode in graph renderers                                  |
| **Branded ID**   | A phantom-typed identifier (e.g., `D2NodeID`, `TreeNodeID`)                           | Prevents mixing different ID types at compile time                           |
| **GraphBuilder** | Write-side builder for graph data (nodes, edges). `Build()` returns immutable `Graph` | Embedded by DOT/Mermaid/PlantUML renderers for shared state                  |
| **ColorMode**    | Terminal color output mode: auto, always, never                                       | Respects `NO_COLOR`, CI env vars, TTY detection                              |
| **Registry**     | Format→marshaler dispatch map for `RenderTable`                                       | `RegisterTableMarshaler(format, fn)` — sub-modules register in `init()`      |
| **Sentinel Error** | An exported `var ErrFoo = errors.New(...)` that consumers match via `errors.Is`   | Root owns `ErrColumnMismatch`, `ErrNilRow`; sub-modules own domain-specific ones |
| **Typed Error**  | An exported struct (`*UnsupportedFormatError`, `*ParseError`) with structured fields | Consumers extract via `errors.AsType[*T](err)` (Go 1.26+)                   |

## Entities

| Term                | Definition                                                                        | Context                                                                                         |
| ------------------- | --------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| **Diagram**         | A D2 diagram with nodes, edges, tables, classes, and styling                      | Lives in `d2/` module — rich domain model                                                       |
| **DOTRenderer**     | A renderer producing DOT/Graphviz output                                          | Lives in `graph/` module                                                                        |
| **MermaidRenderer** | A renderer producing Mermaid flowchart output                                     | Lives in `graph/` module                                                                        |
| **HTMLRenderer**    | A renderer producing styled HTML tables                                           | Supports both table and tree shapes                                                             |
| **Activity**        | Unified source of truth for a workflow activity's state, timing, and visual style | Has branded ID/Label for diagram export; consumed via immutable `ActivitySnapshot` by renderers |
| **DependencyTree**  | Hierarchical activity visualization with priority-based child sorting             | `nom/` module — renders failed/running activities first                                         |
| **ActivityReader**  | Read-only projection of activities as `GraphNode`/`GraphEdge` slices              | `subscriber.Store()` returns this; any `GraphRenderer` can consume it                           |
| **TimingCache**     | Persisted activity duration history (CSV at `~/.cache/nom-timing.csv`)            | Uses median of last <=10 entries for robust estimates                                           |
| **NOMSubscriber**   | Event-driven subscriber that manages activities and drives the dependency tree    | `nom/` module — implements `EventSubscriber` interface                                          |

## Value Objects

| Term          | Definition                                                 | Context                                   |
| ------------- | ---------------------------------------------------------- | ----------------------------------------- |
| **NodeStyle** | Visual style for a graph node (fill, stroke, font color)   | Root-level generic node styling           |
| **EdgeStyle** | Visual style for a graph edge (line style, labels)         | Used in DOT and other graph renderers     |
| **NodeShape** | Visual shape for a graph node: box, diamond, ellipse, etc. | DOT and Mermaid renderers interpret these |
| **Direction** | D2 layout direction: down, right, left, up                 | `d2/` module — default is down            |

## Bounded Contexts

| Context                              | Description                                                                                  |
| ------------------------------------ | -------------------------------------------------------------------------------------------- |
| **Core** (root `package output`)     | Core types, interfaces, enum + envdetect utilities, `RenderTable` dispatch                   |
| **Delimited** (`delimited/`)         | CSV + TSV writers and formatters                                                             |
| **Serialization** (`serialization/`) | JSON + YAML + TOML + JSONL marshaling and renderers                                          |
| **Markup** (`markup/`)               | XML + HTML + AsciiDoc + Streaming HTML renderers                                             |
| **D2** (`d2/` module)                | D2-specific diagram domain — rich types, SQL tables, grid layouts, style classes             |
| **Graph** (`graph/` module)          | DOT and Mermaid graph rendering — flowcharts, directed graphs                                |
| **Table** (`table/` module)          | Terminal table rendering with lipgloss styling — isolated heavy dependency                   |
| **Markdown** (`markdown/` module)    | Markdown table renderer — extracted from root, self-registers via `init()`                   |
| **Tree** (`tree/` module)            | ASCII tree renderer — extracted from root, self-registers via `init()`                       |
| **daghtml** (`daghtml/` module)      | Zero-dep interactive SVG DAG visualization for HTML — independently publishable              |
| **NOM** (`nom/` module)              | NOM-style real-time progress — DependencyTree, InlineRenderer, TimingCache, event subscriber |
| **TUI** (`tui/` module)              | Bubble Tea v2 interactive progress UI — depends on nom + bubbletea + lipgloss                |
| **BDD** (`bdd/` module)              | Behavior-driven test suite (Ginkgo/Gomega) — end-user-focused specs, test-only               |

---

> **How to use this file:**
>
> - Keep terms concise — one clear sentence per definition
> - Update when new domain concepts emerge
> - Use these terms consistently in code, docs, and conversations
> - When in doubt about a word's meaning, check here first
