package graph

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

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

				err := output.RenderTableData(data, tc.format, output.RenderOptions{Writer: &testhelpers.ErrorWriter{}})
				if err == nil {
					t.Fatal("expected error from errWriter")
				}
			})
		})
	}
}

func TestBrandedGraphNodeIDAndLabel(t *testing.T) {
	t.Parallel()

	id := output.NewBrandedID[output.GraphNodeIDBrand]("svc-a")
	if id.Get() != "svc-a" {
		t.Errorf("GraphNodeID Get = %q, want %q", id.Get(), "svc-a")
	}

	label := output.NewBrandedID[output.GraphNodeLabelBrand]("Service A")
	if label.Get() != "Service A" {
		t.Errorf("GraphNodeLabel Get = %q, want %q", label.Get(), "Service A")
	}
}
