package nom

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// RenderWithSnapshots generates the tree rendering using immutable activity
// snapshots instead of reading the shared *Activity pointer. This is the
// canonical render path: snapshots are taken under the subscriber's read lock
// once, then the tree walk reads only immutable data. A nil snapshots map is
// treated as "all activities pending with blank labels" — useful for tests
// that build trees structurally without a subscriber.
func (dt *DependencyTree) RenderWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, maxWidth int,
) string {
	dt.mu.RLock()
	needsBuild := !dt.loaded
	dt.mu.RUnlock()

	if needsBuild {
		err := dt.Build()
		if err != nil {
			return fmt.Sprintf("Error building tree: %v", err)
		}
	}

	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if maxHeight <= 0 {
		maxHeight = len(dt.nodes)
	}

	if len(dt.roots) == 0 {
		return msgNoActivitiesToDisplay
	}

	visible := dt.collectVisibleNodes(snapshots, maxHeight)

	if len(visible) == 0 {
		return msgNoActivitiesToDisplay
	}

	var lines []string

	for _, entry := range visible {
		lines = append(lines, dt.renderLine(entry, snapshots, maxWidth))
	}

	return strings.Join(lines, "\n")
}

type visibleEntry struct {
	node      *ActivityNode
	prefix    string
	connector string
	isRoot    bool
}

