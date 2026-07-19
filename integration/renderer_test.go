// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/tree"
)

func TestTableFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()

	result := renderTableFormat(projects)
	testhelpers.AssertContains(t, result, "Name", "Table should contain header 'Name'")
	testhelpers.AssertContains(t, result, "Alpha", "Table should contain project name 'Alpha'")
	testhelpers.AssertContains(t, result, "Beta", "Table should contain project name 'Beta'")
}

func TestJSONFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProject()

	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent failed: %v", err)
	}

	result := string(data)
	testhelpers.AssertContains(t, result, "Alpha", "JSON should contain project name 'Alpha'")
	testhelpers.AssertContains(t, result, "90", "JSON should contain health value 90")
}

func TestMarkdownTableContent(t *testing.T) {
	t.Parallel()

	result := renderSampleMarkdownTable()

	testhelpers.AssertContains(t, result, "| Name", "Markdown should contain header cell")
	testhelpers.AssertContains(t, result, "| Alpha", "Markdown should contain row data")
	testhelpers.AssertContains(t, result, "|---", "Markdown should contain separator row")
}

func TestCSVFormatContent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := delimited.NewCSVWriter(&buf)

	err := w.WriteHeader([]string{"Name", "Health"})
	if err != nil {
		t.Fatalf("WriteHeader failed: %v", err)
	}

	err = w.WriteRow([]string{"Alpha", "90"})
	if err != nil {
		t.Fatalf("WriteRow failed: %v", err)
	}

	w.Flush()

	result := buf.String()
	testhelpers.AssertContains(t, result, "Name,Health", "CSV should contain header row")
	testhelpers.AssertContains(t, result, "Alpha,90", "CSV should contain data row")
}

func TestYAMLFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProject()

	data, err := serialization.MarshalYAML(projects)
	if err != nil {
		t.Fatalf("MarshalYAML failed: %v", err)
	}

	result := string(data)
	testhelpers.AssertContains(t, result, "Alpha", "YAML should contain project name 'Alpha'")
}

func TestHTMLFormatContent(t *testing.T) {
	t.Parallel()

	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health"})
	html.AddRow([]string{"Alpha", "90%"})

	testhelpers.RenderAssert(t, html, "<table", "Alpha")
}

func TestHTMLFullPage(t *testing.T) {
	t.Parallel()

	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name"})
	html.AddRow([]string{"Test"})

	result, err := html.RenderFullHTML("Test Page")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	testhelpers.AssertContains(t, result, "<html", "Full HTML should contain html tag")
	testhelpers.AssertContains(
		t,
		result,
		"<title>Test Page</title>",
		"Full HTML should contain title",
	)
}

func TestTreeFormatContent(t *testing.T) {
	t.Parallel()

	tree := tree.NewASCIITreeRenderer()
	root := output.NewTreeNode("root", "Projects")
	child := output.NewTreeNode("child1", "Alpha")
	root.AddChild(child)
	tree.SetRoot(root)

	result, err := tree.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if result == "" {
		t.Error("Tree render should not be empty")
	}
}

func TestD2FormatContent(t *testing.T) {
	t.Parallel()

	d2Diagram := d2.NewDiagram()
	d2Diagram.AddTable("test", []d2.Column{
		{Name: "name", Type: "string"},
	})

	result, err := d2Diagram.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if result == "" {
		t.Error("D2 render should not be empty")
	}
}

// testRendererNotEmpty tests that a renderer produces non-empty output.
func testRendererNotEmpty[R output.Renderer](
	t *testing.T,
	createRenderer func(*output.Table) R,
	name string,
) {
	t.Helper()

	data := output.NewTable([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	r := createRenderer(data)

	result, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if result == "" {
		t.Error(name + " render should not be empty")
	}
}

func TestMermaidFormatContent(t *testing.T) {
	t.Parallel()

	testRendererNotEmpty(t, graph.NewMermaidFromTable, "Mermaid")
}

func TestDOTFormatContent(t *testing.T) {
	t.Parallel()

	testRendererNotEmpty(t, graph.NewDOTFromTable, "DOT")
}
