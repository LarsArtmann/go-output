# Status Report: go-output — 2026-04-30 Session 2

**Date:** 2026-04-30 05:02 CEST
**Session Focus:** Adoption review action items — escape export, sort refactoring, ByField helper
**Previous Report:** `2026-04-30_04-53_adoption-review-action-items.md`

---

## a) Fully Done

| #   | Task                                                   | Commit    | Files                                         | Impact                                                                        |
| --- | ------------------------------------------------------ | --------- | --------------------------------------------- | ----------------------------------------------------------------------------- |
| 1   | Export `escape/` from `internal/` to public            | `9e93afa` | 12 files                                      | External consumers can now `import "github.com/larsartmann/go-output/escape"` |
| 2   | Update all 9 internal imports                          | `9e93afa` | d2, html, xml, streaming, dot, mermaid, tests | Zero breakage, git detected rename                                            |
| 3   | Fix depguard config for `escape/` and remove `reflect` | `9e93afa` | `.golangci.yml`                               | Linter passes with 0 issues                                                   |
| 4   | Remove reflection from `sort.Sorter[T]`                | `2a8e841` | 4 files                                       | Deleted 70+ lines of reflection code, zero reflection overhead                |
| 5   | Add `ByField[T, F cmp.Ordered]` helper                 | `64b7493` | 5 files                                       | Ergonomic: `sort.ByField(func(p Project) string { return p.Name })`           |
| 6   | Update all callers to use `ByField`                    | `64b7493` | sort tests, userjourney, integration          | Consistent API usage                                                          |
| 7   | Update README with escape fix + sort example           | `e7d908a` | README.md                                     | Docs match reality, shows idiomatic usage                                     |
| 8   | Write status report                                    | `80f5bbb` | docs/status/                                  | Session documentation                                                         |
| 9   | Update AGENTS.md                                       | `80f5bbb` | AGENTS.md                                     | Memory for future sessions                                                    |
| 10  | Fix PLAN.md stale references                           | `e259a52` | PLAN.md                                       | Removed "internal/escape" and "reflect-based"                                 |
| 11  | Split sort_test.go under 350-line limit                | `e259a52` | sort/sort_test.go, sort/compare_test.go       | 430 → 334 lines                                                               |

## b) Partially Done

Nothing. All tasks started are fully complete.

## c) Not Started

| #   | Task                                                           | Priority | Effort | Notes                                                        |
| --- | -------------------------------------------------------------- | -------- | ------ | ------------------------------------------------------------ |
| 1   | Add `ByTimeField` helper for `time.Time` sorting               | Low      | 15min  | `time.Time` is not `cmp.Ordered`, needs manual `LessFunc`    |
| 2   | Consider `slices.SortFunc` vs `sort.SliceStable`               | Low      | 30min  | `SliceStable` is correct default but slower; benchmark first |
| 3   | Add integration guide to README for external consumers         | Medium   | 30min  | Based on adoption review findings                            |
| 4   | Example: real-world indexer integration                        | Low      | 1h     | Show how indexer would use escape.HTML + table.New()         |
| 5   | Benchmark escape package performance                           | Low      | 20min  | Compare escape.HTML vs html.EscapeString from stdlib         |
| 6   | Consider exporting `mode` type from escape for custom escaping | Low      | 10min  | Only if consumers need custom apostrophe handling            |

## d) Totally Fucked Up

Nothing. All changes compile, test, and lint clean.

## e) What We Should Improve

### Code Quality

1. **sort/sort_test.go:88** — `compareName` takes an unused `_ bool` parameter. This is a leftover from the refactoring where it was simplified to use `ByField` but kept the same signature for compatibility with `compareCount`. Should be cleaned up.

2. **sort/sort_test.go:78** — `compareCount` still uses manual `a.Count < b.Count` instead of `ByField`. Could be simplified to `ByField(func(item testItem) int { return item.Count })` but the ascending/descending logic makes it different enough to keep.

3. **escape/escape.go** — `htmlMode`/`xmlMode` globals require `//nolint:gochecknoglobals` comments. Could be refactored to `const` strings passed directly, eliminating the `mode` struct entirely.

4. **Sorter.By field** — The `SortBy` field on `Sorter` is now purely informational (used for logging/documentation). Could add a `String()` method to `Sorter` that includes `By` for debugging.

### Architecture

5. **sort.Sorter vs stdlib** — Now that Sorter requires an explicit `LessFunc`, it's essentially `sort.SliceStable` with a builder pattern. The `SortBy` field provides no functional value. Consider whether the Sorter abstraction still earns its keep vs just recommending `slices.SortFunc`.

6. **escape package scope** — Currently handles HTML, XML, D2, DOT, and Mermaid. Could be split into format-specific sub-packages (e.g., `escape/html.go`, `escape/d2.go`) as it grows, but at 104 lines it doesn't need it yet.

7. **ByField returns `func(a, b T) bool`** — Could return a named type `LessFunc[T]` for better documentation and extensibility.

### Type Model Improvements

8. **BrandedID could use `cmp.Ordered` constraint** — Currently `BrandedID[Brand]` wraps `string`. If the underlying type were `cmp.Ordered`, `ByField` could work directly with branded IDs.

9. **TableData could use generics** — Currently `TableData` uses `[][]string`. A generic `TableData[T]` with a `RenderFunc(T) string` could eliminate string conversion boilerplate.

10. **Sorter could use `slices.SortStableFunc`** from Go 1.21+ instead of `sort.SliceStable` — slightly cleaner API, same semantics.

