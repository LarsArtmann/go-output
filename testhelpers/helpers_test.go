package testhelpers

import (
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
}

func TestAssertContains(t *testing.T) {
	t.Parallel()

	AssertContains(t, "hello world", "world", "should contain world")
	AssertContains(t, "hello", "hello", "should contain hello")
}

func TestAssertEqual(t *testing.T) {
	t.Parallel()

	AssertEqual(t, "int", "input", 42, 42)
	AssertEqual(t, "string", "input", "same", "same")
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
}

func TestIntField(t *testing.T) {
	t.Parallel()

	IntField("count", 5, 5)(t)
}
