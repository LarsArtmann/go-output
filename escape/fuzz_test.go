package escape

import (
	"strings"
	"testing"
)

// assertEscaped fails the test if the input contains `input` but the result
// does not contain its escaped form. Used by fuzz tests to assert escape
// invariants in a single line.
func assertEscaped(t *testing.T, fn, s, result, input, output string) {
	t.Helper()

	if strings.Contains(s, input) && !strings.Contains(result, output) {
		t.Errorf("%s(%q) = %q, %s not escaped", fn, s, result, input)
	}
}

// assertStripped fails the test if the input contains `ch` and the result
// also contains it — i.e. the raw character survived transformation. Used
// by fuzz tests that verify characters are rewritten (not just escaped) by
// the function under test.
func assertStripped(t *testing.T, fn, s, result, ch string) {
	t.Helper()

	if strings.Contains(s, ch) && strings.Contains(result, ch) {
		t.Errorf("%s(%q) = %q, %s not stripped", fn, s, result, ch)
	}
}

func FuzzD2(f *testing.F) {
	f.Add("")
	f.Add(`hello "world"`)
	f.Add(`back\slash`)
	f.Add("new\nline\ttab")
	f.Add(`mixed "quote\ with\n newlines"`)

	f.Fuzz(func(t *testing.T, s string) {
		result := D2(s)

		assertEscaped(t, "D2", s, result, `\`, `\\`)
		assertEscaped(t, "D2", s, result, `"`, `\"`)
	})
}

func FuzzXML(f *testing.F) {
	f.Add("")
	f.Add("<tag>content</tag>")
	f.Add(`it's "quoted" & <encoded>`)
	f.Add("normal text")

	f.Fuzz(func(t *testing.T, s string) {
		result := XML(s)

		assertEscaped(t, "XML", s, result, "<", "&lt;")
		assertEscaped(t, "XML", s, result, "&", "&amp;")
	})
}

func FuzzHTML(f *testing.F) {
	f.Add("")
	f.Add("<script>alert('xss')</script>")
	f.Add(`it's "quoted" & <encoded>`)

	f.Fuzz(func(t *testing.T, s string) {
		result := HTML(s)

		assertEscaped(t, "HTML", s, result, "<", "&lt;")
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

		assertStripped(t, "MermaidText", s, result, `"`)
		assertStripped(t, "MermaidText", s, result, "[")
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
