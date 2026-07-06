package delimited

import (
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Registers TSV Table marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatTSV, output.ShapeTable)
	output.RegisterTableMarshaler(output.FormatTSV,
		func(w io.Writer, data *output.Table, _ output.RenderOptions) error {
			return renderDelimitedTable(w, data, MarshalTSVFromTable, "tsv")
		})
}

// TSVWriter writes TSV (Tab-Separated Values) output.
type TSVWriter struct {
	writer *DelimitedWriter
}

// NewTSVWriter creates a new TSVWriter.
func NewTSVWriter(w io.Writer) *TSVWriter {
	return &TSVWriter{
		writer: NewDelimitedWriter(w, '\t', "tsv"),
	}
}

// WriteHeader writes the header row.
func (t *TSVWriter) WriteHeader(cols []string) error {
	return t.writer.WriteRow(cols, "tsv header")
}

// WriteRow writes a single row.
func (t *TSVWriter) WriteRow(values []string) error {
	return t.writer.WriteRow(values, "tsv row")
}

// WriteFooter writes a footer row (semantically equivalent to WriteRow for TSV).
func (t *TSVWriter) WriteFooter(values []string) error {
	return t.writer.WriteRow(values, "tsv footer")
}

// WriteRows writes multiple rows.
func (t *TSVWriter) WriteRows(values [][]string) error {
	return t.writer.WriteRows(values, "tsv")
}

// Flush flushes the writer.
func (t *TSVWriter) Flush() {
	t.writer.Flush()
}

// Error returns any error from the writer.
func (t *TSVWriter) Error() error {
	return t.writer.Error()
}

// MarshalTSVFromTable marshals Table as TSV with a header row.
func MarshalTSVFromTable(data *output.Table) ([]byte, error) {
	return marshalFromTable(data, "tsv", func(w io.Writer) tableDataWriter {
		return NewTSVWriter(w)
	})
}
