package nom

import (
	"bytes"
	"os"
	"testing"

	"golang.org/x/term"
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

// colorDetectionVars are every environment variable detectNoColor consults.
// Mirrors root color.go — see integration color-agreement test for the contract
// that both detectors must agree on these.
var colorDetectionVars = []string{
	"NO_COLOR", "TERM", "CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "BUILDKITE",
}

// clearColorDetectionEnv clears every env var detectNoColor consults so a test
// starts from a known clean state. Each is restored automatically by t.Setenv.
func clearColorDetectionEnv(t *testing.T) {
	t.Helper()
	for _, env := range colorDetectionVars {
		t.Setenv(env, "")
	}
}

// TestDetectNoColor covers the environment-driven suppression logic of
// detectNoColor. It has zero coverage otherwise, yet it is the function that
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

			if !detectNoColor() {
				t.Errorf("detectNoColor() with %s=%q = false, want true (color suppressed)", tc.env, tc.val)
			}
		})
	}

	t.Run("all_clear_matches_terminal", func(t *testing.T) {
		clearColorDetectionEnv(t)

		//nolint:gosec // File descriptors are always small positive integers.
		want := !term.IsTerminal(int(os.Stdout.Fd()))
		if got := detectNoColor(); got != want {
			t.Errorf("detectNoColor() with no suppressors = %v, want %v (=!term.IsTerminal(stdout))", got, want)
		}
	})
}
