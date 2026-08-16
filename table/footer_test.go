package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestFromTableWithFooter(t *testing.T) {
	t.Parallel()

	data := &testTableWithFooter{
		headers: []string{"Name", "Count"},
		rows: [][]string{
			{"Alice", "10"},
			{"Bob", "20"},
		},
		footer: []string{"Total", "30"},
	}

	tbl := FromTable(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Name", "should contain header")
	testhelpers.AssertContains(t, output, "Alice", "should contain row")
	testhelpers.AssertContains(t, output, "Total", "should contain footer row")
	testhelpers.AssertContains(t, output, "30", "should contain footer value")
}

func TestFromTableWithEmptyFooter(t *testing.T) {
	t.Parallel()

	data := &testTableWithFooter{
		headers: []string{"Name"},
		rows:    [][]string{{"Alice"}},
		footer:  nil,
	}

	tbl := FromTable(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Alice", "should contain row")
}

func TestFromTableWithRealTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Item", "Qty"})
	data.AddRow([]string{"Apple", "5"})
	data.Footer = []string{"Sum", "5"}

	tbl := FromTable(data)

	output, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, output, "Apple", "should contain data row")
	testhelpers.AssertContains(t, output, "Sum", "should contain footer from output.Table")
}

func TestTableSetFooter(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Count")
	tbl.AddRow("Alice", "10")
	tbl.AddRow("Bob", "20")
	tbl.SetFooter("Total", "30")

	testhelpers.RenderAssert(t, tbl, "Alice", "Total", "30")
}

func TestTableSetFooter_MultipleCalls(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeAlways))
	tbl.SetHeaders("Name", "Count")
	tbl.AddRow("Alice", "10")
	tbl.SetFooter("Old", "0")
	tbl.SetFooter("Total", "30")

	if tbl.footerRowIndex != 2 {
		t.Errorf("footerRowIndex = %d, want 2 (last footer row)", tbl.footerRowIndex)
	}

	sf := tbl.buildStyleFunc(tbl.footerRowIndex)
	footerStyle := sf(2, 0)
	if !footerStyle.GetBold() {
		t.Error("footer row (index 2) should be bold-styled")
	}

	if sf(1, 0).GetBold() {
		t.Error("data row (index 1) should not be footer-styled")
	}

	testhelpers.RenderAssert(t, tbl, "Alice", "Old", "Total")
}

func TestWithFooterStyle(t *testing.T) {
	t.Parallel()

	tbl := New(
		WithColorMode(output.ColorModeAlways),
		WithFooterStyle(func(s lipgloss.Style) lipgloss.Style {
			return s.Foreground(lipgloss.Color("196")).Italic(true)
		}),
	)
	tbl.SetHeaders("Name", "Count")
	tbl.AddRow("Alice", "10")
	tbl.SetFooter("Total", "10")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "Total", "should contain footer text")

	if !strings.Contains(got, "\033[") {
		t.Error("WithFooterStyle should produce ANSI output with ColorModeAlways")
	}
}
