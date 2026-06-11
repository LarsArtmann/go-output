package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestRenderTableData_NilData(t *testing.T) {
	t.Parallel()

	err := output.RenderTableData(nil, output.FormatMarkdown, output.RenderOptions{})
	if err != nil {
		t.Errorf("RenderTableData with nil data should return nil, got: %v", err)
	}
}

func TestRenderTableData_InvalidFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"A", "B"})
	data.AddRow([]string{"1", "2"})
	data.SetFooter([]string{"total"})

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf})
	if err == nil {
		t.Fatal("Expected error for footer with wrong column count")
	}

	if !strings.Contains(err.Error(), "footer column count") {
		t.Errorf("Expected footer column mismatch error, got: %v", err)
	}
}

func TestRenderTableData_PreviouslyUnsupportedFormats(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"A"})
	data.AddRow([]string{"1"})

	formats := []output.Format{
		output.FormatD2,
		output.FormatDOT,
		output.FormatMermaid,
		output.FormatTable,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := output.RenderTableData(data, format, output.RenderOptions{Writer: &buf})
			if err != nil {
				t.Fatalf("Expected %s to be supported, got: %v", format, err)
			}

			if buf.Len() == 0 {
				t.Errorf("Expected non-empty output for %s", format)
			}
		})
	}
}

func TestCreateRowEdges_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		var data *output.TableData

		edges := data.CreateRowEdges()
		if edges != nil {
			t.Errorf("Expected nil edges for nil data, got %v", edges)
		}
	})

	t.Run("zero rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"A"})

		edges := data.CreateRowEdges()
		if edges != nil {
			t.Errorf("Expected nil edges for zero rows, got %v", edges)
		}
	})

	t.Run("single row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"A"})
		data.AddRow([]string{"1"})

		edges := data.CreateRowEdges()
		if edges != nil {
			t.Errorf("Expected nil edges for single row, got %v", edges)
		}
	})
}

func TestRenderTableData_WithWriter(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alpha"})

	t.Run("nil writer defaults to stdout", func(t *testing.T) {
		t.Parallel()

		err := output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{})
		if err != nil {
			t.Fatalf("RenderTableData with nil writer: %v", err)
		}
	})

	t.Run("writer error", func(t *testing.T) {
		t.Parallel()

		err := output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{
			Writer: &testhelpers.ErrorWriter{},
		})
		if err == nil {
			t.Error("Expected error when writer fails")
		}
	})
}
