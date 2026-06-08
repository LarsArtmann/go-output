package output

import (
	"testing"
)

func TestTableData(t *testing.T) {
	t.Parallel()
	runSubtest(t, "RowCount and ColCount", testTableDataRowColCount)
	runSubtest(t, "CreateRowEdges", testTableDataCreateRowEdges)
	runSubtest(t, "ToMapSlice", testTableDataToMapSlice)
	runSubtest(t, "Footer", testTableDataFooter)
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

func testAssertToMapSliceNil(t *testing.T, data *TableData, desc string) {
	t.Helper()

	if got := data.ToMapSlice(); got != nil {
		t.Errorf("ToMapSlice() on %s = %v, want nil", desc, got)
	}
}

func assertMapFields(t *testing.T, got, want map[string]string) {
	t.Helper()

	for k, v := range want {
		if got[k] != v {
			t.Errorf("map[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func testTableDataToMapSlice(t *testing.T) {
	t.Helper()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		var data *TableData
		testAssertToMapSliceNil(t, data, "nil")
	})

	t.Run("no headers", func(t *testing.T) {
		t.Parallel()

		data := &TableData{Rows: [][]string{{"a"}}}
		testAssertToMapSliceNil(t, data, "no headers")
	})

	t.Run("maps rows to headers", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"Name", "Age"})
		data.AddRow([]string{"Alice", "30"})
		data.AddRow([]string{"Bob", "25"})

		got := data.ToMapSlice()

		if len(got) != 2 {
			t.Fatalf("ToMapSlice() returned %d maps, want 2", len(got))
		}

		assertMapFields(t, got[0], map[string]string{"Name": "Alice", "Age": "30"})
		assertMapFields(t, got[1], map[string]string{"Name": "Bob", "Age": "25"})
	})

	t.Run("short row omits missing cells", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"A", "B", "C"})
		data.AddRow([]string{"1"})

		got := data.ToMapSlice()

		if got[0]["A"] != "1" {
			t.Errorf("ToMapSlice()[0][A] = %q, want 1", got[0]["A"])
		}

		if _, ok := got[0]["B"]; ok {
			t.Errorf("ToMapSlice()[0][B] = %q, want absent", got[0]["B"])
		}
	})
}

func testTableDataFooter(t *testing.T) {
	t.Helper()

	t.Run("no footer by default", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"Name", "Value"})
		if data.HasFooter() {
			t.Error("HasFooter() = true, want false")
		}

		if got := data.GetFooter(); got != nil {
			t.Errorf("GetFooter() = %v, want nil", got)
		}
	})

	t.Run("set and get footer", func(t *testing.T) {
		t.Parallel()

		data := NewTableData([]string{"Name", "Count"})
		data.AddRow([]string{"Alice", "10"})
		data.AddRow([]string{"Bob", "20"})
		data.Footer = []string{"Total", "30"}

		if !data.HasFooter() {
			t.Error("HasFooter() = false, want true")
		}

		footer := data.GetFooter()
		if len(footer) != 2 || footer[0] != "Total" || footer[1] != "30" {
			t.Errorf("GetFooter() = %v, want [Total 30]", footer)
		}
	})
}

func TestTableDataNilMethods(t *testing.T) {
	t.Parallel()

	var data *TableData

	if got := data.RowCount(); got != 0 {
		t.Errorf("RowCount() on nil = %d, want 0", got)
	}

	if got := data.ColCount(); got != 0 {
		t.Errorf("ColCount() on nil = %d, want 0", got)
	}

	if got := data.GetHeaders(); got != nil {
		t.Errorf("GetHeaders() on nil = %v, want nil", got)
	}

	if got := data.GetRows(); got != nil {
		t.Errorf("GetRows() on nil = %v, want nil", got)
	}

	if got := data.GetFooter(); got != nil {
		t.Errorf("GetFooter() on nil = %v, want nil", got)
	}

	if data.HasFooter() {
		t.Error("HasFooter() on nil = true, want false")
	}

	// SetFooter on nil should not panic
	data.SetFooter([]string{"a", "b"})
}

func TestTableDataValidate(t *testing.T) {
	t.Parallel()

	buildData := func(headers, footer []string) *TableData {
		data := NewTableData(headers)
		if footer != nil {
			data.SetFooter(footer)
		}

		return data
	}

	for _, tt := range []struct {
		name    string
		data    *TableData
		wantErr bool
	}{
		{name: "nil data", data: nil, wantErr: false},
		{name: "no footer", data: buildData([]string{"A", "B"}, nil), wantErr: false},
		{name: "matching columns", data: buildData([]string{"A", "B"}, []string{"1", "2"}), wantErr: false},
		{name: "empty footer cells", data: buildData([]string{"A", "B"}, []string{"", ""}), wantErr: false},
		{name: "empty headers with footer", data: &TableData{Footer: []string{"x"}}, wantErr: false},
		{name: "footer longer than headers", data: buildData([]string{"A"}, []string{"1", "2", "3"}), wantErr: true},
		{name: "footer shorter than headers", data: buildData([]string{"A", "B", "C"}, []string{"1"}), wantErr: true},
		{name: "nil row", data: &TableData{Headers: []string{"A"}, Rows: [][]string{nil}}, wantErr: true},
		{name: "nil row after valid row", data: &TableData{Headers: []string{"A", "B"}, Rows: [][]string{{"1", "2"}, nil}}, wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.data.Validate()
			if tt.wantErr && err == nil {
				t.Error("Validate() = nil, want error")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Validate() = %v, want nil", err)
			}
		})
	}
}
