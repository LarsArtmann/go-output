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

	if node.Shape != output.NodeShapeBox {
		t.Errorf("Shape = %v, want %v", node.Shape, output.NodeShapeBox)
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

func TestParseNodeShape(t *testing.T) {
	tests := []testhelpers.ParseEnumTestCase[output.NodeShape]{
		{Name: "box", Input: "box", Want: output.NodeShapeBox},
		{Name: "ellipse", Input: "ellipse", Want: output.NodeShapeEllipse},
		{Name: "diamond", Input: "diamond", Want: output.NodeShapeDiamond},
		{Name: "circle", Input: "circle", Want: output.NodeShapeCircle},
		{Name: "cylinder", Input: "cylinder", Want: output.NodeShapeCylinder},
		{Name: "hexagon", Input: "hexagon", Want: output.NodeShapeHexagon},
		{Name: "parallelogram", Input: "parallelogram", Want: output.NodeShapeParallelogram},
		{Name: "rect", Input: "rect", Want: output.NodeShapeRect},
		{Name: "invalid", Input: "invalid", WantErr: true},
		{Name: "empty", Input: "", WantErr: true},
	}
	testhelpers.TestParseEnum(
		t,
		"output.ParseNodeShape",
		output.ParseNodeShape,
		tests,
		func(a, b output.NodeShape) bool { return a == b },
	)
}

func TestNodeShapeString(t *testing.T) {
	tests := []testhelpers.StringEnumTestCase[output.NodeShape]{
		{Value: output.NodeShapeBox, Want: "box"},
		{Value: output.NodeShapeEllipse, Want: "ellipse"},
		{Value: output.NodeShapeDiamond, Want: "diamond"},
		{Value: output.NodeShapeCircle, Want: "circle"},
	}

	testhelpers.TestEnumString(
		t,
		"output.NodeShape.String",
		tests,
		func(s output.NodeShape) string { return s.String() },
	)
}

func TestNodeShapeAllowedValues(t *testing.T) {
	got := output.NodeShapeBox.AllowedValues()
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

	testhelpers.TestAllowedValues(t, "AllowedValues", got, want)
}

func TestNodeShapeIsValid(t *testing.T) {
	t.Parallel()

	testhelpers.TestEnumIsValid(t, []output.NodeShape{
		output.NodeShapeBox,
		output.NodeShapeEllipse,
		output.NodeShapeDiamond,
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
		Fill:      "red",
		Stroke:    "blue",
		FontColor: "green",
		FontSize:  fontSize,
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
		Line:      output.LineStyleDashed,
		ArrowHead: "arrow",
		ArrowTail: "arrow",
	}

	testhelpers.TestStructFields(
		t,
		testhelpers.StringField("Color", style.Color, "black"),
		testhelpers.StringField("Line", style.Line.String(), "dashed"),
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
		testhelpers.StringField("Fill", style.Fill, "red"),
		testhelpers.StringField("Stroke", style.Stroke, "blue"),
		testhelpers.StringField("FontColor", style.FontColor, "green"),
		testhelpers.IntField("FontSize", style.FontSize, wantFontSize),
	)
}
