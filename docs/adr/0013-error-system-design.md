# ADR 013: Error System Design

**Date:** 2026-07-30
**Status:** Accepted
**Frozen API impact:** Additive only (new exported symbols `ErrColumnMismatch`, `ErrNilRow`)

## Context

An `erraudit` run on the root package surfaced 27 violations. After analysis, ~90% were false positives — the tool recommends anti-patterns for a Go library. However, one real bug was discovered: `AddRowChecked`'s doc comment promised `ErrColumnMismatch` (exported) while the variable was `errColumnMismatch` (unexported). This is a lying doc comment — the most insidious error-system defect because consumers trust the documented contract.

The project needed a deliberate error system design, not reactive patching, to prevent future degradation.

## Decision

Adopt a **three-tier error model** following Go convention:

### Tier 1: Sentinel Errors (`errors.Is`)

```go
var ErrColumnMismatch = errors.New("column count does not match headers")
var ErrNilRow         = errors.New("nil row in data")
```

- For conditions consumers match programmatically via `errors.Is(err, ErrFoo)`.
- Always wrapped with `fmt.Errorf("%w: ...", ErrFoo, ...)` to preserve the chain.
- Root owns the contract; sub-modules wrap, never redefine.

### Tier 2: Typed Error Structs (`errors.AsType`)

```go
type UnsupportedFormatError struct{ Format Format }
type ParseError struct{ Value, Allowed string }
```

- For structured error data consumers extract via `errors.AsType[*T](err)` (Go 1.26+ generic).
- NOT `errors.As` — that is legacy in Go 1.26+.
- Each has a proper `Error() string` method with enough context to be self-describing.

### Tier 3: Context Wrapping (`fmt.Errorf`)

```go
return fmt.Errorf("render %s: %w", formatName, err)
```

- Adds operation/format context without losing the error chain.
- Never dumps render output (`out=%v`) — it's the failed result, always garbage on error path.
- Never dumps raw structs (`data=%v`, `opts=%v`) — leaks internal representation.

## What We Explicitly Rejected

| Rejected | Why |
|---|---|
| Concrete return types (`func() *MarshalError`) | Couples consumers to implementation; Go returns `error` |
| Error hierarchy (`RenderError` interface) | Over-engineering; flat model is sufficient and idiomatic |
| `out=%v` in error messages | Failed render result is garbage; leaks noise, breaks grep |
| Merging d2/graph error types | Different domains (D2 shapes vs graph shapes) |
| `errors.As(err, &target)` | Legacy in Go 1.26+; use `errors.AsType[*T](err)` |

## Consequences

- **Positive:** Consumers can now `errors.Is(err, output.ErrColumnMismatch)` — the documented contract is real.
- **Positive:** Contract tests (`TestSentinelErrors_MatchThroughWrapping`) lock the `%w` wrapping chain.
- **Positive:** `erraudit` findings are classified — `context_loss` and `generic_return` are documented false positives.
- **Neutral:** Two new exported symbols (`ErrColumnMismatch`, `ErrNilRow`) — additive, not breaking.
