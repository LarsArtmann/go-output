package envdetect

import "testing"

func TestIsCI(t *testing.T) {
	t.Run("none set", func(t *testing.T) {
		clearCIEnv(t)

		if IsCI() {
			t.Error("IsCI() with no CI env vars = true, want false")
		}
	})

	for _, name := range CIEnvVars {
		t.Run(name, func(t *testing.T) {
			clearCIEnv(t)
			t.Setenv(name, "1")

			if !IsCI() {
				t.Errorf("IsCI() with %s=1 = false, want true", name)
			}
		})
	}
}

func TestIsNoColor(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv("NO_COLOR", "")
		t.Setenv("TERM", "xterm-256color")

		if IsNoColor() {
			t.Error("IsNoColor() with default TERM = true, want false")
		}
	})

	for _, tc := range []struct {
		name string
		env  string
		val  string
	}{
		{"NO_COLOR set", "NO_COLOR", "1"},
		{"TERM dumb", "TERM", "dumb"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(tc.env, tc.val)

			if !IsNoColor() {
				t.Errorf("IsNoColor() with %s=%q = false, want true", tc.env, tc.val)
			}
		})
	}
}

// clearCIEnv clears all CI-related environment variables and restores them
// when the test ends.
func clearCIEnv(t *testing.T) {
	t.Helper()

	for _, name := range CIEnvVars {
		t.Setenv(name, "")
	}
}