## f) Top 25 Things We Should Get Done Next

Sorted by **impact × effort** (high impact + low effort first):

| #   | Task                                                              | Impact | Effort | Category                  |
| --- | ----------------------------------------------------------------- | ------ | ------ | ------------------------- |
| 1   | Clean up `compareName` unused `_ bool` param in sort_test.go      | Low    | 2min   | Code hygiene              |
| 2   | Remove `mode` struct from escape.go, use inline `const` strings   | Low    | 10min  | Code simplicity           |
| 3   | Add doc comment examples to `sort.ByField` (`// Example:`)        | Medium | 5min   | Documentation             |
| 4   | Add `String()` method to `Sorter` for debugging                   | Low    | 5min   | Developer experience      |
| 5   | Add `ByTimeField` helper for `time.Time`                          | Medium | 15min  | API completeness          |
| 6   | Refactor escape globals to eliminate nolint comments              | Low    | 10min  | Code cleanliness          |
| 7   | Benchmark escape.HTML vs stdlib html.EscapeString                 | Medium | 20min  | Performance validation    |
| 8   | Add fuzz tests for escape package                                 | Medium | 20min  | Security hardening        |
| 9   | Replace `sort.SliceStable` with `slices.SortStableFunc`           | Low    | 15min  | Modern stdlib usage       |
| 10  | Named `LessFunc[T]` type instead of raw `func(a, b T) bool`       | Medium | 20min  | Type safety               |
| 11  | Add integration guide section to README                           | Medium | 30min  | Adoption enablement       |
| 12  | Add `.golangci.yml` lint rule for 350-line file limit enforcement | Medium | 10min  | Prevent future violations |
| 13  | Add CI check for file size limits                                 | Medium | 15min  | Quality gate              |
| 14  | Review all docs/status/ for stale references                      | Low    | 15min  | Documentation hygiene     |
| 15  | Add `Sorter.MustSort()` that panics if LessFunc is nil            | Low    | 10min  | API ergonomics            |
| 16  | TableData generic `TableData[T]` with renderer                    | High   | 2h     | Architecture              |
| 17  | BrandedID with `cmp.Ordered` constraint                           | High   | 1h     | Type model                |
| 18  | Escape package benchmarks (all 7 functions)                       | Medium | 30min  | Performance               |
| 19  | Add real-world example: CLI table output with sorting             | Medium | 30min  | Documentation             |
| 20  | Consider `escape.HTMLEntity(string) string` for named entities    | Low    | 20min  | Feature                   |
| 21  | Add `RenderAllFormats()` convenience function                     | Medium | 30min  | Feature                   |
| 22  | Investigate streaming writer for CSV/TSV (not just HTML)          | Medium | 1h     | Feature                   |
| 23  | Add Go doc examples (`func ExampleByField()`)                     | Medium | 20min  | Documentation             |
| 24  | Consider `sort.Desc(field)` shorthand                             | Low    | 15min  | API sugar                 |
| 25  | Propose go-output as dependency in indexer project                | High   | 2h     | Adoption                  |

## g) Top #1 Question I Cannot Figure Out Myself

**Should `sort.Sorter` exist at all now that it requires explicit `LessFunc`?**

The `Sorter` is now a thin wrapper around `sort.SliceStable` with a builder pattern and a `SortBy` field that has no functional effect. The stdlib `slices.SortFunc` provides the same capability more directly:

```go
// Current go-output Sorter
sort.New(items, output.SortByName, false).
    WithLessFunc(sort.ByField(func(p Project) string { return p.Name })).
    Sort()

// vs stdlib directly
slices.SortFunc(items, cmp.Less[Project, string](func(p Project) string { return p.Name }))
```

The Sorter still adds value via:

1. **Builder pattern** — chaining with `.WithLessFunc().Sort()`
2. **`Desc` flag** — single boolean for direction reversal
3. **`SortBy` metadata** — informational, for logging/debugging

But the adoption review said the indexer already uses `slices.SortFunc` directly and skipped go-output's sort. **Should we deprecate the `sort` package entirely, keep it as a thin convenience, or invest in making it genuinely more useful than stdlib?**

This is a product decision that depends on whether other consumers find the builder pattern valuable enough to justify the dependency.

---

## Quality Gates

| Gate                        | Result                                                   |
| --------------------------- | -------------------------------------------------------- |
| `go test -count=1 ./...`    | All pass                                                 |
| `go test -cover ./...`      | 90.1% (root), 100% (cmdguard, enum, escape, sort, table) |
| `golangci-lint run ./...`   | 0 issues                                                 |
| `go vet ./...`              | Clean                                                    |
| `go build ./...`            | Clean                                                    |
| File size limit (350 lines) | All under                                                |

## Session Metrics

| Metric                | Value                                                                              |
| --------------------- | ---------------------------------------------------------------------------------- |
| Total commits         | 6                                                                                  |
| Net lines changed     | -270 (code reduced)                                                                |
| Files changed         | 20                                                                                 |
| New files             | 3 (escape/escape.go, escape/escape_test.go, sort/compare.go, sort/compare_test.go) |
| Deleted files         | 2 (internal/escape/escape.go, internal/escape/escape_test.go)                      |
| Reflection eliminated | 100%                                                                               |
| New public APIs       | 2 (`escape.HTML/XML/...`, `sort.ByField`)                                          |
| Zero new dependencies | Confirmed                                                                          |
