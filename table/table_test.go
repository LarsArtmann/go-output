package table

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

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
	if !strings.Contains(output, "Name") {
		t.Error("Render() should contain header 'Name'")
	}

	if !strings.Contains(output, "Value") {
		t.Error("Render() should contain header 'Value'")
	}

	if !strings.Contains(output, "Count") {
		t.Error("Render() should contain header 'Count'")
	}
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
	if !strings.Contains(output, "Alice") {
		t.Error("Render() should contain row value 'Alice'")
	}

	if !strings.Contains(output, "30") {
		t.Error("Render() should contain row value '30'")
	}
}

func TestAddRowMultiple(t *testing.T) {
	t.Parallel()

	tbl := New()
	tbl.SetHeaders("Name", "Score")
	tbl.AddRow("Alice", "100")
	tbl.AddRow("Bob", "90")
	tbl.AddRow("Charlie", "85")

	output := tbl.Render()
	if !strings.Contains(output, "Alice") {
		t.Error("Render() should contain 'Alice'")
	}

	if !strings.Contains(output, "Bob") {
		t.Error("Render() should contain 'Bob'")
	}

	if !strings.Contains(output, "Charlie") {
		t.Error("Render() should contain 'Charlie'")
	}
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
	if !strings.Contains(output, "Test") {
		t.Error("Render() should contain 'Test'")
	}
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

	if !strings.Contains(output, "Name") || !strings.Contains(output, "Status") {
		t.Error("Render() should contain headers")
	}

	if !strings.Contains(output, "Project A") || !strings.Contains(output, "Active") {
		t.Error("Render() should contain row data")
	}
}

func TestChaining(t *testing.T) {
	t.Parallel()

	output := New().
		SetHeaders("ID", "Name").
		AddRow("1", "First").
		AddRow("2", "Second").
		Render()

	if !strings.Contains(output, "ID") || !strings.Contains(output, "Name") {
		t.Error("Chained Render() should contain headers")
	}

	if !strings.Contains(output, "First") || !strings.Contains(output, "Second") {
		t.Error("Chained Render() should contain row data")
	}
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

	if !strings.Contains(output, "Only") || !strings.Contains(output, "Headers") {
		t.Error("Render() should contain headers even without rows")
	}
}
