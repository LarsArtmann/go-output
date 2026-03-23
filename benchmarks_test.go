package output

import (
	"testing"
)

func BenchmarkASCIITreeRenderer(b *testing.B) {
	root := NewTreeNode("root", "Root")
	for i := 0; i < 100; i++ {
		child := NewTreeNode("child", "Child")
		for j := 0; j < 10; j++ {
			child.AddChild(NewTreeNode("leaf", "Leaf"))
		}
		root.AddChild(child)
	}
	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(root)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	for i := 0; i < 100; i++ {
		row := make([]string, 10)
		for j := range row {
			row[j] = "Cell"
		}
		renderer.AddRow(row)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.Render()
	}
}

func BenchmarkMermaidRenderer(b *testing.B) {
	renderer := NewMermaidRenderer()
	nodes := make([]GraphNode, 100)
	for i := range nodes {
		nodes[i] = GraphNode{ID: "node", Label: "Node"}
	}
	renderer.SetNodes(nodes)
	edges := make([]GraphEdge, 99)
	for i := range edges {
		edges[i] = GraphEdge{From: "node", To: "node"}
	}
	renderer.SetEdges(edges)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.Render()
	}
}

func BenchmarkDOTRenderer(b *testing.B) {
	renderer := NewDOTRenderer()
	nodes := make([]GraphNode, 100)
	for i := range nodes {
		nodes[i] = GraphNode{ID: "node", Label: "Node"}
	}
	renderer.SetNodes(nodes)
	edges := make([]GraphEdge, 99)
	for i := range edges {
		edges[i] = GraphEdge{From: "node", To: "node"}
	}
	renderer.SetEdges(edges)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.Render()
	}
}

func BenchmarkTableDataCreateRowEdges(b *testing.B) {
	data := NewTableData([]string{"A", "B", "C", "D", "E"})
	for i := 0; i < 1000; i++ {
		data.AddRow([]string{"1", "2", "3", "4", "5"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.CreateRowEdges()
	}
}
