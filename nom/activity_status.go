package nom

import "image/color"

// ============================================================================
// ACTIVITY STATUS ENUM
// ============================================================================
// ActivityStatus represents the display status of an activity.
type ActivityStatus int

const (
	ActivityStatusPending ActivityStatus = iota
	ActivityStatusRunning
	ActivityStatusCompleted
	ActivityStatusFailed
	ActivityStatusPaused
)

// String returns the string representation of activity status.
func (as ActivityStatus) String() string {
	switch as {
	case ActivityStatusPending:
		return "pending"
	case ActivityStatusRunning:
		return "running"
	case ActivityStatusCompleted:
		return "completed"
	case ActivityStatusFailed:
		return "failed"
	case ActivityStatusPaused:
		return "paused"
	default:
		return StatusStringUnknown
	}
}

// StatusStringUnknown is the fallback for unrecognized activity statuses.
// Mirrors tui.WorkflowStateStringUnknown — both use "unknown" for the same semantic purpose.
const StatusStringUnknown = "unknown"

// GetSymbol returns the NOM-style symbol for the status.
func (as ActivityStatus) GetSymbol() string {
	switch as {
	case ActivityStatusPending:
		return SymbolPaused
	case ActivityStatusRunning:
		return SymbolRunning
	case ActivityStatusCompleted:
		return SymbolCompleted
	case ActivityStatusFailed:
		return SymbolFailed
	case ActivityStatusPaused:
		return SymbolPaused
	default:
		return "?"
	}
}

// GetColor returns the lipgloss color for the status.
func (as ActivityStatus) GetColor() color.Color {
	switch as {
	case ActivityStatusPending:
		return ColorPaused
	case ActivityStatusRunning:
		return ColorRunning
	case ActivityStatusCompleted:
		return ColorCompleted
	case ActivityStatusFailed:
		return ColorFailed
	case ActivityStatusPaused:
		return ColorPaused
	default:
		return ColorInfo
	}
}

// Interest returns the display priority for sorting: lower = more interesting.
// Order: failed > running > paused > pending > completed.
func (as ActivityStatus) Interest() int {
	switch as {
	case ActivityStatusFailed:
		return 0
	case ActivityStatusRunning:
		return 1
	case ActivityStatusPaused:
		return 2
	case ActivityStatusPending:
		return 3
	case ActivityStatusCompleted:
		return 4
	default:
		return 5
	}
}
