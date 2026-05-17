# Status Report — Modularization Docs Review & Fixes

**Date:** 2026-05-17 03:49
**Branch:** `master`
**Commits since last status:** 3 (97da859, 44c71fb, b149aef)

---

## Executive Summary

Completed a thorough review of the modularization documentation suite (`PROPOSAL.md`, `EXECUTION_PLAN.md`, `DEPENDENCY_GRAPH.md`), fixing 4 critical errors that would cause build failures during execution, resolving an internal contradiction, correcting all LOC counts, and adding a comprehensive "Future Improvements" section with type model analysis. Also fixed a pre-existing `examples/go.mod` bug. **No production code changed.**

The modularization plan is now ready for execution. The execution plan has 5 steps (Step 0 done, Steps 1-4 pending).

---

## A) FULLY DONE

| Item                                    | Commit            | Details                                                                                                                                                    |
| --------------------------------------- | ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Fix `examples/go.mod` missing table dep | 97da859           | Added `table` to require + replace blocks. Was a pre-existing bug.                                                                                         |
| Fix PROPOSAL §3.2 contradiction         | b149aef           | Section described extracting gentest/testutils as modules. Self-review already decided they stay in root. Struck through the section.                      |
| Fix PROPOSAL §4.2 Step 2                | b149aef           | Said "Move graph.go, dot.go, mermaid.go → graph/" but `graph.go` stays in root (core types used by D2). Corrected.                                         |
| Fix EXECUTION_PLAN userjourney_test.go  | b149aef           | Was listed for move to d2/ but tests JSON/CSV/Markdown/YAML/sort, not D2. Removed from move list.                                                          |
| Fix EXECUTION_PLAN empty code block     | b149aef           | Step 1 action 3 had empty `bash` block after edit. Restored D2 test file git mv command.                                                                   |
| Fix EXECUTION_PLAN benchmarks_test.go   | b149aef           | Added missing action: update `benchmarks_test.go` for DOT/Mermaid constructor moves (Step 2).                                                              |
| Fix EXECUTION_PLAN integration imports  | b149aef           | Added missing actions: update integration test import statements (Steps 1 & 2).                                                                            |
| Fix EXECUTION_PLAN README update        | b149aef           | Upgraded from "verify" to explicit actions: change all `output.D2*`/`output.DOT*`/`output.Mermaid*` references.                                            |
| Fix EXECUTION_PLAN Step 0 status        | b149aef           | Marked as DONE (examples/go.mod bug already fixed).                                                                                                        |
| Fix all LOC counts across all 3 docs    | b149aef + 97da859 | D2: 815→833, Graph: 566→568, Core: 681→734, Root: 3,515→3,587, Slimmed: 2,130→2,183                                                                        |
| Fix DEPENDENCY_GRAPH examples row       | 97da859           | Added d2+graph deps, updated DAG proof chains                                                                                                              |
| Remove stale "buggy/Leaky" labels       | b149aef           | examples/go.mod bug is fixed — removed from both DEPENDENCY_GRAPH and PROPOSAL                                                                             |
| Fix PROPOSAL §3.1 examples deps         | b149aef           | Changed from "root" to "root, table" (matching reality)                                                                                                    |
| Fix PROPOSAL §3.5 test dep table        | b149aef           | Removed gentest/testutils as separate modules, updated examples to include d2+graph                                                                        |
| Add PROPOSAL §9 Future Improvements     | b149aef           | 8 items: duplicated tree conversion, EdgeStyle naming, mixin value-embed, RowEdge inconsistency, AddTreeNodes mutation, enum codegen, library replacements |
| Fix PROPOSAL §4.2 Step 2 (graph.go)     | b149aef           | Corrected to say only dot.go + mermaid.go move; graph.go stays                                                                                             |
| Fix PROPOSAL §7 CI/CD                   | b149aef           | Removed gentest/testutils from CI/CD parallel jobs list                                                                                                    |

---

## B) PARTIALLY DONE

| Item                     | Status                                     | What Remains                                               |
| ------------------------ | ------------------------------------------ | ---------------------------------------------------------- |
| Modularization execution | Docs reviewed & fixed, Step 0 done         | Steps 1-4 not started (actual code extraction)             |
| PROPOSAL accuracy        | All known errors fixed                     | May surface more issues during execution                   |
| ADR 002 status           | Code fully implemented & shipped in v0.4.0 | ADR doc still says "PROPOSED" — needs update to "ACCEPTED" |

---

## C) NOT STARTED

