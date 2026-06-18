package nom

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestTimingCache(path string, loaded bool) *TimingCache {
	return &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: path,
		loaded:   loaded,
	}
}

// newTempTimingCache returns a TimingCache backed by a per-test temp directory,
// so Record()/saveAsync() writes never touch the real ~/.cache/nom-timing.csv.
// It also waits for any pending async saves on cleanup to avoid goroutine leaks.
func newTempTimingCache(t *testing.T) *TimingCache {
	t.Helper()
	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: filepath.Join(t.TempDir(), cacheFilename),
		loaded:   true,
	}
	t.Cleanup(tc.waitPendingSaves)

	return tc
}

// assertCacheFileExists fails the test if path does not exist.
func assertCacheFileExists(t *testing.T, path, message string) {
	t.Helper()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal(message)
	}
}

func TestNewTimingCache(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	if tc == nil {
		t.Fatal("NewTimingCache() returned nil")
	}

	if tc.IsLoaded() {
		t.Error("new cache should not be loaded")
	}

	if tc.GetFilePath() == "" {
		t.Error("filePath should not be empty")
	}
}

func TestTimingCache_RecordAndGetMedian(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)

	tc.Record("build", 5*time.Second)
	tc.Record("build", 3*time.Second)

	median := tc.GetMedian("build")
	if median == 0 {
		t.Error("expected non-zero median after recording")
	}
}

func TestTimingCache_GetMedian_RobustToOutlier(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)

	tc.Record("build", 3*time.Second)
	tc.Record("build", 3*time.Second)
	tc.Record("build", 3*time.Second)
	tc.Record("build", 3*time.Second)
	tc.Record("build", 60*time.Second) // outlier

	median := tc.GetMedian("build")
	if median != 3*time.Second {
		t.Errorf("median = %v, want %v (outlier ignored)", median, 3*time.Second)
	}
}

func TestTimingCache_GetMedian_NoHistory(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	median := tc.GetMedian("nonexistent")
	if median != 0 {
		t.Errorf("expected 0 for nonexistent activity, got %v", median)
	}
}

func TestTimingCache_GetAll(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)
	tc.Record("build", 2*time.Second)
	tc.Record("test", 4*time.Second)

	all := tc.GetAll()
	if len(all) != 2 {
		t.Errorf("GetAll() returned %d entries, want 2", len(all))
	}
}

func TestTimingCache_GetHistory(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)
	tc.Record("build", 1*time.Second)
	tc.Record("build", 2*time.Second)

	history := tc.getHistory("build")
	if len(history) != 2 {
		t.Fatalf("GetHistory() returned %d entries, want 2", len(history))
	}

	history[0] = 0

	original := tc.getHistory("build")
	if original[0] == 0 {
		t.Error("modifying returned slice should not affect cache")
	}
}

func TestTimingCache_GetHistory_NonExistent(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	history := tc.getHistory("nonexistent")
	if len(history) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(history))
	}
}

func TestTimingCache_Clear(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)
	tc.Record("build", 5*time.Second)
	tc.Clear()

	if tc.GetMedian("build") != 0 {
		t.Error("expected 0 after Clear()")
	}
}

func TestTimingCache_Remove(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)
	tc.Record("build", 5*time.Second)
	tc.remove("build")

	if tc.GetMedian("build") != 0 {
		t.Error("expected 0 after Remove()")
	}
}

func TestTimingCache_MaxEntries(t *testing.T) {
	t.Parallel()

	tc := newTempTimingCache(t)

	for i := range 15 {
		tc.Record("build", time.Duration(i+1)*time.Second)
	}

	history := tc.getHistory("build")
	if len(history) > maxCachedEntries {
		t.Errorf("history has %d entries, want at most %d", len(history), maxCachedEntries)
	}
}

func TestTimingCache_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "test-timing.csv")

	tc := newTestTimingCache(cachePath, true)

	tc.Record("build", 5*time.Second)
	tc.Record("test", 10*time.Second)

	tc.waitPendingSaves()

	if err := tc.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	assertCacheFileExists(t, cachePath, "cache file was not created")

	tc2 := newTestTimingCache(cachePath, false)

	if err := tc2.Load(); err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if !tc2.IsLoaded() {
		t.Error("cache should be loaded after Load()")
	}

	median := tc2.GetMedian("build")
	if median == 0 {
		t.Error("expected non-zero median after loading from file")
	}

	time.Sleep(50 * time.Millisecond)
}

func TestTimingCache_Load_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "nonexistent", "timing.csv")

	tc := newTestTimingCache(cachePath, false)

	if err := tc.Load(); err != nil {
		t.Fatalf("Load() on nonexistent file should not error: %v", err)
	}

	if !tc.IsLoaded() {
		t.Error("cache should be marked loaded even if file doesn't exist")
	}
}

func TestTimingCache_EnsureLoaded(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "timing.csv")

	tc := newTestTimingCache(cachePath, false)
	if err := tc.EnsureLoaded(); err != nil {
		t.Fatalf("EnsureLoaded() error: %v", err)
	}

	if !tc.IsLoaded() {
		t.Error("cache should be loaded after EnsureLoaded()")
	}
}

func TestWriteCacheToFile_InvalidPath(t *testing.T) {
	t.Parallel()

	data := map[string][]time.Duration{
		"build": {5 * time.Second},
	}

	err := writeCacheToFile("/dev/null/impossible/path/timing.csv", data)
	if err == nil {
		t.Error("expected error writing to impossible path")
	}
}

func TestWriteCacheToFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "timing.csv")

	data := map[string][]time.Duration{
		"build": {5 * time.Second, 3 * time.Second},
		"test":  {1 * time.Second},
	}

	if err := writeCacheToFile(cachePath, data); err != nil {
		t.Fatalf("writeCacheToFile() error: %v", err)
	}

	assertCacheFileExists(t, cachePath, "cache file was not created")
}

func TestRecord_TriggersAsyncSave(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "timing.csv")

	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   true,
	}

	tc.Record("build", 5*time.Second)
	tc.waitPendingSaves()

	assertCacheFileExists(t, cachePath, "async save should have created cache file")
}

func TestRecord_AsyncSaveFailureDoesNotBlock(t *testing.T) {
	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: "/dev/null/impossible/path/timing.csv",
		loaded:   true,
	}

	tc.Record("build", 5*time.Second)
	tc.waitPendingSaves()
}
