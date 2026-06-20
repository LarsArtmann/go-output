#!/usr/bin/env bash
# Pre-tag verification: run before cutting any release tag.
# Exits non-zero if any check fails. Run from the repo root (where go.work lives).
#
# Usage: scripts/pre-tag-check.sh [vX.Y.Z]
#   If a version is passed, the script also verifies the working tree is clean
#   and that the version tag does not already exist.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

VERSION="${1:-}"
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

fail() { echo -e "${RED}FAIL:${NC} $*" >&2; exit 1; }
info() { echo -e ":: $*"; }
pass() { echo -e "${GREEN}ok${NC}  $*"; }

# ---------------------------------------------------------------------------
# Pre-flight: clean tree + tag-not-present (only when a version is supplied)
# ---------------------------------------------------------------------------
if [ -n "$VERSION" ]; then
  info "validating preconditions for $VERSION"

  if ! git diff-index --quiet HEAD --; then
    fail "working tree is dirty — commit or stash before tagging"
  fi
  pass "working tree is clean"

  if git rev-parse -q --verify "refs/tags/$VERSION" >/dev/null; then
    fail "tag $VERSION already exists"
  fi
  pass "tag $VERSION does not exist yet"
fi

if [ ! -f go.work ]; then
  fail "go.work not found — run scripts/setup-workspace.sh first"
fi
pass "go.work present"

# ---------------------------------------------------------------------------
# Module list (mirrors flake.nix `modules`)
# ---------------------------------------------------------------------------
MODULES=(
  "." "bdd" "d2" "delimited" "enum" "envdetect" "escape" "examples"
  "graph" "integration" "markdown" "markup" "nom" "plantuml"
  "serialization" "table" "testhelpers" "testhelpers/graphtest" "tree" "tui"
)

# Modules that carry concurrency-sensitive code and warrant -race.
RACE_MODULES=("nom" "tui" "integration")

echo
info "go version: $(go version)"

# ---------------------------------------------------------------------------
# 1. Build every module
# ---------------------------------------------------------------------------
echo
info "building all modules (go build ./...)"
for mod in "${MODULES[@]}"; do
  (cd "$mod" && GOWORK=off go build ./...) || fail "build failed in $mod"
done
pass "all modules build"

# ---------------------------------------------------------------------------
# 2. Test every module
# ---------------------------------------------------------------------------
echo
info "testing all modules (go test ./...)"
for mod in "${MODULES[@]}"; do
  (cd "$mod" && GOWORK=off CGO_ENABLED=1 go test -count=1 ./...) \
    || fail "tests failed in $mod"
done
pass "all modules pass tests"

# ---------------------------------------------------------------------------
# 3. Race tests on concurrency-sensitive modules
# ---------------------------------------------------------------------------
echo
info "race-testing: ${RACE_MODULES[*]}"
for mod in "${RACE_MODULES[@]}"; do
  (cd "$mod" && GOWORK=off CGO_ENABLED=1 go test -race -count=1 ./...) \
    || fail "race test failed in $mod"
done
pass "race tests clean"

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
echo
echo -e "${GREEN}========================================================${NC}"
echo -e "${GREEN}  All pre-tag checks passed${NC}"
echo -e "${GREEN}========================================================${NC}"
