package output

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
	"github.com/larsartmann/go-output/testhelpers"
)

// Re-export generic helpers for use by package output tests.
type (
	ExpectedOutput         = gentest.ExpectedOutput
	htmlEscapeTestRenderer = gentest.HTMLEscapeTestRenderer
)

//nolint:gochecknoglobals // Re-exported test helpers for package-local use
var (
	assertContains     = testhelpers.AssertContains
	assertMarshalError = gentest.AssertMarshalError
)

// Re-export shared test helpers from testhelpers package.
//
//nolint:gochecknoglobals // Re-exported test helper for package-local use
var assertStringSliceEqual = testhelpers.AssertStringSliceEqual

func testHTMLEscapeShared(
	t *testing.T,
	newRenderer func() gentest.HTMLEscapeTestRenderer,
	name string,
) {
	gentest.AssertHTMLEscape(t, newRenderer, name)
}

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

// newTestNodeWithShape creates a GraphNode with the given ID, label, and shape for testing.
func newTestNodeWithShape(id, label string, shape GraphShape) GraphNode {
	return GraphNode{
		ID:    NewBrandedID[GraphNodeIDBrand](id),
		Label: NewBrandedID[GraphNodeLabelBrand](label),
		Shape: shape,
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

// testEdgesABC returns edges A->B and B->C.
func testEdgesABC() []GraphEdge {
	return []GraphEdge{
		testEdgeAB(""),
		{From: NewBrandedID[GraphNodeIDBrand]("B"), To: NewBrandedID[GraphNodeIDBrand]("C")},
	}
}

// testEmptyRendererOutput verifies that an empty renderer produces valid output structure.
func testEmptyRendererOutput(
	t *testing.T,
	renderer Renderer,
	expectedOutputs []gentest.ExpectedOutput,
) {
	t.Helper()

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected.Substring) {
			t.Error(expected.Message)
		}
	}
}

// AssertTreeNodeDepth verifies the depth of tree nodes in a hierarchy.
func AssertTreeNodeDepth(t *testing.T, root, child, grandchild *TreeNode) {
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

// testExpectedOutputs builds a []ExpectedOutput from alternating substring/message pairs.
func testExpectedOutputs(pairs ...string) []ExpectedOutput {
	if len(pairs)%2 != 0 {
		panic("testExpectedOutputs requires even number of arguments")
	}

	out := make([]ExpectedOutput, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		// gosec G602: Safe - i increments by 2, and len(pairs) is guaranteed even
		out = append(out, ExpectedOutput{Substring: pairs[i], Message: pairs[i+1]}) //nolint:gosec
	}

	return out
}

// testHTMLEmptyExpected returns the expected substrings for an empty HTML renderer.
func testHTMLEmptyExpected() []ExpectedOutput {
	return testExpectedOutputs(
		"<table", "Empty table should still be valid HTML",
		"</table>", "Empty table should have closing tag",
	)
}
