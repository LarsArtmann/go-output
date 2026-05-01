# go-output — Comprehensive Status Report

**Date:** 2026-04-30 00:49  
**Session:** Render() (string, error) migration — self-review and cleanup  
**Commits since last report:** 3 (7bb7c5a..b512f07)  
**Net change:** +753 / -168 lines across 37 files  
**Status:** All tests pass, lint clean, race clean, pushed to origin/master

---

## A) FULLY DONE

### Critical: Renderer.Render() → Render() (string, error)

|| # | Change | Impact | Commit |
||---|--------|--------|--------|
|| 1 | Change `Renderer` interface to return `(string, error)` | Breaking API change | `531454e` |
|| 2 | Update all 10 production renderer implementations | Interface compliance | `531454e` |
|| 3 | Update `RenderFullHTML()` to return `(string, error)` | Secondary API surface | `531454e` |
|| 4 | Update `StreamingHTMLRenderer.Render()` to propagate Stream() errors | Error propagation | `531454e` |
|| 5 | Update `adapterRenderer.Render()` to wrap delegate errors | wrapcheck compliance | `531454e` |
|| 6 | Update ~200 test call sites across 20 test files | Test compliance | `531454e` |
|| 7 | Add `MustRender(r Renderer) string` helper | Boilerplate reduction | `4b3122f` |
|| 8 | Fix `fmt.Println(x.Render())` in examples printing `<nil>` | Runtime bug | `4b3122f` |
|| 9 | Fix README.md discarding error with `out, _ :=` | Bad example | `4b3122f` |
|| 10 | Fix PLAN.md code examples | Broken docs | `4b3122f` |
|| 11 | Update FORMAT_ARCHITECTURE.md interface signature | Doc accuracy | `531454e` |
|| 12 | Update AGENTS.md with MustRender reference | Doc accuracy | `b512f07` |

### Previously Done (from 2026-04-29 session)

|| # | Change | Commit |
||---|--------|--------|
|| 13 | Sort descending stability violation fix | `e882804` |
|| 14 | CI Go 1.23→1.26 + job consolidation | `bf543fb` |
|| 15 | `golang.org/x/term` direct dependency fix | `fe79fec` |
|| 16 | Recursive `joinStrings` → iterative fix | `fe79fec` |
|| 17 | D2 enums: allocating functions → package-level vars | `4aa86a8` |
|| 18 | Split d2.go → d2.go + d2_enum.go | `6f274fa` |
|| 19 | Remove dead `ColorMode.ToANSI()` | `3b1aeea` |
|| 20 | Remove dead `GraphNode.GetStyle()` | `3b1aeea` |
|| 21 | Document Registry as opt-in plugin system | `b7d5cdf` |
|| 22 | Delete 15 stale docs (-4,693 lines) | `1e36a4d` |
|| 23 | Rewrite PLAN.md | `2c91e7e` |
|| 24 | Full AGENTS.md rewrite | `9a943a8` |
|| 25 | Fix README.md deps/imports | `9a943a8` |

---

## B) PARTIALLY DONE

### MustRender() Adoption

**Status:** Helper exists but is unused in tests. 106 test call sites still use the verbose pattern:

```go
got, err := x.Render()
if err != nil {
    t.Fatalf("Render() error = %v", err)
}
```

Could be replaced with:

```go
got := output.MustRender(x)
```

This is safe for tests because:

- No current renderer actually returns an error
- Tests would catch regressions via `MustRender` panicking on unexpected errors
- Benchmarks would still need the verbose pattern (to avoid panic overhead in hot loops)

### Error Path Test Coverage

**Status:** `Render() (string, error)` is the interface, but no renderer currently exercises error paths in tests. The coverage dropped from 91.0% to 90.1% due to untested error handling branches in:

- `streaming.go` — `Render()` error propagation from `Stream()`
- `html.go` — `RenderFullHTML()` wrapping `Render()` error
- `format.go` — `MustRender()` panic-on-error path

---

## C) NOT STARTED

