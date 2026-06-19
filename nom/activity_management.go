package nom

import (
	"strconv"
	"strings"
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

// CompletionPercent returns the percentage of activities that have reached a
// terminal state (completed or failed). Returns 0 when there are no activities.
func (c ActivityCounts) CompletionPercent() int {
	total := c.Total()
	if total == 0 {
		return 0
	}

	return (c.Completed + c.Failed) * 100 / total
}

// Summary renders the counts as a NOM-style status string using typed symbols,
// e.g. "⏵1 ✔2 ⚠3 ⏸4". Empty categories are omitted. Returns "" when all counts
// are zero. This is the single source of truth for count formatting — both
// InlineRenderer.renderSummary and tui.buildActivityCountsSummary delegate here.
func (c ActivityCounts) Summary() string {
	var parts []string

	if c.Running > 0 {
		parts = append(parts, string(SymbolRunning)+strconv.Itoa(c.Running))
	}

	if c.Completed > 0 {
		parts = append(parts, string(SymbolCompleted)+strconv.Itoa(c.Completed))
	}

	if c.Failed > 0 {
		parts = append(parts, string(SymbolFailed)+strconv.Itoa(c.Failed))
	}

	if c.Pending > 0 {
		parts = append(parts, string(SymbolPaused)+strconv.Itoa(c.Pending))
	}

	return strings.Join(parts, " ")
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
