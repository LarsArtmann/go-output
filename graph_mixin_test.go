package output

import (
	"testing"
)

func TestParseGraphShape(t *testing.T) {
	t.Parallel()

	got, err := ParseGraphShape("rect")
	if err != nil {
		t.Fatalf("ParseGraphShape error: %v", err)
	}

	if string(got) != "rect" {
		t.Errorf("ParseGraphShape = %q, want %q", got, "rect")
	}
}

func TestParseGraphShape_Invalid(t *testing.T) {
	t.Parallel()

	_, err := ParseGraphShape("invalid")
	if err == nil {
		t.Error("expected error for invalid shape")
	}
}

func TestGraphShape_AllowedValues(t *testing.T) {
	t.Parallel()

	values := ShapeRect.AllowedValues()

	if len(values) == 0 {
		t.Error("AllowedValues should return non-empty slice")
	}
}

func TestGraphShape_IsValid(t *testing.T) {
	t.Parallel()

	if !ShapeRect.IsValid() {
		t.Error("ShapeRect should be valid")
	}

	invalid := GraphShape("invalid")

	if invalid.IsValid() {
		t.Error("invalid shape should not be valid")
	}
}

func TestNewGraphEdge(t *testing.T) {
	t.Parallel()

	edge := NewGraphEdge("a", "b")

	if edge.From.Get() != "a" {
		t.Errorf("From = %q, want %q", edge.From.Get(), "a")
	}

	if edge.To.Get() != "b" {
		t.Errorf("To = %q, want %q", edge.To.Get(), "b")
	}
}

func TestDefaultGraphNodeLabel(t *testing.T) {
	t.Parallel()

	got := DefaultGraphNodeLabel("Name", "Alpha")

	if got != "Name: Alpha" {
		t.Errorf("DefaultGraphNodeLabel = %q, want %q", got, "Name: Alpha")
	}
}

func TestAddTreeNodes(t *testing.T) {
	t.Parallel()

	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	root.AddChild(child)

	var nodes []GraphNode

	var edges []GraphEdge

	AddTreeNodes(
		&nodes, &edges, root, "",
		func(n *TreeNode) string { return n.ID.Get() },
		ShapeRect,
	)

	if len(nodes) != 2 {
		t.Errorf("nodes count = %d, want 2", len(nodes))
	}

	if len(edges) != 1 {
		t.Errorf("edges count = %d, want 1", len(edges))
	}

	if edges[0].To.Get() != "child" {
		t.Errorf("edge To = %q, want %q", edges[0].To.Get(), "child")
	}
}

func TestNodesFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Name", "Value"})
	data.AddRow([]string{"Alpha", "100"})
	data.AddRow([]string{"Beta", "200"})

	nodes := NodesFromTableData(data, DefaultGraphNodeLabel)

	if len(nodes) != 2 {
		t.Fatalf("nodes count = %d, want 2", len(nodes))
	}

	want := "Name: Alpha\nValue: 100"

	if nodes[0].Label.Get() != want {
		t.Errorf("label = %q, want %q", nodes[0].Label.Get(), want)
	}
}

func TestNodesFromTableData_Nil(t *testing.T) {
	t.Parallel()

	nodes := NodesFromTableData(nil, DefaultGraphNodeLabel)

	if nodes != nil {
		t.Error("expected nil for nil data")
	}
}

func TestNewGraphRendererMixin(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()

	if m.Nodes() == nil {
		t.Error("Nodes() should not be nil")
	}

	if m.Edges() == nil {
		t.Error("Edges() should not be nil")
	}

	if len(m.Nodes()) != 0 {
		t.Errorf("Nodes() len = %d, want 0", len(m.Nodes()))
	}
}

func TestGraphRendererMixin_SetNodes(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()
	m.SetNodes(testNodesAB())

	if len(m.Nodes()) != 2 {
		t.Errorf("Nodes() len = %d, want 2", len(m.Nodes()))
	}
}

func TestGraphRendererMixin_SetEdges(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()
	m.SetEdges(testEdgesAB())

	if len(m.Edges()) != 1 {
		t.Errorf("Edges() len = %d, want 1", len(m.Edges()))
	}
}

func TestGraphRendererMixin_NodesPtr(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()

	if m.NodesPtr() == nil {
		t.Error("NodesPtr() should not be nil")
	}
}

func TestGraphRendererMixin_EdgesPtr(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()

	if m.EdgesPtr() == nil {
		t.Error("EdgesPtr() should not be nil")
	}
}

func TestGraphRendererMixin_AddRowEdges(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()
	m.SetNodes(testNodesABC())

	data := NewTableData([]string{"A"})
	data.AddRow([]string{"1"})
	data.AddRow([]string{"2"})

	m.AddRowEdges(data)

	if len(m.Edges()) != 1 {
		t.Errorf("Edges() len = %d, want 2", len(m.Edges()))
	}
}

func TestGraphRendererMixin_SetNodesFromTableData(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererMixin()
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"Alpha"})
	data.AddRow([]string{"Beta"})

	m.SetNodesFromTableData(data, func(i int, n *GraphNode) {})

	if len(m.Nodes()) != 2 {
		t.Errorf("Nodes() len = %d, want 2", len(m.Nodes()))
	}

	if len(m.Edges()) != 1 {
		t.Errorf("Edges() len = %d, want 1", len(m.Edges()))
	}
}

func TestUnsupportedFormatError_Error(t *testing.T) {
	t.Parallel()

	err := &UnsupportedFormatError{Format: FormatD2}

	if err.Error() == "" {
		t.Error("Error() should not be empty")
	}
}

func TestTableData_GetHeaders_GetRows(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"A", "B"})
	data.AddRow([]string{"1", "2"})

	headers := data.GetHeaders()

	if len(headers) != 2 {
		t.Errorf("GetHeaders() len = %d, want 2", len(headers))
	}

	rows := data.GetRows()

	if len(rows) != 1 {
		t.Errorf("GetRows() len = %d, want 1", len(rows))
	}
}
