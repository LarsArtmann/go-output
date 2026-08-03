module github.com/larsartmann/go-output/markdown

go 1.26.5

require (
	github.com/larsartmann/go-output v0.36.0
	github.com/larsartmann/go-output/escape v0.36.0
	github.com/larsartmann/go-output/testhelpers v0.36.0
)

require (
	github.com/larsartmann/go-branded-id v0.5.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)

replace github.com/larsartmann/go-output/tree => ../tree
