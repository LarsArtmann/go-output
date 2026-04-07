package output

import (
	"fmt"
	"strings"
	"testing"
)

func TestHTMLRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name", "Age"})
	renderer.AddRow([]string{"Alice", "30"})
	renderer.AddRow([]string{"Bob", "25"})

	output := renderer.Render()

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

	output := renderer.RenderFullHTML("Test Title")

	assertContains(t, output, "<!DOCTYPE html>", "Full HTML should contain DOCTYPE")
	assertContains(t, output, "<title>Test Title</title>", "Full HTML should contain title")
	assertContains(t, output, "<table", "Full HTML should contain table")
}

func TestHTMLRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	testEmptyRendererOutput(t, renderer, []ExpectedOutput{
		{Substring: "<table", Message: "Empty table should still be valid HTML"},
		{Substring: "</table>", Message: "Empty table should have closing tag"},
	})
}

// htmlEscapeTestRenderer is an interface for HTML renderers that support escaping tests.
type htmlEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() string
}

// testHTMLEscapeShared is a shared helper for testing HTML escaping across renderer implementations.
func testHTMLEscapeShared(t *testing.T, newRenderer func() htmlEscapeTestRenderer, name string) {
	t.Helper()

	r := newRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"<script>alert('xss')</script>"})

	got := r.Render()

	if strings.Contains(got, "<script>") {
		t.Errorf("%s: Render() should escape script tags", name)
	}

	assertContains(t, got, "&lt;script&gt;", fmt.Sprintf("%s: Render() should contain escaped script tag", name))
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

	output := renderer.Render()

	assertContains(t, output, "<ul", "Output should contain <ul>")
	assertContains(t, output, "<li>", "Output should contain <li>")
	assertContains(t, output, "Root", "Output should contain 'Root'")
	assertContains(t, output, "Child", "Output should contain 'Child'")
}

func TestHTMLTreeRendererFullDocument(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()
	renderer.SetRoot(NewTreeNode("root", "Test Tree"))

	output := renderer.RenderFullHTML("Tree Title")

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

	output := renderer.Render()

	assertContains(t, output, "<th>A", "Output should contain header 'A'")
	assertContains(t, output, "<td>1", "Output should contain cell '1'")
}

func TestHTMLRendererAddRowWithoutSetHeaders(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLRenderer()
	// Call AddRow without first calling SetHeaders - should initialize data
	renderer.AddRow([]string{"test"})

	output := renderer.Render()
	assertContains(t, output, "test", "Output should contain 'test'")
}

func TestHTMLTreeRendererEmpty(t *testing.T) {
	t.Parallel()

	renderer := NewHTMLTreeRenderer()
	// Don't set root - should return empty tree
	output := renderer.Render()

	assertContains(t, output, "<ul", "Empty tree should contain <ul>")
	if strings.Contains(output, "<li>") {
		t.Error("Empty tree should not contain <li>")
	}
}
