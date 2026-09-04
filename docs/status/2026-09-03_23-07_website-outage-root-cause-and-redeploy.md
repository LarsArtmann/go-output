# Status Report: Website Outage Root Cause & Redeploy — go-output.lars.software

**Date:** 2026-09-03 23:07 CEST
**Session scope:** `https://go-output.lars.software/format-matrix/` returned Firebase's "Site Not Found" page. Diagnosed, fixed, rebuilt, redeployed, verified.
**Trigger:** User report. **Skill:** website-launch (loaded and followed; go-live runbook + firebase-rest-api references).

---

## Executive Summary

The site was down because the Firebase Hosting site `go-output` had **zero deployed releases** — not because of DNS or domain misconfiguration (both were already correct and fully active). Getting a release out required fixing a website build that had been silently broken since a hand-edited `package.json` (out-of-sync lockfile, incompatible dependency pins). The site is now live and verified end-to-end. **The systemic gap remains: there is no CI/CD pipeline for the website, which is exactly why the deploy was forgotten.**

| Category          | Count | Summary                                                 |
| ----------------- | ----- | ------------------------------------------------------- |
| Fully done        | 8     | Diagnosis, build fix, gates, deploy, verification, docs |
| Partially done    | 4     | Launch DoD, commit discipline, dep refresh, lockfile    |
| Not started       | 6     | CI/CD, metadata, demo video, dead-config cleanup, more  |
| Totally fucked up | 2     | The original hand-bump; my 5-attempt debugging detour   |

---

## Timeline (what actually happened)

1. User reported "Site Not Found" on `/format-matrix/`. _(✅ done 2026-09-04 — `.github/workflows/website.yml` (frozen install, astro check, build, html-validate, deploy + dual-domain smoke))_
2. Loaded `website-launch` skill → `website/` exists → maintenance mode. _(✅ done 2026-09-04 — `FIREBASE_SERVICE_ACCOUNT` secret created from firebase-adminsdk key; `gh secret list` confirms)_
3. **DNS verified correct:** `go-output.lars.software` → CNAME `go-output.web.app.` → 199.36.158.100 (local + 1.1.1.1). _(✅ done 2026-09-04 — smoke step in website.yml asserts 200 + content marker on both domains)_
4. **Firebase verified correct:** site `go-output` exists in `lars-software`; custom domain `go-output.lars.software` attached with `HOST_ACTIVE` / `OWNERSHIP_ACTIVE` / `CERT_ACTIVE` (via `customDomains` REST API); `_acme-challenge.go-output` TXT record resolves and matches the Terraform-staged value. _(✅ done 2026-09-04 — frozen install is the FIRST CI step)_
5. **Root cause isolated:** `https://go-output.web.app/` also 404s → the hosting site has **no releases at all**. Firebase serves the "You haven't deployed an app yet" page under exactly this condition. _(✅ done 2026-09-04 — release.yml `deploy-website` job calls website.yml via workflow_call on every root v* tag)_
6. Local `dist/` (Aug 4) predated today's `website/src/` touch → rebuild required before deploy. _(✅ done 2026-09-04 — `.github/workflows/uptime.yml` pings 14 launched `*.lars.software` sites every 30 min, deduplicated [uptime] issue)_
7. Rebuild hit a wall of pnpm v11 + dependency-pin failures (details in section d/e). _(✅ done 2026-09-04 — `website/README.md` full runbook)_
8. Build fixed → gates run → `firebase deploy --only hosting:go-output` → release complete. _(✅ done 2026-09-04 — `scripts/pre-deploy-check.sh`, all gates verified green)_
9. Verified: `.web.app` 200, custom domain 200 with full Format Matrix content, landing page 200, headless-Chromium screenshot of the live page visually QA'd (dark theme, sidebar, matrix table all correct). _(✅ done 2026-09-04 — AGENTS.md "Website" section)_
10. Gotchas recorded in `AGENTS.md`; all changes auto-committed by the daemon (final state `e12f3b1`, working tree clean). _(✅ done 2026-09-04 — SSL/preload note in website/README.md; lars.software confirmed `preloaded` via hstspreload API)_

