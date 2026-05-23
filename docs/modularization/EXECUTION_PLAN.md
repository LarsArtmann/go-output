# Execution Plan — go-output Modularization (v2)

**Date:** 2026-05-23
**Branch:** `modularize/extract-d2-graph`
**Status:** COMPLETED
**Supersedes:** 2026-05-16 execution plan

---

## Overview

Extract `d2/` and `graph/` as independent Go modules from root `package output`. Five ordered steps, each independently committable and buildable.

---

## Step 0: Create Branch & Baseline

**Effort:** 1 minute

### Actions

1. Create feature branch: `git checkout -b modularize/extract-d2-graph`
2. Verify clean state: `git status`

### Verification

- [ ] Branch created
- [ ] Working tree clean

---

## Step 1: Extract D2 Module

**Impact:** 1% → 51% — Removes D2 diagram rendering (835 LOC, 5 production + 6 test files) from root.

**Effort:** 25–35 minutes

### Actions

1. **Create `d2/` directory**

2. **Move D2 production files** (preserve git history):

   ```bash
   git mv d2.go d2_enum.go d2_render.go d2_write.go d2_convert.go d2/
   ```

3. **Move D2 test files**:

   ```bash
   git mv d2_test.go d2_node_test.go d2_enum_test.go d2_render_test.go d2_convert_test.go d2_edge_test.go d2/
   ```

   Note: `userjourney_test.go` stays in root — it tests JSON/CSV/Markdown/YAML user journeys.

4. **Create `d2/go.mod`**:

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

6. **Update imports in moved files** — add `"github.com/larsartmann/go-output"` import and prefix root types:
   - `output.GraphNode`, `output.GraphEdge`, `output.TableData`, `output.TreeNode`
   - `output.NewBrandedID`, `output.GraphRenderer`, `output.Renderer`
   - `output.AddTreeNodes`, `output.NodesFromTableData`, `output.DefaultGraphNodeLabel`
   - `output.GraphShape`, `output.GraphStyle`, `output.GraphNodeIDBrand`, `output.GraphNodeLabelBrand`
   - `output.D2NodeIDBrand`, `output.D2NodeLabelBrand`
   - `enum` and `escape` stay as direct imports (no prefix change)

7. **Remove `renderD2TableData` from `render_tabledata.go`**:
   - Remove the `renderD2TableData` function
   - Change the `FormatD2` case in `RenderTableData` switch to fall through to `default` (returns `UnsupportedFormatError`)
   - Or explicitly: `case FormatD2: return &UnsupportedFormatError{Format: format}`

8. **Move D2 test from `render_tabledata_test.go`** to `d2/` (the `TestRenderTableDataD2` test)

9. **Update test files in `d2/`**:
   - Add `"github.com/larsartmann/go-output"` import
   - Prefix root types with `output.`
   - Prefix test helper calls: `testhelpers.AssertEqual` stays, `testNodesAB` etc. become inline or defined locally

10. **Update `examples/go.mod`** and **`examples/` code**:
    - `examples/basic/main.go`: `output.D2Diagram` → `d2.D2Diagram`, `output.NewD2Diagram` → `d2.NewD2Diagram`, etc.
    - `examples/shared/shared.go`: Same updates
    - Add `d2` to `examples/go.mod` require and replace blocks

11. **Update `integration/go.mod`** and **`integration/` code**:
    - Add `d2` to require and replace blocks
    - Update integration test files that reference D2 types

12. **Run `go mod tidy` in d2, root, examples, integration**

### Verification

- [ ] `cd d2 && go build ./...` passes
- [ ] `cd d2 && go test ./...` passes
- [ ] `cd d2 && go vet ./...` clean
- [ ] `cd d2 && go mod tidy` changes nothing
- [ ] Root `go build ./...` passes
- [ ] Root `go test ./...` passes
- [ ] Root does not import `d2/` in production code
- [ ] `cd examples && go build ./...` passes
- [ ] `cd integration && go test ./...` passes

### Rollback

```bash
git revert HEAD
```

---

## Step 2: Extract Graph Module

**Impact:** 4% → 64% — Removes DOT + Mermaid rendering (319 LOC, 2 production + 3 test files) from root.

**Effort:** 25–35 minutes

### Actions

1. **Create `graph/` directory**

2. **Move graph production files**:

   ```bash
   git mv dot.go mermaid.go graph/
   ```

   Note: `graph.go` stays entirely in root — core types (`GraphNode`, `GraphEdge`, `GraphRenderer`, `GraphRendererMixin`) used by D2, JSON, YAML too.

3. **Move graph test files**:

   ```bash
   git mv dot_test.go mermaid_test.go graph_test.go graph/
   ```

4. **Create `graph/go.mod`**:

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

5. **Rename package** in moved files: `package output` → `package graph`

6. **Update imports in moved files** — add `"github.com/larsartmann/go-output"` and prefix:
   - `output.GraphNode`, `output.GraphEdge`, `output.TableData`, `output.TreeNode`
   - `output.NewBrandedID`, `output.GraphRenderer`, `output.Renderer`
   - `output.GraphRendererMixin`, `output.NewGraphRendererMixin`
   - `output.GraphShape`, `output.GraphStyle`
   - `output.GraphNodeIDBrand`, `output.GraphNodeLabelBrand`
   - `output.TreeNodeID`, `output.TreeNodeLabelBrand`
   - `output.AddTreeNodes`, `output.NodesFromTableData`, `output.DefaultGraphNodeLabel`
   - `enum` and `escape` stay as direct imports

