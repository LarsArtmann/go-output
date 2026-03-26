package output

import (
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
		if !containsString(got, want) {
			t.Errorf("Error() = %q, should contain %q", got, want)
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestTableData(t *testing.T) {
	t.Parallel()
	t.Run("RowCount and ColCount", func(t *testing.T) {
		t.Parallel()
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
	})

	t.Run("CreateRowEdges", func(t *testing.T) {
		t.Parallel()
		t.Run("nil data", func(t *testing.T) {
			t.Parallel()
			var data *TableData
			if edges := data.CreateRowEdges(); edges != nil {
				t.Errorf("CreateRowEdges() on nil = %v, want nil", edges)
			}
		})

		t.Run("empty rows", func(t *testing.T) {
			t.Parallel()
			data := NewTableData([]string{"Name"})
			if edges := data.CreateRowEdges(); edges != nil {
				t.Errorf("CreateRowEdges() on empty = %v, want nil", edges)
			}
		})

		t.Run("single row", func(t *testing.T) {
			t.Parallel()
			data := NewTableData([]string{"Name"})
			data.AddRow([]string{"a"})
			if edges := data.CreateRowEdges(); edges != nil {
				t.Errorf("CreateRowEdges() on single row = %v, want nil", edges)
			}
		})

		t.Run("multiple rows", func(t *testing.T) {
			t.Parallel()
			data := NewTableData([]string{"Name"})
			data.AddRow([]string{"a"})
			data.AddRow([]string{"b"})
			data.AddRow([]string{"c"})
			edges := data.CreateRowEdges()
			if len(edges) != 2 {
				t.Fatalf("CreateRowEdges() returned %d edges, want 2", len(edges))
			}
			if edges[0].From != "row0" || edges[0].To != "row1" {
				t.Errorf("First edge = {%s, %s}, want {row0, row1}", edges[0].From, edges[0].To)
			}
			if edges[1].From != "row1" || edges[1].To != "row2" {
				t.Errorf("Second edge = {%s, %s}, want {row1, row2}", edges[1].From, edges[1].To)
			}
		})
	})
}

func TestGraphNode(t *testing.T) {
	t.Parallel()
	node := NewGraphNode("test-id", "Test Label")
	if node.ID.Get() != "test-id" {
		t.Errorf("ID = %q, want %q", node.ID, "test-id")
	}
	if node.Label.Get() != "Test Label" {
		t.Errorf("Label = %q, want %q", node.Label, "Test Label")
	}
	if node.Shape != ShapeBox {
		t.Errorf("Shape = %v, want %v", node.Shape, ShapeBox)
	}
	if node.Metadata == nil {
		t.Error("Metadata is nil")
	}
}

func TestGraphEdge(t *testing.T) {
	t.Parallel()
	edge := NewGraphEdge("from-node", "to-node")
	if edge.From.Get() != "from-node" {
		t.Errorf("From = %q, want %q", edge.From, "from-node")
	}
	if edge.To.Get() != "to-node" {
		t.Errorf("To = %q, want %q", edge.To, "to-node")
	}
}

func TestParseGraphShape(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    GraphShape
		wantErr bool
	}{
		{"box", "box", ShapeBox, false},
		{"ellipse", "ellipse", ShapeEllipse, false},
		{"diamond", "diamond", ShapeDiamond, false},
		{"circle", "circle", ShapeCircle, false},
		{"cylinder", "cylinder", ShapeCylinder, false},
		{"hexagon", "hexagon", ShapeHexagon, false},
		{"parallelogram", "parallelogram", ShapeParallelogram, false},
		{"rect", "rect", ShapeRect, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseGraphShape(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGraphShape() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseGraphShape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphShapeString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape GraphShape
		want  string
	}{
		{ShapeBox, "box"},
		{ShapeEllipse, "ellipse"},
		{ShapeDiamond, "diamond"},
		{ShapeCircle, "circle"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.shape.String(); got != tt.want {
				t.Errorf("GraphShape.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphShapeIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape GraphShape
		want  bool
	}{
		{ShapeBox, true},
		{ShapeEllipse, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.shape), func(t *testing.T) {
			t.Parallel()
			if got := tt.shape.IsValid(); got != tt.want {
				t.Errorf("GraphShape.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGraphShapeAllowedValues(t *testing.T) {
	t.Parallel()
	got := ShapeBox.AllowedValues()
	want := []string{
		"box",
		"ellipse",
		"diamond",
		"circle",
		"cylinder",
		"hexagon",
		"parallelogram",
		"rect",
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

func FuzzParseGraphShape(f *testing.F) {
	for _, shape := range ShapeBox.AllowedValues() {
		f.Add(shape)
	}
	f.Fuzz(func(t *testing.T, input string) {
		got, err := ParseGraphShape(input)
		if err != nil {
			if got != "" {
				t.Errorf("ParseGraphShape(%q) returned error but non-empty shape: %v", input, got)
			}
			return
		}
		if !got.IsValid() {
			t.Errorf("ParseGraphShape(%q) returned invalid shape: %v", input, got)
		}
	})
}