func (dt *DependencyTree) collectVisibleNodes(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []visibleEntry {
	var visible []visibleEntry

	for _, root := range dt.roots {
		dt.walkSubtree(root, "", true, true, &visible, snapshots, maxHeight)

		if len(visible) >= maxHeight {
			break
		}
	}

	return visible
}

func (dt *DependencyTree) elideCompletedUnderPressure(
	children []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, visibleCount int,
) []*ActivityNode {
	if maxHeight <= 0 {
		return children
	}

	remaining := maxHeight - visibleCount
	if remaining >= len(children) {
		return children
	}

	var active []*ActivityNode

	for _, child := range children {
		snap := lookupSnapshot(snapshots, child.ID)
		if snap.Status == ActivityStatusCompleted {
			continue
		}

		active = append(active, child)
	}

	if len(active) < len(children) {
		return active
	}

	return children
}

func (dt *DependencyTree) walkSubtree(
	node *ActivityNode,
	prefix string,
	isLastSibling bool,
	isRoot bool,
	visible *[]visibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) {
	if len(*visible) >= maxHeight {
		return
	}

	entry := visibleEntry{
		node:      node,
		prefix:    prefix,
		connector: "",
		isRoot:    isRoot,
	}

	if !isRoot {
		if isLastSibling {
			entry.connector = "└── "
		} else {
			entry.connector = "├── "
		}
	}

	*visible = append(*visible, entry)

	children := dt.childPriority(node, snapshots)
	if len(children) == 0 {
		return
	}

	children = dt.elideCompletedUnderPressure(children, snapshots, maxHeight, len(*visible))

	var childIndent string

	if isRoot {
		childIndent = ""
	} else if isLastSibling {
		childIndent = prefix + "    "
	} else {
		childIndent = prefix + "│   "
	}

	for i, child := range children {
		if len(*visible) >= maxHeight {
			return
		}

		dt.walkSubtree(
			child,
			childIndent,
			i == len(children)-1,
			false,
			visible,
			snapshots,
			maxHeight,
		)
	}
}

// formatActivityLabel builds the core display string for a single activity
// snapshot: phase-aware symbol + label + timing info. Returns the unstyled
// display and the status-derived color for the caller to apply.
func formatActivityLabel(snap ActivitySnapshot, id ActivityID) (display string, c color.Color) {
	symbol := snap.Symbol
	c = snap.Color

	if isPhaseID(id) {
		symbol = SymbolPhase
		c = Colors.Phase
	}

	display = fmt.Sprintf("%s %s", symbol, snap.Label)

	timingInfo := FormatActivityNodeTiming(
		snap.Status,
		snap.CurrentElapsed,
		snap.EstimatedTime,
	)

	if timingInfo != "" {
		display += " " + timingInfo
	}

	return display, c
}

func (dt *DependencyTree) renderLine(
	entry visibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxWidth int,
) string {
	node := entry.node
	snap := lookupSnapshot(snapshots, node.ID)

	activityDisplay, color := formatActivityLabel(snap, node.ID)

	if len(node.SecondaryParents) > 0 {
		depNames := make([]string, len(node.SecondaryParents))

		for i, depID := range node.SecondaryParents {
			depSnap := lookupSnapshot(snapshots, depID)

			depNames[i] = depSnap.Label
			if depNames[i] == "" {
				depNames[i] = depID.String()
			}
		}

		activityDisplay += lipgloss.NewStyle().
			Foreground(Colors.Info).
			Render(" ⬅ depends on " + strings.Join(depNames, ", "))
	}

	fullPrefix := entry.prefix + entry.connector

	if maxWidth > 0 {
		available := maxWidth - ansi.StringWidth(fullPrefix)
		activityDisplay = TruncateVisible(activityDisplay, available)
	}

	rendered := activityNodeStyle(color).Render(fullPrefix + activityDisplay)

	if maxWidth > 0 && VisibleWidth(rendered) > maxWidth {
		rendered = TruncateVisible(rendered, maxWidth)
	}

	return rendered
}

// VisibleNodesWithSnapshots returns the ordered list of tree nodes that would
// be displayed for the given maxHeight, in priority order. Uses snapshots for
// status-based sorting.
func (dt *DependencyTree) VisibleNodesWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []*ActivityNode {
	dt.mu.RLock()
	needsBuild := !dt.loaded
	dt.mu.RUnlock()

	if needsBuild {
		if err := dt.Build(); err != nil {
			return nil
		}
	}

	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if maxHeight <= 0 {
		maxHeight = len(dt.nodes)
	}

	visible := dt.collectVisibleNodes(snapshots, maxHeight)
	nodes := make([]*ActivityNode, len(visible))

	for i, entry := range visible {
		nodes[i] = entry.node
	}

	return nodes
}

// RenderNode renders a single node for external consumers (e.g., TUI mouse
// click highlight). Uses the snapshot for label/color/symbol.
func (dt *DependencyTree) RenderNode(
	node *ActivityNode,
	_ []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) string {
	snap := lookupSnapshot(snapshots, node.ID)
	display, color := formatActivityLabel(snap, node.ID)

	return activityNodeStyle(color).Render(display)
}

func activityNodeStyle(color color.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color).
		Width(0).
		Inline(true)
}

// IsPhase reports whether this node represents a workflow phase.
func (n *ActivityNode) IsPhase() bool {
	return isPhaseID(n.ID)
}

// isPhaseID checks whether an activity ID follows the "phase:" prefix convention.
func isPhaseID(id ActivityID) bool {
	return strings.HasPrefix(string(id), "phase:")
}

// lookupSnapshot returns the snapshot for id, or a blank pending snapshot if
// the activity hasn't been registered yet (e.g. a structural placeholder for
// a dependency that hasn't received its activity.started event).
func lookupSnapshot(snapshots map[ActivityID]ActivitySnapshot, id ActivityID) ActivitySnapshot {
	if snapshots == nil {
		return blankActivitySnapshot()
	}

	if snap, ok := snapshots[id]; ok {
		return snap
	}

	return blankActivitySnapshot()
}

// blankActivitySnapshot returns the default snapshot for unregistered activities:
// pending status, empty label, pending symbol/color. Kept as a function (not a
// global var) so the snapshot is never accidentally mutated by callers.
func blankActivitySnapshot() ActivitySnapshot {
	return ActivitySnapshot{
		Label:  "",
		Status: ActivityStatusPending,
		Symbol: SymbolPending,
		Color:  Colors.Pending,
	}
}
