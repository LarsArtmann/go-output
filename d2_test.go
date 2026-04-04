package output

import (
	"strings"
	"testing"
)

func TestD2Diagram(t *testing.T) {
	t.Parallel()
	t.Run("basic diagram", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()
		d.AddTable("users", []D2Column{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "string"},
		})

		got := d.Render()

		if !strings.Contains(got, "users:") {
			t.Error("Render() should contain table name")
		}

		if !strings.Contains(got, "id: int") {
			t.Error("Render() should contain column definitions")
		}
	})

	t.Run("chaining", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().
			AddTable("users", []D2Column{
				{Name: "id", Type: "int"},
			})

		if d == nil {
			t.Error("Method chaining should return non-nil")
		}
	})

	t.Run("multiple tables", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()
		d.AddTable("users", []D2Column{{Name: "id", Type: "int"}})
		d.AddTable("posts", []D2Column{{Name: "id", Type: "int"}})

		got := d.Render()

		if !strings.Contains(got, "users:") {
			t.Error("Render() should contain users table")
		}

		if !strings.Contains(got, "posts:") {
			t.Error("Render() should contain posts table")
		}
	})
}

func TestNewD2Diagram(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	// Verify diagram is initialized properly
	_ = d.tables // Just ensure field is accessible
}

func TestD2Diagram_AddNode(t *testing.T) {
	t.Parallel()
	t.Run("AddNode", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNode(
			//nolint:exhaustruct // Test uses minimal required fields
			D2Node{
				ID:    NewBrandedID[D2NodeIDBrand]("server"),
				Label: NewBrandedID[D2NodeLabelBrand]("Web Server"),
			},
		)
		if result != d {
			t.Error("AddNode should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "server") {
			t.Error("Render() should contain node ID")
		}
	})

	t.Run("AddNodeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNodeSimple("db", "Database")
		if result != d {
			t.Error("AddNodeSimple should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "db") || !strings.Contains(got, "Database") {
			t.Error("Render() should contain simple node")
		}
	})

	t.Run("AddNodeWithShape", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNodeWithShape("cache", "Cache", D2ShapeCircle)
		if result != d {
			t.Error("AddNodeWithShape should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "cache") {
			t.Error("Render() should contain node")
		}

		if !strings.Contains(got, ":circle") {
			t.Error("Render() should contain shape attribute")
		}
	})
}

func TestD2Diagram_AddEdge(t *testing.T) {
	t.Parallel()
	t.Run("AddEdge", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()
		d.AddNodeSimple("a", "Node A")
		d.AddNodeSimple("b", "Node B")

		result := d.AddEdge(
			//nolint:exhaustruct // Test uses minimal required fields
			D2Edge{From: NewBrandedID[D2NodeIDBrand]("a"), To: NewBrandedID[D2NodeIDBrand]("b")},
		)
		if result != d {
			t.Error("AddEdge should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "->") {
			t.Error("Render() should contain edge arrow")
		}
	})

	t.Run("AddEdgeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddEdgeSimple("x", "y")
		if result != d {
			t.Error("AddEdgeSimple should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "x -> y") {
			t.Errorf("Render() should contain 'x -> y', got: %q", got)
		}
	})

	t.Run("AddLabeledEdge", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddLabeledEdge("src", "dst", "connects")
		if result != d {
			t.Error("AddLabeledEdge should return diagram for chaining")
		}

		got := d.Render()
		if !strings.Contains(got, "connects") {
			t.Error("Render() should contain edge label")
		}
	})
}

func TestD2FromTableData(t *testing.T) {
	t.Parallel()
	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		d := D2FromTableData(nil)
		if d == nil {
			t.Error("D2FromTableData(nil) should return non-nil diagram")
		}
	})

	t.Run("with data", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"Name", "Value"})
		data.AddRow([]string{"test", "123"})

		d := D2FromTableData(data)
		if d == nil {
			t.Fatal("D2FromTableData should return non-nil diagram")
		}

		got := d.Render()
		if !strings.Contains(got, "Name") {
			t.Error("Render() should contain header Name")
		}
	})
}

func TestD2NodeShapes(t *testing.T) {
	t.Parallel()

	shapes := []D2NodeShape{
		D2ShapeRectangle,
		D2ShapeCircle,
		D2ShapeDiamond,
		D2ShapeCloud,
		D2ShapeCylinder,
		D2ShapeOval,
	}

	for _, shape := range shapes {
		t.Run(string(shape), func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.AddNode(
				//nolint:exhaustruct // Test uses minimal required fields
				D2Node{
					ID:    NewBrandedID[D2NodeIDBrand]("node"),
					Label: NewBrandedID[D2NodeLabelBrand]("Test"),
					Shape: shape,
				},
			)

			got := d.Render()
			if !strings.Contains(got, "node") {
				t.Error("Render() should contain node ID")
			}
		})
	}
}

func TestD2NodeWithStyle(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(
		//nolint:exhaustruct // Test uses minimal required fields
		D2Node{
			ID:    NewBrandedID[D2NodeIDBrand]("styled"),
			Label: NewBrandedID[D2NodeLabelBrand]("Styled Node"),
			Shape: D2ShapeRectangle,
			Style: D2NodeStyle{
				Fill:        "blue",
				Stroke:      "black",
				StrokeWidth: 2,
				FontSize:    14,
			},
		})

	got := d.Render()
	if !strings.Contains(got, "fill:blue") {
		t.Error("Render() should contain fill style")
	}

	if !strings.Contains(got, "stroke:black") {
		t.Error("Render() should contain stroke style")
	}
}

func TestD2EdgeWithArrows(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddEdge(
		//nolint:exhaustruct // Test uses minimal required fields
		D2Edge{
			From:        NewBrandedID[D2NodeIDBrand]("a"),
			To:          NewBrandedID[D2NodeIDBrand]("b"),
			Label:       NewBrandedID[D2NodeLabelBrand]("test"),
			SourceArrow: D2ArrowDiamond,
			TargetArrow: D2ArrowTriangle,
		},
	)

	got := d.Render()
	if !strings.Contains(got, "-diamond") {
		t.Error("Render() should contain source arrow")
	}

	if !strings.Contains(got, "-triangle") {
		t.Error("Render() should contain target arrow")
	}
}

func TestD2NodeNested(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(
		//nolint:exhaustruct // Test uses minimal required fields
		D2Node{
			ID:     NewBrandedID[D2NodeIDBrand]("parent"),
			Label:  NewBrandedID[D2NodeLabelBrand]("Parent"),
			Shape:  D2ShapeRectangle,
			Nested: "child: inner\n",
		},
	)

	got := d.Render()
	if !strings.Contains(got, "child: inner") {
		t.Error("Render() should contain nested content")
	}
}
