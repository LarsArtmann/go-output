package nom

import (
	"time"
)

// GetDependencyTree returns the dependency tree.
//
// IMPORTANT: The returned pointer is SHARED state. The tree's internal RWMutex
// guards the STRUCTURE (nodes/roots/children), but the nodes embed the shared
// *Activity pointer, whose fields (Status, Symbol, Color, timing) are mutated
// by event handlers under the subscriber's ns.mu — NOT the tree's lock. Reading
// those fields via Render/RenderWithWidth WITHOUT holding ns.mu therefore races
// concurrent SetRunning/SetCompleted/SetFailed calls.
//
// For read-only rendering, prefer RenderTree, which takes ns.mu.RLock for you.
// Callers that grab the raw tree must hold the subscriber's read lock
// themselves while walking it.
func (ns *NOMStyleSubscriber) GetDependencyTree() *DependencyTree {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.dependencyTree
}

// RenderTree renders the dependency tree while holding the subscriber's read
// lock, so the shared Activity fields cannot be mutated mid-render by a
// concurrent event handler (SetRunning/SetCompleted/SetFailed update Status,
// Symbol, Color and timing on the same *Activity pointer). Reading those
// fields unlocked is a data race that produces garbled/sort-inconsistent
// frames. Returns (rendered, false) when there is no tree yet.
func (ns *NOMStyleSubscriber) RenderTree(maxHeight, maxWidth int) (string, bool) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	if ns.dependencyTree == nil {
		return "", false
	}

	return ns.dependencyTree.RenderWithWidth(maxHeight, maxWidth), true
}

// WithSubscriberRLock runs fn while holding the subscriber's read lock, so the
// shared *Activity fields (Status/Symbol/Color/timing) cannot be mutated by a
// concurrent event handler (SetRunning/SetCompleted/SetFailed) during fn.
//
// Use this when you already hold a *DependencyTree reference (for example a
// cached pointer that aliases ns.dependencyTree in production) and need to
// read its Activity fields race-free. Holding this lock is what makes the read
// safe — the tree pointer you read inside fn is the caller's choice.
//
// fn MUST NOT call back into the subscriber's own RLock-taking methods
// (GetActivityCounts, RenderTree, UpdateRunningActivityElapsed, ...):
// sync.RWMutex recursive RLock against a waiting writer deadlocks. Perform any
// such calls outside fn.
func (ns *NOMStyleSubscriber) WithSubscriberRLock(fn func()) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	fn()
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
