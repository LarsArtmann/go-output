// Package tui provides text-based user interface components for universal-workflow.
// This file contains shared rendering helpers extracted to eliminate duplication between
// NOM-style and universal-style progress rendering.
package tui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output/nom"
)

// formatElapsedTime formats elapsed time for display.
func formatElapsedTime(elapsed time.Duration) string {
	return fmt.Sprintf("%.1fs", elapsed.Seconds())
}

// createSummaryStyle creates the base style for summary bars.
func createSummaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(lipgloss.Color("12"))
}

// buildUniversalSummary builds a universal-style summary string.
func buildUniversalSummary(
	inProgress, completed int,
	elapsed time.Duration,
	progress float64,
) string {
	summary := fmt.Sprintf("📊 Activities: %d▶%d✓", inProgress, completed)
	summary += fmt.Sprintf(" | Progress: %.1f%%", progress)
	summary += " | ⏱️ " + formatElapsedTime(elapsed)

	return summary
}

// buildActivityCountsSummary builds a summary string with activity counts using NOM symbols.
func buildActivityCountsSummary(running, completed, failed, pending int) string {
	var parts []string
	if running > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", nom.SymbolRunning, running))
	}

	if completed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", nom.SymbolCompleted, completed))
	}

	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", nom.SymbolFailed, failed))
	}

	if pending > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", nom.SymbolPaused, pending))
	}

	return strings.Join(parts, " ")
}

// buildNOMSummary builds a NOM-style summary string.
func buildNOMSummary(running, completed, failed, pending int, elapsed time.Duration) string {
	summary := buildActivityCountsSummary(running, completed, failed, pending)
	if summary != "" {
		summary += " | "
	}

	summary += fmt.Sprintf("%s%s", nom.SymbolTiming, formatElapsedTime(elapsed))

	return summary
}

// getStateStyle returns the appropriate style for a workflow state.
func getStateStyle(state WorkflowState) (string, color.Color) {
	switch state {
	case WorkflowStateIdle:
		return "⏳ Workflow Idle | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", lipgloss.Color("8")
	case WorkflowStateRunning:
		return "", nil
	case WorkflowStateCompleted:
		return "✅ Workflow Complete: {completed}✓ | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", lipgloss.Color(
			"10",
		)
	case WorkflowStateErrored:
		return "❌ Workflow Error: {completed}✓ | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", lipgloss.Color(
			"9",
		)
	default:
		return "", lipgloss.Color("12")
	}
}

// applyStateSummary applies state-specific formatting to a summary.
func applyStateSummary(
	summary string,
	state WorkflowState,
	completed int,
	elapsed time.Duration,
) (string, lipgloss.Style) {
	baseStyle := createSummaryStyle()

	stateSummary, color := getStateStyle(state)
	if stateSummary != "" {
		// Replace placeholders
		stateSummary = strings.ReplaceAll(stateSummary, "{time}", formatElapsedTime(elapsed))
		stateSummary = strings.ReplaceAll(stateSummary, "{completed}", strconv.Itoa(completed))

		return stateSummary, baseStyle.Foreground(color)
	}

	return summary, baseStyle
}
