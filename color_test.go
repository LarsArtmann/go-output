package output

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestIsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "") // Clear NO_COLOR

	if isNoColor() {
		t.Error("isNoColor() should return false when NO_COLOR is not set")
	}

	t.Setenv("NO_COLOR", "1")

	if !isNoColor() {
		t.Error("isNoColor() should return true when NO_COLOR is set")
	}
}

func TestIsCI(t *testing.T) {
	ciVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "BUILDKITE"}
	for _, v := range ciVars {
		t.Setenv(v, "") // Clear CI vars
	}

	if isCI() {
		t.Error("isCI() should return false when no CI env vars are set")
	}

	t.Setenv("CI", "true")

	if !isCI() {
		t.Error("isCI() should return true when CI is set")
	}

	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "true")

	if !isCI() {
		t.Error("isCI() should return true when GITHUB_ACTIONS is set")
	}
}

func TestIsTerminalByEnv(t *testing.T) {
	t.Setenv("FORCE_COLOR", "") // Clear FORCE_COLOR

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
	// Test explicit modes
	if !ColorModeAlways.ShouldColor() {
		t.Error("ColorModeAlways.ShouldColor() should return true")
	}

	if ColorModeNever.ShouldColor() {
		t.Error("ColorModeNever.ShouldColor() should return false")
	}

	// Auto mode depends on environment
	_ = ColorModeAuto.ShouldColor() // Just ensure it doesn't panic
}

func TestColorModeShouldColorDefault(t *testing.T) {
	t.Parallel()

	cm := ColorMode("unknown")
	if cm.ShouldColor() {
		t.Error("unknown ColorMode.ShouldColor() should return false")
	}
}

func TestIsStdoutTerminalWithForceColor(t *testing.T) {
	t.Setenv("GO_OUTPUT_FORCE_COLOR", "1")

	if !isStdoutTerminal() {
		t.Error("isStdoutTerminal() should return true with GO_OUTPUT_FORCE_COLOR=1")
	}
}
