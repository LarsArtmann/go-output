# Status Report — 2026-06-19 03:15

## go-output: V1.0 Composability & Hardening Sprint — COMPLETE

**Session scope:** Pareto-planned and executed root god-package extraction + type safety + docs + CI hardening.
**HEAD:** `ff3d9f4` · **Branch:** `master` (up to date with origin)
**Commits this session:** 10 (`ab59137` → `ff3d9f4`)
**Files changed:** 80 (+1617 / −826)

---

## a) FULLY DONE

| Item | Status |
|------|--------|
| **Pareto plan written** with 1%/4%/20% breakdown + mermaid execution graph | Done (`docs/planning/2026-06-19_02-10_*.md`) |
| **Bug #1 fixed: timing-cache test isolation** — `WithCachePath` option on `TimingCache` + `NOMStyleSubscriber`; `newTestSubscriber(t)` helper injected across 34 test sites; suite is now hermetic (no `~/.cache` writes) | Done |
| **Bug #2 fixed: render_tabledata_test.go 355→265+98 lines** — registry tests extracted to `render_registry_test.go` | Done |
| **`go mod tidy` across all modules** — escape version residue confirmed resolved | Done |
| **Stale misleading proposal deleted** (`docs/modularization/2026-06-18_PROPOSAL.md`) — described the reverted merge | Done |
| **AGENTS.md `internal/gentest`/`internal/testutils` false references corrected** — no `internal/` dir exists | Done |
| **TODO_LIST.md false claims corrected** — timing-cache isolation now covers subscriber tests; test-file limit now includes the 5th file | Done |
| **`markdown/` module extracted** — 289 lines moved out of root; self-registers `FormatMarkdown` via `init()`; all consumers updated | Done |
| **`tree/` module extracted** — 229 lines moved out of root (ASCIITreeRenderer); TreeNode stays in root (shared by 15+ modules); self-registers `FormatTree` via `init()` | Done |
| **Root no longer registers ANY format** — all 16 formats self-register from their sub-modules via `init()` | Done |
| **AGENTS.md module map updated** (18→20 modules) | Done |
| **Typed `Symbol` constants** — `type Symbol string`; 10 constants typed; propagated through `Activity.Symbol`, `GetSymbol()`, `OperationSymbol()`, `formatTimingWithSymbol()`; all callers updated | Done |
| **`ActivityStatus` enum methods** — `ParseActivityStatus()`, `IsValid()`, `AllowedValues()`, `AllActivityStatuses` slice, `InvalidActivityStatusError` type | Done |
| **README updated** — all 15 user-facing sub-modules listed with descriptions, grouped by category | Done |
| **`RELEASE.md` written** — 20-module mono-version bump workflow, pre-release checklist, tagging, proxy verification, rollback | Done |
| **`docs/MIGRATION_v0.12_to_v1.0.md` written** — all breaking changes documented (markdown/tree extraction, typed Symbol, GetSymbol return type) | Done |
| **ADR 008 written** — dedup workflow decision (5-step checklist, t=24 threshold rationale) | Done |
| **CI module loops fixed** — markdown + tree were MISSING from all 9 CI loops (build, test, coverage, tidy, govulncheck in ci.yml + 4 in release.yml); now added everywhere | Done |
| **art-dupl CI gate added** — `duplication` job runs art-dupl at t=50 (non-blocking warning) | Done |
| **All 20 modules build** (workspace + isolated GOWORK=off for markdown/tree) | Verified |
| **All 20 modules test green** (763+ test functions) | Verified |
| **Race tests clean** (nom + tui) | Verified |
| **Lint: 0 issues across all 20 modules** | Verified |
| **Zero TODO/FIXME/HACK in prod code** | Verified |
| **All test files ≤350 lines** | Verified |

---

## b) PARTIALLY DONE

