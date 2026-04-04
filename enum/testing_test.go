package enum

import "testing"

type parseEnumTestCase[T comparable] struct {
	name    string
	input   string
	want    T
	wantErr bool
}

func testParseEnum[T comparable](
	t *testing.T,
	name string,
	values []T,
	toString func(T) string,
	testCases []parseEnumTestCase[T],
	equalFunc func(T, T) bool,
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got, err := Parse(values, tc.input, toString)

				if (err != nil) != tc.wantErr {
					t.Errorf("%s() error = %v, wantErr %v", name, err, tc.wantErr)

					return
				}

				if !equalFunc(got, tc.want) {
					t.Errorf("%s() = %v, want %v", name, got, tc.want)
				}
			})
		}
	})
}
