# Root Split Execution Plan — Strategy B

**Date:** 2026-05-25
**Approach:** Group by Domain — 3 new modules + registry-based dispatch
**Constraint:** Each task ≤ 12 min

---

## Architecture Decision

### New Module Map

```
go-output/                      # Root: core types + lightweight formatters (~25 files)
├── format.go, shape.go, renderer.go, ids.go, registry.go
├── sort.go, color.go, slices.go, tabledata.go
├── tree.go, graph.go, graph_tabledata.go
├── markdown.go                     # stays — zero external deps, 204 lines
├── marshal.go                      # stays — exported helpers used by 3 modules
├── render_tabledata.go             # stays — registry-based dispatch (no sub-module imports)
│
├── delimited/                      # NEW: CSV + TSV + DelimitedWriter
│   ├── delimited.go
│   ├── csv.go
│   ├── tsv.go
│   └── go.mod → root, escape (via root)
│
├── serialization/                  # NEW: JSON + YAML + renderers
│   ├── json.go, json_renderers.go
│   ├── yaml.go, yaml_renderers.go
│   └── go.mod → root
│
├── markup/                         # NEW: XML + HTML + Streaming + markup helpers
│   ├── markup.go, xml.go
│   ├── html.go, streaming.go
│   └── go.mod → root, escape
```

### Critical Design Decisions

1. **`marshal.go` stays in root** — exports `MarshalFormat()`, `UnmarshalFormat()`, `BrandedValue()` (currently unexported). Used by serialization/ AND markup/ (xml.go). Moving it to either would create inter-module dep.

2. **`render_tabledata.go` stays in root** — refactored to use **registry-based dispatch** instead of direct function calls. Each format module registers a `TableDataMarshaler` via `init()`. Root imports ZERO sub-modules. Consistent with existing pattern (D2/DOT/Mermaid already return `UnsupportedFormatError`).

3. **`tableDataBase` exported as `TableDataBase`** — currently unexported, embedded by JSON, YAML, HTML, Streaming renderers. Must be exported so sub-modules can embed it.

4. **`markdown.go` stays in root** — zero external deps, 204 lines, used by `render_tabledata.go`. Not worth a separate module.

5. **Test helpers** — each new module creates its own `_test.go` helpers. Shared assertions already live in `testhelpers/`. Graph/tree node helpers (currently in `output_test_helpers_test.go`) get duplicated in each module that needs them (~10 lines each).

### Dependency Graph (After)

```
root (output) → enum, escape, yaml, x/term, go-branded-id
  └── exports: TableDataBase, MarshalFormat, UnmarshalFormat, BrandedValue, RenderTableData

delimited/    → root (TableData, TableDataBase, Format)
serialization/ → root (Renderer, TableRenderer, TableDataBase, TreeNode, GraphRendererMixin, MarshalFormat, BrandedValue)
markup/       → root (Renderer, TableRenderer, TableDataBase, TreeNode), escape

integration/  → root, delimited, serialization, markup, d2, graph, table
examples/     → root, delimited, serialization, markup, d2, graph, table
```

**Root has ZERO sub-module imports.** Preserved.

---

## Phase 0: Foundation — Export Previously-Unexported Symbols

These unblock ALL extraction work. Must be done first, tested, committed.

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 01 | Export `tableDataBase` → `TableDataBase` in `tabledata.go`. Update all 4 embedders in root (`json.go`, `yaml.go`, `html.go`, `streaming.go`) to use new name. | 8m | HIGH — blocks extraction | `tabledata.go`, `json.go`, `yaml.go`, `html.go`, `streaming.go` |
| 02 | Export `marshal()` → `MarshalFormat()`, `unmarshal()` → `UnmarshalFormat()` in `marshal.go`. Update all callers in root: `json.go`, `yaml.go`, `xml.go`, `json_renderers.go`, `yaml_renderers.go`. | 8m | HIGH — blocks extraction | `marshal.go`, `json.go`, `yaml.go`, `xml.go`, `json_renderers.go`, `yaml_renderers.go` |
| 03 | Export `brandedValue()` → `BrandedValue()` in `marshal.go`. Update callers: `json_renderers.go`, `yaml_renderers.go`. | 4m | HIGH — blocks serialization/ | `marshal.go`, `json_renderers.go`, `yaml_renderers.go` |
| 04 | Add `TableDataMarshaler` type + `RegisterTableDataMarshaler()` to `render_tabledata.go`. Define: `type TableDataMarshaler func(w io.Writer, data *TableData, opts RenderOptions) error` and a `map[Format]TableDataMarshaler` with register function. | 6m | HIGH — enables registry dispatch | `render_tabledata.go` |
| 05 | Refactor `RenderTableData()` to use registered marshalers. Replace direct function calls with map lookup. Keep markdown/tree as direct calls (they stay in root). Register markdown/tree marshalers in `init()` within render_tabledata.go. | 10m | HIGH — core dispatch refactor | `render_tabledata.go` |
| 06 | Run full root test suite. Verify zero regressions from Phase 0 changes. | 3m | HIGH — gate | — |

