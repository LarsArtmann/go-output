package delimited

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

// sampleCQRSTable returns a stable dataset for CQRS golden tests.
func sampleCQRSTable() *output.Table {
	data := output.NewTable([]string{"Name", "Status", "Duration"})
	data.AddRow([]string{"Build", "completed", "1.2s"})
	data.AddRow([]string{"Test", "running", "0.5s"})
	data.SetFooter([]string{"Total", "", "1.7s"})

	return data
}

// TestGolden_CQRS_CSV locks in the exact byte output of WriteCSV.
func TestGolden_CQRS_CSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteCSV(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestGolden_CQRS_TSV locks in the exact byte output of WriteTSV.
func TestGolden_CQRS_TSV(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteTSV(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}
