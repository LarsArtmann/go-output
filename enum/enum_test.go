package enum

import "testing"

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

	testAssertStringSliceEqual(t, "AllowedStrings", got, want)
}

func testAssertStringSliceEqual(t *testing.T, name string, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("%s() returned %d values, want %d", name, len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("%s()[%d] = %v, want %v", name, i, v, want[i])
		}
	}
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
