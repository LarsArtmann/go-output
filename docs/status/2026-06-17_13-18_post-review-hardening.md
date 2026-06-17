# Comprehensive Status Report: Post-Review Hardening Pass

**Date:** 2026-06-17 13:18  
**Scope:** Full `go-output` project — nom/, tui/, root, serialization, docs  
**Branch:** master (in sync with origin/master)  
**Last commit:** `6c91833` docs: update AGENTS.md coverage table and add CHANGELOG entries  
**Session commits:** 12 total (`cc140cb` → `6c91833`), spanning two phases

---

## a) FULLY DONE

### Phase 1: Top-25 Execution Sprint (commits `cc140cb` → `bf54893`)

8 commits implementing 17 items from the Pareto-ranked plan:

1. **TOML round-trip fix** — Wrapped `[]map[string]string` under `[[row]]` array-of-tables key; unblocked `nix run .#test`.
2. **TUI data races eliminated** — `newTestReporter()` helper prevents real Bubble Tea program from starting in tests; all 18 race failures resolved.
3. **Hand-rolled ANSI parser replaced** — `StripANSI`/`VisibleWidth`/`TruncateVisible` now delegate to `charmbracelet/x/ansi`; net-removed 128 lines of custom scanner code; eliminated direct `runewidth` dependency and `unicode/utf8` depguard entry.
4. **`RenderNode` signature made honest** — Ghost `_ []*TreeNode` parameter named `visibleNodes` and documented.
5. **`go test -race` CI gate** — New `nix run .#test-race` app for nom + tui concurrency modules.
6. **`slog` replaces `log.Printf`** — TUI lifecycle uses structured logging.
7. **Channel-based test rendezvous** — `renderNotify` channel replaces `time.Sleep(50ms)` in refresh test.
8. **Render benchmarks** — `BenchmarkDependencyTree_RenderWithWidth` and `_ChildPriority` added.
9. **iTerm2 synchronized updates** — `\x1b[?2026h/l` wrapping for flicker-free redraws.
10. **Completion percentage** — `(N%)` in summary bar where N = (completed + failed) / total.
11. **`ActivityStatus.Interest()` method** — Priority ordering is now a property of the enum type.
12. **Completed-subtree collapsing** — `elideCompletedUnderPressure()` prioritizes active work under height limits.
13. **E2E smoke test** — `TestInlineRenderer_EndToEnd_Lifecycle` exercises full Start → Finish lifecycle.
14. **`Finish()` race fix** — Now calls `Stop()` before accessing the tree.
15. **`FORMAT_ARCHITECTURE.md` refreshed** — Updated from 12-format to 16-format matrix.
16. **Mouse-click math documented** — Comment explains truncation prevents wrapping.
17. **Examples audited + x/term investigation** — Both clean; decided to keep `golang.org/x/term`.

### Phase 2: Self-Review Hardening Pass (commits `e7a8fcf` → `6c91833`)

3 commits fixing issues found during brutal self-review:

18. **Dead code removed** — `elideCompletedUnderPressure` had unreachable guard `visibleCount+maxHeight <= 0`; simplified to `maxHeight <= 0` only.
19. **Dedicated priority tests** — `tree_priority_test.go`: 6 table-driven cases for `elideCompletedUnderPressure` + ordering invariant tests for `ActivityStatus.Interest()`.
20. **TUI coverage gaps filled** — `updateWorkflowCompletionState()` (was 0%) tested with 3 sub-cases (failed→Errored, completed→Completed, running→Running). `Stop()` nil-program path tested. TUI coverage: 86.6%→88.8%.
21. **AGENTS.md coverage table refreshed** — All 13 entries updated with verified numbers from current codebase.
22. **CHANGELOG `[Unreleased]` populated** — Comprehensive entries for all user-visible changes (Added/Changed/Fixed sections per Keep a Changelog format).
23. **`go mod tidy` run** across all 16 modules.

### Verification (all green)

| Check                 | Result                                   |
| --------------------- | ---------------------------------------- |
| `nix run .#lint`      | **16/16 modules** — 0 issues each        |
| `nix run .#test`      | **15/15 modules** — all pass, 0 failures |
| `nix run .#test-race` | **2/2 modules** (nom, tui) — 0 races     |

### Coverage (verified 2026-06-17)

