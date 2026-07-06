package table

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestRenderStyledTable(t *testing.T) {
	t.Parallel()

	t.Run("writes table output to writer", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		var buf bytes.Buffer

		err := output.RenderTable(data, output.FormatTable, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTable table: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("expected non-empty table output")
		}
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTable(
			data,
			output.FormatTable,
			output.RenderOptions{Writer: &testhelpers.ErrorWriter{}},
		)
		if err == nil {
			t.Fatal("expected error from errWriter")
		}
	})
}