| Item                                                                            | Priority | Effort          |
| ------------------------------------------------------------------------------- | -------- | --------------- |
| **Step 1: Extract d2/ module**                                                  | Critical | 20-30 min       |
| **Step 2: Extract graph/ module**                                               | Critical | 20-30 min       |
| **Step 3: Clean up sort dependency**                                            | High     | 5 min           |
| **Step 4: Update documentation (AGENTS.md, FORMAT_ARCHITECTURE.md, README.md)** | High     | 10 min          |
| Write ADR 003 for d2/graph extraction                                           | Medium   | 10 min          |
| Split format.go God file (373 LOC, 43 symbols)                                  | Medium   | 15 min          |
| Add Parse/IsValid/AllowedValues to Shape enum                                   | Medium   | 5 min           |
| Update ADR 002 status from PROPOSED → ACCEPTED                                  | Low      | 1 min           |
| Delete stale PLAN.md                                                            | Low      | 1 min           |
| Decide sort/ module fate (delete / decouple / keep)                             | Medium   | Decision needed |
| Decide cmdguard/ disposition (keep / extract / remove)                          | Low      | Decision needed |
| Fix 85 phantom LSP errors from gitignored go.work                               | Low      | Config change   |
| Marketing (r/golang, Awesome Go, blog post)                                     | Low      | Hours           |

---

## D) TOTALLY FUCKED UP

| Item              | Severity | Details                                                        |
| ----------------- | -------- | -------------------------------------------------------------- |
| Nothing is broken | —        | Build passes, tests pass (90.2% coverage), all 8 modules clean |

Close calls caught in time:

