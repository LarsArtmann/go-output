package table

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

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

func TestTableSetFooter_MultipleCalls(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Count")
	tbl.AddRow("Alice", "10")
	tbl.SetFooter("Old", "0")
	tbl.SetFooter("Total", "30")

	if tbl.footerRowIndex != 3 {
		t.Errorf("footerRowIndex = %d, want 3 (last footer row)", tbl.footerRowIndex)
	}

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "should contain data row")
	testhelpers.AssertContains(t, output, "Old", "should contain first footer as data")
	testhelpers.AssertContains(t, output, "Total", "should contain last footer")
}
