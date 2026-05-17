# Go Modularization Proposal — go-output

**Date:** 2026-05-16
**Status:** DRAFT
**Deciders:** Lars Artmann

**Status:** REVIEWED (self-review complete, see below)

---

## Self-Review (Phase 4)

### Critical Finding: output_test_helpers.go contains graph-specific helpers

`output_test_helpers.go` (176 LOC) contains graph-specific helpers that reference `GraphNode`, `GraphEdge`, DOT, and Mermaid:

- `testNodesAB()`, `testNodesABC()`, `newTestNode()`, `newTestNodeWithShape()` → GraphNode
- `testEdgeAB()`, `testEdgesAB()`, `testEdgesABC()` → GraphEdge
- `testDOTEmptyExpected()` → DOT renderer
- `testMermaidEmptyExpected()` → Mermaid renderer

**Resolution:** Split `output_test_helpers.go`:

- Graph-related helpers move to `graph/` as internal test helpers
- D2 helpers (none currently in this file — D2 has its own test helpers inline)
- Remaining helpers (HTML escape, expected output, tree depth) stay in root

### Decision: gentest/ and testutils/ extraction may be over-engineering

At 151 LOC and 149 LOC respectively, extracting test helpers as full modules adds 2 more `go.mod` files and replace directives for minimal benefit. The cross-module `internal/` leak is a minor inconvenience, not a blocking issue.

**Revised approach:** Keep `internal/gentest` and `internal/testutils` in root module. The `internal/` visibility works correctly — sibling modules can import them. The real value is in extracting `d2/` and `graph/`.

### Final Proposal (Revised)

Extract only 2 new modules:

1. **`d2/`** — D2 diagram rendering (815 LOC → 5 files)
2. **`graph/`** — DOT + Mermaid rendering (566 LOC → 3 files)

Keep `internal/gentest` and `internal/testutils` in root. Fix the sort dependency.

---

## 1. Executive Summary

go-output is already a **partially modularized** multi-module Go workspace (ADR 001 accepted 2026-05-07). The initial split achieved the primary goal: lipgloss isolation in `table/`. However, the **root module remains a god-package** at 3,587 production LOC across 28 files, containing all core types, 7 table formatters, tree rendering, 3 graph renderers, a full D2 diagram system, streaming, and shared infrastructure.

This proposal extracts the remaining natural module boundaries while preserving the existing 7-module structure and the "root IS the core" design from ADR 001.

**What changes:** Extract `d2/` and `graph/` as new modules. Remove the deprecated `sort/` dependency from root's production `go.mod`.

**What stays:** Root module remains `package output` — the formatters (JSON, CSV, TSV, Markdown, XML, YAML, HTML), streaming, tree, core types (Format, TableData, TreeNode, Renderer interfaces), and registry all stay in root. The existing `enum/`, `escape/`, `cmdguard/`, `table/`, `sort/` modules are unchanged.

**Expected benefits:**

- D2 users get an independent module (833 LOC with zero coupling to root formatters)
- Graph (DOT + Mermaid) users get an independent module (568 LOC, only needs root core types)
- Root module shrinks from 3,587 → ~2,186 production LOC
- Clean DAG: `d2/` and `graph/` depend on root, never the reverse
- Each new module can be versioned independently

---

## 2. Current State Analysis

### 2.1 Existing Module Landscape

| Module                  | Path             | Internal Deps               | External Deps                         | Replace Directives | State                                              |
| ----------------------- | ---------------- | --------------------------- | ------------------------------------- | ------------------ | -------------------------------------------------- |
| Root (`package output`) | `./`             | enum, escape, sort          | go-faster/yaml, x/term, go-branded-id | 5 replace          | **Leaky** — sort is test-only dep listed as prod   |
| `enum/`                 | `./enum/`        | None                        | None                                  | None               | Clean                                              |
| `escape/`               | `./escape/`      | None                        | None                                  | None               | Clean                                              |
| `cmdguard/`             | `./cmdguard/`    | None (tests: root, gentest) | None                                  | None               | Clean                                              |
| `sort/`                 | `./sort/`        | root                        | None                                  | 2 replace          | Deprecated                                         |
| `table/`                | `./table/`       | root                        | lipgloss/v2                           | 3 replace          | Clean                                              |
| `integration/`          | `./integration/` | root, sort, table           | None (transitive)                     | 5 replace          | Clean                                              |
| `examples/`             | `./examples/`    | root, table                 | lipgloss (transitive)                 | 4 replace          | **Leaky** — go.mod missing table replace directive |