---

## Phase 1: Extract `delimited/` Module

Cleanest extraction — `DelimitedWriter` has zero root deps. Only `csv.go`/`tsv.go` need `TableData`.

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 07 | Create `delimited/` directory. Create `go.mod` with replace directives for root. | 3m | MED — creates module | `delimited/go.mod` |
| 08 | Copy `delimited.go` → `delimited/delimited.go`. Change package to `delimited`. No import changes needed (pure stdlib). | 3m | MED | `delimited/delimited.go` |
| 09 | Copy `csv.go` → `delimited/csv.go`. Change package to `delimited`. Add `import "github.com/larsartmann/go-output"` for `TableData`. Fix `DelimitedWriter` refs (now local). Fix `TableData` → `output.TableData`. | 8m | MED | `delimited/csv.go` |
| 10 | Copy `tsv.go` → `delimited/tsv.go`. Same treatment as csv.go. | 8m | MED | `delimited/tsv.go` |
| 11 | Add `init()` to delimited module: register CSV and TSV as `TableDataMarshaler` via `output.RegisterTableDataMarshaler()`. Move `renderCSVTableData`/`renderTSVTableData` logic from root `render_tabledata.go` into delimited's init(). | 8m | HIGH — wires dispatch | `delimited/csv.go`, `delimited/tsv.go` |
| 12 | Copy `csv_test.go` + `tsv_test.go` → `delimited/`. Update package to `delimited_test`. Fix imports: `output.X` → `output.X` (root types), local types stay. Create local test helpers if needed. | 10m | MED | `delimited/csv_test.go`, `delimited/tsv_test.go` |
| 13 | Run `go mod tidy` + `go build ./...` + `go test ./...` in `delimited/`. Fix compilation errors. | 5m | HIGH — gate | — |
| 14 | Delete original `csv.go`, `tsv.go`, `delimited.go` from root. Update root `render_tabledata.go` to remove CSV/TSV direct calls (now registry-based). | 5m | HIGH — cleanup | root |

---

## Phase 2: Extract `markup/` Module

Medium complexity — `markup.go` is zero-dep, `xml.go`/`html.go`/`streaming.go` need root types + escape.

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 15 | Create `markup/` directory. Create `go.mod` with replace directives for root + escape. | 3m | MED | `markup/go.mod` |
| 16 | Copy `markup.go` → `markup/markup.go`. Change package to `markup`. Functions stay unexported (only used within this package). | 3m | MED | `markup/markup.go` |
| 17 | Copy `xml.go` → `markup/xml.go`. Change package to `markup`. Fix imports: `output.TableData`, `output.MarshalFormat()`, `escape.XMLEscape`. Fix `marshal()` → `output.MarshalFormat()`. Fix `writeMarkupRow` → local. | 8m | MED | `markup/xml.go` |
| 18 | Copy `html.go` → `markup/html.go`. Change package to `markup`. Fix imports: `output.Renderer`, `output.TableRenderer`, `output.TableDataBase`, `output.TreeNode`. Fix embedded types. Fix `writeMarkupRow` → local. | 10m | MED | `markup/html.go` |
| 19 | Copy `streaming.go` → `markup/streaming.go`. Change package to `markup`. Fix imports: `output.Renderer`, `output.TableRenderer`, `output.TableDataBase`, `output.StreamingRenderer`. Fix embedded types. | 8m | MED | `markup/streaming.go` |
| 20 | Add `init()` to markup module: register XML and HTML as `TableDataMarshaler`. Move `renderXMLTableData`/`renderHTMLTableData` logic from root. | 8m | HIGH | `markup/xml.go`, `markup/html.go` |
| 21 | Copy test files: `xml_test.go`, `html_test.go`, `streaming_test.go`, `markup_test.go` → `markup/`. Update package to `markup_test`. Create local test helpers for graph/tree nodes (thin wrappers calling `output.NewBrandedID`). | 12m | MED | `markup/*_test.go` |
| 22 | Run `go mod tidy` + `go build ./...` + `go test ./...` in `markup/`. Fix compilation errors. | 5m | HIGH — gate | — |
| 23 | Delete original `xml.go`, `html.go`, `streaming.go`, `markup.go` from root. Update root `render_tabledata.go` to remove XML/HTML direct calls. | 5m | HIGH | root |

