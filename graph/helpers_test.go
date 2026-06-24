package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

type expectedOutput struct {
	Substring string
	Message   string
}

//nolint:gochecknoglobals // Re-exported test helper for package-local use
var assertContains = testhelpers.AssertContains

var newTestNode = graphtest.NewTestNode

var newTestNodeWithShape = graphtest.NewTestNodeWithShape

var testNodesAB = graphtest.TestNodesAB

var testNodesABC = graphtest.TestNodesABC

var testEdgeAB = graphtest.TestEdgeAB

var testEdgesAB = graphtest.TestEdgesAB

var testEdgesABC = graphtest.TestEdgesABC

func testEmptyRendererOutput(
	t *testing.T,
	renderer output.Renderer,
	expectedOutputs []expectedOutput,
) {
	t.Helper()

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(got, expected.Substring) {
			t.Error(expected.Message)
		}
	}
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

func testExpectedOutputs(t *testing.T, pairs ...string) []expectedOutput {
	t.Helper()

	if len(pairs)%2 != 0 {
		t.Fatalf("testExpectedOutputs requires even number of arguments, got %d", len(pairs))
	}

	out := make([]expectedOutput, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, expectedOutput{Substring: pairs[i], Message: pairs[i+1]}) //nolint:gosec
	}

	return out
}

func testDOTEmptyExpected(t *testing.T) []expectedOutput {
	return testExpectedOutputs(
		t,
		"digraph G {", "Empty DOT should still have digraph declaration",
		"rankdir=TB", "Empty DOT should have default attributes",
	)
}

func testMermaidEmptyExpected(t *testing.T) []expectedOutput {
	return testExpectedOutputs(
		t,
		"```mermaid", "Empty mermaid should still have fence",
		"flowchart TD", "Empty mermaid should still have flowchart declaration",
	)
}

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
