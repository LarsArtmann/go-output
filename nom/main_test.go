package nom

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// TestMain redirects the timing cache to a per-process temp directory so that
// parallel tests don't race on the shared ~/.cache/nom-timing.csv. Every test
// that calls NewNOMStyleSubscriber() without WithCachePath gets an isolated
// cache file via cachePathOverride.
func TestMain(m *testing.M) {
	dir := filepath.Join(os.TempDir(), "nom-test-cache-"+strconv.Itoa(os.Getpid()))
	if err := os.MkdirAll(dir, 0o750); err != nil {
		// If we can't create the dir, tests will fall back to the default
		// cache path — the race only affects tests that don't use WithCachePath.
		os.Exit(m.Run())
	}

	cachePathOverride = filepath.Join(dir, cacheFilename)

	code := m.Run()

	_ = os.RemoveAll(dir)

	os.Exit(code)
}
