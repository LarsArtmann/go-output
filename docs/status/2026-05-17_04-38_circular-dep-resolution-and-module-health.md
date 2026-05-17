# go-output — Comprehensive Status Report

**Date:** 2026-05-17 04:38
**Session:** Session 5 (continuation of interrupted Session 4)
**Branch:** master (clean, up to date with origin)
**Version:** v0.4.0 (latest tag)
**Coverage:** 90.3% (root), 100% (enum, escape, sort, cmdguard, table)

---

## A. FULLY DONE ✅

### 1. Public Launch (Session 1)
- [x] MIT license (was PROPRIETARY)
- [x] README.md rewritten for general audience
- [x] 27 doc comments on exported symbols
- [x] GitHub release v0.4.0 created
- [x] Pushed to remote

### 2. Shape Capability Matrix (Session 2)
- [x] `Shape` type with `ShapeTable`/`ShapeTree`/`ShapeGraph` constants
- [x] `formatCapabilities` map (single source of truth)
- [x] `Supports(Shape)`, `Shapes()`, `FormatsForShape()` methods
- [x] Full `ParseShape`/`IsValid`/`AllowedValues`/`String` enum methods
- [x] `FormatCategory`/`IsTableFormat`/`IsTreeFormat`/`IsGraphFormat`/`Category()` deprecated
- [x] ADR 002 written and committed
- [x] All tests updated (root, integration)

### 3. Structural Refactoring (Sessions 3-4)
- [x] `TreeNode`/`TreeOutputRenderer` extracted to `tree.go`
- [x] `TableData`/`RowEdge`/`tableDataBase` extracted to `tabledata.go`
- [x] `GraphRendererMixin` moved from `dot.go` to `graph.go`
- [x] Stale `PLAN.md` deleted
- [x] All files under 350-line limit

### 4. Circular Dependency Resolution (Session 5 — THIS SESSION)
- [x] Deleted `sort/sorter.go` (deprecated `Sorter[T]` with circular dep on root)
- [x] Deleted `sort/sort_test.go` (415 lines of tests for deprecated type)
- [x] Modernized `sort/compare.go` — `ByField` now returns `int` (for `slices.SortStableFunc`) instead of `bool` (for old `sort.SliceStable`)
- [x] Rewrote `sort/compare_test.go` — self-contained, zero external deps
- [x] `sort/go.mod` now has ZERO dependencies (was importing root + 8 transitive deps)
- [x] Removed `sort` from root `go.mod` requires and replace directives
- [x] Removed `sort` from `integration/go.mod`
- [x] Updated `userjourney_test.go` — uses stdlib `slices.SortStableFunc` + `cmp.Compare`
- [x] Updated `integration/workflow_test.go` — same migration

### 5. cmdguard/ Fix (Session 5 — THIS SESSION)
- [x] Fixed empty `cmdguard/go.mod` — added all required deps + replace directives
- [x] Generated proper `cmdguard/go.sum`
- [x] Removed `internal/gentest` import from `cmdguard/cmdguard_test.go`
- [x] Inlined `assertStringSliceEqual` helper (was the only gentest function used)

### 6. enum/ Cross-Module Fix (Session 5 — THIS SESSION)
- [x] Removed `internal/gentest` import from `enum/enum_test.go`
- [x] Inlined `assertStringSliceEqual` helper
- [x] `enum/go.mod` now correctly has ZERO external deps (already had none, but test was broken without go.work)

### 7. table/ go.sum Fix (Session 5 — THIS SESSION)
- [x] Added missing `go-branded-id` to `table/go.mod` via `go mod tidy`

### 8. Quality Verification
- [x] **7/7 modules build:** root, enum, escape, sort, cmdguard, table, integration
- [x] **7/7 modules pass tests:** all green
- [x] **7/7 modules pass lint:** 0 issues across all
- [x] **7/7 modules race-clean:** no data races
- [x] **Root coverage:** 90.3%

---

## B. PARTIALLY DONE 🟡

### 1. CHANGELOG.md
- **Current state:** Exists but stale at v0.1.0. Has `[Unreleased]` section with TSV, XML, enum, and deprecation entries. Missing 30+ features since January.
- **What's missing:** D2, Mermaid, DOT, HTML Tree, GraphRenderer, streaming, ColorMode, branded IDs, escape module, multi-module workspace, Shape capability matrix, public launch, structural refactoring, circular dep resolution, cmdguard fix.

### 2. FORMAT_ARCHITECTURE.md
- **Current state:** Still references old "Format Categories" (Table/Tree/Graph as exclusive categories) instead of Shape capability matrix.
- **What's missing:** Update to reflect Shape system, `Supports()` API, `formatCapabilities` map.

---

## C. NOT STARTED ⬜

