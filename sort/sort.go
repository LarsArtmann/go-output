package sort

import (
	"reflect"
	"sort"
	"time"

	"github.com/larsartmann/go-output"
)

type Comparator func(a, b any) int

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

func toInt(v any) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int8:
		return int(val), true
	case int16:
		return int(val), true
	case int32:
		return int(val), true
	case int64:
		return int(val), true
	case uint:
		return int(val), true
	case uint8:
		return int(val), true
	case uint16:
		return int(val), true
	case uint32:
		return int(val), true
	case uint64:
		return int(val), true
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

type Sorter[T any] struct {
	Items    []T
	By       output.SortBy
	Desc     bool
	LessFunc func(a, b T) bool
}

func New[T any](items []T, by output.SortBy, desc bool) *Sorter[T] {
	return &Sorter[T]{
		Items: items,
		By:    by,
		Desc:  desc,
	}
}

func (s *Sorter[T]) WithLessFunc(fn func(a, b T) bool) *Sorter[T] {
	s.LessFunc = fn
	return s
}

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

	fieldName := string(s.By)
	field := aVal.FieldByName(fieldName)

	if !field.IsValid() {
		return false
	}

	switch field.Kind() {
	case reflect.String:
		return field.String() < bVal.FieldByName(fieldName).String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() < bVal.FieldByName(fieldName).Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() < bVal.FieldByName(fieldName).Uint()
	case reflect.Struct:
		if aTime, ok := field.Interface().(time.Time); ok {
			if bTime, ok := bVal.FieldByName(fieldName).Interface().(time.Time); ok {
				return aTime.Before(bTime)
			}
		}
	}

	return false
}
