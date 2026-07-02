package nom

import "github.com/larsartmann/go-output"

// NodeShape returns the root output.NodeShape that best represents this
// status in a diagram export (DOT, Mermaid, D2, PlantUML).
func (s ActivityStatus) NodeShape() output.NodeShape {
	if def, ok := LookupStatus(s); ok {
		return def.NodeShape
	}

	return output.NodeShapeBox
}

// GraphStyle returns the root output.GraphStyle (hex fill/stroke/font colors)
// for diagram export. These are hex strings consumed by DOT/Mermaid/D2 — they
// are separate from the lipgloss color.Color values used for terminal display.
func (s ActivityStatus) GraphStyle() output.GraphStyle {
	if def, ok := LookupStatus(s); ok {
		return def.GraphStyle
	}

	return output.GraphStyle{}
}
