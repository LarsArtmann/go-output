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

	AssertMarshalError(t, "no error", nil, false)
	AssertMarshalError(t, "has error", assertErr, true)

	t.Run("unexpected error", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertMarshalError(mock, "unexpected", assertErr, false)

		if !mock.Failed() {
			t.Error("expected failure for unexpected error")
		}
	})

	t.Run("expected error missing", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertMarshalError(mock, "missing", nil, true)

		if !mock.Failed() {
			t.Error("expected failure for missing error")
		}
	})
}

func TestTestAllowedValues(t *testing.T) {
	t.Parallel()

	TestAllowedValues(t, "formats", []string{"json", "yaml"}, []string{"json", "yaml"})
}

var assertErr = errors.New("test error")

func parseLower(s string) (string, error) {
	if s == "" {
		return "", assertErr
	}

	return strings.ToLower(s), nil
}

func TestTestParseEnum(t *testing.T) {
	t.Parallel()

	TestParseEnum(t, "lower",
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

	TestEnumString(t, "upper",
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
