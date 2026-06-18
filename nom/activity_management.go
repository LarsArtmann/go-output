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

// SetActivityState sets an activity's state (for testing purposes).
func (ns *NOMStyleSubscriber) SetActivityState(id ActivityID, activity *Activity) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.activities[id] = activity
}