### 2.2 Root Module Concern Clusters

| Cluster                        | Files                                                                                               | LOC | External Deps               | Internal Deps                               |
| ------------------------------ | --------------------------------------------------------------------------------------------------- | --- | --------------------------- | ------------------------------------------- |
| **Core types + interfaces**    | format.go, ids.go, color.go, sort.go, slices.go, registry.go, format_deprecated.go                  | 681 | enum, go-branded-id, x/term | —                                           |
| **Table formatters**           | json.go, csv.go, tsv.go, markdown.go, html.go, yaml.go, xml.go, delimited.go, markup.go, marshal.go | 945 | escape, go-faster/yaml      | Core types                                  |
| **Tree formatter**             | tree.go                                                                                             | 130 | —                           | Core types                                  |
| **Graph formatters (generic)** | graph.go, dot.go, mermaid.go                                                                        | 566 | enum, escape                | Core types                                  |
| **D2 (specialized graph)**     | d2.go, d2_enum.go, d2_render.go, d2_write.go, d2_convert.go                                         | 815 | enum, escape                | Core types (GraphNode, TreeNode, BrandedID) |
| **Streaming**                  | streaming.go                                                                                        | 198 | escape                      | Core types                                  |
| **Test helpers**               | output_test_helpers.go                                                                              | 176 | —                           | All formatters                              |

**Total root production LOC: 3,587**

### 2.3 Dependency Graph (Current)

```
enum ─────────────────────────────────────┐
escape ───────────────────────────────────┤
                                         ▼
sort ──► root (package output) ◄── cmdguard (tests)
          │
          ├── json, csv, tsv, markdown, yaml, xml
          ├── html (table + tree)
          ├── tree
          ├── streaming
          ├── graph.go (GraphNode, GraphEdge, GraphRenderer)
          ├── dot.go (DOTRenderer, GraphRendererMixin)
          ├── mermaid.go (MermaidRenderer)
          └── d2*.go (D2Diagram, 5 files)
                │
table ◄──── root ────► integration
                │
            examples
```

### 2.4 Coupling Hotspots

1. **D2 ↔ Root core types** — D2's `d2_convert.go` imports `GraphNode`, `GraphEdge`, `TableData`, `TreeNode`, `GraphShape`, `GraphStyle` from root. But D2 has its own rich types (D2Node, D2Edge, D2Direction, etc.) that are fully self-contained.

2. **DOT/Mermaid ↔ Root core types** — DOT's `GraphRendererMixin`, `MermaidRenderer` both implement `GraphRenderer` and `Renderer` interfaces from root. They import `GraphNode`, `GraphEdge`, `BrandedID`, `TableData`, `TreeNode`.

3. **`internal/gentest` cross-module leak** — `enum/` and `cmdguard/` test files import `internal/gentest` from the root module. This creates a test-time dependency on the full root module.

4. **`internal/testutils` cross-module leak** — `table/` and `integration/` test files import `internal/testutils` from root.

5. **sort is test-only in root** — root's `go.mod` lists `go-output/sort` as a production dependency, but only `userjourney_test.go` imports it. Should be test-only or removed.

6. **No god-packages in sub-modules** — Each sub-module is small and focused (enum: 64 LOC, escape: 76 LOC, cmdguard: 53 LOC, table: 92 LOC). The root module IS the god-package.

### 2.5 God-Package Analysis: Root Module

The root module (`package output`) has **28 production files** spanning **6 distinct concerns**:

