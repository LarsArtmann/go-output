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
	// Host is an optional annotation naming where the activity runs (e.g. a
	// build machine). Rendered as a dim right-aligned tag when non-empty;
	// dormant otherwise. Mirrors NOM's host column.
	Host string
	// Download is an optional byte-progress indicator. Rendered as a compact
	// progress bar when Total > 0; dormant otherwise. Mirrors NOM's per-
	// activity download bars.
	Download DownloadProgress
}

// DownloadProgress is an optional byte-level progress indicator for an activity
// (e.g. a dependency download). Zero-value means "no download in progress".
type DownloadProgress struct {
	Downloaded int64 // bytes transferred
	Total      int64 // total bytes (0 = unknown/streaming)
}

// HasDownload reports whether download progress should be rendered.
func (d DownloadProgress) HasDownload() bool { return d.Downloaded > 0 || d.Total > 0 }

// Fraction returns the completion fraction in [0,1], or 0 when total is unknown.
func (d DownloadProgress) Fraction() float64 {
	if d.Total <= 0 {
		return 0
	}

	if d.Downloaded >= d.Total {
		return 1
	}

	return float64(d.Downloaded) / float64(d.Total)
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
			Host:           activity.Host,
			Download:       activity.Download,
		}
	}

	return snapshots
}
