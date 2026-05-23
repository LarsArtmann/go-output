# Go Modularization Proposal — go-output (v2)

**Date:** 2026-05-23
**Status:** DRAFT
**Deciders:** Lars Artmann
**Supersedes:** 2026-05-16 proposal (stale — referenced deleted cmdguard/, internal/testutils/)

---

## 1. Executive Summary

go-output is a **partially modularized** multi-module Go workspace (ADR 001, accepted 2026-05-07). The initial split achieved lipgloss isolation in `table/`. However, the **root module remains a god-package** at 4,345 production LOC across 31 files, containing all core types, 7 table formatters, tree rendering, 3 graph renderers, a full D2 diagram system, streaming, and shared infrastructure.

This proposal extracts two natural module boundaries (`d2/` and `graph/`) while solving the **critical coupling issue** in `render_tabledata.go` that the previous proposal missed.

**What changes:**

- Extract `d2/` as new module (835 LOC, 5 production files)
- Extract `graph/` as new module (319 LOC, 2 production files for DOT+Mermaid)
- Introduce `TableDataRenderFunc` registry extension to break the root→d2/graph cycle in `render_tabledata.go`
- Move graph-specific test helpers from `output_test_helpers.go` to `graph/`

**What stays:**

- Root module remains `package output` — core types, interfaces, simple formatters, streaming, tree, registry
- `enum/`, `escape/`, `testhelpers/`, `sort/`, `table/` unchanged
- Core types (`GraphNode`, `GraphEdge`, `GraphRenderer`, `TreeNode`, `TableData`, `Format`, `Shape`, branded IDs) stay in root

**Expected benefits:**

- D2 users get an independent module with zero coupling to root formatters
- Graph (DOT + Mermaid) users get an independent module
- Root module shrinks from 4,345 → ~3,191 production LOC (excluding test helpers)
- Clean DAG: `d2/` and `graph/` depend on root, never the reverse
- No circular dependencies via registry-based dispatch

---

## 2. Current State Analysis

### 2.1 Existing Module Landscape (as of 2026-05-23)

| Module                  | Path             | Internal Deps                                       | External Deps                         | Replace Directives | State                                    |
| ----------------------- | ---------------- | --------------------------------------------------- | ------------------------------------- | ------------------ | ---------------------------------------- |
| Root (`package output`) | `./`             | enum, escape, testhelpers                           | go-faster/yaml, x/term, go-branded-id | 4 replace          | **Leaky** — sort in go.mod but test-only |
| `enum/`                 | `./enum/`        | testhelpers (tests)                                 | None                                  | 1 replace          | Clean                                    |
| `escape/`               | `./escape/`      | None                                                | None                                  | None               | Clean                                    |
| `testhelpers/`          | `./testhelpers/` | None                                                | None                                  | None               | Clean                                    |
| `sort/`                 | `./sort/`        | None                                                | None                                  | None               | Deprecated, clean                        |
| `table/`                | `./table/`       | root (prod), enum, escape, testhelpers (transitive) | lipgloss/v2                           | 4 replace          | Clean                                    |
| `integration/`          | `./integration/` | root, table                                         | None (transitive)                     | 5 replace          | Clean                                    |
| `examples/`             | `./examples/`    | root, table                                         | None (transitive)                     | 5 replace          | Clean                                    |

**Changes since v1 proposal:**

- `cmdguard/` deleted (commit 10e8104) — zero consumers
- `internal/testutils/` deleted (commit 291639a) — replaced by public `testhelpers/`
- `internal/gentest/` still exists but assertions migrated to `testhelpers/`
- Root production LOC grew from 3,587 → 4,345
- Module count: 8 (was 10 in v1 proposal due to removed cmdguard)

### 2.2 Root Module Concern Clusters

| Cluster                            | Files                                                                                               | LOC    |
| ---------------------------------- | --------------------------------------------------------------------------------------------------- | ------ |
| **Core types + interfaces**        | format.go, ids.go, color.go, sort.go, slices.go, registry.go, format_deprecated.go                  | ~700   |
| **Table data infrastructure**      | tabledata.go, render_tabledata.go                                                                   | ~440   |
| **Table formatters**               | json.go, csv.go, tsv.go, markdown.go, html.go, yaml.go, xml.go, delimited.go, markup.go, marshal.go | ~1,359 |
| **JSON/YAML tree+graph renderers** | json_renderers.go, yaml_renderers.go                                                                | ~308   |
| **Tree formatter**                 | tree.go                                                                                             | ~181   |
| **Graph core types**               | graph.go                                                                                            | ~249   |
| **D2 diagram rendering**           | d2.go, d2_enum.go, d2_render.go, d2_write.go, d2_convert.go                                         | ~835   |
| **DOT + Mermaid rendering**        | dot.go, mermaid.go                                                                                  | ~319   |
| **Streaming**                      | streaming.go                                                                                        | ~198   |
| **Test helpers**                   | output_test_helpers.go                                                                              | ~180   |

