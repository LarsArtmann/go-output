package delimited

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

func renderDelimitedTableData(
	w io.Writer,
	data *output.TableData,
	marshalFunc func(*output.TableData) ([]byte, error),
	formatName string,
) error {
	b, err := marshalFunc(data)
	if err != nil {
		return fmt.Errorf("render %s: %w", formatName, err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("write %s bytes: %w", formatName, err)
	}

	return nil
}

//nolint:gochecknoinits // Registers CSV TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatCSV,
		func(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
			return renderDelimitedTableData(w, data, MarshalCSVFromTableData, "csv")
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

// tableDataWriter is the common interface for CSV and TSV writers used by marshalFromTableData.
type tableDataWriter interface {
	WriteHeader(cols []string) error
	WriteRow(values []string) error
	Flush()
	Error() error
}

// marshalFromTableData marshals TableData using any delimited writer (CSV or TSV).
func marshalFromTableData(data *output.TableData, name string, newWriter func(io.Writer) tableDataWriter) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var builder strings.Builder

	w := newWriter(&builder)

	if len(data.Headers) > 0 {
		if err := w.WriteHeader(data.Headers); err != nil {
			return nil, fmt.Errorf("write %s header: %w", name, err)
		}
	}

	for _, row := range data.Rows {
		if err := w.WriteRow(row); err != nil {
			return nil, fmt.Errorf("write %s row: %w", name, err)
		}
	}

	if data.HasFooter() {
		if err := w.WriteRow(data.Footer); err != nil {
			return nil, fmt.Errorf("write %s footer: %w", name, err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("flush %s writer: %w", name, err)
	}

	return []byte(builder.String()), nil
}

// MarshalCSVFromTableData marshals TableData as CSV with a header row.
func MarshalCSVFromTableData(data *output.TableData) ([]byte, error) {
	return marshalFromTableData(data, "csv", func(w io.Writer) tableDataWriter {
		return NewCSVWriter(w)
	})
}
