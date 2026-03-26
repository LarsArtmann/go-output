package output

import (
	"testing"
)

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

func TestGraphShapeIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		shape GraphShape
		want  bool
	}{
		{ShapeBox, true},
		{ShapeEllipse, true},
		{ShapeDiamond, true},
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

func TestGraphStyle(t *testing.T) {
	t.Parallel()
	style := GraphStyle{
		FillColor:   "red",
		StrokeColor: "blue",
		FontColor:   "green",
		FontSize:    12,
	}

	if style.FillColor != "red" {
		t.Errorf("FillColor = %q, want %q", style.FillColor, "red")
	}
	if style.StrokeColor != "blue" {
		t.Errorf("StrokeColor = %q, want %q", style.StrokeColor, "blue")
	}
	if style.FontColor != "green" {
		t.Errorf("FontColor = %q, want %q", style.FontColor, "green")
	}
	if style.FontSize != 12 {
		t.Errorf("FontSize = %d, want %d", style.FontSize, 12)
	}
}

func TestEdgeStyle(t *testing.T) {
	t.Parallel()
	style := EdgeStyle{
		Color:     "black",
		Style:     "dashed",
		ArrowHead: "arrow",
		ArrowTail: "arrow",
	}

	if style.Color != "black" {
		t.Errorf("Color = %q, want %q", style.Color, "black")
	}
	if style.Style != "dashed" {
		t.Errorf("Style = %q, want %q", style.Style, "dashed")
	}
	if style.ArrowHead != "arrow" {
		t.Errorf("ArrowHead = %q, want %q", style.ArrowHead, "arrow")
	}
	if style.ArrowTail != "arrow" {
		t.Errorf("ArrowTail = %q, want %q", style.ArrowTail, "arrow")
	}
}

func TestContainsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "xyz", false},
		{"hello", "hello", true},
		{"hi", "hello", false},
		{"", "", true},
		{"abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.substr, func(t *testing.T) {
			t.Parallel()
			if got := containsString(tt.s, tt.substr); got != tt.want {
				t.Errorf("containsString(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
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

func TestFormatStrings(t *testing.T) {
	t.Parallel()
	formats := []Format{FormatTable, FormatJSON, FormatCSV}
	got := formatStrings(formats)
	want := "table, json, csv"

	if got != want {
		t.Errorf("formatStrings() = %q, want %q", got, want)
	}
}
