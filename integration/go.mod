module github.com/larsartmann/go-output/integration

go 1.26.4

require (
	github.com/go-faster/yaml v0.4.6
	github.com/larsartmann/go-output v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/d2 v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/delimited v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/graph v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/markup v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/nom v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/plantuml v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/serialization v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/table v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/testhelpers v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/tui v0.0.0-00010101000000-000000000000
)

require (
	charm.land/bubbletea/v2 v2.0.8 // indirect
	charm.land/lipgloss/v2 v2.0.5 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260720091822-7cc6674724ac // indirect
	github.com/charmbracelet/x/ansi v0.11.7 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/go-faster/errors v0.8.0 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.2 // indirect
	github.com/larsartmann/go-output/escape v0.0.0-00010101000000-000000000000 // indirect
	github.com/larsartmann/go-output/markdown v0.0.0-00010101000000-000000000000
	github.com/larsartmann/go-output/tree v0.0.0-00010101000000-000000000000
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/d2 => ../d2
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/graph => ../graph
	github.com/larsartmann/go-output/markup => ../markup
	github.com/larsartmann/go-output/nom => ../nom
	github.com/larsartmann/go-output/plantuml => ../plantuml
	github.com/larsartmann/go-output/serialization => ../serialization
	github.com/larsartmann/go-output/table => ../table
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
	github.com/larsartmann/go-output/tui => ../tui
)

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
