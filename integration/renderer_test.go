// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

func TestTableFormatContent(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()

	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, formatHealth(p.Health), formatComplexity(p.Complexity))
	}

	result := tbl.Render()
	if !strings.Contains(result, "Name") {
		t.Error("Table should contain header 'Name'")
	}

	if !strings.Contains(result, "Alpha") {
		t.Error("Table should contain project name 'Alpha'")
	}

	if !strings.Contains(result, "Beta") {
		t.Error("Table should contain project name 'Beta'")
	}
}

func TestJSONFormatContent(t *testing.T) {
	t.Parallel()

	projects := []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
	}

	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent failed: %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "Alpha") {
		t.Error("JSON should contain project name 'Alpha'")
	}

	if !strings.Contains(result, "90") {
		t.Error("JSON should contain health value 90")
	}
}

func TestMarkdownTableContent(t *testing.T) {
	t.Parallel()

	result := renderMarkdownTable([]string{"Name", "Health"}, [][]string{{"Alpha", "90%"}})

	if !strings.Contains(result, "| Name") {
		t.Error("Markdown should contain header cell")
	}

	if !strings.Contains(result, "| Alpha") {
		t.Error("Markdown should contain row data")
	}

	if !strings.Contains(result, "|---") {
		t.Error("Markdown should contain separator row")
	}
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
	if !strings.Contains(result, "Name,Health") {
		t.Error("CSV should contain header row")
	}

	if !strings.Contains(result, "Alpha,90") {
		t.Error("CSV should contain data row")
	}
}

func TestYAMLFormatContent(t *testing.T) {
	t.Parallel()

	projects := []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
	}

	data, err := output.MarshalYAML(projects)
	if err != nil {
		t.Fatalf("MarshalYAML failed: %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "Alpha") {
		t.Error("YAML should contain project name 'Alpha'")
	}
}

func TestHTMLFormatContent(t *testing.T) {
	t.Parallel()

	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health"})
	html.AddRow([]string{"Alpha", "90%"})

	result := html.Render()
	if !strings.Contains(result, "<table") {
		t.Error("HTML should contain table tag")
	}

	if !strings.Contains(result, "Alpha") {
		t.Error("HTML should contain project name 'Alpha'")
	}
}

func TestHTMLFullPage(t *testing.T) {
	t.Parallel()

	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name"})
	html.AddRow([]string{"Test"})

	result := html.RenderFullHTML("Test Page")
	if !strings.Contains(result, "<html") {
		t.Error("Full HTML should contain html tag")
	}

	if !strings.Contains(result, "<title>Test Page</title>") {
		t.Error("Full HTML should contain title")
	}
}

func TestTreeFormatContent(t *testing.T) {
	t.Parallel()

	tree := output.NewASCIITreeRenderer()
	root := output.NewTreeNode("root", "Projects")
	child := output.NewTreeNode("child1", "Alpha")
	root.AddChild(child)
	tree.SetRoot(root)

	result := tree.Render()
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

	result := d2.Render()
	if result == "" {
		t.Error("D2 render should not be empty")
	}
}

type renderer interface{ Render() string }

// testRendererNotEmpty tests that a renderer produces non-empty output.
func testRendererNotEmpty[R renderer](t *testing.T, createRenderer func(*output.TableData) R, name string) {
	t.Helper()

	data := output.NewTableData([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	r := createRenderer(data)
	result := r.Render()

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
