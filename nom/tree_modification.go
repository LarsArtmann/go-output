package nom

// AddActivity registers an activity ID and its dependency edges in the tree.
// The tree stores ONLY the ID and structure (Parent/Children) — all mutable
// Activity fields are accessed via ActivitySnapshot at render time.
func (dt *DependencyTree) AddActivity(activityID ActivityID, dependencies []ActivityID) error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	node, exists := dt.nodes[activityID]
	if !exists {
		node = &ActivityNode{
			ID:          activityID,
			Children:    make([]*ActivityNode, 0),
			IsDisplayed: true,
		}
		dt.nodes[activityID] = node
	}

	for i, depID := range dependencies {
		depNode, exists := dt.nodes[depID]
		if !exists {
			depNode = newActivityNode(depID, depID.String())
			dt.nodes[depID] = depNode
		}

		if i == 0 {
			if node.Parent != nil && node.Parent != depNode {
				node.Parent.removeChild(node.ID)
			}

			node.Parent = depNode
			if !depNode.hasChild(node.ID) {
				depNode.Children = append(depNode.Children, node)
			}
		} else if !node.hasSecondaryParent(depID) {
			node.SecondaryParents = append(node.SecondaryParents, depID)
		}
	}

	dt.loaded = false

	return nil
}
