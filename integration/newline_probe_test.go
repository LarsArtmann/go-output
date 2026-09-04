package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
	_ "github.com/larsartmann/go-output/delimited"
	_ "github.com/larsartmann/go-output/d2"
	_ "github.com/larsartmann/go-output/graph"
	_ "github.com/larsartmann/go-output/markdown"
	_ "github.com/larsartmann/go-output/markup"
	_ "github.com/larsartmann/go-output/plantuml"
	_ "github.com/larsartmann/go-output/serialization"
	_ "github.com/larsartmann/go-output/table"
	_ "github.com/larsartmann/go-output/tree"
)

type nlWriter struct {
	b    []byte
	nl   int
	done bool
}

func (w *nlWriter) Write(p []byte) (int, error) {
	w.b = append(w.b, p...)
	return len(p), nil
}

func TestNewlineProbe(t *testing.T) {
	data := output.NewTableWithRow([]string{"A", "B"}, "1", "2")

	for _, f := range output.AllFormats {
		w := &nlWriter{}
		if err := output.RenderTable(data, f, output.RenderOptions{Writer: w}); err != nil {
			t.Logf("%-10s registry: ERR %v", f, err)
			continue
		}
		suffix := 0
		for suffix < len(w.b) && w.b[len(w.b)-1-suffix] == '\n' {
			suffix++
		}
		t.Logf("%-10s trailing_newlines=%d", f, suffix)
	}
}
