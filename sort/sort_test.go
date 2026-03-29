// Package sort_test provides tests for the sort package.
package sort

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
)

func TestCompareString(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := CompareString(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareString(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCompareInt(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			if got := CompareInt(tt.a, tt.b); got != tt.want {
				t.Errorf("CompareInt(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestCompareTime(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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

func testItemsAB() []testItem {
	return []testItem{
		{Name: "b", Count: 2, When: time.Time{}},
		{Name: "a", Count: 1, When: time.Time{}},
	}
}

func TestSorter_New(t *testing.T) {
	t.Parallel()
	items := testItemsAB()
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
	t.Parallel()
	items := testItemsAB()
	sorter := New(items, output.SortByName, false)
	result := sorter.WithLessFunc(func(a, b testItem) bool { return a.Count < b.Count })
	if result != sorter {
		t.Error("WithLessFunc() should return the same sorter")
	}
	if sorter.LessFunc == nil {
		t.Error("WithLessFunc() did not set LessFunc")
	}
}

func TestSorter_Sort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		items         []testItem
		sortBy        output.SortBy
		desc          bool
		lessFunc      func(a, b testItem) bool
		expectedNames []string
		expectedCounts []int
	}{
		{
			name:     "SortByName ascending",
			items:    []testItem{{Name: "charlie", Count: 3, When: time.Time{}}, {Name: "alpha", Count: 1, When: time.Time{}}, {Name: "bravo", Count: 2, When: time.Time{}}},
			sortBy:   output.SortByName,
			desc:     false,
			expectedNames: []string{"alpha", "bravo", "charlie"},
		},
		{
			name:     "SortByName descending",
			items:    []testItem{{Name: "charlie", Count: 3, When: time.Time{}}, {Name: "alpha", Count: 1, When: time.Time{}}, {Name: "bravo", Count: 2, When: time.Time{}}},
			sortBy:   output.SortByName,
			desc:     true,
			expectedNames: []string{"charlie", "bravo", "alpha"},
		},
		{
			name:     "SortByCount ascending",
			items:    []testItem{{Name: "charlie", Count: 30, When: time.Time{}}, {Name: "alpha", Count: 10, When: time.Time{}}, {Name: "bravo", Count: 20, When: time.Time{}}},
			sortBy:   output.SortBy("Count"),
			desc:     false,
			expectedCounts: []int{10, 20, 30},
		},
		{
			name:     "SortByWhen ascending",
			items:    []testItem{{Name: "now", Count: 0, When: time.Now()}, {Name: "later", Count: 0, When: time.Now().Add(2 * time.Hour)}, {Name: "earlier", Count: 0, When: time.Now().Add(-2 * time.Hour)}},
			sortBy:   output.SortBy("When"),
			desc:     false,
			expectedNames: []string{"earlier", "now", "later"},
		},
		{
			name:     "CustomLessFunc count descending",
			items:    []testItem{{Name: "b", Count: 20, When: time.Time{}}, {Name: "a", Count: 10, When: time.Time{}}, {Name: "c", Count: 30, When: time.Time{}}},
			sortBy:   output.SortByName,
			desc:     false,
			lessFunc: func(a, b testItem) bool { return a.Count > b.Count },
			expectedCounts: []int{30, 20, 10},
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sorter := New(tt.items, tt.sortBy, tt.desc)
			if tt.lessFunc != nil {
				sorter.WithLessFunc(tt.lessFunc)
			}
			sorter.Sort()

			if len(tt.expectedNames) > 0 {
				for i, expectedName := range tt.expectedNames {
					if got := sorter.Items[i].Name; got != expectedName {
						t.Errorf("Items[%d].Name = %v, want %v", i, got, expectedName)
					}
				}
			}

			if len(tt.expectedCounts) > 0 {
				for i, expectedCount := range tt.expectedCounts {
					if got := sorter.Items[i].Count; got != expectedCount {
						t.Errorf("Items[%d].Count = %v, want %v", i, got, expectedCount)
					}
				}
			}
		})
	}
}

func TestSorter_Sort_EdgeCases(t *testing.T) {
	t.Parallel()

	// Empty slice
	empty := []testItem{}
	New(empty, output.SortByName, false).Sort()
	if len(empty) != 0 {
		t.Errorf("Sort() on empty slice should remain empty")
	}

	// Single item
	single := []testItem{{Name: "only", Count: 1, When: time.Time{}}}
	New(single, output.SortByName, false).Sort()
	if len(single) != 1 || single[0].Name != "only" {
		t.Errorf("Sort() changed single item")
	}

	// Invalid field - should be stable
	invalid := []testItem{
		{Name: "b", Count: 0, When: time.Time{}},
		{Name: "a", Count: 0, When: time.Time{}},
	}
	New(invalid, output.SortBy("NonExistentField"), false).Sort()
	if invalid[0].Name != "b" || invalid[1].Name != "a" {
		t.Errorf("Sort() with invalid field should be stable")
	}
}

func TestToInt(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  any
		want   int
		wantOk bool
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
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
			_, ok := toTime(tt.input)
			if ok != tt.wantOk {
				t.Errorf("toTime(%v) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
		})
	}
}
