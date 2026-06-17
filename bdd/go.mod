module github.com/larsartmann/go-output/bdd

go 1.26.3

require (
	github.com/larsartmann/go-output v0.11.0
	github.com/larsartmann/go-output/delimited v0.11.0
	github.com/larsartmann/go-output/serialization v0.11.0
	github.com/onsi/ginkgo/v2 v2.23.4
	github.com/onsi/gomega v1.36.3
)

require (
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output/enum v0.11.0 // indirect
	github.com/pelletier/go-toml/v2 v2.4.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/tools v0.46.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/enum => ../enum
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/serialization => ../serialization
	github.com/larsartmann/go-output/table => ../table
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)
