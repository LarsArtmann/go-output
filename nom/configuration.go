package nom

import (
	"time"
)

// SetEnabled enables or disables the subscriber.
func (ns *NOMSubscriber) SetEnabled(enabled bool) {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.enabled = enabled
}

// IsEnabled returns true if the subscriber is enabled.
func (ns *NOMSubscriber) IsEnabled() bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.enabled
}

// Reset clears all state for a new.
func (ns *NOMSubscriber) Reset() {
	ns.mu.Lock()
	defer ns.mu.Unlock()

	ns.activities = make(map[ActivityID]*Activity)
	ns.dependencyTree.Clear()
	ns.counts = ActivityCounts{}
	ns.workflowID = ""
	ns.workflowName = WorkflowName("")
	ns.startTime = time.Time{}
	ns.isRunning = false
}
