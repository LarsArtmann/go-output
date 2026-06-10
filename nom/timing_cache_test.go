package nom

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

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

func TestTimingCache_RecordAndGetAverage(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	tc.Record("build", 5*time.Second)
	tc.Record("build", 3*time.Second)

	avg := tc.GetAverage("build")
	if avg == 0 {
		t.Error("expected non-zero average after recording")
	}
}

func TestTimingCache_GetAverage_NoHistory(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	avg := tc.GetAverage("nonexistent")
	if avg != 0 {
		t.Errorf("expected 0 for nonexistent activity, got %v", avg)
	}
}

func TestTimingCache_GetAll(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	tc.Record("build", 2*time.Second)
	tc.Record("test", 4*time.Second)

	all := tc.GetAll()
	if len(all) != 2 {
		t.Errorf("GetAll() returned %d entries, want 2", len(all))
	}
}

func TestTimingCache_GetHistory(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	tc.Record("build", 1*time.Second)
	tc.Record("build", 2*time.Second)

	history := tc.GetHistory("build")
	if len(history) != 2 {
		t.Fatalf("GetHistory() returned %d entries, want 2", len(history))
	}

	history[0] = 0

	original := tc.GetHistory("build")
	if original[0] == 0 {
		t.Error("modifying returned slice should not affect cache")
	}
}

func TestTimingCache_GetHistory_NonExistent(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	history := tc.GetHistory("nonexistent")
	if len(history) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(history))
	}
}

func TestTimingCache_Clear(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	tc.Record("build", 5*time.Second)
	tc.Clear()

	if tc.GetAverage("build") != 0 {
		t.Error("expected 0 after Clear()")
	}
}

func TestTimingCache_Remove(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	tc.Record("build", 5*time.Second)
	tc.Remove("build")

	if tc.GetAverage("build") != 0 {
		t.Error("expected 0 after Remove()")
	}
}

func TestTimingCache_MaxEntries(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()

	for i := range 15 {
		tc.Record("build", time.Duration(i+1)*time.Second)
	}

	history := tc.GetHistory("build")
	if len(history) > maxCachedEntries {
		t.Errorf("history has %d entries, want at most %d", len(history), maxCachedEntries)
	}
}

func TestTimingCache_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "test-timing.csv")

	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   true,
	}

	tc.Record("build", 5*time.Second)
	tc.Record("test", 10*time.Second)

	tc.WaitPendingSaves()

	if err := tc.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("cache file was not created")
	}

	tc2 := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   false,
	}

	if err := tc2.Load(); err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if !tc2.IsLoaded() {
		t.Error("cache should be loaded after Load()")
	}

	avg := tc2.GetAverage("build")
	if avg == 0 {
		t.Error("expected non-zero average after loading from file")
	}

	time.Sleep(50 * time.Millisecond)
}

func TestTimingCache_Load_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	cachePath := filepath.Join(tmpDir, "nonexistent", "timing.csv")

	tc := &TimingCache{
		cache:    make(map[string][]time.Duration),
		filePath: cachePath,
		loaded:   false,
	}

	if err := tc.Load(); err != nil {
		t.Fatalf("Load() on nonexistent file should not error: %v", err)
	}

	if !tc.IsLoaded() {
		t.Error("cache should be marked loaded even if file doesn't exist")
	}
}

func TestTimingCache_EnsureLoaded(t *testing.T) {
	t.Parallel()

	tc := NewTimingCache()
	if err := tc.EnsureLoaded(); err != nil {
		t.Fatalf("EnsureLoaded() error: %v", err)
	}

	if !tc.IsLoaded() {
		t.Error("cache should be loaded after EnsureLoaded()")
	}
}
