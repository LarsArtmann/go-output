package d2

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

func TestGolden_D2_SimpleDiagram(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.AddNodeSimple("a", "Alpha")
	d.AddNodeSimple("b", "Beta")
	d.AddNodeSimple("c", "Gamma")
	d.AddEdgeSimple("a", "b")
	d.AddEdgeSimple("b", "c")

	got, err := d.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_D2_ShapedNodes(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.SetDirection(DirRight)
	d.AddNodeWithShape("start", "Start", ShapeCircle)
	d.AddNodeWithShape("process", "Process", ShapeRectangle)
	d.AddNodeWithShape("decision", "Decision", ShapeDiamond)
	d.AddNodeWithShape("end", "End", ShapeCircle)
	d.AddEdgeSimple("start", "process")
	d.AddEdgeSimple("process", "decision")
	d.AddEdgeSimple("decision", "end")

	got, err := d.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_D2_Empty(t *testing.T) {
	t.Parallel()

	d := NewDiagram()

	got, err := d.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
