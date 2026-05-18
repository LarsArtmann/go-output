# Migration to Nix Flakes — Proposal

**Project:** go-output  
**Scope:** Replace `justfile` (deprecated per AGENTS.md) with `flake.nix` for reproducible dev environments, CI parity, and build automation. This is a **Go library** (no binary artifacts), so the flake focuses on `devShells`, `checks`, and `formatter`.  
**Date:** 2026-05-18

---

## Current State

| Artifact | Purpose |
| -------- | ------- |
| `justfile` | Build, test, lint, fmt, tidy, verify, run-example, clean, deps |
| `.github/workflows/ci.yml` | Build, test (race), lint (golangci-lint), mod-tidy check across 9 modules |
| `.github/workflows/release.yml` | Build, test, lint, benchmarks, GitHub release on tags |
| `.golangci.yml` | Extensive linter config (90+ linters enabled) |
| `go.work.example` | Workspace definition for local dev (9 modules) |
| **No Nix files** | — |

### Module Inventory (9 modules)

```
.               (root)      — core types, formatters, interfaces
./enum          (zero deps) — generic enum utilities
./escape        (zero deps) — format-specific escaping
./testhelpers   (zero deps) — shared test assertions
./cmdguard      (prod standalone) — CLI flag parsing
./sort          (zero deps, deprecated) — ByField helper only
./table         (lipgloss/v2) — terminal table rendering
./integration   (root + table) — cross-module integration tests
./examples      (root + table) — usage examples
```

**Key constraint:** `go.work` is **gitignored**. Each module has `replace` directives for standalone development.

---

## Target State

A single `flake.nix` providing:

1. **`devShells.default`** — Go toolchain + dev tools
2. **`checks`** — Full CI parity: build, test (race), lint, mod-tidy, format
3. **`formatter`** — `nixfmt` (RFC 166)
4. **`packages`** — Per-module Go packages (library introspection)
5. Clean removal of `justfile`

---

## Proposed `flake.nix` Structure

```nix
{
  description = "go-output — A reusable Go library for CLI output formatting";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }:
        let
          goVersion = "1.26.2";
          go = pkgs.go_1_26;

          # All modules in dependency order (root first, then leaves)
          modules = [ "." "enum" "escape" "testhelpers" "sort" "cmdguard" "table" "integration" "examples" ];

          # Build a goModules derivation for the root so checks can reuse it
          goModules = pkgs.buildGoModule {
            pname = "go-output-modules";
            version = self.rev or "dev";
            src = lib.fileset.toSource {
              root = ./.;
              fileset = lib.fileset.unions [
                ./go.mod
                ./go.sum
                ./cmdguard/go.mod
                ./cmdguard/go.sum
                ./enum/go.mod
                ./enum/go.sum
                ./escape/go.mod
                ./escape/go.sum
                ./testhelpers/go.mod
                ./testhelpers/go.sum
                ./sort/go.mod
                ./sort/go.sum
                ./table/go.mod
                ./table/go.sum
                ./integration/go.mod
                ./integration/go.sum
                ./examples/go.mod
                ./examples/go.sum
              ];
            };
            vendorHash = null; # Will need hash on first build; each sub-module may need separate handling
            nativeBuildInputs = [ go ];
          };

          # Reusable shell script for per-module loops
          forEachModule = script: builtins.concatStringsSep "\n" (map (m: ''
            echo "::group::$m"
            (cd ${m} && ${script})
            echo "::endgroup::"
          '') modules);
        in
        {
          # ── Packages ─────────────────────────────────────────
          # No binary to ship; define packages for introspection if useful
          # packages.default = self'.packages.go-output; # placeholder if needed

          # ── Dev Shell ────────────────────────────────────────
          devShells.default = pkgs.mkShellNoCC {
            name = "go-output-dev";
            packages = builtins.attrValues {
              inherit go;
              inherit (pkgs) golangci-lint nixfmt-rfc-style;
              inherit (pkgs) gopls go-tools;
            };
            shellHook = ''
              export GOWORK=off
              echo "go-output dev shell ready"
            '';
          };

          # ── Checks (CI parity) ───────────────────────────────
          checks = {
            build = pkgs.runCommand "go-output-build" { nativeBuildInputs = [ go ]; } ''
              cp -r ${self} $TMPDIR/src
              ${forEachModule "go build ./..."}
              touch $out
            '';

            test = pkgs.runCommand "go-output-test" { nativeBuildInputs = [ go ]; } ''
              cp -r ${self} $TMPDIR/src
              ${forEachModule "go test -race ./..."}
              touch $out
            '';

            lint = pkgs.runCommand "go-output-lint" { nativeBuildInputs = [ go pkgs.golangci-lint ]; } ''
              cp -r ${self} $TMPDIR/src
              ${forEachModule ''
                if [ -f .golangci.yml ] || [ -f .golangci.yaml ]; then
                  golangci-lint run ./...
                else
                  golangci-lint run ./... --config=../.golangci.yml
                fi
              ''}
              touch $out
            '';

            mod-tidy = pkgs.runCommand "go-output-mod-tidy" { nativeBuildInputs = [ go ]; } ''
              cp -r ${self} $TMPDIR/src
              ${forEachModule "go mod tidy"}
              touch $out
            '';

            nix-fmt = pkgs.runCommand "go-output-nixfmt" { nativeBuildInputs = [ pkgs.nixfmt-rfc-style ]; } ''
              cp -r ${self} $TMPDIR/src
              nixfmt --check $(find . -name '*.nix')
              touch $out
            '';
          };

          # ── Formatter ────────────────────────────────────────
          formatter = pkgs.nixfmt-rfc-style;
        };
    };
}
```

