package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output/nom"
)

// terminalColors groups all ANSI colors used by the TUI into a single
// cohesive, swappable theme. Immutable after initialization.
//
// The 4 semantic colors (success, warning, err, dim) delegate to nom.Colors
// so there is exactly ONE source of truth for the activity-state palette.
// The remaining fields are TUI-specific (title, selection, help) and have no
// nom equivalent. See split-brain finding M1.
type terminalColors struct {
	info     color.Color
	title    color.Color
	success  color.Color
	warning  color.Color
	dim      color.Color
	err      color.Color
	selectBG color.Color
	selectFG color.Color
	helpFG   color.Color
	helpBG   color.Color
}

// colors is the default terminal color theme, set once at package init.
// success/warning/err/dim delegate to nom.Colors (single source of truth);
// the rest are TUI-specific values with no nom equivalent.
//
//nolint:gochecknoglobals // immutable theme configuration (ANSI color constants)
var colors = terminalColors{
	info:     lipgloss.Color("12"),
	title:    lipgloss.Color("39"),
	success:  nom.Colors.Completed, // ANSI 10 — green
	warning:  nom.Colors.Running,   // ANSI 11 — yellow
	dim:      nom.Colors.Pending,   // ANSI 8 — gray
	err:      nom.Colors.Failed,    // ANSI 9 — red
	selectBG: lipgloss.Color("62"),
	selectFG: lipgloss.Color("230"),
	helpFG:   lipgloss.Color("15"),
	helpBG:   lipgloss.Color("0"),
}

const (
	progressBarWidth  = 40
	minWidthThreshold = 80
	widthSubtraction  = 30
	chromeLines       = 8
	defaultTreeHeight = 20
	defaultHelpWidth  = 80
	defaultHelpHeight = 24
)
