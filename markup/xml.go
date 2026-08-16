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

// marshalOrError wraps a stdlib marshal call with the canonical "marshal <what> %T: %w"
// error message that the public MarshalXML function emits. Caller supplies the
// raw (data, err) return pair of the stdlib encoder.
func marshalOrError(label string, v any, data []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf("%s %T: %w", label, v, err)
	}

	return data, nil
}

// MarshalXML encodes v to XML.
func MarshalXML(v any) ([]byte, error) {
	data, err := xml.Marshal(v)
	return marshalOrError("marshal xml", v, data, err)
}

// MarshalXMLIndent encodes v to indented XML.
func MarshalXMLIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := xml.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("marshal xml with prefix=%q indent=%q: %w", prefix, indent, err)
	}

	return data, nil
}

// XMLWriter writes XML output to an io.Writer.
type XMLWriter struct {
	Writer io.Writer
}

// NewXMLWriter creates a new XMLWriter.
func NewXMLWriter(w io.Writer) *XMLWriter {
	return &XMLWriter{Writer: w}
}

// WriteHeader writes the XML header and opening tags. The <headers> block
// is omitted when cols is empty, matching MarshalXMLFromTable.
func (x *XMLWriter) WriteHeader(cols []string) error {
	if err := writeBytes(x.Writer, "write xml header", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"); err != nil {
		return err
	}

	if err := writeBytes(x.Writer, "write table open", "<table>\n"); err != nil {
		return err
	}

	if len(cols) > 0 {
		if err := writeBytes(x.Writer, "write headers open", "  <headers>\n"); err != nil {
			return err
		}

		if err := writeMarkupColumns(x.Writer, cols, "    ", escape.XML); err != nil {
			return fmt.Errorf("write columns: %w", err)
		}

		if err := writeBytes(x.Writer, "write headers close", "  </headers>\n"); err != nil {
			return err
		}
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

// WriteFooter closes the rows block, writes the optional footer row, and
// closes the table element. Pass nil (or an empty slice) when the table has
// no footer — the <footer> block is emitted only for non-empty footers.
func (x *XMLWriter) WriteFooter(footer []string) error {
	if err := writeBytes(x.Writer, "write rows close", "  </rows>\n"); err != nil {
		return err
	}

	if len(footer) > 0 {
		if err := writeBytes(x.Writer, "write footer open", "  <footer>\n"); err != nil {
			return err
		}

		if err := writeMarkupRow(x.Writer, footer, "row", "cell", "    ", escape.XML); err != nil {
			return fmt.Errorf("write footer: %w", err)
		}

		if err := writeBytes(x.Writer, "write footer close", "  </footer>\n"); err != nil {
			return err
		}
	}

	return writeBytes(x.Writer, "write table close", "</table>\n")
}

// MarshalXMLFromTable marshals Table to XML. Non-nil tables delegate to
// WriteXML so the streaming writer and this function are byte-identical by
// construction (single formatting implementation).
func MarshalXMLFromTable(data *output.Table) ([]byte, error) {
	if data == nil {
		return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<table/>\n"), nil
	}

	var b strings.Builder
	if err := WriteXML(&b, data); err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

func renderXMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return WriteXML(w, data)
}