|| # | Task | Impact | Effort | Notes |
||---|------|--------|--------|-------|
|| 1 | Refactor tests to use `MustRender()` | Medium | Medium | 106 call sites across 20 test files |
|| 2 | Add test for `MustRender()` panicking | Medium | Low | Currently 0% coverage on panic path |
|| 3 | Add error path tests for streaming/html | Medium | Low | Cover new error branches |
|| 4 | Add godoc Example functions for top 5 types | High | Medium | pkg.go.dev shows no examples |
|| 5 | Add round-trip marshal/unmarshal integration tests | High | Low | MarshalJSON → UnmarshalJSON equality |
|| 6 | Update CHANGELOG.md with recent fixes | Medium | Low | Hasn't been updated |
|| 7 | Add `cmp` to depguard allowed list | Medium | Low | Modern Go patterns |
|| 8 | Add fuzz tests for escape.D2, escape.DOT, escape.Mermaid | High | Medium | Only ParseFormat fuzzed |
|| 9 | Add benchmarks for XML, Tree, HTML tree | Medium | Low | Coverage gaps |
|| 10 | Add cmdguard usage example to examples/ | Medium | Low | examples/basic uses raw os.Args |
|| 11 | Add functional options for D2Node (reduce exhaustruct nolints) | Medium | High | 9 exhaustruct nolints remain |
|| 12 | Deterministic output ordering for D2 diagrams | Medium | Medium | Go maps are unordered |
|| 13 | Unify tree conversion pattern | Medium | High | 3 renderer-specific addTreeNodes |
|| 14 | DOT attribute quoting for special characters | Low | Low | Robustness |
|| 15 | Mermaid ID sanitization edge cases | Low | Low | Unicode, reserved words |
|| 16 | Markdown table captions | Low | Low | Nice-to-have |
|| 17 | `TableData.RemoveRow()` | Low | Low | No consumer requesting it |
|| 18 | `Format.AutoDetect()` from content sniffing | Low | Medium | Questionable value |
|| 19 | CSV/TSV reading (currently write-only) | Medium | Medium | Would expand library scope |
|| 20 | Add `format_deprecated.go` removal timeline | Low | Low | Breaking change |
|| 21 | Add `MarshalTSVFromTableData` | Medium | Low | Feature parity with XML |
|| 22 | Add `MarshalCSVFromTableData` | Medium | Low | Feature parity with XML |
|| 23 | Create .goreleaser.yml if needed | Low | Low | Release workflow adjustment |
|| 24 | Add `go vet` as explicit CI step | Low | Low | Currently implicit |
|| 25 | Explore Nix flake migration | Low | High | Tooling |

---

## D) TOTALLY FUCKED UP

Nothing is fucked up:

- **0 lint issues** (golangci-lint clean)
- **0 vet issues**
- **0 race conditions** (all tests pass with -race)
- **0 build errors**
- **CI fixed** (Go 1.26 matching go.mod)
- **All docs accurate** (README, PLAN, AGENTS.md, FORMAT_ARCHITECTURE)
- **Examples produce correct output** (verified D2/tree/table/markdown)
- **Clean working tree** (nothing uncommitted)

---

## E) WHAT WE SHOULD IMPROVE

### Architecture

1. **MustRender is unused** — 106 test sites still use verbose error handling. Adding MustRender but not using it is half-done work.
2. **No renderers actually return errors** — The interface supports errors but no renderer exercises them. Consider what real failure modes exist (invalid UTF-8? disk full on streaming?) and add error returns for those cases.
3. **Exhaustruct as design smell** — 9 nolint directives for `exhaustruct` in production code. D2Node/D2Edge have too many optional fields. Functional option pattern would be better.
4. **GraphRendererMixin vs D2 parallel types** — D2 has richer types but can't use the mixin. Intentional but creates parallel hierarchy.

### Testing

5. **Coverage regression** — 91.0% → 90.1% due to untested error paths in streaming.go and html.go. MustRender() at 0% coverage.
6. **No integration test for round-trip marshal/unmarshal** — MarshalJSON → UnmarshalJSON should produce equal data.
7. **Escape functions lack adversarial input testing** — fuzz tests would catch edge cases.

### Developer Experience

8. **No godoc examples** — pkg.go.dev shows no example usage for any type.
9. **cmdguard integration not shown in examples** — examples/basic uses raw os.Args.
10. **CHANGELOG.md stale** — Hasn't been updated with recent fixes.

### Type Model Improvements

11. **D2Node/D2Edge functional options** — Replace 9 exhaustruct nolints with `NewD2Node(id, label string, opts ...D2NodeOption)` pattern. Same for D2Edge.
12. **RenderError type** — Consider a dedicated `RenderError` type wrapping renderer name + context, instead of generic `fmt.Errorf`.
13. **Option pattern for renderers** — HTML/DOT/Mermaid renderers have configuration that could use functional options instead of direct field access.