- **userjourney_test.go move to d2/** — would have broken root tests (tests JSON/CSV/Markdown/YAML/sort, not D2). Caught during review.
- **Empty code block in EXECUTION_PLAN** — Step 1 action 3 had no git mv command after edit. Would have caused confusion during execution. Caught and fixed.
- **GraphRendererMixin location wrong in plan** — plan said "split graph.go" but mixin is in dot.go, not graph.go. Would have caused wrong operation during execution. Caught and fixed.

---

## E) WHAT WE SHOULD IMPROVE

### Code Quality

1. **format.go is a God file** (373 LOC, 43 type/func/var/const declarations). Contains 5 misplaced types that should be in their own files:
   - `TreeNode`, `TreeOutputRenderer` → `tree.go`
   - `TableData`, `RowEdge` → `tabledata.go`
   - `Shape` constants + methods → `shape.go`

2. **Duplicated tree→graph conversion** — `D2Diagram.addTreeNodes()` re-implements the tree walk from `AddTreeNodes()` because it uses D2-specific types. Should use `AddTreeNodes()` + conversion.

3. **`EdgeStyle.Style` field** — naming is confusing (`edge.Style.Style`). Should be `LineStyle` or `DashStyle`.

4. **`GraphRendererMixin` is value-embedded** in `DOTRenderer` but has pointer-receiver methods. Copying a `DOTRenderer` value would share slice state silently.

5. **`RowEdge` uses plain `string`** for `From`/`To` while `GraphEdge` uses branded IDs — inconsistent when RowEdge data flows into GraphEdge.

6. **`AddTreeNodes` mutates via pointer params** (`*[]GraphNode`, `*[]GraphEdge`) — unidiomatic Go. Should return new slices.

### Documentation

7. **ADR 002 is stale** — shows "PROPOSED" but code is fully implemented and shipped in v0.4.0.

8. **PLAN.md is stale** — duplicates AGENTS.md content. Should be deleted.

9. **Shape enum is incomplete** — `Shape` type has no `Parse`, `IsValid`, `AllowedValues` methods (unlike every other enum in the codebase).

### Architecture

10. **sort/ ↔ root circular dependency** — sort imports root's `SortBy` type while root's go.mod lists sort as a production dependency. The deeper issue isn't just the go.mod line — it's the bidirectional coupling.

11. **cmdguard/ is completely disconnected** — zero imports either direction. No clear purpose in the repo.

12. **go.work is gitignored** — causes 85 phantom LSP errors for contributors. The replace directive approach works but has a DX cost.

---

## F) Top 25 Things to Do Next

Sorted by **impact × effort** (highest first):

| #   | Task                                                                      | Impact   | Effort   | Category       |
| --- | ------------------------------------------------------------------------- | -------- | -------- | -------------- |
| 1   | **Execute Step 1: Extract d2/ module**                                    | Critical | 30 min   | Modularization |
| 2   | **Execute Step 2: Extract graph/ module**                                 | Critical | 30 min   | Modularization |
| 3   | **Execute Step 3: Clean up sort dependency**                              | High     | 5 min    | Modularization |
| 4   | **Execute Step 4: Update docs**                                           | High     | 10 min   | Modularization |
| 5   | **Split format.go** (move TreeNode/TableData/RowEdge/Shape to own files)  | High     | 15 min   | Code Quality   |
| 6   | **Write ADR 003** (d2/graph extraction decision)                          | Medium   | 10 min   | Architecture   |
| 7   | **Update ADR 002 status** → ACCEPTED                                      | Low      | 1 min    | Documentation  |
| 8   | **Delete stale PLAN.md**                                                  | Low      | 1 min    | Cleanup        |
| 9   | **Add Parse/IsValid/AllowedValues to Shape enum**                         | Medium   | 5 min    | Code Quality   |
| 10  | **Decide sort/ fate** (delete / decouple SortBy / keep)                   | Medium   | Decision | Architecture   |
| 11  | **Deduplicate D2Diagram.addTreeNodes** — use AddTreeNodes + conversion    | Medium   | 10 min   | Code Quality   |
| 12  | **Rename EdgeStyle.Style** → LineStyle or DashStyle                       | Low      | 10 min   | Naming         |
| 13  | **Decide cmdguard/ disposition**                                          | Low      | Decision | Architecture   |
| 14  | **Make GraphRendererMixin pointer-embed** or document requirement         | Low      | 5 min    | Safety         |
| 15  | **Change AddTreeNodes** to return slices instead of mutating              | Low      | 10 min   | API Quality    |
| 16  | **Add branded IDs to RowEdge** (From/To)                                  | Low      | 5 min    | Consistency    |
| 17  | **Remove deprecated FormatCategory** and old methods                      | Low      | 5 min    | Cleanup        |
| 18  | **Remove BrandedID re-export** from ids.go (use go-branded-id directly)   | Low      | 5 min    | Cleanup        |
| 19  | **Evaluate emicklei/dot** as replacement for hand-rolled DOT generation   | Low      | Research | Dependencies   |
| 20  | **Set up CI** for multi-module repo (parallel jobs per module)            | Medium   | 30 min   | DevOps         |
| 21  | **Fix go.work DX** (commit it or improve contributor guide)               | Low      | 5 min    | DX             |
| 22  | **Add JSON/YAML Renderer implementations** (gap in declared capabilities) | Medium   | 1 hr     | Features       |
| 23  | **Submit to r/golang** and Awesome Go                                     | Low      | 30 min   | Marketing      |
| 24  | **Build shape-specific renderers** (ADR 002 Phase 2)                      | Medium   | 2-3 hr   | Features       |
| 25  | **Consider enum code generation** (stringer or custom)                    | Low      | 1 hr     | Code Quality   |

---

## Project Metrics

| Metric              | Value                                                                       |
| ------------------- | --------------------------------------------------------------------------- |
| Root production LOC | 3,584                                                                       |
| Total Go files      | 81                                                                          |
| Test files          | 28                                                                          |
| Root test coverage  | 90.2%                                                                       |
| Modules             | 8 (root + enum + escape + cmdguard + sort + table + integration + examples) |
| Planned modules     | 10 (+ d2, graph)                                                            |
| Open ADRs           | 2 (001 ACCEPTED, 002 PROPOSED but implemented)                              |
| Status reports      | 13 (including this one)                                                     |
| Git commits today   | 4                                                                           |
| Build status        | Passing                                                                     |
| Lint status         | Clean                                                                       |

---

## G) Top #1 Question I Cannot Answer Myself

**What should happen to the `sort/` module?**

The modularization plan says "remove from root's production go.mod" — but that's only the surface fix. The deeper problem is architectural:

- `sort/` imports `output.SortBy` from root (sort depends on root)
- Root's `go.mod` lists `sort` as a production dependency (root depends on sort)
- This is a **circular dependency** that Go resolves via replace directives, but it's conceptually wrong

Three options:

1. **Delete sort/ entirely** — it's deprecated, stdlib `slices.SortStableFunc` + `cmp.Compare` replaces it (Go 1.21+). But `userjourney_test.go` uses it as an integration test for deprecated functionality.
2. **Decouple: move `SortBy` to sort/** — sort becomes self-contained, root no longer depends on it. But `SortBy` is in root's public API and used by cmdguard tests.
3. **Keep as-is** — accept the circular dependency, just remove from root's production go.mod.

The decision affects whether `sort/` has a future in this repo at all. I can execute any of the three options but need your call on the strategic direction.

---

_Generated with Crush_
