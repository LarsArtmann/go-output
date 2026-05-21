package gentest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-faster/yaml"
)

// ExpectedOutput contains a substring to check and its corresponding error message.
type ExpectedOutput struct {
	Substring string
	Message   string
}

// AssertOutputContains checks that output contains substr, failing with a descriptive error.
func AssertOutputContains(t *testing.T, output, substr string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Errorf("output should contain %q, got %q", substr, output)
	}
}

// AssertValidJSON checks that output is valid JSON by unmarshaling it.
func AssertValidJSON(t *testing.T, output string) {
	t.Helper()

	var result map[string]any
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid JSON: %v, got %q", err, output)
	}
}

// AssertValidYAML checks that output is valid YAML by unmarshaling it.
func AssertValidYAML(t *testing.T, output string) {
	t.Helper()

	var result map[string]any
	if err := yaml.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid YAML: %v, got %q", err, output)
	}
}

// HTMLEscapeTestRenderer is an interface for HTML renderers that support escaping tests.
type HTMLEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() (string, error)
}

// AssertHTMLEscape verifies that a renderer properly escapes HTML content.
func AssertHTMLEscape(t *testing.T, newRenderer func() HTMLEscapeTestRenderer, name string) {
	t.Helper()

	r := newRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"<script>alert('xss')</script>"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "<script>") {
		t.Errorf("%s: Render() should escape script tags", name)
	}

	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("%s: Render() should contain escaped script tag", name)
	}
}

// AssertMarshalError checks that a marshal function returns the expected error.
func AssertMarshalError(t *testing.T, name string, err error, wantErr bool) {
	if (err != nil) != wantErr {
		t.Errorf("%s() error = %v, wantErr %v", name, err, wantErr)
	}
}
