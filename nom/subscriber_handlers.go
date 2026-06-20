package nom

import (
	"context"
	"time"
)

// OnEvent implements EventSubscriber interface
// Handles workflow and activity events to update NOM-style visualization.
// Uses string-based event type routing instead of concrete type switching,
// enabling extraction into an independent module without circular dependencies.
func (ns *NOMStyleSubscriber) OnEvent(ctx context.Context, event Event) error {
	switch event.GetEventType() {
	case EventWorkflowStarted:
		return ns.handleWorkflowStarted(ctx, event)
	case EventWorkflowCompleted:
		return ns.handleWorkflowCompleted(ctx, event)
	case EventWorkflowFailed:
		return ns.handleWorkflowFailed(ctx, event)
	case EventActivityStarted:
		return ns.handleActivityStarted(ctx, event)
	case EventActivityRegistered:
		return ns.handleActivityRegistered(ctx, event)
	case EventActivityCompleted:
		return ns.handleActivityCompleted(ctx, event)
	case EventActivityFailed:
		return ns.handleActivityFailed(ctx, event)
	default:
		return nil
	}
}

// WorkflowEventAccessor extracts workflow-level fields from any event.
type WorkflowEventAccessor interface {
	GetWorkflowID() WorkflowID
}

// WorkflowNameAccessor extracts the workflow name from start.
type WorkflowNameAccessor interface {
	GetWorkflowName() WorkflowName
}

// ActivityEventAccessor extracts activity-level fields from any activity event.
type ActivityEventAccessor interface {
	GetActivityID() ActivityID
	GetActivityName() ActivityName
}

// DurationAccessor extracts duration from completed/failed.
type DurationAccessor interface {
	GetDuration() time.Duration
}

// ErrorAccessor extracts error from failed.
type ErrorAccessor interface {
	GetError() error
}

// DependenciesAccessor extracts parent activity IDs for tree structure.
// When implemented on an "activity.started" event, the subscriber uses these
// IDs as the activity's parents in the dependency tree. The first ID becomes
// the primary parent (tree edge); all IDs get the activity added as a child.
type DependenciesAccessor interface {
	GetDependencies() []ActivityID
}

// HostAccessor optionally names where an activity runs (e.g. a build machine).
// Rendered as a dim tag when non-empty. Backward-compatible: events that don't
// implement it simply leave Host unset.
type HostAccessor interface {
	GetHost() string
}

// DownloadAccessor optionally reports byte-progress for an activity (e.g. a
// dependency download). Rendered as a progress bar when active. Backward-
// compatible: events that don't implement it leave Download at zero.
type DownloadAccessor interface {
	GetDownload() DownloadProgress
}

// handleWorkflowStarted handles workflow started event.
func (ns *NOMStyleSubscriber) handleWorkflowStarted(
	_ context.Context,
	event Event,
) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	wa, ok := event.(WorkflowEventAccessor)
	if !ok {
		return nil
	}

	ns.workflowID = wa.GetWorkflowID()
	ns.startTime = time.Now()
	ns.isRunning = true

	if na, ok := event.(WorkflowNameAccessor); ok {
		ns.workflowName = na.GetWorkflowName()
	}

	// Preserve pre-registered activities and tree structure (e.g. from
	// ProgressBridge.Start()). Only initialize empty maps on a fresh start
	// so callers can register phases/steps before workflow.started.
	if ns.activities == nil {
		ns.activities = make(map[ActivityID]*Activity)
	}

	if ns.dependencyTree == nil {
		ns.dependencyTree = NewDependencyTree()
	}

	return ns.timingCache.EnsureLoaded()
}

// handleWorkflowFinished handles common logic for workflow completed/failed events.
// Stops the workflow and persists timing data.
func (ns *NOMStyleSubscriber) handleWorkflowFinished(_ Event) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.isRunning = false

	return ns.timingCache.Save()
}

// handleWorkflowCompleted handles workflow completed event.
func (ns *NOMStyleSubscriber) handleWorkflowCompleted(
	_ context.Context,
	event Event,
) error {
	return ns.handleWorkflowFinished(event)
}

