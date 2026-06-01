module github.com/larsartmann/go-output/serialization

go 1.26.3

require (
	github.com/go-faster/yaml v0.4.6
	github.com/larsartmann/go-output v0.6.2
	github.com/larsartmann/go-output/testhelpers v0.6.2
	github.com/larsartmann/go-output/testhelpers/graphtest v0.6.2
	github.com/pelletier/go-toml/v2 v2.3.1
)

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/larsartmann/go-branded-id v0.3.0 // indirect
	github.com/larsartmann/go-output/enum v0.6.2 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/enum => ../enum
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)
