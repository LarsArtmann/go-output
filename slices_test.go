package output

import "testing"

func TestFilledStrings(t *testing.T) {
	t.Parallel()

	t.Run("creates slice of correct size", func(t *testing.T) {
		t.Parallel()

		got := FilledStrings(3, "x")
		if len(got) != 3 {
			t.Errorf("FilledStrings(3, x) returned %d elements, want 3", len(got))
		}

		for i, v := range got {
			if v != "x" {
				t.Errorf("got[%d] = %q, want %q", i, v, "x")
			}
		}
	})

	t.Run("zero length", func(t *testing.T) {
		t.Parallel()

		got := FilledStrings(0, "x")
		if len(got) != 0 {
			t.Errorf("FilledStrings(0, x) returned %d elements, want 0", len(got))
		}
	})

	t.Run("empty value", func(t *testing.T) {
		t.Parallel()

		got := FilledStrings(2, "")
		if len(got) != 2 {
			t.Errorf("FilledStrings(2, empty) returned %d elements, want 2", len(got))
		}
	})
}
