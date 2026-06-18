package output

import (
	"fmt"
	"io"
	"os"
)

// RenderOptions configures optional behavior for RenderTableData.
type RenderOptions struct {
	// Title is used as the document title for HTML output and as a header for Markdown.
	Title string

	// Writer overrides the default os.Stdout output destination.
	Writer io.Writer

	// ColorMode controls terminal color output. Defaults to ColorModeAuto.
	ColorMode ColorMode
}

// TableDataRenderer renders TableData in a specific format to a writer.
type TableDataRenderer func(w io.Writer, data *TableData, opts RenderOptions) error

//nolint:gochecknoglobals // Registry for TableData marshalers, populated by sub-module init().
var tableDataRegistry = newFormatRegistry[TableDataRenderer]()

// RegisterTableDataRenderer registers a marshaler for a format.
// Sub-modules call this from their init() to enable RenderTableData dispatch.
func RegisterTableDataRenderer(format Format, marshaler TableDataRenderer) {
	tableDataRegistry.register(format, marshaler)
}

func getTableDataMarshaler(format Format) (TableDataRenderer, bool) {
	return tableDataRegistry.get(format)
}

//nolint:gochecknoinits // Registers Markdown and Tree TableData marshalers for registry-based dispatch.
func init() {
	RegisterTableDataRenderer(FormatMarkdown, renderMarkdownTableData)
	RegisterTableDataRenderer(FormatTree, renderTreeTableData)
}

// RenderTableData renders TableData in the given format and writes to w (or os.Stdout).
// It supports all registered formats when respective sub-modules are imported.
// With all sub-modules imported, all 16 formats are available:
// table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.
func RenderTableData(data *TableData, format Format, opts RenderOptions) error {
	if data == nil {
		return nil
	}

	if err := data.Validate(); err != nil {
		return fmt.Errorf("render table data: %w", err)
	}

	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	if m, ok := getTableDataMarshaler(format); ok {
		return m(w, data, opts)
	}

	return &UnsupportedFormatError{Format: format}
}

// UnsupportedFormatError is returned when RenderTableData cannot handle a format.
type UnsupportedFormatError struct {
	Format Format
}

func (e *UnsupportedFormatError) Error() string {
	return fmt.Sprintf("render table data: format %q not supported", e.Format)
}

func renderMarkdownTableData(w io.Writer, data *TableData, opts RenderOptions) error {
	if opts.Title != "" {
		_, err := fmt.Fprintf(w, "# %s\n\n", opts.Title)
		if err != nil {
			return fmt.Errorf("write markdown title: %w", err)
		}

		_, err = fmt.Fprintf(w, "%d rows\n\n", data.RowCount())
		if err != nil {
			return fmt.Errorf("write markdown row count: %w", err)
		}
	}

	mdTable := NewMarkdownTableFromData(data)
	mdTable.SetColorMode(opts.ColorMode)

	out, err := mdTable.Render()
	if err != nil {
		return fmt.Errorf("render markdown: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}

	return nil
}

func renderTreeTableData(w io.Writer, data *TableData, opts RenderOptions) error {
	renderer := TreeRendererFromTableData(data)
	renderer.SetColorMode(opts.ColorMode)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render tree: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write tree output: %w", err)
	}

	return nil
}

// AnyDataRenderer renders arbitrary data (any) in a specific format to a writer.
type AnyDataRenderer func(w io.Writer, data any, opts RenderOptions) error

//nolint:gochecknoglobals // Registry for any-data marshalers, populated by sub-module init().
var anyDataRegistry = newFormatRegistry[AnyDataRenderer]()

// RegisterAnyDataRenderer registers a marshaler for arbitrary (non-TableData) data.
// Sub-modules call this from their init() to enable RenderAnyData dispatch.
func RegisterAnyDataRenderer(format Format, marshaler AnyDataRenderer) {
	anyDataRegistry.register(format, marshaler)
}

func getAnyDataMarshaler(format Format) (AnyDataRenderer, bool) {
	return anyDataRegistry.get(format)
}

// RenderAnyData renders arbitrary data in the given format and writes to w (or os.Stdout).
// Supports formats that registered an AnyDataRenderer (typically JSON, YAML, TOML).
// Returns UnsupportedFormatError if no marshaler is registered for the format.
func RenderAnyData(data any, format Format, opts RenderOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	if m, ok := getAnyDataMarshaler(format); ok {
		return m(w, data, opts)
	}

	return &UnsupportedFormatError{Format: format}
}

// RegisteredTableDataFormats returns all formats with registered TableDataMarshalers.
func RegisteredTableDataFormats() []Format {
	return tableDataRegistry.formats()
}

// RegisteredAnyDataFormats returns all formats with registered AnyDataMarshalers.
func RegisteredAnyDataFormats() []Format {
	return anyDataRegistry.formats()
}
