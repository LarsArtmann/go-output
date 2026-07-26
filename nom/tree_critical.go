package nom

import (
	"maps"
	"slices"
	"time"
)

// computeCriticalPath returns the set of activity IDs that lie on the longest
// estimated-time path through the DAG. It walks forward from roots to compute
// the longest path ending at each node, then back-tracks from the node(s) with
// the maximum total duration.
//
// The caller must hold dt.mu.RLock (or Lock). This function is O(n · e) in the
// worst case via fixpoint iteration, but typical DAGs converge in 2-3 passes.
//
//nolint:gocognit,cyclop // fixpoint DAG longest-path algorithm is inherently sequential
func (dt *DependencyTree) computeCriticalPath(
	snapshots map[ActivityID]ActivitySnapshot,
) map[ActivityID]bool {
	longestTo := make(map[ActivityID]time.Duration, len(dt.nodes))

	// Fixpoint: longestTo[node] = nodeWeight(node) + max(longestTo[deps]).
	maxIter := len(dt.nodes) + 1

	for range maxIter {
		changed := false

		for _, node := range dt.nodes {
			maxDep := time.Duration(0)

			for _, depID := range node.Deps {
				if depDur, ok := longestTo[depID]; ok && depDur > maxDep {
					maxDep = depDur
				}
			}

			newTotal := maxDep + nodeWeight(node.ID, snapshots)
			if longestTo[node.ID] != newTotal {
				longestTo[node.ID] = newTotal
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	// Find the maximum total duration.
	maxTotal := slices.Max(slices.Collect(maps.Values(longestTo)))

	if maxTotal <= 0 {
		return nil
	}

	// Back-track from every node that achieves the maximum.
	onPath := make(map[ActivityID]bool, len(dt.nodes))

	var frontier []ActivityID

	for id, d := range longestTo {
		if d == maxTotal {
			frontier = append(frontier, id)
		}
	}

	// Sort for deterministic traversal in tests.
	slices.Sort(frontier)

	for len(frontier) > 0 {
		id := frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]

		if onPath[id] {
			continue
		}

		onPath[id] = true

		node := dt.nodes[id]
		nodeTotal := longestTo[id]
		weight := nodeWeight(node.ID, snapshots)

		for _, depID := range node.Deps {
			if onPath[depID] {
				continue
			}

			depTotal := longestTo[depID]
			if depTotal > 0 && depTotal+weight == nodeTotal {
				frontier = append(frontier, depID)
			}
		}
	}

	return onPath
}

// nodeWeight returns the duration weight used for critical-path computation.
// Terminal activities use their actual elapsed time; running activities use
// the larger of elapsed or estimated time; pending activities use their
// estimated time (0 if unknown).
func nodeWeight(id ActivityID, snapshots map[ActivityID]ActivitySnapshot) time.Duration {
	snap := lookupSnapshot(snapshots, id)

	switch snap.Status {
	case ActivityStatusPending:
		return snap.EstimatedTime
	case ActivityStatusCompleted, ActivityStatusFailed:
		return snap.CurrentElapsed
	case ActivityStatusRunning:
		elapsed := snap.CurrentElapsed
		if snap.EstimatedTime > elapsed {
			return snap.EstimatedTime
		}

		return elapsed
	default:
		return snap.EstimatedTime
	}
}

// EstimatedCriticalPathRemaining returns the remaining time along the longest
// path from any root to any leaf, considering only work that is not yet
// finished. Completed/failed nodes contribute 0; running nodes contribute
// max(0, estimated - elapsed); pending nodes contribute their full estimate.
// This is the DAG-aware replacement for a naïve sum of all remaining work.
func (dt *DependencyTree) EstimatedCriticalPathRemaining(
	snapshots map[ActivityID]ActivitySnapshot,
) time.Duration {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if len(dt.nodes) == 0 {
		return 0
	}

	remainingTo := make(map[ActivityID]time.Duration, len(dt.nodes))

	maxIter := len(dt.nodes) + 1

	for range maxIter {
		changed := false

		for _, node := range dt.nodes {
			maxDep := time.Duration(0)

			for _, depID := range node.Deps {
				if depRem, ok := remainingTo[depID]; ok && depRem > maxDep {
					maxDep = depRem
				}
			}

			newTotal := maxDep + nodeRemaining(node.ID, snapshots)
			if remainingTo[node.ID] != newTotal {
				remainingTo[node.ID] = newTotal
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	var maxRemaining time.Duration

	for _, d := range remainingTo {
		if d > maxRemaining {
			maxRemaining = d
		}
	}

	return maxRemaining
}

// CriticalPathIDs returns the set of activity IDs that lie on the longest
// estimated-time path through the DAG. Returns nil if the DAG is empty or
// all activities have zero estimated time. Useful for filtering views to
// show only bottleneck activities.
func (dt *DependencyTree) CriticalPathIDs(
	snapshots map[ActivityID]ActivitySnapshot,
) map[ActivityID]bool {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if len(dt.nodes) == 0 {
		return nil
	}

	return dt.computeCriticalPath(snapshots)
}

func nodeRemaining(id ActivityID, snapshots map[ActivityID]ActivitySnapshot) time.Duration {
	snap := lookupSnapshot(snapshots, id)

	switch snap.Status {
	case ActivityStatusPending:
		return snap.EstimatedTime
	case ActivityStatusCompleted, ActivityStatusFailed:
		return 0
	case ActivityStatusRunning:
		if snap.EstimatedTime <= 0 {
			return 0
		}

		remaining := snap.EstimatedTime - snap.CurrentElapsed
		if remaining < 0 {
			return 0
		}

		return remaining
	default:
		return 0
	}
}
