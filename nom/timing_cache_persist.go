package nom

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
)

// Load loads the cache from disk (acquires lock).
func (tc *TimingCache) Load() error {
	newCache, err := readCacheFile(tc.filePath)
	if err != nil {
		return err
	}

	tc.publishCache(newCache)

	return nil
}

// readCacheFile reads and parses the cache file without holding any lock.
// Returns an empty (non-nil) map if the file doesn't exist.
func readCacheFile(filePath string) (map[string][]time.Duration, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return make(map[string][]time.Duration), nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cache file: %w", err)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to read cache file: %w, close failed: %w", err, closeErr)
		}

		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	newCache := make(map[string][]time.Duration)

	for _, record := range records {
		if len(record) != 2 {
			continue
		}

		activityName := record[0]

		nanos, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			continue
		}

		duration := time.Duration(nanos)

		history := newCache[activityName]
		// Cap at maxCachedEntries during load to prevent unbounded growth from hand-edited files
		newCache[activityName] = capHistory(append(history, duration))
	}

	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close cache file %q: %w", filePath, err)
	}

	return newCache, nil
}

// snapshotData creates a deep copy of the cache map under RLock, then releases the lock.
func (tc *TimingCache) snapshotData() (map[string][]time.Duration, string) {
	tc.mu.RLock()

	snapshot := make(map[string][]time.Duration, len(tc.cache))
	for name, history := range tc.cache {
		historyCopy := make([]time.Duration, 0, len(history))
		historyCopy = append(historyCopy, history...)
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
// The write is atomic (sibling temp file + rename): cache instances default
// to one shared path, so a plain truncate-and-write would tear the file under
// concurrent readers and fail their CSV parse. Readers then see either the
// complete previous file or the complete new one.
func writeCacheToFile(filePath string, data map[string][]time.Duration) error {
	cacheDirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(cacheDirPath, 0o750); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	tempFile, err := os.CreateTemp(cacheDirPath, filepath.Base(filePath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp cache file: %w", err)
	}
	tempPath := tempFile.Name()

	keepTemp := false
	defer func() {
		if keepTemp {
			return
		}

		_ = os.Remove(tempPath)
	}()

	writer := csv.NewWriter(tempFile)

	names := make([]string, 0, len(data))
	for name := range data {
		names = append(names, name)
	}

	slices.Sort(names)

	for _, activityName := range names {
		for _, duration := range data[activityName] {
			record := []string{
				activityName,
				strconv.FormatInt(duration.Nanoseconds(), 10),
			}
			if err := writer.Write(record); err != nil {
				_ = tempFile.Close()
				return fmt.Errorf("failed to write cache record: %w", err)
			}
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		_ = tempFile.Close()
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close cache file: %w", err)
	}

	if err := os.Rename(tempPath, filePath); err != nil {
		return fmt.Errorf("failed to replace cache file: %w", err)
	}

	keepTemp = true

	return nil
}
