package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
)

func testTableData() *TableData {
	data := NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	return data
}

func assertOutputContainsBoth(t *testing.T, out, a, b, format string) {
	t.Helper()

	if !strings.Contains(out, a) || !strings.Contains(out, b) {
		t.Errorf("%s output missing expected content: %q", format, out)
	}
}

func TestRenderTableData_CSV(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatCSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData csv: %v", err)
	}

	out := buf.String()
	assertOutputContainsBoth(t, out, "Name", "Alice", "csv")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
}

func TestRenderTableData_TSV(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatTSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData tsv: %v", err)
	}

	out := buf.String()
	gentest.AssertOutputContains(t, out, "Alice")

	if !strings.Contains(out, "\t") {
		t.Error("tsv output should contain tabs")
	}
}

func TestRenderTableData_Markdown(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatMarkdown, RenderOptions{Writer: &buf, Title: "Test"})
	if err != nil {
		t.Fatalf("RenderTableData markdown: %v", err)
	}

	out := buf.String()
	gentest.AssertOutputContains(t, out, "# Test")

	gentest.AssertOutputContains(t, out, "| Name")
}

func TestRenderTableData_Tree(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatTree, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData tree: %v", err)
	}

	out := buf.String()
	gentest.AssertOutputContains(t, out, "Alice")
}

func TestRenderTableData_UnsupportedFormats(t *testing.T) {
	data := testTableData()

	for _, f := range []Format{FormatTable, FormatJSON, FormatD2, FormatMermaid, FormatDOT} {
		var buf bytes.Buffer

		err := RenderTableData(data, f, RenderOptions{Writer: &buf})
		if err == nil {
			t.Errorf("expected error for format %q, got nil", f)
		}

		var unsupportedErr *UnsupportedFormatError
		if !errors.As(err, &unsupportedErr) {
			t.Errorf("expected UnsupportedFormatError for %q, got %T: %v", f, err, err)
		}
	}
}

func TestRenderTableData_NilData(t *testing.T) {
	var buf bytes.Buffer

	err := RenderTableData(nil, FormatCSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData with nil data should not error: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
}

func TestRenderTableData_EmptyRows(t *testing.T) {
	var buf bytes.Buffer

	data := NewTableData([]string{"A", "B"})

	err := RenderTableData(data, FormatCSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData empty rows csv: %v", err)
	}

	out := buf.String()
	gentest.AssertOutputContains(t, out, "A")
}

func TestRenderTableData_MarkdownWriterError(t *testing.T) {
	data := testTableData()

	err := RenderTableData(data, FormatMarkdown, RenderOptions{
		Writer: &errorWriter{},
		Title:  "Test",
	})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestRenderTableData_TreeWriterError(t *testing.T) {
	data := testTableData()

	err := RenderTableData(data, FormatTree, RenderOptions{Writer: &errorWriter{}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestRenderTableData_MarkdownTitleWriteError(t *testing.T) {
	data := testTableData()

	err := RenderTableData(data, FormatMarkdown, RenderOptions{
		Writer: &writeNThenFailWriter{Remaining: 0},
		Title:  "Test",
	})
	if err == nil {
		t.Fatal("expected error on title write")
	}
}

func TestRenderTableData_MarkdownRowCountWriteError(t *testing.T) {
	data := testTableData()

	err := RenderTableData(data, FormatMarkdown, RenderOptions{
		Writer: &writeNThenFailWriter{Remaining: 1},
		Title:  "Test",
	})
	if err == nil {
		t.Fatal("expected error on row count write")
	}
}
