module github.com/larsartmann/go-output/testhelpers/graphtest

go 1.26.4

require github.com/larsartmann/go-output v0.0.0-00010101000000-000000000000

replace github.com/larsartmann/go-output => ../..

require (
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace github.com/larsartmann/go-output/escape => ../../escape

replace github.com/larsartmann/go-output/markdown => ../../markdown

replace github.com/larsartmann/go-output/tree => ../../tree
