package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markdown"
	"github.com/larsartmann/go-output/serialization"
)

const (
	headerName   = "Name"
	headerHealth = "Health"
	sampleHealth = "90%"
	sampleAlpha  = "Alpha"
)

// assertTable verifies that the table has the expected number of columns and rows.
func assertTable(t *testing.T, data *output.Table, expectedCols, expectedRows int) {
	t.Helper()

	if data == nil {
		t.Fatal("Table is nil")
	}

	if got := data.ColCount(); got != expectedCols {
		t.Errorf("Table has %d columns, want %d", got, expectedCols)
	}

	if got := data.RowCount(); got != expectedRows {
		t.Errorf("Table has %d rows, want %d", got, expectedRows)
	}
}

// renderMarkdownTable renders a markdown table with the given headers and rows,
// failing the test if rendering errors.
func renderMarkdownTable(t *testing.T, headers []string, rows [][]string) string {
	t.Helper()

	md := markdown.NewMarkdownTable()
	md.SetHeaders(headers)

	for _, row := range rows {
		md.AddRow(row)
	}

	return mustRenderMarkdown(t, md)
}

// renderSampleMarkdownTable returns a rendered markdown table with sample project data.
func renderSampleMarkdownTable(t *testing.T) string {
	t.Helper()

	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{headerName, headerHealth})
	md.AddRow([]string{sampleAlpha, sampleHealth})

	return mustRenderMarkdown(t, md)
}

// mustRenderMarkdown renders a MarkdownTable, failing the test on error so
// content assertions report the real cause instead of "missing X".
func mustRenderMarkdown(t *testing.T, md *markdown.MarkdownTable) string {
	t.Helper()

	out, err := md.Render()
	if err != nil {
		t.Fatalf("markdown render failed: %v", err)
	}

	return out
}

// runEmptyDataRendersJSONWithoutPanic runs the "empty data renders without panic" test sub-case.
func runEmptyDataRendersJSONWithoutPanic(t *testing.T) {
	t.Run("empty data renders without error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{})

		if _, err := serialization.MarshalJSON(data); err != nil {
			t.Errorf("MarshalJSON on empty table failed: %v", err)
		}
	})
}
