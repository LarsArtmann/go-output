# Status Report — CQRS Architecture: Streaming Fix + Integration Tests

**Date:** 2026-07-06 06:58
**Branch:** master (pushed)
**Last commit:** `375f353` — feat: complete CQRS streaming, table/ module, integration tests, and doc fixes
**Session commits:** `a70b46d` → `ade1816` → `375f353` (3 commits, 26 files, +1887 lines)

---


> **✅ Resolved (2026-08-04):**
>
> Golden-file tests for CQRS streaming output added in session 09:52. Registry dispatch rewired to CQRS streaming functions (byte-for-byte identical output proven by `TestCQRS_StreamVsRegistry_JSON/CSV`). TableToGraph refactored to functional options. graph.WriteDOT writes directly from Graph. v0.30.0 tagged. HTML/AsciiDoc still buffer (no streaming writer exists — documented).

---

## A. FULLY DONE ✓

### CQRS Streaming Architecture (All 7 Formats)

Every `WriteXxx(w, data)` function now streams directly to `io.Writer` via the standard library's streaming encoder — no intermediate `[]byte`/`string` allocation:

| Format | Old Path                          | New Path                                       | Reuses                 |
| ------ | --------------------------------- | ---------------------------------------------- | ---------------------- |
| JSON   | `MarshalIndent` → `Fprint`        | `json.NewEncoder(w).Encode()`                  | stdlib                 |
| YAML   | `yaml.Marshal` → `Fprint`         | `yaml.NewEncoder(w).Encode()`                  | go-faster/yaml         |
| TOML   | `toml.Marshal` → `Write`          | `toml.NewEncoder(w).Encode()`                  | pelletier/go-toml/v2   |
| JSONL  | `MarshalJSONLFromTable` → `Write` | `NewJSONLWriter(w)` per-row `Encode` + `Flush` | existing `JSONLWriter` |
| CSV    | `MarshalCSVFromTable` → `Write`   | `NewCSVWriter(w)` row-by-row                   | existing `CSVWriter`   |
| TSV    | `MarshalTSVFromTable` → `Write`   | `NewTSVWriter(w)` row-by-row                   | existing `TSVWriter`   |
| XML    | `MarshalXMLFromTable` → `Write`   | `NewXMLWriter(w)` element-by-element           | existing `XMLWriter`   |

### CQRS API Surface (Complete for v0.30.0)

