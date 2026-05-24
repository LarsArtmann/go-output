package graph

import (
	"testing"

	"github.com/larsartmann/go-output"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges(testEdgesAB())

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph G {", "Output should contain digraph declaration")
	assertContains(t, out, "\"A\"", "Output should contain node A in quotes")
	assertContains(t, out, "label=\"Node A\"", "Output should contain label for node A")
	assertContains(t, out, "\"A\" -> \"B\"", "Output should contain directed edge A -> B")
	assertContains(t, out, "}", "Output should close with }")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTUndirectedRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewUndirectedDOTRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges(testEdgesAB())

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "graph G {", "Undirected graph should use 'graph' keyword")
	assertContains(t, out, "\"A\" -- \"B\"", "Undirected edge should use --")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRendererWithStyles(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes([]output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Styled Node"),
			Shape: output.ShapeEllipse,
			Style: output.GraphStyle{
				FillColor:   "#ff0000",
				StrokeColor: "#000000",
			},
		},
	})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "shape=ellipse", "Output should contain shape attribute")
	assertContains(t, out, "fillcolor=#ff0000", "Output should contain fillcolor")
}

func TestDOTRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetEdges([]output.GraphEdge{testEdgeAB("uses")})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "label=\"uses\"", "Output should contain edge label")
}

func TestDOTFromTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"ID", "Name"})
	data.AddRow([]string{"1", "Alice"})
	data.AddRow([]string{"2", "Bob"})

	renderer := DOTFromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph", "Output should be a digraph")
	assertContains(t, out, "row0", "Output should contain row0 node")
	assertContains(
		t,
		out,
		"\"row0\" -> \"row1\"",
		"Output should contain edge from row0 to row1",
	)
}

func TestDOTFromTree(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child", "Child"))

	renderer := DOTFromTree(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph", "Output should be a digraph")
	assertContains(t, out, "Root", "Output should contain 'Root' label")
	assertContains(t, out, "Child", "Output should contain 'Child' label")
}

func TestDOTRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	testEmptyRendererOutput(t, renderer, testDOTEmptyExpected())
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTSetGraphID(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetGraphID("MyGraph")
	renderer.SetNodes(
		[]output.GraphNode{newTestNode("A", "A")},
	)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph MyGraph {", "Output should use custom graph ID")
}

func TestDOTFromTreeWithEmptyID(t *testing.T) {
	t.Parallel()

	root := &output.TreeNode{
		Label: output.NewBrandedID[output.TreeNodeLabelBrand]("My Root"),
		Children: []*output.TreeNode{
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Child Node")},
		},
	}

	renderer := DOTFromTree(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "My_Root", "empty ID should use label slug")
	assertContains(t, out, "Child_Node", "empty ID should use label slug")
}

func TestDOTFromTableDataNil(t *testing.T) {
	t.Parallel()

	renderer := DOTFromTableData(nil)
	if renderer == nil {
		t.Fatal("DOTFromTableData(nil) should return non-nil renderer")
	}
}

func TestDOTFromTreeNil(t *testing.T) {
	t.Parallel()

	renderer := DOTFromTree(nil)
	if renderer == nil {
		t.Fatal("DOTFromTree(nil) should return non-nil renderer")
	}
}
