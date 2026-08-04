# Status Report — Backlog Burn-Down & Self-Review

**Date:** 2026-08-04 02:02 CEST
**Session scope:** Resolve 9 open backlog items (P0 CI fixes, P1 release hygiene/security, P2 error system polish) + 3 remaining P3 items. Then brutal self-review.
**Honesty mode:** BRUTAL. No spin.

---

## TL;DR — What happened this session

Resolved **12 of 12** actionable backlog items. All 19 modules build and test green. CI should turn green on next push (art-dupl pinned, TUI deadlock fixed). Found and fixed a significant doc-drift miss (`docs/ERROR_SYSTEM.md` still referenced deleted d2 sentinels) during self-review. 10 unpushed commits.

---

## a) FULLY DONE ✅ (12 items)

| #   | Item                                     | Evidence                                                                                                                                             |
| --- | ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Fix TUI teatest deadlock (P0)            | `tui/teatest_vt_test.go` — VT reconstruction hoisted out of polling loop. Passes 3× under `-race`. Commit `30f67be`.                                 |
| 2   | Fix art-dupl CI installation (P0)        | `.github/workflows/ci.yml` — pinned `@v0.6.0`. Commit `d750c16`.                                                                                     |
| 3   | Pin GitHub Actions to commit SHAs (P1)   | `ci.yml` + `release.yml` — checkout, setup-go, golangci-lint-action, softprops/action-gh-release all SHA-pinned. Commit `d750c16`.                   |
| 4   | Retract v0.34.0 tag (P1)                 | `go.mod` — `retract v0.34.0` added. Commit `d750c16`.                                                                                                |
| 5   | Create/verify GitHub Releases (P1)       | v0.34.0 retracted (no release needed). v0.35.0 + v0.36.0 confirmed have releases.                                                                    |
| 6   | Root-cause bogus-tag creator (P1)        | No automation creates tags in this repo. Manual error pointing at `194441b`. Retract directives are permanent mitigation.                            |
| 7   | Fix 10 dependabot vulnerabilities (P1)   | `website/package.json` + `package-lock.json` — astro v7.1.6, vite v8, esbuild 0.28.1. `npm audit` → 0 vulnerabilities. Commits `db92aa6`, `1f5de4c`. |
| 8   | Migrate d2 sentinels → typed errors (P2) | `d2/d2_enum.go` — 5 typed structs. `d2/error_contract_test.go` — `errors.AsType` tests. Commit `f92227b`.                                            |
| 9   | Cross-module error integration test (P2) | `integration/cross_module_error_test.go` — 5 test groups across root/d2/graph.                                                                       |
| 10  | Fix ADR numbering collision (P3)         | File already renamed. Fixed stale `ADR 0011` → `ADR 014` in AGENTS.md. Updated d2 error-system description.                                          |
| 11  | CQRS happy-path tests (P3)               | Added `WriteMermaid` (graph), `WriteGraph` + `Write` (d2), `WriteMarkdown` (tree) tests. All 14 `WriteXxx` covered.                                  |
| 12  | Flush() in TUI shutdown (P3)             | `tui/lifecycle.go` — `Stop()` calls `nomSubscriber.Flush()` before quit signal. Added test.                                                          |

**All 19 modules build green. All 19 test green (non-race). TUI passes under `-race`.**

---

## b) PARTIALLY DONE ⚠️

| #   | Item                            | Why partial                                                                                                                                                                                                                       |
| --- | ------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **docs/ERROR_SYSTEM.md update** | Fixed during self-review (d2 sentinel table → typed error table, code example updated). But only caught this because the user asked for a self-review — I should have updated it as part of the d2 migration. Points deducted.    |
| 2   | **CI verification**             | CI fixes are committed but NOT pushed. The next push will be the real test. I cannot verify CI is green until it runs. The art-dupl pin and TUI deadlock fix are high-confidence, but "works on my machine" is not "CI is green." |

---

## c) NOT STARTED ❌

None. All actionable items resolved. Remaining 2 TODO items (community tasks) are owner-dependent.

---

## d) TOTALLY FUCKED UP 💥

