# Status Report — Post-Modularization Polish & Dead Code Removal

**Date:** 2026-05-25 00:43
**Branch:** `modularize/extract-d2-graph`
**Author:** Crush (AI Assistant)
**Commits this session:** 5 (already pushed) + 1 pending

---

## Current System State

### Module Verification Matrix

| Module          | Build | Test | Lint     | Vet | Coverage |
| --------------- | ----- | ---- | -------- | --- | -------- |
| root (`output`) | ✅    | ✅   | 0 issues | ✅  | 95.3%    |
| d2              | ✅    | ✅   | 0 issues | ✅  | 100.0%   |
| graph           | ✅    | ✅   | 0 issues | ✅  | 97.6%    |
| enum            | ✅    | ✅   | 0 issues | ✅  | 100.0%   |
| escape          | ✅    | ✅   | 0 issues | ✅  | 100.0%   |
| testhelpers     | ✅    | ✅   | 0 issues | ✅  | 93.8%    |
| table           | ✅    | ✅   | 0 issues | ✅  | 100.0%   |
| integration     | ✅    | ✅   | —        | ✅  | 82.8%    |
| examples        | ✅    | —    | —        | ✅  | —        |

**DAG verification:** Root has ZERO sub-module imports ✅
**Module count:** 9 (down from 10 — sort/ removed)
**Root production LOC:** ~2,874 (down from ~4,345 pre-modularization)
**Total Go files:** 96

---

## a) FULLY DONE

### Phase 1: Unblock CI (1% → 51%)

| #   | Task                                                                                         | Commit    |
| --- | -------------------------------------------------------------------------------------------- | --------- |
| 1   | Fix 3 goimports failures (color.go, table.go, table_test.go) — 3-group import ordering       | `07b4f47` |
| 2   | Fix `.golangci.yml`: remove conflicting gci formatter, add testhelpers to replace-allow-list | `07b4f47` |
| 3   | Fix `release.yml`: add d2+graph to all 4 module loops                                        | `07b4f47` |
| 4   | Fix `integration/go.mod`: normalize d2 pseudo-version                                        | `07b4f47` |

### Phase 2: Remove Dead Code (4% → 64%)

| #   | Task                                                                                        | LOC Removed | Commit    |
| --- | ------------------------------------------------------------------------------------------- | ----------- | --------- |
| 5   | Delete `format_deprecated.go` + `format_deprecated_test.go`                                 | 353         | `eaabcac` |
| 6   | Delete `sort/` module entirely (compare.go, compare_test.go, go.mod)                        | ~60         | `eaabcac` |
| 7   | Delete deprecated `MermaidFlowchartRenderer`/`MermaidTreeRenderer` wrappers                 | ~13         | `eaabcac` |
| 8   | Rename `output_test_helpers.go` → `_test.go` (fix production test dep)                      | —           | `eaabcac` |
| 9   | Delete `TestContainsString` (tests stdlib)                                                  | ~27         | `eaabcac` |
| 10  | Update `.golangci.yml`: remove sort/ exclusion rules, depguard entries, sort_test exclusion | —           | `eaabcac` |

### Phase 3: Quality Improvements (20% → 80%)

| #   | Task                                                                                      | Commit    |
| --- | ----------------------------------------------------------------------------------------- | --------- |
| 11  | JSON/YAML graph renderers now embed `GraphRendererMixin` (eliminates ~30 LOC duplication) | `455c6e4` |
| 12  | Add `UnsupportedFormatError.Unwrap()` for error chain support                             | `455c6e4` |

### Phase 4: Documentation & Cleanup

| #   | Task                                                               | Commit    |
| --- | ------------------------------------------------------------------ | --------- |
| 13  | Prune 4 stale status docs (keep latest)                            | `499de5b` |
| 14  | Archive 3 stale planning docs                                      | `499de5b` |
| 15  | Update AGENTS.md (9 modules, sort/ removed, coverage table)        | `b8dcd01` |
| 16  | Update CHANGELOG.md (Removed section with all breaking changes)    | `b8dcd01` |
| 17  | Update README.md (MermaidFlowchartRenderer → MermaidFromTableData) | `b8dcd01` |

### This commit (post-status fixes)

