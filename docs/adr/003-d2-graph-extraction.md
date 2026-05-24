# ADR 003: Extract D2 and Graph Renderers into Sub-Modules

**Date:** 2026-05-23
**Status:** ACCEPTED & IMPLEMENTED
**Deciders:** Lars Artmann
**Supersedes:** ADR 001 "Not Done Yet" section

## Context

After ADR 001 established the multi-module workspace with 8 modules, the root `package output` remained a god-package at 4,345 production LOC across 31 files. It contained all core types, 7 table formatters, tree rendering, 3 graph renderers (DOT, Mermaid, D2), a full D2 diagram system with rich domain types, streaming, and shared infrastructure.

Two natural module boundaries existed:

- **D2** (835 LOC, 5 production files) — rich domain model with D2Node, D2Edge, D2Direction, D2NodeShape, D2ArrowType, D2Constraint
- **Graph** (319 LOC, 2 production files) — DOT and Mermaid renderers sharing `GraphRendererMixin`

The key coupling issue: `render_tabledata.go` called D2 and graph constructors directly, creating a root→sub-module dependency cycle if extracted naively.

## Decision

Extract `d2/` and `graph/` as independent Go modules, breaking the `render_tabledata.go` cycle by returning `UnsupportedFormatError` for D2, Mermaid, and DOT formats.

### Specific changes:

1. **D2 module** — Move `d2.go`, `d2_enum.go`, `d2_render.go`, `d2_write.go`, `d2_convert.go` to `d2/` with `package d2`. Re-export `D2NodeID` and `D2NodeLabel` as type aliases for ergonomic use.

2. **Graph module** — Move `dot.go` and `mermaid.go` to `graph/` with `package graph`. Keep `GraphRendererMixin` in root (operates on root types).

3. **Accessor methods** — Add `Nodes()`, `Edges()`, `NodesPtr()`, `EdgesPtr()` to `GraphRendererMixin` so sub-modules can access unexported `nodes`/`edges` fields.

4. **render_tabledata.go decoupling** — D2/Mermaid/DOT cases return `UnsupportedFormatError`. Callers import sub-modules directly.

5. **No init() registration** — Sub-modules provide convenience constructors. Registry stays format-agnostic in root.

## Consequences

**Positive:**

- Root module reduced from 4,345 to ~2,800 production LOC
- `go get github.com/larsartmann/go-output` pulls ZERO lipgloss, d2, or graph deps
- d2/ and graph/ can evolve and be versioned independently
- DAG verified: root has zero imports from any sub-module

**Negative:**

- Breaking API change: `output.D2Node` → `d2.D2Node`, `output.DOTFromTableData` → `graph.DOTFromTableData`
- `render_tabledata.go` can no longer render D2/Mermaid/DOT (returns error)
- Test helpers duplicated in `graph/helpers_test.go` (sub-modules can't import `internal/gentest`)
- Sub-modules depend on root for core types (`GraphNode`, `GraphEdge`, `GraphRenderer`)

## Alternatives Considered

1. **Registry-based dispatch in render_tabledata.go** — Root registers factory functions, sub-modules register themselves via init(). Rejected: global state side effects, import cycles still possible.
2. **Keep everything in root** — Rejected: defeats purpose of ADR 001, root stays god-package.
3. **Move GraphRendererMixin to graph/** — Rejected: creates circular dependency (graph/ imports root for types, root imports graph/ for mixin).
