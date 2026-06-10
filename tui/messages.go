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

// TickMsg represents a timer tick for real-time updates.
type TickMsg time.Time

// CancelMsg signals the TUI to shut down gracefully.
type CancelMsg struct{}
