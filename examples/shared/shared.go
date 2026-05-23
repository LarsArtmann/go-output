// Package shared provides common utilities for go-output examples.
package shared

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output/d2"
)

// HandleError prints the error to stderr and exits with code 1.
func HandleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// NewServiceD2Diagram creates a D2 diagram with a service class pre-configured.
func NewServiceD2Diagram(title string) *d2.D2Diagram {
	return d2.NewD2Diagram().
		SetDirection(d2.D2DirRight).
		SetTitle(title).
		AddClass("service", d2.D2NodeStyle{
			Fill:     "lightblue",
			Stroke:   "navy",
			FontSize: 16,
		})
}
