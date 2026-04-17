package output

import (
	"strings"
	"testing"
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

	got := d.Render()
	assertContains(t, got, "id: int {constraint: primary_key}", "should contain primary key")
	assertContains(t, got, "email: string {constraint: unique}", "should contain unique")
	assertContains(t, got, "org_id: int {constraint: foreign_key}", "should contain foreign key")
	assertContains(t, got, "name: string\n", "should contain column without constraint")

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

func TestD2FromTableData(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		d := D2FromTableData(nil)
		if d == nil {
			t.Fatal("D2FromTableData(nil) should return non-nil diagram")
		}

		got := d.Render()
		if got != "" {
			t.Errorf("nil data should render empty diagram, got %q", got)
		}
	})

	t.Run("with rows", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"Name", "Value"})
		data.AddRow([]string{"Alpha", "100"})
		data.AddRow([]string{"Beta", "200"})

		d := D2FromTableData(data)
		got := d.Render()
		assertContains(t, got, "row0", "should contain first row node")
		assertContains(t, got, "row1", "should contain second row node")
		assertContains(t, got, "->", "should contain edges")
		assertContains(t, got, "Name: Alpha", "should contain label content")
	})
}

func TestD2FromTree(t *testing.T) {
	t.Parallel()

	t.Run("nil root", func(t *testing.T) {
		t.Parallel()

		d := D2FromTree(nil)
		if d == nil {
			t.Fatal("D2FromTree(nil) should return non-nil diagram")
		}

		got := d.Render()
		if got != "" {
			t.Errorf("nil root should render empty diagram, got %q", got)
		}
	})

	t.Run("tree with children", func(t *testing.T) {
		t.Parallel()

		root := NewTreeNode("root", "Root")
		child1 := NewTreeNode("child1", "Child 1")
		child2 := NewTreeNode("child2", "Child 2")

		root.AddChild(child1)
		root.AddChild(child2)

		d := D2FromTree(root)
		got := d.Render()
		assertContains(t, got, "root: Root", "should contain root node")
		assertContains(t, got, "child1: Child 1", "should contain first child")
		assertContains(t, got, "child2: Child 2", "should contain second child")
		assertContains(t, got, "root -> child1", "should contain edge to first child")
		assertContains(t, got, "root -> child2", "should contain edge to second child")
	})

	t.Run("deeply nested tree", func(t *testing.T) {
		t.Parallel()

		root := NewTreeNode("root", "Root")
		child := NewTreeNode("child", "Child")
		grandchild := NewTreeNode("grandchild", "Grandchild")

		root.AddChild(child)
		child.AddChild(grandchild)

		d := D2FromTree(root)
		got := d.Render()
		assertContains(t, got, "root -> child", "should contain root->child edge")
		assertContains(t, got, "child -> grandchild", "should contain child->grandchild edge")
	})
}

func TestD2GraphRendererInterface(t *testing.T) {
	t.Parallel()

	var _ GraphRenderer = NewD2Diagram()

	d := NewD2Diagram()
	d.SetNodes([]GraphNode{
		*NewGraphNode("a", "Node A"),
		*NewGraphNode("b", "Node B"),
	})
	d.SetEdges([]GraphEdge{
		*NewGraphEdge("a", "b"),
	})

	got := d.Render()
	assertContains(t, got, "a: Node A", "should contain node A")
	assertContains(t, got, "b: Node B", "should contain node B")
	assertContains(t, got, "a -> b", "should contain edge")
}

func TestD2GraphShapeConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		shape GraphShape
		want  string
	}{
		{"box implicit rectangle", ShapeBox, ""},
		{"rect implicit rectangle", ShapeRect, ""},
		{"ellipse to oval", ShapeEllipse, "oval"},
		{"diamond to diamond", ShapeDiamond, "diamond"},
		{"circle to circle", ShapeCircle, "circle"},
		{"cylinder to cylinder", ShapeCylinder, "cylinder"},
		{"hexagon to hexagon", ShapeHexagon, "hexagon"},
		{"parallelogram to parallelogram", ShapeParallelogram, "parallelogram"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.SetNodes([]GraphNode{newTestNodeWithShape("node", "Test", tt.shape)})

			got := d.Render()

			if tt.want == "" {
				if strings.Contains(got, "shape:") {
					t.Errorf(
						"shape %q should be implicit, got shape in output",
						tt.shape,
					)
				}
			} else {
				msg := "should convert " + tt.shape.String() + " to " + tt.want
				assertContains(t, got, "shape: "+tt.want, msg)
			}
		})
	}
}

func TestD2GraphStyleConversion(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.SetNodes([]GraphNode{
		{
			ID:    NewBrandedID[GraphNodeIDBrand]("styled"),
			Label: NewBrandedID[GraphNodeLabelBrand]("Styled"),
			Style: GraphStyle{
				FillColor:   "blue",
				StrokeColor: "black",
				FontSize:    14,
			},
		},
	})

	got := d.Render()
	assertContains(t, got, "style.fill: blue", "should convert fill color")
	assertContains(t, got, "style.stroke: black", "should convert stroke color")
	assertContains(t, got, "style.font-size: 14", "should convert font size")
}
