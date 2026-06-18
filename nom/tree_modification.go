package nom

import (
	"sync"

	"github.com/larsartmann/go-output"
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
		node = newActivityNode(activityID, activityName)
		dt.nodes[activityID] = node
	}

	node.Label = output.NewBrandedID[output.GraphNodeLabelBrand](activityName)
	// Add dependency relationships — use first dependency as primary parent for tree structure
	for i, depID := range dependencies {
		depNode, exists := dt.nodes[depID]
		if !exists {
			depNode = newActivityNode(depID, depID.String())
			dt.nodes[depID] = depNode
		}
		// Primary parent is the first dependency; others are secondary parents (display only)
		if i == 0 {
			// Remove from old parent's Children to prevent phantom edges on re-parenting
			if node.Parent != nil && node.Parent != depNode {
				node.Parent.removeChild(ActivityID(node.ID.Get()))
			}

			node.Parent = depNode
			if !depNode.hasChild(ActivityID(node.ID.Get())) {
				depNode.Children = append(depNode.Children, node)
			}
		} else if !node.hasSecondaryParent(depID) {
			node.SecondaryParents = append(node.SecondaryParents, depID)
		}
	}

	dt.loaded = false // Need to rebuild
	dt.buildOnce = sync.Once{}

	return nil
}
