package serialization

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var assertContains = testhelpers.AssertContains

type errorWriter struct{}

var errWrite = errors.New("write error")

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errWrite
}

var _ io.Writer = (*errorWriter)(nil)

func assertOutputContains(t *testing.T, output, substr string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Errorf("output should contain %q, got %q", substr, output)
	}
}

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

func assertMarshalError(t *testing.T, name string, err error, wantErr bool) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("%s() error = %v, wantErr %v", name, err, wantErr)
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
