# Status Report — CQRS Architecture Completion (Session 2)

**Date:** 2026-07-06 06:38
**Branch:** master (NOT committed — 22 files uncommitted)
**Base commit:** `c98ef0a` — docs(status): v0.30.0 breaking changes execution status report

---

## A. FULLY DONE ✓

### TableBuilder + TreeBuilder (root package)

- [x] `table_builder.go` — `NewTableBuilder()`, `SetHeaders()`, `AddRow()`, `SetFooter()`, `Build()→*Table` (fluent API, copies slices)
- [x] `tree_builder.go` — `NewTreeBuilder()`, `SetRoot()`, `AddChild()`, `Build()→*TreeNode` (fluent API, unknown-parent guard)
- [x] Full test coverage: `table_builder_test.go` (3 tests), `tree_builder_test.go` (3 tests)
- [x] Lint clean, build clean

### Cross-Shape Projections (root package)

- [x] `projections.go` — `TableToGraph(t)→Graph`, `GraphToTree(g)→*TreeNode`, `GraphToTable(g)→*Table`
- [x] `projections_test.go` — 8 tests covering linear chains, branching, cycles, empty inputs, nil guards
- [x] Cycle guard in `GraphToTree` (visited set prevents infinite recursion)

### CQRS Writer/Render API (delimited/, serialization/, markup/)

- [x] `delimited/cqrs.go` — `WriteCSV`/`RenderCSV`/`WriteTSV`/`RenderTSV` (4 functions)
- [x] `serialization/cqrs.go` — `WriteJSON`/`RenderJSON`/`WriteYAML`/`RenderYAML`/`WriteTOML`/`RenderTOML`/`WriteJSONL`/`RenderJSONL` (8 functions)
- [x] `markup/cqrs.go` — `WriteXML`/`RenderXML`/`WriteAsciiDoc`/`RenderAsciiDoc`/`WriteHTML`/`RenderHTML` (6 functions)
- [x] Tests for all 18 functions across 3 `cqrs_test.go` files
- [x] Wrapcheck compliance fixed (errors wrapped with context)

### Godoc Examples

- [x] `cqrs_example_test.go` in root (5 examples: TableBuilder, TreeBuilder, GraphBuilder, TableToGraph, GraphToTree)
- [x] `graph/cqrs_example_test.go` (RenderDOT, RenderMermaid)
- [x] `d2/cqrs_example_test.go` (RenderGraph)
- [x] `plantuml/cqrs_example_test.go` (Render)
- [x] `tree/cqrs_example_test.go` (RenderASCII)
- [x] `delimited/cqrs_example_test.go` (RenderCSV)
- [x] `serialization/cqrs_example_test.go` (RenderJSON)

### Documentation

- [x] `README.md` — CQRS quick-start section with Build→Freeze→Render flow, cross-shape projection examples, full migration guide with type rename tables
- [x] `CHANGELOG.md` — New "Added — CQRS Architecture" section listing all builders, projections, and pure-function APIs
- [x] `AGENTS.md` — Three new pattern entries: CQRS builders, CQRS pure-function renderers, cross-shape projections

### Verification — ALL GREEN

- [x] Build: 19/19 modules ✓
- [x] Tests: 19/19 modules ✓
- [x] Lint: 0 issues ✓
- [x] Race tests: nom + tui race-free ✓
- [x] Govulncheck: 0 vulnerabilities ✓

---

## B. PARTIALLY DONE ⚠️

### The "Freeze" Is Only Half-Implemented

The CQRS vision (Appendix A) calls for **immutable value types** at the freeze boundary. Only `Graph` achieves this:

| Type        | Truly Immutable? | Problem                                                                         |
| ----------- | :--------------: | ------------------------------------------------------------------------------- |
| `Graph`     |        ✓         | Unexported fields, only `Nodes()`/`Edges()` read accessors                      |
| `Table`     |      **NO**      | Exported fields `Headers`, `Rows`, `Footer` — anyone can mutate after `Build()` |
| `*TreeNode` |      **NO**      | Exported `Children` slice, exported `AddChild()` method — tree is fully mutable |

