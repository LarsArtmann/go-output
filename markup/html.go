package markup

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*HTMLRenderer)(nil)
	_ output.Renderer      = (*HTMLTreeRenderer)(nil)
	_ output.TableRenderer = (*HTMLRenderer)(nil)
)

//nolint:gochecknoinits // Registers HTML TableData marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatHTML, output.ShapeTable, output.ShapeTree)
	output.RegisterTableDataRenderer(output.FormatHTML, renderHTMLTableData)
}

// HTMLRenderer implements the Renderer interface for HTML table output.
type HTMLRenderer struct {
	output.TableDataStore
}

// NewHTMLRenderer creates a new HTMLRenderer.
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

// Render returns the HTML table as a string by delegating to StreamingHTMLRenderer.
func (r *HTMLRenderer) Render() (string, error) {
	var b strings.Builder

	err := streamHTMLTable(&b, r.Data())
	if err != nil {
		return "", fmt.Errorf("render html table: %w", err)
	}

	return b.String(), nil
}

// RenderFullHTML returns a complete HTML document with the table.
func (r *HTMLRenderer) RenderFullHTML(title string) (string, error) {
	return r.renderFullHTMLWithStyles(title, tableStyles)
}

func renderHTMLWithStyles(r output.Renderer, title, styles, errContext string) (string, error) {
	content, err := r.Render()
	if err != nil {
		return "", fmt.Errorf("%s for title=%q: %w", errContext, title, err)
	}

	return renderFullHTMLDocument(title, styles, content)
}

func (r *HTMLRenderer) renderFullHTMLWithStyles(title, styles string) (string, error) {
	return renderHTMLWithStyles(r, title, styles, "render html table")
}

// HTMLTreeRenderer renders a tree structure as HTML with nested lists.
type HTMLTreeRenderer struct {
	root *output.TreeNode
}

// NewHTMLTreeRenderer creates a new HTMLTreeRenderer.
func NewHTMLTreeRenderer() *HTMLTreeRenderer {
	return &HTMLTreeRenderer{}
}

// SetRoot sets the root node of the tree.
func (r *HTMLTreeRenderer) SetRoot(node *output.TreeNode) {
	r.root = node
}

// Render returns the tree as an HTML string.
func (r *HTMLTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return `<ul class="tree"></ul>`, nil
	}

	var b strings.Builder

	err := treeTemplate.Execute(&b, r.root)
	if err != nil {
		return "", fmt.Errorf("render html tree: %w", err)
	}

	return b.String(), nil
}

// RenderFullHTML returns a complete HTML document with the tree.
func (r *HTMLTreeRenderer) RenderFullHTML(title string) (string, error) {
	return r.renderFullHTMLWithStyles(title, treeStyles)
}

func (r *HTMLTreeRenderer) renderFullHTMLWithStyles(title, styles string) (string, error) {
	return renderHTMLWithStyles(r, title, styles, "render html tree")
}

const tableStyles = `
.data-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  border-radius: 8px;
  overflow: hidden;
}
.data-table th, .data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}
.data-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
}
.data-table tr:hover {
  background: #f8f9fa;
}
.data-table tr:last-child td {
  border-bottom: none;
}
`

const treeStyles = `
.tree {
  list-style: none;
  padding-left: 1.5rem;
}
.tree li {
  margin: 0.5rem 0;
  padding: 0.5rem;
  background: white;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}
.tree ul {
  list-style: none;
  padding-left: 1rem;
  margin-top: 0.5rem;
}
`

//nolint:gochecknoglobals // Parsed template for full HTML document.
var fullHTMLTemplate = template.Must(template.New("fullHTML").Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  margin: 2rem;
  background: #f5f5f5;
}
{{.Styles}}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
{{.Content}}
</body>
</html>`))

type fullHTMLData struct {
	Title   string
	Styles  string
	Content template.HTML
}

func renderFullHTMLDocument(title, styles, content string) (string, error) {
	var b strings.Builder

	// #nosec G203 — Content is already-rendered HTML from the table/tree renderer.
	// The template auto-escapes Title and Styles; Content is trusted HTML output.
	err := fullHTMLTemplate.Execute(&b, fullHTMLData{
		Title:   title,
		Styles:  styles,
		Content: template.HTML(content),
	})
	if err != nil {
		return "", fmt.Errorf("render full html document: %w", err)
	}

	return b.String(), nil
}

//nolint:gochecknoglobals // Parsed template for HTML tree rendering.
var treeTemplate = template.Must(template.New("treeNode").Parse(
	`<ul class="tree">
{{template "treeNodeRec" .}}
</ul>
` + treeNodeRecTemplate,
))

const treeNodeRecTemplate = `{{define "treeNodeRec"}}<li>{{.Label}}
{{- if .Children}}
<ul>
{{- range .Children}}
{{template "treeNodeRec" .}}
{{- end}}
</ul>
{{- end}}
</li>
{{end}}`

func renderHTMLTableData(w io.Writer, data *output.TableData, opts output.RenderOptions) error {
	renderer := NewHTMLRenderer()
	renderer.SetData(data)

	title := opts.Title
	if title == "" {
		title = fmt.Sprintf("Data - %d rows", data.RowCount())
	}

	out, err := renderer.RenderFullHTML(title)
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write html output: %w", err)
	}

	return nil
}
