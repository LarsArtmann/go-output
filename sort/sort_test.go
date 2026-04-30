// Package sort_test provides tests for the sort package.
package sort

import (
	"sort"
	"testing"
	"time"

	output "github.com/larsartmann/go-output"
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

func compareName(_ bool) func(a, b testItem) bool {
	return ByField(func(item testItem) string { return item.Name })
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
			lessFunc:       compareName(true),
			expectedNames:  []string{"alpha", "bravo", "charlie"},
			expectedCounts: []int{1, 2, 3},
		},
		{
			name:           "SortByName descending",
			items:          testItemsUnsorted(),
			sortBy:         output.SortByName,
			desc:           true,
			lessFunc:       compareName(true),
			expectedNames:  []string{"charlie", "bravo", "alpha"},
			expectedCounts: []int{3, 2, 1},
		},
		{
			name:           "SortByCount ascending",
			items:          testItemsWithCounts(10, 20, 30),
			sortBy:         output.SortBy("Count"),
			desc:           false,
			lessFunc:       ByField(func(item testItem) int { return item.Count }),
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
			sortBy: output.SortBy("When"),
			desc:   false,
			lessFunc: func(a, b testItem) bool {
				return a.When.Before(b.When)
			},
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

	// Empty slice — no-op
	empty := []testItem{}
	New(empty, output.SortByName, false).Sort()

	if len(empty) != 0 {
		t.Errorf("Sort() on empty slice should remain empty")
	}

	// Single item — no-op
	single := []testItem{{Name: "only", Count: 1, When: time.Time{}}}
	New(single, output.SortByName, false).Sort()

	if len(single) != 1 || single[0].Name != "only" {
		t.Errorf("Sort() changed single item")
	}

	// No LessFunc — should be stable (no-op)
	items := testItemsAB()
	New(items, output.SortByName, false).Sort()

	if items[0].Name != "b" || items[1].Name != "a" {
		t.Errorf("Sort() without LessFunc should be stable (no-op)")
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

	New(items, output.SortBy("Size"), false).
		WithLessFunc(ByField(func(item unsignedItem) uint64 { return item.Size })).
		Sort()

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

type stableItem struct {
	Name  string
	Order int
}

func TestSorter_Sort_DescStability(t *testing.T) {
	t.Parallel()

	items := []stableItem{
		{Name: "a", Order: 1},
		{Name: "a", Order: 2},
		{Name: "a", Order: 3},
		{Name: "b", Order: 4},
	}

	New(items, output.SortByName, true).
		WithLessFunc(ByField(func(item stableItem) string { return item.Name })).
		Sort()

	if items[0].Name != "b" {
		t.Fatalf("desc sort first = %s, want b", items[0].Name)
	}

	// Equal-name items must preserve original insertion order
	if items[1].Order != 1 || items[2].Order != 2 || items[3].Order != 3 {
		t.Errorf("desc stable sort: orders = [%d, %d, %d], want [1, 2, 3]",
			items[1].Order, items[2].Order, items[3].Order)
	}
}

func TestSorter_Sort_DescCount(t *testing.T) {
	t.Parallel()

	items := testItemsWithCounts(10, 20, 30)

	New(items, output.SortBy("Count"), true).WithLessFunc(compareCount(true)).Sort()

	if items[0].Name != "charlie" {
		t.Errorf("desc sort first = %s, want charlie", items[0].Name)
	}

	if items[2].Name != "alpha" {
		t.Errorf("desc sort last = %s, want alpha", items[2].Name)
	}
}

func TestSorter_Sort_NonStructInput(t *testing.T) {
	t.Parallel()

	items := []string{"b", "a", "c"}

	// No LessFunc — should be stable (no-op)
	New(items, output.SortByName, false).Sort()

	if items[0] != "b" {
		t.Errorf("sort without LessFunc should be stable (no-op), got %s", items[0])
	}

	// With LessFunc — should sort
	New(items, output.SortByName, false).WithLessFunc(func(a, b string) bool {
		return a < b
	}).Sort()

	if items[0] != "a" {
		t.Errorf("sort with LessFunc first = %s, want a", items[0])
	}
}

func TestSorter_Sort_UsesSliceStable(t *testing.T) {
	t.Parallel()

	// Verify that Sort() delegates to sort.SliceStable by checking stability
	items := []testItem{
		{Name: "same", Count: 1, When: time.Time{}},
		{Name: "same", Count: 2, When: time.Time{}},
		{Name: "same", Count: 3, When: time.Time{}},
	}

	New(items, output.SortByName, false).
		WithLessFunc(ByField(func(item testItem) string { return item.Name })).
		Sort()

	// All names are equal, so relative order must be preserved
	if items[0].Count != 1 || items[1].Count != 2 || items[2].Count != 3 {
		t.Errorf("stable sort not preserved: counts = [%d, %d, %d], want [1, 2, 3]",
			items[0].Count, items[1].Count, items[2].Count)
	}
}

func TestNew_WithLessFuncChaining(t *testing.T) {
	t.Parallel()

	items := testItemsWithCounts(30, 10, 20)

	New(items, output.SortBy("Count"), false).
		WithLessFunc(ByField(func(item testItem) int { return item.Count })).
		Sort()

	if items[0].Name != "bravo" || items[0].Count != 10 {
		t.Errorf(
			"chained sort first = %s (count %d), want bravo (count 10)",
			items[0].Name, items[0].Count,
		)
	}
}

func TestSorter_Internals(t *testing.T) {
	t.Parallel()

	// Verify Sorter uses sort.SliceStable under the hood
	// by checking the standard library's stability guarantee
	items := make([]int, 100)
	for i := range items {
		items[i] = i % 10 // 10 groups of 10
	}

	New(items, output.SortByName, false).WithLessFunc(ByField(
		func(item int) int { return item },
	)).Sort()

	if !sort.SliceIsSorted(items, func(i, j int) bool {
		return items[i] < items[j]
	}) {
		t.Error("items should be sorted ascending")
	}
}

func TestByField(t *testing.T) {
	t.Parallel()

	t.Run("string field", func(t *testing.T) {
		t.Parallel()

		items := testItemsUnsorted()
		New(items, output.SortByName, false).
			WithLessFunc(ByField(func(item testItem) string { return item.Name })).
			Sort()

		if items[0].Name != "alpha" {
			t.Errorf("ByField string first = %s, want alpha", items[0].Name)
		}
	})

	t.Run("int field", func(t *testing.T) {
		t.Parallel()

		items := testItemsWithCounts(30, 10, 20)
		New(items, output.SortBy("Count"), false).
			WithLessFunc(ByField(func(item testItem) int { return item.Count })).
			Sort()

		if items[0].Name != "bravo" {
			t.Errorf("ByField int first = %s, want bravo", items[0].Name)
		}
	})

	t.Run("with desc", func(t *testing.T) {
		t.Parallel()

		items := testItemsWithCounts(10, 20, 30)
		New(items, output.SortBy("Count"), true).
			WithLessFunc(ByField(func(item testItem) int { return item.Count })).
			Sort()

		if items[0].Name != "charlie" {
			t.Errorf("ByField desc first = %s, want charlie", items[0].Name)
		}
	})
}
