package output

import (
	"testing"
)

func newTestD2Edge(opts ...func(*D2Edge)) D2Edge {
	edge := D2Edge{ //nolint:exhaustruct // Test uses minimal required fields
		From: NewBrandedID[D2NodeIDBrand]("a"),
		To:   NewBrandedID[D2NodeIDBrand]("b"),
	}
	for _, opt := range opts {
		opt(&edge)
	}

	return edge
}

func TestD2EdgeWithArrows(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddEdge(newTestD2Edge(func(e *D2Edge) {
		e.Label = NewBrandedID[D2NodeLabelBrand]("test")
		e.SourceArrow = D2ArrowDiamond
		e.TargetArrow = D2ArrowTriangle
	}))

	got := d.Render()
	assertContains(t, got, "source-arrowhead.shape: diamond", "should contain source arrow")
	assertContains(t, got, "target-arrowhead.shape: triangle", "should contain target arrow")
}

func TestD2AllArrowTypes(t *testing.T) {
	t.Parallel()

	arrows := []D2ArrowType{
		D2ArrowArrow,
		D2ArrowTriangle,
		D2ArrowDiamond,
		D2ArrowCircle,
		D2ArrowFilled,
		D2ArrowBox,
		D2ArrowCross,
		D2ArrowCFOne,
		D2ArrowCFMany,
		D2ArrowCFOneRequired,
		D2ArrowCFManyRequired,
	}

	for _, arrow := range arrows {
		t.Run(string(arrow), func(t *testing.T) {
			t.Parallel()

			d := NewD2Diagram()
			d.AddEdge(newTestD2Edge(func(e *D2Edge) {
				e.TargetArrow = arrow
			}))

			got := d.Render()
			assertContains(t, got, "target-arrowhead.shape: "+string(arrow),
				"should contain arrow type "+string(arrow))
		})
	}
}

func TestD2EdgeWithFilledArrow(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddEdge(newTestD2Edge(func(e *D2Edge) {
		e.TargetArrow = D2ArrowFilled
	}))

	got := d.Render()
	assertContains(t, got, "target-arrowhead.shape: filled", "should contain filled arrow")
}

func TestD2EdgeStyle(t *testing.T) {
	t.Parallel()

	d := NewD2Diagram()
	d.AddEdge(D2Edge{
		From:  NewBrandedID[D2NodeIDBrand]("a"),
		To:    NewBrandedID[D2NodeIDBrand]("b"),
		Label: NewBrandedID[D2NodeLabelBrand]("styled"),
		Style: D2EdgeStyle{
			Stroke:      "red",
			StrokeWidth: 3,
			StrokeDash:  5,
			Animated:    true,
			FontColor:   "blue",
			FontSize:    12,
		},
	})

	got := d.Render()
	assertContains(t, got, "style.stroke: red", "should contain edge stroke")
	assertContains(t, got, "style.stroke-width: 3", "should contain edge stroke-width")
	assertContains(t, got, "style.stroke-dash: 5", "should contain stroke-dash")
	assertContains(t, got, "style.animated: true", "should contain animated")
	assertContains(t, got, "style.font-color: blue", "should contain edge font-color")
	assertContains(t, got, "style.font-size: 12", "should contain edge font-size")
}
