package tui
import (
"fmt"
"strings"
"time"
tea "charm.land/bubbletea/v2"
"charm.land/lipgloss/v2"
"github.com/larsartmann/go-output/nom"
)
func (m *ProgressModel) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}
	var content string
	if m.displayMode == DisplayModeNOM {
		content = m.renderNOMStyle()
	} else {
		content = m.renderUniversalWorkflowProgress()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}
// renderUniversalWorkflowProgress creates a display for universal workflow execution.
func (m *ProgressModel) renderUniversalWorkflowProgress() string {
	// Build content sections
	var contentSections []string
	// Render steps in tree format like nh darwin switch
	if len(m.steps) > 0 {
		stepsSection := m.renderSteps()
		contentSections = append(contentSections, stepsSection)
	}
	// Progress bar
	if m.currentProgress > 0 {
		progressSection := m.renderProgressBar()
		contentSections = append(contentSections, progressSection)
	}
	// Assemble all sections
	return m.assembleProgressSections(
		"🚀 Universal Workflow Execution",
		contentSections,
		m.renderSummaryBar,
	)
}
// renderCurrentMessage renders the current message with a consistent style across all rendering modes.
func (m *ProgressModel) renderCurrentMessage() string {
	if m.currentMessage == "" {
		return ""
	}
	messageStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("12"))
	return messageStyle.Render(m.currentMessage)
}
// renderTitle creates a styled title header.
func (m *ProgressModel) renderTitle(text string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)
	return titleStyle.Render(text)
}
// assembleProgressSections assembles the common progress display structure
// Handles: title, message, content sections, and summary.
func (m *ProgressModel) assembleProgressSections(
	title string,
	contentSections []string,
	summaryFunc func() string,
) string {
	var sections []string
	// Title header
	sections = append(sections, m.renderTitle(title))
	// Current message
	message := m.renderCurrentMessage()
	if message != "" {
		sections = append(sections, message)
	}
	// Mode-specific content
	for _, section := range contentSections {
		if section != "" {
			sections = append(sections, section)
		}
	}
	// Summary footer
	if summaryFunc != nil {
		summary := summaryFunc()
		if summary != "" {
			sections = append(sections, summary)
		}
	}
	return strings.Join(sections, "\n\n")
}
// renderSteps creates the tree-style step display like nh darwin switch.
func (m *ProgressModel) renderSteps() string {
	lines := make([]string, 0, len(m.steps))
	for i, step := range m.steps {
		line := m.renderStep(step, i == len(m.steps)-1)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
// renderStep renders a single step with nh-style formatting.
func (m *ProgressModel) renderStep(step ProgressStep, isLast bool) string {
	// Choose prefix
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	// Status icon and color
	var (
		icon  string
		style lipgloss.Style
	)
	if step.CompletedAt != nil {
		icon = "✅"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	} else if step.IsActive {
		icon = "🔄"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
	} else {
		icon = "⏸️"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
	}
	// Timing information like nh darwin switch
	var timing string
	if step.CompletedAt != nil {
		duration := step.CompletedAt.Sub(step.StartTime)
		timing = fmt.Sprintf(TimingFormat, duration.Seconds())
	} else if step.IsActive {
		elapsed := time.Since(step.StartTime)
		timing = fmt.Sprintf(TimingFormat, elapsed.Seconds())
	}
	// Step progress info
	stepInfo := ""
	if step.Total > 0 {
		stepInfo = fmt.Sprintf("(%d/%d)", step.Current, step.Total)
	}
	// Format the line
	line := fmt.Sprintf("%s%s %s", prefix, icon, step.Message)
	if stepInfo != "" {
		line = fmt.Sprintf("%s %s", line, stepInfo)
	}
	if timing != "" {
		line = fmt.Sprintf("%s %s", line, timing)
	}
	return style.Render(line)
}
// renderProgressBar creates a progress bar.
func (m *ProgressModel) renderProgressBar() string {
	width := 40
	if m.width > 0 && m.width < 80 {
		width = m.width - 30
	}
	filled := int((m.currentProgress / 100.0) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12"))
	return progressStyle.Render(fmt.Sprintf("Progress: [%s] %.1f%%", bar, m.currentProgress))
}
// renderSummaryBar creates the summary like nh darwin switch.
func (m *ProgressModel) renderSummaryBar() string {
	// Calculate summary statistics
	var completedSteps, inProgressSteps int
	for _, step := range m.steps {
		if step.CompletedAt != nil {
			completedSteps++
		} else if step.IsActive {
			inProgressSteps++
		}
	}
	elapsed := time.Since(m.startTime)
	// Build summary using helper
	summary := buildUniversalSummary(inProgressSteps, completedSteps, elapsed, m.currentProgress)
	// Apply state-specific formatting
	finalSummary, style := applyStateSummary(summary, m.workflowState, completedSteps, elapsed)
	return style.Render(finalSummary)
}
// ============================================================================
// NOM-STYLE RENDERING METHODS
// ============================================================================
// renderNOMStyle creates a NOM-style display with dependency tree nom.
func (m *ProgressModel) renderNOMStyle() string {
	var contentSections []string
	// Render dependency tree if available
	if m.dependencyTree != nil {
		treeSection := m.renderDependencyTree()
		if treeSection != "" {
			contentSections = append(contentSections, treeSection)
		}
	}
	// Assemble with NOM summary
	return m.assembleProgressSections(
		"🔄 NOM Workflow Execution",
		contentSections,
		m.renderNOMSummaryBar,
	)
}
// renderDependencyTree renders the activity dependency tree in NOM style.
func (m *ProgressModel) renderDependencyTree() string {
	if m.dependencyTree == nil {
		return ""
	}
	// Use the tree's built-in Render method with max height based on terminal
	maxHeight := m.height
	if maxHeight <= 0 {
		maxHeight = 20 // Default reasonable height
	}
	return m.dependencyTree.Render(maxHeight)
}
// renderNOMSummaryBar creates the NOM-style summary bar.
func (m *ProgressModel) renderNOMSummaryBar() string {
	running, completed, failed, pending := m.getActivityCounts()
	elapsed := time.Since(m.startTime)
	summary := buildNOMSummary(running, completed, failed, pending, elapsed)
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(lipgloss.Color("12"))
	return style.Render(summary)
}
// getActivityCounts returns counts of activities in each state.
func (m *ProgressModel) getActivityCounts() (running, completed, failed, pending int) {
	for _, activity := range m.activities {
		switch activity.Status {
		case nom.ActivityStatusRunning:
			running++
		case nom.ActivityStatusCompleted:
			completed++
		case nom.ActivityStatusFailed:
			failed++
		case nom.ActivityStatusPending, nom.ActivityStatusPaused:
			pending++
		}
	}
	return running, completed, failed, pending
}
