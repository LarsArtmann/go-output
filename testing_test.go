package output

import "testing"

// runSubtest runs a subtest with the standard parallel pattern.
func runSubtest(t *testing.T, name string, testFunc func(*testing.T)) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		testFunc(t)
	})
}

type parseEnumTestCase[T any] struct {
	name    string
	input   string
	want    T
	wantErr bool
}

func testParseEnum[T any](
	t *testing.T,
	name string,
	parseFunc func(string) (T, error),
	testCases []parseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				got, err := parseFunc(testCase.input)

				if (err != nil) != testCase.wantErr {
					t.Errorf("%s() error = %v, wantErr %v", name, err, testCase.wantErr)

					return
				}

				if !equalFunc(got, testCase.want) {
					t.Errorf("%s() = %v, want %v", name, got, testCase.want)
				}
			})
		}
	})
}

type stringEnumTestCase[T any] struct {
	value T
	want  string
}

func testEnumString[T any](
	t *testing.T,
	name string,
	testCases []stringEnumTestCase[T],
	stringFunc func(T) string,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, testCase := range testCases {
			t.Run(testCase.want, func(t *testing.T) {
				t.Parallel()

				if got := stringFunc(testCase.value); got != testCase.want {
					t.Errorf("%s() = %v, want %v", name, got, testCase.want)
				}
			})
		}
	})
}

func testAllowedValues(
	t *testing.T,
	name string,
	got []string,
	want []string,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		assertStringSliceEqual(t, name, got, want)
	})
}
