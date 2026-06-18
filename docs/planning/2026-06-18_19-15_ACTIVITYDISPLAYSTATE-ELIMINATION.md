# ActivityDisplayState Elimination — Comprehensive Plan

**Created:** 2026-06-18 19:15
**Goal:** Eliminate the `ActivityDisplayState` / `DisplayState` / `SyncActivityTimingToTree` split-brain by making `Activity` the single source of truth, shared via pointer between the subscriber and the dependency tree.

**Design:** `ActivityNode` embeds `*Activity` (pointer, not value). Go auto-promotes fields, so `node.Status`, `node.Symbol`, etc. continue to work with zero render-code changes. The subscriber's `activities` map stores `*Activity`. Both maps point to the same instances. Mutations via `SetRunning()` etc. are instantly visible to both. `SyncActivityTimingToTree` becomes dead code.

---

## Complete Touchpoint Map (28 files)

| File                            | Change                                                                                                   |
| ------------------------------- | -------------------------------------------------------------------------------------------------------- |
| `nom/activity.go`               | Add `CurrentElapsed` to `SetCompleted`/`SetFailed`; add `SetPaused`, `setOperationType`, `addDependency` |
| `nom/activity_display.go`       | **DELETE**                                                                                               |
| `nom/activity_accessors.go`     | `GetActivities`/`GetActivity` return `*Activity`                                                         |
| `nom/activity_management.go`    | Map type; delete sync methods; update counts/elapsed                                                     |
| `nom/tree.go`                   | `ActivityNode` embeds `*Activity`; update `newActivityNode`                                              |
| `nom/tree_modification.go`      | `AddActivity` takes `*Activity`                                                                          |
| `nom/status_updates.go`         | **DELETE** `UpdateActivityStatus`                                                                        |
| `nom/nom_subscriber.go`         | Map type; simplify `subscriberView.Nodes()`                                                              |
| `nom/subscriber_handlers.go`    | `getOrCreateActivity`→`*Activity`; remove `UpdateActivityStatus` calls                                   |
| `nom/format.go`                 | `FormatTimingInfo` takes `*Activity`                                                                     |
| `nom/configuration.go`          | `Reset` map init                                                                                         |
| `nom/inline_renderer.go`        | Remove `SyncActivityTimingToTree()` call                                                                 |
| `nom/activity_display_test.go`  | **Rewrite** as `activity` tests                                                                          |
| `nom/tree_test.go`              | Replace `UpdateActivityStatus`                                                                           |
| `nom/subscriber_test.go`        | Type refs; delete sync test                                                                              |
| `nom/format_test.go`            | `FormatTimingInfo` uses `*Activity`                                                                      |
| `nom/configuration_test.go`     | Map type                                                                                                 |
| `tui/model.go`                  | Remove `SyncActivityTimingToTree()` call                                                                 |
| `tui/model_test.go`             | `addTestActivity` uses `NewActivity`                                                                     |
| `tui/event_sequence_test.go`    | `SetActivityState` with `NewActivity`                                                                    |
| `tui/view_test.go`              | `UpdateActivityStatus` → direct Activity mutation                                                        |
| `integration/nom_tui_test.go`   | Remove sync; update `mustUpdateActivityStatus`                                                           |
| `examples/nom_progress/main.go` | Remove `SyncActivityTimingToTree()` call                                                                 |
| `CHANGELOG.md`                  | Breaking change entry                                                                                    |
| `AGENTS.md`                     | Update patterns                                                                                          |
| `docs/FORMAT_ARCHITECTURE.md`   | Single source of truth                                                                                   |
| `docs/adr/007-*.md`             | Mark migration done                                                                                      |
| `TODO_LIST.md`                  | Record resolution                                                                                        |

---

## Task Breakdown (≤12 min each, sorted by dependency → impact → effort)

### Wave 1: Additive changes to Activity (non-breaking)

| ID   | Task                                                                                        | Impact | Effort |
| ---- | ------------------------------------------------------------------------------------------- | ------ | ------ |
| W1-1 | Add `CurrentElapsed` computation to `Activity.SetCompleted()` + `SetFailed()` (correctness) | 🔥🔥🔥 | 5m     |
| W1-2 | Add `SetPaused()`, `setOperationType()`, `addDependency()` methods to `Activity`            | 🔥     | 8m     |
| W1-3 | Verify Activity builds + existing tests green                                               | 🔥     | 3m     |

### Wave 2: Core type changes (breaks build, fixed in Wave 3)

| ID   | Task                                                                     | Impact | Effort |
| ---- | ------------------------------------------------------------------------ | ------ | ------ |
| W2-1 | `tree.go`: `ActivityNode` embeds `*Activity` instead of `Activity` value | 🔥🔥🔥 | 5m     |
| W2-2 | `tree.go`: update `newActivityNode` to store `*Activity`                 | 🔥🔥🔥 | 3m     |
| W2-3 | `nom_subscriber.go`: `activities` map → `map[ActivityID]*Activity`       | 🔥🔥🔥 | 5m     |
| W2-4 | `tree_modification.go`: `AddActivity(id, *Activity, deps)` signature     | 🔥🔥🔥 | 8m     |
| W2-5 | `subscriber_handlers.go`: `getOrCreateActivity` → `*Activity`            | 🔥🔥🔥 | 5m     |
| W2-6 | `activity_management.go`: `SetActivityState` → `*Activity`               | 🔥🔥   | 3m     |
| W2-7 | `activity_accessors.go`: `GetActivities`/`GetActivity` → `*Activity`     | 🔥🔥   | 5m     |
| W2-8 | `format.go`: `FormatTimingInfo(state *Activity)`                         | 🔥🔥   | 3m     |

