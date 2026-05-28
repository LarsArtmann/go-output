package integration

import (
	"bytes"
	"errors"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestRenderTableData_NilData(t *testing.T) {
	t.Parallel()

	err := output.RenderTableData(nil, output.FormatMarkdown)
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

	if !bytes.Contains([]byte(err.Error()), []byte("footer column count")) {
		t.Errorf("Expected footer column mismatch error, got: %v", err)
	}
}

func TestRenderTableData_UnsupportedFormat(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"A"})
	data.AddRow([]string{"1"})

	var unsupportedFormats = []output.Format{
		output.FormatD2,
		output.FormatDOT,
		output.FormatMermaid,
		output.FormatTable,
		output.FormatJSON,
	}

	for _, format := range unsupportedFormats {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			err := output.RenderTableData(data, format)
			if err == nil {
				t.Fatalf("Expected UnsupportedFormatError for %s", format)
			}

			var unsupportedErr *output.UnsupportedFormatError
			if !errors.As(err, &unsupportedErr) {
				t.Errorf("Expected UnsupportedFormatError, got: %T: %v", err, err)
			}
		})
	}
}

func TestMustRender_PanicOnFailure(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected MustRender to panic on failing renderer")
		}
	}()

	output.MustRender(&testhelpers.ErrorRenderer{})
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
