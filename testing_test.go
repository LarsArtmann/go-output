package output

import (
	"testing"
)

func runSubtest(t *testing.T, name string, testFunc func(*testing.T)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		testFunc(t)
	})
}
