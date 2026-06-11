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

func TestRenderDOTTableData(t *testing.T) {
	t.Parallel()

	t.Run("writes DOT output to writer", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatDOT, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData dot: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("expected non-empty DOT output")
		}
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTableData(data, output.FormatDOT, output.RenderOptions{Writer: &errWriter{}})
		if err == nil {
			t.Fatal("expected error from errWriter")
		}
	})
}

func TestRenderMermaidTableData(t *testing.T) {
	t.Parallel()

	t.Run("writes Mermaid output to writer", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatMermaid, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData mermaid: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("expected non-empty Mermaid output")
		}
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTableData(data, output.FormatMermaid, output.RenderOptions{Writer: &errWriter{}})
		if err == nil {
			t.Fatal("expected error from errWriter")
		}
	})
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
