# ADR 001: Multi-Module Workspace with Opt-In Heavy Dependencies

**Date:** 2026-05-07
**Status:** ACCEPTED
**Deciders:** Lars Artmann

## Context

go-output has a single `go.mod` with 3 third-party dependencies:
- `charm.land/lipgloss/v2` (heavy — many transitive deps, used only by `table/`)
- `github.com/go-faster/yaml` (medium, used only by `yaml.go`)
- `golang.org/x/term` (light, used only by `color.go`)

Most CLI apps using go-output for JSON/YAML/CSV output don't need lipgloss terminal tables.
But with a single module, every user pulls in lipgloss transitively.

Several sub-packages (enum, escape, cmdguard) have zero dependencies and are reusable independently.

## Decision

Split into 7 independent Go modules using `go.work` for local development:

| Module | Deps | Isolation benefit |
|---|---|---|
| Root (`go-output`) | enum, escape, yaml, x/term | Core formatters, no lipgloss |
| `enum/` | None | Reusable enum utilities |
| `escape/` | None | Reusable escaping (D2, DOT, Mermaid) |
| `cmdguard/` | None | Generic CLI flag parsing |
| `table/` | root, lipgloss | **Lipgloss isolated** — biggest win |
| `sort/` | root | Deprecated — points to stdlib |
| `integration/` | root, sort, table | Cross-module tests |
| `examples/` | root, table | Usage examples |

Root stays as `package output` — no core/ directory, no package rename.

## Key Design Choices

1. **Root IS the core module** — `package output` stays, no file moves for formatters
2. **No go.work committed** — gitignored per Go convention, replaced by `replace` directives in each go.mod
3. **Replace directives in every consuming module** — allows `cd table && go test ./...` standalone
4. **Leaf modules first** — enum, escape, cmdguard have zero deps and zero risk
5. **sort/ deprecated** — stdlib `slices.SortStableFunc` + `cmp.Compare` do the same job

## Consequences

**Positive:**
- Users who only need JSON/YAML/CSV get zero lipgloss deps
- enum, escape, cmdguard can be imported independently
- Each module can be versioned independently (future)
- Follows go-cqrs-lite workspace pattern

**Negative:**
- More go.mod files to maintain
- Replace directives needed in every consuming module for standalone dev
- d2/ and graph/ modules not yet extracted (future work)

## Not Done Yet (Future ADRs)

- Extract `d2/` as module (5 files moved, package renamed)
- Extract `graph/` as module (DOT + Mermaid + GraphRendererMixin)
- CI setup for multi-module repo
