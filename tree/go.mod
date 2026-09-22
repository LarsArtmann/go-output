module github.com/larsartmann/go-output/tree

go 1.27.1

require (
	github.com/charmbracelet/x/exp/golden v0.0.0-20260705004817-2cc9a8fe1146
	github.com/larsartmann/go-output v0.38.1
	github.com/larsartmann/go-output/escape v0.38.1
	github.com/larsartmann/go-output/testhelpers v0.38.1
)

require (
	github.com/aymanbagabas/go-udiff v0.4.1 // indirect
	github.com/larsartmann/go-branded-id v0.6.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)

replace github.com/larsartmann/go-output/markdown => ../markdown