| Module        | Coverage | Previous (AGENTS.md) |
| ------------- | -------- | -------------------- |
| root          | 96.5%    | 96.3%                |
| delimited     | 91.7%    | 90.2%                |
| serialization | 91.0%    | 91.4%                |
| markup        | 94.3%    | 93.9%                |
| d2            | 97.0%    | 100%                 |
| graph         | 96.5%    | 96.0%                |
| enum          | 100%     | 100%                 |
| escape        | 100%     | 100%                 |
| table         | 98.6%    | 100%                 |
| testhelpers   | 90.7%    | 91.3%                |
| nom           | 92.3%    | 93.1%                |
| tui           | 88.8%    | 84.2%                |
| integration   | 95.5%    | 95.5%                |

### Benchmark baseline (2026-06-17)

```
BenchmarkDependencyTree_Render/100Nodes          20922    57108 ns/op    24088 B/op    514 allocs/op
BenchmarkDependencyTree_Render/500Nodes           4311   278711 ns/op   114202 B/op   2518 allocs/op
BenchmarkDependencyTree_VisibleNodes/100Nodes   326607     3483 ns/op     9976 B/op    121 allocs/op
BenchmarkDependencyTree_VisibleNodes/500Nodes    69181    15664 ns/op    44280 B/op    563 allocs/op
BenchmarkDependencyTree_RenderWithWidth/100Nodes 15664    76350 ns/op    24088 B/op    514 allocs/op
BenchmarkDependencyTree_RenderWithWidth/500Nodes  3264   379887 ns/op   114206 B/op   2518 allocs/op
BenchmarkDependencyTree_ChildPriority           786171     2595 ns/op      952 B/op      3 allocs/op
```

---

## b) PARTIALLY DONE

1. **TUI coverage at 88.8%** — Up from 84.2% (pre-session) and 86.6% (mid-session), but still below the 90% target. The remaining gap is mostly in `ensureStarted()` (25%) and `lifecycle.go:Stop()` (0% on the `pr.program != nil` branch) — both require a real terminal program to test.
2. **`nom/` coverage at 92.3%** — Above 90% target but slightly below the 93.1% reported in the previous AGENTS.md. The subtree-collapsing feature and `elideCompletedUnderPressure` added new code paths; the dedicated tests in `tree_priority_test.go` now directly exercise the elision logic.
3. **`renderNotify` is an unexported test hook on a production struct** — Pragmatic but slightly pollutes the type. An alternative would be dependency injection of a notification callback, but that adds complexity for minimal gain.
4. **TODO_LIST.md is stale** — Still says "Last updated: 2026-06-15" and references "17 modules" (now 16). The open items (nom internal packaging, naming renames, v1.0.0 tag) are still valid but haven't been re-verified against current code.
5. **No `.github/workflows/ci.yml`** — CI relies entirely on pre-commit hooks and manual `nix run` commands. The `test-race` app exists but is not automated.

---

## c) NOT STARTED

