package nom

import (
	"time"
)

// GetActivityCounts returns counts of activities by status.
func (ns *NOMStyleSubscriber) GetActivityCounts() (running, completed, failed, pending int) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	for _, activity := range ns.activities {
		switch activity.Status {
		case ActivityStatusRunning:
			running++
		case ActivityStatusCompleted:
			completed++
		case ActivityStatusFailed:
			failed++
		case ActivityStatusPending:
			pending++
		case ActivityStatusPaused:
			pending++ // Paused activities counted as pending
		}
	}

	return running, completed, failed, pending
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
// timing) from ActivityDisplayState to its corresponding TreeNode.
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

		// Copy display state from activity to tree node.
		node.Status = activity.Status
		node.Symbol = activity.Symbol
		node.Color = activity.Color
		// Copy timing information from activity to tree node.
		node.StartTime = activity.StartTime
		node.CurrentElapsed = activity.CurrentElapsed
		node.EstimatedTime = activity.EstimatedTime
	}
}
