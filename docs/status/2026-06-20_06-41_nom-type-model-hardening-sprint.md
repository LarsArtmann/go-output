# Status Report — 2026-06-20 06:41

**Branch:** `master` @ `3f60328` (pushed)
**Session scope:** Type-model hardening sprint on `nom/` — 3 refactors, comparison analysis vs upstream nix-output-monitor

---

## Executive Summary

The session began with a comparative analysis of the local `nom/` package against upstream `nix-output-monitor` and `nh`'s integration pattern. What started as a "how do we make it better" discussion uncovered three type-model weaknesses: a string-prefix phase convention, a string-dispatched open event interface, and a domain type leaking framework concerns. All three have been refactored and pushed. **The codebase builds clean across all 20 modules, all tests pass (141 in `nom/` alone), and `-race` is clean.**

**Net code reduction:** −259 lines (525 added, 784 removed) across 30 files — the sealed-event refactor alone deleted 398 net lines of accessor boilerplate.

---

## a) FULLY DONE (this session)

### 1. Comparative analysis: `nom/` vs upstream nix-output-monitor vs `nh`

Documented how upstream NOM is a standalone Haskell binary consuming Nix's `internal-json` wire format, while our `nom/` is an in-process Go library driven by Go events. The two share the _visual presentation layer_ (symbols, timing cache, priority sort, inline ANSI redraw) but differ fundamentally in the input/data layer. Confirmed the design is a correct, honest abstraction — not a reimplementation. Report delivered in the conversation.

### 2. `ActivityKind` enum replaces `"phase:"` ID-prefix convention (`6ac679f`)

**Before:** Phase detection was driven by `strings.HasPrefix(id, "phase:")` — a lying name that claimed to signal kind but was really an opaque string matched at render time. Any user activity literally named `"phase:cleanup"` would wrongly render as a phase.

**After:** Typed `ActivityKind` enum (`ActivityKindTask`, `ActivityKindPhase`) set at construction via `NewPhase()`, threaded through `ActivitySnapshot`, and read via `snap.IsPhase()`. The `"phase:"` magic string is entirely gone from the codebase (only an explanatory comment in `activity_kind.go` references it as historical context).

| Change                             | Detail                                                                    |
| ---------------------------------- | ------------------------------------------------------------------------- |
| New file                           | `nom/activity_kind.go` (66 lines)                                         |
| `Activity.Kind` field              | Set at construction, never changes                                        |
| `NewPhase(id, name)` constructor   | Creates a Phase-kind activity                                             |
| `ActivitySnapshot.Kind`            | Threaded through for render-time phase detection                          |
| `KindAccessor` → deleted           | Replaced by direct field on `ActivityRegistered`/`ActivityStarted` events |
| `isPhaseID()` → deleted            | Was the HasPrefix smell                                                   |
| `ActivityNode.IsPhase()` → deleted | Had zero callers (dead code)                                              |
| Test helpers                       | Added `registerPhase()` and `snapshotBuilder.setPhase()`                  |

### 3. Sealed `Event` sum type replaces string-dispatched open interface (`db8ab6f`)

**Before:** `Event` was an open interface with `GetEventType() string`, dispatched via a string `switch` in `OnEvent`. Callers had to implement 5 optional accessor interfaces (`WorkflowEventAccessor`, `ActivityEventAccessor`, `DurationAccessor`, `ErrorAccessor`, `DependenciesAccessor`). This made silent typos in event routing possible ("workflow.strted" → silent no-op). Seven near-identical local event structs existed across nom/tui/integration/examples (`testEvent`, `nomTestEvent`, `benchEvent`, `hostDownloadEvent`, `diagramTestEvent`, `workflowEvent`, `diagramEvent`, `smokeTestEvent`, `multiTestEvent`, `minimalEvent`) — each reimplementing the same accessors.

**After:** `Event` is sealed: 7 concrete structs (`WorkflowStarted`, `WorkflowCompleted`, `WorkflowFailed`, `ActivityStarted`, `ActivityRegistered`, `ActivityCompleted`, `ActivityFailed`) with an unexported `isEvent()` marker that prevents external implementations. `OnEvent` dispatches via exhaustive Go type switch. All accessor interfaces are deleted; handlers read event fields directly. The silent-typo failure mode is now unrepresentable — the compiler rejects unhandled event types.