7. **Split `output_test_helpers.go`**:
   - **Move to `graph/`**: `testDOTEmptyExpected`, `testMermaidEmptyExpected`, `testNodesAB`, `testNodesABC`, `newTestNode`, `newTestNodeWithShape`, `testEdgeAB`, `testEdgesAB`, `testEdgesABC`
   - **Stay in root**: `testHTMLEscapeShared`, `testEmptyRendererOutput`, `testSanitizeFunc`, `AssertTreeNodeDepth`, `testExpectedOutputs`, `testHTMLEmptyExpected`, gentest/testhelpers re-exports

8. **Remove graph render functions from `render_tabledata.go`**:
   - Remove `renderMermaidTableData` and `renderDOTTableData`
   - Add explicit unsupported cases: `case FormatMermaid:`, `case FormatDOT:` → return `UnsupportedFormatError`
   - Remove `RenderOptions.GraphID` if only used by DOT (check first)

9. **Move Mermaid/DOT test cases from `render_tabledata_test.go`** to `graph/`

10. **Update test files in `graph/`**:
    - Add `"github.com/larsartmann/go-output"` import
    - Prefix root types with `output.`
    - Ensure moved test helpers work in `package graph` context

11. **Update `benchmarks_test.go`**:
    - `NewMermaidRenderer()` → `graph.NewMermaidRenderer()` (requires import)
    - `NewDOTRenderer()` → `graph.NewDOTRenderer()` (requires import)
    - Add `"github.com/larsartmann/go-output/graph"` import

12. **Update `examples/` code**:
    - `output.MermaidFlowchartRenderer()` → `graph.MermaidFlowchartRenderer()`
    - `output.DOTFromTableData()` → `graph.DOTFromTableData()`
    - Add `graph` to `examples/go.mod`

13. **Update `integration/` code**:
    - Add `graph` to require and replace blocks
    - Update integration test files that reference DOT/Mermaid types

14. **Run `go mod tidy` in graph, root, examples, integration**

### Verification

- [ ] `cd graph && go build ./...` passes
- [ ] `cd graph && go test ./...` passes
- [ ] `cd graph && go vet ./...` clean
- [ ] `cd graph && go mod tidy` changes nothing
- [ ] Root `go build ./...` passes
- [ ] Root `go test ./...` passes
- [ ] Root does not import `graph/` in production code
- [ ] D2 module still builds: `cd d2 && go test ./...`
- [ ] `cd examples && go build ./...` passes
- [ ] `cd integration && go test ./...` passes

### Rollback

```bash
git revert HEAD
```

---

## Step 3: Clean Up Sort Dependency

**Impact:** 20% → 80% — Fixes dep leak in root go.mod.

**Effort:** 5 minutes

### Actions

1. Check if `sort` is still in root's `go.mod` as a production dependency
2. If yes: remove from `require` block, keep in `replace` block if needed by integration
3. Run `go mod tidy` in root

### Verification

- [ ] Root `go.mod` does not list `sort` as direct production dependency
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

### Rollback

```bash
git revert HEAD
```

---

## Step 4: Update Documentation

**Impact:** Polish — ensures docs reflect new structure.

**Effort:** 10 minutes

### Actions

1. **Update `AGENTS.md`**:
   - Add `d2/` and `graph/` to module table
   - Update root module description (slimmed)
   - Update line counts
   - Update dependency graph

2. **Update `docs/modularization/DEPENDENCY_GRAPH.md`**:
   - Show proposed graph with d2/ and graph/

3. **Update `README.md`** (if it references `output.D2Diagram` etc.)

### Verification

- [ ] AGENTS.md reflects 10-module workspace
- [ ] No stale references to `cmdguard/`
- [ ] No stale references to `internal/testutils/`

### Rollback

```bash
git revert HEAD
```

---

## Step 5: Final Verification

**Effort:** 5 minutes

### Actions

1. Run full test suite across all modules
2. Run `go vet ./...` across all modules
3. Run `go mod tidy` in all modules — verify no changes
4. Verify DAG: root has zero imports from d2/ or graph/

```bash
for mod in . d2 graph enum escape testhelpers sort table integration examples; do
  echo "=== $mod ==="
  (cd "$mod" && go build ./... && go test ./... && go vet ./... && go mod tidy)
done
```

### Success Criteria

- [ ] All 10 modules build independently
- [ ] All tests pass workspace-wide
- [ ] All modules pass `go vet`
- [ ] All modules pass `go mod tidy` with no changes
- [ ] Root has zero production imports from d2/ or graph/
- [ ] No circular dependencies

---

## Final Module Count

| #   | Module                  | Status                           |
| --- | ----------------------- | -------------------------------- |
| 1   | `root` (package output) | Existing, slimmed                |
| 2   | `enum/`                 | Existing, unchanged              |
| 3   | `escape/`               | Existing, unchanged              |
| 4   | `testhelpers/`          | Existing, unchanged              |
| 5   | `sort/`                 | Existing, unchanged (deprecated) |
| 6   | `table/`                | Existing, unchanged              |
| 7   | `integration/`          | Existing, updated deps           |
| 8   | `examples/`             | Existing, updated deps           |
| 9   | `d2/`                   | **NEW**                          |
| 10  | `graph/`                | **NEW**                          |
