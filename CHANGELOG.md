# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.11.0] - 2026-06-17

### Added

- **`nix run .#test-race`** — new nix app running `go test -race -count=1` for the `nom/` and `tui/` concurrency-sensitive modules.
- **`ActivityStatus.Interest()`** — priority ordering is now a method on the enum type, replacing the standalone `activityInterest()` function.
- **Completed-subtree collapsing** — under height pressure, completed children are elided to prioritize active work (running, failed, pending).
- **Completion percentage in summary bar** — shows `(N%)` where N = (completed + failed) / total.
- **iTerm2 synchronized updates** — inline renderer frames wrapped in `\x1b[?2026h/l` for flicker-free redraws.
- **E2E smoke test** for inline renderer lifecycle (Start → events → Refresh → Stop → Finish).
- **Benchmarks** for `RenderWithWidth` and `childPriority` hot paths.
- **Dedicated tests** for `elideCompletedUnderPressure` and `ActivityStatus.Interest()`.

### Changed

- **TOML renderer wraps rows in `[[row]]`** — TOML cannot encode a bare `[]map[string]string` as document root; rows are now nested under a named key producing valid array-of-tables syntax.
- **ANSI parsing delegated to `charmbracelet/x/ansi`** — `StripANSI`, `VisibleWidth`, `TruncateVisible` now delegate to `ansi.Strip`, `ansi.StringWidth`, `ansi.Truncate`, removing 128 lines of hand-rolled scanner code. Direct `runewidth` dependency eliminated.
- **`Finish()` now calls `Stop()` first** — prevents concurrent tree access between the background render goroutine and the final render.
- **`RenderNode` signature documented** — the second parameter is now named `visibleNodes` (was `_ []*TreeNode`), honestly documenting it as reserved for future width-aware truncation.
- **`log.Printf` → `slog.Error`** in TUI lifecycle per how-to-golang policy.
- **Test refresh rendezvous** — replaced `time.Sleep(50ms)` with deterministic `renderNotify` channel in `TestInlineRenderer_Refresh_TriggersRender`.

### Fixed

- **TOML round-trip test failure** — `nix run .#test` was red on integration/serialization due to `toml: cannot encode a []map[string]string as a document root`.
- **TUI data races** — `go test -race` in tui/ failed on 18 tests; reporter tests now use `newTestReporter()` helper that prevents the real Bubble Tea program from starting.
- **Registry test pollution** — `TestRegisterTableDataMarshaler_ConcurrentAccess` leaked "race-test-\*" formats into the global registry; `TestRegisteredTableDataFormats` now checks known-good formats instead of asserting all are valid.
- **Dead code in `elideCompletedUnderPressure`** — removed unreachable guard `visibleCount+maxHeight <= 0`.

## [0.9.0] - 2026-06-12

### Added

- **`RenderAnyData`** — new registry-based dispatch for rendering arbitrary `any` data (not just `TableData`). JSON, YAML, and TOML register `AnyDataMarshaler` handlers via `init()`.
- **`RegisterAnyDataMarshaler`** / **`RegisteredAnyDataFormats`** — register and introspect any-data marshalers.
- **`RegisteredTableDataFormats`** — returns all formats with registered `TableDataMarshaler`s.
- **`TableData.AddRowChecked(row []string) error`** — fail-fast row addition that returns `ErrColumnMismatch` when column count differs from headers.
- **`nom.Event*` constants** — typed event string constants (`EventWorkflowStarted`, `EventActivityCompleted`, etc.) replacing bare string literals for safer event dispatch.
- **All 16 formats register as `TableDataMarshaler`** — D2, DOT, Mermaid, and PlantUML now register via `init()`, making `RenderTableData()` dispatch all 16 formats when sub-modules are imported (previously only 10 were registered).
- `testhelpers.AssertLineCount`, `AssertLastLineContains`, `AssertErrorContains` — new shared test assertions.
- `tui/colors.go` — extracted all color and layout constants from scattered inline values.
- `escape/fuzz_test.go` — fuzz tests for all escape functions.

