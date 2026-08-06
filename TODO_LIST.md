# TODO_LIST.md — go-output

**Last updated:** 2026-08-06
**Open items:** 9

---

## Release Process

The v0.37.0 tag was cut without submodule tags, without a release-prepare commit, and as a lightweight tag — repeating every v0.36.0 mistake. **Items #1–#6 were resolved on 2026-08-06.** The tag family was created (17 annotated tags), root tag converted to annotated, `release.yml` gained auto-tagging, `scripts/tag-release.sh` and `docs/RELEASE_CHECKLIST.md` were created, and `pre-tag-check.sh` now verifies tag-family parity.

| #  | Task                                                                                                                                                  | Effort  | Status   |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | -------- |
| 1  | **Create v0.37.0 submodule tags** (16 submodules) — annotated, on `dd051b4`, parity with v0.36.0 verified.                                             | 5 min   | ✅ Done  |
| 2  | **Fix v0.37.0 root tag** — converted lightweight → annotated. Now `tag` type on `dd051b4`.                                                             | 5 min   | ✅ Done  |
| 3  | **Add submodule auto-tagging step to `release.yml`** — new "Create and push submodule tags" step + `daghtml` added to all module lists (was missing).  | 1 hr    | ✅ Done  |
| 4  | **Create `docs/RELEASE_CHECKLIST.md`** — 8-step sequence with tag convention, manual/automated paths, and recovery procedures.                         | 30 min  | ✅ Done  |
| 5  | **Add annotated-tag enforcement** — `pre-tag-check.sh` now verifies latest release has all annotated tags (parity check). release.yml creates annotated.| 30 min  | ✅ Done  |
| 6  | **Create `scripts/tag-release.sh`** — full release wrapper: fetch, verify clean tree, pre-tag-check, tag root + 16 submodules, verify parity.          | 45 min  | ✅ Done  |

## Test Infrastructure

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 7  | **Move `pollTeatestOutput` to `teatest_helpers_test.go`** — currently lives in `teatest_vt_test.go` but is used by the helpers file; fragile coupling. | 5 min   | Open   |
| 8  | **Strengthen `waitForVisible` conditions** — 14 call sites pass `"s"` which matches any English text; should wait for actual content (`"Build Module"`, etc.) | 20 min  | Open   |
| 9  | **Add goroutine-leak test** — `runtime.NumGoroutine()` before/after teatest cycle to catch leaks without relying on `-race` to surface them as hangs. | 20 min  | Open   |
| 10 | **Add `nix run .#test-race-all`** — race-test ALL 19 modules, not just nom + tui. Other modules could have races that only surface in CI.              | 15 min  | Open   |

## Code Quality

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 11 | **Align `go.mod` versions** — `daghtml/`, `escape/`, `testhelpers/` use `go 1.26.4` while all other modules use `go 1.26.5`. Align to `1.26.5`.       | 5 min   | Open   |
| 12 | **Fix `docs/ERROR_SYSTEM.md` contributor example** — uses unexported `joinStrings` (line 147); contributors can't use it. Replace with `strings.Join(output.EnumAllowedValues(...))` | 5 min   | Open   |

## CI / DevOps

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 13 | **Add `.github/dependabot.yml`** for GitHub Actions version updates — actions are now SHA-pinned; dependabot can track new releases.                   | 10 min  | Open   |
| 14 | **Add quality gates to `scripts/pre-tag-check.sh`** — currently only build/test/race. Add `art-dupl -t 4`, `govulncheck`, `golangci-lint`, golden-file freshness check. | 30 min  | Open   |

## Community (owner-dependent)

| #  | Task                                                                                                                                                  | Effort  | Status                        |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ----------------------------- |
| 15 | **Post to r/golang, submit to Awesome Go**                                                                                                            | 30 min  | Open (needs owner account)    |
| 16 | **Cut `v1.0.0` tag** — API frozen per ADR 006; all v0.30.x–v0.37.x breaking changes shipped.                                                           | 2 min   | Prepared — awaiting owner tag |
