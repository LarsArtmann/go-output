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
	"github.com/larsartmann/go-output/markdown"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
	"github.com/larsartmann/go-output/tree"
)

// alphaToBetaEdge is a sample graph edge shared by the Mermaid and DOT
// renderers, demonstrating the same edge construction across formats.
var alphaToBetaEdge = output.GraphEdge{
	From: output.NewBrandedID[output.GraphNodeIDBrand]("Alpha"),
	To:   output.NewBrandedID[output.GraphNodeIDBrand]("Beta"),
}

func renderTable(projects []Project) {
	tbl := table.New(table.WithColorMode(colorMode))
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, strconv.Itoa(p.Health)+"%", strconv.Itoa(p.Complexity)+"/10")
	}

	tbl.SetFooter("TOTAL", strconv.Itoa(len(projects)), "-")

	shared.RenderAndPrint(tbl)
}

func renderJSON(projects []Project) {
	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(data))
}

func renderMarkdown(projects []Project) {
	md := markdown.NewMarkdownTable().SetColorMode(colorMode)
	md.SetHeaders([]string{"Name", "Health", "Complexity"})

	for _, p := range projects {
		md.AddRow(
			[]string{p.Name, fmt.Sprintf("%d%%", p.Health), fmt.Sprintf("%d/10", p.Complexity)},
		)
	}

	shared.RenderAndPrint(md)
}

func renderCSV(projects []Project) {
	w := delimited.NewCSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

func renderTSV(projects []Project) {
	w := delimited.NewTSVWriter(os.Stdout)
	renderDelimited(w, projects)
}

func renderDelimited(w delimited.Writer, projects []Project) {
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

	footer := []string{"TOTAL", strconv.Itoa(len(projects)), "-"}

	if err := w.WriteFooter(footer); err != nil {
		shared.HandleError(err)
	}

	w.Flush()

	err = w.Error()
	if err != nil {
		shared.HandleError(err)
	}
}

func renderXML(projects []Project) {
	data := projectsToTable(projects)

	xmlData, err := markup.MarshalXMLFromTable(data)
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
		AddTable("projects", []d2.Column{
			{Name: "id", Type: "serial", Constraint: d2.ConstraintPrimary},
			{Name: "name", Type: "varchar(255)"},
			{Name: "health", Type: "int"},
			{Name: "complexity", Type: "int"},
		})

	for _, p := range projects {
		d2Diagram.AddNode(d2.Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.D2NodeLabelBrand](p.Name),
			Shape: d2.ShapeCircle,
			Class: "service",
		})
		d2Diagram.AddEdge(d2.Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("projects"),
			To:          output.NewBrandedID[output.D2NodeIDBrand](p.Name),
			TargetArrow: d2.ArrowCFMany,
		})
	}

	shared.RenderAndPrint(d2Diagram)
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
	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		projNode := output.NewTreeNode("proj-"+p.Name, p.Name)
		projNode.Metadata["health"] = strconv.Itoa(p.Health) + "%"
		projNode.Metadata["complexity"] = strconv.Itoa(p.Complexity)
		root.AddChild(projNode)
	}

	// CQRS: stream the rendered tree straight to stdout.
	if err := tree.WriteASCII(os.Stdout, root, tree.WithColorMode(colorMode)); err != nil {
		shared.HandleError(err)
	}
}

func renderMermaid(projects []Project) {
	builder := output.NewGraphBuilder()

	for _, p := range projects {
		builder.AddNode(output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](p.Name),
			Style: output.NodeStyle{
				Fill:      "#e8a838",
				Stroke:    "#4a4030",
				FontColor: "#14110d",
			},
		})
	}

	// The same edge added twice demonstrates the DedupEdges() feature.
	for range 2 {
		builder.AddEdge(alphaToBetaEdge)
	}

	// Removes the duplicate Alpha -> Beta edge before freezing the graph.
	builder.DedupEdges()

	g := builder.Build()

	// CQRS render without a surrounding ```mermaid code fence.
	out, err := graph.RenderMermaid(g, graph.WithCodeFence(false))
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Print(out)
}

func renderDOT(projects []Project) {
	builder := output.NewGraphBuilder()

	for _, p := range projects {
		builder.AddNode(output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand](p.Name),
		})
	}

	builder.AddEdge(alphaToBetaEdge)

	// CQRS render with layout options — the functional equivalents of the
	// legacy renderer's SetRankDir/SetSplines/SetNodeSep/SetRankSep chain.
	err := graph.WriteDOT(
		os.Stdout,
		builder.Build(),
		graph.WithDOTRankDir(graph.RankDirLR),
		graph.WithDOTSplines(graph.SplineSpline),
		graph.WithNodeSep("0.8"),
		graph.WithRankSep("1.0"),
	)
	if err != nil {
		shared.HandleError(err)
	}
}

func renderJSONL(projects []Project) {
	data := projectsToTable(projects)

	b, err := serialization.MarshalJSONLFromTable(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Print(string(b))
}

func renderAsciiDoc(projects []Project) {
	data := projectsToTable(projects)

	b, err := markup.MarshalAsciiDocFromTable(data)
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(string(b))
}

func renderTOML(projects []Project) {
	data := projectsToTable(projects)

	b, err := serialization.MarshalTOMLFromTable(data)
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
			Style: output.NodeStyle{
				Fill:   "#e8a838",
				Stroke: "#4a4030",
			},
		})
		diagram.AddEdge(output.GraphEdge{
			From:  output.NewBrandedID[output.GraphNodeIDBrand]("projects"),
			To:    output.NewBrandedID[output.GraphNodeIDBrand](p.Name),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("contains"),
		})
	}

	shared.RenderAndPrint(diagram)
}
