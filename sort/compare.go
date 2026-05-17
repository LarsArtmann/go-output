// Package sort provides sorting utilities for output data.
//
// Deprecated: Use the standard library instead. Go 1.21+ provides:
//
//	slices.SortStableFunc(items, func(a, b T) int {
//	    return cmp.Compare(a.Field, b.Field) // ascending
//	})
//
// For descending, swap the operands: cmp.Compare(b.Field, a.Field).
package sort

import "cmp"

// ByField returns a comparison function that extracts a comparable field from
// each item and returns cmp.Compare(extract(a), extract(b)). Use with
// slices.SortStableFunc:
//
//	slices.SortStableFunc(items, sort.ByField(func(p T) string { return p.Name }))
func ByField[T any, F cmp.Ordered](extract func(T) F) func(a, b T) int {
	return func(a, b T) int {
		return cmp.Compare(extract(a), extract(b))
	}
}
