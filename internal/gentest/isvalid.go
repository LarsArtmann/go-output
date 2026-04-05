package gentest

import (
	"testing"
)

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
