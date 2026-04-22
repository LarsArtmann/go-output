# go-output — Comprehensive Status Report

**Date:** 2026-04-21 05:46 CEST  
**Branch:** `master` (up to date with `origin/master`)  
**Commits this session:** 22 (all pushed)  
**Coverage:** 94.7% main, 100% cmdguard/enum/escape/table, 94.6% sort  
**Lint:** 0 issues (`golangci-lint run --timeout 5m ./...`)  
**Tests:** All pass with race detector (`go test -race ./...`)

---

## A. FULLY DONE (22 tasks completed)

| #   | Task                                                                 | Commit        |
| --- | -------------------------------------------------------------------- | ------------- |
| 1   | Fix gofumpt formatting in d2_test.go and ids_test.go                 | `8b6f4c3`     |
| 2   | Fix wsl_v5 warning in ids_test.go                                    | `8b6f4c3`     |
| 3   | Extract test-id string constant in graph_test.go                     | `8b6f4c3`     |
| 4   | Fix gochecknoglobals: move D2 enum slices inside AllowedValues()     | `8b6f4c3`     |
| 5   | Remove dead pkg/errors/ package                                      | `ede4889`     |
| 6   | Replace deprecated ParseOutputFormat in examples/basic/main.go       | `6e605c8`     |
| 7   | Replace deprecated ParseOutputFormat in userjourney_test.go          | `6e605c8`     |
| 8   | Replace deprecated ParseOutputFormat in integration/workflow_test.go | `6e605c8`     |
| 9   | Fix README.md to show Format.Parse API                               | prior session |
| 10  | Split integration_test.go (399→≤350)                                 | `f80ea72`     |
| 11  | Split d2_test.go (362→≤350)                                          | `7e2c71d`     |
| 12  | Split sort_test.go (391→≤350)                                        | `03a6049`     |
| 13  | Split format_test.go (354→≤350)                                      | `6e69a01`     |
| 14  | Migrate html.go to escape.HTML                                       | `84a6c3d`     |
| 15  | Replace inline escape logic in mermaid.go and dot.go                 | `5b8ade7`     |
| 16  | Extract deprecated format aliases → format_deprecated.go             | `2f5bea4`     |
| 17  | Extract shared AddTreeNodes helper                                   | `e9fd2dd`     |
| 18  | Deduplicate graph test helpers                                       | `7630fb0`     |
| 19  | Fix registry test isolation                                          | `45dc766`     |
| 20  | Extract shared SetNodesFromTableData on GraphRendererMixin           | `0321d68`     |
| 21  | Reduce d2_write.go style duplication                                 | `69552b2`     |
| 22  | Add D2 showcase to examples                                          | `1d5282c`     |
| 23  | Document GraphRendererMixin not shared with D2                       | `b416e42`     |

---

## B. PARTIALLY DONE

| Item             | Status                                                                                              | What's Left                                                                                |
| ---------------- | --------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| `.gitignore`     | Fixed malformed merge on line 55, not yet committed                                                 | Commit the fix                                                                             |
| LSP warnings (3) | `golangci-lint` clean, but LSP still flags: `ids_test.go:26,67` wsl_v5; `escape_test.go:16` golines | These are LSP-only, not caught by `golangci-lint run`. May need fixing if LSP rules align. |

---

## C. NOT STARTED — Identified Improvement Opportunities

These are discoveries from the comprehensive audit. None have been started.

---

## D. TOTALLY FUCKED UP — Nothing!

No regressions, no broken tests, no merge conflicts, no failing builds.

---

## E. WHAT WE SHOULD IMPROVE — Deep Audit Findings

### Critical Architecture Issues (🔴)

1. **3 coexisting enum patterns** — `Format`/`SortBy`/`ColorMode` use `enum` package; `GraphShape`/`D2Direction`/`D2NodeShape`/`D2ArrowType` use hand-rolled `slices.Contains`; `FormatCategory` is bare `int` iota with no validation. This is the **single biggest consistency gap** in the codebase.

2. **10+ unused branded types/brands** in `ids.go` — `DOTGraphID`, `DOTNodeID`, `TreeParentID`, `MermaidNodeID`, `MermaidParentID`, `HTMLTitle`, `D2EdgeID`, and several brand types without aliases. Dead code increases cognitive load.

### Medium Issues (🟡)

