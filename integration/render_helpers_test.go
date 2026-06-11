package integration

import (
	"bytes"
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
)

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
	case output.FormatJSONL:
		return renderJSONLFormat(projects)
	case output.FormatAsciiDoc:
		return renderAsciiDocFormat(projects)
	case output.FormatTOML:
		return renderTOMLFormat(projects)
	case output.FormatPlantUML:
		return renderPlantUMLFormat(projects)
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
	data, _ := serialization.MarshalJSON(projects)

	return string(data)
}

func renderMarkdownFormat(projects []TestProject) string {
	headers := []string{"Name", "Health", "Complexity"}

	return renderMarkdownTable(headers, formatProjectsToRows(projects))
}

func renderCSVFormat(projects []TestProject) string {
	var buf bytes.Buffer

	w := delimited.NewCSVWriter(&buf)

	_ = w.WriteHeader([]string{"Name", "Health", "Complexity"})
	for _, row := range formatProjectsToRows(projects) {
		_ = w.WriteRow(row)
	}

	w.Flush()

	return buf.String()
}

func renderTSVFormat(projects []TestProject) string {
	var buf bytes.Buffer

	w := delimited.NewTSVWriter(&buf)

	_ = w.WriteHeader([]string{"Name", "Health", "Complexity"})
	for _, row := range formatProjectsToRows(projects) {
		_ = w.WriteRow(row)
	}

	w.Flush()

	return buf.String()
}

func newProjectTableData(projects []TestProject) *output.TableData {
	return &output.TableData{
		Headers: []string{"Name", "Health", "Complexity"},
		Rows:    formatProjectsToRows(projects),
	}
}

func renderXMLFormat(projects []TestProject) string {
	data, _ := markup.MarshalXMLFromTableData(newProjectTableData(projects))

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
	data, _ := serialization.MarshalYAML(projects)

	return string(data)
}

func renderHTMLFormat(projects []TestProject) string {
	html := markup.NewHTMLRenderer()
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
	d2Diagram := d2.NewD2Diagram()
	d2Diagram.AddTable("projects", []d2.D2Column{
		{Name: "name", Type: "string"},
	})

	for _, p := range projects {
		d2Diagram.AddNodeWithShape(p.Name, p.Name, d2.D2ShapeCircle)
	}

	out, err := d2Diagram.Render()
	if err != nil {
		return ""
	}

	return out
}

func renderD2FromTableData(projects []TestProject) string {
	data := newGraphTableData(projects)

	out, err := d2.D2FromTableData(data).Render()
	if err != nil {
		return ""
	}

	return out
}

func renderD2FromTree(projects []TestProject) string {
	out, err := d2.D2FromTree(buildProjectTree(projects)).Render()
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
	out, err := graph.DOTFromTableData(newGraphTableData(projects)).Render()
	if err != nil {
		return ""
	}

	return out
}

func renderMermaidFormat(projects []TestProject) string {
	out, err := graph.MermaidFromTableData(newGraphTableData(projects)).Render()
	if err != nil {
		return ""
	}

	return out
}

func renderJSONLFormat(projects []TestProject) string {
	data := newGraphTableData(projects)
	b, _ := serialization.MarshalJSONLFromTableData(data)

	return string(b)
}

func renderAsciiDocFormat(projects []TestProject) string {
	b, _ := markup.MarshalAsciiDocFromTableData(newProjectTableData(projects))

	return string(b)
}

func renderTOMLFormat(projects []TestProject) string {
	b, _ := serialization.MarshalTOMLFromTableData(newProjectTableData(projects))

	return string(b)
}

func renderPlantUMLFormat(projects []TestProject) string {
	out, err := plantuml.PlantUMLFromTableData(newGraphTableData(projects)).Render()
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
