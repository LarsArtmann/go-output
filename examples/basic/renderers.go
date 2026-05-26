// Package main provides format-specific renderer functions for the basic example.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/examples/shared"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
)

func renderTable(projects []Project) {
	tbl := table.New(table.WithColorMode(colorMode))
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
	md := output.NewMarkdownTable().SetColorMode(colorMode)
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

func renderCSV(projects []Project) {
	w := delimited.NewCSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

func renderTSV(projects []Project) {
	w := delimited.NewTSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

type delimitedWriter interface {
	WriteHeader(cols []string) error
	WriteRow(values []string) error
	Flush()
	Error() error
}

func renderDelimited(w delimitedWriter, projects []Project) {
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

	xmlData, err := markup.MarshalXMLFromTableData(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(xmlData))
}

func renderYAML(projects []Project) {
	data, err := serialization.MarshalYAML(projects)
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
	html := markup.NewHTMLRenderer()
	html.SetHeaders(projectHeaders)

	for _, p := range projects {
		html.AddRow(projectToRow(p))
	}

	fmt.Println(html.RenderFullHTML("Project Health Report"))
}

func renderTree(projects []Project) {
	tree := output.NewASCIITreeRenderer()
	tree.SetColorMode(colorMode)

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

func renderJSONL(projects []Project) {
	data := projectsToTableData(projects)

	b, err := serialization.MarshalJSONLFromTableData(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Print(string(b))
}

func renderAsciiDoc(projects []Project) {
	data := projectsToTableData(projects)

	b, err := markup.MarshalAsciiDocFromTableData(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(b))
}

func renderTOML(projects []Project) {
	data := projectsToTableData(projects)

	b, err := serialization.MarshalTOMLFromTableData(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Print(string(b))
}

func renderPlantUML(projects []Project) {
	diagram := plantuml.NewPlantUMLDiagram()
	diagram.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("projects"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Projects"),
	})

	for _, p := range projects {
		diagram.AddNode(output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](p.Name),
		})
		diagram.AddEdge(output.GraphEdge{
			From:  output.NewBrandedID[output.GraphNodeIDBrand]("projects"),
			To:    output.NewBrandedID[output.GraphNodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("contains"),
		})
	}

	out, err := diagram.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}
