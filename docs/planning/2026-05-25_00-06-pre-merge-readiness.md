# Pre-Merge Readiness Plan — Branch `modularize/extract-d2-graph` → `master`

**Date:** 2026-05-25
**Branch:** `modularize/extract-d2-graph` (64 commits, 121 files changed)
**Target:** `master` (zero commits ahead — clean rebase target)
**Status:** No merge conflicts. No open PRs/issues. 3 lint errors + 1 broken release pipeline block merge.

---

## Current State

| Check | Status | Detail |
|---|---|---|
| Merge conflicts | ✅ None | `origin/master` has zero commits not in this branch |
| GitHub PRs/Issues | ✅ None open | Clean repo |
| Build (10 modules) | ✅ All pass | `go build ./...` |
| Tests (10 modules) | ✅ All pass | `go test ./...` |
| `go vet` (10 modules) | ✅ All pass | Zero warnings |
| Lint — root | ❌ 1 issue | `color.go:8` goimports |
| Lint — table | ❌ 2 issues | `table/table_test.go:7` goimports |
| Lint — other 8 modules | ✅ Clean | Zero issues |
| Release workflow | ❌ BROKEN | Missing `d2` + `graph` in all 5 loops |
| CI workflow | ✅ Correct | All 10 modules covered |
| Nix flake build | ❌ Fails | Known sandbox issue (CI handles Go) |
| FEATURES.md | ❌ Stale | `NewD2Renderer` known issue already fixed |
| TODO_LIST.md | 26/40 done | 13 open (P6 future), 1 needs decision |

---

## Pareto Breakdown

| Tier | Tasks | Impact | Est. Time |
|---|---|---|---|
| **1% → 51%** | #1-#4 | Blocks merge. CI red, release broken. | ~16 min |
| **4% → 64%** | #5-#9 | Correctness & reliability gate. | ~19 min |
| **20% → 80%** | #10-#13 | Polish, confidence, reviewability. | ~30 min |
| **Post-merge** | #14-#17 | Future work. Does not block merge. | ~23 min |

---

## Execution Plan

### Phase 1: Blockers (1% → 51% impact)

| # | Task | Impact | Effort | Category | Depends |
|---|---|---|---|---|---|
| **1** | Fix 3 goimports errors: `color.go:8`, `table/table_test.go:7` (2 issues) | CRITICAL — lint fails, CI red | 3min | Build Blocker | — |
| **2** | Fix `release.yml` — add `d2` + `graph` to all 5 loops (build, test, govulncheck, benchmark, lint) | CRITICAL — release ships untested d2+graph | 5min | Build Blocker | — |
| **3** | Fix stale `FEATURES.md` — remove `NewD2Renderer` known issue (line 235, line 260) — already fixed in README | HIGH — misleading docs | 3min | Docs Accuracy | — |
| **4** | Verify all 10 modules: build + test + lint clean after fixes | HIGH — gate before commit | 5min | Verification | #1 |

### Phase 2: Correctness (4% → 64% impact)

