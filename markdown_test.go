package output

import (
	"testing"
)

func TestMarkdownTable(t *testing.T) {
	t.Parallel()
	runSubtest(t, "basic table", testMarkdownBasicTable)
	runSubtest(t, "empty headers", testMarkdownEmptyHeaders)
	runSubtest(t, "alignment", testMarkdownAlignment)
	runSubtest(t, "chaining", testMarkdownChaining)
}

func testMarkdownBasicTable(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		AddRow([]string{"Alice", "30"}).
		AddRow([]string{"Bob", "25"})

	got := m.Render()

	assertContains(t, got, "Name", "Render() should contain header text")
	assertContains(t, got, "Alice", "Render() should contain data row")
}

func testMarkdownEmptyHeaders(t *testing.T) {
	t.Helper()

	m := NewMarkdownTable()

	got := m.Render()

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

	got := m.Render()

	assertContains(t, got, "|--", "Render() should contain separator row")
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
	got := m.Render()

	assertContains(t, got, "Name", "should contain header 'Name'")
	assertContains(t, got, "Status", "should contain header 'Status'")
	assertContains(t, got, "Project A", "should contain row 'Project A'")
	assertContains(t, got, "Active", "should contain row 'Active'")
}

func TestNewMarkdownTableFromDataEmpty(t *testing.T) {
	t.Parallel()

	data := NewTableData(nil)
	m := NewMarkdownTableFromData(data)

	got := m.Render()
	if got != "" {
		t.Errorf("FromData with empty data should render empty, got %q", got)
	}
}
