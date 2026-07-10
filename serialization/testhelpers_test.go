package serialization

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

var assertContains = testhelpers.AssertContains

var assertAllContained = testhelpers.AssertAllContained

var assertOutputContains = testhelpers.AssertOutputContains

var assertMarshalError = testhelpers.AssertMarshalError

type errorWriter = testhelpers.ErrorWriter

var newTestNodeWithShape = graphtest.NewTestNodeWithShape

var newTestNode = graphtest.NewTestNode

var testNodesAB = graphtest.TestNodesAB

var testNodesABC = graphtest.TestNodesABC

var testEdgeAB = graphtest.TestEdgeAB

var testEdgesAB = graphtest.TestEdgesAB

var testEdgesABC = graphtest.TestEdgesABC

func assertValidJSON(t *testing.T, output string) {
	t.Helper()

	if !jsontext.Value([]byte(output)).IsValid() {
		t.Errorf("output should be valid JSON, got %q", output)
	}
}

func assertValidYAML(t *testing.T, output string) {
	t.Helper()

	if output == "" {
		t.Error("output should not be empty for valid YAML check")
	}
}

func testGraphRendererNodeWithShape(t *testing.T, r output.GraphRenderer, wantShape string) {
	t.Helper()

	r.SetNodes([]output.GraphNode{newTestNodeWithShape("A", "Node A", output.NodeShapeDiamond)})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, wantShape)
}

// renderWithABNodes loads the canonical A/B node+edge fixtures into the
// renderer and returns its output. Centralises the "SetNodes(testNodesAB()) +
// SetEdges(testEdgesAB()) + Render() + t.Fatalf" sequence shared by every
// "render a graph with two nodes" test in this package.
func renderWithABNodes(t *testing.T, r output.GraphRenderer) string {
	t.Helper()

	r.SetNodes(testNodesAB())
	r.SetEdges(testEdgesAB())

	out, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	return out
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
