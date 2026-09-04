package integration

import (
	"strings"
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

func TestNewlineProbe(t *testing.T) {
	data := output.NewTableWithRow([]string{"A", "B"}, "1", "2")

	for _, f := range output.AllFormats {
		var reg string
		if err := output.RenderTable(data, f, output.RenderOptions{Writer: &stringWriter{w: &reg}}); err != nil {
			t.Logf("%-10s registry: ERR %v", f, err)
			continue
		}
		t.Logf("%-10s registry: trailingNL=%d bytes=%d", f, countNL(reg), len(reg))
	}
}

type stringWriter struct{ w *string }

func (s *stringWriter) Write(p []byte) (int, error) {
	*s.w += string(p)
	return len(p), nil
}

func countNL(s string) int {
	n := 0
	for _, r := range s {
		if r == '\n' {
			n++
		}
	}
	_ = strings.TrimRight
	return n
}
