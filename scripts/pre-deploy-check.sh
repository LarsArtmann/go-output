#!/usr/bin/env bash
# Pre-deploy verification for the website: run before any manual Firebase
# deploy. Mirrors the CI gates in .github/workflows/website.yml — if this
# passes locally, CI and the deploy will too.
# Exits non-zero if any check fails. Run from the repo root.
#
# Usage: scripts/pre-deploy-check.sh [--skip-install]
#   --skip-install  Skip pnpm install (node_modules already present and
#                   lockfile untouched — e.g. an immediate re-check).
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR/website"

SKIP_INSTALL=false
[ "${1:-}" = "--skip-install" ] && SKIP_INSTALL=true

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

fail() {
	echo -e "${RED}FAIL:${NC} $*" >&2
	exit 1
}
info() { echo -e ":: $*"; }
pass() { echo -e "${GREEN}ok${NC}  $*"; }

run_gate() {
	local name="$1"
	shift
	info "$name"
	if CI=true nix shell nixpkgs#nodejs -c "$@"; then
		pass "$name"
	else
		fail "$name"
	fi
}

[ -f package.json ] || fail "run from the repo root (website/ not found)"
[ -f pnpm-lock.yaml ] || fail "pnpm-lock.yaml missing — the site cannot build reproducibly"

if [ "$SKIP_INSTALL" = false ]; then
	# Gate 1: frozen install — the exact failure mode of the 2026-09-03
	# outage (manifest/lockfile drift made the site unbuildable for weeks).
	run_gate "pnpm install --frozen-lockfile" pnpm install --frozen-lockfile
else
	info "skipping install (--skip-install)"
fi

# Gate 2-4: typecheck, build, HTML validation.
run_gate "astro check" pnpm run typecheck
run_gate "astro build (+fix-csp)" pnpm run build
run_gate "html-validate" pnpm exec html-validate --config .htmlvalidate.json "dist/**/*.html"

# Post-build sanity: the pages whose absence signal a broken build.
for page in dist/index.html dist/format-matrix/index.html dist/404.html; do
	[ -s "$page" ] || fail "build output missing: $page"
done
pass "dist contains index, format-matrix, 404"

if ls dist/og/*.png >/dev/null 2>&1; then
	pass "og images generated (canvaskit prerender OK)"
else
	fail "dist/og/*.png missing — OG route prerender failed (canvaskit-wasm pin broken?)"
fi

echo
echo -e "${GREEN}All website pre-deploy gates passed.${NC} Deploy with:"
echo "  nix shell nixpkgs#nodejs nixpkgs#firebase-tools -c firebase deploy --only hosting:go-output --project lars-software"
