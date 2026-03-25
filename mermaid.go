package output

import (
	"fmt"
	"strings"
)

// MermaidRenderer implements the GraphRenderer interface for Mermaid diagrams.
type MermaidRenderer struct {
	nodes []GraphNode
	edges []GraphEdge
}

// NewMermaidRenderer creates a new MermaidRenderer.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		nodes: make([]GraphNode, 0),
		edges: make([]GraphEdge, 0),
	}
}

// SetNodes sets the graph nodes.
func (r *MermaidRenderer) SetNodes(nodes []GraphNode) {
	r.nodes = nodes
}

// SetEdges sets the graph edges.
func (r *MermaidRenderer) SetEdges(edges []GraphEdge) {
	r.edges = edges
}

// Render returns the Mermaid diagram as a string.
func (r *MermaidRenderer) Render() string {
	var b strings.Builder

	b.WriteString("```mermaid\n")
	b.WriteString("flowchart TD\n")

	// Write nodes
	for _, node := range r.nodes {
		prefix, suffix := r.getMermaidShape(node.Shape)
		label := r.escapeMermaidLabel(node.Label.Get())
		_, _ = fmt.Fprintf(&b, "    %s%s%s%s\n", node.ID.Get(), prefix, label, suffix)
	}

	// Write edges
	for _, edge := range r.edges {
		label := ""
		if !edge.Label.IsEmpty() {
			label = fmt.Sprintf("|%s|", r.escapeMermaidLabel(edge.Label.Get()))
		}
		_, _ = fmt.Fprintf(&b, "    %s -->%s %s\n", edge.From.Get(), label, edge.To.Get())
	}

	// Write styling
	b.WriteString("\n    %% Styling\n")
	b.WriteString("    classDef default fill:#f9f,stroke:#333,stroke-width:4px\n")

	b.WriteString("```\n")

	return b.String()
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

func (r *MermaidRenderer) escapeMermaidLabel(s string) string {
	// Escape quotes and brackets for Mermaid
	s = strings.ReplaceAll(s, "\"", "'")
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	s = strings.ReplaceAll(s, "{", "(")
	s = strings.ReplaceAll(s, "}", ")")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return s
}

// MermaidFlowchartRenderer creates a Mermaid flowchart from table data.
func MermaidFlowchartRenderer(data *TableData) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if data == nil {
		return renderer
	}

	// Create nodes for each row
	for i, row := range data.Rows {
		var labelParts []string
		for j, cell := range row {
			if j < len(data.Headers) {
				labelParts = append(labelParts, data.Headers[j]+": "+cell)
			} else {
				labelParts = append(labelParts, cell)
			}
		}
		label := strings.Join(labelParts, "<br>")
		renderer.nodes = append(renderer.nodes, GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i)),
			Label: NewBrandedID[GraphNodeLabelBrand](label),
			Shape: ShapeBox,
		})
	}

	// Create edges between consecutive rows using shared helper
	for _, edge := range data.CreateRowEdges() {
		renderer.edges = append(renderer.edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](edge.From),
			To:   NewBrandedID[GraphNodeIDBrand](edge.To),
		})
	}

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

func (r *MermaidRenderer) addTreeNodes(node *TreeNode, parentID string) {
	nodeID := sanitizeMermaidID(node.ID.Get())
	if nodeID == "" {
		nodeID = sanitizeMermaidLabel(node.Label.Get())
	}

	graphNodeID := NewBrandedID[GraphNodeIDBrand](nodeID)
	graphNodeLabel := NewBrandedID[GraphNodeLabelBrand](node.Label.Get())

	r.nodes = append(r.nodes, GraphNode{
		ID:    graphNodeID,
		Label: graphNodeLabel,
		Shape: ShapeBox,
	})

	if parentID != "" {
		r.edges = append(r.edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](parentID),
			To:   graphNodeID,
		})
	}

	for _, child := range node.Children {
		r.addTreeNodes(child, nodeID)
	}
}

func sanitizeMermaidID(id string) string {
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	if result.Len() == 0 {
		return "node"
	}
	return result.String()
}

func sanitizeMermaidLabel(label string) string {
	result := strings.ReplaceAll(label, " ", "_")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, "/", "_")
	return result
}
