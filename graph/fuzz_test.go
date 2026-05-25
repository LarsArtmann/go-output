package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

func FuzzDOTEscape(f *testing.F) {
	f.Add("hello")
	f.Add(`"quoted"`)
	f.Add(`back\slash`)
	f.Add("new\nline")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		escaped := escape.DOT(input)

		if strings.Contains(input, `"`) && !strings.Contains(escaped, `\"`) {
			t.Errorf("DOT(%q) = %q, quotes not escaped", input, escaped)
		}

		if strings.Contains(input, `\`) && !strings.Contains(escaped, `\\`) {
			t.Errorf("DOT(%q) = %q, backslashes not escaped", input, escaped)
		}
	})
}

func FuzzMermaidTextEscape(f *testing.F) {
	f.Add("hello")
	f.Add(`"quoted"`)
	f.Add("brackets[here]")
	f.Add("braces{here}")
	f.Add("new\nline")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		escaped := escape.MermaidText(input)

		assertEscaped(t, input, escaped, `"`, "double quotes", "MermaidText")
		assertEscaped(t, input, escaped, "[", "brackets", "MermaidText")
	})
}

func assertEscaped(t *testing.T, input, escaped, char, desc, fn string) {
	t.Helper()

	if strings.Contains(input, char) && strings.Contains(escaped, char) {
		t.Errorf("%s(%q) = %q, %s not escaped", fn, input, escaped, desc)
	}
}

func FuzzMermaidID(f *testing.F) {
	f.Add("valid_id123")
	f.Add("has spaces")
	f.Add("special!@#$chars")
	f.Add("")
	f.Add("café")

	f.Fuzz(func(t *testing.T, input string) {
		result := escape.MermaidID(input)

		for _, r := range result {
			isValid := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '_'

			if !isValid {
				t.Errorf("MermaidID(%q) = %q, invalid rune %c", input, result, r)
			}
		}
	})
}

func fuzzTestNodes(idA, labelA, idB, labelB string) []output.GraphNode {
	return []output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](idA),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](labelA),
		},
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](idB),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](labelB),
		},
	}
}

func FuzzDOTRendererRender(f *testing.F) {
	f.Add("A", "Node A", "B", "Node B")
	f.Add("", "", "", "")
	f.Add(`"quoted"`, `back\slash`, "normal", "ok")

	f.Fuzz(func(t *testing.T, idA, labelA, idB, labelB string) {
		renderer := NewDOTRenderer()
		renderer.SetNodes(fuzzTestNodes(idA, labelA, idB, labelB))

		got, err := renderer.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		if !strings.Contains(got, "digraph G {") {
			t.Error("DOT output missing digraph header")
		}
	})
}

func FuzzMermaidRendererRender(f *testing.F) {
	f.Add("A", "Node A", "B", "Node B")
	f.Add("node1", "simple", "node2", "test")
	f.Add("", "", "", "")

	f.Fuzz(func(t *testing.T, idA, labelA, idB, labelB string) {
		renderer := NewMermaidRenderer()
		renderer.SetNodes(fuzzTestNodes(idA, labelA, idB, labelB))

		got, err := renderer.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		if !strings.Contains(got, "```mermaid") {
			t.Error("Mermaid output missing fence")
		}
	})
}
