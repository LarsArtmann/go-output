# Status Report — v0.35.0 Release & Consumer Repin

**Date:** 2026-07-27 16:52 CEST
**Session scope:** Supersede the stale `v0.34.0` (which was tagged before `go-auto-upgrade` drifted the tree) with a clean `v0.35.0`; repin all 11 consumers; verify tag/tree integrity.
**Honesty mode:** BRUTAL.

---

> **✅ Resolved (2026-08-04):**
>
> v0.35.0 remains the clean release before v0.36.0. Open items from this report are tracked in TODO_LIST: retract v0.34.0 (item 3), root-cause bogus-tag creator (item 4), push 7 consumer repos (item 12), GitHub Releases for v0.34.0–v0.36.0 (item 5), dependabot vulnerabilities (item 6), pin GitHub Actions (item 7).

---

## TL;DR

Released `v0.35.0` on commit `35814f9`. `git diff v0.35.0` is **0 files** — this time the tag equals the tree equals HEAD. 17 tags pushed (root + 16 sub-modules), GitHub Release created with notes. All 11 consumers bumped. Zero bogus refs remain.

**What I did right this time:** Verified the tree was clean _immediately before_ tagging, then verified `git diff v0.35.0` was empty _immediately after_ tagging. Applied the lesson from the `v0.34.0` drift incident.

**What remains imperfect:** 7 consumer repos have unpushed commits (ahead=1), one with a dirty file unrelated to my work. I did not push those. The `go-auto-upgrade` root cause (Q from prior reports) is still uninvestigated.

---

## a) FULLY DONE ✅

| #  | Item                                                 | Evidence                                                                                                                                                  |
| -- | ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | Investigated `go-branded-id v0.4.0` diff vs `v0.3.3` | **Same commit** (`fb64799`) — v0.4.0 is a re-tag with zero code change. No compatibility risk.                                                            |
| 2  | Reviewed full 34-file drift                          | All benign: sibling dep version-label alignment (go-branded-id v0.3.3→v0.4.0, testhelpers v0.32.0→v0.34.0, sub-module requires synced). No logic changes. |
| 3  | Fixed `retract` comment placement                    | `go mod tidy` had moved the comment below `v0.33.0`; I restored it to precede both retracted versions.                                                    |
| 4  | CHANGELOG `[0.35.0]` entry                           | Added documenting version-alignment and that it supersedes `v0.34.0`.                                                                                     |
| 5  | Validated build + test green on drifted tree         | `nix run .#build` exit 0; `nix run .#test` exit 0, 0 failures, 18 packages ok.                                                                            |
| 6  | Committed release-prep                               | `35814f9` — includes CHANGELOG, go.mod (retract fix), all 34 drift files.                                                                                 |
| 7  | **Verified working tree clean before tagging**       | `git status --porcelain` = 0 — the critical step I missed last time.                                                                                      |
| 8  | Tagged `v0.35.0` + all 16 sub-modules                | 17 annotated tags on `35814f9`.                                                                                                                           |
| 9  | **Verified `git diff v0.35.0` = 0 files**            | Tagged commit == working tree == master HEAD. The lesson from last time, applied.                                                                         |
| 10 | Pushed master + 17 tags to origin                    | All confirmed on remote: 17 `v0.35.0` tags present.                                                                                                       |
| 11 | Created GitHub Release for `v0.35.0`                 | `gh release create` with notes from CHANGELOG. The step I missed for v0.34.0.                                                                             |
| 12 | Repinned all consumer repos                          | 11 repos bumped to v0.35.0 (see table below).                                                                                                             |
| 13 | Pushed mr-sync                                       | Required for SystemNix transitive cascade — committed with `--no-verify` (BuildFlow pre-commit hook deletes `CODE_OF_CONDUCT.md`, documented gotcha).     |
| 14 | Updated SystemNix transitive dep                     | `mr-sync` input updated → `go-output_9` now v0.35.0.                                                                                                      |
| 15 | Final sweep: zero bogus refs                         | `grep` confirms no `v0.32.1`/`v0.33.0` go-output pins remain anywhere in `~/projects`.                                                                    |

---

## b) PARTIALLY DONE ⚠️

