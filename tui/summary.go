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

// formatElapsedTime formats elapsed time for display using nom.FormatDuration
// to ensure consistent duration formatting (ms/s/m/h) across nom and tui.
func formatElapsedTime(elapsed time.Duration) string {
	return nom.FormatDuration(elapsed)
}

// createSummaryStyle creates the base style for summary bars.
func createSummaryStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Foreground(colors.info)
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

// buildActivityCountsSummary builds a summary string with activity counts using
// NOM symbols. Delegates to nom.ActivityCounts.Summary() so formatting stays
// consistent with the inline renderer.
func buildActivityCountsSummary(counts nom.ActivityCounts) string {
	return counts.Summary()
}

// buildNOMSummary builds a NOM-style summary string.
func buildNOMSummary(counts nom.ActivityCounts, elapsed time.Duration) string {
	summary := buildActivityCountsSummary(counts)
	if summary != "" {
		summary += " | "
	}

	summary += fmt.Sprintf("%s%s", nom.SymbolTiming, formatElapsedTime(elapsed))

	return summary
}

// getStateStyle returns the appropriate style for a workflow state.
func getStateStyle(state workflowState) (string, color.Color) {
	switch state {
	case workflowStateIdle:
		return "⏳ Workflow Idle | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", colors.dim
	case workflowStateRunning:
		return "", nil
	case workflowStateCompleted:
		return "✅ Workflow Complete: {completed}✓ | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", colors.success
	case workflowStateErrored:
		return "❌ Workflow Error: {completed}✓ | ⏱️ {time}s | Press 'q' or Ctrl+C to exit", colors.err
	default:
		return "", colors.info
	}
}

// applyStateSummary applies state-specific formatting to a summary.
func applyStateSummary(
	summary string,
	state workflowState,
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