| Change                                | Detail                                                                                                                        |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `event.go` rewritten                  | Sealed interface + 7 concrete event types                                                                                     |
| `subscriber_handlers.go` rewritten    | Type switch + direct field access (281 → 175 lines)                                                                           |
| All 10 local event structs deleted    | Across nom, tui, integration, examples                                                                                        |
| `EventWorkflowStarted` etc. constants | Kept for logging/metrics IDs (no longer used for dispatch)                                                                    |
| `event_test.go` comment updated       | Now guards against duplicate _names_ (not silent misrouting)                                                                  |
| 2 obsolete tests deleted              | `TestNOMStyleSubscriber_UnknownEventType` and `_EventWithoutAccessor` — the failure modes they guarded are now compile errors |

**Verified safe before executing:** `GetEventType()` had zero external consumers; accessor interfaces had zero production implementors (only tests/examples). ADR 006 explicitly marks NOM APIs experimental pre-v1.0.

### 4. Decouple `Activity` from `output.GraphNode` (`3f60328`)

**Before:** `Activity` embedded `output.GraphNode`, so it carried `Shape`, `Style`, `Metadata` — diagram-export framework concerns. `applyVisualStyle()` mutated these graph fields on every lifecycle transition (SetRunning/SetCompleted/SetFailed), violating the project's own AGENTS.md principle: _"Leaky implementation details — framework concerns in domain types. Model them separately."_

**After:** `Activity` is a pure domain type: `ID` and `Label` are direct branded-type fields; `Shape`/`Style`/`Metadata` are gone. `applyVisualStyle()` becomes the slim `applyDisplayStyle()` (terminal `Symbol`/`Color` only — those _are_ domain display concerns). Diagram export projects `GraphNode{Shape: Status.NodeShape(), Style: Status.GraphStyle()}` at `subscriberView.Nodes()`, so the framework coupling lives at the boundary instead of inside the domain.

### 5. Verification gates passed

| Gate                                          | Result              |
| --------------------------------------------- | ------------------- |
| `nix run .#build` (all 20 modules)            | ✅                  |
| `go test ./...` on nom                        | ✅ 141 tests pass   |
| `go test -race ./...` on nom                  | ✅ clean            |
| `go test ./...` on tui, integration, examples | ✅ all pass         |
| `go vet ./...` on nom                         | ✅                  |
| BuildFlow pre-commit hook                     | ✅ 17/17 steps pass |

---

## b) PARTIALLY DONE

### Documentation not yet updated for the 3 refactors

The following docs still describe the _old_ APIs and should be updated before the next release:

| Doc                                              | What's stale                                                                                                                |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| `AGENTS.md`                                      | "nom phase detection is ID-prefix-based" — now `ActivityKind`; "shared-pointer model" — Activity no longer embeds GraphNode |
| `FEATURES.md`                                    | No mention of `ActivityKind`, `NewPhase`, sealed events                                                                     |
| `CHANGELOG.md` `[Unreleased]`                    | Empty — the 3 refactors aren't logged                                                                                       |
| `docs/adr/007-nom-composition-via-root-types.md` | May reference the GraphNode embed                                                                                           |

### `CurrentElapsed` is still per-tick mutated

`UpdateRunningActivityElapsed()` walks all running activities every tick (100ms) to stamp `CurrentElapsed = now - StartTime`. This is a mutation hotspot that could be derived on-demand from `StartTime` at render time, eliminating the write path entirely. Not done — deferred to avoid scope creep in a type-model sprint.

---

## c) NOT STARTED

### Incremental subtree aggregates (upstream NOM's real algorithmic advantage)

`GetActivityCounts()` still iterates **all** activities every tick (`activity_management.go:60`). Upstream NOM maintains a `DependencySummary` monoid propagated up parents only on state transitions — O(changed) per update vs our O(n) per frame. This is the single scalability ceiling. Flagged in the comparative analysis but explicitly out of scope for a type-model sprint.

### `tui/` over-exposed API (~15 zero-caller exports)

Documented in the prior session's status report (`2026-06-20_00-49`). Message types (`ProgressUpdateMsg`, `ErrorMsg`, `StepUpdateMsg`, `StateTransitionMsg`, `CancelMsg`), update types, and state-check methods (`CanAcceptTicks`, `CanAcceptUpdates`, `CanTransitionTo`) have zero callers outside `tui/`. Should be unexported or moved to an internal subpackage before v1.0.

### `TableData` dual API (exported fields + getters)

ADR 006 issue #15 (blocked on owner decision): `TableData` exposes both `Headers` field and `GetHeaders()` method. Needs a v1 decision.

### v1.0.0 tag

API declared frozen (ADR 006) but still at v0.16.0. The 3 refactors this session are breaking changes to experimental APIs (explicitly allowed), so v1.0.0 is not blocked by them.

