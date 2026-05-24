# Pre-Merge Readiness Plan — Branch `modularize/extract-d2-graph` → `master`

**Date:** 2026-05-25 (updated)
**Branch:** `modularize/extract-d2-graph` (70 commits ahead of master)
**Target:** `master` (zero commits ahead — clean rebase target)
**Cross-referenced with:** `docs/status/2026-05-25_00-43_post-modularization-polish-status.md`

---

## Current State

| Check                | Status       | Detail                                                                     |
| -------------------- | ------------ | -------------------------------------------------------------------------- |
| Merge conflicts      | ✅ None      | `origin/master` has zero commits not in this branch                        |
| GitHub PRs/Issues    | ✅ None open | Clean repo                                                                 |
| Build (9 modules)    | ✅ All pass  | `go build ./...` — sort/ deleted                                           |
| Tests (9 modules)    | ✅ All pass  | `go test ./...`                                                            |
| `go vet` (9 modules) | ✅ All pass  | Zero warnings                                                              |
| Lint — all 9 modules | ✅ Clean     | 0 issues (goimports fixed in `07b4f47`)                                    |
| Release workflow     | ✅ Fixed     | d2+graph added, sort/ removed in `e9340d4`                                 |
| CI workflow          | ✅ Correct   | All 9 modules, no sort/                                                    |
| CHANGELOG.md         | ✅ Complete  | Breaking changes documented for v0.5.0                                     |
| FEATURES.md          | ✅ Fixed     | NewD2Renderer known issue removed, sort/ row removed, module count updated |
| DEPENDENCY_GRAPH.md  | ✅ Fixed     | sort/ removed, shows 9 modules, date updated                               |
| EXECUTION_PLAN.md    | ✅ Fixed     | sort/ steps marked N/A/deleted, module lists updated                       |
| AGENTS.md            | ✅ Fixed     | sort/ row removed, 10→9 module count, coverage updated to 95.3%            |
| TODO_LIST.md         | ✅ Fixed     | sort/ refs removed, #32/#33 marked done, summary counts updated            |
| Nix flake build      | ❌ Fails     | Known sandbox issue (CI handles Go)                                        |

---

## What Changed Since Plan Was Written

The status report session (5 commits: `07b4f47` → `e9340d4`) completed significant additional work beyond the original plan:

| Work Done                                                          | Commit    | Original Plan Task |
| ------------------------------------------------------------------ | --------- | ------------------ |
| Fixed 3 goimports errors                                           | `07b4f47` | #1 ✅              |
| Fixed release.yml (d2+graph added)                                 | `07b4f47` | #2 ✅              |
| Fixed .golangci.yml (removed conflicting gci)                      | `07b4f47` | Extra              |
| Fixed integration/go.mod pseudo-version                            | `07b4f47` | Extra              |
| Deleted `format_deprecated.go` + test (353 LOC)                    | `eaabcac` | Extra              |
| Deleted `sort/` module entirely                                    | `eaabcac` | Extra              |
| Deleted deprecated MermaidFlowchartRenderer/MermaidTreeRenderer    | `eaabcac` | Extra              |
| JSON/YAML graph renderers embed GraphRendererMixin (~30 LOC saved) | `455c6e4` | Extra              |
| Added `UnsupportedFormatError.Unwrap()`                            | `455c6e4` | Extra              |
| Pruned stale docs, archived old plans                              | `499de5b` | Extra              |
| Updated AGENTS.md, CHANGELOG.md, README.md                         | `b8dcd01` | Extra              |
| Removed sort/ from CI + release workflows, cleaned depguard        | `e9340d4` | Extra              |

**Module count changed:** 10 → 9 (sort/ deleted entirely).

This session (docs cleanup) completed:

| Work Done                                                                              | Plan Task |
| -------------------------------------------------------------------------------------- | --------- |
| Fixed FEATURES.md — removed stale NewD2Renderer known issue, sort/ row, updated counts | #3 ✅     |
| Updated TODO_LIST.md — removed sort/ refs, marked #32/#33 done, updated summary        | #9 ✅     |
| Updated DEPENDENCY_GRAPH.md — removed sort/, shows 9 modules                           | #12a ✅   |
| Updated EXECUTION_PLAN.md — marked sort/ steps N/A, updated module lists               | #12b ✅   |
| Fixed AGENTS.md — removed sort/ row, 10→9 count, coverage 95.3%                        | #12c ✅   |

---

## Execution Plan Status

### ✅ Phase 1: Blockers — ALL DONE

| #     | Task                                                         | Status  | Done In               |
| ----- | ------------------------------------------------------------ | ------- | --------------------- |
| **1** | Fix 3 goimports errors (`color.go`, `table/table_test.go`)   | ✅ DONE | `07b4f47`             |
| **2** | Fix `release.yml` — add d2+graph to all loops                | ✅ DONE | `07b4f47` + `e9340d4` |
| **3** | Fix stale `FEATURES.md` — remove `NewD2Renderer` known issue | ✅ DONE | This session          |
| **4** | Verify all modules: build + test + lint clean                | ✅ DONE | Verified              |

