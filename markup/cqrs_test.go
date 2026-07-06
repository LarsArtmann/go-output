package markup

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestCQRS_RenderXML(t *testing.T) {
	tbl := output.NewTable([]string{"Name", "Age"})
	tbl.AddRow([]string{"Alice", "30"})

	got, err := RenderXML(tbl)
	if err != nil {
		t.Fatalf("RenderXML failed: %v", err)
	}

	if !strings.Contains(got, "Alice") || !strings.Contains(got, "<") {
		t.Errorf("expected XML to contain Alice and tags, got %q", got)
	}
}

func TestCQRS_WriteXML(t *testing.T) {
	tbl := output.NewTable([]string{"A"})
	tbl.AddRow([]string{"1"})

	var buf strings.Builder

	if err := WriteXML(&buf, tbl); err != nil {
		t.Fatalf("WriteXML failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestCQRS_RenderAsciiDoc(t *testing.T) {
	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	got, err := RenderAsciiDoc(tbl)
	if err != nil {
		t.Fatalf("RenderAsciiDoc failed: %v", err)
	}

	if !strings.Contains(got, "Alice") {
		t.Errorf("expected AsciiDoc to contain Alice, got %q", got)
	}
}

func TestCQRS_RenderHTML(t *testing.T) {
	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	got, err := RenderHTML(tbl)
	if err != nil {
		t.Fatalf("RenderHTML failed: %v", err)
	}

	if !strings.Contains(got, "Alice") || !strings.Contains(got, "<") {
		t.Errorf("expected HTML to contain Alice and tags, got %q", got)
	}
}
