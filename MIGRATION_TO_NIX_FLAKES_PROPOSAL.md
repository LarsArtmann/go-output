# Migration to Nix Flakes — Proposal

**Project:** `go-output` (`github.com/larsartmann/go-output`)
**Date:** 2026-04-09
**Status:** Draft — Pending Approval

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State Analysis](#2-current-state-analysis)
3. [Why Nix Flakes](#3-why-nix-flakes)
4. [Proposed Architecture](#4-proposed-architecture)
5. [Detailed Migration Plan](#5-detailed-migration-plan)
6. [File-by-File Specification](#6-file-by-file-specification)
7. [CI Migration Strategy](#7-ci-migration-strategy)
8. [Risk Assessment](#8-risk-assessment)
9. [Migration Checklist](#9-migration-checklist)
10. [Open Questions](#10-open-questions)
11. [References](#11-references)

---

## 1. Executive Summary

This proposal outlines migrating `go-output` to [Nix Flakes](https://nixos.wiki/wiki/Flakes) for **reproducible builds, declarative development environments, and pinned tooling**. The migration is additive — no existing workflow is removed. Nix becomes the single source of truth for tool versions, while the Justfile remains the task runner.

**Key outcomes:**

- One-command reproducible dev shell: `nix develop`
- Pinned, locked versions for Go, golangci-lint, just, and all dev tools
- CI with Nix caching (optional Cachix integration)
- `nix flake check` replaces ad-hoc lint/test/build scripts
- `flake.lock` guarantees bit-for-bit reproducibility across machines

---

## 2. Current State Analysis

### 2.1 Build & Task System

| Tool | Purpose | Location |
|------|---------|----------|
| `just` | Task runner (build, test, lint, verify) | `justfile` |
| `go` | Compiler, test runner, formatter | System-installed |
| `golangci-lint` | Linting (60+ linters configured) | System-installed or `go install @latest` |
| `pre-commit` | Git hook management | `.pre-commit-config.yaml` |

### 2.2 Justfile Recipes

```
build        → go build ./...
test         → go test ./...
test-v       → go test -v ./...
test-cover   → go test -cover ./...
lint         → golangci-lint run --fix ./...
fmt          → go fmt ./...
tidy         → go mod tidy
verify       → build + test + lint
run-example  → go run ./examples/basic/main.go
clean        → rm -f coverage.out; go clean
deps         → go mod graph
```

### 2.3 CI Pipeline (GitHub Actions)

**`ci.yml`** — 3 jobs on every push/PR to main/master:

| Job | Steps |
|-----|-------|
| `test` | checkout → setup-go@v5 (1.23) → `go mod download` → `go build` → `go test -v -race` → `golangci-lint run` |
| `verify` | checkout → setup-go@v5 (1.23) → build + test + lint (duplicated) |
| `lint` | checkout → setup-go@v5 (1.23) → golangci-lint → `go mod tidy` diff check |

**`release.yml`** — On tag push `v*`:

| Job | Steps |
|-----|-------|
| `release` | checkout → setup-go@v5 (1.23) → verify tag → build → test → lint → benchmarks → GitHub Release |
| `goreleaser` | checkout → setup-go@v5 (1.23) → GoReleaser (no `.goreleaser.yml` found!) |

### 2.4 Pre-commit Hooks (`.pre-commit-config.yaml`)

| Hook | Source |
|------|--------|
| trailing-whitespace, end-of-file-fixer, check-yaml, check-added-large-files, check-merge-conflict, check-toml | `pre-commit-hooks` v4.5.0 |
| go-fmt, go-vet, go-mod-tidy | `pre-commit-golang` v1.0.0 |
| golangci-lint (local) | System-installed |
| go-test (local) | System-installed |

### 2.5 Go Module

```
module github.com/larsartmann/go-output
go 1.26.0

Direct deps:
  charm.land/lipgloss/v2 v2.0.2
  github.com/go-faster/yaml v0.4.6
```

### 2.6 Problems with Current State

1. **No tool version pinning** — `go install @latest` in CI yields different golangci-lint versions across runs
2. **CI Go version mismatch** — `ci.yml` uses `go-version: "1.23"` but `go.mod` requires `go 1.26.0` and `.golangci.yml` says `go: 1.26.1`. This is a **build-breaking bug**
3. **Duplicated CI jobs** — `test`, `verify`, and `lint` jobs overlap heavily
4. **No lockfile for dev tools** — Every developer gets different tool versions
5. **GoReleaser misconfigured** — `release.yml` references GoReleaser but no `.goreleaser.yml` exists
6. **Pre-commit depends on system tools** — `golangci-lint` and `go` must be pre-installed
7. **No caching strategy** — Dependencies re-downloaded on every CI run

---

## 3. Why Nix Flakes

### 3.1 Benefits

| Benefit | How |
|---------|-----|
| **Reproducibility** | `flake.lock` pins exact versions of Go, golangci-lint, just, and all transitive dependencies |
| **Declarative dev environment** | `nix develop` gives identical shells on macOS, Linux, CI |
| **No system pollution** | Tools live in `/nix/store`, not `/usr/local/bin` |
| **CI caching** | Nix store is cacheable; Cachix provides remote binary cache |
| **Single source of truth** | `flake.nix` defines all tooling; Justfile defines tasks |
| **Hermetic builds** | `nix build` produces bit-for-bit identical outputs |
| **Flake checks** | `nix flake check` runs build + test + lint in one command |
| **Composability** | Other flakes can consume `go-output` as an input |

### 3.2 Trade-offs

| Concern | Mitigation |
|---------|------------|
| Nix learning curve | Proposal includes full templates; copy-paste to start |
| CI cold-start time | Warm Nix store via Cachix or `actions/cache` |
| Team members without Nix | `just` commands continue to work identically |
| Windows support | Nix works on WSL2; native Windows not a concern for this Go library |
| CI complexity increase | Nix simplifies CI long-term; short-term migration cost |

### 3.3 Approach: `buildGoModule` vs `gomod2nix`

| Approach | Pros | Cons |
|----------|------|------|
| **`buildGoModule`** (recommended) | Built into nixpkgs; no extra input; uses `vendorHash`; simple | Must update `vendorHash` on dependency changes |
| **`gomod2nix`** | Generates `gomod2nix.toml` from `go.mod`; no hash guessing | Extra flake input; extra tool to maintain; another lockfile |

**Recommendation:** Use `buildGoModule` with `vendorHash`. It's the standard nixpkgs approach, requires fewer moving parts, and the hash update workflow is well-documented (`lib.fakeHash` → real hash). Since this is a library (no `main` package), `buildGoModule` is simpler because we primarily need `go build ./...` to compile, not produce a binary.

---

## 4. Proposed Architecture

### 4.1 New File Structure

```
go-output/
├── flake.nix                    # NEW: Main flake definition
├── flake.lock                   # NEW: Pinned dependency versions (auto-generated)
├── .gitignore                   # MODIFIED: Add .direnv/
├── justfile                     # UNCHANGED: Task runner (uses Nix-provided tools)
├── go.mod                       # UNCHANGED
├── go.sum                       # UNCHANGED
├── .golangci.yml                # UNCHANGED
├── .pre-commit-config.yaml      # REPLACED by pre-commit-hooks.nix in flake.nix
├── .github/workflows/ci.yml     # MODIFIED: Use Nix
├── .github/workflows/release.yml # MODIFIED: Use Nix (keep GoReleaser if desired)
├── .envrc                       # NEW (optional): Auto-enter dev shell with direnv
└── ...
```

### 4.2 Layered Design

```
┌─────────────────────────────────────────┐
│            Developer Workflow            │
│  nix develop → just {build,test,lint}   │
├─────────────────────────────────────────┤
│              Justfile                    │
│  Task definitions (unchanged)           │
│  Calls: go, golangci-lint, etc.        │
├─────────────────────────────────────────┤
│             flake.nix                    │
│  Provides: devShell, checks, packages   │
│  Pins: Go 1.26, golangci-lint, just     │
├─────────────────────────────────────────┤
│           nixpkgs (via flake.lock)       │
│  Binary cache: cache.nixos.org          │
└─────────────────────────────────────────┘
```

---

## 5. Detailed Migration Plan

### Phase 1: Foundation (Non-Breaking)

**Goal:** Add Nix support without changing any existing workflow.

| Step | Action | Risk |
|------|--------|------|
| 1.1 | Create `flake.nix` with `devShells.default` and `packages.default` | None — additive |
| 1.2 | Run `nix develop` to verify dev shell works | None |
| 1.3 | Verify `just build`, `just test`, `just lint` inside `nix develop` | None |
| 1.4 | Add `.direnv/` and `result` to `.gitignore` | None |
| 1.5 | Run `nix flake check` to verify all checks pass | None |
| 1.6 | Commit `flake.nix` + `flake.lock` | None |

### Phase 2: CI Migration

**Goal:** Replace GitHub Actions setup-go with Nix.

| Step | Action | Risk |
|------|--------|------|
| 2.1 | Add `cachix/install-nix-action` to CI workflow | Low — fallback to current |
| 2.2 | Replace `test` job with `nix flake check` | Low — same logic, different runner |
| 2.3 | Remove `verify` job (redundant with `nix flake check`) | Low — already redundant |
| 2.4 | Keep `lint` job as `go mod tidy` diff check (Nix doesn't cover this) | None |
| 2.5 | (Optional) Add Cachix for binary caching | None — opt-in |

### Phase 3: Pre-commit Integration

**Goal:** Replace `.pre-commit-config.yaml` with Nix-managed hooks.

| Step | Action | Risk |
|------|--------|------|
| 3.1 | Add `pre-commit-hooks.nix` flake input | Low |
| 3.2 | Configure hooks in `flake.nix` (gofmt, golangci-lint, go-test, go-mod-tidy) | Low |
| 3.3 | Remove `.pre-commit-config.yaml` | Low — `pre-commit-hooks.nix` is a superset |
| 3.4 | Update documentation | None |

### Phase 4: Polish & Documentation

| Step | Action | Risk |
|------|--------|------|
| 4.1 | Add `.envrc` for direnv auto-shell (optional) | None |
| 4.2 | Update `README.md` with Nix instructions | None |
| 4.3 | Update `AGENTS.md` with Nix commands | None |
| 4.4 | Fix CI Go version mismatch (1.23 → 1.26) | **This is a separate bug fix** |

---

## 6. File-by-File Specification

### 6.1 `flake.nix` (New)

```nix
{
  description = "go-output — A Go library for CLI output formatting across 12 formats";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
    pre-commit-hooks.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      pre-commit-hooks,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        goVersion = "go_1_26";
      in
      {
        # ─── Packages ────────────────────────────────────────────
        # This is a library, so the "package" is really just verifying
        # that `go build ./...` succeeds. No binary is produced.
        packages.default = pkgs.buildGoModule {
          pname = "go-output";
          version = "0.1.0";

          src = pkgs.lib.cleanSource self;

          # Update this hash when go.sum changes:
          #   1. Set vendorHash = pkgs.lib.fakeHash;
          #   2. Run: nix build
          #   3. Replace with the expected hash from the error message
          vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

          # Build all packages (library has no main)
          subPackages = [ "." ];

          # Only compile-check; no main binary to install
          installPhase = ''
            mkdir -p $out
          '';
        };

        # ─── Checks (nix flake check) ───────────────────────────
        checks = {
          # Build check
          build = self.packages.${system}.default;

          # Go vet
          go-vet = pkgs.runCommand "go-vet-check" {
            nativeBuildInputs = [ pkgs.${goVersion} ];
            src = pkgs.lib.cleanSource self;
          } ''
            cd $src
            go vet ./...
            touch $out
          '';

          # Go test with race detector
          go-test = pkgs.runCommand "go-test-check" {
            nativeBuildInputs = [ pkgs.${goVersion} ];
            src = pkgs.lib.cleanSource self;
          } ''
            export HOME=$TMPDIR
            cd $src
            go test -race -count=1 ./...
            touch $out
          '';

          # golangci-lint
          go-lint = pkgs.runCommand "go-lint-check" {
            nativeBuildInputs = [
              pkgs.${goVersion}
              pkgs.golangci-lint
            ];
            src = pkgs.lib.cleanSource self;
          } ''
            export GOLANGCI_LINT_CACHE=$TMPDIR/golangci-cache
            export HOME=$TMPDIR
            cd $src
            golangci-lint run ./...
            touch $out
          '';

          # Go format check
          go-fmt = pkgs.runCommand "go-fmt-check" {
            nativeBuildInputs = [ pkgs.${goVersion} ];
            src = pkgs.lib.cleanSource self;
          } ''
            cd $src
            test -z "$(gofmt -l .)"
            touch $out
          '';

          # Pre-commit hooks (managed by pre-commit-hooks.nix)
          pre-commit-check = pre-commit-hooks.lib.${system}.run {
            src = pkgs.lib.cleanSource self;
            hooks = {
              gofmt.enable = true;
              golangci-lint.enable = true;
              go-vet.enable = true;
              trailing-whitespace-fixer.enable = true;
              end-of-file-fixer.enable = true;
              check-yaml.enable = true;
              check-added-large-files.enable = true;
              check-merge-conflict.enable = true;
            };
          };
        };

        # ─── Dev Shell (nix develop) ────────────────────────────
        devShells.default = pkgs.mkShell {
          name = "go-output-dev";

          packages = with pkgs; [
            # Go toolchain
            ${goVersion}
            gotools
            gopls
            delve

            # Linting & formatting
            golangci-lint
            gofumpt
            golines

            # Task runner
            just

            # GoReleaser (for releases)
            goreleaser

            # Pre-commit
            pre-commit

            # YAML tooling
            yamllint
          ];

          shellHook = ''
            echo "go-output development shell"
            echo "Go: $(go version)"
            echo "golangci-lint: $(golangci-lint version)"
            echo "just: $(just --version)"
            echo ""
            echo "Commands: just --list"
          '';
        };

        # ─── Formatter (nix fmt) ────────────────────────────────
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}
```

### 6.2 `.gitignore` (Modified)

Add these lines:

```gitignore
# Nix
.direnv/
result
result-*
```

### 6.3 `.envrc` (New — Optional)

```bash
# Auto-enter Nix dev shell with direnv
# Install: brew install direnv && echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
use flake
```

### 6.4 `.github/workflows/ci.yml` (Modified)

```yaml
name: CI

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]

jobs:
  nix-check:
    name: Nix Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Nix
        uses: cachix/install-nix-action@v25
        with:
          nix_path: nixpkgs=channel:nixos-unstable

      - name: Run all checks
        run: nix flake check --print-build-logs

  tidy-check:
    name: Go Mod Tidy
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Nix
        uses: cachix/install-nix-action@v25
        with:
          nix_path: nixpkgs=channel:nixos-unstable

      - name: Check go.mod is tidy
        run: |
          cp go.mod go.mod.bak
          cp go.sum go.sum.bak
          nix develop --command go mod tidy
          diff go.mod go.mod.bak || (echo "go.mod is not tidy" && exit 1)
          diff go.sum go.sum.bak || (echo "go.sum is not tidy" && exit 1)
```

### 6.5 `.github/workflows/release.yml` (Modified)

```yaml
name: Release

on:
  push:
    tags:
      - "v*"

permissions:
  contents: write

jobs:
  verify:
    name: Verify
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install Nix
        uses: cachix/install-nix-action@v25
        with:
          nix_path: nixpkgs=channel:nixos-unstable

      - name: Run all checks
        run: nix flake check --print-build-logs

      - name: Run benchmarks
        run: nix develop --command go test -run=^$ -bench=. -benchmem ./...

  release:
    name: Release
    runs-on: ubuntu-latest
    needs: verify
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Install Nix
        uses: cachix/install-nix-action@v25
        with:
          nix_path: nixpkgs=channel:nixos-unstable

      - name: Verify tag version
        run: |
          TAG="${GITHUB_REF#refs/tags/}"
          VERSION="${TAG#v}"
          echo "Releasing version: $VERSION"

      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          draft: false
          prerelease: ${{ contains(github.ref, 'alpha') || contains(github.ref, 'beta') || contains(github.ref, 'rc') }}
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 7. CI Migration Strategy

### 7.1 Before (Current)

```
3 jobs × (setup-go@v5 + go install @latest + manual steps)
  = 3 Go installations, 3 golangci-lint installations, no caching
  = ~6-8 minutes per CI run
```

### 7.2 After (Nix)

```
2 jobs × (install-nix-action + nix flake check)
  = 1 Nix installation, cached builds, deduplicated checks
  = ~3-5 minutes per CI run (first run), ~1-2 minutes (cached)
```

### 7.3 Caching Strategy

**Tier 1 (Free, zero-config):**
- `cache.nixos.org` provides pre-built nixpkgs binaries
- `/nix/store` is cached per CI runner

**Tier 2 (Recommended):**
- Add [Cachix](https://cachix.org/) for project-specific binary cache
- Add to CI:

```yaml
- uses: cachix/cachix-action@v14
  with:
    name: go-output
    signingKey: '${{ secrets.CACHIX_SIGNING_KEY }}'
```

**Tier 3 (Optional):**
- GitHub Actions `actions/cache` on `/nix/store` paths

---

## 8. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `vendorHash` mismatch on `go.sum` changes | High | Low | Documented workflow: `lib.fakeHash` → rebuild → copy real hash |
| Nix not available on contributor machines | Medium | Low | `just` commands still work with system-installed Go |
| `pre-commit-hooks.nix` conflicts with `.pre-commit-config.yaml` | Medium | Low | Remove `.pre-commit-config.yaml` in Phase 3 |
| CI slower on first run | High | Low | Cachix caching; one-time cost |
| nixpkgs `go_1_26` not yet available | Low | High | Fall back to `go_1_24` or override with `buildGoModule`'s `go` parameter |
| `charm.land/lipgloss/v2` not in nixpkgs | Low | Low | It's a Go module, fetched via `go.sum` — Nix doesn't need it in nixpkgs |

### Critical Finding: CI Go Version Mismatch

The current `ci.yml` uses `go-version: "1.23"` but:
- `go.mod` declares `go 1.26.0`
- `.golangci.yml` declares `go: 1.26.1`

**This is a build-breaking bug.** The Nix migration fixes this by pinning the Go version in `flake.nix`, ensuring consistency across local dev and CI.

---

## 9. Migration Checklist

### Phase 1 — Foundation

- [ ] Create `flake.nix` with `devShells.default`, `packages.default`, `checks`
- [ ] Set `vendorHash = lib.fakeHash;` and run `nix build` to get the real hash
- [ ] Verify `nix develop` provides working Go, golangci-lint, just
- [ ] Verify `just build && just test && just lint` inside dev shell
- [ ] Run `nix flake check` — all checks pass
- [ ] Add `.direnv/` and `result` to `.gitignore`
- [ ] Commit `flake.nix`, `flake.lock`, updated `.gitignore`

### Phase 2 — CI Migration

- [ ] Replace `ci.yml` with Nix-based workflow
- [ ] Replace `release.yml` with Nix-based workflow
- [ ] Verify CI passes on a feature branch
- [ ] (Optional) Add Cachix binary cache

### Phase 3 — Pre-commit Integration

- [ ] Add `pre-commit-hooks.nix` to flake inputs
- [ ] Configure hooks in `flake.nix`
- [ ] Remove `.pre-commit-config.yaml`
- [ ] Update `AGENTS.md` instructions

### Phase 4 — Polish

- [ ] Add `.envrc` for direnv (optional)
- [ ] Update `README.md` with Nix quickstart
- [ ] Update `AGENTS.md` with Nix commands
- [ ] Fix Go version mismatch in any remaining non-Nix CI configs

---

## 10. Open Questions

1. **GoReleaser config missing** — The `release.yml` references GoReleaser but no `.goreleaser.yml` exists. Should we create one, or remove the GoReleaser step?

2. **CI job deduplication** — The current `ci.yml` has 3 overlapping jobs (`test`, `verify`, `lint`). Should we consolidate to a single `nix flake check` job + a `tidy-check` job, or keep separate jobs for better failure isolation?

3. **`go_1_26` availability in nixpkgs** — If `go_1_26` is not yet in `nixos-unstable`, we may need to use `go_1_24` temporarily or provide a Go overlay. Verify before starting Phase 1.

4. **Cachix** — Should we set up a Cachix cache for this project? This requires creating a Cachix account and adding signing keys as GitHub secrets.

5. **Library-only `packages.default`** — Since `go-output` is a library (no `main` package), `buildGoModule` will only verify compilation. Is this sufficient, or should we also build the `examples/basic` binary as an additional package?

6. **golangci-lint version** — The `.golangci.yml` config is extensive (60+ linters). Should the Nix flake pin a specific golangci-lint version, or use whatever nixpkgs provides? (Recommendation: pin via nixpkgs flake.lock.)

---

## 11. References

- [Nix Flakes — NixOS Wiki](https://nixos.wiki/wiki/Flakes)
- [nix.dev — Flakes concepts](https://nix.dev/concepts/flakes)
- [buildGoModule — nixpkgs manual](https://nixos.org/manual/nixpkgs/stable/#sec-language-go)
- [gomod2nix — GitHub](https://github.com/nix-community/gomod2nix)
- [pre-commit-hooks.nix — GitHub](https://github.com/cachix/pre-commit-hooks.nix)
- [Cachix — Binary cache hosting](https://cachix.org/)
- [install-nix-action — GitHub](https://github.com/cachix/install-nix-action)
- [flake-utils — GitHub](https://github.com/numtide/flake-utils)
- [Xe Iaso — Building Go programs with Nix Flakes](https://xeiaso.net/blog/nix-flakes-go-programs/)
- [Haseeb Majid — Go dev shell with Nix Flakes](https://haseebmajid.dev/posts/2023-10-26-how-to-setup-a-go-development-shell-with-nix-flakes/)

---

_This proposal is designed to be incremental and non-breaking. Each phase can be completed independently and rolled back if issues arise._
