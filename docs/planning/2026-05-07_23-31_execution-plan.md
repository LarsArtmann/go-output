# go-output — Comprehensive Execution Plan

**Created:** 2026-05-07 23:31 CEST
**Status:** ALL TODOs, max 12 min each, sorted by impact ÷ effort
**Source:** Status report + multi-module proposal v3 + reinventing-the-wheel audit + 2026-04-29 improvement plan

---

## Master TODO List — 25 Tasks, Sorted by Impact/Effort

Each task is self-contained, takes ≤12 min, and leaves the project in a working, committable state.

| #      | Task                                                                                      | Effort | Impact      | Category     | Prereq |
| ------ | ----------------------------------------------------------------------------------------- | ------ | ----------- | ------------ | ------ |
| **1**  | Create `go.work` at root with `use .`                                                     | 3 min  | 🔴 Critical | Multi-module | —      |
| **2**  | Add `go.mod` to `enum/` (zero deps)                                                       | 5 min  | 🔴 High     | Multi-module | #1     |
| **3**  | Add `go.mod` to `escape/` (zero deps)                                                     | 5 min  | 🔴 High     | Multi-module | #1     |
| **4**  | Add `go.mod` to `cmdguard/` (zero deps)                                                   | 5 min  | 🔴 High     | Multi-module | #1     |
| **5**  | Update root `go.mod` with `replace` for enum/escape/cmdguard                              | 5 min  | 🔴 High     | Multi-module | #2-4   |
| **6**  | Replace `escape.HTML()` body with `html.EscapeString()`                                   | 5 min  | 🟡 Medium   | Reinvention  | —      |
| **7**  | Replace `escape.XML()` body with stdlib or fix `&apos;`                                   | 5 min  | 🟡 Medium   | Reinvention  | —      |
| **8**  | Add deprecation notice to `sort/sorter.go`                                                | 3 min  | 🟡 Medium   | Reinvention  | —      |
| **9**  | Add deprecation notice to `sort/compare.go`                                               | 2 min  | 🟡 Medium   | Reinvention  | —      |
| **10** | Add `go.mod` to `table/` with replace → root + lipgloss                                   | 8 min  | 🔴 High     | Multi-module | #1     |
| **11** | Audit `SortBy` enum: find all consumers, decide keep/deprecate                            | 10 min | 🟡 Medium   | Reinvention  | #8     |
| **12** | Fix depguard violations: add `examples/shared` to allowlist in `.golangci.yml`            | 5 min  | 🟢 Low      | Quality      | —      |
| **13** | Inline `FilledStrings`, remove `slices.go`                                                | 8 min  | 🟢 Low      | Quality      | —      |
| **14** | Unify `stringEnum` in fuzz_test.go → use `gentest.StringEnum`                             | 5 min  | 🟢 Low      | DRY          | —      |
| **15** | Create `d2/` dir + `go.mod` (depends: root, enum, escape)                                 | 5 min  | 🔴 High     | Multi-module | #1     |
| **16** | Move `d2*.go` files from root → `d2/`, change `package output` → `package d2`             | 10 min | 🔴 High     | Multi-module | #15    |
| **17** | Update d2 file imports: same-package refs → `output "github.com/larsartmann/go-output"`   | 10 min | 🔴 High     | Multi-module | #16    |
| **18** | Move d2 test files to `d2/` + update imports                                              | 10 min | 🟡 Medium   | Multi-module | #17    |
| **19** | Create `graph/` dir + `go.mod` (depends: root, escape)                                    | 5 min  | 🔴 High     | Multi-module | #1     |
| **20** | Move `dot.go`+`mermaid.go` → `graph/`, extract `GraphRendererMixin` into `graph_mixin.go` | 10 min | 🔴 High     | Multi-module | #19    |
| **21** | Update graph file imports: same-package → `output`, change package name                   | 10 min | 🔴 High     | Multi-module | #20    |
| **22** | Move graph test files to `graph/` + update imports                                        | 10 min | 🟡 Medium   | Multi-module | #21    |
| **23** | Write ADR 001: `docs/adr/001-multi-module-split.md`                                       | 10 min | 🟡 Medium   | Docs         | #5     |
| **24** | Update `AGENTS.md` with final multi-module structure                                      | 8 min  | 🟡 Medium   | Docs         | #22    |
| **25** | Update `README.md` with module paths and new examples                                     | 10 min | 🟡 Medium   | Docs         | #24    |

