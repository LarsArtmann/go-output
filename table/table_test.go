package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// assertTableContains checks that output contains substr, failing with msg if not.
func assertTableContains(t *testing.T, output, substr, msg string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Error(msg)
	}
}

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

	output := tbl.Render()
	assertTableContains(t, output, "Name", "Render() should contain header 'Name'")
	assertTableContains(t, output, "Value", "Render() should contain header 'Value'")
	assertTableContains(t, output, "Count", "Render() should contain header 'Count'")
}

func TestAddRow(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Value")
	result := tbl.AddRow("Alice", "30")

	if result != tbl {
		t.Error("AddRow() should return the same table for chaining")
	}

	output := tbl.Render()
	assertTableContains(t, output, "Alice", "Render() should contain row value 'Alice'")
	assertTableContains(t, output, "30", "Render() should contain row value '30'")
}

func TestAddRowMultiple(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Score")
	tbl.AddRow("Alice", "100")
	tbl.AddRow("Bob", "90")
	tbl.AddRow("Charlie", "85")

	output := tbl.Render()
	assertTableContains(t, output, "Alice", "Render() should contain 'Alice'")
	assertTableContains(t, output, "Bob", "Render() should contain 'Bob'")
	assertTableContains(t, output, "Charlie", "Render() should contain 'Charlie'")
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

	output := tbl.Render()
	assertTableContains(t, output, "Test", "Render() should contain 'Test'")
}

func TestRender(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Status")
	tbl.AddRow("Project A", "Active")
	tbl.AddRow("Project B", "Inactive")

	output := tbl.Render()

	if output == "" {
		t.Error("Render() should not return empty string")
	}

	assertTableContains(t, output, "Name", "Render() should contain headers")
	assertTableContains(t, output, "Status", "Render() should contain headers")
	assertTableContains(t, output, "Project A", "Render() should contain row data")
	assertTableContains(t, output, "Active", "Render() should contain row data")
}

func TestChaining(t *testing.T) {
	t.Parallel()

	output := New().
		SetHeaders("ID", "Name").
		AddRow("1", "First").
		AddRow("2", "Second").
		Render()

	assertTableContains(t, output, "ID", "Chained Render() should contain headers")
	assertTableContains(t, output, "Name", "Chained Render() should contain headers")
	assertTableContains(t, output, "First", "Chained Render() should contain row data")
	assertTableContains(t, output, "Second", "Chained Render() should contain row data")
}

func TestEmptyTable(t *testing.T) {
	t.Parallel()

	tbl := New()
	output := tbl.Render()

	// Empty table (no headers, no rows) returns empty string
	if output != "" {
		t.Errorf("Render() of empty table should be empty, got %q", output)
	}
}

func TestHeadersOnlyNoRows(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Only", "Headers")
	output := tbl.Render()

	assertTableContains(t, output, "Only", "Render() should contain headers even without rows")
	assertTableContains(t, output, "Headers", "Render() should contain headers even without rows")
}
