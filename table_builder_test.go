package output

import (
	"testing"
)

func TestTableBuilder_FluentAPI(t *testing.T) {
	tbl := NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		SetFooter("Total", "2 tasks").
		Build()

	if len(tbl.Headers) != 2 {
		t.Fatalf("expected 2 headers, got %d", len(tbl.Headers))
	}

	if tbl.Headers[0] != "Name" {
		t.Errorf("expected first header 'Name', got %q", tbl.Headers[0])
	}

	if len(tbl.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(tbl.Rows))
	}

	if tbl.Rows[0][0] != "Compile" {
		t.Errorf("expected first row first cell 'Compile', got %q", tbl.Rows[0][0])
	}

	if len(tbl.Footer) != 2 {
		t.Fatalf("expected 2 footer cells, got %d", len(tbl.Footer))
	}
}

func TestTableBuilder_BuildIsSnapshot(t *testing.T) {
	b := NewTableBuilder().
		SetHeaders("A", "B").
		AddRow("1", "2")

	t1 := b.Build()

	b.AddRow("3", "4")

	t2 := b.Build()

	if len(t1.Rows) != 1 {
		t.Errorf("first build should have 1 row, got %d", len(t1.Rows))
	}

	if len(t2.Rows) != 2 {
		t.Errorf("second build should have 2 rows, got %d", len(t2.Rows))
	}
}

func TestTableBuilder_Empty(t *testing.T) {
	tbl := NewTableBuilder().Build()

	if tbl == nil {
		t.Fatal("Build() should never return nil")
	}

	if len(tbl.Headers) != 0 {
		t.Errorf("expected 0 headers, got %d", len(tbl.Headers))
	}

	if len(tbl.Rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(tbl.Rows))
	}

	if tbl.Footer != nil {
		t.Errorf("expected nil footer, got %v", tbl.Footer)
	}
}
