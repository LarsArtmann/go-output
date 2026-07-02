package nom

import (
	"context"
	"time"
)

// Event is a sealed sum type representing a workflow or activity lifecycle
// event. Concrete implementations live in this package only — the unexported
// isEvent marker prevents external implementations, so OnEvent dispatch is an
// exhaustive Go type switch. This makes the silent-typo-drop failure mode of
// the old string-based GetEventType() dispatch unrepresentable: the compiler
// rejects unhandled event types.
type Event interface {
	isEvent()
}

// Event name constants identify each event type for logging and metrics. They
// are no longer used for dispatch (the type switch handles that), but are kept
// as stable identifiers and guarded for uniqueness by event_test.go.
const (
	EventWorkflowStarted    = "workflow.started"
	EventWorkflowCompleted  = "workflow.completed"
	EventWorkflowFailed     = "workflow.failed"
	EventActivityStarted    = "activity.started"
	EventActivityCompleted  = "activity.completed"
	EventActivityFailed     = "activity.failed"
	EventActivityRegistered = "activity.registered"
	EventActivityProgress   = "activity.progress"
	EventActivityRetrying   = "activity.retrying"
)

// --- Workflow lifecycle events ---

// WorkflowStarted signals that a workflow has begun.
type WorkflowStarted struct {
	ID   WorkflowID
	Name WorkflowName
}

func (WorkflowStarted) isEvent() {}

// WorkflowCompleted signals successful workflow completion.
type WorkflowCompleted struct {
	ID WorkflowID
}

func (WorkflowCompleted) isEvent() {}

// WorkflowFailed signals workflow failure; Err carries the cause.
type WorkflowFailed struct {
	ID  WorkflowID
	Err error
}

func (WorkflowFailed) isEvent() {}

// --- Activity lifecycle events ---

// ActivityStarted signals that an activity has begun running. Kind, Deps,
// Host, Download, and Category are optional (zero-valued = absent).
type ActivityStarted struct {
	ID       ActivityID
	Name     ActivityName
	Kind     ActivityKind
	Deps     []ActivityID
	Host     string
	Download DownloadProgress
	Category ActivityCategory
}

func (ActivityStarted) isEvent() {}

// ActivityRegistered pre-creates an activity in the tree as pending, without
// transitioning it to running. Used for pre-declaring structure (e.g. phases
// and their children) before the work actually starts.
type ActivityRegistered struct {
	ID       ActivityID
	Name     ActivityName
	Kind     ActivityKind
	Deps     []ActivityID
	Category ActivityCategory
}

func (ActivityRegistered) isEvent() {}

// ActivityCompleted signals successful activity completion. Duration is the
// observed run time (also recorded in the timing cache).
type ActivityCompleted struct {
	ID       ActivityID
	Name     ActivityName
	Duration time.Duration
}

func (ActivityCompleted) isEvent() {}

// ActivityFailed signals activity failure; Err carries the cause.
type ActivityFailed struct {
	ID       ActivityID
	Name     ActivityName
	Err      error
	Duration time.Duration
}

func (ActivityFailed) isEvent() {}

// ActivityProgress signals a live progress update for a running activity.
// This enables sub-step visibility: a single activity like "go-mod-tidy" can
// report "Tidying module [2/26]: modules/gitignore" as it iterates. The
// Message is rendered as a dim sub-line beneath the activity label.
//
// Empty Message clears any prior progress message.
//
// Throttling: the renderer redraws on every progress event, so callers that
// iterate fast (e.g. per-file in a large tree) should rate-limit updates to
// roughly 1/sec to avoid excessive redraws. The simplest pattern is a time-based
// guard: only send when the message changed AND at least a second elapsed.
type ActivityProgress struct {
	ID      ActivityID
	Name    ActivityName
	Message string
}

func (ActivityProgress) isEvent() {}

// ActivityRetrying signals that a failed activity is being retried. The
// Attempt field is the retry attempt number (1 = first retry). The activity
// transitions back to Running and a ⟳ suffix with the attempt count is
// rendered. Reason optionally explains the retry cause (e.g. "timeout",
// "network") and renders as "⟳2 (timeout)" when non-empty.
type ActivityRetrying struct {
	ID      ActivityID
	Name    ActivityName
	Attempt int
	Reason  string
}

func (ActivityRetrying) isEvent() {}

// EventSubscriber receives lifecycle events. Implementations dispatch via Go
// type switch on the concrete event types above.
type EventSubscriber interface {
	OnEvent(ctx context.Context, event Event) error
}
