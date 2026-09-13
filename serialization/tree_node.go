package serialization

import (
	"github.com/larsartmann/go-output"
)

type treeNode struct {
	ID       string            `json:"id"                 yaml:"id"                 toml:"id"`
	Label    string            `json:"label"              yaml:"label"              toml:"label"`
	Children []treeNode        `json:"children,omitempty" yaml:"children,omitempty" toml:"children,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" toml:"metadata,omitempty"`
}

func toTreeNode(node *output.TreeNode) treeNode {
	result := treeNode{
		ID:       node.ID.Get(),
		Label:    node.Label.Get(),
		Metadata: node.Metadata,
	}

	if len(node.Children) > 0 {
		result.Children = make([]treeNode, 0, len(node.Children))
		for _, child := range node.Children {
			result.Children = append(result.Children, toTreeNode(child))
		}
	}

	return result
}
