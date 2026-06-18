# Comprehensive Plan — All TODOs from 2026-06-17 & 2026-06-18

**Created:** 2026-06-18 18:30
**Scope:** Every TODO surfaced across all docs dated 2026-06-17 and 2026-06-18, deduplicated and verified against actual code at `6ae2d0f` (clean tree).
**Rule:** Every task ≤ 12 minutes. Sorted by **impact → effort → customer value** (desc).

---

## Session Results (executed 2026-06-18 18:30–19:10)

**9 tasks completed, all verified green (17 modules: build ✅ test ✅ race ✅ lint 0-issues ✅):**

| Task | What was done                                                                |
| ---- | ---------------------------------------------------------------------------- |
| H1   | Removed `reflect` depguard violation — `fmt.Sprintf("%T")` instead           |
| H2   | Renamed `getTableDataMarshaler`→`getTableDataRenderer` (+ comments + test)   |
| H3   | `newTempTimingCache(t)` — timing-cache tests no longer touch real `~/.cache` |
| H8   | Typed `ActivityNode.IsPhase()` replaces `isPhaseNode` magic-string func      |
| H10  | `go mod tidy` across all 18 modules (no diff — already tidy)                 |
| D1   | `TODO_LIST.md` synced — M4 marked done, O8 added, session fixes logged       |
| D6   | AGENTS.md: dedup workflow (ADR 005 checklist) + IsPhase convention           |
| D12  | ADR 007 M4 status ✅; FORMAT_ARCHITECTURE.md nom architecture corrected      |
| P5   | Full verification sweep — 17 modules build+test+lint+race clean              |

**Bonus:** fixed a pre-existing `exhaustive` lint failure in `tui/model_test.go` (missing Pending/Paused cases).

**Deferred (needs design decision, not autonomous-safe):** Section 2 — `ActivityDisplayState` elimination requires the "tree-as-view" redesign (DependencyTree must stop owning nodes and resolve activities from the subscriber store at render time). Migration analysis complete; see Section 2 notes.

---

## Verification Baseline (code-checked at plan time)

| Claim in docs                            | Actual code state                                                                          | Verdict     |
| ---------------------------------------- | ------------------------------------------------------------------------------------------ | ----------- |
| M4: `InlineRenderer.Render()` → `Draw()` | ✅ `nom/inline_renderer.go:123 func (r *InlineRenderer) Draw()` — DONE                     | **DONE**    |
| `NOTE(split-brain)` markers              | ✅ Zero matches in any `.go` file                                                          | **DONE**    |
| `getTableDataMarshaler` stale naming     | ❌ Still present (`render_tabledata.go:33,61` + test)                                      | **PENDING** |
| `reflect` depguard violation             | ❌ Still present (`integration/distinctness_test.go:5`)                                    | **PENDING** |
| `ActivityDisplayState` split-brain       | ❌ Present across 10 nom/ files                                                            | **PENDING** |
| `SyncActivityTimingToTree()`             | ❌ Called from tui/, nom/, integration/                                                    | **PENDING** |
| Test files > 350 lines                   | ❌ 4 files: subscriber_test (547), roundtrip_test (528), event_seq (487), model_test (383) | **PENDING** |
| `isPhaseNode` magic string (`"phase:"`)  | ❌ `nom/tree_render.go:306`                                                                | **PENDING** |
| nom/ imports root output types           | ✅ ADR 007 refactor landed; `Activity`, `ActivityStore` exist                              | **DONE**    |

---

## Legend

- **Impact:** 🔥🔥🔥 critical / 🔥🔥 high / 🔥 medium / 🟡 low / 🟢 nice-to-have
- **Effort:** minutes (each ≤ 12)
- **CV:** customer value (release blocker = ★★★)
- **Status:** ⛔ blocked / ⬜ ready / ✅ done

---

## SECTION 1 — BLOCKED ON OWNER DECISION (cannot execute autonomously)

