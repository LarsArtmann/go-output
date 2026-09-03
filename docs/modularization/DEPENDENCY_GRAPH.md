# Dependency Graph — go-output

**Date:** 2026-05-25

## Current Dependency Graph

```
Level 0 — Leaf modules (zero internal deps)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌──────────┐  ┌──────────┐  ┌──────────────┐
    │  enum/    │  │ escape/   │  │ testhelpers/ │
    │ 64 LOC    │  │ 76 LOC    │  │ shared tests │
    │ zero deps │  │ zero deps │  │  zero deps   │
    └──────────┘  └──────────┘  └──────────────┘

Level 2 — Core module
━━━━━━━━━━━━━━━━━━━━━
    ┌─────────────────────────────────────────────┐
    │        root (package output)                │
    │        ~5,200 production LOC                │
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
    │  Internal packages: gentest/                 │
    │  Test helpers: output_test_helpers.go        │
    └──────────────────┬──────────────────────────┘
                       │

Level 2 — Format modules (depend on root core types)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌──────────────┐  ┌──────────────┐  ┌────────────┐  ┌──────────────┐
    │    d2/        │  │   graph/     │  │   table/    │  │  delimited/  │
    │  ~850 LOC     │  │  ~320 LOC    │  │  93 LOC     │  │  CSV + TSV   │
    │ →root,escape, │  │ →root,escape,│  │ →root,      │  │ →root        │
    │   testhelpers │  │  testhelpers │  │   lipgloss  │  │              │
    └──────────────┘  └──────────────┘  └────────────┘  └──────────────┘
    ┌──────────────────┐  ┌──────────────────┐
    │  serialization/  │  │     markup/      │
    │  JSON + YAML     │  │  XML + HTML +    │
    │ →root, go-faster │  │  Streaming HTML  │
    │   yaml           │  │ →root, escape    │
    └──────────────────┘  └──────────────────┘

Level 5 — Consumers
━━━━━━━━━━━━━━━━━━━
    ┌───────────────┐  ┌────────────────┐
    │ integration/   │  │   examples/    │
    │ →root, table,  │  │ →root, table,  │
    │  d2, graph     │  │  d2, graph     │
    └───────────────┘  └────────────────┘
```

## Module Dependency Matrix

| ↓ depends on →    | enum | escape | root | d2 | graph | table | delimited | serialization | markup |
| ----------------- | ---- | ------ | ---- | -- | ----- | ----- | --------- | ------------- | ------ |
| **enum**          | —    | —      | —    | —  | —     | —     | —         | —             | —      |
| **escape**        | —    | —      | —    | —  | —     | —     | —         | —             | —      |
| **testhelpers**   | —    | —      | —    | —  | —     | —     | —         | —             | —      |
| **root**          | ✅   | —      | —    | —  | —     | —     | —         | —             | —      |
| **d2**            | —    | ✅     | ✅   | —  | —     | —     | —         | —             | —      |
| **graph**         | —    | ✅     | ✅   | —  | —     | —     | —         | —             | —      |
| **table**         | —    | —      | ✅   | —  | —     | —     | —         | —             | —      |
| **delimited**     | —    | —      | ✅   | —  | —     | —     | —         | —             | —      |
| **serialization** | —    | —      | ✅   | —  | —     | —     | —         | —             | —      |
| **markup**        | —    | ✅     | ✅   | —  | —     | —     | —         | —             | —      |
| **integration**   | —    | —      | ✅   | ✅ | ✅    | ✅    | ✅        | ✅            | ✅     |
| **examples**      | —    | —      | ✅   | ✅ | ✅    | ✅    | ✅        | ✅            | ✅     |

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
2. **`go get github.com/larsartmann/go-output`** pulls only root + enum + x/term + branded-id — zero lipgloss, zero yaml, zero escape
3. **Each format module is independently versionable** — d2, graph, table can evolve at their own pace
4. **sort/ is deleted** — was deprecated, now fully removed. Use `slices.SortStableFunc` + `cmp.Compare` from stdlib
5. **registry.go is deleted** — renderer factory registry removed, replaced by `TableDataMarshaler` dispatch
6. **testhelpers/ is shared** — zero deps, used by d2, graph, and root for test assertions
