// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestTableFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()

	result := renderTableFormat(projects)
	assertContains(t, result, "Name", "Table should contain header 'Name'")
	assertContains(t, result, "Alpha", "Table should contain project name 'Alpha'")
	assertContains(t, result, "Beta", "Table should contain project name 'Beta'")
}

func TestJSONFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProject()

	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent failed: %v", err)
	}

	result := string(data)
	assertContains(t, result, "Alpha", "JSON should contain project name 'Alpha'")
	assertContains(t, result, "90", "JSON should contain health value 90")
}

func TestMarkdownTableContent(t *testing.T) {
	t.Parallel()

	result := renderSampleMarkdownTable()

	assertContains(t, result, "| Name", "Markdown should contain header cell")
	assertContains(t, result, "| Alpha", "Markdown should contain row data")
	assertContains(t, result, "|---", "Markdown should contain separator row")
}

func TestCSVFormatContent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := output.NewCSVWriter(&buf)

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
	assertContains(t, result, "Name,Health", "CSV should contain header row")
	assertContains(t, result, "Alpha,90", "CSV should contain data row")
}

func TestYAMLFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProject()

	data, err := output.MarshalYAML(projects)
	if err != nil {
		t.Fatalf("MarshalYAML failed: %v", err)
	}

	result := string(data)
	assertContains(t, result, "Alpha", "YAML should contain project name 'Alpha'")
}

func TestHTMLFormatContent(t *testing.T) {
	t.Parallel()

	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health"})
	html.AddRow([]string{"Alpha", "90%"})

	result, err := html.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, result, "<table", "HTML should contain table tag")
	assertContains(t, result, "Alpha", "HTML should contain project name 'Alpha'")
}

func TestHTMLFullPage(t *testing.T) {
	t.Parallel()

	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name"})
	html.AddRow([]string{"Test"})

	result, err := html.RenderFullHTML("Test Page")
	if err != nil {
		t.Fatalf("RenderFullHTML() error = %v", err)
	}

	assertContains(t, result, "<html", "Full HTML should contain html tag")
	assertContains(
		t,
		result,
		"<title>Test Page</title>",
		"Full HTML should contain title",
	)
}

func TestTreeFormatContent(t *testing.T) {
	t.Parallel()

	tree := output.NewASCIITreeRenderer()
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

	d2 := output.NewD2Diagram()
	d2.AddTable("test", []output.D2Column{
		{Name: "name", Type: "string"},
	})

	result, err := d2.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if result == "" {
		t.Error("D2 render should not be empty")
	}
}

type renderer interface{ Render() (string, error) }

// testRendererNotEmpty tests that a renderer produces non-empty output.
func testRendererNotEmpty[R renderer](
	t *testing.T,
	createRenderer func(*output.TableData) R,
	name string,
) {
	t.Helper()

	data := output.NewTableData([]string{"Name", "Health"})
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

	testRendererNotEmpty(t, output.MermaidFlowchartRenderer, "Mermaid")
}

func TestDOTFormatContent(t *testing.T) {
	t.Parallel()

	testRendererNotEmpty(t, output.DOTFromTableData, "DOT")
}
