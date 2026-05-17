package sort

import (
	"cmp"
	"slices"
	"testing"
)

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

		if items[0].Name != "alpha" || items[1].Name != "bravo" || items[2].Name != "charlie" {
			t.Errorf("sorted order = %v, want [alpha, bravo, charlie]",
				[]string{items[0].Name, items[1].Name, items[2].Name})
		}
	})

	t.Run("int field ascending", func(t *testing.T) {
		t.Parallel()

		type item struct {
			Name  string
			Count int
		}

		items := []item{
			{Name: "charlie", Count: 30},
			{Name: "alpha", Count: 10},
			{Name: "bravo", Count: 20},
		}

		slices.SortStableFunc(items, ByField(func(i item) int {
			return i.Count
		}))

		if items[0].Name != "alpha" || items[2].Name != "charlie" {
			t.Errorf("sorted order = %v, want [alpha, bravo, charlie]",
				[]string{items[0].Name, items[1].Name, items[2].Name})
		}
	})

	t.Run("uint64 field", func(t *testing.T) {
		t.Parallel()

		type item struct {
			Name string
			Size uint64
		}

		items := []item{
			{Name: "large", Size: 18_446_744_073_709_551_615},
			{Name: "small", Size: 1},
			{Name: "medium", Size: 100},
		}

		slices.SortStableFunc(items, ByField(func(i item) uint64 {
			return i.Size
		}))

		if items[0].Name != "small" || items[2].Name != "large" {
			t.Errorf("sorted order = %v, want [small, medium, large]",
				[]string{items[0].Name, items[1].Name, items[2].Name})
		}
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
