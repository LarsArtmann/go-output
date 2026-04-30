package sort

import "cmp"

// ByField returns a LessFunc that extracts a comparable field from each item
// and compares using cmp.Less. This eliminates boilerplate for the common case
// of sorting by a single field.
//
// Usage:
//
//	sorter.WithLessFunc(sort.ByField(func(p Project) string { return p.Name }))
func ByField[T any, F cmp.Ordered](extract func(T) F) func(a, b T) bool {
	return func(a, b T) bool {
		return cmp.Less(extract(a), extract(b))
	}
}
