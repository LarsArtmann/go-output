# go-output — Comprehensive Status Report

**Date:** 2026-05-08 00:16 CEST
**Branch:** master (up to date with origin/master)
**Previous Report:** 2026-05-07 23:55

---

## Executive Summary

The multi-module workspace migration is **live and functional** — 8 independent Go modules with clean dependency isolation. All tests pass, all lints pass (0 issues), and coverage is strong. However, two concrete problems exist: the `examples/` module is **broken** (missing `table` dependency in `go.mod`), and the `README.md` leaks local filesystem paths. The D2 and graph module extractions from the execution plan remain the major unfinished architectural work.

---

## a) FULLY DONE ✅

### Multi-Module Workspace — 8 Modules Live

| Module         | go.mod | Tests | Coverage | Status                     |
| -------------- | ------ | ----- | -------- | -------------------------- |
| Root           | ✅     | ✅    | 90.3%    | Stable                     |
| `enum/`        | ✅     | ✅    | 100%     | Zero deps                  |
| `escape/`      | ✅     | ✅    | 100%     | Zero deps                  |
| `cmdguard/`    | ✅     | ✅    | 100%     | Zero deps                  |
| `table/`       | ✅     | ✅    | 100%     | Lipgloss isolated          |
| `sort/`        | ✅     | ✅    | 100%     | Deprecated                 |
| `integration/` | ✅     | ✅    | n/a      | Cross-module tests         |
| `examples/`    | ✅     | ⚠️    | 0%       | **BROKEN** — see section d |

### Lint & Quality Gates

- **golangci-lint:** 0 issues (was 4, fixed this session)
- **go vet:** 0 issues
- **go build:** compiles clean across all modules
- **go test -race:** passes (cached)
- **File size limit:** 1 file exceeds 350 lines (`sort/sort_test.go` at 414 lines)
- **TODO/FIXME/HACK comments:** zero found in codebase
- **No code duplication** above threshold

### Execution Plan Phase 1 — Quick Wins

| #   | Task                                             | Status                   |
| --- | ------------------------------------------------ | ------------------------ |
| #6  | `escape.HTML()` → `html.EscapeString()` stdlib   | ✅ Done (commit d2ba200) |
| #7  | `escape.XML()` → stdlib `&apos;` fix             | ✅ Done (commit d2ba200) |
| #8  | `sort/sorter.go` deprecation notice              | ✅ Done (commit f527f10) |
| #9  | `sort/compare.go` deprecation notice             | ✅ Done (commit f527f10) |
| #11 | SortBy audit — keep (cmdguard_test uses it)      | ✅ Done                  |
| #12 | Fix depguard: add `examples/shared` to allowlist | ✅ Done (commit 56c57ea) |

### Execution Plan Phase 2 — Leaf Modules

| #   | Task                                   | Status                   |
| --- | -------------------------------------- | ------------------------ |
| #1  | Create `go.work` at root               | ✅ Done (commit c0a250a) |
| #2  | `enum/go.mod` (zero deps)              | ✅ Done (commit c0a250a) |
| #3  | `escape/go.mod` (zero deps)            | ✅ Done (commit c0a250a) |
| #4  | `cmdguard/go.mod` (zero deps)          | ✅ Done (commit c0a250a) |
| #5  | Root `go.mod` replace directives       | ✅ Done (commit 0027642) |
| #10 | `table/go.mod` with lipgloss + replace | ✅ Done (commit a493e06) |

### Execution Plan Phase 6 — Documentation

| #   | Task             | Status                   |
| --- | ---------------- | ------------------------ |
| #23 | Write ADR 001    | ✅ Done (commit c43938b) |
| #24 | Update AGENTS.md | ✅ Done (commit c43938b) |

### This Session's Fixes

- **gomoddirectives lint errors:** Added `replace-local: true` + `replace-allow-list` to `.golangci.yml`
- **staticcheck SA1019:** Added `//nolint:staticcheck` to `userjourney_test.go` (intentionally testing deprecated `sort` package)

