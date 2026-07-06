package d2

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func newTestEdge(opts ...func(*Edge)) Edge {
	edge := Edge{ //nolint:exhaustruct // Test uses minimal required fields
		From: output.NewBrandedID[output.D2NodeIDBrand]("a"),
		To:   output.NewBrandedID[output.D2NodeIDBrand]("b"),
	}
	for _, opt := range opts {
		opt(&edge)
	}

	return edge
}

func TestD2EdgeWithArrows(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.AddEdge(newTestEdge(func(e *Edge) {
		e.Label = output.NewBrandedID[output.D2NodeLabelBrand]("test")
		e.SourceArrow = ArrowDiamond
		e.TargetArrow = ArrowTriangle
	}))

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(
		t,
		got,
		"source-arrowhead.shape: diamond",
		"should contain source arrow",
	)
	testhelpers.AssertContains(
		t,
		got,
		"target-arrowhead.shape: triangle",
		"should contain target arrow",
	)
}

func TestD2AllArrowTypes(t *testing.T) {
	t.Parallel()

	arrows := []ArrowType{
		ArrowArrow,
		ArrowTriangle,
		ArrowDiamond,
		ArrowCircle,
		ArrowFilled,
		ArrowBox,
		ArrowCross,
		ArrowCFOne,
		ArrowCFMany,
		ArrowCFOneRequired,
		ArrowCFManyRequired,
	}

	for _, arrow := range arrows {
		t.Run(string(arrow), func(t *testing.T) {
			t.Parallel()

			d := NewDiagram()
			d.AddEdge(newTestEdge(func(e *Edge) {
				e.TargetArrow = arrow
			}))

			got, err := d.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			testhelpers.AssertContains(t, got, "target-arrowhead.shape: "+string(arrow),
				"should contain arrow type "+string(arrow))
		})
	}
}

func TestD2EdgeWithFilledArrow(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.AddEdge(newTestEdge(func(e *Edge) {
		e.TargetArrow = ArrowFilled
	}))

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(
		t,
		got,
		"target-arrowhead.shape: filled",
		"should contain filled arrow",
	)
}

func TestEdgeStyle(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.AddEdge(Edge{
		From:  output.NewBrandedID[output.D2NodeIDBrand]("a"),
		To:    output.NewBrandedID[output.D2NodeIDBrand]("b"),
		Label: output.NewBrandedID[output.D2NodeLabelBrand]("styled"),
		Style: EdgeStyle{
			StrokeStyle: StrokeStyle{
				Stroke:      "red",
				StrokeWidth: 3,
				StrokeDash:  5,
				FontColor:   "blue",
				FontSize:    12,
			},
			Animated: true,
		},
	})

	got, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "style.stroke: red", "should contain edge stroke")
	testhelpers.AssertContains(t, got, "style.stroke-width: 3", "should contain edge stroke-width")
	testhelpers.AssertContains(t, got, "style.stroke-dash: 5", "should contain stroke-dash")
	testhelpers.AssertContains(t, got, "style.animated: true", "should contain animated")
	testhelpers.AssertContains(t, got, "style.font-color: blue", "should contain edge font-color")
	testhelpers.AssertContains(t, got, "style.font-size: 12", "should contain edge font-size")
}
