package nom

import (
	"time"
)

// GetDependencyTree returns the dependency tree.
//
// IMPORTANT: The returned pointer is SHARED state. The caller is responsible
// for synchronization when mutating it. The tree itself uses an internal
// RWMutex, so individual method calls (AddActivity, UpdateActivityStatus,
// Render, etc.) are safe. However, sequences of operations that read and
// then write are not atomic with respect to concurrent subscribers.
//
// Typical use cases:
//   - Read-only rendering: safe to call Render() at any time.
//   - Mutation: prefer calling subscriber methods that route through the
//     subscriber's lock instead of mutating the tree directly.
func (ns *NOMStyleSubscriber) GetDependencyTree() *DependencyTree {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.dependencyTree
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
func (ns *NOMStyleSubscriber) GetWorkflowID() string {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.workflowID.String()
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
