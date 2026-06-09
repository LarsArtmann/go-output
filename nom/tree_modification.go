package nom

import (
	"sync"
)

// AddActivity adds an activity to the tree with its dependencies.
func (dt *DependencyTree) AddActivity(
	activityID ActivityID,
	activityName string,
	dependencies []ActivityID,
) error {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	// Create or update the node
	node, exists := dt.nodes[activityID]
	if !exists {
		node = newTreeNode(activityID, activityName)
		dt.nodes[activityID] = node
	}

	node.ActivityName = activityName
	// Add dependency relationships — use first dependency as primary parent for tree structure
	for i, depID := range dependencies {
		depNode, exists := dt.nodes[depID]
		if !exists {
			depNode = newTreeNode(depID, depID.String())
			dt.nodes[depID] = depNode
		}
		// Primary parent is the first dependency; others are secondary parents (display only)
		if i == 0 {
			node.Parent = depNode
			if !depNode.hasChild(node.ActivityID) {
				depNode.Children = append(depNode.Children, node)
			}
		} else {
			node.SecondaryParents = append(node.SecondaryParents, depID)
		}
	}

	dt.loaded = false // Need to rebuild
	dt.buildOnce = sync.Once{}

	return nil
}
