package nom

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// layeredHeaderWidth is the fixed width of the layer label column ("Layer 99 │").
const layeredHeaderWidth = 10

// collectLayeredEntries groups nodes by their DAG depth and returns one
// VisibleEntry per line: layer headers, one-node activity rows, and separator
// lines. Each VisibleEntry is a single scroll line in the TUI; horizontal
// wrapping is left to render time so callers can supply a width budget.
//
//nolint:cyclop // depth grouping with collapse and priority sorting is inherently branchy
func (dt *DependencyTree) collectLayeredEntries(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []VisibleEntry {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if maxHeight <= 0 {
		maxHeight = len(dt.nodes)*3 + 1
	}

	byDepth := make(map[int][]*ActivityNode)
	maxDepth := -1

	depSet := dt.dependencyIDsLocked()

	for _, node := range dt.nodes {
		// Skip structural placeholders that were never registered as real
		// activities. They were created only to satisfy another node's
		// dependency edge and have no snapshot label, so they would render as
		// a blank pending line.
		if dt.isPlaceholderNode(node, snapshots, depSet) {
			continue
		}

		byDepth[node.Depth] = append(byDepth[node.Depth], node)
		if node.Depth > maxDepth {
			maxDepth = node.Depth
		}
	}

	if maxDepth < 0 {
		return nil
	}

	var entries []VisibleEntry

	remainingLayers := maxDepth + 1

	for depth := 0; depth <= maxDepth; depth++ {
		if len(entries) >= maxHeight {
			break
		}

		nodes := byDepth[depth]
		if len(nodes) == 0 {
			remainingLayers--

			continue
		}

		dt.sortNodesByPriority(nodes, snapshots)

		// Future-layer hiding: when enabled, collapse layers where ALL nodes
		// are still pending (not yet started) into a single summary line.
		if dt.hideFutureLayers && dt.allNodesPending(nodes, snapshots) {
			entries = append(entries, VisibleEntry{
				LayerHeader: fmt.Sprintf("%s Layer %d: %d pending", SymbolPending, depth, len(nodes)),
			})
			remainingLayers--

			continue
		}

		// Height-pressure collapse: if all nodes in this layer are terminal and
		// we still have more layers to fit than rows left, summarize the layer.
		allTerminal := dt.allNodesTerminal(nodes, snapshots)
		if allTerminal && remainingLayers >= maxHeight-len(entries) {
			entries = append(entries, VisibleEntry{
				LayerHeader: fmt.Sprintf("%s Layer %d: %d done", SymbolCompleted, depth, len(nodes)),
			})
			remainingLayers--

			continue
		}

		entries = append(entries, VisibleEntry{LayerHeader: fmt.Sprintf("Layer %d", depth)})
		remainingLayers--

		if len(entries) >= maxHeight {
			break
		}

		entries = append(entries, dt.collectLayeredNodeEntries(nodes, snapshots, maxHeight-len(entries))...)

		// Separator between layers, but not after the last one.
		if depth < maxDepth && len(entries) < maxHeight {
			entries = append(entries, VisibleEntry{LayerHeader: layeredSeparator(maxDepth), separator: true})
		}
	}

	return entries
}

// isPlaceholderNode reports whether a node is a structural placeholder that
// was never registered as a real activity. Placeholders are created only to
// satisfy another node's dependency edge, have no snapshot label, and are
// skipped in layered mode so they do not render as blank lines.
func (dt *DependencyTree) isPlaceholderNode(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	depSet map[ActivityID]bool,
) bool {
	snap := lookupSnapshot(snapshots, node.ID)
	return snap.Label == "" && depSet[node.ID]
}

// dependencyIDsLocked returns the set of IDs that appear as dependencies of
// other nodes in the tree. Caller must hold dt.mu.RLock.
func (dt *DependencyTree) dependencyIDsLocked() map[ActivityID]bool {
	deps := make(map[ActivityID]bool)

	for _, node := range dt.nodes {
		for _, dep := range node.Deps {
			deps[dep] = true
		}
	}

	return deps
}

// sortNodesByPriority orders nodes by display urgency using the status
// registry's Interest ranking (lower = more urgent: failed < running <
// pending < completed), then by ID for stability. This matches tree-mode
// ordering (sortKey.interest) so a failed activity surfaces first in BOTH
// display modes, and custom registered statuses get their Interest for
// free instead of sinking below completed work.
func (dt *DependencyTree) sortNodesByPriority(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) {
	slices.SortStableFunc(nodes, func(a, b *ActivityNode) int {
		pa := lookupSnapshot(snapshots, a.ID).Status.Interest()

		pb := lookupSnapshot(snapshots, b.ID).Status.Interest()
		if pa != pb {
			return pa - pb // lower interest value = more urgent, first
		}

		return cmp.Compare(string(a.ID), string(b.ID))
	})
}

func (dt *DependencyTree) allNodesTerminal(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) bool {
	for _, node := range nodes {
		snap := lookupSnapshot(snapshots, node.ID)
		if snap.Status != ActivityStatusCompleted && snap.Status != ActivityStatusFailed {
			return false
		}
	}

	return true
}

// allNodesPending reports whether every node in the layer is still pending
// (not yet started). Used for future-layer hiding.
func (dt *DependencyTree) allNodesPending(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) bool {
	for _, node := range nodes {
		snap := lookupSnapshot(snapshots, node.ID)
		if snap.Status != ActivityStatusPending {
			return false
		}
	}

	return true
}

// categoryGroup tracks a collapsed category's name and node count.
type categoryGroup struct {
	name  string
	count int
}

// collectLayeredNodeEntries produces VisibleEntry rows for the nodes in a
// single layer, applying category collapse when enabled. maxHeight is the
// remaining entry budget before the caller must stop.
func (dt *DependencyTree) collectLayeredNodeEntries(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	maxEntries int,
) []VisibleEntry {
	if dt.collapseCategories && dt.showCategory {
		return dt.collectCategoryCollapsedEntries(nodes, snapshots, maxEntries)
	}

	var entries []VisibleEntry

	for _, node := range nodes {
		if len(entries) >= maxEntries {
			break
		}

		entries = append(entries, VisibleEntry{LayerNodes: []*ActivityNode{node}})
	}

	return entries
}

// collectCategoryCollapsedEntries groups nodes by category, emitting summary
// lines for categories with 2+ nodes and individual entries for the rest.
func (dt *DependencyTree) collectCategoryCollapsedEntries(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	maxEntries int,
) []VisibleEntry {
	collapsed, individuals := groupByCategory(nodes, snapshots)

	var entries []VisibleEntry

	for _, cat := range collapsed {
		if len(entries) >= maxEntries {
			break
		}

		entries = append(entries, VisibleEntry{
			LayerHeader: fmt.Sprintf("%s %d %s tasks", SymbolPhase, cat.count, cat.name),
		})
	}

	for _, node := range individuals {
		if len(entries) >= maxEntries {
			break
		}

		entries = append(entries, VisibleEntry{LayerNodes: []*ActivityNode{node}})
	}

	return entries
}

// groupByCategory partitions nodes into collapsed groups (2+ nodes sharing a
// category) and individual nodes (no category or unique category). Returns
// (collapsed groups, individual nodes).
func groupByCategory(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) ([]categoryGroup, []*ActivityNode) {
	counts := make(map[string]int)

	for _, node := range nodes {
		snap := lookupSnapshot(snapshots, node.ID)
		if snap.Category != "" {
			counts[string(snap.Category)]++
		}
	}

	var collapsed []categoryGroup

	var individuals []*ActivityNode

	for _, node := range nodes {
		snap := lookupSnapshot(snapshots, node.ID)
		cat := string(snap.Category)

		if cat != "" && counts[cat] >= 2 {
			continue // will be collapsed
		}

		individuals = append(individuals, node)
	}

	// Sort category names so collapsed groups render deterministically —
	// map iteration order is random.
	cats := make([]string, 0, len(counts))
	for cat := range counts {
		cats = append(cats, cat)
	}

	sort.Strings(cats)

	for _, cat := range cats {
		if count := counts[cat]; count >= 2 {
			collapsed = append(collapsed, categoryGroup{name: cat, count: count})
		}
	}

	return collapsed, individuals
}

func layeredSeparator(maxDepth int) string {
	// The ┼ must sit exactly under the header row's │. That column is
	// max(layeredHeaderWidth-1, len("Layer N")) — Sprintf pads the label to
	// layeredHeaderWidth-1 but never truncates deeper layer numbers.
	width := max(layeredHeaderWidth-1, len("Layer "+strconv.Itoa(maxDepth)))

	return strings.Repeat("─", width) + "┼" + strings.Repeat("─", 12)
}

// renderLayered builds the full layered-mode string.
func (dt *DependencyTree) renderLayered(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, maxWidth int,
) string {
	entries := dt.collectLayeredEntries(snapshots, maxHeight)
	if len(entries) == 0 {
		return MsgNoActivities
	}

	var lines []string

	for _, entry := range entries {
		lines = append(lines, dt.renderLayeredLine(entry, snapshots, maxWidth))
	}

	return strings.Join(lines, "\n")
}

func (dt *DependencyTree) renderLayeredLine(
	entry VisibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxWidth int,
) string {
	if entry.Kind() == KindSeparator {
		return dt.renderLayeredSeparator(entry.LayerHeader, maxWidth)
	}

	if entry.LayerHeader != "" {
		return dt.renderLayeredHeader(entry.LayerHeader, maxWidth)
	}

	if len(entry.LayerNodes) == 0 {
		return ""
	}

	var parts []string

	for _, node := range entry.LayerNodes {
		snap := lookupSnapshot(snapshots, node.ID)
		label, color := formatActivityLabelWithOptions(snap, dt.theme, labelOptions{
			ShowCategory: dt.showCategory,
		})
		styled := activityNodeStyle(color).Render(label)
		parts = append(parts, styled)
	}

	header := fmt.Sprintf("%-*s", layeredHeaderWidth-1, "") + "│ "
	line := header + strings.Join(parts, "  ")

	if maxWidth > 0 {
		line = TruncateVisible(line, maxWidth)
	}

	return line
}

func (dt *DependencyTree) renderLayeredSeparator(line string, maxWidth int) string {
	if maxWidth > 0 {
		return TruncateVisible(line, maxWidth)
	}

	return line
}

func (dt *DependencyTree) renderLayeredHeader(header string, maxWidth int) string {
	left := fmt.Sprintf("%-*s", layeredHeaderWidth-1, header) + "│"

	if maxWidth > 0 {
		return TruncateVisible(left, maxWidth)
	}

	return left
}
