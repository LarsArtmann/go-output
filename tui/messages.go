package tui

import (
	"time"
)

// ============================================================================
// BUBBLE TEA MESSAGE TYPES
// ============================================================================
// updateType represents the type of progress update.
type updateType int

const (
	progressUpdate updateType = iota
	messageUpdate
)

// progressUpdateMsg represents updates to the progress display.
type progressUpdateMsg struct {
	Type     updateType
	Progress float64
	Message  string
}

// stepUpdateMsg carries step-based progress data (current/total counters + message).
// Processed exclusively on the TUI goroutine via model.Update.
type stepUpdateMsg struct {
	Current uint
	Total   uint
	Message string
}

// errorMsg carries an error to display and triggers transition to Errored state.
type errorMsg struct {
	Err error
}

// stateTransitionMsg requests the model to transition to a new workflow state.
type stateTransitionMsg struct {
	NewState workflowState
}

// tickMsg represents a timer tick for real-time updates.
type tickMsg time.Time
