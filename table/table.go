package table

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

type Table struct {
	t *table.Table
}

func New() *Table {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
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

func (t *Table) SetHeaders(headers ...string) *Table {
	t.t.Headers(headers...)
	return t
}

func (t *Table) AddRow(row ...string) *Table {
	t.t.Row(row...)
	return t
}

func (t *Table) StyleFunc(fn func(row, col int) lipgloss.Style) *Table {
	t.t.StyleFunc(fn)
	return t
}

func (t *Table) Render() string {
	return t.t.String()
}
