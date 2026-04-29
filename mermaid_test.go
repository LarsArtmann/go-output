package output

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/escape"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes(testNodesABC())
	renderer.SetEdges(testEdgesABC())

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "```mermaid", "Output should contain mermaid code fence")
	assertContains(t, output, "flowchart TD", "Output should contain flowchart declaration")
	// Mermaid output format is: "    A[Node A]\n"
	assertContains(t, output, "A[Node A]", "Output should contain node A with label")
	assertContains(t, output, "A --> B", "Output should contain edge A --> B")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererWithDiamond(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{newTestNodeWithShape("decision", "Decision", ShapeDiamond)})
	renderer.SetEdges([]GraphEdge{})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Diamond uses {} syntax: decision{Decision}
	assertContains(t, output, "decision{Decision}", "Diamond shape should use {} syntax")
}

func TestMermaidRendererFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Step", "Description"})
	data.AddRow([]string{"Start", "Begin process"})
	data.AddRow([]string{"Step 1", "Do something"})
	data.AddRow([]string{"End", "Finish"})

	renderer := MermaidFlowchartRenderer(data)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "flowchart TD", "Output should be a flowchart")
	assertContains(t, output, "Start", "Output should contain 'Start'")
	assertContains(t, output, "End", "Output should contain 'End'")
}

func TestMermaidTreeRenderer(t *testing.T) {
	t.Parallel()

	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child1", "Child 1"))
	root.AddChild(NewTreeNode("child2", "Child 2"))

	renderer := MermaidTreeRenderer(root)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "flowchart TD", "Output should be a flowchart")
	assertContains(t, output, "Child 1", "Output should contain 'Child 1'")
	assertContains(t, output, "Child 2", "Output should contain 'Child 2'")
}

func TestMermaidRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	testEmptyRendererOutput(t, renderer, testMermaidEmptyExpected())
}

func TestSanitizeMermaidID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"valid_id", "valid_id"},
		{"has spaces", "hasspaces"},
		{"has-dashes", "hasdashes"},
		{"special!@#$chars", "specialchars"},
		{"123numbers", "123numbers"},
		{"", "node"},
	}

	testSanitizeFunc(t, "escape.MermaidID", escape.MermaidID, tests)
}

func TestSanitizeMermaidLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"has spaces", "has_spaces"},
		{"has-dash", "has_dash"},
		{"path/to/file", "path_to_file"},
		{"multi word test", "multi_word_test"},
	}

	testSanitizeFunc(t, "escape.MermaidSlug", escape.MermaidSlug, tests)
}

func TestMermaidRendererAllShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		shape     GraphShape
		wantLeft  string
		wantRight string
	}{
		{"Box", ShapeBox, "[", "]"},
		{"Rect", ShapeRect, "[", "]"},
		{"Diamond", ShapeDiamond, "{", "}"},
		{"Ellipse", ShapeEllipse, "(", ")"},
		{"Circle", ShapeCircle, "((", "))"},
		{"Hexagon", ShapeHexagon, "{{", "}}"},
		{"Cylinder", ShapeCylinder, "[(", ")]"},
		{"Parallelogram", ShapeParallelogram, "[/", "/]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			renderer := NewMermaidRenderer()
			renderer.SetNodes([]GraphNode{newTestNodeWithShape("n", "Test", tt.shape)})
			renderer.SetEdges([]GraphEdge{})

			output, err := renderer.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if !strings.Contains(output, tt.wantLeft) || !strings.Contains(output, tt.wantRight) {
				t.Errorf(
					"Shape %v should produce %q...%q, got: %s",
					tt.shape,
					tt.wantLeft,
					tt.wantRight,
					output,
				)
			}
		})
	}
}

func TestMermaidRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges([]GraphEdge{testEdgeAB("connects")})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "|connects|", "Output should contain edge label |connects|")
}

func TestMermaidTreeRendererNilRoot(t *testing.T) {
	t.Parallel()

	renderer := MermaidTreeRenderer(nil)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "flowchart TD", "Nil root should still produce valid flowchart")
}

func TestMermaidTreeRendererWithEmptyID(t *testing.T) {
	t.Parallel()
	// TreeNode with empty ID should use label
	root := NewTreeNode("", "RootLabel")
	renderer := MermaidTreeRenderer(root)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "RootLabel", "Output should contain label when ID is empty")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererEscapeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{newTestNode("A", `test "quoted" text`)})
	renderer.SetEdges([]GraphEdge{})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(output, `"quoted"`) {
		t.Error("Quotes should be escaped")
	}

	assertContains(t, output, "'quoted'", "Quotes should be replaced with single quotes")
}
