package nom

import (
	"image/color"
	"slices"
)

// ============================================================================
// ACTIVITY STATUS ENUM
// ============================================================================
// ActivityStatus represents the lifecycle stage of a single activity within a
// workflow. It is intentionally SEPARATE from tui.WorkflowState (which tracks
// the whole-workflow lifecycle) — do NOT merge them. They share Running and
// Completed values but model different domains: an activity can be Pending
// while the workflow is still Running; a workflow Errored while individual
// activities report Failed. Split-brain finding m3: documented to prevent a
// future "unify these enums" mistake.
type ActivityStatus int

const (
	ActivityStatusPending ActivityStatus = iota
	ActivityStatusRunning
	ActivityStatusCompleted
	ActivityStatusFailed
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
	default:
		return StatusStringUnknown
	}
}

// StatusStringUnknown is the fallback for unrecognized activity statuses.
// Mirrors tui.WorkflowStateStringUnknown — both use "unknown" for the same semantic purpose.
const StatusStringUnknown = "unknown"

// GetSymbol returns the NOM-style symbol for the status.
func (as ActivityStatus) GetSymbol() Symbol {
	switch as {
	case ActivityStatusPending:
		return SymbolPending
	case ActivityStatusRunning:
		return SymbolRunning
	case ActivityStatusCompleted:
		return SymbolCompleted
	case ActivityStatusFailed:
		return SymbolFailed
	default:
		return "?"
	}
}

// GetColor returns the lipgloss color for the status.
func (as ActivityStatus) GetColor() color.Color {
	switch as {
	case ActivityStatusPending:
		return Colors.Pending
	case ActivityStatusRunning:
		return Colors.Running
	case ActivityStatusCompleted:
		return Colors.Completed
	case ActivityStatusFailed:
		return Colors.Failed
	default:
		return Colors.Info
	}
}

// Interest returns the display priority for sorting: lower = more interesting.
// Order: failed > running > pending > completed.
func (as ActivityStatus) Interest() int {
	switch as {
	case ActivityStatusFailed:
		return 0
	case ActivityStatusRunning:
		return 1
	case ActivityStatusPending:
		return 2
	case ActivityStatusCompleted:
		return 3
	default:
		return 4
	}
}

// AllActivityStatuses is the complete list of valid ActivityStatus values.
//
//nolint:gochecknoglobals // Global variable used for value iteration.
var AllActivityStatuses = []ActivityStatus{
	ActivityStatusPending,
	ActivityStatusRunning,
	ActivityStatusCompleted,
	ActivityStatusFailed,
}

// ParseActivityStatus parses a string into an ActivityStatus.
// Returns an error for unrecognized values.
func ParseActivityStatus(s string) (ActivityStatus, error) {
	for _, status := range AllActivityStatuses {
		if status.String() == s {
			return status, nil
		}
	}

	return ActivityStatusPending, &InvalidActivityStatusError{Value: s}
}

// IsValid returns true if the status is a recognized ActivityStatus value.
func (as ActivityStatus) IsValid() bool {
	return slices.Contains(AllActivityStatuses, as)
}

// AllowedValues returns all valid status strings for CLI help text and config.
func (ActivityStatus) AllowedValues() []string {
	out := make([]string, len(AllActivityStatuses))
	for i, s := range AllActivityStatuses {
		out[i] = s.String()
	}

	return out
}

// InvalidActivityStatusError represents an invalid activity status.
type InvalidActivityStatusError struct {
	Value string
}

func (e *InvalidActivityStatusError) Error() string {
	return "invalid activity status: " + e.Value
}
