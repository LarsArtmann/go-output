package escape

import "testing"

func TestHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"no escaping", "hello", "hello"},
		{"angle brackets", "<div>", "&lt;div&gt;"},
		{"ampersand", "a&b", "a&amp;b"},
		{"quotes", `"hello"`, "&quot;hello&quot;"},
		{"apostrophe", "it's", "it&#39;s"},
		{"combined", `<a>it's</a>`, "&lt;a&gt;it&#39;s&lt;/a&gt;"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := HTML(tt.input)
			if got != tt.want {
				t.Errorf("HTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestXML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"no escaping", "hello", "hello"},
		{"angle brackets", "<root>", "&lt;root&gt;"},
		{"ampersand", "a&b", "a&amp;b"},
		{"apostrophe", "it's", "it&apos;s"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := XML(tt.input)
			if got != tt.want {
				t.Errorf("XML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHTMLVsXMLDifference(t *testing.T) {
	t.Parallel()

	input := "it's"
	htmlResult := HTML(input)
	xmlResult := XML(input)

	if htmlResult == xmlResult {
		t.Error("HTML and XML should escape apostrophe differently")
	}

	if htmlResult != "it&#39;s" {
		t.Errorf("HTML apostrophe = %q, want %q", htmlResult, "it&#39;s")
	}

	if xmlResult != "it&apos;s" {
		t.Errorf("XML apostrophe = %q, want %q", xmlResult, "it&apos;s")
	}
}

func TestD2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"no escaping", "hello", "hello"},
		{"quotes", `"hi"`, `\"hi\"`},
		{"newline", "a\nb", `a\nb`},
		{"tab", "a\tb", `a\tb`},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := D2(tt.input)
			if got != tt.want {
				t.Errorf("D2(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDOT(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"no escaping", "hello", "hello"},
		{"quotes", `"hi"`, `\"hi\"`},
		{"newline", "a\nb", `a\nb`},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DOT(tt.input)
			if got != tt.want {
				t.Errorf("DOT(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMermaidID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"alphanumeric", "abc123", "abc123"},
		{"with underscores", "my_node", "my_node"},
		{"has spaces", "has spaces", "hasspaces"},
		{"special!@#$chars", "special!@#$chars", "specialchars"},
		{"empty", "", "node"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MermaidID(tt.input)
			if got != tt.want {
				t.Errorf("MermaidID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMermaidSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"simple", "simple", "simple"},
		{"has spaces", "has spaces", "has_spaces"},
		{"has-dash", "has-dash", "has_dash"},
		{"path/to/file", "path/to/file", "path_to_file"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MermaidSlug(tt.input)
			if got != tt.want {
				t.Errorf("MermaidSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMermaidText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, input, want string
	}{
		{"no escaping", "hello", "hello"},
		{"quotes", `"hi"`, `'hi'`},
		{"brackets", "a[b]c", "a(b)c"},
		{"braces", "a{b}c", "a(b)c"},
		{"newline", "a\nb", "a<br>b"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MermaidText(tt.input)
			if got != tt.want {
				t.Errorf("MermaidText(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