| # | Task | Impact | Effort | Category | Depends |
|---|---|---|---|---|---|
| **5** | Commit fixes (#1-#3) with detailed message | HIGH — clean commit before rebase | 3min | Git Hygiene | #1-#4 |
| **6** | Verify `CHANGELOG.md` [Unreleased] section is complete and accurate for v0.5.0 | HIGH — release notes source of truth | 5min | Docs Accuracy | — |
| **7** | Verify CI workflow (`ci.yml`) covers all 10 modules correctly (spot-check) | HIGH — CI is the gate | 3min | Verification | — |
| **8** | Run `go vet ./...` across all 10 modules (already verified clean, re-confirm) | MEDIUM — catch subtle issues | 3min | Verification | — |
| **9** | Update `TODO_LIST.md` — mark release.yml as fixed, update FEATURES.md status line | MEDIUM — keep docs honest | 5min | Docs Accuracy | #2, #3 |

### Phase 3: Polish (20% → 80% impact)

| # | Task | Impact | Effort | Category | Depends |
|---|---|---|---|---|---|
| **10** | Clean git history — squash 64 commits into logical groups (~10-15 commits) | MEDIUM — reviewability for merge | 12min | Git Hygiene | #5 |
| **11** | Verify no broken references in `README.md` (imports, API examples, links) | MEDIUM — first thing users see | 5min | Docs Accuracy | — |
| **12** | Check for stale references across all docs (DOMAIN_LANGUAGE, FORMAT_ARCHITECTURE, ADRs, CONTRIBUTING) | MEDIUM — consistency | 8min | Docs Accuracy | — |
| **13** | Final full-suite verification: build + test + lint + vet across all 10 modules | HIGH — final gate before merge | 5min | Verification | #1-#12 |

### Phase 4: Post-Merge / Decisions (not blocking)

| # | Task | Impact | Effort | Category | Depends |
|---|---|---|---|---|---|
| **14** | Decision needed: move `internal/gentest` → `testhelpers/gentest`? (TODO #20) | LOW — future DX | 2min | Decision | merge |
| **15** | Pre-commit hooks: configure BuildFlow ignore rules or document `--no-verify` workaround (TODO #24) | LOW — DX friction | 8min | DX | merge |
| **16** | Tag `v0.5.0` after merge to master (TODO #31) | LOW — post-merge release | 5min | Release | merge |
| **17** | Post to r/golang, submit to Awesome Go (TODO #40) | LOW — community | 8min | Community | #16 |

---

## Summary

| Metric | Value |
|---|---|
| Total tasks | 17 |
| Blocking merge (Phase 1-3) | 13 |
| Post-merge (Phase 4) | 4 |
| Estimated total effort (Phase 1-3) | ~65 min |
| Estimated total effort (all) | ~88 min |

---

## D2 Execution Graph

```d2
direction: down

title: {
  label: Pre-Merge Readiness — Execution Flow
  shape: text
  near: top-center
}

phase1: {
  label: Phase 1: Blockers (1% → 51%)
  shape: rectangle
  style: {
    fill: "#FF6B6B"
    font-color: white
    bold: true
  }

  t1: Fix goimports (3 errors)
  t2: Fix release.yml (add d2+graph)
  t3: Fix stale FEATURES.md

  t1 -> t2 -> t3
}

phase1 -> t4

t4: {
  label: Verify all 10 modules pass
  shape: diamond
  style: {
    fill: "#FFD93D"
    stroke: "#333"
  }
}

t4 -> phase2: {label: "✅ Pass"}
t4 -> phase1: {label: "❌ Fail"; style: {stroke-dash: 3}}

phase2: {
  label: Phase 2: Correctness (4% → 64%)
  shape: rectangle
  style: {
    fill: "#6BCB77"
    font-color: white
    bold: true
  }

  t5: Commit fixes
  t6: Verify CHANGELOG
  t7: Verify CI workflow
  t8: go vet all modules
  t9: Update TODO_LIST

  t5 -> t6 -> t7 -> t8 -> t9
}

phase2 -> phase3

phase3: {
  label: Phase 3: Polish (20% → 80%)
  shape: rectangle
  style: {
    fill: "#4D96FF"
    font-color: white
    bold: true
  }

  t10: Squash git history
  t11: Verify README refs
  t12: Verify all docs
  t13: Final full-suite verify

  t10 -> t11 -> t12 -> t13
}

phase3 -> gate

gate: {
  label: Ready to Merge?
  shape: diamond
  style: {
    fill: "#FFD93D"
    stroke: "#333"
    bold: true
  }
}

gate -> merge: {label: "✅ Ship it"}
gate -> phase3: {label: "❌ Issues found"; style: {stroke-dash: 3}}

merge: {
  label: Rebase → Master
  shape: hexagon
  style: {
    fill: "#9B59B6"
    font-color: white
    bold: true
  }
}

merge -> phase4

phase4: {
  label: Phase 4: Post-Merge
  shape: rectangle
  style: {
    fill: "#95A5A6"
    font-color: white
  }

  t14: Decision: gentest move?
  t15: Fix pre-commit hooks
  t16: Tag v0.5.0
  t17: Community posting

  t14 -> t15 -> t16 -> t17
}
```

---

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Squashing 64 commits reveals forgotten changes | Low | Medium | Review diff before/after squash |
| Release.yml fix misses a loop | Low | High | Compare with ci.yml (which is correct) |
| Nix flake check fails | Certain | Low | Known sandbox issue; CI handles Go |
| goimports errors recur after other edits | Low | Low | Run lint as final gate (#13) |

---

## TODO_LIST.md Cross-Reference

| TODO # | Description | Status | Plan Task |
|---|---|---|---|
| #20 | Move `internal/gentest` to `testhelpers/gentest`? | Needs Decision | #14 |
| #24 | Pre-commit hooks BuildFlow false positives | Open | #15 |
| #31 | Tag next release (v0.5.0) | Open | #16 |
| #32-#40 | Future features (TOML, JSONL, PlantUML, AsciiDoc, etc.) | P6 Future | Not in scope |
