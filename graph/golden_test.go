package graph

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_DOT_SimpleGraph(t *testing.T) {
	t.Parallel()

	r := NewDOTRenderer()
	r.SetGraphID("G")
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Alpha"),
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("b"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Beta"),
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("c"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Gamma"),
	})
	r.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("b"),
	})
	r.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("b"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("c"),
	})

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_DOT_StyledNodes(t *testing.T) {
	t.Parallel()

	r := NewDOTRenderer()
	r.SetGraphID("pipeline")
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("build"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Build"),
		Shape: output.NodeShapeBox,
		Style: output.NodeStyle{Fill: "#4CAF50", Stroke: "#2E7D32"},
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("test"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
		Shape: output.NodeShapeDiamond,
		Style: output.NodeStyle{Fill: "#FF9800"},
	})
	r.AddEdge(output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("build"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("test"),
		Style: output.EdgeStyle{Line: output.LineStyleDashed},
	})

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_Mermaid_SimpleGraph(t *testing.T) {
	t.Parallel()

	r := NewMermaidRenderer()
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("start"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Start"),
		Shape: output.NodeShapeCircle,
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("process"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Process"),
		Shape: output.NodeShapeBox,
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("end"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("End"),
		Shape: output.NodeShapeCircle,
	})
	r.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("start"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("process"),
	})
	r.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("process"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("end"),
	})

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

//nolint:exhaustruct // Test uses minimal fields
func TestGolden_Mermaid_NoCodeFence(t *testing.T) {
	t.Parallel()

	r := NewMermaidRenderer().SetCodeFence(false)
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Alpha"),
	})
	r.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("b"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Beta"),
	})
	r.AddEdge(output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		To:   output.NewBrandedID[output.GraphNodeIDBrand]("b"),
	})

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
