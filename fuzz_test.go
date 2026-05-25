package output

import (
	"strings"
	"testing"
)

func FuzzParseColorMode(f *testing.F) {
	f.Add("auto")
	f.Add("always")
	f.Add("never")
	f.Add("AUTO")
	f.Add("")
	f.Add("invalid")

	f.Fuzz(func(t *testing.T, input string) {
		mode, err := ParseColorMode(input)
		if err != nil {
			return
		}

		if !mode.IsValid() {
			t.Errorf("ParseColorMode(%q) returned invalid mode %q", input, mode)
		}

		if mode.String() != string(mode) {
			t.Errorf("ColorMode.String() = %q, want %q", mode.String(), string(mode))
		}
	})
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

		result, err := md.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		// Should not panic and should produce output
		if result == "" && len(headers) > 0 {
			t.Error("MarkdownTable produced empty output with headers")
		}
	})
}
