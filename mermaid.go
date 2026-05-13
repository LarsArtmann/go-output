package output

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ Renderer      = (*MermaidRenderer)(nil)
	_ GraphRenderer = (*MermaidRenderer)(nil)
)

// MermaidRenderer implements the GraphRenderer interface for Mermaid diagrams.
type MermaidRenderer struct {
	GraphRendererMixin
}

// NewMermaidRenderer creates a new MermaidRenderer.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		GraphRendererMixin: NewGraphRendererMixin(),
	}
}

// Render returns the Mermaid diagram as a string.
func (r *MermaidRenderer) Render() (string, error) {
	var b strings.Builder

	b.WriteString("```mermaid\n")
	b.WriteString("flowchart TD\n")

	// Write nodes
	for _, node := range r.nodes {
		prefix, suffix := r.getMermaidShape(node.Shape)
		label := escape.MermaidText(node.Label.Get())
		_, _ = fmt.Fprintf(&b, "    %s%s%s%s\n", node.ID.Get(), prefix, label, suffix)
	}

	// Write edges
	for _, edge := range r.edges {
		label := ""
		if !edge.Label.IsZero() {
			label = fmt.Sprintf("|%s|", escape.MermaidText(edge.Label.Get()))
		}

		_, _ = fmt.Fprintf(&b, "    %s -->%s %s\n", edge.From.Get(), label, edge.To.Get())
	}

	// Write styling
	b.WriteString("\n    %% Styling\n")
	b.WriteString("    classDef default fill:#f9f,stroke:#333,stroke-width:4px\n")

	b.WriteString("```\n")

	return b.String(), nil
}

// getMermaidShape returns the prefix and suffix for a Mermaid shape.
func (r *MermaidRenderer) getMermaidShape(shape GraphShape) (string, string) {
	switch shape {
	case ShapeDiamond:
		return "{", "}"
	case ShapeEllipse:
		return "(", ")"
	case ShapeCircle:
		return "((", "))"
	case ShapeHexagon:
		return "{{", "}}"
	case ShapeCylinder:
		return "[(", ")]"
	case ShapeParallelogram:
		return "[/", "/]"
	case ShapeBox, ShapeRect:
		return "[", "]"
	default:
		return "[", "]"
	}
}

// MermaidFlowchartRenderer creates a Mermaid flowchart from table data.
func MermaidFlowchartRenderer(data *TableData) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if data == nil {
		return renderer
	}

	renderer.SetNodesFromTableData(data, func(_ int, n *GraphNode) {
		n.Shape = ShapeBox
		n.Label = NewBrandedID[GraphNodeLabelBrand](escape.MermaidText(n.Label.Get()))
	})

	return renderer
}

// MermaidTreeRenderer converts a TreeNode to Mermaid.
func MermaidTreeRenderer(root *TreeNode) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if root == nil {
		return renderer
	}

	renderer.addTreeNodes(root, "")

	return renderer
}

func mermaidTreeNodeID(node *TreeNode) string {
	if id := escape.MermaidID(node.ID.Get()); id != "" {
		return id
	}

	return escape.MermaidSlug(node.Label.Get())
}

func (r *MermaidRenderer) addTreeNodes(node *TreeNode, parentID string) {
	AddTreeNodes(&r.nodes, &r.edges, node, parentID, mermaidTreeNodeID, ShapeBox)
}
