package output

import (
	"testing"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges(testEdgesAB())

	output := renderer.Render()

	assertContains(t, output, "digraph G {", "Output should contain digraph declaration")
	assertContains(t, output, "\"A\"", "Output should contain node A in quotes")
	assertContains(t, output, "label=\"Node A\"", "Output should contain label for node A")
	assertContains(t, output, "\"A\" -> \"B\"", "Output should contain directed edge A -> B")
	assertContains(t, output, "}", "Output should close with }")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTUndirectedRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewUndirectedDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges(testEdgesAB())

	output := renderer.Render()

	assertContains(t, output, "graph G {", "Undirected graph should use 'graph' keyword")
	assertContains(t, output, "\"A\" -- \"B\"", "Undirected edge should use --")
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

	assertContains(t, output, "shape=ellipse", "Output should contain shape attribute")
	assertContains(t, output, "fillcolor=#ff0000", "Output should contain fillcolor")
}

func TestDOTRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetEdges([]GraphEdge{testEdgeAB("uses")})

	output := renderer.Render()

	assertContains(t, output, "label=\"uses\"", "Output should contain edge label")
}

func TestDOTFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"ID", "Name"})
	data.AddRow([]string{"1", "Alice"})
	data.AddRow([]string{"2", "Bob"})

	renderer := DOTFromTableData(data)
	output := renderer.Render()

	assertContains(t, output, "digraph", "Output should be a digraph")
	assertContains(t, output, "row0", "Output should contain row0 node")
	assertContains(
		t,
		output,
		"\"row0\" -> \"row1\"",
		"Output should contain edge from row0 to row1",
	)
}

func TestDOTFromTree(t *testing.T) {
	t.Parallel()

	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child", "Child"))

	renderer := DOTFromTree(root)
	output := renderer.Render()

	assertContains(t, output, "digraph", "Output should be a digraph")
	assertContains(t, output, "Root", "Output should contain 'Root' label")
	assertContains(t, output, "Child", "Output should contain 'Child' label")
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

	assertContains(t, output, "digraph MyGraph {", "Output should use custom graph ID")
}