| Concern                         | Exports               | Cohesion                           |
| ------------------------------- | --------------------- | ---------------------------------- |
| Format enum + core types        | ~35                   | High — core abstractions           |
| Table formatters (7)            | ~60                   | High — all implement TableRenderer |
| Tree rendering                  | ~7                    | High — self-contained              |
| Graph rendering (DOT + Mermaid) | ~20                   | High — implements GraphRenderer    |
| D2 diagram rendering            | ~33 enum + ~45 render | High — rich domain model           |
| Streaming HTML                  | ~19                   | Medium — depends on HTML + core    |

All concerns are cohesive internally but coupled through `package output` — they share the same namespace, the same `go.mod`, and the same import path.

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

#### Existing Modules (Unchanged)

| Module         | Path             | Purpose                  | Production Deps   |
| -------------- | ---------------- | ------------------------ | ----------------- |
| `enum/`        | `./enum/`        | Generic enum utilities   | None              |
| `escape/`      | `./escape/`      | Format-specific escaping | None              |
| `cmdguard/`    | `./cmdguard/`    | CLI flag parsing         | None              |
| `sort/`        | `./sort/`        | Deprecated sorting       | root              |
| `table/`       | `./table/`       | Lipgloss terminal tables | root, lipgloss    |
| `integration/` | `./integration/` | Cross-module tests       | root, sort, table |
| `examples/`    | `./examples/`    | Usage examples           | root              |

#### New Modules

| Module       | Path       | Purpose                 | Production Deps                      | Public API                                                |
| ------------ | ---------- | ----------------------- | ------------------------------------ | --------------------------------------------------------- |
| **`d2/`**    | `./d2/`    | D2 diagram rendering    | root (core types only), enum, escape | D2Diagram, D2Node, D2Edge, D2Direction, D2NodeShape, etc. |
| **`graph/`** | `./graph/` | DOT + Mermaid rendering | root (core types only), enum, escape | DOTRenderer, MermaidRenderer, GraphRendererMixin          |

**NOT extracted (revised after self-review):**

- `gentest/` and `testutils/` stay as `internal/` packages in root — 300 LOC total, extraction adds complexity for minimal benefit
- Streaming stays in root — only 198 LOC, deeply coupled to HTML renderer

#### Root Module (Slimmed)

After extraction, root `package output` contains:

| Cluster                 | Files                                                                                               | LOC       |
| ----------------------- | --------------------------------------------------------------------------------------------------- | --------- |
| Core types + interfaces | format.go, ids.go, color.go, sort.go, slices.go, registry.go, format_deprecated.go                  | 734       |
| Table formatters        | json.go, csv.go, tsv.go, markdown.go, html.go, yaml.go, xml.go, delimited.go, markup.go, marshal.go | 945       |
| Tree formatter          | tree.go                                                                                             | 130       |
| Streaming               | streaming.go                                                                                        | 198       |
| Test helpers            | output_test_helpers.go                                                                              | 176       |
| **Total**               |                                                                                                     | **2,183** |

### ~~3.2 Test Helper Extraction~~ (Superseded)

> **Decision reversed during self-review (Phase 4).** `internal/gentest` and `internal/testutils` stay in root as `internal/` packages. At 151 LOC and 149 LOC respectively, extraction adds 2 more `go.mod` files and replace directives for minimal benefit. The cross-module `internal/` leak is a minor inconvenience, not a blocking issue.

### 3.3 DAG Verification

Proposed dependency direction:

```
Level 0 (zero deps):     enum, escape
Level 1 (root):          root → enum, escape, go-branded-id, x/term, go-faster/yaml
Level 2 (format mods):   d2 → root, enum, escape
                          graph → root, enum, escape
                          table → root, lipgloss
Level 3 (consumers):     integration → root, sort, table, d2, graph
                          examples → root, table, d2, graph
                          cmdguard (zero prod deps; tests: root)
                          sort → root (deprecated)
```

**Cycles:** None. All arrows point downward.

**Key invariant:** Root never imports `d2/`, `graph/`, `table/`, `sort/`, `integration/`, or `examples/`. These are always downstream consumers.

