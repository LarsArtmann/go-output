package d2

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
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

		graphtest.AssertEscape(t, "D2", input, escaped, `"`, "quotes")
		graphtest.AssertEscape(t, "D2", input, escaped, `\`, "backslashes")
	})
}

func FuzzParseD2Direction(f *testing.F) {
	fuzzTestParseEnum(f, d2DirectionValues, ParseD2Direction)
}

func FuzzParseD2NodeShape(f *testing.F) {
	fuzzTestParseEnum(f, d2NodeShapeValues, ParseD2NodeShape)
}

func FuzzParseD2ArrowType(f *testing.F) {
	fuzzTestParseEnum(f, d2ArrowTypeValues, ParseD2ArrowType)
}

func FuzzParseD2Constraint(f *testing.F) {
	fuzzTestParseEnum(f, allD2Constraints, ParseD2Constraint)
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
		shape, err := output.ParseNodeShape(shapeStr)
		if err != nil {
			return
		}

		node := graphtest.NewTestNodeWithShape(id, label, shape)

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

// FuzzD2NodeStyleRendering verifies that arbitrary style values (Fill, Stroke,
// FontColor, TextTransform) are escaped through the D2 render pipeline. The
// invariant: newlines in the input must be escaped (never create new lines in
// the output) — this prevents D2 statement injection.
func FuzzD2NodeStyleRendering(f *testing.F) {
	f.Add("red")
	f.Add(`red"inject`)
	f.Add("red\nstyle.fill: green")
	f.Add(`red\malicious`)
	f.Add("#ff0000")
	f.Add("")

	f.Fuzz(func(t *testing.T, styleVal string) {
		d := NewD2Diagram()
		d.AddNode(D2Node{ //nolint:exhaustruct // Fuzz test uses minimal fields
			ID:    output.NewBrandedID[output.D2NodeIDBrand]("n"),
			Label: output.NewBrandedID[output.D2NodeLabelBrand]("Test"),
			Style: D2NodeStyle{
				Fill: styleVal,
				D2StrokeStyle: D2StrokeStyle{
					Stroke:    styleVal,
					FontColor: styleVal,
				},
				TextTransform: D2TextTransform(styleVal),
			},
		})

		got, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		// Count newlines in the raw value. Each unescaped newline in the
		// style value would create additional lines in the output, enabling
		// D2 statement injection. The escape function converts \n to literal
		// \n (backslash-n), so the rendered output should have ZERO newlines
		// contributed by the style value.
		rawNewlines := strings.Count(styleVal, "\n")
		if rawNewlines > 0 {
			// The node itself produces a fixed number of lines. If the style
			// value contributed any newlines, the output would have extra
			// lines beyond what the escaped output should produce.
			escapedNewlines := strings.Count(got, "\n")
			// Without the style value, a single-node D2 diagram with 4 style
			// fields produces exactly: node header(1) + fill(1) + stroke(1)
			// + font-color(1) + text-transform(1) + closing brace(1) = 6 lines
			// → 5 newlines. If styleVal added raw newlines, we'd see more.
			expectedNewlines := 5
			if escapedNewlines != expectedNewlines {
				t.Errorf(
					"style value with %d newlines produced %d output newlines, expected %d (newlines not escaped):\n%s",
					rawNewlines, escapedNewlines, expectedNewlines, got,
				)
			}
		}
	})
}
