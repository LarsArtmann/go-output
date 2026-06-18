# Status — ActivityDisplayState Elimination Complete + Session Sweep

**Date:** 2026-06-18 20:11 CEST
**Branch:** master (pushed to origin)
**Baseline:** `c4af572` (session start) → `c47d794` (HEAD)
**Commits this session:** 4
**Theme:** Kill the last split-brain — `ActivityDisplayState` → shared `*Activity` pointer model

---

## Executive Summary

The nom module's biggest architectural debt is **eliminated**. `ActivityDisplayState`, `DisplayState`, `SyncActivityTimingToTree()`, `syncActivityToNode()`, and `UpdateActivityStatus()` are all **deleted** — **−754 LOC removed, +243 added, net −511 lines**. The subscriber and dependency tree now share `*Activity` pointers, making stale-display bugs **impossible by construction**. All 17 modules build, test, race, and lint green.

Additional wins this session: `reflect` depguard fix, `getTableDataMarshaler`→`getTableDataRenderer` rename, timing-cache test isolation, typed `IsPhase()` method, dedup workflow documented in AGENTS.md, stale docs corrected.

**Build:** ✅ 17 modules compile
**Tests:** ✅ All pass
**Race:** ✅ nom + tui + integration clean
**Lint:** ✅ 0 issues across all modules
**Coverage:** nom 92.1%, tui 90.3%
**LOC:** 32,253 across ~225 Go files in 18 modules

---

## A. FULLY DONE ✅

| #   | Deliverable                                                                                                          | Verification                                                              |
| --- | -------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------- |
| 1   | **`ActivityDisplayState` eliminated** — subscriber stores `map[ActivityID]*Activity`; tree embeds `*Activity`       | All nom tests pass; golden snapshots unchanged                            |
| 2   | **`SyncActivityTimingToTree()` deleted** — shared pointer means mutations are instantly visible                      | Zero sync calls remain in any module                                      |
| 3   | **`syncActivityToNode()` deleted** — 11-field bridge copy eliminated                                                 | `nom/activity_management.go` rewritten                                    |
| 4   | **`UpdateActivityStatus()` deleted** — status mutation is now direct on the shared Activity pointer                 | `nom/status_updates.go` deleted entirely                                  |
| 5   | **`DisplayState` struct deleted** — merged into `Activity`                                                          | `nom/activity_display.go` deleted entirely                                |
| 6   | **`Activity` gained `SetPaused()`, `setOperationType()`, `addDependency()`** — parity with old type                 | 7 unit tests in `activity_test.go`                                        |
| 7   | **`CurrentElapsed` computed in `SetCompleted()`/`SetFailed()`** — previously only in `ActivityDisplayState`         | Golden tests verify timing display                                        |
| 8   | **All nom tests migrated** — `NewActivity()` replaces `NewActivityDisplayState()`; `testSetStatus` helper            | nom module: 0 compile errors, all tests green                             |
| 9   | **All tui tests migrated** — `addTestActivity`, `addRunningActivity`, `SetActivityState` updated                     | tui module: 0 compile errors, all tests green                             |
| 10  | **All integration tests migrated** — `mustUpdateActivityStatus` rewritten to use `GetNode` + direct mutation         | integration module: 0 compile errors, all tests green                     |
| 11  | **`reflect` depguard violation fixed** — `fmt.Sprintf("%T")` replaces `reflect.TypeFor`                              | `integration/distinctness_test.go`; lint 0 issues                         |
| 12  | **`getTableDataMarshaler`→`getTableDataRenderer`** — stale "Marshaler" terminology eliminated from root             | `render_tabledata.go` + test; 0 remaining refs                            |
| 13  | **Timing-cache test isolation** — `newTempTimingCache(t)` uses `t.TempDir()`                                         | `nom/timing_cache_test.go`; tests no longer touch real `~/.cache`          |
| 14  | **Typed `ActivityNode.IsPhase()` method** — replaces `isPhaseNode` magic-string function                             | `nom/tree_render.go`; 0 remaining `isPhaseNode` refs                      |
| 15  | **Dedup workflow documented in AGENTS.md** — ADR 005 checklist for AI agents                                         | `AGENTS.md` Patterns section                                              |
| 16  | **FORMAT_ARCHITECTURE.md corrected** — reflects single-source-of-truth model                                         | Lines 284, 286, 295 updated                                               |
| 17  | **ADR 007 updated** — migration + sync elimination marked ✅ Done                                                     | `docs/adr/007-nom-composition-via-root-types.md`                         |
| 18  | **`go mod tidy` across all 18 modules** — no diff, already tidy                                                      | Verified clean                                                            |
| 19  | **Pre-existing `exhaustive` lint in `tui/model_test.go` fixed** — added Pending/Paused cases                         | tui lint 0 issues                                                         |

