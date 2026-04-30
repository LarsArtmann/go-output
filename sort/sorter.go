// Package sort provides sorting utilities for output data.
package sort

import (
	"sort"

	output "github.com/larsartmann/go-output"
)

// Sorter sorts items by a specified field using a provided comparison function.
type Sorter[T any] struct {
	Items    []T
	By       output.SortBy
	Desc     bool
	LessFunc func(a, b T) bool
}

// New creates a new Sorter with the given items, sort field, and direction.
// A LessFunc must be provided via WithLessFunc before calling Sort.
func New[T any](items []T, by output.SortBy, desc bool) *Sorter[T] {
	return &Sorter[T]{
		Items:    items,
		By:       by,
		Desc:     desc,
		LessFunc: nil,
	}
}

// WithLessFunc sets the comparison function and returns the sorter for chaining.
// This must be called before Sort.
func (s *Sorter[T]) WithLessFunc(fn func(a, b T) bool) *Sorter[T] {
	s.LessFunc = fn

	return s
}

// Sort sorts the items using the provided LessFunc.
// If no LessFunc is set, Sort is a no-op (items remain in original order).
func (s *Sorter[T]) Sort() {
	if s.LessFunc == nil {
		return
	}

	sort.SliceStable(s.Items, func(i, j int) bool {
		a, b := s.Items[i], s.Items[j]
		if s.Desc {
			a, b = b, a
		}

		return s.LessFunc(a, b)
	})
}
