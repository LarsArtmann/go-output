# Status: go-output/nom BuildFlow Integration Enhancements

**Date:** 2026-07-01 12:33 CEST
**Session Focus:** Add ActivityProgress + ActivityRetrying events, external estimate injection API, and `~Xm left` summary bar to nom — closing the top gaps from BuildFlow's status reports
**Overall Health:** ✅ All 18 modules passing · 0 lint issues · 12 new tests · 0 regressions

---


> **✅ Resolved (2026-08-04):**
>
> All in-repo items closed by session 2 (18:40): retry reason display (`⟳2 (timeout)`), subscriber-owned `EstimatedTotalRemaining()`, TUI summary bar wired, concurrency/race tests. Critical-path ETA implemented in v0.23.0 (`EstimatedCriticalPathRemaining`). Golden test snapshots for progress/retry exist. Remaining genuinely open: structured progress type (ROADMAP), BuildFlow wiring (external).

---

## Executive Summary

This session addressed the top gaps BuildFlow reported across 5 status documents (dated 2026-07-01). The core problem: nom lacked first-class APIs for four capabilities BuildFlow's progress UI desperately needs — **sub-step progress** (multi-module steps are invisible black boxes), **retry visibility** (retried steps show no live feedback), **external estimate injection** (BuildFlow's SQLite timing store can't feed the tree), and **estimated-remaining display** (the `~Xm left` feature). All four are now implemented as sealed event types and typed public APIs, fully tested, with zero lint issues.

**10 production files changed, 1 test file added, 175 insertions, 2 deletions. 12 new tests, all passing.**

---

## a) FULLY DONE ✅

### 1. ActivityProgress Event — Sub-Step Visibility (BuildFlow's #1 Gap)

**Problem:** BuildFlow's multi-module operations (e.g. `go-mod-tidy` iterating 26 Go modules) are invisible in the inline NOM tree. A long-running step appears as a single line with no indication of what's happening inside. BuildFlow added `StepProgress(name, message)` to their `ProgressBridge` but nom had no event type to receive it — the `inlineRendererAdapter.stepProgress()` was a no-op that just called `Refresh()`.

**What shipped:**

- **New sealed event type** `ActivityProgress{ID, Name, Message}` in `event.go` — part of the sealed sum type (unexported `isEvent()` marker), so the compiler guarantees exhaustive type-switch handling
- **Handler** `handleActivityProgress` in `subscriber_handlers.go` — sets the progress message on the activity
- **`Progress string` field** on `Activity` (`activity.go`) and `ActivitySnapshot` (`activity_snapshot.go`) — threaded through the snapshot path for race-free rendering
- **Auto-clear semantics**: `SetRunning()`, `SetCompleted()`, `SetFailed()` all clear `Progress` — terminal state has no sub-step, and a fresh retry starts with no stale message
- **Direct accessor** `SetActivityProgress(id, message)` on `NOMStyleSubscriber` for callers that don't use the event dispatch path
- **Rendering** in `tree_render.go`: progress renders as a dim `→ Tidying module [2/26]` sub-line beneath the activity label (only while `Status == Running`)
- **New symbols**: `SymbolProgress = "→"` in `symbols.go`

**Files:** `event.go`, `subscriber_handlers.go`, `activity.go`, `activity_snapshot.go`, `activity_accessors.go`, `tree_render.go`, `symbols.go`

**BuildFlow integration path:**

```go
// In ForEachGoModule callback:
_ = subscriber.OnEvent(ctx, nom.ActivityProgress{
    ID: stepID, Name: stepName,
    Message: "Tidying module [2/26]: modules/gitignore",
})
// Tree now shows:
// ⏵ go-mod-tidy 4.2s
//   → Tidying module [2/26]: modules/gitignore
```

### 2. ActivityRetrying Event — Retry Visibility

**Problem:** When a step fails and retries, BuildFlow's `ProgressBridge.stepRetrying()` had no nom event to send. The activity stayed visually "failed" until it completed. BuildFlow's status reports listed "Add retry live flush" as a repeated P2 item.

**What shipped:**

- **New sealed event type** `ActivityRetrying{ID, Name, Attempt}` in `event.go`
- **Handler** `handleActivityRetrying` in `subscriber_handlers.go` — transitions the activity from `Failed` back to `Running`, increments `RetryCount`, and updates the counts cache (decrements Failed, increments Running) via `applyCountsDelta`
- **`RetryCount int` field** on `Activity` and `ActivitySnapshot`
- **Rendering** in `tree_render.go`: renders `⟳2` suffix in Info color (cyan) when `RetryCount > 0`
- **New symbol**: `SymbolRetrying = "⟳"` in `symbols.go`
- **Attempt tracking**: uses `max(e.Attempt, activity.RetryCount+1)` to handle both explicit-attempt and incremental-count callers

