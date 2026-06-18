package nom

import "github.com/larsartmann/go-output"

// NodeShape returns the root output.NodeShape that best represents this
// status in a diagram export (DOT, Mermaid, D2, PlantUML).
func (s ActivityStatus) NodeShape() output.NodeShape {
	switch s {
	case ActivityStatusFailed:
		return output.NodeShapeDiamond // diamonds signal attention
	case ActivityStatusCompleted:
		return output.NodeShapeRect // rect = done, stable
	case ActivityStatusRunning:
		return output.NodeShapeBox // box = active work
	case ActivityStatusPaused:
		return output.NodeShapeHexagon // hexagon = interrupted
	case ActivityStatusPending:
		return output.NodeShapeEllipse // ellipse = waiting
	default:
		return output.NodeShapeBox
	}
}

// GraphStyle returns the root output.GraphStyle (hex fill/stroke/font colors)
// for diagram export. These are hex strings consumed by DOT/Mermaid/D2 — they
// are separate from the lipgloss color.Color values used for terminal display.
func (s ActivityStatus) GraphStyle() output.GraphStyle {
	switch s {
	case ActivityStatusFailed:
		return output.GraphStyle{
			Fill:      "#dc2626", // red-600
			Stroke:    "#991b1b", // red-800
			FontColor: "#ffffff",
		}
	case ActivityStatusRunning:
		return output.GraphStyle{
			Fill:      "#16a34a", // green-600
			Stroke:    "#15803d", // green-700
			FontColor: "#ffffff",
		}
	case ActivityStatusCompleted:
		return output.GraphStyle{
			Fill:      "#6b7280", // gray-500
			Stroke:    "#4b5563", // gray-600
			FontColor: "#ffffff",
		}
	case ActivityStatusPaused:
		return output.GraphStyle{
			Fill:      "#f59e0b", // amber-500
			Stroke:    "#d97706", // amber-600
			FontColor: "#000000",
		}
	case ActivityStatusPending:
		return output.GraphStyle{
			Fill:      "#e5e7eb", // gray-200
			Stroke:    "#9ca3af", // gray-400
			FontColor: "#374151", // gray-700
		}
	default:
		return output.GraphStyle{}
	}
}
