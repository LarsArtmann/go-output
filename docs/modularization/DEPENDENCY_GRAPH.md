# Dependency Graph — go-output

**Date:** 2026-05-16

## Current Dependency Graph

```
                    ┌─────────────────────────────────────────────┐
                    │            External Dependencies            │
                    │  go-faster/yaml · x/term · go-branded-id   │
                    └─────────────────┬───────────────────────────┘
                                      │
                    ┌─────────────────▼───────────────────────────┐
                    │         root (package output)               │
                    │         3,587 production LOC                │
                    │                                             │
                    │  Core: Format, Renderer, TableData,         │
                    │         TreeNode, GraphNode, GraphEdge      │
                    │  Formatters: JSON, CSV, TSV, Markdown,      │
                    │              HTML, YAML, XML                 │
                    │  Graph: DOT, Mermaid, GraphRendererMixin     │
                    │  D2: D2Diagram, 5 files, rich types          │
                    │  Tree: ASCIITreeRenderer                     │
                    │  Streaming: StreamingHTMLRenderer            │
                    │  Test: output_test_helpers.go                │
                    │  Internal: gentest/, testutils/              │
                    └──┬──────┬──────┬──────────────────────────┘
                       │      │      │
            ┌──────────▼┐  ┌──▼──┐  ┌──▼──────┐
            │   enum/    │  │sort/│  │ escape/  │
            │   64 LOC   │  │dep. │  │ 76 LOC   │
            │   zero deps│  │→root│  │ zero deps│
            └────────────┘  └──────┘  └─────────┘


    ┌───────────────┐     ┌──────────────┐     ┌───────────────┐
    │   table/       │     │ integration/ │     │   examples/   │
    │   92 LOC       │     │ tests only   │     │  examples     │
    │ →root, lipgloss│     │→root,sort,   │     │ →root, table  │
    │               │     │  table       │     │ (go.mod buggy)│
    └───────────────┘     └──────────────┘     └───────────────┘

    ┌───────────────┐
    │  cmdguard/     │
    │  53 LOC        │
    │  zero deps     │
    │ (tests: root)  │
    └───────────────┘
```

### Cross-Module internal/ Leaks (Minor)

```
enum/enum_test.go ───► internal/gentest (in root module)
cmdguard/cmdguard_test.go ───► internal/gentest (in root module)
table/table_test.go ───► internal/testutils (in root module)
integration/*.go ───► internal/testutils (in root module)
```

These `internal/` packages are inside the root module but consumed by sibling modules' tests. Go's `internal/` visibility rules allow this (same repo root). **Decision: leave as-is** after self-review — extraction adds more complexity (2 new go.mod files, more replace directives) than the minor coupling warrants.

## Proposed Dependency Graph (After Extraction)

```
Level 0 — Leaf modules (zero internal deps)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌────────┐  ┌─────────┐
    │  enum/  │  │ escape/  │
    │ 64 LOC  │  │ 76 LOC  │
    └────────┘  └─────────┘

Level 2 — Core module
━━━━━━━━━━━━━━━━━━━━━
    ┌─────────────────────────────────────────────┐
    │        root (package output)                │
    │        ~2,183 production LOC                │
    │                                             │
    │  External: go-faster/yaml, x/term,          │
    │            go-branded-id                     │
    │  Internal: → enum, escape                   │
    │                                             │
    │  Core types: Format, Renderer, TableData,   │
    │    TreeNode, GraphNode, GraphEdge,           │
    │    GraphRenderer, BrandedID                  │
    │  Formatters: JSON, CSV, TSV, Markdown,      │
    │    HTML, YAML, XML, Tree, Streaming          │
    │  Internal packages: gentest/, testutils/     │
    │  Test helpers: output_test_helpers.go        │
    └──────────────────┬──────────────────────────┘
                       │
Level 3 — Format modules (depend on root core types)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    ┌──────────┐  ┌──────────┐  ┌───────────┐
    │   d2/     │  │  graph/   │  │  table/    │
    │  833 LOC  │  │  568 LOC  │  │  92 LOC   │
    │→root,enum, │  │→root,enum, │  │→root,     │
    │  escape   │  │  escape   │  │  lipgloss  │
    └──────────┘  └──────────┘  └───────────┘

Level 5 — Consumers
━━━━━━━━━━━━━━━━━━━
    ┌───────────────┐  ┌──────────────┐  ┌──────────┐
    │ integration/   │  │  examples/   │  │ cmdguard/ │
    │→root,sort,     │  │ →root, table │  │ 53 LOC    │
    │  table,d2,graph│  │  d2, graph   │  │ zero deps │
    └───────────────┘  └──────────────┘  └──────────┘

    ┌──────────┐
    │  sort/    │
    │deprecated│
    │ →root    │
    └──────────┘
```

## Module Dependency Matrix (Proposed)

| ↓ depends on →  | enum | escape | root | d2  | graph | table | sort | lipgloss |
| --------------- | ---- | ------ | ---- | --- | ----- | ----- | ---- | -------- |
| **enum**        | —    | —      | —    | —   | —     | —     | —    | —        |
| **escape**      | —    | —      | —    | —   | —     | —     | —    | —        |
| **root**        | ✅   | ✅     | —    | —   | —     | —     | —    | —        |
| **d2**          | ✅   | ✅     | ✅   | —   | —     | —     | —    | —        |
| **graph**       | ✅   | ✅     | ✅   | —   | —     | —     | —    | —        |
| **table**       | —    | —      | ✅   | —   | —     | —     | —    | ✅       |
| **sort**        | —    | —      | ✅   | —   | —     | —     | —    | —        |
| **cmdguard**    | —    | —      | —    | —   | —     | —     | —    | —        |
| **integration** | —    | —      | ✅   | ✅  | ✅    | ✅    | ✅   | —        |
| **examples**    | —    | —      | ✅   | —   | —     | ✅    | —    | —        |

**Cycles:** None. All dependencies point downward (higher row → lower column).

## DAG Proof

Each module's dependencies form a strict partial order:

```
enum < root < d2
enum < root < graph
enum < root < table < integration
root < sort < integration
root < integration
root < examples
```

No module appears on both sides of `<` in any chain. Therefore the graph is a DAG.
