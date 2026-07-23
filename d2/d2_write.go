package d2

import (
	"fmt"
	"strings"
)

func (d *Diagram) writeStyleAttrs(b *strings.Builder, s NodeStyle, indent string) {
	if s.Fill != "" {
		fmt.Fprintf(b, "%sstyle.fill: %s\n", indent, d2Quote(s.Fill))
	}

	d.writeStyleColors(b, s.StrokeStyle, indent)
	d.writeStyleEffects(b, s, indent)
}

func (*Diagram) writeStyleColors(b *strings.Builder, s StrokeStyle, indent string) {
	if s.Stroke != "" {
		fmt.Fprintf(b, "%sstyle.stroke: %s\n", indent, d2Quote(s.Stroke))
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
		fmt.Fprintf(b, "%sstyle.font-color: %s\n", indent, d2Quote(s.FontColor))
	}
}

func (*Diagram) writeStyleEffects(b *strings.Builder, s NodeStyle, indent string) {
	if s.Opacity > 0 {
		opacity := s.Opacity
		if opacity > 1 {
			opacity = 1
		} else if opacity < 0 {
			opacity = 0
		}

		fmt.Fprintf(b, "%sstyle.opacity: %g\n", indent, opacity)
	}

	if s.Shadow {
		fmt.Fprintf(b, "%sshadow: true\n", indent)
	}

	if s.BorderRadius > 0 {
		fmt.Fprintf(b, "%sborder-radius: %d\n", indent, s.BorderRadius)
	}

	if s.TextTransform != TextTransformNone {
		fmt.Fprintf(b, "%sstyle.text-transform: %s\n", indent, d2Quote(string(s.TextTransform)))
	}
}

func (d *Diagram) writeEdge(b *strings.Builder, edge Edge) {
	from := d2Quote(edge.From.Get())
	to := d2Quote(edge.To.Get())

	label := ""
	if !edge.Label.IsZero() {
		label = ": " + d2Quote(edge.Label.Get())
	}

	if !edge.hasBlockAttrs() {
		fmt.Fprintf(b, "%s -> %s%s\n", from, to, label)
		return
	}

	fmt.Fprintf(b, "%s -> %s%s {\n", from, to, label)
	d.writeEdgeBlockAttrs(b, edge)
	b.WriteString("}\n")
}

func (d *Diagram) writeEdgeBlockAttrs(b *strings.Builder, edge Edge) {
	d.writeStyleColors(b, edge.Style.StrokeStyle, "  ")

	if edge.Style.Animated {
		b.WriteString("  style.animated: true\n")
	}

	if edge.SourceArrow != "" {
		fmt.Fprintf(b, "  source-arrowhead.shape: %s\n", edge.SourceArrow)
	}

	if edge.TargetArrow != "" {
		fmt.Fprintf(b, "  target-arrowhead.shape: %s\n", edge.TargetArrow)
	}
}
