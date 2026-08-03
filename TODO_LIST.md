# TODO_LIST.md — go-output

**Last updated:** 2026-08-04
**Open items:** 15
**Blocked:** 0

---

## Open Items

### P0 — CI is broken (red on every push since July 6)

| #   | Task                                                                             | Effort | Status |
| --- | -------------------------------------------------------------------------------- | ------ | ------ |
| 1   | **Fix TUI test deadlock** — teatest E2E test hangs in CI (passes locally)        | Medium | Open   |
| 2   | **Fix art-dupl CI installation** — dedup tool fails to install in GitHub Actions | Low    | Open   |

### P1 — Release hygiene & security

| #   | Task                                                                                                                                    | Effort | Status |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 3   | **Retract v0.34.0 tag** — superseded by v0.35.0 (code-identical but stale tree); `retract` directive missing from go.mod                | Low    | Open   |
| 4   | **Root-cause bogus-tag creator** — what process created v0.32.1/v0.33.0 on a stale June commit? Asked 3× across reports, never answered | Medium | Open   |
| 5   | **Create GitHub Releases for v0.34.0–v0.36.0** — tags pushed but no release notes on GitHub                                             | Low    | Open   |
| 6   | **Address 8 dependabot vulnerabilities** — 3 high, 3 moderate, 2 low (GitHub alerts)                                                    | Medium | Open   |
| 7   | **Pin GitHub Actions to commit SHAs** — `actions/checkout@v4` and `actions/setup-go@v5` are mutable tag refs                            | Low    | Open   |

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
