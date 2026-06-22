module github.com/larsartmann/go-output/tui

go 1.26.3

require (
	charm.land/bubbletea/v2 v2.0.7
	charm.land/lipgloss/v2 v2.0.4
	github.com/charmbracelet/x/ansi v0.11.7
	github.com/larsartmann/go-output/nom v0.17.0
)

require (
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260622092850-f39628c8a989 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output v0.17.0 // indirect
	github.com/larsartmann/go-output/enum v0.17.0 // indirect
	github.com/larsartmann/go-output/envdetect v0.17.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/nom => ../nom
)

replace github.com/larsartmann/go-output/enum => ../enum

replace github.com/larsartmann/go-output/envdetect => ../envdetect

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
