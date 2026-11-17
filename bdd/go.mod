module github.com/larsartmann/go-output/bdd

go 1.27

require (
	github.com/larsartmann/go-output v0.38.5
	github.com/larsartmann/go-output/delimited v0.38.5
	github.com/larsartmann/go-output/markdown v0.38.5
	github.com/larsartmann/go-output/serialization v0.38.5
	github.com/larsartmann/go-output/tree v0.38.5
	github.com/onsi/ginkgo/v2 v2.34.0
	github.com/onsi/gomega v1.44.0
)

require (
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/gkampitakis/go-snaps v0.5.25 // indirect
	github.com/go-faster/errors v0.8.0 // indirect
	github.com/go-faster/jx v1.2.0 // indirect
	github.com/go-faster/yaml v0.4.6 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20261105003125-4bf881d7fb0b // indirect
	github.com/larsartmann/go-branded-id v0.7.0 // indirect
	github.com/larsartmann/go-output/escape v0.38.5 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/rogpeppe/go-internal v1.16.0 // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	github.com/tidwall/gjson v1.20.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e // indirect
	golang.org/x/mod v0.42.0 // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/term v0.47.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/tools v0.52.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/serialization => ../serialization
	github.com/larsartmann/go-output/table => ../table
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)

replace github.com/larsartmann/go-output/escape => ../escape

replace github.com/larsartmann/go-output/markdown => ../markdown

replace github.com/larsartmann/go-output/tree => ../tree