### Changed

- **Generic `formatRegistry[T]`** replaces 3 separate mutex+map registry patterns (shape capabilities, table-data marshalers, any-data marshalers). Eliminates duplicate mutex boilerplate.
- `RenderTableData()` doc updated: all 16 formats are now dispatchable with sub-modules imported.
- NOM subscriber handlers simplified — reduced complexity in event routing.
- NOM `DisplayState` embedded directly, removing indirection.
- TUI reporter and view rendering refactored for clarity.
- D2, graph, plantuml `dt.Build()` errors now propagated instead of silently discarded.
- NOM timing cache async save errors now logged instead of swallowed.

### Deprecated

- `delimited.MarshalTSV(data any)` — use `MarshalTSVFromTableData` or `TSVWriter` directly.
- `graph.NewGraphNodeID` / `graph.NewGraphNodeLabel` — use `output.NewBrandedID` directly.
- `nom.ParseActivityID` / `nom.ParseWorkflowID` — use direct type conversion with manual validation.

### Internal

- Coverage improvements across 5 modules (nom tree/render/subscriber, table registry, serialization, testhelpers).
- Test deduplication across 22 files — shared helpers extracted to `testhelpers/`, table-driven patterns.
- Registry benchmarks for generic `formatRegistry[T]`.
- Comprehensive hardening sprint status reports in `docs/status/`.
- Graph DOMAIN_LANGUAGE updated with registry dispatch documentation.

## [0.8.0] - 2026-06-10

### Added

- **TUI keyboard navigation** — arrow keys, mouse scroll, click selection on dependency tree nodes.
- **`tui.SetCancelFunc`** — allows Ctrl+C to cancel workflow context via cancel function.
- **NOM secondary dependencies** — `DependenciesAccessor` and secondary parent labels in tree rendering.
- **NOM `SetDisplayMode`** — switch between NOM and Universal display modes.
- NOM/TUI golden file tests and event sequence tests.

### Changed

- All `panic()` calls eliminated from production code — replaced with error returns.
- NOM `AddActivity` made idempotent — prevents duplicate children.
- NOM dependency tree locking fixed — `GetDependencyTree` now properly synchronized.

### Fixed

- Timing cache race condition in async save.
- TUI help overlay and viewport scroll issues.
- TUI ticks now allowed in Idle state for NOM sync.
- Pre-registered activities preserved correctly.

## [0.7.0] - 2026-06-09

### Added

- **`nom/` module** — NOM-style real-time progress visualization with dependency trees, timing cache (CSV-persisted), and event-driven activity tracking. Zero imports from root module. `NOMStyleSubscriber` implements `EventSubscriber` with string-based event routing and type-assertion accessor interfaces.
- **`tui/` module** — Bubble Tea interactive TUI for workflow progress display. `BubbleTeaProgressReporter` with state machine (`WorkflowState`: Idle→Running→Completed/Errored), step-based progress, and NOM integration. Depends on `nom/`.
- `TableData.Validate()` now rejects nil rows — catches `nil` in `Rows [][]string` before rendering.
- Race test for `RegisterFormatShapes` confirms thread-safety of the shape capability matrix.

### Changed

- **BREAKING**: `RenderTableData()` changed from variadic `opts ...RenderOptions` to single `opts RenderOptions` — only `opts[0]` was ever used, the variadic signature was misleading.
- `GraphRendererMixin` API refactored for cleaner method signatures.
- String escaping optimized across modules — `strings.NewReplacer` for single-pass replacements.
- File renames for naming consistency across the project.
- Transitive dependencies updated across all 15 modules.

### Fixed

- Gosec and wrapcheck lint violations resolved.
- Critical bugs and config gaps from comprehensive quality audit.
- Flaky `TestTimingCache_SaveAndLoad` — async `saveAsync()` goroutine now completes before test cleanup.

