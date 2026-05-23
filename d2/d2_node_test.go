package d2

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/escape"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestD2AllNodeShapes(t *testing.T) {
	t.Parallel()

	shapes := d2NodeShapeValues

	for _, shape := range shapes {
		t.Run(string(shape), func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.AddNode(D2Node{ //nolint:exhaustruct // Test uses minimal required fields
				ID:    NewBrandedID[D2NodeIDBrand]("node"),
				Label: NewBrandedID[D2NodeLabelBrand]("Test"),
				Shape: shape,
			})

			got, err := d.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "shape: circle", "nested node should support shape")
	assertContains(t, got, "child: Inner", "should contain nested content")
}

func TestD2NodeWithSpecialChars(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddNodeSimple("node", `has "quotes" and\nnewlines`)

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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
		{"mixed", "\"hello\"\nworld", `\"hello\"\nworld`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := escape.D2(tt.input)
			testhelpers.AssertEqual(t, "escape.D2", tt.input, got, tt.want)
		})
	}
}
