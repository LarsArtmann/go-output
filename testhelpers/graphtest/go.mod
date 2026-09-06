module github.com/larsartmann/go-output/testhelpers/graphtest

go 1.26.7

require github.com/larsartmann/go-output v0.38.0

replace github.com/larsartmann/go-output => ../..

require (
	github.com/larsartmann/go-branded-id v0.5.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
)

replace github.com/larsartmann/go-output/escape => ../../escape

replace github.com/larsartmann/go-output/markdown => ../../markdown

replace github.com/larsartmann/go-output/tree => ../../tree
