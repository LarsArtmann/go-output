# Execution Plan — go-output Modularization

**Date:** 2026-05-16
**Branch:** `modularize/extract-d2-graph`

---

## Overview

Extract `d2/` and `graph/` as independent Go modules from the root `package output`. Five ordered steps, each independently committable and buildable.

---

## Step 0: Fix Pre-Existing `examples/go.mod` Bug

**Impact:** Prerequisite — `examples/go.mod` is missing the `table` dependency that `examples/basic/main.go` imports.

**Effort:** 1 minute

### Actions

1. **Fix `examples/go.mod`**
   ```bash
   cd examples && go mod tidy
   ```
   This adds the missing `github.com/larsartmann/go-output/table` require and replace directives.

### Verification

- [ ] `cd examples && go build ./...` passes
- [ ] `cd examples && go test ./...` passes (if examples have tests)

---

## Step 1: Extract D2 Module

**Impact:** 1% → 51% — Removes the largest single concern (833 LOC) from root.

**Effort:** 20–30 minutes

### Actions

1. **Create `d2/` directory**

   ```bash
   mkdir d2
   ```

2. **Move D2 production files** (preserve git history)

   ```bash
   git mv d2.go d2_enum.go d2_render.go d2_write.go d2_convert.go d2/
   ```

3. **Move D2 test files**

   ```bash

   ```

4. **Create `d2/go.mod`**

   ```
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

5. **Rename package** in all moved files: `package output` → `package d2`

6. **Update imports in moved files:**
   - Root types need `output` import: `output.GraphNode`, `output.GraphEdge`, `output.TableData`, `output.TreeNode`, `output.NewBrandedID`, `output.GraphRenderer`, `output.Renderer`, `output.AddTreeNodes`, `output.NodesFromTableData`, `output.DefaultGraphNodeLabel`, `output.GraphShape`, `output.GraphStyle`, `output.GraphNodeIDBrand`, `output.GraphNodeLabelBrand`, `output.D2NodeIDBrand`, `output.D2NodeLabelBrand`
   - `enum` stays as direct import
   - `escape` stays as direct import

7. **Update examples**
   - `examples/basic/main.go`: Change `output.NewD2Diagram()` → `d2.NewD2Diagram()`, update D2 type references
   - `examples/d2/main.go`: Update all `output.D2*` references to `d2.*`
   - `examples/shared/shared.go`: Change `output.NewD2Diagram()` → `d2.NewD2Diagram()`, update return type to `*d2.D2Diagram`
   - Add `d2` to `examples/go.mod` require and replace blocks

8. **Update integration tests**
   - `integration/d2_test.go`: Change `output.NewD2Diagram()` → `d2.NewD2Diagram()`
   - Update `integration/go.mod`: add `d2` to require and replace blocks
   - Add `"github.com/larsartmann/go-output/d2"` import to affected files

9. **Update `benchmarks_test.go`**
   - No D2 benchmarks currently — skip for this step

10. **Run in d2 directory:**

    ```bash
    cd d2 && go mod tidy && go build ./... && go test ./...
    ```

11. **Run workspace-wide:**
    ```bash
    go build ./... && go test ./...
    ```

### Verification

- [ ] `cd d2 && go build ./...` passes
- [ ] `cd d2 && go test ./...` passes
- [ ] `cd d2 && go vet ./...` reports no issues
- [ ] `cd d2 && go mod tidy` changes nothing
- [ ] Root-level `go build ./...` passes
- [ ] Root-level `go test ./...` passes
- [ ] No production dependency on d2 from root module

### Rollback

```bash
git revert HEAD
```

### What stays in root after this step

- D2 format constant (`FormatD2`) stays in `format.go`
- D2 is in `tableFormats` and `graphFormats` maps in `format.go`
- D2 branded ID brands (`D2NodeIDBrand`, `D2NodeLabelBrand`) stay in `ids.go` — they're type definitions used by d2 module through root import
- `AddTreeNodes`, `NodesFromTableData` stay in root (used by graph module too)

---

## Step 2: Extract Graph Module

**Impact:** 4% → 64% — Removes graph rendering (DOT + Mermaid, 568 LOC) from root.

**Effort:** 20–30 minutes

### Actions

1. **Create `graph/` directory**

   ```bash
   mkdir graph
   ```

2. **Move graph production files**

   ```bash
   git mv dot.go mermaid.go graph/
   ```

   Note: `graph.go` stays entirely in root — it contains only core types (`GraphNode`, `GraphEdge`, `GraphRenderer` interface) and utility functions. `GraphRendererMixin` is in `dot.go`, so it moves with `dot.go` to `graph/`.

3. **Extract `GraphRendererMixin` from `dot.go` (optional)**
   - `GraphRendererMixin` is currently at `dot.go:21` and moves with `dot.go` to `graph/`
   - Optionally extract into `graph/mixin.go` for organization — not required, purely cosmetic
   - No changes to `graph.go` needed — it stays in root as-is

4. **Move graph test files**

   ```bash
   git mv dot_test.go mermaid_test.go graph_test.go graph/
   ```

5. **Create `graph/go.mod`**

   ```
   module github.com/larsartmann/go-output/graph

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

