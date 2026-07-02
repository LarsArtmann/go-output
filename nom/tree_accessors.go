package nom

// SetRenderMode selects tree or layered display for this dependency tree.
func (dt *DependencyTree) SetRenderMode(mode RenderMode) {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.renderMode = mode
}

// RenderMode returns the current display mode.
func (dt *DependencyTree) RenderMode() RenderMode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	return dt.renderMode
}

// DAGNode is a read-only snapshot of a node in the dependency tree, suitable
// for external consumers (e.g. DOT export) that should not mutate tree state.
type DAGNode struct {
	ID     ActivityID
	Deps   []ActivityID
	Depth  int
	IsRoot bool
}

// AllNodes returns a snapshot of every node in the dependency tree as read-only
// DAGNode values. The slice is safe to iterate without holding the tree lock.
func (dt *DependencyTree) AllNodes() []DAGNode {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	nodes := make([]DAGNode, 0, len(dt.nodes))

	for _, node := range dt.nodes {
		deps := make([]ActivityID, len(node.Deps))
		copy(deps, node.Deps)

		nodes = append(nodes, DAGNode{
			ID:     node.ID,
			Deps:   deps,
			Depth:  node.Depth,
			IsRoot: node.IsRoot,
		})
	}

	return nodes
}

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
	snapshot := make([]*ActivityNode, 0, len(roots))
	snapshot = append(snapshot, roots...)

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
