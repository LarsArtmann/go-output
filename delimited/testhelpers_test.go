package delimited

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

type errorWriter = testhelpers.ErrorWriter

var assertContains = testhelpers.AssertContains

func testRenderTableData(t *testing.T, format output.Format, name string) {
	t.Helper()

	t.Run("renders via registry dispatch", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Value"})
		data.AddRow([]string{"Alpha", "100"})

		var buf bytes.Buffer

		opts := output.RenderOptions{Writer: &buf}

		err := output.RenderTableData(data, format, opts)
		if err != nil {
			t.Fatalf("RenderTableData(%s) error = %v", name, err)
		}

		result := buf.String()
		assertContains(t, result, "Name", name+" render should contain header")
		assertContains(t, result, "Alpha", name+" render should contain data")
	})

	t.Run("nil data returns nil", func(t *testing.T) {
		t.Parallel()

		err := output.RenderTableData(nil, format)
		if err != nil {
			t.Fatalf("RenderTableData(nil) error = %v", err)
		}
	})
}