**Files:** `event.go`, `subscriber_handlers.go`, `activity.go`, `activity_snapshot.go`, `tree_render.go`, `symbols.go`

### 3. External Estimate Injection API

**Problem:** BuildFlow has a SQLite timing store with P50 durations per step, but nom's only estimate path is the internal CSV-based `TimingCache`. BuildFlow's `inlineRendererAdapter.setEstimator()` was documented as a no-op — there was no API to inject external estimates.

**What shipped:**

- **`SetEstimatedTime(activityID, estimated time.Duration)`** on `NOMStyleSubscriber` in `activity_accessors.go` — directly sets the predicted duration on an activity, sourced from any external store. The estimate renders as `∅4.2s` on pending activities in the tree.

**Files:** `activity_accessors.go`

**BuildFlow integration path:**

```go
// After loading P50 estimates from SQLite:
for stepID, p50 := range estimates {
    subscriber.SetEstimatedTime(stepID, p50)
}
```

### 4. `~Xm left` in Summary Bar

**Problem:** BuildFlow's `TimeEstimator.EstimatedTotalRemaining()` computes the projected remaining time for all pending steps, but there was no rendering path to show it. The inline renderer's summary bar showed only counts and elapsed time. BuildFlow listed "Show `~Xm left`" as a P1 item across 3 status reports.

**What shipped:**

- **`SetEstimatedRemainingFunc(fn func() time.Duration)`** on `InlineRenderer` in `inline_renderer.go` — a callback-based mechanism that decouples the estimate source from the renderer. The callback is invoked once per frame under `renderMu` (not `tickMu`), so it must not acquire renderer locks.
- **Summary bar rendering** in `inline_renderer_summary.go`: when the callback returns > 0, appends `~2m left` after the elapsed-time segment. Shows nothing when nil or returning 0.

**Files:** `inline_renderer.go`, `inline_renderer_summary.go`

**BuildFlow integration path:**

```go
renderer.SetEstimatedRemainingFunc(func() time.Duration {
    return estimator.EstimatedTotalRemaining()
})
// Summary bar now shows:
// ╭──────────────────────────────╮
// │ ⏵2 ✔5 ○3 12.4s ~2m left ∑10 (70%) │
// ╰──────────────────────────────╯
```

### 5. Test Suite — 12 New Tests, All Passing

**File:** `nom/progress_events_test.go` (NEW — 344 LOC)

| Test                                   | Verifies                                                        |
| -------------------------------------- | --------------------------------------------------------------- |
| `TestActivityProgressEvent`            | Progress event sets message on activity + snapshot              |
| `TestActivityProgressClearOnComplete`  | Completing clears stale progress                                |
| `TestActivityProgressEmptyMessage`     | Empty message clears prior progress                             |
| `TestActivityProgressSetDirect`        | Direct accessor `SetActivityProgress` works                     |
| `TestActivityRetryingEvent`            | Retry transitions failed→running, count=1, counts cache updated |
| `TestActivityRetryingMultipleAttempts` | Count increments correctly across 3 retries                     |
| `TestSetEstimatedTimeDirect`           | External estimate injection sets `EstimatedTime`                |
| `TestEstimatedRemainingInSummary`      | Summary bar contains `~2m left` when callback set               |
| `TestEstimatedRemainingZero`           | No `left` shown when callback returns 0                         |
| `TestEstimatedRemainingNil`            | No `left` shown when callback is nil                            |
| `TestProgressRendersInTree`            | Tree output contains progress message                           |
| `TestRetryRendersInTree`               | Tree output contains `⟳` symbol                                 |

### 6. Event Uniqueness Guards Updated

`event_test.go` updated to include `EventActivityProgress` and `EventActivityRetrying` in both the uniqueness test (`TestEventConstantsUnique`) and the non-empty test (`TestEventConstantsNonEmpty`).

---

## Verification Results

| Check                             | Result                                                      |
| --------------------------------- | ----------------------------------------------------------- |
| `go build ./...` (nom)            | ✅ Clean                                                    |
| `go vet ./...` (nom)              | ✅ Clean                                                    |
| `go test ./...` (nom)             | ✅ 0.269s, all pass                                         |
| `nix run .#test` (all 18 modules) | ✅ All pass                                                 |
| `nix run .#lint` (all 18 modules) | ✅ 0 issues across all modules                              |
| Race tests                        | ✅ Pass (no regressions)                                    |
| Golden files                      | ✅ Not broken (no output format changes for existing paths) |

---

## b) PARTIALLY DONE 🟡

### 1. Golden File Regeneration

