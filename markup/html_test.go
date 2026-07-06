package markup

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestHTMLRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name", "Age"})
	renderer.AddRow([]string{"Alice", "30"})
	renderer.AddRow([]string{"Bob", "25"})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "<table", "Output should contain <table>")
	assertContains(t, out, "</table>", "Output should contain </table>")
	assertContains(t, out, "<th>", "Output should contain <th> for headers")
	assertContains(t, out, "<td>", "Output should contain <td> for data cells")
	assertContains(t, out, "Alice", "Output should contain 'Alice'")
	assertContains(t, out, "Bob", "Output should contain 'Bob'")
}

func TestHTMLRendererFullDocument(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Col1", "Col2"})
	renderer.AddRow([]string{"A", "B"})

	out, err := renderer.RenderFullHTML("Test Title")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	assertContains(t, out, "<!DOCTYPE html>", "Full HTML should contain DOCTYPE")
	assertContains(t, out, "<title>Test Title</title>", "Full HTML should contain title")
	assertContains(t, out, "<table", "Full HTML should contain table")
}

func TestHTMLRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	testEmptyRendererOutput(t, renderer, testHTMLEmptyExpected())
}

func TestHTMLRendererEscaping(t *testing.T) {
	t.Parallel()

	testHTMLEscape(t, func() htmlEscapeTestRenderer {
		return NewHTMLRenderer()
	}, "HTMLRenderer")
}

func TestHTMLTreeRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()

	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child", "Child"))
	renderer.SetRoot(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "<ul", "Output should contain <ul>")
	assertContains(t, out, "<li>", "Output should contain <li>")
	assertContains(t, out, "Root", "Output should contain 'Root'")
	assertContains(t, out, "Child", "Output should contain 'Child'")
}

func TestHTMLTreeRendererFullDocument(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()
	renderer.SetRoot(output.NewTreeNode("root", "Test Tree"))

	out, err := renderer.RenderFullHTML("Tree Title")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	assertContains(t, out, "<!DOCTYPE html>", "Full HTML should contain DOCTYPE")
	assertContains(t, out, "<title>Tree Title</title>", "Full HTML should contain title")
}

func TestHTMLRendererSetData(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetData(&output.Table{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "2"}},
	})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "<th>A", "Output should contain header 'A'")
	assertContains(t, out, "<td>1", "Output should contain cell '1'")
}

func TestHTMLRendererWithFooter(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name", "Count"})
	renderer.AddRow([]string{"Alice", "10"})
	renderer.SetFooter([]string{"Total", "10"})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "<tfoot>", "Output should contain <tfoot>")
	assertContains(t, out, `<td class="footer-cell">Total</td>`, "Output should contain footer cell with class")
	assertContains(t, out, "</tfoot>", "Output should contain </tfoot>")

	if strings.Contains(out, "<tfoot>") && !strings.Contains(out, "</tbody>") {
		t.Error("<tfoot> should come after </tbody>")
	}
}

func TestHTMLRendererNoFooter(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name"})
	renderer.AddRow([]string{"Alice"})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(out, "<tfoot>") {
		t.Error("Output should not contain <tfoot> when no footer is set")
	}
}

func TestHTMLRendererAddRowWithoutSetHeaders(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.AddRow([]string{"test"})

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "test", "Output should contain 'test'")
}

func TestHTMLTreeRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "<ul", "Empty tree should contain <ul>")

	if strings.Contains(out, "<li>") {
		t.Error("Empty tree should not contain <li>")
	}
}

func TestRenderHTMLWithStylesError(t *testing.T) {
	t.Parallel()

	_, err := renderHTMLWithStyles(
		&errorRenderer{}, "title", "styles", "test context",
	)
	if err == nil {
		t.Fatal("expected error from errorRenderer")
	}

	assertContains(t, err.Error(), "test context", "error should mention context")
}
