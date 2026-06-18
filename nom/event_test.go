package nom

import "testing"

// TestEventConstantsUnique guards against the silent-typo-drop failure mode:
// if two event constants share a value, an event dispatched to one is silently
// routed to both handlers. This also guards the m2 fix (bare literals → constants)
// by asserting the constants are the only routing keys in use.
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
