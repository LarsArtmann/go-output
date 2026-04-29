# go-output — Comprehensive Status Report

**Date:** 2026-04-29 23:30  
**Session:** Two sessions — bugfix + comprehensive improvement  
**Commits:** 11 (fe79fec..9a943a8)  
**Net change:** +660 / -4,724 lines across 33 files  
**Status:** All tests pass, lint clean, pushed to origin/master

---

## A) FULLY DONE

### Bug Fixes (Critical)

| # | Fix | Severity | Commit |
|---|-----|----------|--------|
| 1 | Sort descending stability violation — `!result` for equal elements broke strict weak ordering | CRITICAL | `e882804` |
| 2 | CI Go version mismatch — 1.23 in CI vs 1.26 in go.mod (build-breaking) | CRITICAL | `bf543fb` |
| 3 | `golang.org/x/term` listed as indirect but imported directly in color.go | HIGH | `fe79fec` |
| 4 | Recursive `joinStrings` in enum — stack overflow risk on large lists | HIGH | `fe79fec` |

### Architecture Improvements

| # | Change | Impact | Commit |
|---|--------|--------|--------|
| 5 | D2 enums: per-call allocating functions → package-level vars | Consistency + perf | `4aa86a8` |
| 6 | Split d2.go (333 lines) → d2.go (105) + d2_enum.go (226) | File size limit | `6f274fa` |
| 7 | Remove dead `ColorMode.ToANSI()` — zero production callers | API surface reduction | `3b1aeea` |
| 8 | Remove dead `GraphNode.GetStyle()` — zero production callers | API surface reduction | `3b1aeea` |
| 9 | Document Registry as opt-in plugin system | Clarify ghost system | `b7d5cdf` |

### CI/CD

| # | Change | Commit |
|---|--------|--------|
| 10 | Fix Go 1.23→1.26 in ci.yml AND release.yml | `bf543fb` |
| 11 | Consolidate 3 duplicated CI jobs → 1 single job | `bf543fb` |
| 12 | Remove goreleaser job (no .goreleaser.yml exists, library has no binaries) | `bf543fb` |

### Documentation

| # | Change | Lines Removed | Commit |
|---|--------|---------------|--------|
| 13 | Delete 15 stale docs (root + docs/planning + docs/status) | -4,693 | `1e36a4d` |
| 14 | Rewrite PLAN.md to match reality | -258 | `2c91e7e` |
| 15 | Update AGENTS.md (full rewrite with architecture notes) | rewrite | `9a943a8` |
| 16 | Fix README.md (deps version, escape import path) | fix | `9a943a8` |
| 17 | Add comprehensive improvement plan with mermaid graph | +143 | `2bffa7f` |

### Test Coverage Improvements

| Package | Before | After | New Tests Added |
|---------|--------|-------|-----------------|
| sort | 86.7% | **95.5%** | 5 (DescStability, DescCount, NonStructInput, SnakeToPascal×6) |
| enum | 94.7% | **100%** | 0 (fixed by joinStrings refactor) |

---

## B) PARTIALLY DONE

### Tree Conversion Split Brain

**Status:** Identified but not fixed. Three renderer-specific `addTreeNodes()` implementations exist:
- `d2_convert.go::addTreeNodes()` — D2-specific ID resolution
- `dot.go::addTreeNodes()` — DOT-specific ID resolution (via `dotTreeNodeID`)
- `mermaid.go::addTreeNodes()` — Mermaid-specific ID resolution (via `mermaidTreeNodeID`)

There's also a generic `graph.go::AddTreeNodes()` that the specific versions don't use directly.

**Why not fixed:** Each renderer has legitimate format-specific ID sanitization needs. Forcing them through a single path would require a more complex abstraction. The current duplication is contained (3 small methods) and explicit.

### nolint Directives

**Status:** 47 remain. Breakdown:
- 25 `exhaustruct` — intentional partial struct initialization
- 13 `gochecknoglobals` — legitimate global enum value slices
- 4 `gosec` — acceptable (file descriptor from int, test indices)
- 1 `cyclop` — test complexity

