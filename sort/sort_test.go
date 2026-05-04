// Package sort_test provides tests for the sort package.
package sort

import (
	"testing"
	"time"

	output "github.com/larsartmann/go-output"
)

// sortByNameField sorts items by the Name field using ByField.
func sortByNameField[T any](items []T, getName func(T) string) {
	New(items, output.SortByName, false).
		WithLessFunc(ByField(getName)).
		Sort()
}

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

// assertOrderSequence checks that consecutive items have expected order values.
func assertOrderSequence(t *testing.T, items []stableItem, startIdx int, expected ...int) {
	for i, exp := range expected {
		idx := startIdx + i
		if idx >= len(items) {
			t.Errorf("assertOrderSequence: index %d out of bounds", idx)
			return
		}

		if items[idx].Order != exp {
			t.Errorf("items[%d].Order = %d, want %d", idx, items[idx].Order, exp)
		}
	}
}

// assertItemByName asserts the item at the given position matches.
func assertItemByName(t *testing.T, items []testItem, pos string, idx int, expected string) {
	l := len(items)
	if idx < 0 || idx >= l {
		t.Fatalf("assertItemByName: index %d out of bounds for length %d", idx, l)
		return
	}

	if items[idx].Name != expected {
		t.Errorf("%s item = %s, want %s", pos, items[idx].Name, expected)
	}
}

// assertLastItemByName asserts the last item's name matches.
func assertLastItemByName(t *testing.T, items []testItem, expected string) {
	assertItemByName(t, items, "last", len(items)-1, expected)
}

// assertFirstAndLast checks both first and last items in a single call.
func assertFirstAndLast[T any](
	t *testing.T,
	items []T,
	firstExpected string,
	lastExpected string,
	getVal func(T) string,
) {
	l := len(items)
	if l < 2 {
		t.Fatalf("assertFirstAndLast: slice length %d < 2", l)
		return
	}

	if got := getVal(items[0]); got != firstExpected {
		t.Errorf("first item = %s, want %s", got, firstExpected)
	}

	if got := getVal(items[l-1]); got != lastExpected {
		t.Errorf("last item = %s, want %s", got, lastExpected)
	}
}

// assertItemAt checks item at specific index.
func assertItemAt[T any](
	t *testing.T,
	items []T,
	pos string,
	idx int,
	getVal func(T) string,
	expected string,
) {
	if idx < 0 || idx >= len(items) {
		t.Errorf("assertItemAt: index %d out of bounds for slice length %d", idx, len(items))
		return
	}

	if got := getVal(items[idx]); got != expected {
		t.Errorf("%s item = %s, want %s", pos, got, expected)
	}
}

func sortByCount(items []testItem, desc bool) {
	New(items, output.SortBy("Count"), desc).
		WithLessFunc(ByField(lessByCount)).
		Sort()
}

func assertFirstItem[T comparable](t *testing.T, items []T, expected T) {
	if items[0] != expected {
		t.Errorf("first item = %v, want %v", items[0], expected)
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
	return ByField(lessByName)
}

func lessByName(item testItem) string         { return item.Name }
func lessByCount(item testItem) int           { return item.Count }
func lessByNameStable(item stableItem) string { return item.Name }

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
			lessFunc:       ByField(lessByCount),
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
		WithLessFunc(ByField(lessByNameStable)).
		Sort()

	if items[0].Name != "b" {
		t.Fatalf("desc sort first = %s, want b", items[0].Name)
	}

	// Equal-name items must preserve original insertion order
	assertOrderSequence(t, items, 1, 1, 2, 3)
}

func TestSorter_Sort_DescCount(t *testing.T) {
	t.Parallel()

	items := testItemsWithCounts(10, 20, 30)

	New(items, output.SortBy("Count"), true).WithLessFunc(compareCount(true)).Sort()

	assertItemByName(t, items, "first", 0, "charlie")
	assertLastItemByName(t, items, "alpha")
}

func TestSorter_Sort_NonStructInput(t *testing.T) {
	t.Parallel()

	items := []string{"b", "a", "c"}

	// No LessFunc — should be stable (no-op)
	New(items, output.SortByName, false).Sort()
	assertFirstItem(t, items, "b")

	// With LessFunc — should sort
	New(items, output.SortByName, false).WithLessFunc(func(a, b string) bool {
		return a < b
	}).Sort()
	assertFirstItem(t, items, "a")
}

func TestSorter_Sort_UsesSliceStable(t *testing.T) {
	t.Parallel()

	// Verify that Sort() delegates to sort.SliceStable by checking stability
	items := []testItem{
		{Name: "same", Count: 1, When: time.Time{}},
		{Name: "same", Count: 2, When: time.Time{}},
		{Name: "same", Count: 3, When: time.Time{}},
	}

	sortByNameField(items, lessByName)

	// All names are equal, so relative order must be preserved
	if items[0].Count != 1 || items[1].Count != 2 || items[2].Count != 3 {
		t.Errorf("stable sort not preserved: counts = [%d, %d, %d], want [1, 2, 3]",
			items[0].Count, items[1].Count, items[2].Count)
	}
}

func TestNew_WithLessFuncChaining(t *testing.T) {
	t.Parallel()

	items := testItemsWithCounts(30, 10, 20)

	sortByCount(items, false)

	if items[0].Name != "bravo" || items[0].Count != 10 {
		t.Errorf(
			"chained sort first = %s (count %d), want bravo (count 10)",
			items[0].Name, items[0].Count,
		)
	}
}
