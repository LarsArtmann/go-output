# TODO_LIST.md — go-output

**Last updated:** 2026-08-04
**Open items:** 6
**Blocked:** 0

---

## Resolved Items (2026-08-04 session)

| #  | Task                                                                                            | Resolution                                                                                                                                                                                             |
| -- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1  | **Fix TUI test deadlock** — teatest E2E test hangs in CI                                        | **Fixed.** `vtScreenFromBytes` was called inside the `teatest.WaitFor` polling loop, creating 100+ VT emulators under `-race`. Refactored to use ANSI-strip for polling, single VT reconstruction after. |
| 2  | **Fix art-dupl CI installation** — dedup tool fails to install                                  | **Fixed.** Pinned to `@v0.6.0` (v0.6.1 has a broken build: `printer/html.go` references undefined `cloneGroupFull` etc.).                                                                             |
| 3  | **Retract v0.34.0 tag** — superseded by v0.35.0                                                 | **Fixed.** Added `retract v0.34.0` to `go.mod` (stale tag drift: sibling dep versions misaligned at tag time, fixed same day by v0.35.0).                                                             |
| 4  | **Root-cause bogus-tag creator** — what created v0.32.1/v0.33.0?                                | **Resolved.** No automation in this repo creates git tags (release.yml fires *on* tag push). Tags were created manually during a session, erroneously pointing at stale commit `194441b`. No process to stop — retract directives are the permanent mitigation. |
| 5  | **Create GitHub Releases for v0.34.0–v0.36.0**                                                  | **Resolved.** v0.34.0 is now retracted (no release needed). v0.35.0 and v0.36.0 already have GitHub Releases.                                                                                         |
| 7  | **Pin GitHub Actions to commit SHAs** — mutable tag refs                                        | **Fixed.** All actions in ci.yml + release.yml pinned to commit SHAs with `# vN` comments: checkout, setup-go, golangci-lint-action, softprops/action-gh-release.                                     |

---

## Open Items

### P1 — Release hygiene & security

| #   | Task                                                                                                 | Effort | Status |
| --- | ---------------------------------------------------------------------------------------------------- | ------ | ------ |
| 6   | **Address 10 dependabot vulnerabilities** — all in `website/` npm deps (astro, vite, sharp, etc.)   | Medium | Open   |

### P2 — Error system polish

| #   | Task                                                                                                                    | Effort | Status |
| --- | ----------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 8   | **Migrate d2 module from sentinels to typed error structs** — root/graph/nom use typed structs; d2 still uses sentinels | Medium | Open   |
| 9   | **Add cross-module error integration test** — verify `errors.Is`/`errors.AsType` works across module boundaries         | Low    | Open   |

### P3 — Code quality & documentation

| #   | Task                                                                                                                                         | Effort | Status |
| --- | -------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 10  | **Fix ADR numbering collision** — `011-status-registry.md` and `0011-api-stability-tiers.md` both claim ADR 011; renumber one                | Low    | Open   |
| 11  | **Add happy-path `err == nil` assertion tests for all CQRS `WriteXxx` functions** — gap identified in post-v0.30.0 self-review               | Medium | Open   |
| 12  | **Push 7 consumer repos with unpushed commits** — go-wizard-sdk, index, projects-management-automation, etc. (committed locally, not pushed) | Low    | Open   |
| 13  | **Add `Flush()` call to TUI shutdown path** — TimingCache.Flush() exists but TUI never calls it on clean shutdown                            | Low    | Open   |

### Community (owner-dependent)

| #   | Task                                                                                      | Effort | Status                        |
| --- | ----------------------------------------------------------------------------------------- | ------ | ----------------------------- |
| 14  | **Post to r/golang, submit to Awesome Go**                                                | Low    | Open (needs owner account)    |
| 15  | **Cut `v1.0.0` tag** — API frozen (ADR 006); all v0.30.x–v0.36.x breaking changes shipped | Low    | Prepared — awaiting owner tag |