1. **New formats:** TOML, JSONL, PlantUML, AsciiDoc
2. **Shape-specific renderers:** `NewJSONTableRenderer(data *TableData) Renderer`, `NewYAMLTableRenderer(data *TableData) Renderer` — capability matrix declares JSON/YAML support Table/Tree/Graph but there are no typed renderers
3. **CHANGELOG.md update** — needs complete rewrite from v0.1.0 to current
4. **FORMAT_ARCHITECTURE.md update** — Shape system
5. **Visibility/marketing:** r/golang post, Awesome Go submission, blog post
6. **LSP workspace errors:** 120+ phantom errors from gitignored `go.work` — could un-gitignore or configure LSP
7. **CONTRIBUTING.md review** — references CHANGELOG.md which exists but is stale
8. **AGENTS.md update** — reflect sort/ circular dep resolution, cmdguard fix, new module dependency graph

---

## D. TOTALLY FUCKED UP 💥

### 1. `internal/gentest` Package Cross-Module Contamination
**Status:** Fixed in this session, but the root cause remains.

The `internal/gentest` package is defined in the root module but was imported by:
- `cmdguard/cmdguard_test.go` — **FIXED** (inlined helper)
- `enum/enum_test.go` — **FIXED** (inlined helper)

Go's `internal` package restriction means sub-modules literally cannot import it (the compiler blocks it). The only reason it ever worked was because `go.work` merged everything into one module at the Go toolchain level. Without `go.work` (which is gitignored), the sub-modules are completely broken.

**Root cause:** `internal/` packages are scoped to their module, not their directory tree. In a multi-module workspace, `internal/gentest` is private to the root module only.

**Remaining contamination:** `internal/testutils` is still imported by root tests, integration tests, and user journey tests. It works because those tests either run in root's module context or integration has root as a direct dep. But `gentest.AssertStringSliceEqual` is now duplicated as `assertStringSliceEqual` in both `cmdguard/` and `enum/`. This is a DRY violation.

**Proper fix:** Extract shared test helpers into a non-internal package (e.g., `testhelpers/`) that any module can import.

### 2. No `go.work` in Repo
**Status:** By design (gitignored), but causes 120+ LSP errors.

Every IDE/LSP shows broken red squiggles for `cmdguard/`, `enum/`, `integration/`, etc. because without `go.work`, gopls can't resolve cross-module imports. The `go build` and `go test` commands work fine individually per-module, but the IDE experience is terrible.

---

## E. WHAT WE SHOULD IMPROVE 🔧

### Architecture
1. **Extract shared test helpers** from `internal/gentest` → `testhelpers/` (non-internal, importable by sub-modules)
2. **Create `go.work` at dev time** — keep it gitignored but add a `go.work.example` or `Makefile` target
3. **JSON/YAML renderers** — the capability matrix says they support Table/Tree/Graph but they're just `Marshal(v any)` functions. Add typed renderers.
4. **Consider removing `sort/` entirely** — now that it's zero-dep with just `ByField`, it's barely worth a separate module. Could move `ByField` to root or delete entirely (stdlib is sufficient).

### Documentation
5. **CHANGELOG.md** is embarrassingly stale — 4 months of unreleased changes
6. **FORMAT_ARCHITECTURE.md** references old category system
7. **AGENTS.md** needs update for sort/cmdguard fixes
8. **CONTRIBUTING.md** references CHANGELOG.md update process — ensure it's accurate

### Code Quality
9. **`format_deprecated.go`** still has 30 lines of deprecated code — plan removal timeline
10. **Registry (`registry.go`)** is opt-in but undocumented in README
11. **`SortBy` enum** in root is now only used by `cmdguard` tests — consider its purpose
12. **`internal/testutils/`** and `internal/gentest/` have overlapping responsibilities

### DevEx
13. **LSP errors** make IDE unusable for sub-modules without manual `go.work` creation
14. **No `go.work.example`** — new contributors can't figure out why their IDE is broken
15. **Module-level `go mod tidy`** must be run per-module — no single command

---

## F. TOP 25 THINGS TO DO NEXT 🎯

