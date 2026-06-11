package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*MermaidRenderer)(nil)
	_ output.GraphRenderer = (*MermaidRenderer)(nil)
)

//nolint:gochecknoinits // Registers Mermaid format capabilities and TableDataMarshaler.
func init() {
	output.RegisterFormatShapes(output.FormatMermaid, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableDataMarshaler(output.FormatMermaid, renderMermaidTableData)
}

func renderMermaidTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	out, err := MermaidFromTableData(data).Render()
	if err != nil {
		return fmt.Errorf("render Mermaid: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write Mermaid output: %w", err)
	}

	return nil
}

// MermaidRenderer implements the GraphRenderer interface for Mermaid diagrams.
type MermaidRenderer struct {
	output.GraphRendererState
}

// NewMermaidRenderer creates a new MermaidRenderer.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		GraphRendererState: output.NewGraphRendererState(),
	}
}

// Render returns the Mermaid diagram as a string.
func (r *MermaidRenderer) Render() (string, error) {
	var b strings.Builder

	b.WriteString("```mermaid\n")
	b.WriteString("flowchart TD\n")

	for _, node := range r.Nodes() {
		prefix, suffix := r.getMermaidShape(node.Shape)
		label := escape.MermaidText(node.Label.Get())
		_, _ = fmt.Fprintf(&b, "    %s%s%s%s\n", node.ID.Get(), prefix, label, suffix)
	}

	for _, edge := range r.Edges() {
		label := ""
		if !edge.Label.IsZero() {
			label = fmt.Sprintf("|%s|", escape.MermaidText(edge.Label.Get()))
		}

		_, _ = fmt.Fprintf(&b, "    %s -->%s %s\n", edge.From.Get(), label, edge.To.Get())
	}

	b.WriteString("\n    %% Styling\n")
	b.WriteString("    classDef default fill:#f9f,stroke:#333,stroke-width:4px\n")

	b.WriteString("```\n")

	return b.String(), nil
}

// getMermaidShape returns the prefix and suffix for a Mermaid shape.
func (r *MermaidRenderer) getMermaidShape(shape output.GraphShape) (string, string) {
	switch shape {
	case output.ShapeDiamond:
		return "{", "}"
	case output.ShapeEllipse:
		return "(", ")"
	case output.ShapeCircle:
		return "((", "))"
	case output.ShapeHexagon:
		return "{{", "}}"
	case output.ShapeCylinder:
		return "[(", ")]"
	case output.ShapeParallelogram:
		return "[/", "/]"
	case output.ShapeBox, output.ShapeRect:
		return "[", "]"
	default:
		return "[", "]"
	}
}

// MermaidFromTableData creates a Mermaid flowchart from table data.
func MermaidFromTableData(data *output.TableData) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if data == nil {
		return renderer
	}

	renderer.SetNodesFromTableData(data, func(_ int, n *output.GraphNode) {
		n.Shape = output.ShapeBox
		n.Label = output.NewBrandedID[output.GraphNodeLabelBrand](escape.MermaidText(n.Label.Get()))
	})

	return renderer
}

// MermaidFromTree converts a TreeNode to Mermaid.
func MermaidFromTree(root *output.TreeNode) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if root == nil {
		return renderer
	}

	renderer.addTreeNodes(root, "")

	return renderer
}

func mermaidTreeNodeID(node *output.TreeNode) string {
	if id := escape.MermaidID(node.ID.Get()); id != "" {
		return id
	}

	return escape.MermaidSlug(node.Label.Get())
}

func (r *MermaidRenderer) addTreeNodes(node *output.TreeNode, parentID string) {
	output.AddTreeNodes(
		&r.GraphRendererState,
		node,
		parentID,
		mermaidTreeNodeID,
		output.ShapeBox,
	)
}
