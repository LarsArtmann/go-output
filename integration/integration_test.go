// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

type TestProject struct {
	Name       string
	Health     int
	Complexity int
}

func TestAllFormatsRender(t *testing.T) {
	t.Parallel()

	projects := []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
	}

	formats := []output.Format{
		output.FormatTable,
		output.FormatJSON,
		output.FormatMarkdown,
		output.FormatCSV,
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

func TestFormatParseRoundtrip(t *testing.T) {
	t.Parallel()

	formats := output.FormatTable.AllowedValues()
	for _, formatStr := range formats {
		t.Run(formatStr, func(t *testing.T) {
			t.Parallel()
			format, err := output.ParseFormat(formatStr)
			if err != nil {
				t.Errorf("ParseFormat(%q) failed: %v", formatStr, err)
				return
			}
			if !format.IsValid() {
				t.Errorf("Format %q should be valid after parsing", formatStr)
			}
			if string(format) != formatStr {
				t.Errorf("Format.String() = %q, want %q", format, formatStr)
			}
		})
	}
}

func TestTableFormatContent(t *testing.T) {
	t.Parallel()

	projects := []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
	}

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

	md := output.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Health"})
	md.AddRow([]string{"Alpha", "90%"})

	result, err := md.Render()
	if err != nil {
		t.Fatalf("Markdown render failed: %v", err)
	}

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
	if err := w.WriteHeader([]string{"Name", "Health"}); err != nil {
		t.Fatalf("WriteHeader failed: %v", err)
	}
	if err := w.WriteRow([]string{"Alpha", "90"}); err != nil {
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

func TestMermaidFormatContent(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	mermaid := output.MermaidFlowchartRenderer(data)
	result := mermaid.Render()

	if result == "" {
		t.Error("Mermaid render should not be empty")
	}
}

func TestDOTFormatContent(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	dot := output.DOTFromTableData(data)
	result := dot.Render()

	if result == "" {
		t.Error("DOT render should not be empty")
	}
}

func TestStreamingRenderer(t *testing.T) {
	t.Parallel()

	html := output.NewStreamingHTMLRenderer()
	html.SetHeaders([]string{"Name"})
	html.AddRow([]string{"Alpha"})

	var buf bytes.Buffer
	if err := html.Stream(&buf); err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "<table") {
		t.Error("Streaming HTML should contain table tag")
	}
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

	if root.Depth() != 0 {
		t.Errorf("Root depth should be 0, got %d", root.Depth())
	}
	if child.Depth() != 1 {
		t.Errorf("Child depth should be 1, got %d", child.Depth())
	}
	if grandchild.Depth() != 2 {
		t.Errorf("Grandchild depth should be 2, got %d", grandchild.Depth())
	}
}

func TestInvalidFormatError(t *testing.T) {
	t.Parallel()

	err := &output.InvalidFormatError{
		Value:   "invalid",
		Allowed: []output.Format{output.FormatTable},
	}
	result := err.Error()

	if !strings.Contains(result, "invalid format") {
		t.Error("Error message should contain 'invalid format'")
	}
	if !strings.Contains(result, "invalid") {
		t.Error("Error message should contain the invalid value")
	}
}

func TestFormatCategories(t *testing.T) {
	t.Parallel()

	tableFormats := []output.Format{
		output.FormatTable,
		output.FormatJSON,
		output.FormatCSV,
		output.FormatMarkdown,
		output.FormatYAML,
	}

	for _, f := range tableFormats {
		if !f.IsTableFormat() {
			t.Errorf("Format %s should be a table format", f)
		}
	}

	treeFormats := []output.Format{
		output.FormatTree,
		output.FormatHTML,
	}

	for _, f := range treeFormats {
		if !f.IsTreeFormat() {
			t.Errorf("Format %s should be a tree format", f)
		}
	}

	graphFormats := []output.Format{
		output.FormatD2,
		output.FormatMermaid,
		output.FormatDOT,
	}

	for _, f := range graphFormats {
		if !f.IsGraphFormat() {
			t.Errorf("Format %s should be a graph format", f)
		}
	}
}

func TestFormatRegistry(t *testing.T) {
	t.Parallel()

	customFormat := output.Format("custom")
	if err := output.Register(customFormat, func() output.Renderer {
		return nil
	}); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer output.Unregister(customFormat)

	if !output.IsRegistered(customFormat) {
		t.Error("Custom format should be registered")
	}

	formats := output.RegisteredFormats()
	found := false
	for _, f := range formats {
		if f == customFormat {
			found = true
			break
		}
	}
	if !found {
		t.Error("Custom format should be in registered formats list")
	}
}

// renderProject renders projects in the specified format.
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
	case output.FormatYAML:
		return renderYAMLFormat(projects)
	case output.FormatHTML:
		return renderHTMLFormat(projects)
	case output.FormatTree:
		return renderTreeFormat(projects)
	case output.FormatD2:
		return renderD2Format()
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
	return tbl.Render()
}

func renderJSONFormat(projects []TestProject) string {
	data, _ := output.MarshalJSONIndent(projects, "", "  ")
	return string(data)
}

func renderMarkdownFormat(projects []TestProject) string {
	md := output.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		md.AddRow([]string{p.Name, formatHealth(p.Health), formatComplexity(p.Complexity)})
	}
	result, _ := md.Render()
	return result
}

func renderCSVFormat(projects []TestProject) string {
	var buf bytes.Buffer
	w := output.NewCSVWriter(&buf)
	_ = w.WriteHeader([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		_ = w.WriteRow([]string{p.Name, formatHealth(p.Health), formatComplexity(p.Complexity)})
	}
	w.Flush()
	return buf.String()
}

func renderYAMLFormat(projects []TestProject) string {
	data, _ := output.MarshalYAML(projects)
	return string(data)
}

func renderHTMLFormat(projects []TestProject) string {
	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		html.AddRow([]string{p.Name, formatHealth(p.Health), formatComplexity(p.Complexity)})
	}
	return html.Render()
}

func renderTreeFormat(projects []TestProject) string {
	tree := output.NewASCIITreeRenderer()
	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		root.AddChild(output.NewTreeNode(p.Name, p.Name))
	}
	tree.SetRoot(root)
	return tree.Render()
}

func renderD2Format() string {
	d2 := output.NewD2Diagram()
	d2.AddTable("projects", []output.D2Column{
		{Name: "name", Type: "string"},
	})
	return d2.Render()
}

func renderMermaidFormat(projects []TestProject) string {
	data := output.NewTableData([]string{"Name"})
	for _, p := range projects {
		data.AddRow([]string{p.Name})
	}
	return output.MermaidFlowchartRenderer(data).Render()
}

func renderDOTFormat(projects []TestProject) string {
	data := output.NewTableData([]string{"Name"})
	for _, p := range projects {
		data.AddRow([]string{p.Name})
	}
	return output.DOTFromTableData(data).Render()
}

func formatHealth(h int) string {
	return fmt.Sprintf("%d%%", h)
}

func formatComplexity(c int) string {
	return fmt.Sprintf("%d/10", c)
}
