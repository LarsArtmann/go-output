package sort

import (
	"testing"

	output "github.com/larsartmann/go-output"
)

type compareTestUnsignedItem struct {
	Name string
	Size uint64
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

	t.Run("uint64 field", func(t *testing.T) {
		t.Parallel()

		items := []compareTestUnsignedItem{
			{Name: "large", Size: 18_446_744_073_709_551_615},
			{Name: "small", Size: 1},
			{Name: "medium", Size: 100},
		}

		New(items, output.SortBy("Size"), false).
			WithLessFunc(ByField(func(item compareTestUnsignedItem) uint64 { return item.Size })).
			Sort()

		if items[0].Name != "small" {
			t.Errorf("ByField uint64 first = %s, want small", items[0].Name)
		}

		if items[2].Name != "large" {
			t.Errorf("ByField uint64 last = %s, want large", items[2].Name)
		}
	})

	t.Run("stability preserved", func(t *testing.T) {
		t.Parallel()

		items := []stableItem{
			{Name: "a", Order: 1},
			{Name: "a", Order: 2},
			{Name: "b", Order: 3},
		}

		New(items, output.SortByName, false).
			WithLessFunc(ByField(func(item stableItem) string { return item.Name })).
			Sort()

		if items[0].Order != 1 || items[1].Order != 2 {
			t.Errorf(
				"ByField stability: orders = [%d, %d], want [1, 2]",
				items[0].Order, items[1].Order,
			)
		}
	})
}
