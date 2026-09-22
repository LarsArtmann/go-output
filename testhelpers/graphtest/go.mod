module github.com/larsartmann/go-output/testhelpers/graphtest

go 1.27.1

require github.com/larsartmann/go-output v0.38.1

require (
	github.com/larsartmann/go-branded-id v0.6.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

replace github.com/larsartmann/go-output => ../..

replace github.com/larsartmann/go-output/escape => ../../escape

replace github.com/larsartmann/go-output/testhelpers => ../

replace github.com/larsartmann/go-output/markdown => ../../markdown

replace github.com/larsartmann/go-output/tree => ../../tree
