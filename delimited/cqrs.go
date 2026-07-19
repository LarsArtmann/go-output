package delimited

import (
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteCSV writes a Table as CSV directly to the provided writer using
// NewCSVWriter — true row-level streaming via encoding/csv.
func WriteCSV(w io.Writer, data *output.Table) error {
	return writeDelimited(w, data, "csv", func(w io.Writer) tableDataWriter {
		return NewCSVWriter(w)
	})
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
	return writeDelimited(w, data, "tsv", func(w io.Writer) tableDataWriter {
		return NewTSVWriter(w)
	})
}

// RenderTSV renders a Table as a TSV string.
func RenderTSV(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteTSV(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