6. **Rename package** in moved files: `package output` → `package graph`

7. **Update imports in moved files:**
   - Root types: `output.GraphNode`, `output.GraphEdge`, `output.TableData`, `output.TreeNode`, `output.NewBrandedID`, `output.GraphRenderer`, `output.Renderer`, `output.GraphShape`, `output.GraphStyle`, `output.GraphNodeIDBrand`, `output.GraphNodeLabelBrand`, `output.TreeNodeID`, `output.TreeNodeLabelBrand`
   - DOT convenience constructors that move: `NewDOTRenderer`, `NewUndirectedDOTRenderer`, `DOTFromTableData`, `DOTFromTree`
   - Mermaid convenience constructors that move: `NewMermaidRenderer`, `MermaidFlowchartRenderer`, `MermaidTreeRenderer`

8. **Split `output_test_helpers.go`:**
   - **Move to `graph/`:** `testDOTEmptyExpected`, `testMermaidEmptyExpected`, `testNodesAB`, `testNodesABC`, `newTestNode`, `newTestNodeWithShape`, `testEdgeAB`, `testEdgesAB`, `testEdgesABC`
   - **Stay in root:** `testHTMLEscapeShared`, `testEmptyRendererOutput`, `testSanitizeFunc`, `AssertTreeNodeDepth`, `testExpectedOutputs`, `testHTMLEmptyExpected`, gentest re-exports

9. **Update `benchmarks_test.go`**
   - Change `NewMermaidRenderer()` → `graph.NewMermaidRenderer()`
   - Change `NewDOTRenderer()` → `graph.NewDOTRenderer()`
   - Add `\"github.com/larsartmann/go-output/graph\"` import

10. **Update `examples/basic/main.go`**
    - Change `output.MermaidFlowchartRenderer()` → `graph.MermaidFlowchartRenderer()`
    - Change `output.DOTFromTableData()` → `graph.DOTFromTableData()`
    - Update `examples/go.mod` with `graph` (and `d2`) replace directives

11. **Update integration tests**
    - `integration/renderer_test.go`: Change `output.MermaidFlowchartRenderer` → `graph.MermaidFlowchartRenderer`, `output.DOTFromTableData` → `graph.DOTFromTableData`
    - `integration/integration_test.go`: Same updates for DOT/Mermaid calls
    - Update `integration/go.mod`: add `graph` to require and replace blocks
    - Add `\"github.com/larsartmann/go-output/graph\"` import to affected files

12. **Run in graph directory:**

    ```bash
    cd graph && go mod tidy && go build ./... && go test ./...
    ```

13. **Run workspace-wide:**
    ```bash
    go build ./... && go test ./...
    ```

### Verification

