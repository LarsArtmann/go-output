package markup

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func BenchmarkMarshalXMLFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalXMLFromTableData(data)
	}
}

func BenchmarkXMLWriter(b *testing.B) {
	headers := []string{"Name", "Age", "Email"}
	rows := make([][]string, 100)
	for i := range rows {
		rows[i] = []string{"Alice", "30", "alice@example.com"}
	}

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer
		w := NewXMLWriter(&buf)
		_ = w.WriteHeader(headers)
		_ = w.WriteRows(rows)
		_ = w.WriteFooter()
	}
}

func BenchmarkHTMLRenderer(b *testing.B) {
	r := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		r.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		renderer := NewHTMLRenderer()
		renderer.SetHeaders(r.GetHeaders())
		for _, row := range r.GetRows() {
			renderer.AddRow(row)
		}
		_, _ = renderer.Render()
	}
}

func BenchmarkStreamingHTMLRenderer(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		renderer := NewStreamingHTMLRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	}
}

func BenchmarkStreamingHTMLRenderer_Stream(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age"})
	for range 100 {
		data.AddRow([]string{"Alice", "30"})
	}

	b.ResetTimer()

	for b.Loop() {
		var buf strings.Builder
		renderer := NewStreamingHTMLRenderer()
		renderer.SetData(data)
		_ = renderer.Stream(&buf)
	}
}