**Total root production LOC: 4,345** (31 files)

### 2.3 Dependency Graph (Current)

```
enum ─────────────────────────────┐
escape ───────────────────────────┤
testhelpers ──────────────────────┤
                                  ▼
sort ──► root (package output)    │
          │                       │
          ├── Core: Format, Renderer, TableData, TreeNode
          ├── GraphNode, GraphEdge, GraphRenderer, BrandedID
          ├── Formatters: JSON, CSV, TSV, Markdown, HTML, YAML, XML
          ├── Tree: ASCIITreeRenderer
          ├── Graph: DOT, Mermaid, GraphRendererMixin
          ├── D2: D2Diagram (5 files, rich domain model)
          ├── Streaming: StreamingHTMLRenderer
          └── Dispatcher: RenderTableData → calls D2, DOT, Mermaid directly
                │
table ◄──── root ────► integration
                │
            examples
```

### 2.4 Critical Coupling: `render_tabledata.go`

The **previous proposal missed** that `render_tabledata.go` directly calls:

- `D2FromTableData(data)` → defined in `d2_convert.go`
- `MermaidFlowchartRenderer(data)` → defined in `mermaid.go`
- `DOTFromTableData(data)` → defined in `dot.go`

If D2 and graph move to separate modules, root would need to import `d2` and `graph` — creating a **cycle** since those modules depend on root for core types.

**Resolution:** Introduce a `TableDataRenderFunc` registration mechanism (see Section 3.5).

### 2.5 God-Package Analysis

The root module has **31 production files** spanning **7 distinct concerns**. All are cohesive internally but coupled through `package output`.

The most extractable concerns are:

1. **D2** (835 LOC, 5 files) — richest domain model, most self-contained
2. **DOT + Mermaid** (319 LOC, 2 files) — share `GraphRendererMixin`, identical dependency profile

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

#### New Modules

| Module       | Path       | Purpose                 | Production Deps    | Public API                                                               |
| ------------ | ---------- | ----------------------- | ------------------ | ------------------------------------------------------------------------ |
| **`d2/`**    | `./d2/`    | D2 diagram rendering    | root, enum, escape | D2Diagram, D2Node, D2Edge, D2Direction, D2NodeShape, all D2 constructors |
| **`graph/`** | `./graph/` | DOT + Mermaid rendering | root, enum, escape | DOTRenderer, MermaidRenderer, GraphRendererMixin, all graph constructors |

#### Root Module (After Extraction)

| Cluster                        | Files                                                                                               | LOC        |
| ------------------------------ | --------------------------------------------------------------------------------------------------- | ---------- |
| Core types + interfaces        | format.go, ids.go, color.go, sort.go, slices.go, registry.go, format_deprecated.go                  | ~700       |
| Table data infrastructure      | tabledata.go, render_tabledata.go                                                                   | ~440       |
| Table formatters               | json.go, csv.go, tsv.go, markdown.go, html.go, yaml.go, xml.go, delimited.go, markup.go, marshal.go | ~1,359     |
| JSON/YAML tree+graph renderers | json_renderers.go, yaml_renderers.go                                                                | ~308       |
| Tree formatter                 | tree.go                                                                                             | ~181       |
| Graph core types               | graph.go (stays — core types used by D2 too)                                                        | ~249       |
| Streaming                      | streaming.go                                                                                        | ~198       |
| Test helpers                   | output_test_helpers.go (graph helpers removed)                                                      | ~80        |
| **Total**                      |                                                                                                     | **~3,515** |

`GraphRendererMixin` moves to `graph/` — only used by DOT and Mermaid.

### 3.2 What Stays in Root

These types must remain in root because they're used by multiple downstream modules:

- `GraphNode`, `GraphEdge`, `GraphShape`, `GraphStyle` — used by D2, DOT, Mermaid, JSON, YAML
- `GraphRenderer` interface — implemented by D2, DOT, Mermaid
- `TreeNode`, `TreeOutputRenderer` — used by DOT, Mermaid, D2, tree
- `TableData` — used by all formatters
- `Format`, `Shape`, `Renderer`, `TableRenderer` — core interfaces
- `BrandedID` re-exports and brand types (`D2NodeIDBrand`, `GraphNodeIDBrand`, etc.) — type definitions
- `AddTreeNodes`, `NodesFromTableData` — utility functions used by DOT, Mermaid, D2
- Format constants (`FormatD2`, `FormatMermaid`, `FormatDOT`) — enum values

### 3.3 DAG Verification

