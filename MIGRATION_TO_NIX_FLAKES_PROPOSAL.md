# Migration to Nix Flakes — Proposal

**Project:** go-output  
**Scope:** Replace `justfile` (deprecated per AGENTS.md) with `flake.nix` for reproducible dev environments, CI parity, and build automation. This is a **Go library** (no binary artifacts), so the flake focuses on `devShells`, `checks`, and `formatter`.  
**Date:** 2026-05-18  
**Status:** ✅ Implemented — `flake.nix` created and verified

---

## What Was Implemented

### `flake.nix` — Production-ready Nix flake

**Stack:**

| Tool | Purpose | Why |
|------|---------|-----|
| `flake-parts` | Module composition | Clean perSystem handling, native integration with treefmt-nix and git-hooks.nix |
| `treefmt-nix` | Formatting (nix fmt) | Nix formatting via nixfmt (RFC 166) + deadnix + statix |
| `git-hooks.nix` | Pre-commit hooks | Replaces broken buildflow hook + manual .pre-commit-config.yaml |
| ~~`gomod2nix`~~ | ~~Go deps → Nix~~ | **Skipped** — library has no binary artifacts; 9 modules with replace directives would need 9 configs |

### What Works

```bash
nix develop          # Dev shell with Go 1.26.2, golangci-lint, gopls
nix fmt              # Format .nix files (nixfmt, deadnix, statix)
nix flake check      # Formatting check + pre-commit validation
```

### What's NOT in the flake (intentional)

| Feature | Why not in flake | Where instead |
|---------|-----------------|---------------|
| `go build ./...` | Nix sandbox blocks network for `go mod download` | CI (`.github/workflows/ci.yml`) |
| `go test -race ./...` | Same sandbox issue | CI |
| `golangci-lint run` | Same sandbox issue + slow | CI |
| `go mod tidy` check | Same sandbox issue | CI |
| Go formatting (gofmt/gofumpt) | Handled by golangci-lint formatters | CI + manual |

**Rationale:** This is a Go library with 9 independent modules. Each module has its own `go.mod` with `replace` directives. Running Go checks in Nix would require either `gomod2nix` (9 configs) or `buildGoModule` (9 vendorHashes). The complexity is not worth it — CI already handles these checks reliably.

---

## Remaining Steps

| # | Task | Status |
|---|------|--------|
| 1 | Create `flake.nix` | ✅ Done |
| 2 | Create `.envrc` | ✅ Done |
| 3 | Add `.gitignore` entries (result, .direnv/) | ✅ Done |
| 4 | Run `nix flake check` | ✅ Passes |
| 5 | Update `CHANGELOG.md` | ⬜ |
| 6 | Update `README.md` with nix commands | ⬜ |
| 7 | Update `AGENTS.md` (nix commands, drop justfile refs) | ⬜ |
| 8 | Handle `.pre-commit-config.yaml` conflict | ⬜ |
| 9 | Delete `justfile` | ⬜ |
| 10 | Update CI to use nix (optional) | ⬜ |

---

## Module Inventory (9 modules, unchanged)

```
.               (root)      — core types, formatters, interfaces
./enum          (zero deps) — generic enum utilities
./escape        (zero deps) — format-specific escaping
./testhelpers   (zero docs) — shared test assertions
./cmdguard      (prod standalone) — CLI flag parsing
./sort          (zero deps, deprecated) — ByField helper only
./table         (lipgloss/v2) — terminal table rendering
./integration   (root + table) — cross-module integration tests
./examples      (root + table) — usage examples
```

## Commands Reference (post-migration)

| Task | Command |
|------|---------|
| Enter dev environment | `nix develop` |
| Format Nix files | `nix fmt` |
| Check formatting + hooks | `nix flake check` |
| Build all modules | `nix develop --command bash -c "for m in . enum escape testhelpers sort cmdguard table integration examples; do (cd \$m && go build ./...); done"` |
| Test all modules | `nix develop --command bash -c "for m in . enum escape testhelpers sort cmdguard table integration examples; do (cd \$m && go test -race ./...); done"` |
| Lint | `nix develop --command golangci-lint run ./...` |
| Format Go code | `nix develop --command bash -c "for m in . enum escape testhelpers sort cmdguard table integration examples; do (cd \$m && go fmt ./...); done"` |
