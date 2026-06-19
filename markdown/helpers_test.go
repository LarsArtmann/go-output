package markdown

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

// Shared with the root package; aliased here so moved tests read identically.
var assertContains = testhelpers.AssertContains

func runSubtest(t *testing.T, name string, testFunc func(*testing.T)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		testFunc(t)
	})
}
