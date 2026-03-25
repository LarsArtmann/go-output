package output

import (
	"os"
	"testing"
)

func TestIsNoColor(t *testing.T) {
	t.Parallel()

	orig := os.Getenv("NO_COLOR")
	defer func() { _ = os.Setenv("NO_COLOR", orig) }() // nolint:errcheck
	_ = os.Unsetenv("NO_COLOR")                        // nolint:errcheck
	if isNoColor() {
		t.Error("isNoColor() should return false when NO_COLOR is not set")
	}

	_ = os.Setenv("NO_COLOR", "1") // nolint:errcheck
	if !isNoColor() {
		t.Error("isNoColor() should return true when NO_COLOR is set")
	}
}

func TestIsCI(t *testing.T) {
	t.Parallel()

	ciVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "BUILDKITE"}
	origVals := make(map[string]string)
	for _, v := range ciVars {
		origVals[v] = os.Getenv(v)
		_ = os.Unsetenv(v) // nolint:errcheck
	}
	defer func() {
		for _, v := range ciVars {
			_ = os.Setenv(v, origVals[v]) // nolint:errcheck
		}
	}()

	if isCI() {
		t.Error("isCI() should return false when no CI env vars are set")
	}

	_ = os.Setenv("CI", "true") // nolint:errcheck
	if !isCI() {
		t.Error("isCI() should return true when CI is set")
	}

	_ = os.Unsetenv("CI")                   // nolint:errcheck
	_ = os.Setenv("GITHUB_ACTIONS", "true") // nolint:errcheck
	if !isCI() {
		t.Error("isCI() should return true when GITHUB_ACTIONS is set")
	}
}

func TestIsTerminalByEnv(t *testing.T) {
	t.Parallel()

	orig := os.Getenv("FORCE_COLOR")
	defer func() { _ = os.Setenv("FORCE_COLOR", orig) }() // nolint:errcheck
	_ = os.Unsetenv("FORCE_COLOR")                        // nolint:errcheck
	if isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return false when env is not set")
	}

	_ = os.Setenv("FORCE_COLOR", "0") // nolint:errcheck
	if isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return false when env is '0'")
	}

	_ = os.Setenv("FORCE_COLOR", "1") // nolint:errcheck
	if !isTerminalByEnv("FORCE_COLOR") {
		t.Error("isTerminalByEnv() should return true when env is '1'")
	}
}

func TestParseColorMode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ColorMode
		wantErr bool
	}{
		{"auto", "auto", ColorModeAuto, false},
		{"always", "always", ColorModeAlways, false},
		{"never", "never", ColorModeNever, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColorMode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColorMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseColorMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorModeString(t *testing.T) {
	tests := []struct {
		mode ColorMode
		want string
	}{
		{ColorModeAuto, "auto"},
		{ColorModeAlways, "always"},
		{ColorModeNever, "never"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("ColorMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorModeAllowedValues(t *testing.T) {
	got := ColorModeAuto.AllowedValues()
	want := []string{"auto", "always", "never"}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestColorModeIsValid(t *testing.T) {
	tests := []struct {
		mode ColorMode
		want bool
	}{
		{ColorModeAuto, true},
		{ColorModeAlways, true},
		{ColorModeNever, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			if got := tt.mode.IsValid(); got != tt.want {
				t.Errorf("ColorMode.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorModeShouldColor(t *testing.T) {
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

func TestColorModeToANSI(t *testing.T) {
	// When color is disabled, should return empty string
	if ColorModeNever.ToANSI() != "" {
		t.Error("ColorModeNever.ToANSI() should return empty string")
	}

	// When color is enabled, should return ANSI prefix
	if ColorModeAlways.ToANSI() != "\033[" {
		t.Error("ColorModeAlways.ToANSI() should return ANSI prefix")
	}
}