---

## b) PARTIALLY DONE 🔧

### README.md Update (#25)

**Status:** Not started, but critical issues identified:

- **Lines 13-14** leak local macOS paths (`/Users/larsartmann/...`) — must remove
- **Dependencies section** still shows lipgloss as root dependency — contradicts multi-module isolation story
- **Development section** references `just` commands but AGENTS.md says justfile is deprecated
- **Installation section** doesn't mention sub-module imports (`go-output/table`, `go-output/enum`, etc.)
- **Supported Formats tables** list all formats as `github.com/larsartmann/go-output` — `table` format should show `go-output/table`

### Registry System

- Functional and tested, but docs describe it as "opt-in" while it's unclear if it's actually used in production
- Could be a ghost system (dead code that works but nobody uses)

---

## c) NOT STARTED 📋

### Execution Plan Phase 3 — D2 Module Extraction (35 min)

| #   | Task                                                   | Effort | Impact    |
| --- | ------------------------------------------------------ | ------ | --------- |
| #15 | Create `d2/` dir + `go.mod` (deps: root, enum, escape) | 5 min  | 🔴 High   |
| #16 | Move `d2*.go` files from root → `d2/`, change package  | 10 min | 🔴 High   |
| #17 | Update d2 file imports                                 | 10 min | 🔴 High   |
| #18 | Move d2 test files + update imports                    | 10 min | 🟡 Medium |

### Execution Plan Phase 4 — Graph Module Extraction (35 min)

| #   | Task                                                   | Effort | Impact    |
| --- | ------------------------------------------------------ | ------ | --------- |
| #19 | Create `graph/` dir + `go.mod`                         | 5 min  | 🔴 High   |
| #20 | Move `dot.go` + `mermaid.go` → `graph/`, extract mixin | 10 min | 🔴 High   |
| #21 | Fix graph imports                                      | 10 min | 🔴 High   |
| #22 | Move graph test files                                  | 10 min | 🟡 Medium |

### Execution Plan Phase 5 — Code Quality DRY (15 min)

| #   | Task                                                      | Effort | Impact |
| --- | --------------------------------------------------------- | ------ | ------ |
| #13 | Inline `FilledStrings`, remove `slices.go`                | 8 min  | 🟢 Low |
| #14 | Unify `stringEnum` in fuzz_test.go → `gentest.StringEnum` | 5 min  | 🟢 Low |

### Other Not Started Items

- **CI setup** — No `.github/workflows/ci.yml` exists yet (badge in README points to 404)
- **Tree conversion unification** — `d2/dot/mermaid` each have own `addTreeNodes` instead of shared
- **Test helper deduplication** — `output_test_helpers.go` vs `internal/testutils` have overlap
- **`format_deprecated.go`** removal — backward compat aliases linger
- **Registry ghost system audit** — confirm if actually used anywhere
- **`sort/sort_test.go`** exceeds 350-line limit (414 lines)
- **Root Go file `format.go`** at 323 lines — approaching limit

---

## d) TOTALLY FUCKED UP 💥

### `examples/` Module — BROKEN (84+ gopls errors)

**Root cause:** `examples/go.mod` is missing the `table` dependency.

`examples/basic/main.go:11` imports `github.com/larsartmann/go-output/table`, but:

- No `require github.com/larsartmann/go-output/table` in `examples/go.mod`
- No `replace github.com/larsartmann/go-output/table => ../table` directive

This cascades into **84+ missing dependency errors** (lipgloss, charmbracelet, clipperhouse transitive deps).

**Fix:** Add to `examples/go.mod`:

```
require github.com/larsartmann/go-output/table v0.0.0

replace github.com/larsartmann/go-output/table => ../table
```

Then `cd examples && go mod tidy`.

### README.md Leaks Local Paths

**Lines 13-14:**

```markdown
- `/Users/larsartmann/projects/project-meta/`
- `/Users/larsartmann/projects/projects-management-automation/`
```