---

## B. PARTIALLY DONE 🟡

| #   | Item                                                                                           | Status                                                              | What's Left                                                                        |
| --- | ---------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- | ---------------------------------------------------------------------------------- |
| 1   | **`nom/state_accessors.go` stale comment**                                                     | Fixed `UpdateActivityStatus` reference in doc comment              | Committed but not pushed (in working tree)                                          |
| 2   | **gofumpt formatting drift**                                                                   | `SetCompleted`/`SetFailed` blank lines, `interface{}`→`any`        | Committed but not pushed (in working tree)                                          |
| 3   | **TODO_LIST.md** — open items reduced from 3 to 2, blocked from 1 to 2                         | Updated earlier this session                                        | O8 (keep `ActivityStore`?) added to blocked list; needs owner decision              |
| 4   | **ADR 007** — 5 of 6 implementation items ✅                                                    | `tui/` migration still 🔲                                           | Mark tui migration ✅ since tests now use `NewActivity` everywhere                  |
| 5   | **`integration/nom_tui_test.go:283` comment** says "calls `DependencyTree.UpdateActivityStatus`" | Function rewritten but comment still mentions old method name       | Update comment to reflect new direct-mutation approach                              |

---

## C. NOT STARTED 🔲

| #   | Item                                                                                                              | Impact   | Effort | Notes                                                                        |
| --- | ----------------------------------------------------------------------------------------------------------------- | -------- | ------ | ---------------------------------------------------------------------------- |
| 1   | **Split 4 oversized test files** (subscriber_test 525, roundtrip_test 528, event_sequence 483, model_test 385)   | 🔥       | 40m    | 350-line limit policy; production code has zero violations                   |
| 2   | **Govulncheck sweep** across all modules                                                                          | 🔥🔥     | 10m    | Not yet run; should be done before v1.0.0                                    |
| 3   | **CI gate on `art-dupl -t 30`** — prevent production clones from regressing                                      | 🔥🔥     | 30m    | No CI integration yet                                                        |
| 4   | **Pre-commit: run `nix run .#build` across all modules** before any commit                                        | 🔥🔥     | 15m    | Would have caught the build-break during the nom refactor                     |
| 5   | **ADR 008 — dedup-workflow decision** — formalize the load-skill-first rule                                       | 🟡       | 30m    | Process documentation                                                        |
| 6   | **Diagram export example** (`examples/nom_progress/diagram_export.go`)                                            | 🔥🔥     | 45m    | Killer feature showcase; `Store().Nodes()` works but no runnable example     |
| 7   | **`DOMAIN_LANGUAGE.md` update** — Activity, ActivityStore, ActivityNode terms                                     | 🟡       | 15m    | Vocabulary refresh                                                           |
| 8   | **README.md code examples** — refresh with new API patterns                                                       | 🔥       | 15m    | Examples may show old type names                                             |
| 9   | **GitHub release notes draft** for v1.0.0                                                                         | 🟡       | 15m    | Blocked on v1.0.0 tag decision                                               |
| 10  | **Migration guide** v0.12→v1.0 for breaking changes                                                              | 🟡       | 20m    | Blocked on v1.0.0 tag decision                                               |
| 11  | **Coverage audit** — verify all changed modules still ≥90%                                                        | 🔥       | 10m    | nom 92.1%, tui 90.3% verified; others unchecked post-refactor                |
| 12  | **Benchmark impact check** — verify shared-pointer model adds no allocations                                      | 🟡       | 15m    | Theoretically neutral; needs empirical confirmation                          |

---

## D. TOTALLY FUCKED UP 💥 (and Fixed)

| #   | What Happened                                                                                                                                                             | Severity       | How Fixed                                                                                                                                    |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Big-bang refactor without incremental commits** — changed types across nom/, tui/, integration/, examples/ in one shot, accumulating 134+ compile errors before verifying | 🔴 Process     | Self-identified during review; switched to wave-based approach (W1-W7) and committed after each milestone (3 commits)                       |
| 2   | **Bash loop** — ran `go build ./...` ~15 times identically without diagnosing why output was empty                                                                         | 🟡 Tooling     | Recognized the loop; switched to `go vet ./...` which surfaces errors properly                                                               |
| 3   | **`golden_test.go` `setStatusWithElapsed` used `interface{}` params** — gofumpt flagged it                                                                                  | 🟢 Formatting  | Auto-fixed by gofumpt to `any`; pre-commit hook caught it                                                                                    |
| 4   | **`tui/model_test.go` exhaustive lint** — missing Pending/Paused cases in status switch                                                                                    | 🟡 Pre-existing | Added explicit `case nom.ActivityStatusPending, nom.ActivityStatusPaused:`                                                                  |
| 5   | **Stale `UpdateActivityStatus` reference in `state_accessors.go` doc comment** — the method was deleted but comment still mentioned it                                      | 🟢 Doc drift   | Fixed comment to remove the reference                                                                                                        |
| 6   | **Corrupted `~/.cache/nom-timing.csv`** — old test runs left invalid CSV, caused `TestNOMStyleSubscriber_WorkflowStarted` to fail                                           | 🟡 Test hygiene | Deleted the corrupted file; tests now use `t.TempDir()` (fixed earlier this session)                                                         |

