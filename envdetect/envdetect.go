// Package envdetect provides shared environment-variable inspection helpers
// for output rendering decisions. Centralizing the CI / NO_COLOR detection
// logic here keeps root, nom, and any other module that needs to suppress
// color output in lockstep.
package envdetect

import "os"

// CIEnvVars are the environment variables that indicate execution inside a
// common CI provider. Exported so tests and downstream modules can keep their
// own list of "known suppressors" aligned with this one.
var CIEnvVars = []string{
	"CI",
	"GITHUB_ACTIONS",
	"GITLAB_CI",
	"JENKINS_URL",
	"BUILDKITE",
}

// IsCI reports whether the process appears to be running inside a CI provider
// by checking the well-known CI environment variables.
//
// This is the single source of truth for CI detection across the go-output
// modules. Root uses it to suppress color in ColorModeAuto, and nom/ uses it
// for the same purpose in the inline renderer. The shared helper prevents the
// two implementations from drifting.
func IsCI() bool {
	for _, name := range CIEnvVars {
		if os.Getenv(name) != "" {
			return true
		}
	}

	return false
}

// IsNoColor reports whether color output should be suppressed per the NO_COLOR
// convention (https://no-color.org) or a TERM=dumb terminal.
func IsNoColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return true
	}

	return os.Getenv("TERM") == "dumb"
}
