package output

import (
	"strings"

	"github.com/larsartmann/go-output/internal/escape"
)

// Compile-time interface checks.
var (
	_ Renderer      = (*HTMLRenderer)(nil)
	_ Renderer      = (*HTMLTreeRenderer)(nil)
	_ TableRenderer = (*HTMLRenderer)(nil)
)

// HTMLRenderer implements the Renderer interface for HTML table output.
type HTMLRenderer struct {
	tableDataBase
}

// tableDataBase provides common table data storage for renderers.
type tableDataBase struct {
	data *TableData
}

// ensureData initializes data if nil.
func (b *tableDataBase) ensureData() {
	if b.data == nil {
		b.data = &TableData{}
	}
}

// SetHeaders sets the column headers.
func (b *tableDataBase) SetHeaders(headers []string) {
	b.ensureData()
	b.data.Headers = headers
}

// AddRow adds a data row.
func (b *tableDataBase) AddRow(row []string) {
	b.ensureData()
	b.data.Rows = append(b.data.Rows, row)
}

// SetData sets the table data directly.
func (b *tableDataBase) SetData(data *TableData) {
	b.data = data
}

// NewHTMLRenderer creates a new HTMLRenderer.
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

// Render returns the HTML table as a string.
func (r *HTMLRenderer) Render() string {
	if r.data == nil {
		return "<table class=\"data-table\"></table>"
	}

	var b strings.Builder

	b.WriteString(`<table class="data-table">
<thead>
<tr>
`)

	for _, h := range r.data.Headers {
		b.WriteString("<th>")
		b.WriteString(escape.HTML(h))
		b.WriteString("</th>\n")
	}

	b.WriteString(`</tr>
</thead>
<tbody>
`)

	for _, row := range r.data.Rows {
		_ = writeMarkupRow(&b, row, "tr", "td", "", escape.HTML)
	}

	b.WriteString(`</tbody>
</table>
`)

	return b.String()
}

// RenderFullHTML returns a complete HTML document with the table.
func (r *HTMLRenderer) RenderFullHTML(title string) string {
	return renderFullHTMLDocument(title, `
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
`, r.Render())
}

// HTMLTreeRenderer renders a tree structure as HTML with nested lists.
type HTMLTreeRenderer struct {
	root *TreeNode
}

// NewHTMLTreeRenderer creates a new HTMLTreeRenderer.
func NewHTMLTreeRenderer() *HTMLTreeRenderer {
	return &HTMLTreeRenderer{ //nolint:exhaustruct // root is set via SetRoot
	}
}

// SetRoot sets the root node of the tree.
func (r *HTMLTreeRenderer) SetRoot(node *TreeNode) {
	r.root = node
}

// Render returns the tree as an HTML string.
func (r *HTMLTreeRenderer) Render() string {
	if r.root == nil {
		return "<ul class=\"tree\"></ul>"
	}

	var b strings.Builder
	b.WriteString(`<ul class="tree">
`)
	r.renderNode(&b, r.root)
	b.WriteString(`</ul>
`)

	return b.String()
}

func (r *HTMLTreeRenderer) renderNode(b *strings.Builder, node *TreeNode) {
	b.WriteString("<li>")
	b.WriteString(escape.HTML(node.Label.Get()))

	if len(node.Children) > 0 {
		b.WriteString("\n<ul>\n")

		for _, child := range node.Children {
			r.renderNode(b, child)
		}

		b.WriteString("</ul>\n")
	}

	b.WriteString("</li>\n")
}

// RenderFullHTML returns a complete HTML document with the tree.
func (r *HTMLTreeRenderer) RenderFullHTML(title string) string {
	return renderFullHTMLDocument(title, `
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
`, r.Render())
}

// renderFullHTMLDocument creates a complete HTML document with the given styles and content.
func renderFullHTMLDocument(title, styles, content string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + escape.HTML(title) + `</title>
<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  margin: 2rem;
  background: #f5f5f5;
}
` + styles + `
</style>
</head>
<body>
<h1>` + escape.HTML(title) + `</h1>
` + content + `
</body>
</html>`
}