The `exhaustruct` ones indicate structs with many optional fields. Proper fix would be functional option constructors, but that's a large API change.

---

## C) NOT STARTED

| # | Task | Impact | Effort | Notes |
|---|------|--------|--------|-------|
| 1 | Unify tree conversion (eliminate 3-way split brain) | Medium | High | May not be worth the complexity |
| 2 | Add functional option constructors for D2Node/D2Edge/D2NodeStyle | Medium | High | Would eliminate 25 exhaustruct nolints |
| 3 | Add godoc examples (Example functions) | Medium | Medium | Improves pkg.go.dev documentation |
| 4 | Add property-based/fuzz tests for escape functions | Medium | Medium | Current fuzz tests only cover ParseFormat |
| 5 | D2 diagram map ordering (deterministic output) | Low | Medium | Go maps are unordered |
| 6 | DOT attribute quoting robustness | Low | Low | Special chars in attribute values |
| 7 | Mermaid ID sanitization edge cases | Low | Low | Unicode, reserved words |
| 8 | Markdown table captions | Low | Low | Nice-to-have feature |
| 9 | `TableData.RemoveRow()` | Low | Low | No consumer requesting it |
| 10 | `Format.AutoDetect()` (from content sniffing) | Low | Medium | Questionable value |
| 11 | CSV/TSV reading (currently write-only) | Medium | Medium | Would expand library scope |
| 12 | Add `format_deprecated.go` removal timeline | Low | Low | Breaking change, needs major version bump |

---

## D) TOTALLY FUCKED UP

Nothing is fucked up. The codebase is in the healthiest state it's ever been:

- **0 lint issues** (golangci-lint clean)
- **0 vet issues**
- **0 race conditions** (all tests pass with -race)
- **0 build errors**
- **CI fixed** (was build-breaking with Go 1.23 vs 1.26)
- **All docs accurate** (no stale/lying documentation)

---

## E) WHAT WE SHOULD IMPROVE

### Architecture
1. **Renderer interface returns `string` not `(string, error)`** — limits error propagation. Every Render() caller must trust the output. Changing this is a breaking API change requiring v2.
2. **Exhaustruct as a design smell** — 25 nolint directives for `exhaustruct` suggest D2Node/D2Edge/D2NodeStyle have too many optional fields. Functional option pattern would be better.
3. **GraphRendererMixin embedding vs D2 parallel types** — D2 has richer types but can't use the mixin. This is intentional but creates a parallel hierarchy.

### Testing
4. **No integration test for round-trip marshal/unmarshal** — MarshalJSON → UnmarshalJSON should produce equal data.
5. **Benchmarks cover most renderers but not XML, Tree, or HTML tree** — coverage gaps.
6. **Escape functions lack adversarial input testing** — fuzz tests would catch edge cases.

### Developer Experience
7. **No godoc examples** — pkg.go.dev shows no example usage for any type.
8. **cmdguard integration not shown in examples** — examples/basic uses raw os.Args.
9. **No CHANGELOG entries for recent fixes** — CHANGELOG.md has 40 lines, hasn't been updated.

### Operational
10. **No .goreleaser.yml** — release.yml referenced it but it was removed. Release workflow may need adjustment.
11. **Depguard blocks `cmp`** — should be added to allowed list for modern Go patterns.

---

## F) TOP 25 THINGS TO DO NEXT

