package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// terminalColors groups all ANSI colors used by the TUI into a single
// cohesive, swappable theme. Immutable after initialization.
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
//
//nolint:gochecknoglobals // immutable theme configuration (ANSI color constants)
var colors = terminalColors{
	info:     lipgloss.Color("12"),
	title:    lipgloss.Color("39"),
	success:  lipgloss.Color("10"),
	warning:  lipgloss.Color("11"),
	dim:      lipgloss.Color("8"),
	err:      lipgloss.Color("9"),
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
