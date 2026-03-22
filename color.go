// Package output provides consistent output formatting for CLI applications.
package output

import (
	"fmt"
	"os"
	"slices"
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

// ParseColorMode parses a color mode string.
func ParseColorMode(s string) (ColorMode, error) {
	if slices.Contains(colorModeValues, ColorMode(s)) {
		return ColorMode(s), nil
	}
	return "", fmt.Errorf("invalid color mode: %q (allowed: %v)", s, colorModeValues)
}

func (c ColorMode) String() string {
	return string(c)
}

// AllowedValues returns all valid color mode values.
func (c ColorMode) AllowedValues() []string {
	values := make([]string, len(colorModeValues))
	for i, v := range colorModeValues {
		values[i] = string(v)
	}
	return values
}

// IsValid checks if the color mode is valid.
func (c ColorMode) IsValid() bool {
	return slices.Contains(colorModeValues, c)
}

// ShouldColor returns true if colors should be enabled.
func (c ColorMode) ShouldColor() bool {
	switch c {
	case ColorModeAlways:
		return true
	case ColorModeNever:
		return false
	case ColorModeAuto:
		return isTerminal() && !isNoColor() && !isCI()
	default:
		return false
	}
}

// ToANSI returns the ANSI escape sequence prefix if colors are enabled.
func (c ColorMode) ToANSI() string {
	if !c.ShouldColor() {
		return ""
	}
	return "\033["
}

func isTerminal() bool {
	return isStdoutTerminal() || isStderrTerminal()
}

func isStdoutTerminal() bool {
	return isTerminalByEnv("GO_OUTPUT_FORCE_COLOR", "FORCE_COLOR")
}

func isStderrTerminal() bool {
	return false
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
