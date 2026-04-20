package enum

import (
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
)

type testEnum string

const (
	testEnumA testEnum = "a"
	testEnumB testEnum = "b"
	testEnumC testEnum = "c"
)

//nolint:gochecknoglobals // Test data
var testEnumValues = []testEnum{testEnumA, testEnumB, testEnumC}

func testEnumString(v testEnum) string {
	return string(v)
}

func (e testEnum) String() string {
	return string(e)
}

func TestParse(t *testing.T) {
	tests := []parseEnumTestCase[testEnum]{
		{"valid_a", "a", testEnumA, false},
		{"valid_b", "b", testEnumB, false},
		{"valid_c", "c", testEnumC, false},
		{"invalid", "x", "", true},
		{"empty", "", "", true},
	}
	testParseEnum(
		t,
		"Parse",
		testEnumValues,
		testEnumString,
		tests,
		func(a, b testEnum) bool { return a == b },
	)
}

func TestContains(t *testing.T) {
	t.Parallel()

	if !Contains(testEnumValues, testEnumA) {
		t.Error("Contains() should return true for valid value")
	}

	if Contains(testEnumValues, testEnum("invalid")) {
		t.Error("Contains() should return false for invalid value")
	}
}

func TestAllowedStrings(t *testing.T) {
	t.Parallel()

	got := AllowedStrings(testEnumValues, testEnumString)
	want := []string{"a", "b", "c"}

	gentest.AssertStringSliceEqual(t, "AllowedStrings", got, want)
}

func TestAllowedValues(t *testing.T) {
	t.Parallel()

	got := AllowedValues(testEnumValues)
	want := []string{"a", "b", "c"}

	gentest.AssertStringSliceEqual(t, "AllowedValues", got, want)
}

func TestParseError(t *testing.T) {
	t.Parallel()

	err := &ParseError[testEnum]{Value: "invalid", Values: testEnumValues}
	got := err.Error()

	want := `invalid value: "invalid"`
	if got != want {
		t.Errorf("ParseError.Error() = %q, want %q", got, want)
	}
}
