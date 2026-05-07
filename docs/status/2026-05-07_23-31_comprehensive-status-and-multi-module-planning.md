# go-output — Full Comprehensive Status Report

**Date:** 2026-05-07 23:31 CEST
**Branch:** master
**Status:** Clean working tree, all tests passing, 90.3% coverage
**Recent work:** Multi-module migration planning (3 proposal iterations)

---

## A) FULLY DONE ✅

### Multi-Module Migration Planning

| Item                          | Details                                                    | Commit                 |
| ----------------------------- | ---------------------------------------------------------- | ---------------------- |
| v1 Proposal                   | Initial 9-module split with dependency graph               | `5670077`              |
| v2 Proposal                   | Self-critique, coupling analysis, 7-phase plan             | `eda0152`              |
| v3 Proposal (current)         | No core/ dir, root stays as-is, deprecate sort/, 8 modules | `30c34b1`              |
| Dependency graph analysis     | Full cross-file import map for all 28 root .go files       | Included in proposal   |
| go-cqrs-lite pattern research | Module paths, replace directives, testhelpers pattern      | Referenced in proposal |
| Go workspace research         | go.work, use/replace directive patterns                    | Referenced in proposal |

### "Reinventing the Wheel" Audit

| Item                             | Verdict              | Details                                                                                                         |
| -------------------------------- | -------------------- | --------------------------------------------------------------------------------------------------------------- |
| `sort/` package                  | **DEPRECATE**        | `slices.SortStableFunc` + `cmp.Compare` (Go 1.21+) does the same. `SortBy` stored but never used in sort logic. |
| `escape.HTML()` / `escape.XML()` | **USE STDLIB**       | `html.EscapeString()` / `xml.EscapeText()` exist. Only `&apos;` vs `&#39;` difference.                          |
| `color.go` terminal detection    | **PARTIAL**          | `termenv.EnvColorProfile()` does the same. Our CI env var list is incomplete.                                   |
| `yaml.go`                        | **THINNEST WRAPPER** | Literally just `marshal("yaml", yaml.Marshal, v)`. Zero logic.                                                  |

### Previous Plan Items (from 2026-04-29)

| Item                                                                  | Status              |
| --------------------------------------------------------------------- | ------------------- |
| Commit pending fixes                                                  | ✅ Done             |
| Remove ToANSI, GetStyle                                               | ✅ Done             |
| Delete stale root docs (IMPROVEMENTS_SUMMARY, BDD_TESTS_REVIEW, etc.) | ✅ Done             |
| Split d2.go, extract D2 enums into d2_enum.go                         | ✅ Done             |
| Extract shared example utilities                                      | ✅ Done (`3442bce`) |
| Extract renderFullHTMLWithStyles helper                               | ✅ Done (`751da61`) |

---

## B) PARTIALLY DONE 🟡

| Item                                 | What's done                                                                     | What's left                                   |
| ------------------------------------ | ------------------------------------------------------------------------------- | --------------------------------------------- |
| Multi-module proposal                | v3 written, all analysis complete                                               | Zero phases executed — purely a planning doc  |
| `sort/` deprecation                  | Identified as reinventing the wheel                                             | No deprecation notice added, no code changed  |
| `escape.HTML/XML` stdlib replacement | Identified as redundant                                                         | No code changed                               |
| SortBy enum audit                    | Found: only used by `sort/`, `cmdguard_test.go` (as example enum), 2 test files | No decision made on whether to keep/deprecate |
| Improvement plan execution           | 6 of 30 tasks done                                                              | 24 tasks remain (see NOT STARTED)             |

---

## C) NOT STARTED ❌

### Multi-Module Migration (0 of 7 phases executed)

| Phase   | Description                                        | Status      |
| ------- | -------------------------------------------------- | ----------- |
| Phase 1 | Leaf modules: go.mod for enum/, escape/, cmdguard/ | Not started |
| Phase 2 | Deprecate sort/ package                            | Not started |
| Phase 3 | Extract table/ as module (lipgloss isolation)      | Not started |
| Phase 4 | Extract d2/ as module (5 files moved)              | Not started |
| Phase 5 | Extract graph/ as module (DOT+Mermaid)             | Not started |
| Phase 6 | Integration + Examples go.mod files                | Not started |
| Phase 7 | Polish: ADR, justfile, AGENTS.md, README           | Not started |

### Code Quality (from 2026-04-29 plan, still open)