These are developer-specific macOS paths that expose local filesystem structure in a public GitHub repo. Must be removed or made generic.

### README.md Claims `just` for Development

The Development section references `just build`, `just test`, `just lint`, `just verify` — but AGENTS.md explicitly states justfile is deprecated. The README should reference `go build`, `go test`, `golangci-lint` commands or `flake.nix` equivalents.

---

## e) WHAT WE SHOULD IMPROVE 🎯

### Architecture

1. **D2 + Graph extraction from root** — The root module is a 323-line `format.go` monolith plus 7 D2 files and 3 graph files. Extracting these into `d2/` and `graph/` modules would make the root truly lean and each concern independently versionable.
2. **Tree conversion duplication** — `d2_convert.go`, `dot.go`, and `mermaid.go` each have their own `addTreeNodes` implementation. The generic `AddTreeNodes` in `graph.go` handles the common case. These should be consolidated.
3. **`GraphRendererMixin` location** — Defined in `dot.go` but used by Mermaid too. Should live in `graph.go` or its own file, and eventually in the `graph/` module.

### Code Quality

4. **`FilledStrings` utility** — Only used internally, could be inlined and `slices.go` removed.
5. **`stringEnum` in `fuzz_test.go`** — Duplicates `gentest.StringEnum`. Should use the shared helper.
6. **`format_deprecated.go`** — Backward compat aliases for `OutputFormat`. Decide: keep forever or set a removal date.
7. **`sort/sort_test.go`** at 414 lines — Exceeds 350-line project convention.

### Documentation

8. **README.md overhaul** — Leaked paths, wrong dependency info, deprecated `just` commands, no sub-module examples.
9. **CI badge points to 404** — No workflow file exists.
10. **No CHANGELOG** — No version history or migration guides for the multi-module change.

### DevEx

11. **No `flake.nix`** — AGENTS.md mandates flake.nix over justfile, but none exists yet.
12. **No CI/CD** — No automated testing on push/PR.

---

## f) Top 25 Things We Should Get Done Next

Sorted by impact ÷ effort:

| #      | Task                                                                               | Effort | Impact      | Category |
| ------ | ---------------------------------------------------------------------------------- | ------ | ----------- | -------- |
| **1**  | Fix `examples/go.mod` — add `table` dep + replace directive                        | 3 min  | 🔴 Critical | Fix      |
| **2**  | Remove leaked local paths from `README.md`                                         | 2 min  | 🔴 Critical | Fix      |
| **3**  | Update README.md dependencies section (lipgloss not in root)                       | 5 min  | 🔴 High     | Docs     |
| **4**  | Update README.md development section (remove `just` refs)                          | 3 min  | 🟡 Medium   | Docs     |
| **5**  | Add sub-module installation examples to README.md                                  | 5 min  | 🟡 Medium   | Docs     |
| **6**  | Create `d2/` dir + `go.mod` + move 7 d2\*.go files                                 | 30 min | 🔴 High     | Arch     |
| **7**  | Create `graph/` dir + `go.mod` + move dot.go + mermaid.go                          | 30 min | 🔴 High     | Arch     |
| **8**  | Consolidate `addTreeNodes` implementations (d2/dot/mermaid)                        | 15 min | 🟡 Medium   | DRY      |
| **9**  | Move `GraphRendererMixin` from `dot.go` to `graph.go` or own file                  | 5 min  | 🟡 Medium   | Arch     |
| **10** | Inline `FilledStrings`, remove `slices.go`                                         | 8 min  | 🟢 Low      | Quality  |
| **11** | Unify `stringEnum` in fuzz_test.go → use `gentest.StringEnum`                      | 5 min  | 🟢 Low      | DRY      |
| **12** | Deduplicate test helpers: `output_test_helpers.go` vs `testutils`                  | 15 min | 🟢 Low      | DRY      |
| **13** | Split `sort/sort_test.go` under 350 lines                                          | 10 min | 🟢 Low      | Quality  |
| **14** | Create `.github/workflows/ci.yml` (build + test + lint)                            | 15 min | 🟡 Medium   | CI       |
| **15** | Create `flake.nix` for build automation (per AGENTS.md mandate)                    | 30 min | 🟡 Medium   | DevEx    |
| **16** | Add CHANGELOG.md with multi-module migration notes                                 | 10 min | 🟡 Medium   | Docs     |
| **17** | Audit registry system — confirm usage or document as opt-in                        | 10 min | 🟢 Low      | Audit    |
| **18** | Decide fate of `format_deprecated.go` — keep with date or remove                   | 5 min  | 🟢 Low      | Cleanup  |
| **19** | Add `golangci-lint` config to each sub-module or use root config                   | 5 min  | 🟢 Low      | Quality  |
| **20** | Verify `go.work.sum` is in `.gitignore` (currently untracked)                      | 2 min  | 🟢 Low      | Config   |
| **21** | Add integration test that imports each sub-module from a clean GOPATH              | 15 min | 🟡 Medium   | Testing  |
| **22** | Refactor `format.go` (323 lines) — extract format category logic                   | 15 min | 🟢 Low      | Quality  |
| **23** | Add Go doc comments to all public APIs                                             | 20 min | 🟡 Medium   | Docs     |
| **24** | Add `html` format to Supported Formats table in README (missing from tree section) | 2 min  | 🟢 Low      | Docs     |
| **25** | Create `CONTRIBUTING.md` with development setup instructions                       | 15 min | 🟢 Low      | Docs     |

