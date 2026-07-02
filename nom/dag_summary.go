package nom

import (
	"fmt"
	"strings"
	"time"
)

// DAGSummary describes the structural shape of the dependency DAG.
type DAGSummary struct {
	Nodes    int
	Edges    int
	MaxDepth int
	MaxWidth int
	Roots    int
	Leaves   int
	Phases   int
	Critical time.Duration
}

// String renders a compact one-line summary like "12 nodes · 18 edges · 4 layers · 3 parallel".
func (s DAGSummary) String() string {
	parts := []string{
		fmt.Sprintf("%d nodes", s.Nodes),
		fmt.Sprintf("%d edges", s.Edges),
	}

	if s.MaxDepth > 0 {
		parts = append(parts, fmt.Sprintf("%d layers", s.MaxDepth+1))
	}

	if s.MaxWidth > 1 {
		parts = append(parts, fmt.Sprintf("%d wide", s.MaxWidth))
	}

	if s.Phases > 0 {
		parts = append(parts, fmt.Sprintf("%d phases", s.Phases))
	}

	return joinParts(parts)
}

// DAGSummary computes structural metrics about the dependency tree: node count,
// edge count, maximum depth, maximum layer width, and root/leaf counts.
// Thread-safe.
func (dt *DependencyTree) DAGSummary() DAGSummary {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	summary := DAGSummary{
		Nodes: len(dt.nodes),
	}

	if len(dt.nodes) == 0 {
		return summary
	}

	byDepth := make(map[int]int)
	depSet := make(map[ActivityID]bool)

	for _, node := range dt.nodes {
		summary.Edges += len(node.Deps)

		for _, dep := range node.Deps {
			depSet[dep] = true
		}

		byDepth[node.Depth]++

		if node.Depth > summary.MaxDepth {
			summary.MaxDepth = node.Depth
		}

		if node.IsRoot || len(node.Deps) == 0 {
			summary.Roots++
		}

		if len(node.Children) == 0 && !depSet[node.ID] {
			summary.Leaves++
		} else if len(node.Children) == 0 {
			summary.Leaves++
		}
	}

	for _, count := range byDepth {
		if count > summary.MaxWidth {
			summary.MaxWidth = count
		}
	}

	return summary
}

// DAGSummaryWithSnapshots computes the structural summary plus the critical-path
// remaining time and phase count using the given activity snapshots. Thread-safe.
func (dt *DependencyTree) DAGSummaryWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
) DAGSummary {
	summary := dt.DAGSummary()
	summary.Critical = dt.EstimatedCriticalPathRemaining(snapshots)

	for _, snap := range snapshots {
		if snap.IsPhase() {
			summary.Phases++
		}
	}

	return summary
}

func joinParts(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " · ")
}
