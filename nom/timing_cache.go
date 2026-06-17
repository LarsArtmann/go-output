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
	cache        map[string][]time.Duration // activity name -> duration history
	filePath     string                     // Path to cache file
	loaded       bool                       // Whether cache has been loaded
	pendingSaves sync.WaitGroup             // tracks in-flight saveAsync goroutines
}

// NewTimingCache creates a new timing cache.
func NewTimingCache() *TimingCache {
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

	return &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   false,
	}
}

// Record records a duration for an activity.
func (tc *TimingCache) Record(activityName string, duration time.Duration) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	// Add duration to history
	history := tc.cache[activityName]
	history = append(history, duration)
	// Keep only last maxCachedEntries entries
	if len(history) > maxCachedEntries {
		history = history[len(history)-maxCachedEntries:]
	}

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

// EnsureLoaded loads the cache if not already loaded.
func (tc *TimingCache) EnsureLoaded() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if tc.loaded {
		return nil
	}

	return tc.loadLocked()
}

func (tc *TimingCache) waitPendingSaves() {
	tc.pendingSaves.Wait()
}
