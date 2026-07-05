package nom

import (
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

// ============================================================================
// TIMING CACHE
// ============================================================================.
const (
	// maxCachedEntries is the maximum number of entries to keep per activity.
	maxCachedEntries = 10
	// cacheFilename is the name of the cache file.
	cacheFilename = "nom-timing.csv"
	// cacheDir is the directory where the cache is stored.
	cacheDir = ".cache"
)

// TimingCache manages activity duration history and averages.
type TimingCache struct {
	mu       sync.RWMutex
	saveMu   sync.Mutex                 // serializes file writes to prevent concurrent truncation
	cache    map[string][]time.Duration // activity name -> duration history
	filePath string                     // Path to cache file
	loaded   bool                       // Whether cache has been loaded

	// Background saver — a single goroutine replaces per-call go saveAsync()
	// to prevent unbounded goroutine spawning (C4 fix). saveCh is buffered(1):
	// non-blocking sends coalesce rapid Record() calls into one disk write.
	// saverMu guards the lifecycle (start/stop); saveMu serializes disk writes.
	saveCh   chan struct{}
	saveDone chan struct{}
	saverMu  sync.Mutex
	saveErr  error // last async save error (guarded by saveMu)
}

// TimingCacheOption configures a TimingCache at construction time.
type TimingCacheOption func(*TimingCache)

// withFilePath returns an option that overrides the default cache file path.
// Used to isolate the cache to a temp directory in tests.
func withFilePath(path string) TimingCacheOption {
	return func(tc *TimingCache) { tc.filePath = path }
}

// cachePathOverride allows test binaries to redirect the default cache to a
// per-process temp directory, avoiding races when tests run in parallel.
// Production code never sets this — it defaults to empty, falling through to
// the standard ~/.cache/nom-timing.csv path.
//
//nolint:gochecknoglobals // intentional test-injection point
var cachePathOverride string

// NewTimingCache creates a new timing cache backed by ~/.cache/nom-timing.csv
// by default. Pass options to override (e.g. tests inject a temp path).
func NewTimingCache(opts ...TimingCacheOption) *TimingCache {
	cachePath := cachePathOverride
	if cachePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil || homeDir == "" {
			homeDir = os.TempDir()
		}

		cachePath = filepath.Join(homeDir, cacheDir, cacheFilename)

		// Validate the parent directory is writable, not /dev/null or similar
		dir := filepath.Dir(filepath.Dir(cachePath))
		if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
			cachePath = filepath.Join(os.TempDir(), cacheDir, cacheFilename)
		}
	}

	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   false,
	}

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

// capHistory trims history to at most maxCachedEntries entries, keeping the
// most recent samples. Used by both Record (in-memory) and load (disk) so the
// cap is enforced consistently regardless of where the data came from.
func capHistory(history []time.Duration) []time.Duration {
	if len(history) > maxCachedEntries {
		history = history[len(history)-maxCachedEntries:]
	}

	return history
}

// Record records a duration for an activity.
func (tc *TimingCache) Record(activityName string, duration time.Duration) error {
	tc.mu.Lock()
	tc.cache[activityName] = capHistory(append(tc.cache[activityName], duration))
	tc.mu.Unlock()

	// Trigger a background save via the single saver goroutine (non-blocking,
	// coalesced — no goroutine per call).
	tc.triggerSave()

	return nil
}

// triggerSave starts the background saver if needed and sends a non-blocking
// save signal. If a save is already pending, the signal is coalesced.
// Caller-safe; uses saverMu to guard the channel lifecycle.
func (tc *TimingCache) triggerSave() {
	tc.saverMu.Lock()
	defer tc.saverMu.Unlock()

	if tc.saveCh == nil {
		tc.saveCh = make(chan struct{}, 1)
		tc.saveDone = make(chan struct{})
		ch := tc.saveCh
		done := tc.saveDone

		go func() {
			defer close(done)

			for range ch {
				tc.performAsyncSave()
			}
		}()
	}

	select {
	case tc.saveCh <- struct{}{}:
	default: // a save is already pending; coalesce
	}
}

