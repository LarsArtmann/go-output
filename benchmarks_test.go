package output

import (
	"io"
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

	headers := FilledStrings(10, "Header")
	renderer.SetHeaders(headers)

	for range 100 {
		row := FilledStrings(10, "Cell")
		renderer.AddRow(row)
	}

	b.ResetTimer()

	for b.Loop() {
		renderer.Render()
	}
}

// generateBenchmarkNodes creates a slice of GraphNode for benchmarking.
func generateBenchmarkNodes(n int) []GraphNode {
	nodes := make([]GraphNode, n)
	for i := range nodes {
		nodes[i] = GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand]("node"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Node"),
		}
	}

	return nodes
}

// generateBenchmarkEdges creates a slice of GraphEdge for benchmarking.
func generateBenchmarkEdges(n int) []GraphEdge {
	edges := make([]GraphEdge, n)
	for i := range edges {
		edges[i] = GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand]("node"),
			To:   NewBrandedID[GraphNodeIDBrand]("node"),
		}
	}

	return edges
}

// benchmarkGraphRenderer sets up nodes and edges for a graph renderer benchmark.
func benchmarkGraphRenderer(b *testing.B, renderer GraphRenderer) {
	nodes := generateBenchmarkNodes(100)
	renderer.SetNodes(nodes)

	edges := generateBenchmarkEdges(99)
	renderer.SetEdges(edges)

	b.ResetTimer()

	for b.Loop() {
		renderer.Render()
	}
}

func BenchmarkMermaidRenderer(b *testing.B) {
	renderer := NewMermaidRenderer()
	benchmarkGraphRenderer(b, renderer)
}

func BenchmarkDOTRenderer(b *testing.B) {
	renderer := NewDOTRenderer()
	benchmarkGraphRenderer(b, renderer)
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
	const (
		headerCell = "Header"
		dataCell   = "Cell"
	)

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

	benchmarkTableWriter(b, headers, rows, func(w io.Writer) TableWriter {
		return NewCSVWriter(w)
	})
}

func BenchmarkMarkdownTable(b *testing.B) {
	md := NewMarkdownTable()

	const (
		headerCell = "Header"
		dataCell   = "Cell"
	)

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

// BenchmarkStruct is used for JSON unmarshal benchmarks.
type BenchmarkStruct struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Items     []string `json:"items"`
	Count     int      `json:"count"`
	Active    bool     `json:"active"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// BenchmarkYAMLStruct is used for YAML unmarshal benchmarks.
type BenchmarkYAMLStruct struct {
	ID        int      `yaml:"id"`
	Name      string   `yaml:"name"`
	Items     []string `yaml:"items"`
	Count     int      `yaml:"count"`
	Active    bool     `yaml:"active"`
	CreatedAt string   `yaml:"created_at"`
	UpdatedAt string   `yaml:"updated_at"`
}

func NewBenchmarkData() BenchmarkData {
	return BenchmarkData{
		ID:        12345,
		Name:      "Test Project Alpha",
		Items:     []string{"item1", "item2", "item3", "item4", "item5"},
		Count:     100,
		Active:    true,
		CreatedAt: "2026-03-22T10:00:00Z",
		UpdatedAt: "2026-03-22T12:00:00Z",
	}
}
