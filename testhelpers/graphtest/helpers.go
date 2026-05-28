package graphtest

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

// NewTestNode creates a GraphNode with the given ID and label for testing.
func NewTestNode(id, label string) output.GraphNode {
	return output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

// NewTestNodeWithShape creates a GraphNode with shape for testing.
func NewTestNodeWithShape(id, label string, shape output.GraphShape) output.GraphNode {
	return output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
		Shape: shape,
	}
}

// TestNodesAB returns a slice of GraphNode with nodes A and B for testing.
func TestNodesAB() []output.GraphNode {
	return []output.GraphNode{
		NewTestNode("A", "Node A"),
		NewTestNode("B", "Node B"),
	}
}

// TestNodesABC returns a slice of GraphNode with nodes A, B, and C.
func TestNodesABC() []output.GraphNode {
	nodes := TestNodesAB()

	return append(nodes, NewTestNode("C", "Node C"))
}

// TestEdgeAB returns a GraphEdge connecting A to B with the given label.
func TestEdgeAB(label string) output.GraphEdge {
	return output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("A"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("B"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

// TestEdgesAB returns a slice with a single edge from A to B.
func TestEdgesAB() []output.GraphEdge {
	return []output.GraphEdge{TestEdgeAB("")}
}

// TestEdgesABC returns edges A→B and B→C.
func TestEdgesABC() []output.GraphEdge {
	return []output.GraphEdge{
		TestEdgeAB(""),
		{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("B"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("C"),
		},
	}
}

// AssertEscape checks that if input contains a character, the escaped output has it backslash-escaped.
func AssertEscape(t *testing.T, fnName, input, escaped, char, desc string) {
	t.Helper()

	if strings.Contains(input, char) && !strings.Contains(escaped, `\`+char) {
		t.Errorf("%s(%q) = %q, %s not escaped", fnName, input, escaped, desc)
	}
}