The new `Progress` and `RetryCount` fields are additive — existing golden tests pass unchanged because no existing test sends the new events. However, **new golden files should be added** showing the rendered output with progress and retry annotations for documentation purposes. The tests verify via `strings.Contains` which is sufficient for correctness but doesn't serve as visual reference.

### 2. TUI Module Integration

The `tui/` module was not modified. The new `ActivityProgress`/`ActivityRetrying` events flow through the sealed `OnEvent` type switch, so if the TUI's subscriber receives them, they'll be stored. But the TUI's own rendering code (`tui/view.go`, `tui/summary.go`) doesn't yet render the `Progress` or `RetryCount` fields — it reads `ActivitySnapshot` but only renders Label/Symbol/Color/Elapsed. The fields are in the snapshot, ready to be consumed.

### 3. Examples Not Updated

`examples/nom_progress/` doesn't demonstrate the new events. BuildFlow's integration will be the first real consumer.

---

## c) NOT STARTED ⬜

1. **TUI rendering of Progress/RetryCount** — the TUI view layer needs to read the new snapshot fields and render them. The data path is complete (snapshot → TUI); only the view layer needs updating.
2. **Golden test snapshots for new event types** — visual reference files showing rendered output.
3. **Multi-subscriber fan-out test for progress events** — verify progress propagates through `MultiSubscriber`.
4. **Progress throttling** — no rate limiting on progress message updates; a caller sending 1000 updates/sec would cause 1000 renders. A debounce would help, but the caller (BuildFlow) controls cadence.
5. **Progress as structured data** — currently `Progress` is a bare `string`. A structured type (e.g. `{Current, Total int, Label string}`) would enable rendering a progress bar instead of text. Defer until a consumer needs it.
6. **EstimatedRemaining via subscriber** — the `~Xm left` callback is on the renderer, not the subscriber. A subscriber-level API (e.g. `SetEstimatedRemainingFunc` on `NOMStyleSubscriber`) would let the subscriber compute it from activity estimates. Currently the caller must compute the sum externally.

---

## d) TOTALLY FUCKED UP 💥 (Honest Assessment)

### Nothing

All changes compile, all 18 modules pass, 0 lint issues, 0 regressions. The implementation follows existing patterns exactly (sealed event types, snapshot threading, handler under write lock, counts cache maintenance).

---

## e) WHAT WE SHOULD IMPROVE 🔧

### Architecture

1. **TUI module should consume Progress/RetryCount** — the data is in the snapshot; `tui/view.go` should render it. This would make the TUI mode as feature-complete as inline mode for these new capabilities.

2. **Progress message should support structured sub-step progress** — a bare string can't render a progress bar. Consider `ProgressDetail struct { Current, Total int; Label string }` alongside or instead of `Progress string` so renderers can show `[████░░░░░░] 3/26` instead of text.

3. **EstimatedRemaining should be computable from subscriber state** — the subscriber has all activity estimates (via `SetEstimatedTime` or timing cache). A `subscriber.EstimatedTotalRemaining()` method would let the renderer compute `~Xm left` internally, eliminating the callback. The callback approach was chosen first because BuildFlow already has this computation in `TimeEstimator` — but the subscriber could own it.

4. **Golden tests for visual documentation** — the `testdata/` golden files should include examples of progress and retry rendering so contributors can see what the output looks like without running code.

### Code Quality

5. **Progress sub-line rendering** — the `formatActivityLabel` function now conditionally appends a `\n` + dim sub-line. This multi-line output interacts with `PhysicalLineCount` and `buildRedrawOutput` correctly (they split on `\n`), but it's the first multi-line activity label. Worth watching for edge cases with `maxHeight` truncation.

