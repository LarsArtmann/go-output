package serialization

import (
	"fmt"

	"github.com/go-faster/yaml"
	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer           = (*YAMLTreeRenderer)(nil)
	_ output.TreeOutputRenderer = (*YAMLTreeRenderer)(nil)
	_ output.Renderer           = (*YAMLGraphRenderer)(nil)
	_ output.GraphRenderer      = (*YAMLGraphRenderer)(nil)
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

type yamlTreeNode struct {
	ID       string            `yaml:"id"`
	Label    string            `yaml:"label"`
	Children []yamlTreeNode    `yaml:"children,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

// Render returns the tree as a YAML string.
func (r *YAMLTreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "null\n", nil
	}

	node := r.toYAMLNode(r.root)

	data, err := yaml.Marshal(node)
	if err != nil {
		return "", fmt.Errorf("marshal yaml tree: %w", err)
	}

	return string(data), nil
}

func (r *YAMLTreeRenderer) toYAMLNode(node *output.TreeNode) yamlTreeNode {
	result := yamlTreeNode{
		ID:       node.ID.Get(),
		Label:    node.Label.Get(),
		Metadata: node.Metadata,
	}

	if len(node.Children) > 0 {
		result.Children = make([]yamlTreeNode, 0, len(node.Children))
		for _, child := range node.Children {
			result.Children = append(result.Children, r.toYAMLNode(child))
		}
	}

	return result
}

// YAMLGraphRenderer renders graph nodes and edges as YAML.
type YAMLGraphRenderer struct {
	output.GraphRendererMixin
}

// NewYAMLGraphRenderer creates a new YAMLGraphRenderer.
func NewYAMLGraphRenderer() *YAMLGraphRenderer {
	return &YAMLGraphRenderer{
		GraphRendererMixin: output.NewGraphRendererMixin(),
	}
}

type yamlGraph struct {
	Nodes []yamlGraphNode `yaml:"nodes"`
	Edges []yamlGraphEdge `yaml:"edges"`
}

type yamlGraphNode struct {
	ID       string            `yaml:"id"`
	Label    string            `yaml:"label"`
	Shape    string            `yaml:"shape,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

type yamlGraphEdge struct {
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Label string `yaml:"label,omitempty"`
}

// Render returns the graph as a YAML string.
func (r *YAMLGraphRenderer) Render() (string, error) {
	graph := yamlGraph{
		Nodes: make([]yamlGraphNode, 0, len(r.Nodes())),
		Edges: make([]yamlGraphEdge, 0, len(r.Edges())),
	}

	for _, node := range r.Nodes() {
		n := yamlGraphNode{
			ID:       node.ID.Get(),
			Label:    node.Label.Get(),
			Shape:    string(node.Shape),
			Metadata: node.Metadata,
		}

		graph.Nodes = append(graph.Nodes, n)
	}

	for _, edge := range r.Edges() {
		e := yamlGraphEdge{
			From:  edge.From.Get(),
			To:    edge.To.Get(),
			Label: brandedEdgeLabel(edge.Label),
		}

		graph.Edges = append(graph.Edges, e)
	}

	data, err := yaml.Marshal(graph)
	if err != nil {
		return "", fmt.Errorf("marshal yaml graph: %w", err)
	}

	return string(data), nil
}
