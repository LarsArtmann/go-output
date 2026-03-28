package output

import (
	"strings"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"table", "table", OutputFormatTable, false},
		{"json", "json", OutputFormatJSON, false},
		{"csv", "csv", OutputFormatCSV, false},
		{"markdown", "markdown", OutputFormatMarkdown, false},
		{"d2", "d2", OutputFormatD2, false},
		{"yaml", "yaml", OutputFormatYAML, false},
		{"html", "html", FormatHTML, false},
		{"tree", "tree", FormatTree, false},
		{"mermaid", "mermaid", FormatMermaid, false},
		{"dot", "dot", FormatDOT, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format OutputFormat
		want   string
	}{
		{OutputFormatTable, "table"},
		{OutputFormatJSON, "json"},
		{OutputFormatCSV, "csv"},
		{OutputFormatMarkdown, "markdown"},
		{OutputFormatD2, "d2"},
		{OutputFormatYAML, "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.format.String(); got != tt.want {
				t.Errorf("OutputFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatAllowedValues(t *testing.T) {
	t.Parallel()
	got := OutputFormatTable.AllowedValues()
	want := []string{
		"table",
		"json",
		"csv",
		"tsv",
		"markdown",
		"d2",
		"yaml",
		"html",
		"tree",
		"mermaid",
		"dot",
	}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestOutputFormatIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format OutputFormat
		want   bool
	}{
		{OutputFormatTable, true},
		{OutputFormatJSON, true},
		{OutputFormatCSV, true},
		{OutputFormatMarkdown, true},
		{OutputFormatD2, true},
		{OutputFormatYAML, true},
		{FormatHTML, true},
		{FormatTree, true},
		{FormatMermaid, true},
		{FormatDOT, true},
		{OutputFormat("invalid"), false},
		{OutputFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			if got := tt.format.IsValid(); got != tt.want {
				t.Errorf("OutputFormat(%q).IsValid() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func FuzzParseOutputFormat(f *testing.F) {
	f.Add("table")
	f.Add("json")
	f.Add("csv")
	f.Add("markdown")
	f.Add("d2")
	f.Add("yaml")
	f.Add("html")
	f.Add("tree")
	f.Add("mermaid")
	f.Add("dot")
	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		format, err := ParseOutputFormat(s)
		if err != nil {
			if format != "" {
				t.Errorf("ParseOutputFormat(%q) returned error but non-empty format: %q", s, format)
			}
		}
		if format.IsValid() && err == nil {
			if string(format) != s {
				t.Errorf("ParseOutputFormat(%q) = %q, but IsValid() was true", s, format)
			}
		}
	})
}

func TestFormatIsTableFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format Format
		want   bool
	}{
		{FormatTable, true},
		{FormatJSON, true},
		{FormatCSV, true},
		{FormatMarkdown, true},
		{FormatD2, true},
		{FormatYAML, true},
		{FormatHTML, false},
		{FormatTree, false},
		{FormatMermaid, false},
		{FormatDOT, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			if got := tt.format.IsTableFormat(); got != tt.want {
				t.Errorf("Format(%q).IsTableFormat() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestFormatIsTreeFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format Format
		want   bool
	}{
		{FormatTable, false},
		{FormatJSON, false},
		{FormatTree, true},
		{FormatHTML, true},
		{FormatMermaid, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			if got := tt.format.IsTreeFormat(); got != tt.want {
				t.Errorf("Format(%q).IsTreeFormat() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestFormatIsGraphFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		format Format
		want   bool
	}{
		{FormatTable, false},
		{FormatJSON, false},
		{FormatD2, true},
		{FormatMermaid, true},
		{FormatDOT, true},
		{FormatTree, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()
			if got := tt.format.IsGraphFormat(); got != tt.want {
				t.Errorf("Format(%q).IsGraphFormat() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestInvalidFormatError(t *testing.T) {
	t.Parallel()
	err := &InvalidFormatError{
		Value:   "invalid",
		Allowed: []Format{FormatTable, FormatJSON},
	}

	got := err.Error()
	wantContains := []string{"invalid format", "invalid", "table", "json"}
	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("Error() = %q, should contain %q", got, want)
		}
	}
}

func TestTableData(t *testing.T) {
	t.Parallel()
	t.Run("RowCount and ColCount", func(t *testing.T) {
		t.Parallel()
		testTableDataRowColCount(t)
	})
	t.Run("CreateRowEdges", func(t *testing.T) {
		t.Parallel()
		testTableDataCreateRowEdges(t)
	})
}

func testTableDataRowColCount(t *testing.T) {
	t.Helper()
	data := NewTableData([]string{"Name", "Value", "Count"})
	if data.ColCount() != 3 {
		t.Errorf("ColCount() = %d, want 3", data.ColCount())
	}
	if data.RowCount() != 0 {
		t.Errorf("RowCount() = %d, want 0", data.RowCount())
	}

	data.AddRow([]string{"a", "b", "c"})
	data.AddRow([]string{"d", "e", "f"})
	if data.RowCount() != 2 {
		t.Errorf("RowCount() = %d, want 2", data.RowCount())
	}
}

func testTableDataCreateRowEdges(t *testing.T) {
	t.Helper()
	t.Run("nil data", testCreateRowEdgesNil)
	t.Run("empty rows", testCreateRowEdgesEmpty)
	t.Run("single row", testCreateRowEdgesSingle)
	t.Run("multiple rows", testCreateRowEdgesMultiple)
}

func testCreateRowEdgesNil(t *testing.T) {
	var data *TableData
	if edges := data.CreateRowEdges(); edges != nil {
		t.Errorf("CreateRowEdges() on nil = %v, want nil", edges)
	}
}

func testCreateRowEdgesEmpty(t *testing.T) {
	data := NewTableData([]string{"Name"})
	if edges := data.CreateRowEdges(); edges != nil {
		t.Errorf("CreateRowEdges() on empty = %v, want nil", edges)
	}
}

func testCreateRowEdgesSingle(t *testing.T) {
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"a"})
	if edges := data.CreateRowEdges(); edges != nil {
		t.Errorf("CreateRowEdges() on single row = %v, want nil", edges)
	}
}

func testCreateRowEdgesMultiple(t *testing.T) {
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"a"})
	data.AddRow([]string{"b"})
	data.AddRow([]string{"c"})
	edges := data.CreateRowEdges()
	if len(edges) != 2 {
		t.Fatalf("CreateRowEdges() returned %d edges, want 2", len(edges))
	}
	verifyEdge := func(idx int, from, to string) {
		if edges[idx].From != from || edges[idx].To != to {
			t.Errorf(
				"Edge %d = {%s, %s}, want {%s, %s}",
				idx,
				edges[idx].From,
				edges[idx].To,
				from,
				to,
			)
		}
	}
	verifyEdge(0, "row0", "row1")
	verifyEdge(1, "row1", "row2")
}
