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
	if def, ok := LookupStatus(as); ok {
		return def.Name
	}

	return StatusStringUnknown
}

// StatusStringUnknown is the fallback for unrecognized activity statuses.
// Mirrors tui.WorkflowStateStringUnknown — both use "unknown" for the same semantic purpose.
const StatusStringUnknown = "unknown"

// GetSymbol returns the NOM-style symbol for the status.
func (as ActivityStatus) GetSymbol() Symbol {
	if def, ok := LookupStatus(as); ok {
		return def.Symbol
	}

	return "?"
}

// GetColor returns the lipgloss color for the status.
func (as ActivityStatus) GetColor() color.Color {
	if def, ok := LookupStatus(as); ok {
		return def.Color
	}

	return Colors.Fallback
}

// Interest returns the display priority for sorting: lower = more interesting.
// Order: failed > running > pending > completed.
func (as ActivityStatus) Interest() int {
	if def, ok := LookupStatus(as); ok {
		return def.Interest
	}

	return 4
}

// AllActivityStatuses returns the complete list of valid ActivityStatus values
// in ascending ID order. The list is dynamic: custom statuses registered via
// RegisterStatus are included automatically.

// ParseActivityStatus parses a string into an ActivityStatus.
// Returns an error for unrecognized values.
func ParseActivityStatus(s string) (ActivityStatus, error) {
	for _, status := range AllActivityStatuses() {
		if status.String() == s {
			return status, nil
		}
	}

	return ActivityStatusPending, &InvalidActivityStatusError{Value: s}
}

// IsValid returns true if the status is a recognized ActivityStatus value.
func (as ActivityStatus) IsValid() bool {
	return slices.Contains(AllActivityStatuses(), as)
}

// AllowedValues returns all valid status strings for CLI help text and config.
func (ActivityStatus) AllowedValues() []string {
	statuses := AllActivityStatuses()

	out := make([]string, 0, len(statuses))
	for _, s := range statuses {
		out = append(out, s.String())
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
