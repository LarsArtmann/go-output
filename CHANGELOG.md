# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Changed

- **BREAKING**: CSV and TSV writers moved to `delimited/` sub-module — `output.NewCSVWriter` → `delimited.NewCSVWriter`, `output.NewTSVWriter` → `delimited.NewTSVWriter`. Import `github.com/larsartmann/go-output/delimited`.
- **BREAKING**: JSON and YAML marshalers/renderers moved to `serialization/` sub-module — `output.MarshalJSON` → `serialization.MarshalJSON`, `output.NewJSONTableRenderer` → `serialization.NewJSONTableRenderer`, etc. Import `github.com/larsartmann/go-output/serialization`.
- **BREAKING**: XML, HTML, and StreamingHTML renderers moved to `markup/` sub-module — `output.NewHTMLRenderer` → `markup.NewHTMLRenderer`, `output.MarshalXMLFromTableData` → `markup.MarshalXMLFromTableData`, etc. Import `github.com/larsartmann/go-output/markup`.
- **BREAKING**: D2 diagram types moved to `d2/` sub-module — `output.D2Node` → `d2.D2Node`, `output.NewD2Renderer` → `d2.NewD2Diagram`, etc. Import `github.com/larsartmann/go-output/d2`.
- **BREAKING**: DOT and Mermaid renderers moved to `graph/` sub-module — `output.DOTFromTableData` → `graph.DOTFromTableData`, `output.MermaidFromTableData` → `graph.MermaidFromTableData`. Import `github.com/larsartmann/go-output/graph`.
- `RenderTableData` now returns `UnsupportedFormatError` for D2, Mermaid, and DOT formats (use sub-module constructors directly).
- `RenderTableData` uses registry-based dispatch via `TableDataMarshaler` — sub-modules register via `init()`. Root has zero sub-module imports.
- `tableDataBase` exported as `TableDataBase` with `Data()` getter — enables cross-module embedding.
- `marshal()`, `unmarshal()`, `brandedValue()` exported as `MarshalFormat()`, `UnmarshalFormat()`, `BrandedValue()` — used by serialization/ and markup/.
- Multi-module workspace: 12 independent modules (see ADR 001, ADR 003).
- Root production code has zero imports from sub-modules (`delimited`, `serialization`, `markup`, `d2`, `graph`, `table`).
- Root production code has zero `go-faster/yaml` and zero `escape` imports (isolated in `serialization/` and `markup/`).

### Removed

- **BREAKING**: `format_deprecated.go` removed — `OutputFormat`, `FormatCategory`, `ParseOutputFormat`, `IsTableFormat()`, `IsTreeFormat()`, `IsGraphFormat()`, `Category()` all removed. Use `Format`, `Shape`, `ParseFormat()`, `Supports(Shape*)`, `Shapes()` instead.
- **BREAKING**: `sort/` module removed — `ByField` is trivially replaceable with `slices.SortStableFunc` + `cmp.Compare` from stdlib.
- **BREAKING**: `registry.go` removed — `Register()`, `Create()`, `Unregister()`, `RegisteredFormats()`, `IsRegistered()` all removed. Use direct constructors (`d2.NewD2Diagram()`, `serialization.NewJSONTableRenderer()`, etc.) instead.
- **BREAKING**: `SortBy` enum removed from root — `sort.go`, `ParseSortBy()`, `SortBy*` constants all removed. Zero external callers.
- **BREAKING**: `FilledStrings()` removed — use `slices.Repeat` (Go 1.26 stdlib) instead.
- **BREAKING**: `BrandedValue()` removed from `marshal.go` — zero callers after sub-module extraction.
- **BREAKING**: `MermaidFlowchartRenderer` and `MermaidTreeRenderer` removed — use `MermaidFromTableData` and `MermaidFromTree` instead.
- `gci` formatter removed from `.golangci.yml` (conflicted with `goimports` on local-prefix grouping in sub-modules).
- `ErrUnsupportedFormat` — renamed to `UnsupportedFormatError` (breaking change).
- `TestContainsString` in graph/ — tested stdlib `strings.Contains` (zero value).

### Added

- `JSONTableRenderer` — renders `TableData` as a JSON array of objects (implements `Renderer` + `TableRenderer`)
- `YAMLTableRenderer` — renders `TableData` as a YAML sequence of mappings (implements `Renderer` + `TableRenderer`)
- `TableData.ToMapSlice()` — converts tabular data to `[]map[string]string` for serialization
- `UnsupportedFormatError` — renamed from `ErrUnsupportedFormat` (follows Go naming conventions)
- `TableDataMarshaler` registry — sub-modules register via `init()`, root has zero sub-module imports
- `TableDataBase` — exported from root for cross-module embedding
- `RenderTableData` now accepts `RenderOptions.ColorMode` to control ANSI color output for terminal renderers.
- `table.New()` accepts `WithColorMode(ColorMode)` functional option — lipgloss styles conditionally applied based on terminal detection.
- `ASCIITreeRenderer.SetColorMode(ColorMode)` — depth-based ANSI color cycling, bold labels, dim connectors, cyan metadata.
- `MarkdownTable.SetColorMode(ColorMode)` — bold headers, dim separators when terminal detected.
- ColorMode auto-detection respects `NO_COLOR`, `CI`, `GITHUB_ACTIONS`, `GITLAB_CI`, `JENKINS_URL`, `BUILDKITE`, `GO_OUTPUT_FORCE_COLOR`, `FORCE_COLOR` env vars.
- `MarshalFormat()`, `UnmarshalFormat()`, `MarshalJSONIndent()` — exported helpers used by serialization/ and markup/ (BrandedValue removed)
- `flake.nix` — Nix flake with devShell (Go 1.26.2, golangci-lint, gopls), treefmt-nix formatter, git-hooks.nix
- `.envrc` — direnv integration for automatic `nix develop` on cd
- Depguard whitelist for `examples/` module (scoped `examples/**/*.go` rule)
- CI: golangci-lint v2 (`github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

### Changed

- `RenderTableData` — all writer errors now wrapped with `fmt.Errorf("write X: %w", err)` for pinpoint failure reporting
- `RenderTableData` — uses registry-based dispatch via `TableDataMarshaler` instead of direct function calls
- `FilledStrings` — uses `slices.Repeat` (Go 1.26 stdlib) instead of manual make+for loop
- `NewBrandedID` — simplified from `id.NewID[Brand, string](value)` to `id.NewID[Brand](value)` (inferred type arg)
- Added `Nodes()`, `Edges()`, `NodesPtr()`, `EdgesPtr()` accessor methods to `GraphRendererMixin` for cross-package use
- `enum/enum_test.go` no longer imports `internal/gentest` (inlined helper)
- `.gitignore` — added `result` and `.direnv/` for Nix artifacts
- Multi-module workspace expanded from 9 to 12 independent modules

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
