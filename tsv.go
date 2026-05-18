package output

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// TSVWriter writes TSV (Tab-Separated Values) output.
type TSVWriter struct {
	writer *DelimitedWriter
}

// ErrUnsupportedType is returned when an unsupported type is provided for TSV marshaling.
var ErrUnsupportedType = errors.New("unsupported type")

// NewTSVWriter creates a new TSVWriter.
func NewTSVWriter(w io.Writer) *TSVWriter {
	return &TSVWriter{
		writer: NewDelimitedWriter(w, '\t', "tsv"),
	}
}

// WriteHeader writes the header row.
func (t *TSVWriter) WriteHeader(cols []string) error {
	return t.writer.WriteRow(cols, "tsv header")
}

// WriteRow writes a single row.
func (t *TSVWriter) WriteRow(values []string) error {
	return t.writer.WriteRow(values, "tsv row")
}

// WriteRows writes multiple rows.
func (t *TSVWriter) WriteRows(values [][]string) error {
	return t.writer.WriteRows(values, "tsv")
}

// Flush flushes the writer.
func (t *TSVWriter) Flush() {
	t.writer.Flush()
}

// Error returns any error from the writer.
func (t *TSVWriter) Error() error {
	return t.writer.Error()
}

// MarshalTSV marshals data as TSV.
func MarshalTSV(data any) ([]byte, error) {
	var builder strings.Builder

	tsvWriter := NewTSVWriter(&builder)

	err := writeTSVData(tsvWriter, data)
	if err != nil {
		return nil, fmt.Errorf("write tsv data: %w", err)
	}

	tsvWriter.Flush()

	err = tsvWriter.Error()
	if err != nil {
		return nil, fmt.Errorf("flush tsv writer: %w", err)
	}

	return []byte(builder.String()), nil
}

func writeTSVData(w *TSVWriter, data any) error {
	switch v := data.(type) {
	case [][]string:
		for _, row := range v {
			err := w.WriteRow(row)
			if err != nil {
				return fmt.Errorf("write row: %w", err)
			}
		}
	case []string:
		err := w.WriteRow(v)
		if err != nil {
			return fmt.Errorf("write single row: %w", err)
		}
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedType, data)
	}

	return nil
}

// MarshalTSVFromTableData marshals TableData as TSV with a header row.
func MarshalTSVFromTableData(data *TableData) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	var builder strings.Builder

	tsvWriter := NewTSVWriter(&builder)

	if len(data.Headers) > 0 {
		if err := tsvWriter.WriteHeader(data.Headers); err != nil {
			return nil, fmt.Errorf("write tsv header: %w", err)
		}
	}

	for _, row := range data.Rows {
		if err := tsvWriter.WriteRow(row); err != nil {
			return nil, fmt.Errorf("write tsv row: %w", err)
		}
	}

	tsvWriter.Flush()

	if err := tsvWriter.Error(); err != nil {
		return nil, fmt.Errorf("flush tsv writer: %w", err)
	}

	return []byte(builder.String()), nil
}
