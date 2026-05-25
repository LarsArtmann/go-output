module github.com/larsartmann/go-output/plantuml

go 1.26.2

require github.com/larsartmann/go-output v0.0.0

require (
	github.com/larsartmann/go-branded-id v0.1.0 // indirect
	github.com/larsartmann/go-output/enum v0.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.3.1 // indirect
	golang.org/x/sys v0.44.0 // indirect
	golang.org/x/term v0.43.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/enum => ../enum
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/serialization => ../serialization
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)
