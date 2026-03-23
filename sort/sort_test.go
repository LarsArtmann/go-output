package sort

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
)

func TestCompareString(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a < b", "apple", "banana", -1},
		{"a > b", "banana", "apple", 1},
		{"a == b", "same", "same", 0},
		{"a not string", 123, "banana", 0},
		{"b not string", "apple", 123, 0},
		{"both not string", 123, 456, 0},
		{"empty strings", "", "", 0},
		{"empty vs non-empty", "", "a", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareString(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareString(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCompareInt(t *testing.T) {
	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a < b", 1, 2, -1},
		{"a > b", 2, 1, 1},
		{"a == b", 5, 5, 0},
		{"int8", int8(1), int8(2), -1},
		{"int16", int16(1), int16(2), -1},
		{"int32", int32(1), int32(2), -1},
		{"int64", int64(1), int64(2), -1},
		{"uint", uint(1), uint(2), -1},
		{"uint8", uint8(1), uint8(2), -1},
		{"uint16", uint16(1), uint16(2), -1},
		{"uint32", uint32(1), uint32(2), -1},
		{"uint64", uint64(1), uint64(2), -1},
		{"mixed types", int(1), int64(2), -1},
		{"a not int", "string", 2, 0},
		{"b not int", 1, "string", 0},
		{"both not int", "a", "b", 0},
		{"negative values", -5, 5, -1},
		{"zero values", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareInt(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareInt(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCompareTime(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-time.Hour)
	later := now.Add(time.Hour)

	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a before b", earlier, now, -1},
		{"a after b", later, now, 1},
		{"a equal b", now, now, 0},
		{"pointer to time", &earlier, now, -1},
		{"both pointers", &earlier, &later, -1},
		{"nil pointer", nil, now, 0},
		{"a not time", "string", now, 0},
		{"b not time", now, "string", 0},
		{"both not time", "a", "b", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareTime(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareTime(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

type testItem struct {
	Name  string
	Count int
	When  time.Time
}

func TestSorter_New(t *testing.T) {
	items := []testItem{
		{Name: "b", Count: 2},
		{Name: "a", Count: 1},
	}
	sorter := New(items, output.SortByName, false)
	if sorter == nil {
		t.Fatal("New() returned nil")
	}
	if len(sorter.Items) != 2 {
		t.Errorf("New() items length = %d, want 2", len(sorter.Items))
	}
	if sorter.By != output.SortByName {
		t.Errorf("New() By = %v, want %v", sorter.By, output.SortByName)
	}
	if sorter.Desc != false {
		t.Errorf("New() Desc = %v, want false", sorter.Desc)
	}
}

func TestSorter_WithLessFunc(t *testing.T) {
	items := []testItem{
		{Name: "b", Count: 2},
		{Name: "a", Count: 1},
	}
	sorter := New(items, output.SortByName, false)
	result := sorter.WithLessFunc(func(a, b testItem) bool {
		return a.Count < b.Count
	})
	if result != sorter {
		t.Error("WithLessFunc() should return the same sorter")
	}
	if sorter.LessFunc == nil {
		t.Error("WithLessFunc() did not set LessFunc")
	}
}

func TestSorter_Sort_ByName(t *testing.T) {
	items := []testItem{
		{Name: "charlie", Count: 3},
		{Name: "alpha", Count: 1},
		{Name: "bravo", Count: 2},
	}

	sorter := New(items, output.SortByName, false)
	sorter.Sort()

	if items[0].Name != "alpha" {
		t.Errorf("Sort() first item = %s, want alpha", items[0].Name)
	}
	if items[1].Name != "bravo" {
		t.Errorf("Sort() second item = %s, want bravo", items[1].Name)
	}
	if items[2].Name != "charlie" {
		t.Errorf("Sort() third item = %s, want charlie", items[2].Name)
	}
}

func TestSorter_Sort_ByNameDesc(t *testing.T) {
	items := []testItem{
		{Name: "charlie", Count: 3},
		{Name: "alpha", Count: 1},
		{Name: "bravo", Count: 2},
	}

	sorter := New(items, output.SortByName, true)
	sorter.Sort()

	if items[0].Name != "charlie" {
		t.Errorf("Sort() desc first item = %s, want charlie", items[0].Name)
	}
	if items[1].Name != "bravo" {
		t.Errorf("Sort() desc second item = %s, want bravo", items[1].Name)
	}
	if items[2].Name != "alpha" {
		t.Errorf("Sort() desc third item = %s, want alpha", items[2].Name)
	}
}

func TestSorter_Sort_ByCount(t *testing.T) {
	items := []testItem{
		{Name: "charlie", Count: 30},
		{Name: "alpha", Count: 10},
		{Name: "bravo", Count: 20},
	}

	sorter := New(items, output.SortBy("Count"), false)
	sorter.Sort()

	if items[0].Count != 10 {
		t.Errorf("Sort() by count first = %d, want 10", items[0].Count)
	}
	if items[1].Count != 20 {
		t.Errorf("Sort() by count second = %d, want 20", items[1].Count)
	}
	if items[2].Count != 30 {
		t.Errorf("Sort() by count third = %d, want 30", items[2].Count)
	}
}

func TestSorter_Sort_ByTime(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-2 * time.Hour)
	later := now.Add(2 * time.Hour)

	items := []testItem{
		{Name: "now", When: now},
		{Name: "later", When: later},
		{Name: "earlier", When: earlier},
	}

	sorter := New(items, output.SortBy("When"), false)
	sorter.Sort()

	if items[0].Name != "earlier" {
		t.Errorf("Sort() by time first = %s, want earlier", items[0].Name)
	}
	if items[1].Name != "now" {
		t.Errorf("Sort() by time second = %s, want now", items[1].Name)
	}
	if items[2].Name != "later" {
		t.Errorf("Sort() by time third = %s, want later", items[2].Name)
	}
}

func TestSorter_Sort_CustomLessFunc(t *testing.T) {
	items := []testItem{
		{Name: "b", Count: 20},
		{Name: "a", Count: 10},
		{Name: "c", Count: 30},
	}

	sorter := New(items, output.SortByName, false)
	sorter.WithLessFunc(func(a, b testItem) bool {
		return a.Count > b.Count // Sort by count descending
	})
	sorter.Sort()

	if items[0].Count != 30 {
		t.Errorf("Custom LessFunc first = %d, want 30", items[0].Count)
	}
	if items[1].Count != 20 {
		t.Errorf("Custom LessFunc second = %d, want 20", items[1].Count)
	}
	if items[2].Count != 10 {
		t.Errorf("Custom LessFunc third = %d, want 10", items[2].Count)
	}
}

func TestSorter_Sort_EmptySlice(t *testing.T) {
	items := []testItem{}

	sorter := New(items, output.SortByName, false)
	sorter.Sort() // Should not panic

	if len(items) != 0 {
		t.Errorf("Sort() on empty slice should remain empty, got %d items", len(items))
	}
}

func TestSorter_Sort_SingleItem(t *testing.T) {
	items := []testItem{{Name: "only", Count: 1}}

	sorter := New(items, output.SortByName, false)
	sorter.Sort()

	if len(items) != 1 {
		t.Errorf("Sort() changed slice length")
	}
	if items[0].Name != "only" {
		t.Errorf("Sort() changed single item")
	}
}

func TestSorter_Sort_InvalidField(t *testing.T) {
	items := []testItem{
		{Name: "b", Count: 2},
		{Name: "a", Count: 1},
	}

	sorter := New(items, output.SortBy("NonExistentField"), false)
	sorter.Sort() // Should not panic, order should remain stable

	// Since field doesn't exist, sort should be stable (no change)
	if items[0].Name != "b" || items[1].Name != "a" {
		t.Errorf("Sort() with invalid field should be stable")
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantOk  bool
	}{
		{"int", int(42), 42, true},
		{"int8", int8(42), 42, true},
		{"int16", int16(42), 42, true},
		{"int32", int32(42), 42, true},
		{"int64", int64(42), 42, true},
		{"uint", uint(42), 42, true},
		{"uint8", uint8(42), 42, true},
		{"uint16", uint16(42), 42, true},
		{"uint32", uint32(42), 42, true},
		{"uint64", uint64(42), 42, true},
		{"string", "42", 0, false},
		{"float", 3.14, 0, false},
		{"nil", nil, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := toInt(tt.input)
			if ok != tt.wantOk {
				t.Errorf("toInt(%v) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("toInt(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		input  any
		wantOk bool
	}{
		{"time.Time", now, true},
		{"*time.Time", &now, true},
		{"nil pointer", (*time.Time)(nil), false},
		{"string", "2023-01-01", false},
		{"int", 123, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := toTime(tt.input)
			if ok != tt.wantOk {
				t.Errorf("toTime(%v) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
		})
	}
}