| #   | What                                                                                                              | Severity | Impact                                                                                                                                                                                                                                                                          |
| --- | ----------------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Updated `error_contract_test.go` with custom `contains`/`findSubstring` helpers instead of `strings.Contains`** | Medium   | Wrote 15 lines of reimplementation when 1 import would do. Caught it immediately and rewrote. But the first instinct was wrong — this is the "reinvent the standard library" anti-pattern.                                                                                      |
| 2   | **Cross-module test used `d2Err, _ :=` instead of `_, d2Err :=`**                                                 | Low      | Compile error — the first return value of `ParseDirection` is `Direction`, not `error`. Fixed on first test run. Classic destructuring mistake.                                                                                                                                 |
| 3   | **Mermaid test asserted `"graph"` in output when Mermaid uses `"flowchart TD"`**                                  | Low      | Didn't read the renderer source before writing the assertion. Fixed on first run.                                                                                                                                                                                               |
| 4   | **Website vite override fight**                                                                                   | Low      | Left the `vite: "7.3.2"` override in place when upgrading to astro v7 (which needs vite v8). Caused a build failure I had to debug. Should have checked astro v7's peer deps first.                                                                                             |
| 5   | **Didn't update `docs/ERROR_SYSTEM.md` during the d2 migration**                                                  | HIGH     | The consumer-facing error reference documented deleted sentinels and had a code example using `d2.ErrInvalidDirection` that would **fail to compile**. This is a living doc, not a historical doc. A consumer reading it would get broken code. Only caught during self-review. |

---

## e) WHAT WE SHOULD IMPROVE

### Process

1. **Living docs must be updated in the same commit as the code change.** The d2 sentinel→typed migration was committed (`f92227b`) without touching `docs/ERROR_SYSTEM.md`. That's a split brain. The AGENTS.md rule "update at the moment of discovery" was violated.

2. **Read the renderer source before writing assertions.** The Mermaid `flowchart` miss and the `d2Err, _` destructuring error both stem from assuming instead of reading. 30 seconds of reading saves a test-run-debug cycle.

3. **Don't fight overrides you don't understand.** The vite override battle wasted 3 minutes. The right move was: check what astro v7 expects, then adjust/remove the override in one shot.

4. **The `contains`/`findSubstring` reinvention** is a symptom of not reaching for `strings.Contains` first. Always check the standard library before writing utility helpers.

5. **Run lint (`golangci-lint`) after changes**, not just `go vet` + `go test`. I only ran `go vet` because golangci-lint wasn't in my shell. The Nix flake has `nix run .#lint` — I should have used it.

### Technical Debt Remaining

6. **10 unpushed commits.** The auto-daemon committed everything, but nothing is pushed to origin. CI hasn't run on any of these changes.

7. **Historical docs still reference old d2 sentinels.** `docs/planning/2026-07-06_*` and `docs/status/2026-07-30_*` mention `ErrInvalidDirection` etc. These are point-in-time docs and SHOULD NOT be changed — but they create search noise.

8. **`docs/ERROR_SYSTEM.md` contributor example uses `joinStrings`** which is an unexported root helper. Contributors can't actually use it. Should reference `strings.Join(output.EnumAllowedValues(...))` or `output.EnumAllowedStrings(...)`.

9. **No `nix flake check` was run.** The Nix formatting + pre-commit hooks haven't been validated.

10. **The d2 `error_contract_test.go` golden pattern** — I tested error messages include `(allowed:)` but didn't verify the exact format matches root/graph output. There could be subtle formatting differences (e.g., `strings.Join` vs `joinStrings` separator).

11. **The TUI `Flush()` call logs errors via `slog.Error` but doesn't surface them to the caller.** A caller can't tell if timing data was lost. Consider returning the error or adding a callback.

12. **npm binary was found via `/nix/store/32f9l3m7...`** — a fragile hardcoded path. The website should have npm available through the devShell or a documented `nix run` command.

---

## f) Up to 50 Things We Should Get Done Next

### Immediate (blocks green CI / next release)

| #   | Task                                                                           | Effort |
| --- | ------------------------------------------------------------------------------ | ------ |
| 1   | **Push the 10 unpushed commits** and watch CI turn green                       | 2 min  |
| 2   | **Run `nix run .#lint`** across all modules — I only ran `go vet`              | 5 min  |
| 3   | **Run `nix flake check`** — verify Nix formatting + hooks pass                 | 2 min  |
| 4   | **Verify `nix run .#test`** (the flake app) passes across all 19 modules       | 5 min  |
| 5   | **Run `nix run .#test-race`** — the race test app tests nom + tui specifically | 5 min  |

