// Package sort provides sorting utilities for output data.
package sort

import (
	"reflect"
	"sort"
	"time"

	output "github.com/larsartmann/go-output"
)

// Comparator compares two values, returning -1, 0, or 1.
type Comparator func(a, b any) int

// CompareString compares two string values.
func CompareString(a, b any) int {
	sa, ok := a.(string)
	if !ok {
		return 0
	}

	sb, ok := b.(string)
	if !ok {
		return 0
	}

	if sa < sb {
		return -1
	}

	if sa > sb {
		return 1
	}

	return 0
}

// CompareInt compares two integer values.
func CompareInt(a, b any) int {
	ia, ok := toInt(a)
	if !ok {
		return 0
	}

	ib, ok := toInt(b)
	if !ok {
		return 0
	}

	if ia < ib {
		return -1
	}

	if ia > ib {
		return 1
	}

	return 0
}

// CompareTime compares two time values.
func CompareTime(a, b any) int {
	ta, ok := toTime(a)
	if !ok {
		return 0
	}

	tb, ok := toTime(b)
	if !ok {
		return 0
	}

	if ta.Before(tb) {
		return -1
	}

	if ta.After(tb) {
		return 1
	}

	return 0
}

//nolint:cyclop
func toInt(v any) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true //nolint:gosec // G115: safe truncation for sorting purposes
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true //nolint:gosec // G115: safe truncation for sorting purposes
	default:
		return 0, false
	}
}

func toTime(v any) (time.Time, bool) {
	switch val := v.(type) {
	case time.Time:
		return val, true
	case *time.Time:
		if val != nil {
			return *val, true
		}

		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}

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

	// Convert snake_case to PascalCase for field lookup
	// e.g., "name" -> "Name", "created_at" -> "CreatedAt"
	fieldName := snakeToPascal(string(s.By))
	field := aVal.FieldByName(fieldName)

	if !field.IsValid() {
		return false
	}

	return compareFieldValues(field, bVal.FieldByName(fieldName))
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return a.Int() < b.Int()
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
