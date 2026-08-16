package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()

	out := renderWithABNodes(t, renderer)

	assertContains(t, out, "digraph G {", "Output should contain digraph declaration")
	assertContains(t, out, "\"A\"", "Output should contain node A in quotes")
	assertContains(t, out, "label=\"Node A\"", "Output should contain label for node A")
	assertContains(t, out, "\"A\" -> \"B\"", "Output should contain directed edge A -> B")
	assertContains(t, out, "}", "Output should close with }")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTUndirectedRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewUndirectedDOTRenderer()

	out := renderWithABNodes(t, renderer)

	assertContains(t, out, "graph G {", "Undirected graph should use 'graph' keyword")
	assertContains(t, out, "\"A\" -- \"B\"", "Undirected edge should use --")
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTRendererWithStyles(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes([]output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Styled Node"),
			Shape: output.NodeShapeEllipse,
			Style: output.NodeStyle{
				Fill:   "#ff0000",
				Stroke: "#000000",
			},
		},
	})

	testhelpers.RenderAssert(t, renderer,
		`shape="ellipse"`,
		`fillcolor="#ff0000"`,
	)
}

func TestDOTRendererWithEdgeLabel(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetEdges([]output.GraphEdge{testEdgeAB("uses")})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "label=\"uses\"", "Output should contain edge label")
}

func TestNewDOTFromTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"ID", "Name"})
	data.AddRow([]string{"1", "Alice"})
	data.AddRow([]string{"2", "Bob"})

	renderer := NewDOTFromTable(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph", "Output should be a digraph")
	assertContains(t, out, "row0", "Output should contain row0 node")
	assertContains(
		t,
		out,
		"\"row0\" -> \"row1\"",
		"Output should contain edge from row0 to row1",
	)
}

func TestNewDOTFromTree(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child", "Child"))

	renderer := NewDOTFromTree(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph", "Output should be a digraph")
	assertContains(t, out, "Root", "Output should contain 'Root' label")
	assertContains(t, out, "Child", "Output should contain 'Child' label")
}

func TestDOTRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	testEmptyRendererOutput(t, renderer, testDOTEmptyExpected(t))
}

//nolint:exhaustruct // Test files use partial struct initialization
func TestDOTSetGraphID(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetGraphID("MyGraph")
	renderer.SetNodes(
		[]output.GraphNode{newTestNode("A", "A")},
	)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "digraph MyGraph {", "Output should use custom graph ID")
}

func TestDOTFromTreeWithEmptyID(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("", "Root Node")
	root.AddChild(&output.TreeNode{
		Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Leaf"),
	})

	renderer := NewDOTFromTree(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "Root_Node", "empty ID should use label slug")
	assertContains(t, out, "Leaf", "empty ID should use label slug")
}

func TestNewDOTFromTableNil(t *testing.T) {
	t.Parallel()

	renderer := NewDOTFromTable(nil)
	if renderer == nil {
		t.Fatal("NewDOTFromTable(nil) should return non-nil renderer")
	}
}

func TestDOTFromTreeNil(t *testing.T) {
	t.Parallel()

	renderer := NewDOTFromTree(nil)
	if renderer == nil {
		t.Fatal("NewDOTFromTree(nil) should return non-nil renderer")
	}
}

func TestDOTRendererConfigurableLayout(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer().
		SetRankDir(RankDirLR).
		SetSplines(SplineSpline).
		SetNodeSep("0.8").
		SetRankSep("1.0")

	out := renderWithABNodes(t, renderer)

	assertContains(t, out, "rankdir=LR", "Output should use custom rankdir")
	assertContains(t, out, "splines=spline", "Output should use custom splines")
	assertContains(t, out, "nodesep=0.8", "Output should use custom nodesep")
	assertContains(t, out, "ranksep=1.0", "Output should use custom ranksep")
}

func TestDOTRendererDefaultLayout(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes(testNodesAB())

	testhelpers.RenderAssert(t, renderer, "rankdir=TB", "splines=ortho")
}

func TestDOTRendererSetDirection(t *testing.T) {
	t.Parallel()

	cases := []struct {
		dir  output.Direction
		want string
	}{
		{output.DirectionDown, "rankdir=TB"},
		{output.DirectionRight, "rankdir=LR"},
		{output.DirectionUp, "rankdir=BT"},
		{output.DirectionLeft, "rankdir=RL"},
	}
	for _, tc := range cases {
		renderer := NewDOTRenderer().SetDirection(tc.dir)
		renderer.SetNodes(testNodesAB())

		out, err := renderer.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, out, tc.want, "SetDirection should produce expected rankdir")
	}
}

// TestDOTNodeStyleEscapesInjection verifies that malicious style values
// (double quotes, backslashes, newlines) are escaped through the DOT render
// pipeline. If escape.DOT were removed from writeNodeAttr, a double quote in
// a style value would break out of the quoted attribute, allowing attribute
// injection.
func TestDOTNodeStyleEscapesInjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{"double quote in Fill", `red"; injected=true`},
		{"newline in Fill", "red\ninjected_attr"},
		{"backslash in Fill", `red\malicious`},
		{"double quote in Stroke", `#000"breakout`},
		{"backslash in Stroke", `#000\evil`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			renderer := NewDOTRenderer()
			renderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Test uses minimal fields
				{
					ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
					Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
					Shape: output.NodeShapeBox,
					Style: output.NodeStyle{
						Fill:   tt.value,
						Stroke: tt.value,
					},
				},
			})

			out, err := renderer.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if strings.Contains(out, tt.value) {
				t.Errorf("raw malicious value %q leaked unescaped into DOT output", tt.value)
			}
		})
	}
}

// TestDOTNodeStyleEscapeOutput verifies the exact escaped sequences appear in
// DOT output, complementing the "raw value doesn't leak" check above.
func TestDOTNodeStyleEscapeOutput(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Test uses minimal fields
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("A"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
			Shape: output.NodeShapeBox,
			Style: output.NodeStyle{
				Fill: `a"b\c` + "\n" + `d`,
			},
		},
	})

	testhelpers.RenderAssert(t, renderer,
		`a\"b\\c`,
		`\nd`,
	)
}

func TestDOTGraphIDInjectionNeutralized(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetGraphID("x\nlabel=\"injected\"; //")
	renderer.AddNode(output.GraphNode{ //nolint:exhaustruct // Injection test
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("A"),
	})

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// The header must stay on ONE line: the quoted form may contain the
	// word "injected" as inert ID text, but the raw newline must not split
	// it into a second statement line.
	firstLine := strings.SplitN(got, "\n", 2)[0]
	if !strings.HasSuffix(firstLine, "{") {
		t.Errorf("graph ID broke the DOT header onto multiple lines:\n%s", got)
	}

	if !strings.HasPrefix(firstLine, `digraph "x\nlabel=`) {
		t.Errorf("hostile graph ID should be quoted+escaped, got header %q", firstLine)
	}
}

func TestDOTInvalidRankDirNotEmitted(t *testing.T) {
	t.Parallel()

	renderer := NewDOTRenderer()
	renderer.SetRankDir(RankDir("TB\nlabel=x"))
	renderer.AddNode(output.GraphNode{ //nolint:exhaustruct // Injection test
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("A"),
	})

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "label=x") {
		t.Errorf("invalid rankdir value emitted:\n%s", got)
	}

	if strings.Contains(got, "rankdir=") {
		t.Errorf("invalid rankdir should be omitted entirely:\n%s", got)
	}
}
