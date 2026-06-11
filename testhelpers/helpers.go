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

// AssertMarshalError checks that a marshal function returns the expected error.
func AssertMarshalError(t *testing.T, name string, err error, wantErr bool) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("%s() error = %v, wantErr %v", name, err, wantErr)
	}
}

// AssertLineCount checks that output, when split on \n and trimmed, has exactly want lines.
func AssertLineCount(t *testing.T, name, output string, want int) {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != want {
		t.Errorf("%s: got %d lines, want %d\noutput: %q", name, len(lines), want, output)
	}
}

// AssertLastLineContains checks that the last non-empty line of output contains substr.
func AssertLastLineContains(t *testing.T, name, output, substr string) {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		t.Errorf("%s: output is empty", name)

		return
	}

	if !strings.Contains(lines[len(lines)-1], substr) {
		t.Errorf("%s: last line %q should contain %q", name, lines[len(lines)-1], substr)
	}
}

// AssertErrorContains checks that err's message contains substr, failing with msg.
func AssertErrorContains(t *testing.T, err error, substr, msg string) {
	t.Helper()

	if err == nil {
		t.Errorf("%s: expected error, got nil", msg)

		return
	}

	if !strings.Contains(err.Error(), substr) {
		t.Errorf("%s: error %q should contain %q", msg, err.Error(), substr)
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

// TestAllowedValues runs a subtest checking that got matches want string slice.
func TestAllowedValues(t *testing.T, name string, got, want []string) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		AssertStringSliceEqual(t, name, got, want)
	})
}

// IntField creates a FieldCheck for an int field.
func IntField(name string, got, want int) FieldCheck {
	return equalField(name, got, want)
}

// ParseEnumTestCase is a test case for enum Parse functions.
type ParseEnumTestCase[T any] struct {
	Name    string
	Input   string
	Want    T
	WantErr bool
}

// TestParseEnum runs table-driven tests for a generic enum Parse function.
func TestParseEnum[T any](
	t *testing.T,
	name string,
	parseFunc func(string) (T, error),
	testCases []ParseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, tc := range testCases {
			t.Run(tc.Name, func(t *testing.T) {
				t.Parallel()

				got, err := parseFunc(tc.Input)

				if (err != nil) != tc.WantErr {
					t.Errorf("%s() error = %v, wantErr %v", name, err, tc.WantErr)

					return
				}

				if !equalFunc(got, tc.Want) {
					t.Errorf("%s() = %v, want %v", name, got, tc.Want)
				}
			})
		}
	})
}

// StringEnumTestCase is a test case for enum String functions.
type StringEnumTestCase[T any] struct {
	Value T
	Want  string
}

// TestEnumString runs table-driven tests for a generic enum String function.
func TestEnumString[T any](
	t *testing.T,
	name string,
	testCases []StringEnumTestCase[T],
	stringFunc func(T) string,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, tc := range testCases {
			t.Run(tc.Want, func(t *testing.T) {
				t.Parallel()

				if got := stringFunc(tc.Value); got != tc.Want {
					t.Errorf("%s() = %v, want %v", name, got, tc.Want)
				}
			})
		}
	})
}
