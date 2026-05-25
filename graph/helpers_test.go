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

type parseEnumTestCase[T any] = testhelpers.ParseEnumTestCase[T]

func testParseEnum[T any](
	t *testing.T,
	name string,
	parseFunc func(string) (T, error),
	testCases []parseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	testhelpers.TestParseEnum(t, name, parseFunc, testCases, equalFunc)
}

type stringEnumTestCase[T any] = testhelpers.StringEnumTestCase[T]

func testEnumString[T any](
	t *testing.T,
	name string,
	testCases []stringEnumTestCase[T],
	stringFunc func(T) string,
) {
	testhelpers.TestEnumString(t, name, testCases, stringFunc)
}

var testAllowedValues = testhelpers.TestAllowedValues
