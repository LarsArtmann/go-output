package nom

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// Render generates NOM-style tree rendering using depth-first forest walk.
// Ports the algorithm from nix-output-monitor's showForest:
// walks each root tree recursively, building proper box-drawing prefixes (├──, └──, │).
func (dt *DependencyTree) Render(maxHeight int) string {
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
		lines = append(lines, dt.renderLine(entry))
	}

	return strings.Join(lines, "\n")
}

type visibleEntry struct {
	node      *TreeNode
	prefix    string
	connector string
	isRoot    bool
}

// collectVisibleNodes walks the forest depth-first, collecting nodes to display.
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

// walkSubtree recursively walks the tree depth-first, building proper NOM-style prefixes.
func (dt *DependencyTree) walkSubtree(
	node *TreeNode,
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

	if len(node.Children) == 0 {
		return
	}

	var childIndent string

	if isRoot {
		childIndent = ""
	} else if isLastSibling {
		childIndent = prefix + "    "
	} else {
		childIndent = prefix + "│   "
	}

	for i, child := range node.Children {
		if len(*visible) >= maxHeight {
			return
		}

		dt.walkSubtree(
			child,
			childIndent,
			i == len(node.Children)-1,
			false,
			visible,
			maxHeight,
		)
	}
}

// renderLine renders a single line for a visible entry.
func (dt *DependencyTree) renderLine(entry visibleEntry) string {
	node := entry.node
	symbol := node.Symbol
	color := node.Color

	if isPhaseNode(node) {
		symbol = SymbolPhase
		color = ColorPhase
	}

	activityDisplay := fmt.Sprintf("%s %s", symbol, node.ActivityName)

	timingInfo := FormatTreeNodeTiming(
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
				depNames[i] = depNode.ActivityName
			} else {
				depNames[i] = depID.String()
			}
		}

		activityDisplay += lipgloss.NewStyle().
			Foreground(ColorInfo).
			Render(fmt.Sprintf(" ⬅ depends on %s", strings.Join(depNames, ", ")))
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Width(0).
		Inline(true)

	fullPrefix := entry.prefix + entry.connector

	return style.Render(fullPrefix + activityDisplay)
}

// VisibleNodes returns the ordered list of tree nodes that would be displayed
// for the given maxHeight, in depth-first tree order.
func (dt *DependencyTree) VisibleNodes(maxHeight int) []*TreeNode {
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
	nodes := make([]*TreeNode, len(visible))

	for i, entry := range visible {
		nodes[i] = entry.node
	}

	return nodes
}

// RenderNode renders a single node for external consumers (e.g., TUI mouse click highlight).
func (dt *DependencyTree) RenderNode(node *TreeNode, _ []*TreeNode) string {
	symbol := node.Symbol
	color := node.Color

	if isPhaseNode(node) {
		symbol = SymbolPhase
		color = ColorPhase
	}

	activityDisplay := fmt.Sprintf("%s %s", symbol, node.ActivityName)

	timingInfo := FormatTreeNodeTiming(
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

func isPhaseNode(node *TreeNode) bool {
	return len(node.ActivityID) > 6 && node.ActivityID[:6] == "phase:"
}