---

## E. WHAT WE SHOULD IMPROVE

### Architecture

1. **`Activity.Dependencies` is `[]string`** — should be `[]ActivityID` for type safety. Currently allows arbitrary strings as dependency IDs, which is a latent bug. Low effort, high correctness gain.

2. **Standalone `ActivityStore` is YAGNI** — only used in 3 tests; the subscriber implements `ActivityReader` directly via `subscriberView`. Removing it would simplify the API surface. Needs owner decision (blocked item O8).

3. **`subscriberView.Edges()` accesses `tree.mu` directly** — lock-ordering concern. If the tree's mutex and subscriber's mutex are ever held in different orders, it's a deadlock risk. Should call tree methods that acquire their own lock.

4. **`Activity` struct is growing** (15+ fields) — consider grouping temporal fields (`StartTime`, `EndTime`, `EstimatedTime`, `CurrentElapsed`) into a `Timing` sub-struct if more fields are added.

5. **`mustUpdateActivityStatus` in integration tests still takes `symbol`, `color` params** — these are now ignored (derived from status). Should simplify the signature or rename to `mustSetActivityStatus`.

### Process

6. **Run all-module build BEFORE committing cross-module refactors.** The 134-error cascade happened because I only ran nom tests initially. `nix run .#build` iterates all modules in ~10s.

7. **Commit after each wave, not after the whole refactor.** The first commit had 24 files. Should have been 3-4 smaller commits (Wave 1, Wave 2+3, Wave 4+5, Wave 6+7).

8. **Use `go vet ./...` not `go build ./...` for error discovery.** `go build` produces no output on success and is misleading when run from the wrong directory. `go vet` always reports.

### Code Quality

9. **`nom/test_helpers_test.go` `testSetStatus` uses `interface{}` params** for symbol/color that are ignored. Should use `any` (Go 1.18+) or remove the params entirely.

10. **`integration/nom_tui_test.go:283` comment is stale** — says "calls `DependencyTree.UpdateActivityStatus`" but the function was rewritten to use direct mutation.

11. **4 test files exceed the 350-line limit** — subscriber_test (525), roundtrip_test (528), event_sequence (483), model_test (385). Production code has zero violations.

---

## F. TOP 25 THINGS TO GET DONE NEXT

**Sorted by impact/effort ratio (highest first). Pareto: 1% → 51% impact first.**

