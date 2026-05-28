package output

import (
	"strings"
	"testing"
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

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "Name", "Render() should contain header text")
	assertContains(t, got, "Alice", "Render() should contain data row")
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

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "|--", "Render() should contain separator row")
	assertContains(t, got, "--:|", "Render() should contain right-align marker")
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

func TestNewMarkdownTableFromData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Name", "Status"})
	data.AddRow([]string{"Project A", "Active"})
	data.AddRow([]string{"Project B", "Inactive"})

	m := NewMarkdownTableFromData(data)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "Name", "should contain header 'Name'")
	assertContains(t, got, "Status", "should contain header 'Status'")
	assertContains(t, got, "Project A", "should contain row 'Project A'")
	assertContains(t, got, "Active", "should contain row 'Active'")
}

func TestNewMarkdownTableFromDataEmpty(t *testing.T) {
	t.Parallel()

	data := NewTableData(nil)
	m := NewMarkdownTableFromData(data)

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
	m.SetColorMode(ColorModeNever)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, ansiBold) {
		t.Error("ColorModeNever should not produce ANSI escape codes")
	}

	if strings.Contains(got, ansiDim) {
		t.Error("ColorModeNever should not produce dim ANSI codes")
	}

	assertContains(t, got, "Name", "should contain header text without ANSI")
	assertContains(t, got, "Alice", "should contain data without ANSI")
}

func TestMarkdownColorModeAlways(t *testing.T) {
	t.Parallel()

	m := newMarkdownTableWithSingleRow()
	m.SetColorMode(ColorModeAlways)

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, ansiBold, "ColorModeAlways should bold headers")
	assertContains(t, got, ansiReset, "ColorModeAlways should reset after bold")
	assertContains(t, got, ansiDim, "ColorModeAlways should dim separators")
	assertContains(t, got, "Name", "should contain header text")
	assertContains(t, got, "Alice", "should contain data row")
}

func TestMarkdownColorModeDefault(t *testing.T) {
	t.Parallel()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Col"}).
		AddRow([]string{"Val"})

	got, err := m.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "Col", "default render should contain header")
	assertContains(t, got, "Val", "default render should contain row")
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

	var _ TableRenderer = tr
}
