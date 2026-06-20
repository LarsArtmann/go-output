package nom

import (
	"image/color"
	"time"
)

// ActivitySnapshot is an immutable value copy of the mutable fields of an
// Activity, captured under the subscriber's read lock. The dependency tree
// renderer consumes snapshots instead of reading the shared *Activity pointer
// directly, eliminating the data race between event handlers (SetRunning etc.)
// and rendering.
type ActivitySnapshot struct {
	Label          string
	Status         ActivityStatus
	Symbol         Symbol
	Color          color.Color
	CurrentElapsed time.Duration
	EstimatedTime  time.Duration
}

// SnapshotActivities returns immutable copies of every registered activity's
// mutable fields. Thread-safe: acquires the subscriber's read lock internally,
// copies the fields, then releases. The returned map is safe to read without
// any further locking. This is the race-free way to obtain activity state for
// rendering — the tree walk reads only the immutable snapshot data.
func (ns *NOMStyleSubscriber) SnapshotActivities() map[ActivityID]ActivitySnapshot {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	snapshots := make(map[ActivityID]ActivitySnapshot, len(ns.activities))

	for id, activity := range ns.activities {
		snapshots[id] = ActivitySnapshot{
			Label:          activity.Label.Get(),
			Status:         activity.Status,
			Symbol:         activity.Symbol,
			Color:          activity.Color,
			CurrentElapsed: activity.CurrentElapsed,
			EstimatedTime:  activity.EstimatedTime,
		}
	}

	return snapshots
}
