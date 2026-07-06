package markup

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteXML writes a Table as XML directly to the provided writer using
// NewXMLWriter — true element-level streaming.
func WriteXML(w io.Writer, data *output.Table) error {
	if data == nil {
		return nil
	}

	xw := NewXMLWriter(w)

	if err := xw.WriteHeader(data.Headers); err != nil {
		return fmt.Errorf("write xml header: %w", err)
	}

	if err := xw.WriteRows(data.Rows); err != nil {
		return fmt.Errorf("write xml rows: %w", err)
	}

	if err := xw.WriteFooter(); err != nil {
		return fmt.Errorf("write xml footer: %w", err)
	}

	return nil
}

// RenderXML renders a Table as an XML string.
func RenderXML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteXML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteAsciiDoc writes a Table as AsciiDoc to the provided writer.
// Note: AsciiDoc buffers in memory (no row-level streaming writer exists
// for this format). The output is identical to MarshalAsciiDocFromTable.
func WriteAsciiDoc(w io.Writer, data *output.Table) error {
	b, err := MarshalAsciiDocFromTable(data)
	if err != nil {
		return fmt.Errorf("marshal asciidoc: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write asciidoc output: %w", err)
	}

	return nil
}

// RenderAsciiDoc renders a Table as an AsciiDoc string.
func RenderAsciiDoc(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteAsciiDoc(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteHTML writes a Table as HTML to the provided writer.
// Note: HTML buffers in memory via renderHTMLTable (no row-level streaming
// writer exists for this format). For large tables, prefer CSV/JSON/JSONL
// which stream row-by-row via standard encoders.
func WriteHTML(w io.Writer, data *output.Table) error {
	return renderHTMLTable(w, data, output.RenderOptions{})
}

// RenderHTML renders a Table as an HTML string.
func RenderHTML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteHTML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
