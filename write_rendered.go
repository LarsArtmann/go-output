package output

import (
	"fmt"
	"io"
	"strings"
)

// WriteRendered writes out to w with a trailing newline and wraps any write
// error as "write <format> output: %w". It is the shared tail of every
// diagram marshaler registered by the d2 / graph (DOT, Mermaid) / plantuml /
// table / tree sub-modules' renderTable — they render, this writes.
// Centralising it removes a five-line block duplicated across each
// sub-module. Returns nil on a successful write.
func WriteRendered(w io.Writer, formatName, out string) error {
	if _, err := fmt.Fprintln(w, out); err != nil {
		return fmt.Errorf("write %s output: %w", formatName, err)
	}

	return nil
}

// WriteRenderedRaw writes out to w WITHOUT a trailing newline and wraps any
// write error as "<format> output: %w". It is the shared tail of every
// sub-module's CQRS Write function (markdown.Write, plantuml.Write,
// tree.WriteASCII) where the rendered payload already carries its own
// trailing whitespace and an extra newline would corrupt the output.
// Returns nil on a successful write.
func WriteRenderedRaw(w io.Writer, formatName, out string) error {
	if _, err := io.WriteString(w, out); err != nil {
		return fmt.Errorf("%s output: %w", formatName, err)
	}

	return nil
}

// RenderFromWrite adapts a streaming Write-style function (one that takes
// an io.Writer) into the matching Render convenience function (one that
// returns the full string). It is the shared body of every sub-module's
// Render(data, opts...) (string, error) wrapper around its Write(w, data,
// opts...) error sibling. Centralising it removes a 5-line copy/paste
// between markdown.Render, table.Render, plantuml.Render, tree.RenderASCII,
// etc.
func RenderFromWrite(write func(io.Writer) error) (string, error) {
	var buf strings.Builder
	if err := write(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteRenderedRawFrom is the shared body of every CQRS Write function that
// takes a configured renderer, calls its Render method, wraps any render
// error with a "<format>:" prefix, and writes the rendered output via
// WriteRenderedRaw. It collapses the canonical
//
//	out, err := renderer.Render()
//	if err != nil { return fmt.Errorf("%s: %w", err) }
//	return output.WriteRenderedRaw(w, formatName, out)
//
// pattern into a single call.
//
// Pass the renderer's Render method as a method value (e.g. `m.Render`);
// formatName is the human-readable format string ("markdown", "plantuml",
// "ascii tree"); formatErrorPrefix is the wrapping key for render errors
// (typically the same as formatName).
func WriteRenderedRawFrom(
	w io.Writer,
	render func() (string, error),
	formatName, formatErrorPrefix string,
) error {
	out, err := render()
	if err != nil {
		return fmt.Errorf("%s: %w", formatErrorPrefix, err)
	}

	return WriteRenderedRaw(w, formatName, out)
}

// WriteRenderedFrom mirrors WriteRenderedRawFrom but routes the rendered
// output through WriteRendered (which appends a trailing newline). It is
// the shared body of CQRS Write functions whose rendered payload does NOT
// carry its own trailing newline — e.g. table.Write, tree/registry's
// renderTreeTable, and any renderers writing plain text without an
// embedded \n.
//
// Same semantics as WriteRenderedRawFrom: pass the renderer's Render method
// as a method value, formatName for the trailing error wrap, and
// formatErrorPrefix for the render-time error wrap.
func WriteRenderedFrom(
	w io.Writer,
	render func() (string, error),
	formatName, formatErrorPrefix string,
) error {
	out, err := render()
	if err != nil {
		return fmt.Errorf("%s: %w", formatErrorPrefix, err)
	}

	return WriteRendered(w, formatName, out)
}