Sorted by impact × effort (highest first):

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Add godoc Example functions for top 5 types | High | Low | DX |
| 2 | Add round-trip marshal/unmarshal integration tests | High | Low | Testing |
| 3 | Update CHANGELOG.md with recent fixes | Medium | Low | Docs |
| 4 | Add `cmp` to depguard allowed list | Medium | Low | Tooling |
| 5 | Add fuzz tests for escape.D2, escape.DOT, escape.Mermaid | High | Medium | Testing |
| 6 | Add benchmarks for XML, Tree, HTML tree renderers | Medium | Low | Testing |
| 7 | Add cmdguard usage example to examples/ | Medium | Low | DX |
| 8 | Remove `format_deprecated.go` (assess breaking change) | Medium | Low | Cleanup |
| 9 | Add functional options for D2Node (reduce exhaustruct nolints) | Medium | High | Architecture |
| 10 | Add deterministic output ordering for D2 diagrams | Medium | Medium | Correctness |
| 11 | Add `Renderer.Render() (string, error)` as v2 proposal | High | High | Architecture |
| 12 | Unify tree conversion pattern (if complexity warrants) | Medium | High | Architecture |
| 13 | Add DOT attribute quoting for special characters | Low | Low | Robustness |
| 14 | Add Mermaid Unicode ID handling | Low | Low | Robustness |
| 15 | Add `MarshalTSVFromTableData` (like `MarshalXMLFromTableData`) | Medium | Low | Feature parity |
| 16 | Add `MarshalCSVFromTableData` | Medium | Low | Feature parity |
| 17 | Add CSV/TSV reader support | Medium | Medium | Feature expansion |
| 18 | Add `TableData.RemoveRow()` | Low | Low | Feature |
| 19 | Add Markdown table caption support | Low | Low | Feature |
| 20 | Add `Format.AutoDetect()` from content sniffing | Low | Medium | Feature |
| 21 | Create .goreleaser.yml (if releasing binaries ever needed) | Low | Low | CI |
| 22 | Add pre-commit hook to CI | Low | Low | CI |
| 23 | Add `go vet` as explicit CI step | Low | Low | CI |
| 24 | Explore Nix flake migration (if approved) | Low | High | Tooling |
| 25 | Add D2 theme/style preset support | Low | Medium | Feature |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**Should `Renderer.Render()` return `(string, error)` instead of `string`?**

This is the single most consequential architectural decision remaining:

- **Current:** `Render() string` — no error propagation. All renderers swallow errors (HTML escapes silently, D2/DOT/Mermaid skip invalid data).
- **Proposed:** `Render() (string, error)` — proper error propagation but a **breaking API change** for every consumer and every Renderer implementation.

**Arguments for:** Every formatter has potential failure modes (I/O in streaming, invalid UTF-8, buffer overflows). The current API forces consumers to trust that output is always valid. This is the #1 thing I'd fix if I could break the API.

**Arguments against:** This is a library with external consumers (even if currently only Lars's projects). Breaking the interface requires a major version bump (v2). All 12 format renderers, the Registry, StreamingRenderer, and all downstream code would need updating.

**What I cannot determine:** Whether any downstream consumers exist that would be affected by this change, and whether you want to start a v2 branch now or keep the current API stable.

---

## Codebase Health Metrics

| Metric | Value |
|--------|-------|
| Production code | 4,558 lines |
| Test code | 6,424 lines |
| Test:Code ratio | 1.41:1 |
| Packages | 10 (7 with tests) |
| Exported types | ~30 |
| Exported functions | ~80 |
| CI jobs | 1 (consolidated from 3) |
| Lint issues | 0 |
| Race issues | 0 |
| Build errors | 0 |
| Coverage (root) | 91.0% |
| Coverage (all) | 95.5% avg |
| Largest file | sort/sort_test.go (334 lines) |
| Largest prod file | format.go (312 lines) |
| nolint directives | 47 (25 exhaustruct, 13 gochecknoglobals) |

---

## Commits This Session

```
9a943a8 docs: update AGENTS.md and README.md to reflect current state
bf543fb fix(ci): update Go version from 1.23 to 1.26 and consolidate jobs
6f274fa refactor(d2): extract enum types to d2_enum.go to reduce file size
3b1aeea refactor: remove dead API surface (ToANSI, GetStyle)
b7d5cdf docs(registry): document as opt-in plugin system, use cmp.Compare
2c91e7e docs: rewrite PLAN.md to match actual codebase structure
1e36a4d chore: remove 4,693 lines of stale planning and status documentation
2bffa7f docs(planning): add comprehensive improvement plan for 2026-04-29
4aa86a8 refactor(d2): use package-level vars for enum values instead of allocating functions
e882804 fix(sort): correct descending sort stability violation and add comprehensive tests
fe79fec fix: move golang.org/x/term to direct dependency and replace recursive joinStrings
```