### Error System Polish

| #   | Task                                                                                                                     | Effort |
| --- | ------------------------------------------------------------------------------------------------------------------------ | ------ |
| 6   | Fix `docs/ERROR_SYSTEM.md` contributor example: replace `joinStrings` with `strings.Join(output.EnumAllowedValues(...))` | 5 min  |
| 7   | Add d2 typed errors to the root `error_contract_test.go` patterns doc (if any consumers reference d2 errors)             | 15 min |
| 8   | Consider whether `d2.AllDirections()` should filter empty-string `DirDown` from the `Allowed` list in error messages     | 10 min |
| 9   | Verify all 5 d2 typed error messages are byte-identical in format to root/graph typed errors                             | 10 min |

### TUI / nom

| #   | Task                                                                                                                | Effort |
| --- | ------------------------------------------------------------------------------------------------------------------- | ------ |
| 10  | Surface `Flush()` errors from TUI to caller (return error or callback) instead of just `slog.Error`                 | 30 min |
| 11  | Add a teatest E2E test that verifies the full lifecycle: start → activities → stop → timing cache persisted         | 1 hour |
| 12  | Add `Flush()` to the ctrl+c / error shutdown path, not just the clean `Stop()` path                                 | 15 min |
| 13  | Verify the `Flush()` call doesn't deadlock if the TUI program is still rendering (bubbletea goroutine vs Flush I/O) | 30 min |
| 14  | Profile nom `TimingCache.Flush()` under load to ensure it doesn't block the UI thread                               | 30 min |

### Test Coverage

| #   | Task                                                                                                                                             | Effort |
| --- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------ |
| 15  | Add error-path tests for the 4 newly-tested CQRS Write functions (not just happy path)                                                           | 30 min |
| 16  | Add integration test for `d2.Write` + `d2.WriteGraph` with options (WithDirection, WithLayout, WithTitle)                                        | 20 min |
| 17  | Add integration test for `graph.WriteMermaid` with options                                                                                       | 15 min |
| 18  | Add fuzz test for `d2.ParseDirection`/`ParseNodeShape` etc. (enum fuzz coverage)                                                                 | 20 min |
| 19  | Verify golden files are still byte-identical after the d2 error migration (golden tests check rendered output, not errors, but worth confirming) | 10 min |
| 20  | Add `TestCQRS_StreamVsRegistry` for YAML/TOML/HTML/AsciiDoc (currently only JSON + CSV have this cross-check)                                    | 1 hour |

### CI / DevOps

| #   | Task                                                                                                       | Effort |
| --- | ---------------------------------------------------------------------------------------------------------- | ------ |
| 21  | Pin `golang.org/x/vuln/cmd/govulncheck@latest` to a specific version in CI                                 | 5 min  |
| 22  | Add a `name:` comment to the retract directives in go.mod for future grep-ability                          | 2 min  |
| 23  | Consider adding `go mod verify` to CI to catch corrupted module caches                                     | 10 min |
| 24  | Add CI job for `nix flake check` (if Nix is available in CI)                                               | 30 min |
| 25  | Consider adding dependabot config for Go modules (currently only npm alerts exist)                         | 15 min |
| 26  | Add a `.github/dependabot.yml` for GitHub Actions version updates                                          | 10 min |
| 27  | Verify the SHA-pinned actions are the latest stable (they were checked at session time, but versions move) | 5 min  |

### Documentation

| #   | Task                                                                                                                 | Effort |
| --- | -------------------------------------------------------------------------------------------------------------------- | ------ |
| 28  | Update `docs/adr/0013-error-system-design.md` to reflect d2 typed-error migration (it may still reference sentinels) | 15 min |
| 29  | Add a "Migration Guide" section to CHANGELOG.md for the d2 sentinel→typed breaking change                            | 20 min |
| 30  | Update `README.md` error handling section if it references d2 sentinels                                              | 10 min |
| 31  | Verify `FEATURES.md` error system row reflects the current state (typed errors across all modules)                   | 10 min |
| 32  | Add `docs/CQRS_TEST_COVERAGE.md` documenting which WriteXxx functions have tests and what patterns are used          | 30 min |

### Architecture / Code Quality

