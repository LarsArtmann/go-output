# Dependency Graph — go-output

**Date:** 2026-05-23

## Current Dependency Graph

```
Level 0 — Leaf modules (zero internal deps)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌──────────┐  ┌──────────┐  ┌──────────────┐  ┌──────────┐
    │  enum/    │  │ escape/   │  │ testhelpers/ │  │  sort/    │
    │ 64 LOC    │  │ 76 LOC    │  │ shared tests │  │deprecated│
    │ zero deps │  │ zero deps │  │  zero deps   │  │ zero deps│
    └──────────┘  └──────────┘  └──────────────┘  └──────────┘

Level 2 — Core module
━━━━━━━━━━━━━━━━━━━━━
    ┌─────────────────────────────────────────────┐
    │        root (package output)                │
    │        ~1,400 production LOC                │
    │                                             │
    │  External: go-faster/yaml, x/term,          │
    │            go-branded-id                     │
    │  Internal: → enum, escape, testhelpers      │
    │                                             │
    │  Core types: Format, Renderer, TableData,   │
    │    TreeNode, GraphNode, GraphEdge,           │
    │    GraphRenderer, GraphRendererMixin,        │
    │    BrandedID, SortBy, ColorMode              │
    │  Formatters: JSON, CSV, TSV, Markdown,      │
    │    HTML, YAML, XML, Tree, Streaming          │
    │  Internal packages: gentest/, testutils/     │
    │  Test helpers: output_test_helpers.go        │
    └──────────────────┬──────────────────────────┘
                       │

Level 3 — Format modules (depend on root core types)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌──────────────┐  ┌──────────────┐  ┌────────────┐
    │    d2/        │  │   graph/     │  │   table/    │
    │  833 LOC      │  │  ~350 LOC    │  │  92 LOC     │
    │ →root,escape, │  │ →root,escape,│  │ →root,      │
    │   testhelpers │  │  testhelpers │  │   lipgloss  │
    └──────────────┘  └──────────────┘  └────────────┘

Level 5 — Consumers
━━━━━━━━━━━━━━━━━━━
    ┌───────────────┐  ┌────────────────┐
    │ integration/   │  │   examples/    │
    │ →root, table,  │  │ →root, table,  │
    │  d2, graph     │  │  d2, graph     │
    └───────────────┘  └────────────────┘
```

## Module Dependency Matrix

| ↓ depends on →  | enum | escape | root | d2  | graph | table | lipgloss |
| --------------- | ---- | ------ | ---- | --- | ----- | ----- | -------- |
| **enum**        | —    | —      | —    | —   | —     | —     | —        |
| **escape**      | —    | —      | —    | —   | —     | —     | —        |
| **testhelpers** | —    | —      | —    | —   | —     | —     | —        |
| **root**        | ✅   | ✅     | —    | —   | —     | —     | —        |
| **d2**          | —    | ✅     | ✅   | —   | —     | —     | —        |
| **graph**       | —    | ✅     | ✅   | —   | —     | —     | —        |
| **table**       | —    | —      | ✅   | —   | —     | —     | ✅       |
| **sort**        | —    | —      | —    | —   | —     | —     | —        |
| **integration** | —    | —      | ✅   | ✅  | ✅    | ✅    | —        |
| **examples**    | —    | —      | ✅   | ✅  | ✅    | ✅    | —        |

**Cycles:** None. All dependencies point downward (higher row → lower column).

## DAG Proof

Each module's dependencies form a strict partial order:

```
enum < root < d2
enum < root < graph
enum < root < table < integration
root < d2 < integration
root < graph < integration
root < table < integration
root < d2 < examples
root < graph < examples
root < table < examples
```

No module appears on both sides of `<` in any chain. Therefore the graph is a DAG.

## Key Properties

1. **Root has ZERO imports from sub-modules** — verified by `go mod graph`
2. **`go get github.com/larsartmann/go-output`** pulls only root + enum + escape + yaml + x/term + branded-id — zero lipgloss, zero d2, zero graph
3. **Each format module is independently versionable** — d2, graph, table can evolve at their own pace
4. **sort/ is deprecated** — zero deps, only `ByField` helper remains
5. **testhelpers/ is shared** — zero deps, used by d2, graph, and root for test assertions