| ID  | Task                                              | Impact | Effort | CV  | Status | Why blocked                                                        |
| --- | ------------------------------------------------- | ------ | ------ | --- | ------ | ------------------------------------------------------------------ |
| B1  | Decide `TableData` fields vs getters for v1 (#15) | 🔥🔥🔥 | 5m     | ★★★ | ⛔     | Needs owner call: A (fields) / B (getters+setters) / C (keep both) |
| B2  | Cut `v1.0.0` tag (#16)                            | 🔥🔥🔥 | 10m    | ★★★ | ⛔     | Depends on B1; needs explicit owner go-ahead                       |
| B3  | Submit to r/golang + Awesome Go (#14)             | 🔥🔥   | 30m    | ★★☆ | ⛔     | Needs owner's Reddit/GitHub account                                |

---

## SECTION 2 — SPLIT-BRAIN ELIMINATION: `ActivityDisplayState` (highest-impact code work)

The single largest pending refactor. Eliminates the dual-state `Activity` vs `ActivityDisplayState`, kills `SyncActivityTimingToTree`, `syncActivityToNode`, and `DisplayState` — ~150 LOC of duplicated state and an entire class of "forgot to sync → stale display" bugs.

| ID  | Task                                                                  | Impact | Effort | CV  | Status | Depends |
| --- | --------------------------------------------------------------------- | ------ | ------ | --- | ------ | ------- |
| S1  | Map every `ActivityDisplayState` field read/write across 10 files     | 🔥🔥🔥 | 10m    | ★☆☆ | ⬜     | —       |
| S2  | Confirm subscriber primary store + projection path                    | 🔥🔥🔥 | 8m     | ★☆☆ | ⬜     | S1      |
| S3  | Refactor subscriber to store `map[ActivityID]*Activity`               | 🔥🔥🔥 | 12m    | ★★★ | ⬜     | S2      |
| S4  | Rewrite `handleActivityStarted` → Activity                            | 🔥🔥   | 10m    | ★★☆ | ⬜     | S3      |
| S5  | Rewrite `handleActivityRegistered` → Activity                         | 🔥🔥   | 8m     | ★★☆ | ⬜     | S3      |
| S6  | Rewrite `handleActivityCompleted` → Activity                          | 🔥🔥   | 10m    | ★★☆ | ⬜     | S3      |
| S7  | Rewrite `handleActivityFailed` → Activity                             | 🔥🔥   | 8m     | ★★☆ | ⬜     | S3      |
| S8  | Rewrite `handleWorkflow*` handlers → Activity meta                    | 🔥     | 8m     | ★☆☆ | ⬜     | S3      |
| S9  | Make `DependencyTree` read temporal data from `Activity`              | 🔥🔥🔥 | 12m    | ★★★ | ⬜     | S4-S7   |
| S10 | Make `InlineRenderer` read from `Activity`                            | 🔥🔥   | 10m    | ★★☆ | ⬜     | S9      |
| S11 | Update `format.go` (`FormatTimingInfo`) for `*Activity`               | 🔥     | 10m    | ★☆☆ | ⬜     | S3      |
| S12 | Update `configuration.go` for `*Activity`                             | 🔥     | 8m     | ★☆☆ | ⬜     | S3      |
| S13 | Update `activity_management.go` (counts/lookups)                      | 🔥🔥   | 10m    | ★★☆ | ⬜     | S3      |
| S14 | Update `tui/model.go` — drop `SyncActivityTimingToTree()` call        | 🔥🔥🔥 | 8m     | ★★★ | ⬜     | S9      |
| S15 | Update `integration/nom_tui_test.go` sync calls                       | 🔥🔥   | 10m    | ★★☆ | ⬜     | S14     |
| S16 | Migrate `activity_display_test.go` → test `Activity`                  | 🔥🔥   | 12m    | ★☆☆ | ⬜     | S13     |
| S17 | Migrate `subscriber_test.go` assertions                               | 🔥🔥   | 12m    | ★☆☆ | ⬜     | S13     |
| S18 | Migrate `format_test.go` + `configuration_test.go`                    | 🔥     | 12m    | ★☆☆ | ⬜     | S11,S12 |
| S19 | Delete `ActivityDisplayState` + `DisplayState` + `syncActivityToNode` | 🔥🔥🔥 | 8m     | ★★★ | ⬜     | S16,S17 |
| S20 | Delete `SyncActivityTimingToTree()` method                            | 🔥🔥🔥 | 5m     | ★★★ | ⬜     | S19     |
| S21 | Verify: nom build + test green                                        | 🔥🔥🔥 | 5m     | ★★★ | ⬜     | S20     |
| S22 | Verify: tui + integration + race green                                | 🔥🔥🔥 | 8m     | ★★★ | ⬜     | S21     |
| S23 | Verify: full workspace `nix run .#build && .#test && .#lint`          | 🔥🔥🔥 | 10m    | ★★★ | ⬜     | S22     |

---

## SECTION 3 — CODE HYGIENE & CORRECTNESS (quick wins)

| ID  | Task                                                                      | Impact | Effort | CV  | Status | Depends  |
| --- | ------------------------------------------------------------------------- | ------ | ------ | --- | ------ | -------- |
| H1  | Fix `reflect` depguard in `integration/distinctness_test.go`              | 🔥🔥   | 8m     | ★☆☆ | ⬜     | —        |
| H2  | Rename `getTableDataMarshaler`→`getTableDataRenderer` (+ comments + test) | 🔥🔥   | 8m     | ★☆☆ | ⬜     | —        |
| H3  | Fix nom timing-cache test isolation (`t.TempDir()`)                       | 🔥🔥   | 10m    | ★★☆ | ⬜     | —        |
| H4  | Split `nom/subscriber_test.go` (547 → 2 files)                            | 🔥     | 10m    | ★☆☆ | ⬜     | —        |
| H5  | Split `integration/roundtrip_test.go` (528 → 2 files)                     | 🔥     | 10m    | ★☆☆ | ⬜     | —        |
| H6  | Split `tui/event_sequence_test.go` (487 → 2 files)                        | 🔥     | 10m    | ★☆☆ | ⬜     | —        |
| H7  | Split `tui/model_test.go` (383 → 2 files)                                 | 🔥     | 10m    | ★☆☆ | ⬜     | —        |
| H8  | Add typed `ActivityNode.IsPhase()` — deprecate `"phase:"` string          | 🔥     | 12m    | ★☆☆ | ⬜     | —        |
| H9  | Refactor `subscriberView.Edges()` to avoid direct `tree.mu`               | 🟡     | 10m    | ★☆☆ | ⬜     | —        |
| H10 | Run `go mod tidy` workspace-wide                                          | 🟢     | 10m    | ★☆☆ | ⬜     | —        |
| H11 | Govulncheck sweep across modules                                          | 🔥🔥   | 10m    | ★★☆ | ⬜     | —        |
| H12 | Coverage audit — verify all modules still ≥90%                            | 🔥     | 10m    | ★☆☆ | ⬜     | H-series |
| H13 | Benchmark impact check — renames added no allocs                          | 🟡     | 10m    | ★☆☆ | ⬜     | —        |
| H14 | Verify `go.work.example` matches current modules                          | 🟢     | 5m     | ★☆☆ | ⬜     | —        |
| H15 | Review `graph/CHANGELOG.md` stale "Marshaler" ref                         | 🟢     | 5m     | ★☆☆ | ⬜     | —        |

---

## SECTION 4 — DOCUMENTATION SYNC

| ID  | Task                                                             | Impact | Effort | CV  | Status | Depends  |
| --- | ---------------------------------------------------------------- | ------ | ------ | --- | ------ | -------- |
| D1  | Sync `TODO_LIST.md` — mark M4 done (code verified)               | 🔥🔥   | 8m     | ★☆☆ | ⬜     | —        |
| D2  | Update `ADR 006` — remove `RenderOptions.GraphID`, refresh types | 🔥     | 12m    | ★☆☆ | ⬜     | —        |
| D3  | Update `FORMAT_ARCHITECTURE.md` — Activity embeds GraphNode      | 🔥     | 10m    | ★☆☆ | ⬜     | S-series |
| D4  | Update `FORMAT_ARCHITECTURE.md` — diagram-export section         | 🔥     | 10m    | ★★☆ | ⬜     | S-series |
| D5  | Update `DOMAIN_LANGUAGE.md` — Activity, ActivityStore terms      | 🟡     | 10m    | ★☆☆ | ⬜     | S-series |
| D6  | Update `AGENTS.md` nom patterns + dedup workflow section         | 🔥🔥   | 12m    | ★☆☆ | ⬜     | —        |
| D7  | Add `extractDependencies`/`capHistory` to nom patterns           | 🟡     | 8m     | ★☆☆ | ⬜     | D6       |
| D8  | Update `FEATURES.md` — diagram export + MultiSubscriber          | 🔥     | 10m    | ★★☆ | ⬜     | —        |
| D9  | Update `README.md` — code examples with new type names           | 🔥     | 10m    | ★★☆ | ⬜     | —        |
| D10 | Update `README.md` — v1 readiness + dedup policy                 | 🟡     | 10m    | ★★☆ | ⬜     | D9       |
| D11 | Add `ADR 008` — dedup-workflow decision                          | 🟡     | 12m    | ★☆☆ | ⬜     | D6       |
| D12 | Audit docs for stale `Render()` refs post-M4                     | 🔥🔥   | 12m    | ★☆☆ | ⬜     | —        |
| D13 | Refresh `examples/` to new API patterns                          | 🟡     | 12m    | ★★☆ | ⬜     | —        |
| D14 | Godoc review — new exported types documented                     | 🟡     | 10m    | ★☆☆ | ⬜     | —        |
| D15 | Draft GitHub release notes for v1.0.0                            | 🟡     | 12m    | ★★☆ | ⬜     | B2       |
| D16 | Draft migration guide v0.12→v1.0                                 | 🟡     | 12m    | ★★☆ | ⬜     | B2       |

---

## SECTION 5 — TEST DEDUPLICATION (ADR 005 follow-up)

| ID  | Task                                                      | Impact | Effort | CV  | Status | Depends  |
| --- | --------------------------------------------------------- | ------ | ------ | --- | ------ | -------- |
| T1  | Extract `fireEvents(sub, ctx, events...)` in tui tests    | 🔥     | 12m    | ★☆☆ | ⬜     | —        |
| T2  | Extract `roundtrip(t, fmt, headers, rows)` in integration | 🔥     | 12m    | ★☆☆ | ⬜     | —        |
| T3  | Extract `setUpFiveActivityFixture()` for golden tests     | 🟡     | 12m    | ★☆☆ | ⬜     | —        |
| T4  | Extract enum-test helpers via `testhelpers/`              | 🟡     | 12m    | ★☆☆ | ⬜     | —        |
| T5  | Golden test for DOT diagram export output                 | 🟡     | 12m    | ★★☆ | ⬜     | S-series |
| T6  | Integration test: full workflow → DOT export              | 🟡     | 12m    | ★★☆ | ⬜     | S-series |

---

## SECTION 6 — PROCESS / CI

| ID  | Task                                                        | Impact | Effort | CV  | Status | Depends |
| --- | ----------------------------------------------------------- | ------ | ------ | --- | ------ | ------- |
| P1  | CI gate on `art-dupl -t 30` (prevent prod-clone regression) | 🔥🔥   | 12m    | ★☆☆ | ⬜     | —       |
| P2  | Pre-commit: run `nix run .#build` across all modules        | 🔥🔥   | 10m    | ★☆☆ | ⬜     | —       |
| P3  | Per-module pre-commit lint                                  | 🟡     | 12m    | ★☆☆ | ⬜     | —       |
| P4  | Audit `image/color` transitive usage in nom/                | 🟢     | 8m     | ★☆☆ | ⬜     | —       |
| P5  | Final pre-v1.0 sweep — all `nix run .#*` apps end-to-end    | 🔥🔥   | 10m    | ★★★ | ⬜     | all     |

---

## SECTION 7 — OPTIONAL / SPECULATIVE (defer unless time)

| ID  | Task                                                          | Impact | Effort | CV  | Status   |
| --- | ------------------------------------------------------------- | ------ | ------ | --- | -------- |
| O1  | Composite `Event` accessor interface (drop 5 type assertions) | 🟡     | —      | ★☆☆ | ⬜       |
| O2  | Replace `switch event.GetEventType()` with handler map        | 🟡     | —      | ★☆☆ | ⬜       |
| O3  | `DependencyTree.SetStore()` store-backed mode                 | 🟡     | —      | ★☆☆ | ⬜       |
| O4  | `RenderOptions` for diagram export (title, theme)             | 🟢     | —      | ★☆☆ | ⬜       |
| O5  | `Theme` struct for custom symbols/colors                      | 🟢     | —      | ★☆☆ | ⬜       |
| O6  | `examples/nom_progress/diagram_export.go` demo                | 🟢     | 12m    | ★★☆ | ⬜       |
| O7  | `examples/tui_progress/diagram_export.go` demo                | 🟢     | 12m    | ★☆☆ | ⬜       |
| O8  | Decide: keep standalone `ActivityStore` or remove (YAGNI)     | 🟡     | 5m     | ★☆☆ | ⛔ owner |

---

## Summary Counts

| Section                   | Tasks  | Ready  | Blocked |
| ------------------------- | ------ | ------ | ------- |
| 1 Blocked (owner)         | 3      | 0      | 3       |
| 2 Split-brain elimination | 23     | 23     | 0       |
| 3 Hygiene/correctness     | 15     | 15     | 0       |
| 4 Documentation           | 16     | 16     | 0       |
| 5 Test dedup              | 6      | 6      | 0       |
| 6 Process/CI              | 5      | 5      | 0       |
| 7 Optional/speculative    | 8      | 7      | 1       |
| **Total**                 | **76** | **72** | **4**   |

**Execution order:** Section 3 quick wins (H1-H3, H10-H11) → Section 4 D1/D6/D12 → Section 2 split-brain (S1-S23) → remaining hygiene/docs/tests → P5 final sweep.

---

_Generated 2026-06-18 18:30 · Verified against `6ae2d0f` (clean tree)_
