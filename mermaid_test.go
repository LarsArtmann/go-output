package output

import (
	"strings"
	"testing"
)

func TestMermaidRenderer(t *testing.T) {
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

func TestMermaidRendererWithDiamond(t *testing.T) {
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
