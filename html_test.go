package output

import (
	"strings"
	"testing"
)

func TestHTMLRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name", "Age"})
	renderer.AddRow([]string{"Alice", "30"})
	renderer.AddRow([]string{"Bob", "25"})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "<table", "Output should contain <table>")
	assertContains(t, output, "</table>", "Output should contain </table>")
	assertContains(t, output, "<th>", "Output should contain <th> for headers")
	assertContains(t, output, "<td>", "Output should contain <td> for data cells")
	assertContains(t, output, "Alice", "Output should contain 'Alice'")
	assertContains(t, output, "Bob", "Output should contain 'Bob'")
}

func TestHTMLRendererFullDocument(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Col1", "Col2"})
	renderer.AddRow([]string{"A", "B"})

	output, err := renderer.RenderFullHTML("Test Title")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	assertContains(t, output, "<!DOCTYPE html>", "Full HTML should contain DOCTYPE")
	assertContains(t, output, "<title>Test Title</title>", "Full HTML should contain title")
	assertContains(t, output, "<table", "Full HTML should contain table")
}

func TestHTMLRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	testEmptyRendererOutput(t, renderer, testHTMLEmptyExpected())
}

func TestHTMLRendererEscaping(t *testing.T) {
	t.Parallel()
	testHTMLEscapeShared(
		t,
		func() htmlEscapeTestRenderer { return NewHTMLRenderer() },
		"HTMLRenderer",
	)
}

func TestHTMLTreeRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()

	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child", "Child"))
	renderer.SetRoot(root)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "<ul", "Output should contain <ul>")
	assertContains(t, output, "<li>", "Output should contain <li>")
	assertContains(t, output, "Root", "Output should contain 'Root'")
	assertContains(t, output, "Child", "Output should contain 'Child'")
}

func TestHTMLTreeRendererFullDocument(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()
	renderer.SetRoot(NewTreeNode("root", "Test Tree"))

	output, err := renderer.RenderFullHTML("Tree Title")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	assertContains(t, output, "<!DOCTYPE html>", "Full HTML should contain DOCTYPE")
	assertContains(t, output, "<title>Tree Title</title>", "Full HTML should contain title")
}

func TestHTMLRendererSetData(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetData(&TableData{
		Headers: []string{"A", "B"},
		Rows:    [][]string{{"1", "2"}},
	})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "<th>A", "Output should contain header 'A'")
	assertContains(t, output, "<td>1", "Output should contain cell '1'")
}

func TestHTMLRendererAddRowWithoutSetHeaders(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	// Call AddRow without first calling SetHeaders - should initialize data
	renderer.AddRow([]string{"test"})

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "test", "Output should contain 'test'")
}

func TestHTMLTreeRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()
	// Don't set root - should return empty tree
	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "<ul", "Empty tree should contain <ul>")

	if strings.Contains(output, "<li>") {
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
