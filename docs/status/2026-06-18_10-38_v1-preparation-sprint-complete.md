# Status Report: v1 Preparation Sprint Complete

**Date:** 2026-06-18 10:38
**Branch:** master
**Commits since v0.12.0:** 41

---

## A. FULLY DONE (11 TODO items resolved this sprint)

| #     | Item                                                   | Impact                                                                                                                                     | Files Changed                                                                                                     |
| ----- | ------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------- |
| #22   | TUI reporter data race eliminated                      | Critical — reporter owns mutex-protected workflowState; all model mutations via `send()`→`model.Update()`; 3 concurrent stress tests added | `tui/reporter.go`, `tui/model.go`, `tui/messages.go`, `tui/lifecycle.go`, `tui/reporter_concurrent_test.go` (new) |
| #13   | `RenderOptions.GraphID` dead code removed              | Medium — zero callers ever read it; clean for v1                                                                                           | `render_tabledata.go`                                                                                             |
| #27   | `EdgeStyle.Style` → typed `LineStyle` enum             | Medium — solid/dashed/dotted now compile-time safe; field renamed `Style`→`Line`                                                           | `graph.go`, `graph/dot.go`, `graph/graph_test.go`                                                                 |
| #28   | `GetActivityCounts()` → `ActivityCounts` struct        | High — eliminates arg-swap bugs; named fields + `Total()` method                                                                           | 13 files across nom, tui, integration, examples                                                                   |
| #9    | `ColorModeAuto.ShouldColor()` made testable            | Medium — detection functions are overridable vars; deterministic test covering all 4 combinations                                          | `color.go`, `color_test.go`                                                                                       |
| M5    | `GraphShape` → `NodeShape` disambiguation              | Medium — constants no longer collide with `ShapeTable`/`ShapeTree`/`ShapeGraph`                                                            | 12 files across root, graph, d2, serialization, testhelpers                                                       |
| M8    | `GraphStyle.FillColor`/`StrokeColor` → `Fill`/`Stroke` | Medium — aligned with `D2NodeStyle` naming; converter is near 1:1 copy                                                                     | 11 files across root, graph, d2, plantuml, examples                                                               |
| M6/M7 | `output.Direction` enum + branded-ID doc               | Medium — bridges D2 ("down"/"right") ↔ DOT ("TB"/"LR"); canonical import path documented                                                   | `direction.go` (new), `ids.go`, `d2/d2_enum.go`                                                                   |
| #12   | `Marshaler` → `Renderer` terminology unification       | Low — `TableDataMarshaler`→`TableDataRenderer`, `AnyDataMarshaler`→`AnyDataRenderer`; 51 references across all sub-modules                 | 17 files                                                                                                          |
| #8    | Closed as won't-fix                                    | —                                                                                                                                          | `HandleError` is honest; `Must` implies panic (wrong for `os.Exit`)                                               |
| Lint  | `gochecknoglobals` on `lineStyleValues`                | Low — added nolint directive matching existing pattern                                                                                     | `graph.go`                                                                                                        |

**Verification:** All 16 modules pass build ✓, test ✓, lint 0 issues ✓.

---

## B. PARTIALLY DONE

Nothing. All started items are fully completed.

---

## C. NOT STARTED (3 open items)

| #   | Item                                                                            | Effort | Why deferred                                                                                                                      |
| --- | ------------------------------------------------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------- |
| M4  | Rename `InlineRenderer.Render()`→`Draw()`, `DependencyTree.Render()`→`Format()` | Medium | Pure naming refactor, low impact, many call sites across nom/tui/examples. Last remaining "deferred" item from split-brain audit. |
| #14 | Post to r/golang, submit to Awesome Go                                          | Low    | Needs owner's Reddit/GitHub account                                                                                               |
| #16 | Cut v1.0.0 tag                                                                  | Low    | Needs owner's explicit go-ahead                                                                                                   |

---

## D. TOTALLY FUCKED UP

Nothing. Zero regressions. All tests pass. Zero lint issues across all 16 modules.

**One pre-existing issue found (not caused by this sprint):**

- `~/.cache/nom-timing.csv` can become corrupted with wrong field counts (integration test failure). Trashing the file fixes it. Root cause: concurrent writes to the cache from multiple test processes. The `saveMu` fix in the nom module serializes writes within a single process, but cross-process corruption is still possible. This should be addressed by using a unique temp path per test process.

---

## E. WHAT WE SHOULD IMPROVE

1. **M4 Render→Draw/Format rename** is the last split-brain residue. Should be done before v1.0.0 freeze.
2. **`go.work` is gitignored** — new contributors must run `nix run .#setup-workspace` or manually create it. Consider committing `go.work` or making setup more discoverable.
3. **Nom timing cache path** should be per-process in tests to prevent cross-test corruption.
4. **ADR 006 needs updating** — the frozen API symbol list has changed (renamed types, removed GraphID, new Direction/LineStyle/ActivityCounts types).
5. **CHANGELOG.md** needs a v1.0.0 section documenting all breaking changes.
6. **Stale gopls diagnostics** — 107 phantom LSP errors persist from renames that the actual build doesn't have. Needs `gopls restart`.
7. **Documentation drift** — `docs/FORMAT_ARCHITECTURE.md`, `docs/DOMAIN_LANGUAGE.md`, `README.md`, and `FEATURES.md` still reference old names (`GraphShape`, `FillColor`, `Marshaler`, `ShapeBox`).

