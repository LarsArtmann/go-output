package output

import (
	"testing"
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
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("ColorMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorModeAllowedValues(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			if got := tt.mode.IsValid(); got != tt.want {
				t.Errorf("ColorMode.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
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

func TestColorModeToANSI(t *testing.T) {
	t.Parallel()
	// When color is disabled, should return empty string
	if ColorModeNever.ToANSI() != "" {
		t.Error("ColorModeNever.ToANSI() should return empty string")
	}

	// When color is enabled, should return ANSI prefix
	if ColorModeAlways.ToANSI() != "\033[" {
		t.Error("ColorModeAlways.ToANSI() should return ANSI prefix")
	}
}
