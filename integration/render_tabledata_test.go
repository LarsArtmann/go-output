package integration

import (
	"bytes"
	"strings"
	"testing"

	output "github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func renderTableFixture() *output.Table {
	data := output.NewTable([]string{"Name", "Age", "City"})
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

func TestRenderTable_CSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()

	err := output.RenderTable(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable csv: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Name", "Alice", "csv")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
}

func TestRenderTable_TSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()

	err := output.RenderTable(data, output.FormatTSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable tsv: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "Alice")

	if !strings.Contains(out, "\t") {
		t.Error("tsv output should contain tabs")
	}
}

func TestRenderTable_Markdown(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()

	err := output.RenderTable(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf, Title: "Test"})
	if err != nil {
		t.Fatalf("RenderTable markdown: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "# Test")

	testhelpers.AssertOutputContains(t, out, "| Name")
}

func TestRenderTable_Tree(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()

	err := output.RenderTable(data, output.FormatTree, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable tree: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "Alice")
}

func TestRenderTable_EmptyRows(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := output.NewTable([]string{"A", "B"})

	err := output.RenderTable(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable empty rows csv: %v", err)
	}

	out := buf.String()
	testhelpers.AssertOutputContains(t, out, "A")
}

func TestRenderTable_MarkdownWriterError(t *testing.T) {
	t.Parallel()

	testRenderTableWriterError(t, output.FormatMarkdown, output.RenderOptions{
		Writer: &testhelpers.ErrorWriter{},
		Title:  "Test",
	}, "expected error from errorWriter")
}

func TestRenderTable_TreeWriterError(t *testing.T) {
	t.Parallel()

	testRenderTableWriterError(t, output.FormatTree, output.RenderOptions{Writer: &testhelpers.ErrorWriter{}},
		"expected error from errorWriter")
}

func TestRenderTable_MarkdownPartialWriteErrors(t *testing.T) {
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

			data := renderTableFixture()

			err := output.RenderTable(data, output.FormatMarkdown, tt.opts)
			if err == nil {
				t.Fatal("expected write error")
			}
		})
	}
}

func testRenderTableWriterError(t *testing.T, format output.Format, opts output.RenderOptions, msg string) {
	t.Helper()

	data := renderTableFixture()

	err := output.RenderTable(data, format, opts)
	if err == nil {
		t.Fatal(msg)
	}
}

func TestRenderTable_MarkdownWithFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := output.RenderTable(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable markdown with footer: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Alice", "Total", "markdown with footer")
}

func TestRenderTable_CSVWithFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	data := renderTableFixture()
	data.Footer = []string{"Total", "55", "NYC+LA"}

	err := output.RenderTable(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable csv with footer: %v", err)
	}

	out := buf.String()
	assertRenderContainsBoth(t, out, "Alice", "Total", "csv with footer")

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header + 2 rows + footer), got %d", len(lines))
	}
}
