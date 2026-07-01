package nom

import (
	"context"
	"time"
)

// OnEvent dispatches a typed lifecycle event to the matching handler via an
// exhaustive Go type switch. Because Event is sealed (unexported isEvent
// marker), unhandled event types are a compile error, not a silent runtime
// no-op. This replaces the old string-based GetEventType() dispatch.
func (ns *NOMStyleSubscriber) OnEvent(_ context.Context, event Event) error {
	switch e := event.(type) {
	case WorkflowStarted:
		return ns.handleWorkflowStarted(e)
	case WorkflowCompleted:
		return ns.handleWorkflowFinished()
	case WorkflowFailed:
		return ns.handleWorkflowFinished()
	case ActivityStarted:
		return ns.handleActivityStarted(e)
	case ActivityRegistered:
		return ns.handleActivityRegistered(e)
	case ActivityCompleted:
		return ns.handleActivityCompleted(e)
	case ActivityFailed:
		return ns.handleActivityFailed(e)
	case ActivityProgress:
		return ns.handleActivityProgress(e)
	case ActivityRetrying:
		return ns.handleActivityRetrying(e)
	default:
		return nil
	}
}

// handleWorkflowStarted records the workflow identity and loads the timing
// cache. Pre-registered activities/tree are preserved on a fresh start so
// callers can register phases/steps before workflow.started.
func (ns *NOMStyleSubscriber) handleWorkflowStarted(e WorkflowStarted) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.workflowID = e.ID
	ns.startTime = time.Now()
	ns.isRunning = true
	ns.workflowName = e.Name

	if ns.activities == nil {
		ns.activities = make(map[ActivityID]*Activity)
	}

	if ns.dependencyTree == nil {
		ns.dependencyTree = NewDependencyTree()
	}

	return ns.timingCache.EnsureLoaded()
}

// handleWorkflowFinished marks the workflow as not running and persists the
// timing cache. Shared by completed and failed.
func (ns *NOMStyleSubscriber) handleWorkflowFinished() error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.isRunning = false

	return ns.timingCache.Save()
}

// getOrCreateActivity retrieves an existing activity or creates a new one.
// The kind only applies on first creation; an existing activity keeps its
// original kind (set at construction, never changes).
// Must be called while holding ns.mu lock.
func (ns *NOMStyleSubscriber) getOrCreateActivity(
	activityID ActivityID,
	activityName ActivityName,
	kind ActivityKind,
) *Activity {
	activity, exists := ns.activities[activityID]
	if !exists {
		if kind.IsPhase() {
			activity = NewPhase(string(activityID), activityName.String())
		} else {
			activity = NewActivity(string(activityID), activityName.String())
		}

		ns.activities[activityID] = activity
		// New activities start Pending; maintain the count cache.
		ns.counts.Pending++
	}

	return activity
}

// handleActivityStarted transitions the activity to running, applies optional
// host/download annotations, and records the dependency edges in the tree.
func (ns *NOMStyleSubscriber) handleActivityStarted(e ActivityStarted) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(e.ID, e.Name, e.Kind)
	applyCountsDelta(&ns.counts, activity.Status, ActivityStatusRunning)
	activity.SetRunning()
	activity.Host = e.Host
	activity.Download = e.Download

	medianDuration := ns.timingCache.GetMedian(e.Name.String())
	if medianDuration > 0 {
		activity.SetEstimatedTime(medianDuration)
	}

	return ns.dependencyTree.AddActivity(e.ID, e.Deps)
}

// handleActivityRegistered pre-creates the activity in the tree as pending,
// without transitioning it to running. Used for declaring structure before
// work starts.
func (ns *NOMStyleSubscriber) handleActivityRegistered(e ActivityRegistered) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.getOrCreateActivity(e.ID, e.Name, e.Kind)

	return ns.dependencyTree.AddActivity(e.ID, e.Deps)
}

// handleActivityCompleted transitions the activity to completed and records
// the observed duration in the timing cache.
func (ns *NOMStyleSubscriber) handleActivityCompleted(e ActivityCompleted) error {
	ns.transitionTask(e.ID, e.Name, ActivityStatusCompleted, func(a *Activity) {
		a.SetCompleted()
	})

	return ns.recordDuration(e.Name, e.Duration)
}

// handleActivityFailed transitions the activity to failed and records the
// observed duration in the timing cache.
func (ns *NOMStyleSubscriber) handleActivityFailed(e ActivityFailed) error {
	ns.transitionTask(e.ID, e.Name, ActivityStatusFailed, func(a *Activity) {
		a.SetFailed(e.Err)
	})

	return ns.recordDuration(e.Name, e.Duration)
}

// transitionTask looks up (or creates) the named task activity, adjusts the
// status counts, and invokes apply to apply status-specific fields (e.g.
// error for failure). The entire sequence runs under ns.mu so callers don't
// need to hold any lock — apply sees the activity after SetCompleted/SetFailed
// would have been called, eliminating a race where the activity's Status
// field was read by SnapshotActivities mid-transition.
func (ns *NOMStyleSubscriber) transitionTask(
	id ActivityID, name ActivityName, target ActivityStatus, apply func(*Activity),
) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(id, name, ActivityKindTask)
	applyCountsDelta(&ns.counts, activity.Status, target)
	apply(activity)
}

// recordDuration stores the observed duration in the timing cache if positive.
func (ns *NOMStyleSubscriber) recordDuration(name ActivityName, duration time.Duration) error {
	if duration > 0 {
		return ns.timingCache.Record(name.String(), duration)
	}

	return nil
}

// handleActivityProgress sets a live progress message on a running activity.
// This enables sub-step visibility: a single activity like "go-mod-tidy" can
// report "Tidying module [2/26]: modules/gitignore" while iterating. An empty
// Message clears any prior progress. The activity is created if it doesn't
// exist yet (progress events may arrive before started in some orderings).
func (ns *NOMStyleSubscriber) handleActivityProgress(e ActivityProgress) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(e.ID, e.Name, ActivityKindTask)
	activity.Progress = e.Message

	return nil
}

// handleActivityRetrying transitions a failed activity back to running and
// increments the retry count. The attempt number is rendered as a ⟳ suffix.
// The counts cache is updated: the activity moves from Failed back to Running.
func (ns *NOMStyleSubscriber) handleActivityRetrying(e ActivityRetrying) error {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	activity := ns.getOrCreateActivity(e.ID, e.Name, ActivityKindTask)
	applyCountsDelta(&ns.counts, activity.Status, ActivityStatusRunning)
	activity.SetRunning()

	if e.Attempt > activity.RetryCount {
		activity.RetryCount = e.Attempt
	} else {
		activity.RetryCount++
	}

	activity.RetryReason = e.Reason

	return nil
}
