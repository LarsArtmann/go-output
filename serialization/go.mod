module github.com/larsartmann/go-output/serialization

go 1.26.3

require (
	github.com/go-faster/yaml v0.4.6
	github.com/larsartmann/go-output v0.17.2
	github.com/larsartmann/go-output/testhelpers v0.13.0
	github.com/larsartmann/go-output/testhelpers/graphtest v0.13.0
	github.com/pelletier/go-toml/v2 v2.4.1
)

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
