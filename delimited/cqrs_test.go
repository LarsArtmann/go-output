package delimited

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestCQRS_RenderCSV(t *testing.T) {
	tbl := output.NewTable([]string{"Name", "Age"})
	tbl.AddRow([]string{"Alice", "30"})
	tbl.AddRow([]string{"Bob", "25"})

	got, err := RenderCSV(tbl)
	if err != nil {
		t.Fatalf("RenderCSV failed: %v", err)
	}

	if !strings.Contains(got, "Name") || !strings.Contains(got, "Alice") {
		t.Errorf("expected CSV to contain Name and Alice, got %q", got)
	}
}

func TestCQRS_WriteCSV(t *testing.T) {
	tbl := output.NewTable([]string{"A", "B"})
	tbl.AddRow([]string{"1", "2"})

	var buf strings.Builder
	if err := WriteCSV(&buf, tbl); err != nil {
		t.Fatalf("WriteCSV failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output from WriteCSV")
	}
}

func TestCQRS_RenderTSV(t *testing.T) {
	tbl := output.NewTable([]string{"Name", "Age"})
	tbl.AddRow([]string{"Alice", "30"})

	got, err := RenderTSV(tbl)
	if err != nil {
		t.Fatalf("RenderTSV failed: %v", err)
	}

	if !strings.Contains(got, "Alice") {
		t.Errorf("expected TSV to contain Alice, got %q", got)
	}
}
