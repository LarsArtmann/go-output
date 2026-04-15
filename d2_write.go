package output

import (
	"fmt"
	"strings"
)

func (d *D2Diagram) writeStyleAttrs(b *strings.Builder, s D2NodeStyle, indent string) {
	d.writeStyleColors(b, s, indent)
	d.writeStyleEffects(b, s, indent)
}

func (*D2Diagram) writeStyleColors(b *strings.Builder, s D2NodeStyle, indent string) {
	if s.Fill != "" {
		fmt.Fprintf(b, "%sstyle.fill: %s\n", indent, s.Fill)
	}

	if s.Stroke != "" {
		fmt.Fprintf(b, "%sstyle.stroke: %s\n", indent, s.Stroke)
	}

	if s.StrokeWidth > 0 {
		fmt.Fprintf(b, "%sstyle.stroke-width: %d\n", indent, s.StrokeWidth)
	}

	if s.StrokeDash > 0 {
		fmt.Fprintf(b, "%sstyle.stroke-dash: %d\n", indent, s.StrokeDash)
	}

	if s.FontSize > 0 {
		fmt.Fprintf(b, "%sstyle.font-size: %d\n", indent, s.FontSize)
	}

	if s.FontColor != "" {
		fmt.Fprintf(b, "%sstyle.font-color: %s\n", indent, s.FontColor)
	}
}

func (*D2Diagram) writeStyleEffects(b *strings.Builder, s D2NodeStyle, indent string) {
	if s.Opacity > 0 {
		fmt.Fprintf(b, "%sstyle.opacity: %g\n", indent, s.Opacity)
	}

	if s.Shadow {
		fmt.Fprintf(b, "%sshadow: true\n", indent)
	}

	if s.BorderRadius > 0 {
		fmt.Fprintf(b, "%sborder-radius: %d\n", indent, s.BorderRadius)
	}

	if s.TextTransform != "" {
		fmt.Fprintf(b, "%sstyle.text-transform: %s\n", indent, s.TextTransform)
	}
}

func (d *D2Diagram) writeEdge(b *strings.Builder, edge D2Edge) {
	from := escapeD2(edge.From.Get())
	to := escapeD2(edge.To.Get())

	label := ""
	if !edge.Label.IsEmpty() {
		label = ": " + escapeD2(edge.Label.Get())
	}

	if !edge.hasBlockAttrs() {
		fmt.Fprintf(b, "%s -> %s%s\n", from, to, label)
		return
	}

	fmt.Fprintf(b, "%s -> %s%s {\n", from, to, label)
	d.writeEdgeBlockAttrs(b, edge)
	b.WriteString("}\n")
}

func (*D2Diagram) writeEdgeBlockAttrs(b *strings.Builder, edge D2Edge) {
	s := edge.Style
	if s.Stroke != "" {
		fmt.Fprintf(b, "  style.stroke: %s\n", s.Stroke)
	}

	if s.StrokeWidth > 0 {
		fmt.Fprintf(b, "  style.stroke-width: %d\n", s.StrokeWidth)
	}

	if s.StrokeDash > 0 {
		fmt.Fprintf(b, "  style.stroke-dash: %d\n", s.StrokeDash)
	}

	if s.Animated {
		b.WriteString("  style.animated: true\n")
	}

	if s.FontColor != "" {
		fmt.Fprintf(b, "  style.font-color: %s\n", s.FontColor)
	}

	if s.FontSize > 0 {
		fmt.Fprintf(b, "  style.font-size: %d\n", s.FontSize)
	}

	if edge.SourceArrow != "" {
		fmt.Fprintf(b, "  source-arrowhead.shape: %s\n", edge.SourceArrow)
	}

	if edge.TargetArrow != "" {
		fmt.Fprintf(b, "  target-arrowhead.shape: %s\n", edge.TargetArrow)
	}
}

// escapeD2 escapes special characters for safe inclusion in D2 output.
func escapeD2(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)

	return s
}
