package output

import (
	"os"

	"golang.org/x/term"

	"github.com/larsartmann/go-output/enum"
)

// ANSI escape codes for terminal coloring.
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiCyan    = "\033[36m"
	ansiBlue    = "\033[34m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

// ColorMode controls terminal color output.
type ColorMode string

// Terminal color output modes.
const (
	ColorModeAuto   ColorMode = "auto"
	ColorModeAlways ColorMode = "always"
	ColorModeNever  ColorMode = "never"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var colorModeValues = []ColorMode{
	ColorModeAuto,
	ColorModeAlways,
	ColorModeNever,
}

// InvalidColorModeError is returned when an invalid color mode is provided.
type InvalidColorModeError struct {
	Value string
}

// Error returns a descriptive error message for the invalid color mode.
func (e *InvalidColorModeError) Error() string {
	return "invalid color mode: " + e.Value
}

// ParseColorMode converts a string to ColorMode, returning an error if invalid.
func ParseColorMode(s string) (ColorMode, error) {
	v, err := enum.Parse(colorModeValues, s, func(c ColorMode) string { return string(c) })
	if err != nil {
		return "", &InvalidColorModeError{Value: s}
	}

	return v, nil
}

// String returns the string representation of the color mode.
func (c ColorMode) String() string {
	return string(c)
}

// AllowedValues returns all valid color mode values.
func (c ColorMode) AllowedValues() []string {
	return enum.AllowedValues(colorModeValues)
}

// IsValid checks if the color mode is valid.
func (c ColorMode) IsValid() bool {
	return enum.Contains(colorModeValues, c)
}

// ShouldColor returns true if colors should be enabled.
func (c ColorMode) ShouldColor() bool {
	switch c {
	case ColorModeAlways:
		return true
	case ColorModeNever:
		return false
	case ColorModeAuto:
		return isStdoutTerminal() && !isNoColor() && !isCI()
	default:
		return false
	}
}

func isStdoutTerminal() bool {
	if isTerminalByEnv("GO_OUTPUT_FORCE_COLOR", "FORCE_COLOR") {
		return true
	}

	//nolint:gosec // File descriptors are always small positive integers.
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func isTerminalByEnv(envVars ...string) bool {
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" && val != "0" {
			return true
		}
	}

	return false
}

// isNoColor reports whether color output should be suppressed via NO_COLOR or TERM=dumb.
// Mirrors nom.detectNoColor — keep in sync.
func isNoColor() bool {
	return os.Getenv("NO_COLOR") != "" ||
		os.Getenv("TERM") == "dumb"
}

func isCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" ||
		os.Getenv("JENKINS_URL") != "" ||
		os.Getenv("BUILDKITE") != ""
}
