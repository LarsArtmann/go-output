package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestRenderDOT_PureFunction(t *testing.T) {
	t.Parallel()

	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("a", "Node A"))
	b.AddNode(*output.NewGraphNode("b", "Node B"))
	b.AddEdge(*output.NewGraphEdge("a", "b"))
	g := b.Build()

	got, err := RenderDOT(g)
	if err != nil {
		t.Fatalf("RenderDOT error: %v", err)
	}

	if !strings.Contains(got, "digraph") {
		t.Error("expected digraph in output")
	}

	if !strings.Contains(got, "Node A") {
		t.Error("expected node label in output")
	}

	if !strings.Contains(got, `"a" -> "b"`) {
		t.Error("expected edge in output")
	}
}

func TestWriteDOT_PureFunction(t *testing.T) {
	t.Parallel()

	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("x", "Node X"))
	g := b.Build()

	var buf strings.Builder
	if err := WriteDOT(&buf, g); err != nil {
		t.Fatalf("WriteDOT error: %v", err)
	}

	if !strings.Contains(buf.String(), "Node X") {
		t.Error("expected node label in output")
	}
}

func TestRenderMermaid_PureFunction(t *testing.T) {
	t.Parallel()

	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("a", "Node A"))
	b.AddNode(*output.NewGraphNode("b", "Node B"))
	b.AddEdge(*output.NewGraphEdge("a", "b"))
	g := b.Build()

	got, err := RenderMermaid(g, WithCodeFence(false))
	if err != nil {
		t.Fatalf("RenderMermaid error: %v", err)
	}

	if strings.Contains(got, "```mermaid") {
		t.Error("expected no code fence")
	}

	if !strings.Contains(got, "Node A") {
		t.Error("expected node label in output")
	}
}
