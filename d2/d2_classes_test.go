package d2

import (
	"strings"
	"testing"
)

func TestD2ClassesDeterministic(t *testing.T) {
	t.Parallel()

	d := NewDiagram()
	d.AddClass("zebra", NodeStyle{Fill: "black"})
	d.AddClass("alpha", NodeStyle{Fill: "red"})
	d.AddClass("beta", NodeStyle{Fill: "green"})
	d.AddNode(newClassNode("n1", "N1", "alpha"))

	got1, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	got2, err := d.Render()
	if err != nil {
		t.Fatalf("Render() second call error = %v", err)
	}

	if got1 != got2 {
		t.Error("D2 output with multiple classes is not deterministic across Render() calls")
	}

	// Verify sorted order: alpha before beta before zebra
	alphaIdx := strings.Index(got1, "alpha:")
	betaIdx := strings.Index(got1, "beta:")
	zebraIdx := strings.Index(got1, "zebra:")

	if alphaIdx >= betaIdx || betaIdx >= zebraIdx {
		t.Errorf("classes not sorted alphabetically: alpha@%d beta@%d zebra@%d", alphaIdx, betaIdx, zebraIdx)
	}
}
