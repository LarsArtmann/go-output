package delimited

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteCSV writes a Table as CSV directly to the provided writer using
// NewCSVWriter — true row-level streaming via encoding/csv.
func WriteCSV(w io.Writer, data *output.Table) error {
	if data == nil {
		return nil
	}

	cw := NewCSVWriter(w)

	if len(data.Headers) > 0 {
		if err := cw.WriteHeader(data.Headers); err != nil {
			return fmt.Errorf("write csv header: %w", err)
		}
	}

	for _, row := range data.Rows {
		if err := cw.WriteRow(row); err != nil {
			return fmt.Errorf("write csv row: %w", err)
		}
	}

	if data.HasFooter() {
		if err := cw.WriteFooter(data.Footer); err != nil {
			return fmt.Errorf("write csv footer: %w", err)
		}
	}

	cw.Flush()

	if err := cw.Error(); err != nil {
		return fmt.Errorf("flush csv writer: %w", err)
	}

	return nil
}

// RenderCSV renders a Table as a CSV string.
func RenderCSV(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteCSV(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteTSV writes a Table as TSV directly to the provided writer using
// NewTSVWriter — true row-level streaming via encoding/csv.
func WriteTSV(w io.Writer, data *output.Table) error {
	if data == nil {
		return nil
	}

	tw := NewTSVWriter(w)

	if len(data.Headers) > 0 {
		if err := tw.WriteHeader(data.Headers); err != nil {
			return fmt.Errorf("write tsv header: %w", err)
		}
	}

	for _, row := range data.Rows {
		if err := tw.WriteRow(row); err != nil {
			return fmt.Errorf("write tsv row: %w", err)
		}
	}

	if data.HasFooter() {
		if err := tw.WriteFooter(data.Footer); err != nil {
			return fmt.Errorf("write tsv footer: %w", err)
		}
	}

	tw.Flush()

	if err := tw.Error(); err != nil {
		return fmt.Errorf("flush tsv writer: %w", err)
	}

	return nil
}

// RenderTSV renders a Table as a TSV string.
func RenderTSV(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteTSV(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
