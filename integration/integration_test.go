// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

type TestProject struct {
	Name       string
	Health     int
	Complexity int
}

func SampleProjects() []TestProject {
	return []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
	}
}

// SampleProject returns a single sample project for testing.
func SampleProject() []TestProject {
	return []TestProject{{Name: "Alpha", Health: 90, Complexity: 7}}
}

// sharedTestData contains common test data used across workflow tests.
func sharedTestData() (headers []string, rows [][]string) {
	return []string{"Name", "Value"}, [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}
}

func TestAllFormatsRender(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()

	formats := []output.Format{
		output.FormatTable,
		output.FormatJSON,
		output.FormatMarkdown,
		output.FormatCSV,
		output.FormatTSV,
		output.FormatXML,
		output.FormatYAML,
		output.FormatHTML,
		output.FormatTree,
		output.FormatD2,
		output.FormatMermaid,
		output.FormatDOT,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			result := renderProject(format, projects)
			if result == "" {
				t.Errorf("Format %s returned empty output", format)
			}
		})
	}
}

func TestStreamingRenderer(t *testing.T) {
	r := output.NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Alpha"})

	var buf bytes.Buffer

	err := r.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	result := buf.String()
	assertContains(t, result, "<table", "Streaming HTML should contain table tag")
}

func TestTableDataRowEdges(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Row0"})
	data.AddRow([]string{"Row1"})
	data.AddRow([]string{"Row2"})

	edges := data.CreateRowEdges()
	if len(edges) != 2 {
		t.Errorf("Expected 2 edges for 3 rows, got %d", len(edges))
	}
}

func TestTreeNodeDepth(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	grandchild := output.NewTreeNode("grandchild", "Grandchild")

	root.AddChild(child)
	child.AddChild(grandchild)

	output.AssertTreeNodeDepth(t, root, child, grandchild)
}

// renderProject renders projects in the specified format.
//
//nolint:cyclop // Complexity is inherent to format handling
func renderProject(format output.Format, projects []TestProject) string {
	switch format {
	case output.FormatTable:
		return renderTableFormat(projects)
	case output.FormatJSON:
		return renderJSONFormat(projects)
	case output.FormatMarkdown:
		return renderMarkdownFormat(projects)
	case output.FormatCSV:
		return renderCSVFormat(projects)
	case output.FormatTSV:
		return renderTSVFormat(projects)
	case output.FormatXML:
		return renderXMLFormat(projects)
	case output.FormatYAML:
		return renderYAMLFormat(projects)
	case output.FormatHTML:
		return renderHTMLFormat(projects)
	case output.FormatTree:
		return renderTreeFormat(projects)
	case output.FormatD2:
		return renderD2Format(projects)
	case output.FormatMermaid:
		return renderMermaidFormat(projects)
	case output.FormatDOT:
		return renderDOTFormat(projects)
	default:
		return ""
	}
}

func renderTableFormat(projects []TestProject) string {
	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, formatHealth(p.Health), formatComplexity(p.Complexity))
	}

	out, err := tbl.Render()
	if err != nil {
		return ""
	}

	return out
}

func renderJSONFormat(projects []TestProject) string {
	data, _ := output.MarshalJSONIndent(projects, "", "  ")

	return string(data)
}

func renderMarkdownFormat(projects []TestProject) string {
	headers := []string{"Name", "Health", "Complexity"}

	return renderMarkdownTable(headers, formatProjectsToRows(projects))
}

func renderCSVFormat(projects []TestProject) string {
	var buf bytes.Buffer

	w := output.NewCSVWriter(&buf)

	_ = w.WriteHeader([]string{"Name", "Health", "Complexity"})
	for _, row := range formatProjectsToRows(projects) {
		_ = w.WriteRow(row)
	}

	w.Flush()

	return buf.String()
}

func renderTSVFormat(projects []TestProject) string {
	var buf bytes.Buffer

	w := output.NewTSVWriter(&buf)

	_ = w.WriteHeader([]string{"Name", "Health", "Complexity"})
	for _, row := range formatProjectsToRows(projects) {
		_ = w.WriteRow(row)
	}

	w.Flush()

	return buf.String()
}

func renderXMLFormat(projects []TestProject) string {
	data, _ := output.MarshalXMLFromTableData(&output.TableData{
		Headers: []string{"Name", "Health", "Complexity"},
		Rows:    formatProjectsToRows(projects),
	})

	return string(data)
}

func formatProjectsToRows(projects []TestProject) [][]string {
	rows := make([][]string, len(projects))
	for i, p := range projects {
		rows[i] = []string{p.Name, formatHealth(p.Health), formatComplexity(p.Complexity)}
	}

	return rows
}

func renderYAMLFormat(projects []TestProject) string {
	data, _ := output.MarshalYAML(projects)

	return string(data)
}

func renderHTMLFormat(projects []TestProject) string {
	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health", "Complexity"})

	for _, row := range formatProjectsToRows(projects) {
		html.AddRow(row)
	}

	out, err := html.Render()
	if err != nil {
		return ""
	}

	return out
}

func buildProjectTree(projects []TestProject) *output.TreeNode {
	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		root.AddChild(output.NewTreeNode(p.Name, p.Name))
	}

	return root
}

func renderTreeFormat(projects []TestProject) string {
	tree := output.NewASCIITreeRenderer()
	tree.SetRoot(buildProjectTree(projects))

	out, err := tree.Render()
	if err != nil {
		return ""
	}

	return out
}

func renderD2Format(projects []TestProject) string {
	d2 := output.NewD2Diagram()
	d2.AddTable("projects", []output.D2Column{
		{Name: "name", Type: "string"},
	})

	for _, p := range projects {
		d2.AddNodeWithShape(p.Name, p.Name, output.D2ShapeCircle)
	}

	out, err := d2.Render()
	if err != nil {
		return ""
	}

	return out
}

func renderD2FromTableData(projects []TestProject) string {
	data := newGraphTableData(projects)

	out, err := output.D2FromTableData(data).Render()
	if err != nil {
		return ""
	}

	return out
}

func renderD2FromTree(projects []TestProject) string {
	out, err := output.D2FromTree(buildProjectTree(projects)).Render()
	if err != nil {
		return ""
	}

	return out
}

func newGraphTableData(projects []TestProject) *output.TableData {
	data := output.NewTableData([]string{"Name"})
	for _, p := range projects {
		data.AddRow([]string{p.Name})
	}

	return data
}

func renderDOTFormat(projects []TestProject) string {
	out, err := output.DOTFromTableData(newGraphTableData(projects)).Render()
	if err != nil {
		return ""
	}

	return out
}

func renderMermaidFormat(projects []TestProject) string {
	out, err := output.MermaidFlowchartRenderer(newGraphTableData(projects)).Render()
	if err != nil {
		return ""
	}

	return out
}

func formatHealth(h int) string {
	return fmt.Sprintf("%d%%", h)
}

func formatComplexity(c int) string {
	return fmt.Sprintf("%d/10", c)
}
