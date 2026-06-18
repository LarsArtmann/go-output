package nom

import "image/color"

// ============================================================================
// ACTIVITY STATUS ENUM
// ============================================================================
// ActivityStatus represents the lifecycle stage of a single activity within a
// workflow. It is intentionally SEPARATE from tui.WorkflowState (which tracks
// the whole-workflow lifecycle) — do NOT merge them. They share Running and
// Completed values but model different domains: an activity can be Paused while
// the workflow is still Running; a workflow Errored while individual activities
// report Failed. Split-brain finding m3: documented to prevent a future
// "unify these enums" mistake.
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
		return Colors.Paused
	case ActivityStatusRunning:
		return Colors.Running
	case ActivityStatusCompleted:
		return Colors.Completed
	case ActivityStatusFailed:
		return Colors.Failed
	case ActivityStatusPaused:
		return Colors.Paused
	default:
		return Colors.Info
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
