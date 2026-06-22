# ADR 009: Pattern B Versioning — Committed Replace Directives

**Date:** 2026-06-22
**Status:** ACCEPTED & IMPLEMENTED
**Deciders:** Lars Artmann

## Context

go-output used a "mono-version lockstep" policy (ADR 001): all 20 modules tagged at the same `vX.Y.Z`, with each `go.mod` requiring siblings at real published versions. In practice this created two problems:

1. **Version rot hidden by `replace`**: Every `go.mod` had `replace => ../sibling` directives for local development. These `replace` directives masked stale version pins — local builds were always green, but external consumers resolved a Frankenstein mix of versions that had never been built or tested together. By v0.17.2, sibling pins ranged from `v0.13.0` (testhelpers, frozen 4 releases ago) to `v0.17.1` (one minor behind), and `nom/v0.17.2`'s go.mod declared a dependency on root `v0.17.1` — a version behind its own release.

2. **Maintenance burden with zero benefit**: Every release required bumping 100+ sibling version pins across 20 `go.mod` files. This was tedious, error-prone, and delivered no value — no consumer was independently `go get`-ing individual sub-modules.

Research into established Go projects revealed two legitimate patterns:

- **Pattern A** (google-cloud-go): Real version pins, no `replace`. Each sub-module is independently consumable. High maintenance cost but maximum flexibility.
- **Pattern B** (block/blip, deliveryhero/asya): `v0.0.0-00010101000000-000000000000` + committed `replace`. Sub-modules are build-boundary only. Zero pin maintenance but sub-modules are not independently consumable.

## Decision

Adopt **Pattern B** (committed `replace` + `v0.0.0-00010101000000-000000000000`) for all inter-module dependencies.

### Scope

- **Root** (`github.com/larsartmann/go-output`) remains independently `go get`-able. Its `go.mod` requires only external deps (`go-branded-id`, `x/term`) and `testhelpers`.
- **`testhelpers/`** keeps real published versions — it's zero-dep, used by every module's tests, and genuinely independently useful.
- **All other sub-modules** use `v0.0.0-00010101000000-000000000000` + `replace => ../path`. They are consumed via clone + `go.work`.

### Prerequisite: enum and envdetect merged into root

Root previously depended on `enum/` and `envdetect/` sub-modules for production code (shape.go, format.go, graph.go, color.go). Under Pattern B, these would become `v0.0.0-...` in root's go.mod, breaking external consumers (since `replace` doesn't propagate).

Solution: merge `enum/` and `envdetect/` packages into root's `package output`. Their functions (`ParseEnum`, `ContainsEnum`, `IsCI`, `IsNoColor`, `CIEnvVars`, etc.) are now in root. The `enum/` and `envdetect/` directories were deleted.

### Prerequisite: cross-module tests moved to integration/

Root's `userjourney_test.go` and `render_tabledata_test.go` imported sub-modules (delimited, markdown, tree). These would add `v0.0.0-...` deps to root's go.mod. Both files were moved to `integration/` where all modules are available.

## Consequences

**Positive:**

- Zero sibling version-pin maintenance — no more bumping 100+ lines per release
- No more version rot — the entire concept of "sibling version drift" is eliminated
- Root's `go.mod` is clean and honest: only external deps + testhelpers
- `go.work` + `go.work.example` handle local development (already existed)
- The `v0.0.0-00010101000000-000000000000` sentinel is Go's canonical zero version for locally-replaced deps — it's what `go mod tidy` generates

**Negative:**

- `go get github.com/larsartmann/go-output/table` no longer works (must clone repo + `go.work`)
- Two fewer independently consumable modules (`enum/` and `envdetect/` merged into root)
- Sub-module consumers must use `replace` directives in their own `go.mod`

## Rejected Alternatives

### Pattern A (real pins, no replace)

Requires maintaining 100+ version pins across 20 `go.mod` files. No consumer was independently using sub-modules. The maintenance cost was pure overhead.

### Hybrid (D0b: leaf modules tagged, heavy modules on Pattern B)

Would still require maintaining version pins for leaf modules. The split policy adds cognitive overhead ("which modules are tagged?") for minimal benefit.

### D0c (full Pattern B including root)

Would make root itself not independently consumable. This contradicts the library's primary value proposition: `go get github.com/larsartmann/go-output` should work.
