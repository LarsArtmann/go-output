package nom

import "testing"

// TestEventConstantsUnique guards against duplicate event-name constants.
// Event dispatch is now an exhaustive type switch on sealed concrete types
// (not string matching), so a duplicate constant can no longer cause silent
// misrouting — but duplicate names would still confuse logging/metrics, so
// this guard remains valuable.
func TestEventConstantsUnique(t *testing.T) {
	t.Parallel()

	events := []struct {
		name string
		val  string
	}{
		{"EventWorkflowStarted", EventWorkflowStarted},
		{"EventWorkflowCompleted", EventWorkflowCompleted},
		{"EventWorkflowFailed", EventWorkflowFailed},
		{"EventActivityStarted", EventActivityStarted},
		{"EventActivityCompleted", EventActivityCompleted},
		{"EventActivityFailed", EventActivityFailed},
		{"EventActivityRegistered", EventActivityRegistered},
	}

	seen := make(map[string]string, len(events))
	for _, e := range events {
		if prev, dup := seen[e.val]; dup {
			t.Errorf("event constant %q duplicates %q (both = %q)", e.name, prev, e.val)
		}

		seen[e.val] = e.name
	}

	if len(seen) != len(events) {
		t.Errorf("unique event values = %d, want %d", len(seen), len(events))
	}
}

// TestEventConstantsNonEmpty ensures no constant is the empty string, which
// would collide with a zero-value Event.GetEventType() and silently route
// malformed events to that handler.
func TestEventConstantsNonEmpty(t *testing.T) {
	t.Parallel()

	for _, e := range []string{
		EventWorkflowStarted, EventWorkflowCompleted, EventWorkflowFailed,
		EventActivityStarted, EventActivityCompleted, EventActivityFailed,
		EventActivityRegistered,
	} {
		if e == "" {
			t.Errorf("event constant is empty; empty string collides with zero-value GetEventType()")
		}
	}
}
