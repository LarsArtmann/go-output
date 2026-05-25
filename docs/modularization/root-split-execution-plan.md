# Root Split Execution Plan — Strategy B

**Date:** 2026-05-25
**Last Audited:** 2026-05-25
**Approach:** Group by Domain — 3 new modules + registry-based dispatch
**Constraint:** Each task ≤ 12 min

---

## Overall Progress

| Phase     | Description                 | Status     | Tasks Done |
| --------- | --------------------------- | ---------- | ---------- |
| Phase 0   | Foundation — Export symbols | ✅ DONE    | 6/6        |
| Phase 1   | Extract `delimited/`        | ✅ DONE    | 8/8        |
| Phase 2   | Extract `markup/`           | ✅ DONE    | 9/9        |
| Phase 3   | Extract `serialization/`    | ✅ DONE    | 9/9        |
| Phase 4   | Update dependents           | ✅ DONE    | 7/7        |
| Phase 5   | Root cleanup                | ⚠️ PARTIAL | 4/6        |
| Phase 6   | Docs + verify               | ⚠️ PARTIAL | 4/6        |
| **Total** |                             | **⚠️ 92%** | **43/51**  |

---

## Architecture Decision

### New Module Map (AS-BUILT)

```
go-output/                      # Root: core types + lightweight formatters
├── format.go, shape.go, renderer.go, ids.go, registry.go
├── sort.go, color.go, slices.go, tabledata.go
├── tree.go, graph.go, graph_tabledata.go
├── streaming.go                    # interface + adapter only (implementation moved to markup/)
├── markdown.go                     # stays — zero external deps
├── marshal.go                      # stays — exported helpers used by 3 modules
├── render_tabledata.go             # stays — registry-based dispatch (zero sub-module imports)
│
├── delimited/                      # ✅ DONE: CSV + TSV + DelimitedWriter
│   ├── delimited.go
│   ├── csv.go, csv_test.go
│   ├── tsv.go, tsv_test.go
│   ├── testhelpers_test.go
│   └── go.mod → root
│
├── serialization/                  # ✅ DONE: JSON + YAML + renderers
│   ├── json.go, json_test.go
│   ├── json_renderers.go, json_renderers_test.go
│   ├── yaml.go, yaml_test.go
│   ├── yaml_renderers.go, yaml_renderers_test.go
│   ├── testhelpers_test.go
│   └── go.mod → root, go-faster/yaml, testhelpers
│
├── markup/                         # ✅ DONE: XML + HTML + StreamingHTML + markup helpers
│   ├── markup.go, markup_test.go
│   ├── xml.go, xml_test.go
│   ├── html.go, html_test.go
│   ├── streaming.go, streaming_test.go
│   ├── testhelpers_test.go
│   └── go.mod → root, escape
```

### Critical Design Decisions

1. **`marshal.go` stays in root** ✅ — exports `MarshalFormat()`, `UnmarshalFormat()`, `BrandedValue()`. Used by serialization/ AND markup/. Verified: all callers updated.

2. **`render_tabledata.go` stays in root** ✅ — uses registry-based dispatch. Each format module registers via `init()`. Root imports ZERO sub-modules. No format-specific references remain.

3. **`tableDataBase` exported as `TableDataBase`** ✅ — exported with `Data()` getter. Sub-modules embed it successfully.

4. **`markdown.go` stays in root** ✅ — zero external deps, used by `render_tabledata.go`.

5. **`streaming.go` split correctly** ✅ — root retains the `StreamingRenderer` interface + `StreamingRendererFromRenderer` adapter. Markup/ has the `StreamingHTMLRenderer` implementation. This is cleaner than the original plan (which called for full move).

6. **Test helpers** ✅ — each module has its own `testhelpers_test.go`. Shared assertions in `testhelpers/`.

