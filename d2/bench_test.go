package d2

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
)

func generateBenchmarkD2Nodes(n int) []Node {
	nodes := make([]Node, 0, n)
	for i := range n {
		nodes = append(nodes, Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
			Label: output.NewBrandedID[output.D2NodeLabelBrand](fmt.Sprintf("Node %d", i)),
		})
	}

	return nodes
}

func generateBenchmarkD2Edges(n int) []Edge {
	edges := make([]Edge, 0, n)
	for i := range n {
		edges = append(edges, Edge{
			From: output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
			To:   output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i+1)),
		})
	}

	return edges
}

func generateBenchmarkD2Tables(n int) []Table {
	tables := make([]Table, 0, n)
	for i := range n {
		tables = append(tables, Table{
			Name: fmt.Sprintf("table%d", i),
			Columns: []Column{
				{Name: "id", Type: "int"},
				{Name: "name", Type: "string"},
				{Name: "email", Type: "string"},
			},
		})
	}

	return tables
}

func benchmarkDiagram(b *testing.B, setup func(*Diagram)) {
	d := NewDiagram()
	setup(d)

	b.ResetTimer()

	for b.Loop() {
		_, _ = d.Render()
	}
}

func BenchmarkD2DiagramEmpty(b *testing.B) {
	benchmarkDiagram(b, func(_ *Diagram) {})
}

func BenchmarkD2DiagramNodes(b *testing.B) {
	benchmarkDiagram(b, func(d *Diagram) {
		for _, node := range generateBenchmarkD2Nodes(100) {
			d.AddNode(node)
		}
	})
}

func BenchmarkD2DiagramEdges(b *testing.B) {
	benchmarkDiagram(b, func(d *Diagram) {
		for _, node := range generateBenchmarkD2Nodes(100) {
			d.AddNode(node)
		}

		for _, edge := range generateBenchmarkD2Edges(99) {
			d.AddEdge(edge)
		}
	})
}

func BenchmarkD2DiagramTables(b *testing.B) {
	benchmarkDiagram(b, func(d *Diagram) {
		addBenchmarkD2Tables(d, 50)
	})
}

func BenchmarkD2DiagramStyledNodes(b *testing.B) {
	benchmarkDiagram(b, func(d *Diagram) {
		for i := range 100 {
			d.AddNode(Node{
				ID:    output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
				Label: output.NewBrandedID[output.D2NodeLabelBrand](fmt.Sprintf("Node %d", i)),
				Shape: ShapeCircle,
				Style: NodeStyle{
					Fill: "#f0f0f0",
					StrokeStyle: StrokeStyle{
						Stroke:   "#333333",
						FontSize: 14,
					},
					Shadow: true,
				},
				Tooltip: fmt.Sprintf("Tooltip for node %d", i),
				Link:    fmt.Sprintf("https://example.com/node/%d", i),
			})
		}
	})
}

func BenchmarkD2DiagramFullConfig(b *testing.B) {
	benchmarkDiagram(b, func(d *Diagram) {
		d.SetDirection(DirRight).
			SetTitle("Benchmark Diagram").
			SetLayout("elk")

		d.AddClass("highlight", NodeStyle{
			Fill:        "#ffcc00",
			StrokeStyle: StrokeStyle{Stroke: "#ff9900"},
		})

		for _, node := range generateBenchmarkD2Nodes(50) {
			d.AddNode(node)
		}

		for _, edge := range generateBenchmarkD2Edges(49) {
			d.AddEdge(edge)
		}

		addBenchmarkD2Tables(d, 10)
	})
}

func addBenchmarkD2Tables(d *Diagram, n int) {
	for _, table := range generateBenchmarkD2Tables(n) {
		d.AddTable(table.Name, table.Columns)
	}
}
