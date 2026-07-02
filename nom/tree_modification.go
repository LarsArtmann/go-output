package nom

import "slices"

// AddActivity registers an activity ID and its dependency edges in the DAG.
// All dependencies are recorded as real edges in node.Deps — none are
// demoted to "secondary" status. The display tree structure (Parent /
// Children) is assigned later by Build(), which picks the deepest dep as
// the display parent for each node (matching nom's strategy).
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

	// Replace deps only when the caller provides new ones. An empty slice
	// means "not declaring deps on this event" (e.g. ActivityStarted after
	// ActivityRegistered), not "this activity has no deps". This preserves
	// backward compatibility with the event flow.
	if len(dependencies) > 0 {
		node.Deps = make([]ActivityID, 0, len(dependencies))

		for _, depID := range dependencies {
			if _, depExists := dt.nodes[depID]; !depExists {
				dt.nodes[depID] = newActivityNode(depID, depID.String())
			}

			if !slices.Contains(node.Deps, depID) {
				node.Deps = append(node.Deps, depID)
			}
		}
	}

	dt.loaded = false

	return nil
}
