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
		diagram.AddEdge(graphtest.NewTestEdge("node", "node", "edge"))
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = diagram.Render()
	}
}

func BenchmarkNewPlantUMLFromTable(b *testing.B) {
	data := output.NewTable([]string{"Name", "Value"})
	for range 100 {
		data.AddRow([]string{"item", "value"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = NewPlantUMLFromTable(data).Render()
	}
}
