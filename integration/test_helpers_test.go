package integration

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestAssertTableData(t *testing.T) {
	t.Parallel()

	t.Run("valid data", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"A", "B"})
		data.AddRow([]string{"1", "2"})

		assertTableData(t, data, 2, 1)
	})

	mismatchTests := []struct {
		name    string
		headers int
		rows    int
	}{
		{name: "mismatched columns", headers: 5, rows: 1},
		{name: "mismatched rows", headers: 1, rows: 99},
	}

	for _, tt := range mismatchTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := &testing.T{}

			data := output.NewTableData([]string{"A"})
			data.AddRow([]string{"1"})

			assertTableData(mock, data, tt.headers, tt.rows)

			if !mock.Failed() {
				t.Errorf("expected failure for %s", tt.name)
			}
		})
	}
}

func TestRenderMarkdownTable(t *testing.T) {
	t.Parallel()

	got := renderMarkdownTable([]string{"X", "Y"}, [][]string{{"1", "2"}})

	testhelpers.AssertAllContained(t, got, "X", "1")
}

func TestRenderSampleMarkdownTable(t *testing.T) {
	t.Parallel()

	got := renderSampleMarkdownTable()

	if !strings.Contains(got, sampleAlpha) || !strings.Contains(got, sampleHealth) {
		t.Errorf("renderSampleMarkdownTable() = %q, expected sample data", got)
	}
}

func TestRunEmptyDataRendersJSONWithoutPanic(t *testing.T) {
	t.Parallel()

	runEmptyDataRendersJSONWithoutPanic(t)
}
