package nom

import (
	"os"
	"testing"
)

func TestDetectAutoTheme_NoEnvVar(t *testing.T) {
	t.Setenv("COLORFGBG", "")

	theme := detectAutoTheme()
	if theme.Colors != ThemeDefault.Colors {
		t.Error("expected ThemeDefault when COLORFGBG is unset")
	}
}

func TestDetectAutoTheme_DarkBackground(t *testing.T) {
	t.Setenv("COLORFGBG", "15;0")

	theme := detectAutoTheme()
	if theme.Colors != ThemeDefault.Colors {
		t.Error("expected ThemeDefault for dark background (bg=0)")
	}
}

func TestDetectAutoTheme_LightBackground(t *testing.T) {
	t.Setenv("COLORFGBG", "0;15")

	theme := detectAutoTheme()
	if theme.Colors != ThemeHighContrast.Colors {
		t.Error("expected ThemeHighContrast for light background (bg=15)")
	}
}

func TestDetectAutoTheme_Malformed(t *testing.T) {
	t.Setenv("COLORFGBG", "garbage")

	theme := detectAutoTheme()
	if theme.Colors != ThemeDefault.Colors {
		t.Error("expected ThemeDefault for malformed COLORFGBG")
	}
}

func TestDetectAutoTheme_PartialMalformed(t *testing.T) {
	t.Setenv("COLORFGBG", "0;abc")

	theme := detectAutoTheme()
	if theme.Colors != ThemeDefault.Colors {
		t.Error("expected ThemeDefault for partially malformed COLORFGBG")
	}
}

func TestWithAutoTheme_Applied(t *testing.T) {
	t.Setenv("COLORFGBG", "0;15")

	sub := NewNOMStyleSubscriber(WithAutoTheme())

	if sub.theme.Colors != ThemeHighContrast.Colors {
		t.Error("expected ThemeHighContrast applied to subscriber")
	}
}

func TestWithAutoTheme_UnsetDefaultsToDefault(t *testing.T) {
	os.Unsetenv("COLORFGBG")

	sub := NewNOMStyleSubscriber(WithAutoTheme())

	if sub.theme.Colors != ThemeDefault.Colors {
		t.Error("expected ThemeDefault when COLORFGBG unset")
	}
}
