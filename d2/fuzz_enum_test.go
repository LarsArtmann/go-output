package d2

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func fuzzTestParseEnum[E testhelpers.StringEnum](f *testing.F, values []E, parse func(string) (E, error)) {
	for _, v := range values {
		f.Add(string(v))
	}

	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		got, err := parse(input)
		if err != nil {
			if got != "" {
				t.Errorf("parse(%q) returned non-empty on error: %q", input, got)
			}

			return
		}

		if !got.IsValid() {
			t.Errorf("parse(%q) = %q, not valid", input, got)
		}
	})
}
