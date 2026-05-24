package integration

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestAssertTableData(t *testing.T) {
	t.Parallel()

	t.Run("valid data", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"A", "B"})
		data.AddRow([]string{"1", "2"})

		assertTableData(t, data, 2, 1)
	})

	t.Run("mismatched columns", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		data := output.NewTableData([]string{"A"})
		data.AddRow([]string{"1"})

		assertTableData(mock, data, 5, 1)

		if !mock.Failed() {
			t.Error("expected failure for column count mismatch")
		}
	})

	t.Run("mismatched rows", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		data := output.NewTableData([]string{"A"})
		data.AddRow([]string{"1"})

		assertTableData(mock, data, 1, 99)

		if !mock.Failed() {
			t.Error("expected failure for row count mismatch")
		}
	})
}

func TestRenderMarkdownTable(t *testing.T) {
	t.Parallel()

	got := renderMarkdownTable([]string{"X", "Y"}, [][]string{{"1", "2"}})

	if !strings.Contains(got, "X") || !strings.Contains(got, "1") {
		t.Errorf("renderMarkdownTable() = %q, expected table content", got)
	}
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
