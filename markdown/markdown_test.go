package markdown

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestMarkdownTable(t *testing.T) {
	t.Parallel()
	runSubtest(t, "basic table", testMarkdownBasicTable)
	runSubtest(t, "empty headers", testMarkdownEmptyHeaders)
	runSubtest(t, "alignment", testMarkdownAlignment)
	runSubtest(t, "center alignment", testMarkdownCenterAlignment)
	runSubtest(t, "chaining", testMarkdownChaining)
	runSubtest(t, "footer", testMarkdownFooter)
}

func testMarkdownBasicTable(t *testing.T) {
	t.Helper()

	m := newMarkdownTableWithData()

	testhelpers.RenderAssert(t, m, "Name", "Alice")
}

func newMarkdownTableWithSingleRow() *MarkdownTable {
	return newMarkdownTableWithData()
}

func newMarkdownTableWithData() *MarkdownTable {
	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		AddRow([]string{"Alice", "30"}).
		AddRow([]string{"Bob", "25"})

	return m
}

func testMarkdownEmptyHeaders(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable()

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "" {
		t.Error("Render() should return empty string for empty headers")
	}
}

func testMarkdownAlignment(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		SetAlign(1, AlignRight).
		AddRow([]string{"Alice", "30"})

	testhelpers.RenderAssert(t, m, "|--", "--:|")
}

func TestMarkdownAlignmentMarkers(t *testing.T) {
	t.Parallel()

	t.Run("left align marker", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdownTable()
		m.SetHeaders([]string{"A"}).SetAlign(0, AlignLeft).AddRow([]string{"x"})

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if strings.Contains(got, ":--") {
			t.Error("Left align should not have colon prefix in separator")
		}
	})

	t.Run("right align marker", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdownTable()
		m.SetHeaders([]string{"A"}).SetAlign(0, AlignRight).AddRow([]string{"x"})

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, "--:|", "Right align should have colon suffix")
	})

	t.Run("center align marker", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdownTable()
		m.SetHeaders([]string{"A"}).SetAlign(0, AlignCenter).AddRow([]string{"x"})

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, ":--:", "Center align should have colons on both sides")
	})
}

func testMarkdownChaining(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable().
		SetHeaders([]string{"Name"}).
		AddRow([]string{"Test"})

	if m == nil {
		t.Error("Method chaining should return non-nil")
	}
}

func testMarkdownCenterAlignment(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		SetAlign(0, AlignCenter).
		SetAlign(1, AlignCenter).
		AddRow([]string{"A", "30"})

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "|", "Render() should contain pipe delimiters")
	assertContains(t, got, "A", "Render() should contain cell value 'A'")
	assertContains(t, got, "30", "Render() should contain cell value '30'")
}

func TestNewMarkdownTable(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	// Verify table is initialized properly
	_ = m.headers // Just ensure fields are accessible
	_ = m.rows
	_ = m.align
}

func TestNewMarkdownTableFromTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Status"})
	data.AddRow([]string{"Project A", "Active"})
	data.AddRow([]string{"Project B", "Inactive"})

	m := NewMarkdownTableFromTable(data)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "Name", "should contain header 'Name'")
	assertContains(t, got, "Status", "should contain header 'Status'")
	assertContains(t, got, "Project A", "should contain row 'Project A'")
	assertContains(t, got, "Active", "should contain row 'Active'")
}

func TestNewMarkdownTableEmpty(t *testing.T) {
	t.Parallel()

	data := output.NewTable(nil)
	m := NewMarkdownTableFromTable(data)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "" {
		t.Errorf("FromData with empty data should render empty, got %q", got)
	}
}

func TestMarkdownColorModeNever(t *testing.T) {
	t.Parallel()

	m := newMarkdownTableWithSingleRow()
	m.SetColorMode(output.ColorModeNever)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, escape.ANSIBold) {
		t.Error("output.ColorModeNever should not produce ANSI escape codes")
	}

	if strings.Contains(got, escape.ANSIDim) {
		t.Error("output.ColorModeNever should not produce dim ANSI codes")
	}

	assertContains(t, got, "Name", "should contain header text without ANSI")
	assertContains(t, got, "Alice", "should contain data without ANSI")
}

