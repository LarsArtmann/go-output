package testutils

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/internal/gentest"
)

// Re-export generic helpers from gentest for convenience.
type (
	ExpectedOutput         = gentest.ExpectedOutput
	HTMLEscapeTestRenderer = gentest.HTMLEscapeTestRenderer
)

var (
	AssertContains   = gentest.AssertContains
	AssertHTMLEscape = gentest.AssertHTMLEscape
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

// CreateTestEdge creates a GraphEdge from one node ID to another.
func CreateTestEdge(from, to string) output.GraphEdge {
	return output.GraphEdge{
		From: output.NewBrandedID[output.GraphNodeIDBrand](from),
		To:   output.NewBrandedID[output.GraphNodeIDBrand](to),
	}
}

// CreateTestEdgeAB creates the common test edge from A to B used in DOT and Mermaid tests.
func CreateTestEdgeAB() []output.GraphEdge {
	return []output.GraphEdge{CreateTestEdge("A", "B")}
}

// CreateTestEdgeABWithLabel creates a test edge from A to B with the given label.
func CreateTestEdgeABWithLabel(label string) output.GraphEdge {
	edge := CreateTestEdge("A", "B")
	edge.Label = output.NewBrandedID[output.GraphNodeLabelBrand](label)

	return edge
}

// CreateTestNodesABC creates the common test nodes A, B, and C used in Mermaid tests.
func CreateTestNodesABC() []output.GraphNode {
	nodes := CreateTestNodesAB()
	return append(nodes, CreateTestNode("C", "Node C"))
}

// CreateTestEdgesABC creates the common test edges A->B and B->C used in Mermaid tests.
func CreateTestEdgesABC() []output.GraphEdge {
	return []output.GraphEdge{
		CreateTestEdge("A", "B"),
		CreateTestEdge("B", "C"),
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

// AssertTableData verifies that the table has the expected number of columns and rows.
func AssertTableData(t *testing.T, data *output.TableData, expectedCols, expectedRows int) {
	t.Helper()

	if data == nil {
		t.Fatal("TableData is nil")
		return
	}

	if got := data.ColCount(); got != expectedCols {
		t.Errorf("TableData has %d columns, want %d", got, expectedCols)
	}

	if got := data.RowCount(); got != expectedRows {
		t.Errorf("TableData has %d rows, want %d", got, expectedRows)
	}
}

// RenderMarkdownTable renders a markdown table with the given headers and rows.
func RenderMarkdownTable(headers []string, rows [][]string) string {
	md := output.NewMarkdownTable()
	md.SetHeaders(headers)

	for _, row := range rows {
		md.AddRow(row)
	}

	out, err := md.Render()
	if err != nil {
		return ""
	}

	return out
}

// RenderSampleMarkdownTable returns a rendered markdown table with sample project data.
func RenderSampleMarkdownTable() string {
	return RenderMarkdownTable([]string{"Name", "Health"}, [][]string{{"Alpha", "90%"}})
}

// AssertEmptyRendererOutput verifies that an empty renderer produces valid output structure.
func AssertEmptyRendererOutput(
	t *testing.T,
	renderer output.Renderer,
	expectedOutputs []ExpectedOutput,
) {
	t.Helper()

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(out, expected.Substring) {
			t.Error(expected.Message)
		}
	}
}
