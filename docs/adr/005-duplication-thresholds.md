# ADR 005: Code Duplication Thresholds and Acceptable Categories

**Status:** Accepted _(decision still holds; **counts below are outdated as of 2026-07-26** — see [Update](#update-2026-07-26) at the end)_
**Date:** 2026-05-28

## Context

Running `art-dupl --semantic -t 15` reported approximately 50 clone groups in the codebase at the time. _(As of 2026-07-26 the codebase is far cleaner: `art-dupl -t 4` = **0** groups, `-t 3` = **2**, `-t 2` = **16**, `-t 1` = **20** accepted. See the [Update](#update-2026-07-26) appendix.)_ We need clear guidelines on which clones to fix vs accept.

## Decision

### Threshold Policy

_**(Counts in this table are the original 2026-05-28 snapshot — see [Update](#update-2026-07-26) for current figures.)**_

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

## Update (2026-07-26)

**The decision and category policy above still hold — only the counts have drifted.**
After multiple dedup sweeps (v0.30–v0.32), the current `art-dupl` figures are:

| Threshold                    | Clone Groups | Notes                                                                                                |
| ---------------------------- | ------------ | --------------------------------------------------------------------------------------------------- |
| **t=4** (production gate)    | **0**        | Gate is clean — zero actionable clones.                                                             |
| **t=3**                      | **2**        | Both accepted minimum idioms: thread-safe time-read lock scope; `strings.Builder` opener.           |
| **t=2**                      | **16**       | All accepted: test idioms, module boundaries, examples, single-line patterns, minimum Go idioms.    |
| **t=1** (strict type-aware)  | **20**       | All accepted per the category policy above.                                                         |

The category definitions (Accept: B/C/D/E; Fix: 1–4) are unchanged and remain the
governing rules for all future dedup judgment calls. Current accepted-group
enumeration lives in `AGENTS.md` ("Current dedup state" bullet).
