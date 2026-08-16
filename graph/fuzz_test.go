package graph

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

func FuzzDOTEscape(f *testing.F) {
	f.Add("hello")
	f.Add(`"quoted"`)
	f.Add(`back\slash`)
	f.Add("new\nline")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		escaped := escape.DOT(input)

		graphtest.AssertEscape(t, "DOT", input, escaped, `"`, "quotes")
		graphtest.AssertEscape(t, "DOT", input, escaped, `\`, "backslashes")
	})
}

func FuzzMermaidTextEscape(f *testing.F) {
	f.Add("hello")
	f.Add(`"quoted"`)
	f.Add("brackets[here]")
	f.Add("braces{here}")
	f.Add("new\nline")
	f.Add("<script>alert('xss')</script>")
	f.Add("#60;script#62;")
	f.Add("#ff0000")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		escaped := escape.MermaidText(input)

		assertEscaped(t, input, escaped, `"`, "double quotes", "MermaidText")
		assertEscaped(t, input, escaped, "[", "brackets", "MermaidText")
		assertEscaped(t, input, escaped, "<", "angle brackets", "MermaidText")
		assertEscaped(t, input, escaped, ">", "angle brackets", "MermaidText")
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

		if escape.MermaidID(result) != result {
			t.Errorf("MermaidID(%q) = %q, not idempotent", input, result)
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

// FuzzDOTNodeStyleNewlines verifies that newlines in DOT style values (Fill,
// Stroke) never create additional lines in the output. Unescaped newlines
// would allow attribute injection in DOT syntax.
func FuzzDOTNodeStyleNewlines(f *testing.F) {
	f.Add("red")
	f.Add("red\ninjected_attr")
	f.Add("#000\nline:evil")
	f.Add("")

	f.Fuzz(func(t *testing.T, styleVal string) {
		renderer := NewDOTRenderer()
		renderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Fuzz test
			{
				ID:    output.NewBrandedID[output.GraphNodeIDBrand]("n"),
				Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
				Shape: output.NodeShapeBox,
				Style: output.NodeStyle{
					Fill:   styleVal,
					Stroke: styleVal,
				},
			},
		})

		got, err := renderer.Render()
		if err != nil {
			t.Fatalf("Render() error: %v", err)
		}

		// If the style value has newlines, verify they don't create extra
		// output lines. Compare against the same renderer with a safe value.
		rawNewlines := strings.Count(styleVal, "\n")
		if rawNewlines > 0 {
			safeRenderer := NewDOTRenderer()
			safeRenderer.SetNodes([]output.GraphNode{ //nolint:exhaustruct // Fuzz test
				{
					ID:    output.NewBrandedID[output.GraphNodeIDBrand]("n"),
					Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
					Shape: output.NodeShapeBox,
					Style: output.NodeStyle{
						Fill:   "safe",
						Stroke: "safe",
					},
				},
			})

			safeOut, err := safeRenderer.Render()
			if err != nil {
				t.Fatalf("safe Render() error: %v", err)
			}

			maliciousLines := strings.Count(got, "\n")
			safeLines := strings.Count(safeOut, "\n")

			if maliciousLines != safeLines {
				t.Errorf(
					"style value with %d newlines produced %d output lines, expected %d (newlines not escaped)",
					rawNewlines, maliciousLines, safeLines,
				)
			}
		}
	})
}