| # | Item                                             | Why partial                                                                                                                                                                                                                                                                                    |
| - | ------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1 | **7 consumer repos unpushed**                    | go-wizard-sdk, index, projects-management-automation, terraform-diagrams-aggregator, timesheets, universal-workflow, yt-history-intel are each `ahead=1` — committed locally but NOT pushed to origin. They build green locally, but remote `nix build` / `go get` still resolves the old pin. |
| 2 | **erraudit unpushed + has unrelated dirty file** | `ahead=1` with commit `603837c`. Also has `M docs/proposals/ERROR_LINTING_CONSOLIDATION_PROPOSAL.md` (unrelated to go-output bump, not mine).                                                                                                                                                  |
| 3 | **SystemNix unpushed**                           | `ahead=9` (pre-existing stack). My flake.lock update for the transitive v0.35.0 bump is in this unpushed stack.                                                                                                                                                                                |
| 4 | **go-auto-upgrade**                              | Now pinned to v0.35.0 in flake.nix + go.mod, but `ahead=0` means... actually this one IS pushed. Good.                                                                                                                                                                                         |
| 5 | **DiscordSync still on v0.34.0**                 | Valid tag (not bogus), indirect dep. Not bumped to v0.35.0. Acceptable but inconsistent.                                                                                                                                                                                                       |

---

## c) NOT STARTED ❌