// saveSnapshot serializes a disk write with any concurrent async saver via
// saveMu, records any write error in saveErr, and returns the (possibly
// overwritten) saved error. Both performAsyncSave and Flush funnel through
// here so the serialize-write-and-record-error sequence lives in one place.
func (tc *TimingCache) saveSnapshot() error {
	tc.saveMu.Lock()
	defer tc.saveMu.Unlock()

	dataCopy, filePath := tc.snapshotData()

	if err := writeCacheToFile(filePath, dataCopy); err != nil {
		tc.saveErr = err
	}

	return tc.saveErr
}

// performAsyncSave snapshots the in-memory cache and writes it to disk.
// Any error is stored in saveErr for later retrieval via Flush().
func (tc *TimingCache) performAsyncSave() {
	_ = tc.saveSnapshot()
}

// stopSaver closes the background saver goroutine and waits for it to exit.
// After stopSaver, the next Record() call starts a fresh saver.
func (tc *TimingCache) stopSaver() {
	tc.saverMu.Lock()
	ch := tc.saveCh
	done := tc.saveDone
	tc.saveCh = nil
	tc.saveDone = nil
	tc.saverMu.Unlock()

	if ch != nil {
		close(ch)
		<-done
	}
}

// Flush ensures all pending data is persisted to disk and returns the last
// async save error (if any). Safe to call concurrently with Record().
// The synchronous write serializes with any in-flight async save via saveMu.
func (tc *TimingCache) Flush() error {
	err := tc.saveSnapshot()

	tc.saveMu.Lock()
	tc.saveErr = nil
	tc.saveMu.Unlock()

	return err
}

// medianDuration returns the median of a slice of durations. The input is not
// mutated.
func medianDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, 0, len(durations))
	sorted = append(sorted, durations...)
	slices.Sort(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return sorted[mid]
	}

	return (sorted[mid-1] + sorted[mid]) / 2
}

// GetMedian returns the median duration for an activity. Median is more robust
// than mean when one run is an outlier (e.g. cold cache).
func (tc *TimingCache) GetMedian(activityName string) time.Duration {
	history := tc.getHistory(activityName)
	if len(history) == 0 {
		return 0
	}

	return medianDuration(history)
}

// GetAll returns all cached medians.
func (tc *TimingCache) GetAll() map[string]time.Duration {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	medians := make(map[string]time.Duration)

	for name, history := range tc.cache {
		if len(history) > 0 {
			medians[name] = medianDuration(history)
		}
	}

	return medians
}

func (tc *TimingCache) getHistory(activityName string) []time.Duration {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	history, exists := tc.cache[activityName]
	if !exists {
		return make([]time.Duration, 0)
	}
	// Return a copy to prevent external modification
	result := make([]time.Duration, 0, len(history))
	result = append(result, history...)

	return result
}

// Clear removes all entries from the cache.
func (tc *TimingCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cache = make(map[string][]time.Duration)
}

func (tc *TimingCache) remove(activityName string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	delete(tc.cache, activityName)
}

// IsLoaded returns true if the cache has been loaded from disk.
func (tc *TimingCache) IsLoaded() bool {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	return tc.loaded
}

// GetFilePath returns the cache file path.
func (tc *TimingCache) GetFilePath() string {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	return tc.filePath
}

// publishCache installs newCache atomically if not already loaded.
// Caller must have produced newCache outside this lock (e.g. via readCacheFile).
// Returns true if installed, false if another goroutine already loaded.
func (tc *TimingCache) publishCache(newCache map[string][]time.Duration) bool {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.loaded {
		return false
	}

	tc.cache = newCache
	tc.loaded = true

	return true
}

// EnsureLoaded loads the cache if not already loaded. File I/O happens outside
// the lock so concurrent GetMedian/GetAll calls are not blocked during disk reads.
func (tc *TimingCache) EnsureLoaded() error {
	tc.mu.RLock()

	if tc.loaded {
		tc.mu.RUnlock()
		return nil
	}

	filePath := tc.filePath
	tc.mu.RUnlock()

	// Read and parse the file without holding any lock.
	newCache, err := readCacheFile(filePath)
	if err != nil {
		return err
	}

	// Publish under the write lock. Another goroutine may have loaded
	// concurrently — last-writer-wins on the file, but both read the same
	// file so it doesn't matter.
	tc.publishCache(newCache)

	return nil
}

// waitPendingSaves flushes pending data and stops the background saver goroutine.
// Test-only helper; production code should use Flush().
func (tc *TimingCache) waitPendingSaves() {
	_ = tc.Flush()
	tc.stopSaver()
}
