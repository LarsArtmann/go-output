package output

import (
	"strings"
	"testing"
)

func TestDOTRenderer(t *testing.T) {
	renderer := NewDOTRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "A", Label: "Node A"},
		{ID: "B", Label: "Node B"},
	})
	renderer.SetEdges([]GraphEdge{
		{From: "A", To: "B"},
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

func TestDOTUndirectedRenderer(t *testing.T) {
	renderer := NewUndirectedDOTRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "A", Label: "Node A"},
		{ID: "B", Label: "Node B"},
	})
	renderer.SetEdges([]GraphEdge{
		{From: "A", To: "B"},
	})

	output := renderer.Render()

	if !strings.Contains(output, "graph G {") {
		t.Error("Undirected graph should use 'graph' keyword")
	}
	if !strings.Contains(output, "\"A\" -- \"B\"") {
		t.Error("Undirected edge should use --")
	}
}

func TestDOTRendererWithStyles(t *testing.T) {
	renderer := NewDOTRenderer()
	renderer.SetNodes([]GraphNode{
		{
			ID:    "A",
			Label: "Styled Node",
			Shape: ShapeEllipse,
			Style: GraphStyle{
				FillColor: "#ff0000",
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

func TestDOTRendererWithEdgeLabel(t *testing.T) {
	renderer := NewDOTRenderer()
	renderer.SetEdges([]GraphEdge{
		{From: "A", To: "B", Label: "uses"},
	})

	output := renderer.Render()

	if !strings.Contains(output, "label=\"uses\"") {
		t.Error("Output should contain edge label")
	}
}

func TestDOTFromTableData(t *testing.T) {
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
	renderer := NewDOTRenderer()
	output := renderer.Render()

	if !strings.Contains(output, "digraph G {") {
		t.Error("Empty DOT should still have digraph declaration")
	}
	if !strings.Contains(output, "rankdir=TB") {
		t.Error("Empty DOT should have default attributes")
	}
}

func TestDOTSetGraphID(t *testing.T) {
	renderer := NewDOTRenderer()
	renderer.SetGraphID("MyGraph")
	renderer.SetNodes([]GraphNode{{ID: "A", Label: "A"}})

	output := renderer.Render()

	if !strings.Contains(output, "digraph MyGraph {") {
		t.Error("Output should use custom graph ID")
	}
}
