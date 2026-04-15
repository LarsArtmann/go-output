package gentest

import (
	"strings"
	"testing"
)

// ExpectedOutput contains a substring to check and its corresponding error message.
type ExpectedOutput struct {
	Substring string
	Message   string
}

// AssertContains checks that output contains substr, failing with msg if not.
func AssertContains(t *testing.T, output, substr, msg string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Error(msg)
	}
}

// HTMLEscapeTestRenderer is an interface for HTML renderers that support escaping tests.
type HTMLEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() string
}

// AssertHTMLEscape verifies that a renderer properly escapes HTML content.
func AssertHTMLEscape(t *testing.T, newRenderer func() HTMLEscapeTestRenderer, name string) {
	t.Helper()

	r := newRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"<script>alert('xss')</script>"})

	got := r.Render()

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

// AssertStringSliceEqual checks that got and want are equal, failing with descriptive error.
func AssertStringSliceEqual(t *testing.T, name string, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("%s returned %d values, want %d", name, len(got), len(want))

		return
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("%s[%d] = %v, want %v", name, i, v, want[i])
		}
	}
}

// AssertEqual checks that got equals want, failing with descriptive error.
func AssertEqual[T comparable](t *testing.T, name string, input any, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("%s(%v) = %v, want %v", name, input, got, want)
	}
}
