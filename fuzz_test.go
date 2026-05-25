package output

import (
	"strings"
	"testing"
)

// stringEnum is a constraint for string-based enum types used in fuzz testing.
type stringEnum interface {
	~string
	IsValid() bool
}

func fuzzEnumTest[E stringEnum](
	t *testing.T,
	s string,
	parse func(string) (E, error),
	typeName string,
) {
	result, err := parse(s)
	if err != nil {
		if result != "" {
			t.Errorf("%s(%q) returned error but non-empty result: %q", typeName, s, result)
		}
	}

	if result.IsValid() && err == nil {
		if string(result) != s {
			t.Errorf("%s(%q) = %q, but IsValid() was true", typeName, s, result)
		}
	}
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
