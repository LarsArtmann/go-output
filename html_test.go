package output

import (
	"strings"
	"testing"
)

func TestHTMLRenderer(t *testing.T) {
	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name", "Age"})
	renderer.AddRow([]string{"Alice", "30"})
	renderer.AddRow([]string{"Bob", "25"})

	output := renderer.Render()

	// Check for HTML table structure
	if !strings.Contains(output, "<table") {
		t.Error("Output should contain <table>")
	}
	if !strings.Contains(output, "</table>") {
		t.Error("Output should contain </table>")
	}
	if !strings.Contains(output, "<th>") {
		t.Error("Output should contain <th> for headers")
	}
	if !strings.Contains(output, "<td>") {
		t.Error("Output should contain <td> for data cells")
	}
	if !strings.Contains(output, "Alice") {
		t.Error("Output should contain 'Alice'")
	}
	if !strings.Contains(output, "Bob") {
		t.Error("Output should contain 'Bob'")
	}
}

func TestHTMLRendererFullDocument(t *testing.T) {
	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Col1", "Col2"})
	renderer.AddRow([]string{"A", "B"})

	output := renderer.RenderFullHTML("Test Title")

	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("Full HTML should contain DOCTYPE")
	}
	if !strings.Contains(output, "<title>Test Title</title>") {
		t.Error("Full HTML should contain title")
	}
	if !strings.Contains(output, "<table") {
		t.Error("Full HTML should contain table")
	}
}

func TestHTMLRendererEmpty(t *testing.T) {
	renderer := NewHTMLRenderer()
	output := renderer.Render()

	if !strings.Contains(output, "<table") {
		t.Error("Empty table should still be valid HTML")
	}
	if !strings.Contains(output, "</table>") {
		t.Error("Empty table should have closing tag")
	}
}

func TestHTMLRendererEscaping(t *testing.T) {
	renderer := NewHTMLRenderer()
	renderer.SetHeaders([]string{"Name"})
	renderer.AddRow([]string{"<script>alert('xss')</script>"})

	output := renderer.Render()

	// Should escape HTML
	if strings.Contains(output, "<script>") {
		t.Error("Output should escape script tags")
	}
	if !strings.Contains(output, "&lt;script&gt;") {
		t.Error("Output should contain escaped script tag")
	}
}

func TestHTMLTreeRenderer(t *testing.T) {
	renderer := NewHTMLTreeRenderer()

	root := NewTreeNode("root", "Root")
	root.AddChild(NewTreeNode("child", "Child"))
	renderer.SetRoot(root)

	output := renderer.Render()

	if !strings.Contains(output, "<ul") {
		t.Error("Output should contain <ul>")
	}
	if !strings.Contains(output, "<li>") {
		t.Error("Output should contain <li>")
	}
	if !strings.Contains(output, "Root") {
		t.Error("Output should contain 'Root'")
	}
	if !strings.Contains(output, "Child") {
		t.Error("Output should contain 'Child'")
	}
}

func TestHTMLTreeRendererFullDocument(t *testing.T) {
	renderer := NewHTMLTreeRenderer()
	renderer.SetRoot(NewTreeNode("root", "Test Tree"))

	output := renderer.RenderFullHTML("Tree Title")

	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("Full HTML should contain DOCTYPE")
	}
	if !strings.Contains(output, "<title>Tree Title</title>") {
		t.Error("Full HTML should contain title")
	}
}