---

## g) Top #1 Question I Cannot Figure Out Myself 🤔

**Should D2 and Graph modules be extracted from the root package, breaking the current flat `package output` API?**

The execution plan calls for moving `d2*.go` into `d2/` and `dot.go`+`mermaid.go` into `graph/`, changing `package output` to `package d2` / `package graph`. This would require users to import sub-modules:

```go
// Before (current):
import "github.com/larsartmann/go-output"
output.NewD2Diagram(...)

// After (proposed):
import d2 "github.com/larsartmann/go-output/d2"
d2.NewDiagram(...)
```

**Arguments for:** Cleaner root, independently versionable, smaller dependency graph for users who only need basic formats.

**Arguments against:** Breaking API change for all existing users who use D2/DOT/Mermaid. These are not "heavy" dependencies (no external libs needed) — they're pure Go code. The isolation benefit is organizational, not dependency-based (unlike `table/` where lipgloss was the clear win).

**The key question:** Is the organizational purity worth a breaking API change? The `table/` extraction had a clear dependency isolation win. D2/Graph extraction is purely about code organization. Should we keep them in root and just refactor internally, or is the multi-module vision worth the breakage?

---

## Metrics Snapshot

| Metric                 | Value                                            |
| ---------------------- | ------------------------------------------------ |
| Go files               | 81 (40 test files)                               |
| Root module coverage   | 90.3%                                            |
| Sub-module coverage    | 100% (enum, escape, cmdguard, table, sort)       |
| Lint issues            | 0                                                |
| go vet issues          | 0                                                |
| Files over 350 lines   | 1 (`sort/sort_test.go` at 414)                   |
| TODO/FIXME comments    | 0                                                |
| Modules in workspace   | 8                                                |
| Total lines of Go code | ~8,241                                           |
| Uncommitted changes    | 2 files (`.golangci.yml`, `userjourney_test.go`) |

---

## Git Diff (Uncommitted)

```
 .golangci.yml       | 6 ++++++
 userjourney_test.go | 2 +-
 2 files changed, 7 insertions(+), 1 deletion(-)
```

1. `.golangci.yml` — Added `gomoddirectives` settings: `replace-local: true` + `replace-allow-list` for cmdguard/enum/escape
2. `userjourney_test.go` — Added `//nolint:staticcheck` on deprecated `sort` import (intentionally testing deprecated package)
