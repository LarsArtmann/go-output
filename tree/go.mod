module github.com/larsartmann/go-output/tree

go 1.26.4

require (
	github.com/charmbracelet/x/exp/golden v0.0.0-20260705004817-2cc9a8fe1146
	github.com/larsartmann/go-output v0.30.2
	github.com/larsartmann/go-output/testhelpers v0.30.2
)

require (
	github.com/aymanbagabas/go-udiff v0.4.1 // indirect
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)

replace github.com/larsartmann/go-output/markdown => ../markdown
