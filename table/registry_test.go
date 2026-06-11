package table

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

func TestRenderStyledTableData(t *testing.T) {
	t.Parallel()

	t.Run("writes table output to writer", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatTable, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData table: %v", err)
		}

		if buf.Len() == 0 {
			t.Error("expected non-empty table output")
		}
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTableData(data, output.FormatTable, output.RenderOptions{Writer: &errWriter{}})
		if err == nil {
			t.Fatal("expected error from errWriter")
		}
	})
}