3. **`StreamingHTMLRenderer` duplicates rendering logic** from `HTMLRenderer` — the `Stream()` method re-implements the same HTML structure line-by-line instead of composing with `HTMLRenderer`.

4. **`ASCIITreeRenderer` not goroutine-safe** — `r.builder` is a mutable field shared across `Render()` calls. Two concurrent renders would race.

5. **`RegisteredFormats` returns unsorted** — map iteration order is non-deterministic, so output varies between calls.

6. **`FormatCategory` is bare `int` iota** — no `String()`, `IsValid()`, `Parse()`. Inconsistent with all other enums.

7. **`CreateRowEdges` returns anonymous struct** — `struct{From, To string}` should be a named type.

8. **`table/` subpackage disconnected from `TableData`** — `table.New()` is standalone lipgloss wrapper; `TableData` feeds all other renderers. No bridge between them.

9. **`color.go` not integrated with lipgloss** — `ShouldColor()` and lipgloss's own terminal detection may conflict. `ToANSI()` returns incomplete ANSI prefix, `isStderrTerminal()` is dead code.

10. **`MarkdownTable` stores own headers/rows** — doesn't use `TableData`, duplicating that data model.

### Low Issues (🟢)

11. **`htmlEscape`/`xmlEscape` trivial wrappers** in `markup.go` — add nothing over calling `escape.HTML`/`escape.XML` directly.

12. **3 deprecated cmdguard files** — `format.go`, `sort.go`, `color.go` still exist as deprecated wrappers.

13. **Inconsistent error wrapping** in `marshal.go` — `MarshalJSONIndent`, `MarshalXML`, and `marshal()` use 3 different error message formats.

14. **`fmt.Fprintf` with `%s`** — `ids.go:49` uses `fmt.Fprintf(s, "%s", id.value)` where `s.WriteString(id.value)` would be faster.

15. **Simple string concatenation via `fmt.Sprintf`** — several calls like `fmt.Sprintf("|%s|", x)` could be `"|" + x + "|"`.

16. **`writeMarkupColumns` hardcodes `<column>` tag** — not parameterized like `writeMarkupRow`.

17. **Escape functions incomplete** — `D2` missing `\\`, `$`, `{`/`}`; `DOT` missing `\\`, `|`, `{`/`}`; `MermaidText` missing `&`, `#`, `<>`.

18. **`HTMLTreeRenderer` CSS promises "collapsible"** but no collapse behavior.

19. **No `CONTRIBUTING.md`** — missing standard OSS doc.

20. **README missing** streaming API, D2 API, Mermaid/DOT API, tree API, branded IDs, sort subpackage.

---

## F. Top #25 Things We Should Get Done Next

Sorted by **impact × (1 / work_required)** — high impact, low effort first.

