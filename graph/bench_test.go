package graph

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func generateBenchmarkNodes(n int) []output.GraphNode {
	nodes := make([]output.GraphNode, n)
	for i := range nodes {
		nodes[i] = newTestNode("node", "Node")
	}

	return nodes
}

func generateBenchmarkEdges(n int) []output.GraphEdge {
	edges := make([]output.GraphEdge, n)
	for i := range edges {
		edges[i] = output.GraphEdge{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("node"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("node"),
		}
	}

	return edges
}

func benchmarkGraphRenderer(b *testing.B, renderer output.GraphRenderer) {
	nodes := generateBenchmarkNodes(100)
	renderer.SetNodes(nodes)

	edges := generateBenchmarkEdges(99)
	renderer.SetEdges(edges)

	b.ResetTimer()

	for b.Loop() {
		_, _ = renderer.Render()
	}
}

func BenchmarkDOTRenderer(b *testing.B) {
	renderer := NewDOTRenderer()
	benchmarkGraphRenderer(b, renderer)
}

func BenchmarkMermaidRenderer(b *testing.B) {
	renderer := NewMermaidRenderer()
	benchmarkGraphRenderer(b, renderer)
}
