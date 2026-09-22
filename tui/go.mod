module github.com/larsartmann/go-output/tui

go 1.26.0

require (
	charm.land/bubbletea/v2 v2.0.9
	charm.land/lipgloss/v2 v2.0.6
	github.com/charmbracelet/x/ansi v0.11.8
	github.com/charmbracelet/x/exp/teatest/v2 v2.0.0-20260629091435-9c70f75e26a4
	github.com/charmbracelet/x/vt v0.0.0-20260629091435-9c70f75e26a4
	github.com/larsartmann/go-output/nom v0.38.1
)

require (
	github.com/aymanbagabas/go-udiff v0.4.1 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260910203606-6c9e17dc7a16 // indirect
	github.com/charmbracelet/x/exp/golden v0.0.0-20260705004817-2cc9a8fe1146 // indirect
	github.com/charmbracelet/x/exp/ordered v0.1.0 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/larsartmann/go-branded-id v0.6.0 // indirect
	github.com/larsartmann/go-output v0.38.1 // indirect
	github.com/larsartmann/go-output/escape v0.38.1
	github.com/larsartmann/go-output/markdown v0.38.1
	github.com/larsartmann/go-output/testhelpers v0.38.1
	github.com/larsartmann/go-output/tree v0.38.1
	github.com/larsartmann/go-output/tui v0.38.1
	github.com/lucasb-eyer/go-colorful v1.4.1 // indirect
	github.com/mattn/go-runewidth v0.0.30 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v1.2.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/nom => ../nom
)

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/testhelpers => ../testhelpers

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
