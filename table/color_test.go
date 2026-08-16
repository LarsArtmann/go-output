package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

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

// TestTable_ColoredFooterBoldOnFooterLineOnly asserts end-to-end that the
// bold style lands on the line whose CONTENT is the footer — not merely on
// whatever row index footerRowIndex happens to name. The style-level tests
// in TestBuildStyleFunc_DirectCall read tbl.footerRowIndex itself, so they
// cannot catch an off-by-one in the rendered table (the footer-index bug
// class fixed in the v0.38 review).
func TestTable_ColoredFooterBoldOnFooterLineOnly(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeAlways))
	tbl.SetHeaders("Name", "Value")
	tbl.AddRow("Alice", "30")
	tbl.AddRow("Bob", "40")
	tbl.SetFooter("TOTAL", "70")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	bold := "\x1b[1m"

	var footerLine, dataLine string

	for _, line := range strings.Split(got, "\n") {
		switch {
		case strings.Contains(line, "TOTAL"):
			footerLine = line
		case strings.Contains(line, "Alice"):
			dataLine = line
		}
	}

	if footerLine == "" || dataLine == "" {
		t.Fatalf("rendered table missing footer or data line:\n%s", got)
	}

	if !strings.Contains(footerLine, bold) {
		t.Errorf("footer line should carry bold styling:\n%q", footerLine)
	}

	if strings.Contains(dataLine, bold) {
		t.Errorf("data row should not carry footer bold styling:\n%q", dataLine)
	}
}

func TestTableColorModeDefault(t *testing.T) {
	t.Parallel()

	tbl := New()
	if tbl.colorMode != output.ColorModeAuto {
		t.Errorf("default ColorMode = %v, want %v", tbl.colorMode, output.ColorModeAuto)
	}
}

func TestAsTableRenderer(t *testing.T) {
	t.Parallel()

	tbl := New()
	tr := tbl.AsTableRenderer()

	tr.SetHeaders([]string{"A", "B"})
	tr.AddRow([]string{"1", "2"})

	got, err := tr.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "A", "should contain header A")
	testhelpers.AssertContains(t, got, "1", "should contain cell 1")

	_ = tr
}

func TestBuildStyleFunc_AllBranches(t *testing.T) {
	t.Parallel()

	t.Run("always color header footer even odd rows", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeAlways))
		tbl.SetHeaders("A", "B")
		tbl.AddRow("1", "2")
		tbl.AddRow("3", "4")
		tbl.SetFooter("T", "6")

		got, err := tbl.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(got, "\x1b[") {
			t.Error("ColorModeAlways should produce ANSI codes")
		}

		testhelpers.AssertContains(t, got, "1", "should contain row")
		testhelpers.AssertContains(t, got, "3", "should contain row")
		testhelpers.AssertContains(t, got, "T", "should contain footer")
	})

	t.Run("never color header footer even odd rows", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeNever))
		tbl.SetHeaders("A", "B")
		tbl.AddRow("1", "2")
		tbl.AddRow("3", "4")
		tbl.SetFooter("T", "6")

		got, err := tbl.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if strings.Contains(got, "\x1b[") {
			t.Error("ColorModeNever should produce no ANSI codes")
		}

		testhelpers.AssertContains(t, got, "1", "should contain row")
		testhelpers.AssertContains(t, got, "3", "should contain row")
		testhelpers.AssertContains(t, got, "T", "should contain footer")
	})
}

func TestBuildStyleFunc_DirectCall(t *testing.T) {
	t.Parallel()

	t.Run("color always all branches", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeAlways))
		tbl.SetFooter("T") // sets footerRowIndex = 1 (after 0 rows)

		sf := tbl.buildStyleFunc(tbl.footerRowIndex)

		headerStyle := sf(table.HeaderRow, 0)
		if !headerStyle.GetBold() {
			t.Error("color always: header should be bold")
		}

		footerStyle := sf(tbl.footerRowIndex, 0)
		if !footerStyle.GetBold() {
			t.Error("color always: footer should be bold")
		}

		evenStyle := sf(0, 0)
		if evenStyle.String() == "" {
			t.Error("color always: even row style should not be empty")
		}

		oddStyle := sf(1, 0)
		_ = oddStyle
	})

	t.Run("color never all branches", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeNever))
		tbl.SetFooter("T")

		sf := tbl.buildStyleFunc(tbl.footerRowIndex)

		headerStyle := sf(table.HeaderRow, 0)
		if headerStyle.GetBold() {
			t.Error("color never: header should not be bold")
		}

		footerStyle := sf(tbl.footerRowIndex, 0)
		if footerStyle.GetBold() {
			t.Error("color never: footer should not be bold")
		}

		evenStyle := sf(0, 0)
		_ = evenStyle

		oddStyle := sf(1, 0)
		_ = oddStyle
	})

	t.Run("color always with custom footer style", func(t *testing.T) {
		t.Parallel()

		tbl := New(
			WithColorMode(output.ColorModeAlways),
			WithFooterStyle(func(s lipgloss.Style) lipgloss.Style {
				return s.Italic(true)
			}),
		)
		tbl.SetFooter("T")

		sf := tbl.buildStyleFunc(tbl.footerRowIndex)

		footerStyle := sf(tbl.footerRowIndex, 0)
		if !footerStyle.GetItalic() {
			t.Error("custom footer style should be italic")
		}
	})

	t.Run("no footer", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeAlways))
		tbl.SetHeaders("A")

		sf := tbl.buildStyleFunc(0)

		headerStyle := sf(table.HeaderRow, 0)
		_ = headerStyle

		evenStyle := sf(0, 0)
		_ = evenStyle
	})
}
