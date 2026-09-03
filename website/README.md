# go-output Website

Astro + Starlight documentation site, deployed to Firebase Hosting (project
`lars-software`, hosting target/site `go-output`):

- https://go-output.lars.software (custom domain)
- https://go-output.web.app (Firebase default — authoritative right after deploy)

## Toolchain

Node 24 + pnpm 11 (pinned via `packageManager` in `package.json`). On this
machine neither is on PATH — always run through nix-shell:

```bash
nix shell nixpkgs#nodejs -c pnpm <command>
```

`CI=true` makes pnpm default to `--frozen-lockfile`; keep it set for
deterministic installs.

## Build & gates (local)

Run from this directory, in order. Each step is a CI gate — a local pass
predicts CI exactly.

```bash
# 1. Install (clean checkout / lockfile already in sync)
CI=true nix shell nixpkgs#nodejs -c pnpm install --frozen-lockfile

# 1b. Only when the manifest changed and the lockfile needs regenerating:
CI=true nix shell nixpkgs#nodejs -c pnpm install --no-frozen-lockfile

# 2. Typecheck
CI=true nix shell nixpkgs#nodejs -c pnpm run typecheck

# 3. Build (includes the fix-csp.mjs post-build step)
CI=true nix shell nixpkgs#nodejs -c pnpm run build

# 4. HTML validation
CI=true nix shell nixpkgs#nodejs -c pnpm exec html-validate --config .htmlvalidate.json "dist/**/*.html"

# All four in one target:
CI=true nix shell nixpkgs#nodejs -c pnpm run verify
```

**pnpm v11 traps** (cost real debugging time once — don't repeat them):

- Plain `pnpm install` in an interactive shell aborts on the purge prompt
  (no TTY for confirmation). Always set `CI=true`.
- npm-style top-level `overrides` in `package.json` are silently IGNORED by
  pnpm. Dependency pins go in `package.json` directly (exact versions) or
  `pnpm-workspace.yaml`.
- The `allowScripts` field in `package.json` is non-standard; pnpm's native
  mechanism is `onlyBuiltDependencies` in `pnpm-workspace.yaml` (that is
  what this repo uses).

**canvaskit-wasm must stay a DIRECT dependency** with the exact pin
`0.41.1`. `astro-og-canvas` pulls `canvaskit-wasm@^0.42.0` transitively,
whose emscripten Node build references `__dirname` and crashes Astro's ESM
prerender on `/og/*.png` routes with a misleading error blaming pnpm
hoisting. Installing the known-good version directly (and exactly) is the
documented remedy. Re-widen only after upstream fixes the ESM issue.

## Deploy

### CI (preferred)

`.github/workflows/website.yml` builds + deploys on every push to `master`
that touches `website/**`, and `.github/workflows/release.yml` redeploys on
every root version tag. It needs the `FIREBASE_SERVICE_ACCOUNT` repo secret
(firebase-adminsdk key for `lars-software`).

### Manual

```bash
nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c \
  firebase deploy --only hosting:go-output --project lars-software
```

Run all four gates before a manual deploy. CI runs them for you.

## Verify after deploy

1. `https://go-output.web.app/` returns 200 — proves the new release is
   actually live (the custom domain rides the same release but adds CDN
   variables).
2. `https://go-output.lars.software/` and `/format-matrix/` return 200 with
   real content (not Firebase's "Site Not Found" page — that is a 404).
3. Screenshot spot-check: `nix shell nixpkgs#chromium -c chromium --headless
   --no-sandbox --screenshot=/tmp/site.png --window-size=1440,900 <url>`

CI does 1–2 automatically (smoke step) and fails the run on mismatch.

## Rollback

```bash
nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c \
  firebase hosting:releases:list --site go-output --project lars-software

nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c \
  firebase hosting:rollback --site go-output --project lars-software
```

## Analytics

Intentionally none — no trackers, no cookies, no telemetry. The newsletter
form is the only data collection and it posts to its own service. Keep it
this way unless there is a concrete product reason not to.

## Custom domain & SSL

`go-output.lars.software` is a CNAME to `go-output.web.app` (Terraform in
the `domains` repo), with an ACME TXT record for certificate issuance.
Domain/cert state is queryable via the Firebase customDomains REST API (see
the website-launch skill's `firebase-rest-api.md` reference).

**HSTS preload is set** (`max-age=63072000; includeSubDomains; preload`).
If certificate renewal ever fails, the domain does not just downgrade — it
bricks for all modern browsers (preload = no click-through). If the custom
domain ever serves cert errors, treat it as P0: check the ACME TXT record
first, then cert state via the REST API.

## History note

The site "launched" on 2026-07-13 but never actually shipped a release — a
hand-bumped manifest broke the lockfile, rebuilds became impossible, and the
custom domain served Firebase's "Site Not Found" from birth until the
2026-09-03 incident response rebuilt the pins, deployed the first working
release, and added the CI pipeline + uptime monitor that now guard it. Full
post-mortem:
`docs/status/2026-09-03_23-07_website-outage-root-cause-and-redeploy.md`.
