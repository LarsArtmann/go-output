# ADR: Footer Row — Data on TableData vs RenderOptions

**Date:** 2026-05-28
**Status:** ACCEPTED & IMPLEMENTED

## Context

Tabular output formats need a "footer" row — a summary/totals row rendered after data rows. This feature spans 7 modules (root, delimited, markup, serialization, table, d2, integration).

Two design alternatives:

1. **Data on `TableData`** — Store `Footer []string` directly on `TableData`, propagate to sub-module renderers via `Data() *TableData`.
2. **Data on `RenderOptions`** — Pass footer as a `RenderOptions` field, require callers to pass it separately from data.

Additionally, each renderer needs to decide *how* to style the footer (bold, colored, CSS class, XML element).

## Decision

### 1. Footer data lives on `TableData`

```go
type TableData struct {
    Headers []string
    Rows    [][]string
    Footer  []string   // optional summary row
}
```

**Rationale:**
- Footer is *data*, not *render configuration*. It belongs with the data model.
- Sub-modules already receive `*TableData` — zero API changes needed.
- `RenderTableData()` callers don't need to pass footer separately.
- `TableData.Validate()` can check footer column count matches headers.

### 2. Styling is renderer-specific

Each renderer decides how to style the footer:

| Renderer | Styling |
|----------|---------|
| Markdown | Second separator + bold footer (inherits column alignment) |
| HTML | `<tfoot>` with `footer-cell` CSS class |
| XML | `<footer>` element wrapping `<cell>` elements |
| AsciiDoc | Bold cells via `*text*` |
| CSV/TSV | Plain row (no special styling — delimiter format) |
| Table (lipgloss) | Bold via `buildStyleFunc` + `WithFooterStyle` option |
| JSON/YAML/TOML/JSONL | **Skipped** — data formats don't have visual footers |
| Tree/Graph | **Skipped** — non-tabular shapes |

### 3. Validation is centralized

`TableData.Validate()` returns an error if footer column count ≠ header count. `RenderTableData()` calls `Validate()` before dispatching.

### 4. TableRenderer adapter pattern

Both `MarkdownTable` and `table.Table` use fluent/builder APIs (returning `*Self`) incompatible with the void-returning `TableRenderer` interface. Solution: `AsTableRenderer()` adapter methods that wrap the builder with void-returning delegates.

### 5. `WithFooterStyle` for table customization

The `table` module provides `WithFooterStyle(func(lipgloss.Style) lipgloss.Style)` for composable footer styling without exposing lipgloss internals.

## Consequences

- **Positive:** All tabular formats support footer with zero breaking API changes.
- **Positive:** Data formats skip footer cleanly — no special casing needed.
- **Positive:** Validation prevents malformed footers early.
- **Trade-off:** Footer on `TableData` means data formats receive but ignore it. Acceptable because the data is correct; it's just not rendered.
- **Trade-off:** Adapter pattern adds a small allocation per `AsTableRenderer()` call. Negligible.
