# Footer Row Polish — Execution Plan

**Date:** 2026-05-28
**Session:** Round 4 — Comprehensive Polish

---

## What We Forgot / Could Do Better

1. **9 packages missing `doc.go`** — pkg.go.dev shows nothing for these packages
2. **Root coverage at 88.6%** — below 90% target
3. **Serialization coverage at 88.3%** — below 90% target
4. **13 tests missing `t.Parallel()`** in render_tabledata_test.go
5. **`markup/xml_test.go` at 341 lines** — 9 lines from 350 limit
6. **40+ exported struct fields** missing GoDoc (graph, tree, d2 types)
7. **Inconsistent error types** — InvalidFormatError includes Allowed, others don't
8. **Missing GoDoc examples** for Validate(), ColorMode, Shape
9. **TODO_LIST.md** not updated with footer work completion

## Decisions Made (Skip List)

- **Skip generic enum macro** — Current explicit types are better for GoDoc/discoverability
- **Skip generic SerializationRenderer** — Specific types are better for API; internal duplication is acceptable
- **Skip MarkdownTable→TableRenderer adapter** — Would require breaking changes
- **Skip RenderOptions functional options** — Breaking API change, not justified
- **Skip moving internal/gentest** — Needs user decision, documented as TODO #20
- **Skip graph FromTableData dedup** — Helper functions are small, not worth refactoring risk

## Execution Plan (18 tasks, ~10min each)

| # | Task | Impact | Effort | Module |
|---|------|--------|--------|--------|
| 1 | Add doc.go to 9 packages | HIGH | LOW | all |
| 2 | Add t.Parallel() to 13 tests | LOW | LOW | root |
| 3 | Split xml_test.go (341→2 files) | LOW | LOW | markup |
| 4 | Root coverage 88.6→90%+ | MED | LOW | root |
| 5 | Serialization coverage 88.3→90%+ | MED | LOW | serialization |
| 6 | Footer edge-case tests | MED | LOW | root |
| 7 | GoDoc on exported struct fields (root: graph, tree) | MED | MED | root |
| 8 | GoDoc on exported struct fields (d2) | MED | MED | d2 |
| 9 | Make InvalidError types consistent | MED | LOW | root |
| 10 | GoDoc examples for Validate, ColorMode, Shape | MED | LOW | root |
| 11 | Update FEATURES.md footer matrix | MED | LOW | docs |
| 12 | Consolidate brandedEdgeLabel | LOW | LOW | serialization |
| 13 | Update CHANGELOG.md | MED | LOW | docs |
| 14 | Update TODO_LIST.md | LOW | LOW | docs |
| 15 | Update AGENTS.md | MED | LOW | docs |
| 16 | Write final status report | LOW | LOW | docs |
| 17 | Integration test coverage bump | MED | LOW | integration |
| 18 | Final build/test/lint + push | HIGH | LOW | all |
