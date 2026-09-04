# TODO_LIST.md — go-output

**Last updated:** 2026-09-04
**Open items:** 7

---

## Website & Release Engineering (from the 2026-09-03 outage response)

| #  | Task                                                                                                                                                                                                                            | Effort   | Source                                            |
| -- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------------------------------------------- |
| 1  | **Produce the demo video** — 20–30s HyperFrames piece (hook 0-3s pain, value claim, 2-3 evidence scenes, install CTA), render to `public/demo.mp4` + poster + og:image, add ShowcaseSection to the landing page, README deep-link. Full spec: website-launch skill `references/demo-video.md`; outage report item 36. | 30–60 min | `docs/status/2026-09-03_23-07_website-outage*.md` (f)36 |
| 2  | **Deep accessibility pass** — computed contrast ratios on both themes, full keyboard-nav walk, focus-trap check on search dialog. Baseline (lang, skip-link, named controls, single h1, focus-visible) already verified green via CDP. | 45 min   | Outage report (f)33                                |
| 3  | **Lighthouse perf follow-up** — headless-software-rendered TBT is inflated by the 120px hero blur filters; evaluate pre-rendered radial gradients or reduced blur radius on mobile breakpoints. Perf was 70 / a11y 100 / BP 100 / SEO 100.   | 30 min   | 2026-09-04 session, Lighthouse JSON                |
| 4  | **Launch or de-list fleet stragglers** — `cmdguard.lars.software` (TLS cert altname invalid: domain not attached in Firebase) and `typespec-asyncapi.lars.software` (no DNS record). They are excluded from `uptime.yml` SITES until launched; their entries are commented in the workflow. | 1–2 h    | `uptime.yml` comments; 2026-09-04 uptime runs      |
| 5  | **Dependabot PR triage** — the astro 7.2.10 bump PR reproduces the outage-adjacent manifest drift and correctly fails CI; close it or merge deliberately with a regenerated lockfile after the canvaskit question settles. canvaskit-wasm bumps are now permanently ignored via dependabot config. | 10 min   | Dependabot runs 2026-09-04 10:45 UTC               |
| 6  | **Fleet-level header audit** — remaining 12 sibling sites vs the go-output/gogenfilter hardened `firebase.json` baseline (HSTS preload, CORP/COOP, cache rules). go-output + gogenfilter verified identical 2026-09-04.                | 1 h      | Outage report (f)49                                |
| 7  | **Cut `v1.0.0`** — API frozen per ADR 006; all v0.30.x–v0.37.x breaking changes shipped. Run `scripts/pre-tag-check.sh v1.0.0`, then `scripts/tag-release.sh v1.0.0` + push. **Owner decision** (public API-freeze declaration).          | 15 min   | Prior TODO_LIST #10                                |

## Community (owner-dependent)

| # | Task                                        | Effort | Status                     |
| - | ------------------------------------------- | ------ | -------------------------- |
| 8 | **Post to r/golang, submit to Awesome Go**  | 30 min | Open (needs owner account) |
