package nom

import (
	"time"
)

// GetActivities returns a copy of all activities.
func (ns *NOMStyleSubscriber) GetActivities() map[ActivityID]*Activity {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	result := make(map[ActivityID]*Activity, len(ns.activities))
	for id, activity := range ns.activities {
		result[id] = activity.Copy()
	}

	return result
}

// GetActivity returns a specific activity.
func (ns *NOMStyleSubscriber) GetActivity(
	activityID ActivityID,
) *Activity {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	activity, exists := ns.activities[activityID]
	if !exists {
		return nil
	}

	return activity.Copy()
}

// SetActivityProgress sets a live progress message on an activity for sub-step
// visibility (e.g. "Tidying module [2/26]"). This is the direct-call equivalent
// of sending an ActivityProgress event — useful for callers that don't use the
// event dispatch path but still want progress rendered in the tree. Pass an
// empty message to clear.
func (ns *NOMStyleSubscriber) SetActivityProgress(activityID ActivityID, message string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if activity, exists := ns.activities[activityID]; exists {
		activity.Progress = message
	}
}

// SetEstimatedTime sets the predicted duration for a specific activity, sourced
// from an external timing store (e.g. a SQLite timings database). This enables
// external estimate injection without going through the nom timing cache: the
// caller loads estimates from their own store and injects them per-activity.
// The estimate is rendered as ∅Xs on pending activities in the tree.
func (ns *NOMStyleSubscriber) SetEstimatedTime(activityID ActivityID, estimated time.Duration) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	if activity, exists := ns.activities[activityID]; exists {
		activity.SetEstimatedTime(estimated)
	}
}