```
Level 0 (zero deps):    enum, escape, testhelpers
Level 1 (root):         root → enum, escape, testhelpers, go-branded-id, x/term, go-faster/yaml
Level 2 (format mods):  d2 → root, enum, escape
                         graph → root, enum, escape
                         table → root, lipgloss
Level 3 (consumers):    integration → root, table
                         examples → root, table
                         sort (deprecated, zero deps)
```

**Cycles:** None. Root never imports `d2/`, `graph/`, `table/`, `integration/`, or `examples/`.

### 3.4 Replace / Workspace Strategy

**Keep current approach:** Replace directives in each consuming module's `go.mod`. No committed `go.work`.

**Rationale:**

- `cd d2 && go test ./...` works standalone
- Published versions are clean (no go.work)
- Each module is self-contained for development

**New files:**

- `d2/go.mod` — replace for root, enum, escape, testhelpers
- `graph/go.mod` — replace for root, enum, escape, testhelpers

**Updated files:**

- `integration/go.mod` — add d2, graph to replace
- `examples/go.mod` — add d2, graph to replace

### 3.5 TableDataRenderFunc Registry (Breaking the Cycle)

**Problem:** `render_tabledata.go` directly calls `D2FromTableData()`, `MermaidFlowchartRenderer()`, and `DOTFromTableData()`. If these move to separate modules, root would need to import them → cycle.

**Solution:** Extend the registry with a `TableDataRenderFunc` type:

```go
// In root's registry.go (or render_tabledata.go):
type TableDataRenderFunc func(w io.Writer, data *TableData, opts RenderOptions) error

var tableDataRenderers = make(map[Format]TableDataRenderFunc)

func RegisterTableDataRenderer(format Format, fn TableDataRenderFunc) {
    // register
}

// RenderTableData checks tableDataRenderers for D2/Mermaid/DOT cases
```

Then in `d2/` and `graph/` modules, each registers its render function:

```go
// In d2/init.go or via explicit call:
func init() {
    output.RegisterTableDataRenderer(output.FormatD2, renderD2TableData)
}
```

**Alternative (simpler, preferred):** Don't use `init()`. Instead:

1. Root's `RenderTableData` returns `UnsupportedFormatError` for D2/Mermaid/DOT after extraction
2. D2 and graph modules provide their own `RenderTableData(data) error` convenience functions
3. Callers who need D2/Mermaid/DOT import the module directly

**Decision: Option 2 (UnsupportedFormatError approach).**

**Rationale:**

- Simpler — no registry extension needed
- Honest — `RenderTableData` in root only handles root-owned formats
- D2/graph modules already have `D2FromTableData()`, `DOTFromTableData()`, `MermaidFlowchartRenderer()` which accept `TableData` directly
- Users who need D2/graph already import those modules
- Avoids `init()` side effects and global state

**Impact:** `render_tabledata.go` loses ~50 LOC (3 render functions). `render_tabledata_test.go` loses ~35 LOC (3 test cases). These move to their respective modules as tests.

### 3.6 Test Dependency Isolation

| Module         | Production Deps                                                  | Test-Only Deps               |
| -------------- | ---------------------------------------------------------------- | ---------------------------- |
| `enum/`        | —                                                                | testhelpers                  |
| `escape/`      | —                                                                | —                            |
| `testhelpers/` | —                                                                | —                            |
| `root`         | enum, escape, testhelpers, go-branded-id, x/term, go-faster/yaml | —                            |
| `d2/`          | root, enum, escape                                               | testhelpers (for assertions) |
| `graph/`       | root, enum, escape                                               | testhelpers (for assertions) |
| `table/`       | root, lipgloss                                                   | —                            |
| `sort/`        | —                                                                | —                            |
| `integration/` | root, table                                                      | —                            |
| `examples/`    | root, table                                                      | —                            |

### 3.7 Interface Extraction

No interface extraction needed. The current design uses clean interfaces already:

- `Renderer`, `GraphRenderer`, `TableRenderer`, `TreeOutputRenderer` — all defined in root
- D2 and graph modules import root's interfaces and implement them
- Root never imports D2 or graph

### 3.8 Versioning Strategy

**Shared versioning** (unchanged):

- All modules share a single git tag `v1.2.3`
- Root module is the primary published artifact
- Sub-modules are opt-in extensions
- No independent semver — modules are tightly coupled to root's core types

---

## 4. Migration Strategy

### 4.1 Prerequisites

- All tests pass: `go test ./...` in each module
- Clean lint: `golangci-lint run ./...`
- Create branch: `git checkout -b modularize/extract-d2-graph`

### 4.2 Ordered Steps

#### Step 1: Extract `d2/` Module (Highest Value)

