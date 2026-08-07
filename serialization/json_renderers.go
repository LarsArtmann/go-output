package serialization

import (
	"encoding/json/jsontext"
	"encoding/json/v2"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*JSONTreeRenderer)(nil)
	_ output.TreeRenderer  = (*JSONTreeRenderer)(nil)
	_ output.Renderer      = (*JSONGraphRenderer)(nil)
	_ output.GraphRenderer = (*JSONGraphRenderer)(nil)
)

// JSONTreeRenderer renders a TreeNode hierarchy as JSON.
type JSONTreeRenderer struct {
	root *output.TreeNode
}

// NewJSONTreeRenderer creates a new JSONTreeRenderer.
func NewJSONTreeRenderer() *JSONTreeRenderer {
	return &JSONTreeRenderer{}
}

// SetRoot sets the root node of the tree.
func (r *JSONTreeRenderer) SetRoot(node *output.TreeNode) {
	r.root = node
}

// Render returns the tree as a JSON string.
func (r *JSONTreeRenderer) Render() (string, error) {
	return renderTreeNode(r.root, "null", "json", func(v any) ([]byte, error) {
		return json.Marshal(v, json.Deterministic(true), jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))
	})
}

// JSONGraphRenderer renders graph nodes and edges as JSON.
type JSONGraphRenderer struct {
	output.GraphBuilder
}

// NewJSONGraphRenderer creates a new JSONGraphRenderer.
func NewJSONGraphRenderer() *JSONGraphRenderer {
	return &JSONGraphRenderer{
		GraphBuilder: *output.NewGraphBuilder(),
	}
}

// Render returns the graph as a JSON string.
func (r *JSONGraphRenderer) Render() (string, error) {
	graph := buildGraphView(r.GraphBuilder)

	data, err := json.Marshal(graph, json.Deterministic(true), jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))

	return stringFromBytes("json", "graph", data, err)
}

func brandedEdgeLabel(label output.GraphNodeLabel) string {
	if label.IsZero() {
		return ""
	}

	return label.Get()
}
