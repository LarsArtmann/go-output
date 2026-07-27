module github.com/larsartmann/go-output/plantuml

go 1.26.5

require (
	github.com/charmbracelet/x/exp/golden v0.0.0-20260705004817-2cc9a8fe1146
	github.com/larsartmann/go-output v0.34.0
	github.com/larsartmann/go-output/escape v0.34.0
	github.com/larsartmann/go-output/testhelpers v0.34.0
	github.com/larsartmann/go-output/testhelpers/graphtest v0.34.0
)

require (
	github.com/aymanbagabas/go-udiff v0.4.1 // indirect
	github.com/larsartmann/go-branded-id v0.4.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
