// Package sort_test provides tests for the sort package.
package sort

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/internal/gentest"
)

func compareTest[T any](
	t *testing.T,
	cmpName, testName string,
	cmp func(a, b T) int,
	a, b T,
	want int,
) {
	t.Run(testName, func(t *testing.T) {
		t.Parallel()

		if got := cmp(a, b); got != want {
			t.Errorf("%s(%v, %v) = %v, want %v", cmpName, a, b, got, want)
		}
	})
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
		compareTest(t, "CompareString", tt.name, CompareString, tt.a, tt.b, tt.want)
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
		compareTest(t, "CompareInt", tt.name, CompareInt, tt.a, tt.b, tt.want)
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
		compareTest(t, "CompareTime", tt.name, CompareTime, tt.a, tt.b, tt.want)
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
			assertOkBool(t, "toInt", tt.input, ok, tt.wantOk)
			gentest.AssertEqual(t, "toInt", tt.input, got, tt.want)
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
			assertOkBool(t, "toTime", tt.input, ok, tt.wantOk)
		})
	}
}

func assertOkBool(t *testing.T, name string, input any, ok, wantOk bool) {
	if ok != wantOk {
		t.Errorf("%s(%v) ok = %v, want %v", name, input, ok, wantOk)
	}
}