### Internal

- `go.work.example` and `AGENTS.md` updated for 15 modules (was 13).
- Coverage table updated: nom 93.1%, tui 84.2%.
- Nix configuration standardized with format and build checks.

## [0.6.3] - 2026-06-03

### Changed

- Updated `charmbracelet/ultraviolet` from `v0.6.0-20260525` to `v0.6.0-20260601` — includes fix for modified Kitty keyboard navigation/function key releases (affects `table/`, `examples/`, `integration/` modules via lipgloss/v2 transitive dependency).
- Updated `golang.org/x/exp` from `v0.0.0-20260527` to `v0.0.0-20260529`.
- Updated nixpkgs flake input to latest revision.

### Internal

- Eliminated all code clones at `art-dupl` threshold 30 across test files in `markup/`, `serialization/`, `table/`, `graph/`, `plantuml/`, `testhelpers/`, and root module.

## [0.6.2] - 2026-06-01

### Fixed

- All internal cross-module `go.mod` references upgraded from `v0.0.0` pseudo-versions to canonical `v0.6.2` tags across all 14 modules. Resolves the chicken-and-egg release issue where v0.6.1 was tagged before dependency versions were bumped.
- `integration/go.mod`: fixed `d2`, `graph`, `markup`, `plantuml`, `table` references from `v0.0.0` to `v0.6.2`.
- `plantuml/go.mod`: fixed `testhelpers/graphtest` reference from `v0.0.0-00010101000000-000000000000` to `v0.6.2`.
- `d2/go.mod`, `graph/go.mod`, `serialization/go.mod`: fixed `testhelpers/graphtest` references from `v0.0.0` to `v0.6.2`.
- `enum/go.mod`, `testhelpers/graphtest/go.mod`: fixed root and testhelpers references from `v0.0.0` to `v0.6.2`.

### Changed

- Mono-version tagging policy documented in `AGENTS.md` — all 14 modules release in lockstep.

### Added

- `docs/research/go-error-family-adoption-report.html` — comprehensive PRO/CONTRA analysis evaluating `github.com/larsartmann/go-error-family v0.3.0` adoption. Verdict: Do Not Adopt (Yet); Strategy B recommended (add to `examples/` module only).

## [0.6.1] - 2026-05-30

### Added

- **Footer row** (`TableData.Footer []string`) — optional totals/summary row on `TableData`. Tabular renderers render it visually: CSV/TSV append as last row, HTML uses `<tfoot>` with `footer-cell` CSS class, XML uses `<footer>`, Markdown adds separator + bold row (with column alignment), AsciiDoc appends footer row, Terminal Table uses bold styling. Data formats (JSON/YAML/TOML/JSONL) skip footer.
- `TableData.GetFooter()`, `TableData.HasFooter()`, `TableDataStore.SetFooter()`, `TableDataStore.HasFooter()` — accessor methods for footer row.
- `TableData.Validate()` — validates footer column count matches headers. Returns `errColumnMismatch` on mismatch. Wired into `RenderTableData()` for automatic validation.
- `MarkdownTable.SetFooter()` — sets footer row on Markdown table renderer (inherits column alignment).
- `table.SetFooter()` — adds bold-styled footer row to lipgloss terminal table. Tracks `footerRowIndex` for correct bold styling on multiple calls.
- `table.FooterProvider` — optional interface checked by `FromTableData()` for automatic footer styling.
- `CSVWriter.WriteFooter()`, `TSVWriter.WriteFooter()` — explicit footer methods for streaming delimited output.
- Package doc.go files for 8 packages (output, delimited, d2, graph, markup, plantuml, serialization, testhelpers) — pkg.go.dev now shows proper package documentation.
- GoDoc examples for `Format.IsValid`, `ParseFormat`, `ColorMode`, `Shape`, `TableData.Validate`, `MustRender`.
- GoDoc comments on all exported struct fields in graph, tree, and d2 types (40+ fields).

