module github.com/larsartmann/go-output/plantuml

go 1.26.3

require (
	github.com/larsartmann/go-output v0.12.0
	github.com/larsartmann/go-output/escape v0.12.0
	github.com/larsartmann/go-output/testhelpers v0.12.0
	github.com/larsartmann/go-output/testhelpers/graphtest v0.12.0
)

require (
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output/enum v0.12.0 // indirect
	github.com/larsartmann/go-output/envdetect v0.12.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
	github.com/larsartmann/go-output/testhelpers/graphtest => ../testhelpers/graphtest
)

replace github.com/larsartmann/go-output/enum => ../enum

replace github.com/larsartmann/go-output/envdetect => ../envdetect

replace github.com/larsartmann/go-output/escape => ../escape
