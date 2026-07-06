package d2

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestDiagram(t *testing.T) {
	t.Parallel()

	t.Run("basic diagram with table", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()
		d.AddTable("users", []Column{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "string"},
		})

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "users:", "should contain table name")
		testhelpers.AssertContains(t, got, "shape: sql_table", "should contain sql_table shape")
		testhelpers.AssertContains(t, got, "id: int", "should contain column definitions")
	})

	t.Run("chaining", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram().
			AddTable("users", []Column{{Name: "id", Type: "int"}})

		if d == nil {
			t.Error("Method chaining should return non-nil")
		}
	})

	t.Run("multiple tables", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()
		d.AddTable("users", []Column{{Name: "id", Type: "int"}})
		d.AddTable("posts", []Column{{Name: "id", Type: "int"}})

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "users:", "should contain users table")
		testhelpers.AssertContains(t, got, "posts:", "should contain posts table")
	})

	t.Run("empty diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "" {
			t.Errorf("empty diagram should render empty string, got %q", got)
		}
	})
}

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	if d == nil {
		t.Fatal("NewDiagram() should return non-nil")
	}
}

func TestD2Config(t *testing.T) {
	t.Parallel()

	t.Run("direction right", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram().SetDirection(DirRight)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "direction: right", "should contain direction")
	})

	t.Run("title", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram().SetTitle("My Diagram")

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "title:", "should contain title block")
		testhelpers.AssertContains(t, got, "My Diagram", "should contain title value")
	})

	t.Run("layout engine", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram().SetLayout("elk")

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "layout: elk", "should contain layout engine")
	})

	t.Run("combined config", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram().SetDirection(DirRight).SetLayout("elk").SetTitle("Architecture")

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "direction: right", "should contain direction")
		testhelpers.AssertContains(t, got, "layout: elk", "should contain layout")
		testhelpers.AssertContains(t, got, "Architecture", "should contain title")
	})
}

func TestD2Diagram_AddNode(t *testing.T) {
	t.Parallel()

	t.Run("AddNode", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddNode(Node{ //nolint:exhaustruct // Test uses minimal required fields
			ID:    output.NewBrandedID[output.D2NodeIDBrand]("server"),
			Label: output.NewBrandedID[output.D2NodeLabelBrand]("Web Server"),
		})
		if result != d {
			t.Error("AddNode should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "server: Web Server", "should contain node")
	})

	t.Run("AddNodeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddNodeSimple("db", "Database")
		if result != d {
			t.Error("AddNodeSimple should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "db: Database", "should contain simple node")
	})

	t.Run("AddNodeWithShape", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddNodeWithShape("cache", "Cache", ShapeCircle)
		if result != d {
			t.Error("AddNodeWithShape should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(
			t,
			got,
			"cache: Cache {",
			"should use block syntax for shaped node",
		)
		testhelpers.AssertContains(t, got, "shape: circle", "should contain shape attribute")
	})
}

func TestD2Diagram_AddEdge(t *testing.T) {
	t.Parallel()

	t.Run("AddEdge", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddEdge(Edge{ //nolint:exhaustruct // Test uses minimal required fields
			From: output.NewBrandedID[output.D2NodeIDBrand]("a"),
			To:   output.NewBrandedID[output.D2NodeIDBrand]("b"),
		})
		if result != d {
			t.Error("AddEdge should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "a -> b", "should contain edge")
	})

	t.Run("AddEdgeSimple", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddEdgeSimple("x", "y")
		if result != d {
			t.Error("AddEdgeSimple should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "x -> y", "should contain simple edge")
	})

	t.Run("AddLabeledEdge", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram()

		result := d.AddLabeledEdge("src", "dst", "connects")
		if result != d {
			t.Error("AddLabeledEdge should return diagram for chaining")
		}

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertContains(t, got, "src -> dst: connects", "should contain labeled edge")
	})
}