- Move `d2.go`, `d2_enum.go`, `d2_render.go`, `d2_write.go`, `d2_convert.go` → `d2/`
- Move D2 test files → `d2/`
- Rename package: `output` → `d2`
- Create `d2/go.mod`
- Update imports (root types → `output.` prefix)
- Remove `renderD2TableData` from `render_tabledata.go` — return `UnsupportedFormatError` for `FormatD2`
- Move D2 render test case from `render_tabledata_test.go` → `d2/`
- Update `examples/`, `integration/` go.mod files
- Update `benchmarks_test.go` if needed (D2 benchmarks inline)

**Verification:** `cd d2 && go test ./...` and root `go test ./...` both pass

#### Step 2: Extract `graph/` Module

- Move `dot.go`, `mermaid.go` → `graph/` (graph.go stays in root)
- Move `GraphRendererMixin` with `dot.go` to `graph/`
- Move graph test files (`dot_test.go`, `mermaid_test.go`, `graph_test.go`) → `graph/`
- Split `output_test_helpers.go`: graph-related helpers move to `graph/`
- Rename package: `output` → `graph`
- Create `graph/go.mod`
- Update imports
- Remove `renderMermaidTableData`, `renderDOTTableData` from `render_tabledata.go`
- Move Mermaid/DOT render test cases from `render_tabledata_test.go` → `graph/`
- Update `examples/`, `integration/` go.mod files
- Update `benchmarks_test.go`: `NewMermaidRenderer()` → `graph.NewMermaidRenderer()`, `NewDOTRenderer()` → `graph.NewDOTRenderer()`

**Verification:** `cd graph && go test ./...` and root `go test ./...` both pass

#### Step 3: Clean Up Sort Dependency

- Remove sort from root's production `go.mod` if present (it's test-only in `userjourney_test.go`)
- Run `go mod tidy`

**Verification:** Root `go.mod` clean

#### Step 4: Update Documentation

- Update `AGENTS.md` with new module table (8 → 10 modules)
- Update `docs/modularization/DEPENDENCY_GRAPH.md`
- Update README examples if needed
- Update `docs/FORMAT_ARCHITECTURE.md`

---

## 5. Risk Assessment

| Risk                                               | Likelihood | Impact | Mitigation                                             |
| -------------------------------------------------- | ---------- | ------ | ------------------------------------------------------ |
| D2 package rename breaks consumers                 | Medium     | High   | Pre-v1 library; breaking change acceptable             |
| Graph package rename breaks consumers              | Medium     | High   | Same — pre-v1                                          |
| `RenderTableData` no longer handles D2/Mermaid/DOT | Medium     | Medium | Document clearly; modules provide equivalent functions |
| Integration tests fail after moves                 | Low        | Medium | Step-by-step extraction with test verification         |
| Root module still large after extraction           | Low        | Low    | ~3,515 LOC acceptable for core library with 9 formats  |
| go.mod replace directive sprawl                    | Medium     | Low    | 10 modules is manageable; consider go.work if it grows |

---

## 6. Key Decisions

1. **Root IS the core** (ADR 001) — unchanged
2. **D2 gets its own module** — 835 LOC, rich domain model
3. **DOT + Mermaid share `graph/` module** — share `GraphRendererMixin`, identical deps
4. **`GraphNode`/`GraphEdge` stay in root** — used by D2, JSON, YAML too
5. **`GraphRendererMixin` moves to `graph/`** — only DOT/Mermaid use it
6. **`render_tabledata.go` drops D2/Mermaid/DOT cases** — returns `UnsupportedFormatError` for those formats; avoids cycle
7. **Package rename required** — `package d2`, `package graph` (breaking, acceptable pre-v1)
8. **No `init()` registration** — simpler, no global state side effects
9. **`sort/` dependency cleaned from root** — test-only
10. **Test helpers split** — graph helpers move to `graph/`, rest stay in root

---

## 7. Build System Impact

- **flake.nix** — Update to build/test `d2/` and `graph/` independently
- **CI/CD** — Add parallel jobs for `d2/`, `graph/`
- **golangci-lint** — Already workspace-wide, no changes
- **go.work** — Add new modules if re-created for local dev

---

## 8. What Does NOT Change

- `enum/`, `escape/`, `testhelpers/` — unchanged
- `table/` — unchanged
- `sort/` — unchanged (still deprecated)
- `internal/gentest/` — stays in root
- Root module name — still `github.com/larsartmann/go-output`, `package output`
- Format enum — still in root
- `Renderer`/`TableRenderer`/`GraphRenderer` interfaces — still in root
- Core types (Format, TableData, TreeNode, GraphNode, GraphEdge) — still in root
- JSON, CSV, TSV, Markdown, HTML, YAML, XML formatters — still in root
- JSON/YAML tree and graph renderers — still in root
- Streaming — still in root
- Registry — still in root
