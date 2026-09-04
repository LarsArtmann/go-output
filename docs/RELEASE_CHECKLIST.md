# Release Checklist

The authoritative 8-step sequence for cutting a release. Every step is
mandatory. Skipping any step has caused broken releases in the past
(v0.36.0 and v0.37.0 both shipped without submodule tags).

## Tag Convention

- **17 tags per release**: 1 root (`vX.Y.Z`) + 16 submodule tags (`<module>/vX.Y.Z`).
- **All annotated** (never lightweight — `git cat-file -t <tag>` must return `tag`).
- **All on the same commit**.
- **Excluded from tagging**: `examples/` and `integration/` (test-only modules).
- **Tag message**: `vX.Y.Z` (or `vX.Y.Z: <summary>` for manual tags).

### The 16 Tagged Submodules

```
bdd  d2  daghtml  delimited  escape  graph  markdown  markup
nom  plantuml  serialization  table  testhelpers  testhelpers/graphtest  tree  tui
```

## The 8-Step Sequence

### 1. Promote CHANGELOG

Move entries from `[Unreleased]` to `[X.Y.Z] - YYYY-MM-DD` in `CHANGELOG.md`.

```bash
# Verify the version isn't already in the CHANGELOG
grep -c "\[X.Y.Z\]" CHANGELOG.md  # should be 0 before, 1 after
```

### 2. Create release-prepare commit

```bash
git add CHANGELOG.md
git commit -m "chore(release): prepare vX.Y.Z changelog"
```

### 2b. Re-bump sibling pins to the new version

Every module's `go.mod` must pin the NEW released version for all
`github.com/larsartmann/go-output/*` requires (Pattern B, ADR 009
amendment). The pins point at the PREVIOUS release right now; after
tagging they must move to the new one so `go mod tidy` and CI keep
resolving published artifacts and no sentinel ever reappears.

```bash
# After the tags exist locally (step 5), bump every sibling pin:
NEW="vX.Y.Z"
for f in go.mod */go.mod */*/go.mod; do
  sed -i -E "s|(github.com/larsartmann/go-output/[a-z/-]+) v[0-9]+\.[0-9]+\.[0-9]+|\1 $NEW|g" "$f"
  sed -i -E "s|(github.com/larsartmann/go-output) v[0-9]+\.[0-9]+\.[0-9]+|\1 $NEW|g" "$f"
done
nix run .#tidy      # must preserve the pins — if any go.mod reverts, investigate
nix run .#build && nix run .#test
git add -A && git commit -m "chore(deps): bump internal modules to $NEW"
```

A `v0.0.0-00010101…` sentinel anywhere is a release blocker.

### 3. Verify CI is green

Check the latest CI run on the release-prepare commit's branch.
Do NOT proceed if CI is red.

### 4. Run pre-tag checks

```bash
nix run .#setup-workspace   # ensure go.work exists
scripts/pre-tag-check.sh vX.Y.Z
```

This builds, tests, and race-tests all 19 modules. It also verifies the
working tree is clean and the tag does not already exist.

### 5. Create the tag family

**Option A — automated (recommended):**

```bash
scripts/tag-release.sh vX.Y.Z
```

This runs the full sequence: fetch, verify clean tree, pre-tag-check, tag
root + 16 submodules, verify parity.

**Option B — manual (if you must):**

```bash
git fetch --tags
git ls-remote --tags origin | grep vX.Y.Z  # must be empty

COMMIT=$(git rev-parse HEAD)
MSG="vX.Y.Z: <summary from CHANGELOG>"

# Root tag
git tag -a vX.Y.Z -m "$MSG" "$COMMIT"

# 16 submodule tags
for mod in bdd d2 daghtml delimited escape graph markdown markup \
           nom plantuml serialization table testhelpers testhelpers/graphtest tree tui; do
  git tag -a "${mod}/vX.Y.Z" -m "$MSG" "$COMMIT"
done
```

### 6. Push tags

```bash
git push origin vX.Y.Z
git push origin "refs/tags/*vX.Y.Z"
```

Or let the auto-git daemon push them.

The `release.yml` workflow will trigger on the root `v*` tag, run build/test/lint,
auto-create any missing submodule tags as a safety net, and create the GitHub Release.

### 7. Verify GitHub Release

- Go to the repository Releases page.
- Confirm a release for `vX.Y.Z` exists with auto-generated notes.
- If missing, the `release.yml` step failed — check the Actions tab.

### 8. Verify module resolution

```bash
# From a clean directory (not this repo)
GOPROXY=off GOWORK=off go mod download github.com/larsartmann/go-output@vX.Y.Z
GOPROXY=off GOWORK=off go mod download github.com/larsartmann/go-output/testhelpers@vX.Y.Z
```

Both should succeed. The root module is the only one external consumers
`go get` directly; `testhelpers` is the only independently versioned sub-module.

## Post-Release Verification

```bash
# Verify all 17 tags exist on origin, annotated, on same commit
git fetch --tags
git tag -l '*vX.Y.Z' | wc -l   # must be 17
for t in $(git tag -l '*vX.Y.Z'); do
  echo "$(git cat-file -t $t) $(git rev-list -n1 $t | cut -c1-7) $t"
done
```

All must show `tag` (annotated) and the same commit SHA.

## What To Do If Something Goes Wrong

### Missing submodule tags (the recurring failure)

If the root tag was pushed but submodule tags are missing:

```bash
git fetch --tags
COMMIT=$(git rev-list -n1 vX.Y.Z)
MSG="vX.Y.Z: <summary>"

for mod in bdd d2 daghtml delimited escape graph markdown markup \
           nom plantuml serialization table testhelpers testhelpers/graphtest tree tui; do
  SUBTAG="${mod}/vX.Y.Z"
  if ! git rev-parse -q --verify "refs/tags/${SUBTAG}" >/dev/null; then
    git tag -a "$SUBTAG" -m "$MSG" "$COMMIT"
    echo "created $SUBTAG"
  fi
done
git push origin "refs/tags/*vX.Y.Z"
```

### Root tag is lightweight instead of annotated

```bash
COMMIT=$(git rev-list -n1 vX.Y.Z)
git tag -d vX.Y.Z
git tag -a vX.Y.Z -m "vX.Y.Z: <summary>" "$COMMIT"
git push --force-with-lease origin vX.Y.Z
```

This changes only the tag metadata (commit is unchanged). The Go module
proxy resolves tags to commits, so module content is unaffected.

### Retracting a bad release

Add to root `go.mod`:

```
retract vX.Y.Z
```

Then cut a new version (vX.Y.(Z+1)) that supersedes it. Document the
retraction in the CHANGELOG.
