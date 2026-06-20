module github.com/larsartmann/go-output/graph

go 1.26.3

require (
	github.com/larsartmann/go-output v0.16.0
	github.com/larsartmann/go-output/enum v0.13.0
	github.com/larsartmann/go-output/escape v0.13.0
	github.com/larsartmann/go-output/testhelpers v0.13.0
	github.com/larsartmann/go-output/testhelpers/graphtest v0.13.0
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)

require (
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output/envdetect v0.13.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace github.com/larsartmann/go-output/enum => ../enum

replace github.com/larsartmann/go-output/envdetect => ../envdetect

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