| Module         | Write Function                                      | Render Function                                         | Status               |
| -------------- | --------------------------------------------------- | ------------------------------------------------------- | -------------------- |
| delimited/     | `WriteCSV`, `WriteTSV`                              | `RenderCSV`, `RenderTSV`                                | ✓                    |
| serialization/ | `WriteJSON`, `WriteYAML`, `WriteTOML`, `WriteJSONL` | `RenderJSON`, `RenderYAML`, `RenderTOML`, `RenderJSONL` | ✓                    |
| markup/        | `WriteXML`, `WriteAsciiDoc`, `WriteHTML`            | `RenderXML`, `RenderAsciiDoc`, `RenderHTML`             | ✓                    |
| graph/         | `WriteDOT`, `WriteMermaid`                          | `RenderDOT`, `RenderMermaid`                            | ✓                    |
| d2/            | `Write`, `WriteGraph`                               | `Render`, `RenderGraph`                                 | ✓                    |
| plantuml/      | `Write`                                             | `Render`                                                | ✓                    |
| tree/          | `WriteASCII`                                        | `RenderASCII`                                           | ✓                    |
| markdown/      | `Write`                                             | `Render`                                                | ✓                    |
| **table/**     | `Write`                                             | `Render`                                                | ✓ (NEW this session) |

### Builders + Projections (Root Package)

- [x] `TableBuilder` — `SetHeaders`/`AddRow`/`AddRows`/`SetFooter`/`Build()→*Table`
- [x] `TreeBuilder` — `SetRoot`/`AddChild`/`Build()→*TreeNode`
- [x] `GraphBuilder.Build()→Graph` (truly immutable — unexported fields)
- [x] `TableToGraph(t, labelFn...)` — optional label function
- [x] `GraphToTree(g)→*TreeNode` — cycle-guarded
- [x] `GraphToTable(g)→*Table`

### Tests + Benchmarks

- [x] 14 builder/projection unit tests (root)
- [x] 18 CQRS function tests (3 modules: delimited, serialization, markup)
- [x] 6 integration tests (pipeline end-to-end + streaming vs registry consistency)
- [x] 3 benchmarks (TableToGraph, GraphBuilder.Build freeze, TableBuilder.Build)
- [x] 7 Godoc example files across all renderer modules

### Documentation Fixes

- [x] `TableBuilder.Build()` doc: dropped "freeze" lie, honest about mutability
- [x] `TreeBuilder` doc: same honest language
- [x] CHANGELOG: removed stale self-referential "Deprecated" section (symbols were deleted, not deprecated)
- [x] README: CQRS quick-start + migration guide added

### Verification — ALL GREEN

- [x] Build: 19/19 modules ✓
- [x] Tests: 19/19 modules ✓
- [x] Lint: 0 issues ✓
- [x] Race tests: nom + tui race-free ✓
- [x] Govulncheck: 0 vulnerabilities ✓
- [x] Git: clean, pushed to origin/master ✓

---

## B. PARTIALLY DONE ⚠️

### Registry Dispatch Still Uses RenderOptions Struct

The root `output.RenderTable(data, fmt, RenderOptions{Writer: w})` dispatch still uses the rigid options struct. The CQRS vision (D2 decision) calls for functional options everywhere. This is a medium-priority cleanup — the struct works fine, it's just inconsistent with the functional-options pattern used by every CQRS function.

### CQRS Functions Don't Accept RenderOptions

The new `WriteJSON(w, data)`, `WriteCSV(w, data)`, etc. don't accept any options. The registry dispatch path (`RenderTable`) passes `RenderOptions{ColorMode: ...}` but the CQRS path has no way to receive options. For serialization/delimited/markup this is fine (no color settings), but for `table.Render(data, opts...)` the options ARE passed through. The API is inconsistent — some CQRS functions take `opts...`, some don't.

### HTML/AsciiDoc Write Functions Still Buffer

`WriteHTML` delegates to `renderHTMLTable` which goes through `renderViaRenderer` (buffer-then-write). `WriteAsciiDoc` goes through `MarshalAsciiDocFromTable` → `[]byte` → `Write`. These two formats don't have streaming writers in the existing codebase. They're lower priority since HTML/AsciiDoc are less commonly streamed.

---

## C. NOT STARTED ❌

### v0.30.0 Release Blockers

- [ ] **Tag testhelpers/ v0.30.0** (D4 decision — testhelpers is independently versioned)
- [ ] **Tag root v0.30.0** — all code is ready, just needs the tag
- [ ] **Registry: wire CQRS functions** — `output.RegisterFormat(fmt, RenderFunc)` so registry dispatch uses streaming path (currently registry still calls the old renderer structs internally)

### Architectural

- [ ] **Delete old renderer structs** (DOTRenderer, MermaidRenderer, etc.) — v0.31.0 decision; currently kept as implementation detail backing CQRS functions
- [ ] **Make `graph.RenderDOT` read from `Graph` without intermediate DOTRenderer** — the CQRS function creates a DOTRenderer internally, populates it from Graph, then calls Render(). Should write DOT directly from Graph nodes/edges.
- [ ] **FrozenTable / FrozenTree** — real immutability for Table/Tree types (unexported fields + accessor methods). This is a v1.0.0 architectural decision.

### Missing CQRS

- [ ] `tree.RenderMarkdown(root)` / `tree.WriteMarkdown(w, root)` — no markdown tree renderer exists
- [ ] `daghtml.Render(g)` / `daghtml.Write(w, g)` — daghtml module has no CQRS functions

### Missing Tests

- [ ] **Golden-file tests for CQRS functions** — the streaming refactoring changed output formatting (encoders may add trailing newlines differently). No golden-file regression test exists for CQRS output.
- [ ] **Streaming vs non-streaming output equivalence test** — the integration test checks content contains expected substrings, but doesn't verify byte-for-byte equivalence between old registry path and new CQRS path. They may differ in trailing newlines.
- [ ] **Projection edge case tests** — `GraphToTree` with disconnected nodes, `TableToGraph` with custom label functions tested in benchmark only, not unit test

### Documentation

- [ ] **Update AGENTS.md** with the streaming architecture pattern (encoders vs marshalers)
- [ ] **Update FEATURES.md** with the complete CQRS API listing
- [ ] **ADR for CQRS streaming decision** — document why `json.NewEncoder` over `json.MarshalIndent`

---

## D. TOTALLY FUCKED UP 💥

### None — No regressions, no broken state, no data loss.

### Potential Issues Spotted (but not confirmed):

1. **Streaming output may differ from registry output in trailing newlines.** `json.Encoder.Encode()` appends `\n`, while `json.MarshalIndent` does not. The integration test `TestCQRS_StreamVsRegistry_JSON` only checks `strings.Contains`, not exact equality. If a consumer parses the output with a strict parser that rejects trailing newlines, this is a behavior change. **Not tested for exact equivalence.**
2. **`TableToGraph` variadic label func signature is unusual.** `TableToGraph(data *Table, labelFn ...GraphNodeLabelFunc) Graph` — a variadic optional function is a Go anti-pattern. The idiomatic approach would be a functional option: `TableToGraph(data, WithLabelFunc(fn))`. This should be refactored before v1.0.0.

---

## E. WHAT WE SHOULD IMPROVE 🎯

### 1. Output Equivalence Testing Is Missing

The streaming refactoring changed the code path for every serialization format. JSON's `Encoder.Encode()` adds a trailing `\n` that `MarshalIndent` doesn't. The integration tests check substring containment, not exact output. We need at least one golden-file test per format to catch subtle output differences.

### 2. The Registry Dispatch Is Now the Bottleneck

CQRS `WriteJSON` streams directly via `json.NewEncoder`. But `output.RenderTable(data, FormatJSON, opts)` still calls the old `renderJSONTable` → `renderViaRenderer` → `JSONTableRenderer.Render()` → `MarshalIndent` path. The registry dispatch is now SLOWER and uses MORE MEMORY than the CQRS path. The registry should be updated to call the CQRS functions.

### 3. `TableToGraph`'s Variadic Label Func Is an Anti-Pattern

```go
func TableToGraph(data *Table, labelFn ...GraphNodeLabelFunc) Graph
```

This should be:

```go
func TableToGraph(data *Table, opts ...TableToGraphOption) Graph
```

with `WithLabelFunc(fn)` as the option. More extensible, idiomatic Go.

### 4. HTML and AsciiDoc Don't Stream

All other formats stream. HTML and AsciiDoc still buffer. This is an inconsistency that should be resolved — either by adding streaming writers for those formats or documenting that they buffer.

### 5. The "Freeze" Boundary Is Still a Lie for Table/Tree

`Graph` is properly immutable. `*Table` and `*TreeNode` have exported fields. The docs now say this honestly ("callers SHOULD treat as read-only"), but the type system doesn't enforce it. This is acceptable for v0.30.0 but should be addressed before v1.0.0.

### 6. TableBuilder.AddRows Missing Test

Added `AddRows` method but no dedicated test for it. The existing `TestTableBuilder_FluentAPI` doesn't test bulk row addition.

---

## F. TOP 25 THINGS TO DO NEXT 📋

| #   | Task                                                                                | Impact   | Effort | Deps |
| --- | ----------------------------------------------------------------------------------- | -------- | ------ | ---- |
| 1   | Add golden-file test for CQRS streaming output (JSON, YAML, CSV)                    | High     | 10m    | —    |
| 2   | Test `TableBuilder.AddRows` with dedicated unit test                                | Low      | 5m     | —    |
| 3   | Wire CQRS functions into registry dispatch (`output.RenderTable` calls `WriteJSON`) | High     | 15m    | —    |
| 4   | Refactor `TableToGraph` variadic labelFunc → functional options                     | Med      | 10m    | —    |
| 5   | Add exact output equivalence test (CQRS vs registry, byte-for-byte)                 | High     | 10m    | 1    |
| 6   | Make `graph.RenderDOT` write directly from Graph (no DOTRenderer intermediary)      | High     | 20m    | —    |
| 7   | Add streaming HTML writer (or document that HTML buffers)                           | Low      | 15m    | —    |
| 8   | Add streaming AsciiDoc writer (or document that AsciiDoc buffers)                   | Low      | 15m    | —    |
| 9   | Update AGENTS.md with streaming architecture pattern                                | Med      | 5m     | —    |
| 10  | Update FEATURES.md with complete CQRS API listing                                   | Med      | 10m    | —    |
| 11  | Add `tree.RenderMarkdown(root)` CQRS                                                | Low      | 10m    | —    |
| 12  | Add `daghtml.Render(g)` CQRS                                                        | Low      | 10m    | —    |
| 13  | Add `GraphToTree` test for disconnected nodes (forest)                              | Low      | 5m     | —    |
| 14  | Add `TableToGraph` test with custom label function                                  | Low      | 5m     | 4    |
| 15  | Consider FrozenTable/FrozenTree types for real immutability                         | High     | 30m    | —    |
| 16  | Delete old renderer structs (DOTRenderer etc.) — v0.31.0                            | High     | 30m    | 6    |
| 17  | Unify `RenderOptions` struct → functional options for root dispatch                 | Med      | 15m    | 3    |
| 18  | Tag testhelpers/ v0.30.0                                                            | Low      | 3m     | —    |
| 19  | Tag root v0.30.0                                                                    | Critical | 2m     | 1-17 |
| 20  | Add streaming benchmarks (compare old marshal vs new encoder)                       | Med      | 15m    | —    |
| 21  | Add CQRS example to `examples/` directory                                           | Med      | 10m    | —    |
| 22  | Add ADR for CQRS streaming decision                                                 | Low      | 10m    | —    |
| 23  | Add `TreeBuilder.AddChildren(parentID, children...)` bulk method                    | Low      | 5m     | —    |
| 24  | Investigate stale gopls cache (NodeShapeRect, physicalLines warnings)               | Low      | 10m    | —    |
| 25  | Add `TableBuilder.SetFooterRow([]string)` alternative API                           | Low      | 5m     | —    |

---

## G. TOP #1 QUESTION 🤔

**Should the registry dispatch (`output.RenderTable`) be rewired to call the CQRS streaming functions, or should it keep calling the old renderer structs?**

Currently:

```
output.RenderTable(data, FormatJSON, opts)
  → renderJSONTable(w, data, opts)
    → renderViaRenderer(w, data, NewJSONTableRenderer(), "json")
      → JSONTableRenderer.Render() → json.MarshalIndent → string → Fprint
```

After rewiring:

```
output.RenderTable(data, FormatJSON, opts)
  → serialization.WriteJSON(w, data)
    → json.NewEncoder(w).Encode()
```

The rewire is cleaner, faster, and uses less memory. BUT: the output may differ slightly (trailing newlines from Encoder). This could break consumers who do exact-output comparisons.

**I cannot resolve this without knowing whether any consumer depends on the exact byte output of the registry path.** If trailing newlines changed, golden-file tests across downstream projects would break.

---

## Verification Summary

| Check                   | Result                |
| ----------------------- | --------------------- |
| `nix run .#build`       | 19/19 ✓               |
| `nix run .#test`        | 19/19 ✓               |
| `nix run .#lint`        | 0 issues ✓            |
| `nix run .#test-race`   | nom + tui race-free ✓ |
| `nix run .#govulncheck` | 0 vulnerabilities ✓   |
| Git status              | Clean, pushed ✓       |

---

## Session Inventory

**3 commits, 26 files, +1887 lines:**

- `a70b46d` — TableBuilder, TreeBuilder, projections, CQRS Writer/Render for 3 modules, Godoc examples, docs
- `ade1816` — Streaming refactoring: all WriteXxx now use encoders directly
- `375f353` — table/ CQRS, integration tests, benchmarks, doc fixes, AddRows, CHANGELOG cleanup

