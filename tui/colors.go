package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	colorInfo     color.Color = lipgloss.Color("12")
	colorTitle    color.Color = lipgloss.Color("39")
	colorSuccess  color.Color = lipgloss.Color("10")
	colorWarning  color.Color = lipgloss.Color("11")
	colorDim      color.Color = lipgloss.Color("8")
	colorError    color.Color = lipgloss.Color("9")
	colorCyan     color.Color = lipgloss.Color("14")
	colorSelectBG color.Color = lipgloss.Color("62")
	colorSelectFG color.Color = lipgloss.Color("230")
	colorHelpFG   color.Color = lipgloss.Color("15")
	colorHelpBG   color.Color = lipgloss.Color("0")
)

const (
	progressBarWidth  = 40
	minWidthThreshold = 80
	widthSubtraction  = 30
	chromeLines       = 8
	defaultTreeHeight = 20
	defaultHelpWidth  = 80
	defaultHelpHeight = 24
)
