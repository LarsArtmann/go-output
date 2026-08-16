package integration

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
	"github.com/larsartmann/go-output/tree"
)

//nolint:cyclop // Complexity is inherent to format handling
func renderProject(t *testing.T, format output.Format, projects []TestProject) string {
	t.Helper()

	switch format {
	case output.FormatTable:
		return renderTableFormat(t, projects)
	case output.FormatJSON:
		return renderJSONFormat(t, projects)
	case output.FormatMarkdown:
		return renderMarkdownFormat(t, projects)
	case output.FormatCSV:
		return renderCSVFormat(t, projects)
	case output.FormatTSV:
		return renderTSVFormat(t, projects)
	case output.FormatXML:
		return renderXMLFormat(t, projects)
	case output.FormatYAML:
		return renderYAMLFormat(t, projects)
	case output.FormatHTML:
		return renderHTMLFormat(t, projects)
	case output.FormatTree:
		return renderTreeFormat(t, projects)
	case output.FormatD2:
		return renderD2Format(t, projects)
	case output.FormatMermaid:
		return renderMermaidFormat(t, projects)
	case output.FormatDOT:
		return renderDOTFormat(t, projects)
	case output.FormatJSONL:
		return renderJSONLFormat(t, projects)
	case output.FormatAsciiDoc:
		return renderAsciiDocFormat(t, projects)
	case output.FormatTOML:
		return renderTOMLFormat(t, projects)
	case output.FormatPlantUML:
		return renderPlantUMLFormat(t, projects)
	default:
		t.Fatalf("unhandled format %s in renderProject", format)

		return ""
	}
}

// mustRender fails the test on a render error, so content assertions report
// the real cause instead of degrading to "returned empty output".
func mustRender(t *testing.T, format output.Format, out string, err error) string {
	t.Helper()

	if err != nil {
		t.Fatalf("render %s failed: %v", format, err)
	}

	return out
}

func renderTableFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, formatHealth(p.Health), formatComplexity(p.Complexity))
	}

	out, err := tbl.Render()

	return mustRender(t, output.FormatTable, out, err)
}

func renderJSONFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	data, err := serialization.MarshalJSON(projects)

	return mustRender(t, output.FormatJSON, string(data), err)
}

func renderMarkdownFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	headers := []string{"Name", "Health", "Complexity"}

	return renderMarkdownTable(t, headers, formatProjectsToRows(projects))
}

func renderCSVFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	var buf bytes.Buffer

	w := delimited.NewCSVWriter(&buf)

	if err := w.WriteHeader([]string{"Name", "Health", "Complexity"}); err != nil {
		t.Fatalf("write CSV header: %v", err)
	}

	for _, row := range formatProjectsToRows(projects) {
		if err := w.WriteRow(row); err != nil {
			t.Fatalf("write CSV row: %v", err)
		}
	}

	w.Flush()

	return buf.String()
}

func renderTSVFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	var buf bytes.Buffer

	w := delimited.NewTSVWriter(&buf)

	if err := w.WriteHeader([]string{"Name", "Health", "Complexity"}); err != nil {
		t.Fatalf("write TSV header: %v", err)
	}

	for _, row := range formatProjectsToRows(projects) {
		if err := w.WriteRow(row); err != nil {
			t.Fatalf("write TSV row: %v", err)
		}
	}

	w.Flush()

	return buf.String()
}

func newProjectTable(projects []TestProject) *output.Table {
	return &output.Table{
		Headers: []string{"Name", "Health", "Complexity"},
		Rows:    formatProjectsToRows(projects),
	}
}

func renderXMLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	data, err := markup.MarshalXMLFromTable(newProjectTable(projects))

	return mustRender(t, output.FormatXML, string(data), err)
}

func formatProjectsToRows(projects []TestProject) [][]string {
	rows := make([][]string, 0, len(projects))
	for _, p := range projects {
		rows = append(rows, []string{p.Name, formatHealth(p.Health), formatComplexity(p.Complexity)})
	}

	return rows
}

func renderYAMLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	data, err := serialization.MarshalYAML(projects)

	return mustRender(t, output.FormatYAML, string(data), err)
}

func renderHTMLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health", "Complexity"})

	for _, row := range formatProjectsToRows(projects) {
		html.AddRow(row)
	}

	out, err := html.Render()

	return mustRender(t, output.FormatHTML, out, err)
}

func buildProjectTree(projects []TestProject) *output.TreeNode {
	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		root.AddChild(output.NewTreeNode(p.Name, p.Name))
	}

	return root
}

func renderTreeFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	tr := tree.NewASCIITreeRenderer()
	tr.SetRoot(buildProjectTree(projects))

	out, err := tr.Render()

	return mustRender(t, output.FormatTree, out, err)
}

func renderD2Format(t *testing.T, projects []TestProject) string {
	t.Helper()

	d2Diagram := d2.NewDiagram()
	d2Diagram.AddTable("projects", []d2.Column{
		{Name: "name", Type: "string"},
	})

	for _, p := range projects {
		d2Diagram.AddNodeWithShape(p.Name, p.Name, d2.ShapeCircle)
	}

	out, err := d2Diagram.Render()

	return mustRender(t, output.FormatD2, out, err)
}

func renderNewD2FromTable(t *testing.T, projects []TestProject) string {
	t.Helper()

	data := newGraphTable(projects)

	out, err := d2.NewD2FromTable(data).Render()

	return mustRender(t, output.FormatD2, out, err)
}

func renderNewD2FromTree(t *testing.T, projects []TestProject) string {
	t.Helper()

	out, err := d2.NewD2FromTree(buildProjectTree(projects)).Render()

	return mustRender(t, output.FormatD2, out, err)
}

func newGraphTable(projects []TestProject) *output.Table {
	data := output.NewTable([]string{"Name"})
	for _, p := range projects {
		data.AddRow([]string{p.Name})
	}

	return data
}

func renderDOTFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	out, err := graph.NewDOTFromTable(newGraphTable(projects)).Render()

	return mustRender(t, output.FormatDOT, out, err)
}

func renderMermaidFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	out, err := graph.NewMermaidFromTable(newGraphTable(projects)).Render()

	return mustRender(t, output.FormatMermaid, out, err)
}

func renderJSONLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	data := newGraphTable(projects)
	b, err := serialization.MarshalJSONLFromTable(data)

	return mustRender(t, output.FormatJSONL, string(b), err)
}

func renderAsciiDocFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	b, err := markup.MarshalAsciiDocFromTable(newProjectTable(projects))

	return mustRender(t, output.FormatAsciiDoc, string(b), err)
}

func renderTOMLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	b, err := serialization.MarshalTOMLFromTable(newProjectTable(projects))

	return mustRender(t, output.FormatTOML, string(b), err)
}

func renderPlantUMLFormat(t *testing.T, projects []TestProject) string {
	t.Helper()

	out, err := plantuml.NewPlantUMLFromTable(newGraphTable(projects)).Render()

	return mustRender(t, output.FormatPlantUML, out, err)
}

func formatHealth(h int) string {
	return fmt.Sprintf("%d%%", h)
}

func formatComplexity(c int) string {
	return fmt.Sprintf("%d/10", c)
}
