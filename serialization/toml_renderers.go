package serialization

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*TOMLTreeRenderer)(nil)
	_ output.TreeRenderer  = (*TOMLTreeRenderer)(nil)
	_ output.Renderer      = (*TOMLGraphRenderer)(nil)
	_ output.GraphRenderer = (*TOMLGraphRenderer)(nil)
)

// TOMLGraphRenderer renders graph nodes and edges as TOML.
type TOMLGraphRenderer struct {
	output.GraphBuilder
}

// NewTOMLGraphRenderer creates a new TOMLGraphRenderer.
func NewTOMLGraphRenderer() *TOMLGraphRenderer {
	return &TOMLGraphRenderer{
		GraphBuilder: *output.NewGraphBuilder(),
	}
}

// Render returns the graph as a TOML string.
func (r *TOMLGraphRenderer) Render() (string, error) {
	graph := buildGraphView(r.GraphBuilder)

	data, err := toml.Marshal(graph)
	if err != nil {
		return "", fmt.Errorf("marshal toml graph: %w", err)
	}

	return string(data), nil
}

// TOMLTreeRenderer renders a TreeNode hierarchy as TOML.
type TOMLTreeRenderer struct {
	root *output.TreeNode
}

// NewTOMLTreeRenderer creates a new TOMLTreeRenderer.
func NewTOMLTreeRenderer() *TOMLTreeRenderer {
	return &TOMLTreeRenderer{}
}

// SetRoot sets the root node of the tree.
func (r *TOMLTreeRenderer) SetRoot(node *output.TreeNode) {
	r.root = node
}

// Render returns the tree as a TOML string.
func (r *TOMLTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "", nil
	}

	node := toTreeNode(r.root)

	data, err := toml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("marshal toml tree: %w", err)
	}

	return string(data), nil
}