// handleWorkflowFailed handles workflow failed event.
func (ns *NOMStyleSubscriber) handleWorkflowFailed(
	_ context.Context,
	event Event,
) error {
	return ns.handleWorkflowFinished(event)
}

// getOrCreateActivity retrieves an existing activity or creates a new one.
// Must be called while holding ns.mu lock.
func (ns *NOMStyleSubscriber) getOrCreateActivity(
	activityID ActivityID,
	activityName ActivityName,
) *Activity {
	activity, exists := ns.activities[activityID]
	if !exists {
		activity = NewActivity(string(activityID), activityName.String())
		ns.activities[activityID] = activity
	}

	return activity
}

// extractDependencies returns the activity's dependencies if the event implements
// DependenciesAccessor; otherwise it returns nil. Both handleActivityStarted and
// handleActivityRegistered need this for AddActivity.
func extractDependencies(event Event) []ActivityID {
	var deps []ActivityID
	if da, ok := event.(DependenciesAccessor); ok {
		deps = da.GetDependencies()
	}

	return deps
}

// handleActivityStarted handles activity started event.
func (ns *NOMStyleSubscriber) handleActivityStarted(
	_ context.Context,
	event Event,
) error {
	aa, ok := event.(ActivityEventAccessor)
	if !ok {
		return nil
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(aa.GetActivityID(), aa.GetActivityName())
	activity.SetRunning()

	// Optional host/download annotations — dormant unless the event carries them.
	if ha, ok := event.(HostAccessor); ok {
		activity.Host = ha.GetHost()
	}

	if da, ok := event.(DownloadAccessor); ok {
		activity.Download = da.GetDownload()
	}

	medianDuration := ns.timingCache.GetMedian(aa.GetActivityName().String())
	if medianDuration > 0 {
		activity.SetEstimatedTime(medianDuration)
	}

	// The tree stores only the activity ID; the subscriber's activities map
	// (populated by getOrCreateActivity above) is the source of truth for
	// mutable fields, accessed via SnapshotActivities at render time.
	return ns.dependencyTree.AddActivity(
		aa.GetActivityID(),
		extractDependencies(event),
	)
}

// handleActivityRegistered registers an activity in the dependency tree as pending.
// Unlike handleActivityStarted, this does NOT set the activity to running —
// it only creates the tree node with pending status for pre-registration.
func (ns *NOMStyleSubscriber) handleActivityRegistered(
	_ context.Context,
	event Event,
) error {
	aa, ok := event.(ActivityEventAccessor)
	if !ok {
		return nil
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.getOrCreateActivity(aa.GetActivityID(), aa.GetActivityName())

	return ns.dependencyTree.AddActivity(
		aa.GetActivityID(),
		extractDependencies(event),
	)
}

// handleActivityCompleted handles activity completed event.
func (ns *NOMStyleSubscriber) handleActivityCompleted(
	_ context.Context,
	event Event,
) error {
	aa, ok := event.(ActivityEventAccessor)
	if !ok {
		return nil
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(aa.GetActivityID(), aa.GetActivityName())
	activity.SetCompleted()

	return ns.recordTimingAndPersist(event, aa)
}

// handleActivityFailed handles activity failed event.
func (ns *NOMStyleSubscriber) handleActivityFailed(
	_ context.Context,
	event Event,
) error {
	aa, ok := event.(ActivityEventAccessor)
	if !ok {
		return nil
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(aa.GetActivityID(), aa.GetActivityName())

	var eventErr error
	if ea, ok := event.(ErrorAccessor); ok {
		eventErr = ea.GetError()
	}

	activity.SetFailed(eventErr)

	return ns.recordTimingAndPersist(event, aa)
}

// recordTimingAndPersist extracts duration from the event and records it
// to the timing cache. The activity's state was already updated by the caller
// via the shared Activity pointer, so no tree sync is needed.
func (ns *NOMStyleSubscriber) recordTimingAndPersist(
	event Event,
	aa ActivityEventAccessor,
) error {
	if da, ok := event.(DurationAccessor); ok {
		duration := da.GetDuration()
		if duration > 0 {
			if err := ns.timingCache.Record(aa.GetActivityName().String(), duration); err != nil {
				return err
			}
		}
	}

	return nil
}
