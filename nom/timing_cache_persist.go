package nom

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Load loads the cache from disk (acquires lock).
func (tc *TimingCache) Load() error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	return tc.loadLocked()
}

// loadLocked loads the cache from disk (caller must hold tc.mu).
func (tc *TimingCache) loadLocked() error {
	if _, err := os.Stat(tc.filePath); os.IsNotExist(err) {
		tc.loaded = true
		return nil
	}

	file, err := os.Open(tc.filePath)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return fmt.Errorf("failed to read cache file: %w, close failed: %w", err, closeErr)
		}

		return fmt.Errorf("failed to read cache file: %w", err)
	}

	newCache := make(map[string][]time.Duration)

	for _, record := range records {
		if len(record) != 2 {
			continue
		}

		activityName := record[0]

		var duration time.Duration

		_, err := fmt.Sscanf(record[1], "%d", &duration)
		if err != nil {
			continue
		}

		history := newCache[activityName]
		history = append(history, duration)
		newCache[activityName] = history
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close cache file %q: %w", tc.filePath, err)
	}

	tc.cache = newCache
	tc.loaded = true

	return nil
}

// snapshotData creates a deep copy of the cache map under RLock, then releases the lock.
func (tc *TimingCache) snapshotData() (map[string][]time.Duration, string) {
	tc.mu.RLock()

	snapshot := make(map[string][]time.Duration, len(tc.cache))
	for name, history := range tc.cache {
		historyCopy := make([]time.Duration, len(history))
		copy(historyCopy, history)
		snapshot[name] = historyCopy
	}

	filePath := tc.filePath
	tc.mu.RUnlock()

	return snapshot, filePath
}

// Save saves the cache to disk (acquires read lock).
func (tc *TimingCache) Save() error {
	cacheSnapshot, filePath := tc.snapshotData()

	return writeCacheToFile(filePath, cacheSnapshot)
}

// writeCacheToFile writes the cache data to the specified file path.
func writeCacheToFile(filePath string, data map[string][]time.Duration) error {
	cacheDirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(cacheDirPath, 0o750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	file, err := os.Create(filePath) //nolint:gosec // G304: path is validated upstream
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}

	writer := csv.NewWriter(file)

	for activityName, history := range data {
		for _, duration := range history {
			record := []string{
				activityName,
				strconv.FormatInt(duration.Nanoseconds(), 10),
			}
			if err := writer.Write(record); err != nil {
				_ = file.Close()
				return fmt.Errorf("failed to write cache record: %w", err)
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close cache file: %w", err)
	}

	return nil
}

// saveAsync saves the cache asynchronously (non-blocking).
// Must be called from a goroutine — snapshots data under RLock, then writes without holding any lock.
func (tc *TimingCache) saveAsync() {
	defer tc.pendingSaves.Done()

	dataCopy, filePath := tc.snapshotData()

	if err := writeCacheToFile(filePath, dataCopy); err != nil {
		log.Printf("nom: async cache save failed: %v", err)
	}
}
