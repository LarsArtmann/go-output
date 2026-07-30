package output

import (
	"os"

	"golang.org/x/term"
)

// ColorMode controls terminal color output.
type ColorMode string

// Terminal color output modes.
const (
	ColorModeAuto   ColorMode = "auto"
	ColorModeAlways ColorMode = "always"
	ColorModeNever  ColorMode = "never"
)

// AllColorModes contains all valid color mode values.
//
//nolint:gochecknoglobals // Global variable used for value iteration.
var AllColorModes = []ColorMode{
	ColorModeAuto,
	ColorModeAlways,
	ColorModeNever,
}

// InvalidColorModeError is returned when an invalid color mode is provided.
type InvalidColorModeError struct {
	Value   string
	Allowed []ColorMode
}

// Error returns a descriptive error message for the invalid color mode.
func (e *InvalidColorModeError) Error() string {
	return "invalid color mode: " + e.Value + " (allowed: " + joinStrings(EnumAllowedValues(e.Allowed)) + ")"
}

// ParseColorMode converts a string to ColorMode, returning an error if invalid.
func ParseColorMode(s string) (ColorMode, error) {
	v, err := ParseEnum(AllColorModes, s, func(c ColorMode) string { return string(c) })
	if err != nil {
		return "", &InvalidColorModeError{Value: s, Allowed: AllColorModes}
	}

	return v, nil
}

// String returns the string representation of the color mode.
func (c ColorMode) String() string {
	return string(c)
}

// AllowedValues returns all valid color mode values.
func (c ColorMode) AllowedValues() []string {
	return EnumAllowedValues(AllColorModes)
}

// IsValid checks if the color mode is valid.
func (c ColorMode) IsValid() bool {
	return ContainsEnum(AllColorModes, c)
}

// ShouldColor returns true if colors should be enabled.
func (c ColorMode) ShouldColor() bool {
	switch c {
	case ColorModeAlways:
		return true
	case ColorModeNever:
		return false
	case ColorModeAuto:
		return stdoutIsTerminal() && !noColorEnv() && !ciEnv()
	default:
		return false
	}
}

// ColorConfig is the shared color-mode configuration that every sub-module's
// CQRS Config embeds. It carries the ColorMode field and the canonical
// default value. Embedding it lets a module's Config reuse the same
// color-mode behavior across markdown, tree, table, and any future
// module that needs to honor WithColorMode(...) without re-declaring the
// field or its default.
type ColorConfig struct {
	ColorMode ColorMode
}

// DefaultColorConfig returns a ColorConfig pre-initialized to ColorModeAuto.
// Use it as the starting value in your module's CQRS Config literal.
func DefaultColorConfig() ColorConfig {
	return ColorConfig{ColorMode: ColorModeAuto}
}

// Detection functions are overridable variables for deterministic testing.
// Tests can swap these to control ShouldColor() output without relying on
// the real TTY, env vars, or CI environment.
//
//nolint:gochecknoglobals // Overridable for test determinism (#9).
var (
	stdoutIsTerminal = func() bool {
		if isTerminalByEnv("GO_OUTPUT_FORCE_COLOR", "FORCE_COLOR") {
			return true
		}

		//nolint:gosec // File descriptors are always small positive integers.
		return term.IsTerminal(int(os.Stdout.Fd()))
	}

	noColorEnv = IsNoColor

	ciEnv = IsCI
)

func isTerminalByEnv(envVars ...string) bool {
	for _, env := range envVars {
		if val := os.Getenv(env); val != "" && val != "0" {
			return true
		}
	}

	return false
}