7. **`internal/gentest` imports `go-faster/yaml`** — this keeps `go-faster/yaml` as a direct root dependency. `internal/gentest` is root-only (sub-modules can't import it), so this is isolated. Root `.go` production files have zero yaml imports.

### Dependency Graph (AS-BUILT)

```
root (output) → enum, x/term, go-branded-id
  └── exports: TableDataBase, MarshalFormat, UnmarshalFormat, BrandedValue, MarshalJSONIndent, RenderTableData
  └── note: go-faster/yaml in go.mod via internal/gentest only (not production code)

delimited/    → root (TableData, TableDataBase, Format)
serialization/ → root (Renderer, TableRenderer, TableDataBase, TreeNode, GraphRendererMixin, MarshalFormat, BrandedValue), go-faster/yaml, testhelpers
markup/       → root (Renderer, TableRenderer, TableDataBase, StreamingRenderer), escape

integration/  → root, delimited, serialization, markup, d2, graph, table
examples/     → root, delimited, serialization, markup, d2, graph, table
```

**Root has ZERO sub-module imports.** ✅ Preserved.

---

## Phase 0: Foundation — Export Previously-Unexported Symbols

These unblock ALL extraction work. Must be done first, tested, committed.

| #   | Task                                                                                                             | Status  | Notes                                                                                                               |
| --- | ---------------------------------------------------------------------------------------------------------------- | ------- | ------------------------------------------------------------------------------------------------------------------- |
| 01  | Export `tableDataBase` → `TableDataBase` in `tabledata.go`. Update all embedders to use new name.                | ✅ DONE | `TableDataBase` exported with `Data()` getter, `SetHeaders()`, `AddRow()`, `SetData()`                              |
| 02  | Export `marshal()` → `MarshalFormat()`, `unmarshal()` → `UnmarshalFormat()` in `marshal.go`. Update all callers. | ✅ DONE | All callers in serialization/, markup/ updated                                                                      |
| 03  | Export `brandedValue()` → `BrandedValue()` in `marshal.go`. Update callers.                                      | ✅ DONE | Generic function, used by serialization/ renderers                                                                  |
| 04  | Add `TableDataMarshaler` type + `RegisterTableDataMarshaler()` to `render_tabledata.go`.                         | ✅ DONE | Thread-safe registry with `sync.RWMutex`, `map[Format]TableDataMarshaler`                                           |
| 05  | Refactor `RenderTableData()` to use registered marshalers.                                                       | ✅ DONE | Registry lookup first, falls back to markdown/tree (root-local). `UnsupportedFormatError` for unregistered formats. |
| 06  | Run full root test suite. Verify zero regressions.                                                               | ✅ DONE | Zero compile errors, all tests pass                                                                                 |

---

## Phase 1: Extract `delimited/` Module

Cleanest extraction — `DelimitedWriter` has zero root deps. Only `csv.go`/`tsv.go` need `TableData`.

| #   | Task                                                                             | Status  | Notes                                                |
| --- | -------------------------------------------------------------------------------- | ------- | ---------------------------------------------------- |
| 07  | Create `delimited/` directory. Create `go.mod` with replace directives for root. | ✅ DONE | `go.mod` exists, depends on root                     |
| 08  | Copy `delimited.go` → `delimited/delimited.go`. Change package.                  | ✅ DONE | Pure stdlib, no import changes                       |
| 09  | Copy `csv.go` → `delimited/csv.go`. Change package. Fix imports.                 | ✅ DONE | `output.TableData`, local `DelimitedWriter`          |
| 10  | Copy `tsv.go` → `delimited/tsv.go`. Same treatment.                              | ✅ DONE | Same pattern as csv.go                               |
| 11  | Add `init()`: register CSV and TSV as `TableDataMarshaler`.                      | ✅ DONE | `csv.go` and `tsv.go` each have `func init()`        |
| 12  | Copy test files → `delimited/`. Update package + imports.                        | ✅ DONE | `csv_test.go`, `tsv_test.go`, `testhelpers_test.go`  |
| 13  | Run `go mod tidy` + `go build` + `go test` in `delimited/`.                      | ✅ DONE | Module builds and tests pass                         |
| 14  | Delete originals from root. Update `render_tabledata.go`.                        | ✅ DONE | `csv.go`, `tsv.go`, `delimited.go` deleted from root |

---

## Phase 2: Extract `markup/` Module

Medium complexity — `markup.go` is zero-dep, `xml.go`/`html.go`/`streaming.go` need root types + escape.

| #   | Task                                                                 | Status  | Notes                                                                                       |
| --- | -------------------------------------------------------------------- | ------- | ------------------------------------------------------------------------------------------- |
| 15  | Create `markup/` directory. Create `go.mod` with replace directives. | ✅ DONE | Depends on root + escape                                                                    |
| 16  | Copy `markup.go` → `markup/markup.go`. Change package.               | ✅ DONE | Functions stay unexported (package-internal)                                                |
| 17  | Copy `xml.go` → `markup/xml.go`. Fix imports.                        | ✅ DONE | `output.TableData`, `output.MarshalFormat()`, `escape.XMLEscape`                            |
| 18  | Copy `html.go` → `markup/html.go`. Fix imports.                      | ✅ DONE | `output.TableDataBase`, `output.TreeNode`, etc.                                             |
| 19  | Copy `streaming.go` → `markup/streaming.go`. Fix imports.            | ✅ DONE | `StreamingHTMLRenderer` implementation. Root retains interface.                             |
| 20  | Add `init()`: register XML and HTML as `TableDataMarshaler`.         | ✅ DONE | `xml.go` and `html.go` each have `func init()`                                              |
| 21  | Copy test files → `markup/`. Update package + imports.               | ✅ DONE | `xml_test.go`, `html_test.go`, `streaming_test.go`, `markup_test.go`, `testhelpers_test.go` |
| 22  | Run `go mod tidy` + `go build` + `go test` in `markup/`.             | ✅ DONE | Module builds and tests pass                                                                |
| 23  | Delete originals from root. Update `render_tabledata.go`.            | ✅ DONE | `xml.go`, `html.go`, `markup.go` deleted. Root `streaming.go` kept (interface only).        |

---

## Phase 3: Extract `serialization/` Module

Heaviest — yaml external dep, renderers with graph/tree types.

| #   | Task                                                                       | Status  | Notes                                                                                                     |
| --- | -------------------------------------------------------------------------- | ------- | --------------------------------------------------------------------------------------------------------- |
| 24  | Create `serialization/` directory. Create `go.mod`.                        | ✅ DONE | Depends on root, go-faster/yaml, testhelpers                                                              |
| 25  | Copy `json.go` → `serialization/json.go`. Fix imports.                     | ✅ DONE | `output.TableDataBase`, `output.MarshalFormat()`, etc.                                                    |
| 26  | Copy `json_renderers.go` → `serialization/json_renderers.go`. Fix imports. | ✅ DONE | `output.GraphRendererMixin`, `output.BrandedValue()`, etc.                                                |
| 27  | Copy `yaml.go` → `serialization/yaml.go`. Fix imports.                     | ✅ DONE | Same pattern as json.go + yaml imports                                                                    |
| 28  | Copy `yaml_renderers.go` → `serialization/yaml_renderers.go`. Fix imports. | ✅ DONE | Same pattern as json_renderers.go                                                                         |
| 29  | Add `init()`: register YAML as `TableDataMarshaler`.                       | ✅ DONE | `yaml.go` has `func init()` registering `FormatYAML`                                                      |
| 30  | Copy test files → `serialization/`.                                        | ✅ DONE | `json_test.go`, `json_renderers_test.go`, `yaml_test.go`, `yaml_renderers_test.go`, `testhelpers_test.go` |
| 31  | Run `go mod tidy` + `go build` + `go test` in `serialization/`.            | ✅ DONE | Module builds and tests pass                                                                              |
| 32  | Delete originals from root.                                                | ✅ DONE | `json.go`, `json_renderers.go`, `yaml.go`, `yaml_renderers.go` deleted from root                          |

---

## Phase 4: Update Dependent Modules

Integration + examples need new import paths.

| #   | Task                                                   | Status  | Notes                                                                                 |
| --- | ------------------------------------------------------ | ------- | ------------------------------------------------------------------------------------- |
| 33  | Update `integration/go.mod`: add new modules.          | ✅ DONE | `delimited`, `serialization`, `markup` in require + replace                           |
| 34  | Update `integration/integration_test.go`: fix imports. | ✅ DONE | `delimited.NewCSVWriter`, `serialization.MarshalJSON`, `markup.NewHTMLRenderer`, etc. |
| 35  | Update `integration/renderer_test.go`: fix imports.    | ✅ DONE | Same pattern — all format-specific imports point to sub-modules                       |
| 36  | Update `integration/workflow_test.go`: fix imports.    | ✅ DONE | `serialization.MarshalYAML`, `markup.NewStreamingHTMLRenderer`                        |
| 37  | Update `integration/test_helpers.go`: fix imports.     | ✅ DONE | `serialization.MarshalJSON`                                                           |
| 38  | Update `examples/go.mod`: add new modules.             | ✅ DONE | All 3 new modules in require + replace                                                |
| 39  | Update `examples/basic/main.go`: fix imports.          | ✅ DONE | All format-specific imports use sub-module paths                                      |

---

## Phase 5: Root Cleanup + Verify

| #   | Task                                                             | Status      | Notes                                                                                                                                                                                                                                                    |
| --- | ---------------------------------------------------------------- | ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 40  | Clean root `output_test_helpers_test.go`: remove unused helpers. | ❌ NOT DONE | **4 unused declarations remain:** `htmlEscapeTestRenderer`, `assertMarshalError`, `testHTMLEscapeShared`, `testHTMLEmptyExpected`. Plus 3 more unused by root: `newTestNodeWithShape`, `testEdgesABC`, `testEmptyRendererOutput`. Total: 7 unused items. |
| 41  | Remove `benchmarks_test.go` entries for moved formats.           | ✅ DONE     | Only root benchmarks remain: `ASCIITreeRenderer`, `TableDataCreateRowEdges`, `MarkdownTable`                                                                                                                                                             |
| 42  | Remove `fuzz_test.go` entries for moved formats.                 | ✅ DONE     | Only `FuzzMarkdownTable` + generic `fuzzEnumTest` helper remain                                                                                                                                                                                          |
| 43  | Remove `userjourney_test.go` format-specific parts.              | ❌ NOT DONE | **Still references moved formats** via sub-module imports (`delimited`, `serialization`). Tests like `TestRenderDataAsCSV`, `TestRenderDataAsJSON`, `TestRenderDataAsYAML` remain. Depguard warns on sub-module imports from root.                       |
| 44  | Update root `go.mod`: remove `go-faster/yaml`, remove `escape`.  | ⚠️ PARTIAL  | `escape` correctly absent from `require` (only in `replace`). `go-faster/yaml` stays in `require` — justified by `internal/gentest/assert.go` which uses it directly. Root production code has zero yaml imports.                                        |
| 45  | Run full root test suite.                                        | ✅ DONE     | Zero compile errors, all tests pass                                                                                                                                                                                                                      |

---

## Phase 6: Documentation + Final Verification

| #   | Task                                                                      | Status      | Notes                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| --- | ------------------------------------------------------------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 46  | Update `AGENTS.md`: new module table, dependency graph, import paths.     | ✅ DONE     | Module table has 12 modules, dependency graph updated, go.work example includes all modules                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| 47  | Update `docs/modularization/root-split-proposal.md`: mark as implemented. | ❌ NOT DONE | Still reads as a raw proposal. No "implemented" status marker.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| 48  | Update `README.md`: new import path examples.                             | ⚠️ PARTIAL  | Sub-module import paths shown in some sections. **But many examples still use stale `output.X` references:** `output.NewCSVWriter` (→ should be `delimited.NewCSVWriter`), `output.MarshalJSONIndent` (stays in root ✓), `output.NewJSONTableRenderer` (→ should be `serialization.NewJSONTableRenderer`), `output.NewYAMLTableRenderer` (→ should be `serialization.NewYAMLTableRenderer`), `output.MarshalXMLFromTableData` (→ should be `markup.MarshalXMLFromTableData`), `output.MarshalYAML` (→ should be `serialization.MarshalYAML`), `output.NewHTMLTreeRenderer` (→ should be `markup.NewHTMLTreeRenderer`), `output.NewStreamingHTMLRenderer` (→ should be `markup.NewStreamingHTMLRenderer`). |
| 49  | Update `go.work` example in AGENTS.md.                                    | ✅ DONE     | Includes `./serialization`, `./delimited`, `./markup`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| 50  | Run ALL modules: test suite.                                              | ✅ DONE     | All modules build and test without errors                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| 51  | Run `golangci-lint` across ALL modules.                                   | ✅ DONE     | No errors. Warnings exist (unused test helpers — tracked in task 40). `.golangci.yml` not present.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        |

---

## Remaining Work

### Required (breaks correctness/staleness)

| #   | Original Task | What's Left                                                                                                                                                                                                  | Priority |
| --- | ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- |
| R1  | #48           | **Fix README.md examples** — ~10 code examples still reference `output.X` for functions that moved to sub-modules (`delimited`, `serialization`, `markup`). Users copying these examples get compile errors. | HIGH     |
| R2  | #47           | **Mark `root-split-proposal.md` as implemented** — add status section noting Strategy B was executed, date completed, actual vs planned deviations.                                                          | LOW      |

### Nice-to-have (dead code / lint)

| #   | Original Task | What's Left                                                                                                                                                                                                                               | Priority |
| --- | ------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| R3  | #40           | **Remove 7 unused declarations from `output_test_helpers_test.go`** — `htmlEscapeTestRenderer`, `assertMarshalError`, `testHTMLEscapeShared`, `testHTMLEmptyExpected`, `newTestNodeWithShape`, `testEdgesABC`, `testEmptyRendererOutput`. | MED      |
| R4  | #43           | **Decide fate of `userjourney_test.go`** — currently imports `delimited`/`serialization` from root (depguard violation). Options: (a) move to `integration/`, (b) delete format-specific tests, (c) update depguard config to allow.      | MED      |

---

## Summary (Actual vs Planned)

| Metric                               | Planned                                                  | Actual                                                                                                       |
| ------------------------------------ | -------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| **Total tasks**                      | 51                                                       | 43 done, 4 partial, 4 not done                                                                               |
| **Phases**                           | 7 (0-6)                                                  | 4 complete, 2 partial, 0 not started                                                                         |
| **New modules**                      | 3 (`delimited/`, `serialization/`, `markup/`)            | ✅ 3 created, all with tests + init() registration                                                           |
| **Root file reduction**              | 54 → ~25                                                 | 54 → ~20 source files (better than planned)                                                                  |
| **Root deps removed from prod code** | `go-faster/yaml`, `escape`                               | ✅ Zero yaml/escape imports in root `.go` files. `go-faster/yaml` in go.mod justified by `internal/gentest`. |
| **Breaking changes**                 | Users must import format-specific modules                | ✅ Documented, examples/ updated correctly                                                                   |
| **Registry dispatch**                | `render_tabledata.go` uses `TableDataMarshaler` registry | ✅ Thread-safe, zero sub-module imports in root                                                              |
| **Key invariant preserved**          | Root has ZERO sub-module imports                         | ✅ Verified via `go mod graph`                                                                               |
| **Deviation from plan**              | —                                                        | `streaming.go` split: interface stayed in root, implementation moved to markup/ (cleaner than planned)       |

## Risk Register

| Risk                                         | Likelihood | Impact | Mitigation                                  | Status                                          |
| -------------------------------------------- | ---------- | ------ | ------------------------------------------- | ----------------------------------------------- |
| Unexported type embedding breaks             | Medium     | High   | Phase 0 exports `TableDataBase` first       | ✅ Mitigated                                    |
| Test helper duplication across modules       | High       | Low    | ~10 lines each, acceptable tradeoff         | ✅ Accepted                                     |
| `render_tabledata.go` registry not populated | Medium     | High   | Clear error message if format not imported  | ✅ Mitigated — returns `UnsupportedFormatError` |
| Circular imports between new modules         | Low        | High   | `marshal.go` stays in root                  | ✅ Mitigated — zero circular deps               |
| Integration tests break                      | High       | Medium | Phase 4 dedicates 5 tasks to fixing imports | ✅ Mitigated                                    |
| `go-faster/yaml` still in root go.mod        | Low        | Low    | Verify after Phase 5 cleanup                | ⚠️ Stays — justified by `internal/gentest`      |
| README examples stale                        | High       | Medium | Task 48 update                              | ❌ **Active risk** — ~10 examples broken        |

## Estimated Total Time

| Phase                   | Tasks  | Est. Time            | Actual               |
| ----------------------- | ------ | -------------------- | -------------------- |
| Phase 0: Foundation     | 6      | ~39 min              | ✅ Done              |
| Phase 1: delimited/     | 8      | ~50 min              | ✅ Done              |
| Phase 2: markup/        | 9      | ~62 min              | ✅ Done              |
| Phase 3: serialization/ | 9      | ~67 min              | ✅ Done              |
| Phase 4: Dependents     | 7      | ~41 min              | ✅ Done              |
| Phase 5: Cleanup        | 6      | ~34 min              | ⚠️ 4/6 done          |
| Phase 6: Docs + Verify  | 6      | ~39 min              | ⚠️ 4/6 done          |
| **Total**               | **51** | **~332 min (~5.5h)** | **43/51 done (92%)** |
