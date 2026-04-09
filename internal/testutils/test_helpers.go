package testutils

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

// CreateTestNode creates a GraphNode with the given ID and label.
func CreateTestNode(id, label string) output.GraphNode {
	return output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

// CreateTestNodesAB creates the common test nodes A and B used in DOT and Mermaid tests.
func CreateTestNodesAB() []output.GraphNode {
	return []output.GraphNode{
		CreateTestNode("A", "Node A"),
		CreateTestNode("B", "Node B"),
	}
}

// CreateTestEdgeAB creates the common test edge from A to B used in DOT and Mermaid tests.
func CreateTestEdgeAB() []output.GraphEdge {
	return []output.GraphEdge{
		{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("B"),
		},
	}
}

// CreateTestEdgeABWithLabel creates a test edge from A to B with the given label.
func CreateTestEdgeABWithLabel(label string) output.GraphEdge {
	return output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("A"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("B"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

// CreateTestNodesABC creates the common test nodes A, B, and C used in Mermaid tests.
func CreateTestNodesABC() []output.GraphNode {
	nodes := CreateTestNodesAB()
	return append(nodes, CreateTestNode("C", "Node C"))
}

// CreateTestEdgesABC creates the common test edges A->B and B->C used in Mermaid tests.
func CreateTestEdgesABC() []output.GraphEdge {
	return []output.GraphEdge{
		{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("B"),
		},
		{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("B"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("C"),
		},
	}
}

// AssertEmptyDataRendersJSONWithoutPanic verifies that empty data renders as JSON without panic or error.
func AssertEmptyDataRendersJSONWithoutPanic(t *testing.T) {
	data := output.NewTableData([]string{})

	_, err := output.MarshalJSON(data)
	if err != nil {
		t.Errorf("MarshalJSON on empty data should not error: %v", err)
	}
}

// RunEmptyDataRendersJSONWithoutPanic runs the "empty data renders without panic" test sub-case.
func RunEmptyDataRendersJSONWithoutPanic(t *testing.T) {
	t.Run("empty data renders without panic", func(t *testing.T) {
		t.Parallel()
		AssertEmptyDataRendersJSONWithoutPanic(t)
	})
}

// RenderMarkdownTable renders a markdown table with the given headers and rows.
func RenderMarkdownTable(headers []string, rows [][]string) string {
	md := output.NewMarkdownTable()
	md.SetHeaders(headers)

	for _, row := range rows {
		md.AddRow(row)
	}

	return md.Render()
}

// ExpectedOutput contains a substring to check and its corresponding error message.
type ExpectedOutput struct {
	Substring string
	Message   string
}

// AssertContains checks that output contains substr, failing with msg if not.
func AssertContains(t *testing.T, output, substr, msg string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Error(msg)
	}
}

// AssertEmptyRendererOutput verifies that an empty renderer produces valid output structure.
func AssertEmptyRendererOutput(
	t *testing.T,
	renderer output.Renderer,
	expectedOutputs []ExpectedOutput,
) {
	t.Helper()

	out := renderer.Render()
	for _, expected := range expectedOutputs {
		if !strings.Contains(out, expected.Substring) {
			t.Error(expected.Message)
		}
	}
}

// HTMLEscapeTestRenderer is an interface for HTML renderers that support escaping tests.
type HTMLEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() string
}

// AssertHTMLEscape verifies that a renderer properly escapes HTML content.
func AssertHTMLEscape(t *testing.T, newRenderer func() HTMLEscapeTestRenderer, name string) {
	t.Helper()

	r := newRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"<script>alert('xss')</script>"})

	got := r.Render()

	if strings.Contains(got, "<script>") {
		t.Errorf("%s: Render() should escape script tags", name)
	}

	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("%s: Render() should contain escaped script tag", name)
	}
}
