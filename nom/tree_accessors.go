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

func (dt *DependencyTree) snapshotRoots() []*ActivityNode {
	roots := dt.GetRootNodes()
	snapshot := make([]*ActivityNode, len(roots))
	copy(snapshot, roots)

	return snapshot
}

func (dt *DependencyTree) findNodesByStatus(
	status ActivityStatus,
	snapshots map[ActivityID]ActivitySnapshot,
) []*ActivityNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	var result []*ActivityNode

	for _, node := range dt.nodes {
		snap := lookupSnapshot(snapshots, node.ID)
		if snap.Status == status {
			result = append(result, node)
		}
	}

	sortNodesByID(result)

	return result
}
