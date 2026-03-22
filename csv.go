package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

// CSVWriter writes CSV output.
type CSVWriter struct {
	writer *csv.Writer
}

// NewCSVWriter creates a new CSVWriter.
func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		writer: csv.NewWriter(w),
	}
}

// WriteHeader writes the header row.
func (c *CSVWriter) WriteHeader(cols []string) error {
	if err := c.writer.Write(cols); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	return nil
}

// WriteRow writes a single row.
func (c *CSVWriter) WriteRow(values []string) error {
	if err := c.writer.Write(values); err != nil {
		return fmt.Errorf("write csv row: %w", err)
	}
	return nil
}

// WriteRows writes multiple rows.
func (c *CSVWriter) WriteRows(values [][]string) error {
	if err := c.writer.WriteAll(values); err != nil {
		return fmt.Errorf("write csv rows: %w", err)
	}
	return nil
}

// Flush flushes the writer.
func (c *CSVWriter) Flush() {
	c.writer.Flush()
}

// Error returns any error encountered during writing.
func (c *CSVWriter) Error() error {
	if err := c.writer.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}
	return nil
}
