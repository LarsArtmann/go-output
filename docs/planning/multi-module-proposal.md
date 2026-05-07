# go-output: Multi-Module Migration Proposal

**Created:** 2026-05-07
**Status:** DRAFT — Refined Proposal v3
**Based on:** Deep analysis of current codebase + patterns from go-cqrs-lite + Go workspace best practices

---

## Table of Contents

- [Self-Critique of Previous Versions](#self-critique-of-previous-versions)
- [Key Design Decisions](#key-design-decisions)
- [Current State Analysis](#current-state-analysis)
- [Proposed Module Structure](#proposed-module-structure)
- [Dependency Graph](#dependency-graph)
- [Detailed Module Specs](#detailed-module-specs)
- [What Moves vs What Stays](#what-moves-vs-what-stays)
- [Key Coupling Points & Resolutions](#key-coupling-points--resolutions)
- [Pareto-Sorted Execution Plan](#pareto-sorted-execution-plan)
- [Open Questions](#open-questions)

---

## Self-Critique of Previous Versions

| Issue                                         | Version | Why it matters                                                                                                                                                                                                                                                    | Resolution                                                                      |
| --------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| `enum/` and `escape/` shown with no .go files | v1      | They already have `.go` files — proposal was misleading                                                                                                                                                                                                           | They just need `go.mod` added                                                   |
| Missed reverse dependencies                   | v1      | `sort/sorter.go` → root (`SortBy`), `table/table.go` → root (`Renderer`)                                                                                                                                                                                          | Must update import paths                                                        |
| Missed unexported type coupling               | v1      | `streaming.go` embeds unexported `tableDataBase` from `html.go`                                                                                                                                                                                                   | These files must stay together (in root)                                        |
| `GraphRendererMixin` defined in wrong file    | v1      | In `dot.go` but used by `mermaid.go`                                                                                                                                                                                                                              | Extract to `graph_mixin.go` in graph module                                     |
| No `replace` directives                       | v1      | Needed for standalone module development                                                                                                                                                                                                                          | Added, following go-cqrs-lite pattern                                           |
| **Created unnecessary `core/` directory**     | v2      | Moving 28 files and renaming `package output` → `package core` is massive churn for zero user benefit                                                                                                                                                             | **Root IS the core module.** Keep `package output`, keep all formatters in root |
| **Kept `sort/` package**                      | v2      | `Sorter[T]` is a thin wrapper over `sort.SliceStable` (stdlib has `slices.SortStableFunc` since Go 1.21). `ByField[T]` is literally `cmp.Less(extract(a), extract(b))` — two lines. The `SortBy` enum is stored but **never used in sort logic** — dead metadata. | **Deprecate `sort/` entirely.** Stdlib does the same job                        |

---

## Key Design Decisions

### 1. Root IS the core module — no `core/` directory

Moving 28 files into a `core/` subdirectory and renaming `package output` to `package core` is massive, risky churn:

- Every file changes package declaration
- Every internal import path changes
- Users must change `import "github.com/larsartmann/go-output"` to `import "github.com/larsartmann/go-output/core"`
- Tests, examples, docs all break

Instead, the root `go.mod` stays. It IS the library. Types, interfaces, and formatters live in `package output` at the root. Only truly independent subsystems get their own module.

This mirrors go-cqrs-lite's pattern: `core/` is the zero-dep foundation. Our root plays that same role.

### 2. Deprecate `sort/` — reinventing the wheel

The entire `sort/` package provides no value over Go 1.21+ stdlib:

```go
// sort/ package — unnecessary wrapper
sort.New(items, output.SortByName, true).
    WithLessFunc(sort.ByField(func(p Project) string { return p.Name })).
    Sort()

// stdlib equivalent — simpler, no extra dependency
slices.SortStableFunc(items, func(a, b Project) int {
    return cmp.Compare(b.Name, a.Name) // desc = swap b,a
})
```

**Action:** Deprecate `sort/` with a Go doc comment pointing to `slices.SortStableFunc` + `cmp.Compare`. Remove in next major version. Remove `SortBy` enum from root if it has no other consumers.

---

## Current State Analysis

### Existing Sub-Packages

| Package     | Files                     | Imports from root                   | External deps        |
| ----------- | ------------------------- | ----------------------------------- | -------------------- |
| `enum/`     | `enum.go`                 | **None**                            | None                 |
| `escape/`   | `escape.go`               | **None**                            | None                 |
| `sort/`     | `sorter.go`, `compare.go` | `output.SortBy`                     | None — **DEPRECATE** |
| `table/`    | `table.go`                | `output.Renderer` (interface check) | `lipgloss/v2`        |
| `cmdguard/` | `flag.go`                 | **None**                            | None                 |

### Third-Party Dependencies (only 3!)

| Dependency                  | Used by               | Heavy?                         |
| --------------------------- | --------------------- | ------------------------------ |
| `charm.land/lipgloss/v2`    | `table/table.go` only | **Yes** — many transitive deps |
| `github.com/go-faster/yaml` | `yaml.go` only        | Medium                         |
| `golang.org/x/term`         | `color.go` only       | Light                          |

### Unexported Coupling in Root (forces files to stay together)

| Unexported symbol                           | Used by                   | Must stay in same package               |
| ------------------------------------------- | ------------------------- | --------------------------------------- |
| `tableDataBase` (struct)                    | `html.go`, `streaming.go` | Yes — `StreamingHTMLRenderer` embeds it |
| `marshal()` / `unmarshal()`                 | `json.go`, `yaml.go`      | Yes                                     |
| `writeMarkupRow()` / `writeMarkupColumns()` | `xml.go`, `html.go`       | Yes                                     |

---

## Proposed Module Structure

```
go-output/
├── go.mod                           # ROOT MODULE — package output (stays as-is)
├── go.work                          # Workspace file
│
├── format.go                        # Format enum, Renderer, TableData, TreeNode
├── format_deprecated.go             # OutputFormat backward compat
├── color.go                         # ColorMode + terminal detection
├── ids.go                           # Branded ID phantom types
├── registry.go                      # Opt-in renderer registry
├── slices.go                        # FilledStrings
├── sort.go                          # SortBy enum (kept if still useful without sort/)
├── marshal.go, json.go, yaml.go     # JSON/YAML formatters
├── csv.go, tsv.go, delimited.go     # CSV/TSV formatters
├── xml.go, markdown.go, markup.go   # XML/Markdown formatters
├── html.go, streaming.go            # HTML formatters (coupled via tableDataBase)
├── tree.go                          # ASCII tree renderer
├── graph.go                         # GraphNode, GraphEdge, GraphRenderer, AddTreeNodes
├── ...root test files
├── internal/                        # Test helpers (stays in root)
│   ├── gentest/
│   └── testutils/
│
├── enum/                            # Module: Generic enum utilities
│   ├── go.mod                       # ZERO deps
│   └── enum.go
│
├── escape/                          # Module: Format-specific escaping
│   ├── go.mod                       # ZERO deps
│   └── escape.go
│
├── cmdguard/                        # Module: CLI flag integration
│   ├── go.mod                       # ZERO deps
│   └── flag.go
│
├── d2/                              # Module: D2 diagram subsystem
│   ├── go.mod                       # depends: root, enum, escape
│   ├── d2.go                        # D2 domain types (D2Node, D2Edge, D2Table)
│   ├── d2_enum.go                   # D2-specific enums
│   ├── d2_render.go                 # D2Diagram + Render()
│   ├── d2_write.go                  # Style/edge writing helpers
│   ├── d2_convert.go                # TableData/Tree → D2 conversion
│   └── ...d2 test files
│
├── graph/                           # Module: DOT + Mermaid graph renderers
│   ├── go.mod                       # depends: root, escape
│   ├── graph_mixin.go               # GraphRendererMixin (extracted from dot.go)
│   ├── dot.go                       # DOT renderer
│   ├── mermaid.go                   # Mermaid renderer
│   └── ...graph test files
│
├── table/                           # Module: Lipgloss terminal tables
│   ├── go.mod                       # depends: root, lipgloss
│   ├── table.go
│   └── ...table test files
│
├── integration/                     # Module: Cross-module integration tests
│   ├── go.mod                       # depends: root, d2, graph, table
│   └── ...integration test files
│
└── examples/                        # Module: Usage examples
    ├── go.mod                       # depends: root, d2, graph, table, cmdguard
    └── ...example files
```

---

## Dependency Graph

```
  enum/      escape/     cmdguard/
  (leaf)      (leaf)      (isolated)
    ↑           ↑
    │           │
    └─────┬─────┘
          │
    ┌─────┴────────────────────────────────────────┐
    │            ROOT (package output)               │
    │                                                │
    │  Format, Renderer, TableData, TreeNode         │
    │  ColorMode, BrandedIDs, SortBy                 │
    │  json, csv, tsv, yaml, xml, markdown            │
    │  html, tree, streaming                          │
    │  GraphNode, GraphEdge, GraphRenderer            │
    │  marshal, delimited, markup, slices             │
    │  registry                                       │
    └──┬──────────┬──────────┬───────────────────────┘
       │          │          │
    ┌──┴──┐  ┌───┴────┐  ┌──┴──────┐
    │ d2/ │  │ graph/  │  │ table/  │  ← only module pulling lipgloss
    └─────┘  │ (dot,   │  └─────────┘
             │  mermaid)│
             └─────────┘

    integration/ → root, d2, graph, table
    examples/    → root, d2, graph, table, cmdguard
```

---

## Detailed Module Specs

### Root Module — `github.com/larsartmann/go-output`

|                    |                                                                                       |
| ------------------ | ------------------------------------------------------------------------------------- |
| **Package**        | `package output` (unchanged)                                                          |
| **Deps**           | `enum`, `escape`, `go-faster/yaml`, `x/term`                                          |
| **What stays**     | All types, interfaces, text/structured formatters, graph types, tree, HTML, streaming |
| **What moves out** | `d2*.go` → `d2/`, `dot.go`+`mermaid.go` → `graph/`                                    |
| **`internal/`**    | Stays in root — test helpers are for root package                                     |

**No changes to package name, no file renames for root files.** Root stays `package output`.

### Module: `enum/` — Generic Enum Utilities

|                 |                                         |
| --------------- | --------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/enum` |
| **Deps**        | None (stdlib only)                      |
| **Change**      | Just add `go.mod`                       |

### Module: `escape/` — Format-Specific Escaping

|                 |                                           |
| --------------- | ----------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/escape` |
| **Deps**        | None (stdlib only)                        |
| **Change**      | Just add `go.mod`                         |

### Module: `cmdguard/` — CLI Flag Integration

|                 |                                                                 |
| --------------- | --------------------------------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/cmdguard`                     |
| **Deps**        | None (stdlib only)                                              |
| **Change**      | Just add `go.mod`                                               |
| **Note**        | Generic `EnumFlag[T]` works with ANY enum type — fully isolated |

### Module: `d2/` — D2 Diagram Subsystem

|                           |                                                                       |
| ------------------------- | --------------------------------------------------------------------- |
| **Module path**           | `github.com/larsartmann/go-output/d2`                                 |
| **Deps**                  | root (`output`), `enum`, `escape`                                     |
| **Files moved from root** | `d2.go`, `d2_enum.go`, `d2_render.go`, `d2_write.go`, `d2_convert.go` |
| **Package name**          | `package d2` (changed from `package output`)                          |

`d2_convert.go` converts `output.TableData`/`output.TreeNode`/`output.GraphNode` → D2 types. One-way dependency: `d2 → root`. Clean.

### Module: `graph/` — DOT + Mermaid Graph Renderers

|                           |                                                               |
| ------------------------- | ------------------------------------------------------------- |
| **Module path**           | `github.com/larsartmann/go-output/graph`                      |
| **Deps**                  | root (`output`), `escape`                                     |
| **Files moved from root** | `dot.go`, `mermaid.go`                                        |
| **New file**              | `graph_mixin.go` (extract `GraphRendererMixin` from `dot.go`) |
| **Package name**          | `package graph` (changed from `package output`)               |

`GraphNode`/`GraphEdge`/`GraphRenderer` types stay in root (shared by d2). The `graph/` module contains only the DOT and Mermaid renderers plus the `GraphRendererMixin` they share.

### Module: `table/` — Lipgloss Terminal Tables

|                 |                                          |
| --------------- | ---------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/table` |
| **Deps**        | root (`output`), `lipgloss`              |
| **Change**      | Add `go.mod`, update import path         |

**Biggest win:** Only module pulling in heavy lipgloss dependency. Users who don't need terminal tables skip it entirely.

### Module: `integration/` — Cross-Module Integration Tests

|                 |                                                |
| --------------- | ---------------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/integration` |
| **Deps**        | root, `d2`, `graph`, `table`                   |

### Module: `examples/` — Usage Examples

|                 |                                             |
| --------------- | ------------------------------------------- |
| **Module path** | `github.com/larsartmann/go-output/examples` |
| **Deps**        | root, `d2`, `graph`, `table`, `cmdguard`    |

### Package: `sort/` — **DEPRECATED**

|                   |                                                                              |
| ----------------- | ---------------------------------------------------------------------------- |
| **Action**        | Add deprecation notice, remove in next major version                         |
| **Reason**        | `slices.SortStableFunc` + `cmp.Compare` (stdlib, Go 1.21+) does the same job |
| **`SortBy` enum** | Evaluate if it has consumers outside `sort/`. If not, deprecate too          |

---

## What Moves vs What Stays

### Files that MOVE to a new module (package name changes)

| From root       | To module | New package     |
| --------------- | --------- | --------------- |
| `d2.go`         | `d2/`     | `package d2`    |
| `d2_enum.go`    | `d2/`     | `package d2`    |
| `d2_render.go`  | `d2/`     | `package d2`    |
| `d2_write.go`   | `d2/`     | `package d2`    |
| `d2_convert.go` | `d2/`     | `package d2`    |
| `dot.go`        | `graph/`  | `package graph` |
| `mermaid.go`    | `graph/`  | `package graph` |

### Files that STAY in root (no changes)

All other root `.go` files. `package output` stays. No renames, no moves.

### Test files that move

| From root                                                                                   | To module |
| ------------------------------------------------------------------------------------------- | --------- |
| `d2_test.go`, `d2_enum_test.go`, `d2_edge_test.go`, `d2_node_test.go`, `d2_convert_test.go` | `d2/`     |
| `dot_test.go`, `mermaid_test.go`, `graph_test.go`                                           | `graph/`  |

---

## Key Coupling Points & Resolutions

### 1. `tableDataBase` (unexported) — shared by html.go and streaming.go

Both stay in root. No issue.

### 2. `marshal()` / `unmarshal()` (unexported) — used by json.go, yaml.go

All stay in root. No issue.

### 3. `writeMarkupRow()` (unexported) — used by xml.go, html.go

All stay in root. No issue.

### 4. `GraphRendererMixin` defined in `dot.go` — used by `mermaid.go`

Both move to `graph/` module. Extract `GraphRendererMixin` into `graph_mixin.go`. Clean.

### 5. `table/table.go` imports root `output.Renderer`

After split, `table/table.go` stays as-is — root `go.mod` is still `github.com/larsartmann/go-output`. The import path doesn't change. Only `table/go.mod` is added.

### 6. `d2_convert.go` uses `GraphNode`, `GraphShape`, `NodesFromTableData` from root

After split, imports change from same-package references to `output "github.com/larsartmann/go-output"`. Clean one-way dep: `d2 → root`.

---

## `go.work` File

```go
go 1.26.2

use (
    .
    ./enum
    ./escape
    ./cmdguard
    ./d2
    ./graph
    ./table
    ./integration
    ./examples
)
```

The `.` entry includes the root module itself. With `go.work`, `go test ./...` at root tests the entire workspace.

## `replace` Directives Pattern

Each module that depends on siblings uses `replace` directives for standalone development:

```go
// Example: d2/go.mod
module github.com/larsartmann/go-output/d2

go 1.26.2

require (
    github.com/larsartmann/go-output v0.0.0
    github.com/larsartmann/go-output/enum v0.0.0
    github.com/larsartmann/go-output/escape v0.0.0
)

replace (
    github.com/larsartmann/go-output => ../
    github.com/larsartmann/go-output/enum => ../enum
    github.com/larsartmann/go-output/escape => ../escape
)
```

```go
// Example: graph/go.mod
module github.com/larsartmann/go-output/graph

go 1.26.2

require (
    github.com/larsartmann/go-output v0.0.0
    github.com/larsartmann/go-output/escape v0.0.0
)

replace (
    github.com/larsartmann/go-output => ../
    github.com/larsartmann/go-output/escape => ../escape
)
```

```go
// Example: table/go.mod
module github.com/larsartmann/go-output/table

go 1.26.2

require (
    github.com/larsartmann/go-output v0.0.0
    charm.land/lipgloss/v2 v2.0.3
)

replace (
    github.com/larsartmann/go-output => ../
)
```

**Note:** With `go.work` present, `replace` directives are redundant during workspace development. But they allow `cd d2 && go test ./...` to work standalone — same pattern as go-cqrs-lite.

---

## Pareto-Sorted Execution Plan

Each step is self-contained and leaves the project in a working, committed state.

### Phase 1: Leaf Modules (zero risk, immediate value)

| Step    | What                                                         | Effort | Impact                 |
| ------- | ------------------------------------------------------------ | ------ | ---------------------- |
| **1.1** | Create `go.work` at root                                     | 5 min  | Foundation             |
| **1.2** | Add `go.mod` to `enum/`                                      | 5 min  | Isolates zero-dep leaf |
| **1.3** | Add `go.mod` to `escape/`                                    | 5 min  | Isolates zero-dep leaf |
| **1.4** | Add `go.mod` to `cmdguard/`                                  | 5 min  | Isolates zero-dep leaf |
| **1.5** | Update root `go.mod` with `replace` for enum/escape/cmdguard | 5 min  | Root still works       |
| **1.6** | Verify: `go build ./...` + `go test ./...` pass              | 5 min  | Confidence             |

### Phase 2: Deprecate `sort/` (low risk, removes dead code)

| Step    | What                                                             | Effort | Impact                                |
| ------- | ---------------------------------------------------------------- | ------ | ------------------------------------- |
| **2.1** | Add deprecation notice to `sort/sorter.go` and `sort/compare.go` | 5 min  | Users know to migrate                 |
| **2.2** | Audit `SortBy` enum consumers — deprecate if only used by sort/  | 10 min | Remove dead enum                      |
| **2.3** | Update integration tests that use sort/ to use stdlib directly   | 15 min | Tests don't depend on deprecated code |
| **2.4** | Verify: `go test ./...` passes                                   | 5 min  | Confidence                            |

### Phase 3: Extract `table/` as module (medium risk, biggest dep win)

| Step    | What                                                             | Effort | Impact                          |
| ------- | ---------------------------------------------------------------- | ------ | ------------------------------- |
| **3.1** | Add `go.mod` to `table/` with replace directive pointing to root | 5 min  | Lipgloss becomes opt-in         |
| **3.2** | Update root `go.mod` with `replace` for table                    | 5 min  | Root works with table as module |
| **3.3** | Verify: `go test ./table/...` and `go test ./...` pass           | 5 min  | Confidence                      |

### Phase 4: Extract `d2/` as module (medium effort, clean domain boundary)

| Step    | What                                                                            | Effort | Impact          |
| ------- | ------------------------------------------------------------------------------- | ------ | --------------- |
| **4.1** | Create `d2/` directory + `go.mod`                                               | 5 min  | Foundation      |
| **4.2** | Move `d2*.go` files from root to `d2/`                                          | 5 min  | D2 isolated     |
| **4.3** | Change `package output` → `package d2` in moved files                           | 10 min | Package rename  |
| **4.4** | Update imports: same-package refs → `output "github.com/larsartmann/go-output"` | 15 min | Fix compilation |
| **4.5** | Move d2 test files to `d2/` + update their imports                              | 10 min | Tests follow    |
| **4.6** | Verify: `go test ./d2/...` passes                                               | 5 min  | Confidence      |

### Phase 5: Extract `graph/` as module (medium effort, DOT/Mermaid isolation)

| Step    | What                                                                      | Effort | Impact               |
| ------- | ------------------------------------------------------------------------- | ------ | -------------------- |
| **5.1** | Create `graph/` directory + `go.mod`                                      | 5 min  | Foundation           |
| **5.2** | Move `dot.go` and `mermaid.go` to `graph/`                                | 5 min  | DOT/Mermaid isolated |
| **5.3** | Extract `GraphRendererMixin` from `dot.go` into `graph_mixin.go`          | 10 min | Clean separation     |
| **5.4** | Change `package output` → `package graph` + fix imports                   | 15 min | Compilation          |
| **5.5** | Move graph test files (`dot_test.go`, `mermaid_test.go`, `graph_test.go`) | 10 min | Tests follow         |
| **5.6** | Verify: `go test ./graph/...` passes                                      | 5 min  | Confidence           |

### Phase 6: Integration & Examples (low effort, high confidence)

| Step    | What                                                      | Effort | Impact                 |
| ------- | --------------------------------------------------------- | ------ | ---------------------- |
| **6.1** | Add `go.mod` to `integration/` with all module deps       | 10 min | Integration tests work |
| **6.2** | Update integration test imports for new module paths      | 15 min | Fix compilation        |
| **6.3** | Add `go.mod` to `examples/` + update imports              | 15 min | Examples work          |
| **6.4** | Full workspace verify: `go build ./...` + `go test ./...` | 10 min | Complete confidence    |

### Phase 7: Polish

| Step    | What                                            | Effort | Impact            |
| ------- | ----------------------------------------------- | ------ | ----------------- |
| **7.1** | Write ADR: `docs/adr/001-multi-module-split.md` | 15 min | Document decision |
| **7.2** | Update justfile for multi-module commands       | 10 min | Dev workflow      |
| **7.3** | Update AGENTS.md with new structure             | 10 min | AI context        |
| **7.4** | Update README.md with new module paths          | 10 min | User-facing docs  |
| **7.5** | Final verify: build + test + lint               | 10 min | Done              |

---

## Open Questions

1. **Should `SortBy` enum stay in root after deprecating `sort/`?** If it has no consumers outside `sort/`, it can go too. Need audit.

2. **Should `GraphNode`/`GraphEdge` stay in root or move to `graph/`?** D2 module also uses them. If in `graph/`, then `d2/` depends on both root and `graph/`. If in root, `graph/` depends on root only. **Recommendation:** Keep in root — shared vocabulary like `TableData`/`TreeNode`.

3. **Version strategy?** All modules share the same git tag, or independent semver? **Recommendation:** Shared tags (same as go-cqrs-lite). Simpler for a small library.

4. **How to handle public API constructors that return D2/DOT/Mermaid types?** E.g., `NewDOTRenderer()` is currently in root but returns a `DOTRenderer`. After move, root can't construct `graph.DOTRenderer` directly. Options: (a) Users import the graph module directly, (b) root re-exports. **Recommendation:** (a) — users import the module they need.

---

## Key Benefits Summary

| User wants...                        | They import...                  | Transitive deps            |
| ------------------------------------ | ------------------------------- | -------------------------- |
| JSON/YAML/CSV/XML/Markdown/HTML/Tree | `go-output` (root)              | enum, escape, yaml, x/term |
| Terminal tables                      | `go-output` + `go-output/table` | + lipgloss (heavy)         |
| D2 diagrams                          | `go-output` + `go-output/d2`    | + enum, escape             |
| DOT/Mermaid graphs                   | `go-output` + `go-output/graph` | + escape                   |
| CLI flag parsing                     | `go-output/cmdguard`            | **ZERO** transitive deps   |
| Enum utilities                       | `go-output/enum`                | **ZERO** transitive deps   |
| HTML escaping                        | `go-output/escape`              | **ZERO** transitive deps   |

**Biggest win: `table/` isolation.** Lipgloss is by far the heaviest dependency. Most CLI apps using go-output for JSON/YAML/CSV don't need it at all.

**Second biggest win: D2/graph isolation.** Users who only need structured data output (JSON/YAML/CSV) don't pull in diagram rendering code.

**Simplest win: Deprecating `sort/`.** Removes a reinvented wheel and its coupling to the `SortBy` enum.
