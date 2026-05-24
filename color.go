package output

import (
	"os"

	"golang.org/x/term"

	"github.com/larsartmann/go-output/enum"
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

func isNoColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

func isCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" ||
		os.Getenv("JENKINS_URL") != "" ||
		os.Getenv("BUILDKITE") != ""
}
