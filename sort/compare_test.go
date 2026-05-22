package sort

import (
	"cmp"
	"slices"
	"testing"
)

func extractNames[T any](items []T, name func(T) string) []string {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = name(item)
	}
	return names
}

func assertSortedOrder(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("assertSortedOrder: len(got)=%d != len(want)=%d", len(got), len(want))
	}

	for i := range got {
		if got[i] != want[i] {
			t.Errorf("sorted order = %v, want %v", got, want)

			return
		}
	}
}

func TestByField(t *testing.T) {
	t.Parallel()

	t.Run("string field ascending", func(t *testing.T) {
		t.Parallel()

		items := []struct{ Name string }{
			{Name: "charlie"},
			{Name: "alpha"},
			{Name: "bravo"},
		}

		slices.SortStableFunc(items, ByField(
			func(item struct{ Name string }) string { return item.Name },
		))

		assertSortedOrder(t, extractNames(items, func(i struct{ Name string }) string { return i.Name }),
			[]string{"alpha", "bravo", "charlie"})
	})

	t.Run("numeric field ascending", func(t *testing.T) {
		t.Parallel()

		t.Run("int field", func(t *testing.T) {
			t.Parallel()

			type orderedItem struct {
				Name  string
				Order int
			}

			items := []orderedItem{
				{Name: "charlie", Order: 30},
				{Name: "alpha", Order: 10},
				{Name: "bravo", Order: 20},
			}

			slices.SortStableFunc(items, ByField(func(i orderedItem) int { return i.Order }))
			assertSortedOrder(t, extractNames(items, func(i orderedItem) string { return i.Name }), []string{"alpha", "bravo", "charlie"})
		})

		t.Run("uint64 field", func(t *testing.T) {
			t.Parallel()

			type sizedItem struct {
				Name string
				Size uint64
			}

			items := []sizedItem{
				{Name: "large", Size: 18_446_744_073_709_551_615},
				{Name: "small", Size: 1},
				{Name: "medium", Size: 100},
			}

			slices.SortStableFunc(items, ByField(func(i sizedItem) uint64 { return i.Size }))
			assertSortedOrder(t, extractNames(items, func(i sizedItem) string { return i.Name }), []string{"small", "medium", "large"})
		})
	})

	t.Run("stability preserved", func(t *testing.T) {
		t.Parallel()

		type item struct {
			Name  string
			Order int
		}

		items := []item{
			{Name: "a", Order: 1},
			{Name: "a", Order: 2},
			{Name: "b", Order: 3},
		}

		slices.SortStableFunc(items, ByField(func(i item) string { return i.Name }))

		if items[0].Order != 1 || items[1].Order != 2 {
			t.Errorf("stability not preserved: orders = [%d, %d], want [1, 2]",
				items[0].Order, items[1].Order)
		}
	})
}

func TestByFieldMatchesCompare(t *testing.T) {
	t.Parallel()

	type p struct{ Name string }

	byName := ByField(func(i p) string { return i.Name })

	if byName(p{"a"}, p{"b"}) != cmp.Compare("a", "b") {
		t.Error("ByField should match cmp.Compare")
	}

	byNameDesc := func(a, b p) int { return cmp.Compare(b.Name, a.Name) }

	if byNameDesc(p{"a"}, p{"b"}) != 1 {
		t.Error("descending comparison should return positive")
	}
}
