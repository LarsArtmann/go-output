package delimited

import (
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteCSV writes a Table as CSV to the provided writer.
func WriteCSV(w io.Writer, data *output.Table) error {
	return renderDelimitedTable(w, data, MarshalCSVFromTable, "csv")
}

// RenderCSV renders a Table as a CSV string.
func RenderCSV(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteCSV(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteTSV writes a Table as TSV to the provided writer.
func WriteTSV(w io.Writer, data *output.Table) error {
	return renderDelimitedTable(w, data, MarshalTSVFromTable, "tsv")
}

// RenderTSV renders a Table as a TSV string.
func RenderTSV(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteTSV(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
