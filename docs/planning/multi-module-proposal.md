# go-output: Multi-Module Structure Proposal

**Created:** 2026-05-07

---

## Proposed Module Split

```
go-output/
├── go.work
│
├── enum/                    # Module 1: Generic enum utilities
│   └── go.mod               # github.com/larsartmann/go-output/enum
│
├── escape/                  # Module 2: Format-specific escaping
│   └── go.mod               # github.com/larsartmann/go-output/escape
│
├── core/                    # Module 3: Core types + text/structured formatters
│   └── go.mod               # depends: enum, escape, go-faster/yaml, x/term
│   ├── format.go            # Format enum, Renderer, TableData, TreeNode, TableRenderer
│   ├── color.go             # ColorMode + terminal detection
│   ├── ids.go               # Branded IDs
│   ├── sort.go              # SortBy enum
│   ├── registry.go          # Opt-in renderer registry
│   ├── slices.go            # FilledStrings
│   ├── json.go, csv.go, tsv.go, yaml.go, xml.go
│   ├── markdown.go, html.go, tree.go, streaming.go
│   ├── delimited.go, markup.go, marshal.go
│   └── format_deprecated.go
│
├── d2/                      # Module 4: D2 diagram subsystem
│   └── go.mod               # depends: core, enum, escape
│   ├── d2.go                # D2 domain types (D2Node, D2Edge, D2Table)
│   ├── d2_enum.go           # D2-specific enums (Direction, Shape, ArrowType…)
│   ├── d2_render.go         # D2Diagram builder + Render()
│   ├── d2_write.go          # Style/edge writing helpers
│   └── d2_convert.go        # TableData/Tree → D2 (depends on core types)
│
├── graph/                   # Module 5: Generic graph + DOT/Mermaid renderers
│   └── go.mod               # depends: core, enum, escape
│   ├── graph.go             # GraphNode, GraphEdge, GraphRenderer interface, AddTreeNodes
│   ├── dot.go               # DOT renderer + GraphRendererMixin
│   └── mermaid.go           # Mermaid renderer
│
├── sort/                    # Module 6: Generic sorting
│   └── go.mod               # depends: core (SortBy type)
│   ├── sorter.go
│   └── compare.go
│
├── table/                   # Module 7: Lipgloss terminal tables
│   └── go.mod               # depends: core (Renderer/TableData), lipgloss
│   └── table.go
│
├── cmdguard/                # Module 8: CLI flag integration
│   └── go.mod               # ZERO deps (fully isolated)
│   └── flag.go
│
└── examples/                # Module 9: Usage examples
    └── go.mod               # depends: all modules above
```

## Dependency Graph

```
     enum/    escape/    cmdguard/
       ↑         ↑          (isolated)
       │         │
    ┌──┴─────────┴──┐
    │     core/      │   ← Format, Renderer, TableData, TreeNode, ColorMode, SortBy
    │  (json, csv,   │      + json, csv, tsv, yaml, xml, markdown, html, tree, streaming
    │   yaml, xml,   │
    │   markdown,    │
    │   html, tree)  │
    └──┬─────┬──────┬┘
       │     │      │
    ┌──┴──┐ ┌┴───┐ ┌┴──────┐
    │ d2/ │ │graph│ │ sort/ │  ← table/ also depends on core + lipgloss
    └─────┘ │(dot,│ └───────┘
            │merm)│
            └─────┘
```

## Rationale

| Module | Why independent? |
|---|---|
| **enum** | Leaf package, zero deps, widely reusable beyond this project |
| **escape** | Leaf package, zero deps, pure string transformations |
| **core** | Defines the shared vocabulary (Renderer, TableData, TreeNode, Format) — everything depends on this |
| **d2** | Rich self-contained domain (5 files, its own enums, types, render logic). Only cross-module dependency is `core` for TableData/TreeNode conversion |
| **graph** | DOT + Mermaid share `GraphRendererMixin`. Naturally clusters with `graph.go` generic types. Isolates the graph concern from text formatters |
| **sort** | Generic utility, only needs `SortBy` from core. Users who just need sorting don't pull in formatters |
| **table** | Only module pulling in heavy **lipgloss** dependency. Users who don't need terminal tables skip it entirely |
| **cmdguard** | Zero deps on everything. Utility for CLI apps that may not even use the output library |

## Key Benefit: Dependency Isolation

```
User wants just JSON output?    → core only (no lipgloss, no D2, no graph deps)
User wants D2 diagrams?         → core + d2
User wants terminal tables?     → core + table (lipgloss pulled in)
User wants CLI flag parsing?    → cmdguard only (zero transitive deps)
User wants enum utilities?      → enum only
```

## `go.work`

```go
go 1.26.2

use (
    ./enum
    ./escape
    ./core
    ./d2
    ./graph
    ./sort
    ./table
    ./cmdguard
    ./examples
)
```

## Caveat

The biggest challenge is `d2_convert.go` — it converts `core.TableData`/`core.TreeNode` → D2 types, creating a dependency from `d2/` → `core/`. This is unidirectional and clean. Similarly, `graph.go`'s `AddTreeNodes` converts `TreeNode` → `GraphNode`, so `graph/` → `core/` is also clean. No circular deps in this layout.

## Open Questions

- Should `core` stay as a flat package or be split further (e.g., separate text formatters from core types)?
- Should `d2_convert.go` live in `d2/` (depends on core types) or in `core/` (would reverse the dependency)?
- Where do integration tests live — a top-level `integration/` module or distributed per-module?
- Should `internal/testutils` and `internal/gentest` become shared test helpers across modules?
