# TODO_LIST.md — go-output

**Last updated:** 2026-06-15
**Open items:** 13
**Blocked:** 1 (needs owner decision)

All items below are **verified open** against the current code (2026-06-15). Completed/stale items are filtered out. Items resolved this session are listed at the bottom — do not re-do them.

---

## P1 — Production Code (non-breaking)

| #   | Task                                                                                                                                                                               | Effort | Status |
| --- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 1   | **`tui/colors.go`: convert 10 mutable package globals → immutable style struct** — 11th global (`colorCyan`) was dead and removed. Remaining 10 flagged by `gochecknoglobals`.     | Medium | Open   |
| 2   | **`nom/` internal decomposition** — split subscriber/tree/cache/render into `internal/` sub-packages keeping a thin public API. Improves locality & navigability (35 files today). | High   | Open   |

---

## P2 — Test Quality (low risk)

| #   | Task                                                                                                                           | Effort | Status |
| --- | ------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 3   | **`tui/model_test.go`: replace deprecated `EnsureBuild()` → `GetRootNodes()`/`VisibleNodes()`** (3 SA1019)                     | Low    | Open   |
| 4   | **`tui/event_sequence_test.go`: guard 3 unchecked type assertions** (forcetypeassert)                                          | Low    | Open   |
| 5   | **Test errcheck sweep** — 31 unchecked `AddActivity`/`Record`/`OnEvent`/`UpdateActivityStatus` returns (nom, tui, integration) | Medium | Open   |
| 6   | **err113 test sweep** — 7 dynamic `errors.New(...)` in tests → wrapped static errors (integration, testhelpers, tui)           | Low    | Open   |

---

## P3 — Naming & Polish (low value)

| #   | Task                                                                                               | Effort | Status |
| --- | -------------------------------------------------------------------------------------------------- | ------ | ------ |
| 7   | **`nom`: rename `GetOperationSymbol` → `OperationSymbol`**                                         | Low    | Open   |
| 8   | **`examples/shared`: rename `HandleError` → `Must`**                                               | Low    | Open   |
| 9   | **Make `ColorModeAuto.ShouldColor()` deterministic/testable** — currently reads env+TTY at runtime | Low    | Open   |
| 10  | **Cache `TreeNode.Depth()`** — currently O(n) parent-chain walk                                    | Low    | Open   |
| 11  | **Add bounds validation for `D2NodeStyle.Opacity`**                                                | Low    | Open   |

---

## P4 — Deferred (post-v1, blocked by ADR 006 API freeze)

| #   | Task                                                                                                                                                        | Why deferred |
| --- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| 12  | **Unify `Marshaler` → `Renderer` terminology** — registry types `TableDataMarshaler`/`AnyDataMarshaler` vs everything-else `Renderer`                       | Breaking     |
| 13  | **`RenderOptions.GraphID` is dead code** — no marshaler reads it. Implement its purpose or remove at next major version (verified: zero readers 2026-06-15) | Breaking     |

### Blocked — Needs Owner Decision

| #   | Question                                                                                                                                                                                                                                                                                                                                                                                                                   |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 15  | **Should `TableData` use exported fields or getters for v1?** Current: both exist (`Headers` + `GetHeaders()`). Option A: exported fields only. Option B: unexported + validated setters (impossible-state-proof). Option C: keep both for v0.x. Affects every consumer; v1 stability commitment. _(See also: improve-codebase-architecture deepening report — unexporting fields makes column-mismatch unrepresentable.)_ |

### Release / Community

| #   | Task                                                                         | Effort | Status |
| --- | ---------------------------------------------------------------------------- | ------ | ------ |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                        | Low    | Open   |
| 16  | **Cut `v1.0.0` tag** — API declared frozen/ready (ADR 006), still at v0.10.x | Low    | Open   |

---

## Resolved This Session (2026-06-15) — Do Not Redo

| Task                                                                                    | Resolution                                                                                                  |
| --------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| depguard config not portable (24 false positives)                                       | **Fixed** — added `files` filter to `main` rule + missing deps (log, charmbracelet/x/\*, runewidth, ginkgo) |
| `nom/inline_renderer.go` write errors + nestif                                          | **Fixed** — added best-effort `write`/`writef` helpers, flattened status block                              |
| `nom/tree_render.go` perfsprint                                                         | **Fixed** — `fmt.Sprintf` → concatenation                                                                   |
| `nom` embedded field ordering (activity_display.go, tree.go)                            | **Fixed** — moved `DisplayState` first                                                                      |
| `nom` wsl_v5 whitespace (inline_renderer, subscriber_handlers, tree_render)             | **Fixed**                                                                                                   |
| `graph/registry_test.go` DOT/Mermaid duplication                                        | **Fixed** — table-driven `TestRenderGraphTableData`                                                         |
| `tui/view.go` + `reporter.go` wsl_v5 whitespace                                         | **Fixed**                                                                                                   |
| `tui/colors.go` dead `colorCyan` global                                                 | **Removed**                                                                                                 |
| `nom/activity_display_test.go` `copy` builtin shadow                                    | **Fixed** → `copied`                                                                                        |
| **BDD module added** (`bdd/`) — 19 Ginkgo specs, wired into flake.nix + go.work.example | **Done**                                                                                                    |
| `#11` pre-commit structure-linter                                                       | **Done** — `.structure-linter.yml` skips `root-package-files`                                               |
| `#12` gomod2nix                                                                         | **Closed** — decided against (AGENTS.md)                                                                    |
| `FEATURES.md` missing nom/tui                                                           | **Fixed** — added Progress Visualization sections                                                           |
| `DOMAIN_LANGUAGE.md` GraphRendererMixin + missing nom/tui                               | **Fixed** → GraphRendererState + new contexts                                                               |

---

## Completed Earlier (Summary — details in git history)

- **2026-06-08 Architecture & Naming Sprint:** SlugifyID extracted, GraphRendererMixin→State, TableDataBase→Store, DTO suffix removed, formatCapabilities inverted, HTML table unified, html/template adopted, NodesPtr/EdgesPtr→AddNode/AddEdge, escape.D2/MermaidText optimized with NewReplacer, AsciiDoc escaping completed, lipgloss style cached, D2 writeClasses sorted, D2ArrowNone added, FormatJSON registered, RenderTableData variadic→single RenderOptions.
- **2026-05-28 Round 6:** Footer row feature, pre-v1 API stability audit (ADR 006, 228 symbols frozen), round-trip integration tests, root coverage 82%→96%.
- **2026-05-25 Modularization:** D2/graph/table/delimited/serialization/markup/plantuml extraction, JSONL+AsciiDoc+TOML+PlantUML added (12→16 formats), Shape capability matrix (ADR 002), zero transitive deps, deduplication to 0 actionable clones.
- **FormatCategory + OutputFormat aliases removed** (verified gone 2026-06-15).
- Branded ID phantom types, ColorMode wiring, streaming writers, registry dispatch, Nix flake.
