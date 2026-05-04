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
		sortByNameField(items, func(item testItem) string { return item.Name })

		assertItemByName(t, items, "first", 0, "alpha")
	})

	t.Run("int field", func(t *testing.T) {
		t.Parallel()

		items := testItemsWithCounts(30, 10, 20)
		sortByCount(items, false)

		assertItemByName(t, items, "first", 0, "bravo")
	})

	t.Run("with desc", func(t *testing.T) {
		t.Parallel()

		items := testItemsWithCounts(10, 20, 30)
		New(items, output.SortBy("Count"), true).
			WithLessFunc(ByField(lessByCount)).
			Sort()

		assertItemByName(t, items, "first", 0, "charlie")
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

		assertFirstAndLast(t, items, "small", "large",
			func(item compareTestUnsignedItem) string { return item.Name })
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

		assertOrderSequence(t, items, 0, 1, 2)
	})
}
