package output

import (
	"encoding/csv"
	"io"
)

type CSVWriter struct {
	writer *csv.Writer
}

func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{
		writer: csv.NewWriter(w),
	}
}

func (c *CSVWriter) WriteHeader(cols []string) error {
	return c.writer.Write(cols)
}

func (c *CSVWriter) WriteRow(values []string) error {
	return c.writer.Write(values)
}

func (c *CSVWriter) WriteRows(values [][]string) error {
	return c.writer.WriteAll(values)
}

func (c *CSVWriter) Flush() {
	c.writer.Flush()
}

func (c *CSVWriter) Error() error {
	return c.writer.Error()
}
