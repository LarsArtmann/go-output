package output

import (
	"bytes"
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

//nolint:exhaustruct // Benchmark uses minimal struct initialization
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

//nolint:exhaustruct // Benchmark uses minimal struct initialization
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

func BenchmarkCSVWriter(b *testing.B) {
	var buf bytes.Buffer
	const headerCell = "Header"
	const dataCell = "Cell"
	headers := make([]string, 10)
	for i := range headers {
		headers[i] = headerCell
	}
	rows := make([][]string, 100)
	for i := range rows {
		row := make([]string, 10)
		for j := range row {
			row[j] = dataCell
		}
		rows[i] = row
	}

	b.ResetTimer()
	for b.Loop() {
		buf.Reset()
		w := NewCSVWriter(&buf)
		_ = w.WriteHeader(headers)
		for _, row := range rows {
			_ = w.WriteRow(row)
		}
		w.Flush()
	}
}

func BenchmarkMarkdownTable(b *testing.B) {
	md := NewMarkdownTable()
	const headerCell = "Header"
	const dataCell = "Cell"
	headers := make([]string, 10)
	for i := range headers {
		headers[i] = headerCell
	}
	md.SetHeaders(headers)
	rows := make([][]string, 100)
	for i := range rows {
		row := make([]string, 10)
		for j := range row {
			row[j] = dataCell
		}
		rows[i] = row
		md.AddRow(row)
	}

	b.ResetTimer()
	for b.Loop() {
		md.Render()
	}
}
