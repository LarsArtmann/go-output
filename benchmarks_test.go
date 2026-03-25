package output

import (
	"testing"
)

func BenchmarkASCIITreeRenderer(b *testing.B) {
	root := NewTreeNode("root", "Root")
	for i := range 100 {
		child := NewTreeNode("child", "Child")
		for j := range 10 {
			_ = j // suppress unused
			child.AddChild(NewTreeNode("leaf", "Leaf"))
		}
		root.AddChild(child)
		_ = i // suppress unused
	}
	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(root)

	b.ResetTimer()
	for b.Loop() {
		renderer.Render()
	}
}

func BenchmarkHTMLRenderer(b *testing.B) {
	renderer := NewHTMLRenderer()
	headers := make([]string, 10)
	for i := range headers {
		headers[i] = "Header"
	}
	renderer.SetHeaders(headers)
	for range 100 {
		row := make([]string, 10)
		for j := range row {
			row[j] = "Cell"
		}
		renderer.AddRow(row)
	}

	b.ResetTimer()
	for b.Loop() {
		renderer.Render()
	}
}

func BenchmarkMermaidRenderer(b *testing.B) {
	renderer := NewMermaidRenderer()
	nodes := make([]GraphNode, 100)
	for i := range nodes {
		nodes[i] = GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand]("node"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Node"),
		}
		_ = i // suppress unused
	}
	renderer.SetNodes(nodes)
	edges := make([]GraphEdge, 99)
	for i := range edges {
		edges[i] = GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand]("node"),
			To:   NewBrandedID[GraphNodeIDBrand]("node"),
		}
		_ = i // suppress unused
	}
	renderer.SetEdges(edges)

	b.ResetTimer()
	for b.Loop() {
		renderer.Render()
	}
}

func BenchmarkDOTRenderer(b *testing.B) {
	renderer := NewDOTRenderer()
	nodes := make([]GraphNode, 100)
	for i := range nodes {
		nodes[i] = GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand]("node"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Node"),
		}
		_ = i // suppress unused
	}
	renderer.SetNodes(nodes)
	edges := make([]GraphEdge, 99)
	for i := range edges {
		edges[i] = GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand]("node"),
			To:   NewBrandedID[GraphNodeIDBrand]("node"),
		}
		_ = i // suppress unused
	}
	renderer.SetEdges(edges)

	b.ResetTimer()
	for b.Loop() {
		renderer.Render()
	}
}

func BenchmarkTableDataCreateRowEdges(b *testing.B) {
	data := NewTableData([]string{"A", "B", "C", "D", "E"})
	for range 1000 {
		data.AddRow([]string{"1", "2", "3", "4", "5"})
	}

	b.ResetTimer()
	for b.Loop() {
		data.CreateRowEdges()
	}
}
