package serialization

import (
	"encoding/json"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var assertContains = testhelpers.AssertContains

var assertOutputContains = testhelpers.AssertOutputContains

var assertMarshalError = testhelpers.AssertMarshalError

type errorWriter = testhelpers.ErrorWriter

func assertValidJSON(t *testing.T, output string) {
	t.Helper()

	if !json.Valid([]byte(output)) {
		t.Errorf("output should be valid JSON, got %q", output)
	}
}

func assertValidYAML(t *testing.T, output string) {
	t.Helper()

	if output == "" {
		t.Error("output should not be empty for valid YAML check")
	}
}

func testNodesAB() []output.GraphNode {
	return []output.GraphNode{
		newTestNode("A", "Node A"),
		newTestNode("B", "Node B"),
	}
}

func newTestNode(id, label string) output.GraphNode {
	return output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

func newTestNodeWithShape(id, label string, shape output.GraphShape) output.GraphNode {
	return output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
		Shape: shape,
	}
}

func testEdgeAB(label string) output.GraphEdge {
	return output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("A"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("B"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
	}
}

func testEdgesAB() []output.GraphEdge {
	return []output.GraphEdge{testEdgeAB("")}
}

func testNodesABC() []output.GraphNode {
	nodes := testNodesAB()
	return append(nodes, newTestNode("C", "Node C"))
}

func testEdgesABC() []output.GraphEdge {
	return []output.GraphEdge{
		testEdgeAB(""),
		{
			From: output.NewBrandedID[output.GraphNodeIDBrand]("B"),
			To:   output.NewBrandedID[output.GraphNodeIDBrand]("C"),
		},
	}
}

type graphRenderer interface {
	output.Renderer
	SetNodes([]output.GraphNode)
	SetEdges([]output.GraphEdge)
}

func testGraphRendererNodeWithShape(t *testing.T, r graphRenderer, wantShape string) {
	t.Helper()

	r.SetNodes([]output.GraphNode{newTestNodeWithShape("A", "Node A", output.ShapeDiamond)})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, wantShape)
}

type unmarshalTestCase struct {
	name    string
	data    string
	wantErr bool
}

func testUnmarshalCases(
	t *testing.T,
	tests []unmarshalTestCase,
	unmarshal func([]byte, any) error,
	funcName string,
) {
	t.Helper()

	for _, tt := range tests {
		testUnmarshalError(t, tt.name, tt.data, tt.wantErr, unmarshal, funcName)
	}
}

func testUnmarshalError(
	t *testing.T,
	name, data string,
	wantErr bool,
	unmarshal func([]byte, any) error,
	funcName string,
) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		var got any

		err := unmarshal([]byte(data), &got)
		if (err != nil) != wantErr {
			t.Errorf("%s() error = %v, wantErr %v", funcName, err, wantErr)
		}
	})
}
