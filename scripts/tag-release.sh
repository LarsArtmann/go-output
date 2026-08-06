#!/usr/bin/env bash
# tag-release.sh — create a complete release tag family (root + 16 submodules).
#
# This script eliminates the recurring "missing submodule tags" failure mode
# by atomically creating all annotated tags on the same commit. It does NOT
# push — the auto-git daemon or the release.yml workflow handles that.
#
# Usage: scripts/tag-release.sh <version>
#   version   e.g. v0.38.0 (must start with 'v')
#
# What it does (in order):
#   1. Fetch tags + verify remote state (avoid clobbering existing tags)
#   2. Verify working tree is clean
#   3. Verify the version does not already exist locally or on origin
#   4. Run pre-tag-check.sh (build + test + race all 19 modules)
#   5. Create 1 root + 16 submodule annotated tags on HEAD
#   6. Verify tag family parity (17 tags, all annotated, all on same commit)
#
# What it does NOT do:
#   - Push (safety: never push without explicit user action)
#   - Promote CHANGELOG (do that in a release-prepare commit first)
#   - Create GitHub Releases (release.yml handles that)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

fail() { echo -e "${RED}FAIL:${NC} $*" >&2; exit 1; }
info() { echo -e ":: $*"; }
pass()  { echo -e "${GREEN}ok${NC}  $*"; }

# ---------------------------------------------------------------------------
# Submodules that get release tags (excludes test-only: examples, integration).
# Keep in sync with .github/workflows/release.yml TAGGED_SUBMODULES.
# ---------------------------------------------------------------------------
TAGGED_SUBMODULES=(
  bdd d2 daghtml delimited escape graph markdown markup
  nom plantuml serialization table testhelpers testhelpers/graphtest tree tui
)

# ---------------------------------------------------------------------------
# 0. Parse arguments
# ---------------------------------------------------------------------------
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
  fail "usage: scripts/tag-release.sh <version>  (e.g. v0.38.0)"
fi
case "$VERSION" in
  v*) ;;
  *) fail "version must start with 'v' (got: $VERSION)" ;;
esac

VER_NUM="${VERSION#v}"

info "tagging release $VERSION (root + ${#TAGGED_SUBMODULES[@]} submodules)"
echo

# ---------------------------------------------------------------------------
# 1. Fetch tags + verify remote state
# ---------------------------------------------------------------------------
info "fetching tags from origin"
git fetch --tags
pass "tags synced"

info "checking remote for existing $VERSION tags"
REMOTE_TAGS=$(git ls-remote --tags origin | grep -F "/$VERSION" || true)
if [ -n "$REMOTE_TAGS" ]; then
  echo "$REMOTE_TAGS"
  fail "$VERSION tags already exist on origin — use a new version or retract"
fi
pass "no $VERSION tags on origin"

# ---------------------------------------------------------------------------
# 2. Verify clean working tree
# ---------------------------------------------------------------------------
info "checking working tree"
if ! git diff-index --quiet HEAD --; then
  fail "working tree is dirty — commit or stash before tagging"
fi
pass "working tree is clean"

# ---------------------------------------------------------------------------
# 3. Verify version does not exist locally
# ---------------------------------------------------------------------------
if git rev-parse -q --verify "refs/tags/$VERSION" >/dev/null; then
  fail "root tag $VERSION already exists locally"
fi
for mod in "${TAGGED_SUBMODULES[@]}"; do
  SUBTAG="${mod}/${VERSION}"
  if git rev-parse -q --verify "refs/tags/${SUBTAG}" >/dev/null; then
    fail "submodule tag $SUBTAG already exists locally"
  fi
done
pass "no $VERSION tags exist locally"

# ---------------------------------------------------------------------------
# 4. Run pre-tag-check.sh (build + test + race)
# ---------------------------------------------------------------------------
info "running pre-tag checks (build + test + race all modules)"
bash "$SCRIPT_DIR/pre-tag-check.sh"
echo

# ---------------------------------------------------------------------------
# 5. Create root + submodule annotated tags
# ---------------------------------------------------------------------------
COMMIT=$(git rev-parse HEAD)
info "creating annotated tags on $COMMIT"

git tag -a "$VERSION" -m "$VERSION" "$COMMIT"
pass "created root $VERSION"

for mod in "${TAGGED_SUBMODULES[@]}"; do
  SUBTAG="${mod}/${VERSION}"
  git tag -a "$SUBTAG" -m "$VERSION" "$COMMIT"
  pass "created $SUBTAG"
done

# ---------------------------------------------------------------------------
# 6. Verify tag family parity
# ---------------------------------------------------------------------------
echo
info "verifying tag family ($((1 + ${#TAGGED_SUBMODULES[@]})) tags, all annotated, all on same commit)"

TAG_COUNT=$(git tag -l "*$VERSION" | wc -l)
EXPECTED=$((1 + ${#TAGGED_SUBMODULES[@]}))
if [ "$TAG_COUNT" -ne "$EXPECTED" ]; then
  fail "expected $EXPECTED tags, found $TAG_COUNT"
fi

NON_ANNOTATED=0
WRONG_COMMIT=0
for t in $(git tag -l "*$VERSION"); do
  [ "$(git cat-file -t "$t")" = "tag" ] || NON_ANNOTATED=$((NON_ANNOTATED + 1))
  [ "$(git rev-list -n1 "$t")" = "$COMMIT" ] || WRONG_COMMIT=$((WRONG_COMMIT + 1))
done

[ "$NON_ANNOTATED" -eq 0 ] || fail "$NON_ANNOTATED tag(s) are not annotated"
[ "$WRONG_COMMIT" -eq 0 ] || fail "$WRONG_COMMIT tag(s) point to wrong commit"

pass "all $EXPECTED tags annotated, all on $COMMIT"

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
echo
echo -e "${GREEN}========================================================${NC}"
echo -e "${GREEN}  $VERSION tag family created ($EXPECTED tags)${NC}"
echo -e "${GREEN}========================================================${NC}"
echo
echo "Tags are LOCAL only. To publish:"
echo "  git push origin $VERSION"
echo "  git push origin refs/tags/*$VERSION"
echo
echo "Or let the auto-git daemon push them."
