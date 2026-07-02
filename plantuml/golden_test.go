package plantuml

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_PlantUML_SimpleDiagram(t *testing.T) {
	t.Parallel()

	d := NewPlantUMLDiagram()
	d.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Alpha"),
	})
	d.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("b"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Beta"),
	})
	d.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("b"),
	})

	got, err := d.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_PlantUML_ShapedNodes(t *testing.T) {
	t.Parallel()

	d := NewPlantUMLDiagram()
	d.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("start"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Start"),
		Shape: output.NodeShapeCircle,
	})
	d.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("end"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("End"),
		Shape: output.NodeShapeCircle,
	})
	d.AddEdge(output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("start"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("end"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("flows to"),
	})

	got, err := d.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