| #   | Task                                                                                                                 | Impact | Work    | Score |
| --- | -------------------------------------------------------------------------------------------------------------------- | ------ | ------- | ----- | ----- |
| 1   | **Migrate GraphShape to `enum` package** — eliminates the most visible enum inconsistency                            | High   | Low     | ★★★★★ |
| 2   | **Migrate D2 enums to `enum` package** (D2Direction, D2NodeShape, D2ArrowType) — removes 3 more hand-rolled patterns | High   | Low     | ★★★★★ |
| 3   | **Remove unused branded types from ids.go** (DOTGraphID, DOTNodeID, MermaidNodeID, etc.) — reduces dead code         | Med    | Low     | ★★★★☆ |
| 4   | **Fix `RegisteredFormats` to return sorted** — deterministic output                                                  | Med    | Trivial | ★★★★☆ |
| 5   | **Add `String()` to `FormatCategory`** — minimal fix for bare iota                                                   | Med    | Trivial | ★★★★☆ |
| 6   | **Name the `CreateRowEdges` return type** — replace anonymous struct with `RowEdge`                                  | Med    | Low     | ★★★★☆ |
| 7   | **Replace `fmt.Fprintf(b, "%s", x)` with `b.WriteString(x)`** in ids.go                                              | Low    | Trivial | ★★★☆☆ |
| 8   | **Replace `fmt.Sprintf` with string concat** where no formatting needed (mermaid, dot, d2_write)                     | Low    | Low     | ★★★☆☆ |
| 9   | **Remove trivial `htmlEscape`/`xmlEscape` wrappers** — call `escape.HTML`/`escape.XML` directly                      | Low    | Trivial | ★★★☆☆ |
| 10  | **Fix `.gitignore` formatting** and commit                                                                           | Low    | Trivial | ★★★☆☆ |
| 11  | **Consolidate error wrapping in marshal.go** — single consistent pattern                                             | Med    | Low     | ★★★☆☆ |
| 12  | **Delete deprecated cmdguard files** (format.go, sort.go, color.go) — they just wrap `NewEnumFlag`                   | Med    | Trivial | ★★★☆☆ |
| 13  | **Fix `ASCIITreeRenderer` goroutine safety** — localize `builder` in `Render()`                                      | Med    | Low     | ★★★☆☆ |
| 14  | **Add `D2Constraint` Parse/IsValid** — only enum missing validation                                                  | Low    | Low     | ★★★☆☆ |
| 15  | **Bridge `table/` subpackage to `TableData`** — add `RenderFromTableData(data *TableData) string`                    | High   | Med     | ★★★☆☆ |
| 16  | **Make `MarkdownTable` accept `TableData`** — reduce data model duplication                                          | Med    | Med     | ★★☆☆☆ |
| 17  | **Refactor `StreamingHTMLRenderer` to compose with HTMLRenderer** — eliminate duplicate rendering                    | Med    | Med     | ★★☆☆☆ |
| 18  | **Integrate `color.go` with lipgloss** — remove duplicate terminal detection                                         | Med    | Med     | ★★☆☆☆ |
| 19  | **Remove `isStderrTerminal()` dead code** in color.go                                                                | Low    | Trivial | ★★☆☆☆ |
| 20  | **Add missing escape characters** (D2 `\\`/`$`, DOT `\\`/`                                                           | `)     | Med     | Low   | ★★☆☆☆ |
| 21  | **Update README** — add streaming, D2, Mermaid/DOT, tree, sort, branded IDs sections                                 | High   | Med     | ★★☆☆☆ |
| 22  | **Parameterize `writeMarkupColumns`** tag name (like `writeMarkupRow`)                                               | Low    | Low     | ★★☆☆☆ |
| 23  | **Add `CONTRIBUTING.md`**                                                                                            | Med    | Med     | ★☆☆☆☆ |
| 24  | **Fix `HTMLTreeRenderer` CSS** — remove "collapsible" claim or add JS                                                | Low    | Low     | ★☆☆☆☆ |
| 25  | **Sorter error returns** — add error return instead of silent mis-sort                                               | Med    | High    | ★☆☆☆☆ |

---

## G. Top #1 Question I Cannot Figure Out Myself

**Should `GraphShape` use the `enum` package exactly like `Format`/`SortBy`/`ColorMode` (i.e., string-based with `enum.Parse[T]`), or should we generalize the `enum` package to also support int-based enums (like `FormatCategory`)?**

- Currently `enum.Parse[T comparable]` works for `string`-based enums only because `ParseError` wraps a value with `%v` and the allowed values list.
- `FormatCategory` is `int`-based — it can't use `enum.Parse` without either: (a) converting to `string`, or (b) adding a second generic `ParseInt` to the enum package.
- `GraphShape` is `string`-based but uses hand-rolled code that's identical to what `enum.Parse` provides.
- **My recommendation:** Migrate `GraphShape` and all D2 enums to `enum.Parse` first (clear win, no design question). Then decide on `FormatCategory` separately — either make it `string`-based or extend the `enum` package.

---

## Architecture Health Summary

| Dimension             | Rating            | Notes                            |
| --------------------- | ----------------- | -------------------------------- |
| Test coverage         | ✅ 94.7%          | Strong; target 95%+              |
| Lint cleanliness      | ✅ 0 issues       | golangci-lint clean              |
| Enum consistency      | 🟡 Mixed          | 3 patterns coexist               |
| Dead code             | 🟡 ~10% of ids.go | Unused brands                    |
| API documentation     | 🔴 Sparse         | README missing many APIs         |
| Type safety           | ✅ Strong         | Branded IDs, generic enum        |
| Renderer architecture | ✅ Clean          | Interface-based, composable      |
| Streaming             | 🟡 Limited        | Only HTML streams; others buffer |
| Thread safety         | 🟡 1 issue        | ASCIITreeRenderer.builder race   |
| Escape completeness   | 🟡 Gaps           | D2/DOT/Mermaid incomplete        |

---

_Report generated by comprehensive codebase audit on 2026-04-21._
