package output

import (
	"encoding/csv"
	"fmt"
	"io"
)

// DelimitedWriter handles common logic for CSV and TSV writing.
type DelimitedWriter struct {
	writer *csv.Writer
	name   string
}

// NewDelimitedWriter creates a new DelimitedWriter.
func NewDelimitedWriter(w io.Writer, delimiter rune, name string) *DelimitedWriter {
	writer := csv.NewWriter(w)
	writer.Comma = delimiter

	return &DelimitedWriter{
		writer: writer,
		name:   name,
	}
}

// WriteRow writes a single row.
func (d *DelimitedWriter) WriteRow(cols []string, description string) error {
	err := d.writer.Write(cols)
	if err != nil {
		return fmt.Errorf("write %s %s: %w", description, cols, err)
	}

	return nil
}

// WriteRows writes multiple rows.
func (d *DelimitedWriter) WriteRows(values [][]string, description string) error {
	err := d.writer.WriteAll(values)
	if err != nil {
		return fmt.Errorf("write %s rows (count=%d): %w", description, len(values), err)
	}

	return nil
}

// Flush flushes the writer.
func (d *DelimitedWriter) Flush() {
	d.writer.Flush()
}

// Error returns any error encountered during writing.
func (d *DelimitedWriter) Error() error {
	err := d.writer.Error()
	if err != nil {
		return fmt.Errorf("%s writer error: %w", d.name, err)
	}

	return nil
}
