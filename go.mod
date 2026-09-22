module github.com/larsartmann/go-output

go 1.26

require (
	github.com/larsartmann/go-branded-id v0.6.0
	github.com/larsartmann/go-output/testhelpers v0.38.1
	golang.org/x/term v0.46.0
)

replace (
	github.com/larsartmann/go-output => ./
	github.com/larsartmann/go-output/delimited => ./delimited
	github.com/larsartmann/go-output/serialization => ./serialization
	github.com/larsartmann/go-output/testhelpers => ./testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ./testhelpers/graphtest
)

require (
	github.com/larsartmann/go-output v0.38.1
	github.com/larsartmann/go-output/delimited v0.38.1
	github.com/larsartmann/go-output/markdown v0.38.1
	github.com/larsartmann/go-output/markup v0.38.1
	github.com/larsartmann/go-output/serialization v0.38.1
	github.com/larsartmann/go-output/testhelpers/graphtest v0.38.1
	github.com/larsartmann/go-output/tree v0.38.1
	golang.org/x/sys v0.48.0 // indirect
)

retract (
	// Stale tag drift: sibling dep versions were misaligned at tag time; v0.35.0 realigned them the same day.
	v0.34.0
	// Bogus tags: pointed at a stale June commit, never real releases. Deleted from git; retracted here to poison proxy cache.
	v0.33.0
	// Same incident as v0.33.0 — bogus tag on stale commit, deleted from git; retracted to poison proxy cache.
	v0.32.1
)

replace github.com/larsartmann/go-output/markdown => ./markdown

replace github.com/larsartmann/go-output/markup => ./markup

replace github.com/larsartmann/go-output/tree => ./tree
