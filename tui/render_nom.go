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
func (m *ProgressModel) renderDependencyTree() string {
	if m.dependencyTree == nil {
		return ""
	}

	if m.scrollOffset > 0 {
		return m.dependencyTree.RenderWithWidth(0, m.width)
	}

	treeHeight := m.height - chromeLines
	if treeHeight <= 0 {
		treeHeight = defaultTreeHeight
	}

	m.visibleNodes = m.dependencyTree.VisibleNodes(treeHeight)

	if len(m.visibleNodes) == 0 {
		return MsgNoActivitiesToDisplay
	}

	var lines []string

	for _, node := range m.visibleNodes {
		line := m.dependencyTree.RenderNode(node, m.visibleNodes)
		if m.selectedNode != "" && nom.ActivityID(node.ID.Get()) == m.selectedNode {
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
