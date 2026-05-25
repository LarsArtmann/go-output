package serialization

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer           = (*TOMLTreeRenderer)(nil)
	_ output.TreeOutputRenderer = (*TOMLTreeRenderer)(nil)
)

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

type tomlTreeNode struct {
	ID       string            `toml:"id"`
	Label    string            `toml:"label"`
	Children []tomlTreeNode    `toml:"children,omitempty"`
	Metadata map[string]string `toml:"metadata,omitempty"`
}

// Render returns the tree as a TOML string.
func (r *TOMLTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "", nil
	}

	node := r.toTOMLNode(r.root)

	data, err := toml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("marshal toml tree: %w", err)
	}

	return string(data), nil
}

func (r *TOMLTreeRenderer) toTOMLNode(node *output.TreeNode) tomlTreeNode {
	result := tomlTreeNode{
		ID:       node.ID.Get(),
		Label:    node.Label.Get(),
		Metadata: node.Metadata,
	}

	if len(node.Children) > 0 {
		result.Children = make([]tomlTreeNode, 0, len(node.Children))
		for _, child := range node.Children {
			result.Children = append(result.Children, r.toTOMLNode(child))
		}
	}

	return result
}
