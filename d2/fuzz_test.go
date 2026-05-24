package d2

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

func FuzzD2Escape(f *testing.F) {
	f.Add("hello")
	f.Add(`"quoted"`)
	f.Add(`back\slash`)
	f.Add("new\nline")
	f.Add("tab\there")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		escaped := escape.D2(input)

		if strings.Contains(input, `"`) && !strings.Contains(escaped, `\"`) {
			t.Errorf("D2(%q) = %q, quotes not escaped", input, escaped)
		}

		if strings.Contains(input, `\`) && !strings.Contains(escaped, `\\`) {
			t.Errorf("D2(%q) = %q, backslashes not escaped", input, escaped)
		}
	})
}

func FuzzParseD2Direction(f *testing.F) {
	for _, dir := range d2DirectionValues {
		f.Add(string(dir))
	}

	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got, err := ParseD2Direction(input)
		if err != nil {
			if got != "" {
				t.Errorf("ParseD2Direction(%q) returned non-empty on error: %q", input, got)
			}

			return
		}

		if !got.IsValid() {
			t.Errorf("ParseD2Direction(%q) = %q, not valid", input, got)
		}
	})
}

func FuzzParseD2NodeShape(f *testing.F) {
	for _, shape := range d2NodeShapeValues {
		f.Add(string(shape))
	}

	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got, err := ParseD2NodeShape(input)
		if err != nil {
			if got != "" {
				t.Errorf("ParseD2NodeShape(%q) returned non-empty on error: %q", input, got)
			}

			return
		}

		if !got.IsValid() {
			t.Errorf("ParseD2NodeShape(%q) = %q, not valid", input, got)
		}
	})
}

func FuzzParseD2ArrowType(f *testing.F) {
	for _, arrow := range d2ArrowTypeValues {
		f.Add(string(arrow))
	}

	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got, err := ParseD2ArrowType(input)
		if err != nil {
			if got != "" {
				t.Errorf("ParseD2ArrowType(%q) returned non-empty on error: %q", input, got)
			}

			return
		}

		if !got.IsValid() {
			t.Errorf("ParseD2ArrowType(%q) = %q, not valid", input, got)
		}
	})
}

func FuzzParseD2Constraint(f *testing.F) {
	for _, c := range allD2Constraints {
		f.Add(string(c))
	}

	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got, err := ParseD2Constraint(input)
		if err != nil {
			if got != "" {
				t.Errorf("ParseD2Constraint(%q) returned non-empty on error: %q", input, got)
			}

			return
		}

		if !got.IsValid() {
			t.Errorf("ParseD2Constraint(%q) = %q, not valid", input, got)
		}
	})
}

func FuzzD2DiagramRender(f *testing.F) {
	f.Add("nodeA", "labelA", "nodeB", "labelB")
	f.Add("", "", "", "")
	f.Add(`"quoted"`, `\slash`, "normal", "ok")

	f.Fuzz(func(t *testing.T, idA, labelA, idB, labelB string) {
		d := NewD2Diagram()
		d.AddNodeSimple(idA, labelA)
		d.AddNodeSimple(idB, labelB)
		d.AddEdgeSimple(idA, idB)

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		if idA != "" && !strings.Contains(got, escape.D2(idA)) {
			t.Errorf("Render() output missing node id %q: %s", idA, got)
		}

		if idB != "" && !strings.Contains(got, escape.D2(idB)) {
			t.Errorf("Render() output missing node id %q: %s", idB, got)
		}
	})
}

func FuzzGraphNodeRoundTrip(f *testing.F) {
	f.Add("test-id", "test-label", "box")
	f.Add("node1", "My Node", "circle")
	f.Add("", "", "")

	f.Fuzz(func(t *testing.T, id, label, shapeStr string) {
		shape, err := output.ParseGraphShape(shapeStr)
		if err != nil {
			return
		}

		node := output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](id),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](label),
			Shape: shape,
		}

		d := NewD2Diagram()
		d.SetNodes([]output.GraphNode{node})

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		if id != "" && !strings.Contains(got, escape.D2(id)) {
			t.Errorf("round-trip: output missing id %q", id)
		}
	})
}
