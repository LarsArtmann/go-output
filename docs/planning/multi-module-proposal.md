# go-output: Multi-Module Migration Proposal

**Created:** 2026-05-07
**Status:** DRAFT — Comprehensive Refined Proposal v2
**Based on:** Deep analysis of current codebase + patterns from go-cqrs-lite + Go workspace best practices

---

## Table of Contents

- [Self-Critique of v1 Proposal](#self-critique-of-v1-proposal)
- [Current State Analysis](#current-state-analysis)
- [Proposed Module Structure](#proposed-module-structure)
- [Dependency Graph](#dependency-graph)
- [Detailed Module Specs](#detailed-module-specs)
- [What Stays in Root vs What Moves](#what-stays-in-root-vs-what-moves)
- [Key Coupling Points & Resolutions](#key-coupling-points--resolutions)
- [Pareto-Sorted Execution Plan](#pareto-sorted-execution-plan)
- [Open Questions](#open-questions)

---

## Self-Critique of v1 Proposal

The first draft had several significant issues:

| Issue | Why it matters | Resolution in v2 |
|---|---|---|
| **`enum/` and `escape/` shown with no .go files** | They ALREADY have .go files — the proposal was misleading | They're already packages; they just need their own `go.mod` |
| **Missed reverse dependencies** | `sort/sorter.go` imports root (for `SortBy`), `table/table.go` imports root (for `Renderer`) | Must extract shared types first, or move them into sub-modules that depend on a core |
| **Missed unexported type coupling** | `streaming.go` embeds unexported `tableDataBase` from `html.go` — can't split across modules | Keep them together or export the type |
| **`GraphRendererMixin` defined in wrong file** | Defined in `dot.go` but used by `mermaid.go` — belongs with graph layer | Move to `graph.go` (or new file in graph module) |
| **Missed test infrastructure** | `internal/gentest/` (zero deps) and `internal/testutils/` (depends on root) need homes | Dedicated `testhelpers` module, like go-cqrs-lite pattern |
| **Missed `format_deprecated.go`** | Backward-compat shim that must stay with `format.go` | Documented explicitly |
| **No `replace` directives** | go-cqrs-lite uses them for standalone module development | Include in plan |
| **No migration path** | Jumped to final state with no steps | Pareto-sorted execution plan with small increments |

---

## Current State Analysis

### Existing Sub-Packages (already have directories + .go files)

| Package | Files | Imports from root | Imports from siblings | External deps |
|---|---|---|---|---|
| `enum/` | `enum.go` | **None** | None | None |
| `escape/` | `escape.go` | **None** | None | None |
| `sort/` | `sorter.go`, `compare.go` | `output.SortBy` | None | None |
| `table/` | `table.go` | `output.Renderer` (interface check) | None | `lipgloss/v2` |
| `cmdguard/` | `flag.go` | **None** (only uses `fmt`) | None | None |
| `internal/gentest/` | `struct.go`, `assert.go`, `isvalid.go` | **None** | None | None |
| `internal/testutils/` | `test_helpers.go` | Root types (`GraphNode`, `TableData`, etc.) | `gentest` | None |

### Root Package Internal Dependencies (within same package)

```
                    ids.go  ← foundation (zero deps)
                      ↑
              format.go    ← Renderer, TableData, TreeNode, Format, interfaces
             ↗  ↑  ↑  ↑  ↑
     registry  tree  markdown  html.go ← defines tableDataBase (unexported)
                         ↑           ↑
                     streaming.go ───┘ (embeds tableDataBase)
                    
              graph.go     ← GraphNode, GraphEdge, GraphRenderer
             ↗    ↑    ↖
        dot.go  d2_convert  mermaid.go   [GraphRendererMixin in dot.go — wrong place!]
                  ↑
          d2_render.go ← d2_write.go
              ↑
          d2_enum.go + d2.go ← ids.go
                    
              marshal.go ← json.go, yaml.go      [unexported marshal()/unmarshal()]
              delimited.go ← csv.go, tsv.go       [exported DelimitedWriter]
              markup.go ← xml.go, html.go          [unexported writeMarkupRow()]
```

### Third-Party Dependencies (only 3!)

| Dependency | Used by | Heavy? |
|---|---|---|
| `charm.land/lipgloss/v2` | `table/table.go` only | **Yes** — pulls in many transitive deps |
| `github.com/go-faster/yaml` | `yaml.go` only | Medium |
| `golang.org/x/term` | `color.go` only | Light |

---

## Proposed Module Structure

```
go-output/
├── go.work                          # Workspace file tying all modules together
│
├── enum/                            # Module 1: Generic enum utilities
│   ├── go.mod                       # ZERO deps
│   └── enum.go                      # (already exists — just add go.mod)
│
├── escape/                          # Module 2: Format-specific escaping
│   ├── go.mod                       # ZERO deps
│   └── escape.go                    # (already exists — just add go.mod)
│
├── cmdguard/                        # Module 3: CLI flag integration
│   ├── go.mod                       # ZERO deps
│   └── flag.go                      # (already exists — just add go.mod)
│
├── core/                            # Module 4: Shared types + all text/structured formatters
│   ├── go.mod                       # depends: enum, escape, go-faster/yaml, x/term
│   ├── format.go                    # Format enum, Renderer, TableData, TreeNode, TableRenderer
│   ├── format_deprecated.go         # OutputFormat backward compat
│   ├── color.go                     # ColorMode + terminal detection
│   ├── ids.go                       # Branded ID phantom types
│   ├── sort.go                      # SortBy enum
│   ├── registry.go                  # Opt-in renderer registry
│   ├── slices.go                    # FilledStrings
│   ├── marshal.go                   # Shared marshal/unmarshal error wrapping
│   ├── json.go, csv.go, tsv.go, yaml.go
│   ├── xml.go, markdown.go
│   ├── delimited.go                 # Shared CSV/TSV writer
│   ├── markup.go                    # Shared XML/HTML row helpers
│   ├── html.go                      # HTML table + tree + tableDataBase
│   ├── tree.go                      # ASCII tree renderer
│   ├── streaming.go                 # Streaming HTML renderer
│   ├── graph.go                     # GraphNode, GraphEdge, GraphRenderer, AddTreeNodes
│   ├── output_test_helpers.go       # White-box test helpers
│   ├── testing_test.go, benchmarks_test.go, fuzz_test.go
│   ├── ...other root test files
│   └── internal/                    # stays internal to core
│       ├── gentest/                 # Generic test helpers (zero deps)
│       └── testutils/               # Domain-aware test helpers
│
├── d2/                              # Module 5: D2 diagram subsystem
│   ├── go.mod                       # depends: core, enum, escape
│   ├── d2.go                        # D2 domain types
│   ├── d2_enum.go                   # D2-specific enums
│   ├── d2_render.go                 # D2Diagram + Render()
│   ├── d2_write.go                  # Style/edge writing helpers
│   ├── d2_convert.go                # TableData/Tree → D2 conversion
│   └── ...d2 test files
│
├── graph/                           # Module 6: DOT + Mermaid graph renderers
│   ├── go.mod                       # depends: core, escape
│   ├── graph_mixin.go               # GraphRendererMixin (moved from dot.go!)
│   ├── dot.go                       # DOT renderer
│   ├── mermaid.go                   # Mermaid renderer
│   └── ...graph test files
│
├── sort/                            # Module 7: Generic sorting
│   ├── go.mod                       # depends: core (for SortBy type)
│   ├── sorter.go                    # (already exists — just add go.mod)
│   ├── compare.go                   # (already exists)
│   └── ...sort test files
│
├── table/                           # Module 8: Lipgloss terminal tables
│   ├── go.mod                       # depends: core (for Renderer), lipgloss
│   ├── table.go                     # (already exists — just add go.mod)
│   └── ...table test files
│
├── integration/                     # Module 9: Cross-module integration tests
│   ├── go.mod                       # depends: core, d2, graph, sort, table
│   └── ...integration test files
│
└── examples/                        # Module 10: Usage examples
    ├── go.mod                       # depends: all modules
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
    ┌─────┴──────────────────────────────────────┐
    │                  core/                       │
    │  Format, Renderer, TableData, TreeNode       │
    │  SortBy, ColorMode, BrandedIDs               │
    │  json, csv, tsv, yaml, xml, markdown          │
    │  html, tree, streaming, graph types           │
    │  marshal, delimited, markup, slices           │
    │  registry, internal/gentest, internal/testutils│
    └──┬──────┬──────┬──────┬──────────────────────┘
       │      │      │      │
    ┌──┴──┐ ┌─┴───┐ ┌┴────┐ ┌┴──────┐
    │ d2/ │ │graph│ │sort/│ │ table/ │ ← only module pulling lipgloss
    └─────┘ │(dot,│ └─────┘ └────────┘
            │merm)│
            └─────┘

    integration/ → core, d2, graph, sort, table
    examples/    → core, d2, graph, sort, table, cmdguard
```

---

## Detailed Module Specs

### Module 1: `enum/` — Generic Enum Utilities

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/enum` |
| **Deps** | None (stdlib only) |
| **Files** | `enum.go` (already exists) |
| **Change needed** | Just add `go.mod` |
| **Why independent** | Leaf package, zero deps. Reusable beyond this project. |

### Module 2: `escape/` — Format-Specific Escaping

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/escape` |
| **Deps** | None (stdlib only) |
| **Files** | `escape.go` (already exists) |
| **Change needed** | Just add `go.mod` |
| **Why independent** | Leaf package, zero deps. Pure string transformations. |

### Module 3: `cmdguard/` — CLI Flag Integration

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/cmdguard` |
| **Deps** | None (stdlib only) |
| **Files** | `flag.go` (already exists) |
| **Change needed** | Just add `go.mod` |
| **Why independent** | Fully isolated. Generic `EnumFlag[T]` works with ANY enum type, not just go-output types. |

### Module 4: `core/` — Shared Types + Formatters

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/core` |
| **Deps** | `enum`, `escape`, `go-faster/yaml`, `x/term` |
| **Files from root** | 28 .go files (see structure above) |
| **Internal packages** | `internal/gentest/`, `internal/testutils/` stay inside core |
| **Package name** | Changes from `package output` to `package core` |

**What lives here and why:**

| Group | Files | Reason |
|---|---|---|
| Core types | `format.go`, `ids.go`, `sort.go`, `color.go` | Shared vocabulary everything depends on |
| Interfaces | `Renderer`, `TableRenderer`, `TreeOutputRenderer`, `GraphRenderer` | Defined in `format.go` and `graph.go` |
| Text formatters | `json.go`, `csv.go`, `tsv.go`, `yaml.go`, `xml.go`, `markdown.go` | No cross-module deps, only depend on core types |
| Tree formatters | `tree.go`, `html.go`, `streaming.go` | `html` and `streaming` share unexported `tableDataBase` — must stay together |
| Graph types | `graph.go` | Generic `GraphNode`/`GraphEdge`/`GraphRenderer` — needed by d2/, graph/ modules |
| Shared helpers | `delimited.go`, `markup.go`, `marshal.go`, `slices.go` | Used by formatters within core |
| Registry | `registry.go` | Depends on `Format` + `Renderer` (both in core) |
| Backward compat | `format_deprecated.go` | Aliases to `format.go` types |

### Module 5: `d2/` — D2 Diagram Subsystem

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/d2` |
| **Deps** | `core`, `enum`, `escape` |
| **Files moved from root** | `d2.go`, `d2_enum.go`, `d2_render.go`, `d2_write.go`, `d2_convert.go` |
| **Package name** | Changes from `package output` to `package d2` |

**Key consideration:** `d2_convert.go` converts `core.TableData`/`core.TreeNode`/`core.GraphNode` → D2 types. This is a one-way dependency (d2 → core), which is clean.

### Module 6: `graph/` — DOT + Mermaid Graph Renderers

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/graph` |
| **Deps** | `core`, `escape` |
| **Files moved from root** | `dot.go`, `mermaid.go` |
| **New file** | `graph_mixin.go` (extract `GraphRendererMixin` from `dot.go`) |
| **Package name** | Changes from `package output` to `package graph` |

**Key change:** `GraphRendererMixin` is currently defined in `dot.go` but used by `mermaid.go`. It must move to its own file (`graph_mixin.go`) or into `graph.go` in core. Since it references `GraphNode`/`GraphEdge` which would live in core, the mixin could live in core too — but since DOT and Mermaid are the only consumers, it makes more sense in the `graph/` module.

**Option A:** Move `GraphNode`/`GraphEdge` types to `graph/` module. Then core only has `Renderer` interface, and graph module owns its own types.
**Option B:** Keep `GraphNode`/`GraphEdge` in core (since `d2_convert.go` also uses them). Graph module depends on core for types.

**Recommended: Option B** — keeps graph types as shared vocabulary, avoids d2 needing to depend on graph module.

### Module 7: `sort/` — Generic Sorting

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/sort` |
| **Deps** | `core` (for `SortBy` type) |
| **Files** | `sorter.go`, `compare.go` (already exist) |
| **Change needed** | Add `go.mod`, change import from `output.SortBy` to `core.SortBy` |

**Key consideration:** `sort/sorter.go` currently imports `output.SortBy`. After split, it would import `core.SortBy`. This is a clean one-way dependency.

### Module 8: `table/` — Lipgloss Terminal Tables

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/table` |
| **Deps** | `core` (for `Renderer` interface check), `lipgloss` |
| **Files** | `table.go` (already exists) |
| **Change needed** | Add `go.mod`, change import from `output.Renderer` to `core.Renderer` |

**This is the biggest win:** Only module pulling in heavy lipgloss dependency. Users who don't need terminal tables skip it entirely.

### Module 9: `integration/` — Cross-Module Integration Tests

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/integration` |
| **Deps** | `core`, `d2`, `graph`, `sort`, `table` |
| **Files** | All files from existing `integration/` directory |
| **Pattern** | Same as go-cqrs-lite — dedicated module for cross-cutting tests |

### Module 10: `examples/` — Usage Examples

| | |
|---|---|
| **Module path** | `github.com/larsartmann/go-output/examples` |
| **Deps** | `core`, `d2`, `graph`, `sort`, `table`, `cmdguard` |
| **Files** | All files from existing `examples/` directory |

---

## What Stays in Root vs What Moves

### Root directory AFTER migration

```
go-output/
├── go.work
├── go.mod + go.sum          # Root module becomes EMPTY or is removed
├── LICENSE, README.md, CHANGELOG.md, AUTHORS, CONTRIBUTING.md
├── .golangci.yml, .gitignore, .gitattributes
├── docs/
│   ├── FORMAT_ARCHITECTURE.md
│   ├── adr/                 # NEW: Architecture Decision Records
│   │   └── 001-multi-module-split.md
│   └── planning/
│       └── multi-module-proposal.md  (this file)
├── enum/                     # Module 1 (just add go.mod)
├── escape/                   # Module 2 (just add go.mod)
├── cmdguard/                 # Module 3 (just add go.mod)
├── core/                     # Module 4 (files moved here)
├── d2/                       # Module 5 (files moved here)
├── graph/                    # Module 6 (files moved here)
├── sort/                     # Module 7 (just add go.mod)
├── table/                    # Module 8 (just add go.mod)
├── integration/              # Module 9 (just add go.mod)
└── examples/                 # Module 10 (just add go.mod)
```

**The root `go.mod` should be REMOVED.** The `go.work` file replaces it for workspace-level builds.

---

## Key Coupling Points & Resolutions

### 1. `tableDataBase` (unexported) — shared by html.go and streaming.go

**Status:** Must stay in same module (core). Both `HTMLRenderer` and `StreamingHTMLRenderer` embed it. No issue since both go to core.

### 2. `marshal()` / `unmarshal()` (unexported) — used by json.go, yaml.go

**Status:** Must stay in same module (core). All three files move together. No issue.

### 3. `writeMarkupRow()` (unexported) — used by xml.go, html.go

**Status:** Must stay in same module (core). Both files move together. No issue.

### 4. `GraphRendererMixin` defined in `dot.go` — used by `mermaid.go`

**Resolution:** Extract `GraphRendererMixin` + `NewGraphRendererMixin` into a new file `graph/graph_mixin.go`. Both `dot.go` and `mermaid.go` move to `graph/` module. Clean.

### 5. `sort/sorter.go` imports root `output.SortBy`

**Resolution:** After split, `sort/sorter.go` imports `core.SortBy`. Change import path. Clean one-way dep.

### 6. `table/table.go` imports root `output.Renderer`

**Resolution:** After split, `table/table.go` imports `core.Renderer`. The `var _ output.Renderer = (*Table)(nil)` becomes `var _ core.Renderer = (*Table)(nil)`. Clean one-way dep.

### 7. `d2_convert.go` uses `GraphNode`, `GraphShape`, `NodesFromTableData` from `graph.go`

**Resolution:** `d2_convert.go` imports these from `core` (since `graph.go` types live in core). Clean one-way dep: `d2 → core`.

### 8. `internal/gentest/` and `internal/testutils/`

**Resolution:** Stay inside `core/internal/`. They're test infrastructure for core. Other modules (d2, graph, sort) write their own test helpers or duplicate minimal assertions. This follows go-cqrs-lite pattern where `testhelpers` is a separate module — but for go-output the volume doesn't justify a full module.

---

## `go.work` File

```go
go 1.26.2

use (
    ./enum
    ./escape
    ./cmdguard
    ./core
    ./d2
    ./graph
    ./sort
    ./table
    ./integration
    ./examples
)
```

## `replace` Directives Pattern

Each module that depends on siblings uses `replace` directives for standalone development:

```go
// Example: d2/go.mod
module github.com/larsartmann/go-output/d2

go 1.26.2

require (
    github.com/larsartmann/go-output/core v0.0.0
    github.com/larsartmann/go-output/enum v0.0.0
    github.com/larsartmann/go-output/escape v0.0.0
)

replace (
    github.com/larsartmann/go-output/core => ../core
    github.com/larsartmann/go-output/enum => ../enum
    github.com/larsartmann/go-output/escape => ../escape
)
```

**Note:** With `go.work` present, `replace` directives are technically redundant during workspace development. But they allow `cd d2 && go test ./...` to work standalone — same pattern as go-cqrs-lite.

---

## Pareto-Sorted Execution Plan

Sorted by: **highest impact, lowest effort first.** Each step is self-contained and leaves the project in a working state.

### Phase 1: Leaf Modules (zero risk, immediate value)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **1.1** | Create `go.work` at root | 5 min | Foundation for all future work | `go.work` |
| **1.2** | Add `go.mod` to `enum/` | 5 min | Isolates zero-dep leaf module | `enum/go.mod` |
| **1.3** | Add `go.mod` to `escape/` | 5 min | Isolates zero-dep leaf module | `escape/go.mod` |
| **1.4** | Add `go.mod` to `cmdguard/` | 5 min | Isolates zero-dep leaf module | `cmdguard/go.mod` |
| **1.5** | Update root `go.mod` to use `replace` for enum/escape/cmdguard | 5 min | Root still works with new modules | `go.mod` |
| **1.6** | Verify: `go build ./...` and `go test ./...` pass | 5 min | Confidence | — |

**Phase 1 result:** 3 leaf modules extracted. Root still works. Users can now `go get github.com/larsartmann/go-output/enum` independently.

### Phase 2: Create `core/` Module (medium effort, high impact)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **2.1** | Create `core/` directory + `go.mod` | 5 min | Foundation for the big move | `core/go.mod` |
| **2.2** | Move all root .go files to `core/` (except tests) | 15 min | Core module gets all formatter code | 28 .go files moved |
| **2.3** | Change `package output` → `package core` in all moved files | 15 min | Required for new package name | All moved files |
| **2.4** | Update all internal imports (`go-output` → `go-output/core`) | 20 min | Fix compilation | All moved files |
| **2.5** | Move root test files to `core/` | 10 min | Tests follow the code | ~15 test files |
| **2.6** | Move `internal/` to `core/internal/` | 5 min | Test helpers follow | `internal/` dir |
| **2.7** | Update `core/go.mod` with deps + replace directives | 10 min | Core can build standalone | `core/go.mod` |
| **2.8** | Verify: `go build ./core/...` and `go test ./core/...` pass | 10 min | Confidence | — |
| **2.9** | Remove old root `go.mod` (or make it a thin re-export module) | 10 min | Clean up | Root `go.mod` |

**Phase 2 result:** All formatter code lives in `core/`. Root is empty or a compat shim.

### Phase 3: Extract Graph Module (medium effort, medium impact)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **3.1** | Create `graph/` directory + `go.mod` | 5 min | Foundation | `graph/go.mod` |
| **3.2** | Move `dot.go` and `mermaid.go` to `graph/` | 5 min | DOT/Mermaid isolated | 2 files |
| **3.3** | Extract `GraphRendererMixin` from `dot.go` into `graph_mixin.go` | 10 min | Clean separation | `graph_mixin.go` |
| **3.4** | Change package name → `package graph` + fix imports | 15 min | Compilation | All graph files |
| **3.5** | Remove graph files from `core/` | 5 min | Core shrinks | — |
| **3.6** | Move graph test files | 10 min | Tests follow | `dot_test.go`, `mermaid_test.go`, `graph_test.go` |
| **3.7** | Verify: `go test ./graph/...` passes | 5 min | Confidence | — |

### Phase 4: Extract D2 Module (medium effort, medium impact)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **4.1** | Create `d2/` directory + `go.mod` | 5 min | Foundation | `d2/go.mod` |
| **4.2** | Move `d2*.go` files from `core/` to `d2/` | 5 min | D2 isolated | 5 files |
| **4.3** | Change package name → `package d2` + fix imports | 15 min | Compilation | All d2 files |
| **4.4** | Move d2 test files | 10 min | Tests follow | ~6 test files |
| **4.5** | Verify: `go test ./d2/...` passes | 5 min | Confidence | — |

### Phase 5: Update sort/ and table/ (low effort, medium impact)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **5.1** | Add `go.mod` to `sort/`, change import to `core.SortBy` | 10 min | sort/ no longer depends on monolith | `sort/go.mod`, `sort/sorter.go` |
| **5.2** | Add `go.mod` to `table/`, change import to `core.Renderer` | 10 min | table/ isolated with lipgloss | `table/go.mod`, `table/table.go` |
| **5.3** | Verify: `go test ./sort/... ./table/...` passes | 5 min | Confidence | — |

### Phase 6: Integration & Examples (low effort, high confidence)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **6.1** | Create `integration/go.mod` with all module deps | 10 min | Integration tests work | `integration/go.mod` |
| **6.2** | Update integration test imports | 15 min | Fix compilation | ~5 test files |
| **6.3** | Create `examples/go.mod` | 10 min | Examples work | `examples/go.mod` |
| **6.4** | Update example imports | 10 min | Fix compilation | ~4 files |
| **6.5** | Full workspace verify: `go build ./...` + `go test ./...` | 10 min | Complete confidence | — |

### Phase 7: Polish (low effort, high polish)

| Step | What | Effort | Impact | Files Changed |
|------|------|--------|--------|---------------|
| **7.1** | Write ADR: `docs/adr/001-multi-module-split.md` | 15 min | Document the decision | New file |
| **7.2** | Update justfile for multi-module commands | 10 min | Dev workflow | `justfile` |
| **7.3** | Update AGENTS.md with new structure | 10 min | AI context | `AGENTS.md` |
| **7.4** | Update README.md with new module paths | 10 min | User-facing docs | `README.md` |
| **7.5** | Remove root `go.mod` (or convert to re-export compat module) | 10 min | Clean slate | Root |
| **7.6** | Final full verify: build + test + lint | 10 min | Done | — |

---

## Open Questions

1. **Root module as re-export shim?** Should the root `go.mod` re-export all types from `core/` for backward compat (`import "github.com/larsartmann/go-output"` still works), or fully remove it? **Recommendation:** Remove it. Breaking change is fine pre-v1.0.

2. **Should `GraphNode`/`GraphEdge` live in `core` or `graph`?** D2 module also uses them for conversion. If in `graph/`, then `d2/` depends on both `core` and `graph`. If in `core`, `graph/` depends on `core` only. **Recommendation:** Keep in `core` — they're shared vocabulary like `TableData`/`TreeNode`.

3. **Should `internal/gentest/` become a shared `testhelpers` module?** go-cqrs-lite has this pattern. But go-output's test helpers are small. **Recommendation:** Keep inside `core/internal/`. Not worth a separate module.

4. **Package naming: `core` vs something more descriptive?** Other options: `output`, `format`, `formats`. Since the project is called "go-output", the core module being `go-output/core` is clear. **Recommendation:** `core`.

5. **Version strategy?** All modules share the same git tag initially, or independent semver? **Recommendation:** Shared tag strategy (same as go-cqrs-lite). Simpler for a small library.

6. **`Format` enum currently in root — where does it go after split?** It's needed by `cmdguard/` (for EnumFlag), `sort/` (not directly, but for registry), `table/` (indirectly). Should it live in `core`? **Recommendation:** Yes, `core`. `cmdguard/` is generic and doesn't actually import `Format` — it's `EnumFlag[T]` that works with any enum.

---

## Key Benefits Summary

| User wants... | They import... | Transitive deps |
|---|---|---|
| JSON output only | `core` | enum, escape, go-faster/yaml |
| Terminal tables | `core` + `table` | + lipgloss (heavy) |
| D2 diagrams | `core` + `d2` | + enum, escape |
| DOT/Mermaid graphs | `core` + `graph` | + escape |
| Sorting utility | `core` + `sort` | + enum |
| CLI flag parsing | `cmdguard` | **ZERO** transitive deps |
| Enum utilities | `enum` | **ZERO** transitive deps |
| HTML escaping | `escape` | **ZERO** transitive deps |

**The biggest win: `table/` isolation.** Lipgloss is by far the heaviest dependency. Most CLI apps using go-output for JSON/YAML/CSV don't need it at all.
