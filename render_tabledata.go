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

	// GraphID is used as the graph identifier for DOT output.
	GraphID string

	// Writer overrides the default os.Stdout output destination.
	Writer io.Writer
}

// RenderTableData renders TableData in the given format and writes to w (or os.Stdout).
// It supports all tabular formats: csv, tsv, markdown, xml, yaml, d2, html, tree, mermaid, dot.
// Table and JSON formats are NOT handled — those require per-command customization
// (table for lipgloss styling, json for full struct marshaling).
//
// Returns ErrUnsupportedFormat if the format is table or json (caller should handle those).
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

	switch format {
	case FormatCSV:
		return renderCSVTableData(w, data)
	case FormatTSV:
		return renderTSVTableData(w, data)
	case FormatMarkdown:
		return renderMarkdownTableData(w, data, o)
	case FormatXML:
		return renderXMLTableData(w, data)
	case FormatYAML:
		return renderYAMLTableData(w, data)
	case FormatD2:
		return renderD2TableData(w, data)
	case FormatMermaid:
		return renderMermaidTableData(w, data)
	case FormatDOT:
		return renderDOTTableData(w, data, o)
	case FormatHTML:
		return renderHTMLTableData(w, data, o)
	case FormatTree:
		return renderTreeTableData(w, data)
	default:
		return &ErrUnsupportedFormat{Format: format}
	}
}

// ErrUnsupportedFormat is returned when RenderTableData cannot handle a format.
type ErrUnsupportedFormat struct {
	Format Format
}

func (e *ErrUnsupportedFormat) Error() string {
	return fmt.Sprintf("render table data: format %q not supported (handle table/json in caller)", e.Format)
}

func renderCSVTableData(w io.Writer, data *TableData) error {
	b, err := MarshalCSVFromTableData(data)
	if err != nil {
		return fmt.Errorf("render csv: %w", err)
	}

	_, err = w.Write(b)

	return err
}

func renderTSVTableData(w io.Writer, data *TableData) error {
	b, err := MarshalTSVFromTableData(data)
	if err != nil {
		return fmt.Errorf("render tsv: %w", err)
	}

	_, err = w.Write(b)

	return err
}

func renderMarkdownTableData(w io.Writer, data *TableData, opts RenderOptions) error {
	if opts.Title != "" {
		fmt.Fprintf(w, "# %s\n\n", opts.Title)
		fmt.Fprintf(w, "%d rows\n\n", data.RowCount())
	}

	mdTable := NewMarkdownTableFromData(data)
	out, err := mdTable.Render()
	if err != nil {
		return fmt.Errorf("render markdown: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}

func renderXMLTableData(w io.Writer, data *TableData) error {
	b, err := MarshalXMLFromTableData(data)
	if err != nil {
		return fmt.Errorf("render xml: %w", err)
	}

	fmt.Fprintln(w, string(b))

	return nil
}

func renderYAMLTableData(w io.Writer, data *TableData) error {
	renderer := NewYAMLTableRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render yaml: %w", err)
	}

	fmt.Fprint(w, out)

	return nil
}

func renderD2TableData(w io.Writer, data *TableData) error {
	diagram := D2FromTableData(data)
	out, err := diagram.Render()
	if err != nil {
		return fmt.Errorf("render d2: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}

func renderMermaidTableData(w io.Writer, data *TableData) error {
	renderer := MermaidFlowchartRenderer(data)
	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render mermaid: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}

func renderDOTTableData(w io.Writer, data *TableData, opts RenderOptions) error {
	renderer := DOTFromTableData(data)
	if opts.GraphID != "" {
		renderer.SetGraphID(opts.GraphID)
	}

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render dot: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}

func renderHTMLTableData(w io.Writer, data *TableData, opts RenderOptions) error {
	renderer := NewHTMLRenderer()
	renderer.SetData(data)

	title := opts.Title
	if title == "" {
		title = fmt.Sprintf("Data - %d rows", data.RowCount())
	}

	out, err := renderer.RenderFullHTML(title)
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}

func renderTreeTableData(w io.Writer, data *TableData) error {
	renderer := TreeRendererFromTableData(data)
	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render tree: %w", err)
	}

	fmt.Fprintln(w, out)

	return nil
}
