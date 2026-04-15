package output

import (
	"strings"
	"testing"
)

func TestD2Diagram(t *testing.T) {
	t.Parallel()

	t.Run("basic diagram with table", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()
		d.AddTable("users", []D2Column{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "string"},
		})

		got := d.Render()
		assertContains(t, got, "users:", "should contain table name")
		assertContains(t, got, "shape: sql_table", "should contain sql_table shape")
		assertContains(t, got, "id: int", "should contain column definitions")
	})

	t.Run("chaining", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().
			AddTable("users", []D2Column{{Name: "id", Type: "int"}})

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
		assertContains(t, got, "users:", "should contain users table")
		assertContains(t, got, "posts:", "should contain posts table")
	})

	t.Run("empty diagram", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		got := d.Render()
		if got != "" {
			t.Errorf("empty diagram should render empty string, got %q", got)
		}
	})
}

func TestNewD2Diagram(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	if d == nil {
		t.Fatal("NewD2Diagram() should return non-nil")
	}
}

func TestD2Config(t *testing.T) {
	t.Parallel()

	t.Run("direction right", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().SetDirection(D2DirRight)
		got := d.Render()
		assertContains(t, got, "direction: right", "should contain direction")
	})

	t.Run("title", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().SetTitle("My Diagram")
		got := d.Render()
		assertContains(t, got, "title:", "should contain title block")
		assertContains(t, got, "My Diagram", "should contain title value")
	})

	t.Run("layout engine", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().SetLayout("elk")
		got := d.Render()
		assertContains(t, got, "layout: elk", "should contain layout engine")
	})

	t.Run("combined config", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram().SetDirection(D2DirRight).SetLayout("elk").SetTitle("Architecture")
		got := d.Render()
		assertContains(t, got, "direction: right", "should contain direction")
		assertContains(t, got, "layout: elk", "should contain layout")
		assertContains(t, got, "Architecture", "should contain title")
	})
}

func TestD2Diagram_AddNode(t *testing.T) {
	t.Parallel()

	t.Run("AddNode", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
			ID:    NewBrandedID[D2NodeIDBrand]("server"),
			Label: NewBrandedID[D2NodeLabelBrand]("Web Server"),
		})
		if result != d {
			t.Error("AddNode should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "server: Web Server", "should contain node")
	})

	t.Run("AddNodeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNodeSimple("db", "Database")
		if result != d {
			t.Error("AddNodeSimple should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "db: Database", "should contain simple node")
	})

	t.Run("AddNodeWithShape", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddNodeWithShape("cache", "Cache", D2ShapeCircle)
		if result != d {
			t.Error("AddNodeWithShape should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "cache: Cache {", "should use block syntax for shaped node")
		assertContains(t, got, "shape: circle", "should contain shape attribute")
	})
}

func TestD2Diagram_AddEdge(t *testing.T) {
	t.Parallel()

	t.Run("AddEdge", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddEdge(D2Edge{ //nolint:exhaustruct // Test uses minimal required fields
			From: NewBrandedID[D2NodeIDBrand]("a"),
			To:   NewBrandedID[D2NodeIDBrand]("b"),
		})
		if result != d {
			t.Error("AddEdge should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "a -> b", "should contain edge")
	})

	t.Run("AddEdgeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddEdgeSimple("x", "y")
		if result != d {
			t.Error("AddEdgeSimple should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "x -> y", "should contain simple edge")
	})

	t.Run("AddLabeledEdge", func(t *testing.T) {
		t.Parallel()

		d := NewD2Diagram()

		result := d.AddLabeledEdge("src", "dst", "connects")
		if result != d {
			t.Error("AddLabeledEdge should return diagram for chaining")
		}

		got := d.Render()
		assertContains(t, got, "src -> dst: connects", "should contain labeled edge")
	})
}

func TestD2AllNodeShapes(t *testing.T) {
	t.Parallel()

	shapes := []D2NodeShape{
		D2ShapeRectangle,
		D2ShapeSquare,
		D2ShapeCircle,
		D2ShapeDiamond,
		D2ShapeHexagon,
		D2ShapeCloud,
		D2ShapeCylinder,
		D2ShapePerson,
		D2ShapeQueue,
		D2ShapeOval,
		D2ShapeParallelogram,
		D2ShapeTriangle,
		D2ShapeSQLTable,
		D2ShapeImage,
		D2ShapeCode,
		D2ShapeText,
		D2ShapeClass,
		D2ShapePage,
		D2ShapeStep,
		D2ShapeStoredData,
	}

	for _, shape := range shapes {
		t.Run(string(shape), func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
				ID:    NewBrandedID[D2NodeIDBrand]("node"),
				Label: NewBrandedID[D2NodeLabelBrand]("Test"),
				Shape: shape,
			})

			got := d.Render()
			assertContains(t, got, "node:", "should contain node ID")

			if shape != D2ShapeRectangle {
				assertContains(t, got, "shape: "+string(shape), "should contain shape")
			}
		})
	}
}

func TestD2NodeRectangleImplicit(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:    NewBrandedID[D2NodeIDBrand]("node"),
		Label: NewBrandedID[D2NodeLabelBrand]("Simple"),
		Shape: D2ShapeRectangle,
	})

	got := d.Render()
	if strings.Contains(got, "shape:") {
		t.Error("rectangle shape should be implicit, not explicitly rendered")
	}

	assertContains(t, got, "node: Simple", "should render as simple node")
}