func TestMarkdownColorModeAlways(t *testing.T) {
	t.Parallel()

	m := newMarkdownTableWithSingleRow()
	m.SetColorMode(output.ColorModeAlways)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, escape.ANSIBold, "output.ColorModeAlways should bold headers")
	assertContains(t, got, escape.ANSIReturn, "output.ColorModeAlways should reset after bold")
	assertContains(t, got, escape.ANSIDim, "output.ColorModeAlways should dim separators")
	assertContains(t, got, "Name", "should contain header text")
	assertContains(t, got, "Alice", "should contain data row")
}

func TestMarkdownColorModeDefault(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Col"}).
		AddRow([]string{"Val"})

	testhelpers.RenderAssert(t, m, "Col", "Val")
}

func TestMarkdownTableGetAlignmentOutOfBounds(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"A", "B", "C"}).SetAlign(0, AlignRight)
	m.AddRow([]string{"x", "y", "z"})

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "x", "should contain cell x")
	assertContains(t, got, "y", "should contain cell y")
	assertContains(t, got, "z", "should contain cell z")
}

func TestMarkdownTable_AsTableRenderer(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	tr := m.AsTableRenderer()

	tr.SetHeaders([]string{"A", "B"})
	tr.AddRow([]string{"1", "2"})

	got, err := tr.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "A", "should contain header A")
	assertContains(t, got, "1", "should contain cell 1")

	_ = tr
}

func TestMarkdownTable_EscapesCellContent(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Expr", "Note"})
	m.AddRow([]string{"a|b", "plain"})
	m.AddRow([]string{"line1\nline2", `back\slash`})

	out, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	// Header + separator + 2 data rows = 4 lines: no cell newline may
	// create an extra rendered row.
	if len(lines) != 4 {
		t.Fatalf("escaping must not add rows: got %d lines:\n%s", len(lines), out)
	}

	if !strings.Contains(lines[2], `a\|b`) {
		t.Errorf("pipe not escaped in data row: %q", lines[2])
	}

	if !strings.Contains(lines[3], "line1<br>line2") {
		t.Errorf("newline not escaped as <br>: %q", lines[3])
	}

	if !strings.Contains(lines[3], "back"+`\`+`\`+"slash") {
		t.Errorf("backslash not escaped: %q", lines[3])
	}

	// Each line must contain exactly 3 UNESCAPED pipes (2 columns).
	for i, line := range lines {
		unescaped := strings.Count(line, "|") - strings.Count(line, `\|`)
		if unescaped != 3 {
			t.Errorf("line %d has %d unescaped pipes, want 3 (cell content leaked):\n%s", i, unescaped, out)
		}
	}
}

func TestMarkdownTable_EscapedWidthsStayAligned(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"A", "B"})
	m.AddRow([]string{"x|y", "z"})

	out, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected header/sep/row, got:\n%s", out)
	}

	header, separator, row := lines[0], lines[1], lines[2]

	// Header and data rows share the same per-cell " | " framing, so their
	// raw lengths must match — widths must be computed on the ESCAPED text
	// (x\|y is 4 runes), otherwise the wider cell skews row padding.
	if len(header) != len(row) {
		t.Errorf("header (%d) and row (%d) lengths diverge:\n%s", len(header), len(row), out)
	}

	// Column 0 width = len("x\\|y") = 4: header "A" pads to 4, the escaped
	// cell is followed by exactly one pad space, and the separator run is
	// width+1 = 5 dashes.
	if !strings.Contains(header, "| A    |") {
		t.Errorf("header not padded to escaped width (want A + 4 spaces):\n%s", out)
	}

	if !strings.Contains(row, "| x\\|y |") {
		t.Errorf("row not padded to escaped width (want x\\|y + 1 space):\n%s", out)
	}

	if !strings.Contains(separator, "|-----|") {
		t.Errorf("separator dashes must be escaped-width+1 = 5:\n%s", out)
	}
}
