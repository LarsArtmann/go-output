package output

import (
	"encoding/json"
	"fmt"
)

// Compile-time interface checks.
var (
	_ Renderer           = (*JSONTreeRenderer)(nil)
	_ TreeOutputRenderer = (*JSONTreeRenderer)(nil)
	_ Renderer           = (*JSONGraphRenderer)(nil)
	_ GraphRenderer      = (*JSONGraphRenderer)(nil)
)

// JSONTreeRenderer renders a TreeNode hierarchy as JSON.
// Each node becomes a JSON object with "id", "label", and optional "children".
type JSONTreeRenderer struct {
	root *TreeNode
}

// NewJSONTreeRenderer creates a new JSONTreeRenderer.
func NewJSONTreeRenderer() *JSONTreeRenderer {
	return &JSONTreeRenderer{} //nolint:exhaustruct // root is set via SetRoot
}

// SetRoot sets the root node of the tree.
func (r *JSONTreeRenderer) SetRoot(node *TreeNode) {
	r.root = node
}

// jsonTreeNode is the JSON representation of a TreeNode.
type jsonTreeNode struct {
	ID       string            `json:"id"`
	Label    string            `json:"label"`
	Children []jsonTreeNode    `json:"children,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Render returns the tree as a JSON string.
func (r *JSONTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "null", nil
	}

	node := r.toJSONNode(r.root)

	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json tree: %w", err)
	}

	return string(data), nil
}

func (r *JSONTreeRenderer) toJSONNode(node *TreeNode) jsonTreeNode {
	result := jsonTreeNode{
		ID:       node.ID.Get(),
		Label:    node.Label.Get(),
		Metadata: node.Metadata,
	}

	if len(node.Children) > 0 {
		result.Children = make([]jsonTreeNode, 0, len(node.Children))
		for _, child := range node.Children {
			result.Children = append(result.Children, r.toJSONNode(child))
		}
	}

	return result
}

// JSONGraphRenderer renders graph nodes and edges as JSON.
type JSONGraphRenderer struct {
	GraphRendererMixin
}

// NewJSONGraphRenderer creates a new JSONGraphRenderer.
func NewJSONGraphRenderer() *JSONGraphRenderer {
	return &JSONGraphRenderer{
		GraphRendererMixin: NewGraphRendererMixin(),
	}
}

// jsonGraph is the JSON representation of a graph.
type jsonGraph struct {
	Nodes []jsonGraphNode `json:"nodes"`
	Edges []jsonGraphEdge `json:"edges"`
}

// jsonGraphNode is the JSON representation of a GraphNode.
type jsonGraphNode struct {
	ID       string            `json:"id"`
	Label    string            `json:"label"`
	Shape    string            `json:"shape,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// jsonGraphEdge is the JSON representation of a GraphEdge.
type jsonGraphEdge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label,omitempty"`
}

// Render returns the graph as a JSON string.
func (r *JSONGraphRenderer) Render() (string, error) {
	graph := jsonGraph{
		Nodes: make([]jsonGraphNode, 0, len(r.Nodes())),
		Edges: make([]jsonGraphEdge, 0, len(r.Edges())),
	}

	for _, node := range r.Nodes() {
		n := jsonGraphNode{
			ID:       node.ID.Get(),
			Label:    node.Label.Get(),
			Shape:    string(node.Shape),
			Metadata: node.Metadata,
		}
		graph.Nodes = append(graph.Nodes, n)
	}

	for _, edge := range r.Edges() {
		e := jsonGraphEdge{
			From:  edge.From.Get(),
			To:    edge.To.Get(),
			Label: brandedValue(edge.Label),
		}

		graph.Edges = append(graph.Edges, e)
	}

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json graph: %w", err)
	}

	return string(data), nil
}
