package nom

import (
	"bytes"
	"os"
	"testing"
)

func TestVisibleWidth(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  int
	}{
		{"", 0},
		{"hello", 5},
		{"\x1b[31mhello\x1b[0m", 5},
		{"⏵ run", 5},
		{"…", 1},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			if got := VisibleWidth(tc.input); got != tc.want {
				t.Errorf("VisibleWidth(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestVisibleLineCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		line      string
		termWidth int
		want      int
	}{
		{"", 80, 0},
		{"hello", 80, 1},
		{"hello world", 5, 3},
		{"\x1b[31mhello\x1b[0m", 3, 2},
	}

	for _, tc := range cases {
		t.Run(tc.line, func(t *testing.T) {
			t.Parallel()

			if got := VisibleLineCount(tc.line, tc.termWidth); got != tc.want {
				t.Errorf("VisibleLineCount(%q, %d) = %d, want %d", tc.line, tc.termWidth, got, tc.want)
			}
		})
	}
}

func TestTruncateVisible(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input     string
		maxWidth  int
		wantWidth int
	}{
		{"hello world", 5, 5},
		{"hello", 10, 5},
		{"\x1b[31mhello world\x1b[0m", 5, 5},
		{"", 5, 0},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got := TruncateVisible(tc.input, tc.maxWidth)
			if w := VisibleWidth(got); w != tc.wantWidth {
				t.Errorf(
					"TruncateVisible(%q, %d) = %q (width %d), want width %d",
					tc.input,
					tc.maxWidth,
					got,
					w,
					tc.wantWidth,
				)
			}
		})
	}
}

func TestGetTerminalWidth_DefaultForBuffer(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if got := GetTerminalWidth(&buf); got != defaultTerminalWidth {
		t.Errorf("GetTerminalWidth(bytes.Buffer) = %d, want %d", got, defaultTerminalWidth)
	}
}

func TestGetTerminalWidth_Stderr(t *testing.T) {
	t.Parallel()

	// Stderr may or may not be a terminal; we just verify it doesn't panic
	// and returns a positive width.
	got := GetTerminalWidth(os.Stderr)
	if got <= 0 {
		t.Errorf("GetTerminalWidth(os.Stderr) = %d, want > 0", got)
	}
}

func TestPhysicalLineCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		frame     string
		termWidth int
		want      int
	}{
		{"a\nb", 80, 2},
		{"hello world", 5, 3},
		{"hello\n\nworld", 80, 3},
	}

	for _, tc := range cases {
		t.Run(tc.frame, func(t *testing.T) {
			t.Parallel()

			if got := PhysicalLineCount(tc.frame, tc.termWidth); got != tc.want {
				t.Errorf("PhysicalLineCount(%q, %d) = %d, want %d", tc.frame, tc.termWidth, got, tc.want)
			}
		})
	}
}
