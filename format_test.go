package output

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []parseEnumTestCase[Format]{
		{"table", "table", FormatTable, false},
		{"json", "json", FormatJSON, false},
		{"csv", "csv", FormatCSV, false},
		{"markdown", "markdown", FormatMarkdown, false},
		{"d2", "d2", FormatD2, false},
		{"yaml", "yaml", FormatYAML, false},
		{"html", "html", FormatHTML, false},
		{"tree", "tree", FormatTree, false},
		{"mermaid", "mermaid", FormatMermaid, false},
		{"dot", "dot", FormatDOT, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}
	testParseEnum(
		t,
		"ParseOutputFormat",
		ParseOutputFormat,
		tests,
		func(a, b Format) bool { return a == b },
	)
}

func TestOutputFormatString(t *testing.T) {
	tests := []stringEnumTestCase[Format]{
		{FormatTable, "table"},
		{FormatJSON, "json"},
		{FormatCSV, "csv"},
		{FormatMarkdown, "markdown"},
		{FormatD2, "d2"},
		{FormatYAML, "yaml"},
	}
	testEnumString(t, "OutputFormat.String", tests, func(f Format) string { return f.String() })
}

func TestOutputFormatAllowedValues(t *testing.T) {
	testAllowedValues(
		t,
		"AllowedValues",
		OutputFormatTable.AllowedValues(),
		[]string{
			"table",
			"json",
			"csv",
			"tsv",
			"markdown",
			"xml",
			"d2",
			"yaml",
			"html",
			"tree",
			"mermaid",
			"dot",
		},
	)
}

func TestOutputFormatIsValid(t *testing.T) {
	t.Parallel()

	gentest.TestEnumIsValid[OutputFormat](t, []OutputFormat{
		OutputFormatTable,
		OutputFormatJSON,
		OutputFormatCSV,
		OutputFormatMarkdown,
		OutputFormatD2,
		OutputFormatYAML,
		FormatHTML,
		FormatTree,
		FormatMermaid,
		FormatDOT,
		OutputFormat("invalid"),
		OutputFormat(""),
	}, []bool{
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
	})
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
		fuzzEnumTest(t, s, ParseOutputFormat, "ParseOutputFormat")
	})
}

func TestFormatIsTableFormat(t *testing.T) {
	testBoolMethod(
		t,
		"Format",
		"IsTableFormat",
		[]boolMethodTestCase[Format]{
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
		},
		func(f Format) bool { return f.IsTableFormat() },
		func(f Format) string { return string(f) },
	)
}

func TestFormatIsTreeFormat(t *testing.T) {
	testBoolMethod(
		t,
		"Format",
		"IsTreeFormat",
		[]boolMethodTestCase[Format]{
			{FormatTable, false},
			{FormatJSON, false},
			{FormatTree, true},
			{FormatHTML, true},
			{FormatMermaid, false},
		},
		func(f Format) bool { return f.IsTreeFormat() },
		func(f Format) string { return string(f) },
	)
}

func TestFormatIsGraphFormat(t *testing.T) {
	testBoolMethod(
		t,
		"Format",
		"IsGraphFormat",
		[]boolMethodTestCase[Format]{
			{FormatTable, false},
			{FormatJSON, false},
			{FormatD2, true},
			{FormatMermaid, true},
			{FormatDOT, true},
			{FormatTree, false},
		},
		func(f Format) bool { return f.IsGraphFormat() },
		func(f Format) string { return string(f) },
	)
}

func TestInvalidFormatError(t *testing.T) {
	t.Parallel()

	err := &InvalidFormatError{
		Value:   "invalid",
		Allowed: nil,
	}

	got := err.Error()

	wantContains := []string{"invalid format", "invalid"}
	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("Error() = %q, should contain %q", got, want)
		}
	}
}

func TestFormatCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format    Format
		wantTable bool
		wantTree  bool
		wantGraph bool
	}{
		{FormatTable, true, false, false},
		{FormatJSON, true, false, false},
		{FormatCSV, true, false, false},
		{FormatTSV, true, false, false},
		{FormatMarkdown, true, false, false},
		{FormatYAML, true, false, false},
		{FormatD2, true, false, true}, // D2 is both table and graph
		{FormatHTML, false, true, false},
		{FormatTree, false, true, false},
		{FormatMermaid, false, false, true},
		{FormatDOT, false, false, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			testBoolValue(
				t,
				string(tt.format),
				"IsTableFormat",
				tt.format.IsTableFormat(),
				tt.wantTable,
			)
			testBoolValue(
				t,
				string(tt.format),
				"IsTreeFormat",
				tt.format.IsTreeFormat(),
				tt.wantTree,
			)
			testBoolValue(
				t,
				string(tt.format),
				"IsGraphFormat",
				tt.format.IsGraphFormat(),
				tt.wantGraph,
			)
		})
	}
}

func TestTableData(t *testing.T) {
	t.Parallel()
	runSubtest(t, "RowCount and ColCount", testTableDataRowColCount)
	runSubtest(t, "CreateRowEdges", testTableDataCreateRowEdges)
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
