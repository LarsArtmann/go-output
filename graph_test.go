package output

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
)

const testGraphNodeID = "test-id"

func TestGraphNode(t *testing.T) {
	t.Parallel()

	node := NewGraphNode(testGraphNodeID, "Test Label")
	if node.ID.Get() != testGraphNodeID {
		t.Errorf("ID = %q, want %q", node.ID, testGraphNodeID)
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
	tests := []parseEnumTestCase[GraphShape]{
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
	testParseEnum(
		t,
		"ParseGraphShape",
		ParseGraphShape,
		tests,
		func(a, b GraphShape) bool { return a == b },
	)
}

func TestGraphShapeString(t *testing.T) {
	tests := []stringEnumTestCase[GraphShape]{
		{ShapeBox, "box"},
		{ShapeEllipse, "ellipse"},
		{ShapeDiamond, "diamond"},
		{ShapeCircle, "circle"},
	}

	testEnumString(t, "GraphShape.String", tests, func(s GraphShape) string { return s.String() })
}

func TestGraphShapeAllowedValues(t *testing.T) {
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

	testAllowedValues(t, "AllowedValues", got, want)
}

func TestGraphShapeIsValid(t *testing.T) {
	t.Parallel()

	gentest.TestEnumIsValid(t, []GraphShape{
		ShapeBox,
		ShapeEllipse,
		ShapeDiamond,
		"invalid",
		"",
	}, []bool{
		true,
		true,
		true,
		false,
		false,
	})
}

func TestGraphStyle(t *testing.T) {
	t.Parallel()

	style := GraphStyle{
		FillColor:   "red",
		StrokeColor: "blue",
		FontColor:   "green",
		FontSize:    12,
	}

	gentest.TestStructFields(t,
		gentest.StringField("FillColor", style.FillColor, "red"),
		gentest.StringField("StrokeColor", style.StrokeColor, "blue"),
		gentest.StringField("FontColor", style.FontColor, "green"),
		gentest.IntField("FontSize", style.FontSize, 12),
	)
}

func TestEdgeStyle(t *testing.T) {
	t.Parallel()

	style := EdgeStyle{
		Color:     "black",
		Style:     "dashed",
		ArrowHead: "arrow",
		ArrowTail: "arrow",
	}

	gentest.TestStructFields(t,
		gentest.StringField("Color", style.Color, "black"),
		gentest.StringField("Style", style.Style, "dashed"),
		gentest.StringField("ArrowHead", style.ArrowHead, "arrow"),
		gentest.StringField("ArrowTail", style.ArrowTail, "arrow"),
	)
}

func TestGetStyle(t *testing.T) {
	t.Parallel()

	node := &GraphNode{
		ID:    NewBrandedID[GraphNodeIDBrand]("id"),
		Label: NewBrandedID[GraphNodeLabelBrand]("label"),
		Style: GraphStyle{FillColor: "red", StrokeColor: "blue", FontColor: "green", FontSize: 14},
	}

	style := node.GetStyle()

	gentest.TestStructFields(t,
		gentest.StringField("FillColor", style.FillColor, "red"),
		gentest.StringField("StrokeColor", style.StrokeColor, "blue"),
		gentest.StringField("FontColor", style.FontColor, "green"),
		gentest.IntField("FontSize", style.FontSize, 14),
	)
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

			if got := strings.Contains(tt.s, tt.substr); got != tt.want {
				t.Errorf("containsString(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}
