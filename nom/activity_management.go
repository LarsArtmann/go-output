package nom

import (
	"time"
)

// ActivityCounts holds counts of activities grouped by status.
type ActivityCounts struct {
	Running   int
	Completed int
	Failed    int
	Pending   int
}

// Total returns the sum of all activity counts.
func (c ActivityCounts) Total() int {
	return c.Running + c.Completed + c.Failed + c.Pending
}

// GetActivityCounts returns counts of activities by status.
func (ns *NOMStyleSubscriber) GetActivityCounts() ActivityCounts {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var c ActivityCounts

	for _, activity := range ns.activities {
		switch activity.Status {
		case ActivityStatusRunning:
			c.Running++
		case ActivityStatusCompleted:
			c.Completed++
		case ActivityStatusFailed:
			c.Failed++
		case ActivityStatusPending:
			c.Pending++
		case ActivityStatusPaused:
			c.Pending++ // Paused activities counted as pending
		}
	}

	return c
}

// UpdateRunningActivityElapsed updates elapsed time for all currently running activities.
// This should be called periodically (e.g., on each tick) to ensure timing displays are current.
func (ns *NOMStyleSubscriber) UpdateRunningActivityElapsed() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	now := time.Now()

	for _, activity := range ns.activities {
		if activity.Status == ActivityStatusRunning && !activity.StartTime.IsZero() {
			activity.CurrentElapsed = now.Sub(activity.StartTime)
		}
	}
}

// SetActivityState sets an activity's display state (for testing purposes).
func (ns *NOMStyleSubscriber) SetActivityState(activity *ActivityDisplayState) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.activities[activity.ActivityID] = activity
}

// SyncActivityTimingToTree synchronizes display state (status, symbol, color,
// timing) from ActivityDisplayState to its corresponding ActivityNode.
//
// Both structures duplicate the display fields for their respective consumers
// (TUI status list vs. dependency tree rendering). This helper keeps them
// aligned; callers should invoke it after any state mutation and before
// rendering. The tree's own mutex is acquired internally.
func (ns *NOMStyleSubscriber) SyncActivityTimingToTree() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	for activityID, activity := range ns.activities {
		node := ns.dependencyTree.GetNode(activityID)
		if node == nil {
			// Tree node doesn't exist yet - skip
			continue
		}

		// Sync fields from ActivityDisplayState to ActivityNode's embedded Activity.
		syncActivityToNode(node, activity)
	}
}

// syncActivityToNode copies display fields from ActivityDisplayState to an
// ActivityNode's embedded Activity, including Shape/Style for diagram export.
// This is the bridge sync that will be eliminated once the subscriber stores
// Activities directly.
func syncActivityToNode(node *ActivityNode, ads *ActivityDisplayState) {
	node.Status = ads.Status
	node.Symbol = ads.Symbol
	node.Color = ads.Color
	node.StartTime = ads.StartTime
	node.EstimatedTime = ads.EstimatedTime
	node.CurrentElapsed = ads.CurrentElapsed
	node.EndTime = ads.EndTime
	node.OperationType = ads.OperationType
	node.Err = ads.Error
	node.Shape = ads.Status.NodeShape()
	node.Style = ads.Status.GraphStyle()
}
