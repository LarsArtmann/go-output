module github.com/larsartmann/go-output/delimited

go 1.26.3

require (
	github.com/larsartmann/go-output v0.10.1
	github.com/larsartmann/go-output/testhelpers v0.10.1
)

require (
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output/enum v0.10.1 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)
