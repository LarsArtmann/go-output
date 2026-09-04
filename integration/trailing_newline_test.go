package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
	_ "github.com/larsartmann/go-output/d2"
	_ "github.com/larsartmann/go-output/delimited"
	_ "github.com/larsartmann/go-output/graph"
	_ "github.com/larsartmann/go-output/markdown"
	_ "github.com/larsartmann/go-output/markup"
	_ "github.com/larsartmann/go-output/plantuml"
	_ "github.com/larsartmann/go-output/serialization"
	_ "github.com/larsartmann/go-output/table"
	_ "github.com/larsartmann/go-output/tree"
)

// capturingWriter records everything written so tests can assert on the
// exact bytes a marshaler produced.
type capturingWriter struct {
	b []byte
}

func (w *capturingWriter) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

// TestRegistryOutput_SingleTrailingNewline pins the unified newline rule
// across ALL formats: every registry-dispatched render ends with exactly
// one trailing newline — the Go/POSIX canonical for text output. Before
// the rule, markdown/tree/dot/mermaid emitted two and asciidoc emitted
// zero, so consumers appending their own newline got ragged output
// depending on format.
func TestRegistryOutput_SingleTrailingNewline(t *testing.T) {
	t.Parallel()

	data := output.NewTableWithRow([]string{"Name", "Status"}, "Compile", "done")

	for _, f := range output.AllFormats {
		f := f
		t.Run(string(f), func(t *testing.T) {
			t.Parallel()

			w := &capturingWriter{}
			if err := output.RenderTable(data, f, output.RenderOptions{Writer: w}); err != nil {
				t.Fatalf("RenderTable(%s): %v", f, err)
			}

			if len(w.b) == 0 {
				t.Fatalf("RenderTable(%s) produced no output", f)
			}

			if w.b[len(w.b)-1] != '\n' {
				t.Errorf("RenderTable(%s) output does not end with a newline: %q", f, tail(w.b))
			}

			if len(w.b) >= 2 && w.b[len(w.b)-2] == '\n' {
				t.Errorf("RenderTable(%s) output ends with multiple newlines: %q", f, tail(w.b))
			}
		})
	}
}

func tail(b []byte) string {
	const n = 16
	if len(b) > n {
		return string(b[len(b)-n:])
	}
	return string(b)
}