| Item                                                                 | Details                                                                                                                      |
| -------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Registry is a ghost system                                           | `Register()`/`Create()`/`IsRegistered()` exist but only used in tests. Real API surface with zero production consumers.      |
| Tree conversion not unified                                          | `d2_convert.go`, `dot.go`, `mermaid.go` each have their own `addTreeNodes()` — should use shared `graph.AddTreeNodes()`      |
| `format_deprecated.go` still exists                                  | Backward-compat shim — planned for removal                                                                                   |
| Depguard violations in examples/                                     | `examples/basic` and `examples/d2` import `examples/shared` which is not in allowlist                                        |
| Test helper duplication                                              | 7 items in `output_test_helpers.go` duplicate `internal/testutils/` (structural necessity for same-package tests, but messy) |
| `fuzz_test.go` stringEnum constraint duplicates `gentest.StringEnum` | Could unify into gentest                                                                                                     |
| PLAN.md outdated                                                     | References old structure, wrong file counts                                                                                  |
| No CI workflows                                                      | `.github/workflows/` does not exist                                                                                          |

### "Reinventing the Wheel" Fixes

| Item                                                        | Status      |
| ----------------------------------------------------------- | ----------- |
| Replace `escape.HTML()` with `html.EscapeString()`          | Not started |
| Replace `escape.XML()` with `xml.EscapeText()`              | Not started |
| Evaluate `color.go` → `termenv` for detection               | Not started |
| Evaluate `yaml.go` — inline or add value                    | Not started |
| Evaluate `slices.go` — inline `FilledStrings` at call sites | Not started |

---

## D) TOTALLY FUCKED UP 💥

| Item                                           | Severity        | Details                                                                                                             |
| ---------------------------------------------- | --------------- | ------------------------------------------------------------------------------------------------------------------- |
| **No CI**                                      | Critical        | Zero `.github/workflows/`. No automated build/test/lint on push. Entire quality gate is manual.                     |
| **v1 proposal had empty dirs**                 | Embarrassing    | Proposed enum/ and escape/ as modules "with no .go files" — they already had .go files. User caught this.           |
| **v2 proposed core/ for no reason**            | Wasted effort   | Moving 28 files and renaming package output → core would be massive churn with zero user benefit. Root IS the core. |
| **Depguard violations ignored**                | Ongoing         | 2 active warnings in examples/ that nobody fixed                                                                    |
| **30-item improvement plan at 20% completion** | Process failure | Created 2026-04-29, 8 days ago. Only 6 of 30 items done. No follow-through.                                         |

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Design

1. **Type model improvement:** `TableData` is the central data structure — consider making it an interface so different formatters can optimize for their use case (streaming, chunked, etc.)
2. **Renderer interface is minimal:** `Render() (string, error)` — no `io.Writer` support except for streaming HTML. Should consider `RenderTo(io.Writer) error` as standard.
3. **GraphNode/Edge types are too simple:** Only `ID`, `Label`, `Shape`, `Style`, `Metadata`. D2 has much richer types. The conversion layer (`d2_convert.go`) is lossy.
4. **No error types standardization:** Each formatter wraps errors differently. Should have `RenderError`, `MarshalError`, `WriteError` types.
5. **Registry has no consumers:** Dead API surface. Either integrate it properly or remove it.

### Process & Quality

6. **No CI:** This is the biggest process gap. Every commit should run build+test+lint automatically.
7. **No ADRs:** Architecture decisions (multi-module, format categories, branded IDs) are not formally documented.
8. **Stale docs:** PLAN.md, some status reports, and improvement plans are outdated.
9. **No versioning strategy:** Pre-v1.0 with no release process defined.

### Code Quality

10. **Test helper duplication** between `output_test_helpers.go` and `internal/testutils/`
11. **GraphRendererMixin in wrong file** (dot.go instead of graph.go or its own file)
12. **FilledStrings** is a trivial 8-line helper in its own file — inline it
13. **stringEnum** in fuzz_test.go duplicates `gentest.StringEnum`

---

## F) Top 25 Things to Get Done Next

Sorted by: **impact × urgency ÷ effort** (highest first)

### Tier 1: High Impact, Low Effort (do immediately)

| #   | Task                                                             | Effort | Impact | Rationale                                          |
| --- | ---------------------------------------------------------------- | ------ | ------ | -------------------------------------------------- |
| 1   | Replace `escape.HTML()` with `html.EscapeString()`               | 10 min | Medium | Removes reinvented code, uses battle-tested stdlib |
| 2   | Replace `escape.XML()` with `xml.EscapeText()` or similar        | 10 min | Medium | Same as above                                      |
| 3   | Add deprecation notice to `sort/sorter.go` and `sort/compare.go` | 5 min  | Medium | Signals intent, costs nothing                      |
| 4   | Add `go.mod` to `enum/`, `escape/`, `cmdguard/` (leaf modules)   | 15 min | High   | First real multi-module step, zero risk            |
| 5   | Create `go.work` at root                                         | 5 min  | High   | Enables workspace development                      |
| 6   | Fix depguard violations in examples/                             | 10 min | Low    | Clean lint output                                  |
| 7   | Inline `FilledStrings` at call sites, remove `slices.go`         | 10 min | Low    | Removes trivial file                               |
| 8   | Unify `stringEnum` in fuzz_test.go → use `gentest.StringEnum`    | 5 min  | Low    | DRY                                                |