### 3.4 Replace / Workspace Strategy

**Current:** Replace directives in every consuming module's `go.mod`. No committed `go.work`.

**Proposed:** Keep current approach. It works well:

- Replace directives allow `cd d2 && go test ./...` to work standalone
- No committed `go.work` means published versions are clean
- Each module is self-contained for development

**Changes needed:**

- `d2/go.mod` — new, with replace for root, enum, escape
- `graph/go.mod` — new, with replace for root, enum, escape
- Root `go.mod` — remove `sort` from prod deps (test-only or remove entirely)
- `integration/go.mod` — add `d2`, `graph` to replace directives
- `examples/go.mod` — add `d2`, `graph`, `table` to replace directives (table is currently missing)

### 3.5 Test Dependency Isolation

| Module         | Production Deps                                     | Test-Only Deps          |
| -------------- | --------------------------------------------------- | ----------------------- |
| `enum/`        | —                                                   | —                       |
| `escape/`      | —                                                   | —                       |
| `cmdguard/`    | —                                                   | root                    |
| `root`         | enum, escape, go-branded-id, x/term, go-faster/yaml | sort (or remove), table |
| `d2/`          | root, enum, escape                                  | —                       |
| `graph/`       | root, enum, escape                                  | —                       |
| `table/`       | root, lipgloss                                      | —                       |
| `sort/`        | root                                                | —                       |
| `integration/` | root, sort, table                                   | —                       |
| `examples/`    | root, table, d2, graph                              | —                       |

### 3.6 Interface Extraction

No interface extraction needed. The current design already uses clean interfaces:

- `Renderer` — defined in root, implemented by all renderers
- `GraphRenderer` — defined in root, implemented by DOT/Mermaid/D2
- `TableRenderer` — defined in root, implemented by table formatters
- `TreeOutputRenderer` — defined in root, implemented by tree renderer

D2 and graph modules import root's interfaces and implement them. Root never imports D2 or graph. This is the correct Go pattern.

### 3.7 Versioning Strategy

**Shared versioning** (same as current):

- All modules share a single git tag `v1.2.3`
- Root module is the primary published artifact
- Sub-modules (`d2/`, `graph/`, `table/`, etc.) are opt-in extensions
- No independent semver — the modules are too tightly coupled to root's core types

**Rationale:** go-output is a single library, not a platform. Users import `go-output` and optionally `go-output/table`, `go-output/d2`, etc. Independent versioning would create compatibility headaches for no real benefit — there are no independent consumers of `enum/` or `escape/` outside this repo.

---

## 4. Migration Strategy

### 4.1 Prerequisites

- Ensure all tests pass: `go test ./...`
- Ensure clean lint: `golangci-lint run --fix ./...`
- Create branch: `git checkout -b modularize/extract-d2-graph`

### 4.2 Ordered Steps

Each step is independently committable and leaves the project buildable.

#### Step 1: Extract `d2/` module (Highest value)

- Move `d2.go`, `d2_enum.go`, `d2_render.go`, `d2_write.go`, `d2_convert.go` → `d2/`
- Rename package from `output` to `d2`
- Create `d2/go.mod` (depends on root, enum, escape)
- Move `d2_*_test.go` → `d2/`
- D2 imports root's `GraphNode`, `GraphEdge`, `TableData`, `TreeNode`, `BrandedID`, `Renderer`, `GraphRenderer`
- Format enum stays in root — D2 is a consumer of `output.FormatD2`
- Add `d2/` replace directive in integration's go.mod
- **Verification:** `go test ./d2/... && go test ./...`

#### Step 2: Extract `graph/` module

- Move `graph.go`, `dot.go`, `mermaid.go` → `graph/`
- Rename package from `output` to `graph`
- Create `graph/go.mod` (depends on root, enum, escape)
- Move `graph_test.go`, `dot_test.go`, `mermaid_test.go` → `graph/`
- Split `output_test_helpers.go`: graph-related helpers (`testNodesAB`, `testDOTEmptyExpected`, etc.) move to `graph/`
- Keep `GraphRenderer`, `GraphNode`, `GraphEdge`, `GraphShape`, `GraphStyle` in root — they're core types
- Move `GraphRendererMixin` to `graph/` — currently defined in `dot.go:21`, only used by DOT and Mermaid
- `AddTreeNodes` stays in root — used by D2, graph, and potentially external consumers
- Add `graph/` replace directive in integration's go.mod
- **Verification:** `go test ./graph/... && go test ./...`

