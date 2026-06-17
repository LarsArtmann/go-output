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

	// codeFence controls whether output is wrapped in a ```mermaid fence.
	// Default true for backwards compatibility. Set false via SetCodeFence
	// for raw flowchart syntax (.mmd files, Mermaid CLI, embedded diagrams).
	codeFence bool
}

// NewMermaidRenderer creates a new MermaidRenderer with the code fence enabled.
func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{
		GraphRendererState: output.NewGraphRendererState(),
		codeFence:          true,
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
		_, _ = fmt.Fprintf(&b, "    %s%s%s%s\n", node.ID.Get(), prefix, label, suffix)
	}

	for _, edge := range r.Edges() {
		label := ""
		if !edge.Label.IsZero() {
			label = fmt.Sprintf("|%s|", escape.MermaidText(edge.Label.Get()))
		}

		_, _ = fmt.Fprintf(&b, "    %s -->%s %s\n", edge.From.Get(), label, edge.To.Get())
	}

	r.writeNodeStyles(&b)

	if r.codeFence {
		b.WriteString("```\n")
	}

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

// writeNodeStyles emits per-node Mermaid style directives for nodes that have
// a non-zero GraphStyle. This replaces the previous hardcoded pink classDef,
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

		_, _ = fmt.Fprintf(b, "    style %s %s\n", node.ID.Get(), strings.Join(parts, ","))
	}
}

// mermaidStyleParts converts a GraphStyle into Mermaid style key-value pairs.
func mermaidStyleParts(s output.GraphStyle) []string {
	var parts []string

	if s.FillColor != "" {
		parts = append(parts, "fill:"+s.FillColor)
	}

	if s.StrokeColor != "" {
		parts = append(parts, "stroke:"+s.StrokeColor)
	}

	if s.FontColor != "" {
		parts = append(parts, "color:"+s.FontColor)
	}

	if s.FontSize > 0 {
		parts = append(parts, fmt.Sprintf("font-size:%dpx", s.FontSize))
	}

	return parts
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