---

## Execution Phases (grouped for commit granularity)

### Phase 1: Quick Wins — Reinvention Fixes (30 min, 5 commits)

```
#6  escape.HTML() → html.EscapeString()          [5 min]  → commit
#7  escape.XML() → stdlib                          [5 min]  → commit
#8  sort/sorter.go deprecation notice              [3 min]  ┐
#9  sort/compare.go deprecation notice             [2 min]  ┤→ 1 commit
#11 SortBy audit: keep (cmdguard_test uses it as example enum) [10 min] ┘
#12 Fix depguard: add examples/shared to allowlist [5 min]  → commit
```

**Verify after each commit:** `go build ./... && go test ./...`

### Phase 2: Leaf Modules (20 min, 2 commits)

```
#1  Create go.work                                [3 min]  ┐
#2  enum/go.mod                                   [5 min]  ┤
#3  escape/go.mod                                 [5 min]  ├→ 1 commit
#4  cmdguard/go.mod                               [5 min]  ┤
#5  Root go.mod replace directives                [5 min]  ┘
#10 table/go.mod with lipgloss + replace          [8 min]  → commit
```

**Verify after each commit:** `go build ./... && go test ./...`

### Phase 3: D2 Module Extraction (35 min, 2 commits)

```
#15 Create d2/ dir + go.mod                       [5 min]  ┐
#16 Move d2*.go → d2/, package output → d2        [10 min] ├→ 1 commit
#17 Fix d2 imports: output.* refs                 [10 min] ┤
#18 Move d2 test files + fix imports              [10 min] ┘→ 1 commit
```

**Verify after each commit:** `go test ./d2/... && go test ./...`

### Phase 4: Graph Module Extraction (35 min, 2 commits)

```
#19 Create graph/ dir + go.mod                    [5 min]  ┐
#20 Move dot.go+mermaid.go → graph/, extract mixin [10 min] ├→ 1 commit
#21 Fix graph imports                             [10 min] ┤
#22 Move graph test files + fix imports           [10 min] ┘→ 1 commit
```

**Verify after each commit:** `go test ./graph/... && go test ./...`

### Phase 5: Code Quality DRY (15 min, 2 commits)

```
#13 Inline FilledStrings, remove slices.go        [8 min]  → commit
#14 Unify stringEnum → gentest.StringEnum         [5 min]  → commit
```

### Phase 6: Documentation (25 min, 1 commit)

```
#23 Write ADR 001                                 [10 min] ┐
#24 Update AGENTS.md                              [8 min]  ├→ 1 commit
#25 Update README.md                              [10 min] ┘
```

**Final verify:** `go build ./... && go test ./... && golangci-lint run ./...`

---

## Total Estimate

| Phase                     | Tasks  | Time         | Commits        |
| ------------------------- | ------ | ------------ | -------------- |
| Phase 1: Quick Wins       | 6      | 30 min       | 3              |
| Phase 2: Leaf Modules     | 6      | 31 min       | 2              |
| Phase 3: D2 Extraction    | 4      | 35 min       | 2              |
| Phase 4: Graph Extraction | 4      | 35 min       | 2              |
| Phase 5: Code Quality     | 2      | 13 min       | 2              |
| Phase 6: Documentation    | 3      | 28 min       | 1              |
| **Total**                 | **25** | **~3 hours** | **12 commits** |

---

## What's NOT in This Plan (future work)

- CI setup (`.github/workflows/ci.yml`)
- Registry ghost system: integrate or remove
- Unify tree conversion (d2/dot/mermaid each have own `addTreeNodes`)
- `color.go` → `termenv` evaluation
- `yaml.go` evaluation (add value or inline)
- Test helper deduplication (output_test_helpers.go vs testutils)
- Remove `format_deprecated.go`
- `integration/` and `examples/` as standalone modules (do after d2/graph/table work)
- PLAN.md cleanup
- Stale docs/status/ cleanup
