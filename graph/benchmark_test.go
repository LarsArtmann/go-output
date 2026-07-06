package graph

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
)

func BenchmarkDOTRenderer_100Nodes150Edges(b *testing.B) {
	renderer := NewDOTRenderer()
	for i := 0; i < 100; i++ {
		renderer.AddNode(*output.NewGraphNode(fmt.Sprintf("node%d", i), fmt.Sprintf("Node %d", i)))
	}
	for i := 0; i < 150; i++ {
		renderer.AddEdge(*output.NewGraphEdge(fmt.Sprintf("node%d", i%100), fmt.Sprintf("node%d", (i+7)%100)))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = renderer.Render()
	}
}

func BenchmarkMermaidRenderer_100Nodes150Edges(b *testing.B) {
	renderer := NewMermaidRenderer()
	for i := 0; i < 100; i++ {
		renderer.AddNode(*output.NewGraphNode(fmt.Sprintf("node%d", i), fmt.Sprintf("Node %d", i)))
	}
	for i := 0; i < 150; i++ {
		renderer.AddEdge(*output.NewGraphEdge(fmt.Sprintf("node%d", i%100), fmt.Sprintf("node%d", (i+7)%100)))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = renderer.Render()
	}
}
