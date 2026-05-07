module github.com/larsartmann/go-output

go 1.26.2

require (
	github.com/go-faster/yaml v0.4.6
	github.com/larsartmann/go-output/enum v0.0.0
	github.com/larsartmann/go-output/escape v0.0.0
	golang.org/x/term v0.42.0
)

replace (
	github.com/larsartmann/go-output/cmdguard => ./cmdguard
	github.com/larsartmann/go-output/enum => ./enum
	github.com/larsartmann/go-output/escape => ./escape
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/sys v0.43.0 // indirect
)
