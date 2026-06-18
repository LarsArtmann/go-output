# Domain Language

A **Unified Language** for `go-output` — shared across Customer, Product Owner, Developer, and AI.
Inspired by Domain-Driven Design (DDD) Ubiquitous Language.

Every term below should mean the **same thing** to everyone who reads it.

## Glossary

| Term                   | Definition                                                             | Context                                                                      |
| ---------------------- | ---------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| **Format**             | A string enum identifying an output format (e.g., "json", "csv", "d2") | Used for CLI flags, runtime dispatch, capability queries                     |
| **Shape**              | A data shape an output format can render: table, tree, or graph        | `ShapeTable`, `ShapeTree`, `ShapeGraph` — formats declare which they support |
| **Renderer**           | The core interface: `Render() (string, error)`                         | Every format implements this — JSON, CSV, HTML, D2, DOT, Mermaid, etc.       |
| **TableData**          | Tabular data structure with headers and rows (`[][]string`)            | Central data type shared across all table-capable formats                    |
| **TreeNode**           | A node in a tree hierarchy with id, label, children, metadata          | Used by Tree, JSON Tree, YAML Tree, HTML Tree renderers                      |
| **GraphNode**          | A node in a graph with branded ID, label, and optional shape           | Used by DOT, Mermaid, JSON Graph, YAML Graph renderers                       |
| **GraphEdge**          | A directed edge between two GraphNodes with optional label and style   | Used alongside GraphNode in graph renderers                                  |
| **Branded ID**         | A phantom-typed identifier (e.g., `D2NodeID`, `TreeNodeID`)            | Prevents mixing different ID types at compile time                           |
| **GraphRendererState** | Shared state holder for graph renderers (nodes, edges)                 | Embedded by DOTRenderer and MermaidRenderer in `graph/` module               |
| **ColorMode**          | Terminal color output mode: auto, always, never                        | Respects `NO_COLOR`, CI env vars, TTY detection                              |
| **Registry**           | Format→marshaler dispatch map for `RenderTableData`                    | `RegisterTableDataRenderer(format, fn)` — sub-modules register in `init()`   |

## Entities

| Term                   | Definition                                                                        | Context                                                                              |
| ---------------------- | --------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| **D2Diagram**          | A D2 diagram with nodes, edges, tables, classes, and styling                      | Lives in `d2/` module — rich domain model                                            |
| **DOTRenderer**        | A renderer producing DOT/Graphviz output                                          | Lives in `graph/` module                                                             |
| **MermaidRenderer**    | A renderer producing Mermaid flowchart output                                     | Lives in `graph/` module                                                             |
| **HTMLRenderer**       | A renderer producing styled HTML tables                                           | Supports both table and tree shapes                                                  |
| **Activity**           | Unified source of truth for a workflow activity's state, timing, and visual style | Embeds `GraphNode` for diagram export; shared by pointer between subscriber and tree |
| **DependencyTree**     | Hierarchical activity visualization with priority-based child sorting             | `nom/` module — renders failed/running activities first                              |
| **ActivityReader**     | Read-only projection of activities as `GraphNode`/`GraphEdge` slices              | `subscriber.Store()` returns this; any `GraphRenderer` can consume it                |
| **TimingCache**        | Persisted activity duration history (CSV at `~/.cache/nom-timing.csv`)            | Uses median of last ≤10 entries for robust estimates                                 |
| **NOMStyleSubscriber** | Event-driven subscriber that manages activities and drives the dependency tree    | `nom/` module — implements `EventSubscriber` interface                               |

## Value Objects

| Term             | Definition                                                         | Context                                   |
| ---------------- | ------------------------------------------------------------------ | ----------------------------------------- |
| **NodeShape**    | Visual shape for a graph node: rectangle, diamond, ellipse, etc.   | DOT and Mermaid renderers interpret these |
| **EdgeStyle**    | Visual style for a graph edge: solid, dashed, dotted               | Used in DOT renderer                      |
| **D2NodeShape**  | D2-specific node shape: rectangle, circle, hexagon, cylinder, etc. | 20 shapes supported in `d2/` module       |
| **D2ArrowType**  | D2-specific arrow head style                                       | 11 arrow types in `d2/` module            |
| **D2Constraint** | SQL column constraint: primary, unique, foreign                    | Used in D2 SQL table rendering            |

## Bounded Contexts

| Context                              | Description                                                                                  |
| ------------------------------------ | -------------------------------------------------------------------------------------------- |
| **Core** (root `package output`)     | Core types, interfaces, formatters (Markdown, Tree) + ColorMode, RenderTableData dispatch    |
| **Delimited** (`delimited/`)         | CSV + TSV writers and formatters                                                             |
| **Serialization** (`serialization/`) | JSON + YAML marshaling and renderers (go-faster/yaml isolated here)                          |
| **Markup** (`markup/`)               | XML + HTML + Streaming HTML renderers (escape utilities isolated here)                       |
| **D2** (`d2/` module)                | D2-specific diagram domain — rich types, SQL tables, grid layouts, style classes             |
| **Graph** (`graph/` module)          | DOT and Mermaid graph rendering — flowcharts, directed graphs                                |
| **Table** (`table/` module)          | Terminal table rendering with lipgloss styling — isolated heavy dependency                   |
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
