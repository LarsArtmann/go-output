# Status: TUI Deadlock Real Fix + CI Hardening

**Generated:** 2026-08-04 05:47 CEST
**Session Focus:** Continue from prior session's CI health work — correct stale docs, harden CI, find and fix the REAL TUI deadlock root cause
**Overall Health:** All 19 modules passing locally · 0 lint issues · 0 vulnerabilities · 0 duplication (t=4) · **2 commits unpushed** · CI still red (last run used the incomplete fix)

---

## a) FULLY DONE

### 1. TUI Test Deadlock — REAL FIX (the P0 blocker)

**What happened:** The prior session claimed the TUI deadlock was fixed by replacing `teatest.WaitFor` with `pollTeatestOutput` in the VT test only. But `waitForVisible` — used by **5 of 9 teatest tests** — STILL called `teatest.WaitFor`, which internally uses `io.ReadAll(io.TeeReader(...))` that blocks indefinitely when the program's tick loop writes continuously (the output buffer never empties for EOF). Additionally, test cleanup called only `tm.Quit()` (fire-and-forget) without `tm.WaitFinished()`, leaking the program goroutine.

**The actual fix (two parts, `tui/teatest_helpers_test.go`):**

1. `waitForVisible` → now delegates to `pollTeatestOutput` (bounded single-`Read` polling, already existed in `teatest_vt_test.go`) instead of `teatest.WaitFor`
2. `newTeatestModel` cleanup → added `tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))` after `tm.Quit()` to drain the program goroutine

**Verification:**

- All 10 `TestTeatest_*` tests pass under `-race` in **2.0 seconds** (was 600-second hang → 10-minute CI timeout)
- `nix run .#test-race` passes (nom 1.6s + tui 2.1s)
- `sync.Once` in `teatest.TestModel.waitDone()` makes double-`WaitFinished` calls safe (tests that call it in body + cleanup)

**Key lesson:** Never trust a prior session's diagnosis. `grep -rn "teatest.WaitFor"` should have been the FIRST thing run.

### 2. CI Test Timeout Added

- Added `-timeout 120s` to `go test -v -race` in **both** `ci.yml` and `release.yml`
- Previously: no timeout → Go's 10-minute default → CI burned 11+ minutes before failing
- Now: deadlocks surface in 2 minutes max

### 3. Art-dupl Threshold Lowered (t=50 → t=4)

- `ci.yml` duplication check lowered from t=50 (cosmetic) to **t=4** (the production gate per ADR 008)
- Verified locally: **0 clone groups at t=4** (codebase is clean)
- CI will now surface real duplication instead of rubber-stamping everything

### 4. GitHub Actions Bumped to Node.js 24

All four pinned actions in both `ci.yml` and `release.yml` updated to eliminate Node.js 20 deprecation warnings:

| Action                          | Old            | New            |
| ------------------------------- | -------------- | -------------- |
| `actions/checkout`              | v4 (`11d5960`) | v7 (`3d3c42e`) |
| `actions/setup-go`              | v5 (`40f1582`) | v7 (`b7ad1da`) |
| `golangci/golangci-lint-action` | v7 (`9fae48a`) | v9 (`ba0d7d2`) |
| `softprops/action-gh-release`   | v1 (`de2c0eb`) | v3 (`3d0d988`) |

SHAs fetched from GitHub API (lightweight tags used directly; annotated tags dereferenced to commit SHAs).

### 5. CHANGELOG Corrected with Real Root Causes

- TUI deadlock entry: now documents both root causes (`waitForVisible` still using `teatest.WaitFor` + missing `WaitFinished` in cleanup) instead of the prior session's wrong "vtScreenFromBytes in polling loop" narrative
- art-dupl entry: now documents the real root cause (`report_templ.go` globally gitignored) instead of "v0.6.1 has a broken build"
- Added entries for: CI test timeout, art-dupl threshold, Node.js 24 bump, nom lint fixes, d2 test refactor, 39 annotation conversions

### 6. TODO_LIST Items #1 and #2 Corrected

Same root-cause corrections as CHANGELOG — the prior session's wrong diagnoses replaced with the actual findings.

### 7. AGENTS.md Updated

Added critical test-infrastructure guidance to the teatest/v2 E2E pattern entry:

- "Never use `teatest.WaitFor` directly" with explanation of the `io.ReadAll` deadlock
- `pollTeatestOutput()` is the canonical replacement
- `WaitFinished` in cleanup is mandatory — without it the program goroutine leaks

