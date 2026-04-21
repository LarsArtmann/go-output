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
	// Generate expected list from AllFormats to avoid hardcoding
	want := make([]string, len(AllFormats))
	for i, f := range AllFormats {
		want[i] = string(f)
	}

	testAllowedValues(
		t,
		"AllowedValues",
		OutputFormatTable.AllowedValues(),
		want,
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

func TestFormatCategoryMethod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format Format
		want   FormatCategory
	}{
		{FormatTable, CategoryTable},
		{FormatJSON, CategoryTable},
		{FormatCSV, CategoryTable},
		{FormatYAML, CategoryTable},
		{FormatHTML, CategoryTree},
		{FormatTree, CategoryTree},
		{FormatD2, CategoryGraph},
		{FormatMermaid, CategoryGraph},
		{FormatDOT, CategoryGraph},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			got := tt.format.Category()
			if got != tt.want {
				t.Errorf("Format(%q).Category() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}
