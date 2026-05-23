// Package main demonstrates usage of the go-output library.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/examples/shared"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/table"
)

// Project represents a sample data structure for demonstration.
type Project struct {
	Name       string
	Health     int
	Complexity int
}

// rendererFunc is a function that renders projects in a specific format.
type rendererFunc func([]Project)

// getRenderers returns a map of format to renderer function.
func getRenderers() map[output.Format]rendererFunc {
	return map[output.Format]rendererFunc{
		output.FormatTable:    renderTable,
		output.FormatJSON:     renderJSON,
		output.FormatMarkdown: renderMarkdown,
		output.FormatCSV:      renderCSV,
		output.FormatTSV:      renderTSV,
		output.FormatXML:      renderXML,
		output.FormatYAML:     renderYAML,
		output.FormatD2:       renderD2,
		output.FormatHTML:     renderHTML,
		output.FormatTree:     renderTree,
		output.FormatMermaid:  renderMermaid,
		output.FormatDOT:      renderDOT,
	}
}

func main() {
	projects := []Project{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
		{Name: "Gamma", Health: 85, Complexity: 8},
	}

	// Parse command line format (default to table)
	format := output.FormatTable

	if len(os.Args) > 1 {
		f, err := output.ParseFormat(os.Args[1])
		if err != nil {
			shared.HandleError(err)
		}

		format = f
	}

	// Output in the specified format
	renderOutput(format, projects)
}

func renderOutput(format output.Format, projects []Project) {
	renderers := getRenderers()
	if renderer, ok := renderers[format]; ok {
		renderer(projects)

		return
	}
	// Handle unknown format safely - format is validated by ParseFormat
	//nolint:gosec // format is validated enum type
	fmt.Fprintf(os.Stderr, "Unsupported format: %s\n", format)
	fmt.Fprintf(os.Stderr, "Available formats: %v\n", output.FormatTable.AllowedValues())
	os.Exit(1)
}

func renderTable(projects []Project) {
	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, strconv.Itoa(p.Health)+"%", strconv.Itoa(p.Complexity)+"/10")
	}

	out, err := tbl.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

func renderJSON(projects []Project) {
	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(data))
}

func renderMarkdown(projects []Project) {
	md := output.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Health", "Complexity"})

	for _, p := range projects {
		md.AddRow(
			[]string{p.Name, fmt.Sprintf("%d%%", p.Health), fmt.Sprintf("%d/10", p.Complexity)},
		)
	}

	out, err := md.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

// projectHeaders defines the common headers for project data.
var projectHeaders = []string{"Name", "Health", "Complexity"}

// projectToRow converts a Project to a formatted row slice.
func projectToRow(p Project) []string {
	return []string{
		p.Name,
		strconv.Itoa(p.Health) + "%",
		strconv.Itoa(p.Complexity) + "/10",
	}
}

// projectsToTableData creates TableData from projects.
func projectsToTableData(projects []Project) *output.TableData {
	data := output.NewTableData(projectHeaders)
	for _, p := range projects {
		data.AddRow(projectToRow(p))
	}

	return data
}

func renderCSV(projects []Project) {
	w := output.NewCSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

func renderTSV(projects []Project) {
	w := output.NewTSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

// writer interface matches both CSVWriter and TSVWriter.
type writer interface {
	WriteHeader(cols []string) error
	WriteRow(values []string) error
	Flush()
	Error() error
}

func renderDelimited(w writer, projects []Project) {
	err := w.WriteHeader(projectHeaders)
	if err != nil {
		shared.HandleError(err)
	}

	for _, p := range projects {
		err := w.WriteRow(projectToRow(p))
		if err != nil {
			shared.HandleError(err)
		}
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		shared.HandleError(err)
	}
}

func renderXML(projects []Project) {
	data := projectsToTableData(projects)

	xmlData, err := output.MarshalXMLFromTableData(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(xmlData))
}

func renderYAML(projects []Project) {
	data, err := output.MarshalYAML(projects)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(data))
}

func renderD2(projects []Project) {
	d2Diagram := shared.NewServiceD2Diagram("Project Architecture").
		AddTable("projects", []d2.D2Column{
			{Name: "id", Type: "serial", Constraint: d2.D2ConstraintPrimary},
			{Name: "name", Type: "varchar(255)"},
			{Name: "health", Type: "int"},
			{Name: "complexity", Type: "int"},
		})

	for _, p := range projects {
		d2Diagram.AddNode(d2.D2Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.D2NodeLabelBrand](p.Name),
			Shape: d2.D2ShapeCircle,
			Class: "service",
		})
		d2Diagram.AddEdge(d2.D2Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("projects"),
			To:          output.NewBrandedID[output.D2NodeIDBrand](p.Name),
			TargetArrow: d2.D2ArrowCFMany,
		})
	}

	out, err := d2Diagram.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

func renderHTML(projects []Project) {
	html := output.NewHTMLRenderer()
	html.SetHeaders(projectHeaders)

	for _, p := range projects {
		html.AddRow(projectToRow(p))
	}

	fmt.Println(html.RenderFullHTML("Project Health Report"))
}

func renderTree(projects []Project) {
	tree := output.NewASCIITreeRenderer()

	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		projNode := output.NewTreeNode("proj-"+p.Name, p.Name)
		projNode.Metadata["health"] = strconv.Itoa(p.Health) + "%"
		projNode.Metadata["complexity"] = strconv.Itoa(p.Complexity)
		root.AddChild(projNode)
	}

	tree.SetRoot(root)

	out, err := tree.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

func renderDiagram(projects []Project, createRenderer func(*output.TableData) output.Renderer) {
	data := projectsToTableData(projects)
	renderer := createRenderer(data)

	out, err := renderer.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

func renderMermaid(projects []Project) {
	renderDiagram(projects, func(data *output.TableData) output.Renderer {
		return graph.MermaidFromTableData(data)
	})
}

func renderDOT(projects []Project) {
	renderDiagram(projects, func(data *output.TableData) output.Renderer {
		return graph.DOTFromTableData(data)
	})
}