### 8. All 39 Annotation Conversions Validated

- 0 files still have old `## Resolution` sections (markdown) or `<!-- Resolution -->` comments (HTML)
- 30 markdown files have `> **✅ Resolved:**` blockquotes (consistent formatting spot-checked across 10 files)
- 5 HTML files have visible `<div>` resolution banners after `<body>`
- 0 double conversions (no file has both old and new format)
- 1 file (`2026-08-04_05-00_ci-health-fixes-and-lint-cleanup.md`) mentions `<!-- Resolution -->` in its body text — that's a description of the conversion, not an unconverted annotation

### 9. Full Local Quality Gate Verification

| Gate                    | Result                     |
| ----------------------- | -------------------------- |
| `nix run .#build`       | 19/19 modules compile      |
| `nix run .#test`        | 19/19 pass                 |
| `nix run .#test-race`   | nom + tui pass (2s total)  |
| `nix run .#lint`        | 0 issues across 19 modules |
| `nix run .#govulncheck` | 0 vulnerabilities          |
| `nix flake check`       | All checks passed          |
| `art-dupl -t 4`         | 0 clone groups             |
| Dependabot alerts       | 0 open (12 fixed)          |

---

## b) PARTIALLY DONE

### 1. CI Green Status — NOT YET VERIFIED

The two commits with the real fix (`307c304` + `6e43912`) are **unpushed**. The last CI run (`30874219640`) used the incomplete fix from the prior session and failed at 3m44s (timeout worked — the 120s limit caught the deadlock). Once the daemon pushes, CI should go green for the first time since July 6, but this is **not yet confirmed**.

### 2. Action SHA Verification — NOT TESTED IN CI

