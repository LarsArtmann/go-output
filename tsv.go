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

// WriteHeader writes the header row.
func (t *TSVWriter) WriteHeader(cols []string) error {
	if err := t.writer.Write(cols); err != nil {
		return fmt.Errorf("write tsv header %v: %w", cols, err)
	}
	return nil
}

// WriteRow writes a single row.
func (t *TSVWriter) WriteRow(values []string) error {
	if err := t.writer.Write(values); err != nil {
		return fmt.Errorf("write tsv row %v: %w", values, err)
	}
	return nil
}

// WriteRows writes multiple rows.
func (t *TSVWriter) WriteRows(values [][]string) error {
	if err := t.writer.WriteAll(values); err != nil {
		return fmt.Errorf("write tsv rows (count=%d): %w", len(values), err)
	}
	return nil
}

// Flush flushes the writer.
func (t *TSVWriter) Flush() {
	t.writer.Flush()
}

// Error returns any error encountered during writing.
func (t *TSVWriter) Error() error {
	if err := t.writer.Error(); err != nil {
		return fmt.Errorf("tsv writer error: %w", err)
	}
	return nil
}

// MarshalTSV marshals data as TSV.
func MarshalTSV(data any) ([]byte, error) {
	var b strings.Builder
	w := NewTSVWriter(&b)

	// Handle slices of slices or structs - simplified for common cases
	switch v := data.(type) {
	case [][]string:
		for _, row := range v {
			if err := w.WriteRow(row); err != nil {
				return nil, fmt.Errorf(
					"write tsv row to %s: %w",
					b.String()[:min(50, len(b.String()))],
					err,
				)
			}
		}
	case []string:
		if err := w.WriteRow(v); err != nil {
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
	if err := w.Error(); err != nil {
		return nil, fmt.Errorf(
			"flush tsv writer for %s: %w",
			b.String()[:min(50, len(b.String()))],
			err,
		)
	}

	return []byte(b.String()), nil
}
