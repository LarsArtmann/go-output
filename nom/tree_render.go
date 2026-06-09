package nom

import (
	"fmt"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
)

// Render generates NOM-style tree rendering.
func (dt *DependencyTree) Render(maxHeight int) string {
	// Build tree if not already built (release lock before calling Build to avoid deadlock)
	dt.mu.RLock()
	needsBuild := !dt.loaded
	dt.mu.RUnlock()

	if needsBuild {
		err := dt.Build()
		if err != nil {
			return fmt.Sprintf("Error building tree: %v", err)
		}
	}
	// Now acquire read lock for rendering
	dt.mu.RLock()
	defer dt.mu.RUnlock()
	// If no maxHeight provided, show all nodes
	if maxHeight <= 0 {
		maxHeight = len(dt.nodes)
	}
	// Get display order with smart filtering
	displayNodes := dt.getDisplayNodes(maxHeight)

	var lines []string

	for _, node := range displayNodes {
		line := dt.renderNode(node, displayNodes)
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return msgNoActivitiesToDisplay
	}

	return strings.Join(lines, "\n")
}

// getDisplayNodes returns nodes to display based on smart filtering.
func (dt *DependencyTree) getDisplayNodes(maxHeight int) []*TreeNode {
	// Priority: Running > Failed > Pending > Completed
	var nodes []*TreeNode
	// Categorize nodes by status
	running, failed, pending, completed := dt.categorizeNodes()
	// Sort each category by activity ID
	dt.sortNodesByID(running)
	dt.sortNodesByID(failed)
	dt.sortNodesByID(pending)
	dt.sortNodesByID(completed)
	// Allocate slots based on priority
	slots := maxHeight
	// 1. Add all running activities
	slots = dt.addNodesIfSlotsRemaining(&nodes, running, slots)
	// 2. Add all failed activities
	slots = dt.addNodesIfSlotsRemaining(&nodes, failed, slots)
	// 3. Add pending activities (up to remaining slots)
	slots = dt.addNodesIfSlotsRemaining(&nodes, pending, slots)
	// 4. Add completed activities only if slots remain
	dt.addNodesIfSlotsRemaining(&nodes, completed, slots)
	// Sort final list by depth, then by activity ID
	dt.sortNodesByDepthAndID(nodes)

	return nodes
}

func (dt *DependencyTree) categorizeNodes() (running, failed, pending, completed []*TreeNode) {
	for _, node := range dt.nodes {
		switch node.Status {
		case ActivityStatusRunning:
			running = append(running, node)
		case ActivityStatusFailed:
			failed = append(failed, node)
		case ActivityStatusPending, ActivityStatusPaused:
			pending = append(pending, node)
		case ActivityStatusCompleted:
			completed = append(completed, node)
		}
	}

	return running, failed, pending, completed
}

func (dt *DependencyTree) sortNodesByID(nodes []*TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		return string(nodes[i].ActivityID) < string(nodes[j].ActivityID)
	})
}

func (dt *DependencyTree) addNodesIfSlotsRemaining(
	target *[]*TreeNode,
	nodes []*TreeNode,
	slots int,
) int {
	for _, node := range nodes {
		if slots > 0 {
			*target = append(*target, node)
			slots--
		}
	}

	return slots
}

func (dt *DependencyTree) sortNodesByDepthAndID(nodes []*TreeNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Depth != nodes[j].Depth {
			return nodes[i].Depth < nodes[j].Depth
		}

		return string(nodes[i].ActivityID) < string(nodes[j].ActivityID)
	})
}

// renderNode renders a single node with appropriate tree symbols.
func (dt *DependencyTree) renderNode(node *TreeNode, displayNodes []*TreeNode) string {
	// Build prefix for tree structure
	prefix := dt.buildTreePrefix(node, displayNodes)
	// Build activity display with timing
	activityDisplay := fmt.Sprintf("%s %s", node.Symbol, node.ActivityName)
	// Add timing information based on status
	timingInfo := FormatTreeNodeTiming(
		node.Status,
		node.CurrentElapsed,
		node.EstimatedTime,
	)
	if timingInfo != "" {
		activityDisplay += " " + timingInfo
	}
	// Style with color - force ANSI colors
	style := lipgloss.NewStyle().
		Foreground(node.Color).
		Width(0).    // Don't limit width
		Inline(true) // Inline mode for better compatibility
	// Force color output by setting a color profile
	// lipgloss v2 handles color profile automatically
	return style.Render(prefix + activityDisplay)
}