---

## a) FULLY DONE

1. **Outage root-cause diagnosis** — eliminated DNS, Terraform, domain attachment, and SSL cert as suspects with hard evidence (REST API states, `nslookup` of CNAME + ACME TXT); isolated "no release deployed" as the cause.
2. **Custom-domain health verification** — `customDomains` REST API shows `HOST_ACTIVE`/`OWNERSHIP_ACTIVE`/`CERT_ACTIVE`; ACME TXT in DNS matches `lars.software.tf` staging.
3. **Website build repaired** — `pnpm-lock.yaml` re-synced; `esbuild` build script approved (fixed the literal `set this to true or false` placeholder left in `pnpm-workspace.yaml` by an earlier session).
4. **canvaskit `__dirname` prerender crash fixed** — `astro-og-canvas` pinned to `0.13.0` and `canvaskit-wasm@0.41.1` added as a **direct** dependency (the documented pnpm remedy); OG images (`/og/*.png`, incl. `format-matrix.png`) now generate.
5. **`astro check` crash fixed** — `typescript` pinned from broken `7.0.2` to `6.0.3` (language-server compatible).
6. **Dead config removed** — `param: "slug"` (no longer part of the `OGImageRoute` API in 0.13.x) deleted from `src/pages/og/[...slug].ts`, clearing the last type error.
7. **All quality gates green** — `astro check`: 0 errors / 0 warnings / 0 hints; `astro build`: 15 pages + sitemap + pagefind index; `fix-csp.mjs`: 15/15 patched; `html-validate`: exit 0.
8. **Deploy + end-to-end verification** — upload endpoint pre-checked, `firebase deploy --only hosting:go-output` release complete; `/format-matrix/` and `/` verified 200 with full content on **both** `go-output.web.app` and `go-output.lars.software`; visual QA via headless screenshot of the live page.
9. **Knowledge captured** — AGENTS.md Gotchas entry covering the exact pins, the misleading error text, and the three pnpm v11 traps (`CI=true` forces frozen-lockfile; top-level `overrides` ignored; no-TTY abort).

## b) PARTIALLY DONE

1. **"Website launch" Definition of Done (skill Phase 6–7)** — the site itself is live (the core), but GitHub metadata verification, demo video (§3.11), and CI/CD phases were never executed this session; several were never executed at launch either.
2. **Commit discipline** — all work IS committed (working tree clean, `e12f3b1`), but by the auto-commit daemon in five heuristic `chore: auto-commit N file(s)` chunks interleaved mid-debugging (e.g. lockfile churn split across `96688df`/`a0e2b90`) — not atomic, intentionally-worded commits. History is functional but noisy.
3. **Today's repo-wide "toolchain + dependency refresh" commit** — its `website/` portion is now effectively reverted by my exact pins (astro `7.2.1`, og-canvas `0.13.0`, canvaskit `0.41.1`, typescript `6.0.3`). The Go-side refresh stands; the website-side intent (newer majors) is deferred, not accomplished.
4. **AGENTS.md accuracy** — updated with the new gotcha, but the Commands section still doesn't document the website build/deploy flow (`pnpm` via nix shell, `firebase deploy` invocation) — only the gotcha was added.

## c) NOT STARTED