---

## d) TOTALLY FUCKED UP

**Nothing.** The codebase is in excellent shape:

- ✅ Zero data races (`-race` clean on nom + tui)
- ✅ All 20 modules build
- ✅ All test suites pass
- ✅ BuildFlow pre-commit hook passes 17/17 steps
- ✅ Zero deprecated symbols remain (`ParseActivityID`, `ColorRunning`, `applyVisualStyle`, all 10 local event structs — all gone)
- ✅ The 3 refactors compound: typed events + typed kind + decoupled domain = a cleaner foundation than upstream NOM's Haskell model in several respects (compile-time exhaustiveness, no string dispatch)

**One process mistake worth noting:** My first response in this session analyzed a _stale dirty working tree_ and claimed 6 issues, 4 of which were already fixed in committed code. I failed to `git log` and read `docs/status/` before recommending changes. Corrected after the user pushed back, but it wasted a round-trip.

---

## e) WHAT WE SHOULD IMPROVE

### Type-model opportunities still open

1. **`CurrentElapsed` → derived field.** Eliminate the per-tick mutation; compute from `StartTime` at render time. Smaller surface, fewer race surfaces, less code.
2. **Incremental `ActivityCounts`.** Port upstream's `DependencySummary` monoid. Biggest perf win at scale.
3. **`ActivityStatus.Interest()` → `SortOrder` enum.** The `int`-based priority is a magic number; a named enum would be self-documenting.
4. **`TimingCache` keyed on `string`.** Could use a branded type (`ActivityName` already exists) for compile-time safety.

### Process opportunities

5. **Read `docs/status/` and `git log` before analyzing.** This session's first-response failure cost a round-trip. Should be reflexive.
6. **Document the sealed-event migration in an ADR.** The decision to break the open `Event` interface is significant and reversible only with effort; an ADR would anchor the reasoning.
7. **Update docs in the same commit as the refactor.** AGENTS.md and FEATURES.md are now stale — this creates knowledge drift for the next session.

### Architecture opportunities

8. **`nom/` is 2,186 LOC across 35 files.** The prior status report flagged it as the coarsest module. Splitting into `internal/` sub-packages (e.g. `event`, `render`, `tree`, `cache`) would improve locality without expanding the public API.
9. **The `tui/` package tightly couples to `nom/` internals** (e.g., `WithSubscriberRLock`). A narrower interface contract would make the boundary explicit.

---

## f) Top 25 things to get done next

**Pareto-sorted: impact / work ratio, highest first.**

### Tier 1 — High impact, low work (do first)

| #   | Task                                                                                               | Impact | Work | Why                                                                 |
| --- | -------------------------------------------------------------------------------------------------- | ------ | ---- | ------------------------------------------------------------------- |
| 1   | **Update `CHANGELOG.md` `[Unreleased]`** with the 3 refactors                                      | med    | 10m  | Knowledge hygiene; the Unreleased section is empty                  |
| 2   | **Update `AGENTS.md`** Patterns + Gotchas for ActivityKind/sealed events/decoupled Activity        | med    | 20m  | Prevents the next session from re-discovering what changed          |
| 3   | **Update `FEATURES.md`** with `ActivityKind`, `NewPhase`, sealed events                            | low    | 15m  | Public feature inventory accuracy                                   |
| 4   | **Write ADR 009: sealed event sum type**                                                           | med    | 30m  | Anchors the breaking-change reasoning for future contributors       |
| 5   | **Derive `CurrentElapsed` from `StartTime`** at render time; delete `UpdateRunningActivityElapsed` | high   | 45m  | Eliminates a per-tick write path + a race surface + a public method |

### Tier 2 — High impact, medium work

| #   | Task                                                                       | Impact | Work | Why                                                                    |
| --- | -------------------------------------------------------------------------- | ------ | ---- | ---------------------------------------------------------------------- |
| 6   | **Incremental `ActivityCounts`** via subtree-aggregate monoid              | high   | 2h   | The single scalability ceiling; upstream NOM's real algorithmic edge   |
| 7   | **Unexport ~15 zero-caller `tui/` symbols**                                | med    | 1h   | Pre-v1 API hygiene; flagged in prior status report                     |
| 8   | **`ActivityStatus.Interest()` → named `SortOrder` enum**                   | low    | 30m  | Removes magic numbers from the sort path                               |
| 9   | **Branded type for `TimingCache` key** (currently `string`)                | low    | 30m  | Compile-time safety against mixing ActivityName with arbitrary strings |
| 10  | **`nom/` internal sub-package split** (`event`, `render`, `tree`, `cache`) | med    | 2h   | The coarsest module (35 files); improves navigability                  |

