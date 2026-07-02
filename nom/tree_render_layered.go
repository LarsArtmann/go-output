package nom

import (
	"cmp"
	"fmt"
	"slices"
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

	for _, node := range dt.nodes {
		// Skip structural placeholders that were never registered as real
		// activities. They have no snapshot label and would render as a blank
		// pending line, which is confusing in layered mode.
		if dt.isPlaceholderNode(node, snapshots) {
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

		for _, node := range nodes {
			if len(entries) >= maxHeight {
				break
			}

			entries = append(entries, VisibleEntry{LayerNodes: []*ActivityNode{node}})
		}

		// Separator between layers, but not after the last one.
		if depth < maxDepth && len(entries) < maxHeight {
			entries = append(entries, VisibleEntry{LayerHeader: layeredSeparator(maxDepth)})
		}
	}

	return entries
}

// isPlaceholderNode reports whether a node is a structural placeholder that
// was never registered as a real activity. Placeholders have no snapshot label
// and are created only to satisfy a dependency edge before the activity event
// arrives. In layered mode they are skipped because they would render as a
// blank line.
func (dt *DependencyTree) isPlaceholderNode(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) bool {
	snap := lookupSnapshot(snapshots, node.ID)
	return snap.Label == "" && len(node.Deps) == 0
}

// sortNodesByPriority orders nodes by activity status priority (running >
// pending > failed > completed), then by ID for stability. This surfaces live
// work before completed work in each layer.
func (dt *DependencyTree) sortNodesByPriority(
	nodes []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) {
	slices.SortStableFunc(nodes, func(a, b *ActivityNode) int {
		pa := statusPriority(lookupSnapshot(snapshots, a.ID).Status)

		pb := statusPriority(lookupSnapshot(snapshots, b.ID).Status)
		if pa != pb {
			return pb - pa // higher priority first
		}

		return cmp.Compare(string(a.ID), string(b.ID))
	})
}

// statusPriority returns a sort rank for status ordering: running > pending >
// failed > completed. Lower values are less urgent.
func statusPriority(s ActivityStatus) int {
	switch s {
	case ActivityStatusRunning:
		return 4
	case ActivityStatusPending:
		return 3
	case ActivityStatusFailed:
		return 2
	case ActivityStatusCompleted:
		return 1
	default:
		return 0
	}
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

func layeredSeparator(maxDepth int) string {
	width := layeredHeaderWidth - 1 // "Layer N │" width is header width

	// Extend the separator to account for deeper layer numbers (e.g. "Layer 12").
	if maxDepth >= 10 {
		width += len(strconv.Itoa(maxDepth)) - 1
	}

	return strings.Repeat("─", width) + "┼" + strings.Repeat("─", 12)
}

// renderLayered builds the full layered-mode string.
func (dt *DependencyTree) renderLayered(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, maxWidth int,
) string {
	entries := dt.collectLayeredEntries(snapshots, maxHeight)
	if len(entries) == 0 {
		return msgNoActivitiesToDisplay
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
	if entry.LayerHeader != "" {
		return dt.renderLayeredHeader(entry.LayerHeader, maxWidth)
	}

	if len(entry.LayerNodes) == 0 {
		return ""
	}

	var parts []string

	for _, node := range entry.LayerNodes {
		snap := lookupSnapshot(snapshots, node.ID)
		label, color := formatActivityLabelWithOptions(snap, dt.theme, labelOptions{})
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

func (dt *DependencyTree) renderLayeredHeader(header string, maxWidth int) string {
	isSeparator := strings.ContainsRune(header, '┼') || strings.ContainsRune(header, '─')
	if isSeparator {
		if maxWidth > 0 {
			return TruncateVisible(header, maxWidth)
		}

		return header
	}

	left := fmt.Sprintf("%-*s", layeredHeaderWidth-1, header) + "│"

	line := left

	if maxWidth > 0 {
		line = TruncateVisible(line, maxWidth)
	}

	return line
}
