# Comprehensive Improvement Plan — 2026-07-05

**Source:** Brutal self-review `docs/reviews/2026-07-05_20-57_brutal-self-review.html`
**59 tasks** — each ≤12 min. Sorted by impact ÷ effort (Pareto).
**Resolution (2026-08-04):** All P0–P7 tasks shipped in v0.30.0. Concurrency bugs fixed, dead code deleted, split brains resolved, deprecated APIs removed. See CHANGELOG `[0.30.0]` for details.

## Priority Key

| Band   | Meaning                                                                                |
| ------ | -------------------------------------------------------------------------------------- |
| **P0** | Correctness bug — crashes, races, wrong output. Consumer-visible.                      |
| **P1** | Trust violation — doc/API lies that mislead consumers or contributors.                 |
| **P2** | Split brain — two sources of truth for one concept.                                    |
| **P3** | Ghost system — shipped feature with zero consumers (prove or remove before v1 freeze). |
| **P4** | Dead code — provably unreachable (unexported or never-produced).                       |
| **P5** | Pattern smell — stringly-typed, false-promise API, complexity, missing robustness.     |
| **P6** | Naming polish — breaking rename; deprecate for v1, remove in v2.                       |
| **P7** | Housekeeping — doc archive, test gaps, community.                                      |

Within each band: sorted by effort ascending (quickest wins first).

---

## Full Task Table

