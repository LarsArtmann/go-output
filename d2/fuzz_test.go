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
