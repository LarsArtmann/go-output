# TODO_LIST.md — go-output

**Last updated:** 2026-06-15
**Open items:** 7
**Blocked:** 1 (needs owner decision)

All items below are **verified open** against the current code (2026-06-15).
Lint status: **ZERO issues** across all 17 modules.

---

## P1 — Production Code (non-breaking)

| #   | Task                                                                                                                                                                               | Effort | Status |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 2   | **`nom/` internal decomposition** — split subscriber/tree/cache/render into `internal/` sub-packages keeping a thin public API. Improves locality & navigability (35 files today). | High   | Open   |

---

## P3 — Naming & Polish (low value)

| #   | Task                                                                                     | Effort | Status |
| --- | ---------------------------------------------------------------------------------------- | ------ | ------ |
| 7   | **`nom`: rename `GetOperationSymbol` → `OperationSymbol`**                               | Low    | Open   |
| 8   | **`examples/shared`: rename `HandleError` → `Must`**                                     | Low    | Open   |
| 9   | **Make `ColorModeAuto.ShouldColor()` deterministic/testable** — reads env+TTY at runtime | Low    | Open   |
| 10  | **Cache `TreeNode.Depth()`** — currently O(n) parent-chain walk                          | Low    | Open   |
| 11  | **Add bounds validation for `D2NodeStyle.Opacity`**                                      | Low    | Open   |

---

## P4 — Deferred (post-v1, blocked by ADR 006 API freeze)

| #   | Task                                                                                                                                  | Why deferred |
| --- | ------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| 12  | **Unify `Marshaler` → `Renderer` terminology** — registry types `TableDataMarshaler`/`AnyDataMarshaler` vs everything-else `Renderer` | Breaking     |
| 13  | **`RenderOptions.GraphID` is dead code** — no marshaler reads it (documented; wire or remove at v1)                                   | Breaking     |

### Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                                          |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only. Option B: unexported + validated setters (impossible-state-proof). Option C: keep both for v0.x. Affects every consumer; v1 stability commitment. |

### Release / Community

| #   | Task                                                                         | Effort | Status |
| --- | ---------------------------------------------------------------------------- | ------ | ------ |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                        | Low    | Open   |
| 16  | **Cut `v1.0.0` tag** — API declared frozen/ready (ADR 006), still at v0.10.x | Low    | Open   |

---

## Resolved This Session (2026-06-15) — Do Not Redo

| Task                                                          | Resolution                                                               |
| ------------------------------------------------------------- | ------------------------------------------------------------------------ |
| depguard config not portable (24 false positives)             | **Fixed** — `files` filter on `main` rule + missing deps                 |
| `nom/inline_renderer.go` write errors + nestif                | **Fixed** — `write`/`writef` helpers, flattened status block             |
| `nom/tree_render.go` perfsprint                               | **Fixed** — `fmt.Sprintf` → concatenation                                |
| `nom` embedded field ordering + whitespace                    | **Fixed** — `DisplayState` first, wsl_v5 auto-fixed                      |
| `graph/registry_test.go` DOT/Mermaid duplication              | **Fixed** — table-driven `TestRenderGraphTableData`                      |
| `tui/colors.go`: 10 mutable globals → `terminalColors` struct | **Fixed** — cohesive theme struct (was gochecknoglobals #1)              |
| `tui/colors.go` dead `colorCyan` global                       | **Removed**                                                              |
| `tui` whitespace (view.go, reporter.go, tests)                | **Fixed** — auto-fixed via `golangci-lint --fix`                         |
| Deprecated `EnsureBuild()` in tui tests (3 SA1019)            | **Fixed** → `GetRootNodes()`                                             |
| forcetypeassert in tui tests (3)                              | **Excluded** from test files via config (standard Go practice)           |
| Test errcheck/err113/gosec G104 sweep                         | **Excluded** from `_test.go` via config                                  |
| All 23 remaining wsl_v5 whitespace issues                     | **Fixed** — auto-fixed via `golangci-lint --fix`                         |
| `nom/activity_display_test.go` `copy` builtin shadow          | **Fixed** → `copied`                                                     |
| **BDD module added** (`bdd/`) — 19 Ginkgo specs               | **Done** — wired into flake.nix + go.work.example                        |
| `#11` pre-commit structure-linter                             | **Done** — `.structure-linter.yml` skips `root-package-files`            |
| `#12` gomod2nix                                               | **Closed** — decided against (AGENTS.md)                                 |
| `FEATURES.md` missing nom/tui                                 | **Fixed** — added Progress Visualization sections                        |
| `DOMAIN_LANGUAGE.md` GraphRendererMixin + missing nom/tui     | **Fixed** → GraphRendererState + new contexts                            |
| `RenderOptions.GraphID` dead code                             | **Documented** — clarified it's unwired, directs users to `SetGraphID()` |
| Integration dep sync (stale nom/lipgloss versions)            | **Fixed** — `go work sync`                                               |
| **ZERO lint achieved** (was 118 at session start)             | **Done** — 0 issues across all 17 modules                                |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-08 Architecture & Naming Sprint:** SlugifyID extracted, GraphRendererMixin→State, TableDataBase→Store, DTO suffix removed, formatCapabilities inverted, HTML table unified, html/template adopted, NodesPtr/EdgesPtr→AddNode/AddEdge, escape.D2/MermaidText optimized, AsciiDoc escaping completed, lipgloss style cached, D2 writeClasses sorted, D2ArrowNone added, FormatJSON registered, RenderTableData variadic→single RenderOptions.
- **2026-05-28 Round 6:** Footer row feature, pre-v1 API stability audit (ADR 006, 228 symbols frozen), round-trip integration tests, root coverage 82%→96%.
- **2026-05-25 Modularization:** D2/graph/table/delimited/serialization/markup/plantuml extraction, JSONL+AsciiDoc+TOML+PlantUML added (12→16 formats), Shape capability matrix (ADR 002), zero transitive deps, deduplication to 0 actionable clones.
- **FormatCategory + OutputFormat aliases removed** (verified gone 2026-06-15).
- Branded ID phantom types, ColorMode wiring, streaming writers, registry dispatch, Nix flake.
