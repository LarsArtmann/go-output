# TODO_LIST.md — go-output

**Last updated:** 2026-06-08
**Open items:** 13
**Blocked:** 1 (needs owner decision)

---

## P0 — Bugs & Latent Issues

| #   | Task                                                                                                                                                                                                                                                                                         | Effort | Status |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 1   | **Extract `escape.SlugifyID()` and unify all 4 call sites** — `escape.MermaidSlug` (complete), `plantuml.sanitizePlantUMLID` (missing `/`), `graph.dotTreeNodeID` (missing `-`, `/`), `d2.treeNodeID` (missing `-`, `/`). Inconsistent sanitization is a latent bug in DOT/D2 tree node IDs. | 15 min | Open   |
| 2   | **Add race test for `RegisterTableDataMarshaler`** — `sync.RWMutex` in `render_tabledata.go` is untested under concurrent access. Thread safety is assumed but not proven.                                                                                                                   | 10 min | Open   |

---

## P1 — Architecture (Breaking Changes)

| #   | Task                                                                                                                                                                              | Effort | Status |
| --- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 4   | **Invert `formatCapabilities` dependency** — Sub-modules should register their own shapes via `init()` instead of hardcoding in `shape.go`. Biggest remaining architectural seam. | 45 min | Open   |
| 5   | **Merge HTMLRenderer/StreamingHTMLRenderer table generation** — Both produce identical HTML table structure via string concatenation. Extract shared writer.                      | 30 min | Open   |
| 6   | **Use `html/template` for HTML generation** — Replace string concatenation in `markup/html.go` and `markup/streaming.go` with template engine for robust auto-escaping.           | 30 min | Open   |
| 7   | **Inline `marshal.go` wrappers into `serialization/`** — `MarshalFormat`, `UnmarshalFormat`, `MarshalJSONIndent` are shallow wrappers. Core shouldn't own marshaling.             | 20 min | Open   |

---

## P2 — Naming Cleanup (Breaking Changes)

| #   | Task                                                                                                                                     | Effort | Status |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 8   | **Rename `TableDataBase` → `TableDataStore`** — Current name leaks implementation detail ("Base" implies inheritance pattern).           | 15 min | Open   |
| 9   | **Rename `GraphRendererMixin` → `GraphRendererState`** — "Mixin" leaks composition pattern choice to callers.                            | 15 min | Open   |
| 10  | **Remove `DTO` suffix from serialization types** — `treeNodeDTO`, `graphDTO` etc. are Java-isms. Internal-only, so no external breakage. | 15 min | Open   |

---

## P3 — Build & Config

| #   | Task                                                                                                                                                                                                                                                       | Effort | Status |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 11  | **Fix pre-commit `--no-verify` requirement** — BuildFlow's `go-structure-linter` reports 29 false-positive "root-package-files" issues. Root package IS the public API for a Go library. Configure BuildFlow to ignore this rule, or accept `--no-verify`. | 15 min | Open   |
| 12  | **Add `gomod2nix` for reproducible Nix builds** — Nix sandbox blocks `go mod download`. Currently Go deps download at build time.                                                                                                                          | 30 min | Open   |
| 13  | **Investigate `go:generate stringer` for enums** — 7 hand-rolled enum types with identical Parse/IsValid/AllowedValues/String patterns. Code generation could eliminate boilerplate.                                                                       | 20 min | Open   |

---

## P4 — Future & Community

| #   | Task                                                  | Effort | Status |
| --- | ----------------------------------------------------- | ------ | ------ |
| 14  | **Community: Post to r/golang, submit to Awesome Go** | 30 min | Open   |

---

## Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                              | Why Blocked                                          |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only (Go-idiomatic, simpler). Option B: unexported fields + getters (controlled, future-proof). Option C: keep both for v0.x, decide at v1. | Affects every consumer. v1 API stability commitment. |

---

## Completed Work (Summary)

All items below are DONE. Do not re-do. Details in git history.

### 2026-06-08 — Comprehensive Audit Sprint

- D2 `writeClasses` sorted for deterministic output
- `D2ArrowNone` added to `D2ArrowType` values
- `FormatJSON` registered in `RenderTableData` dispatch
- `TableData` nil-receiver safety with comprehensive tests
- `RenderTableData` signature: variadic → single `RenderOptions` (BREAKING)
- `NodesPtr`/`EdgesPtr` removed → `AddNode`/`AddEdge` + `NodeEdgeAppender` interface
- `escape.D2` + `escape.MermaidText` optimized with `strings.NewReplacer`
- AsciiDoc escaping completed: `|`, `*`, `_`, `` ` ``, `~`, `^`
- `lipgloss.NewStyle()` cached in `table.buildStyleFunc`
- `doc.go` files rewritten with real package descriptions
- Missing `replace` directives in `delimited/go.mod`, `markup/go.mod`
- `flake.nix` `checks.format` moved to correct location
- `testhelpers/graphtest` added to CI workflow loops
- 5 research/audit reports + comprehensive status report

### 2026-05-28 — Round 6 Polish

- Footer row feature complete (TableData.Footer, Validate, WriteFooter, CSS, examples)
- Pre-v1 API stability audit (ADR 006, 228 exported symbols frozen)
- Round-trip integration tests for all 16 formats
- Root coverage 82% → 96%, D2 coverage → 100%

### 2026-05-25 — Modularization & Format Expansion

- D2/graph/table/delimited/serialization/markup/plantuml module extraction
- JSONL, AsciiDoc, TOML, PlantUML format additions (12 → 16 formats)
- Zero transitive deps from sub-modules users don't import
- Shape capability matrix (ADR 002)
- Code deduplication to 0 actionable clones

### Earlier

- Full test coverage across all 14 modules (90%+ each)
- Branded ID phantom types for type-safe IDs
- ColorMode wired into table, tree, markdown
- Streaming HTML/XML writers
- Registry-based TableData dispatch
- Nix flake with flake-parts + treefmt-nix
