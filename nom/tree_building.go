package nom

// Build constructs the tree structure and identifies root nodes.
func (dt *DependencyTree) Build() error {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	// Clear existing roots
	dt.roots = make([]*ActivityNode, 0)
	// Find root nodes (no parent)
	for _, node := range dt.nodes {
		if node.Parent == nil {
			node.IsRoot = true
			node.Depth = 0
			dt.roots = append(dt.roots, node)
		}
	}
	// Sort root nodes by activity ID for consistent display
	sortNodesByID(dt.roots)
	// Calculate depths for all nodes
	for _, root := range dt.roots {
		dt.calculateDepth(root, 0)
	}

	dt.loaded = true

	return nil
}

// calculateDepth recursively calculates the depth of each node.
func (dt *DependencyTree) calculateDepth(node *ActivityNode, depth int) {
	node.Depth = depth
	for _, child := range node.Children {
		dt.calculateDepth(child, depth+1)
	}
}
