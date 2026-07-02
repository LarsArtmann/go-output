package daghtml

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

func TestGolden_Render_SimpleDAG(t *testing.T) {
	t.Parallel()

	dag := DAG{
		Nodes: []Node{
			{ID: "config", Label: "config", Color: "var(--success)"},
			{ID: "db", Label: "db", Color: "var(--accent)"},
			{ID: "cache", Label: "cache", Color: "var(--warning)"},
		},
		Edges: []Edge{
			{From: "db", To: "config"},
			{From: "cache", To: "config"},
		},
	}

	got, err := Render(dag, WithTitle("Test DAG"), WithHeight(400))
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Render_WithErrorNode(t *testing.T) {
	t.Parallel()

	dag := DAG{
		Nodes: []Node{
			{ID: "ok", Label: "ok-svc", Color: "var(--success)"},
			{ID: "fail", Label: "fail-svc", Color: "var(--error)", Error: true, Tooltip: "connection refused"},
		},
		Edges: []Edge{
			{From: "fail", To: "ok"},
		},
	}

	got, err := Render(dag, WithTitle("Error DAG"))
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Render_EmptyDAG(t *testing.T) {
	t.Parallel()

	got, err := Render(DAG{}, WithTitle("Empty"))
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_GraphHTML_EmbeddedSection(t *testing.T) {
	t.Parallel()

	dag := DAG{
		Nodes: []Node{
			{ID: "a", Label: "Alpha", Color: "#e8a838"},
			{ID: "b", Label: "Beta", Color: "#38e8a8"},
		},
		Edges: []Edge{
			{From: "b", To: "a"},
		},
	}

	got, err := GraphHTML(dag, WithContainerID("my-graph"), WithDataID("my-data"))
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
