package nom

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// ============================================================================
// NOM-STYLE SYMBOL CONSTANTS
// ============================================================================.

// Symbol is a typed string constant representing a visual status indicator.
// Typing prevents accidental mixing with arbitrary strings.
type Symbol string

// String returns the underlying string value of the symbol.
func (s Symbol) String() string { return string(s) }

const (
	// SymbolRunning represents an activity currently executing.
	SymbolRunning Symbol = "⏵"
	// SymbolCompleted represents a successfully completed activity.
	SymbolCompleted Symbol = "✔"
	// SymbolPending represents a queued / not-yet-started activity.
	SymbolPending Symbol = "○"
	// SymbolFailed represents a failed activity.
	SymbolFailed Symbol = "⚠"
	// SymbolDownload represents a download operation.
	SymbolDownload Symbol = "↓"
	// SymbolUpload represents an upload operation.
	SymbolUpload Symbol = "↑"
	// SymbolAverage represents average duration.
	SymbolAverage Symbol = "∅"
	// SymbolTotal represents total count/summary.
	SymbolTotal Symbol = "∑"
	// SymbolPhase represents a phase/group node in the tree.
	SymbolPhase Symbol = "◈"
	// SymbolRetrying represents a retry indicator.
	SymbolRetrying Symbol = "⟳"
	// SymbolProgress represents a sub-step progress indicator.
	SymbolProgress Symbol = "→"
)

// ============================================================================
// COLOR MAPPING FOR ACTIVITY STATES
// ============================================================================
// SemanticColors holds the ANSI color codes for activity states and phases.
// Immutable after initialization — do not mutate at runtime.
//
//nolint:gochecknoglobals // immutable theme configuration
type SemanticColors struct {
	Running   color.Color
	Completed color.Color
	Pending   color.Color
	Failed    color.Color
	Info      color.Color
	Phase     color.Color
}

// Colors is the default color theme for activity states. Mirrors the 4 semantic
// colors used by tui/colors.go (success≈Completed, warning≈Running,
// err≈Failed, dim≈Pending). See split-brain M1 in SPLIT-BRAIN.html.
//
//nolint:gochecknoglobals // immutable theme configuration
var Colors = SemanticColors{
	Running:   lipgloss.Color("11"),
	Completed: lipgloss.Color("10"),
	Pending:   lipgloss.Color("8"),
	Failed:    lipgloss.Color("9"),
	Info:      lipgloss.Color("14"),
	Phase:     lipgloss.Color("13"),
}
