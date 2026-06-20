package tui

import (
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output/nom"
)

// renderNOMStyle creates a NOM-style display with dependency tree.
func (m *ProgressModel) renderNOMStyle() string {
	m.treeStartLine = 0

	var contentSections []string

	if m.dependencyTree != nil {
		treeSection := m.renderDependencyTree()
		if treeSection != "" {
			contentSections = append(contentSections, treeSection)
		}
	}

	result := m.assembleProgressSections(
		"🔄 NOM Workflow Execution",
		contentSections,
		m.renderNOMSummaryBar,
	)
	m.treeStartLine = 2

	return result
}

// renderDependencyTree renders the activity dependency tree in NOM style.
//
// Uses immutable ActivitySnapshot data for all field reads — no shared
// *Activity pointer is accessed, so event handlers (SetRunning/SetCompleted/
// SetFailed) mutating activities on dispatcher goroutines cannot race the
// render. The snapshot is taken under the subscriber's read lock (brief),
// then the tree walk reads only immutable data (no lock held).
func (m *ProgressModel) renderDependencyTree() string {
	if m.dependencyTree == nil || m.nomSubscriber == nil {
		return ""
	}

	snapshots := m.nomSubscriber.SnapshotActivities()
	tree := m.dependencyTree

	if m.scrollOffset > 0 {
		return tree.RenderWithSnapshots(snapshots, 0, m.width)
	}

	treeHeight := m.height - chromeLines
	if treeHeight <= 0 {
		treeHeight = defaultTreeHeight
	}

	m.visibleNodes = tree.VisibleNodesWithSnapshots(snapshots, treeHeight)

	if len(m.visibleNodes) == 0 {
		return msgNoActivitiesToDisplay
	}

	lines := make([]string, 0, len(m.visibleNodes))

	for _, node := range m.visibleNodes {
		line := tree.RenderNode(node, m.visibleNodes, snapshots)
		if m.selectedNode != "" && node.ID == m.selectedNode {
			line = lipgloss.NewStyle().
				Background(colors.selectBG).
				Foreground(colors.selectFG).
				Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderNOMSummaryBar creates the NOM-style summary bar, colored by workflow state.
func (m *ProgressModel) renderNOMSummaryBar() string {
	counts := m.getActivityCounts()
	elapsed := time.Since(m.startTime)
	summary := buildNOMSummary(counts, elapsed)
	baseStyle := createSummaryStyle()

	switch m.workflowState {
	case workflowStateIdle, workflowStateRunning:
		return baseStyle.Render(summary)
	case workflowStateCompleted:
		return baseStyle.Foreground(colors.success).Render("✅ " + summary)
	case workflowStateErrored:
		return baseStyle.Foreground(colors.err).Render("❌ " + summary)
	default:
		return baseStyle.Render(summary)
	}
}

// getActivityCounts delegates to the subscriber for counts.
func (m *ProgressModel) getActivityCounts() nom.ActivityCounts {
	return m.nomSubscriber.GetActivityCounts()
}

func (m *ProgressModel) renderHelpOverlay() string {
	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(colors.helpFG).
		Background(colors.helpBG)

	shortcuts := []string{
		"  Keyboard Shortcuts",
		"",
		"  j / ↓       Scroll down",
		"  k / ↑       Scroll up",
		"  pgdown      Scroll half page",
		"  pgup        Scroll half page up",
		"  g / Home    Scroll to top",
		"  G / End     Scroll to bottom",
		"  ?           Toggle this help",
		"  q / ctrl+c  Quit",
	}

	helpText := helpStyle.Render(strings.Join(shortcuts, "\n"))

	width := m.width
	height := m.height

	if width <= 0 {
		width = defaultHelpWidth
	}

	if height <= 0 {
		height = defaultHelpHeight
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, helpText)
}