The new SHAs were fetched from the GitHub API but have not been exercised in an actual CI run. If any SHA is wrong (e.g., pointing to a pre-release or a different branch's commit), CI will fail on the action setup step.

---

## c) NOT STARTED

- **Push 2 unpushed commits** — relying on the auto-git daemon; could push manually
- **Cut v0.37.0 tag** — pending CI green confirmation
- **Update prior status report** (`2026-08-04_05-00_ci-health-fixes-and-lint-cleanup.md`) — still has the wrong TUI deadlock diagnosis
- **dependabot.yml for GitHub Actions** — now that actions are pinned to SHAs, dependabot could keep them updated
- **Contributing upstream to teatest/v2** — the `io.ReadAll` deadlock is a library-level bug worth reporting

---

## d) TOTALLY FUCKED UP

### 1. Trusted the Prior Session's Diagnosis Without Verification

**This is the biggest failure of this session.** The context summary said the TUI deadlock was fixed. I started by fixing **documentation** (CHANGELOG, TODO_LIST) instead of verifying the fix was real. I should have run `nix run .#test-race` as my **first action** — it would have immediately exposed that `waitForVisible` still used `teatest.WaitFor`.

Instead, I corrected docs based on a wrong diagnosis, then had to correct them **again** when I discovered the real root causes after CI failed. Double work, double churn.

### 2. Fixed Docs Before Fixing Code

I spent the first 30 minutes on CHANGELOG/TODO_LIST text corrections while the P0 CI blocker was still broken. The correct order was: **verify the fix → fix the code → then document accurately**.

### 3. Didn't Run `-race` Tests Until After CI Failed

`nix run .#test` (without `-race`) passed fine — the deadlock only manifests under `-race`. I should have known this from the context summary ("surfaces under -race") and run `nix run .#test-race` before declaring any task complete.

### 4. Relied on the Daemon to Push

The user's context explicitly flagged the auto-git daemon as an ongoing problem. Two commits are sitting unpushed. I should have pushed manually after verifying the fix.

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Always verify before documenting.** Never write a CHANGELOG entry for a fix you haven't personally confirmed works. The prior session's wrong diagnosis propagated into CHANGELOG, TODO_LIST, and AGENTS.md — three places to correct instead of zero.

2. **Run the exact CI command locally before pushing.** CI uses `go test -v -race`. `nix run .#test` doesn't use `-race`. `nix run .#test-race` only covers nom + tui. There's a gap: the root + other 17 modules aren't race-tested locally. Consider adding `nix run .#test-race-all` that covers every module.

3. **Grep for the pattern, not the symptom.** When someone says "fixed teatest.WaitFor deadlock", grep for `teatest.WaitFor` — if ANY usage remains, the fix is incomplete. This would have taken 2 seconds.

4. **Push manually when the fix is critical.** The daemon runs on an interval. A P0 CI fix should not wait for a daemon — push it immediately after verification.

### Technical

5. **`pollTeatestOutput` lives in `teatest_vt_test.go` but is used by `teatest_helpers_test.go`.** This is fragile — if someone moves or deletes the VT test file, the helpers file breaks. The helper should live in `teatest_helpers_test.go` (the canonical location for shared test infrastructure) and the VT-specific `vtScreenFromBytes` stays in the VT file.

6. **The `waitForVisible` condition `"s"` is extremely weak.** Every test calls `waitForVisible(t, tm, "s")` — the letter "s" appears in almost any English text. This means the condition matches on the first render regardless of whether the actual content (activity labels, summary bar) is present. The tests pass but they're not actually verifying what they claim to verify. They should wait for specific content like `"Build Module"` or `"Run Tests"`.

7. **No test verifies that teatest cleanup actually drains goroutines.** A test that counts `runtime.NumGoroutine()` before and after a teatest test cycle would catch goroutine leaks immediately, without relying on `-race` to surface them as hangs.

---

## f) Up to 50 Things to Do Next

| #   | Task                                                                                                           | Impact   | Effort   |
| --- | -------------------------------------------------------------------------------------------------------------- | -------- | -------- |
| 1   | **Push 2 unpushed commits** (`git push origin master`)                                                         | Critical | 5s       |
| 2   | **Monitor CI to green** — first green run since July 6                                                         | Critical | 5min     |
| 3   | **Cut v0.37.0 tag** after CI confirmed green                                                                   | High     | 2min     |
| 4   | **Update prior status report** (`05-00`) with real root causes                                                 | Medium   | 5min     |
| 5   | **Move `pollTeatestOutput` to `teatest_helpers_test.go`** (canonical location)                                 | Medium   | 5min     |
| 6   | **Strengthen `waitForVisible` conditions** — wait for actual content, not `"s"`                                | Medium   | 15min    |
| 7   | **Add goroutine-leak test** — `runtime.NumGoroutine()` before/after teatest cycle                              | Medium   | 20min    |
| 8   | **Add `nix run .#test-race-all`** — race-test ALL 19 modules, not just nom+tui                                 | Medium   | 15min    |
| 9   | **Add dependabot.yml** for GitHub Actions version updates                                                      | Low      | 10min    |
| 10  | **File upstream bug report** on teatest/v2 `io.ReadAll` deadlock                                               | Low      | 15min    |
| 11  | **Post to r/golang + submit to Awesome Go** (TODO_LIST #14)                                                    | Low      | 30min    |
| 12  | **Cut v1.0.0 tag** (TODO_LIST #15) — API frozen per ADR 006                                                    | Low      | 2min     |
| 13  | **Investigate `integration/roundtrip_test.go` gopls warnings** — `json.Unmarshal` flagged as requiring go1.27  | Low      | 15min    |
| 14  | **Verify GitHub Action SHAs work in CI** — currently untested                                                  | Medium   | passive  |
| 15  | **Consider adding `CODEOWNERS`** file                                                                          | Low      | 5min     |
| 16  | **Consider adding `SECURITY.md`**                                                                              | Low      | 10min    |
| 17  | **Add CI status badge to README.md**                                                                           | Low      | 5min     |
| 18  | **The CI coverage step uses `bc`** — consider a Go-native alternative                                          | Low      | 30min    |
| 19  | **Parallelize CI jobs per module** instead of sequential loop                                                  | Low      | 1hr      |
| 20  | **Run `nix run .#tidy`** to ensure all go.mod files are clean                                                  | Low      | 2min     |
| 21  | **Add pre-commit hook** that greps for `teatest.WaitFor` to prevent reintroduction                             | Low      | 10min    |
| 22  | **Verify `TestTeatest_WindowSize_Propagates` cleanup** — it creates an inline model, not via `newTeatestModel` | Low      | 5min     |
| 23  | **Consider whether 3s `WaitFinished` timeout is sufficient for slow CI runners**                               | Low      | 5min     |
| 24  | **Run full module test suite under `-race` locally** to catch races in other modules                           | Medium   | 10min    |
| 25  | **Archive old status reports** — 50+ files in `docs/status/`                                                   | Low      | 30min    |
| 26  | **Consider whether `-timeout 120s` is too aggressive** for the full 19-module suite in CI                      | Low      | passive  |
| 27  | **Verify release.yml works end-to-end** (only fires on tag push)                                               | Low      | passive  |
| 28  | **Check if `website/` dependabot checks are passing in CI**                                                    | Low      | 5min     |
| 29  | **Consider adding a funding.yml**                                                                              | Low      | 5min     |
| 30  | **Review whether the auto-git daemon should be disabled** (ongoing discussion)                                 | Medium   | decision |
| 31  | **Add `CONTRIBUTING.md`** with GOEXPERIMENT=jsonv2 requirement                                                 | Low      | 20min    |
| 32  | **Verify the 39 annotation conversions render correctly in a browser** (HTML files)                            | Low      | 15min    |
| 33  | **Consider adding a lint rule** for `teatest.WaitFor` usage (custom golangci-lint linter)                      | Low      | 1hr      |
| 34  | **Document the CI workflow architecture** in a docs/ ADR                                                       | Low      | 30min    |
| 35  | **Run `art-dupl --sort total-tokens -t 1 --type-aware`** for a strict dedup audit                              | Low      | 10min    |
| 36  | **Verify the d2 `assertWrappedTypedError[T error]` refactor compiles in CI**                                   | Low      | passive  |
| 37  | **Consider whether nom lint fixes changed behavior** (exhaustive switch, makezero)                             | Low      | 10min    |
| 38  | **Add a benchmark for `pollTeatestOutput`** to verify it's not a performance regression                        | Low      | 20min    |
| 39  | **Check if any other test files in the codebase use `teatest.WaitFor`**                                        | Low      | 2min     |
| 40  | **Consider contributing `pollTeatestOutput` upstream to teatest/v2**                                           | Low      | 30min    |
| 41  | **Add a `.editorconfig`** for consistent formatting across editors                                             | Low      | 10min    |
| 42  | **Review ROADMAP.md** for accuracy against current state                                                       | Low      | 15min    |
| 43  | **Consider whether the coverage threshold (80%) is being met** in CI                                           | Low      | 5min     |
| 44  | **Verify FEATURES.md** is up to date with v0.36.0 changes                                                      | Low      | 15min    |
| 45  | **Consider adding a `Makefile`-to-flake migration guide** for non-Nix contributors                             | Low      | 30min    |
| 46  | **Review whether the CI workflow needs `permissions: contents: read`** explicitly                              | Low      | 5min     |
| 47  | **Check if `GOLANGCI_LINT_VERSION: v2.12` in CI matches the Nix devShell**                                     | Low      | 5min     |
| 48  | **Consider adding a `renovate.json`** as an alternative to dependabot for Go modules                           | Low      | 15min    |
| 49  | **Add a test that runs under `GOMAXPROCS=1`** to catch scheduling-dependent races                              | Low      | 20min    |
| 50  | **Celebrate when CI goes green** — it's been red since July 6                                                  | Low      | deserved |

---

## g) Questions (3)

### Q1: Should I push the 2 unpushed commits manually now?

The auto-git daemon hasn't pushed them yet. The TUI fix is verified locally (all 10 teatest tests pass under `-race` in 2s). Every minute we wait is a minute CI stays red. The daemon may push in 5 minutes or 5 hours — I can't tell. Should I `git push origin master` immediately, or do you want to review the diff first?

### Q2: Should the `waitForVisible` tests be strengthened in this session or deferred?

Currently every `waitForVisible` call waits for `"s"` — a substring so short it matches on the first render regardless of actual content. The tests technically pass but aren't verifying what their names suggest (e.g., `TestTeatest_ProgramStarts_RendersContent` doesn't verify NOM content is rendered, just that the letter "s" appears). Fixing this means changing the condition to specific content like `"Build Module"` — but that may surface rendering timing issues that make the tests flaky. Strengthen now or defer to a follow-up?

### Q3: Should I disable the auto-git daemon?

It has created multiple problems this session and prior sessions: empty commit messages (`d056ab1`), unpredictable push timing (2 critical commits sitting unpushed), and commits that mix unrelated changes. The alternative is committing manually with meaningful messages. The daemon provides convenience (never forget to commit) at the cost of control (can't predict when things land on origin). Should I disable it, or keep it and just push manually when timing matters?
