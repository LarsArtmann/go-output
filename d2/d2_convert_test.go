package d2

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestD2TableWithConstraints(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddTable("users", []D2Column{
		{Name: "id", Type: "int", Constraint: D2ConstraintPrimary},
		{Name: "email", Type: "string", Constraint: D2ConstraintUnique},
		{Name: "org_id", Type: "int", Constraint: D2ConstraintForeign},
		{Name: "name", Type: "string"},
	})

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(
		t,
		got,
		"id: int {constraint: primary_key}",
		"should contain primary key",
	)
	testhelpers.AssertContains(
		t,
		got,
		"email: string {constraint: unique}",
		"should contain unique",
	)
	testhelpers.AssertContains(
		t,
		got,
		"org_id: int {constraint: foreign_key}",
		"should contain foreign key",
	)
	testhelpers.AssertContains(t, got, "name: string\n", "should contain column without constraint")

	if strings.Contains(got, "name: string {constraint") {
		t.Error("column without constraint should not have constraint block")
	}
}

func TestD2ConstraintConstants(t *testing.T) {
	t.Parallel()

	if D2ConstraintPrimary != "primary_key" {
		t.Errorf("D2ConstraintPrimary = %q, want %q", D2ConstraintPrimary, "primary_key")
	}

	if D2ConstraintForeign != "foreign_key" {
		t.Errorf("D2ConstraintForeign = %q, want %q", D2ConstraintForeign, "foreign_key")
	}

	if D2ConstraintUnique != "unique" {
		t.Errorf("D2ConstraintUnique = %q, want %q", D2ConstraintUnique, "unique")
	}
}

func TestNewD2FromTable(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		d := NewD2FromTable(nil)
		if d == nil {
			t.Fatal("NewD2FromTable(nil) should return non-nil diagram")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "" {
			t.Errorf("nil data should render empty diagram, got %q", got)
		}
	})

	t.Run("with rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Value"})
		data.AddRow([]string{"Alpha", "100"})
		data.AddRow([]string{"Beta", "200"})

		d := NewD2FromTable(data)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "row0", "should contain first row node")
		testhelpers.AssertContains(t, got, "row1", "should contain second row node")
		testhelpers.AssertContains(t, got, "->", "should contain edges")
		testhelpers.AssertContains(t, got, "Name: Alpha", "should contain label content")
	})
}

func TestNewD2FromTree(t *testing.T) {
	t.Parallel()

	t.Run("nil root", func(t *testing.T) {
		t.Parallel()

		d := NewD2FromTree(nil)
		if d == nil {
			t.Fatal("NewD2FromTree(nil) should return non-nil diagram")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "" {
			t.Errorf("nil root should render empty diagram, got %q", got)
		}
	})

	t.Run("tree with children", func(t *testing.T) {
		t.Parallel()

		root := output.NewTreeNode("root", "Root")

		root.AddChild(output.NewTreeNode("child1", "Child 1"))
		root.AddChild(output.NewTreeNode("child2", "Child 2"))

		d := NewD2FromTree(root)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "root: Root", "should contain root node")
		testhelpers.AssertContains(t, got, "child1: Child 1", "should contain first child")
		testhelpers.AssertContains(t, got, "child2: Child 2", "should contain second child")
		testhelpers.AssertContains(t, got, "root -> child1", "should contain edge to first child")
		testhelpers.AssertContains(t, got, "root -> child2", "should contain edge to second child")
	})

	t.Run("deeply nested tree", func(t *testing.T) {
		t.Parallel()

		root := output.NewTreeNode("root", "Root")

		root.AddChild(output.NewTreeNode("child", "Child"))
		root.Children[0].AddChild(output.NewTreeNode("grandchild", "Grandchild"))

		d := NewD2FromTree(root)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "root -> child", "should contain root->child edge")
		testhelpers.AssertContains(
			t,
			got,
			"child -> grandchild",
			"should contain child->grandchild edge",
		)
	})

	t.Run("empty ID uses label slug", func(t *testing.T) {
		t.Parallel()

		root := &output.TreeNode{
			Label: output.NewBrandedID[output.TreeNodeLabelBrand]("My Root"),
			Children: []*output.TreeNode{
				{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Child Node")},
			},
		}

		d := NewD2FromTree(root)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "My_Root: My Root", "should use slug for empty ID")
		testhelpers.AssertContains(t, got, "My_Root -> Child_Node", "should use slug in edge")
	})
}

func TestD2GraphRendererInterface(t *testing.T) {
	t.Parallel()

	var _ output.GraphRenderer = NewD2Diagram()

	d := NewD2Diagram()
	d.SetNodes([]output.GraphNode{
		*output.NewGraphNode("a", "Node A"),
		*output.NewGraphNode("b", "Node B"),
	})
	d.SetEdges([]output.GraphEdge{
		*output.NewGraphEdge("a", "b"),
	})

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "a: Node A", "should contain node A")
	testhelpers.AssertContains(t, got, "b: Node B", "should contain node B")
	testhelpers.AssertContains(t, got, "a -> b", "should contain edge")
}

func TestD2NodeShapeConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		shape output.NodeShape
		want  string
	}{
		{"box implicit rectangle", output.NodeShapeBox, ""},
		{"ellipse to oval", output.NodeShapeEllipse, "oval"},
		{"diamond to diamond", output.NodeShapeDiamond, "diamond"},
		{"circle to circle", output.NodeShapeCircle, "circle"},
		{"cylinder to cylinder", output.NodeShapeCylinder, "cylinder"},
		{"hexagon to hexagon", output.NodeShapeHexagon, "hexagon"},
		{"parallelogram to parallelogram", output.NodeShapeParallelogram, "parallelogram"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.SetNodes([]output.GraphNode{
				{
					ID:    output.NewBrandedID[output.GraphNodeIDBrand]("node"),
					Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
					Shape: tt.shape,
				},
			})

			got, err := d.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if tt.want == "" {
				if strings.Contains(got, "shape:") {
					t.Errorf("shape %q implicit but got shape in output", tt.shape)
				}
			} else {
				msg := "should convert " + tt.shape.String() + " to " + tt.want
				testhelpers.AssertContains(t, got, "shape: "+tt.want, msg)
			}
		})
	}
}

func TestD2NodeStyleConversion(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.SetNodes([]output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("styled"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Styled"),
			Style: output.NodeStyle{
				Fill:     "blue",
				Stroke:   "black",
				FontSize: 14,
			},
		},
	})

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "style.fill: blue", "should convert fill color")
	testhelpers.AssertContains(t, got, "style.stroke: black", "should convert stroke color")
	testhelpers.AssertContains(t, got, "style.font-size: 14", "should convert font size")
}