| # | Priority | Task | Impact | Effort |
|---|----------|------|--------|--------|
| 1 | P0 | **Update CHANGELOG.md** — comprehensive update for v0.4.0 → v0.5.0 | High | Low |
| 2 | P0 | **Update FORMAT_ARCHITECTURE.md** — Shape system | Medium | Low |
| 3 | P0 | **Update AGENTS.md** — reflect sort/cmdguard fixes, new dep graph | Medium | Low |
| 4 | P1 | **Create `go.work.example`** — document workspace setup for contributors | High | Trivial |
| 5 | P1 | **Extract test helpers** from `internal/gentest` → `testhelpers/` | Medium | Medium |
| 6 | P1 | **Add `NewJSONTableRenderer(data *TableData) Renderer`** | High | Medium |
| 7 | P1 | **Add `NewYAMLTableRenderer(data *TableData) Renderer`** | High | Medium |
| 8 | P1 | **Commit this session's work** and push | High | Trivial |
| 9 | P2 | **Tag v0.5.0** — circular dep fix, cmdguard fix, ByField modernization | High | Trivial |
| 10 | P2 | **Add TOML format** — closes serialization format trio | Medium | Medium |
| 11 | P2 | **Add JSONL format** — streaming/log use cases | Medium | Medium |
| 12 | P2 | **Decide `sort/` fate** — keep as zero-dep convenience or delete entirely | Low | Low |
| 13 | P2 | **Post to r/golang** — public visibility | High | Low |
| 14 | P2 | **Submit to Awesome Go** — discoverability | High | Low |
| 15 | P2 | **Write blog post** — Shape capability matrix is a unique selling point | Medium | High |
| 16 | P3 | **Add PlantUML format** — diagram variety | Low | Medium |
| 17 | P3 | **Add AsciiDoc format** — documentation pipeline | Low | Medium |
| 18 | P3 | **Plan deprecated code removal** — `FormatCategory`, `OutputFormat` aliases | Low | Low |
| 19 | P3 | **Evaluate `cmdguard/` extraction** — separate repo or keep? | Low | Low |
| 20 | P3 | **Add `FormatsForShape` to registry** — runtime dispatch by shape | Medium | Medium |
| 21 | P4 | **Add D2 `ShapeTable` renderer** — D2FromTableData exists but isn't a Renderer | Low | Low |
| 22 | P4 | **Benchmark comparison** — old Sorter vs stdlib slices.SortStableFunc | Low | Trivial |
| 23 | P4 | **CI/CD pipeline** — GitHub Actions for build/test/lint across all modules | Medium | Medium |
| 24 | P4 | **API stability review** — pre-v1.0 audit of all exported symbols | Medium | Medium |
| 25 | P4 | **CONTRIBUTING.md review** — ensure accuracy for external contributors | Low | Low |

---

## G. TOP #1 QUESTION 🤔

**Should `sort/` continue to exist as a module?**

The `ByField` helper is now the only exported function. It's a 5-line wrapper around `cmp.Compare`:

```go
func ByField[T any, F cmp.Ordered](extract func(T) F) func(a, b T) int {
    return func(a, b T) int {
        return cmp.Compare(extract(a), extract(b))
    }
}
```

Arguments for keeping:
- Convenience — saves users from writing the `cmp.Compare(extract(a), extract(b))` pattern
- Existing users may import it

Arguments for deleting:
- The entire deprecation notice says "use stdlib instead"
- It's barely worth a separate Go module for one 5-line function
- Zero external deps means it could live in root if kept at all
- Removes a module from the workspace maintenance burden

**My recommendation:** Delete the `sort/` module entirely. Move `ByField` to the root package as `CompareByField` (renamed to match `cmp.Compare` convention) or just delete it and let users write the one-liner. The stdlib is sufficient.

---

## Session 5 Changes Summary

### Files Modified
| File | Change | Lines |
|------|--------|-------|
| `sort/sorter.go` | **DELETED** — circular dep eliminated | -61 |
| `sort/sort_test.go` | **DELETED** — tests for deleted type | -415 |
| `sort/compare.go` | Modernized `ByField` return type `bool→int`, updated docs | 25 |
| `sort/compare_test.go` | Full rewrite, zero external deps | 118 |
| `sort/go.mod` | Stripped to 3 lines (zero deps) | -18 |
| `sort/go.sum` | **DELETED** — no longer needed | -12 |
| `go.mod` | Removed sort dep + replace directive | -2 |
| `userjourney_test.go` | Stdlib sort, removed sort import | +12/-16 |
| `integration/go.mod` | Removed sort dep + replace | -2 |
| `integration/workflow_test.go` | Stdlib sort, removed sort import | +5/-4 |
| `cmdguard/go.mod` | Added all deps + replace directives | +21 |
| `cmdguard/cmdguard_test.go` | Removed `internal/gentest` import, inlined helper | +16/-2 |
| `cmdguard/go.sum` | **CREATED** — proper dependency checksums | +12 |
| `enum/enum_test.go` | Removed `internal/gentest` import, inlined helper | +16/-2 |
| `table/go.mod` | Added missing `go-branded-id` | +1 |
| `table/go.sum` | Added `go-branded-id` checksums | +2 |

### Net Impact
- **Circular dependency:** ELIMINATED (sort/ → root cycle broken)
- **Modules with zero external deps:** enum, escape, sort (was: enum, escape only)
- **`internal/gentest` imports from sub-modules:** 0 (was: 2 — cmdguard, enum)
- **LSP errors:** ~118 remaining (all from gitignored go.work, not real issues)
- **Test coverage:** 90.3% root, 100% all sub-modules
- **Lint issues:** 0 across all 7 modules
- **Race conditions:** 0 across all 7 modules

### Dependency Graph (Post-Fix)
```
root (output) → enum, escape, yaml, x/term, go-branded-id
enum          → (none)
escape        → (none)
sort          → (none) ← FIXED: was root→sort→root circular
cmdguard      → root (tests only; prod code is standalone)
table         → root, lipgloss/v2
integration   → root, table
examples      → root, table
```
