package nom

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// ============================================================================
// THEME SYSTEM
// ============================================================================
// Theme bundles all visual decisions (semantic colors, status symbols, and
// category tints) into a single swappable configuration. A theme is immutable
// after construction and is stored on the NOMSubscriber so that snapshot
// time is the single point where status/state is resolved into concrete ANSI
// colors and symbols. Renderers consume ActivitySnapshot and do not need to
// know the theme directly.
type Theme struct {
	// Colors are the semantic ANSI colors used for activity states and phases.
	Colors SemanticColors
	// Symbols overrides the display symbol for specific statuses. A missing
	// entry falls back to the symbol registered for that status.
	Symbols map[ActivityStatus]Symbol
	// CategoryColors provides an optional tint for user-defined activity
	// categories. A missing entry means the category renders without a tint.
	CategoryColors map[ActivityCategory]color.Color
}

// StatusColor returns the theme color for a task status. Unknown statuses
// resolve to the theme's Info color.
func (t Theme) StatusColor(status ActivityStatus) color.Color {
	if def, ok := LookupStatus(status); ok {
		switch def.Name {
		case "pending":
			return t.Colors.Pending
		case "running":
			return t.Colors.Running
		case "completed":
			return t.Colors.Completed
		case "failed":
			return t.Colors.Failed
		}
	}

	return t.Colors.Fallback
}

// ColorFor returns the theme color for an activity, choosing the phase color
// for phase nodes and the status color for task nodes.
func (t Theme) ColorFor(kind ActivityKind, status ActivityStatus) color.Color {
	if kind.IsPhase() {
		return t.Colors.Phase
	}

	return t.StatusColor(status)
}

// SymbolFor returns the theme symbol override for a status, falling back to
// the provided default (typically the symbol registered for the status).
func (t Theme) SymbolFor(status ActivityStatus, fallback Symbol) Symbol {
	if s, ok := t.Symbols[status]; ok {
		return s
	}

	return fallback
}

// CategoryColor returns the tint color for a category, or nil if the theme
// does not define one.
func (t Theme) CategoryColor(cat ActivityCategory) color.Color {
	if t.CategoryColors == nil {
		return nil
	}

	return t.CategoryColors[cat]
}

// ThemeDefault is the current NOM-style color scheme, matching the historic
// Colors global. It is the implicit theme when none is supplied.
//
//nolint:gochecknoglobals // immutable exported theme presets
var ThemeDefault = Theme{
	Colors: Colors,
	Symbols: map[ActivityStatus]Symbol{
		ActivityStatusPending:   SymbolPending,
		ActivityStatusRunning:   SymbolRunning,
		ActivityStatusCompleted: SymbolCompleted,
		ActivityStatusFailed:    SymbolFailed,
	},
	CategoryColors: map[ActivityCategory]color.Color{},
}

// ThemeDracula is a dark, high-contrast theme based on the Dracula palette.
//
//nolint:gochecknoglobals // immutable exported theme presets
var ThemeDracula = Theme{
	Colors: SemanticColors{
		Running:   lipgloss.Color("#50fa7b"),
		Completed: lipgloss.Color("#6272a4"),
		Pending:   lipgloss.Color("#6272a4"),
		Failed:    lipgloss.Color("#ff5555"),
		Fallback:  lipgloss.Color("#8be9fd"),
		Phase:     lipgloss.Color("#bd93f9"),
	},
	Symbols: map[ActivityStatus]Symbol{
		ActivityStatusPending:   SymbolPending,
		ActivityStatusRunning:   SymbolRunning,
		ActivityStatusCompleted: SymbolCompleted,
		ActivityStatusFailed:    SymbolFailed,
	},
	CategoryColors: map[ActivityCategory]color.Color{},
}

// ThemeNord is a polar-inspired low-saturation theme based on the Nord palette.
//
//nolint:gochecknoglobals // immutable exported theme presets
var ThemeNord = Theme{
	Colors: SemanticColors{
		Running:   lipgloss.Color("#a3be8c"),
		Completed: lipgloss.Color("#4c566a"),
		Pending:   lipgloss.Color("#4c566a"),
		Failed:    lipgloss.Color("#bf616a"),
		Fallback:  lipgloss.Color("#88c0d0"),
		Phase:     lipgloss.Color("#b48ead"),
	},
	Symbols: map[ActivityStatus]Symbol{
		ActivityStatusPending:   SymbolPending,
		ActivityStatusRunning:   SymbolRunning,
		ActivityStatusCompleted: SymbolCompleted,
		ActivityStatusFailed:    SymbolFailed,
	},
	CategoryColors: map[ActivityCategory]color.Color{},
}

// ThemeMonochrome is a grayscale theme suitable for printers, plain-text
// capture, or users who prefer no chroma.
//
//nolint:gochecknoglobals // immutable exported theme presets
var ThemeMonochrome = Theme{
	Colors: SemanticColors{
		Running:   lipgloss.Color("#ffffff"),
		Completed: lipgloss.Color("#a0a0a0"),
		Pending:   lipgloss.Color("#808080"),
		Failed:    lipgloss.Color("#ffffff"),
		Fallback:  lipgloss.Color("#c0c0c0"),
		Phase:     lipgloss.Color("#d0d0d0"),
	},
	Symbols: map[ActivityStatus]Symbol{
		ActivityStatusPending:   "○",
		ActivityStatusRunning:   "›",
		ActivityStatusCompleted: "✓",
		ActivityStatusFailed:    "×",
	},
	CategoryColors: map[ActivityCategory]color.Color{},
}

// ThemeHighContrast uses bold, maximally distinct ANSI colors for
// accessibility and bright terminals.
//
//nolint:gochecknoglobals // immutable exported theme presets
var ThemeHighContrast = Theme{
	Colors: SemanticColors{
		Running:   lipgloss.Color("#00ff00"),
		Completed: lipgloss.Color("#808080"),
		Pending:   lipgloss.Color("#ffff00"),
		Failed:    lipgloss.Color("#ff0000"),
		Fallback:  lipgloss.Color("#00ffff"),
		Phase:     lipgloss.Color("#ff00ff"),
	},
	Symbols: map[ActivityStatus]Symbol{
		ActivityStatusPending:   SymbolPending,
		ActivityStatusRunning:   SymbolRunning,
		ActivityStatusCompleted: SymbolCompleted,
		ActivityStatusFailed:    SymbolFailed,
	},
	CategoryColors: map[ActivityCategory]color.Color{},
}
