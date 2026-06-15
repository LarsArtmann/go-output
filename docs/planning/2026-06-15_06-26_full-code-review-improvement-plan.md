# Full Code Review — Comprehensive Improvement Plan

**Date:** 2026-06-15
**Reviewer:** Senior Software Architect pass (full-code-review skill)
**Scope:** All 16 modules · 214 Go files · ~9.9k prod LOC

---

## Executive Summary

The go-output codebase is **architecturally excellent** — correct module boundaries, zero cycles, radical dependency isolation, type-safe enums, branded phantom IDs, and a self-registering registry-dispatch pattern. This review visited every production file across all modules, fixed the clear high-value issues on the spot, and prioritizes the remainder below.

**Fixes applied this session (build + tests green):**
1. `.golangci.yml` depguard config portability — eliminated **24 false-positive** lint errors across 8 modules
2. `nom/inline_renderer.go` — added best-effort write helpers (7 errcheck), flattened nested status block (1 nestif)
3. `nom/tree_render.go` — perfsprint fix
4. `nom/activity_display.go` + `nom/tree.go` — embedded field ordering (2 embeddedstructfieldcheck)
5. `nom/subscriber_handlers.go`, `tui/view.go`, `tui/reporter.go` — whitespace (7 wsl_v5)
6. `nom/activity_display_test.go` — `copy` builtin shadow (predeclared)
7. `graph/registry_test.go` — refactored DOT/Mermaid duplicate tests into table-driven helper (dupl)
8. `tui/colors.go` — removed dead `colorCyan` global

**Result:** production lint issues dropped from ~25 → 10 (remaining 10 are `tui/colors.go` globals). All modules build and test clean.

---

## Pareto Breakdown

### 🔴 Tier 0 — DONE (1% → 51% of value)

| Task | Impact | Status |
| --- | --- | --- |
| Fix depguard config (24 false positives) | Restores lint signal across all modules | ✅ Done |
| nom production lint (errcheck/nestif/perfsprint/embedded) | Cleaner hot-path code | ✅ Done |
| Remove dead `colorCyan` | Dead code removal | ✅ Done |

### 🟠 Tier 1 — 4% → 64% (high value, low effort)

| # | Task | Files | Effort | Value |
| - | --- | --- | --- | --- |
| 1 | `tui/colors.go`: convert 10 mutable globals → immutable style struct | colors.go, view.go, summary.go | Medium | Medium |
| 2 | `tui/model_test.go`: replace deprecated `EnsureBuild()` → `GetRootNodes()` (3 SA1019) | model_test.go | Low | Low |
| 3 | `tui/event_sequence_test.go`: guard 3 type assertions (forcetypeassert) | event_sequence_test.go | Low | Low |

### 🟡 Tier 2 — 20% → 80% (medium value)

| # | Task | Files | Effort | Value |
| - | --- | --- | --- | --- |
| 4 | Decompose `nom/` into `internal/` sub-packages (tree, cache, render, subscriber) keeping thin public API | nom/ (35 files) | High | Medium (locality, navigability) |
| 5 | Test errcheck sweep (31 unchecked `AddActivity`/`Record`/`OnEvent` returns across nom/tui/integration tests) | many _test.go | Medium | Low |
| 6 | err113 test sweep (7 dynamic `errors.New` → wrapped static errors) | tests | Low | Low |

### 🔵 Tier 3 — Polish (long-term health)

| # | Task | Notes |
| - | --- | --- |
| 7 | `TableData` invariant enforcement (unexport fields, validated setters) — **post-v1, ADR 006 conflict** | See improve-codebase-architecture report |
| 8 | Unify `Marshaler` → `Renderer` terminology in registry | Post-v1 API change |
| 9 | Collapse `InlineRenderer` 8-method interface behind smaller render seam | See deepening report |
| 10 | `DOMAIN_LANGUAGE.md` staleness: says `GraphRendererMixin`, code says `GraphRendererState`; missing nom/tui bounded contexts | docs-freshness |

---

## Execution Graph

```d2
direction: down

done: "Tier 0 — DONE" { style.fill: "#3fb95030"; shape: package }
t1: "Tier 1 — tui cleanup" { style.fill: "#d2992230"; shape: package }
t2: "Tier 2 — nom decompose + test sweep" { style.fill: "#58a6ff30"; shape: package }
t3: "Tier 3 — post-v1 polish" { style.fill: "#8b949e30"; shape: package }

done -> t1: next
t1 -> t2
t2 -> t3: deferred (v1 freeze)

t1: |go
  colors → immutable struct
  EnsureBuild → GetRootNodes
  guard type assertions
|

t2: |go
  nom/ internal sub-packages
  test errcheck + err113 sweep
|

t3: |go
  TableData invariants
  Marshaler → Renderer
  InlineRenderer deepening
  DOMAIN_LANGUAGE refresh
|
```

---

## Architect Checklist Reflection

| Concern | Assessment |
| --- | --- |
| **Data flow** | Clean — registry dispatch + renderer interfaces, single direction |
| **Impossible states** | Mostly good (branded IDs, enums). Gap: `TableData` allows column mismatch (post-hoc `Validate`) |
| **Composed architecture** | Excellent — interfaces, composition over inheritance, shared state via embedding |
| **Generics** | Used appropriately (`formatRegistry[T]`, `enum.Parse[T]`, `BrandedID[T]`) |
| **Booleans → enums** | Good — `ColorMode`, `WorkflowState`, `DisplayMode`, `ActivityStatus` are enums |
| **uints** | Not applicable (string-heavy library) |
| **Files > 350 lines** | Only `tui/view.go` (375 lines) — borderline, acceptable |
| **Split brains** | Minor: `Marshaler` vs `Renderer`; `TableData` fields+setters+getters |
| **Errors centralized** | `pkg/errors` referenced in allow-list; sentinel errors in defining modules |
| **External tools wrapped** | lipgloss/bubbletea/yaml/toml isolated in modules (adapters) |
| **Naming** | Excellent — zero Manager/Handler/Impl (see naming-review-2026-06-15.md) |
| **Duplication** | Minimal — art-dupl "Excellent health", 1 clone fixed |
| **Test coverage** | 84-100% across modules; integration + round-trip + user-journey tests |

---

## Verdict

The codebase does not need a rescue — it needs **polish**. The architecture is sound, the types are honest, the boundaries are correct. The remaining work is: (1) finish the tui cleanup, (2) decide whether `nom/` benefits from internal decomposition, and (3) plan the post-v1 deepening (TableData invariants, terminology unification). None of this blocks usage or correctness today.
