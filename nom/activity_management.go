package nom

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

// ActivityCounts holds counts of activities grouped by status. The four
// named buckets cover the core statuses; Other aggregates activities in
// registered custom statuses (e.g. "skipped", "cached") so the open status
// registry can never silently vanish from counts, totals, or percentages.
type ActivityCounts struct {
	Running   int
	Completed int
	Failed    int
	Pending   int
	Other     int
}

// applyCountsDelta adjusts the counts by decrementing the old status and
// incrementing the new status. Called on every activity state transition
// under the subscriber's write lock, so GetActivityCounts can read a
// pre-computed total in O(1) instead of scanning all activities in O(n).
// A no-op when from == to.
func applyCountsDelta(c *ActivityCounts, from, to ActivityStatus) {
	if from == to {
		return
	}

	adjustStatusCount(c, from, -1)
	adjustStatusCount(c, to, +1)
}

// adjustStatusCount applies a delta (+1 or -1) to the count bucket for the
// given status. Statuses outside the four core ones land in Other.
func adjustStatusCount(c *ActivityCounts, status ActivityStatus, delta int) {
	switch status {
	case ActivityStatusRunning:
		c.Running += delta
	case ActivityStatusCompleted:
		c.Completed += delta
	case ActivityStatusFailed:
		c.Failed += delta
	case ActivityStatusPending:
		c.Pending += delta
	default:
		c.Other += delta
	}
}

// Total returns the sum of all activity counts, including custom statuses.
func (c ActivityCounts) Total() int {
	return c.Running + c.Completed + c.Failed + c.Pending + c.Other
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
		parts = append(parts, string(SymbolPending)+strconv.Itoa(c.Pending))
	}

	if c.Other > 0 {
		parts = append(parts, string(SymbolOther)+strconv.Itoa(c.Other))
	}

	return strings.Join(parts, " ")
}

// SummaryColored renders the counts with lipgloss Foreground styling per status.
// Each non-zero count gets its symbol and number colored according to the
// provided SemanticColors palette. Zero counts are omitted, matching Summary().
// Returns "" when all counts are zero.
func (c ActivityCounts) SummaryColored(colors SemanticColors) string {
	var parts []string

	style := lipgloss.NewStyle()

	if c.Running > 0 {
		parts = append(parts, style.Foreground(colors.Running).Render(string(SymbolRunning)+strconv.Itoa(c.Running)))
	}

	if c.Completed > 0 {
		parts = append(
			parts,
			style.Foreground(colors.Completed).Render(string(SymbolCompleted)+strconv.Itoa(c.Completed)),
		)
	}

	if c.Failed > 0 {
		parts = append(parts, style.Foreground(colors.Failed).Render(string(SymbolFailed)+strconv.Itoa(c.Failed)))
	}

	if c.Pending > 0 {
		parts = append(parts, style.Foreground(colors.Pending).Render(string(SymbolPending)+strconv.Itoa(c.Pending)))
	}

	if c.Other > 0 {
		parts = append(parts, style.Foreground(colors.Pending).Render(string(SymbolOther)+strconv.Itoa(c.Other)))
	}

	return strings.Join(parts, "  ")
}

// GetActivityCounts returns counts of activities by status.
// O(1) — reads a pre-computed cache maintained incrementally on every state
// transition (applyDelta in the event handlers). Previously this scanned all
// activities every frame; now it returns the cached aggregate directly.
func (ns *NOMSubscriber) GetActivityCounts() ActivityCounts {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.counts
}

// SetActivityState sets an activity's state (for testing purposes).
// Maintains the count cache: if replacing an existing activity, the old
// status is decremented before the new one is counted. A nil activity is
// ignored — matching the defensive contract of the renderer entry points.
func (ns *NOMSubscriber) SetActivityState(id ActivityID, activity *Activity) {
	if activity == nil {
		return
	}

	ns.mu.Lock()
	defer ns.mu.Unlock()

	if old, exists := ns.activities[id]; exists {
		adjustStatusCount(&ns.counts, old.Status, -1)
	}

	adjustStatusCount(&ns.counts, activity.Status, +1)
	ns.activities[id] = activity
}
