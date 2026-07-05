package nom

import (
	"fmt"
	"time"
)

// ParallelismStats reports how many activities are currently running and how
// many more could start immediately (all dependencies satisfied).
type ParallelismStats struct {
	Running  int
	Possible int
}

// String returns a compact display form like "parallel: 3/4 possible".
func (ps ParallelismStats) String() string {
	return fmt.Sprintf("parallel: %d/%d possible", ps.Running, ps.Possible)
}

// GetDependencyTree returns the dependency tree for structural access (node
// lookup, root listing). The tree nodes store ONLY IDs and tree structure —
// all mutable Activity fields are accessed via SnapshotActivities + the
// snapshot-aware render methods (RenderWithSnapshots).
//
// Deprecated: Use DependencyTree() instead — identical behavior, shorter name.
// This method will be removed in v2.
//
//nolint:staticcheck // kept for backward compatibility, remove in v2
func (ns *NOMStyleSubscriber) GetDependencyTree() *DependencyTree {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.dependencyTree
}

// RenderSnapshot takes a snapshot of all activity fields (thread-safe), then
// renders the tree from that immutable data. No lock is held during the tree
// walk, so event handlers can proceed concurrently. Returns ("", false) when
// no tree exists yet.
func (ns *NOMStyleSubscriber) RenderSnapshot(maxHeight, maxWidth int) (string, bool) {
	tree := ns.DependencyTree()
	if tree == nil {
		return "", false
	}

	snapshots := ns.SnapshotActivities()

	return tree.RenderWithSnapshots(snapshots, maxHeight, maxWidth), true
}

// GetTimingCache returns timing cache.
func (ns *NOMStyleSubscriber) GetTimingCache() *TimingCache {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.timingCache
}

// Flush persists all pending timing-cache writes to disk and returns the last
// save error (if any). Call this at subscriber shutdown to ensure no timing
// data is lost.
func (ns *NOMStyleSubscriber) Flush() error {
	ns.mu.RLock()
	cache := ns.timingCache
	ns.mu.RUnlock()

	return cache.Flush()
}

// IsWorkflowRunning returns true if a workflow is currently running.
func (ns *NOMStyleSubscriber) IsWorkflowRunning() bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.isRunning
}

// GetWorkflowID returns current workflow ID.
func (ns *NOMStyleSubscriber) GetWorkflowID() WorkflowID {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.workflowID
}

// GetWorkflowName returns current workflow name.
func (ns *NOMStyleSubscriber) GetWorkflowName() string {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.workflowName.String()
}

// GetStartTime returns workflow start time.
func (ns *NOMStyleSubscriber) GetStartTime() time.Time {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.startTime
}

// EstimatedTotalRemaining returns the projected remaining time for all
// unfinished work, computed from per-activity estimates:
//   - Pending activities contribute their full EstimatedTime.
//   - Running activities contribute max(0, EstimatedTime - elapsed).
//   - Completed/failed activities contribute nothing.
//
// This is the subscriber-owned estimate source. Both renderers can consume it:
// the TUI uses it directly for its "~Xm left" summary, and the inline renderer's
// SetEstimatedRemainingFunc callback can delegate to it when the caller has no
// external estimator. Returns 0 when no unfinished activity has an estimate.
func (ns *NOMStyleSubscriber) EstimatedTotalRemaining() time.Duration {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	now := time.Now()

	var total time.Duration

	for _, a := range ns.activities {
		if a.Status == ActivityStatusCompleted || a.Status == ActivityStatusFailed {
			continue
		}

		if a.EstimatedTime <= 0 {
			continue
		}

		remaining := a.EstimatedTime

		if a.Status == ActivityStatusRunning {
			elapsed := a.elapsedAt(now)
			if elapsed >= remaining {
				continue
			}

			remaining -= elapsed
		}

		total += remaining
	}

	return total
}

// ParallelismStats returns the number of currently running activities and the
// number of pending activities whose dependencies are all complete (i.e., ready
// to start). The tree is consulted for dependency edges; the subscriber lock is
// held first, then the tree lock, matching the lock-order used by Edges().
func (ns *NOMStyleSubscriber) ParallelismStats() ParallelismStats {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var stats ParallelismStats

	tree := ns.dependencyTree
	if tree == nil {
		return stats
	}

	tree.mu.RLock()
	defer tree.mu.RUnlock()

	for _, activity := range ns.activities {
		switch activity.Status { //nolint:exhaustive // only running/pending affect parallelism
		case ActivityStatusRunning:
			stats.Running++
		case ActivityStatusPending:
			if ns.canStartImmediatelyLocked(ActivityID(activity.ID.String()), tree) {
				stats.Possible++
			}
		}
	}

	return stats
}

func (ns *NOMStyleSubscriber) canStartImmediatelyLocked(
	id ActivityID,
	tree *DependencyTree,
) bool {
	node, ok := tree.nodes[id]
	if !ok || len(node.Deps) == 0 {
		return true
	}

	for _, depID := range node.Deps {
		dep, ok := ns.activities[depID]
		if !ok || dep.Status != ActivityStatusCompleted {
			return false
		}
	}

	return true
}
