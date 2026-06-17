package escape

import (
	"strings"
	"testing"
)

func FuzzD2(f *testing.F) {
	f.Add("")
	f.Add(`hello "world"`)
	f.Add(`back\slash`)
	f.Add("new\nline\ttab")
	f.Add(`mixed "quote\ with\n newlines"`)

	f.Fuzz(func(t *testing.T, s string) {
		result := D2(s)

		if strings.Contains(s, `\`) && !strings.Contains(result, `\\`) {
			t.Errorf("D2(%q) = %q, backslash not escaped", s, result)
		}

		if strings.Contains(s, `"`) && !strings.Contains(result, `\"`) {
			t.Errorf("D2(%q) = %q, quote not escaped", s, result)
		}
	})
}

func FuzzXML(f *testing.F) {
	f.Add("")
	f.Add("<tag>content</tag>")
	f.Add(`it's "quoted" & <encoded>`)
	f.Add("normal text")

	f.Fuzz(func(t *testing.T, s string) {
		result := XML(s)

		if strings.Contains(s, "<") && !strings.Contains(result, "&lt;") {
			t.Errorf("XML(%q) = %q, < not escaped", s, result)
		}

		if strings.Contains(s, "&") && !strings.Contains(result, "&amp;") {
			t.Errorf("XML(%q) = %q, & not escaped", s, result)
		}
	})
}

func FuzzHTML(f *testing.F) {
	f.Add("")
	f.Add("<script>alert('xss')</script>")
	f.Add(`it's "quoted" & <encoded>`)

	f.Fuzz(func(t *testing.T, s string) {
		result := HTML(s)

		if strings.Contains(s, "<") && !strings.Contains(result, "&lt;") {
			t.Errorf("HTML(%q) = %q, < not escaped", s, result)
		}
	})
}

func FuzzMermaidID(f *testing.F) {
	f.Add("")
	f.Add("simple-id")
	f.Add(" spaces & symbols!")
	f.Add("123_abc")

	f.Fuzz(func(t *testing.T, s string) {
		result := MermaidID(s)

		if result == "" {
			t.Errorf("MermaidID(%q) returned empty string (should be 'node' for empty)", s)
		}

		for _, r := range result {
			if !isMermaidIdentRune(r) {
				t.Errorf("MermaidID(%q) = %q, invalid rune %q", s, result, r)
			}
		}
	})
}

func FuzzMermaidText(f *testing.F) {
	f.Add("")
	f.Add(`"quoted"`)
	f.Add("with [brackets] and {braces}")
	f.Add("line1\nline2")

	f.Fuzz(func(t *testing.T, s string) {
		result := MermaidText(s)

		if strings.Contains(s, `"`) && strings.Contains(result, `"`) {
			t.Errorf("MermaidText(%q) = %q, quotes not escaped", s, result)
		}

		if strings.Contains(s, "[") && strings.Contains(result, "[") {
			t.Errorf("MermaidText(%q) = %q, brackets not escaped", s, result)
		}
	})
}

func FuzzSlugifyID(f *testing.F) {
	f.Add("")
	f.Add("simple name")
	f.Add("with-hyphens")
	f.Add("path/to/file")
	f.Add("mix of all-three/types here")

	f.Fuzz(func(t *testing.T, s string) {
		result := SlugifyID(s)

		if strings.Contains(result, " ") {
			t.Errorf("SlugifyID(%q) = %q, contains spaces", s, result)
		}

		if strings.Contains(result, "-") {
			t.Errorf("SlugifyID(%q) = %q, contains hyphens", s, result)
		}

		if strings.Contains(result, "/") {
			t.Errorf("SlugifyID(%q) = %q, contains slashes", s, result)
		}
	})
}