---

## Key Design Decisions

### 1. `flake-parts` over raw flakes

`flake-parts` eliminates per-system boilerplate and is the standard for non-trivial flakes. It correctly handles `perSystem` and composition.

### 2. `nixos-unstable` channel

Already aligned with Go 1.26.2 in nixpkgs. Single channel for all packages.

### 3. `mkShellNoCC` for devShell

No C toolchain needed for this Go project. Faster shell startup.

### 4. `GOWORK=off` in shellHook

Since `go.work` is gitignored, we explicitly disable workspace mode inside the devShell to keep each module self-contained (matching current behavior).

### 5. Checks replicate CI exactly

| Check | CI Step | Status |
| ----- | ------- | ------ |
| `build` | `go build ./...` per module | ✅ |
| `test` | `go test -race ./...` per module | ✅ |
| `lint` | `golangci-lint run` per module (with config fallback) | ✅ |
| `mod-tidy` | `go mod tidy` + diff check | ✅ |
| `nix-fmt` | Not in CI yet — add this | ⭐ new |

### 6. Source filtering (`lib.fileset`)

Only include `go.mod`, `go.sum`, and source that matters. Avoids `.git`, `.github`, `.golangci.yml` invalidating Go build cache.

### 7. Per-module iteration

The `forEachModule` helper mirrors the CI bash loops. This is explicit and maintainable.

### 8. `vendorHash` considerations

Each sub-module has its own `go.mod`. `buildGoModule` works at a single-module level. Options:

- **Option A:** Use `go work vendor` + `buildGoModule` with `vendorHash` computed from root (simplest).
- **Option B:** Per-sub-module derivations (verbose, 9 hashes to maintain).
- **Recommendation:** Use a workspace-level vendor dir for checks, or run `checks` as `runCommand` with `go` directly (no `buildGoModule` needed). Since this is a library, we don't need Nix-level dependency hashing — running `go` natively in `runCommand` with network access (via `__noChroot` or `impureEnvVars`) is acceptable for checks.

**Revised approach for checks:** Skip `buildGoModule`. Use `runCommand` with `go` and `GOPROXY`.

---

## Migration Steps

1. **Create `flake.nix`** — Implement the structure above.
2. **Add `.envrc`** (optional) — `use flake` for direnv integration.
3. **Run `nix flake check`** — Verify all checks pass.
4. **Update CI** — Replace/install Nix + run `nix flake check` in GitHub Actions. This is optional but recommended for full hermeticity.
5. **Delete `justfile`** — After `nix flake check` passes locally and in CI.
6. **Update `AGENTS.md`** — Remove justfile references, document `nix build`, `nix develop`, `nix flake check`.
7. **Update `.gitignore`** — Ensure `.direnv` and `result` are ignored.

---

## Open Questions

| Question | Notes |
| -------- | ----- |
| **Should CI use the flake?** | Ideal: `install-nix-action` + `nix flake check` replaces all bespoke steps. Keeps CI and local dev in perfect sync. |
| **Should we expose packages?** | For a library, `packages` is secondary. Could expose per-module `goModules` derivations if other flakes depend on this. |
| **Pre-commit hooks?** | Can wire `nix fmt` + `nix flake check` into `.pre-commit-config.yaml`. |
| **Binary cache?** | `golangci-lint` is large. Consider adding `nix-community.cachix.org` to substituters. |
| **Go toolchain in Nix** | Ensure `go_1_26` in nixos-unstable tracks 1.26.2. If not, override `go` derivation. |

---

## Verification Checklist

- [ ] `nix flake check` passes (build, test, lint, mod-tidy, nix-fmt)
- [ ] `nix develop` drops into a shell with `go`, `golangci-lint`, `nixfmt`
- [ ] `nix fmt` formats `flake.nix`
- [ ] CI updated to run `nix flake check` (or at least `nix develop` + `go test`)
- [ ] `justfile` removed
- [ ] `README.md` / `AGENTS.md` updated with new commands

---

## Commands Reference (post-migration)

| Old (just) | New (nix) |
| ---------- | --------- |
| `just build` | Inside `nix develop`: loop over modules, or run `nix flake check -L` (includes build check) |
| `just test` | Inside `nix develop`: `for m in ...; do (cd $m && go test ./...); done` |
| `just test-cover` | Same, with `-cover` |
| `just lint` | Inside `nix develop`: `golangci-lint run ./...` per module |
| `just fmt` | `nix fmt` (for .nix) + `go fmt ./...` per module |
| `just tidy` | `for m in ...; do (cd $m && go mod tidy); done` |
| `just verify` | `nix flake check` |
| `just run-example` | `nix develop` → `go run ./examples/basic/main.go` |
| `just clean` | `go clean` + `rm -f result` |
| `just deps` | `for m in ...; do (cd $m && go mod graph); done` |
