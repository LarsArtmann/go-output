# ADR 008: Dedup Workflow Decision

**Date:** 2026-06-18
**Status:** Accepted

## Context

The project uses `art-dupl` for code duplication detection. The standard production threshold is `t=4`; `--type-aware -t 1` is the strict audit mode used to review minimum idioms. When clones are reported, a judgment-based workflow determines whether to act on them.

## Decision

We adopt a 5-step checklist for evaluating every clone group reported by `art-dupl -t 4`, with periodic `art-dupl --sort total-tokens -t 1 --type-aware` audits:

1. **Generated/single-line/interface-compliance?** → Accept. No action needed.
2. **Structural or semantic?** If it's an idiomatic Go test pattern (`strings.Contains`, `t.Errorf`, table-driven) → Accept.
3. **Would abstraction help readability?** If no → Accept.
4. **Drift likely?** If no → Accept.
5. **Only fix production clones** where identical logic is duplicated. Never abstract test code for the sake of dedup metrics.

## Rationale

At threshold `t=15`, approximately 100 clone groups are reported. Analysis showed that ~95% are either:

- Table-driven test idioms (acceptable duplication)
- Module-boundary re-declarations (structural, not semantic)
- Interface conformance (each implementation must satisfy the interface)

Chasing these creates harmful abstractions that make the code harder to read, not easier.

At `t=4`, zero clone groups remain. The strict type-aware `t=1` audit reported 24 groups after the July 2026 sweep; all are minimum Go idioms, module-boundary contracts, self-contained examples, or short error-handling patterns where abstraction would reduce clarity.

## Consequences

- `art-dupl -t 4` is the project standard; use `art-dupl --sort total-tokens -t 1 --type-aware` for strict audits.
- Keep `examples/` in strict scans as visible documentation-code context, but classify its self-contained clones as acceptable under ADR 005 rather than refactoring or excluding them.
- Test code duplication is explicitly accepted as a project pattern.
- CI should gate on `art-dupl -t 4` to catch production regressions.
- Zero report lines is NOT the goal; zero HARMFUL duplication is.