`TableBuilder.Build()` claims to "freeze" but returns `*Table` with all fields exported. The slice copy is a snapshot, but the caller can still append rows or change headers. `TreeBuilder.Build()` returns a `*TreeNode` that has `AddChild()` — the "frozen" tree can grow new branches.

**Why this happened:** `Table` and `TreeNode` predate CQRS. They are the primary data types used by every renderer. Making their fields unexported would break every consumer and every test. The "freeze" is aspirational, not enforced.

### serialization/ CQRS Functions Are NOT Streaming

`WriteJSON(w, data)` calls `renderJSONTable(w, data, opts)` which internally calls `renderViaRenderer()` → `renderer.Render()` (builds full string) → `fmt.Fprint(w, out)`. This is **string-then-write**, not true streaming. The CQRS vision says "Writer-primary gives streaming for free" but the implementation just wraps the old string-building approach behind a `WriteXxx` signature.

The delimited/ and markup/ modules have the same issue — they call `MarshalXxxFromTable()` which returns `[]byte`, then write it. No incremental streaming.

### Registry Dispatch Not Unified

The CQRS vision shows `output.Render(t, output.FormatJSON)` as a clean unified dispatch. The current `RenderTable(data, fmt, RenderOptions{...})` still uses the rigid options struct. No functional-options variant exists at the root level.

---

## C. NOT STARTED ❌

### Missing CQRS Functions

- [ ] `table.Render(t, opts...)` / `table.Write(w, t, opts...)` — lipgloss table module has no CQRS pure functions
- [ ] `tree.RenderMarkdown(root)` / `tree.WriteMarkdown(w, root)` — only `RenderASCII` was added
- [ ] `daghtml.Render(g)` / `daghtml.Write(w, g)` — not in scope but mentioned in decision record

### Missing Tests

- [ ] **Integration tests** — No cross-module test proving `TableToGraph(tbl) → graph.RenderDOT(g)` works end-to-end through the public API
- [ ] **Benchmarks for new CQRS functions** — Only old API has benchmarks (graph/, table/)
- [ ] **Projection edge cases** — `GraphToTree` with disconnected nodes, `TableToGraph` with custom label functions

### Missing Operations

- [ ] **D4: Tag testhelpers/ v0.30.0** — not done
- [ ] **Commit the work** — 22 files uncommitted (19 new, 3 modified)
- [ ] **Tag v0.30.0** — blocked on commit + user decision on old renderer structs

---

## D. TOTALLY FUCKED UP 💥

### None — No regressions, no broken state, no data loss.

But there IS a **naming lie** in the API:

`TableBuilder.Build()` doc says "freezes the builder state into a *Table snapshot" — but `*Table` is fully mutable. The word "freeze" is a lie. The doc should say "copies" or the type should actually be immutable. This is the kind of thing that erodes trust in an API.

### Near-misses:

1. **Wrapcheck violations** — Initial `markup/cqrs.go` and `serialization/cqrs.go` returned raw `w.Write()` errors without wrapping. Fixed after lint caught it. Should have wrapped from the start — the project's wrapcheck rule is documented.
2. **Stale gopls diagnostics** — `NodeShapeRect` deprecated warning on `status_registry.go:72` is stale (line shows `NodeShapeHexagon`). The `physicalLines` unusedparam on `inline_renderer.go:500` is also stale (used on line 544). Did not investigate why gopls cache is stale — just verified via lint that the actual code is correct.

---

## E. WHAT WE SHOULD IMPROVE 🎯

### 1. The Freeze Boundary Is a Lie for Table and Tree

`Graph` is properly immutable (unexported fields, read accessors only). `Table` and `TreeNode` are NOT — they have exported mutable fields. Options:

