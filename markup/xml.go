package markup

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

//nolint:gochecknoinits // Registers XML Table marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatXML, output.ShapeTable)
	output.RegisterTableMarshaler(output.FormatXML, renderXMLTable)
}

// MarshalXML encodes v to XML.
func MarshalXML(v any) ([]byte, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal xml %T: %w", v, err)
	}

	return data, nil
}

// MarshalXMLIndent encodes v to indented XML.
func MarshalXMLIndent(v any, prefix, indent string) (result []byte, err error) {
	result, err = xml.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("marshal xml with prefix=%q indent=%q: %w", prefix, indent, err)
	}

	return result, nil
}

// XMLWriter writes XML output to an io.Writer.
type XMLWriter struct {
	Writer io.Writer
}

// NewXMLWriter creates a new XMLWriter.
func NewXMLWriter(w io.Writer) *XMLWriter {
	return &XMLWriter{Writer: w}
}

// WriteHeader writes the XML header and opening tags.
func (x *XMLWriter) WriteHeader(cols []string) error {
	if err := writeBytes(x.Writer, "write xml header", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"); err != nil {
		return err
	}

	if err := writeBytes(x.Writer, "write table open", "<table>\n"); err != nil {
		return err
	}

	if err := writeBytes(x.Writer, "write headers open", "  <headers>\n"); err != nil {
		return err
	}

	if err := writeMarkupColumns(x.Writer, cols, "    ", escape.XML); err != nil {
		return fmt.Errorf("write columns: %w", err)
	}

	if err := writeBytes(x.Writer, "write headers close", "  </headers>\n"); err != nil {
		return err
	}

	if err := writeBytes(x.Writer, "write rows open", "  <rows>\n"); err != nil {
		return err
	}

	return nil
}

// WriteRow writes a single row.
func (x *XMLWriter) WriteRow(values []string) error {
	if err := writeMarkupRow(x.Writer, values, "row", "cell", "    ", escape.XML); err != nil {
		return fmt.Errorf("write row: %w", err)
	}

	return nil
}

// WriteRows writes multiple rows.
func (x *XMLWriter) WriteRows(values [][]string) error {
	for i, row := range values {
		err := x.WriteRow(row)
		if err != nil {
			return fmt.Errorf("write xml row %d of %d (values=%v): %w", i, len(values), row, err)
		}
	}

	return nil
}

// WriteFooter writes the closing tags.
func (x *XMLWriter) WriteFooter() error {
	if err := writeBytes(x.Writer, "write rows close", "  </rows>\n"); err != nil {
		return err
	}

	return writeBytes(x.Writer, "write table close", "</table>\n")
}

// MarshalXMLFromTable marshals Table to XML.
func MarshalXMLFromTable(data *output.Table) ([]byte, error) {
	if data == nil {
		return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<table/>\n"), nil
	}

	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<table>\n")

	if len(data.Headers) > 0 {
		b.WriteString("  <headers>\n")

		if err := writeMarkupColumns(&b, data.Headers, "    ", escape.XML); err != nil {
			return nil, fmt.Errorf("write columns: %w", err)
		}

		b.WriteString("  </headers>\n")
	}

	b.WriteString("  <rows>\n")

	for _, row := range data.Rows {
		if err := writeMarkupRow(&b, row, "row", "cell", "    ", escape.XML); err != nil {
			return nil, fmt.Errorf("write row: %w", err)
		}
	}

	b.WriteString("  </rows>\n")

	if data.HasFooter() {
		b.WriteString("  <footer>\n")

		if err := writeMarkupRow(&b, data.Footer, "row", "cell", "    ", escape.XML); err != nil {
			return nil, fmt.Errorf("write footer: %w", err)
		}

		b.WriteString("  </footer>\n")
	}

	b.WriteString("</table>\n")

	return []byte(b.String()), nil
}

func renderXMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return WriteXML(w, data)
}
