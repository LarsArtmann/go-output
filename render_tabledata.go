package output

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// RenderOptions configures optional behavior for RenderTableData.
type RenderOptions struct {
	// Title is used as the document title for HTML output and as a header for Markdown.
	Title string

	// GraphID is used as the graph identifier for DOT output.
	GraphID string

	// Writer overrides the default os.Stdout output destination.
	Writer io.Writer

	// ColorMode controls terminal color output. Defaults to ColorModeAuto.
	ColorMode ColorMode
}

// TableDataMarshaler renders TableData in a specific format to a writer.
type TableDataMarshaler func(w io.Writer, data *TableData, opts RenderOptions) error

var (
	//nolint:gochecknoglobals // Registry for TableData marshalers, populated by sub-module init().
	tableDataMarshalers = map[Format]TableDataMarshaler{}
	//nolint:gochecknoglobals // Mutex protects concurrent access to tableDataMarshalers.
	tableDataMarshalersMu sync.RWMutex
)

// RegisterTableDataMarshaler registers a marshaler for a format.
// Sub-modules call this from their init() to enable RenderTableData dispatch.
func RegisterTableDataMarshaler(format Format, marshaler TableDataMarshaler) {
	tableDataMarshalersMu.Lock()
	defer tableDataMarshalersMu.Unlock()

	tableDataMarshalers[format] = marshaler
}

func getTableDataMarshaler(format Format) (TableDataMarshaler, bool) {
	tableDataMarshalersMu.RLock()
	defer tableDataMarshalersMu.RUnlock()

	m, ok := tableDataMarshalers[format]

	return m, ok
}

// RenderTableData renders TableData in the given format and writes to w (or os.Stdout).
// It supports: csv, tsv, markdown, xml, yaml, html, tree (when respective sub-modules are imported).
//
// D2, Mermaid, and DOT are NOT handled — those require importing the d2 or graph
// sub-modules directly. Table and JSON formats also require per-command customization
// (table for lipgloss styling, json for full struct marshaling).
//
// Returns UnsupportedFormatError for unsupported formats (d2, mermaid, dot, table, json).
//
//nolint:cyclop,exhaustive // Dispatcher function with many format cases.
func RenderTableData(data *TableData, format Format, opts ...RenderOptions) error {
	if data == nil {
		return nil
	}

	var o RenderOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	w := o.Writer
	if w == nil {
		w = os.Stdout
	}

	// Registry-based dispatch for sub-module formats.
	if m, ok := getTableDataMarshaler(format); ok {
		return m(w, data, o)
	}

	// Direct dispatch for formats that live in root.
	switch format {
	case FormatMarkdown:
		return renderMarkdownTableData(w, data, o)
	case FormatTree:
		return renderTreeTableData(w, data, o)
	case FormatD2, FormatMermaid, FormatDOT:
		return &UnsupportedFormatError{Format: format}
	default:
		return &UnsupportedFormatError{Format: format}
	}
}

// UnsupportedFormatError is returned when RenderTableData cannot handle a format.
type UnsupportedFormatError struct {
	Format Format
}

func (e *UnsupportedFormatError) Error() string {
	return fmt.Sprintf("render table data: format %q not supported", e.Format)
}

func (e *UnsupportedFormatError) Unwrap() error {
	return nil
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
