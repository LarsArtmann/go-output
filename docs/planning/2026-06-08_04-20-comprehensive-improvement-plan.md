# Comprehensive Improvement Plan — go-output

**Date:** 2026-06-08\
**Author:** Post-audit reflection\
**Method:** Pareto prioritization (1% → 51%, 4% → 64%, 20% → 80%)

---

## Part 1: Honest Reflection — What Was Forgotten

### What I forgot in the first audit pass:

1. **No tests for my 4 critical fixes** — D2 determinism, D2ArrowType, JSON registry, and TableData nil safety all have zero tests verifying the fix. Coverage dropped from 96.1% → 94.7% on root because new nil-safe code paths are untested.
2. **Never ran `go test -race` on sub-modules** — Only ran root. Formatter modules also need race verification.
3. **Variadic `RenderOptions` anti-pattern still present** — HIGH issue #4 from code review, unfixed.
4. **`NodesPtr/EdgesPtr` encapsulation break still present** — HIGH issue #6, unfixed.
5. **Auto-generated `doc.go` files are garbage** — Pre-commit created "Package X provides ..." stubs. Worse than nothing.
6. **Integration tests don't test JSON via `RenderTableData`** — `renderJSONFormat` in integration uses `serialization.MarshalJSON` directly.
7. **`CODE_OF_CONDUCT.md` was auto-deleted** without me noticing during the pre-commit pass.
8. **Didn't check existing code before implementing** — e.g., `strings.NewReplacer` (stdlib) already exists and would optimize `escape.D2`.
9. **Didn't think about type model improvements before editing** — Adding nil checks to `TableData` is a band-aid; the real question is whether `TableData` should allow nil at all.
10. **AGENTS.md not updated** with new findings from the audit.

### What could be better:

- I stopped at "CRITICAL + HIGH" fixes but left MEDIUM issues untouched. Many MEDIUM issues are 5-minute fixes with real value.
- I generated reports but didn't turn findings into concrete TODOs in `TODO_LIST.md`.
- I should have asked: "does this fix need a test?" before every edit.

### What could still improve (from reports):

- **Type models:** `RenderOptions` variadic is wrong abstraction. `GraphRendererMixin` leaks mutable state. `TableData` has both exported fields and getters (design tension).
- **Performance:** `escape.D2` chains 4 `strings.ReplaceAll` → 4 allocations per call. `table.buildStyleFunc` allocates `lipgloss.NewStyle()` per row. `MermaidRenderer` uses `fmt.Fprintf` per node.
- **Escaping:** AsciiDoc and Mermaid escaping is incomplete.
- **Architecture:** `formatCapabilities` hardcoded in core. `HTMLRenderer` and `StreamingHTMLRenderer` duplicate generation.

---

## Part 2: Execution Plan

### Tier 1 — 1% Effort → 51% Impact (Do First)

| #   | Task                                           | Files                                                                           | Effort | Impact | Rationale                                                      |
| --- | ---------------------------------------------- | ------------------------------------------------------------------------------- | ------ | ------ | -------------------------------------------------------------- |
| 1.1 | Add nil-safety tests for `TableData`           | `tabledata_test.go`                                                             | 10 min | High   | Coverage dropped to 94.7%; restore to 96%+                     |
| 1.2 | Add D2 deterministic classes test              | `d2/d2_render_test.go`                                                          | 10 min | High   | Verify writeClasses sort fix works                             |
| 1.3 | Add D2ArrowType empty-string parse test        | `d2/d2_enum_test.go`                                                            | 5 min  | High   | Verify ParseD2ArrowType("") now returns D2ArrowNone            |
| 1.4 | Add JSON registry integration test             | `integration/render_helpers_test.go` + `integration_test.go`                    | 15 min | High   | Verify RenderTableData works with FormatJSON                   |
| 1.5 | Add race test for `RegisterTableDataMarshaler` | `render_tabledata_test.go`                                                      | 10 min | High   | `sync.RWMutex` is untested under concurrency                   |
| 1.6 | Fix variadic `RenderOptions` → single value    | `render_tabledata.go`, all callers                                              | 15 min | High   | Removes misleading API; every caller already passes 0 or 1 arg |
| 1.7 | Fix garbage `doc.go` files with real content   | `integration/doc.go`, `internal/gentest/doc.go`, `testhelpers/graphtest/doc.go` | 10 min | Medium | Replace "provides ..." with actual package descriptions        |

### Tier 2 — 4% Effort → 64% Impact (Do Next)

