package nom
import (
	"sort"
)
// GetDisplayActivities returns the list of activity IDs in display order.
func (dt *DependencyTree) GetDisplayActivities() []ActivityID {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return dt.order
}
// GetNode returns a node by activity ID.
func (dt *DependencyTree) GetNode(activityID ActivityID) *TreeNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return dt.nodes[activityID]
}
// GetRootNodes returns all root nodes, building the tree first if needed.
// Uses double-checked locking to avoid building under a read lock (deadlock prevention).
func (dt *DependencyTree) GetRootNodes() []*TreeNode {
	dt.mu.RLock()
	loaded := dt.loaded
	dt.mu.RUnlock()
	if !loaded {
		_ = dt.Build()
	}
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	return dt.roots
}
// EnsureBuild guarantees the tree is built at most once between AddActivity calls.
// AddActivity resets buildOnce so the next GetRootNodes/EnsureBuild will rebuild.
func (dt *DependencyTree) EnsureBuild() {
	dt.mu.RLock()
	loaded := dt.loaded
	dt.mu.RUnlock()
	if !loaded {
		_ = dt.Build()
	}
}
// SnapshotRoots returns a snapshot of root nodes safe for concurrent traversal.
func (dt *DependencyTree) SnapshotRoots() []*TreeNode {
	roots := dt.GetRootNodes()
	snapshot := make([]*TreeNode, len(roots))
	copy(snapshot, roots)
	return snapshot
}
// FindNodesByStatus returns all nodes matching the given status.
func (dt *DependencyTree) FindNodesByStatus(status ActivityStatus) []*TreeNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	var result []*TreeNode
	for _, node := range dt.nodes {
		if node.Status == status {
			result = append(result, node)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return string(result[i].ActivityID) < string(result[j].ActivityID)
	})
	return result
}
