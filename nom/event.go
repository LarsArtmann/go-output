package nom

import "context"

// Event type constants. Use these instead of bare string literals to avoid
// silent typos in event dispatch. The string-based interface is preserved
// for backward compatibility with callers that implement Event externally.
const (
	EventWorkflowStarted      = "workflow.started"
	EventWorkflowCompleted    = "workflow.completed"
	EventWorkflowFailed       = "workflow.failed"
	EventActivityStarted      = "activity.started"
	EventActivityCompleted    = "activity.completed"
	EventActivityFailed       = "activity.failed"
	EventActivityRegistered   = "activity.registered"
)

type Event interface {
	GetEventType() string
}
type EventSubscriber interface {
	OnEvent(ctx context.Context, event Event) error
}
