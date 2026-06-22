package integration

import (
	"bytes"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func renderTableDataFixture() *output.TableData {
	data := output.NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	return data
}

func assertRenderContainsBoth(t *testing.T, out, a, b, format string) {
	t.Helper()

	if !strings.Contains(out, a) || !strings.Contains(out, b) {
		t.Errorf("%s output missing expected content: %q", format, out)
	}
}

func TestRenderTableData_CSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()

	err := output.RenderTableData(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData csv: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Name", "Alice", "csv")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
}

func TestRenderTableData_TSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()

	err := output.RenderTableData(data, output.FormatTSV, output.RenderOptions{Writer: &buf})
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
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()

	err := output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf, Title: "Test"})
	if err != nil {
		t.Fatalf("RenderTableData markdown: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "# Test")

	testhelpers.AssertOutputContains(t, out, "| Name")
}

func TestRenderTableData_Tree(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()

	err := output.RenderTableData(data, output.FormatTree, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData tree: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "Alice")
}

func TestRenderTableData_EmptyRows(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := output.NewTableData([]string{"A", "B"})

	err := output.RenderTableData(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData empty rows csv: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "A")
}

func TestRenderTableData_MarkdownWriterError(t *testing.T) {
	t.Parallel()

	testRenderTableDataWriterError(t, output.FormatMarkdown, output.RenderOptions{
		Writer: &testhelpers.ErrorWriter{},
		Title:  "Test",
	}, "expected error from errorWriter")
}

func TestRenderTableData_TreeWriterError(t *testing.T) {
	t.Parallel()

	testRenderTableDataWriterError(t, output.FormatTree, output.RenderOptions{Writer: &testhelpers.ErrorWriter{}},
		"expected error from errorWriter")
}

func TestRenderTableData_MarkdownPartialWriteErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts output.RenderOptions
	}{
		{
			name: "title write error",
			opts: output.RenderOptions{Writer: &testhelpers.WriteNThenFailWriter{Remaining: 0}, Title: "Test"},
		},
		{
			name: "row count write error",
			opts: output.RenderOptions{Writer: &testhelpers.WriteNThenFailWriter{Remaining: 1}, Title: "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := renderTableDataFixture()

			err := output.RenderTableData(data, output.FormatMarkdown, tt.opts)
			if err == nil {
				t.Fatal("expected write error")
			}
		})
	}
}

func testRenderTableDataWriterError(t *testing.T, format output.Format, opts output.RenderOptions, msg string) {
	t.Helper()

	data := renderTableDataFixture()

	err := output.RenderTableData(data, format, opts)
	if err == nil {
		t.Fatal(msg)
	}
}

func TestRenderTableData_MarkdownWithFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData markdown with footer: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Alice", "Total", "markdown with footer")
}

func TestRenderTableData_CSVWithFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableDataFixture()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := output.RenderTableData(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData csv with footer: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Alice", "Total", "csv with footer")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header + 2 rows + footer), got %d", len(lines))
	}
}
