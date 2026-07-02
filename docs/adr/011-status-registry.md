# ADR 011: Status Registry — Extensible Activity Statuses

**Date:** 2026-07-02
**Status:** ACCEPTED & IMPLEMENTED
**Deciders:** Lars Artmann

## Context

Activity statuses were hardcoded as a fixed enum (`ActivityStatusPending`, `ActivityStatusRunning`, `ActivityStatusCompleted`, `ActivityStatusFailed`) with every status-to-symbol and status-to-color mapping implemented via `switch` statements scattered across the codebase. This created two problems:

1. **No extensibility**: Consumers (e.g. BuildFlow) that needed custom statuses (e.g. `ActivityStatusCached`, `ActivityStatusSkipped`) had no way to register them — they'd have to fork the enum or use stringly-typed workarounds.

2. **Scattered switch statements**: Every time a new status was added, the `GetSymbol()`, `GetColor()`, `String()`, and `IsValid()` switch statements all needed updating, in multiple files. Missing one case was a silent bug.

## Decision

Introduce a thread-safe `StatusDef` registry as the single source of truth for status metadata.

### Design:

- **`StatusDef` struct**: Contains `Status ActivityStatus`, `Name string`, `Symbol Symbol`, `Color color.Color`, `Priority int`.
- **Global registry**: Package-level `statusRegistry` map protected by `sync.RWMutex`. Pre-populated with the 4 built-in statuses at `init()` time.
- **`RegisterStatus(def StatusDef)`**: Adds or updates a status. Thread-safe.
- **`LookupStatus(status) (StatusDef, bool)`**: O(1) map lookup replaces all `switch` statements.
- **`AllActivityStatuses() []StatusDef`**: Returns all registered statuses in ascending ID order.
- **`GetSymbol()` / `GetColor()`**: Now delegate to `LookupStatus()` internally — no more scattered switches.
- **Snapshot-time resolution**: The subscriber resolves status → symbol/color via the registry at snapshot time (`SnapshotActivities()`), so custom statuses flow through to renderers automatically.

### Registration ergonomics:

```go
nom.RegisterStatus(nom.StatusDef{
    Status:   nom.ActivityStatus(99),
    Name:     "cached",
    Symbol:   nom.Symbol("◇"),
    Priority: 3,
})
```

## Consequences

**Positive:**

- Adding a custom status is a single `RegisterStatus()` call — no code changes needed in the library.
- All `switch` statements on status are replaced by O(1) map lookups.
- The registry is thread-safe, so consumers can register statuses at any time (including at `init()`).

**Negative:**

- Global mutable state (the registry) is inherently less testable than pure functions. Mitigated by thread-safety and the `ResetStatusRegistry()` test helper.
- Custom statuses registered by one package affect all consumers of the same process. This is acceptable for CLI tools (single binary) but would be problematic in a library consumed by multiple packages.

**Neutral:**

- The 4 built-in statuses are always pre-registered, so the default behavior is unchanged. Custom statuses are purely additive.
