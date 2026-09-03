# Status Report — 2026-06-19 00:25

## go-output: Modularization Audit + Damage Recovery

---

## a) FULLY DONE

| Item                                                                       | Status   |
| -------------------------------------------------------------------------- | -------- |
| Phase 1-3: Module landscape analysis (18 modules, DAG, coupling)           | Done     |
| Phase 4: Self-review identified version drift, missing require, stale docs | Done     |
| Phase 6: Version drift fixed (v0.6.3→v0.12.0 across 10 modules)            | Done     |
| Phase 6: graphtest missing root require fixed                              | Done     |
| Phase 6: enum/go.mod testhelpers version fixed (v0.6.3→v0.12.0)            | Done     |
| **REVERT**: Ill-advised merge of enum/envdetect/escape into root — UNDONE  | Done     |
| All 18 modules restored as standalone                                      | Verified |
| All 18 modules build (workspace + isolated)                                | Verified |
| All 18 modules lint at 0 issues                                            | Verified |
| All tests pass                                                             | Verified |
| Commit `8af32f2` with BuildFlow pre-commit passing all 36 checks           | Done     |

---

## b) PARTIALLY DONE

| Item                                                  | Status                                         | Remaining     |
| ----------------------------------------------------- | ---------------------------------------------- | ------------- |
| Root bloat analysis                                   | Identified 5 concern clusters (1908 lines)     | Not yet split |
| go.work + replace redundancy                          | Documented (Failure Mode #4)                   | Not resolved  |
| Stale `internal/` docs in AGENTS.md                   | Identified                                     | Not fixed     |
| `.golangci.yml` allow-lists for enum/envdetect/escape | Still present (correct for standalone modules) | OK            |

---

## c) NOT STARTED

| Item                                      | Impact   | Notes                                                                                        |
| ----------------------------------------- | -------- | -------------------------------------------------------------------------------------------- |
| Split root into core/ + markdown/ + tree/ | **HIGH** | Root has 37 files (1908 prod lines). Core types mixed with format renderers. See plan below. |
| Remove replace directives, commit go.work | Medium   | Correct for published libs but changes dev workflow                                          |
| Tag envdetect on proxy                    | Medium   | envdetect was NEVER tagged — all modules need replace directives for it                      |
| Remove stale proposal doc                 | Low      | `docs/modularization/2026-06-18_PROPOSAL.md` describes the reverted merge                    |

---

## d) TOTALLY FUCKED UP

### The Merge (c0ad787) — REVERTED

**What happened:** I merged enum/, envdetect/, escape/ INTO root, reducing 18→15 modules.

**Why it was wrong:**

1. **Wrong direction** — user wants root SMALLER, not bigger. I added 8 files to root.
2. **Misread the Core Invariant** — the invariant is about dependency _weight_ isolation (no lipgloss/bubbletea/yaml in root), not forbidding zero-dep imports.
3. **Destroyed composability** — consumers can no longer `go get go-output/enum` independently.
4. **Pre-commit hook chaos** — the go-mod-tidy hook fought my changes, sweeping go.mod deletions into unrelated commits (e71a6d2, 071ba9d). Took significant effort to untangle.
5. **Revert was painful** — `git revert` blocked by pre-commit hook, conflicts with later commits, manual restoration needed.

**Status:** Fully reverted in commit `8af32f2`. All 18 modules restored and verified.

---

## e) WHAT WE SHOULD IMPROVE

### 1. Root is a God-Package (37 files, 1908 lines)

Root mixes 5 distinct concerns:

| Cluster                                                                                   | Lines | Files | Could be                                                              |
| ----------------------------------------------------------------------------------------- | ----- | ----- | --------------------------------------------------------------------- |
| **Core types** (Format, Shape, ColorMode, registry, TableData, GraphNode, TreeNode types) | ~556  | 8     | `core/` module — thin types + interfaces                              |
| **Graph state** (GraphRendererState, GraphStyle, edge logic)                              | 359   | 2     | Shared with d2/graph/plantuml — could stay in core or be `graphcore/` |
| **TableData** (TableData struct + render dispatch)                                        | 408   | 2     | Core type — belongs in core/                                          |
| **Markdown renderer** (MarkdownTable)                                                     | 289   | 1     | `markdown/` module — it's a format renderer like `table/`             |
| **Tree renderer** (TreeNode, ASCIITreeRenderer)                                           | 229   | 1     | `tree/` module — it's a shape renderer                                |

### 2. SDK Not Composable Enough

A consumer who wants JUST JSON output currently must import root, which drags in Markdown, Tree, Graph types, ColorMode, etc. The vision:

```
go get go-output/core       # Format, Shape, TableData, interfaces (zero deps)
go get go-output/markdown   # Markdown renderer (depends on core)
go get go-output/serialization  # JSON/YAML/TOML (depends on core)
// Don't want trees? Don't import tree/. Don't want markdown? Don't import markdown/.
```

### 3. envdetect Never Tagged

`envdetect` has NO git tags. Every module that transitively depends on it (through root) needs a replace directive. This is fragile. Should tag it at next release.

### 4. Pre-commit Hook Sweeping Unrelated Changes

The BuildFlow `go-mod-tidy` hook runs `go mod tidy` on commit, which can modify go.mod/go.sum files beyond what was staged. This caused my deletions to be swept into unrelated commits. The hook should either:

- Only tidy staged files' modules, or
- Warn instead of auto-fix, or
- Run in a separate pre-push hook

---

## f) Top 25 Things to Get Done Next

Sorted by **impact / effort ratio** (highest first):

### Tier 1 — High Impact, Low Effort

| # | Task                                                                                                                          | Impact      | Effort |
| - | ----------------------------------------------------------------------------------------------------------------------------- | ----------- | ------ |
| 1 | **Delete stale proposal doc** (`docs/modularization/2026-06-18_PROPOSAL.md`) — describes reverted merge, misleading           | Cleanup     | 1 min  |
| 2 | **Remove stale `internal/` references from AGENTS.md** — references `internal/gentest`, `internal/testutils` that don't exist | Accuracy    | 2 min  |
| 3 | **Tag envdetect/v0.12.0** on next release — eliminates replace directive fragility                                            | Stability   | 5 min  |
| 4 | **Add envdetect to CI release workflow** — it's missing from the tagging automation                                           | Correctness | 10 min |

### Tier 2 — High Impact, Medium Effort

| # | Task                                                                                                                                   | Impact        | Effort  |
| - | -------------------------------------------------------------------------------------------------------------------------------------- | ------------- | ------- |
| 5 | **Extract `markdown/` as standalone module** — MarkdownTable is 289 lines, a format renderer like `table/`. Root shrinks by 289 lines. | Composability | 1-2 hrs |
| 6 | **Extract `tree/` as standalone module** — ASCIITreeRenderer is 229 lines, a shape renderer. Root shrinks by 229 lines.                | Composability | 1-2 hrs |
| 7 | **Move GraphRendererState to core or graphcore/** — shared state used by d2/graph/plantuml. 359 lines. Root shrinks further.           | Composability | 2-3 hrs |
| 8 | **Define `core/` module** — Format, Shape, ColorMode, registry interfaces, TableData/GraphNode/TreeNode types. The thin API surface.   | Composability | 3-4 hrs |

### Tier 3 — Medium Impact, Medium Effort

| #  | Task                                                                                                    | Impact          | Effort |
| -- | ------------------------------------------------------------------------------------------------------- | --------------- | ------ |
| 9  | **Replace replace directives with committed go.work** — cleaner for multi-module repo, eliminates drift | Maintainability | 1 hr   |
| 10 | **Add `go work sync` to setup-workspace** — prevent stale go.work after module changes                  | DX              | 30 min |
| 11 | **Audit all exported types for naming quality** — load `naming-review` skill                            | Quality         | 2 hrs  |
| 12 | **Review if `direction.go` belongs in root** — 40 lines, only used by graph/. Could move.               | Composability   | 30 min |
| 13 | **Review if `ids.go` belongs in root** — 58 lines, branded ID re-exports. Could stay.                   | Composability   | 30 min |

### Tier 4 — Lower Priority

| #  | Task                                                                                        | Impact         | Effort  |
| -- | ------------------------------------------------------------------------------------------- | -------------- | ------- |
| 14 | **Consolidate ColorMode into core** — it's a cross-cutting type used by multiple renderers  | Architecture   | 1 hr    |
| 15 | **Review streaming.go** (53 lines) — is it core or should it be separate?                   | Architecture   | 30 min  |
| 16 | **Consider value-object pattern for Format/Shape** — make impossible states unrepresentable | Type safety    | 2 hrs   |
| 17 | **Run `deduplicate-code` skill** — check for duplication across root + sub-modules          | Quality        | 1 hr    |
| 18 | **Add integration test for isolated module builds** — CI should verify GOWORK=off builds    | CI             | 1 hr    |
| 19 | **Document module dependency DAG** in FORMAT_ARCHITECTURE.md                                | Docs           | 30 min  |
| 20 | **Review if testhelpers/graphtest should merge into testhelpers** — nested module, 2 files  | Simplification | 1 hr    |
| 21 | **Consider `cmd/` or `cli/` module for examples** — examples currently imports everything   | Architecture   | 2 hrs   |
| 22 | **Audit go.sum files for unnecessary entries** — clean up after module changes              | Hygiene        | 30 min  |
| 23 | **Add `make docs` equivalent for module graph visualization** — auto-generate DAG diagram   | DX             | 2 hrs   |
| 24 | **Review nom/ (49 files) for god-package** — largest module, may need internal split        | Quality        | 3 hrs   |
| 25 | **Consider adopting `go-runtime` or `x/exp` patterns** — leverage stdlib ecosystem          | Modernization  | Ongoing |

---

## g) Top #1 Question I Cannot Figure Out

**Should root's shared types (TableData, GraphNode, TreeNode, ColorMode) move to a `core/` module, or stay in root?**

The tension:

- **Moving to `core/`**: Cleaner separation. Root becomes just the registry + dispatch. Consumers import `core/` for types, `root/` for dispatch. But root IS currently the "core" — creating `core/` means root becomes nearly empty (just registry.go + marshal.go = 75 lines), which seems over-modularized.
- **Keeping in root**: Root stays as the "core + shared types" module. Other format renderers (markdown, tree) extract OUT. Root shrinks from 1908 → ~1200 lines by extracting markdown + tree + graph state. Still compositional enough — consumers who don't want markdown just don't import `markdown/`.

**The real question:** Is the goal to make `go get go-output` as thin as possible (→ extract core/), or to make root lean but still the canonical import (→ extract renderers OUT, keep types IN root)?

This determines whether we create 1 new module (core/) or 3 new modules (markdown/, tree/, graphcore/).
