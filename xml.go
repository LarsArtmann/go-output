package output

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// MarshalXML encodes v to XML.
func MarshalXML(v any) ([]byte, error) {
	return marshal("xml", xml.Marshal, v)
}

// MarshalXMLIndent encodes v to indented XML.
func MarshalXMLIndent(v any, prefix, indent string) ([]byte, error) {
	return marshalIndent("xml", xml.MarshalIndent, v, prefix, indent)
}

// XMLWriter writes XML output.
type XMLWriter struct {
	writer   *strings.Builder
	rowCount int
}

// NewXMLWriter creates a new XMLWriter.
func NewXMLWriter() *XMLWriter {
	return &XMLWriter{
		writer:   new(strings.Builder),
		rowCount: 0,
	}
}

// WriteHeader writes the XML header and opening tags.
func (x *XMLWriter) WriteHeader(cols []string) error {
	x.writer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	x.writer.WriteString("<table>\n")
	x.writer.WriteString("  <headers>\n")

	writeMarkupColumns(x.writer, cols, "    ", xmlEscape)

	x.writer.WriteString("  </headers>\n")
	x.writer.WriteString("  <rows>\n")

	return nil
}

// WriteRow writes a single row.
func (x *XMLWriter) WriteRow(values []string) error {
	writeMarkupRow(x.writer, values, "row", "cell", "    ", xmlEscape)
	x.rowCount++

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

// String returns the XML output.
func (x *XMLWriter) String() string {
	x.writer.WriteString("  </rows>\n")
	x.writer.WriteString("</table>\n")

	return x.writer.String()
}

// MarshalXMLFromTableData marshals TableData to XML.
func MarshalXMLFromTableData(data *TableData) ([]byte, error) {
	if data == nil {
		return []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<table/>\n"), nil
	}

	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	b.WriteString("<table>\n")

	if len(data.Headers) > 0 {
		b.WriteString("  <headers>\n")

		writeMarkupColumns(&b, data.Headers, "    ", xmlEscape)

		b.WriteString("  </headers>\n")
	}

	b.WriteString("  <rows>\n")

	for _, row := range data.Rows {
		writeMarkupRow(&b, row, "row", "cell", "    ", xmlEscape)
	}

	b.WriteString("  </rows>\n")
	b.WriteString("</table>\n")

	return []byte(b.String()), nil
}