### Changed

- `delimited.marshalFromTableData()` — extracted shared helper from `MarshalCSVFromTableData` and `MarshalTSVFromTableData`, eliminating ~70 lines of duplication.
- HTML footer cells now use `class="footer-cell"` for CSS targeting (both `HTMLRenderer` and `StreamingHTMLRenderer`).
- All 14 modules unified to `go 1.26.3`.
- `table.SetFooter()` now correctly tracks `footerRowIndex` — only the last footer row receives bold styling.
- Root module test coverage improved from 88.6% to 96.1%.

### Fixed

- `MarkdownTable.AsTableRenderer()` — adapter wrapping fluent MarkdownTable API as void-returning `TableRenderer` interface.
- `table.Table.AsTableRenderer()` — adapter wrapping lipgloss table builder as `TableRenderer` interface.
- `table.WithFooterStyle(func(lipgloss.Style) lipgloss.Style)` — composable footer styling for lipgloss tables.
- Alignment constants (`AlignmentLeft/Right/Center`) unexported to `alignmentLeft/Right/Center` — were documented as unexported but actually exported.
- `UnsupportedFormatError.Unwrap()` removed — returned nil, semantically identical to not having it.
- ADR 004: Footer row design decision record.
- Coverage: gentest 80.8%→96.2%, integration 82.8%→95.5%, table 85.5%→100%, serialization 89.0%→91.4%.
- `table/table_test.go` split into `table/table_test.go` + `table/color_test.go` (391→274 lines, under 350-line limit).
- GoDoc on all exported testhelpers symbols: `ErrTest`, `ErrorRenderer`, `FixedRenderer`, `ErrWrite`, `ErrorWriter`, `WriteNThenFailWriter`.
- Package doc added to `testhelpers/graphtest`.
- `integration/go.mod` root dep reference fixed from `v0.5.0` to `v0.0.0`.
- AsciiDoc renderer now uses `HasFooter()` consistently with other renderers.
- `TestBrandedIDFormat` updated for go-branded-id v0.3.0 `%#v` output.

## [0.6.0] - 2026-05-25

### Added

- **JSONL** (`jsonl` format) — JSON Lines output, one JSON object per line. `serialization.NewJSONLTableRenderer()`, `serialization.MarshalJSONLFromTableData()`, `serialization.NewJSONLWriter()`. Supports `ShapeTable`.
- **AsciiDoc** (`asciidoc` format) — AsciiDoc table output. `markup.NewAsciiDocTableRenderer()`, `markup.MarshalAsciiDocFromTableData()`. Supports `ShapeTable`.
- **TOML** (`toml` format) — TOML serialization with table and tree support. `serialization.MarshalTOML()`, `serialization.UnmarshalTOML()`, `serialization.NewTOMLTableRenderer()`, `serialization.NewTOMLTreeRenderer()`. Supports `ShapeTable`, `ShapeTree`. Uses `github.com/pelletier/go-toml/v2`.
- **PlantUML** (`plantuml` format) — PlantUML component diagrams. `plantuml.NewPlantUMLDiagram()`, `plantuml.PlantUMLFromTableData()`, `plantuml.PlantUMLFromTree()`. Supports `ShapeTable`, `ShapeGraph`. New independent `plantuml/` module with zero external dependencies.
- 4 new `Format` constants: `FormatJSONL`, `FormatAsciiDoc`, `FormatTOML`, `FormatPlantUML` — 16 formats total.
- Shape capability matrix expanded: JSONL and AsciiDoc support `ShapeTable`; TOML supports `ShapeTable` + `ShapeTree`; PlantUML supports `ShapeTable` + `ShapeGraph`.
- `JSONTableRenderer` — renders `TableData` as a JSON array of objects (implements `Renderer` + `TableRenderer`)
- `YAMLTableRenderer` — renders `TableData` as a YAML sequence of mappings (implements `Renderer` + `TableRenderer`)
- `TableData.ToMapSlice()` — converts tabular data to `[]map[string]string` for serialization
- `UnsupportedFormatError` — renamed from `ErrUnsupportedFormat` (follows Go naming conventions)
- `TableDataMarshaler` registry — sub-modules register via `init()`, root has zero sub-module imports
- `TableDataStore` — exported from root for cross-module embedding
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

