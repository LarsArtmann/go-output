module github.com/larsartmann/go-output/nom

go 1.26.4

require (
	charm.land/lipgloss/v2 v2.0.5
	github.com/charmbracelet/colorprofile v0.4.3
	github.com/charmbracelet/x/ansi v0.11.7
	github.com/charmbracelet/x/exp/golden v0.0.0-20260615092313-b57e5e6d29bb
	github.com/charmbracelet/x/vt v0.0.0-20260629091435-9c70f75e26a4
	github.com/larsartmann/go-output v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/testhelpers v0.23.2
	github.com/onsi/gomega v1.42.1
	golang.org/x/term v0.44.0
)

require (
	github.com/aymanbagabas/go-udiff v0.4.1 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260703014108-f5a850f9c2b7 // indirect
	github.com/charmbracelet/x/exp/ordered v0.1.0 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.24 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
)

replace github.com/larsartmann/go-output => ../

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
