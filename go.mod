module github.com/larsartmann/go-output

go 1.26.4

require (
	github.com/larsartmann/go-branded-id v0.3.2
	github.com/larsartmann/go-output/testhelpers v0.31.1
	golang.org/x/term v0.45.0
)

replace (
	github.com/larsartmann/go-output => ./
	github.com/larsartmann/go-output/delimited => ./delimited
	github.com/larsartmann/go-output/serialization => ./serialization
	github.com/larsartmann/go-output/testhelpers => ./testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ./testhelpers/graphtest
)

require golang.org/x/sys v0.47.0 // indirect

replace github.com/larsartmann/go-output/markdown => ./markdown

replace github.com/larsartmann/go-output/tree => ./tree
