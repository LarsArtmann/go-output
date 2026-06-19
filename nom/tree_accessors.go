package nom

// GetNode returns a node by activity ID.
func (dt *DependencyTree) GetNode(activityID ActivityID) *ActivityNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	return dt.nodes[activityID]
}

// GetRootNodes returns all root nodes, building the tree first if needed.
// Uses double-checked locking to avoid building under a read lock (deadlock prevention).
func (dt *DependencyTree) GetRootNodes() []*ActivityNode {
	dt.mu.RLock()
	loaded := dt.loaded
	dt.mu.RUnlock()

	if !loaded {
		if err := dt.Build(); err != nil {
			return nil
		}
	}

	dt.mu.RLock()
	defer dt.mu.RUnlock()

	return dt.roots
}

// EnsureBuild guarantees the tree is built at least once between AddActivity calls.
// AddActivity resets loaded so the next GetRootNodes/EnsureBuild will rebuild.
//
// Deprecated: EnsureBuild is exported for cross-module test use only. Production
// code should use GetRootNodes() or VisibleNodes(), which build implicitly.
func (dt *DependencyTree) EnsureBuild() {
	dt.mu.RLock()
	loaded := dt.loaded
	dt.mu.RUnlock()

	if !loaded {
		_ = dt.Build() // Best-effort; Build currently never errors.
	}
}

func (dt *DependencyTree) snapshotRoots() []*ActivityNode {
	roots := dt.GetRootNodes()
	snapshot := make([]*ActivityNode, len(roots))
	copy(snapshot, roots)

	return snapshot
}

func (dt *DependencyTree) findNodesByStatus(status ActivityStatus) []*ActivityNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	var result []*ActivityNode

	for _, node := range dt.nodes {
		if node.Status == status {
			result = append(result, node)
		}
	}

	sortNodesByID(result)

	return result
}
