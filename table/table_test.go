package table

import (
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	testutils "github.com/larsartmann/go-output/internal/testutils"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tbl := New()
	if tbl == nil {
		t.Fatal("New() returned nil")
	}

	if tbl.t == nil {
		t.Error("New() internal table is nil")
	}
}

func TestSetHeaders(t *testing.T) {
	t.Parallel()

	tbl := New()
	result := tbl.SetHeaders("Name", "Value", "Count")

	if result != tbl {
		t.Error("SetHeaders() should return the same table for chaining")
	}

	output := tbl.Render()
	testutils.AssertContains(t, output, "Name", "Render() should contain header 'Name'")
	testutils.AssertContains(t, output, "Value", "Render() should contain header 'Value'")
	testutils.AssertContains(t, output, "Count", "Render() should contain header 'Count'")
}

func TestAddRow(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Value")
	result := tbl.AddRow("Alice", "30")

	if result != tbl {
		t.Error("AddRow() should return the same table for chaining")
	}

	output := tbl.Render()
	testutils.AssertContains(t, output, "Alice", "Render() should contain row value 'Alice'")
	testutils.AssertContains(t, output, "30", "Render() should contain row value '30'")
}

func TestAddRowMultiple(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Score")
	tbl.AddRow("Alice", "100")
	tbl.AddRow("Bob", "90")
	tbl.AddRow("Charlie", "85")

	output := tbl.Render()
	testutils.AssertContains(t, output, "Alice", "Render() should contain 'Alice'")
	testutils.AssertContains(t, output, "Bob", "Render() should contain 'Bob'")
	testutils.AssertContains(t, output, "Charlie", "Render() should contain 'Charlie'")
}

func TestStyleFunc(t *testing.T) {
	t.Parallel()

	tbl := New()
	result := tbl.StyleFunc(func(row, _ int) lipgloss.Style {
		if row == table.HeaderRow {
			return lipgloss.NewStyle().Bold(true)
		}

		return lipgloss.NewStyle()
	})

	if result != tbl {
		t.Error("StyleFunc() should return the same table for chaining")
	}

	tbl.SetHeaders("Test")
	tbl.AddRow("Value")

	output := tbl.Render()
	testutils.AssertContains(t, output, "Test", "Render() should contain 'Test'")
}

func TestRender(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Status")
	tbl.AddRow("Project A", "Active")
	tbl.AddRow("Project B", "Inactive")

	output := tbl.Render()

	if output == "" {
		t.Error("Render() should not return empty string")
	}

	testutils.AssertContains(t, output, "Name", "Render() should contain headers")
	testutils.AssertContains(t, output, "Status", "Render() should contain headers")
	testutils.AssertContains(t, output, "Project A", "Render() should contain row data")
	testutils.AssertContains(t, output, "Active", "Render() should contain row data")
}

func TestChaining(t *testing.T) {
	t.Parallel()

	output := New().
		SetHeaders("ID", "Name").
		AddRow("1", "First").
		AddRow("2", "Second").
		Render()

	testutils.AssertContains(t, output, "ID", "Chained Render() should contain headers")
	testutils.AssertContains(t, output, "Name", "Chained Render() should contain headers")
	testutils.AssertContains(t, output, "First", "Chained Render() should contain row data")
	testutils.AssertContains(t, output, "Second", "Chained Render() should contain row data")
}

func TestEmptyTable(t *testing.T) {
	t.Parallel()

	tbl := New()
	output := tbl.Render()

	// Empty table (no headers, no rows) returns empty string
	if output != "" {
		t.Errorf("Render() of empty table should be empty, got %q", output)
	}
}

func TestHeadersOnlyNoRows(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Only", "Headers")
	output := tbl.Render()

	testutils.AssertContains(t, output, "Only", "Render() should contain headers even without rows")
	testutils.AssertContains(
		t,
		output,
		"Headers",
		"Render() should contain headers even without rows",
	)
}

type testTableData struct {
	headers []string
	rows    [][]string
}

func (d *testTableData) GetHeaders() []string { return d.headers }
func (d *testTableData) GetRows() [][]string  { return d.rows }

func TestFromTableData(t *testing.T) {
	t.Parallel()

	data := &testTableData{
		headers: []string{"Name", "Status"},
		rows: [][]string{
			{"Project A", "Active"},
			{"Project B", "Inactive"},
		},
	}

	tbl := FromTableData(data)
	if tbl == nil {
		t.Fatal("FromTableData() returned nil")
	}

	output := tbl.Render()
	testutils.AssertContains(t, output, "Name", "should contain header 'Name'")
	testutils.AssertContains(t, output, "Status", "should contain header 'Status'")
	testutils.AssertContains(t, output, "Project A", "should contain row 'Project A'")
	testutils.AssertContains(t, output, "Active", "should contain row 'Active'")
}

func TestFromTableDataNil(t *testing.T) {
	t.Parallel()

	tbl := FromTableData(nil)
	if tbl == nil {
		t.Fatal("FromTableData(nil) should not return nil")
	}

	output := tbl.Render()
	if output != "" {
		t.Errorf("FromTableData(nil) should render empty, got %q", output)
	}
}

func TestFromTableDataEmpty(t *testing.T) {
	t.Parallel()

	data := &testTableData{
		headers: []string{},
		rows:    nil,
	}

	tbl := FromTableData(data)
	output := tbl.Render()

	if output != "" {
		t.Errorf("FromTableData with empty data should render empty, got %q", output)
	}
}
