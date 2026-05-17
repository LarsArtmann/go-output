module github.com/larsartmann/go-output/cmdguard

go 1.26.2

require github.com/larsartmann/go-output v0.0.0

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/larsartmann/go-branded-id v0.1.0 // indirect
	github.com/larsartmann/go-output/enum v0.0.0 // indirect
	github.com/larsartmann/go-output/escape v0.0.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/term v0.42.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/enum => ../enum
	github.com/larsartmann/go-output/escape => ../escape
)