func TestD2NodeWithStyle(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{
		ID:     NewBrandedID[D2NodeIDBrand]("styled"),
		Label:  NewBrandedID[D2NodeLabelBrand]("Styled Node"),
		Width:  200,
		Height: 100,
		Style: D2NodeStyle{
			Fill:          "blue",
			Stroke:        "black",
			StrokeWidth:   2,
			StrokeDash:    3,
			FontSize:      14,
			FontColor:     "white",
			Opacity:       0.8,
			Shadow:        true,
			BorderRadius:  8,
			TextTransform: "uppercase",
		},
	})

	got := d.Render()
	assertContains(t, got, "style.fill: blue", "should contain fill style")
	assertContains(t, got, "style.stroke: black", "should contain stroke style")
	assertContains(t, got, "style.stroke-width: 2", "should contain stroke-width")
	assertContains(t, got, "style.stroke-dash: 3", "should contain stroke-dash")
	assertContains(t, got, "style.font-size: 14", "should contain font-size")
	assertContains(t, got, "style.font-color: white", "should contain font-color")
	assertContains(t, got, "style.opacity: 0.8", "should contain opacity")
	assertContains(t, got, "shadow: true", "should contain shadow")
	assertContains(t, got, "border-radius: 8", "should contain border-radius")
	assertContains(t, got, "width: 200", "should contain width")
	assertContains(t, got, "height: 100", "should contain height")
	assertContains(t, got, "style.text-transform: uppercase", "should contain text-transform")
}

func TestD2NodeWithIcon(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:    NewBrandedID[D2NodeIDBrand]("api"),
		Label: NewBrandedID[D2NodeLabelBrand]("API Server"),
		Icon:  "https://icons.terrastruct.com/essentials/004-cloud.svg",
	})

	got := d.Render()
	assertContains(t, got, "icon:", "should contain icon attribute")
	assertContains(t, got, "004-cloud.svg", "should contain icon URL")
}

func TestD2NodeWithLink(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:    NewBrandedID[D2NodeIDBrand]("docs"),
		Label: NewBrandedID[D2NodeLabelBrand]("Documentation"),
		Link:  "https://example.com/docs",
	})

	got := d.Render()
	assertContains(t, got, "link: https://example.com/docs", "should contain link")
}

func TestD2NodeWithTooltip(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:      NewBrandedID[D2NodeIDBrand]("info"),
		Label:   NewBrandedID[D2NodeLabelBrand]("Info"),
		Tooltip: "Additional information",
	})

	got := d.Render()
	assertContains(t, got, "tooltip: Additional information", "should contain tooltip")
}

func TestD2NodeWithNear(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNodeSimple("server", "Server")
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:    NewBrandedID[D2NodeIDBrand]("label"),
		Label: NewBrandedID[D2NodeLabelBrand]("Label"),
		Near:  "server",
	})

	got := d.Render()
	assertContains(t, got, "near: server", "should contain near attribute")
}

func TestD2NodeWithGrid(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:          NewBrandedID[D2NodeIDBrand]("grid"),
		Label:       NewBrandedID[D2NodeLabelBrand]("Grid Container"),
		GridRows:    3,
		GridColumns: 2,
		GridGap:     10,
	})

	got := d.Render()
	assertContains(t, got, "grid-rows: 3", "should contain grid-rows")
	assertContains(t, got, "grid-columns: 2", "should contain grid-columns")
	assertContains(t, got, "grid-gap: 10", "should contain grid-gap")
}

func TestD2NodeWithClass(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddClass("important", D2NodeStyle{
		Fill:   "red",
		Stroke: "darkred",
	})
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:    NewBrandedID[D2NodeIDBrand]("alert"),
		Label: NewBrandedID[D2NodeLabelBrand]("Alert"),
		Class: "important",
	})

	got := d.Render()
	assertContains(t, got, "classes:", "should contain classes block")
	assertContains(t, got, "important:", "should contain class definition")
	assertContains(t, got, "style.fill: red", "should contain class style")
	assertContains(t, got, "class: important", "should contain class reference")
}

func TestD2NodeNested(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:     NewBrandedID[D2NodeIDBrand]("parent"),
		Label:  NewBrandedID[D2NodeLabelBrand]("Parent"),
		Nested: "  child: Inner\n",
	})

	got := d.Render()
	assertContains(t, got, "child: Inner", "should contain nested content")
	assertContains(t, got, "parent: Parent {", "should contain parent block")
}

func TestD2NodeNestedWithShape(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
		ID:     NewBrandedID[D2NodeIDBrand]("parent"),
		Label:  NewBrandedID[D2NodeLabelBrand]("Parent"),
		Shape:  D2ShapeCircle,
		Nested: "  child: Inner\n",
	})

	got := d.Render()
	assertContains(t, got, "shape: circle", "nested node should support shape")
	assertContains(t, got, "child: Inner", "should contain nested content")
}

func TestD2NodeWithSpecialChars(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNodeSimple("node", `has "quotes" and\nnewlines`)

	got := d.Render()
	if strings.Contains(got, `"quotes"`) {
		t.Error("quotes should be escaped in D2 output")
	}

	assertContains(t, got, `\"quotes\"`, "quotes should be escaped")
}

func TestEscapeD2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"plain", "hello", "hello"},
		{"quotes", `"quoted"`, `\"quoted\"`},
		{"newline", "line1\nline2", `line1\nline2`},
		{"tab", "col1\tcol2", `col1\tcol2`},
		{"mixed", `"hello"\nworld`, `\"hello\"\nworld`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := escapeD2(tt.input)
			if got != tt.want {
				t.Errorf("escapeD2(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