| Item | Status | Gap |
|------|--------|-----|
| **Type-safety hardening** | Symbol type + ActivityStatus enums done | **Branding `ActivityID`/`WorkflowID` via `go-branded-id` deferred** — 293 references across 20 nom files; large refactor requiring careful test updates |
| **GraphStyle typed colors** | Not done | `GraphStyle.Fill`/`Stroke`/`FontColor` still bare `string`; branded `Color` type would cascade across graph/, d2/, plantuml/ |
| **benchstat CI baseline storage** | Benchmarks run in CI but output is discarded | No stored baseline, no benchstat comparison, no regression threshold |
| **README module list** | 15 of 19 sub-modules listed | Omits `testhelpers`, `testhelpers/graphtest`, `bdd`, `integration` — intentionally, these are test/dev-only modules |
| **Pre-commit Go hooks** | BuildFlow hook already runs gofumpt + golangci-lint + gomod-check | No standalone `go vet` step in git-hooks.nix (covered by golangci-lint's `vet` linter) |

---

## c) NOT STARTED

| Item | Impact | Notes |
|------|--------|-------|
| **Brand `ActivityID`/`WorkflowID`** via `go-branded-id` | 🔥🔥 | Eliminates ID-mixing class of bugs at compile time. 293 references. Requires touching every constructor, handler, test. |
| **Typed `Color` for `GraphStyle`** | 🟡 | Branded Color type for Fill/Stroke/FontColor. Cascades to d2/, graph/, plantuml/ conversion code. |
| **`core/` module extraction** | 🔥 | Move shared types (Format, Shape, ColorMode, TableData, GraphNode, registry interfaces) to a thin `core/` module. Root becomes registry + dispatch only. Riskier — may over-modularize. |
| **`graphcore/` module extraction** | 🟡 | Move `GraphRendererState` + graph state out of root. 359 lines. Shared by d2/graph/plantuml. |
| **`direction.go` relocation** | 🟢 | 40 lines; only used by graph modules. Could move to `graphcore/`. |
| **Tag `envdetect/v0.12.0`** | 🟡 | envdetect was NEVER tagged. All modules use replace directives. Tagging eliminates fragility. |
| **benchstat stored baseline** | 🟡 | Benchmarks execute but output is discarded — no comparison, no regression detection. |
| **`go test -race` for ALL modules in CI** | — | Already done! CI already runs `-race` for every module. Previously thought missing; verified present. |
| **BDD spec name verification** | 🟢 | Spec names should match post-extraction type names. Not verified this session. |
| **GitHub release notes draft for v1.0.0** | 🟡 | Template not written. `docs/RELEASE_NOTES_v1.0.0.md` doesn't exist. |
| **`TableData` API decision** (fields vs getters) | 🔥🔥🔥 | **Sole v1.0.0 freeze blocker.** Blocked on owner. Option A (fields only) / B (getters+setters) / C (keep both). |

---

## d) TOTALLY FUCKED UP

**Nothing is fucked up this session.** All changes are verified green.

### What almost went wrong (caught and avoided):

1. **CI module lists missing markdown + tree** — if we hadn't caught this, CI would not build/test/lint the 2 new modules. Caught during Tier 4 verification; fixed before push.
2. **Root replace path `../markdown` instead of `./markdown`** — root is the parent, so the path is `./markdown`, not `../markdown` (which is for sibling sub-modules). Caught immediately by build failure; fixed in same step.
3. **bdd registrations test didn't blank-import markdown** — BDD tests that call `RenderTableData(data, FormatMarkdown, ...)` would fail because markdown's `init()` never ran. Caught by test failure; fixed by adding `_ "github.com/larsartmann/go-output/markdown"` import.
4. **Conflicting replace directives** — sibling modules got `./tree` instead of `../tree` from `go mod edit` default behavior. Caught by gopls diagnostic; fixed for all 14 sibling modules + nested graphtest.
5. **gopls stale diagnostics** — gopls reported duplicate declarations from deleted files throughout the session. Verified actual state with `go build`/`go test` (always green) rather than trusting gopls cache. Restarted gopls twice to clear cache.

---

## e) WHAT WE SHOULD IMPROVE

### 1. Root is still 1401 lines (was 1908 — 27% reduction, but still large)

The 2 biggest remaining clusters:
- **Graph state** (`graph.go` 282L + `graph_tabledata.go` 77L = 359L) — shared by d2/graph/plantuml; candidate for `graphcore/` extraction.
- **TableData** (`tabledata.go` 245L + `render_tabledata.go` 115L = 360L) — core type + dispatch.

### 2. `envdetect` was never tagged

Every module that transitively depends on it (through root) needs a replace directive. This is the most fragile part of the module graph. Tagging at next release eliminates this.

### 3. ID-mixing bugs still possible in nom

`ActivityID` and `WorkflowID` are both `type X string` — they're assignable to each other and to plain `string` at compile time. Graph/d2/tree IDs already use `go-branded-id` phantom types. nom IDs do not. This is a real class of bugs waiting to happen.

### 4. Replace directive explosion (44+ directives)

The markdown + tree extraction added 30+ new replace directives across sibling modules (root's test imports cascade through the graph). The mono-version tagging strategy mitigates this for consumers, but the dev experience is increasingly fragile. Consider committing `go.work` or adopting a workspace-first workflow.

### 5. Benchmarks run but are discarded

CI runs benchmarks but doesn't store baselines. Performance regressions go undetected until a user complains.

---

## f) Top 25 Things to Get Done Next

Sorted by **impact / effort ratio** (highest first):

### Tier 1 — Owner-blocked (cannot execute autonomously)

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | **Decide `TableData` fields vs getters for v1** (Option A/B/C) | Critical | 5m |
| 2 | **Cut `v1.0.0` tag** — API declared frozen (ADR 006) | Critical | 10m |
| 3 | **Tag `envdetect/v0.12.0`** — eliminates replace fragility | High | 5m |
| 4 | **Submit to r/golang + Awesome Go** | High | 30m |

### Tier 2 — High-impact code work

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 5 | **Brand `ActivityID`/`WorkflowID`** via `go-branded-id` (293 refs) | 🔥🔥 | 2 hrs |
| 6 | **Typed `Color` for `GraphStyle`** Fill/Stroke/FontColor | 🟡 | 30m |
| 7 | **Extract `graphcore/`** — GraphRendererState + graph state (359L) | 🔥 | 3 hrs |
| 8 | **Decide: `core/` module or keep types in root** | 🔥 | Design |
| 9 | **BDD spec name verification** — specs match post-extraction names | 🟡 | 15m |

### Tier 3 — CI/process hardening

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 10 | **benchstat CI step** with stored baseline artifact | 🔥 | 30m |
| 11 | **GitHub release notes draft** for v1.0.0 | 🟡 | 20m |
| 12 | **Commit `go.work` or document workspace-first workflow** | 🟡 | 30m |
| 13 | **Add `go work sync` to setup-workspace** app | 🟡 | 15m |
| 14 | **art-dupl threshold tuning** — verify t=50 doesn't false-positive | 🟢 | 30m |
| 15 | **CI: verify markdown/tree actually tested** after module loop fix | 🟢 | 5m |

### Tier 4 — Documentation

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 16 | **Document module dependency DAG** in FORMAT_ARCHITECTURE.md | 🟡 | 30m |
| 17 | **Update `doc.go`** root package doc — remove stale markdown/tree refs | 🟡 | 10m |
| 18 | **CHANGELOG.md** — document markdown/tree extraction + Symbol type | 🔥 | 20m |
| 19 | **Update FEATURES.md** — mark markdown/ and tree/ as standalone modules | 🟡 | 15m |
| 20 | **Godoc review** — new exported types (Symbol, ActivityStatus methods) documented | 🟡 | 15m |

### Tier 5 — Optional / future

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 21 | **`direction.go` relocation** to graphcore/ (40L, only used by graph) | 🟢 | 15m |
| 22 | **`ids.go` relocation review** (58L, branded ID re-exports) | 🟢 | 15m |
| 23 | **`streaming.go` review** (53L, is it core or separate?) | 🟢 | 15m |
| 24 | **Run `deduplicate-code` skill** across root + new modules | 🟢 | 1 hr |
| 25 | **CBOR format** (ROADMAP) — only on real user demand | 🟢 | 3 hrs |

---

## g) Top #1 Question I Cannot Figure Out

**Should root's shared types (TableData, GraphNode, ColorMode, registry interfaces) move to a `core/` module, or stay in root?**

This is the single architectural decision that determines the project's long-term composability shape.

**The tension:**

- **Moving to `core/`**: Root becomes just registry + dispatch (~300 lines). Consumers import `core/` for types, `root/` for dispatch. Maximally composable — a user who wants only JSON doesn't pull in any format-renderer code. BUT: root becomes nearly empty, which may over-modularize. And the migration is breaking: every consumer's `output.TableData` becomes `core.TableData`.

- **Keeping in root**: Root stays as "core + shared types" (current state, 1401 lines). Format renderers (markdown, tree, future ones) extract OUT. Root is still the canonical import. Less composable — `go get go-output` pulls everything — but simpler and less disruptive. The user already said "we have TOO MANY files in root," which suggests extraction is wanted.

**Why I can't decide:** This is an irreversible commitment. Once v1.0.0 ships, the import path is frozen. Each option has real tradeoffs that only the project owner can weigh:

- Option A (`core/` extraction): maximally composable, but creates a 2-import pattern (`core` + `output`) and a nearly-empty root. 20→22 modules.
- Option B (keep types in root): simpler, less disruptive, but root stays large. The user's complaint about "too many files" is only partially addressed.
- Option C (extract `graphcore/` only): middle ground — move graph state out (359L), keep types in root. Root shrinks to ~1042L.

**What I recommend:** Option C for now (graphcore/ extraction only). Defer the full `core/` decision to after v1.0 — the markdown + tree extraction already delivered the main composability win. Revisit `core/` if users actually request thinner `go get` payloads.

---

## Verification Summary

| Check | Status | Notes |
|-------|--------|-------|
| `go build ./...` (all 20 modules) | ✅ | Workspace + isolated (markdown/tree) |
| `go test ./...` (all 20 modules) | ✅ | 763+ test functions |
| `go test -race` (nom + tui) | ✅ | Race-detector clean |
| golangci-lint (all modules) | ✅ | 0 issues across 20 modules |
| TODO/FIXME in prod `*.go` | ✅ | Zero |
| Test files ≤350 lines | ✅ | All compliant |
| Root prod lines | ✅ | 1401 (was 1908, −27%) |
| Module count | ✅ | 20 (root + 19 sub-modules) |
| Working tree | ✅ | Clean, up to date with origin |

---

_Generated 2026-06-19 03:15 · HEAD `ff3d9f4` · 20 modules · All tiers complete_
