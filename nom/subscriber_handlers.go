package nom

import (
	"context"
	"fmt"
	"time"
)

// OnEvent implements EventSubscriber interface
// Handles workflow and activity events to update NOM-style visualization.
// Uses string-based event type routing instead of concrete type switching,
// enabling extraction into an independent module without circular dependencies.
func (ns *NOMStyleSubscriber) OnEvent(ctx context.Context, event Event) error {
	switch event.GetEventType() {
	case "workflow.started":
		return ns.handleWorkflowStarted(ctx, event)
	case "workflow.completed":
		return ns.handleWorkflowCompleted(ctx, event)
	case "workflow.failed":
		return ns.handleWorkflowFailed(ctx, event)
	case "activity.started":
		return ns.handleActivityStarted(ctx, event)
	case "activity.completed":
		return ns.handleActivityCompleted(ctx, event)
	case "activity.failed":
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

	if err := ns.timingCache.EnsureLoaded(); err != nil {
		fmt.Printf("Warning: Failed to load timing cache: %v\n", err)
	}

	ns.activities = make(map[ActivityID]*ActivityDisplayState)
	ns.dependencyTree.Clear()

	return nil
}

// handleWorkflowCompleted handles workflow completed event.
func (ns *NOMStyleSubscriber) handleWorkflowCompleted(
	_ context.Context,
	event Event,
) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.isRunning = false

	if wa, ok := event.(WorkflowEventAccessor); ok && !wa.GetWorkflowID().IsZero() {
		fmt.Printf("✅ Workflow '%s' completed\n", wa.GetWorkflowID().String())
	}

	if err := ns.timingCache.Save(); err != nil {
		fmt.Printf("Warning: Failed to save timing cache: %v\n", err)
	}

	return nil
}

// handleWorkflowFailed handles workflow failed event.
func (ns *NOMStyleSubscriber) handleWorkflowFailed(
	_ context.Context,
	event Event,
) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.isRunning = false

	if ea, ok := event.(ErrorAccessor); ok && ea.GetError() != nil {
		fmt.Printf("❌ Workflow failed: %s\n", ea.GetError().Error())
	}

	if err := ns.timingCache.Save(); err != nil {
		fmt.Printf("Warning: Failed to save timing cache: %v\n", err)
	}

	return nil
}

// getOrCreateActivity retrieves an existing activity or creates a new one.
// Must be called while holding ns.mu lock.
func (ns *NOMStyleSubscriber) getOrCreateActivity(
	activityID ActivityID,
	activityName ActivityName,
) *ActivityDisplayState {
	activity, exists := ns.activities[activityID]
	if !exists {
		activity = NewActivityDisplayState(activityID, activityName)
		ns.activities[activityID] = activity
	}

	return activity
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

	avgDuration := ns.timingCache.GetAverage(aa.GetActivityName().String())
	if avgDuration > 0 {
		activity.SetEstimatedTime(avgDuration)
	}

	var deps []ActivityID
	if da, ok := event.(DependenciesAccessor); ok {
		deps = da.GetDependencies()
	}

	if err := ns.dependencyTree.AddActivity(
		aa.GetActivityID(),
		aa.GetActivityName().String(),
		deps,
	); err != nil {
		fmt.Printf("Warning: Failed to add activity to tree: %v\n", err)
	}

	if err := ns.dependencyTree.UpdateActivityStatus(
		aa.GetActivityID(),
		activity.Status,
		activity.Symbol,
		activity.Color,
		activity.StartTime,
		activity.EstimatedTime,
	); err != nil {
		fmt.Printf("Warning: Failed to update tree status: %v\n", err)
	}

	return nil
}

// updateActivityStateAfterExecution handles common logic for completing/failing activities.
// Records timing data and updates activity status in the dependency tree.
func (ns *NOMStyleSubscriber) updateActivityStateAfterExecution(
	activityID ActivityID,
	activityName ActivityName,
	activity *ActivityDisplayState,
	duration time.Duration,
) {
	if duration > 0 {
		if err := ns.timingCache.Record(activityName.String(), duration); err != nil {
			fmt.Printf("Warning: Failed to record timing: %v\n", err)
		}
	}

	err := ns.dependencyTree.UpdateActivityStatus(
		activityID,
		activity.Status,
		activity.Symbol,
		activity.Color,
		activity.StartTime,
		activity.EstimatedTime,
	)
	if err != nil {
		fmt.Printf("Warning: Failed to update tree status: %v\n", err)
	}
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

	var duration time.Duration
	if da, ok := event.(DurationAccessor); ok {
		duration = da.GetDuration()
	}

	ns.updateActivityStateAfterExecution(
		aa.GetActivityID(),
		aa.GetActivityName(),
		activity,
		duration,
	)

	return nil
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

	var duration time.Duration
	if da, ok := event.(DurationAccessor); ok {
		duration = da.GetDuration()
	}

	ns.updateActivityStateAfterExecution(
		aa.GetActivityID(),
		aa.GetActivityName(),
		activity,
		duration,
	)

	return nil
}
