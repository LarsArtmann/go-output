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

// write writes a row with the given description.
func (c *CSVWriter) write(cols []string, description string) error {
	err := c.writer.Write(cols)
	if err != nil {
		return fmt.Errorf("write %s %s: %w", description, cols, err)
	}

	return nil
}

// WriteHeader writes the header row.
func (c *CSVWriter) WriteHeader(cols []string) error {
	return c.write(cols, "csv header")
}

// WriteRow writes a single row.
func (c *CSVWriter) WriteRow(values []string) error {
	return c.write(values, "csv row")
}

// WriteRows writes multiple rows.
func (c *CSVWriter) WriteRows(values [][]string) error {
	err := c.writer.WriteAll(values)
	if err != nil {
		return fmt.Errorf("write csv rows (count=%d): %w", len(values), err)
	}

	return nil
}

// Flush flushes the writer.
func (c *CSVWriter) Flush() {
	c.writer.Flush()
}

func (c *CSVWriter) Error() error {
	return writerError(c.writer, "csv")
}

// writerError returns any error encountered during writing.
func writerError(w *csv.Writer, formatName string) error {
	err := w.Error()
	if err != nil {
		return fmt.Errorf("%s writer error: %w", formatName, err)
	}

	return nil
}
