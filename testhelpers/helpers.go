package testhelpers

import (
	"strings"
	"testing"
)

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

// AssertContains checks that output contains substr, failing with msg if not.
func AssertContains(t *testing.T, output, substr, msg string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Error(msg)
	}
}

// AssertEqual checks that got equals want, failing with descriptive error.
func AssertEqual[T comparable](t *testing.T, name string, input any, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("%s(%v) = %v, want %v", name, input, got, want)
	}
}

// StringEnum is a constraint for string-based enum types with IsValid().
type StringEnum interface {
	~string
	IsValid() bool
}

// TestEnumIsValid runs table-driven tests for IsValid() on string-based enums.
func TestEnumIsValid[T StringEnum](t *testing.T, values []T, expected []bool) {
	t.Helper()

	if len(values) != len(expected) {
		t.Fatalf("TestEnumIsValid: len(values)=%d != len(expected)=%d", len(values), len(expected))
	}

	for i, val := range values {
		name := string(val)
		if name == "" {
			name = "empty"
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := val.IsValid(); got != expected[i] {
				t.Errorf("%v.IsValid() = %v, want %v", val, got, expected[i])
			}
		})
	}
}

// FieldCheck represents a single field validation function.
type FieldCheck func(t *testing.T)

// TestStructFields runs field validation checks for a struct.
func TestStructFields(t *testing.T, checks ...FieldCheck) {
	t.Helper()

	for _, check := range checks {
		check(t)
	}
}

func equalField[V comparable](name string, got, want V) FieldCheck {
	return func(t *testing.T) {
		if got != want {
			t.Errorf("%s = %v, want %v", name, got, want)
		}
	}
}

// StringField creates a FieldCheck for a string field.
func StringField(name, got, want string) FieldCheck {
	return equalField(name, got, want)
}

// IntField creates a FieldCheck for an int field.
func IntField(name string, got, want int) FieldCheck {
	return equalField(name, got, want)
}
