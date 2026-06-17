package nom

import (
	"sort"
	"time"
)

// sortKey captures the display priority of a node. Lower values sort first.
type sortKey struct {
	interest int
	// For equal interest, longer elapsed time sorts first so long-running work
	// is more visible than work that just started.
	elapsed time.Duration
	// Stable tie-breaker for deterministic output.
	activityID ActivityID
}

// less reports whether this key should sort before other.
func (k sortKey) less(other sortKey) bool {
	if k.interest != other.interest {
		return k.interest < other.interest
	}

	if k.elapsed != other.elapsed {
		return k.elapsed > other.elapsed
	}

	return k.activityID < other.activityID
}

// childPriority sorts a node's children by display priority without mutating the original slice.
func (dt *DependencyTree) childPriority(node *TreeNode) []*TreeNode {
	if len(node.Children) == 0 {
		return nil
	}

	sorted := make([]*TreeNode, len(node.Children))
	copy(sorted, node.Children)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sortKeyForNode(sorted[i]).less(sortKeyForNode(sorted[j]))
	})

	return sorted
}

// sortKeyForNode returns the display sort key for a node.
func sortKeyForNode(node *TreeNode) sortKey {
	return sortKey{
		interest:   node.Status.Interest(),
		elapsed:    node.CurrentElapsed,
		activityID: node.ActivityID,
	}
}
