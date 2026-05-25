module github.com/larsartmann/go-output/graph

go 1.26.2

require (
	github.com/larsartmann/go-output v0.0.0
	github.com/larsartmann/go-output/escape v0.0.0
	github.com/larsartmann/go-output/testhelpers v0.0.0
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)

require (
	github.com/larsartmann/go-branded-id v0.1.0 // indirect
	github.com/larsartmann/go-output/enum v0.0.0 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/term v0.43.0 // indirect
)
