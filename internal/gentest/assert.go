package gentest

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output/testhelpers"
)

type ExpectedOutput = testhelpers.ExpectedOutput

var AssertOutputContains = testhelpers.AssertOutputContains

var AssertMarshalError = testhelpers.AssertMarshalError

type HTMLEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() (string, error)
}

func AssertHTMLEscape(t *testing.T, newRenderer func() HTMLEscapeTestRenderer, name string) {
	t.Helper()

	r := newRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"<script>alert('xss')</script>"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !assertNoRawScript(t, got, name) {
		return
	}

	assertEscapedScript(t, got, name)
}

func assertNoRawScript(t *testing.T, got, name string) bool {
	t.Helper()

	if strings.Contains(got, "<script>") {
		t.Errorf("%s: Render() should escape script tags", name)

		return false
	}

	return true
}

func assertEscapedScript(t *testing.T, got, name string) {
	t.Helper()

	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("%s: Render() should contain escaped script tag", name)
	}
}

func AssertValidJSON(t *testing.T, output string) {
	t.Helper()

	var result map[string]any
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid JSON: %v, got %q", err, output)
	}
}

func AssertValidYAML(t *testing.T, output string) {
	t.Helper()

	var result map[string]any
	if err := yaml.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("output should be valid YAML: %v, got %q", err, output)
	}
}
