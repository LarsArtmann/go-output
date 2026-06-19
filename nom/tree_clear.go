package nom

// Clear removes all nodes from the tree.
func (dt *DependencyTree) Clear() {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.nodes = make(map[ActivityID]*ActivityNode)
	dt.roots = make([]*ActivityNode, 0)
	dt.loaded = false
}
