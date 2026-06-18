package tui

import (
	"time"
)

// ============================================================================
// BUBBLE TEA MESSAGE TYPES
// ============================================================================
// UpdateType represents the type of progress update.
type UpdateType int

const (
	ProgressUpdate UpdateType = iota
	MessageUpdate
	StepUpdate
)

// ProgressUpdateMsg represents updates to the progress display.
type ProgressUpdateMsg struct {
	Type     UpdateType
	Progress float64
	Message  string
	Current  uint
	Total    uint
}

// StepUpdateMsg carries step-based progress data (current/total counters + message).
// Processed exclusively on the TUI goroutine via model.Update.
type StepUpdateMsg struct {
	Current uint
	Total   uint
	Message string
}

// ErrorMsg carries an error to display and triggers transition to Errored state.
type ErrorMsg struct {
	Err error
}

// StateTransitionMsg requests the model to transition to a new workflow state.
type StateTransitionMsg struct {
	NewState WorkflowState
}

// TickMsg represents a timer tick for real-time updates.
type TickMsg time.Time

// CancelMsg signals the TUI to shut down gracefully.
type CancelMsg struct{}
