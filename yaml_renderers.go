package output

import (
	"fmt"

	"github.com/go-faster/yaml"
)

// Compile-time interface checks.
var (
	_ Renderer           = (*YAMLTreeRenderer)(nil)
	_ TreeOutputRenderer = (*YAMLTreeRenderer)(nil)
	_ Renderer           = (*YAMLGraphRenderer)(nil)
	_ GraphRenderer      = (*YAMLGraphRenderer)(nil)
)

// YAMLTreeRenderer renders a TreeNode hierarchy as YAML.
type YAMLTreeRenderer struct {
	root *TreeNode
}

// NewYAMLTreeRenderer creates a new YAMLTreeRenderer.
func NewYAMLTreeRenderer() *YAMLTreeRenderer {
	return &YAMLTreeRenderer{} //nolint:exhaustruct // root is set via SetRoot
}

// SetRoot sets the root node of the tree.
func (r *YAMLTreeRenderer) SetRoot(node *TreeNode) {
	r.root = node
}

// yamlTreeNode is the YAML representation of a TreeNode.
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

func (r *YAMLTreeRenderer) toYAMLNode(node *TreeNode) yamlTreeNode {
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
	nodes []GraphNode
	edges []GraphEdge
}

// NewYAMLGraphRenderer creates a new YAMLGraphRenderer.
func NewYAMLGraphRenderer() *YAMLGraphRenderer {
	return &YAMLGraphRenderer{
		nodes: make([]GraphNode, 0),
		edges: make([]GraphEdge, 0),
	}
}

// SetNodes sets the graph nodes.
func (r *YAMLGraphRenderer) SetNodes(nodes []GraphNode) {
	r.nodes = nodes
}

// SetEdges sets the graph edges.
func (r *YAMLGraphRenderer) SetEdges(edges []GraphEdge) {
	r.edges = edges
}

// yamlGraph is the YAML representation of a graph.
type yamlGraph struct {
	Nodes []yamlGraphNode `yaml:"nodes"`
	Edges []yamlGraphEdge `yaml:"edges"`
}

// yamlGraphNode is the YAML representation of a GraphNode.
type yamlGraphNode struct {
	ID       string            `yaml:"id"`
	Label    string            `yaml:"label"`
	Shape    string            `yaml:"shape,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty"`
}

// yamlGraphEdge is the YAML representation of a GraphEdge.
type yamlGraphEdge struct {
	From  string `yaml:"from"`
	To    string `yaml:"to"`
	Label string `yaml:"label,omitempty"`
}

// Render returns the graph as a YAML string.
func (r *YAMLGraphRenderer) Render() (string, error) {
	graph := yamlGraph{
		Nodes: make([]yamlGraphNode, 0, len(r.nodes)),
		Edges: make([]yamlGraphEdge, 0, len(r.edges)),
	}

	for _, node := range r.nodes {
		n := yamlGraphNode{
			ID:       node.ID.Get(),
			Label:    node.Label.Get(),
			Shape:    string(node.Shape),
			Metadata: node.Metadata,
		}

		graph.Nodes = append(graph.Nodes, n)
	}

	for _, edge := range r.edges {
		e := yamlGraphEdge{
			From: edge.From.Get(),
			To:   edge.To.Get(),
		}

		if !edge.Label.IsZero() {
			e.Label = edge.Label.Get()
		}

		graph.Edges = append(graph.Edges, e)
	}

	data, err := yaml.Marshal(graph)
	if err != nil {
		return "", fmt.Errorf("marshal yaml graph: %w", err)
	}

	return string(data), nil
}
