package tui

// Display Format Constants - Standardized display format strings.
const (
	// timingFormat is the TUI-specific timing format for step durations. No
	// icon is baked in — the variation-selector emoji (⏱️, U+23F1 U+FE0F)
	// renders as a phantom column on many terminals (especially over SSH),
	// making the timing look "half there". A plain duration matches NOM.
	timingFormat = "%.1fs"
)

// Message Constants - Standardized UI message strings.
const (
	// msgNoActivitiesToDisplay mirrors nom.msgNoActivitiesToDisplay — identical
	// string, kept separate because tui renders its own view. See split-brain M3.
	msgNoActivitiesToDisplay = "No activities to display"
)
