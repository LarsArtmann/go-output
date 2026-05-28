package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
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
	testhelpers.AssertOutputContains(t, out, "Alice")

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
	testhelpers.AssertOutputContains(t, out, "# Test")

	testhelpers.AssertOutputContains(t, out, "| Name")
}

func TestRenderTableData_Tree(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatTree, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData tree: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "Alice")
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
	testhelpers.AssertOutputContains(t, out, "A")
}

func TestRenderTableData_MarkdownWriterError(t *testing.T) {
	testRenderTableDataWriterError(t, FormatMarkdown, RenderOptions{
		Writer: &errorWriter{},
		Title:  "Test",
	}, "expected error from errorWriter")
}

func TestRenderTableData_TreeWriterError(t *testing.T) {
	testRenderTableDataWriterError(t, FormatTree, RenderOptions{Writer: &errorWriter{}},
		"expected error from errorWriter")
}

func TestRenderTableData_MarkdownPartialWriteErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts RenderOptions
	}{
		{
			name: "title write error",
			opts: RenderOptions{Writer: &writeNThenFailWriter{Remaining: 0}, Title: "Test"},
		},
		{
			name: "row count write error",
			opts: RenderOptions{Writer: &writeNThenFailWriter{Remaining: 1}, Title: "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := testTableData()

			err := RenderTableData(data, FormatMarkdown, tt.opts)
			if err == nil {
				t.Fatal("expected write error")
			}
		})
	}
}

func testRenderTableDataWriterError(t *testing.T, format Format, opts RenderOptions, msg string) {
	t.Helper()

	data := testTableData()

	err := RenderTableData(data, format, opts)
	if err == nil {
		t.Fatal(msg)
	}
}

func TestRenderTableData_MarkdownWithFooter(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := RenderTableData(data, FormatMarkdown, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData markdown with footer: %v", err)
	}

	out := buf.String()
	assertOutputContainsBoth(t, out, "Alice", "Total", "markdown with footer")
}

func TestRenderTableData_CSVWithFooter(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := RenderTableData(data, FormatCSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData csv with footer: %v", err)
	}

	out := buf.String()
	assertOutputContainsBoth(t, out, "Alice", "Total", "csv with footer")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header + 2 rows + footer), got %d", len(lines))
	}
}

func TestRenderTableData_ValidateRejectsFooterMismatch(t *testing.T) {
	var buf bytes.Buffer

	data := NewTableData([]string{"A", "B"})
	data.Footer = []string{"too", "many", "cols"}

	err := RenderTableData(data, FormatMarkdown, RenderOptions{Writer: &buf})
	if err == nil {
		t.Fatal("RenderTableData should reject footer with wrong column count")
	}

	if !strings.Contains(err.Error(), "footer column count") {
		t.Errorf("error should mention footer column count, got: %v", err)
	}
}
