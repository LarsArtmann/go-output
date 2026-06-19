package tui

// Display Format Constants - Standardized display format strings.
const (
	// timingFormatWithIcon is the TUI-specific timing format that bakes in the
	// ⏱️ icon. This intentionally differs from nom.TimingFormat (which is just
	// "%.1fs") because the TUI step-display integrates the icon into the format
	// string, while nom composes icon + duration separately.
	timingFormatWithIcon = "⏱️ %.1fs"
)

// Message Constants - Standardized UI message strings.
const (
	// msgNoActivitiesToDisplay mirrors nom.msgNoActivitiesToDisplay — identical
	// string, kept separate because tui renders its own view. See split-brain M3.
	msgNoActivitiesToDisplay = "No activities to display"
)