---

## F. TOP 25 THINGS TO DO NEXT (sorted by impact/effort)

| #   | Task                                                                                        | Impact   | Effort |
| --- | ------------------------------------------------------------------------------------------- | -------- | ------ |
| 1   | **Update docs** — FORMAT_ARCHITECTURE.md, DOMAIN_LANGUAGE.md, README.md with new type names | High     | Low    |
| 2   | **Update CHANGELOG.md** with v1.0.0 breaking changes section                                | High     | Low    |
| 3   | **Update FEATURES.md** with new types (Direction, LineStyle, ActivityCounts, NodeShape)     | High     | Low    |
| 4   | **M4: Rename Render→Draw/Format** in nom module — last split-brain item                     | Medium   | Medium |
| 5   | **Update AGENTS.md** — reflect new type names, removed fields, new patterns                 | Medium   | Low    |
| 6   | **Update ADR 006** — regenerate frozen symbol list with new names                           | Medium   | Medium |
| 7   | **Fix nom timing cache test isolation** — use t.TempDir() for cache path                    | Medium   | Low    |
| 8   | **Add Direction unit tests** — test ToD2Direction/ToRankDir for all 4 values                | Medium   | Low    |
| 9   | **Add LineStyle unit tests** — test IsValid/String for all values                           | Low      | Low    |
| 10  | **Add NodeShape test coverage** — verify ParseNodeShape works after rename                  | Low      | Low    |
| 11  | **Rename `NOTE(split-brain)` markers** that remain in code (M4 still has one)               | Low      | Low    |
| 12  | **Coverage check** — verify all changed modules still ≥90%                                  | Medium   | Low    |
| 13  | **Govulncheck** — run across all modules before v1.0.0                                      | High     | Low    |
| 14  | **Cut v1.0.0 tag** — after docs updated and ADR refreshed                                   | Critical | Low    |
| 15  | **Submit to Awesome Go** — draft submission text                                            | Medium   | Low    |
| 16  | **Post to r/golang** — draft show-and-tell post                                             | Medium   | Low    |
| 17  | **Add deprecation guide** — migration doc for v0.12→v1.0 users                              | Medium   | Medium |
| 18  | **Integration test for Direction** — verify D2 and DOT renderers accept Direction           | Medium   | Medium |
| 19  | **Benchmark impact check** — verify renames didn't introduce allocations                    | Low      | Medium |
| 20  | **Godoc review** — verify all new exported types have proper doc comments                   | Medium   | Low    |
| 21  | **Example refresh** — update examples to use new type names (LineStyle, NodeShape, etc.)    | Medium   | Low    |
| 22  | **BDD test refresh** — verify bdd module uses new type names                                | Low      | Low    |
| 23  | **Dependency audit** — verify no new deps were added                                        | Low      | Low    |
| 24  | **README badge update** — coverage %, version                                               | Low      | Low    |
| 25  | **GitHub release notes** — draft for v1.0.0 tag                                             | Medium   | Low    |

---

## G. TOP QUESTION I CANNOT FIGURE OUT MYSELF

**#15: Should `TableData` use exported fields or getters for v1?**

This is the single remaining blocked item. Three valid approaches with big tradeoffs:

- **Option A (fields only):** Remove all `GetHeaders()`, `GetData()`, `GetFooter()` methods. Users access `data.Headers` directly. Simplest, most Go-idiomatic, but no validation on mutation.
- **Option B (getters only):** Unexport `Headers`/`Data`/`Footer`, keep getters, add validated setters. Strongest invariant enforcement, but verbose and breaks every consumer.
- **Option C (keep both):** Current state. Backwards compatible, but redundant API surface for v1.

This affects every consumer of the library. It's a v1 stability commitment. **Which option do you want?**

---

## Sprint Summary

| Metric           | Value                                                              |
| ---------------- | ------------------------------------------------------------------ |
| Items resolved   | 11                                                                 |
| Commits          | 10                                                                 |
| Files changed    | 40+                                                                |
| Modules affected | 16 (all)                                                           |
| Breaking changes | 7 (clean break for v1.0.0)                                         |
| New types added  | 4 (`Direction`, `LineStyle`, `ActivityCounts`, `NodeShape` rename) |
| Tests added      | 4 (3 concurrent + 1 deterministic color)                           |
| Lint issues      | 0 across all 16 modules                                            |
| Build status     | ✓ all modules                                                      |
| Test status      | ✓ all modules                                                      |
| Race detector    | ✓ nom + tui                                                        |