6. **Retry count on non-running activities** — `RetryCount` persists across state transitions (it's only incremented, never cleared). A completed activity that was retried shows `⟳2` forever. This is probably correct (you want to know it was retried), but it's a presentation choice that could be debated.

---

## f) Top 25 Things to Do Next

| #   | Priority | Task                                                                                                                                          |
| --- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **P0**   | **Bump go-output to v0.21.0 in BuildFlow** — push tags, update flake input, remove local replace directives                                   |
| 2   | **P0**   | **Wire `ActivityProgress` events from BuildFlow** — call `subscriber.OnEvent(ctx, nom.ActivityProgress{...})` from `ForEachGoModule` callback |
| 3   | **P0**   | **Wire `ActivityRetrying` events from BuildFlow** — call from retry path in `ProgressBridge.stepRetrying()`                                   |
| 4   | **P0**   | **Wire `SetEstimatedTime` from BuildFlow** — inject SQLite P50 estimates after loading from `dbstore`                                         |
| 5   | **P0**   | **Wire `SetEstimatedRemainingFunc` from BuildFlow** — connect `TimeEstimator.EstimatedTotalRemaining()` to the renderer                       |
| 6   | **P1**   | **TUI rendering of Progress/RetryCount** — update `tui/view.go` to render the new snapshot fields                                             |
| 7   | **P1**   | **Golden test snapshots** — add rendered examples of progress/retry output to `testdata/`                                                     |
| 8   | **P1**   | **Multi-subscriber fan-out test** — verify progress events propagate through `MultiSubscriber`                                                |
| 9   | **P1**   | **Progress throttling guidance** — document recommended cadence for progress updates (e.g. 1/sec)                                             |
| 10  | **P2**   | **Structured progress type** — `ProgressDetail{Current, Total, Label}` for progress-bar rendering                                             |
| 11  | **P2**   | **Subscriber-level `EstimatedTotalRemaining()`** — compute from activity estimates internally                                                 |
| 12  | **P2**   | **Render remaining estimate in TUI** — wire `~Xm left` into `tui/summary.go`                                                                  |
| 13  | **P2**   | **Retry reason display** — `ActivityRetrying` could carry a `Reason string` for `⟳2 (timeout)`                                                |
| 14  | **P2**   | **Progress bar rendering for structured progress** — if structured type is added                                                              |
| 15  | **P2**   | **Clear RetryCount on fresh workflow** — verify `Reset()` clears it (it clears activities entirely, so this works, but worth testing)         |
| 16  | **P3**   | **Examples update** — `examples/nom_progress/` demonstrate progress/retry events                                                              |
| 17  | **P3**   | **Fuzz test for progress events** — rapid progress/retry/complete sequences                                                                   |
| 18  | **P3**   | **Benchmark progress event path** — ensure no latency under high-frequency progress updates                                                   |
| 19  | **P3**   | **Progress during retry** — when retried, should progress from prior attempt persist? Currently cleared by SetRunning                         |
| 20  | **P3**   | **Concurrent progress + retry** — race test sending progress and retry simultaneously                                                         |
| 21  | **P3**   | **Documentation** — update `AGENTS.md` patterns section with new event types                                                                  |
| 22  | **P3**   | **CHANGELOG entry** for v0.21.0                                                                                                               |
| 23  | **P4**   | **Explore `aymanbagabas/go-udiff` for frame diffing** (from BuildFlow backlog)                                                                |
| 24  | **P4**   | **Explore `charmbracelet/x/term` as replacement for `golang.org/x/term`**                                                                     |
| 25  | **P4**   | **Adaptive tree pruning** — nom-style "fill 1/3 of terminal, prune low-priority"                                                              |

---

## g) Top #1 Question I Cannot Answer Myself

**Should `ActivityProgress.Message` remain a bare `string`, or should we add a structured progress type?**

The current `Progress string` field enables text-based sub-step messages like `"Tidying module [2/26]: modules/gitignore"`. This is maximally flexible and trivially integrated. BuildFlow can start using it immediately.

However, nom's reference implementations (nix-output-monitor, nh) show **progress bars** for download/build progress — `▕████░░░░▏ 45%`. A bare string can't render that; it requires structured data (`{Current, Total int}`). We could add a parallel `ProgressDetail` field:

```go
type ProgressDetail struct {
    Current int
    Total   int
    Label   string
}
```

Rendered as `→ Label [████░░░░░░] 3/26` when populated, falling back to `Progress string` when not.

**Arguments FOR structured type:**

- Enables progress-bar rendering (the nom signature visual)
- Self-documenting (no string parsing needed)
- Type-safe (impossible to represent "3 of 26" inconsistently)

**Arguments AGAINST:**

- Over-engineering if no consumer needs it yet
- The string approach works for BuildFlow's current use case
- Adds API surface that must be maintained

**I cannot determine the answer** because it depends on whether BuildFlow (or future consumers) will want visual progress bars in the tree, or whether text messages suffice. The text approach is shipped and working. The structured type can be added later as a non-breaking addition (new field alongside `Progress string`).

---

## Session Statistics

| Metric                   | Value                                                                      |
| ------------------------ | -------------------------------------------------------------------------- |
| Production files changed | 10                                                                         |
| Test files added         | 1                                                                          |
| New event types          | 2 (`ActivityProgress`, `ActivityRetrying`)                                 |
| New public APIs          | 3 (`SetActivityProgress`, `SetEstimatedTime`, `SetEstimatedRemainingFunc`) |
| New symbols              | 2 (`SymbolProgress →`, `SymbolRetrying ⟳`)                                 |
| New tests                | 12                                                                         |
| Insertions               | 175                                                                        |
| Deletions                | 2                                                                          |
| All modules passing      | 18/18                                                                      |
| Lint issues              | 0                                                                          |
| BuildFlow gaps closed    | 4 of top 8 P0/P1 items from 2026-07-01 status reports                      |

