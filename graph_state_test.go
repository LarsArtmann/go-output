package output

import (
	"testing"
)

func TestParseNodeShape(t *testing.T) {
	t.Parallel()

	got, err := ParseNodeShape("box")
	if err != nil {
		t.Fatalf("ParseNodeShape error: %v", err)
	}

	if string(got) != "box" {
		t.Errorf("ParseNodeShape = %q, want %q", got, "box")
	}
}

func TestParseNodeShape_Invalid(t *testing.T) {
	t.Parallel()

	_, err := ParseNodeShape("invalid")
	if err == nil {
		t.Error("expected error for invalid shape")
	}
}

func TestNodeShape_AllowedValues(t *testing.T) {
	t.Parallel()

	values := NodeShapeBox.AllowedValues()

	if len(values) == 0 {
		t.Error("AllowedValues should return non-empty slice")
	}
}

func TestNodeShape_IsValid(t *testing.T) {
	t.Parallel()

	if !NodeShapeBox.IsValid() {
		t.Error("NodeShapeBox should be valid")
	}

	invalid := NodeShape("invalid")

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

	m := NewGraphRendererState()
	AddTreeNodes(
		&m, root, "",
		func(n *TreeNode) string { return n.ID.Get() },
		NodeShapeBox,
	)

	if len(m.Nodes()) != 2 {
		t.Errorf("nodes count = %d, want 2", len(m.Nodes()))
	}

	if len(m.Edges()) != 1 {
		t.Errorf("edges count = %d, want 1", len(m.Edges()))
	}

	if m.Edges()[0].To.Get() != "child" {
		t.Errorf("edge To = %q, want %q", m.Edges()[0].To.Get(), "child")
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

func TestNewGraphRendererState(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()

	if m.Nodes() == nil {
		t.Error("Nodes() should not be nil")
	}

	if m.Edges() == nil {
		t.Error("Edges() should not be nil")
	}
}

func TestGraphRendererState_SetNodes(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.SetNodes(testNodesAB())

	assertSliceLen(t, "Nodes", m.Nodes(), 2)
}

func TestGraphRendererState_SetEdges(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.SetEdges(testEdgesAB())

	assertSliceLen(t, "Edges", m.Edges(), 1)
}

func TestGraphRendererState_AddNode(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.AddNode(GraphNode{ID: NewBrandedID[GraphNodeIDBrand]("a")})

	assertSliceLen(t, "Nodes", m.Nodes(), 1)
}

func TestGraphRendererState_AddEdge(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.AddEdge(newTestEdge("a", "b"))

	assertSliceLen(t, "Edges", m.Edges(), 1)
}

func TestGraphRendererState_DedupEdges(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.AddEdge(newTestEdge("a", "b"))
	m.AddEdge(newTestEdge("a", "b"))
	m.AddEdge(newTestEdge("b", "c"))
	m.AddEdge(newTestEdge("a", "b"))

	m.DedupEdges()

	assertSliceLen(t, "Edges after dedup", m.Edges(), 2)
}

func TestGraphRendererState_DedupEdgesEmpty(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.DedupEdges()

	assertSliceLen(t, "Edges", m.Edges(), 0)
}

func TestGraphRendererState_DedupEdgesSingle(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.AddEdge(newTestEdge("a", "b"))
	m.DedupEdges()

	assertSliceLen(t, "Edges", m.Edges(), 1)
}

func TestGraphRendererState_AddRowEdges(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	m.SetNodes(testNodesABC())

	data := NewTableData([]string{"A"})
	data.AddRow([]string{"1"})
	data.AddRow([]string{"2"})

	m.AddRowEdges(data)

	assertSliceLen(t, "Edges", m.Edges(), 1)
}

func TestGraphRendererState_SetNodesFromTableData(t *testing.T) {
	t.Parallel()

	m := NewGraphRendererState()
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"Alpha"})
	data.AddRow([]string{"Beta"})

	m.SetNodesFromTableData(data, func(i int, n *GraphNode) {})

	assertSliceLen(t, "Nodes", m.Nodes(), 2)
	assertSliceLen(t, "Edges", m.Edges(), 1)
}

func assertSliceLen[T any](t *testing.T, name string, slice []T, want int) {
	t.Helper()

	if len(slice) != want {
		t.Errorf("%s() len = %d, want %d", name, len(slice), want)
	}
}

// newTestEdge builds a GraphEdge with branded IDs from plain strings.
func newTestEdge(from, to string) GraphEdge {
	return GraphEdge{
		From: NewBrandedID[GraphNodeIDBrand](from),
		To:   NewBrandedID[GraphNodeIDBrand](to),
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
