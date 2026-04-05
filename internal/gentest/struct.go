package gentest

import "testing"

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
