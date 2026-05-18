package output

import (
	"fmt"
	"io"
	"strings"
)

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

// MarshalCSVFromTableData marshals TableData as CSV with a header row.
func MarshalCSVFromTableData(data *TableData) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var builder strings.Builder

	csvWriter := NewCSVWriter(&builder)

	if len(data.Headers) > 0 {
		if err := csvWriter.WriteHeader(data.Headers); err != nil {
			return nil, fmt.Errorf("write csv header: %w", err)
		}
	}

	for _, row := range data.Rows {
		if err := csvWriter.WriteRow(row); err != nil {
			return nil, fmt.Errorf("write csv row: %w", err)
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return nil, fmt.Errorf("flush csv writer: %w", err)
	}

	return []byte(builder.String()), nil
}
