package output

import (
	"strings"
	"testing"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRenderer(t *testing.T) {
	t.Parallel()
	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "A", Label: "Node A"},
		{ID: "B", Label: "Node B"},
		{ID: "C", Label: "Node C"},
	})
	renderer.SetEdges([]GraphEdge{
		{From: "A", To: "B"},
		{From: "B", To: "C"},
	})

	output := renderer.Render()

	if !strings.Contains(output, "```mermaid") {
		t.Error("Output should contain mermaid code fence")
	}
	if !strings.Contains(output, "flowchart TD") {
		t.Error("Output should contain flowchart declaration")
	}
	// Mermaid output format is: "    A[Node A]\n"
	if !strings.Contains(output, "A[Node A]") {
		t.Error("Output should contain node A with label")
	}
	if !strings.Contains(output, "A --> B") {
		t.Error("Output should contain edge A --> B")
	}
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererWithDiamond(t *testing.T) {
	t.Parallel()
	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "decision", Label: "Decision", Shape: ShapeDiamond},
	})
	renderer.SetEdges([]GraphEdge{})

	output := renderer.Render()

	// Diamond uses {} syntax: decision{Decision}
	if !strings.Contains(output, "decision{Decision}") {
		t.Error("Diamond shape should use {} syntax")
	}
}

func TestMermaidRendererFromTableData(t *testing.T) {
	t.Parallel()
	data := NewTableData([]string{"Step", "Description"})
	data.AddRow([]string{"Start", "Begin process"})
	data.AddRow([]string{"Step 1", "Do something"})
	data.AddRow([]string{"End", "Finish"})

	renderer := MermaidFlowchartRenderer(data)
	output := renderer.Render()

	if !strings.Contains(output, "flowchart TD") {
		t.Error("Output should be a flowchart")
	}
	if !strings.Contains(output, "Start") {
		t.Error("Output should contain 'Start'")
	}
	if !strings.Contains(output, "End") {
		t.Error("Output should contain 'End'")
	}
}

func TestMermaidTreeRenderer(t *testing.T) {
	t.Parallel()
	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child1", "Child 1"))
	root.AddChild(NewTreeNode("child2", "Child 2"))

	renderer := MermaidTreeRenderer(root)
	output := renderer.Render()

	if !strings.Contains(output, "flowchart TD") {
		t.Error("Output should be a flowchart")
	}
	if !strings.Contains(output, "Child 1") {
		t.Error("Output should contain 'Child 1'")
	}
	if !strings.Contains(output, "Child 2") {
		t.Error("Output should contain 'Child 2'")
	}
}

func TestMermaidRendererEmpty(t *testing.T) {
	t.Parallel()
	renderer := NewMermaidRenderer()
	output := renderer.Render()

	if !strings.Contains(output, "```mermaid") {
		t.Error("Empty mermaid should still have fence")
	}
	if !strings.Contains(output, "flowchart TD") {
		t.Error("Empty mermaid should still have flowchart declaration")
	}
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

	for _, tt := range tests {
		got := sanitizeMermaidID(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeMermaidID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
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

	for _, tt := range tests {
		got := sanitizeMermaidLabel(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeMermaidLabel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
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
			renderer.SetNodes([]GraphNode{{ID: "n", Label: "Test", Shape: tt.shape}})
			renderer.SetEdges([]GraphEdge{})

			output := renderer.Render()
			if !strings.Contains(output, tt.wantLeft) || !strings.Contains(output, tt.wantRight) {
				t.Errorf("Shape %v should produce %q...%q, got: %s", tt.shape, tt.wantLeft, tt.wantRight, output)
			}
		})
	}
}

func TestMermaidRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()
	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "A", Label: "Node A"},
		{ID: "B", Label: "Node B"},
	})
	renderer.SetEdges([]GraphEdge{
		{From: "A", To: "B", Label: "connects"},
	})

	output := renderer.Render()
	if !strings.Contains(output, "|connects|") {
		t.Error("Output should contain edge label |connects|")
	}
}

func TestMermaidTreeRendererNilRoot(t *testing.T) {
	t.Parallel()
	renderer := MermaidTreeRenderer(nil)
	output := renderer.Render()

	if !strings.Contains(output, "flowchart TD") {
		t.Error("Nil root should still produce valid flowchart")
	}
}

func TestMermaidTreeRendererWithEmptyID(t *testing.T) {
	t.Parallel()
	// TreeNode with empty ID should use label
	root := NewTreeNode("", "RootLabel")
	renderer := MermaidTreeRenderer(root)
	output := renderer.Render()

	if !strings.Contains(output, "RootLabel") {
		t.Error("Output should contain label when ID is empty")
	}
}

func TestMermaidRendererEscapeLabel(t *testing.T) {
	t.Parallel()
	renderer := NewMermaidRenderer()
	renderer.SetNodes([]GraphNode{
		{ID: "A", Label: `test "quoted" text`},
	})
	renderer.SetEdges([]GraphEdge{})

	output := renderer.Render()
	if strings.Contains(output, `"quoted"`) {
		t.Error("Quotes should be escaped")
	}
	if !strings.Contains(output, "'quoted'") {
		t.Error("Quotes should be replaced with single quotes")
	}
}
