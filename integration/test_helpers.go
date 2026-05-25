package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/serialization"
)

const (
	headerName   = "Name"
	headerHealth = "Health"
	sampleHealth = "90%"
	sampleAlpha  = "Alpha"
)

// assertTableData verifies that the table has the expected number of columns and rows.
func assertTableData(t *testing.T, data *output.TableData, expectedCols, expectedRows int) {
	t.Helper()

	if data == nil {
		t.Fatal("TableData is nil")

		return
	}

	if got := data.ColCount(); got != expectedCols {
		t.Errorf("TableData has %d columns, want %d", got, expectedCols)
	}

	if got := data.RowCount(); got != expectedRows {
		t.Errorf("TableData has %d rows, want %d", got, expectedRows)
	}
}

// renderMarkdownTable renders a markdown table with the given headers and rows.
func renderMarkdownTable(headers []string, rows [][]string) string {
	md := output.NewMarkdownTable()
	md.SetHeaders(headers)

	for _, row := range rows {
		md.AddRow(row)
	}

	out, err := md.Render()
	if err != nil {
		return ""
	}

	return out
}

// renderSampleMarkdownTable returns a rendered markdown table with sample project data.
func renderSampleMarkdownTable() string {
	md := output.NewMarkdownTable()
	md.SetHeaders([]string{headerName, headerHealth})
	md.AddRow([]string{sampleAlpha, sampleHealth})

	out, err := md.Render()
	if err != nil {
		return ""
	}

	return out
}

// runEmptyDataRendersJSONWithoutPanic runs the "empty data renders without panic" test sub-case.
func runEmptyDataRendersJSONWithoutPanic(t *testing.T) {
	t.Run("empty data renders without panic", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{})

		_, err := serialization.MarshalJSON(data)
		if err != nil {
			t.Errorf("MarshalJSON on empty data should not error: %v", err)
		}
	})
}
