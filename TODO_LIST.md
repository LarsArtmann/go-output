# TODO_LIST.md — go-output

**Last updated:** 2026-08-04
**Open items:** 2
**Blocked:** 0

---

## Resolved Items (2026-08-04 session)

| #   | Task                                                                    | Resolution                                                                                                                                                                                                                                                                                                                    |
| --- | ----------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | **Fix TUI test deadlock** — teatest E2E test hangs in CI                | **Fixed.** `teatest.WaitFor` internally calls `io.ReadAll(io.TeeReader(...))` which blocks forever when the test program writes continuously to its output buffer (surfaces under `-race` at 10-min timeout). Replaced with `pollTeatestOutput` helper using bounded `var buf [8192]byte` reads + hard deadline. Test now passes in 0.06s.                                                                                                                      |
| 2   | **Fix art-dupl CI installation** — dedup tool fails to install          | **Fixed.** Root cause: generated file `report_templ.go` was globally gitignored (`*_templ.go` in `~/.config/git/ignore`), so all released versions (v0.5.0–v0.6.1) shipped without it. Fixed at source: force-added the file in art-dupl repo, removed gitignore rule, tagged **v0.6.2**. CI pinned to `@v0.6.2`.                                                                                                     |
| 3   | **Retract v0.34.0 tag** — superseded by v0.35.0                         | **Fixed.** Added `retract v0.34.0` to `go.mod` (stale tag drift: sibling dep versions misaligned at tag time, fixed same day by v0.35.0).                                                                                                                                                                                     |
| 4   | **Root-cause bogus-tag creator** — what created v0.32.1/v0.33.0?        | **Resolved.** No automation in this repo creates git tags (release.yml fires _on_ tag push). Tags were created manually during a session, erroneously pointing at stale commit `194441b`. No process to stop — retract directives are the permanent mitigation.                                                               |
| 5   | **Create GitHub Releases for v0.34.0–v0.36.0**                          | **Resolved.** v0.34.0 is now retracted (no release needed). v0.35.0 and v0.36.0 already have GitHub Releases.                                                                                                                                                                                                                 |
| 6   | **Address dependabot vulnerabilities** — 10 alerts in website/ npm deps | **Fixed.** Upgraded astro v6→v7.1.6, starlight v0.39→v0.41, astro-og-canvas v0.11→v0.13. Removed stale vite override (astro v7 uses vite v8). Added esbuild override 0.28.1. 0 vulnerabilities.                                                                                                                               |
| 7   | **Pin GitHub Actions to commit SHAs** — mutable tag refs                | **Fixed.** All actions in ci.yml + release.yml pinned to commit SHAs with `# vN` comments: checkout, setup-go, golangci-lint-action, softprops/action-gh-release.                                                                                                                                                             |
| 8   | **Migrate d2 module from sentinels to typed error structs**             | **Fixed.** Replaced 5 sentinel errors with typed structs (`InvalidDirectionError`, `InvalidNodeShapeError`, `InvalidArrowTypeError`, `InvalidConstraintError`, `InvalidTextTransformError`) following the root/graph pattern. Added `AllDirections()`, `AllNodeShapes()`, `AllArrowTypes()`, `AllTextTransforms()` functions. |
| 9   | **Add cross-module error integration test**                             | **Fixed.** Added `integration/cross_module_error_test.go` with 5 test groups: root sentinels via `errors.Is`, root typed errors via `errors.AsType`, sub-module typed errors (d2+graph), error distinctness, and deep wrapping preservation.                                                                                  |
| 10  | **Fix ADR numbering collision**                                         | **Fixed.** ADR file was already renamed `0011-api-stability-tiers.md` → `0014-api-stability-tiers.md` in a prior session. Fixed remaining stale reference in AGENTS.md (`ADR 0011` → `ADR 014`) and updated the d2 error-system description (sentinels → typed structs).                                                      |
| 11  | **Add happy-path tests for all CQRS `WriteXxx` functions**              | **Fixed.** Added `TestWriteMermaid_PureFunction` (graph), `TestWriteGraph_HappyPath` + `TestWrite_HappyPath` (d2), and `TestCQRS_WriteMarkdown_HappyPath` (tree). All 14 CQRS WriteXxx functions now have happy-path `err == nil` coverage.                                                                                   |
| 12  | **Add `Flush()` call to TUI shutdown path**                             | **Fixed.** `BubbleTeaProgressReporter.Stop()` now calls `nomSubscriber.Flush()` after the 100% progress update and before the quit signal, ensuring timing-cache data is persisted to disk on clean shutdown. Added `TestBubbleTeaProgressReporter_Stop_FlushesTimingCache`.                                                  |

---

## Open Items

### Community (owner-dependent)

| #   | Task                                                                                      | Effort | Status                        |
| --- | ----------------------------------------------------------------------------------------- | ------ | ----------------------------- |
| 14  | **Post to r/golang, submit to Awesome Go**                                                | Low    | Open (needs owner account)    |
| 15  | **Cut `v1.0.0` tag** — API frozen (ADR 006); all v0.30.x–v0.36.x breaking changes shipped | Low    | Prepared — awaiting owner tag |
