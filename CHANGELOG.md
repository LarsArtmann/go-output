# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Removed

- **BREAKING**: `format_deprecated.go` removed — `OutputFormat`, `FormatCategory`, `ParseOutputFormat`, `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` all removed. Use `Format`, `Shape`, `ParseFormat()`, `Supports(Shape*)`, `Shapes()` instead.
- **BREAKING**: `sort/` module removed — `ByField` is trivially replaceable with `slices.SortStableFunc` + `cmp.Compare` from stdlib.
- **BREAKING**: `MermaidFlowchartRenderer` and `MermaidTreeRenderer` removed — use `MermaidFromTableData` and `MermaidFromTree` instead.
- `gci` formatter removed from `.golangci.yml` (conflicted with `goimports` on local-prefix grouping in sub-modules).

### Changed

- **BREAKING**: D2 diagram types moved to `d2/` sub-module — `output.D2Node` → `d2.D2Node`, `output.NewD2Renderer` → `d2.NewD2Diagram`, etc. Import `github.com/larsartmann/go-output/d2`.
- **BREAKING**: DOT and Mermaid renderers moved to `graph/` sub-module — `output.DOTFromTableData` → `graph.DOTFromTableData`, `output.MermaidFlowchartRenderer` → `graph.MermaidFlowchartRenderer`. Import `github.com/larsartmann/go-output/graph`.
- `RenderTableData` now returns `UnsupportedFormatError` for D2, Mermaid, and DOT formats (use sub-module constructors directly).
- Added `Nodes()`, `Edges()`, `NodesPtr()`, `EdgesPtr()` accessor methods to `GraphRendererMixin` for cross-package use.
- Multi-module workspace expanded from 8 to 10 independent modules (see ADR 001, ADR 003).

### Added

