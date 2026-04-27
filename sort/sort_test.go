// Package sort_test provides tests for the sort package.
package sort

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
)

func assertItemField[V comparable](
	t *testing.T,
	items []testItem,
	expected []V,
	accessor func(testItem) V,
	fieldName string,
) {
	for i, expectedVal := range expected {
		if got := accessor(items[i]); got != expectedVal {
			t.Errorf("Items[%d].%s = %v, want %v", i, fieldName, got, expectedVal)
		}
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

	result := sorter.WithLessFunc(compareCount(true))
	if result != sorter {
		t.Error("WithLessFunc() should return the same sorter")
	}

	if sorter.LessFunc == nil {
		t.Error("WithLessFunc() did not set LessFunc")
	}
}

func compareCount(ascending bool) func(a, b testItem) bool {
	return func(a, b testItem) bool {
		if ascending {
			return a.Count < b.Count
		}

		return a.Count > b.Count
	}
}

func testItemsUnsorted() []testItem {
	return testItemsWithCounts(1, 2, 3)
}

func testItemsWithCounts(countAlpha, countBravo, countCharlie int) []testItem {
	return []testItem{
		{Name: "charlie", Count: countCharlie, When: time.Time{}},
		{Name: "alpha", Count: countAlpha, When: time.Time{}},
		{Name: "bravo", Count: countBravo, When: time.Time{}},
	}
}

func TestSorter_Sort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		items          []testItem
		sortBy         output.SortBy
		desc           bool
		lessFunc       func(a, b testItem) bool
		expectedNames  []string
		expectedCounts []int
	}{
		{
			name:           "SortByName ascending",
			items:          testItemsUnsorted(),
			sortBy:         output.SortByName,
			desc:           false,
			lessFunc:       nil,
			expectedNames:  []string{"alpha", "bravo", "charlie"},
			expectedCounts: []int{1, 2, 3},
		},
		{
			name:           "SortByName descending",
			items:          testItemsUnsorted(),
			sortBy:         output.SortByName,
			desc:           true,
			lessFunc:       nil,
			expectedNames:  []string{"charlie", "bravo", "alpha"},
			expectedCounts: []int{3, 2, 1},
		},
		{
			name:           "SortByCount ascending",
			items:          testItemsWithCounts(10, 20, 30),
			sortBy:         output.SortBy("Count"),
			desc:           false,
			lessFunc:       nil,
			expectedNames:  []string{"alpha", "bravo", "charlie"},
			expectedCounts: []int{10, 20, 30},
		},
		{
			name: "SortByWhen ascending",
			items: []testItem{
				{Name: "now", Count: 0, When: time.Now()},
				{Name: "later", Count: 0, When: time.Now().Add(2 * time.Hour)},
				{Name: "earlier", Count: 0, When: time.Now().Add(-2 * time.Hour)},
			},
			sortBy:         output.SortBy("When"),
			desc:           false,
			lessFunc:       nil,
			expectedNames:  []string{"earlier", "now", "later"},
			expectedCounts: []int{0, 0, 0},
		},
		{
			name:           "CustomLessFunc count descending",
			items:          testItemsWithCounts(10, 20, 30),
			sortBy:         output.SortByName,
			desc:           false,
			lessFunc:       compareCount(false),
			expectedNames:  []string{"charlie", "bravo", "alpha"},
			expectedCounts: []int{30, 20, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sorter := New(tt.items, tt.sortBy, tt.desc)
			if tt.lessFunc != nil {
				sorter.WithLessFunc(tt.lessFunc)
			}

			sorter.Sort()

			if len(tt.expectedNames) > 0 {
				assertItemField(
					t,
					sorter.Items,
					tt.expectedNames,
					func(item testItem) string { return item.Name },
					"Name",
				)
			}

			if len(tt.expectedCounts) > 0 {
				assertItemField(
					t,
					sorter.Items,
					tt.expectedCounts,
					func(item testItem) int { return item.Count },
					"Count",
				)
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
	invalid := testItemsAB()
	New(invalid, output.SortBy("NonExistentField"), false).Sort()

	if invalid[0].Name != "b" || invalid[1].Name != "a" {
		t.Errorf("Sort() with invalid field should be stable")
	}
}

type unsignedItem struct {
	Name string
	Size uint64
}

func TestSorter_Sort_UnsignedInt(t *testing.T) {
	t.Parallel()

	items := []unsignedItem{
		{Name: "large", Size: 18_446_744_073_709_551_615}, // max uint64
		{Name: "small", Size: 1},
		{Name: "medium", Size: 100},
	}

	New(items, output.SortBy("Size"), false).Sort()

	if items[0].Name != "small" {
		t.Errorf("Sort() unsigned first = %s, want small", items[0].Name)
	}

	if items[1].Name != "medium" {
		t.Errorf("Sort() unsigned second = %s, want medium", items[1].Name)
	}

	if items[2].Name != "large" {
		t.Errorf("Sort() unsigned third = %s, want large", items[2].Name)
	}
}
