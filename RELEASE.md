# Release Process — go-output

This document describes how to cut a release for the go-output workspace (Pattern B: committed-replace model).

## Versioning Model

Only **root** (`github.com/larsartmann/go-output`) and **`testhelpers/`** get real version tags. All other sub-modules use `v0.0.0-00010101000000-000000000000` + `replace` directives — they are build-boundary optimizations, not independently consumable packages.

- Root tag: `vX.Y.Z`
- testhelpers tag: `testhelpers/vX.Y.Z`
- testhelpers/graphtest tag: `testhelpers/graphtest/vX.Y.Z`
- All other sub-modules: **no tags** (they resolve via local `replace`)

## Pre-Release Checklist

1. Ensure all tests pass: `nix run .#test`
2. Ensure lint is clean: `nix run .#lint`
3. Ensure race tests pass: `nix run .#test-race`
4. Ensure no vulnerabilities: `nix run .#govulncheck`
5. Ensure go.mod files are tidy: `nix run .#tidy` (should be no-op)
6. Verify zero TODO/FIXME in prod: `rg "TODO|FIXME|HACK" --type go -g '!*_test.go'`
7. Update `CHANGELOG.md` with the new version and changes
8. Update `FEATURES.md` if new features were added

## Cutting the Release

### 1. Tag root + testhelpers

```bash
git tag vX.Y.Z
git tag "testhelpers/vX.Y.Z"
git tag "testhelpers/graphtest/vX.Y.Z"
```

### 2. Push tags

```bash
git push origin --tags
```

### 3. Create GitHub Release

Draft release notes from `CHANGELOG.md`. Use the `docs/RELEASE_NOTES_vX.Y.Z.md` template if available.

### 4. Verify proxy

After ~5 minutes, verify root is available on the Go proxy:

```bash
GOFLAGS="" go list -m github.com/larsartmann/go-output@vX.Y.Z
GOFLAGS="" go list -m github.com/larsartmann/go-output/testhelpers@vX.Y.Z
```

## CI Automation

The `.github/workflows/release.yml` workflow automates tagging when a tag is pushed. It:

- Runs the full test suite
- Cross-compiles binaries for examples
- Publishes to the Go proxy
- Creates a GitHub release with auto-generated notes

## Emergency Rollback

To retract a broken release:

```bash
go mod edit -retract vX.Y.Z  # in root go.mod
git commit -am "retract vX.Y.Z"
git tag vX.Y.Z1  # new patch with retraction
git push origin --tags
```
