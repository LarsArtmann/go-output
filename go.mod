module github.com/larsartmann/go-output

go 1.26.3

require (
	github.com/larsartmann/go-branded-id v0.3.1
	github.com/larsartmann/go-output/delimited v0.13.0
	github.com/larsartmann/go-output/enum v0.17.1
	github.com/larsartmann/go-output/envdetect v0.17.1
	github.com/larsartmann/go-output/testhelpers v0.13.0
	golang.org/x/term v0.44.0
)

replace (
	github.com/larsartmann/go-output => ./
	github.com/larsartmann/go-output/delimited => ./delimited
	github.com/larsartmann/go-output/enum => ./enum
	github.com/larsartmann/go-output/envdetect => ./envdetect
	github.com/larsartmann/go-output/serialization => ./serialization
	github.com/larsartmann/go-output/testhelpers => ./testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ./testhelpers/graphtest
)

require (
	github.com/larsartmann/go-output/markdown v0.13.0
	github.com/larsartmann/go-output/tree v0.13.0
	golang.org/x/sys v0.46.0 // indirect
)

replace github.com/larsartmann/go-output/markdown => ./markdown

replace github.com/larsartmann/go-output/tree => ./tree
