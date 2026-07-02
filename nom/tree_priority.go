package nom

import (
	"sort"
	"time"
)

type sortKey struct {
	interest       int
	onCriticalPath bool
	elapsed        time.Duration
	activityID     ActivityID
}

func (k sortKey) less(other sortKey) bool {
	if k.interest != other.interest {
		return k.interest < other.interest
	}

	// Critical-path nodes sort before non-critical at the same interest level,
	// so the longest-time-chain activities surface to the top of each subtree.
	if k.onCriticalPath != other.onCriticalPath {
		return other.onCriticalPath // true sorts first (higher priority)
	}

	if k.elapsed != other.elapsed {
		return k.elapsed > other.elapsed
	}

	return k.activityID < other.activityID
}

// childPriority sorts a node's children by display priority using immutable
// snapshots for status/elapsed data. When criticalPath is non-nil, nodes on
// the critical path are boosted in sort order.
func (dt *DependencyTree) childPriority(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	criticalPath map[ActivityID]bool,
) []*ActivityNode {
	if len(node.Children) == 0 {
		return nil
	}

	sorted := make([]*ActivityNode, 0, len(node.Children))
	sorted = append(sorted, node.Children...)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := sortKeyForNode(sorted[i], snapshots, criticalPath)
		kj := sortKeyForNode(sorted[j], snapshots, criticalPath)

		return ki.less(kj)
	})

	return sorted
}

func sortKeyForNode(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	criticalPath map[ActivityID]bool,
) sortKey {
	snap := lookupSnapshot(snapshots, node.ID)

	onPath := false
	if criticalPath != nil {
		onPath = criticalPath[node.ID]
	}

	return sortKey{
		interest:       snap.Status.Interest(),
		onCriticalPath: onPath,
		elapsed:        snap.CurrentElapsed,
		activityID:     node.ID,
	}
}
