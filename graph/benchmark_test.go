package graph

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
)

//nolint:goconst // benchmark constants are fine to repeat
func BenchmarkDOTRenderer_100Nodes150Edges(b *testing.B) {
	renderer := NewDOTRenderer()

	for i := range 100 {
		renderer.AddNode(*output.NewGraphNode(fmt.Sprintf("node%d", i), fmt.Sprintf("Node %d", i)))
	}

	for i := range 150 {
		renderer.AddEdge(*output.NewGraphEdge(fmt.Sprintf("node%d", i%100), fmt.Sprintf("node%d", (i+7)%100)))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = renderer.Render()
	}
}

//nolint:goconst // benchmark constants are fine to repeat
func BenchmarkMermaidRenderer_100Nodes150Edges(b *testing.B) {
	renderer := NewMermaidRenderer()

	for i := range 100 {
		renderer.AddNode(*output.NewGraphNode(fmt.Sprintf("node%d", i), fmt.Sprintf("Node %d", i)))
	}

	for i := range 150 {
		renderer.AddEdge(*output.NewGraphEdge(fmt.Sprintf("node%d", i%100), fmt.Sprintf("node%d", (i+7)%100)))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = renderer.Render()
	}
}
