package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
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

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Name", "Render() should contain header 'Name'")
	testhelpers.AssertContains(t, output, "Value", "Render() should contain header 'Value'")
	testhelpers.AssertContains(t, output, "Count", "Render() should contain header 'Count'")
}

func TestAddRow(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Value")
	result := tbl.AddRow("Alice", "30")

	if result != tbl {
		t.Error("AddRow() should return the same table for chaining")
	}

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "Render() should contain row value 'Alice'")
	testhelpers.AssertContains(t, output, "30", "Render() should contain row value '30'")
}

func TestAddRowMultiple(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Score")
	tbl.AddRow("Alice", "100")
	tbl.AddRow("Bob", "90")
	tbl.AddRow("Charlie", "85")

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "Render() should contain 'Alice'")
	testhelpers.AssertContains(t, output, "Bob", "Render() should contain 'Bob'")
	testhelpers.AssertContains(t, output, "Charlie", "Render() should contain 'Charlie'")
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

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Test", "Render() should contain 'Test'")
}

func TestRender(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Status")
	tbl.AddRow("Project A", "Active")
	tbl.AddRow("Project B", "Inactive")

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if output == "" {
		t.Error("Render() should not return empty string")
	}

	testhelpers.AssertContains(t, output, "Name", "Render() should contain headers")
	testhelpers.AssertContains(t, output, "Status", "Render() should contain headers")
	testhelpers.AssertContains(t, output, "Project A", "Render() should contain row data")
	testhelpers.AssertContains(t, output, "Active", "Render() should contain row data")
}

func TestChaining(t *testing.T) {
	t.Parallel()

	output, err := New().
		SetHeaders("ID", "Name").
		AddRow("1", "First").
		AddRow("2", "Second").
		Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "ID", "Chained Render() should contain headers")
	testhelpers.AssertContains(t, output, "Name", "Chained Render() should contain headers")
	testhelpers.AssertContains(t, output, "First", "Chained Render() should contain row data")
	testhelpers.AssertContains(t, output, "Second", "Chained Render() should contain row data")
}

func TestEmptyTable(t *testing.T) {
	t.Parallel()

	tbl := New()

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Empty table (no headers, no rows) returns empty string
	if output != "" {
		t.Errorf("Render() of empty table should be empty, got %q", output)
	}
}

func TestHeadersOnlyNoRows(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Only", "Headers")

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(
		t,
		output,
		"Only",
		"Render() should contain headers even without rows",
	)
	testhelpers.AssertContains(
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

type testTableDataWithFooter struct {
	headers []string
	rows    [][]string
	footer  []string
}

func (d *testTableDataWithFooter) GetHeaders() []string { return d.headers }
func (d *testTableDataWithFooter) GetRows() [][]string  { return d.rows }
func (d *testTableDataWithFooter) GetFooter() []string  { return d.footer }

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

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Name", "should contain header 'Name'")
	testhelpers.AssertContains(t, output, "Status", "should contain header 'Status'")
	testhelpers.AssertContains(t, output, "Project A", "should contain row 'Project A'")
	testhelpers.AssertContains(t, output, "Active", "should contain row 'Active'")
}

func TestFromTableDataNil(t *testing.T) {
	t.Parallel()

	tbl := FromTableData(nil)
	if tbl == nil {
		t.Fatal("FromTableData(nil) should not return nil")
	}

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

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

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if output != "" {
		t.Errorf("FromTableData with empty data should render empty, got %q", output)
	}
}

func TestTableColorModeNever(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Value")
	tbl.AddRow("Alice", "30")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeNever should produce no ANSI codes, got: %q", got)
	}

	testhelpers.AssertContains(t, got, "Alice", "should contain data even without colors")
}

func TestTableColorModeAlways(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeAlways))
	tbl.SetHeaders("Name", "Value")
	tbl.AddRow("Alice", "30")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeAlways should produce ANSI codes, got: %q", got)
	}
}

func TestTableColorModeDefault(t *testing.T) {
	t.Parallel()

	tbl := New()
	if tbl.colorMode != output.ColorModeAuto {
		t.Errorf("default ColorMode = %v, want %v", tbl.colorMode, output.ColorModeAuto)
	}
}

func TestFromTableDataWithFooter(t *testing.T) {
	t.Parallel()

	data := &testTableDataWithFooter{
		headers: []string{"Name", "Count"},
		rows: [][]string{
			{"Alice", "10"},
			{"Bob", "20"},
		},
		footer: []string{"Total", "30"},
	}

	tbl := FromTableData(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Name", "should contain header")
	testhelpers.AssertContains(t, output, "Alice", "should contain row")
	testhelpers.AssertContains(t, output, "Total", "should contain footer row")
	testhelpers.AssertContains(t, output, "30", "should contain footer value")
}

func TestFromTableDataWithEmptyFooter(t *testing.T) {
	t.Parallel()

	data := &testTableDataWithFooter{
		headers: []string{"Name"},
		rows:    [][]string{{"Alice"}},
		footer:  nil,
	}

	tbl := FromTableData(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "should contain row")
}

func TestFromTableDataWithRealTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Item", "Qty"})
	data.AddRow([]string{"Apple", "5"})
	data.Footer = []string{"Sum", "5"}

	tbl := FromTableData(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Apple", "should contain data row")
	testhelpers.AssertContains(t, output, "Sum", "should contain footer from output.TableData")
}

func TestTableSetFooter(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Count")
	tbl.AddRow("Alice", "10")
	tbl.AddRow("Bob", "20")
	tbl.SetFooter("Total", "30")

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "should contain row")
	testhelpers.AssertContains(t, output, "Total", "should contain footer")
	testhelpers.AssertContains(t, output, "30", "should contain footer value")
}
