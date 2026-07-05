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
	// SymbolDeps represents extra (non-display-parent) dependencies beneath
	// a multi-dependency node. Only shown when showExtraDeps is enabled.
	SymbolDeps Symbol = "↳"
	// SymbolCritical marks nodes on the longest estimated-time path through
	// the DAG (the critical path). Rendered as a prefix to the activity label.
	SymbolCritical Symbol = "◆"
	// SymbolConvergence marks nodes with multiple incoming dependencies
	// (fan-in points), making DAG join points visible in tree mode.
	SymbolConvergence Symbol = "◇"
	// SymbolBlocked introduces the blockage sub-line for pending nodes that
	// have incomplete dependencies.
	SymbolBlocked Symbol = "⊘"
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
	// Fallback is used for unknown/custom statuses that don't match a
	// specific semantic slot. (Was named Info in previous versions.)
	Fallback color.Color
	Phase    color.Color
}

// Colors is the DEFAULT color palette for activity states. Used as the
// fallback when no theme is explicitly set (ThemeDefault embeds these values).
//
// At render time, SnapshotActivities() resolves colors via the subscriber's
// theme, so a custom theme always overrides these defaults. This global is
// NOT a parallel source of truth for rendering — it is the initial value that
// themes build on. See theme.go for the Theme system.
//
//nolint:gochecknoglobals // immutable default theme configuration
var Colors = SemanticColors{
	Running:   lipgloss.Color("11"),
	Completed: lipgloss.Color("10"),
	Pending:   lipgloss.Color("8"),
	Failed:    lipgloss.Color("9"),
	Fallback:  lipgloss.Color("14"),
	Phase:     lipgloss.Color("13"),
}
