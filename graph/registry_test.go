package graph

import (
	"bytes"
	"errors"
	"testing"

	"github.com/larsartmann/go-output"
)

// errWriter always returns an error from Write.
type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errWriteFailed
}

var errWriteFailed = errors.New("write failed")

func TestRenderGraphTableData(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		format output.Format
	}{
		{"DOT", output.FormatDOT},
		{"Mermaid", output.FormatMermaid},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Run("writes output to writer", func(t *testing.T) {
				t.Parallel()

				data := output.NewTableData([]string{"Name"})
				data.AddRow([]string{"Alpha"})

				var buf bytes.Buffer

				err := output.RenderTableData(data, tc.format, output.RenderOptions{Writer: &buf})
				if err != nil {
					t.Fatalf("RenderTableData %s: %v", tc.name, err)
				}

				if buf.Len() == 0 {
					t.Errorf("expected non-empty %s output", tc.name)
				}
			})

			t.Run("propagates writer error", func(t *testing.T) {
				t.Parallel()

				data := output.NewTableData([]string{"Name"})
				data.AddRow([]string{"Alpha"})

				err := output.RenderTableData(data, tc.format, output.RenderOptions{Writer: &errWriter{}})
				if err == nil {
					t.Fatal("expected error from errWriter")
				}
			})
		})
	}
}

func TestNewGraphNodeIDAndLabel(t *testing.T) {
	t.Parallel()

	id := NewGraphNodeID("svc-a")
	if id.Get() != "svc-a" {
		t.Errorf("NewGraphNodeID Get = %q, want %q", id.Get(), "svc-a")
	}

	label := NewGraphNodeLabel("Service A")
	if label.Get() != "Service A" {
		t.Errorf("NewGraphNodeLabel Get = %q, want %q", label.Get(), "Service A")
	}
}