1. **v1.0.0 release tag** — API frozen per ADR 006 but still at v0.10.x. Blocked on owner decision about `TableData` field-vs-getter API (TODO_LIST #15).
2. **Community outreach** — r/golang post, Awesome Go submission (TODO_LIST #14).
3. **`nom/` internal/ packaging** — 25 production files could benefit from sub-packaging (TODO_LIST #2).
4. **`Marshaler` → `Renderer` terminology** — Blocked by ADR 006 API freeze (TODO_LIST #12).
5. **`RenderOptions.GraphID` dead code** — Documented but not removed; blocked by ADR 006 (TODO_LIST #13).
6. **Cache `TreeNode.Depth()`** — Currently O(n) parent-chain walk (TODO_LIST #10).
7. **Bounds validation for `D2NodeStyle.Opacity`** (TODO_LIST #11).
8. **`ColorModeAuto.ShouldColor()` deterministic testing** (TODO_LIST #9).
9. **Naming renames** — `GetOperationSymbol` → `OperationSymbol`, `HandleError` → `Must` (TODO_LIST #7, #8).
10. **GitHub Actions CI workflow** — No `.github/workflows/` directory exists.

---

## d) TOTALLY FUCKED UP!

1. **Nothing is currently broken.** All 16 modules pass lint with zero issues. All 15 test suites pass. Both concurrency-sensitive modules pass `go test -race`. The working tree has one minor uncommitted change (CHANGELOG formatting by a pre-commit hook).
2. **In the first phase I committed a bug** — `elideCompletedUnderPressure` had an unreachable guard clause (`visibleCount+maxHeight <= 0`). This was dead code, not a runtime bug, but it was sloppy. I caught and fixed it during the self-review pass.
3. **In the first phase I forgot docs hygiene** — CHANGELOG `[Unreleased]` was empty, AGENTS.md coverage table was stale, and no dedicated tests existed for `elideCompletedUnderPressure` or `Interest()`. All fixed in the second pass.
4. **Coverage in d2/ and table/ dropped from 100%** — d2 went from 100% to 97.0%, table from 100% to 98.6%. This appears to be from earlier refactoring work, not this session. The AGENTS.md now reflects the correct numbers.

---

## e) WHAT WE SHOULD IMPROVE!

### Architecture / type models

1. **Introduce a `DisplayConfig` value object** for `nom/` that bundles `maxHeight`, `maxWidth`, `colorMode`, and `noColor` — currently passed as ad-hoc parameters. Makes impossible states unrepresentable (e.g., negative width).
2. **Add a `Snapshot()` method on `DependencyTree`** — returns an immutable copy of visible nodes under a read lock. Eliminates the need for `Finish()` to call `Stop()` before rendering.
3. **Cache `TreeNode.Depth()`** — Currently O(n) parent-chain walk per call. Store depth as a field set during `Build()`.
4. **Consider `nom/` internal/ packaging** — 25 production files in a flat package. Sub-packages (`internal/tree/`, `internal/render/`, `internal/cache/`) would improve locality.

### Code quality

5. **Add `.github/workflows/ci.yml`** — Automate `nix run .#lint`, `nix run .#test`, and `nix run .#test-race` on every push/PR. Currently no CI exists.
6. **Add fuzz tests for `TruncateVisible` and `PhysicalLineCount`** — Width-calculation hot paths only tested with hand-picked inputs.
7. **Update `TODO_LIST.md`** — Stale since 2026-06-15; references "17 modules" (now 16); items need re-verification against current code.
8. **Add `CLICOLOR=0` / `CLICOLOR_FORCE` detection** to `detectNoColor()` — Standard env vars for color control.

### Product / UX

9. **Add a visual progress bar** — `[████░░░░░░] 40%` alongside the text percentage. More scannable than `(40%)`.
10. **Add ETA estimate to summary bar** — "ETA 2m 30s" leveraging the existing `TimingCache` median data.

---

## f) Top #25 things we should get done next

Sorted by **impact ÷ effort** (highest first):

| #   | Task                                                           | Module      | Impact           | Effort  | Why                                           |
| --- | -------------------------------------------------------------- | ----------- | ---------------- | ------- | --------------------------------------------- |
| 1   | Add `.github/workflows/ci.yml` with lint+test+race             | `.github/`  | 🔴 CI gate       | Low     | No CI exists; relies on manual runs           |
| 2   | Commit CHANGELOG formatting fix from pre-commit hook           | root        | 🟢 hygiene       | Trivial | Uncommitted diff in working tree              |
| 3   | Cut `v1.0.0` tag                                               | root        | 🔴 release       | Low     | API frozen (ADR 006); still at v0.10.x        |
| 4   | Update `TODO_LIST.md` — re-verify items, fix "17 modules" → 16 | root        | 🟡 docs          | Low     | Stale since 2026-06-15                        |
| 5   | Cache `TreeNode.Depth()` as field set during `Build()`         | `nom/`      | 🟢 performance   | Low     | Currently O(n) parent-chain walk              |
| 6   | Add fuzz tests for `TruncateVisible` / `PhysicalLineCount`     | `nom/`      | 🟡 robustness    | Low     | Hot paths only tested with hand-picked inputs |
| 7   | Restore tui/ coverage to 90%+ — test View() rendering paths    | `tui/`      | 🟡 quality       | Medium  | 88.8% → 90%+ target                           |
| 8   | Add `CLICOLOR=0` / `CLICOLOR_FORCE` detection                  | `nom/`      | 🟢 compatibility | Low     | Standard color detection env vars             |
| 9   | Add visual progress bar in summary                             | `nom/`      | 🟢 UX            | Medium  | More scannable than text percentage           |
| 10  | Rename `GetOperationSymbol` → `OperationSymbol`                | `nom/`      | 🟢 naming        | Low     | Getter prefix on non-getter (TODO_LIST #7)    |
| 11  | Rename `HandleError` → `Must` in examples                      | `examples/` | 🟢 naming        | Low     | Convention alignment (TODO_LIST #8)           |
| 12  | Add ETA estimate to summary bar                                | `nom/`      | 🟢 UX            | Medium  | Leverage existing timing cache                |
| 13  | `nom/` internal/ packaging — split 25 files into sub-packages  | `nom/`      | 🟡 architecture  | High    | Better locality & navigability                |
| 14  | Add bounds validation for `D2NodeStyle.Opacity`                | `d2/`       | 🟡 safety        | Low     | No validation currently                       |
| 15  | `ColorModeAuto.ShouldColor()` deterministic testing            | root        | 🟡 testability   | Low     | Reads env+TTY at runtime                      |
| 16  | Add `Snapshot()` to `DependencyTree` for lock-free rendering   | `nom/`      | 🟡 architecture  | Medium  | Eliminate lock during render                  |
| 17  | Community: Post to r/golang, submit to Awesome Go              | docs        | 🟢 visibility    | Low     | Project is ready for users                    |
| 18  | Add `nom.DisplayConfig` struct for bundled render config       | `nom/`      | 🟢 API           | Medium  | Better discoverability than ad-hoc params     |
| 19  | Unify `Marshaler` → `Renderer` terminology                     | root        | 🟡 consistency   | Medium  | Blocked by ADR 006; plan for v1               |
| 20  | Wire or remove `RenderOptions.GraphID`                         | root        | 🟢 hygiene       | Low     | Dead code; blocked by ADR 006                 |
| 21  | Add BDD tests for nom inline renderer lifecycle                | `bdd/`      | 🟢 quality       | Medium  | Integration-level behavior verification       |
| 22  | Profile 1000+ node trees for render performance                | `nom/`      | 🟢 performance   | Low     | Current bench is 500 nodes                    |
| 23  | Investigate bubbletea v2 program testing helpers               | `tui/`      | 🟡 testability   | Low     | May eliminate need for newTestReporter hack   |
| 24  | Add `--no-color` flag support to InlineRenderer                | `nom/`      | 🟢 UX            | Low     | SetNoColor exists but not CLI-exposed         |
| 25  | Document `renderNotify` test hook in type doc comment          | `nom/`      | 🟢 docs          | Low     | Unexported field, undocumented in type doc    |

---

## g) Top #1 question I cannot figure out myself

**Should we cut `v1.0.0` now, or address the remaining API inconsistencies first?**

The API is frozen (ADR 006). All tests and lint pass. The race detector is clean. But three TODO_LIST items are blocked by the API freeze:

- **#12** — `Marshaler` vs `Renderer` terminology mismatch (breaking rename)
- **#13** — `RenderOptions.GraphID` is dead code (breaking removal)
- **#15** — `TableData` field-vs-getter decision (affects every consumer)

If we ship v1.0 now, these get locked in until v2. If we fix them first, v1.0 is delayed but ships cleaner. I don't have the product authority to weigh "shipping now with known warts" vs "delaying for cleanliness." Please instruct.

---

## Sprint metrics summary

| Metric                   | Start of session      | End of session        |
| ------------------------ | --------------------- | --------------------- |
| Failing test modules     | 2 (integration, root) | 0                     |
| Data race failures       | 18 in tui             | 0                     |
| Lint issues              | 0                     | 0                     |
| `nom/` coverage          | 93.1% (stale)         | 92.3%                 |
| `tui/` coverage          | 84.2% (stale)         | 88.8%                 |
| Root coverage            | 96.3% (stale)         | 96.5%                 |
| Hand-rolled ANSI code    | ~200 lines            | ~70 (wrappers)        |
| Benchmarks               | 2                     | 4                     |
| Flaky test sleeps        | 1                     | 0                     |
| CI race gate             | None                  | `nix run .#test-race` |
| CHANGELOG entries        | 0 (empty)             | Full [Unreleased]     |
| AGENTS.md coverage stale | Yes                   | Updated               |
| Commits this session     | —                     | 12                    |

---

_Report generated at 2026-06-17 13:18._