- **Option A (honest docs):** Rename the concept. Stop calling it "freeze." `TableBuilder.Build()` returns a snapshot copy, not an immutable value. Document this honestly.
- **Option B (real immutability):** Create `FrozenTable` and `FrozenTree` types with unexported fields. `Build()` returns these. Renderers accept `FrozenTable`/`FrozenTree`. This is the full CQRS vision but requires updating every renderer signature.
- **Option C (Go convention):** Accept that Go doesn't have real immutability. Document that `Build()` returns a defensive copy and callers SHOULD treat it as read-only (but nothing enforces this). This is what the current implementation does — it just doesn't say so.

### 2. The "Write" Functions Don't Actually Stream

`serialization.WriteJSON(w, data)` builds the entire JSON string in memory, then writes it. For a 100K-row table, this allocates the full string before writing a single byte. True streaming would use `json.NewEncoder(w).Encode()` directly. The delimited writers (`CSVWriter`, `TSVWriter`) already have proper streaming — the CQRS `WriteCSV` should use them directly, not go through `MarshalCSVFromTable`.

### 3. TableToGraph Lacks Customization

The projection hardcodes `DefaultGraphNodeLabel` and `CreateRowEdges`. A real projection API needs options:

```go
g := output.TableToGraph(tbl,
    output.WithLabelFunc(func(header, cell string) string { return cell }),
    output.WithEdgeStrategy(output.EdgeStrategyNone),
)
```

### 4. CHANGELOG Has Stale Self-Referential Entries

Lines 83-93 of CHANGELOG.md contain entries like "NOMSubscriber → use NOMSubscriber" (same name → same name). These are artifacts of the earlier rename session where the type was renamed and the deprecated alias had the same name as the new type. Should be cleaned up or removed since the aliases are deleted.

### 5. No CQRS for table/ Module

The lipgloss table module is the most user-facing renderer and it has no CQRS pure functions. `table.Render(t)` should exist alongside the existing `table.New(table.WithColorMode(...))` constructor.

---

## F. TOP 25 THINGS TO DO NEXT 📋

| #   | Task                                                                                | Impact   | Effort | Deps |
| --- | ----------------------------------------------------------------------------------- | -------- | ------ | ---- |
| 1   | **Commit the 22 uncommitted files**                                                 | Critical | 2m     | —    |
| 2   | Fix `TableBuilder.Build()` doc: stop saying "freeze" for mutable `*Table`           | High     | 5m     | —    |
| 3   | Make `serialization.WriteJSON` actually stream (use `json.NewEncoder(w)`)           | High     | 15m    | —    |
| 4   | Make `serialization.WriteYAML` actually stream (use `yaml.NewEncoder(w)`)           | High     | 10m    | —    |
| 5   | Make `delimited.WriteCSV` use `CSVWriter` directly instead of `MarshalCSVFromTable` | High     | 10m    | —    |
| 6   | Add `table.Render(t, opts...)` / `table.Write(w, t)` CQRS to table/                 | High     | 10m    | —    |
| 7   | Add `tree.RenderMarkdown(root)` / `tree.WriteMarkdown(w, root)`                     | Med      | 10m    | —    |
| 8   | Add integration test: `TableToGraph → graph.RenderDOT` end-to-end                   | High     | 10m    | 1    |
| 9   | Add `TableToGraph` options (custom label func, edge strategy)                       | Med      | 15m    | —    |
| 10  | Add benchmarks for new CQRS functions (RenderDOT, RenderCSV, etc.)                  | Med      | 15m    | —    |
| 11  | Clean up CHANGELOG stale self-referential "Deprecated" entries                      | Low      | 5m     | —    |
| 12  | Consider `FrozenTable` / `FrozenTree` types for real immutability                   | High     | 30m    | —    |
| 13  | Unify root `RenderTable` dispatch to functional options                             | Med      | 15m    | —    |
| 14  | Tag testhelpers/ v0.30.0 (D4)                                                       | Low      | 3m     | 1    |
| 15  | Add `GraphToTree` test for disconnected nodes (forest)                              | Low      | 5m     | —    |
| 16  | Add `GraphToTable` test for styled nodes (style preservation)                       | Low      | 5m     | —    |
| 17  | Add `daghtml.Render(g)` CQRS function                                               | Low      | 10m    | —    |
| 18  | Delete old renderer structs (DOTRenderer etc.) — v0.31.0 decision                   | High     | 30m    | —    |
| 19  | Make `graph.RenderDOT` read from `Graph` without intermediate DOTRenderer           | High     | 20m    | —    |
| 20  | Add `TableBuilder.AddRows(rows [][]string)` bulk method                             | Low      | 5m     | —    |
| 21  | Add `TreeBuilder.AddChildren(parentID, children...)` bulk method                    | Low      | 5m     | —    |
| 22  | Investigate stale gopls cache (NodeShapeRect, physicalLines warnings)               | Low      | 10m    | —    |
| 23  | Add CQRS example to `examples/` directory                                           | Med      | 10m    | 1    |
| 24  | Wire CQRS functions into registry: `output.RegisterFormat(fmt, RenderFunc)`         | Med      | 20m    | —    |
| 25  | Tag v0.30.0 root module                                                             | Critical | 2m     | 1-12 |

