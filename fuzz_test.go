package output

import (
	"testing"
)

func FuzzParseColorMode(f *testing.F) {
	f.Add("auto")
	f.Add("always")
	f.Add("never")
	f.Add("AUTO")
	f.Add("")
	f.Add("invalid")

	f.Fuzz(func(t *testing.T, input string) {
		mode, err := ParseColorMode(input)
		if err != nil {
			return
		}

		if !mode.IsValid() {
			t.Errorf("ParseColorMode(%q) returned invalid mode %q", input, mode)
		}

		if mode.String() != string(mode) {
			t.Errorf("ColorMode.String() = %q, want %q", mode.String(), string(mode))
		}
	})
}