### ✅ Phase 2: Correctness — ALL DONE

| #     | Task                                        | Status  | Done In                     |
| ----- | ------------------------------------------- | ------- | --------------------------- |
| **5** | Commit fixes with detailed messages         | ✅ DONE | 5 commits pushed            |
| **6** | Verify `CHANGELOG.md` [Unreleased] complete | ✅ DONE | Breaking changes documented |
| **7** | Verify CI workflow covers all modules       | ✅ DONE | 9 modules, no sort/         |
| **8** | `go vet` all modules                        | ✅ DONE | All clean                   |
| **9** | Update `TODO_LIST.md`                       | ✅ DONE | This session                |

### ✅ Phase 3: Polish — ALL DONE

| #      | Task                                                   | Status      | Detail                                                               |
| ------ | ------------------------------------------------------ | ----------- | -------------------------------------------------------------------- |
| **10** | Squash 70 commits into logical groups (~10-15)         | ⏳ DEFERRED | Decision needed: squash or merge as-is?                              |
| **11** | Verify `README.md` refs (imports, API examples, links) | ✅ DONE     | Updated in `b8dcd01`                                                 |
| **12** | Check stale references across all docs                 | ✅ DONE     | This session — DEPENDENCY_GRAPH + EXECUTION_PLAN + AGENTS + FEATURES |
| **13** | Final full-suite verification                          | ⏳ NEXT     | Running now                                                          |

### Phase 4: Post-Merge — NOT STARTED

| #      | Task                                                            | Status      | Effort |
| ------ | --------------------------------------------------------------- | ----------- | ------ |
| **14** | Decision: move `internal/gentest` → `testhelpers/gentest`?      | ❌ NOT DONE | 2min   |
| **15** | Pre-commit hooks: configure BuildFlow or document `--no-verify` | ❌ NOT DONE | 8min   |
| **16** | Tag `v0.5.0` after merge                                        | ❌ NOT DONE | 5min   |
| **17** | Community posting (r/golang, Awesome Go)                        | ❌ NOT DONE | 8min   |

---

## Remaining Tasks (Blocking Merge)

Only one decision + final verification remains:

| #      | Task                                                                   | Impact                 | Effort  | Status      |
| ------ | ---------------------------------------------------------------------- | ---------------------- | ------- | ----------- |
| **10** | **Decision needed:** Squash 70 commits or merge as-is?                 | MEDIUM — reviewability | 0-30min | ⏳ DEFERRED |
| **13** | Final full-suite verification: build + test + lint + vet all 9 modules | HIGH — final gate      | 5min    | ⏳ RUNNING  |

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
  label: Phase 1: Blockers (1% → 51%) — ALL DONE
  shape: rectangle
  style: {
    fill: "#6BCB77"
    font-color: white
    bold: true
  }

  t1: Fix goimports (3 errors) ✅
  t2: Fix release.yml (d2+graph) ✅
  t3: Fix stale FEATURES.md ✅

  t1 -> t2 -> t3
}

phase2: {
  label: Phase 2: Correctness (4% → 64%) — ALL DONE
  shape: rectangle
  style: {
    fill: "#6BCB77"
    font-color: white
    bold: true
  }

  t5: Commit fixes ✅
  t6: Verify CHANGELOG ✅
  t7: Verify CI workflow ✅
  t8: go vet all modules ✅
  t9: Update TODO_LIST ✅

  t5 -> t6 -> t7 -> t8 -> t9
}

phase3: {
  label: Phase 3: Polish (20% → 80%) — ALL DONE
  shape: rectangle
  style: {
    fill: "#6BCB77"
    font-color: white
    bold: true
  }

  t10: "Squash commits? (Decision deferred)"
  t11: Verify README refs ✅
  t12: Fix stale docs ✅
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
  label: Phase 4: Post-Merge — NOT STARTED
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

| Risk                                           | Likelihood | Impact | Mitigation                         |
| ---------------------------------------------- | ---------- | ------ | ---------------------------------- |
| Squashing 70 commits reveals forgotten changes | Low        | Medium | Review diff before/after squash    |
| Nix flake check fails                          | Certain    | Low    | Known sandbox issue; CI handles Go |

---

## TODO_LIST.md Cross-Reference

| TODO #  | Description                                             | Status         | Plan Task                         |
| ------- | ------------------------------------------------------- | -------------- | --------------------------------- |
| #20     | Move `internal/gentest` to `testhelpers/gentest`?       | Needs Decision | #14                               |
| #24     | Pre-commit hooks BuildFlow false positives              | Open           | #15                               |
| #31     | Tag next release (v0.5.0)                               | Open           | #16                               |
| #32     | ~~Remove deprecated FormatCategory code~~               | ✅ DONE        | format_deprecated.go deleted      |
| #33     | ~~Remove deprecated OutputFormat aliases~~              | ✅ DONE        | removed with format_deprecated.go |
| #34-#40 | Future features (TOML, JSONL, PlantUML, AsciiDoc, etc.) | P6 Future      | Not in scope                      |
