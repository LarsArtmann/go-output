package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
		content = m.applyScrollViewport(content)
	}

	if m.showHelp {
		content = m.renderHelpOverlay()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion

	return v
}

// applyScrollViewport clips the content to the visible viewport based on scrollOffset.
func (m *ProgressModel) applyScrollViewport(content string) string {
	if m.height <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	totalLines := len(lines)

	if totalLines <= m.height {
		m.scrollOffset = 0
		return content
	}

	maxOffset := totalLines - m.height
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}

	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	end := min(m.scrollOffset+m.height, totalLines)

	visible := lines[m.scrollOffset:end]

	return strings.Join(visible, "\n")
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
		Foreground(colors.info)

	return messageStyle.Render(m.currentMessage)
}

// renderTitle creates a styled title header.
func (m *ProgressModel) renderTitle(text string) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(colors.title).
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

// stepIconAndStyle returns the icon and lipgloss style for a step based on its state.
func stepIconAndStyle(step progressStep) (string, lipgloss.Style) {
	if step.CompletedAt != nil {
		return "✅", lipgloss.NewStyle().Foreground(colors.success)
	}

	if step.isActive() {
		return "🔄", lipgloss.NewStyle().Foreground(colors.warning)
	}

	return "⏸️", lipgloss.NewStyle().Foreground(colors.dim)
}

// renderStep renders a single step with nh-style formatting.
func (m *ProgressModel) renderStep(step progressStep, isLast bool) string {
	// Choose prefix
	prefix := "├── "
	if isLast {
		prefix = "└── "
	}
	// Status icon and color
	icon, style := stepIconAndStyle(step)
	// Timing information like nh darwin switch
	var timing string

	if step.CompletedAt != nil {
		duration := step.CompletedAt.Sub(step.StartTime)
		timing = fmt.Sprintf(timingFormat, duration.Seconds())
	} else if step.isActive() {
		elapsed := time.Since(step.StartTime)
		timing = fmt.Sprintf(timingFormat, elapsed.Seconds())
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
	width := progressBarWidth
	if m.width > 0 && m.width < minWidthThreshold {
		width = m.width - widthSubtraction
	}

	if width < 1 {
		width = 1
	}

	filled := int((m.currentProgress / percentScale) * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	progressStyle := lipgloss.NewStyle().
		Foreground(colors.info)

	return progressStyle.Render(fmt.Sprintf("Progress: [%s] %.1f%%", bar, m.currentProgress))
}

// renderSummaryBar creates the summary like nh darwin switch.
func (m *ProgressModel) renderSummaryBar() string {
	// Calculate summary statistics
	var completedSteps, inProgressSteps int

	for _, step := range m.steps {
		if step.CompletedAt != nil {
			completedSteps++
		} else if step.isActive() {
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
