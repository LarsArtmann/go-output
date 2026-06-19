# Release Process — go-output

This document describes how to cut a release for the 20-module mono-version workspace.

## Mono-Version Tagging

All 20 modules release in lockstep under the same `vX.Y.Z`. The root module gets the bare tag; each sub-module gets a `submod/vX.Y.Z` tag.

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

### 1. Tag the root module

```bash
git tag vX.Y.Z
```

### 2. Tag each sub-module

Every sub-module needs its own `submod/vX.Y.Z` tag:

```bash
for module in delimited serialization markup table markdown tree d2 graph plantuml nom tui enum escape envdetect testhelpers bdd integration examples; do
  git tag "${module}/vX.Y.Z"
done
```

For nested modules:

```bash
git tag "testhelpers/graphtest/vX.Y.Z"
```

### 3. Push all tags

```bash
git push origin --tags
```

### 4. Create GitHub Release

Draft release notes from `CHANGELOG.md`. Use the `docs/RELEASE_NOTES_vX.Y.Z.md` template if available.

### 5. Verify proxy

After ~5 minutes, verify the modules are available on the Go proxy:

```bash
GOFLAGS="" go list -m github.com/larsartmann/go-output@vX.Y.Z
GOFLAGS="" go list -m github.com/larsartmann/go-output/nom@vX.Y.Z
```

## CI Automation

The `.github/workflows/release.yml` workflow automates tagging when a tag is pushed. It:

- Runs the full test suite
- Cross-compiles binaries for examples
- Publishes to the Go proxy
- Creates a GitHub release with auto-generated notes

## Untagged Modules

Some modules (e.g., `envdetect`) were historically not tagged. If a module has no tag, consumers must use a `replace` directive in their `go.mod` to point to the local path. This is fragile — ensure every module gets tagged at each release.

## Emergency Rollback

To retract a broken release:

```bash
go mod edit -retract vX.Y.Z  # in root go.mod
git commit -am "retract vX.Y.Z"
git tag vX.Y.Z1  # new patch with retraction
git push origin --tags
```
