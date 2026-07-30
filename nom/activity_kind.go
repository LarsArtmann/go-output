package nom

import (
	"slices"
	"strings"

	"github.com/larsartmann/go-output"
)

// ActivityKind classifies what an activity IS in the workflow topology,
// independent of its lifecycle Status. A Task is a concrete unit of work; a
// Phase is a grouping/container node that owns child activities. Kind is set
// at construction and never changes — unlike Status, it is not a lifecycle
// state.
//
// This replaces the older "phase:" ID-prefix convention, which was a lying
// name: the prefix claimed to signal kind but was really just an opaque
// string matched by HasPrefix. A typed Kind makes the intent explicit and
// impossible to confuse with the activity's identity.
type ActivityKind int

const (
	// ActivityKindTask is the default kind: a concrete unit of work.
	ActivityKindTask ActivityKind = iota
	// ActivityKindPhase marks a grouping node whose children are the real
	// deliverables. Phase nodes render with SymbolPhase/Colors.Phase.
	ActivityKindPhase
)

// String returns the lowercase name of the kind.
func (k ActivityKind) String() string {
	switch k {
	case ActivityKindTask:
		return "task"
	case ActivityKindPhase:
		return "phase"
	default:
		return "unknown"
	}
}

// IsPhase reports whether this kind is a phase grouping node.
func (k ActivityKind) IsPhase() bool { return k == ActivityKindPhase }

// AllActivityKinds is the complete list of valid ActivityKind values.
//
//nolint:gochecknoglobals // Global used for value iteration.
var AllActivityKinds = []ActivityKind{
	ActivityKindTask,
	ActivityKindPhase,
}

// ParseActivityKind parses a string into an ActivityKind.
// Returns an error for unrecognized values.
func ParseActivityKind(s string) (ActivityKind, error) {
	for _, k := range AllActivityKinds {
		if k.String() == s {
			return k, nil
		}
	}

	return ActivityKindTask, &InvalidActivityKindError{Value: s, Allowed: AllActivityKinds}
}

// IsValid returns true if the kind is a recognized ActivityKind value.
func (k ActivityKind) IsValid() bool { return slices.Contains(AllActivityKinds, k) }

// InvalidActivityKindError represents an invalid activity kind.
type InvalidActivityKindError struct {
	Value   string
	Allowed []ActivityKind
}

func (e *InvalidActivityKindError) Error() string {
	if len(e.Allowed) == 0 {
		return "invalid activity kind: " + e.Value
	}

	return "invalid activity kind: " + e.Value + " (allowed: " + strings.Join(output.EnumAllowedValues(e.Allowed), ", ") + ")"
}
