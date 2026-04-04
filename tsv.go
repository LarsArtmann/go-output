package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// TSVWriter writes TSV (Tab-Separated Values) output.
type TSVWriter struct {
	writer *csv.Writer
}

// NewTSVWriter creates a new TSVWriter.
func NewTSVWriter(w io.Writer) *TSVWriter {
	// CSV writer with comma delimiter, but we'll write tabs manually
	writer := csv.NewWriter(w)
	writer.Comma = '\t' // Use tab as delimiter

	return &TSVWriter{
		writer: writer,
	}
}

// write writes a row with the given description.
func (t *TSVWriter) write(cols []string, description string) error {
	err := t.writer.Write(cols)
	if err != nil {
		return fmt.Errorf("write %s %s: %w", description, cols, err)
	}

	return nil
}

// WriteHeader writes the header row.
func (t *TSVWriter) WriteHeader(cols []string) error {
	return t.write(cols, "tsv header")
}

// WriteRow writes a single row.
func (t *TSVWriter) WriteRow(values []string) error {
	return t.write(values, "tsv row")
}

// WriteRows writes multiple rows.
func (t *TSVWriter) WriteRows(values [][]string) error {
	err := t.writer.WriteAll(values)
	if err != nil {
		return fmt.Errorf("write tsv rows (count=%d): %w", len(values), err)
	}

	return nil
}

// Flush flushes the writer.
func (t *TSVWriter) Flush() {
	t.writer.Flush()
}

func (t *TSVWriter) Error() error {
	return writerError(t.writer, "tsv")
}

// MarshalTSV marshals data as TSV.
func MarshalTSV(data any) ([]byte, error) {
	var b strings.Builder

	w := NewTSVWriter(&b)

	// Handle slices of slices or structs - simplified for common cases
	switch v := data.(type) {
	case [][]string:
		for _, row := range v {
			err := w.WriteRow(row)
			if err != nil {
				return nil, fmt.Errorf(
					"write tsv row to %s: %w",
					b.String()[:min(50, len(b.String()))],
					err,
				)
			}
		}
	case []string:
		err := w.WriteRow(v)
		if err != nil {
			return nil, fmt.Errorf("write tsv single row %v: %w", v, err)
		}
	default:
		return nil, fmt.Errorf(
			"unsupported type %T for TSV marshaling to %s",
			data,
			b.String()[:min(50, len(b.String()))],
		)
	}

	w.Flush()

	err := w.Error()
	if err != nil {
		return nil, fmt.Errorf(
			"flush tsv writer for %s: %w",
			b.String()[:min(50, len(b.String()))],
			err,
		)
	}

	return []byte(b.String()), nil
}
