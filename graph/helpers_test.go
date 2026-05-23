package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

type expectedOutput struct {
	Substring string
	Message   string
}

//nolint:gochecknoglobals // Re-exported test helper for package-local use
var assertContains = testhelpers.AssertContains

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

func testNodesAB() []output.GraphNode {
	return []output.GraphNode{
		newTestNode("A", "Node A"),
		newTestNode("B", "Node B"),
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

func testExpectedOutputs(pairs ...string) []expectedOutput {
	if len(pairs)%2 != 0 {
		panic("testExpectedOutputs requires even number of arguments")
	}

	out := make([]expectedOutput, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, expectedOutput{Substring: pairs[i], Message: pairs[i+1]}) //nolint:gosec
	}

	return out
}

func testDOTEmptyExpected() []expectedOutput {
	return testExpectedOutputs(
		"digraph G {", "Empty DOT should still have digraph declaration",
		"rankdir=TB", "Empty DOT should have default attributes",
	)
}

func testMermaidEmptyExpected() []expectedOutput {
	return testExpectedOutputs(
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

type parseEnumTestCase[T any] struct {
	name    string
	input   string
	want    T
	wantErr bool
}

func testParseEnum[T any](
	t *testing.T,
	name string,
	parseFunc func(string) (T, error),
	testCases []parseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				got, err := parseFunc(testCase.input)

				if (err != nil) != testCase.wantErr {
					t.Errorf("%s() error = %v, wantErr %v", name, err, testCase.wantErr)

					return
				}

				if !equalFunc(got, testCase.want) {
					t.Errorf("%s() = %v, want %v", name, got, testCase.want)
				}
			})
		}
	})
}

type stringEnumTestCase[T any] struct {
	value T
	want  string
}

func testEnumString[T any](
	t *testing.T,
	name string,
	testCases []stringEnumTestCase[T],
	stringFunc func(T) string,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			t.Run(testCase.want, func(t *testing.T) {
				t.Parallel()

				if got := stringFunc(testCase.value); got != testCase.want {
					t.Errorf("%s() = %v, want %v", name, got, testCase.want)
				}
			})
		}
	})
}

func testAllowedValues(
	t *testing.T,
	name string,
	got []string,
	want []string,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		testhelpers.AssertStringSliceEqual(t, name, got, want)
	})
}
