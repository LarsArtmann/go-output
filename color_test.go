package output

import (
	"testing"
)

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