1. **Website CI/CD pipeline** (skill Phase 7) — no `website-deploy` workflow exists in `.github/workflows/` (only root `ci.yml` + `release.yml`). Deploys are 100% manual. **This is the systemic root cause of the outage recurring.**
2. **`FIREBASE_SERVICE_ACCOUNT` GitHub secret** — never created for this repo (Phase 7 Step 1 prerequisite for #1).
3. **Demo video** (§3.11) — landing page has no `demo.mp4` / ShowcaseSection; the launch shipped without the skill's default sales centerpiece.
4. **GitHub repo metadata check** (Phase 6) — homepage/description/topics for `LarsArtmann/go-output` not verified against the live site this session.
5. **Dead-config cleanup** — the npm-style top-level `overrides` block in `website/package.json` (`brace-expansion`, `devalue`, `esbuild`, `yaml`) is ignored by pnpm v11; still sitting there misleading future readers (I added to it, proved it dead, and removed only my own addition).
6. **TODO_LIST.md harvest** — section (f) of this report has not been routed into `TODO_LIST.md`/`ROADMAP.md` (docs-health HARVEST).

## d) TOTALLY FUCKED UP

1. **The original sin (pre-existing, caused the outage):** an earlier session hand-bumped `website/package.json` ranges (astro `^7.2.10`, typescript `^7.0.2`, …) **without regenerating the lockfile and without ever running a build**. Result: the lockfile violated the manifest (`--frozen-lockfile` impossible), the site could not be rebuilt, and — because deploys were manual anyway — nobody noticed the site had _never once_ been redeployed since launch. The site's `dist/` from Aug 4 had **no `og/` directory**, meaning a successful full build (og route included) may never have happened post-launch at all. Silent, compounding, and only surfaced when a user hit the dead URL.
2. **My debugging detour (this session, ~5 wasted cycles):** the build error message **explicitly printed the fix** ("install `canvaskit-wasm` directly: `pnpm add canvaskit-wasm`") from the very first failure. Instead of following the printed remedy immediately, I went on a version-archaeology spree: overrides pin (dead — pnpm v11 ignores npm-style top-level `overrides`), astro-og-canvas downgrade, astro downgrade (disproved my own hypothesis), `vite.ssr.external` experiment (no effect, reverted). Five attempts before doing what the error said. This violated my own error-handling protocol (read the complete error first, follow the documented remedy, _then_ experiment). Cost: roughly 15–20 minutes and 4 avoidable install/build cycles. Mitigating: each failed attempt did eliminate a hypothesis cleanly and the final pin set is deliberate — but the order of operations was wrong.

**Not fucked up (verified clean):** DNS/Terraform records, Firebase domain/cert state, the deployed site's rendering (screenshot-verified), the Go modules (untouched), HSTS/security headers (in `firebase.json`, live).

## e) WHAT WE SHOULD IMPROVE

1. **Deploy automation** — the outage class is "manual step forgotten." A path-filtered GitHub Actions workflow (build on `website/**` change → deploy to Firebase) eliminates it.
2. **Error-reading discipline** — when a tool prints a remedy, try the remedy before designing experiments. Cheapest signal first.
3. **Post-fix verification breadth** — I verified desktop dark-mode rendering only; mobile/light/console-CSP checks are cheap add-ons to the same screenshot loop.
4. **Lockfile hygiene** — manifest changes must always ship with a regenerated lockfile in the same commit; a CI `pnpm install --frozen-lockfile` gate enforces this for free.
5. **Pin hygiene** — exact pins are correct for the website, but they need a scheduled re-test (upstream `canvaskit-wasm` will eventually fix the ESM `__dirname` bug; the pins should be re-widened deliberately, not by accident).
6. **Config honesty** — dead config (ignored `overrides`, non-standard `allowScripts`) is worse than no config; delete or migrate to `pnpm-workspace.yaml` where pnpm v11 actually reads it.
7. **Deploy runbook as code** — the deploy knowledge currently lives in a skill + this session; a short `website/README.md` (build, gates, deploy, verify) makes it survivable without the skill.
8. **Monitoring** — the outage was discovered by a user, not by tooling. Uptime monitoring on `go-output.lars.software` would have caught a site that was _born dead_ at launch.
9. **Commit granularity** — when the auto-daemon is active, committing intentional milestones myself (skill checkpoint discipline) keeps history reviewable; heuristic chunks interleaved with experiments are hard to bisect.
10. **Stale LSP noise** — vtsls reported phantom syntax errors on `astro.config.mjs` (byte-verified clean, build-proven valid) for the rest of the session; restarting the LSP after config edits avoids trusting/ignoring diagnostics ambiguity.

## f) UP TO 50 THINGS TO GET DONE NEXT

_Pareto note: items 1–10 carry most of the value; 11–50 are ROADMAP-fuel brainstorm and need HARVEST routing rigor._

**CI/CD & reliability (the outage class)**

1. Add `website-deploy.yml` GitHub Actions workflow: build (pnpm frozen install, astro check, build, html-validate) → deploy to Firebase, path-filtered on `website/**`.
2. Create `FIREBASE_SERVICE_ACCOUNT` secret (firebase-adminsdk key for `lars-software`) via `gh secret set`.
3. Add a post-deploy smoke step in CI: fetch `/` and `/format-matrix/`, assert HTTP 200 + content marker.
4. Add `pnpm install --frozen-lockfile` as the FIRST CI step so lockfile/manifest drift fails loudly.
5. Add root `release.yml` hook (or manual checklist step): redeploy website on version tags if docs changed.
6. Set up uptime monitoring for `go-output.lars.software` (and ideally the other 6 sibling sites on the same pattern).
7. Write `website/README.md`: build/gates/deploy/verify runbook + the exact nix-shell invocations.
8. Add a `pre-deploy-check.sh` analog of `scripts/pre-tag-check.sh` (frozen install, typecheck, build, html-validate, optional local preview fetch).
9. Document the website deploy command in AGENTS.md Commands section.
10. Add SSL-renewal note + periodic cert-state check (HSTS preload means a renewal failure bricks the domain).

**Build & dependency health**
11. Delete the dead npm-style `overrides` block from `website/package.json` (or migrate effective pins to `pnpm-workspace.yaml` `overrides:`). _(✅ done 2026-09-04 — dead `overrides` removed; build green without them)_
12. Audit the non-standard `allowScripts` field — verify pnpm v11 actually honors it; migrate to `onlyBuiltDependencies` if not. _(✅ done 2026-09-04 — `allowScripts` removed; `allowBuilds` confirmed as the real pnpm 11 field (pnpm.io/settings))_
13. Schedule a deliberate pin re-widening attempt (astro >7.2.1, canvaskit >0.41.1) once upstream fixes the emscripten ESM `__dirname` issue; until then keep exact pins. _(✅ routed 2026-09-04 — ROADMAP "Dependency pin re-widening"; dependabot ignores canvaskit-wasm)_
14. Add a tiny canary script `pnpm run verify` = typecheck + build + html-validate in one target. _(✅ done 2026-09-04 — `pnpm run verify` target)_
15. Add renovate/dependabot config coverage confirmation for `website/` (repo has `renovate.json` — verify it reaches the pnpm manifest). _(✅ done 2026-09-04 — report claimed renovate.json exists; it does not. dependabot npm ecosystem added for website/ instead)_
16. Pin the pnpm version story (dev machine vs CI vs daemon all use `packageManager: pnpm@11.20.0` — verify CI honors it via corepack). _(✅ done 2026-09-04 — CI enables corepack pinned to pnpm@11.20.0 before setup-node)_
17. Document Node 24 / nix-shell requirement for website builds in `website/README.md`. _(✅ done 2026-09-04 — website/README.md Toolchain section)_
18. Consider a `flake.nix` app `.#website-build` (and `.#website-deploy`) for tooling parity with the Go apps. _(✅ done 2026-09-04 — flake apps `.#website-build` + `.#website-deploy`, build verified)_

**Content truth & docs**
19. Audit the site's Module Map table against the current 19-module reality (site lists 14 rows; missing `bdd/`, `integration/`, `examples/`, `testhelpers/graphtest/` is fine to omit — but verify nothing stale claims deleted things). _(✅ done 2026-09-04 — module map audited: 15 public modules correct, test-only modules annotated)_
20. Re-verify every code example on the site against current source (skill's #1 documented failure mode; I did NOT do a content audit this session — only build/serve verification). _(✅ done 2026-09-04 — every Go example API-verified vs source; fixed 2× `tree.NewTreeRendererFromTable` → `tree.TreeRendererFromTable` (compile-error drift))_
21. Verify Format Matrix page claims vs `docs/FORMAT_ARCHITECTURE.md` (16 formats × 3 shapes matrix unchanged since v0.38 review?). _(✅ done 2026-09-04 — matrix was WRONG: d2/mermaid/dot/plantuml DO support tree (shape.go); fixed + example corrected)_
22. Check `/og/*.png` are actually referenced in page `<meta property="og:image">` tags (images build; reference not verified). _(✅ done 2026-09-04 — og:image meta present; /og/home.png live 200)_
23. Verify sitemap + robots.txt served correctly on the live domain. _(✅ done 2026-09-04 — sitemap-index.xml, sitemap-0.xml, robots.txt all 200)_
24. Verify the `firebase.json` `/docs/:path*` → `/:path*` 301 redirect is still needed (possible leftover from a pre-launch URL structure). _(✅ done 2026-09-04 — /docs/* 301 verified working; kept as legacy-bookmark courtesy (zero cost))_
25. Verify cleanUrls/trailing-slash canonical behavior on the live site (`cleanUrls: true` + `trailingSlash: false`). _(✅ done 2026-09-04 — trailing-slash URLs 301 to canonical non-slash; consistent)_
26. Verify cache headers on live assets (immutable for hashed, must-revalidate for HTML) with a header dump. _(✅ done 2026-09-04 — BUG fixed: `**/*.html` rule never matched clean URLs (pages cached max-age=3600). Cache-Control moved to catch-all must-revalidate, assets override immutable)_
27. Test the custom 404 page on the live domain (dist has one; serving behavior unverified). _(✅ done 2026-09-04 — /nonexistent 404 + correct cache header)_
28. Link-check all outbound links (pkg.go.dev, GitHub, newsletters). _(✅ done 2026-09-04 — all 37 external links fetched OK)_
29. Smoke-test pagefind search UI on the deployed site (index built; UI unverified). _(✅ done 2026-09-04 — pagefind.js/jsontools live 200; search UI assets verified)_
30. Mobile-viewport screenshot QA (only 1440×900 dark checked). _(✅ done 2026-09-04 — REAL BUG fixed: hero overflowed horizontally at 390px (flex min-width:auto + wide pre); CDP-verified scrollW=390 after min-w-0 fix)_
31. Light-theme screenshot QA (dark is default; picker verified in HTML only). _(✅ done 2026-09-04 — light theme screenshot via CDP `html.light`: correct palette, contrast, accent adaptation)_
32. Browser-console CSP violation check on real load (fix-csp patched 15/15; runtime violations unverified). _(✅ done 2026-09-04 — zero CSP console violations on live site (sole console error is a chrome-extension artifact))_
33. Accessibility pass: contrast, focus states, keyboard nav (skip-link exists). _(✅ baseline done 2026-09-04 — lang/skip-link/labels/h1/focus-visible green; deep contrast+keyboard pass in TODO_LIST #2)_
34. Lighthouse/perf quick pass on the live landing page. _(✅ done 2026-09-04 — Lighthouse: perf 70 (software-render inflation, blur-driven TBT), a11y/BP/SEO 100; follow-up TODO_LIST #3)_
35. Landing page "1 Stars" widget — decide keep/fix/remove (undersells; see question 2). _(✅ decided 2026-09-04 — live count only at ≥10 stars, else "Star on GitHub" CTA (also fixes the "1 Stars" grammar bug))_
36. Demo video (§3.11): produce 20–30s HyperFrames video, add ShowcaseSection + `public/demo.mp4`, poster + og:image, README deep-link. _(routed 2026-09-04: TODO_LIST #1 top item — full production spec loaded; deliberate follow-up)_ 
37. Verify/refresh GitHub repo metadata: description, homepage → `https://go-output.lars.software`, topics (Phase 6). _(✅ done 2026-09-04 — verified: description, homepage, 20 topics all correct)_
38. Verify README documentation-link bar matches the template and live URLs (docs drift). _(✅ done 2026-09-04 — README link bar matches template; all 3 URLs live 200)_
39. Add CHANGELOG entry for the website launch/fix (untouched this session). _(✅ done 2026-09-04 — CHANGELOG Unreleased entries added)_
40. HARVEST: route sections (c)/(e)/(f) of this report into `TODO_LIST.md` (bounded tasks) and `ROADMAP.md` (ideas). _(✅ done 2026-09-04 — harvested into TODO_LIST.md (7 items) + ROADMAP.md)_

**Hygiene & follow-through**
41. Re-start the LSP (stale vtsls phantom errors on `astro.config.mjs` pollute every diagnostics readout). _(✅ done 2026-09-04 — LSP restarted)_
42. Confirm the daemon's `e12f3b1` actually contains the og-route fix + AGENTS.md entry (working tree is clean, but verify content, not just cleanliness). _(✅ done 2026-09-04 — verified: e12f3b1 contains og-route fix + AGENTS.md entry)_
43. Consider squashing/chore-noting the five heuristic auto-commits from this session if history hygiene matters to you (rewrite risk — likely skip). _(✅ decided 2026-09-04 — skip (history-rewrite risk outweighs cosmetic gain, per report’s own lean))_
44. gogenfilter parity check: its `astro-og-canvas@0.13.0` + pnpm v10-era node_modules — does gogenfilter still build TODAY? If it shares the fragility, fix both (the pattern is fleet-wide: 6+ sites). _(✅ done 2026-09-04 — gogenfilter Go builds; website hardened (direct canvaskit pin + real allowBuilds value), 18 pages build green)_
45. Add a fleet-level note in the website-launch skill (via skill-creator feedback) about canvaskit direct-dependency remedy + pnpm v11 traps, so the next sibling site doesn't repeat this session. _(✅ done 2026-09-04 — common-pitfalls.md #31/#32 added (canvaskit direct pin, pnpm v11 traps))_
46. Decide analytics stance (none present today — is that intentional?). _(✅ decided 2026-09-04 — no analytics, intentional; documented in website/README.md)_
47. Add `Captures "last deployed at <date> (vX.Y.Z)"` footer stamp to make staleness visible on the page itself. _(✅ done 2026-09-04 — footer stamps "Last deployed YYYY-MM-DD" (+ version when GO_OUTPUT_VERSION set))_
48. Verify HSTS preload submission state for `lars.software` (headers include `preload` — is the domain actually submitted? preload is near-irreversible). _(✅ done 2026-09-04 — hstspreload API: lars.software status "preloaded")_
49. Review whether `Strict-Transport-Security`/header config in `firebase.json` matches the gogenfilter hardened baseline (CSP hash injection exists here — good; compare remaining headers). _(✅ done 2026-09-04 — headers identical to gogenfilter baseline (HSTS/XFO/CORP/COOP); full fleet audit in TODO_LIST #6)_
50. Post-CI/CD: re-annotate this report (docs-health ANNOTATE) with what landed. _(✅ done 2026-09-04 — this annotation)_

## g) QUESTIONS I CANNOT FIGURE OUT MYSELF

1. **Deploy automation intent:** Should I build the `website-deploy.yml` CI pipeline now (requires creating the `FIREBASE_SERVICE_ACCOUNT` GitHub secret for `lars-software` — an action only you can authorize), or do you deliberately want website deploys to stay manual/release-time? This decides whether the outage class is closed or just survived.
2. **The "1 Stars" landing widget:** the landing page renders a live GitHub star count that currently reads "1 Stars". Keep (honest), remove (undersells), or replace with a static social-proof element? Judgment call about public perception — not inferable from code.
3. **Intent of the earlier package.json hand-bump:** astro `^7.2.10` / typescript `^7.0.2` were staged in the manifest but never built or locked. Was that an upgrade you want re-attempted deliberately later (e.g. once canvaskit fixes the ESM bug), or was it accidental drift that the exact pins should permanently overwrite? Determines whether items 3/13 are "restore intent" or "close forever".

---

## Session Artifacts

- Commits (daemon, chronological): `5735e10` → `e6b8bd8` → `96688df` → `a0e2b90` → `e12f3b1` (working tree clean at 23:07)
- Key files: `website/package.json`, `website/pnpm-lock.yaml`, `website/pnpm-workspace.yaml`, `website/astro.config.mjs` (net-unchanged), `website/src/pages/og/[...slug].ts`, `AGENTS.md` (Gotchas)
- Infra state: release live on `go-output.web.app` + `go-output.lars.software`; DNS/Terraform untouched; Firebase domain CERT_ACTIVE
- Verification: `astro check` 0/0/0 · build 15 pages · html-validate exit 0 · dual-domain 200s · screenshot QA
