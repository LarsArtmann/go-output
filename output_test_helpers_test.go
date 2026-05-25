package output

import (
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
	"github.com/larsartmann/go-output/testhelpers"
)

// Re-export generic helpers for use by package output tests.
type ExpectedOutput = gentest.ExpectedOutput

//nolint:gochecknoglobals // Re-exported test helpers for package-local use
var assertContains = testhelpers.AssertContains

// testNodesAB returns a slice of GraphNode with nodes A and B for testing.
func testNodesAB() []GraphNode {
	return []GraphNode{
		newTestNode("A", "Node A"),
		newTestNode("B", "Node B"),
	}
}

// newTestNode creates a GraphNode with the given ID and label for testing.
func newTestNode(id, label string) GraphNode {
	return GraphNode{
		ID:    NewBrandedID[GraphNodeIDBrand](id),
		Label: NewBrandedID[GraphNodeLabelBrand](label),
	}
}

// testEdgeAB returns a GraphEdge connecting A to B with the given label for testing.
func testEdgeAB(label string) GraphEdge {
	return GraphEdge{
		From:  NewBrandedID[GraphNodeIDBrand]("A"),
		To:    NewBrandedID[GraphNodeIDBrand]("B"),
		Label: NewBrandedID[GraphNodeLabelBrand](label),
	}
}

// testEdgesAB returns a slice with a single edge from A to B.
func testEdgesAB() []GraphEdge {
	return []GraphEdge{testEdgeAB("")}
}

// testNodesABC returns a slice of GraphNode with nodes A, B, and C.
func testNodesABC() []GraphNode {
	nodes := testNodesAB()
	return append(nodes, newTestNode("C", "Node C"))
}

// assertTreeNodeDepth verifies the depth of tree nodes in a hierarchy.
func assertTreeNodeDepth(t testing.TB, root, child, grandchild *TreeNode) {
	t.Helper()

	if root.Depth() != 0 {
		t.Errorf("Root depth should be 0, got %d", root.Depth())
	}

	if child.Depth() != 1 {
		t.Errorf("Child depth should be 1, got %d", child.Depth())
	}

	if grandchild.Depth() != 2 {
		t.Errorf("Grandchild depth should be 2, got %d", grandchild.Depth())
	}
}
