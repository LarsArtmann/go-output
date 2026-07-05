package nom

import "errors"

// ErrCycleDetected is returned by Build when the dependency graph contains a
// cycle, making a well-formed display tree impossible.
var ErrCycleDetected = errors.New("dependency cycle detected")

// Build constructs the display tree from the DAG topology and identifies root
// nodes. For each node, the deepest dependency becomes the display parent
// (matching nom's "lowermost dependency" rule: shared deps appear under the
// parent that is furthest from the root). This produces a tree where shared
// dependencies sink toward the leaves — their natural position.
//
// Returns ErrCycleDetected if the dependency graph contains a cycle.
func (dt *DependencyTree) Build() error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	dt.resetDisplayState()

	if cyclic := dt.computeDepths(); cyclic {
		return ErrCycleDetected
	}

	dt.assignDisplayParents()

	dt.loaded = true

	return nil
}

// resetDisplayState clears Parent/Children/IsRoot/Depth for every node.
// Called at the start of Build so that AddActivity changes are reflected.
func (dt *DependencyTree) resetDisplayState() {
	for _, node := range dt.nodes {
		node.Parent = nil
		node.Children = make([]*ActivityNode, 0)
		node.IsRoot = false
		node.Depth = 0
	}
}

// computeDepths assigns each node its longest-path depth from a root via
// fixpoint iteration. Nodes with no deps get depth 0; all others get
// max(dep depths) + 1. Returns true if a cycle is detected (the fixpoint
// did not converge within len(nodes)+1 iterations).
func (dt *DependencyTree) computeDepths() bool {
	maxIter := len(dt.nodes) + 1

	for i := range maxIter {
		changed := false

		for _, node := range dt.nodes {
			maxDepDepth := -1

			for _, depID := range node.Deps {
				if depNode, exists := dt.nodes[depID]; exists {
					if depNode.Depth > maxDepDepth {
						maxDepDepth = depNode.Depth
					}
				}
			}

			newDepth := 0
			if maxDepDepth >= 0 {
				newDepth = maxDepDepth + 1
			}

			if newDepth != node.Depth {
				node.Depth = newDepth
				changed = true
			}
		}

		if !changed {
			return false // converged — no cycle
		}

		// On the last iteration, if we still have changes, it's a cycle.
		if i == maxIter-1 {
			return true
		}
	}

	return false
}

// assignDisplayParents picks the deepest dep as each node's display parent
// and populates Children lists + the roots slice. Ties are broken by Deps
// order (first dep at the max depth wins — deterministic).
func (dt *DependencyTree) assignDisplayParents() {
	dt.roots = make([]*ActivityNode, 0)

	for _, node := range dt.nodes {
		if len(node.Deps) == 0 {
			node.IsRoot = true
			dt.roots = append(dt.roots, node)

			continue
		}

		var deepest *ActivityNode

		for _, depID := range node.Deps {
			if depNode, exists := dt.nodes[depID]; exists {
				if deepest == nil || depNode.Depth > deepest.Depth {
					deepest = depNode
				}
			}
		}

		if deepest != nil {
			node.Parent = deepest
			deepest.Children = append(deepest.Children, node)
		} else {
			node.IsRoot = true
			dt.roots = append(dt.roots, node)
		}
	}

	sortNodesByID(dt.roots)
}
