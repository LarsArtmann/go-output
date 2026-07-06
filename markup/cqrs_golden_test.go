package markup

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

func sampleCQRSTable() *output.Table {
	data := output.NewTable([]string{"Name", "Status"})
	data.AddRow([]string{"Build", "completed"})
	data.AddRow([]string{"Test", "running"})

	return data
}

func TestGolden_CQRS_XML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteXML(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}
