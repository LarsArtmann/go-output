package testhelpers

import (
	"errors"
	"strings"
	"testing"
)

type isValidEnum string

func (e isValidEnum) IsValid() bool {
	return e == "valid"
}

func TestAssertStringSliceEqual(t *testing.T) {
	t.Parallel()

	AssertStringSliceEqual(t, "equal", []string{"a", "b"}, []string{"a", "b"})
	AssertStringSliceEqual(t, "empty", []string{}, []string{})

	t.Run("different lengths", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertStringSliceEqual(mock, "len", []string{"a"}, []string{"a", "b"})

		if !mock.Failed() {
			t.Error("expected failure for different lengths")
		}
	})

	t.Run("different values", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertStringSliceEqual(mock, "val", []string{"a"}, []string{"b"})

		if !mock.Failed() {
			t.Error("expected failure for different values")
		}
	})
}

func TestAssertContains(t *testing.T) {
	t.Parallel()

	AssertContains(t, "hello world", "world", "should contain world")
	AssertContains(t, "hello", "hello", "should contain hello")

	t.Run("missing substring", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertContains(mock, "hello", "missing", "should fail")

		if !mock.Failed() {
			t.Error("expected failure for missing substring")
		}
	})
}

func TestAssertEqual(t *testing.T) {
	t.Parallel()

	AssertEqual(t, "int", "input", 42, 42)
	AssertEqual(t, "string", "input", "same", "same")

	t.Run("unequal values", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertEqual(mock, "int", "input", 1, 2)

		if !mock.Failed() {
			t.Error("expected failure for unequal values")
		}
	})
}

func TestTestEnumIsValid(t *testing.T) {
	t.Parallel()

	TestEnumIsValid(t, []isValidEnum{"valid", "invalid"}, []bool{true, false})
}

func TestTestEnumIsValidEmptyString(t *testing.T) {
	t.Parallel()

	TestEnumIsValid(t, []isValidEnum{""}, []bool{false})
}

func TestTestStructFields(t *testing.T) {
	t.Parallel()

	TestStructFields(
		t,
		StringField("name", "same", "same"),
		IntField("count", 5, 5),
	)
}

func TestStringField(t *testing.T) {
	t.Parallel()

	StringField("name", "same", "same")(t)

	t.Run("mismatch", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		StringField("name", "got", "want")(mock)

		if !mock.Failed() {
			t.Error("expected failure for mismatched string field")
		}
	})
}

func TestAssertOutputContains(t *testing.T) {
	t.Parallel()

	AssertOutputContains(t, "hello world", "world")

	t.Run("missing substring", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertOutputContains(mock, "hello", "missing")

		if !mock.Failed() {
			t.Error("expected failure for missing substring")
		}
	})
}

func TestAssertMarshalError(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name       string
		label      string
		err        error
		wantError  bool
		shouldFail bool
	}{
		{name: "no error", label: "no error", err: nil, wantError: false, shouldFail: false},
		{name: "has error", label: "has error", err: errTestAssert, wantError: true, shouldFail: false},
		{name: "unexpected error", label: "unexpected", err: errTestAssert, wantError: false, shouldFail: true},
		{name: "expected error missing", label: "missing", err: nil, wantError: true, shouldFail: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.shouldFail {
				mock := &testing.T{}
				AssertMarshalError(mock, tt.label, tt.err, tt.wantError)

				if !mock.Failed() {
					t.Errorf("expected failure for %s", tt.name)
				}

				return
			}

			AssertMarshalError(t, tt.label, tt.err, tt.wantError)
		})
	}
}

func TestTestAllowedValues(t *testing.T) {
	t.Parallel()

	TestAllowedValues(t, "formats", []string{"json", "yaml"}, []string{"json", "yaml"})
}

var errTestAssert = errors.New("test error")

func parseLower(s string) (string, error) {
	if s == "" {
		return "", errTestAssert
	}

	return strings.ToLower(s), nil
}

func TestTestParseEnum(t *testing.T) {
	t.Parallel()

	TestParseEnum(
		t, "lower",
		parseLower,
		[]ParseEnumTestCase[string]{
			{Name: "upper", Input: "FOO", Want: "foo", WantErr: false},
			{Name: "error", Input: "", Want: "", WantErr: true},
		},
		func(a, b string) bool { return a == b },
	)
}

func TestTestEnumString(t *testing.T) {
	t.Parallel()

	TestEnumString(
		t, "upper",
		[]StringEnumTestCase[string]{
			{Value: "hello", Want: "HELLO"},
			{Value: "world", Want: "WORLD"},
		},
		strings.ToUpper,
	)
}

func TestIntField(t *testing.T) {
	t.Parallel()

	IntField("count", 5, 5)(t)

	t.Run("mismatch", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		IntField("count", 1, 2)(mock)

		if !mock.Failed() {
			t.Error("expected failure for mismatched int field")
		}
	})
}

func TestAssertLineCount(t *testing.T) {
	t.Parallel()

	t.Run("correct line count", func(t *testing.T) {
		t.Parallel()

		AssertLineCount(t, "three lines", "a\nb\nc", 3)
	})

	t.Run("wrong line count", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertLineCount(mock, "wrong", "a\nb", 3)

		if !mock.Failed() {
			t.Error("expected failure for wrong line count")
		}
	})

	t.Run("trims whitespace", func(t *testing.T) {
		t.Parallel()

		AssertLineCount(t, "trailing newline", "a\nb\n", 2)
	})
}

func TestAssertLastLineContains(t *testing.T) {
	t.Parallel()

	t.Run("last line contains substring", func(t *testing.T) {
		t.Parallel()

		AssertLastLineContains(t, "test", "first\nlast line", "last")
	})

	t.Run("last line missing substring", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertLastLineContains(mock, "test", "first\nlast", "missing")

		if !mock.Failed() {
			t.Error("expected failure for missing substring in last line")
		}
	})
}

func TestAssertErrorContains(t *testing.T) {
	t.Parallel()

	t.Run("error contains substring", func(t *testing.T) {
		t.Parallel()

		AssertErrorContains(t, errors.New("file not found"), "not found", "should contain")
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertErrorContains(mock, nil, "anything", "nil error")

		if !mock.Failed() {
			t.Error("expected failure for nil error")
		}
	})

	t.Run("error missing substring", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertErrorContains(mock, errors.New("other"), "missing", "wrong error")

		if !mock.Failed() {
			t.Error("expected failure for missing substring in error")
		}
	})
}