#### Step 3: Clean up sort dependency

- Remove `sort` from root's production `go.mod` — it's only used in `userjourney_test.go`
- Add `sort` as test-only dependency or remove the test (sort is deprecated)
- **Verification:** `go mod tidy && go build ./...`

#### Step 4: Update documentation

- Update `AGENTS.md` with new module table
- Update `docs/FORMAT_ARCHITECTURE.md` if needed
- Update README if module structure section exists
- **Verification:** Manual review

---

## 5. Risk Assessment

| Risk                                         | Likelihood | Impact | Mitigation                                                                              |
| -------------------------------------------- | ---------- | ------ | --------------------------------------------------------------------------------------- |
| D2 package rename breaks consumers           | Medium     | High   | D2 module re-exports root types; public API is `d2.D2Diagram`, `d2.NewD2Diagram()` etc. |
| Graph package rename breaks consumers        | Medium     | High   | Same approach — `graph.DOTRenderer`, `graph.NewMermaidRenderer()` etc.                  |
| Integration tests fail after moves           | Low        | Medium | Step-by-step extraction with test verification at each step                             |
| Internal package moves break existing tests  | Low        | Low    | `gentest/` and `testutils/` staying in root eliminates this risk                        |
| Root module still too large after extraction | Low        | Low    | 2,183 LOC is acceptable for a core library module with 12 formats                       |
| go.mod replace directive sprawl              | Medium     | Low    | Document pattern; consider `go.work` for local dev if it grows beyond 12 modules        |

---

## 6. Key Decisions

1. **Root IS the core** (ADR 001) — unchanged. Core types, interfaces, and simple formatters stay in root.
2. **D2 gets its own module** — 833 LOC, rich domain model, self-contained rendering.
3. **DOT + Mermaid share a `graph/` module** — they share `GraphRendererMixin` and have identical dependency profiles.
4. **`GraphNode`/`GraphEdge` stay in root** — D2 also uses them for `SetNodes`/`SetEdges`. Moving them to `graph/` would create a D2→graph dependency or require duplication.
5. **`GraphRendererMixin` moves to `graph/`** — only used by DOT and Mermaid, not by D2.
6. **Test helpers stay in root** — `internal/gentest` and `internal/testutils` remain. The cross-module leak is minor and extraction adds more complexity than value.
7. **`output_test_helpers.go` split** — graph-related test helpers move with graph module; HTML/tree helpers stay in root.
8. **`sort/` dependency cleaned from root** — test-only import should not be in production `go.mod`.
9. **Package rename required** — `d2/` becomes `package d2`, `graph/` becomes `package graph`. This is a breaking change for any code using `output.D2Diagram` etc. Acceptable for pre-v1 library.

---

## 7. Build System Impact

- **flake.nix** — Update to build/test each new module independently
- **CI/CD** — Add parallel jobs for `d2/`, `graph/`
- **golangci-lint** — Already runs workspace-wide, no changes needed
- **go.work** — Add new modules to workspace file (if re-created for local dev)

---

## 8. What Does NOT Change

- `enum/`, `escape/`, `cmdguard/` — unchanged
- `table/` — unchanged
- `sort/` — unchanged (still deprecated)
- `internal/gentest/`, `internal/testutils/` — unchanged (stays in root after self-review)
- Root module name — still `github.com/larsartmann/go-output`, still `package output`
- Format enum — still in root
- Renderer/TableRenderer/GraphRenderer interfaces — still in root
- Core types (Format, TableData, TreeNode, GraphNode, GraphEdge) — still in root
- All 7 table formatters — still in root
- Streaming — still in root
- Registry — still in root
