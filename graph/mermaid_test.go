package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes(testNodesABC())
	renderer.SetEdges(testEdgesABC())

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "```mermaid", "Output should contain mermaid code fence")
	assertContains(t, out, "flowchart TD", "Output should contain flowchart declaration")
	// Mermaid output format is: "    A[Node A]\n"
	assertContains(t, out, "A[Node A]", "Output should contain node A with label")
	assertContains(t, out, "A --> B", "Output should contain edge A --> B")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererWithDiamond(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes(
		[]output.GraphNode{newTestNodeWithShape("decision", "Decision", output.ShapeDiamond)},
	)
	renderer.SetEdges([]output.GraphEdge{})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Diamond uses {} syntax: decision{Decision}
	assertContains(t, out, "decision{Decision}", "Diamond shape should use {} syntax")
}

func TestMermaidRendererFromTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Step", "Description"})
	data.AddRow([]string{"Start", "Begin process"})
	data.AddRow([]string{"Step 1", "Do something"})
	data.AddRow([]string{"End", "Finish"})

	renderer := MermaidFlowchartRenderer(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "flowchart TD", "Output should be a flowchart")
	assertContains(t, out, "Start", "Output should contain 'Start'")
	assertContains(t, out, "End", "Output should contain 'End'")
}

func TestMermaidTreeRenderer(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child1", "Child 1"))
	root.AddChild(output.NewTreeNode("child2", "Child 2"))

	renderer := MermaidTreeRenderer(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "flowchart TD", "Output should be a flowchart")
	assertContains(t, out, "Child 1", "Output should contain 'Child 1'")
	assertContains(t, out, "Child 2", "Output should contain 'Child 2'")
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
		shape     output.GraphShape
		wantLeft  string
		wantRight string
	}{
		{"Box", output.ShapeBox, "[", "]"},
		{"Rect", output.ShapeRect, "[", "]"},
		{"Diamond", output.ShapeDiamond, "{", "}"},
		{"Ellipse", output.ShapeEllipse, "(", ")"},
		{"Circle", output.ShapeCircle, "((", "))"},
		{"Hexagon", output.ShapeHexagon, "{{", "}}"},
		{"Cylinder", output.ShapeCylinder, "[(", ")]"},
		{"Parallelogram", output.ShapeParallelogram, "[/", "/]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			renderer := NewMermaidRenderer()
			renderer.SetNodes([]output.GraphNode{newTestNodeWithShape("n", "Test", tt.shape)})
			renderer.SetEdges([]output.GraphEdge{})

			out, err := renderer.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if !strings.Contains(out, tt.wantLeft) || !strings.Contains(out, tt.wantRight) {
				t.Errorf(
					"Shape %v should produce %q...%q, got: %s",
					tt.shape,
					tt.wantLeft,
					tt.wantRight,
					out,
				)
			}
		})
	}
}

func TestMermaidRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes(testNodesAB())
	renderer.SetEdges([]output.GraphEdge{testEdgeAB("connects")})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "|connects|", "Output should contain edge label |connects|")
}

func TestMermaidTreeRendererNilRoot(t *testing.T) {
	t.Parallel()

	renderer := MermaidTreeRenderer(nil)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "flowchart TD", "Nil root should still produce valid flowchart")
}

func TestMermaidTreeRendererWithEmptyID(t *testing.T) {
	t.Parallel()
	// TreeNode with empty ID should use label
	root := output.NewTreeNode("", "RootLabel")
	renderer := MermaidTreeRenderer(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "RootLabel", "Output should contain label when ID is empty")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererEscapeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes([]output.GraphNode{newTestNode("A", `test "quoted" text`)})
	renderer.SetEdges([]output.GraphEdge{})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(out, `"quoted"`) {
		t.Error("Quotes should be escaped")
	}

	assertContains(t, out, "'quoted'", "Quotes should be replaced with single quotes")
}
