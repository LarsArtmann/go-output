# go-output

A reusable Go library for CLI output formatting (16 formats across table/tree/graph shapes) plus NOM-style real-time progress visualization. Multi-module Go workspace; root is intentionally dependency-light.

## The Core Invariant

**Root (`package output`) has ZERO imports of any sub-module.** This is the load-bearing architectural rule. Everything else follows from it.

- `go get github.com/larsartmann/go-output` pulls NO lipgloss, NO bubbletea, NO yaml/toml, NO d2/graph/table/nom/tui deps.
- Sub-modules self-register into root's registries via their own `init()` (see Patterns). Root never imports them back.
- **Never** add an `import ".../table"`, `".../d2"`, `".../nom"`, etc. to a file in the root package. It breaks the entire dependency-isolation guarantee.

## Module Map

18 modules (root + 17 sub-modules). Each has its own `go.mod`; **deps live in `go.mod`, not here** — read the file for ground truth.

| Module                   | Purpose                                                                       |
| ------------------------ | ----------------------------------------------------------------------------- |
| **root** (`output`)      | Core types, Format/Shape enums, registries, Markdown, Tree, Graph state       |
| `delimited/`             | CSV + TSV writers                                                             |
| `serialization/`         | JSON + YAML + TOML + JSONL                                                    |
| `markup/`                | XML + HTML + AsciiDoc + Streaming HTML                                        |
| `table/`                 | Lipgloss terminal tables                                                      |
| `d2/`                    | D2 diagrams (rich domain model: shapes, arrows, SQL tables)                   |
| `graph/`                 | DOT + Mermaid renderers (share root's `GraphRendererState`)                   |
| `plantuml/`              | PlantUML diagrams                                                             |
| `nom/`                   | NOM-style real-time progress (dependency tree, timing cache, inline renderer) |
| `tui/`                   | Bubble Tea interactive TUI (depends on nom)                                   |
| `enum/`                  | Generic enum utilities (zero deps)                                            |
| `escape/`                | Format-specific escaping (zero deps)                                          |
| `envdetect/`             | Shared CI / NO_COLOR env detection (zero deps; used by root + nom)            |
| `testhelpers/`           | Shared test assertions (zero deps)                                            |
| `testhelpers/graphtest/` | Shared graph test fixtures                                                    |
| `bdd/`                   | BDD test suite (Ginkgo/Gomega, test-only)                                     |
| `integration/`           | Cross-module integration tests                                                |
| `examples/`              | Usage examples                                                                |

## Commands

**Nix is the source of truth.** `go build ./...` from root only builds root — always use the flake apps, which iterate every module:

```bash
nix develop                # Dev shell (Go 1.26, golangci-lint, gopls)
nix run .#build            # Build all 18 modules
nix run .#test             # Test all 18 modules
nix run .#test-race        # Race-test nom + tui (concurrency-sensitive)
nix run .#lint             # golangci-lint across all modules
nix run .#tidy             # go mod tidy all modules
nix run .#setup-workspace  # Generate gitignored go.work from go.work.example
nix fmt                    # Format .nix files
nix flake check            # Formatting + pre-commit hooks
```

Go checks are NOT in `nix flake check` (sandbox blocks `go mod download`); CI handles them. `.pre-commit-config.yaml` exists for non-Nix users.

## Patterns

These are non-obvious from reading code alone — the "how does this even work" knowledge.

- **Registry dispatch via `init()`**: Root's `RenderTableData()` / `RenderAnyData()` dispatch to registered marshalers, yet root imports no sub-module. Each sub-module calls `RegisterFormatShapes(...)` / registers its marshaler in its own `init()`. The generic `formatRegistry[T]` backs shape capabilities, tableData, and anyData registries. Importing a sub-module is what activates it.
- **Shape capability matrix**: Each format declares which data shapes (table/tree/graph) it supports via `RegisterFormatShapes()` in `init()`. Query with `f.Supports(shape)` or `FormatsForShape(shape)`. See `docs/FORMAT_ARCHITECTURE.md`.
- **ColorMode wiring — three mechanisms, use the right one**: `table.New(table.WithColorMode(...))` (functional option), `ASCIITreeRenderer.SetColorMode(...)` / `MarkdownTable.SetColorMode(...)` (setter), `RenderTableData(data, fmt, RenderOptions{ColorMode: ...})` (dispatch field). `ColorModeAuto` detects terminal via `x/term`, respects `NO_COLOR`/`CI`/`FORCE_COLOR`.
- **Composition via `GraphRendererState`**: DOT/Mermaid/PlantUML renderers embed root's `GraphRendererState` for shared node/edge state via `AddNode`/`AddEdge`. `AddTreeNodes` accepts the `NodeEdgeAppender` interface (not raw slice pointers). `TableDataStore` exposes `Data()` so sub-modules share the table data field.
- **Branded IDs**: Phantom types (`D2NodeID`, `TreeNodeID`, ...) via `go-branded-id` prevent mixing ID types at compile time. `d2/` re-exports `D2NodeID`/`D2NodeLabel` from root so users need only one import.
- **All renderers implement `Render() (string, error)`**: Use `MustRender(r)` in tests/examples for the panic-on-error shortcut.

## Gotchas

Things that will silently break or that an agent would get wrong from code alone.

- **Never import a sub-module into root** — see Core Invariant above.
- **`testhelpers/` is zero-dep by design** — it cannot import `output`. Cross-module test helpers must stay local to each module or use table-driven patterns.
- **`internal/` is root-only** — Go forbids sub-modules from importing `internal/` packages. `internal/gentest` and `internal/testutils` are root-only; sub-modules inline their own test helpers.
- **Depguard restricts imports** — `.golangci.yml` has explicit allow-lists. When a module gains a new sibling dep, add it to BOTH the `default` and `main` allow-lists or lint fails. Each module has its own `.golangci.yml` section.
- **Every module's `go.mod` needs `replace` directives** for sibling deps, plus add the module to `flake.nix`'s `modules` list and `go.work.example`'s `use (...)` block.
- **Mono-version tagging** — all 18 modules release in lockstep under the same `vX.Y.Z` (root tag + `submod/vX.Y.Z` tags). Never version a module independently.
- **NOM events use `nom.Event*` constants** (e.g., `nom.EventWorkflowStarted`), not bare string literals.
- **`envdetect` centralizes CI/NO_COLOR detection** — root `color.go` and `nom/inline_renderer.go` both delegate to `envdetect.IsCI()` / `IsNoColor()`. Don't re-inline this logic.
- **Code duplication threshold is `art-dupl -t 24`** (project standard). Below t=20, reported clones are almost entirely Go test idioms or module-boundary re-declarations — both acceptable. See ADR 005.

## Pointers

Deeper knowledge lives in dedicated docs — read these before non-trivial work:

- `docs/DOMAIN_LANGUAGE.md` — ubiquitous language (Format, Shape, Renderer, TableData, TreeNode, GraphNode, ...)
- `docs/FORMAT_ARCHITECTURE.md` — the 16 formats × 3 shapes matrix and registry internals
- `docs/adr/` — 7 ADRs: multi-module workspace (001), shape matrix (002), d2/graph extraction (003), footer row (004), duplication thresholds (005), API stability (006), nom composition (007)
- `FEATURES.md`, `TODO_LIST.md`, `CHANGELOG.md`, `README.md`
