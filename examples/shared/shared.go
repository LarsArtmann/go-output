// Package shared provides common utilities for go-output examples.
package shared

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output"
)

// HandleError prints the error to stderr and exits with code 1.
func HandleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// NewServiceD2Diagram creates a D2 diagram with a service class pre-configured.
func NewServiceD2Diagram(title string) *output.D2Diagram {
	return output.NewD2Diagram().
		SetDirection(output.D2DirRight).
		SetTitle(title).
		AddClass("service", output.D2NodeStyle{
			Fill:     "lightblue",
			Stroke:   "navy",
			FontSize: 16,
		})
}
