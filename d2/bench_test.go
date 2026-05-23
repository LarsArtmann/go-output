package d2

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
)

func generateBenchmarkD2Nodes(n int) []D2Node {
	nodes := make([]D2Node, n)
	for i := range nodes {
		nodes[i] = D2Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
			Label: output.NewBrandedID[output.D2NodeLabelBrand](fmt.Sprintf("Node %d", i)),
		}
	}

	return nodes
}

func generateBenchmarkD2Edges(n int) []D2Edge {
	edges := make([]D2Edge, n)
	for i := range edges {
		edges[i] = D2Edge{
			From: output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
			To:   output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i+1)),
		}
	}

	return edges
}

func generateBenchmarkD2Tables(n int) []D2Table {
	tables := make([]D2Table, n)
	for i := range tables {
		tables[i] = D2Table{
			Name: fmt.Sprintf("table%d", i),
			Columns: []D2Column{
				{Name: "id", Type: "int"},
				{Name: "name", Type: "string"},
				{Name: "email", Type: "string"},
			},
		}
	}

	return tables
}

func benchmarkD2Diagram(b *testing.B, setup func(*D2Diagram)) {
	d := NewD2Diagram()
	setup(d)

	b.ResetTimer()

	for b.Loop() {
		_, _ = d.Render()
	}
}

func BenchmarkD2DiagramEmpty(b *testing.B) {
	benchmarkD2Diagram(b, func(_ *D2Diagram) {})
}

func BenchmarkD2DiagramNodes(b *testing.B) {
	benchmarkD2Diagram(b, func(d *D2Diagram) {
		for _, node := range generateBenchmarkD2Nodes(100) {
			d.AddNode(node)
		}
	})
}

func BenchmarkD2DiagramEdges(b *testing.B) {
	benchmarkD2Diagram(b, func(d *D2Diagram) {
		for _, node := range generateBenchmarkD2Nodes(100) {
			d.AddNode(node)
		}

		for _, edge := range generateBenchmarkD2Edges(99) {
			d.AddEdge(edge)
		}
	})
}

func BenchmarkD2DiagramTables(b *testing.B) {
	benchmarkD2Diagram(b, func(d *D2Diagram) {
		for _, table := range generateBenchmarkD2Tables(50) {
			d.AddTable(table.Name, table.Columns)
		}
	})
}

func BenchmarkD2DiagramStyledNodes(b *testing.B) {
	benchmarkD2Diagram(b, func(d *D2Diagram) {
		for i := range 100 {
			d.AddNode(D2Node{
				ID:    output.NewBrandedID[output.D2NodeIDBrand](fmt.Sprintf("node%d", i)),
				Label: output.NewBrandedID[output.D2NodeLabelBrand](fmt.Sprintf("Node %d", i)),
				Shape: D2ShapeCircle,
				Style: D2NodeStyle{
					Fill:     "#f0f0f0",
					Stroke:   "#333333",
					FontSize: 14,
					Shadow:  true,
				},
				Tooltip: fmt.Sprintf("Tooltip for node %d", i),
				Link:    fmt.Sprintf("https://example.com/node/%d", i),
			})
		}
	})
}

func BenchmarkD2DiagramFullConfig(b *testing.B) {
	benchmarkD2Diagram(b, func(d *D2Diagram) {
		d.SetDirection(D2DirRight).
			SetTitle("Benchmark Diagram").
			SetLayout("elk")

		d.AddClass("highlight", D2NodeStyle{Fill: "#ffcc00", Stroke: "#ff9900"})

		for _, node := range generateBenchmarkD2Nodes(50) {
			d.AddNode(node)
		}

		for _, edge := range generateBenchmarkD2Edges(49) {
			d.AddEdge(edge)
		}

		for _, table := range generateBenchmarkD2Tables(10) {
			d.AddTable(table.Name, table.Columns)
		}
	})
}
