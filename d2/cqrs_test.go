package d2

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestWriteGraph_HappyPath(t *testing.T) {
	t.Parallel()

	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("design", "Design")).
		AddNode(*output.NewGraphNode("implement", "Implement")).
		AddEdge(*output.NewGraphEdge("design", "implement")).
		Build()

	var buf strings.Builder
	if err := WriteGraph(&buf, g); err != nil {
		t.Fatalf("WriteGraph should return nil on success, got: %v", err)
	}

	if !strings.Contains(buf.String(), "Design") {
		t.Error("expected node label 'Design' in output")
	}

	if !strings.Contains(buf.String(), "Implement") {
		t.Error("expected node label 'Implement' in output")
	}

	if !strings.Contains(buf.String(), "design") {
		t.Error("expected edge source 'design' in output")
	}
}

func TestWrite_HappyPath(t *testing.T) {
	t.Parallel()

	diagram := NewDiagram()
	diagram.SetNodes([]output.GraphNode{
		*output.NewGraphNode("a", "Alpha"),
		*output.NewGraphNode("b", "Beta"),
	})
	diagram.SetEdges([]output.GraphEdge{
		*output.NewGraphEdge("a", "b"),
	})

	var buf strings.Builder
	if err := Write(&buf, diagram); err != nil {
		t.Fatalf("Write should return nil on success, got: %v", err)
	}

	if !strings.Contains(buf.String(), "Alpha") {
		t.Error("expected node label 'Alpha' in output")
	}
}