---

## Phase 3: Extract `serialization/` Module

Heaviest — yaml external dep, renderers with graph/tree types.

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 24 | Create `serialization/` directory. Create `go.mod` with replace directives for root. | 3m | MED | `serialization/go.mod` |
| 25 | Copy `json.go` → `serialization/json.go`. Change package to `serialization`. Fix imports: `output.Renderer`, `output.TableRenderer`, `output.TableDataBase`, `output.MarshalFormat()`, `output.UnmarshalFormat()`, `output.TableData`. | 8m | MED | `serialization/json.go` |
| 26 | Copy `json_renderers.go` → `serialization/json_renderers.go`. Change package to `serialization`. Fix imports: `output.Renderer`, `output.TreeNode`, `output.TreeOutputRenderer`, `output.GraphRendererMixin`, `output.GraphRenderer`, `output.NewGraphRendererMixin`, `output.BrandedValue()`. | 10m | MED | `serialization/json_renderers.go` |
| 27 | Copy `yaml.go` → `serialization/yaml.go`. Same treatment as json.go but with yaml imports. | 8m | MED | `serialization/yaml.go` |
| 28 | Copy `yaml_renderers.go` → `serialization/yaml_renderers.go`. Same treatment as json_renderers.go. | 10m | MED | `serialization/yaml_renderers.go` |
| 29 | Add `init()` to serialization module: register YAML as `TableDataMarshaler`. Move `renderYAMLTableData` logic from root. JSON is NOT registered (RenderTableData doesn't handle JSON — it returns UnsupportedFormatError). | 6m | HIGH | `serialization/yaml.go` |
| 30 | Copy test files: `json_test.go`, `json_renderers_test.go`, `yaml_test.go`, `yaml_renderers_test.go` → `serialization/`. Update package to `serialization_test`. Create local test helpers for graph/tree nodes. | 12m | MED | `serialization/*_test.go` |
| 31 | Run `go mod tidy` + `go build ./...` + `go test ./...` in `serialization/`. Fix compilation errors. | 5m | HIGH — gate | — |
| 32 | Delete original `json.go`, `json_renderers.go`, `yaml.go`, `yaml_renderers.go` from root. Update root `render_tabledata.go` to remove YAML direct call. | 5m | HIGH | root |

---

## Phase 4: Update Dependent Modules

Integration + examples need new import paths. Other sub-modules (d2, graph, table) likely unaffected.

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 33 | Update `integration/go.mod`: add `delimited`, `serialization`, `markup` to require + replace. | 3m | MED | `integration/go.mod` |
| 34 | Update `integration/integration_test.go`: add imports for `delimited`, `serialization`, `markup`. Fix `output.NewCSVWriter` → `delimited.NewCSVWriter`, `output.MarshalJSON` → `serialization.MarshalJSON`, `output.NewHTMLRenderer` → `markup.NewHTMLRenderer`, etc. | 10m | HIGH — tests must pass | `integration/integration_test.go` |
| 35 | Update `integration/renderer_test.go`: same import fixes as above. | 8m | HIGH | `integration/renderer_test.go` |
| 36 | Update `integration/workflow_test.go`: fix `output.MarshalJSONIndent` → `serialization.MarshalJSONIndent`, `output.NewStreamingHTMLRenderer` → `markup.NewStreamingHTMLRenderer`. | 6m | HIGH | `integration/workflow_test.go` |
| 37 | Update `integration/test_helpers.go`: fix `output.MarshalJSON` → `serialization.MarshalJSON`. | 3m | MED | `integration/test_helpers.go` |
| 38 | Update `examples/go.mod`: add `delimited`, `serialization`, `markup` to require + replace. | 3m | MED | `examples/go.mod` |
| 39 | Update `examples/basic/main.go`: fix all format-specific imports. | 8m | MED | `examples/basic/main.go` |

---

## Phase 5: Root Cleanup + Verify

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 40 | Clean root test files: remove moved test helpers from `output_test_helpers_test.go` (graph/tree helpers only used by moved tests). Keep helpers used by remaining root tests. | 8m | MED | `output_test_helpers_test.go` |
| 41 | Remove `benchmarks_test.go` entries for moved formats (CSV, TSV, JSON, XML benchmarks) or move them to respective modules. | 5m | LOW | `benchmarks_test.go` |
| 42 | Remove `fuzz_test.go` entries for moved formats or move them. | 5m | LOW | `fuzz_test.go` |
| 43 | Remove `userjourney_test.go` format-specific parts or update imports. | 8m | MED | `userjourney_test.go` |
| 44 | Update root `go.mod`: remove `go-faster/yaml` if no longer needed (check `yaml.go` moved). Remove `escape` if no longer needed (check only `xml.go`/`html.go` used it). Run `go mod tidy`. | 5m | HIGH — dep isolation | `go.mod` |
| 45 | Run full root test suite: `go test ./...` in root. Verify all tests pass. | 3m | HIGH — gate | — |

---

## Phase 6: Documentation + Final Verification

| # | Task | Est. | Impact | Files |
|---|------|------|--------|-------|
| 46 | Update `AGENTS.md`: new module table, dependency graph, module count (9→12), new import paths, updated coverage. | 10m | HIGH — future sessions depend on this | `AGENTS.md` |
| 47 | Update `docs/modularization/root-split-proposal.md`: mark as implemented, add actual vs planned notes. | 5m | LOW | `docs/modularization/root-split-proposal.md` |
| 48 | Update `README.md`: new import path examples for serialization/delimited/markup. | 8m | HIGH — user-facing | `README.md` |
| 49 | Update `go.work` example in AGENTS.md: add `./serialization`, `./delimited`, `./markup`. | 3m | MED | `AGENTS.md` |
| 50 | Run ALL modules: `for mod in . delimited serialization markup d2 graph enum escape testhelpers table integration examples; do (cd $mod && go test ./...); done` | 5m | HIGH — final gate | — |
| 51 | Run `golangci-lint` across ALL modules. Fix any issues. | 8m | HIGH — quality gate | — |

---

## Summary

| Metric | Value |
|--------|-------|
| **Total tasks** | 51 |
| **Phases** | 7 (0-6) |
| **New modules** | 3 (`delimited/`, `serialization/`, `markup/`) |
| **Root file reduction** | 54 → ~25 |
| **Root deps removed** | `go-faster/yaml`, `escape` (if fully moved) |
| **Breaking changes** | Users must import format-specific modules: `import "go-output/serialization"` instead of `import "go-output"` for JSON/YAML |
| **Registry dispatch** | `render_tabledata.go` uses `TableDataMarshaler` registry — no sub-module imports in root |
| **Key invariant preserved** | Root has ZERO sub-module imports |

## Estimated Total Time

| Phase | Tasks | Est. Time |
|-------|-------|-----------|
| Phase 0: Foundation | 6 | ~39 min |
| Phase 1: delimited/ | 8 | ~50 min |
| Phase 2: markup/ | 9 | ~62 min |
| Phase 3: serialization/ | 9 | ~67 min |
| Phase 4: Dependents | 7 | ~41 min |
| Phase 5: Cleanup | 6 | ~34 min |
| Phase 6: Docs + Verify | 6 | ~39 min |
| **Total** | **51** | **~332 min (~5.5h)** |

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Unexported type embedding breaks | Medium | High | Phase 0 exports `TableDataBase` first |
| Test helper duplication across modules | High | Low | ~10 lines each, acceptable tradeoff |
| `render_tabledata.go` registry not populated | Medium | High | Clear error message if format not imported |
| Circular imports between new modules | Low | High | `marshal.go` stays in root; modules only import root |
| Integration tests break | High | Medium | Phase 4 dedicates 5 tasks to fixing imports |
| `go-faster/yaml` still in root go.mod | Low | Low | Verify after Phase 5 cleanup |
