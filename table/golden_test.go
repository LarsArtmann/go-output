package table

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

func TestGolden_Table_BasicHeadersRows(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Status", "Duration")
	tbl.AddRow("Build", "Completed", "1.2s")
	tbl.AddRow("Test", "Running", "0.5s")
	tbl.AddRow("Deploy", "Pending", "-")

	got, err := tbl.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Table_WithFooterRow(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Module", "Coverage", "Status")
	tbl.AddRow("core", "92%", "pass")
	tbl.AddRow("cli", "78%", "pass")
	tbl.AddRow("utils", "45%", "fail")
	tbl.SetFooter("Average", "71%", "2/3")

	got, err := tbl.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Table_SingleColumn(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Step")
	tbl.AddRow("Download dependencies")
	tbl.AddRow("Compile sources")
	tbl.AddRow("Run tests")
	tbl.AddRow("Generate docs")

	got, err := tbl.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Table_Empty(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))

	got, err := tbl.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
