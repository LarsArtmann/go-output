# ADR 001: Multi-Module Workspace with Opt-In Heavy Dependencies

**Date:** 2026-05-07
**Status:** ACCEPTED & IMPLEMENTED
**Deciders:** Lars Artmann

## Context

go-output has a single `go.mod` with 3 third-party dependencies:

- `charm.land/lipgloss/v2` (heavy — many transitive deps, used only by `table/`)
- `github.com/go-faster/yaml` (medium, used only by `yaml.go`)
- `golang.org/x/term` (light, used only by `color.go`)

Most CLI apps using go-output for JSON/YAML/CSV output don't need lipgloss terminal tables.
But with a single module, every user pulls in lipgloss transitively.

Several sub-packages (enum, escape, testhelpers) have zero dependencies and are reusable independently.

## Decision

Split into 10 independent Go modules using `go.work` for local development:

| Module             | Deps                       | Isolation benefit                    |
| ------------------ | -------------------------- | ------------------------------------ |
| Root (`go-output`) | enum, escape, yaml, x/term, branded-id, testhelpers | Core formatters, no lipgloss/d2/graph |
| `enum/`            | testhelpers (tests only)   | Reusable enum utilities              |
| `escape/`          | None                       | Reusable escaping (D2, DOT, Mermaid) |
| `testhelpers/`     | None                       | Shared test assertions               |
| `d2/`              | root, escape, testhelpers  | D2 diagram renderer (rich domain)    |
| `graph/`           | root, escape, testhelpers  | DOT + Mermaid renderers              |
| `table/`           | root, lipgloss             | **Lipgloss isolated** — biggest win  |
| `sort/`            | None                       | Deprecated — only ByField helper     |
| `integration/`     | root, table, d2, graph     | Cross-module tests                   |
| `examples/`        | root, table, d2, graph     | Usage examples                       |

Root stays as `package output` — no core/ directory, no package rename.

## Key Design Choices

1. **Root IS the core module** — `package output` stays, no file moves for formatters
2. **No go.work committed** — gitignored per Go convention, replaced by `replace` directives in each go.mod
3. **Replace directives in every consuming module** — allows `cd table && go test ./...` standalone
4. **Leaf modules first** — enum, escape, testhelpers have zero deps and zero risk
5. **sort/ deprecated** — stdlib `slices.SortStableFunc` + `cmp.Compare` do the same job
6. **d2/ and graph/ extracted** — rich domain models moved to own modules (see ADR 003)

## Consequences

**Positive:**

- Users who only need JSON/YAML/CSV get zero lipgloss, zero d2, zero graph deps
- enum, escape, testhelpers can be imported independently
- Each module can be versioned independently (future)
- Follows go-cqrs-lite workspace pattern
- CI tests all 10 modules independently

**Negative:**

- More go.mod files to maintain (10 total)
- Replace directives needed in every consuming module for standalone dev
- `render_tabledata.go` cannot call d2/graph constructors (returns `UnsupportedFormatError`)