| #   | Task                                                                                                                         | Effort  |
| --- | ---------------------------------------------------------------------------------------------------------------------------- | ------- |
| 33  | Run `art-dupl -t 4` to verify the d2 typed-error structs didn't introduce duplication vs root/graph patterns                 | 10 min  |
| 34  | Consider extracting a shared `invalidEnumError[T]` generic type to reduce the 12+ typed error structs across modules         | 2 hours |
| 35  | Run `erraudit` on the new d2 error structs to verify they follow the project conventions                                     | 15 min  |
| 36  | Consider whether `d2.AllDirections()` should return a copy (like `AllConstraints()`) vs the backing slice (current behavior) | 15 min  |
| 37  | Verify the `website/` astro v7 upgrade didn't break any existing content (check CSP headers, OG images, page structure)      | 30 min  |
| 38  | Run `nix run .#govulncheck` to verify no new Go vulnerabilities were introduced                                              | 10 min  |
| 39  | Consider adding `// Deprecated:` comments to the retract directives for IDE warnings                                         | 5 min   |

### Release

| #   | Task                                                                              | Effort |
| --- | --------------------------------------------------------------------------------- | ------ |
| 40  | Cut `v0.37.0` after CI is confirmed green — this session shipped 12 items         | 30 min |
| 41  | Create a GitHub Release for v0.37.0 with the CHANGELOG entries                    | 10 min |
| 42  | Tag `testhelpers/v0.37.0` (only independently versioned module)                   | 5 min  |
| 43  | Consider cutting `v1.0.0` — API is frozen (ADR 006), all breaking changes shipped | 1 hour |

### Website

| #   | Task                                                                                             | Effort |
| --- | ------------------------------------------------------------------------------------------------ | ------ |
| 44  | Deploy the updated website to Firebase Hosting and verify the astro v7 upgrade renders correctly | 30 min |
| 45  | Run Lighthouse audit on the updated website                                                      | 15 min |
| 46  | Verify OG image generation works with astro-og-canvas v0.13                                      | 15 min |

### Misc

| #   | Task                                                                                                                      | Effort |
| --- | ------------------------------------------------------------------------------------------------------------------------- | ------ |
| 47  | Add `npm audit` to CI for the website directory                                                                           | 15 min |
| 48  | Document the npm dependency management approach in AGENTS.md (how to update, what overrides exist and why)                | 20 min |
| 49  | Run `go mod tidy` across all modules to ensure cleanliness                                                                | 5 min  |
| 50  | Clean up `docs/status/` — there are 20+ status reports; consider archiving old ones to a `docs/status/archive/` directory | 30 min |

---

## g) Questions (3)

### Q1: Should I push the 10 unpushed commits to origin/master now?

There are 10 unpushed commits on master. The CI fixes (art-dupl pin, TUI deadlock, SHA pinning) and the d2 error migration are all in there. CI hasn't run on any of them. I did not push because the AGENTS.md says "NEVER push to remote unless explicitly asked." Should I push now so CI can verify, or do you want to review first?

### Q2: Should the d2 sentinel→typed-error migration be a breaking-change CHANGELOG entry for v0.37.0?

The d2 sentinels (`ErrInvalidDirection`, `ErrInvalidNodeShape`, etc.) were exported and consumers may have used `errors.Is(err, d2.ErrInvalidDirection)`. They are now gone — replaced with typed structs matched via `errors.AsType[*d2.InvalidDirectionError]`. This is a breaking change. Should I document it as such in the CHANGELOG under a `[Unreleased]` → `### Changed` → `**BREAKING**` heading, or is this an internal-enough change that it doesn't need a migration note?

### Q3: Should `Flush()` error in TUI `Stop()` be surfaced to the caller or stay as `slog.Error`?

Currently `Stop()` calls `Flush()` and logs errors via `slog.Error` but swallows them (no return value). The caller (`BubbleTeaProgressReporter.Stop()` has no return value either). If timing data fails to persist, the caller has no way to know. Should I:

- **(a)** Leave as-is (timing data is best-effort, not critical), or
- **(b)** Add a `Stop() error` variant or a `FlushOnStop` callback, or
- **(c)** Add a `LastFlushError()` accessor the caller can check after `Stop()`?

---

## Build & Test Status

```
All 19 modules: BUILD OK
All 19 modules: TEST OK (non-race)
TUI (race): TEST OK
go vet (root, d2, tui, integration): CLEAN
npm audit (website): 0 vulnerabilities
CI: NOT VERIFIED (10 unpushed commits)
```
