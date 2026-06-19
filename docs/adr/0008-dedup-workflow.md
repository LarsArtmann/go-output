# ADR 008: Dedup Workflow Decision

**Date:** 2026-06-18
**Status:** Accepted

## Context

The project uses `art-dupl` for code duplication detection. The standard threshold is `t=24` (see ADR 005). When clones are reported, a judgment-based workflow determines whether to act on them.

## Decision

We adopt a 5-step checklist for evaluating every clone group reported by `art-dupl -t 24`:

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

At `t=50`, zero actionable clones exist. The `t=24` threshold is the sweet spot — it surfaces real duplication without drowning in noise.

## Consequences

- `art-dupl -t 24` is the project standard. Do not change the threshold without re-evaluating this ADR.
- Test code duplication is explicitly accepted as a project pattern.
- CI should gate on `art-dupl -t 30` to catch egregious regressions without failing on acceptable patterns.
- Zero actionable clones is NOT the goal — zero HARMFUL duplication is.
