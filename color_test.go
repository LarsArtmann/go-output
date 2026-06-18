package output

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "")

	if noColorEnv() {
		t.Error("noColorEnv() should return false when NO_COLOR is not set")
	}

	t.Setenv("NO_COLOR", "1")

	if !noColorEnv() {
		t.Error("noColorEnv() should return true when NO_COLOR is set")
	}
}

func TestCIEnv(t *testing.T) {
	ciVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "BUILDKITE"}

	for _, v := range ciVars {
		t.Setenv(v, "")
	}

	if ciEnv() {
		t.Error("ciEnv() should return false when no CI env vars are set")
	}

	t.Setenv("CI", "true")

	if !ciEnv() {
		t.Error("ciEnv() should return true when CI is set")
	}

	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "true")

	if !ciEnv() {
		t.Error("ciEnv() should return true when GITHUB_ACTIONS is set")
	}
}

func TestIsTerminalByEnv(t *testing.T) {
	t.Setenv("FORCE_COLOR", "")

	if isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return false when env is not set")
	}

	t.Setenv("FORCE_COLOR", "0")

	if isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return false when env is '0'")
	}

	t.Setenv("FORCE_COLOR", "1")

	if !isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return true when env is '1'")
	}
}

func TestParseColorMode(t *testing.T) {
	tests := []testhelpers.ParseEnumTestCase[ColorMode]{
		{Name: "auto", Input: "auto", Want: ColorModeAuto},
		{Name: "always", Input: "always", Want: ColorModeAlways},
		{Name: "never", Input: "never", Want: ColorModeNever},
		{Name: "invalid", Input: "invalid", WantErr: true},
		{Name: "empty", Input: "", WantErr: true},
	}
	testhelpers.TestParseEnum(
		t,
		"ParseColorMode",
		ParseColorMode,
		tests,
		func(a, b ColorMode) bool { return a == b },
	)
}

func TestColorModeString(t *testing.T) {
	tests := []testhelpers.StringEnumTestCase[ColorMode]{
		{Value: ColorModeAuto, Want: "auto"},
		{Value: ColorModeAlways, Want: "always"},
		{Value: ColorModeNever, Want: "never"},
	}

	testhelpers.TestEnumString(t, "ColorMode.String", tests, func(m ColorMode) string { return m.String() })
}

func TestColorModeAllowedValues(t *testing.T) {
	testhelpers.TestAllowedValues(
		t,
		"AllowedValues",
		ColorModeAuto.AllowedValues(),
		[]string{"auto", "always", "never"},
	)
}

func TestColorModeIsValid(t *testing.T) {
	t.Parallel()

	testhelpers.TestEnumIsValid(t, []ColorMode{
		ColorModeAuto,
		ColorModeAlways,
		ColorModeNever,
		"invalid",
		"",
	}, []bool{
		true,
		true,
		true,
		false,
		false,
	})
}

func TestColorModeShouldColor(t *testing.T) {
	t.Parallel()

	if !ColorModeAlways.ShouldColor() {
		t.Error("ColorModeAlways.ShouldColor() should return true")
	}

	if ColorModeNever.ShouldColor() {
		t.Error("ColorModeNever.ShouldColor() should return false")
	}
}

func TestColorModeShouldColorDefault(t *testing.T) {
	t.Parallel()

	cm := ColorMode("unknown")
	if cm.ShouldColor() {
		t.Error("unknown ColorMode.ShouldColor() should return false")
	}
}

func TestColorModeShouldColorAuto_Deterministic(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "")
	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "")
	t.Setenv("GITLAB_CI", "")
	t.Setenv("JENKINS_URL", "")
	t.Setenv("BUILDKITE", "")
	t.Setenv("GO_OUTPUT_FORCE_COLOR", "")
	t.Setenv("FORCE_COLOR", "")

	// Override detection functions for deterministic testing.
	origTerminal := stdoutIsTerminal
	origNoColor := noColorEnv
	origCI := ciEnv

	t.Cleanup(func() {
		stdoutIsTerminal = origTerminal
		noColorEnv = origNoColor
		ciEnv = origCI
	})

	// Case 1: TTY + no NO_COLOR + not CI → color.
	stdoutIsTerminal = func() bool { return true }
	noColorEnv = func() bool { return false }
	ciEnv = func() bool { return false }

	if !ColorModeAuto.ShouldColor() {
		t.Error("Auto should color when TTY + no NO_COLOR + not CI")
	}

	// Case 2: NO_COLOR set → no color.
	noColorEnv = func() bool { return true }

	if ColorModeAuto.ShouldColor() {
		t.Error("Auto should not color when NO_COLOR is set")
	}

	// Case 3: CI environment → no color.
	noColorEnv = func() bool { return false }
	ciEnv = func() bool { return true }

	if ColorModeAuto.ShouldColor() {
		t.Error("Auto should not color in CI")
	}

	// Case 4: Not a TTY → no color.
	ciEnv = func() bool { return false }
	stdoutIsTerminal = func() bool { return false }

	if ColorModeAuto.ShouldColor() {
		t.Error("Auto should not color when not a TTY")
	}
}

func TestStdoutIsTerminalWithForceColor(t *testing.T) {
	t.Setenv("GO_OUTPUT_FORCE_COLOR", "1")

	if !stdoutIsTerminal() {
		t.Error("stdoutIsTerminal() should return true with GO_OUTPUT_FORCE_COLOR=1")
	}
}
