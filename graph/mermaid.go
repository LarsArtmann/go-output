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

//nolint:gochecknoinits // Registers Mermaid format capabilities and TableRenderer.
func init() {
	output.RegisterFormatShapes(output.FormatMermaid, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableMarshaler(output.FormatMermaid, renderMermaidTable)
}

func renderMermaidTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return output.WriteRenderedFrom(w, NewMermaidFromTable(data).Render, "Mermaid", "render Mermaid")
}

// MermaidRenderer implements the GraphRenderer interface for Mermaid diagrams.
type MermaidRenderer struct {
	output.GraphBuilder

	// codeFence controls whether output is wrapped in a ```mermaid fence.
	// Default true for backwards compatibility. Set false via SetCodeFence
	// for raw flowchart syntax (.mmd files, Mermaid CLI, embedded diagrams).
	codeFence bool
}

// NewMermaidRenderer creates a new MermaidRenderer with the code fence enabled.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		GraphBuilder: *output.NewGraphBuilder(),
		codeFence:    true,
	}
}

// SetCodeFence controls whether Render wraps output in a ```mermaid fence.
// Pass false for raw flowchart syntax consumed by .mmd files, the Mermaid CLI,
// or programmatic APIs that expect bare flowchart syntax.
func (r *MermaidRenderer) SetCodeFence(enabled bool) *MermaidRenderer {
	r.codeFence = enabled
	return r
}

// Render returns the Mermaid diagram as a string.
func (r *MermaidRenderer) Render() (string, error) {
	var b strings.Builder

	if r.codeFence {
		b.WriteString("```mermaid\n")
	}

	b.WriteString("flowchart TD\n")

	for _, node := range r.Nodes() {
		prefix, suffix := r.getMermaidShape(node.Shape)
		label := escape.MermaidText(node.Label.Get())
		_, _ = fmt.Fprintf(&b, "    %s%s%s%s\n", escape.MermaidID(node.ID.Get()), prefix, label, suffix)
	}

	for _, edge := range r.Edges() {
		label := ""
		if !edge.Label.IsZero() {
			label = fmt.Sprintf("|%s|", escape.MermaidText(edge.Label.Get()))
		}

		_, _ = fmt.Fprintf(
			&b,
			"    %s -->%s %s\n",
			escape.MermaidID(edge.From.Get()),
			label,
			escape.MermaidID(edge.To.Get()),
		)
	}

	r.writeNodeStyles(&b)

	if r.codeFence {
		b.WriteString("```\n")
	}

	return b.String(), nil
}

// getMermaidShape returns the prefix and suffix for a Mermaid shape.
func (r *MermaidRenderer) getMermaidShape(shape output.NodeShape) (string, string) {
	switch shape {
	case output.NodeShapeDiamond:
		return "{", "}"
	case output.NodeShapeEllipse:
		return "(", ")"
	case output.NodeShapeCircle:
		return "((", "))"
	case output.NodeShapeHexagon:
		return "{{", "}}"
	case output.NodeShapeCylinder:
		return "[(", ")]"
	case output.NodeShapeParallelogram:
		return "[/", "/]"
	case output.NodeShapeBox:
		return "[", "]"
	default:
		return "[", "]"
	}
}

// writeNodeStyles emits per-node Mermaid style directives for nodes that have
// a non-zero NodeStyle. This replaces the previous hardcoded pink classDef,
// giving consumers full control over node appearance.
func (r *MermaidRenderer) writeNodeStyles(b *strings.Builder) {
	wroteAny := false

	for _, node := range r.Nodes() {
		parts := mermaidStyleParts(node.Style)
		if len(parts) == 0 {
			continue
		}

		if !wroteAny {
			b.WriteString("\n    %% Styling\n")

			wroteAny = true
		}

		_, _ = fmt.Fprintf(b, "    style %s %s\n", escape.MermaidID(node.ID.Get()), strings.Join(parts, ","))
	}
}

// mermaidStyleParts converts a NodeStyle into Mermaid style key-value pairs.
func mermaidStyleParts(s output.NodeStyle) []string {
	var parts []string

	if s.Fill != "" {
		parts = append(parts, "fill:"+escape.MermaidText(s.Fill))
	}

	if s.Stroke != "" {
		parts = append(parts, "stroke:"+escape.MermaidText(s.Stroke))
	}

	if s.FontColor != "" {
		parts = append(parts, "color:"+escape.MermaidText(s.FontColor))
	}

	if s.FontSize > 0 {
		parts = append(parts, fmt.Sprintf("font-size:%dpx", s.FontSize))
	}

	return parts
}

// NewMermaidFromTable creates a Mermaid flowchart from table data.
func NewMermaidFromTable(data *output.Table) *MermaidRenderer {
	renderer := NewMermaidRenderer()
	if data == nil {
		return renderer
	}

	renderer.SetNodesFromTable(data, func(_ int, n *output.GraphNode) {
		n.Shape = output.NodeShapeBox
	})

	return renderer
}

// NewMermaidFromTree converts a TreeNode to Mermaid.
func NewMermaidFromTree(root *output.TreeNode) *MermaidRenderer {
	return output.TreeToRenderer(NewMermaidRenderer, (*MermaidRenderer).addTreeNodes, root)
}

func mermaidTreeNodeID(node *output.TreeNode) string {
	if id := escape.MermaidID(node.ID.Get()); id != "" {
		return id
	}

	return escape.MermaidSlug(node.Label.Get())
}

func (r *MermaidRenderer) addTreeNodes(node *output.TreeNode, parentID string) {
	output.AddTreeNodes(
		&r.GraphBuilder,
		node,
		parentID,
		mermaidTreeNodeID,
		output.NodeShapeBox,
	)
}
