package nom

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// RenderString generates NOM-style tree rendering as a string using depth-first forest walk.
// Ports the algorithm from nix-output-monitor's showForest:
// walks each root tree recursively, building proper box-drawing prefixes (├──, └──, │).
//
// If maxWidth is > 0, long activity lines are truncated with "…" so they fit.
// Split-brain M4 resolved: renamed from Render(maxHeight) to distinguish from
// the output.Renderer.Render() (string, error) contract. This method returns a
// bare string (errors are encoded into the string itself, not returned separately).
func (dt *DependencyTree) RenderString(maxHeight int) string {
	return dt.RenderWithWidth(maxHeight, 0)
}

// RenderWithWidth generates the tree rendering with an optional terminal width
// constraint. A maxWidth of 0 disables truncation.
func (dt *DependencyTree) RenderWithWidth(maxHeight, maxWidth int) string {
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

	visible := dt.collectVisibleNodes(maxHeight)

	if len(visible) == 0 {
		return msgNoActivitiesToDisplay
	}

	var lines []string

	for _, entry := range visible {
		lines = append(lines, dt.renderLine(entry, maxWidth))
	}

	return strings.Join(lines, "\n")
}

type visibleEntry struct {
	node      *ActivityNode
	prefix    string
	connector string
	isRoot    bool
}

// collectVisibleNodes returns the most interesting nodes to display, up to maxHeight,
// preserving tree prefixes. Children at each level are sorted by status (failed > running >
// paused > pending > completed), then by elapsed time, so the user sees the currently
// important concurrent work first when screen space is limited.
func (dt *DependencyTree) collectVisibleNodes(maxHeight int) []visibleEntry {
	var visible []visibleEntry

	for _, root := range dt.roots {
		dt.walkSubtree(root, "", true, true, &visible, maxHeight)

		if len(visible) >= maxHeight {
			break
		}
	}

	return visible
}

// elideCompletedUnderPressure removes completed children when height is limited
// and expanding all children would overflow the available space. This ensures
// that active work (running, failed, pending) is prioritized over completed
// history when the viewport is constrained.
func (dt *DependencyTree) elideCompletedUnderPressure(
	children []*ActivityNode,
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
		if child.Status == ActivityStatusCompleted {
			continue
		}

		active = append(active, child)
	}

	if len(active) < len(children) {
		return active
	}

	return children
}

// walkSubtree recursively walks the tree depth-first, building proper NOM-style prefixes.
// Children are traversed in priority order (failed/running first, completed last).
func (dt *DependencyTree) walkSubtree(
	node *ActivityNode,
	prefix string,
	isLastSibling bool,
	isRoot bool,
	visible *[]visibleEntry,
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

	children := dt.childPriority(node)
	if len(children) == 0 {
		return
	}

	children = dt.elideCompletedUnderPressure(children, maxHeight, len(*visible))

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
			maxHeight,
		)
	}
}

// renderLine renders a single line for a visible entry. If maxWidth > 0, the
// activity name is truncated so the whole line fits the terminal.
func (dt *DependencyTree) renderLine(entry visibleEntry, maxWidth int) string {
	node := entry.node
	symbol := node.Symbol
	color := node.Color

	if isPhaseNode(node) {
		symbol = SymbolPhase
		color = ColorPhase
	}

	activityDisplay := fmt.Sprintf("%s %s", symbol, node.Label.Get())

	timingInfo := FormatActivityNodeTiming(
		node.Status,
		node.CurrentElapsed,
		node.EstimatedTime,
	)

	if timingInfo != "" {
		activityDisplay += " " + timingInfo
	}

	if len(node.SecondaryParents) > 0 {
		depNames := make([]string, len(node.SecondaryParents))

		for i, depID := range node.SecondaryParents {
			if depNode, ok := dt.nodes[depID]; ok {
				depNames[i] = depNode.Label.Get()
			} else {
				depNames[i] = depID.String()
			}
		}

		activityDisplay += lipgloss.NewStyle().
			Foreground(ColorInfo).
			Render(" ⬅ depends on " + strings.Join(depNames, ", "))
	}

	fullPrefix := entry.prefix + entry.connector

	if maxWidth > 0 {
		available := maxWidth - ansi.StringWidth(fullPrefix)
		activityDisplay = TruncateVisible(activityDisplay, available)
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Width(0).
		Inline(true)

	return style.Render(fullPrefix + activityDisplay)
}

// VisibleNodes returns the ordered list of tree nodes that would be displayed
// for the given maxHeight, in priority order (failed > running > paused >
// pending > completed).
func (dt *DependencyTree) VisibleNodes(maxHeight int) []*ActivityNode {
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

	visible := dt.collectVisibleNodes(maxHeight)
	nodes := make([]*ActivityNode, len(visible))

	for i, entry := range visible {
		nodes[i] = entry.node
	}

	return nodes
}

// RenderNode renders a single node for external consumers (e.g., TUI mouse
// click highlight). The visibleNodes parameter provides sibling context for
// callers that need width-aware rendering; it is currently unused by this
// method but accepted to support future width-truncation without a breaking
// API change.
func (dt *DependencyTree) RenderNode(node *ActivityNode, visibleNodes []*ActivityNode) string {
	symbol := node.Symbol
	color := node.Color

	if isPhaseNode(node) {
		symbol = SymbolPhase
		color = ColorPhase
	}

	activityDisplay := fmt.Sprintf("%s %s", symbol, node.Label.Get())

	timingInfo := FormatActivityNodeTiming(
		node.Status,
		node.CurrentElapsed,
		node.EstimatedTime,
	)

	if timingInfo != "" {
		activityDisplay += " " + timingInfo
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Width(0).
		Inline(true).
		Render(activityDisplay)
}

func isPhaseNode(node *ActivityNode) bool {
	return strings.HasPrefix(node.ID.Get(), "phase:")
}
