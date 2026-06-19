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
// The whole walk runs under the subscriber's read lock: the tree nodes embed a
// shared *Activity pointer whose fields event handlers mutate
// (SetRunning/SetCompleted/SetFailed) on dispatcher goroutines, and reading
// them unlocked races those writes — producing garbled frames and failing
// -race. m.dependencyTree is the cached pointer set by syncNOMSubscriber; in
// production it aliases the subscriber's tree, so holding ns.mu makes reading
// its Activity fields safe. We MUST NOT call other subscriber accessors
// (GetActivityCounts, RenderTree, …) inside the callback: recursive RLock
// against a waiting writer deadlocks.
func (m *ProgressModel) renderDependencyTree() string {
	// m.dependencyTree is the "synced at least once" flag (set by
	// syncNOMSubscriber). Reading just the pointer is non-racy under
	// bubbletea's serialized Update/View. The actual Activity fields are read
	// under the subscriber's read lock below.
	if m.dependencyTree == nil || m.nomSubscriber == nil {
		return ""
	}

	var result string

	m.nomSubscriber.WithSubscriberRLock(func() {
		tree := m.dependencyTree

		if m.scrollOffset > 0 {
			result = tree.RenderWithWidth(0, m.width)
			return
		}

		treeHeight := m.height - chromeLines
		if treeHeight <= 0 {
			treeHeight = defaultTreeHeight
		}

		m.visibleNodes = tree.VisibleNodes(treeHeight)

		if len(m.visibleNodes) == 0 {
			result = MsgNoActivitiesToDisplay
			return
		}

		lines := make([]string, 0, len(m.visibleNodes))

		for _, node := range m.visibleNodes {
			line := tree.RenderNode(node, m.visibleNodes)
			if m.selectedNode != "" && nom.ActivityID(node.ID.Get()) == m.selectedNode {
				line = lipgloss.NewStyle().
					Background(colors.selectBG).
					Foreground(colors.selectFG).
					Render(line)
			}

			lines = append(lines, line)
		}

		result = strings.Join(lines, "\n")
	})

	return result
}

// renderNOMSummaryBar creates the NOM-style summary bar, colored by workflow state.
func (m *ProgressModel) renderNOMSummaryBar() string {
	counts := m.getActivityCounts()
	elapsed := time.Since(m.startTime)
	summary := buildNOMSummary(counts, elapsed)
	baseStyle := createSummaryStyle()

	switch m.workflowState {
	case WorkflowStateIdle, WorkflowStateRunning:
		return baseStyle.Render(summary)
	case WorkflowStateCompleted:
		return baseStyle.Foreground(colors.success).Render("✅ " + summary)
	case WorkflowStateErrored:
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
