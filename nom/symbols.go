package nom

import "charm.land/lipgloss/v2"

// ============================================================================
// NOM-STYLE SYMBOL CONSTANTS
// ============================================================================.
const (
	// SymbolRunning represents an activity currently executing.
	SymbolRunning = "⏵"
	// SymbolCompleted represents a successfully completed activity.
	SymbolCompleted = "✔"
	// SymbolPaused represents a paused or waiting activity.
	SymbolPaused = "⏸"
	// SymbolFailed represents a failed activity.
	SymbolFailed = "⚠"
	// SymbolDownload represents a download operation.
	SymbolDownload = "↓"
	// SymbolUpload represents an upload operation.
	SymbolUpload = "↑"
	// SymbolTiming represents timing information.
	SymbolTiming = "⏱️"
	// SymbolAverage represents average duration.
	SymbolAverage = "∅"
	// SymbolTotal represents total count/summary.
	SymbolTotal = "∑"
	// SymbolPhase represents a phase/group node in the tree.
	SymbolPhase = "◈"
)

// ============================================================================
// COLOR MAPPING FOR ACTIVITY STATES
// ============================================================================
// ColorRunning is the color for running activities (yellow, matching NOM).
var ColorRunning = lipgloss.Color("11")

// ColorCompleted is the color for completed activities (green).
var ColorCompleted = lipgloss.Color("10")

// ColorPaused is the color for paused activities (gray).
var ColorPaused = lipgloss.Color("8")

// ColorFailed is the color for failed activities (red).
var ColorFailed = lipgloss.Color("9")

// ColorWarning is the color for warnings (yellow).
var ColorWarning = lipgloss.Color("11")

// ColorInfo is the color for information (cyan).
var ColorInfo = lipgloss.Color("14")

// ColorPhase is the color for phase/group nodes (magenta).
var ColorPhase = lipgloss.Color("13")
