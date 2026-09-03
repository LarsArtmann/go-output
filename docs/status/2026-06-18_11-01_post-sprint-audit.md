# Status Report: Post-Sprint Audit & Cleanup

**Date:** 2026-06-18 11:01
**Branch:** master
**Commits since v0.12.0:** 48
**LOC:** 31,333 across ~225 Go files in 16 modules

---

## A. FULLY DONE

### Code Quality (all verified green)

| Check                             | Result                                               |
| --------------------------------- | ---------------------------------------------------- |
| Build                             | ✅ All 16 modules compile                            |
| Tests                             | ✅ All 16 modules pass                               |
| Lint                              | ✅ 0 issues across all 16 modules (golangci-lint v2) |
| Race detector                     | ✅ nom + tui pass with -race                         |
| NOTE(split-brain) markers         | 2 remaining (M4 only — see section C)                |
| TODO/BUG/FIXME in production code | ✅ Zero                                              |
| Files over 350 lines (production) | ✅ Zero                                              |
| Old type names in .go files       | ✅ Zero                                              |

### Sprint Deliverables (20 commits since code review sprint)

**Critical fixes:**

- TUI reporter data race eliminated — all model mutations via `send()`→`model.Update()`, mutex-protected `workflowState` on reporter, 3 concurrent stress tests added (#22)
- `FormatDuration` boundary bug, `GetActivitySummaryString` slice formatting, DOT injection, escaping gaps (all from prior sprint)

**Type safety improvements:**

- `EdgeStyle.Style` (free-form string) → typed `LineStyle` enum with `LineStyleSolid`/`Dashed`/`Dotted` + `IsValid()`/`String()` (#27)
- `GetActivityCounts()` (4 unnamed ints) → `ActivityCounts` struct with named fields + `Total()` (#28)
- `ColorModeAuto.ShouldColor()` detection → overridable function vars for deterministic testing (#9)

**API alignment (breaking changes for v1.0.0):**

- `GraphShape` → `NodeShape`, `ShapeBox` → `NodeShapeBox` etc. — disambiguated from `ShapeTable`/`ShapeTree`/`ShapeGraph` (M5)
- `GraphStyle.FillColor` → `Fill`, `.StrokeColor` → `Stroke` — aligned with `D2NodeStyle` (M8)
- `TableDataMarshaler` → `TableDataRenderer`, `AnyDataMarshaler` → `AnyDataRenderer` — unified terminology (#12)
- `RenderOptions.GraphID` removed — dead code, zero callers (#13)
- `output.Direction` enum added — bridges D2 ("down"/"right") ↔ DOT ("TB"/"LR") (M7)
- D2NodeID canonical import path documented (M6)

**Test coverage additions:**

- Direction: `ToD2Direction()`, `ToRankDir()`, round-trip distinctness
- LineStyle: `IsValid()`, `String()` for all 3 values
- ActivityCounts: `Total()` for zero, single-status, mixed cases
- ColorMode: deterministic test for all 4 detection combinations
- TUI reporter: 3 concurrent stress tests (200 goroutines × 4 methods)

**Documentation updates:**

- 5 live docs updated (FORMAT_ARCHITECTURE.md, DOMAIN_LANGUAGE.md, AGENTS.md, README.md, FEATURES.md)
- CHANGELOG.md v1.0.0 breaking changes section added
- TODO_LIST.md rewritten: 3 open items (down from 21)
- Stale function/test names renamed (graphShapeToD2→nodeShapeToD2, TestParseGraphShape→TestParseNodeShape)

---

## B. PARTIALLY DONE

| Item                            | Status                                                                                                                                                                                          | What remains                                       |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| Stale unexported function names | 7 references to old "Marshaler" terminology in unexported funcs: `getTableDataMarshaler()`, `getAnyDataMarshaler()`, comments in `render_tabledata.go`, test name in `render_tabledata_test.go` | Mechanical rename to `getTableDataRenderer()` etc. |
| AGENTS.md design patterns       | Updated type names but patterns section still references some old concepts                                                                                                                      | Should review after M4 decision                    |
| ADR 006                         | Still references `RenderOptions.GraphID` and old type names                                                                                                                                     | Needs refresh for v1.0.0                           |

---

## C. NOT STARTED

| #   | Item                                                                            | Effort | Impact   | Notes                                                                                                                                         |
| --- | ------------------------------------------------------------------------------- | ------ | -------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| M4  | Rename `InlineRenderer.Render()`→`Draw()`, `DependencyTree.Render()`→`Format()` | Medium | Low      | Last split-brain item. 2 NOTE markers in `nom/inline_renderer.go:127` and `nom/tree_render.go:16`. Many call sites across nom, tui, examples. |
| #14 | Post to r/golang, submit to Awesome Go                                          | Low    | Medium   | Needs owner's Reddit/GitHub account                                                                                                           |
| #16 | Cut v1.0.0 tag                                                                  | Low    | Critical | Needs owner's explicit go-ahead                                                                                                               |
| #15 | TableData API: fields vs getters for v1                                         | —      | High     | **Blocked on owner decision**                                                                                                                 |

---

## D. TOTALLY FUCKED UP

**Nothing.** Zero regressions, zero data loss, zero broken builds.

**Pre-existing issues found (not caused by this sprint):**

1. **Nom timing cache cross-process corruption** — `~/.cache/nom-timing.csv` can become corrupted when multiple test processes write concurrently. The `saveMu` fix serializes within a single process but not across processes. Fix: use `t.TempDir()` for cache path in tests.
2. **3 test files over 350-line limit** — `integration/roundtrip_test.go` (528), `nom/subscriber_test.go` (523), `tui/event_sequence_test.go` (497). Production code has zero violations.

---

## E. WHAT WE SHOULD IMPROVE

1. **Finish M4** — It's the last split-brain residue. Two NOTE markers remain in nom module.
2. **Rename stale unexported functions** — `getTableDataMarshaler()` etc. still use old "Marshaler" terminology.
3. **Split oversized test files** — 3 test files exceed the 350-line limit.
4. **Fix nom timing cache test isolation** — Use per-test temp directory.
5. **Update ADR 006** — Frozen API symbol list is stale after all renames.
6. **Govulncheck** — Not yet run. Should be done before v1.0.0.
7. **Stale gopls diagnostics** — 96 phantom LSP errors from renames. Need `gopls restart` (attempted but failed).
8. **BDD test coverage** — bdd module wasn't updated for new type names (it may still work via interfaces but should be verified).
9. **Example refresh** — Examples compile but may show outdated API patterns.
10. **GitHub release notes** — Draft for v1.0.0 tag announcement.

---

## F. TOP 25 THINGS TO DO NEXT

| #  | Task                                                                                       | Impact   | Effort |
| -- | ------------------------------------------------------------------------------------------ | -------- | ------ |
| 1  | **Rename `getTableDataMarshaler`→`getTableDataRenderer`** + comments + test name           | Medium   | 5 min  |
| 2  | **M4: Rename Render→Draw/Format** in nom module + update all callers + remove NOTE markers | Medium   | 20 min |
| 3  | **Split `integration/roundtrip_test.go`** (528 lines → 2 files)                            | Low      | 10 min |
| 4  | **Split `nom/subscriber_test.go`** (523 lines → 2 files)                                   | Low      | 10 min |
| 5  | **Split `tui/event_sequence_test.go`** (497 lines → 2 files)                               | Low      | 10 min |
| 6  | **Fix nom timing cache test isolation** — use t.TempDir() for cache path                   | Medium   | 10 min |
| 7  | **Update ADR 006** — regenerate frozen symbol list                                         | Medium   | 15 min |
| 8  | **Govulncheck** across all modules                                                         | High     | 5 min  |
| 9  | **Verify BDD module** uses new type names correctly                                        | Low      | 5 min  |
| 10 | **Refresh examples** to show new API patterns (LineStyle, NodeShape, Direction)            | Medium   | 15 min |
| 11 | **Update README.md** code examples with new type names                                     | Medium   | 10 min |
| 12 | **Cut v1.0.0 tag** — after above items done                                                | Critical | 5 min  |
| 13 | **Draft r/golang show-and-tell post**                                                      | Medium   | 10 min |
| 14 | **Draft Awesome Go submission**                                                            | Medium   | 10 min |
| 15 | **Coverage audit** — verify all changed modules still ≥90%                                 | Medium   | 10 min |
| 16 | **Dependency audit** — verify no new deps added by sprint                                  | Low      | 5 min  |
| 17 | **GitHub release notes draft** for v1.0.0                                                  | Medium   | 15 min |
| 18 | **Migration guide** for v0.12→v1.0 users                                                   | Medium   | 20 min |
| 19 | **Godoc review** — verify all new exported types have proper doc comments                  | Low      | 10 min |
| 20 | **Integration test for Direction** — verify D2/DOT renderers accept Direction              | Medium   | 15 min |
| 21 | **Benchmark impact check** — verify renames didn't add allocations                         | Low      | 10 min |
| 22 | **Update `docs/adr/006-api-stability.md`** — remove GraphID reference, update type list    | Low      | 10 min |
| 23 | **Review `graph/CHANGELOG.md`** — update stale Marshaler reference                         | Low      | 2 min  |
| 24 | **Verify `go.work.example`** is up to date with current modules                            | Low      | 2 min  |
| 25 | **Celebrate** — this is a genuinely excellent codebase                                     | High     | ∞      |

---

## G. TOP QUESTION I CANNOT FIGURE OUT MYSELF

**#15: Should `TableData` use exported fields or getters for v1?**

This is the single remaining blocked item. Three valid approaches:

- **Option A (fields only):** Remove `GetHeaders()`, `GetData()`, `GetFooter()`. Users access `data.Headers` directly. Most Go-idiomatic. No validation on mutation.
- **Option B (getters only):** Unexport fields, keep getters, add validated setters. Strongest invariants. Breaks every consumer.
- **Option C (keep both):** Current state. Backwards compatible. Redundant API for v1.

**Which option do you want?** This affects the entire v1.0.0 API surface.

---

## Sprint Metrics

| Metric                             | Value                                                                                    |
| ---------------------------------- | ---------------------------------------------------------------------------------------- |
| Commits since v0.12.0              | 48                                                                                       |
| TODO items resolved (both sprints) | 34                                                                                       |
| Open TODO items                    | 3 (+1 blocked)                                                                           |
| Breaking changes prepared          | 7                                                                                        |
| New types added                    | 4 (Direction, LineStyle, ActivityCounts, NodeShape rename)                               |
| Tests added                        | 11 (3 concurrent + 1 color + 3 direction + 2 linestyle + 1 activity counts + 1 send_nil) |
| Modules                            | 16 (all green)                                                                           |
| Lint issues                        | 0                                                                                        |
| Production files >350 lines        | 0                                                                                        |
| Test files >350 lines              | 3 (pre-existing)                                                                         |
| LOC                                | 31,333                                                                                   |
