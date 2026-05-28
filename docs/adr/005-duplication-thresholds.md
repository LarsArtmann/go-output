# ADR 005: Code Duplication Thresholds and Acceptable Categories

**Status:** Accepted  
**Date:** 2026-05-28

## Context

Running `art-dupl --semantic -t 15` reports ~50 clone groups in the codebase. We need clear guidelines on which clones to fix vs accept.

## Decision

### Threshold Policy

| Threshold                    | Clone Groups | Action                                  |
| ---------------------------- | ------------ | --------------------------------------- |
| **t=50** (industry standard) | 2            | Fix all                                 |
| **t=30** (meaningful)        | ~11          | Fix production code, evaluate test code |
| **t=15** (aggressive)        | ~50          | Categorize and accept most              |

### Acceptable Clone Categories

**Accept — do NOT attempt to deduplicate:**

1. **Go test idioms (Category B)**: `strings.Contains` checks, `t.Errorf` patterns, `t.Parallel()`, table-driven test struct declarations. These are language-standard testing patterns.

2. **Module boundary (Category C)**: Interface re-declarations, type aliases, helper re-exports across Go modules. Go's module system forces this duplication. The `testhelpers` package is zero-dep by design and cannot import `output`.

3. **Example/documentation code (Category D)**: Full API usage in `example_test.go`, `examples/`, and GoDoc examples. These must be self-contained to serve as documentation.

4. **Single-line patterns (Category E)**: Function signatures implementing the same interface (`TableDataMarshaler`), `init()` registration closures (1-line unique bindings), `var _ Interface = (*Type)(nil)` assertions.

**Fix — deduplicate when found:**

1. **Identical function bodies**: Same logic with only format-name differences (e.g., `renderYAMLTableData` vs `renderTOMLTableData` → extracted `renderViaRenderer`).
2. **Identical test functions**: Same test body testing the same thing (e.g., duplicate `TestMarshalTSVUnsupportedType`).
3. **Duplicated validation predicates**: Character validity checks duplicated from production code into test code (use idempotency checks instead).
4. **Repeated loop bodies**: Identical inner loops in production code (e.g., AsciiDoc row/footer cell writing → `writeAsciiDocCells`).

## Consequences

- At t=50: Zero actionable clones remain
- At t=30: Only cross-module test patterns remain (acceptable)
- Dedup work focuses on production code and truly identical test functions
- `testhelpers` remains zero-dep, preserving the lightweight dependency profile
