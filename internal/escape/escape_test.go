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
