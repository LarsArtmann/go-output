// Package sort provides sorting utilities for output data.
package sort

import (
	"reflect"
	"sort"
	"time"

	output "github.com/larsartmann/go-output"
)

// Sorter sorts items by a specified field.
type Sorter[T any] struct {
	Items    []T
	By       output.SortBy
	Desc     bool
	LessFunc func(a, b T) bool
}

// New creates a new Sorter.
func New[T any](items []T, by output.SortBy, desc bool) *Sorter[T] {
	return &Sorter[T]{
		Items:    items,
		By:       by,
		Desc:     desc,
		LessFunc: nil,
	}
}

// WithLessFunc sets a custom less function.
func (s *Sorter[T]) WithLessFunc(fn func(a, b T) bool) *Sorter[T] {
	s.LessFunc = fn

	return s
}

// Sort sorts the items.
func (s *Sorter[T]) Sort() {
	sort.SliceStable(s.Items, func(i, j int) bool {
		var result bool
		if s.LessFunc != nil {
			result = s.LessFunc(s.Items[i], s.Items[j])
		} else {
			result = s.defaultLess(s.Items[i], s.Items[j])
		}

		if s.Desc {
			return !result
		}

		return result
	})
}

func (s *Sorter[T]) defaultLess(a, b T) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Struct || bVal.Kind() != reflect.Struct {
		return false
	}

	fieldName := snakeToPascal(string(s.By))
	fieldA := aVal.FieldByName(fieldName)
	fieldB := bVal.FieldByName(fieldName)

	if !fieldA.IsValid() || !fieldB.IsValid() {
		return false
	}

	return compareFieldValues(fieldA, fieldB)
}

// snakeToPascal converts snake_case to PascalCase.
func snakeToPascal(s string) string {
	if s == "" {
		return ""
	}

	result := make([]byte, 0, len(s))
	upper := true

	for i := range len(s) {
		c := s[i]
		if c == '_' {
			upper = true

			continue
		}

		if upper {
			if c >= 'a' && c <= 'z' {
				c = c - 'a' + 'A'
			}

			upper = false
		}

		result = append(result, c)
	}

	return string(result)
}

func compareFieldValues(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.String:
		return a.String() < b.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Uint() < b.Uint()
	case reflect.Struct:
		if aTime, ok := a.Interface().(time.Time); ok {
			if bTime, ok := b.Interface().(time.Time); ok {
				return aTime.Before(bTime)
			}
		}
	case reflect.Invalid, reflect.Bool, reflect.Uintptr, reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer,
		reflect.Slice, reflect.UnsafePointer:
		return false
	}

	return false
}
