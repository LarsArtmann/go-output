package nom

import (
	"time"
)

// GetDependencyTree returns the dependency tree for structural access (node
// lookup, root listing). The tree nodes store ONLY IDs and tree structure —
// all mutable Activity fields are accessed via SnapshotActivities + the
// snapshot-aware render methods (RenderWithSnapshots).
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
	tree := ns.GetDependencyTree()
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
