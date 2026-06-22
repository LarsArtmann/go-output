package output

import "os"

// CIEnvVars are the environment variables that indicate execution inside a
// common CI provider. Exported so tests and downstream code can keep their
// own list of "known suppressors" aligned with this one.
//
//nolint:gochecknoglobals // Exported so tests and downstream code can align their CI detection.
var CIEnvVars = []string{
	"CI",
	"GITHUB_ACTIONS",
	"GITLAB_CI",
	"JENKINS_URL",
	"BUILDKITE",
}

// IsCI reports whether the process appears to be running inside a CI provider
// by checking the well-known CI environment variables.
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
