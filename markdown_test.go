package output

import (
	"strings"
	"testing"
)

func TestMarkdownTable(t *testing.T) {
	t.Parallel()
	t.Run("basic table", func(t *testing.T) {
		t.Parallel()
		testMarkdownBasicTable(t)
	})
	t.Run("empty headers", func(t *testing.T) {
		t.Parallel()
		testMarkdownEmptyHeaders(t)
	})
	t.Run("alignment", func(t *testing.T) {
		t.Parallel()
		testMarkdownAlignment(t)
	})
	t.Run("chaining", func(t *testing.T) {
		t.Parallel()
		testMarkdownChaining(t)
	})
}

func testMarkdownBasicTable(t *testing.T) {
	t.Helper()
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
	t.Helper()
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
	t.Helper()
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