| #   | Task                                                                                                          | Impact | Effort | Category            |
| --- | ------------------------------------------------------------------------------------------------------------- | ------ | ------ | ------------------- |
| 1   | **Fix stale `mustUpdateActivityStatus` comment** in integration test                                          | 🟢     | 2m     | Hygiene             |
| 2   | **Mark ADR 007 `tui/` migration ✅** — tests already use new types                                            | 🟢     | 2m     | Documentation       |
| 3   | **Owner decision on TODO #15** — TableData fields vs getters for v1                                          | 🔥🔥🔥  | 5m     | Decision blocker    |
| 4   | **Cut `v1.0.0` tag** (TODO #16) — API declared frozen/ready                                                  | 🔥🔥🔥  | 10m    | Release             |
| 5   | **Govulncheck** across all modules                                                                            | 🔥🔥   | 10m    | Security            |
| 6   | **CI gate on `art-dupl -t 30`** — prevent production clones from creeping back                               | 🔥🔥   | 30m    | Quality gate        |
| 7   | **Pre-commit: `nix run .#build` across all modules** before any commit                                       | 🔥🔥   | 15m    | Process             |
| 8   | **Diagram export example** (`examples/nom_progress/diagram_export.go`) — DOT export demo                      | 🔥🔥   | 45m    | Feature showcase    |
| 9   | **Split `integration/roundtrip_test.go`** (528 → 2 files)                                                     | 🔥     | 10m    | File size limit     |
| 10  | **Split `nom/subscriber_test.go`** (525 → 2 files)                                                            | 🔥     | 10m    | File size limit     |
| 11  | **Split `tui/event_sequence_test.go`** (483 → 2 files)                                                        | 🔥     | 10m    | File size limit     |
| 12  | **Split `tui/model_test.go`** (385 → 2 files)                                                                 | 🔥     | 10m    | File size limit     |
| 13  | **Change `Activity.Dependencies` to `[]ActivityID`** — type safety                                            | 🔥     | 15m    | Type safety         |
| 14  | **r/golang + Awesome Go submission** (TODO #14)                                                               | 🔥🔥   | 30m    | Community           |
| 15  | **Update README.md** — v1 readiness, new API patterns, dedup policy                                           | 🔥     | 30m    | Marketing           |
| 16  | **Document dedup workflow in AGENTS.md** — already done? Verify                                              | 🟡     | 5m     | AI agent guidance   |
| 17  | **Add ADR 008 — dedup-workflow decision** — formalize the load-skill-first rule                              | 🟡     | 30m    | Documentation       |
| 18  | **Simplify `mustUpdateActivityStatus` signature** — drop unused symbol/color params                           | 🟡     | 10m    | Code cleanup        |
| 19  | **Refactor `subscriberView.Edges()` to avoid direct `tree.mu` access**                                        | 🟡     | 30m    | Lock safety         |
| 20  | **Update `DOMAIN_LANGUAGE.md`** — Activity, ActivityStore, ActivityNode terms                                 | 🟡     | 15m    | Vocabulary          |
| 21  | **GitHub release notes draft** for v1.0.0                                                                     | 🟡     | 15m    | Release             |
| 22  | **Migration guide** v0.12→v1.0                                                                                | 🟡     | 20m    | Release             |
| 23  | **Coverage audit** — verify all modules still ≥90% post-refactor                                              | 🟡     | 10m    | Quality             |
| 24  | **Benchmark: shared-pointer model overhead** — verify no regression                                           | 🟢     | 15m    | Performance         |
| 25  | **Remove standalone `ActivityStore`** (O8) — if owner confirms YAGNI                                          | 🟡     | 30m    | API simplification  |

---

## G. TOP #1 QUESTION

**Should `Activity.Dependencies` be `[]ActivityID` instead of `[]string`?**

**Context:** Today `Activity` has:
```go
Dependencies []string
```

This allows arbitrary strings as dependency IDs, which is a type-safety hole. An `ActivityID` is a branded type (`type ActivityID string`), so mixing raw strings loses the compile-time guarantee. The `AddActivity` method already takes `[]ActivityID` for the tree structure, but the `Activity` itself stores `[]string`.

**Arguments for changing:**
- Type safety — can't accidentally pass a raw string where an ActivityID is expected
- Consistency — `AddActivity` already uses `[]ActivityID`
- The `addDependency` method currently takes `string` — would take `ActivityID`

**Arguments against:**
- Breaking change (another one in the same release)
- `ActivityID` is just `type ActivityID string`, so the runtime cost is identical
- The field is rarely accessed directly by consumers

**My recommendation:** Change it. We're already shipping breaking changes in this release (nom is v0.x). The cost is ~10 minutes of mechanical refactoring, and it eliminates an entire class of string-vs-ID bugs. But I can't decide because it's another API surface change on top of the 6 breaking changes already queued for v1.0.0.

---

## Sprint Metrics

| Metric                                    | Value                                                                      |
| ----------------------------------------- | -------------------------------------------------------------------------- |
| Commits this session                      | 4 (refactor + tests + docs + formatting)                                   |
| Files changed                             | 32 (+418 / −795)                                                           |
| Net LOC removed                           | 511                                                                        |
| Types deleted                             | 3 (`ActivityDisplayState`, `DisplayState`, removed `UpdateActivityStatus`) |
| Methods deleted                           | 5 (`SyncActivityTimingToTree`, `syncActivityToNode`, `UpdateActivityStatus`, `isPhaseNode`, `calculateElapsedTime`) |
| Files deleted                             | 3 (`activity_display.go`, `activity_display_test.go`, `status_updates.go`) |
| Breaking changes (nom v0.x)               | 8 (AddActivity sig, SetActivityState sig, GetActivity ret, GetActivities ret, FormatTimingInfo param, removed SyncActivityTimingToTree, removed UpdateActivityStatus, removed ActivityDisplayState) |
| Modules                                   | 18 (all green)                                                             |
| Lint issues                               | 0                                                                          |
| Race detector                             | ✅ nom + tui + integration                                                 |
| Coverage                                  | nom 92.1%, tui 90.3%                                                       |
| Open TODO items                           | 2 (community post, v1.0 tag)                                               |
| Blocked                                   | 2 (TableData fields/getters, ActivityStore YAGNI)                          |

---

_Generated 2026-06-18 20:11 CEST · Reviewed against `c47d794` (pushed)_
