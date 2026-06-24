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
		[]output.GraphNode{newTestNodeWithShape("decision", "Decision", output.NodeShapeDiamond)},
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

	renderer := MermaidFromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "flowchart TD", "Output should be a flowchart")
	assertContains(t, out, "Start", "Output should contain 'Start'")
	assertContains(t, out, "End", "Output should contain 'End'")
}

func TestMermaidFromTree(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child1", "Child 1"))
	root.AddChild(output.NewTreeNode("child2", "Child 2"))

	renderer := MermaidFromTree(root)

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
	testEmptyRendererOutput(t, renderer, testMermaidEmptyExpected(t))
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
		shape     output.NodeShape
		wantLeft  string
		wantRight string
	}{
		{"Box", output.NodeShapeBox, "[", "]"},
		{"Rect", output.NodeShapeRect, "[", "]"}, //nolint:staticcheck // tests deprecated shape rendering
		{"Diamond", output.NodeShapeDiamond, "{", "}"},
		{"Ellipse", output.NodeShapeEllipse, "(", ")"},
		{"Circle", output.NodeShapeCircle, "((", "))"},
		{"Hexagon", output.NodeShapeHexagon, "{{", "}}"},
		{"Cylinder", output.NodeShapeCylinder, "[(", ")]"},
		{"Parallelogram", output.NodeShapeParallelogram, "[/", "/]"},
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

func TestMermaidFromTreeNilRoot(t *testing.T) {
	t.Parallel()

	renderer := MermaidFromTree(nil)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "flowchart TD", "Nil root should still produce valid flowchart")
}

func TestMermaidFromTableDataNil(t *testing.T) {
	t.Parallel()

	renderer := MermaidFromTableData(nil)
	if renderer == nil {
		t.Fatal("MermaidFromTableData(nil) should return non-nil renderer")
	}
}

func TestMermaidRendererNoCodeFence(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetCodeFence(false)

	out := renderWithABNodes(t, renderer)

	if strings.Contains(out, "```mermaid") {
		t.Errorf("Raw output should not contain code fence, got: %s", out)
	}

	assertContains(t, out, "flowchart TD", "Raw output should still contain flowchart declaration")
	assertContains(t, out, "A[Node A]", "Raw output should still contain nodes")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestMermaidRendererWithNodeStyle(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes([]output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Styled"),
			Shape: output.NodeShapeBox,
			Style: output.GraphStyle{
				Fill:      "#e8a838",
				Stroke:    "#4a4030",
				FontColor: "#14110d",
				FontSize:  14,
			},
		},
	})
	renderer.SetEdges([]output.GraphEdge{})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "style A fill:#e8a838,stroke:#4a4030,color:#14110d,font-size:14px",
		"Output should contain per-node style directive")
	assertContains(t, out, "% Styling", "Styling section should be present")
}

func TestMermaidRendererNoStyleNoStylingSection(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()

	out := renderWithABNodes(t, renderer)

	if strings.Contains(out, "Styling") {
		t.Errorf("Nodes without style should not emit styling section, got: %s", out)
	}
}

func TestMermaidFromTreeWithEmptyID(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("", "RootLabel")
	renderer := MermaidFromTree(root)

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

// TestMermaidNodeStyleEscapesInjection verifies that malicious style values
// (brackets, braces, quotes, newlines) are escaped through the Mermaid render
// pipeline. If escape.MermaidText were removed from mermaidStyleParts, brackets
// or newlines in a style value could break the flowchart syntax or inject
// arbitrary Mermaid directives.
func TestMermaidNodeStyleEscapesInjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{"brackets in Fill", `red[injected]`},
		{"braces in Stroke", `#000{injected}`},
		{"double quote in FontColor", `#fff"breakout`},
		{"newline in Fill", "red\nstyle other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			renderer := NewMermaidRenderer()
			renderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Test uses minimal fields
				{
					ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
					Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
					Shape: output.NodeShapeBox,
					Style: output.GraphStyle{
						Fill:      tt.value,
						Stroke:    tt.value,
						FontColor: tt.value,
					},
				},
			})
			renderer.SetEdges([]output.GraphEdge{})

			out, err := renderer.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if strings.Contains(out, tt.value) {
				t.Errorf("raw malicious value %q leaked unescaped into Mermaid output", tt.value)
			}
		})
	}
}

// TestMermaidNodeStyleEscapeOutput verifies the exact escaped sequences appear
// in Mermaid output, complementing the "raw value doesn't leak" check above.
func TestMermaidNodeStyleEscapeOutput(t *testing.T) {
	t.Parallel()

	renderer := NewMermaidRenderer()
	renderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Test uses minimal fields
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
			Shape: output.NodeShapeBox,
			Style: output.GraphStyle{
				Fill: `a"b[c]` + "\n" + `d`,
			},
		},
	})
	renderer.SetEdges([]output.GraphEdge{})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "fill:a'b(c)<br>d",
		"quotes→apostrophes, brackets→parens, newline→<br> in style values")
}
