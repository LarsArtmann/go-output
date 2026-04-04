package output

import (
	"strings"
	"testing"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges([]GraphEdge{
		{From: NewBrandedID[GraphNodeIDBrand]("A"), To: NewBrandedID[GraphNodeIDBrand]("B")},
	})

	output := renderer.Render()

	if !strings.Contains(output, "digraph G {") {
		t.Error("Output should contain digraph declaration")
	}

	if !strings.Contains(output, "\"A\"") {
		t.Error("Output should contain node A in quotes")
	}

	if !strings.Contains(output, "label=\"Node A\"") {
		t.Error("Output should contain label for node A")
	}

	if !strings.Contains(output, "\"A\" -> \"B\"") {
		t.Error("Output should contain directed edge A -> B")
	}

	if !strings.Contains(output, "}") {
		t.Error("Output should close with }")
	}
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTUndirectedRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewUndirectedDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges([]GraphEdge{
		{From: NewBrandedID[GraphNodeIDBrand]("A"), To: NewBrandedID[GraphNodeIDBrand]("B")},
	})

	output := renderer.Render()

	if !strings.Contains(output, "graph G {") {
		t.Error("Undirected graph should use 'graph' keyword")
	}

	if !strings.Contains(output, "\"A\" -- \"B\"") {
		t.Error("Undirected edge should use --")
	}
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRendererWithStyles(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes([]GraphNode{
		{
			ID:    NewBrandedID[GraphNodeIDBrand]("A"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Styled Node"),
			Shape: ShapeEllipse,
			Style: GraphStyle{
				FillColor:   "#ff0000",
				StrokeColor: "#000000",
			},
		},
	})

	output := renderer.Render()

	if !strings.Contains(output, "shape=ellipse") {
		t.Error("Output should contain shape attribute")
	}

	if !strings.Contains(output, "fillcolor=#ff0000") {
		t.Error("Output should contain fillcolor")
	}
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetEdges([]GraphEdge{testEdgeAB("uses")})

	output := renderer.Render()

	if !strings.Contains(output, "label=\"uses\"") {
		t.Error("Output should contain edge label")
	}
}

func TestDOTFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"ID", "Name"})
	data.AddRow([]string{"1", "Alice"})
	data.AddRow([]string{"2", "Bob"})

	renderer := DOTFromTableData(data)
	output := renderer.Render()

	if !strings.Contains(output, "digraph") {
		t.Error("Output should be a digraph")
	}

	if !strings.Contains(output, "row0") {
		t.Error("Output should contain row0 node")
	}

	if !strings.Contains(output, "\"row0\" -> \"row1\"") {
		t.Error("Output should contain edge from row0 to row1")
	}
}

func TestDOTFromTree(t *testing.T) {
	t.Parallel()

	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child", "Child"))

	renderer := DOTFromTree(root)
	output := renderer.Render()

	if !strings.Contains(output, "digraph") {
		t.Error("Output should be a digraph")
	}

	if !strings.Contains(output, "Root") {
		t.Error("Output should contain 'Root' label")
	}

	if !strings.Contains(output, "Child") {
		t.Error("Output should contain 'Child' label")
	}
}

func TestDOTRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	testEmptyRendererOutput(t, renderer, []ExpectedOutput{
		{Substring: "digraph G {", Message: "Empty DOT should still have digraph declaration"},
		{Substring: "rankdir=TB", Message: "Empty DOT should have default attributes"},
	})
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTSetGraphID(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetGraphID("MyGraph")
	renderer.SetNodes(
		[]GraphNode{
			{
				ID:    NewBrandedID[GraphNodeIDBrand]("A"),
				Label: NewBrandedID[GraphNodeLabelBrand]("A"),
			},
		},
	)

	output := renderer.Render()

	if !strings.Contains(output, "digraph MyGraph {") {
		t.Error("Output should use custom graph ID")
	}
}
