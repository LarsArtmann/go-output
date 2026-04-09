package output

import (
	"strings"
	"testing"
)

// ExpectedOutput contains a substring to check and its corresponding error message.
type ExpectedOutput struct {
	Substring string
	Message   string
}

// assertContains checks that output contains substr, failing with msg if not.
func assertContains(t *testing.T, output, substr, msg string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Error(msg)
	}
}

// testNodesAB returns a slice of GraphNode with nodes A and B for testing.
func testNodesAB() []GraphNode {
	return []GraphNode{
		{
			ID:    NewBrandedID[GraphNodeIDBrand]("A"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Node A"),
		},
		{
			ID:    NewBrandedID[GraphNodeIDBrand]("B"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Node B"),
		},
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

// testEmptyRendererOutput verifies that an empty renderer produces valid output structure.
func testEmptyRendererOutput(t *testing.T, renderer Renderer, expectedOutputs []ExpectedOutput) {
	t.Helper()

	output := renderer.Render()
	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected.Substring) {
			t.Error(expected.Message)
		}
	}
}

// testSanitizeFunc is a shared helper for testing sanitization functions.
func testSanitizeFunc(
	t *testing.T,
	name string,
	fn func(string) string,
	tests []struct{ input, want string },
) {
	t.Helper()

	for _, tt := range tests {
		got := fn(tt.input)
		if got != tt.want {
			t.Errorf("%s(%q) = %q, want %q", name, tt.input, got, tt.want)
		}
	}
}
