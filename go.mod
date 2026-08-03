module github.com/larsartmann/go-output

go 1.26.5

require (
	github.com/larsartmann/go-branded-id v0.5.1
	github.com/larsartmann/go-output/testhelpers v0.36.0
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

retract (
	// Bogus tags: pointed at a stale June commit, never real releases. Deleted from git; retracted here to poison proxy cache.
	v0.33.0
	// Same incident as v0.33.0 — bogus tag on stale commit, deleted from git; retracted to poison proxy cache.
	v0.32.1
)

replace github.com/larsartmann/go-output/markdown => ./markdown

replace github.com/larsartmann/go-output/tree => ./tree