- [ ] `cd graph && go build ./...` passes
- [ ] `cd graph && go test ./...` passes
- [ ] `cd graph && go vet ./...` reports no issues
- [ ] `cd graph && go mod tidy` changes nothing
- [ ] Root-level `go build ./...` passes
- [ ] Root-level `go test ./...` passes
- [ ] No production dependency on graph from root module
- [ ] D2 module still builds and passes tests (it uses root's GraphNode etc., not graph/)

### Rollback

```bash
git revert HEAD
```

---

## Step 3: Clean Up Sort Dependency

**Impact:** 20% → 80% — Fixes a real dep leak in root go.mod.

**Effort:** 5 minutes

### Actions

1. **`userjourney_test.go` stays in root** (not moved in Step 1 — it tests JSON/CSV/Markdown/YAML/sort, not D2)
   - Its `sort` import is test-only

2. **Remove sort from root's production go.mod**
   - Currently listed as: `github.com/larsartmann/go-output/sort v0.0.0-20260507215750-c2091663ee59`
   - Remove from `require` block
   - Keep in `replace` block (needed by integration tests transitively)
   - If needed as test-only dep, it will be auto-added by `go mod tidy` with `// indirect`

3. **Run `go mod tidy` in root**
   ```bash
   go mod tidy
   ```

### Verification

- [ ] Root `go.mod` no longer lists `sort` as direct production dependency
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] `go mod tidy` changes nothing (already clean)

### Rollback

```bash
git revert HEAD
```

---

## Step 4: Update Documentation

**Impact:** Polish — ensures docs reflect the new structure.

**Effort:** 10 minutes

### Actions

1. **Update `AGENTS.md`**
   - Add `d2/` and `graph/` to module table
   - Update root module description (slimmed)
   - Update line counts
   - Update "Current Coverage" table if applicable

2. **Update `docs/FORMAT_ARCHITECTURE.md`**
   - Note that D2 is now in `d2/` module
   - Note that DOT/Mermaid are now in `graph/` module
   - Core types remain in root

3. **Update `README.md`**
   - Change `output.NewD2Diagram()` → `d2.NewD2Diagram()` (and all `output.D2*` types)
   - Change `output.DOTFromTableData()` → `graph.DOTFromTableData()`
   - Change `output.MermaidFlowchartRenderer()` → `graph.MermaidFlowchartRenderer()`
   - Add `d2` and `graph` import examples

### Verification

- [ ] AGENTS.md reflects 10-module workspace (root + enum + escape + cmdguard + sort + table + integration + examples + **d2** + **graph**)
- [ ] Build/test commands still accurate
- [ ] No references to old `output.D2Diagram` etc. without noting the module change

### Rollback

```bash
git revert HEAD
```

---

## Dependency Changes Summary

### New `d2/go.mod`

- `require`: root, enum, escape
- `replace`: root → `../`, enum → `../enum`, escape → `../escape`

### New `graph/go.mod`

- `require`: root, enum, escape
- `replace`: root → `../`, enum → `../enum`, escape → `../escape`

### Updated `root/go.mod`

- Remove: d2.go, d2_enum.go, d2_render.go, d2_write.go, d2_convert.go, dot.go, mermaid.go from root package
- Remove: `sort` from production requires (test-only)
- Keep: all current `replace` directives

### Updated `integration/go.mod`

- Add: `d2` and `graph` to `require` and `replace` blocks

### Updated `examples/go.mod`

- Add: `d2` and/or `graph` to `require` and `replace` blocks (if examples use them)

---

## Final Module Count

| #   | Module                  | Status                           |
| --- | ----------------------- | -------------------------------- |
| 1   | `root` (package output) | Existing, slimmed                |
| 2   | `enum/`                 | Existing, unchanged              |
| 3   | `escape/`               | Existing, unchanged              |
| 4   | `cmdguard/`             | Existing, unchanged              |
| 5   | `sort/`                 | Existing, unchanged (deprecated) |
| 6   | `table/`                | Existing, unchanged              |
| 7   | `integration/`          | Existing, updated deps           |
| 8   | `examples/`             | Existing, updated deps           |
| 9   | `d2/`                   | **NEW**                          |
| 10  | `graph/`                | **NEW**                          |

---

## Success Criteria

- [ ] All 10 modules build independently: `go build ./...` in each directory
- [ ] All tests pass workspace-wide: `go test ./...` from root
- [ ] All modules pass `go vet ./...`
- [ ] All modules pass `go mod tidy` with no changes
- [ ] Root module has zero production imports from d2/ or graph/
- [ ] No circular dependencies between any module pair
- [ ] `golangci-lint run ./...` passes
