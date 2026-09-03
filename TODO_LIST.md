# TODO_LIST.md — go-output

**Last updated:** 2026-08-16
**Open items:** 12

---

## Test Infrastructure

| # | Task                                                                                                                                                          | Effort | Status |
| - | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 1 | **Move `pollTeatestOutput` to `teatest_helpers_test.go`** — currently lives in `teatest_vt_test.go` but is used by the helpers file; fragile coupling.        | 5 min  | Done   |
| 2 | **Strengthen `waitForVisible` conditions** — 14 call sites pass `"s"` which matches any English text; should wait for actual content (`"Build Module"`, etc.) | 20 min | Done   |
| 3 | **Add goroutine-leak test** — `runtime.NumGoroutine()` before/after teatest cycle to catch leaks without relying on `-race` to surface them as hangs.         | 20 min | Done   |
| 4 | **Add `nix run .#test-race-all`** — race-test ALL 19 modules, not just nom + tui. Other modules could have races that only surface in CI.                     | 15 min | Done   |

## Code Quality

| # | Task                                                                                                                                                                                 | Effort | Status |
| - | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 5 | **Align `go.mod` versions** — `daghtml/`, `escape/`, `testhelpers/` use `go 1.26.4` while all other modules use `go 1.26.5`. Align to `1.26.5`.                                      | 5 min  | Done   |
| 6 | **Fix `docs/ERROR_SYSTEM.md` contributor example** — uses unexported `joinStrings` (line 147); contributors can't use it. Replace with `strings.Join(output.EnumAllowedValues(...))` | 5 min  | Done   |

## CI / DevOps

| # | Task                                                                                                                                                                                 | Effort | Status |
| - | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------ | ------ |
| 7 | **Add `.github/dependabot.yml`** for GitHub Actions version updates — actions are now SHA-pinned; dependabot can track new releases.                                                 | 10 min | Done   |
| 8 | **Add quality gates to `scripts/pre-tag-check.sh`** — currently only build/test/race + tag parity. Add `art-dupl -t 4`, `govulncheck`, `golangci-lint`, golden-file freshness check. | 30 min | Done   |

## Community (owner-dependent)

| #  | Task                                                                                         | Effort | Status                        |
| -- | -------------------------------------------------------------------------------------------- | ------ | ----------------------------- |
| 9  | **Post to r/golang, submit to Awesome Go**                                                   | 30 min | Open (needs owner account)    |
| 10 | **Cut `v1.0.0` tag** — API frozen per ADR 006; all v0.30.x–v0.37.x breaking changes shipped. | 2 min  | Prepared — awaiting owner tag |

## v0.38.0 — harvested from the 2026-08-16 full code review

Behavior-drift items deliberately deferred (not bugs; API-contract decisions).
Source: `docs/reviews/2026-08-16_12-15_full-code-review.html`.

| #  | Task                                                                                                                                                                                                                                                                  | Effort | Status |
| -- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 11 | **Decide `Finish(err)` parameter contract** — `nom.InlineRenderer.Finish(workflowErr)` accepts the error but never renders it (documented twice as caller-owned). Either render it or drop the parameter.                                                             | 15 min | Open   |
| 12 | **Add `VisibleEntry.Kind` field** — layered-mode separator handling rune-sniffs entry strings to tell separators from rows; a typed kind field removes the guessing.                                                                                                  | 30 min | Open   |
| 13 | **Unify registry-vs-CQRS trailing-newline behavior** — diagram renderTable paths use `WriteRendered` (adds `\n`) while some CQRS writers use `WriteRenderedRaw`; audit all 16 formats for one rule.                                                                   | 45 min | Open   |
| 14 | **Remove `plantuml.WithDiagramType` dead option** — accepted, stored, never read. P3 dead API.                                                                                                                                                                        | 10 min | Open   |
| 15 | **Layered separator alignment ≥10 layers** — separator width doesn't account for double-digit layer numbers ("Layer 10"). Cosmetic.                                                                                                                                   | 15 min | Open   |
| 16 | **Derive format-count tripwires from literal lists** — the magic `16` is copy-pasted in `integration/format_registration_test.go`, `bdd/capability_test.go`, and `integration/doc.go`; derive from explicit format names so failures say what's missing.              | 20 min | Open   |
| 17 | **Slim `TestFormatCategories` matrix re-encode** — `integration/format_test.go` hand-duplicates the whole shape matrix; keep only load-bearing entries (e.g. CSV-not-graph) and assert stable invariants instead.                                                     | 25 min | Open   |
| 18 | **Basic example: one table construction story** — `main.go` builds via `output.NewTable` (root data model) while `renderers.go` uses `table.New` (lipgloss renderer); both valid, but the teaching example should show `output.NewTableBuilder()` once and render it. | 30 min | Open   |
| 19 | **Make nil-writer dispatch test non-chatty** — `integration/error_test.go` "nil writer defaults to stdout" writes a real table to stdout inside a parallel test; inject `io.Discard` where possible or accept and document the noise.                                 | 10 min | Open   |
| 20 | **ADR 009 amendment: versioning model change** — commit `d16650b` moved sibling requires from v0.0.0 sentinels to real `v0.37.0` pins + replace. Document the rationale and the new release-time re-bump step in `docs/RELEASE_CHECKLIST.md`.                         | 30 min | Open   |