- `JSONTableRenderer` — renders `TableData` as a JSON array of objects (implements `Renderer` + `TableRenderer`)
- `YAMLTableRenderer` — renders `TableData` as a YAML sequence of mappings (implements `Renderer` + `TableRenderer`)
- `TableData.ToMapSlice()` — converts tabular data to `[]map[string]string` for serialization
- `UnsupportedFormatError` — renamed from `ErrUnsupportedFormat` (follows Go naming conventions)
- `flake.nix` — Nix flake with devShell (Go 1.26.2, golangci-lint, gopls), treefmt-nix formatter, git-hooks.nix
- `.envrc` — direnv integration for automatic `nix develop` on cd
- `MIGRATION_TO_NIX_FLAKES_PROPOSAL.md` — migration plan from justfile to nix flakes
- Depguard whitelist for `examples/` module (scoped `examples/**/*.go` rule)
- CI: golangci-lint v2 (`github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

### Changed

- `RenderTableData` — all writer errors now wrapped with `fmt.Errorf("write X: %w", err)` for pinpoint failure reporting
- `FilledStrings` — uses `slices.Repeat` (Go 1.26 stdlib) instead of manual make+for loop
- `NewBrandedID` — simplified from `id.NewID[Brand, string](value)` to `id.NewID[Brand](value)` (inferred type arg)
- `sort/` module no longer depends on root — circular dependency eliminated
- `sort/sorter.go` (`Sorter[T]`, `New`, `WithLessFunc`, `Sort`) deleted — use stdlib `slices.SortStableFunc`
- `sort/sort_test.go` deleted (415 lines of tests for removed deprecated type)
- `sort/compare.go` updated: package docs, `ByField` signature, usage examples
- `sort/compare_test.go` rewritten as self-contained with zero external deps
- `sort/go.mod` stripped to 3 lines (zero dependencies)
- `enum/enum_test.go` no longer imports `internal/gentest` (inlined helper)
- `cmdguard/cmdguard_test.go` no longer imports `internal/gentest` (inlined helper)
- Root `go.mod` no longer requires or replaces `sort/` module
- `integration/go.mod` no longer requires or replaces `sort/` module
- `userjourney_test.go` uses stdlib `slices.SortStableFunc` + `cmp.Compare`
- `integration/workflow_test.go` uses stdlib sort instead of deprecated `sort.New`
- `.gitignore` — added `result` and `.direnv/` for Nix artifacts

### Removed

- `ErrUnsupportedFormat` — renamed to `UnsupportedFormatError` (breaking change)
- `sort/sorter.go` — deprecated `Sorter[T]` type (circular dependency with root)
- `sort/sort_test.go` — tests for deleted type
- `sort/go.sum` — no longer needed (zero deps)

## [0.4.0] - 2026-05-17

### Added

- MIT license (replaced PROPRIETARY)
- README.md rewritten for general audience and public launch
- 27 doc comments on exported symbols across `d2_enum.go`, `color.go`, `graph.go`, `sort.go`, `enum/enum.go`
- `Shape` type with `ShapeTable`/`ShapeTree`/`ShapeGraph` constants
- `formatCapabilities` map — single source of truth for format-to-shape mapping
- `Format.Supports(Shape)`, `Format.Shapes()`, `FormatsForShape(Shape)` methods
- `ParseShape()`, `Shape.IsValid()`, `Shape.AllowedValues()`, `Shape.String()` enum methods
- `ErrInvalidShape` sentinel error
- `AllShapes` slice for shape iteration
- ADR 002: Shape capability matrix decision record

### Changed

- `FormatCategory`, `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` deprecated — use `Supports(Shape)` instead
- `TreeNode` and `TreeOutputRenderer` extracted from `format.go` to `tree.go`
- `TableData`, `RowEdge`, and `tableDataBase` extracted from `format.go`/`html.go` to `tabledata.go`
- `GraphRendererMixin` moved from `dot.go` to `graph.go`
- `format.go` reduced from 373 to 291 lines (under 350-line limit)
- `dot.go` reduced from 253 to 199 lines

### Removed

- `PLAN.md` — stale duplicate of AGENTS.md with incorrect info

## [0.3.0] - 2026-05-13

### Added

- Multi-module workspace: `enum/`, `escape/`, `cmdguard/`, `table/`, `sort/`, `integration/`, `examples/`
- `enum` package with generic `Parse`, `Contains`, `AllowedStrings`, `AllowedValues` utilities
- `escape` package with format-specific HTML/XML escaping (using stdlib `html.EscapeString`)
- `cmdguard` package with generic `EnumFlag[T]` CLI flag parser
- `table` package with lipgloss-styled terminal tables (isolated from root)
- Branded ID types via `github.com/larsartmann/go-branded-id` for D2NodeID, TreeNodeID, etc.
- ADR 001: Multi-module workspace decision record
- `go.work` support for local development (gitignored)
- `FormatTSV` constant and TSV formatter implementation
- `FormatXML` constant and XML formatter (`MarshalXML`, `MarshalXMLIndent`, `XMLWriter`, `MarshalXMLFromTableData`)

### Changed

- `ParseFormat`, `ParseSortBy`, `ParseColorMode` refactored to use `enum` helpers
- `AllowedValues()` methods refactored to use `enum` helpers
- `escape/` uses `html.EscapeString()` from stdlib instead of custom implementation
- Go toolchain updated to 1.26+
- `sort/` deprecated with notice pointing to `slices.SortStableFunc` + `cmp.Compare`
- Format classification uses map-based lookup instead of switch statements

### Deprecated

- `OutputFormat` type alias — will be removed in v2.0
- `OutputFormat*` constants — will be removed in v2.0
- `ParseOutputFormat()` function — will be removed in v2.0
- `sort.Sorter[T]` — use stdlib `slices.SortStableFunc` instead

## [0.2.0] - 2026-04-30

### Changed

- Updated dependencies (`charmbracelet/x/exp/golden`)

## [0.1.0] - 2026-01-01

### Added

- Initial release with 12 output formats: Table, JSON, CSV, TSV, Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT
- `Renderer` interface with `Render() (string, error)`
- `TableRenderer` interface with `SetHeaders`/`AddRow`
- `GraphRenderer` interface with `SetNodes`/`SetEdges`
- `TreeOutputRenderer` interface with `SetRoot`
- `StreamingRenderer` interface with `Stream(io.Writer)`
- `ColorMode` enum (Auto/Always/Never) with terminal detection
- `SortBy` enum for sort field selection
- Opt-in renderer registry (`Register`/`Create`)
- D2 diagram renderer with shapes, arrows, SQL tables, classes, user journeys
- DOT/Graphviz renderer with `GraphRendererMixin`
- Mermaid flowchart renderer
- HTML table renderer with escaping
- Streaming HTML renderer for large datasets
- `MustRender(r)` helper for tests/examples