### Wave 3: Fix all callers (restores build)

| ID    | Task                                                                                  | Impact | Effort |
| ----- | ------------------------------------------------------------------------------------- | ------ | ------ |
| W3-1  | Handlers: update `handleActivityStarted` (pass Activity, remove UpdateActivityStatus) | 🔥🔥🔥 | 8m     |
| W3-2  | Handlers: update `handleActivityRegistered` (pass Activity)                           | 🔥🔥   | 5m     |
| W3-3  | Handlers: update `handleActivityCompleted`/`Failed` (remove UpdateActivityStatus)     | 🔥🔥🔥 | 8m     |
| W3-4  | Handlers: update `handleWorkflowStarted` map init                                     | 🔥     | 3m     |
| W3-5  | `nom_subscriber.go`: simplify `subscriberView.Nodes()` (use `a.GraphNode`)            | 🔥🔥   | 5m     |
| W3-6  | `activity_management.go`: update `GetActivityCounts` + `UpdateRunningActivityElapsed` | 🔥🔥   | 5m     |
| W3-7  | `configuration.go`: `Reset` map init                                                  | 🔥     | 2m     |
| W3-8  | `inline_renderer.go`: remove `SyncActivityTimingToTree()` call                        | 🔥🔥🔥 | 2m     |
| W3-9  | `tui/model.go`: remove `SyncActivityTimingToTree()` call                              | 🔥🔥🔥 | 2m     |
| W3-10 | `integration/nom_tui_test.go`: remove sync calls                                      | 🔥🔥   | 3m     |
| W3-11 | `examples/nom_progress/main.go`: remove sync call                                     | 🔥     | 2m     |

### Wave 4: Delete dead code

| ID   | Task                                                                                   | Impact | Effort |
| ---- | -------------------------------------------------------------------------------------- | ------ | ------ |
| W4-1 | Delete `activity_display.go` entirely                                                  | 🔥🔥🔥 | 2m     |
| W4-2 | Delete `UpdateActivityStatus` from `status_updates.go`                                 | 🔥🔥   | 2m     |
| W4-3 | Delete `SyncActivityTimingToTree` + `syncActivityToNode` from `activity_management.go` | 🔥🔥🔥 | 3m     |

### Wave 5: Fix tests

| ID   | Task                                                                                              | Impact | Effort |
| ---- | ------------------------------------------------------------------------------------------------- | ------ | ------ |
| W5-1 | Rewrite `activity_display_test.go` → test `Activity` (constructor, transitions, predicates, copy) | 🔥🔥   | 12m    |
| W5-2 | `tree_test.go`: replace `UpdateActivityStatus` with direct Activity mutations                     | 🔥🔥   | 10m    |
| W5-3 | `subscriber_test.go`: type refs, delete SyncActivityTimingToTree test                             | 🔥🔥   | 12m    |
| W5-4 | `format_test.go`: `FormatTimingInfo` tests use `*Activity`                                        | 🔥     | 8m     |
| W5-5 | `configuration_test.go`: map type update                                                          | 🔥     | 5m     |
| W5-6 | `tui/model_test.go`: `addTestActivity` uses `NewActivity`                                         | 🔥     | 5m     |
| W5-7 | `tui/event_sequence_test.go`: `SetActivityState` with `NewActivity`                               | 🔥     | 8m     |
| W5-8 | `tui/view_test.go`: `addRunningActivity` uses Activity methods                                    | 🔥     | 5m     |
| W5-9 | `integration/nom_tui_test.go`: rewrite `mustUpdateActivityStatus`, fix `GetActivity` callers      | 🔥🔥   | 12m    |

### Wave 6: Verify (every gate must pass)

| ID   | Task                                           | Impact | Effort |
| ---- | ---------------------------------------------- | ------ | ------ |
| W6-1 | Build nom module green                         | 🔥🔥🔥 | 3m     |
| W6-2 | Test nom module green                          | 🔥🔥🔥 | 5m     |
| W6-3 | Race test nom + tui green                      | 🔥🔥🔥 | 5m     |
| W6-4 | Build + test tui green                         | 🔥🔥🔥 | 5m     |
| W6-5 | Build + test integration green                 | 🔥🔥   | 3m     |
| W6-6 | Build + test examples green                    | 🔥     | 3m     |
| W6-7 | Lint all 17 modules clean                      | 🔥🔥🔥 | 5m     |
| W6-8 | Full workspace sweep (all nix apps equivalent) | 🔥🔥🔥 | 5m     |

### Wave 7: Docs

| ID   | Task                                                     | Impact | Effort |
| ---- | -------------------------------------------------------- | ------ | ------ |
| W7-1 | CHANGELOG.md — breaking change entry                     | 🔥🔥   | 5m     |
| W7-2 | AGENTS.md — update nom patterns (single source of truth) | 🔥     | 5m     |
| W7-3 | FORMAT_ARCHITECTURE.md — Activity is the store           | 🔥     | 5m     |
| W7-4 | ADR 007 — mark migration done                            | 🔥     | 3m     |
| W7-5 | TODO_LIST.md — record resolution                         | 🔥     | 3m     |

**Total: 40 tasks · ~3.5h estimated**