| #  | Band | Task                                                                                   | Finding  | Files                                        | Effort | Impact   |
| -- | ---- | -------------------------------------------------------------------------------------- | -------- | -------------------------------------------- | ------ | -------- |
| 1  | P0   | Delete `hasDep()` — unexported, gopls-confirmed unused                                 | D6       | `nom/tree.go:149`                            | 2m     | Low      |
| 2  | P0   | Delete dead test helpers `assertScreenNotContains` / `assertRowEmpty`                  | D10      | `nom/vttest_test.go:127,138`                 | 3m     | Low      |
| 3  | P0   | Fix `r.appName` race: use `snapshotConfig().appName` in `RenderCompletion`             | C1       | `nom/inline_renderer.go:816`                 | 5m     | Critical |
| 4  | P0   | Route `renderNotify` read through `snapshotConfig` or `tickMu`                         | C3       | `nom/inline_renderer.go:728`                 | 8m     | High     |
| 5  | P0   | Fix `showParallelism` unlocked read: snapshot or document immutable                    | C6       | `nom/inline_renderer_summary.go:145`         | 8m     | Medium   |
| 6  | P0   | Drop `ns.mu` before `timingCache.EnsureLoaded()` in `handleWorkflowStarted`            | C2       | `nom/subscriber_handlers.go:57`              | 8m     | Critical |
| 7  | P0   | Drop `ns.mu` before `timingCache.Save()` in `handleWorkflowFinished`                   | C2       | `nom/subscriber_handlers.go:68`              | 8m     | Critical |
| 8  | P0   | Wire `physicalLines` into `buildLogAndTreeOutput` (latent ghost-line bug)              | D9       | `nom/inline_renderer.go:500`                 | 10m    | Critical |
| 9  | P0   | Replace `log.Printf` error swallowing with returned error in `saveAsync`               | C5       | `nom/timing_cache_persist.go:160`            | 10m    | High     |
| 10 | P0   | Extract `shouldSkipDraw()` from `Draw()`                                               | P5       | `nom/inline_renderer.go:357`                 | 10m    | High     |
| 11 | P0   | Extract `drawPlainText()` from `Draw()`                                                | P5       | `nom/inline_renderer.go:357`                 | 12m    | High     |
| 12 | P0   | Extract `drawInline()` from `Draw()`                                                   | P5       | `nom/inline_renderer.go:357`                 | 12m    | High     |
| 13 | P0   | Design single background-saver goroutine + buffered channel                            | C4       | `nom/timing_cache.go:102`                    | 10m    | Critical |
| 14 | P0   | Implement background saver, replace per-call `go saveAsync()`                          | C4       | `nom/timing_cache.go:102`                    | 12m    | Critical |
| 15 | P0   | Add `Flush()` method to subscriber to drain pending cache saves                        | C4       | `nom/timing_cache.go`                        | 8m     | High     |
| 16 | P0   | Wire `Flush()` into subscriber shutdown/Close path                                     | C4       | `nom/nom_subscriber.go`                      | 8m     | High     |
| 17 | P0   | Update `waitPendingSaves` tests to use new background-saver mechanism                  | C4       | `nom/timing_cache_test.go`                   | 8m     | Medium   |
| 18 | P0   | Add race test: concurrent config setters vs Draw/Finish/RenderCompletion               | C1/C3/C6 | `nom/inline_renderer_race_test.go`           | 12m    | High     |
| 19 | P0   | Verify cyclop lint passes after Draw() decomposition                                   | P5       | `nom/inline_renderer.go`                     | 3m     | High     |
| 20 | P0   | Run `nix run .#test-race` after all concurrency fixes                                  | C1-C6    | all nom/                                     | 5m     | Critical |
| 21 | P1   | Fix AGENTS.md: `SymbolOverrides` → `Symbols`                                           | S5/L4    | `AGENTS.md`                                  | 3m     | High     |
| 22 | P1   | Fix TODO_LIST claim: "Zero Deprecated markers remain" is false                         | L1       | `TODO_LIST.md:26`                            | 3m     | High     |
| 23 | P1   | Close TODO_LIST #18: v0.23.0–v0.23.3 already tagged                                    | L2       | `TODO_LIST.md:15`                            | 3m     | Medium   |
| 24 | P1   | Close/clarify TODO_LIST #19: daghtml lints clean (0 issues)                            | L3       | `TODO_LIST.md:16`                            | 5m     | Medium   |
| 25 | P1   | Update TODO_LIST #16: v1.0.0 tag status (post this sprint)                             | L2       | `TODO_LIST.md:14`                            | 5m     | Medium   |
| 26 | P1   | Fix `ErrActivityNotFound`: make `GetNode` return it, OR delete the var                 | D3       | `nom/tree.go:11`, `nom/tree_accessors.go:52` | 10m    | Medium   |
| 27 | P1   | Add `// Deprecated:` to `TimingFormat`, fix lying comment                              | D4       | `nom/format.go:12`                           | 5m     | Medium   |
| 28 | P1   | Add `// Deprecated:` to `Activity.IsPhase()`, point to `snap.IsPhase()`                | D5       | `nom/activity.go:139`                        | 5m     | Medium   |
| 29 | P1   | Decision: keep `SymbolUpload` as palette or deprecate                                  | D8       | `nom/symbols.go:32`                          | 2m     | Low      |
| 30 | P1   | Decision: wire `EdgeStyle.ArrowHead/ArrowTail` into renderers or deprecate             | D7       | `graph.go:189`                               | 5m     | Medium   |
| 31 | P1   | Deprecate `GetDependencyTree()`, add `// Deprecated:` pointing to `DependencyTree()`   | S3       | `nom/state_accessors.go:24`                  | 5m     | Medium   |
| 32 | P1   | Update internal callers from `GetDependencyTree()` → `DependencyTree()`                | S3       | `tui/`, `integration/`                       | 10m    | Medium   |
| 33 | P1   | Update FEATURES.md + CHANGELOG with deprecations from #27-31                           | —        | `FEATURES.md`, `CHANGELOG.md`                | 10m    | Medium   |
| 34 | P2   | Export `MsgNoActivities` in nom package (single source of truth)                       | S1       | `nom/tree.go:17`                             | 5m     | Medium   |
| 35 | P2   | Update tui to import `nom.MsgNoActivities`, delete local constant                      | S1       | `tui/constants.go:16`                        | 5m     | Medium   |
| 36 | P2   | Make `Direction.ToD2Direction()` return `D2Direction` (type-safe bridge)               | S4       | `direction.go:27`                            | 10m    | Medium   |
| 37 | P2   | Audit all direct `Colors` global reads that bypass Theme                               | S2       | `nom/symbols.go:78`                          | 8m     | High     |
| 38 | P2   | Route `status_registry.go` color reads through Theme                                   | S2       | `nom/status_registry.go:44,57,70,83`         | 12m    | High     |
| 39 | P2   | Route `activity.go` color read through Theme                                           | S2       | `nom/activity.go:161`                        | 8m     | Medium   |
| 40 | P2   | Route `activity_status.go` color read through Theme                                    | S2       | `nom/activity_status.go:55`                  | 5m     | Medium   |
| 41 | P2   | Add test verifying themed colors override defaults in all paths                        | S2       | `nom/theme_test.go`                          | 10m    | High     |
| 42 | P3   | Write example binary registering a custom "skipped" status                             | D1       | `examples/`                                  | 12m    | Medium   |
| 43 | P3   | Write integration test rendering custom status end-to-end                              | D1       | `integration/`                               | 12m    | Medium   |
| 44 | P3   | Decision: keep or delete `ThemeNord` and `ThemeMonochrome`                             | D2       | `nom/theme.go:117,139`                       | 2m     | Low      |
| 45 | P3   | If kept: add smoke tests for `ThemeNord` and `ThemeMonochrome`                         | D2       | `nom/theme_test.go`                          | 8m     | Low      |
| 46 | P3   | If kept: add theme to example binary                                                   | D2       | `examples/`                                  | 8m     | Low      |
| 47 | P3   | If wiring arrows (#30=wiring): implement in DOT renderer                               | D7       | `graph/dot.go`                               | 10m    | Medium   |
| 48 | P3   | If wiring arrows: implement in D2 renderer                                             | D7       | `d2/d2_write.go`                             | 10m    | Medium   |
| 49 | P3   | If wiring arrows: implement in Mermaid renderer                                        | D7       | `graph/mermaid.go`                           | 8m     | Medium   |
| 50 | P3   | If wiring arrows: implement in PlantUML renderer                                       | D7       | `plantuml/plantuml.go`                       | 8m     | Medium   |
| 51 | P5   | Rename `StreamingRendererFromRenderer` → `RendererAsWriter` (stop lying)               | P1       | `streaming.go:24`                            | 8m     | Medium   |
| 52 | P5   | Update callers/tests of renamed adapter                                                | P1       | `tabledatastore_test.go`                     | 5m     | Low      |
| 53 | P5   | Create `TextTransform` enum (Parse/IsValid/AllowedValues)                              | P2       | `d2/d2_enum.go`                              | 12m    | Medium   |
| 54 | P5   | Change `D2NodeStyle.TextTransform` from `string` to enum type                          | P2       | `d2/d2.go:47`                                | 5m     | Medium   |
| 55 | P5   | Update D2 render/write code for typed TextTransform                                    | P2       | `d2/d2_write.go`                             | 8m     | Medium   |
| 56 | P5   | Make `Build()` return cycle-detection error, or drop the `error` return                | P3       | `nom/tree_building.go:8`                     | 10m    | Medium   |
| 57 | P5   | Add consecutive-failure counter + self-Stop in `write()`                               | P4       | `nom/inline_renderer.go:594`                 | 10m    | Medium   |
| 58 | P5   | Add test for dead-writer self-Stop behavior                                            | P4       | `nom/inline_renderer_test.go`                | 10m    | Low      |
| 59 | P6   | Deprecate `NOMStyleSubscriber` → type alias `NOMSubscriber` (v2 remove)                | P6       | `nom/nom_subscriber.go:22`                   | 10m    | Low      |
| 60 | P6   | Deprecate `SemanticColors.Info` → alias `.Fallback` (v2 remove)                        | P6       | `nom/symbols.go:69`                          | 8m     | Low      |
| 61 | P7   | Create `docs/archive/status/` directory                                                | DB       | `docs/`                                      | 2m     | Low      |
| 62 | P7   | Move pre-v0.20 status snapshots (63 files) to archive                                  | DB       | `docs/status/`                               | 8m     | Low      |
| 63 | P7   | Move pre-v0.20 reviews/planning snapshots to archive                                   | DB       | `docs/reviews/`, `docs/planning/`            | 8m     | Low      |
| 64 | P7   | Update CHANGELOG Unreleased section with all changes from this sprint                  | —        | `CHANGELOG.md`                               | 10m    | Medium   |
| 65 | P7   | Full regression: `nix run .#build && .#test && .#test-race && .#lint && .#govulncheck` | —        | all                                          | 10m    | Critical |

---

## Summary by Band

| Band                 | Tasks | Total Effort | What It Unlocks                             |
| -------------------- | ----- | ------------ | ------------------------------------------- |
| **P0** Correctness   | 1–20  | ~2.5h        | Race-free nom/, clean lint, no ghost lines  |
| **P1** Trust         | 21–33 | ~1.5h        | Docs match code, honest deprecated surface  |
| **P2** Split brains  | 34–41 | ~1h          | Single source of truth for constants/colors |
| **P3** Ghost systems | 42–50 | ~1.5h        | Proven or removed speculative features      |
| **P5** Patterns      | 51–58 | ~1.5h        | Type-safe D2, honest adapters, robust I/O   |
| **P6** Naming        | 59–60 | ~20m         | v2-ready aliases                            |
| **P7** Housekeeping  | 61–65 | ~40m         | Pruned docs, regression-verified release    |

**Total: ~9.5h** | **The 80/20 cut line: tasks 1–20 (P0) = ~2.5h = all real risk eliminated.**
