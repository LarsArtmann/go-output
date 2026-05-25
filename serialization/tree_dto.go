package serialization

import "github.com/larsartmann/go-output"

type treeNodeDTO struct {
	ID       string            `json:"id"                 toml:"id"                 yaml:"id"`
	Label    string            `json:"label"              toml:"label"              yaml:"label"`
	Children []treeNodeDTO     `json:"children,omitempty" toml:"children,omitempty" yaml:"children,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" toml:"metadata,omitempty" yaml:"metadata,omitempty"`
}

func toTreeNodeDTO(node *output.TreeNode) treeNodeDTO {
	result := treeNodeDTO{
		ID:       node.ID.Get(),
		Label:    node.Label.Get(),
		Metadata: node.Metadata,
	}

	if len(node.Children) > 0 {
		result.Children = make([]treeNodeDTO, 0, len(node.Children))
		for _, child := range node.Children {
			result.Children = append(result.Children, toTreeNodeDTO(child))
		}
	}

	return result
}
