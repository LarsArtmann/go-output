package serialization

import (
	"encoding/json"
	"fmt"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer           = (*JSONTreeRenderer)(nil)
	_ output.TreeOutputRenderer = (*JSONTreeRenderer)(nil)
	_ output.Renderer           = (*JSONGraphRenderer)(nil)
	_ output.GraphRenderer      = (*JSONGraphRenderer)(nil)
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

func (r *JSONTreeRenderer) toJSONNode(node *output.TreeNode) jsonTreeNode {
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
	output.GraphRendererMixin
}

// NewJSONGraphRenderer creates a new JSONGraphRenderer.
func NewJSONGraphRenderer() *JSONGraphRenderer {
	return &JSONGraphRenderer{
		GraphRendererMixin: output.NewGraphRendererMixin(),
	}
}

type jsonGraph struct {
	Nodes []jsonGraphNode `json:"nodes"`
	Edges []jsonGraphEdge `json:"edges"`
}

type jsonGraphNode struct {
	ID       string            `json:"id"`
	Label    string            `json:"label"`
	Shape    string            `json:"shape,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

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
			Label: brandedEdgeLabel(edge.Label),
		}

		graph.Edges = append(graph.Edges, e)
	}

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json graph: %w", err)
	}

	return string(data), nil
}

func brandedEdgeLabel(label output.GraphNodeLabel) string {
	if label.IsZero() {
		return ""
	}

	return label.Get()
}
