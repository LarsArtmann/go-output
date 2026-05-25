package graph

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

const testGraphNodeID = "test-id"

func TestGraphNode(t *testing.T) {
	t.Parallel()

	node := output.NewGraphNode(testGraphNodeID, "Test Label")
	if node.ID.Get() != testGraphNodeID {
		t.Errorf("ID = %q, want %q", node.ID, testGraphNodeID)
	}

	if node.Label.Get() != "Test Label" {
		t.Errorf("Label = %q, want %q", node.Label, "Test Label")
	}

	if node.Shape != output.ShapeBox {
		t.Errorf("Shape = %v, want %v", node.Shape, output.ShapeBox)
	}

	if node.Metadata == nil {
		t.Error("Metadata is nil")
	}
}

func TestGraphEdge(t *testing.T) {
	t.Parallel()

	edge := output.NewGraphEdge("from-node", "to-node")
	if edge.From.Get() != "from-node" {
		t.Errorf("From = %q, want %q", edge.From, "from-node")
	}

	if edge.To.Get() != "to-node" {
		t.Errorf("To = %q, want %q", edge.To, "to-node")
	}
}

func TestParseGraphShape(t *testing.T) {
	tests := []parseEnumTestCase[output.GraphShape]{
		{Name: "box", Input: "box", Want: output.ShapeBox},
		{Name: "ellipse", Input: "ellipse", Want: output.ShapeEllipse},
		{Name: "diamond", Input: "diamond", Want: output.ShapeDiamond},
		{Name: "circle", Input: "circle", Want: output.ShapeCircle},
		{Name: "cylinder", Input: "cylinder", Want: output.ShapeCylinder},
		{Name: "hexagon", Input: "hexagon", Want: output.ShapeHexagon},
		{Name: "parallelogram", Input: "parallelogram", Want: output.ShapeParallelogram},
		{Name: "rect", Input: "rect", Want: output.ShapeRect},
		{Name: "invalid", Input: "invalid", WantErr: true},
		{Name: "empty", Input: "", WantErr: true},
	}
	testParseEnum(
		t,
		"output.ParseGraphShape",
		output.ParseGraphShape,
		tests,
		func(a, b output.GraphShape) bool { return a == b },
	)
}

func TestGraphShapeString(t *testing.T) {
	tests := []stringEnumTestCase[output.GraphShape]{
		{Value: output.ShapeBox, Want: "box"},
		{Value: output.ShapeEllipse, Want: "ellipse"},
		{Value: output.ShapeDiamond, Want: "diamond"},
		{Value: output.ShapeCircle, Want: "circle"},
	}

	testEnumString(
		t,
		"output.GraphShape.String",
		tests,
		func(s output.GraphShape) string { return s.String() },
	)
}

func TestGraphShapeAllowedValues(t *testing.T) {
	got := output.ShapeBox.AllowedValues()
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

	testhelpers.TestEnumIsValid(t, []output.GraphShape{
		output.ShapeBox,
		output.ShapeEllipse,
		output.ShapeDiamond,
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

func newTestGraphStyle(fontSize int) output.GraphStyle {
	return output.GraphStyle{
		FillColor:   "red",
		StrokeColor: "blue",
		FontColor:   "green",
		FontSize:    fontSize,
	}
}

func TestGraphStyle(t *testing.T) {
	t.Parallel()

	testGraphStyleFields(t, newTestGraphStyle(12), 12)
}

func TestEdgeStyle(t *testing.T) {
	t.Parallel()

	style := output.EdgeStyle{
		Color:     "black",
		Style:     "dashed",
		ArrowHead: "arrow",
		ArrowTail: "arrow",
	}

	testhelpers.TestStructFields(
		t,
		testhelpers.StringField("Color", style.Color, "black"),
		testhelpers.StringField("Style", style.Style, "dashed"),
		testhelpers.StringField("ArrowHead", style.ArrowHead, "arrow"),
		testhelpers.StringField("ArrowTail", style.ArrowTail, "arrow"),
	)
}

func TestGraphNodeStyle(t *testing.T) {
	t.Parallel()

	node := &output.GraphNode{Style: newTestGraphStyle(14)}

	testGraphStyleFields(t, node.Style, 14)
}

// testGraphStyleFields tests the common output.GraphStyle fields.
func testGraphStyleFields(t *testing.T, style output.GraphStyle, wantFontSize int) {
	testhelpers.TestStructFields(
		t,
		testhelpers.StringField("FillColor", style.FillColor, "red"),
		testhelpers.StringField("StrokeColor", style.StrokeColor, "blue"),
		testhelpers.StringField("FontColor", style.FontColor, "green"),
		testhelpers.IntField("FontSize", style.FontSize, wantFontSize),
	)
}
