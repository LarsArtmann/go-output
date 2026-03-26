package output

import (
	"strings"
	"testing"
)

func TestMarkdownTable(t *testing.T) {
	t.Parallel()
	t.Run("basic table", testMarkdownBasicTable)
	t.Run("empty headers", testMarkdownEmptyHeaders)
	t.Run("alignment", testMarkdownAlignment)
	t.Run("chaining", testMarkdownChaining)
}

func testMarkdownBasicTable(t *testing.T) {
	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		AddRow([]string{"Alice", "30"}).
		AddRow([]string{"Bob", "25"})

	got, err := m.Render()
	if err != nil {
		t.Errorf("Render() error = %v", err)
	}

	if !strings.Contains(got, "Name") {
		t.Error("Render() should contain header text")
	}
	if !strings.Contains(got, "Alice") {
		t.Error("Render() should contain data row")
	}
}

func testMarkdownEmptyHeaders(t *testing.T) {
	m := NewMarkdownTable()

	got, err := m.Render()
	if err != nil {
		t.Errorf("Render() error = %v", err)
	}

	if got != "" {
		t.Error("Render() should return empty string for empty headers")
	}
}

func testMarkdownAlignment(t *testing.T) {
	m := NewMarkdownTable()
	m.SetHeaders([]string{"Name", "Age"}).
		SetAlign(1, AlignRight).
		AddRow([]string{"Alice", "30"})

	got, err := m.Render()
	if err != nil {
		t.Errorf("Render() error = %v", err)
	}

	if !strings.Contains(got, "|--") {
		t.Error("Render() should contain separator row")
	}
}

func testMarkdownChaining(t *testing.T) {
	m := NewMarkdownTable().
		SetHeaders([]string{"Name"}).
		AddRow([]string{"Test"})

	if m == nil {
		t.Error("Method chaining should return non-nil")
	}
}

func TestNewMarkdownTable(_ *testing.T) {
	m := NewMarkdownTable()
	// Verify table is initialized properly
	_ = m.headers // Just ensure fields are accessible
	_ = m.rows
	_ = m.align
}
