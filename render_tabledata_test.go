package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func testTableData() *TableData {
	data := NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	return data
}

func TestRenderTableData_CSV(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatCSV, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData csv: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Name") || !strings.Contains(out, "Alice") {
		t.Errorf("csv output missing expected content: %q", out)
	}

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
	if !strings.Contains(out, "Alice") {
		t.Errorf("tsv output missing expected content: %q", out)
	}

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
	if !strings.Contains(out, "# Test") {
		t.Errorf("markdown output missing title: %q", out)
	}

	if !strings.Contains(out, "| Name") {
		t.Errorf("markdown output missing table header: %q", out)
	}
}

func TestRenderTableData_XML(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatXML, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData xml: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "<") {
		t.Errorf("xml output missing expected content: %q", out)
	}
}

func TestRenderTableData_YAML(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatYAML, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData yaml: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") {
		t.Errorf("yaml output missing expected content: %q", out)
	}
}

func TestRenderTableData_D2(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatD2, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData d2: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") {
		t.Errorf("d2 output missing expected content: %q", out)
	}
}

func TestRenderTableData_Mermaid(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatMermaid, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData mermaid: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "graph") && !strings.Contains(out, "flowchart") {
		t.Errorf("mermaid output missing graph declaration: %q", out)
	}
}

func TestRenderTableData_DOT(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatDOT, RenderOptions{Writer: &buf, GraphID: "test_graph"})
	if err != nil {
		t.Fatalf("RenderTableData dot: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "digraph") || !strings.Contains(out, "test_graph") {
		t.Errorf("dot output missing expected content: %q", out)
	}
}

func TestRenderTableData_HTML(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatHTML, RenderOptions{Writer: &buf, Title: "My Data"})
	if err != nil {
		t.Fatalf("RenderTableData html: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "<!DOCTYPE html>") || !strings.Contains(out, "Alice") {
		t.Errorf("html output missing expected content: %q", out)
	}
}

func TestRenderTableData_Tree(t *testing.T) {
	var buf bytes.Buffer

	data := testTableData()

	err := RenderTableData(data, FormatTree, RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData tree: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Alice") {
		t.Errorf("tree output missing expected content: %q", out)
	}
}

func TestRenderTableData_UnsupportedFormats(t *testing.T) {
	data := testTableData()

	for _, f := range []Format{FormatTable, FormatJSON} {
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
	if !strings.Contains(out, "A") {
		t.Errorf("expected header in output, got %q", out)
	}
}

func TestMarshalCSVFromTableData(t *testing.T) {
	data := testTableData()

	b, err := MarshalCSVFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalCSVFromTableData: %v", err)
	}

	if !strings.Contains(string(b), "Alice") {
		t.Errorf("csv missing Alice: %q", string(b))
	}
}

func TestMarshalCSVFromTableData_Nil(t *testing.T) {
	b, err := MarshalCSVFromTableData(nil)
	if err != nil {
		t.Fatalf("MarshalCSVFromTableData nil: %v", err)
	}

	if b != nil {
		t.Errorf("expected nil for nil data, got %q", string(b))
	}
}

func TestMarshalTSVFromTableData(t *testing.T) {
	data := testTableData()

	b, err := MarshalTSVFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalTSVFromTableData: %v", err)
	}

	if !strings.Contains(string(b), "Alice") {
		t.Errorf("tsv missing Alice: %q", string(b))
	}

	if !strings.Contains(string(b), "\t") {
		t.Error("tsv should contain tabs")
	}
}

func TestMarshalTSVFromTableData_Nil(t *testing.T) {
	b, err := MarshalTSVFromTableData(nil)
	if err != nil {
		t.Fatalf("MarshalTSVFromTableData nil: %v", err)
	}

	if b != nil {
		t.Errorf("expected nil for nil data, got %q", string(b))
	}
}
