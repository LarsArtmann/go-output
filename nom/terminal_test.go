package nom

import (
	"bytes"
	"os"
	"testing"

	"golang.org/x/term"

	output "github.com/larsartmann/go-output"
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

//nolint:gosmopolitan // intentional wide-rune width test data (CJK + emoji)
func TestVisibleWidth_WideRunes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  int
	}{
		{"中文", 4}, // CJK characters are double-width (2 each)
		{"中a", 3}, // mixed CJK + ASCII
		{"🎉", 2},  // emoji is double-width
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

//nolint:gosmopolitan // intentional wide-rune line-wrap test data
func TestVisibleLineCount_WideRunes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		line      string
		termWidth int
		want      int
	}{
		{"中文", 2, 2},  // width-4 CJK over termWidth 2 -> 2 lines
		{"中文", 4, 1},  // width-4 CJK fits exactly in termWidth 4 -> 1 line
		{"中中中中", 3, 4}, // width-8 CJK over termWidth 3 -> straddle: each line holds 1 wide char (2 cells) + 1 gap = 4 lines
		{"中中中中", 5, 2}, // width-8 CJK over termWidth 5 -> line1: 2 chars (4 cells), line2: 2 chars (4 cells) = 2 lines
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

func TestStripANSI(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"plain", "hello", "hello"},
		{"empty", "", ""},
		{"color", "\x1b[31mhello\x1b[0m", "hello"},
		{"cursor_reset", "\x1b[2K\x1b[1;1Hline", "line"},
		{"nested", "\x1b[1m\x1b[31mbold red\x1b[0m", "bold red"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := StripANSI(tc.input); got != tc.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
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

// colorDetectionVars are every environment variable detectNoColor consults.
// Backed by envdetect.CIEnvVars so root and nom cannot drift apart.
var colorDetectionVars = append([]string{"NO_COLOR", "TERM"}, output.CIEnvVars...)

// clearColorDetectionEnv clears every env var detectNoColor consults so a test
// starts from a known clean state. Each is restored automatically by t.Setenv.
func clearColorDetectionEnv(t *testing.T) {
	t.Helper()

	for _, env := range colorDetectionVars {
		t.Setenv(env, "")
	}
}

// TestDetectNoColor covers the environment-driven suppression logic of
// detectNoColorForWriter. It has zero coverage otherwise, yet it is the function that
// must stay aligned with root output.isNoColor()+isCI() (M2 split-brain fix).
func TestDetectNoColor(t *testing.T) {
	suppressors := []struct {
		name string
		env  string
		val  string
	}{
		{"NO_COLOR", "NO_COLOR", "1"},
		{"TERM_dumb", "TERM", "dumb"},
		{"CI", "CI", "true"},
		{"GITHUB_ACTIONS", "GITHUB_ACTIONS", "true"},
		{"GITLAB_CI", "GITLAB_CI", "true"},
		{"JENKINS_URL", "JENKINS_URL", "https://jenkins.example"},
		{"BUILDKITE", "BUILDKITE", "true"},
	}

	for _, tc := range suppressors {
		t.Run(tc.name, func(t *testing.T) {
			clearColorDetectionEnv(t)
			t.Setenv(tc.env, tc.val)

			if !detectNoColorForWriter(os.Stdout) {
				t.Errorf("detectNoColorForWriter() with %s=%q = false, want true (color suppressed)", tc.env, tc.val)
			}
		})
	}

	t.Run("all_clear_matches_terminal", func(t *testing.T) {
		clearColorDetectionEnv(t)

		//nolint:gosec // File descriptors are always small positive integers.
		stdoutIsTTY := term.IsTerminal(int(os.Stdout.Fd()))

		want := !stdoutIsTTY // non-TTY → no color
		if got := detectNoColorForWriter(os.Stdout); got != want {
			t.Errorf(
				"detectNoColorForWriter(stdout) with no suppressors = %v, want %v (=!term.IsTerminal(stdout))",
				got,
				want,
			)
		}
	})
}
