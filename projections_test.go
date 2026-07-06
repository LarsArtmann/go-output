package output

import (
	"strings"
	"testing"
)

func TestTableToGraph(t *testing.T) {
	tbl := NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		AddRow("Lint", "done").
		Build()

	g := TableToGraph(tbl)

	nodes := g.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	edges := g.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges (consecutive rows), got %d", len(edges))
	}
}

func TestTableToGraph_NilTable(t *testing.T) {
	g := TableToGraph(nil)

	if len(g.Nodes()) != 0 {
		t.Errorf("expected 0 nodes for nil table, got %d", len(g.Nodes()))
	}
}

func TestGraphToTree_LinearChain(t *testing.T) {
	b := NewGraphBuilder()
	b.AddNode(*NewGraphNode("a", "A"))
	b.AddNode(*NewGraphNode("b", "B"))
	b.AddNode(*NewGraphNode("c", "C"))
	b.AddEdge(*NewGraphEdge("a", "b"))
	b.AddEdge(*NewGraphEdge("b", "c"))

	g := b.Build()
	root := GraphToTree(g)

	if root == nil {
		t.Fatal("expected non-nil root")
	}

	if root.ID.Get() != "a" {
		t.Errorf("expected root 'a', got %q", root.ID.Get())
	}

	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(root.Children))
	}

	child := root.Children[0]
	if child.ID.Get() != "b" {
		t.Errorf("expected child 'b', got %q", child.ID.Get())
	}

	if len(child.Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(child.Children))
	}

	if child.Children[0].ID.Get() != "c" {
		t.Errorf("expected grandchild 'c', got %q", child.Children[0].ID.Get())
	}
}

func TestGraphToTree_Branching(t *testing.T) {
	b := NewGraphBuilder()
	b.AddNode(*NewGraphNode("root", "Root"))
	b.AddNode(*NewGraphNode("a", "A"))
	b.AddNode(*NewGraphNode("b", "B"))
	b.AddEdge(*NewGraphEdge("root", "a"))
	b.AddEdge(*NewGraphEdge("root", "b"))

	g := b.Build()
	root := GraphToTree(g)

	if root.ID.Get() != "root" {
		t.Errorf("expected root 'root', got %q", root.ID.Get())
	}

	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
}

func TestGraphToTree_Empty(t *testing.T) {
	g := Graph{}
	root := GraphToTree(g)

	if root != nil {
		t.Errorf("expected nil for empty graph, got %v", root)
	}
}

func TestGraphToTree_CycleGuard(t *testing.T) {
	b := NewGraphBuilder()
	b.AddNode(*NewGraphNode("a", "A"))
	b.AddNode(*NewGraphNode("b", "B"))
	b.AddEdge(*NewGraphEdge("a", "b"))
	b.AddEdge(*NewGraphEdge("b", "a"))

	g := b.Build()
	root := GraphToTree(g)

	if root == nil {
		t.Fatal("expected non-nil root")
	}

	if root.ID.Get() != "a" {
		t.Errorf("expected root 'a', got %q", root.ID.Get())
	}

	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child (cycle guard), got %d", len(root.Children))
	}
}

func TestGraphToTree_DisconnectedForest(t *testing.T) {
	b := NewGraphBuilder()
	b.AddNode(*NewGraphNode("a", "A"))
	b.AddNode(*NewGraphNode("b", "B"))
	b.AddNode(*NewGraphNode("c", "C"))
	b.AddNode(*NewGraphNode("d", "D"))
	b.AddEdge(*NewGraphEdge("a", "b"))
	b.AddEdge(*NewGraphEdge("c", "d"))

	g := b.Build()
	root := GraphToTree(g)

	if root == nil {
		t.Fatal("expected non-nil root")
	}

	// Only the first root found (a) should be in the tree — c/d are disconnected.
	if root.ID.Get() != "a" {
		t.Errorf("expected root 'a', got %q", root.ID.Get())
	}

	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(root.Children))
	}

	if root.Children[0].ID.Get() != "b" {
		t.Errorf("expected child 'b', got %q", root.Children[0].ID.Get())
	}
}

func TestTableToGraph_CustomLabelFunc(t *testing.T) {
	tbl := NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		Build()

	// Use a custom label function that returns only the cell value (not "header: cell").
	customLabel := func(header, cell string) string { return cell }

	g := TableToGraph(tbl, WithGraphNodeLabelFunc(customLabel))

	nodes := g.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}

	// Custom label: "Compile, done" → each part is just the cell value.
	// NodesFromTable joins all cell labels with ", ".
	label := nodes[0].Label.Get()
	if !strings.Contains(label, "Compile") {
		t.Errorf("expected label to contain 'Compile', got %q", label)
	}

	// Verify the default "header:" prefix is NOT present (custom func was used).
	if strings.Contains(label, "Name:") {
		t.Errorf("custom label func should not produce 'Name:' prefix, got %q", label)
	}
}

func TestGraphToTable(t *testing.T) {
	b := NewGraphBuilder()
	b.AddNode(*NewGraphNode("a", "Alpha"))
	b.AddNode(*NewGraphNode("b", "Beta"))

	g := b.Build()
	tbl := GraphToTable(g)

	if tbl == nil {
		t.Fatal("expected non-nil table")
	}

	if len(tbl.Headers) != 2 {
		t.Fatalf("expected 2 headers, got %d", len(tbl.Headers))
	}

	if tbl.Headers[0] != "ID" || tbl.Headers[1] != "Label" {
		t.Errorf("expected headers [ID, Label], got %v", tbl.Headers)
	}

	if len(tbl.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows))
	}

	if tbl.Rows[0][0] != "a" || tbl.Rows[0][1] != "Alpha" {
		t.Errorf("expected first row [a, Alpha], got %v", tbl.Rows[0])
	}
}

func TestGraphToTable_Empty(t *testing.T) {
	g := Graph{}
	tbl := GraphToTable(g)

	if tbl != nil {
		t.Errorf("expected nil for empty graph, got %v", tbl)
	}
}
