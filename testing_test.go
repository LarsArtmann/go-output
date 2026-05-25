package output

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func runSubtest(t *testing.T, name string, testFunc func(*testing.T)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		testFunc(t)
	})
}

type parseEnumTestCase[T any] = testhelpers.ParseEnumTestCase[T]

func testParseEnum[T any](
	t *testing.T,
	name string,
	parseFunc func(string) (T, error),
	testCases []parseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	testhelpers.TestParseEnum(t, name, parseFunc, testCases, equalFunc)
}

type stringEnumTestCase[T any] = testhelpers.StringEnumTestCase[T]

func testEnumString[T any](
	t *testing.T,
	name string,
	testCases []stringEnumTestCase[T],
	stringFunc func(T) string,
) {
	testhelpers.TestEnumString(t, name, testCases, stringFunc)
}

var testAllowedValues = testhelpers.TestAllowedValues