---

## G. TOP #1 QUESTION 🤔

**Should `TableBuilder.Build()` return an immutable `FrozenTable` (unexported fields, read accessors) or keep returning mutable `*Table`?**

This is the core tension:

- **`FrozenTable`** (real CQRS): Clean, honest, safe for concurrent access, cacheable. But requires updating every renderer that currently reads `data.Headers` / `data.Rows` directly. ~50+ call sites.
- **`*Table`** (current): Backward compatible, zero churn. But the "freeze" is a doc lie — nothing prevents mutation. A consumer could `Build()` then `AddRow()` on the result.

The same question applies to `TreeBuilder.Build() → *TreeNode` vs `FrozenTree`.

**I cannot resolve this myself** because it's a fundamental API stability decision. `FrozenTable` means every consumer's `data.Headers` access becomes `data.Headers()` (method call). That's a v1.0.0-level break, not a v0.30.0 polish.

---

## Verification Summary

| Check                   | Result                      |
| ----------------------- | --------------------------- |
| `nix run .#build`       | 19/19 ✓                     |
| `nix run .#test`        | 19/19 ✓                     |
| `nix run .#lint`        | 0 issues ✓                  |
| `nix run .#test-race`   | nom + tui race-free ✓       |
| `nix run .#govulncheck` | 0 vulnerabilities ✓         |
| Git status              | **22 files uncommitted** ⚠️ |

---

## Session Inventory

**New files (19):**

- Root: `table_builder.go`, `table_builder_test.go`, `tree_builder.go`, `tree_builder_test.go`, `projections.go`, `projections_test.go`, `cqrs_example_test.go`
- delimited/: `cqrs.go`, `cqrs_test.go`, `cqrs_example_test.go`
- serialization/: `cqrs.go`, `cqrs_test.go`, `cqrs_example_test.go`
- markup/: `cqrs.go`, `cqrs_test.go`
- graph/: `cqrs_example_test.go`
- d2/: `cqrs_example_test.go`
- plantuml/: `cqrs_example_test.go`
- tree/: `cqrs_example_test.go`

**Modified files (3):**

- `README.md` (+94 lines: CQRS quick-start, migration guide)
- `CHANGELOG.md` (+18 lines: CQRS section)
- `AGENTS.md` (+3 pattern entries)

---

## Resolution (2026-08-04)

All items resolved in session 06:58. Serialization WriteJSON/WriteYAML now stream via stdlib encoders (no intermediate string). TableBuilder.Build() doc corrected (stopped saying "freeze"). Table CQRS added. The 22 uncommitted files were committed (`a70b46d`, `ade1816`, `375f353`). v0.30.0 tagged. The FrozenTable/FrozenTree question for v1.0.0 remains in ROADMAP.
