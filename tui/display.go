package tui

// DisplayMode controls which visualization style the TUI renders.
type DisplayMode int

const (
	DisplayModeUniversal DisplayMode = iota
	DisplayModeNOM
)
