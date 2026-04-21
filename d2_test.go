package output

import (
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

func TestParseD2NodeShape(t *testing.T) {
	t.Parallel()

	t.Run("valid shapes", func(t *testing.T) {
		t.Parallel()

		shapes := []string{
			"rectangle", "circle", "diamond", "hexagon", "cloud",
			"cylinder", "person", "queue", "oval", "parallelogram",
			"triangle", "sql_table", "image", "code", "text", "class",
			"page", "step", "stored_data", "square",
		}

		for _, s := range shapes {
			got, err := ParseD2NodeShape(s)
			if err != nil {
				t.Errorf("ParseD2NodeShape(%q) unexpected error: %v", s, err)
			}

			if string(got) != s {
				t.Errorf("ParseD2NodeShape(%q) = %q, want %q", s, got, s)
			}
		}
	})

	t.Run("invalid shape", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2NodeShape("invalid")
		if err == nil {
			t.Error("ParseD2NodeShape(invalid) should return error")
		}
	})
}

func TestD2NodeShapeValidation(t *testing.T) {
	t.Parallel()

	if !D2ShapeCircle.IsValid() {
		t.Error("D2ShapeCircle should be valid")
	}

	invalid := D2NodeShape("nonexistent")
	if invalid.IsValid() {
		t.Error("nonexistent shape should not be valid")
	}

	values := D2ShapeRectangle.AllowedValues()
	if len(values) != 20 {
		t.Errorf("AllowedValues() returned %d values, want 20", len(values))
	}
}

func TestParseD2ArrowType(t *testing.T) {
	t.Parallel()

	t.Run("valid arrows", func(t *testing.T) {
		t.Parallel()

		arrows := []string{
			"arrow", "triangle", "diamond", "circle", "filled",
			"box", "cross", "cf-one", "cf-many",
			"cf-one-required", "cf-many-required",
		}

		for _, s := range arrows {
			got, err := ParseD2ArrowType(s)
			if err != nil {
				t.Errorf("ParseD2ArrowType(%q) unexpected error: %v", s, err)
			}

			if string(got) != s {
				t.Errorf("ParseD2ArrowType(%q) = %q, want %q", s, got, s)
			}
		}
	})

	t.Run("invalid arrow", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2ArrowType("invalid")
		if err == nil {
			t.Error("ParseD2ArrowType(invalid) should return error")
		}
	})

	t.Run("empty arrow is valid", func(t *testing.T) {
		t.Parallel()

		a := D2ArrowType("")
		if !a.IsValid() {
			t.Error("empty arrow type should be valid (means no arrow)")
		}
	})
}

func TestD2ArrowTypeValidation(t *testing.T) {
	t.Parallel()

	if !D2ArrowArrow.IsValid() {
		t.Error("D2ArrowArrow should be valid")
	}

	values := D2ArrowArrow.AllowedValues()
	if len(values) != 11 {
		t.Errorf("AllowedValues() returned %d values, want 11", len(values))
	}
}

func TestParseD2Direction(t *testing.T) {
	t.Parallel()

	t.Run("valid directions", func(t *testing.T) {
		t.Parallel()

		cases := []struct{ input, want string }{
			{"", ""},
			{"right", "right"},
			{"left", "left"},
			{"up", "up"},
		}

		for _, tc := range cases {
			got, err := ParseD2Direction(tc.input)
			if err != nil {
				t.Errorf("ParseD2Direction(%q) unexpected error: %v", tc.input, err)
			}

			if string(got) != tc.want {
				t.Errorf("ParseD2Direction(%q) = %q, want %q", tc.input, got, tc.want)
			}
		}
	})

	t.Run("invalid direction", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2Direction("diagonal")
		if err == nil {
			t.Error("ParseD2Direction(diagonal) should return error")
		}
	})
}

func TestD2DirectionValidation(t *testing.T) {
	t.Parallel()

	if !D2DirRight.IsValid() {
		t.Error("D2DirRight should be valid")
	}

	invalid := D2Direction("down-under")
	if invalid.IsValid() {
		t.Error("invalid direction should not be valid")
	}

	values := D2DirRight.AllowedValues()
	if len(values) != 4 {
		t.Errorf("AllowedValues() returned %d values, want 4", len(values))
	}
}
