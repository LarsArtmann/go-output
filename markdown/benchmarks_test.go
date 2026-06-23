package markdown

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func BenchmarkMarkdownTableColored(b *testing.B) {
	md := NewMarkdownTable()

	const (
		headerCell = "Header"
		dataCell   = "Cell"
	)

	headers := make([]string, 0, 10)
	for range 10 {
		headers = append(headers, headerCell)
	}

	md.SetHeaders(headers)
	md.SetColorMode(output.ColorModeAlways)

	for range 100 {
		row := make([]string, 0, 10)
		for range 10 {
			row = append(row, dataCell)
		}

		md.AddRow(row)
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = md.Render()
	}
}

func BenchmarkMarkdownTableWithFooter(b *testing.B) {
	md := NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Age", "Email", "City"})
	md.SetColorMode(output.ColorModeAlways)

	for range 100 {
		md.AddRow([]string{"Alice", "30", "alice@example.com", "Berlin"})
	}

	md.SetFooter([]string{"Total", "100", "", ""})

	b.ResetTimer()

	for b.Loop() {
		_, _ = md.Render()
	}
}
