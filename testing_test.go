package output

import (
	"io"
	"strings"
	"testing"
)

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

		if len(got) != len(want) {
			t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))

			return
		}

		for i, v := range got {
			if v != want[i] {
				t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
			}
		}
	})
}

type boolMethodTestCase[T any] struct {
	Value T
	Want  bool
}

func testBoolMethod[T any](
	t *testing.T,
	typeName string,
	methodName string,
	testCases []boolMethodTestCase[T],
	boolFunc func(T) bool,
	nameFunc func(T) string,
) {
	t.Helper()

	for _, tc := range testCases {
		name := nameFunc(tc.Value)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := boolFunc(tc.Value); got != tc.Want {
				t.Errorf("%s(%q).%s() = %v, want %v", typeName, name, methodName, got, tc.Want)
			}
		})
	}
}

func testBoolValue(t *testing.T, valueName, methodName string, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("%s.%s() = %v, want %v", valueName, methodName, got, want)
	}
}

type TableWriter interface {
	WriteHeader(cols []string) error
	WriteRow(values []string) error
	Flush()
}

func benchmarkTableWriter(b *testing.B, headers []string, rows [][]string, newWriter func(io.Writer) TableWriter) {
	var buf strings.Builder

	b.ResetTimer()

	for b.Loop() {
		buf.Reset()
		w := newWriter(&buf)

		_ = w.WriteHeader(headers)
		for _, row := range rows {
			_ = w.WriteRow(row)
		}

		w.Flush()
	}
}