### Operational

14. **Depguard blocks `cmp`** — should be added to allowed list.
15. **No .goreleaser.yml** — release.yml referenced it but it was removed.

---

## F) TOP 25 THINGS TO DO NEXT

Sorted by impact × effort (highest first):

|| # | Task | Impact | Effort | Category |
||---|------|--------|--------|----------|
|| 1 | Refactor 106 test Render() calls to use MustRender() | High | Medium | Testing |
|| 2 | Add MustRender() panic test | High | Low | Testing |
|| 3 | Add error path tests for streaming/html Render() | High | Low | Testing |
|| 4 | Add godoc Example functions for top 5 types | High | Medium | DX |
|| 5 | Add round-trip marshal/unmarshal integration tests | High | Low | Testing |
|| 6 | Update CHANGELOG.md with recent fixes | Medium | Low | Docs |
|| 7 | Add `cmp` to depguard allowed list | Medium | Low | Tooling |
|| 8 | Add fuzz tests for escape.D2, escape.DOT, escape.Mermaid | High | Medium | Testing |
|| 9 | Add benchmarks for XML, Tree, HTML tree renderers | Medium | Low | Testing |
|| 10 | Add cmdguard usage example to examples/ | Medium | Low | DX |
|| 11 | Add functional options for D2Node | Medium | High | Architecture |
|| 12 | Add deterministic output ordering for D2 diagrams | Medium | Medium | Correctness |
|| 13 | Add `MarshalTSVFromTableData` | Medium | Low | Feature parity |
|| 14 | Add `MarshalCSVFromTableData` | Medium | Low | Feature parity |
|| 15 | Unify tree conversion pattern | Medium | High | Architecture |
|| 16 | Add DOT attribute quoting for special characters | Low | Low | Robustness |
|| 17 | Add Mermaid Unicode ID handling | Low | Low | Robustness |
|| 18 | Add `TableData.RemoveRow()` | Low | Low | Feature |
|| 19 | Add Markdown table caption support | Low | Low | Feature |
|| 20 | Add `Format.AutoDetect()` from content sniffing | Low | Medium | Feature |
|| 21 | Create .goreleaser.yml if releasing binaries ever needed | Low | Low | CI |
|| 22 | Add pre-commit hook to CI | Low | Low | CI |
|| 23 | Add `go vet` as explicit CI step | Low | Low | CI |
|| 24 | Add D2 theme/style preset support | Low | Medium | Feature |
|| 25 | Explore Nix flake migration | Low | High | Tooling |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**Should renderers actually return errors for real failure modes now, or is the current "always return nil" adequate?**

The `(string, error)` interface is future-proof, but no renderer currently exercises it. The question is:

1. **StreamingHTMLRenderer** — `Stream()` can fail on I/O. This is a real error path. ✅ Already returns errors.
2. **All other renderers** — They build strings in memory. What could fail?
   - Invalid UTF-8 in input? (Currently just passed through)
   - Encoding issues? (Currently escaped silently)
   - Memory allocation failure? (Go doesn't expose this)

Should we:

- **A)** Leave renderers returning `nil` errors until real failure modes emerge? (current state)
- **B)** Add validation (e.g., UTF-8 checking) that could return errors? (more defensive)
- **C)** Accept that the `error` return is primarily for streaming/I/O and the interface is correctly forward-looking?

I lean toward **C**, but I'm not sure if there are downstream consumers who expect errors from non-streaming renderers.

---

## Codebase Health Metrics

|| Metric | Value |
||--------|-------|
| Production code | 4,642 lines |
| Test code | 6,861 lines |
| Test:Code ratio | 1.48:1 |
| Packages | 10 (7 with tests) |
| Exported types | ~30 |
| Exported functions | ~82 |
| CI jobs | 1 |
| Lint issues | 0 |
| Race issues | 0 |
| Build errors | 0 |
| Coverage (root) | 90.1% |
| Coverage (all avg) | 95.5% |
| Largest prod file | format.go (323 lines) |
| nolint directives (prod) | 25 (9 exhaustruct, 13 gochecknoglobals, 3 gosec) |
| MustRender test coverage | 0% |
| Test Render() calls | 111 |
| MustRender could replace | 106 |

---

## Commits This Session

```
b512f07 docs: add status report and update AGENTS.md for Render() migration
4b3122f fix(examples): handle Render() error returns and add MustRender helper
531454e feat!: change Renderer.Render() to return (string, error)
```