### Tier 2: High Impact, Medium Effort (do soon)

| #   | Task                                                         | Effort | Impact | Rationale                                         |
| --- | ------------------------------------------------------------ | ------ | ------ | ------------------------------------------------- |
| 9   | Extract `table/` as module (lipgloss isolation)              | 30 min | High   | Biggest dependency win — lipgloss is heaviest dep |
| 10  | Deprecate `SortBy` enum if no real consumers                 | 15 min | Medium | Removes dead type after sort/ deprecation         |
| 11  | Extract `d2/` as module (5 files + 6 test files)             | 45 min | High   | Clean domain boundary                             |
| 12  | Extract `graph/` as module (DOT+Mermaid+Mixin)               | 45 min | High   | Diagram code isolated from core                   |
| 13  | Evaluate Registry: integrate or remove                       | 20 min | Medium | Ghost system — decide its fate                    |
| 14  | Unify tree conversion: use `graph.AddTreeNodes()` everywhere | 20 min | Medium | DRY across d2_convert, dot, mermaid               |
| 15  | Write ADR 001: multi-module split decision                   | 15 min | Medium | Formal decision record                            |

### Tier 3: Medium Impact, Medium Effort (plan for next session)

| #   | Task                                                                | Effort | Impact | Rationale                             |
| --- | ------------------------------------------------------------------- | ------ | ------ | ------------------------------------- |
| 16  | Create `.github/workflows/ci.yml` — build+test+lint                 | 30 min | High   | No CI is critical gap                 |
| 17  | Evaluate `color.go` → use `termenv` for detection                   | 20 min | Low    | Better CI detection, less reinvention |
| 18  | Evaluate `yaml.go` — add value or inline                            | 15 min | Low    | Thinnest wrapper, adds almost nothing |
| 19  | Update `integration/` to standalone module with go.mod              | 20 min | Medium | Cross-module integration tests        |
| 20  | Update `examples/` to standalone module with go.mod                 | 15 min | Medium | Examples work independently           |
| 21  | Deduplicate test helpers (output_test_helpers.go vs testutils)      | 30 min | Low    | Cleaner test infrastructure           |
| 22  | Move `GraphRendererMixin` from `dot.go` to `graph.go` (or own file) | 10 min | Low    | Right place for shared code           |

### Tier 4: Polish & Documentation (when time permits)

| #   | Task                                                   | Effort | Impact | Rationale                |
| --- | ------------------------------------------------------ | ------ | ------ | ------------------------ |
| 23  | Update AGENTS.md with final multi-module structure     | 15 min | Medium | AI context stays current |
| 24  | Update README.md with module paths and examples        | 15 min | Medium | User-facing docs         |
| 25  | Remove stale docs/status/ reports (5 files from April) | 10 min | Low    | Clean docs tree          |

---

## G) Top #1 Question I Cannot Figure Out Myself

**Is the multi-module migration actually worth doing right now?**

Arguments for:

- Lipgloss isolation (table/ module) is a real win for users who only need JSON/YAML
- D2/graph isolation keeps the core small for simple use cases
- Follows go-cqrs-lite pattern, which works well

Arguments against:

- This is a pre-v1.0 library with zero known external users
- The migration involves moving files, changing package names, updating imports in ~40 test/example files
- No CI exists to catch regressions during migration
- The improvement plan from 8 days ago is at 20% — we keep planning instead of shipping

**My honest take:** The 3 leaf modules (enum/escape/cmdguard) are zero-risk wins. The table/ extraction is the only one with real dependency benefit. D2/ and graph/ are nice-to-have. We should do the easy stuff first and ship it before planning more.

---

## Metrics Summary

| Metric                 | Value                                |
| ---------------------- | ------------------------------------ |
| Root package coverage  | 90.3%                                |
| Sub-package coverage   | 95-100%                              |
| Total root .go files   | 45 (8,241 lines)                     |
| Third-party deps       | 3 (lipgloss, go-faster/yaml, x/term) |
| Supported formats      | 12                                   |
| Test files             | ~30                                  |
| CI                     | ❌ None                              |
| Open ADRs              | 0                                    |
| Proposal iterations    | 3 (v1 → v2 → v3)                     |
| Actual modules created | 0                                    |