| #   | Task                                                                                               |
| --- | -------------------------------------------------------------------------------------------------- |
| 18  | Remove `sort` from ci.yml module loops (4 occurrences)                                             |
| 19  | Remove `sort` from release.yml module loops (was missed — had d2/graph added but sort not removed) |
| 20  | Remove deleted `internal/testutils` from `.golangci.yml` depguard allow lists (2 occurrences)      |

---

## b) PARTIALLY DONE

| Task                          | Status      | What's left                                                                             |
| ----------------------------- | ----------- | --------------------------------------------------------------------------------------- |
| Update `FEATURES.md`          | Partial     | Known-issue about `NewD2Renderer` is stale (README is already fixed). Needs full audit. |
| Update `flake.nix`            | Not started | Still only builds root module, doesn't verify all 9 modules                             |
| Update `docs/modularization/` | Partial     | DEPENDENCY_GRAPH.md still shows sort/, 10 modules. EXECUTION_PLAN.md shows sort/ steps. |

---

## c) NOT STARTED

| #   | Task                                                                                | Impact                   | Effort | Why deferred                          |
| --- | ----------------------------------------------------------------------------------- | ------------------------ | ------ | ------------------------------------- |
| 1   | Standardize error types (sentinel → rich struct)                                    | API consistency          | 60min  | Medium impact, many callers to update |
| 2   | Normalize renderer naming (`MarkdownTable` → `MarkdownTableRenderer`)               | Naming consistency       | 45min  | Breaking change, cosmetic             |
| 3   | Standardize `FromTableData` naming pattern                                          | Naming consistency       | 30min  | Breaking change, cosmetic             |
| 4   | Add streaming + XML benchmarks                                                      | Perf regression          | 30min  | Nice-to-have                          |
| 5   | Deprecate or remove registry (zero production usage)                                | Removes ~280 LOC         | 30min  | Low urgency                           |
| 6   | Update `.pre-commit-config.yaml` (stale versions, no sub-module coverage)           | Non-nix users            | 30min  | Nix users unaffected                  |
| 7   | Add `EdgeStyle.Style` as defined type (currently bare `string`)                     | Type safety              | 15min  | Minor                                 |
| 8   | Consistent re-export pattern for graph/ branded IDs                                 | Cross-module consistency | 15min  | Low impact                            |
| 9   | Refactor D2's `addTreeNodes` to use shared `output.AddTreeNodes`                    | Dedup                    | 30min  | D2 has different traversal needs      |
| 10  | Extract render helpers from `integration/integration_test.go` (near 350-line limit) | File size                | 15min  | Not yet over limit                    |
| 11  | Add generic test helpers to `testhelpers/`                                          | Dedup across modules     | 30min  | Known tradeoff from modularization    |
| 12  | Update `docs/modularization/DEPENDENCY_GRAPH.md` (still shows sort/)                | Doc accuracy             | 10min  |                                       |
| 13  | Full audit & update `FEATURES.md`                                                   | Doc accuracy             | 30min  |                                       |
| 14  | Update `flake.nix` to build/test all 9 modules                                      | CI coverage              | 45min  |                                       |
| 15  | Squash 69 commits into logical groups for clean merge                               | Reviewability            | 30min  |                                       |

---

## d) TOTALLY FUCKED UP

| #   | What                                                             | Status              | Fix                                                                                                                 |
| --- | ---------------------------------------------------------------- | ------------------- | ------------------------------------------------------------------------------------------------------------------- |
| 1   | `.golangci.yml` gomodguard name                                  | Fixed               | Initially changed `gomodguard_v2` → `gomodguard` (wrong — v2 IS current). Reverted.                                 |
| 2   | Removed exclusion rules from `.golangci.yml`                     | Fixed               | Removed `_test.go$` goconst/paralleltest/thelper exclusions that were needed. Restored.                             |
| 3   | Added gci formatter config that conflicted with goimports        | Fixed               | `gci` used `localmodule` (only current module) while `goimports` used `local-prefixes` (all siblings). Removed gci. |
| 4   | Forgot to remove `sort` from ci.yml and release.yml module loops | Fixed (this commit) | Removed sort from all 8 loop occurrences.                                                                           |
| 5   | Forgot to remove `internal/testutils` from depguard allow lists  | Fixed (this commit) | Module deleted months ago but depguard still listed it.                                                             |

