// Package sort provides sorting utilities for output data.
package sort

import (
	"fmt"
	"reflect"

	"github.com/larsartmann/go-output"
)

// Adapter provides adapter functionality for integrating sort with cmdguard.
type Adapter struct {
	SortBy output.SortBy
	Desc   bool
}

// NewAdapter creates a new Adapter with the specified sort field.
func NewAdapter(by output.SortBy) *Adapter {
	return &Adapter{
		SortBy: by,
		Desc:   false,
	}
}

// WithDescending sets the sort order to descending.
func (a *Adapter) WithDescending(desc bool) *Adapter {
	a.Desc = desc
	return a
}

// Comparator returns a sort function that can be used with the sort package.
func (a *Adapter) Comparator() func(a, b any) int {
	return func(av, bv any) int {
		return compareValues(av, bv, a.SortBy, a.Desc)
	}
}

// ParseAdapter creates an Adapter from a string value.
func ParseAdapter(s string) (*Adapter, error) {
	sb, err := output.ParseSortBy(s)
	if err != nil {
		return nil, fmt.Errorf("parse sort by: %w", err)
	}
	return NewAdapter(sb), nil
}

// AllowedValues returns all allowed sort values for CLI help text.
func AllowedValues() []string {
	// Return values from the output package
	values := make([]string, 0, 6)
	for _, v := range []string{"name", "importance", "created_at", "updated_at", "health", "complexity"} {
		values = append(values, v)
	}
	return values
}

// compareValues compares two values based on the sort field.
func compareValues(a, b any, by output.SortBy, desc bool) int {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	// Get field name
	fieldName := snakeToPascal(string(by))

	aField := aVal.FieldByName(fieldName)
	bField := bVal.FieldByName(fieldName)

	if !aField.IsValid() || !bField.IsValid() {
		return 0
	}

	result := compareReflectValues(aField, bField)
	if desc {
		return -result
	}
	return result
}

// compareReflectValues compares two reflect.Value instances.
func compareReflectValues(a, b reflect.Value) int {
	switch a.Kind() {
	case reflect.String:
		sa := a.String()
		sb := b.String()
		if sa < sb {
			return -1
		}
		if sa > sb {
			return 1
		}
		return 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ia := a.Int()
		ib := b.Int()
		if ia < ib {
			return -1
		}
		if ia > ib {
			return 1
		}
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ua := a.Uint()
		ub := b.Uint()
		if ua < ub {
			return -1
		}
		if ua > ub {
			return 1
		}
		return 0
	case reflect.Float32, reflect.Float64:
		fa := a.Float()
		fb := b.Float()
		if fa < fb {
			return -1
		}
		if fa > fb {
			return 1
		}
		return 0
	default:
		return 0
	}
}
