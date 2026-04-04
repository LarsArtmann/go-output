// Package table provides terminal table output formatting using lipgloss.
package table

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// Table renders formatted tables using lipgloss.
type Table struct {
	t *table.Table
}

// New creates a new Table with default styling.
func New() *Table {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("99")).
					Bold(true).
					Padding(0, 1)
			}

			if row%2 == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("245")).
					Padding(0, 1)
			}

			return lipgloss.NewStyle().
				Padding(0, 1)
		})

	return &Table{t: t}
}

// apply executes fn on the underlying table and returns self for chaining.
func (t *Table) apply(fn func()) *Table {
	fn()

	return t
}

// SetHeaders sets the table headers.
func (t *Table) SetHeaders(headers ...string) *Table {
	return t.apply(func() { t.t.Headers(headers...) })
}

// AddRow adds a row to the table.
func (t *Table) AddRow(row ...string) *Table {
	return t.apply(func() { t.t.Row(row...) })
}

// StyleFunc sets a custom style function.
func (t *Table) StyleFunc(fn func(row, col int) lipgloss.Style) *Table {
	return t.apply(func() { t.t.StyleFunc(fn) })
}

// Render returns the rendered table string.
func (t *Table) Render() string {
	return t.t.String()
}
