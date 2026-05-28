package plantuml

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

func BenchmarkPlantUMLRender(b *testing.B) {
	diagram := NewPlantUMLDiagram()

	for range 50 {
		diagram.AddNode(graphtest.NewTestNode("node", "Node"))
	}

	for range 49 {
		diagram.AddEdge(output.GraphEdge{
			From:  output.NewBrandedID[output.GraphNodeIDBrand]("node"),
			To:    output.NewBrandedID[output.GraphNodeIDBrand]("node"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("edge"),
		})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = diagram.Render()
	}
}

func BenchmarkPlantUMLFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Value"})
	for range 100 {
		data.AddRow([]string{"item", "value"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = PlantUMLFromTableData(data).Render()
	}
}
