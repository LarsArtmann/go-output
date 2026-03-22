package output

import (
	"fmt"
	"os"
)

type ColorMode string

const (
	ColorModeAuto   ColorMode = "auto"
	ColorModeAlways ColorMode = "always"
	ColorModeNever  ColorMode = "never"
)

var colorModeValues = []ColorMode{
	ColorModeAuto,
	ColorModeAlways,
	ColorModeNever,
}

func ParseColorMode(s string) (ColorMode, error) {
	for _, v := range colorModeValues {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid color mode: %q (allowed: %v)", s, colorModeValues)
}

func (c ColorMode) String() string {
	return string(c)
}

func (c ColorMode) AllowedValues() []string {
	values := make([]string, len(colorModeValues))
	for i, v := range colorModeValues {
		values[i] = string(v)
	}
	return values
}

func (c ColorMode) IsValid() bool {
	for _, v := range colorModeValues {
		if c == v {
			return true
		}
	}
	return false
}

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
