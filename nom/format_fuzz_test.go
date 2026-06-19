package nom

import (
	"strings"
	"testing"
	"time"
)

// FuzzFormatDuration guards FormatDuration against panics and malformed output
// across the full int64 nanosecond range (negative, zero, sub-second, and
// multi-hour values). It enforces two invariants for non-negative inputs:
//   - the result is never empty, and
//   - it never renders a stray minus sign.
func FuzzFormatDuration(f *testing.F) {
	f.Add(int64(0))
	f.Add(int64(time.Millisecond))
	f.Add(int64(999 * time.Millisecond))
	f.Add(int64(time.Second))
	f.Add(int64(59*time.Second + 950*time.Millisecond)) // boundary that once rounded to "60.0s"
	f.Add(int64(time.Minute))
	f.Add(int64(time.Minute + 30*time.Second))
	f.Add(int64(time.Hour))
	f.Add(int64(time.Hour + 30*time.Minute))
	f.Add(int64(24 * time.Hour))
	f.Add(int64(-1)) // negative durations must not panic

	f.Fuzz(func(t *testing.T, nanos int64) {
		d := time.Duration(nanos)
		got := FormatDuration(d)

		if d >= 0 {
			if got == "" {
				t.Errorf("FormatDuration(%v) returned empty string", d)
			}

			if strings.Contains(got, "-") {
				t.Errorf("FormatDuration(%v) = %q contains '-' for non-negative input", d, got)
			}
		}
	})
}