- **BREAKING**: CSV and TSV writers moved to `delimited/` sub-module — `output.NewCSVWriter` → `delimited.NewCSVWriter`, `output.NewTSVWriter` → `delimited.NewTSVWriter`. Import `github.com/larsartmann/go-output/delimited`.
- **BREAKING**: JSON and YAML marshalers/renderers moved to `serialization/` sub-module — `output.MarshalJSON` → `serialization.MarshalJSON`, `output.NewJSONTableRenderer` → `serialization.NewJSONTableRenderer`, etc. Import `github.com/larsartmann/go-output/serialization`.
- **BREAKING**: XML, HTML, and StreamingHTML renderers moved to `markup/` sub-module — `output.NewHTMLRenderer` → `markup.NewHTMLRenderer`, `output.MarshalXMLFromTableData` → `markup.MarshalXMLFromTableData`, etc. Import `github.com/larsartmann/go-output/markup`.
- **BREAKING**: D2 diagram types moved to `d2/` sub-module — `output.D2Node` → `d2.D2Node`, `output.NewD2Renderer` → `d2.NewD2Diagram`, etc. Import `github.com/larsartmann/go-output/d2`.
- **BREAKING**: DOT and Mermaid renderers moved to `graph/` sub-module — `output.DOTFromTableData` → `graph.DOTFromTableData`, `output.MermaidFromTableData` → `graph.MermaidFromTableData`. Import `github.com/larsartmann/go-output/graph`.
- `RenderTableData` now returns `UnsupportedFormatError` for D2, Mermaid, and DOT formats (use sub-module constructors directly).
- `RenderTableData` uses registry-based dispatch via `TableDataMarshaler` — sub-modules register via `init()`. Root has zero sub-module imports.
- `RenderTableData` — all writer errors now wrapped with `fmt.Errorf("write X: %w", err)` for pinpoint failure reporting
- `tableDataBase` exported as `TableDataStore` with `Data()` getter — enables cross-module embedding.
- `marshal()`, `unmarshal()`, `brandedValue()` exported as `MarshalFormat()`, `UnmarshalFormat()`, `BrandedValue()` — used by serialization/ and markup/.
- Multi-module workspace: 13 independent modules (see ADR 001, ADR 003).
- Root production code has zero imports from sub-modules (`delimited`, `serialization`, `markup`, `d2`, `graph`, `table`, `plantuml`).
- Root production code has zero `go-faster/yaml`, zero `go-toml/v2`, and zero `escape` imports (isolated in `serialization/` and `markup/`).
- `FilledStrings` — uses `slices.Repeat` (Go 1.26 stdlib) instead of manual make+for loop
- `NewBrandedID` — simplified from `id.NewID[Brand, string](value)` to `id.NewID[Brand](value)` (inferred type arg)
- Added `Nodes()`, `Edges()`, `NodesPtr()`, `EdgesPtr()` accessor methods to `GraphRendererState` for cross-package use
- `enum/enum_test.go` no longer imports `internal/gentest` (inlined helper)
- `.gitignore` — added `result` and `.direnv/` for Nix artifacts
- Deduplication sprints reduced code clones from 44 → 26 (41% total reduction)

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
- `GraphRendererState` moved from `dot.go` to `graph.go`
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
- DOT/Graphviz renderer with `GraphRendererState`
- Mermaid flowchart renderer
- HTML table renderer with escaping
- Streaming HTML renderer for large datasets
- `MustRender(r)` helper for tests/examples
