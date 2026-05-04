# go-output: Public or Private?

**Decision Date:** 2026-05-04
**Status:** ⚖️ OPEN — Conditional recommendation: **Go public, with prerequisites**

---

## Executive Summary

**go-output** is a well-engineered, uniquely positioned Go library with no direct competitor in the ecosystem. It fills a genuine gap (unified 12-format output + diagram generation). The code is clean, tested (90%+ coverage), and free of sensitive data. The case for going public is strong — but two blockers must be resolved first.

**Recommendation:** Make public **after** addressing the two blockers below. Then pursue ecosystem visibility.

---

## PRO — Why Go Public

### 1. Unique Ecosystem Position

No other Go library unifies data serialization (JSON/YAML/XML/CSV/TSV), presentation formatting (Table/Markdown/HTML/Tree), and diagram generation (D2/Mermaid/DOT) behind a single type-safe API.

- `go-pretty` covers 5 table formats but no diagrams or type-safe enums
- `lipgloss` handles terminal rendering but not data format abstraction
- No Go library exists for programmatic D2 diagram generation
- No Go library exists for Mermaid diagram generation from typed Go structs

**This is a genuine ecosystem gap that go-output fills.**

### 2. Production Quality

| Metric               | Score                                             |
| -------------------- | ------------------------------------------------- |
| Test coverage        | 90.3% root, 95-100% subpackages                   |
| Linter               | golangci-lint clean                               |
| Race detector        | Passing                                           |
| Fuzz tests           | Present (ParseFormat, ParseSortBy, MarkdownTable) |
| Benchmarks           | Present                                           |
| CI/CD                | Full pipeline (build, test, lint, tidy check)     |
| Release workflow     | Automated with tag-based releases                 |
| Code duplication     | Actively managed (threshold: 30 tokens)           |
| File size discipline | 350-line limit enforced                           |

### 3. Type Safety as Differentiator

- Branded IDs via phantom types (compile-time ID mixing prevention)
- Type-safe format enums with Parse/Validate/AllowedValues
- Generic Sorter[T] with ByField[T, F cmp.Ordered]
- Format categories (table/tree/graph) for programmatic filtering

This is above-average Go library craftsmanship. The Go community values this.

### 4. Diagram Languages Are Rising

- D2 is gaining adoption (GitHub native support, Terrastruct investment)
- Mermaid is the de-facto diagram format for GitHub/GitLab markdown
- DOT/Graphviz remains the standard for graph visualization
- Having a Go library that generates all three from the same data model is compelling

### 5. Marketing & Portfolio Value

- Demonstrates Go generics mastery, API design skill, and library architecture
- Provides concrete reference for consulting/collaboration opportunities
- MIT license in CONTRIBUTING.md signals openness (need to update LICENSE)

### 6. Low Maintenance Overhead for Public

- No database, no network, no authentication — pure library
- API is already stabilized (only at v0.1.0, clean migration path)
- Registry system allows community extension without core changes

---

## CONTRA — Why Stay Private (or Wait)

### 1. 🚫 BLOCKER: License Mismatch

The LICENSE file says **PROPRIETARY**. The CONTRIBUTING.md says contributors agree to **MIT**. The README badge says **MIT**. This is inconsistent and legally ambiguous.

**Resolution:** Replace LICENSE with MIT before going public. This is non-negotiable.

### 2. 🚫 BLOCKER: Hardcoded Local Paths in README

README.md:13-14 reference personal local paths:

```
- `/Users/larsartmann/projects/project-meta/`
- `/Users/larsartmann/projects/projects-management-automation/`
```

This reveals personal system structure and looks unprofessional for a public library.

**Resolution:** Remove or generalize the "Purpose" section before going public.

### 3. Audience Scope

The README currently positions it as "Standardizes output formatting across personal Go projects." A public library should be framed as solving a general problem for the Go community.

**Resolution:** Reframe README to emphasize the general use case. The personal projects can remain as motivation but shouldn't be the primary framing.

### 4. Pre-v1.0 API Stability

Currently at v0.1.0 with `OutputFormat` deprecation warnings for v2.0. Public users may hesitate to depend on a pre-v1.0 library. The `Render()` → `(string, error)` migration was recently done — expect more breaking changes as the API stabilizes.

**Mitigation:** Not a blocker. v0.x is normal for new libraries. Just be explicit about stability expectations in README.

### 5. Dependency on charm.land/lipgloss/v2

Lipgloss v2 is relatively new. If it has breaking changes, go-output inherits that risk.

**Mitigation:** Not a blocker. Charmbracelet libraries are widely adopted and well-maintained.

### 6. No Godoc / API Reference Site

The library has good README documentation but no generated API reference (pkg.go.dev will auto-generate this, so minimal effort needed).

**Mitigation:** pkg.go.dev auto-indexes public repos. Ensure Go doc comments are good on exported types.

### 7. Internal Status Reports and Planning Docs

The `docs/status/` and `docs/planning/` directories contain internal status reports with detailed session reviews. Not sensitive, but may reveal working style and internal process that you may not want public.

**Mitigation:** Not a blocker — these are normal. Can add to .gitignore if desired.

---

## Conditional Recommendation

### Go public after completing these steps:

- [ ] **Replace LICENSE file** with MIT license text
- [ ] **Remove hardcoded paths** from README.md (lines 13-14)
- [ ] **Reframe README purpose** from "personal project utility" to "general-purpose Go library"
- [ ] **Verify all exported symbols have Go doc comments** (for pkg.go.dev)
- [ ] **Tag v0.2.0** to mark the public release milestone

### Then pursue visibility:

- [ ] Post to r/golang
- [ ] Submit to Awesome Go
- [ ] Share in Go Discord / Gophers Slack
- [ ] Write a blog post: "One Go library, 12 output formats"

---

## Risk Assessment

| Risk                               | Severity   | Mitigation                                                  |
| ---------------------------------- | ---------- | ----------------------------------------------------------- |
| API breaking changes upset users   | Medium     | Use Go module versioning, clear CHANGELOG, v1.0 when stable |
| Maintenance burden from issues/PRs | Low        | Start with clear CONTRIBUTING.md, respond weekly            |
| Adoption is low                    | Low        | No downside — private library with zero users is worse      |
| License confusion (current state)  | High       | Fix before going public                                     |
| Competitor copies the approach     | Negligible | First-mover advantage + execution quality matters more      |

---

## Final Verdict

**The project is a strong candidate for going public.** The code quality is high, the ecosystem gap is real, and the maintenance overhead is low. The only hard blockers are the LICENSE file and the hardcoded local paths — both are 5-minute fixes.

The Go community benefits from libraries like this. Ship it.
