package output

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/larsartmann/go-output/internal/escape"
)

// MarshalXML encodes v to XML.
func MarshalXML(v any) ([]byte, error) {
	data, err := xml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal xml (%T): %w", v, err)
	}
	return data, nil
}

// MarshalXMLIndent encodes v to indented XML.
func MarshalXMLIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := xml.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal xml indent (prefix=%q, indent=%q) for %T: %w",
			prefix,
			indent,
			v,
			err,
		)
	}
	return data, nil
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
	for _, col := range cols {
		x.writer.WriteString("    <column>")
		x.writer.WriteString(escape.XML(col))
		x.writer.WriteString("</column>\n")
	}
	x.writer.WriteString("  </headers>\n")
	x.writer.WriteString("  <rows>\n")
	return nil
}

// WriteRow writes a single row.
func (x *XMLWriter) WriteRow(values []string) error {
	x.writer.WriteString("    <row>\n")
	for _, val := range values {
		x.writer.WriteString("      <cell>")
		x.writer.WriteString(escape.XML(val))
		x.writer.WriteString("</cell>\n")
	}
	x.writer.WriteString("    </row>\n")
	x.rowCount++
	return nil
}

// WriteRows writes multiple rows.
func (x *XMLWriter) WriteRows(values [][]string) error {
	for i, row := range values {
		if err := x.WriteRow(row); err != nil {
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
		for _, h := range data.Headers {
			b.WriteString("    <column>")
			b.WriteString(escape.XML(h))
			b.WriteString("</column>\n")
		}
		b.WriteString("  </headers>\n")
	}

	b.WriteString("  <rows>\n")
	for _, row := range data.Rows {
		b.WriteString("    <row>\n")
		for _, cell := range row {
			b.WriteString("      <cell>")
			b.WriteString(escape.XML(cell))
			b.WriteString("</cell>\n")
		}
		b.WriteString("    </row>\n")
	}
	b.WriteString("  </rows>\n")
	b.WriteString("</table>\n")

	return []byte(b.String()), nil
}