---

## e) WHAT WE SHOULD IMPROVE

1. **Stop guessing at edits** — The golangci.yml mistakes came from guessing instead of reading carefully. Always read full context before editing config files.
2. **Check ALL files that reference deleted modules** — When removing sort/, should have grepped ALL files (ci.yml, release.yml, depguard, AGENTS.md, etc.) in a single pass.
3. **Run the formatter instead of hand-editing imports** — `golangci-lint fmt` is the right tool for import grouping, not manual blank-line insertion.
4. **Commit more atomically** — The Phase 2 commit bundled too many changes. Smaller commits are easier to revert.
5. **Pre-commit hook NOTE false positives** — The BuildFlow `todo-check` step treats `// Note:` as a TODO. Should be configured to only match `// TODO:`, `// FIXME:`, etc.

---

## f) Top 25 Things to Do Next

Sorted by impact × effort (highest first):

| #   | Task                                                                            | Impact            | Effort | Blocks merge? |
| --- | ------------------------------------------------------------------------------- | ----------------- | ------ | ------------- |
| 1   | Squash commits into logical groups (~10-15) for clean merge                     | Reviewability     | 30min  | Yes           |
| 2   | Update `docs/modularization/DEPENDENCY_GRAPH.md` (remove sort/, show 9 modules) | Doc accuracy      | 10min  | Yes           |
| 3   | Update `docs/modularization/EXECUTION_PLAN.md` (mark sort/ steps as N/A)        | Doc accuracy      | 10min  | Yes           |
| 4   | Audit & update `FEATURES.md` (remove stale known-issue, update module count)    | Doc accuracy      | 30min  | Yes           |
| 5   | Verify `CHANGELOG.md` [Unreleased] is complete for v0.5.0                       | Release           | 10min  | Yes           |
| 6   | Final full-suite verify: build + test + lint + vet all 9 modules                | Quality gate      | 10min  | Yes           |
| 7   | Rebase onto master (zero conflicts expected)                                    | Merge prep        | 5min   | Yes           |
| 8   | Tag `v0.5.0` after merge                                                        | Release           | 5min   | Post-merge    |
| 9   | Update `flake.nix` to build/test all 9 modules                                  | CI coverage       | 45min  | No            |
| 10  | Standardize error types (sentinel → rich struct)                                | API consistency   | 60min  | No            |
| 11  | Deprecate registry (zero production usage)                                      | Removes ~280 LOC  | 30min  | No            |
| 12  | Normalize renderer naming (`MarkdownTable` → `MarkdownTableRenderer`)           | Naming            | 45min  | No            |
| 13  | Standardize `FromTableData` naming                                              | Naming            | 30min  | No            |
| 14  | Add streaming + XML benchmarks                                                  | Perf coverage     | 30min  | No            |
| 15  | Consistent re-export pattern for graph/ branded IDs                             | Consistency       | 15min  | No            |
| 16  | Refactor D2's `addTreeNodes` to use shared helper                               | Dedup             | 30min  | No            |
| 17  | Add `EdgeStyle.Style` as defined type                                           | Type safety       | 15min  | No            |
| 18  | Update `.pre-commit-config.yaml`                                                | Non-nix users     | 30min  | No            |
| 19  | Extract render helpers from integration_test.go                                 | File size         | 15min  | No            |
| 20  | Add generic test helpers to testhelpers/                                        | Dedup             | 30min  | No            |
| 21  | Add large-dataset integration test for streaming                                | Streaming quality | 20min  | No            |
| 22  | Make `GraphRendererMixin` exported as value type (currently struct)             | API review        | 15min  | No            |
| 23  | Fix BuildFlow `todo-check` to not match `// Note:`                              | DX                | 10min  | No            |
| 24  | Add `go.work` to `.gitignore` if not already                                    | Config            | 2min   | No            |
| 25  | Post to r/golang, submit to Awesome Go                                          | Community         | 15min  | Post-release  |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we squash the 69 commits into logical groups before merging to master, or merge as-is?**

Arguments for squashing: Clean git history, easier bisect, professional appearance.
Arguments against: Loses granular change history, extra effort, risk of mistakes during interactive rebase.

The pre-merge readiness plan suggests squashing to ~10-15 logical commits. I cannot make this decision — it's a team/project preference.
