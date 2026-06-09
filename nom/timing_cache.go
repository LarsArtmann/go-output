package nom

import (
	"os"
	"path/filepath"
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
	cache    map[string][]time.Duration // activity name -> duration history
	filePath string                     // Path to cache file
	loaded   bool                       // Whether cache has been loaded
}

// NewTimingCache creates a new timing cache.
func NewTimingCache() *TimingCache {
	// Determine cache file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir() // Fallback to temp directory
	}

	cachePath := filepath.Join(homeDir, cacheDir, cacheFilename)

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
	go tc.saveAsync()

	return nil
}

// GetAverage returns the average duration for an activity.
func (tc *TimingCache) GetAverage(activityName string) time.Duration {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	history, exists := tc.cache[activityName]
	if !exists || len(history) == 0 {
		return 0
	}
	// Calculate average
	var sum time.Duration
	for _, d := range history {
		sum += d
	}

	return sum / time.Duration(len(history))
}

// GetAll returns all cached averages.
func (tc *TimingCache) GetAll() map[string]time.Duration {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	averages := make(map[string]time.Duration)

	for name, history := range tc.cache {
		if len(history) > 0 {
			var sum time.Duration
			for _, d := range history {
				sum += d
			}

			averages[name] = sum / time.Duration(len(history))
		}
	}

	return averages
}

// GetHistory returns the duration history for an activity.
func (tc *TimingCache) GetHistory(activityName string) []time.Duration {
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

// Remove removes an activity from the cache.
func (tc *TimingCache) Remove(activityName string) {
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
