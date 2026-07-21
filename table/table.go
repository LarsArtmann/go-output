// Package table provides terminal table output formatting using lipgloss.
package table

import (
	"io"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"

	"github.com/larsartmann/go-output"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*Table)(nil)
	_ output.TableRenderer = (*tableRendererAdapter)(nil)
)

//nolint:gochecknoinits // Registers Table TableRenderer for registry-based dispatch.
func init() {
	output.RegisterTableMarshaler(output.FormatTable, renderStyledTable)
}

func renderStyledTable(w io.Writer, data *output.Table, opts output.RenderOptions) error {
	t := FromTable(data, WithColorMode(opts.ColorMode))

	return output.WriteRenderedFrom(w, t.Render, "table", "render table")
}

// TableProvider defines the interface for types that provide tabular data.
// The root package's Table satisfies this interface implicitly.
type TableProvider interface {
	GetHeaders() []string
	GetRows() [][]string
}

// FooterProvider is an optional interface that TableProvider can implement
// to provide a footer row (e.g., totals). Checked via type assertion in FromTable.
type FooterProvider interface {
	GetFooter() []string
}

// Option configures a Table during construction.
type Option func(*Table)

// WithColorMode sets the color mode for the table renderer.
// ColorModeAuto (default) enables colors when stdout is a terminal.
// ColorModeNever disables all ANSI styling.
// ColorModeAlways forces ANSI styling regardless of terminal state.
func WithColorMode(mode output.ColorMode) Option {
	return func(t *Table) { t.colorMode = mode }
}

// WithFooterStyle sets a custom style function for the footer row.
// The provided function receives a base lipgloss.Style (with padding)
// and returns the styled result. When set, this overrides the default
// bold footer styling.
func WithFooterStyle(fn func(lipgloss.Style) lipgloss.Style) Option {
	return func(t *Table) { t.footerStyleFn = fn }
}

// Table renders formatted tables using lipgloss.
type Table struct {
	t              *table.Table
	colorMode      output.ColorMode
	rowCount       int
	footerRowIndex int
	footerStyleFn  func(lipgloss.Style) lipgloss.Style
}

// New creates a new Table with default styling.
// Pass WithColorMode to control ANSI color output (default: ColorModeAuto).
func New(opts ...Option) *Table {
	tbl := &Table{colorMode: output.ColorModeAuto}

	for _, opt := range opts {
		opt(tbl)
	}

	tbl.t = table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(tbl.buildStyleFunc(0))

	return tbl
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
	t.rowCount++

	return t.apply(func() { t.t.Row(row...) })
}

// SetFooter adds a bold-styled footer row to the table.
// Only the last footer row receives bold styling; previous footer rows
// become unstyled data rows. Use SetFooter once for a single summary row.
func (t *Table) SetFooter(row ...string) *Table {
	t.footerRowIndex = t.rowCount + 1
	t.rowCount++

	t.t.Row(row...)
	t.StyleFunc(t.buildStyleFunc(t.footerRowIndex))

	return t
}

// StyleFunc sets a custom style function.
func (t *Table) StyleFunc(fn func(row, col int) lipgloss.Style) *Table {
	return t.apply(func() { t.t.StyleFunc(fn) })
}

// Render returns the rendered table string.
func (t *Table) Render() (string, error) {
	return t.t.String(), nil
}

// FromTable creates a new Table populated from a TableProvider.
// If data is nil, returns an empty table.
// If data also implements FooterProvider, the footer row is added and bold-styled.
func FromTable(data TableProvider, opts ...Option) *Table {
	if data == nil {
		return New(opts...)
	}

	t := New(opts...)
	t.SetHeaders(data.GetHeaders()...)

	for _, row := range data.GetRows() {
		t.AddRow(row...)
	}

	if fp, ok := data.(FooterProvider); ok {
		footer := fp.GetFooter()
		if len(footer) > 0 {
			t.SetFooter(footer...)
		}
	}

	return t
}

// buildStyleFunc returns a StyleFunc that applies header, footer, and alternating row styles.
// footerRow > 0 enables bold footer styling on that row index.
func (t *Table) buildStyleFunc(footerRow int) func(row, col int) lipgloss.Style {
	useColor := t.colorMode.ShouldColor()
	baseStyle := lipgloss.NewStyle().Padding(0, 1)

	return func(row, _ int) lipgloss.Style {
		if row == table.HeaderRow {
			style := baseStyle
			if useColor {
				style = style.Foreground(lipgloss.Color("99")).Bold(true)
			}

			return style
		}

		if footerRow > 0 && row == footerRow {
			style := baseStyle
			if useColor {
				style = style.Bold(true)
			}

			if t.footerStyleFn != nil {
				style = t.footerStyleFn(style)
			}

			return style
		}

		if row%2 == 0 {
			style := baseStyle
			if useColor {
				style = style.Foreground(lipgloss.Color("245"))
			}

			return style
		}

		return baseStyle
	}
}

// tableRendererAdapter wraps a Table to satisfy the output.TableRenderer interface.
// It adapts the variadic API (...string) to the slice-based TableRenderer methods.
type tableRendererAdapter struct {
	inner *Table
}

func (a *tableRendererAdapter) Render() (string, error)     { return a.inner.Render() }
func (a *tableRendererAdapter) SetHeaders(headers []string) { a.inner.SetHeaders(headers...) }
func (a *tableRendererAdapter) AddRow(row []string)         { a.inner.AddRow(row...) }

// AsTableRenderer returns a TableRenderer that delegates to this Table.
// This adapts the variadic API (...string) to the slice-based TableRenderer interface.
func (t *Table) AsTableRenderer() output.TableRenderer {
	return &tableRendererAdapter{inner: t}
}
