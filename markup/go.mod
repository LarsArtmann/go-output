module github.com/larsartmann/go-output/markup

go 1.26.3

require (
	github.com/larsartmann/go-output v0.6.3
	github.com/larsartmann/go-output/escape v0.6.3
	github.com/larsartmann/go-output/testhelpers v0.6.3
)

require (
	github.com/larsartmann/go-branded-id v0.3.0 // indirect
	github.com/larsartmann/go-output/enum v0.6.3 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/term v0.43.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/delimited => ../delimited
	github.com/larsartmann/go-output/escape => ../escape
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)
