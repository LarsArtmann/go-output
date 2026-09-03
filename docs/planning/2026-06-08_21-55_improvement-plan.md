# Improvement Plan: Post-Dedup Deep Dive

**Date:** 2026-06-08 21:55 CEST
**Reporter:** Crush
**Based on:** `fa6ac0d` (master)
**Status:** 3 of 5 Phase 1/2 items complete

---

## Research Findings

### 1. Error Pattern Inconsistency (Medium Impact, Low Effort)

Root enums use 4 different invalid error patterns:

- `InvalidFormatError{Value, Allowed}` — includes allowed values
- `InvalidShapeError{Value}` — just value
- `InvalidColorModeError{Value}` — just value
- `InvalidGraphShapeError{Value}` — just value

D2 enums use sentinel errors + `fmt.Errorf("%w: %q", ...)`:

- `var ErrInvalidD2Direction = errors.New("invalid D2 direction")`

**Opportunity:** Unify to one pattern. The sentinel error + fmt.Errorf wrapping is cleaner (works with `errors.Is()`). The struct approach requires new types per enum.

**Status:** Deferred — breaking change for consumers who type-assert on `*InvalidFormatError`. Move to v1 decision.

### 2. Missing Race Test for `RegisterFormatShapes` (Medium Impact, Low Effort)

`RegisterTableDataMarshaler` has `TestRegisterTableDataMarshaler_ConcurrentAccess`.
`RegisterFormatShapes` has no equivalent test. The pattern is identical (RWMutex + map).

**Why it matters:** Both registries use `sync.RWMutex` + map. Without a race test, a future refactor could accidentally remove the mutex or introduce unsafe map access. The test documents the thread-safety contract.

**Status:** ✅ COMPLETE — `TestRegisterFormatShapes_ConcurrentAccess` added to `format_shape_test.go`

### 3. `TableData.Validate()` is Minimal (Medium Impact, Low Effort)

Previously only checked footer column count. Could also validate:

- Row column counts match header count (but renderers already handle short rows gracefully — too restrictive to enforce)
- No nil rows (would panic in some renderers that don't check)
- Headers are non-empty (but empty headers are valid for some use cases)

**Why it matters:** Nil rows in `[][]string` are almost certainly bugs. They can cause panics or silent data loss in downstream renderers. Catching them at validation time prevents runtime failures.

**Status:** ✅ COMPLETE — Added `errNilRow` sentinel error; `Validate()` checks for nil rows and returns wrapped error with index.

### 4. `enum` Package is Already Excellent (No Action)

The generic `Parse[T]`, `Contains[T]`, `AllowedValues[T]` already eliminate ~80% of enum boilerplate. Each enum needs only: type, consts, values slice, Parse, IsValid, AllowedValues, String, error type. Code generation would save ~5 lines per enum but add build complexity. **Not worth it.**

### 5. `RenderOptions` Field Usage is Clear (No Action)

- `Title`: markdown, html
- `GraphID`: dot
- `Writer`: all
- `ColorMode`: table, tree, markdown

No architectural issue — each format uses what it needs.

### 6. No Missing Libraries (No Action)

Already using appropriate libraries for each domain:

- `html/template` for HTML (standard, auto-escaping)
- `encoding/json` for JSON (standard)
- `go-faster/yaml` for YAML (fast, maintained)
- `go-toml/v2` for TOML (standard)
- `charm.land/lipgloss/v2` for terminal tables (best-in-class)
- `golang.org/x/term` for terminal detection (standard)

### 7. BuildFlow `library-policy` (External Blocker)

Suggests `github.com/a-h/templ` and `github.com/larsartmann/go-error-family`. Both would be significant refactors. Defer to owner decision.

---

## Execution Results

### Phase 1: Low Effort, Medium Impact

| # | Task                                               | Effort | Impact | Status     | Commit    |
| - | -------------------------------------------------- | ------ | ------ | ---------- | --------- |
| 1 | Add race test for `RegisterFormatShapes`           | 10 min | Medium | ✅ Done    | `86eebac` |
| 2 | Improve `TableData.Validate()` — nil row detection | 15 min | Medium | ✅ Done    | `34570b1` |
| 3 | Unify root error patterns                          | 20 min | Medium | ⏸️ Deferred | —         |

### Phase 2: Medium Effort, Medium Impact

| # | Task                                              | Effort | Impact | Status  | Commit    |
| - | ------------------------------------------------- | ------ | ------ | ------- | --------- |
| 4 | Update AGENTS.md with `RegisterFormatShapes` docs | 10 min | Low    | ✅ Done | `6f6114f` |
| 5 | Verify all `replace` directives                   | 15 min | Low    | ✅ Done | —         |

### Phase 3: Deferred (Needs Owner Decision)

| # | Task                                        | Effort | Impact | Blocker        |
| - | ------------------------------------------- | ------ | ------ | -------------- |
| 6 | Configure BuildFlow `library-policy`        | 15 min | High   | Owner decision |
| 7 | Decide v1 API: exported fields vs getters   | 30 min | High   | Owner decision |
| 8 | Add `gomod2nix` for reproducible Nix builds | 2h     | High   | Nix expertise  |

---

## Type Model Improvement: Generic InvalidValueError

Instead of 4+ error structs, a generic type:

```go
type InvalidValueError[T comparable] struct {
    Name   string
    Value  T
    Values []T
}
```

This would replace `InvalidFormatError`, `InvalidShapeError`, `InvalidColorModeError`, `InvalidGraphShapeError` with one type. However, changing error types is a **breaking change** for consumers who type-assert errors. Defer to v1 decision.

For now, unify the _pattern_ within non-breaking constraints:

- Use sentinel errors + `fmt.Errorf("%w: %q", ...)` for all new enums
- For existing root enums, keep structs but make error messages consistent

---

## What I Considered But Did NOT Do

1. **Strict row width validation** — Renderers intentionally handle short rows (e.g., `ToMapSlice()` skips extra headers). Enforcing exact width would break legitimate use cases.

2. **Enum code generation** — The `enum` package already provides generic `Parse`, `Contains`, `AllowedValues`. Each enum needs only ~20 lines of boilerplate. Code generation would add build complexity for minimal gain.

3. **Error type unification** — Changing `InvalidFormatError` to a sentinel error would break consumers who do `errors.As(err, &invalidErr)`. Defer to v1.

4. **gomod2nix** — Requires evaluating whether it supports multi-module workspaces with 14 modules and extensive `replace` directives. Needs Nix expertise.

---

_Plan generated by Crush on 2026-06-08 21:55 CEST_
_Updated by Crush on 2026-06-08 22:00 CEST_
