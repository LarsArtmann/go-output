package nom

import (
	"sort"
	"time"
)

type sortKey struct {
	interest   int
	elapsed    time.Duration
	activityID ActivityID
}

func (k sortKey) less(other sortKey) bool {
	if k.interest != other.interest {
		return k.interest < other.interest
	}

	if k.elapsed != other.elapsed {
		return k.elapsed > other.elapsed
	}

	return k.activityID < other.activityID
}

// childPriority sorts a node's children by display priority using immutable
// snapshots for status/elapsed data.
func (dt *DependencyTree) childPriority(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) []*ActivityNode {
	if len(node.Children) == 0 {
		return nil
	}

	sorted := make([]*ActivityNode, 0, len(node.Children))
	sorted = append(sorted, node.Children...)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := sortKeyForNode(sorted[i], snapshots)
		kj := sortKeyForNode(sorted[j], snapshots)

		return ki.less(kj)
	})

	return sorted
}

func sortKeyForNode(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) sortKey {
	snap := lookupSnapshot(snapshots, node.ID)

	return sortKey{
		interest:   snap.Status.Interest(),
		elapsed:    snap.CurrentElapsed,
		activityID: node.ID,
	}
}
