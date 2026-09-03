# Status: go-output/nom BuildFlow Integration — Session 2 (Completion + Remaining Gaps)

**Date:** 2026-07-01 18:40 CEST
**Session Focus:** Close the in-repo gaps from the 12:33 NOM BuildFlow status report — retry reason display, subscriber-owned `~Xm left`, TUI summary wiring, concurrency safety tests, and documentation. Correct the report's erroneous TUI claim.
**Overall Health:** ✅ All 18 modules passing · 0 lint issues · 10 new tests (nom) + 3 new tests (tui) · 0 regressions · Race clean

---

> **✅ Resolved (2026-08-04):**
>
> All items resolved. Critical-path ETA (the #1 question: "sum vs critical-path") was answered — both implemented in v0.23.0 (`EstimatedTotalRemaining` for sum, `EstimatedCriticalPathRemaining` for longest path). Reset() clearing, multi-subscriber fan-out tests done. Remaining open: structured progress type, adaptive tree pruning (both in ROADMAP).

---

## Executive Summary

This session continued the NOM BuildFlow integration work from the 12:33 status report. The previous session shipped four new capabilities (ActivityProgress, ActivityRetrying, SetEstimatedTime, SetEstimatedRemainingFunc). This session closed the top in-repo follow-up items: added retry **reason** display (`⟳2 (timeout)`), built a subscriber-owned `EstimatedTotalRemaining()` primitive that powers `~Xm left` in both renderers, wired the TUI summary bar to consume it automatically, wrote comprehensive concurrency/race tests, and updated all documentation and examples.

**Critical correction:** The previous report's #6 P1 item ("TUI rendering of Progress/RetryCount — the data path is complete; only the view layer needs updating") was **wrong**. The TUI already renders both via its `RenderVisibleEntry` delegation to `formatActivityLabel`. This session proved it with dedicated tests rather than adding redundant view-layer code, and documented the delegation pattern in AGENTS.md to prevent future confusion.

**13 production/doc files changed, 2 new test files, +275 insertions, −16 deletions. 13 new tests (10 nom + 3 tui), all passing. 0 lint issues. Race clean.**

---

## a) FULLY DONE ✅

### 1. Retry Reason Display — `⟳2 (timeout)`

**Problem:** `ActivityRetrying` had no way to communicate _why_ the retry happened. BuildFlow's retry path knows the cause (timeout, network error, flaky test) but nom had no field for it.

**What shipped:**

- **`Reason string` field** on `ActivityRetrying` event (`event.go`) — non-breaking, zero-value backward-compatible (no reason = plain `⟳N`)
- **`RetryReason string` field** on `Activity` (`activity.go`) and `ActivitySnapshot` (`activity_snapshot.go`) — threaded through the snapshot path for race-free rendering
- **Handler update** in `subscriber_handlers.go` — `handleActivityRetrying` sets `activity.RetryReason = e.Reason`
- **Rendering** in `tree_render.go:formatActivityLabel` — renders `⟳2 (timeout)` when reason is non-empty, plain `⟳2` when empty

**Files:** `event.go`, `activity.go`, `activity_snapshot.go`, `subscriber_handlers.go`, `tree_render.go`

**Tests:** `TestRetryReasonEvent`, `TestRetryReasonEmpty` (in `progress_events_test.go`)

### 2. Subscriber-Owned `EstimatedTotalRemaining()` — Single Source of Truth

**Problem:** The previous session shipped `SetEstimatedRemainingFunc` on the InlineRenderer (callback-based), but the subscriber had no way to compute the remaining time from its own activity estimates. BuildFlow had to compute the sum externally. The status report listed "Subscriber-level `EstimatedTotalRemaining()`" as P2 item #11.

**What shipped:**

- **`EstimatedTotalRemaining() time.Duration`** on `NOMStyleSubscriber` in `state_accessors.go` — sums remaining estimates:
  - Pending activities: contribute full `EstimatedTime`
  - Running activities: contribute `max(0, EstimatedTime - elapsed)`
  - Completed/failed activities: contribute nothing
  - Returns 0 when no unfinished activity has an estimate
- This is now the **single source of truth** for the `~Xm left` summary segment. The InlineRenderer's callback can delegate to it; the TUI consumes it directly.

**Files:** `state_accessors.go`

**Tests:** `TestEstimatedTotalRemaining`, `TestEstimatedTotalRemainingZero`, `TestEstimatedTotalRemainingRunningElapsed` (verifies elapsed subtraction for running activities)

### 3. TUI `~Xm left` Summary Wiring

**Problem:** The TUI's NOM summary bar (`tui/summary.go:buildNOMSummary`) showed only counts and elapsed time. The InlineRenderer had `~Xm left` support but the TUI didn't. Status report P2 item #12.

**What shipped:**

- **`buildNOMSummary`** gained a `remaining time.Duration` parameter — appends `| ~2m left` when positive
- **`estimatedRemaining()` helper** on `ProgressModel` (`render_nom.go`) — delegates to `nomSubscriber.EstimatedTotalRemaining()`, returns 0 if no subscriber
- **`renderNOMSummaryBar`** updated to pass the subscriber's remaining estimate
- The TUI now shows `~Xm left` **automatically** when activities have estimates — no external wiring needed

**Files:** `tui/summary.go`, `tui/render_nom.go`, `tui/summary_test.go`

**Tests:** `TestTUISummaryShowsRemaining`, `TestBuildNOMSummary` (updated with remaining-estimate subtest)

### 4. TUI Already Renders Progress/Retry — Report Correction + Proof Tests

**Problem:** The 12:33 status report item #6 claimed the TUI needed view-layer updates to render Progress/RetryCount. This was **incorrect** — the TUI delegates line rendering to `nom.RenderVisibleEntry` → `formatActivityLabel`, which already renders both fields.

**What shipped:**

- **3 proof tests** in new file `tui/render_nom_progress_test.go`:
  - `TestTUIRendersProgress` — verifies `→ Tidying [2/26]` appears in TUI tree output
  - `TestTUIRendersRetry` — verifies `⟳` symbol and `(timeout)` reason appear in TUI tree output
  - `TestTUISummaryShowsRemaining` — verifies `~Xm left` in TUI summary bar
- **AGENTS.md pattern documented**: "TUI NOM rendering delegates to nom" — explains the delegation chain so future agents don't重复 this mistake

**Files:** `tui/render_nom_progress_test.go` (NEW), `AGENTS.md`

### 5. Concurrency & Lifecycle Safety Tests

**Problem:** The new progress/retry event handlers mutate shared state (`Progress`, `RetryCount`, `RetryReason`, counts cache) under the subscriber's write lock. The previous session had unit tests but no concurrent stress test. Status report P3 items #19, #20.

**What shipped (new file `nom/progress_lifecycle_test.go`, 5 tests):**

| Test                                         | Verifies                                                                                                                                    |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------- |
| `TestMultiSubscriber_ProgressAndRetryFanout` | Progress + retry events propagate to every subscriber behind MultiSubscriber                                                                |
| `TestResetClearsProgressAndRetry`            | `Reset()` wipes Progress, RetryCount, RetryReason, and zeroes counts                                                                        |
| `TestProgressClearedOnRetry`                 | Retry's `SetRunning()` clears stale progress (fresh start semantics)                                                                        |
| `TestConcurrentProgressAndRetry`             | 8 goroutines × 100 ops each: progress/fail/retry + concurrent reads. Validates counts cache == brute-force recount after storm. Race-clean. |
| `TestEstimatedTotalRemainingRunningElapsed`  | Running activities contribute `estimate - elapsed`, not full estimate                                                                       |

### 6. Documentation & Examples

**AGENTS.md** — 4 pattern/gotcha entries updated or added:

- Sealed Event count corrected: "7 concrete structs" → "9 concrete structs" (listed all 9 by name)
- New pattern: **Progress/Retry events (v0.21.0)** — documents auto-clear semantics, persistence rules, direct accessors
- New pattern: **EstimatedTotalRemaining is subscriber-owned (v0.21.0)** — single source of truth, renderer delegation
- New pattern: **TUI NOM rendering delegates to nom** — prevents the repeated "TUI needs view-layer updates" mistake
- Gotcha updated: NOM events now include `ActivityProgress` and `ActivityRetrying` examples

**CHANGELOG.md** — `[Unreleased]` section with full feature breakdown, changed APIs, and test count

**`examples/nom_progress/main.go`** — demonstrates all three new capabilities:

- `ActivityProgress` sub-step message
- `ActivityRetrying` with reason
- `SetEstimatedTime` + `EstimatedTotalRemaining()`

**`event.go` doc comment** — Added throttling guidance on `ActivityProgress` (recommended ~1/sec cadence, time-based guard pattern)

### 7. Verification — All Green

| Check                             | Result                                                   |
| --------------------------------- | -------------------------------------------------------- |
| `nix run .#build` (18 modules)    | ✅ All build clean                                       |
| `nix run .#test` (18 modules)     | ✅ All pass                                              |
| `nix run .#lint` (18 modules)     | ✅ 0 issues everywhere                                   |
| `nix run .#test-race` (nom + tui) | ✅ Race clean, no regressions                            |
| nom test count                    | 183 PASS (was 168 before session 1, 173 after session 1) |
| tui test count                    | 80 PASS (was 77)                                         |
| Golden files                      | ✅ Not broken (no existing output format changes)        |

---

## b) PARTIALLY DONE 🟡

### 1. Golden File Snapshots for New Event Types

The new `Progress`, `RetryCount`, `RetryReason` fields are additive — existing golden tests pass unchanged. However, **new golden files** showing rendered output with progress (`→ message`), retry (`⟳N (reason)`), and `~Xm left` would serve as visual documentation for contributors. The 13 `strings.Contains`-based tests verify correctness but don't serve as visual reference.

### 2. Structured Progress Type (Deliberately Deferred)

`Progress` is currently a bare `string`. A structured type (`ProgressDetail{Current, Total int; Label string}`) would enable progress-bar rendering (`[████░░░░░░] 3/26`). This was deliberately deferred per YAGNI — no current consumer needs it, and the report's own "top question" acknowledges it's a non-breaking addition that can come later. The bare string works for BuildFlow's current use case.

---

## c) NOT STARTED ⬜

1. **Golden test snapshots** — visual reference files in `nom/testdata/` showing progress/retry/remaining rendered output
2. **Structured progress type** — `ProgressDetail{Current, Total, Label}` for progress-bar rendering (deferred until a consumer needs it)
3. **Progress throttling mechanism** — currently the caller controls cadence; a built-in debounce (e.g. 1/sec max) would prevent excessive redraws from fast-iterating callers. Documented as guidance only.
4. **Fuzz test for progress events** — rapid progress/retry/complete sequences to surface edge cases
5. **Benchmark progress event path** — ensure no latency under high-frequency progress updates
6. **Adaptive tree pruning** — nom-style "fill 1/3 of terminal, prune low-priority" (P4 research)
7. **Explore `aymanbagabas/go-udiff` for frame diffing** (P4 research)
8. **Explore `charmbracelet/x/term` as replacement for `golang.org/x/term`** (P4 research)

---

## d) TOTALLY FUCKED UP 💥 (Honest Assessment)

### Nothing

All changes compile, all 18 modules pass, 0 lint issues, 0 regressions, race-clean. The implementation follows existing patterns exactly (sealed event types, snapshot threading, handler under write lock, counts cache maintenance, two-mutex renderer model).

### Minor mistakes caught and fixed during the session:

1. **Test event ordering bug** — Initial `TestMultiSubscriber_ProgressAndRetryFanout` sent progress _before_ retry, but retry's `SetRunning()` clears progress. Fixed by reordering events (progress last).
2. **`TestResetClearsProgressAndRetry` same issue** — Same root cause. Fixed by sending progress after retry.
3. **Lint formatting** — Two gci/golines issues in test struct alignment and a long line in the concurrent test assertion. Fixed by aligning struct fields and simplifying the assertion to use `ActivityCounts` value equality.
4. **TUI `summary_test.go` struct misalignment** — Added `remaining` field without aligning the struct. Fixed by running through gci/golines alignment.

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Architecture

1. **Golden tests for visual documentation** — The `testdata/` directory should include rendered examples of progress and retry output so contributors can see what the output looks like without running code. The `strings.Contains` tests verify correctness but not visual quality.

2. **Progress message should support structured sub-step progress** — A bare string can't render a progress bar. Consider `ProgressDetail struct { Current, Total int; Label string }` alongside `Progress string` so renderers can show `[████░░░░░░] 3/26` instead of text. The nom reference implementations (nix-output-monitor, nh) show progress bars — nom should support that path.

3. **Progress throttling guidance vs. mechanism** — Currently only documented as a recommendation (~1/sec). A built-in time-based debounce in `handleActivityProgress` would make excessive-update safety the default rather than opt-in. But it adds complexity and a mutable `lastProgressUpdate` field. Tradeoff: safety vs. simplicity.

4. **Retry count on non-running activities** — `RetryCount` persists across state transitions (only incremented, never cleared). A completed activity that was retried shows `⟳2` forever. This is probably correct (you want to know it was retried), but it's a presentation choice worth being aware of.

5. **EstimatedTotalRemaining accuracy** — The sum assumes activities run sequentially. If activities run in parallel (which nom supports via dependency edges), the actual remaining time could be less than the sum. A max-of-critical-path computation would be more accurate but significantly more complex. The sum is a useful upper bound.

### Code Quality

6. **Multi-line activity label edge cases** — `formatActivityLabel` now conditionally appends `\n` + dim sub-line for progress. This multi-line output interacts correctly with `PhysicalLineCount` and `buildRedrawOutput` (they split on `\n`), but it's the first multi-line activity label. Worth watching for edge cases with `maxHeight` truncation — a progress sub-line could get orphaned from its parent if the height cap cuts between them.

7. **TUI delegation pattern should be more prominent** — The fact that TUI tree rendering delegates to nom's `formatActivityLabel` (and thus gets new fields for free) is a significant architectural property. It's now documented in AGENTS.md, but a diagram or more explicit comment in `render_nom.go` would help.

---

## f) Top 25 Things to Do Next

| #  | Priority | Task                                                                                                                                          |
| -- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | **P0**   | **Bump go-output to v0.21.0 in BuildFlow** — push tags, update flake input, remove local replace directives                                   |
| 2  | **P0**   | **Wire `ActivityProgress` events from BuildFlow** — call `subscriber.OnEvent(ctx, nom.ActivityProgress{...})` from `ForEachGoModule` callback |
| 3  | **P0**   | **Wire `ActivityRetrying` events from BuildFlow** — call from retry path in `ProgressBridge.stepRetrying()` with `Reason` field               |
| 4  | **P0**   | **Wire `SetEstimatedTime` from BuildFlow** — inject SQLite P50 estimates after loading from `dbstore`                                         |
| 5  | **P0**   | **Wire `SetEstimatedRemainingFunc` or `EstimatedTotalRemaining()` from BuildFlow** — connect to the renderer or let TUI auto-consume          |
| 6  | **P1**   | **Golden test snapshots** — add rendered examples of progress/retry/reason/remaining output to `nom/testdata/`                                |
| 7  | **P1**   | **Fuzz test for progress events** — rapid progress/retry/complete sequences via Go fuzzing                                                    |
| 8  | **P1**   | **Benchmark progress event path** — ensure no latency under high-frequency progress updates (1000/sec scenario)                               |
| 9  | **P2**   | **Structured progress type** — `ProgressDetail{Current, Total, Label}` for progress-bar rendering                                             |
| 10 | **P2**   | **Progress bar rendering** — `▕████░░░░▏ 45%` style when structured progress is populated                                                     |
| 11 | **P2**   | **Progress throttling mechanism** — built-in debounce in `handleActivityProgress` (1/sec default)                                             |
| 12 | **P2**   | **Clear RetryCount on fresh workflow** — verify `Reset()` behavior is tested (it clears activities entirely — done in this session)           |
| 13 | **P2**   | **Retry reason persistence across completions** — should a completed-then-retried step keep its reason forever? Currently yes.                |
| 14 | **P2**   | **EstimatedTotalRemaining parallel-accuracy** — compute max-of-critical-path instead of sum for parallel activities                           |
| 15 | **P2**   | **Render remaining estimate in universal (non-NOM) TUI mode** — currently only NOM mode shows `~Xm left`                                      |
| 16 | **P3**   | **Progress during retry** — when retried, should progress from prior attempt persist? Currently cleared (tested in this session)              |
| 17 | **P3**   | **Documentation: update `docs/DOMAIN_LANGUAGE.md`** — add `Progress`, `RetryCount`, `RetryReason`, `EstimatedTotalRemaining` to the glossary  |
| 18 | **P3**   | **Documentation: update `FEATURES.md`** — add the new event types and APIs to the feature inventory                                           |
| 19 | **P3**   | **Multi-subscriber with inline renderer test** — verify progress events flow through MultiSubscriber to an InlineRenderer-bound subscriber    |
| 20 | **P3**   | **Examples: `tui_progress` update** — demonstrate progress/retry events in the TUI example (not just inline)                                  |
| 21 | **P3**   | **Progress sub-line truncation with maxWidth** — verify long progress messages truncate correctly                                             |
| 22 | **P4**   | **Explore `aymanbagabas/go-udiff` for frame diffing** (from BuildFlow backlog)                                                                |
| 23 | **P4**   | **Explore `charmbracelet/x/term` as replacement for `golang.org/x/term`**                                                                     |
| 24 | **P4**   | **Adaptive tree pruning** — nom-style "fill 1/3 of terminal, prune low-priority"                                                              |
| 25 | **P4**   | **Progress as structured event log** — emit progress/retry events to an append-only log for post-run analysis                                 |

---

## g) Top #1 Question I Cannot Answer Myself

**Should `EstimatedTotalRemaining()` use a simple sum or a critical-path computation for parallel activities?**

The current implementation sums all unfinished activity estimates. This is correct for sequential execution but overestimates when activities run in parallel (which nom supports via dependency edges). For example, if two independent pending activities each estimate 60s, the sum is 120s — but if they run in parallel, the actual remaining time is ~60s.

**Arguments FOR the simple sum (current approach):**

- Trivially correct, O(n) scan, no graph traversal needed
- Provides a useful _upper bound_ — "at most 2m left"
- BuildFlow's steps are largely sequential (phase-by-phase), so the sum is close to accurate
- The dependency graph may not reflect actual execution parallelism (a dependency edge means "must complete before," not "runs immediately after")

**Arguments FOR critical-path computation:**

- More accurate for parallel builds — matches what nom reference implementations (nix-output-monitor) show
- The subscriber already has the dependency tree — the data is there
- Would show `~1m left` instead of `~2m left` for two parallel 60s steps

**I cannot determine the answer** because it depends on BuildFlow's execution model:

- If BuildFlow runs phases sequentially (phase 1 completes, then phase 2 starts), the sum is correct
- If BuildFlow runs independent steps in parallel within a phase, the sum overestimates
- A hybrid (sum within each phase, max across phases) might be the right middle ground

The sum is shipped and working. A critical-path variant can be added as a non-breaking alternative method (`EstimatedTotalRemainingCriticalPath()`) when a consumer needs the more accurate computation.

---

## Session Statistics

| Metric                       | Value                                                                   |
| ---------------------------- | ----------------------------------------------------------------------- |
| Production files changed     | 10 (nom: 6, tui: 3, examples: 1)                                        |
| Documentation files changed  | 2 (AGENTS.md, CHANGELOG.md)                                             |
| New test files               | 2 (`nom/progress_lifecycle_test.go`, `tui/render_nom_progress_test.go`) |
| New tests (nom)              | 10 (5 lifecycle + 4 retry/remaining + 1 running-elapsed)                |
| New tests (tui)              | 3 (progress render, retry render, summary remaining)                    |
| Total new tests this session | 13                                                                      |
| Insertions                   | 275                                                                     |
| Deletions                    | 16                                                                      |
| All modules passing          | 18/18                                                                   |
| Lint issues                  | 0                                                                       |
| Race test status             | ✅ Clean (nom + tui)                                                    |
| nom total test count         | 183 PASS                                                                |
| tui total test count         | 80 PASS                                                                 |
| New public APIs              | 2 (`EstimatedTotalRemaining()`, `Reason` field on `ActivityRetrying`)   |
| New event fields             | 2 (`RetryReason` on Activity/Snapshot, `Reason` on ActivityRetrying)    |
| Report items closed          | 6 of 8 P1/P2 in-repo items (#6, #7, #8, #9, #11, #12)                   |
| Report items corrected       | 1 (#6 was already done — proved with tests)                             |
