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
	mu           sync.RWMutex
	saveMu       sync.Mutex                 // serializes file writes to prevent concurrent truncation
	cache        map[string][]time.Duration // activity name -> duration history
	filePath     string                     // Path to cache file
	loaded       bool                       // Whether cache has been loaded
	pendingSaves sync.WaitGroup             // tracks in-flight saveAsync goroutines
}

// TimingCacheOption configures a TimingCache at construction time.
type TimingCacheOption func(*TimingCache)

// withFilePath returns an option that overrides the default cache file path.
// Used to isolate the cache to a temp directory in tests.
func withFilePath(path string) TimingCacheOption {
	return func(tc *TimingCache) { tc.filePath = path }
}

// NewTimingCache creates a new timing cache backed by ~/.cache/nom-timing.csv
// by default. Pass options to override (e.g. tests inject a temp path).
func NewTimingCache(opts ...TimingCacheOption) *TimingCache {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		homeDir = os.TempDir()
	}

	cachePath := filepath.Join(homeDir, cacheDir, cacheFilename)

	// Validate the parent directory is writable, not /dev/null or similar
	dir := filepath.Dir(filepath.Dir(cachePath))
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		cachePath = filepath.Join(os.TempDir(), cacheDir, cacheFilename)
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
	defer tc.mu.Unlock()
	// Add duration to history
	history := capHistory(append(tc.cache[activityName], duration))

	tc.cache[activityName] = history
	// Save to disk asynchronously (non-blocking)
	tc.pendingSaves.Add(1)
	go tc.saveAsync()

	return nil
}

// medianDuration returns the median of a slice of durations. The input is not
// mutated.
func medianDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
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
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	history, exists := tc.cache[activityName]
	if !exists || len(history) == 0 {
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
	result := make([]time.Duration, len(history))
	copy(result, history)

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
	// concurrently — take its result if so (last-writer-wins on the file,
	// but both read the same file so it doesn't matter).
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.loaded {
		return nil // another goroutine already loaded
	}

	tc.cache = newCache
	tc.loaded = true

	return nil
}

func (tc *TimingCache) waitPendingSaves() {
	tc.pendingSaves.Wait()
}
