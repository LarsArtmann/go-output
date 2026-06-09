package nom
import (
	"time"
)
// GetDependencyTree returns dependency tree.
func (ns *NOMStyleSubscriber) GetDependencyTree() *DependencyTree {
	return ns.dependencyTree
}
// GetTimingCache returns timing cache.
func (ns *NOMStyleSubscriber) GetTimingCache() *TimingCache {
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
