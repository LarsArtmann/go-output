package output

import (
	"bytes"
	"strings"
	"testing"
)

func FuzzCSVWriter(f *testing.F) {
	// Seed corpus with common cases
	f.Add("Name,Value", "Alice,100")
	f.Add("Header1,Header2,Header3", "a,b,c")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, headerLine, dataLine string) {
		// Parse the CSV-like input
		headers := parseCSVFuzz(headerLine)
		row := parseCSVFuzz(dataLine)

		// Ensure at least some data
		if len(headers) == 0 && len(row) == 0 {
			return
		}

		// Test CSV writing
		var buf bytes.Buffer
		w := NewCSVWriter(&buf)

		if len(headers) > 0 {
			_ = w.WriteHeader(headers)
		}
		if len(row) > 0 {
			_ = w.WriteRow(row)
		}
		w.Flush()

		// Should not panic and should produce output
		result := buf.String()
		if result == "" && (len(headers) > 0 || len(row) > 0) {
			t.Error("CSVWriter produced empty output when it shouldn't")
		}
	})
}

// parseCSVFuzz parses a simple CSV line for fuzz testing
func parseCSVFuzz(line string) []string {
	if line == "" {
		return nil
	}
	return strings.Split(line, ",")
}

func FuzzMarkdownTable(f *testing.F) {
	// Seed corpus
	f.Add("Name", "Alice")
	f.Add("Col1,Col2", "a,b")
	f.Add("Header1,Header2,Header3,Header4", "v1,v2,v3,v4")

	f.Fuzz(func(t *testing.T, headerLine, dataLine string) {
		headers := strings.Split(headerLine, ",")
		row := strings.Split(dataLine, ",")

		// Skip if both are empty
		if len(headers) == 0 && len(row) == 0 {
			return
		}

		// Test Markdown table rendering
		md := NewMarkdownTable()
		md.SetHeaders(headers)
		md.AddRow(row)
		result := md.Render()

		// Should not panic and should produce output
		if result == "" && len(headers) > 0 {
			t.Error("MarkdownTable produced empty output with headers")
		}
	})
}
