# Status Report: Public Presence (README, Website, GitHub) — 2026-07-13

## Context

Goal: make the three public-facing surfaces (README.md, website, GitHub metadata) superb for the `go-output` repo.

---


> **✅ Resolved (2026-08-04):**
>
> The pointer dereference bug was superseded — GraphBuilder changed to fluent chain API in v0.31.0 (`NewGraphBuilder()` returns `*GraphBuilder`, `AddNode`/`AddEdge` return `*GraphBuilder`). The old `*output.NewGraphNode(...)` pattern no longer applies. Website deployed at `go-output.lars.software`. Remaining open items are in TODO_LIST (v1.0.0 tag #15, community launch #14).

---

## a) FULLY DONE

### 1. GitHub Metadata

- **Description**: Updated from stale "12 formats" to accurate "16 output formats across tables, trees, and diagrams, plus NOM-style real-time progress visualization"
- **Homepage URL**: Set to `https://go-output.lars.software`
- **Topics**: Curated to 20/20 — swapped low-value generic topics (`csv`, `xml`, `enum`) for high-value ones (`toml`, `plantuml`, `progress`, `diagrams`, `visualization`). Verified final state.

### 2. Website — Project Structure

- Full Astro 6.x + Starlight 0.39.x + Tailwind v4 project in `website/`
- Config files: `package.json`, `astro.config.mjs`, `firebase.json`, `.firebaserc`, `flake.nix`, `tsconfig.json`, `.node-version`, `.htmlvalidate.json`, `.gitignore`
- Firebase hosting config with security headers (HSTS, X-Frame-Options, etc.)
- CSP fix script (`scripts/fix-csp.mjs`) — same pattern as gogenfilter
- Nix flake with `dev`, `build`, `preview`, `deploy` apps

### 3. Website — Landing Page

- **HeroSection**: Live CQRS code example, GitHub stars badge, `go get` install command, copy-to-clipboard
- **FormatGrid**: Visual matrix of all 16 formats with shape badges (table/tree/graph)
- **FeatureGrid**: 6 feature cards (16 Formats, Type-Safe, Zero Deps, NOM Progress, Streaming, Color Modes)
- **PhaseSection**: 3-step Build/Freeze/Render CQRS visualization
- **ComparisonSection**: go-output vs DIY vs Heavy comparison
- **UseCasesSection**: 4 use cases (CLI Tools, CI/CD, Documentation, Dashboards)
- **CTASection**: Links to docs and format matrix

### 4. Website — Documentation (14 pages)

- Getting Started: Installation, Quick Start
- Guides: Tables, Trees & Graphs, NOM Progress, CQRS Architecture, Color Modes, Streaming, Cross-Shape Conversion
- Reference: Format Matrix (full 16-format capability table)
- Community: Changelog, Contributing, Related Tools

### 5. Website — Infrastructure

- OG image generation (`astro-og-canvas`) with cyan-bordered dark theme
- Sitemap generation (`@astrojs/sitemap`)
- JSON-LD structured data (`SoftwareApplication` schema)
- robots.txt, manifest.json, favicon.svg
- Starlight docs with sidebar navigation, search (Pagefind), dark/light toggle
- Landing layout with responsive header, footer, skip-to-content link

### 6. Website — Build Verified

- `npm run build` succeeds: 15 HTML pages, 5 OG images, sitemap, CSP patched 15/15 files
- `npm run typecheck`: 0 errors, 0 warnings, 0 hints (31 files)
- Pagefind search index built

### 7. README — Minimal Improvements

- Added prominent website documentation links at the top
- Fixed module count references (20 -> 19) in two places

---

## b) PARTIALLY DONE

### 1. README Improvements — SHALLOW

I only made 3 edits (links + module count). The 859-line README still has deeper issues:

- The CQRS section code examples have a **pointer dereference bug** (see section d)
- No table of contents for easy navigation
- The "Development > Go toolchain" section is misleading — bare `go build ./...` doesn't build all modules (AGENTS.md says "always use the flake apps")
- No badges for Go Reference, Go Report Card are broken/stale potentially
- The migration guide references v0.23.x but the latest is v0.30.x — still relevant but could be clearer

### 2. API Accuracy in Website Docs — UNVERIFIED

I wrote code examples from the README and AGENTS.md, **not from verifying against the actual Go source**. Critical issue found during post-hoc verification:

- `NewGraphNode` returns `*GraphNode` but `AddNode` takes a `GraphNode` value — code needs dereference `*output.NewGraphNode(...)`. This bug exists in the pre-existing README AND I propagated it into the website.

### 3. Package Dependency Health

- Used `--legacy-peer-deps` to work around `astro-og-canvas@0.11.1` peer dep conflict with Astro 6.x
- 3 npm audit vulnerabilities (1 low, 1 moderate, 1 high) — not investigated
- The `overrides` section in package.json may need updating

---

## c) NOT STARTED

1. **GitHub social preview image** — no custom social card set (uses default GitHub OG)
2. **Website CI/CD** — no GitHub Actions workflow for auto-deploy on push
3. **Website deployment** — not deployed to Firebase (domain `go-output.lars.software` doesn't resolve yet)
4. **Firebase hosting target** — `go-output` target not verified to exist in the `lars-software` Firebase project
5. **README table of contents** — not added
6. **Website analytics** — no analytics integrated
7. **Dependents/Who Uses page** — gogenfilter has one, we don't
8. **HTML validation** — `html-validate` is in devDeps but no `test:html` script in package.json
9. **Link checking** — no automated check for broken links in built site
10. **Lighthouse/performance audit** — not run on the built site
11. **Accessibility audit** — not run (though basic a11y patterns are in place)
12. **OG image visual verification** — images generated but not visually inspected
13. **Mobile responsive testing** — not done beyond code patterns
14. **`go-output.lars.software` DNS** — not verified/configured

---

## d) TOTALLY FUCKED UP

### 1. CRITICAL: Code Examples Have Pointer Dereference Bug

**The actual API:**

```go
func NewGraphNode(id, label string) *GraphNode  // returns POINTER
func (m *GraphBuilder) AddNode(node GraphNode)   // takes VALUE
```

**What I wrote in the website (and what the README already had):**

```go
b.AddNode(output.NewGraphNode("compile", "Compile"))  // WRONG: passing *GraphNode to GraphNode
```

**Correct code should be:**

```go
b.AddNode(*output.NewGraphNode("compile", "Compile"))  // dereference
```

This appears in:

- `website/src/data/hero-code.ts` (the hero code snippet)
- `website/src/components/HeroSection.astro` (the highlighted code)
- `website/src/content/docs/getting-started/quick-start.mdx`
- `website/src/content/docs/guides/cqrs.mdx`
- `website/src/content/docs/guides/trees-and-graphs.mdx`
- `README.md` (pre-existing bug, lines 92-93)

This is a **public-facing documentation bug** — anyone copy-pasting the hero example gets a compile error.

### 2. CRITICAL: RenderTable Return Type Wrong in README

The actual API:

```go
func RenderTable(data *Table, format Format, opts RenderOptions) error  // returns ERROR only
```

The README at line 75-78 implies it produces output but doesn't show where output goes. The website docs reference it correctly in some places but the API itself is confusing (returns error, writes to... what? No writer parameter visible).

### 3. MODERATE: Used `--legacy-peer-deps` Instead of Fixing Dependencies

This is a band-aid. The real fix is either:

- Upgrade `astro-og-canvas` to a version that supports Astro 6.x/7.x
- Or remove the OG image feature if it's incompatible
- Or pin Astro to a version that works

The `package-lock.json` now has this hack baked in.

### 4. MINOR: No `package-lock.json` Commit Verification

The lockfile exists but I didn't verify it's in a committable state (no `.npmrc` with `legacy-peer-deps=true` setting documented).

---

## e) WHAT WE SHOULD IMPROVE

### Immediate ( correctness )

1. Fix ALL pointer dereference bugs in code examples (website + README)
2. Verify every code example in every doc page compiles against the real API
3. Add `.npmrc` with `legacy-peer-deps=true` to the website (or fix the root cause)
4. Verify `RenderTable` usage is correctly documented everywhere

### Quality

5. Add a README table of contents (it's 859 lines with no navigation)
6. Add `test:html` script using `html-validate` (already a devDep)
7. Add a link-checker script for the built site
8. Run Lighthouse audit and fix findings
9. Create a GitHub social preview image (1200x630)
10. Add `website/dist/` and `website/node_modules/` to root `.gitignore` if not covered by nested `.gitignore`

### Content

11. Add an API reference page (auto-generated from Go source, or curated)
12. Add a "Playground" page where users can try formats interactively
13. Add a benchmarks/comparison page
14. Add screenshots/GIFs of the NOM progress visualization
15. Add screenshots of the TUI
16. Write a "Why NOM?" blog-style guide
17. Add a migration guide page for v0.23.x -> v0.30.x (currently only in README)
18. Add more cross-shape conversion examples with visual diagrams

### Infrastructure

19. Create GitHub Actions workflow for website build + deploy
20. Set up Firebase hosting target `go-output` in Firebase console
21. Configure DNS for `go-output.lars.software`
22. Add `pre-commit-config.yaml` for the website (html-validate, prettier)
23. Add the website to the root project's `flake.nix` as a sub-flake or document the separate flake

### README-specific

24. Restructure README with clearer sections and TOC
25. Add a "Gallery" section with rendered output examples (screenshots)
26. Move migration guide to the website, keep README lean
27. Fix "Go toolchain" section to mention `GOEXPERIMENT=jsonv2`
28. Add a "Comparison with alternatives" section
29. Trim the README — 859 lines is too long for a landing page; move deep content to website

---

## f) NEXT 50 THINGS TO DO

### P0 — Fix correctness issues (do first)

1. Fix pointer dereference in `website/src/data/hero-code.ts`
2. Fix pointer dereference in `website/src/components/HeroSection.astro`
3. Fix pointer dereference in `website/src/content/docs/getting-started/quick-start.mdx`
4. Fix pointer dereference in `website/src/content/docs/guides/cqrs.mdx`
5. Fix pointer dereference in `website/src/content/docs/guides/trees-and-graphs.mdx`
6. Fix pointer dereference in `README.md` lines 92-93
7. Verify `RenderTable` documentation accuracy in README + website
8. Verify all other code examples against actual Go source signatures
9. Add `.npmrc` with `legacy-peer-deps=true` to website
10. Rebuild and verify website after fixes

### P1 — Deploy and connect

11. Create Firebase hosting target `go-output`
12. Deploy website to Firebase
13. Configure DNS for `go-output.lars.software`
14. Verify homepage URL resolves
15. Create GitHub Actions workflow for website CI/CD
16. Add website deployment to push workflow

### P2 — README overhaul

17. Add table of contents to README
18. Fix the "Go toolchain" development section (mention GOEXPERIMENT=jsonv2)
19. Trim README to essentials, move deep content to website
20. Add rendered output examples (not just code blocks)
21. Add a "Comparison" section
22. Review every code example in README for compile correctness
23. Update migration section to reference website page

### P3 — Website content enrichment

24. Add screenshots/GIFs of NOM progress visualization
25. Add TUI screenshots
26. Add interactive playground page
27. Add benchmarks page
28. Add API auto-reference page
29. Add "Who uses go-output" / dependents page
30. Write NOM deep-dive guide
31. Add architecture diagrams (using go-output's own D2 renderer!)
32. Add a "Format Gallery" page with rendered examples side-by-side

### P4 — Website quality

33. Add `test:html` script to package.json
34. Run html-validate and fix all issues
35. Add link-checker (lychee or similar)
36. Run Lighthouse audit
37. Fix Lighthouse findings
38. Accessibility audit (axe-core)
39. Mobile responsive testing
40. Visual OG image inspection
41. Add `prefetch` strategy tuning
42. Add Open Graph image per doc page (already set up, verify quality)

### P5 — GitHub polish

43. Create GitHub social preview image (1200x630)
44. Set up GitHub Discussions (if not already)
45. Create release notes template
46. Add funding/sponsorship link
47. Review and update topic tags for SEO
48. Create GitHub Actions badge for website deploy status
49. Add `website/` to repo description or pinned docs
50. Verify pkg.go.dev documentation renders correctly

---

## g) TOP 2 QUESTIONS I CANNOT ANSWER MYSELF

### Q1: Does the Firebase hosting target `go-output` exist in the `lars-software` project?

The `.firebaserc` references `targets.lars-software.hosting.go-output: ["go-output"]`, but I have no way to verify this target actually exists in the Firebase console. The gogenfilter and go-atomic-write projects each have their own targets in the same Firebase project — does `go-output` need to be created manually? **Blocking for deployment.**

### Q2: Is the domain `go-output.lars.software` configured in DNS?

I set the homepage URL to `https://go-output.lars.software` in both GitHub and the website config, but I cannot verify the DNS record exists or that Firebase has been authorized for this subdomain. If it's not configured, the GitHub homepage link will 404. **The pattern from gogenfilter (`gogenfilter.lars.software`) suggests it's a wildcard `*.lars.software` CNAME, but I cannot confirm.**

---

## Summary Scorecard

| Surface         | Status      | Score    | Key Issue                                      |
| --------------- | ----------- | -------- | ---------------------------------------------- |
| GitHub metadata | Done        | 8/10     | Homepage not verified, no social preview       |
| Website         | Functional  | 6/10     | Code examples have pointer bug, not deployed   |
| README          | Touched     | 4/10     | Only 3 edits, no deep improvement              |
| **Overall**     | **Partial** | **6/10** | **Correctness issues in public code examples** |

