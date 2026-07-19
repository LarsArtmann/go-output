package serialization

import (
	"fmt"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*YAMLTreeRenderer)(nil)
	_ output.TreeRenderer  = (*YAMLTreeRenderer)(nil)
	_ output.Renderer      = (*YAMLGraphRenderer)(nil)
	_ output.GraphRenderer = (*YAMLGraphRenderer)(nil)
)

// YAMLTreeRenderer renders a TreeNode hierarchy as YAML.
type YAMLTreeRenderer struct {
	root *output.TreeNode
}

// NewYAMLTreeRenderer creates a new YAMLTreeRenderer.
func NewYAMLTreeRenderer() *YAMLTreeRenderer {
	return &YAMLTreeRenderer{}
}

// SetRoot sets the root node of the tree.
func (r *YAMLTreeRenderer) SetRoot(node *output.TreeNode) {
	r.root = node
}

// Render returns the tree as a YAML string.
func (r *YAMLTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "null\n", nil
	}

	node := toTreeNode(r.root)

	data, err := yaml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("marshal yaml tree: %w", err)
	}

	return string(data), nil
}

// YAMLGraphRenderer renders graph nodes and edges as YAML.
type YAMLGraphRenderer struct {
	output.GraphBuilder
}

// NewYAMLGraphRenderer creates a new YAMLGraphRenderer.
func NewYAMLGraphRenderer() *YAMLGraphRenderer {
	return &YAMLGraphRenderer{
		GraphBuilder: *output.NewGraphBuilder(),
	}
}

// Render returns the graph as a YAML string.
func (r *YAMLGraphRenderer) Render() (string, error) {
	graph := buildGraphView(r.GraphBuilder)

	data, err := yaml.Marshal(graph)
	if err != nil {
		return "", fmt.Errorf("marshal yaml graph: %w", err)
	}

	return string(data), nil
}
