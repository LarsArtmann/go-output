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

// Render returns the graph as a JSON string.
func (r *JSONGraphRenderer) Render() (string, error) {
	graph := buildGraphDTO(r.GraphRendererMixin)

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
