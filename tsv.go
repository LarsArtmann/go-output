package output

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

// TSVWriter writes TSV (Tab-Separated Values) output.
type TSVWriter struct {
	writer *csv.Writer
}

// ErrUnsupportedType is returned when an unsupported type is provided for TSV marshaling.
var ErrUnsupportedType = errors.New("unsupported type")

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
	var builder strings.Builder

	tsvWriter := NewTSVWriter(&builder)

	if err := writeTSVData(tsvWriter, data); err != nil {
		return nil, fmt.Errorf("write tsv data: %w", err)
	}

	tsvWriter.Flush()

	if err := tsvWriter.Error(); err != nil {
		return nil, fmt.Errorf("flush tsv writer: %w", err)
	}

	return []byte(builder.String()), nil
}

func writeTSVData(w *TSVWriter, data any) error {
	switch v := data.(type) {
	case [][]string:
		for _, row := range v {
			if err := w.WriteRow(row); err != nil {
				return fmt.Errorf("write row: %w", err)
			}
		}
	case []string:
		if err := w.WriteRow(v); err != nil {
			return fmt.Errorf("write single row: %w", err)
		}
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedType, data)
	}
	return nil
}
