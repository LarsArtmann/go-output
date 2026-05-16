package output

import (
	"slices"
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

	gentest.TestEnumIsValid(t, []OutputFormat{
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
			{FormatTSV, true},
			{FormatMarkdown, true},
			{FormatXML, true},
			{FormatD2, true},
			{FormatYAML, true},
			{FormatHTML, true},
			{FormatMermaid, true},
			{FormatDOT, true},
			{FormatTree, false},
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
			{FormatTree, true},
			{FormatHTML, true},
			{FormatJSON, true},
			{FormatYAML, true},
			{FormatTable, false},
			{FormatCSV, false},
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
			{FormatD2, true},
			{FormatMermaid, true},
			{FormatDOT, true},
			{FormatJSON, true},
			{FormatYAML, true},
			{FormatTable, false},
			{FormatTree, false},
			{FormatCSV, false},
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
		{FormatJSON, true, true, true},
		{FormatCSV, true, false, false},
		{FormatTSV, true, false, false},
		{FormatMarkdown, true, false, false},
		{FormatXML, true, false, false},
		{FormatYAML, true, true, true},
		{FormatD2, true, false, true},
		{FormatHTML, true, true, false},
		{FormatTree, false, true, false},
		{FormatMermaid, true, false, true},
		{FormatDOT, true, false, true},
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

	// Category() returns graph > tree > table priority for multi-shape formats
	tests := []struct {
		format Format
		want   FormatCategory
	}{
		{FormatTable, CategoryTable},
		{FormatJSON, CategoryGraph},
		{FormatCSV, CategoryTable},
		{FormatTSV, CategoryTable},
		{FormatXML, CategoryTable},
		{FormatMarkdown, CategoryTable},
		{FormatYAML, CategoryGraph},
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

func TestFormatSupports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format Format
		shape  Shape
		want   bool
	}{
		{FormatJSON, ShapeTable, true},
		{FormatJSON, ShapeTree, true},
		{FormatJSON, ShapeGraph, true},
		{FormatCSV, ShapeTable, true},
		{FormatCSV, ShapeTree, false},
		{FormatCSV, ShapeGraph, false},
		{FormatD2, ShapeTable, true},
		{FormatD2, ShapeTree, false},
		{FormatD2, ShapeGraph, true},
		{FormatTree, ShapeTree, true},
		{FormatTree, ShapeTable, false},
		{FormatTable, ShapeTable, true},
		{FormatTable, ShapeTree, false},
		{Format("invalid"), ShapeTable, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.format)+"_"+string(tt.shape), func(t *testing.T) {
			t.Parallel()

			got := tt.format.Supports(tt.shape)
			if got != tt.want {
				t.Errorf("Format(%q).Supports(%v) = %v, want %v", tt.format, tt.shape, got, tt.want)
			}
		})
	}
}

func TestFormatShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format Format
		want   []Shape
	}{
		{FormatJSON, []Shape{ShapeTable, ShapeTree, ShapeGraph}},
		{FormatCSV, []Shape{ShapeTable}},
		{FormatD2, []Shape{ShapeTable, ShapeGraph}},
		{FormatTree, []Shape{ShapeTree}},
		{FormatTable, []Shape{ShapeTable}},
		{FormatYAML, []Shape{ShapeTable, ShapeTree, ShapeGraph}},
	}
	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			got := tt.format.Shapes()
			if len(got) != len(tt.want) {
				t.Errorf("Format(%q).Shapes() = %v, want %v", tt.format, got, tt.want)

				return
			}

			for i, s := range got {
				if s != tt.want[i] {
					t.Errorf("Format(%q).Shapes()[%d] = %v, want %v", tt.format, i, s, tt.want[i])
				}
			}
		})
	}
}

func TestFormatsForShape(t *testing.T) {
	t.Parallel()

	assertContainsAll(t, "graph", FormatsForShape(ShapeGraph),
		FormatJSON, FormatYAML, FormatD2, FormatMermaid, FormatDOT)
	assertContainsAll(t, "tree", FormatsForShape(ShapeTree),
		FormatJSON, FormatYAML, FormatHTML, FormatTree)
	assertContainsAll(t, "table", FormatsForShape(ShapeTable),
		FormatTable, FormatJSON, FormatCSV, FormatD2, FormatMermaid, FormatDOT)
}

func assertContainsAll(t *testing.T, label string, formats []Format, required ...Format) {
	t.Helper()

	for _, f := range required {
		if !slices.Contains(formats, f) {
			t.Errorf("FormatsForShape(%s) missing %v", label, f)
		}
	}
}

func TestAllShapes(t *testing.T) {
	t.Parallel()

	want := []Shape{ShapeTable, ShapeTree, ShapeGraph}
	if len(AllShapes) != len(want) {
		t.Errorf("AllShapes length = %d, want %d", len(AllShapes), len(want))
	}

	for i, s := range AllShapes {
		if s != want[i] {
			t.Errorf("AllShapes[%d] = %v, want %v", i, s, want[i])
		}
	}
}
