package output

import (
	"fmt"
	"io"
	"os"
)

// RenderOptions configures optional behavior for RenderTable.
type RenderOptions struct {
	// Title is used as the document title for HTML output and as a header for Markdown.
	Title string

	// Writer overrides the default os.Stdout output destination.
	Writer io.Writer

	// ColorMode controls terminal color output. Defaults to ColorModeAuto.
	ColorMode ColorMode
}

// TableMarshaler renders Table in a specific format to a writer.
type TableMarshaler func(w io.Writer, data *Table, opts RenderOptions) error

//nolint:gochecknoglobals // Registry for Table renderers, populated by sub-module init().
var tableDataRegistry = newFormatRegistry[TableMarshaler]()

// RegisterTableMarshaler registers a renderer for a format.
// Sub-modules call this from their init() to enable RenderTable dispatch.
func RegisterTableMarshaler(format Format, renderer TableMarshaler) {
	tableDataRegistry.register(format, renderer)
}

func getTableMarshaler(format Format) (TableMarshaler, bool) {
	return tableDataRegistry.get(format)
}

// All format renderers (Markdown, Tree, CSV, JSON, etc.) self-register from
// their respective sub-modules via init(). Root provides the registry but
// registers no format itself.

// RenderTable renders Table in the given format and writes to w (or os.Stdout).
// It supports all registered formats when respective sub-modules are imported.
// With all sub-modules imported, all 16 formats are available:
// table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.
func RenderTable(data *Table, format Format, opts RenderOptions) error {
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

	if m, ok := getTableMarshaler(format); ok {
		return m(w, data, opts)
	}

	return &UnsupportedFormatError{Format: format}
}

// UnsupportedFormatError is returned when RenderTable cannot handle a format.
type UnsupportedFormatError struct {
	Format Format
}

func (e *UnsupportedFormatError) Error() string {
	return fmt.Sprintf("render table data: format %q not supported", e.Format)
}

// UnknownRenderer renders arbitrary data (any) in a specific format to a writer.
type UnknownRenderer func(w io.Writer, data any, opts RenderOptions) error

//nolint:gochecknoglobals // Registry for any-data renderers, populated by sub-module init().
var unknownRegistry = newFormatRegistry[UnknownRenderer]()

// RegisterUnknownRenderer registers a renderer for arbitrary (non-Table) data.
// Sub-modules call this from their init() to enable RenderUnknown dispatch.
func RegisterUnknownRenderer(format Format, renderer UnknownRenderer) {
	unknownRegistry.register(format, renderer)
}

func getUnknownRenderer(format Format) (UnknownRenderer, bool) {
	return unknownRegistry.get(format)
}

// RenderUnknown renders arbitrary data in the given format and writes to w (or os.Stdout).
// Supports formats that registered an UnknownRenderer (typically JSON, YAML, TOML).
// Returns UnsupportedFormatError if no renderer is registered for the format.
func RenderUnknown(data any, format Format, opts RenderOptions) error {
	w := opts.Writer
	if w == nil {
		w = os.Stdout
	}

	if m, ok := getUnknownRenderer(format); ok {
		return m(w, data, opts)
	}

	return &UnsupportedFormatError{Format: format}
}

// RegisteredTableMarshalFormats returns all formats with registered TableMarshalers.
func RegisteredTableMarshalFormats() []Format {
	return tableDataRegistry.formats()
}

// RegisteredUnknownFormats returns all formats with registered UnknownRenderers.
func RegisteredUnknownFormats() []Format {
	return unknownRegistry.formats()
}