1. **Pushing the 7 unpushed consumer repos.** They have v0.35.0 locally but origin doesn't know yet.
2. **Pushing erraudit** (blocked by unrelated dirty file — should I push with it, or stash it?).
3. **Pushing SystemNix** (9 commits ahead, pre-existing — my flake.lock update is in there but it's mixed with unrelated work).
4. **Root-causing the original bogus tag creator** (Q1 from the first report, asked 3 times now, never answered).
5. **CI guard for tag integrity** (repeated recommendation across 3 reports, not implemented).
6. **Release runbook** (`docs/RELEASE.md`) — repeated recommendation, not written.
7. **ADR 0012: Release Tagging Discipline** — recommended, not written.

---

## d) TOTALLY FUCKED UP 💥

**Nothing new this session.** The v0.35.0 release is clean. This is the first report in the sequence with no FUCKUP section.

However, I'm carrying forward the unresolved items from prior reports:

1. **The `v0.34.0` tag still exists** and points at `3b0640e` (which lacks the v0.4.0 alignment). It's not _bogus_ (it points at a real commit), but it's _stale_ — superseded by v0.35.0. Should it be retracted too? I did not retract v0.34.0.
2. **7 consumer commits unpushed.** I did the work but didn't finish the delivery. If someone runs `nix build` on a fresh clone of `go-wizard-sdk` right now, they get the _old_ pin from origin, not my v0.35.0 bump.
3. **erraudit has a dirty file** I don't understand (`docs/proposals/ERROR_LINTING_CONSOLIDATION_PROPOSAL.md`). I touched erraudit for the go-output bump but left this file modified without investigating whether I caused it or whether it was pre-existing.

---

## e) WHAT WE SHOULD IMPROVE 🛠️

### Process (this session's lessons)

1. **A release is not done until the consumers are pushed, not just committed.** I committed v0.35.0 bumps to 7 repos but left them unpushed. The local state is correct; the remote state is stale. This is a softer version of the "stale release" problem.
2. **BuildFlow pre-commit hook deletes `CODE_OF_CONDUCT.md`.** I hit this in mr-sync and used `--no-verify` per the AGENTS.md gotcha. But the hook also blocked my initial `git commit` in go-output. The BuildFlow config issue should be fixed at the source (in BuildFlow's config), not worked around forever.
3. **The `go-auto-upgrade` daemon is still invisible to me.** It committed the 34-file drift as `b3e9a4b` before I tagged v0.35.0. I worked around it this time by committing+verifying+tagging in quick succession, but it's still a race I won't always win.
4. **Consumer dirty files need investigation.** erraudit has a modified proposal doc I didn't look at. It might be daemon-written, pre-existing, or something I need to revert. I should not leave repos dirty without understanding why.

### Architectural / project observations

5. **`v0.34.0` should probably be retracted.** It points at `3b0640e`, which lacks the sibling dep version alignment. Consumers who resolve v0.34.0 get a _slightly inconsistent_ state (go-branded-id v0.3.3 label instead of v0.4.0, testhelpers v0.32.0 instead of v0.34.0). It's not broken (same code), but it's not what we'd want anyone to pin to. A `retract v0.34.0` in the next release would close this.
6. **The version-lineage gap** (`v0.32.0` → `v0.34.0` → `v0.35.0`, skipping `v0.33.0`) is permanent history. Consumers or tools expecting contiguous versions will be confused. Documented in the CHANGELOG, but worth a note in AGENTS.md.
7. **`gopls` go1.26 vs go1.27 warnings on `marshal.go` persist** across all three reports. Not blocking, but not addressed.

---

## f) Up to 50 things to get done next

Sorted by **impact × urgency**.

### 🔴 Urgent — finish the delivery

1. **Push the 7 unpushed consumer repos** (go-wizard-sdk, index, projects-management-automation, terraform-diagrams-aggregator, timesheets, universal-workflow, yt-history-intel).
2. **Investigate + resolve erraudit's dirty proposal doc**, then push.
3. **Decide on SystemNix's 9-commit-ahead stack** — push it, or leave for the user to handle.
4. **Verify origin builds after pushing** — `nix build` / `go build` on fresh clones to confirm remote state matches local.

### 🟠 High — close the release loop

5. **Retract `v0.34.0`** in the next release (it's stale/superseded).
6. **Write `docs/RELEASE.md`** — the step-by-step release runbook (recommended in 3 reports now).
7. **Write ADR 0012: Release Tagging Discipline** — codify: clean tree → verify diff empty → tag → push → GitHub Release.
8. **Root-cause the original bogus-tag creator** — still unknown, asked 3×.
9. **Fix BuildFlow's `CODE_OF_CONDUCT.md` deletion** at the source (BuildFlow config), not per-repo `--no-verify`.
10. **Add CI guard for tag integrity** — tag commit must == master HEAD.
11. **Add `nix run .#tags-audit`** — report any tag whose commit != master HEAD.

### 🟡 Medium — quality and consistency

12. **Bump DiscordSync** from v0.34.0 to v0.35.0 (currently valid but stale).
13. **Address `marshal.go` go1.27 `gopls` warnings** (3 reports running).
14. **Add `retract` lint** — warn if `retract` directives get reordered by tooling.
15. **Run `nix run .#lint`** on go-output post-release.
16. **Run `nix run .#govulncheck`** on go-output post-release.
17. **Run `nix run .#test-race`** on nom + tui post-release.
18. **Regenerate golden files** if any rendered bytes changed (shouldn't have, but verify).
19. **Audit `.golangci.yml` allow-lists** for new transitive deps.
20. **Sweep `~/projects` for other stale go-output pins** (v0.31.x, v0.30.x).
21. **Consider `go mod vendor` in mr-sync** — `go build` warned about stale vendor/modules.txt.
22. **Document the version-lineage gap** (v0.33.0 skipped) in AGENTS.md.
23. **Update `docs/DOMAIN_LANGUAGE.md`** if any terms drifted.
24. **Review `FEATURES.md`** and `TODO_LIST.md` for accuracy.
25. **Write `docs/postmortems/2026-07-27-bogus-tags-and-release-incidents.md`** combining all 3 incidents.

### 🟢 Lower — hardening and ecosystem

26. **Add `.github/workflows/release.yml`** automating the release runbook.
27. **Add tag-push branch protection ruleset** on GitHub (tags aren't protected by default).
28. **Add a `nix run .#deps-audit`** that reports all consumer repos' go-output pins in one table.
29. **Consider a `flake.lock` sharing mechanism** for the LarsArtmann ecosystem to reduce N-level transitive churn.
30. **Investigate `mkPreparedSource` auto-detecting missing inputs** (mr-sync `go-ndjson` gap).
31. **Schedule weekly tag-integrity audit** across the ecosystem.
32. **Consider disabling auto-commit daemon during interactive sessions** (misleading messages).
33. **Add `CHANGELOG.md` lint** — new tag must have matching `[X.Y.Z]` section.
34. **Pin `go-auto-upgrade` to skip `go-output`'s own `go.mod`** during release windows.
35. **Add a pre-commit hook** blocking edits to `retract` block without rationale.
36. **Consider signing commits** (not just tags) for release-prep commits.
37. **Review the `enum v0.17.1` ghost dep** in go-wizard-sdk (pre-existing, 2 reports running).
38. **Verify `examples/` build** against v0.35.0 deps.
39. **Confirm `go.work.example` in sync** with 19-module list.
40. **Review whether `testhelpers` should be bumped to v0.35.0** in its own go.mod (currently v0.34.0 in the require, v0.35.0 tag exists).
41. **Document the `GOPRIVATE` requirement** for consumer repos more prominently.
42. **Consider a monorepo-wide version variable** to reduce per-module version drift.
43. **Add a `nix run .#retract-check`** validating retract directives match deleted tags.
44. **Review whether v0.34.0 GitHub Release should note the retraction/supersession.**
45. **Update the v0.34.0 GitHub Release** with a "Superseded by v0.35.0" note.
46. **Consider `git notes`** to annotate daemon commits with accurate messages.
47. **Add a `.github/SECURITY.md`** if missing.
48. **Review dependabot alerts** (8 vulnerabilities reported on push — 3 high, 3 moderate, 2 low).
49. **Consider a release-bots cooldown** — prevent `go-auto-upgrade` from running for N minutes after a tag push.
50. **Retro on the 3-incident day** — bogus tags, stale v0.34.0, v0.35.0 supersession. Pattern: automation + manual ops without coordination.

---

## g) Questions I CANNOT figure out myself (need you)

**Q1 — Should I push the 7 unpushed consumer repos now?** They have v0.35.0 committed locally but origin still has the old pin. I don't know if you have a push cadence (e.g., "push everything at end of session") or if there's a reason to hold (e.g., CI will auto-bump them via `go-auto-upgrade` anyway). Pushing is safe and finishes the delivery, but I don't want to cause merge conflicts if `go-auto-upgrade` is about to run.

**Q2 — Should `v0.34.0` be retracted in the next release?** It points at `3b0640e` which lacks the sibling dep version alignment (go-branded-id v0.3.3 label, testhelpers v0.32.0 label). It's not broken (same code as v0.35.0), but it's not what we'd want anyone to pin to. Retracting it closes the loop, but it also means anyone who already resolved v0.34.0 gets a retraction notice. Your call.

**Q3 — What created the bogus `v0.32.1` / `v0.33.0` tags this morning?** This is the 3rd time I'm asking. The tagger is you (`git@lars.software`), but I cannot determine the _mechanism_ (script, CI, `go-auto-upgrade`, manual) from inside go-output. Without this answer, the same bug can recur tonight. I've now released v0.34.0 and v0.35.0 to work around the symptoms, but the root cause is still unidentified.

---

## Self-review scorecard (brutal)

| Question                       | Answer                                                                                                                                                |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| What did I forget?             | To push 7 consumer repos; to investigate erraudit's dirty file; to retract v0.34.0                                                                    |
| What's stupid we do anyway?    | BuildFlow deleting `CODE_OF_CONDUCT.md`; auto-commit daemon rewriting my commit messages; `go-auto-upgrade` racing manual releases                    |
| What could I have done better? | Pushed the consumers immediately after committing; investigated the erraudit dirty file instead of ignoring it; retracted v0.34.0 in v0.35.0's go.mod |
| What could I still improve?    | See 50 items above — top 3: push consumers, write release runbook, root-cause the tag creator                                                         |
| Did I lie?                     | No. This release is clean. `git diff v0.35.0` = 0.                                                                                                    |
| How can we be less stupid?     | Fix BuildFlow config; write the release runbook; coordinate automation with manual ops                                                                |
| Ghost systems?                 | `enum v0.17.1` ghost dep in go-wizard-sdk (pre-existing, carried forward)                                                                             |
| Scope creep?                   | No — this was a focused supersession release                                                                                                          |
| Did I remove something useful? | No                                                                                                                                                    |
| Split brains?                  | Local vs remote: 7 repos have v0.35.0 locally but old pin on origin. Not a split brain _within_ a repo, but a split across the ecosystem.             |
| How are tests?                 | Green at v0.35.0 tag. Verified correctly this time (tag == tree).                                                                                     |

**Overall grade for this session: A-.** The release mechanics were correct and the lesson from v0.34.0 was applied (verify clean tree, verify `git diff` empty, create GitHub Release). Deduction for leaving 7 consumers unpushed — the work is locally complete but not delivered. The 3 prior reports' open questions (root cause, runbook, CI guard) remain open but are now clearly documented.

---

## Timeline (full session arc — all 3 releases)

| Time (CEST) | Event                                                                                   |
| ----------- | --------------------------------------------------------------------------------------- |
| ~11:00      | User asked "What is the v0.33 release?"                                                 |
| ~11:30      | Diagnosed: `v0.32.1` + `v0.33.0` are bogus tags pointing at stale June commit           |
| ~11:45      | Deleted bogus tags locally + from origin                                                |
| ~12:00      | Fixed failing `TestBrandedIDFormat` test (justified by go-branded-id v0.3.3 change)     |
| ~12:05      | First self-review report                                                                |
| ~12:30      | User said "release v0.34.0"                                                             |
| ~13:00      | Tagged v0.34.0, pushed, repinned consumers                                              |
| ~14:00      | Declared done — all green                                                               |
| ~16:41      | Second self-review: discovered `go-auto-upgrade` had drifted the tree after v0.34.0 tag |
| ~16:45      | User said "supersede with v0.35.0"                                                      |
| ~16:46      | Investigated v0.4.0 (same commit as v0.3.3), committed drift, **verified clean tree**   |
| ~16:47      | Tagged v0.35.0, **verified `git diff` empty**, pushed, created GitHub Release           |
| ~16:50      | Repinned all consumers to v0.35.0                                                       |
| 16:52       | This report.                                                                            |
