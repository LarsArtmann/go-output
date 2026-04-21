package output

import (
	"testing"
)

func TestTableData(t *testing.T) {
	t.Parallel()
	runSubtest(t, "RowCount and ColCount", testTableDataRowColCount)
	runSubtest(t, "CreateRowEdges", testTableDataCreateRowEdges)
}

func testTableDataRowColCount(t *testing.T) {
	t.Helper()

	data := NewTableData([]string{"Name", "Value", "Count"})

	cc := data.ColCount()
	if cc != 3 {
		t.Errorf("ColCount() = %d, want 3", cc)
	}

	rc := data.RowCount()
	if rc != 0 {
		t.Errorf("RowCount() = %d, want 0", rc)
	}

	data.AddRow([]string{"a", "b", "c"})
	data.AddRow([]string{"d", "e", "f"})

	rowCount := data.RowCount()
	if rowCount != 2 {
		t.Errorf("RowCount() = %d, want 2", rowCount)
	}
}

func testTableDataCreateRowEdges(t *testing.T) {
	t.Helper()
	t.Run("nil data", testCreateRowEdgesNil)
	t.Run("empty rows", testCreateRowEdgesEmpty)
	t.Run("single row", testCreateRowEdgesSingle)
	t.Run("multiple rows", testCreateRowEdgesMultiple)
}

func testAssertNilRowEdges(t *testing.T, data *TableData, desc string) {
	if edges := data.CreateRowEdges(); edges != nil {
		t.Errorf("CreateRowEdges() on %s = %v, want nil", desc, edges)
	}
}

func testCreateRowEdgesNil(t *testing.T) {
	var data *TableData
	testAssertNilRowEdges(t, data, "nil")
}

func testCreateRowEdgesEmpty(t *testing.T) {
	data := NewTableData([]string{"Name"})
	testAssertNilRowEdges(t, data, "empty")
}

func testCreateRowEdgesSingle(t *testing.T) {
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"a"})

	testAssertNilRowEdges(t, data, "single row")
}

func testCreateRowEdgesMultiple(t *testing.T) {
	data := NewTableData([]string{"Name"})
	data.AddRow([]string{"a"})
	data.AddRow([]string{"b"})
	data.AddRow([]string{"c"})

	edges := data.CreateRowEdges()
	if len(edges) != 2 {
		t.Fatalf("CreateRowEdges() returned %d edges, want 2", len(edges))
	}

	verifyEdge := func(idx int, from, to string) {
		if edges[idx].From != from || edges[idx].To != to {
			t.Errorf(
				"Edge %d = {%s, %s}, want {%s, %s}",
				idx,
				edges[idx].From,
				edges[idx].To,
				from,
				to,
			)
		}
	}
	verifyEdge(0, "row0", "row1")
	verifyEdge(1, "row1", "row2")
}
