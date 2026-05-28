package markup

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*HTMLRenderer)(nil)
	_ output.Renderer      = (*HTMLTreeRenderer)(nil)
	_ output.TableRenderer = (*HTMLRenderer)(nil)
)

//nolint:gochecknoinits // Registers HTML TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatHTML, renderHTMLTableData)
}

// HTMLRenderer implements the Renderer interface for HTML table output.
type HTMLRenderer struct {
	output.TableDataBase
}

// NewHTMLRenderer creates a new HTMLRenderer.
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{}
}

// Render returns the HTML table as a string.
func (r *HTMLRenderer) Render() (string, error) {
	data := r.Data()
	if data == nil {
		return "<table class=\"data-table\"></table>", nil
	}

	var b strings.Builder

	b.WriteString(`<table class="data-table">
<thead>
<tr>
`)

	for _, h := range data.Headers {
		b.WriteString("<th>")
		b.WriteString(escape.HTML(h))
		b.WriteString("</th>\n")
	}

	b.WriteString(`</tr>
</thead>
<tbody>
`)

	for _, row := range data.Rows {
		_ = writeMarkupRow(&b, row, "tr", "td", "", escape.HTML)
	}

	b.WriteString("</tbody>\n")

	if data.HasFooter() {
		b.WriteString("<tfoot>\n<tr>\n")

		for _, cell := range data.Footer {
			b.WriteString(`<td class="footer-cell">`)
			b.WriteString(escape.HTML(cell))
			b.WriteString("</td>\n")
		}

		b.WriteString("</tr>\n</tfoot>\n")
	}

	b.WriteString("</table>\n")

	return b.String(), nil
}

// RenderFullHTML returns a complete HTML document with the table.
func (r *HTMLRenderer) RenderFullHTML(title string) (string, error) {
	return r.renderFullHTMLWithStyles(title, `
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
`)
}

func renderHTMLWithStyles(r output.Renderer, title, styles, errContext string) (string, error) {
	content, err := r.Render()
	if err != nil {
		return "", fmt.Errorf("%s for title=%q styles=%q: %w", errContext, title, styles, err)
	}

	return renderFullHTMLDocument(title, styles, content), nil
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
		return "<ul class=\"tree\"></ul>", nil
	}

	var b strings.Builder
	b.WriteString(`<ul class="tree">
`)
	r.renderNode(&b, r.root)
	b.WriteString(`</ul>
`)

	return b.String(), nil
}

func (r *HTMLTreeRenderer) renderNode(b *strings.Builder, node *output.TreeNode) {
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
func (r *HTMLTreeRenderer) RenderFullHTML(title string) (string, error) {
	return r.renderFullHTMLWithStyles(title, `
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
`)
}

func (r *HTMLTreeRenderer) renderFullHTMLWithStyles(title, styles string) (string, error) {
	return renderHTMLWithStyles(r, title, styles, "render html tree")
}

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
