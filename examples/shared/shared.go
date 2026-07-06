// Package shared provides common utilities for go-output examples.
package shared

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
)

// HandleError prints the error to stderr and exits with code 1.
func HandleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

// NewServiceD2Diagram creates a D2 diagram with a service class pre-configured.
func NewServiceD2Diagram(title string) *d2.Diagram {
	return d2.NewDiagram().
		SetDirection(d2.DirRight).
		SetTitle(title).
		AddClass("service", d2.NodeStyle{
			Fill: "lightblue",
			StrokeStyle: d2.StrokeStyle{
				Stroke:   "navy",
				FontSize: 16,
			},
		})
}

// RenderAndPrint runs r.Render() and prints the result, routing any error
// through HandleError. Centralises the "Render() + check error + println"
// sequence shared by every go-output example entrypoint.
func RenderAndPrint(r output.Renderer) {
	out, err := r.Render()
	if err != nil {
		HandleError(err)
	}

	fmt.Println(out)
}
