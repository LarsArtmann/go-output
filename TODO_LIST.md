# TODO_LIST.md — go-output

**Last updated:** 2026-08-06
**Open items:** 10

---

## Test Infrastructure

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 1  | **Move `pollTeatestOutput` to `teatest_helpers_test.go`** — currently lives in `teatest_vt_test.go` but is used by the helpers file; fragile coupling. | 5 min   | Open   |
| 2  | **Strengthen `waitForVisible` conditions** — 14 call sites pass `"s"` which matches any English text; should wait for actual content (`"Build Module"`, etc.) | 20 min  | Open   |
| 3  | **Add goroutine-leak test** — `runtime.NumGoroutine()` before/after teatest cycle to catch leaks without relying on `-race` to surface them as hangs. | 20 min  | Open   |
| 4  | **Add `nix run .#test-race-all`** — race-test ALL 19 modules, not just nom + tui. Other modules could have races that only surface in CI.              | 15 min  | Open   |

## Code Quality

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 5  | **Align `go.mod` versions** — `daghtml/`, `escape/`, `testhelpers/` use `go 1.26.4` while all other modules use `go 1.26.5`. Align to `1.26.5`.       | 5 min   | Open   |
| 6  | **Fix `docs/ERROR_SYSTEM.md` contributor example** — uses unexported `joinStrings` (line 147); contributors can't use it. Replace with `strings.Join(output.EnumAllowedValues(...))` | 5 min   | Open   |

## CI / DevOps

| #  | Task                                                                                                                                                  | Effort  | Status |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ------ |
| 7  | **Add `.github/dependabot.yml`** for GitHub Actions version updates — actions are now SHA-pinned; dependabot can track new releases.                   | 10 min  | Open   |
| 8  | **Add quality gates to `scripts/pre-tag-check.sh`** — currently only build/test/race + tag parity. Add `art-dupl -t 4`, `govulncheck`, `golangci-lint`, golden-file freshness check. | 30 min  | Open   |

## Community (owner-dependent)

| #  | Task                                                                                                                                                  | Effort  | Status                        |
| -- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | ------- | ----------------------------- |
| 9  | **Post to r/golang, submit to Awesome Go**                                                                                                            | 30 min  | Open (needs owner account)    |
| 10 | **Cut `v1.0.0` tag** — API frozen per ADR 006; all v0.30.x–v0.37.x breaking changes shipped.                                                           | 2 min   | Prepared — awaiting owner tag |