| #   | Task                                                       | Files                                                | Effort | Impact | Rationale                                              |
| --- | ---------------------------------------------------------- | ---------------------------------------------------- | ------ | ------ | ------------------------------------------------------ |
| 2.1 | Remove `NodesPtr/EdgesPtr`, replace with `AddNode/AddEdge` | `graph.go`, `graph_mixin_test.go`, callers in graph/ | 20 min | High   | Breaks encapsulation; callers use it for mutation      |
| 2.2 | Optimize `escape.D2` with `strings.NewReplacer`            | `escape/escape.go`                                   | 5 min  | Medium | Single-pass replace, 1 allocation instead of 4         |
| 2.3 | Optimize `escape.MermaidText` with `strings.NewReplacer`   | `escape/escape.go`                                   | 5 min  | Medium | Same optimization                                      |
| 2.4 | Cache `lipgloss.NewStyle()` in `table.buildStyleFunc`      | `table/table.go`                                     | 15 min | Medium | Eliminates per-row allocation in hot path              |
| 2.5 | Complete AsciiDoc escaping                                 | `markup/asciidoc.go`                                 | 10 min | Medium | Add `*`, `_`, `` ` ``, `~`, `^` escaping               |
| 2.6 | Extract shared `SlugifyID` helper                          | New file or `escape/escape.go`                       | 15 min | Medium | Remove 4 duplications of `ReplaceAll(label, " ", "_")` |

### Tier 3 — 20% Effort → 80% Impact (Do Later)

| #   | Task                                                | Files                                   | Effort | Impact | Rationale                                       |
| --- | --------------------------------------------------- | --------------------------------------- | ------ | ------ | ----------------------------------------------- |
| 3.1 | Rename `TableDataBase` → `TableDataStore`           | `tabledata.go`, all sub-modules         | 30 min | Low    | Name leaks implementation; breaking change      |
| 3.2 | Rename `GraphRendererMixin` → `GraphRendererState`  | `graph.go`, all graph renderers         | 30 min | Low    | Name leaks pattern; breaking change             |
| 3.3 | Remove `DTO` suffix from serialization types        | `serialization/*_dto.go`                | 20 min | Low    | Java-ism in Go code; internal only              |
| 3.4 | Invert `formatCapabilities` dependency              | `shape.go`, all sub-modules             | 45 min | High   | Sub-modules register shapes; cleaner seam       |
| 3.5 | Merge HTMLRenderer/StreamingHTMLRenderer generation | `markup/html.go`, `markup/streaming.go` | 30 min | High   | Single source of truth for HTML table structure |
| 3.6 | Inline `marshal.go` wrappers into `serialization/`  | `marshal.go`, `serialization/*.go`      | 20 min | Low    | Shallow wrappers; core shouldn't own marshaling |
| 3.7 | Use `html/template` for HTML generation             | `markup/html.go`, `markup/streaming.go` | 30 min | Medium | More robust than string concatenation           |

---

## Part 3: Type Model Improvements

### 1. `RenderOptions` — From Variadic to Concrete

**Current:** `RenderTableData(data *TableData, format Format, opts ...RenderOptions) error`\
**Problem:** Only `opts[0]` is used. Variadic implies merge semantics that don't exist.\
**Better:** `RenderTableData(data *TableData, format Format, opts RenderOptions) error`\
**Every caller** already passes 0 or 1 arg. `RenderOptions{}` for default.

### 2. `GraphRendererMixin` — Stop Leaking Mutable State

**Current:** `NodesPtr() *[]GraphNode` / `EdgesPtr() *[]GraphEdge`\
**Problem:** Callers can mutate internal slices directly, bypassing interface contracts.\
**Better:** Remove `NodesPtr/EdgesPtr`. Add `AddNode(node GraphNode)` and `AddEdge(edge GraphEdge)` to the mixin. Callers use controlled mutation.

### 3. `TableData` — Exported Fields + Getters = Tension

**Current:** `Headers`, `Rows`, `Footer` are exported fields. `GetHeaders()`, `GetRows()`, `GetFooter()` are getters.\
**Problem:** Both exist. Which is canonical? Callers can mutate `data.Headers = ...` directly, bypassing any future validation.\
**Better:** Either make fields unexported and getters the only access (breaking change), OR document that direct mutation is the intended API and remove getters (also breaking). For v0.x, document the tension and plan for v1.

### 4. `TableDataMarshaler` — Function Type is Fine

**Current:** `type TableDataMarshaler func(w io.Writer, data *TableData, opts RenderOptions) error`\
**Assessment:** Function type is correct here. An interface would add indirection without benefit. Keep as-is.

---

## Part 4: Library Opportunities (Well-Established)

| Library                        | Already Used?    | Could Use For                                           | Effort       | Impact         |
| ------------------------------ | ---------------- | ------------------------------------------------------- | ------------ | -------------- |
| `strings.NewReplacer` (stdlib) | No               | `escape.D2`, `escape.MermaidText`, `escape.MermaidSlug` | 10 min       | Medium         |
| `html/template` (stdlib)       | No               | `markup/html.go`, `markup/streaming.go`                 | 30 min       | Medium         |
| `sort.Strings` (stdlib)        | Yes (just added) | `d2/d2_render.go` writeClasses                          | Done         | High           |
| `sync.RWMutex` (stdlib)        | Yes              | `render_tabledata.go` registry                          | Already used | —              |
| `github.com/puzpuzpuz/xsync`   | No               | Lock-free registry map                                  | 20 min       | Low (overkill) |

**Verdict:** `strings.NewReplacer` and `html/template` are clear wins. `xsync` is overkill for a write-once-read-many map populated at init().

---

## Part 5: Verification Checklist Per Task

After every task:

- [ ] `go test ./...` passes in affected module(s)
- [ ] `go test -race ./...` passes in affected module(s)
- [ ] `golangci-lint run ./...` passes in affected module(s)
- [ ] Coverage not regressed (check `go test -cover`)
- [ ] `git status && git diff` reviewed
- [ ] `git commit` with detailed message

After all tasks:

- [ ] `nix flake check` passes
- [ ] `go test -race ./...` passes in all 14 modules
- [ ] `golangci-lint run ./...` passes in all 14 modules
- [ ] `git push`
