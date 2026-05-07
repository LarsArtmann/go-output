package escape

import "testing"

// escapeTestCase represents a single escape function test case.
type escapeTestCase struct {
	name  string
	input string
	want  string
}

// Common test cases for escape functions that don't modify "hello" and empty strings.
var commonEscapeTests = []escapeTestCase{
	{"no escaping", "hello", "hello"},
	{"empty", "", ""},
}

// Predefined test case groups for common escape scenarios.
var (
	xmlSpecificTests = xmlTestCases()
	d2EscapeTests    = d2TestCases()
	dotEscapeTests   = dotTestCases()
	mermaidTextTests = mermaidTextTestCases()
)

// Individual test case builders to avoid structural duplication.
func xmlTestCases() []escapeTestCase {
	return []escapeTestCase{
		{"angle brackets", "<root>", "&lt;root&gt;"},
		{"ampersand", "a&b", "a&amp;b"},
		{"apostrophe", "it's", "it&apos;s"},
	}
}

func d2TestCases() []escapeTestCase {
	return []escapeTestCase{
		{"quotes", `"hi"`, `\"hi\"`},
		{"newline", "a\nb", `a\nb`},
		{"tab", "a\tb", `a\tb`},
		{"backslash", `a\b`, `a\\b`},
	}
}

func dotTestCases() []escapeTestCase {
	return []escapeTestCase{
		{"quotes", `"hi"`, `\"hi\"`},
		{"newline", "a\nb", `a\nb`},
		{"backslash", `a\b`, `a\\b`},
	}
}

func mermaidTextTestCases() []escapeTestCase {
	return []escapeTestCase{
		{"quotes", `"hi"`, `'hi'`},
		{"brackets", "a[b]c", "a(b)c"},
		{"braces", "a{b}c", "a(b)c"},
		{"newline", "a\nb", "a<br>b"},
	}
}

func mermaidIDTestCases() []escapeTestCase {
	return []escapeTestCase{
		{"alphanumeric", "abc123", "abc123"},
		{"with underscores", "my_node", "my_node"},
		{"has spaces", "has spaces", "hasspaces"},
		{"special!@#$chars", "special!@#$chars", "specialchars"},
		{"empty returns node", "", "node"},
	}
}

var mermaidSlugTestData = []escapeTestCase{
	{"simple text", "simple", "simple"},
	{"has spaces", "has spaces", "has_spaces"},
	{"has-dash", "has-dash", "has_dash"},
	{"path/to/file", "path/to/file", "path_to_file"},
	{"empty", "", ""},
}

func mermaidSlugTestCases() []escapeTestCase {
	return mermaidSlugTestData
}

// testEscapeFunc is a helper that runs table-driven tests for escape functions.
func testEscapeFunc(t *testing.T, fnName string, fn func(string) string, tests []escapeTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := fn(tt.input)
			if got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", fnName, tt.input, got, tt.want)
			}
		})
	}
}

func TestHTML(t *testing.T) {
	t.Parallel()

	tests := []escapeTestCase{
		{"no escaping", "hello", "hello"},
		{"angle brackets", "<div>", "&lt;div&gt;"},
		{"ampersand", "a&b", "a&amp;b"},
		{"quotes", `"hello"`, "&#34;hello&#34;"},
		{"apostrophe", "it's", "it&#39;s"},
		{"combined", "<a>it's</a>", "&lt;a&gt;it&#39;s&lt;/a&gt;"},
		{"empty", "", ""},
	}

	testEscapeFunc(t, "HTML", HTML, tests)
}

func TestXML(t *testing.T) {
	t.Parallel()

	tests := make([]escapeTestCase, 0, len(commonEscapeTests)+len(xmlSpecificTests))
	tests = append(tests, commonEscapeTests...)
	tests = append(tests, xmlSpecificTests...)
	testEscapeFunc(t, "XML", XML, tests)
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

	tests := make([]escapeTestCase, 0, len(commonEscapeTests)+len(d2EscapeTests))
	tests = append(tests, commonEscapeTests...)
	tests = append(tests, d2EscapeTests...)
	testEscapeFunc(t, "D2", D2, tests)
}

func TestDOT(t *testing.T) {
	t.Parallel()

	tests := make([]escapeTestCase, 0, len(commonEscapeTests)+len(dotEscapeTests))
	tests = append(tests, commonEscapeTests...)
	tests = append(tests, dotEscapeTests...)
	testEscapeFunc(t, "DOT", DOT, tests)
}

func TestMermaidID(t *testing.T) {
	t.Parallel()

	testEscapeFunc(t, "MermaidID", MermaidID, mermaidIDTestCases())
}

func TestMermaidSlug(t *testing.T) {
	t.Parallel()

	testEscapeFunc(t, "MermaidSlug", MermaidSlug, mermaidSlugTestCases())
}

func TestMermaidText(t *testing.T) {
	t.Parallel()

	tests := make([]escapeTestCase, 0, len(commonEscapeTests)+len(mermaidTextTests))
	tests = append(tests, commonEscapeTests...)
	tests = append(tests, mermaidTextTests...)
	testEscapeFunc(t, "MermaidText", MermaidText, tests)
}
