package delimited

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Registers CSV Table marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatCSV, output.ShapeTable)
	output.RegisterTableMarshaler(output.FormatCSV,
		func(w io.Writer, data *output.Table, _ output.RenderOptions) error {
			return WriteCSV(w, data)
		})
}

// CSVWriter writes CSV output.
type CSVWriter struct {
	writer *DelimitedWriter
}

// NewCSVWriter creates a new CSVWriter.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		writer: NewDelimitedWriter(w, ',', "csv"),
	}
}

// WriteHeader writes the header row.
func (c *CSVWriter) WriteHeader(cols []string) error {
	return c.writer.WriteRow(cols, "csv header")
}

// WriteRow writes a single row.
func (c *CSVWriter) WriteRow(values []string) error {
	return c.writer.WriteRow(values, "csv row")
}

// WriteFooter writes a footer row (semantically equivalent to WriteRow for CSV).
func (c *CSVWriter) WriteFooter(values []string) error {
	return c.writer.WriteRow(values, "csv footer")
}

// WriteRows writes multiple rows.
func (c *CSVWriter) WriteRows(values [][]string) error {
	return c.writer.WriteRows(values, "csv")
}

// Flush flushes the writer.
func (c *CSVWriter) Flush() {
	c.writer.Flush()
}

// Error returns any error from the writer.
func (c *CSVWriter) Error() error {
	return c.writer.Error()
}

// Writer is the common interface implemented by CSVWriter and TSVWriter.
// It exposes the row-streaming API (header / row / footer) plus Flush and
// Error for inspecting the underlying buffered writer.
type Writer interface {
	WriteHeader(cols []string) error
	WriteRow(values []string) error
	WriteFooter(values []string) error
	Flush()
	Error() error
}

// tableDataWriter is the unexported alias used internally so the
// marshalFromTable generic helper doesn't depend on the public Writer
// type. tableDataWriter == delimited.Writer structurally; keep them in sync.
type tableDataWriter = Writer

// marshalFromTable writes data to the provided writer using the writer
// constructor supplied. Returns nil when data is nil (callers can treat that
// as "nothing to do"). All row writing, error wrapping, and flush logic live
// here; CSV/TSV callers just supply the right writer constructor.
func marshalFromTable(
	data *output.Table,
	name string,
	newWriter func(io.Writer) tableDataWriter,
) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var builder strings.Builder

	if err := writeDelimited(&builder, data, name, newWriter); err != nil {
		return nil, err
	}

	return []byte(builder.String()), nil
}

// writeDelimited is the streaming core shared by WriteCSV/WriteTSV (which
// pass their target io.Writer directly) and marshalFromTable (which buffers
// through strings.Builder). Pulling the body out eliminates a near-identical
// 22-line copy/paste between the two WriteCSV/WriteTSV entry points.
func writeDelimited(
	w io.Writer,
	data *output.Table,
	name string,
	newWriter func(io.Writer) tableDataWriter,
) error {
	if data == nil {
		return nil
	}

	dw := newWriter(w)

	if len(data.Headers) > 0 {
		if err := dw.WriteHeader(data.Headers); err != nil {
			return fmt.Errorf("write %s header: %w", name, err)
		}
	}

	for _, row := range data.Rows {
		if err := dw.WriteRow(row); err != nil {
			return fmt.Errorf("write %s row: %w", name, err)
		}
	}

	if data.HasFooter() {
		if err := dw.WriteFooter(data.Footer); err != nil {
			return fmt.Errorf("write %s footer: %w", name, err)
		}
	}

	dw.Flush()

	if err := dw.Error(); err != nil {
		return fmt.Errorf("flush %s writer: %w", name, err)
	}

	return nil
}

// MarshalCSVFromTable marshals Table as CSV with a header row.
func MarshalCSVFromTable(data *output.Table) ([]byte, error) {
	return marshalFromTable(data, "csv", func(w io.Writer) tableDataWriter {
		return NewCSVWriter(w)
	})
}
