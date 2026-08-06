# TODO_LIST.md — go-output

**Last updated:** 2026-08-06
**Open items:** 16

---

## Release Process

The v0.37.0 tag was cut without submodule tags, without a release-prepare commit, and as a lightweight tag — repeating every v0.36.0 mistake. These items fix the systemic gap.

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 1  | **Create v0.37.0 submodule tags** (16 submodules) — no `*/v0.37.0` tags exist (`git tag -l '*/v0.37.0'` returns empty). Match v0.36.0 coverage.       | 5 min   | Open   |
| 2  | **Fix v0.37.0 root tag** — convert lightweight → annotated (`git cat-file -t v0.37.0` returns `commit`, should be `tag`). Matches v0.34.0/v0.35.0.    | 5 min   | Open   |
| 3  | **Add submodule auto-tagging step to `release.yml`** — iterate modules, create annotated `<module>/v<version>` tags on root tag push. Eliminates #1.   | 1 hr    | Open   |
| 4  | **Create `docs/RELEASE_CHECKLIST.md`** — the 8-step sequence: CHANGELOG → release-prepare commit → pre-tag-check → CI green → tag root + 16 submodules → push → verify GitHub Release → verify `go get` | 30 min  | Open   |
| 5  | **Add annotated-tag enforcement** — CI check or pre-receive hook that rejects lightweight tags (`git for-each-ref --format='%(objecttype)' refs/tags/v*` must return `tag`) | 30 min  | Open   |
| 6  | **Create `scripts/tag-release.sh`** — wrapper that runs the full release sequence: fetch, verify clean tree, pre-tag-check, tag root + submodules, verify parity | 45 min  | Open   |

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
