module github.com/larsartmann/go-output/testhelpers/graphtest

go 1.26

require github.com/larsartmann/go-output v0.38.1

replace github.com/larsartmann/go-output => ../..

require (
	github.com/larsartmann/go-branded-id v0.6.0 // indirect
	github.com/larsartmann/go-output/escape v0.38.1
	github.com/larsartmann/go-output/markdown v0.38.1
	github.com/larsartmann/go-output/testhelpers v0.38.1
	github.com/larsartmann/go-output/testhelpers/graphtest v0.38.1
	github.com/larsartmann/go-output/tree v0.38.1
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

replace github.com/larsartmann/go-output/escape => ../../escape

replace github.com/larsartmann/go-output/testhelpers => ../

replace github.com/larsartmann/go-output/markdown => ../../markdown

replace github.com/larsartmann/go-output/tree => ../../tree
