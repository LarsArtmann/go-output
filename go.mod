module github.com/larsartmann/go-output

go 1.26.3

require (
	github.com/larsartmann/go-branded-id v0.3.1
	github.com/larsartmann/go-output/delimited v0.12.0
	github.com/larsartmann/go-output/testhelpers v0.12.0
	golang.org/x/term v0.44.0
)

replace (
	github.com/larsartmann/go-output => ./
	github.com/larsartmann/go-output/delimited => ./delimited
	github.com/larsartmann/go-output/serialization => ./serialization
	github.com/larsartmann/go-output/testhelpers => ./testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ./testhelpers/graphtest
)

require golang.org/x/sys v0.46.0 // indirect