### Tier 3 — Medium impact, medium work

| #   | Task                                                                    | Impact | Work           | Why                                                                                                   |
| --- | ----------------------------------------------------------------------- | ------ | -------------- | ----------------------------------------------------------------------------------------------------- |
| 11  | **Cut v1.0.0 tag** (or decide to defer)                                 | high   | 30m            | API is frozen per ADR 006; the 3 refactors are the last planned breaking changes to experimental APIs |
| 12  | **Add `.github/workflows/ci.yml`** running lint+test+race on every push | high   | 1h             | No CI exists; relies entirely on pre-commit hooks and manual runs                                     |
| 13  | **`TableData` dual API decision** (ADR 006 #15)                         | med    | owner decision | Blocks v1.0.0 API freeze                                                                              |
| 14  | **`nom.DependencyTree` snapshot-render path documentation**             | low    | 30m            | The snapshot architecture is non-obvious; a diagram or ADR would help                                 |
| 15  | **Property-based test for sealed event exhaustiveness**                 | low    | 30m            | Guard against a future event type being added without a handler case                                  |

### Tier 4 — Lower impact or higher work

| #   | Task                                                                                      | Impact | Work  | Why                                                                                                                |
| --- | ----------------------------------------------------------------------------------------- | ------ | ----- | ------------------------------------------------------------------------------------------------------------------ |
| 16  | **Port upstream's greedy height-bounded node selection** (`derivationsToShow`)            | med    | 3h    | Current `elideCompletedUnderPressure` is simpler but less optimal than upstream's "fill 1/3 of terminal" heuristic |
| 17  | **Time-weighted timing cache** (decay old samples)                                        | low    | 1h    | With 10 samples it's low-stakes, but stale data dominates if the cap is raised                                     |
| 18  | **`GetMedian` caching** (compute on Record, invalidate on load)                           | low    | 30m   | Currently allocates+sorts every call                                                                               |
| 19  | **Narrow `tui/`→`nom/` coupling** (`WithSubscriberRLock` → interface)                     | low    | 1h    | Makes the module boundary explicit                                                                                 |
| 20  | **`renderNotify` test hook** moved off the production `InlineRenderer` struct             | low    | 20m   | Test-only field on a prod type; should be an option or separate test wrapper                                       |
| 21  | **CLI demo binary** (`cmd/go-output-demo`) showcasing all 16 formats + NOM                | low    | 2h    | Makes the library tangible; flagged in prior status reports                                                        |
| 22  | **`bdd/` package has zero nom event references** — verify BDD coverage of new event types | low    | 1h    | The sealed events may have created a coverage gap                                                                  |
| 23  | **`integration/` test for full workflow → DOT diagram export** with new ActivityKind      | low    | 45m   | Cross-module verification of the decoupled projection                                                              |
| 24  | **Audit `examples/` for stale patterns** post-refactor                                    | low    | 30m   | The examples were rewritten but may have lost illustrative value                                                   |
| 25  | **Investigate `nh`-style pipe mode** (`cmd/nom` binary reading `internal-json`)           | low    | 1 day | Scope expansion — only if external integration is ever wanted                                                      |

---

## g) Top #1 question I cannot figure out myself

**Should `CurrentElapsed` be deleted entirely (derived on-demand from `StartTime`), or does the cached value serve a purpose I'm not seeing?**

`UpdateRunningActivityElapsed()` walks all running activities every 100ms, mutating `CurrentElapsed = now - StartTime` under the subscriber's write lock. This is:

1. **A per-tick write path** — the only periodic mutation in an otherwise event-driven system.
2. **A race surface** — though currently closed by the subscriber lock, it's unnecessary contention with the render loop.
3. **Redundant** — `StartTime` is immutable after `SetRunning()`, and render time is the only consumer. `time.Since(start)` at render gives the same value without any mutation.

But the field has existed since the original NOM port, and `FormatActivityNodeTiming` reads it from `ActivitySnapshot`. Before I delete it, I want to confirm there's no consumer (e.g., a TUI path, a future "pause/resume" feature, a benchmark) that needs the _cached_ value rather than a fresh computation. The code doesn't reveal one, but the design intent might.

**My recommendation if yes:** Make `CurrentElapsed` a render-time computation in `SnapshotActivities()`, delete `UpdateRunningActivityElapsed()`, and remove the field from `Activity` and `ActivitySnapshot`. This simplifies the model and eliminates a write path — a net win. But this is your call.
