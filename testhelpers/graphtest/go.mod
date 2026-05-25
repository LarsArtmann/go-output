module github.com/larsartmann/go-output/testhelpers/graphtest

go 1.26.2

require (
	github.com/larsartmann/go-output v0.0.0
	github.com/larsartmann/go-branded-id v0.1.0
)

replace github.com/larsartmann/go-output => ../..

require golang.org/x/sys v0.44.0 // indirect
