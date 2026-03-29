package testutils

import (
	"github.com/larsartmann/go-output"
)

// CreateTestNodesAB creates the common test nodes A and B used in DOT and Mermaid tests
func CreateTestNodesAB() []output.GraphNode {
	return []output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node A"),
		},
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("B"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node B"),
		},
	}
}

// CreateTestEdgeAB creates the common test edge from A to B used in DOT and Mermaid tests
func CreateTestEdgeAB() []output.GraphEdge {
	return []output.GraphEdge{
		{From: output.NewBrandedID[output.GraphNodeIDBrand]("A"), To: output.NewBrandedID[output.GraphNodeIDBrand]("B")},
	}
}

// CreateTestNodesABC creates the common test nodes A, B, and C used in Mermaid tests
func CreateTestNodesABC() []output.GraphNode {
	return []output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node A"),
		},
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("B"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node B"),
		},
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("C"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node C"),
		},
	}
}

// CreateTestEdgesABC creates the common test edges A->B and B->C used in Mermaid tests
func CreateTestEdgesABC() []output.GraphEdge {
	return []output.GraphEdge{
		{From: output.NewBrandedID[output.GraphNodeIDBrand]("A"), To: output.NewBrandedID[output.GraphNodeIDBrand]("B")},
		{From: output.NewBrandedID[output.GraphNodeIDBrand]("B"), To: output.NewBrandedID[output.GraphNodeIDBrand]("C")},
	}
}