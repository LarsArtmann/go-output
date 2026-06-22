module github.com/larsartmann/go-output/tree

go 1.26.3

require (
	github.com/larsartmann/go-output v0.17.1
	github.com/larsartmann/go-output/testhelpers v0.13.0
)

require (
	github.com/larsartmann/go-branded-id v0.3.1 // indirect
	github.com/larsartmann/go-output/enum v0.17.1 // indirect
	github.com/larsartmann/go-output/envdetect v0.17.1 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/term v0.44.0 // indirect
)

replace (
	github.com/larsartmann/go-output => ../
	github.com/larsartmann/go-output/enum => ../enum
	github.com/larsartmann/go-output/envdetect => ../envdetect
	github.com/larsartmann/go-output/testhelpers => ../testhelpers
)

replace github.com/larsartmann/go-output/markdown => ../markdown
