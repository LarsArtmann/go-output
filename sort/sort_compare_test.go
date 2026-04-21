// Package sort_test provides tests for the sort package.
package sort

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output/internal/gentest"
)

func TestCompareString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a < b", "apple", "banana", -1},
		{"a > b", "banana", "apple", 1},
		{"a == b", "same", "same", 0},
		{"a not string", 123, "banana", 0},
		{"b not string", "apple", 123, 0},
		{"both not string", 123, 456, 0},
		{"empty strings", "", "", 0},
		{"empty vs non-empty", "", "a", -1},
	}

	for _, tt := range tests {
		compareTest(t, "CompareString", tt.name, CompareString, tt.a, tt.b, tt.want)
	}
}

func TestCompareInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a < b", 1, 2, -1},
		{"a > b", 2, 1, 1},
		{"a == b", 5, 5, 0},
		{"int8", int8(1), int8(2), -1},
		{"int16", int16(1), int16(2), -1},
		{"int32", int32(1), int32(2), -1},
		{"int64", int64(1), int64(2), -1},
		{"uint", uint(1), uint(2), -1},
		{"uint8", uint8(1), uint8(2), -1},
		{"uint16", uint16(1), uint16(2), -1},
		{"uint32", uint32(1), uint32(2), -1},
		{"uint64", uint64(1), uint64(2), -1},
		{"mixed types", int(1), int64(2), -1},
		{"a not int", "string", 2, 0},
		{"b not int", 1, "string", 0},
		{"both not int", "a", "b", 0},
		{"negative values", -5, 5, -1},
		{"zero values", 0, 0, 0},
	}

	for _, tt := range tests {
		compareTest(t, "CompareInt", tt.name, CompareInt, tt.a, tt.b, tt.want)
	}
}

func TestCompareTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	earlier := now.Add(-time.Hour)
	later := now.Add(time.Hour)

	tests := []struct {
		name string
		a, b any
		want int
	}{
		{"a before b", earlier, now, -1},
		{"a after b", later, now, 1},
		{"a equal b", now, now, 0},
		{"pointer to time", &earlier, now, -1},
		{"both pointers", &earlier, &later, -1},
		{"nil pointer", nil, now, 0},
		{"a not time", "string", now, 0},
		{"b not time", now, "string", 0},
		{"both not time", "a", "b", 0},
	}

	for _, tt := range tests {
		compareTest(t, "CompareTime", tt.name, CompareTime, tt.a, tt.b, tt.want)
	}
}

func TestToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  any
		want   int
		wantOk bool
	}{
		{"int", int(42), 42, true},
		{"int8", int8(42), 42, true},
		{"int16", int16(42), 42, true},
		{"int32", int32(42), 42, true},
		{"int64", int64(42), 42, true},
		{"uint", uint(42), 42, true},
		{"uint8", uint8(42), 42, true},
		{"uint16", uint16(42), 42, true},
		{"uint32", uint32(42), 42, true},
		{"uint64", uint64(42), 42, true},
		{"string", "42", 0, false},
		{"float", 3.14, 0, false},
		{"nil", nil, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := toInt(tt.input)
			assertOkBool(t, "toInt", tt.input, ok, tt.wantOk)
			gentest.AssertEqual(t, "toInt", tt.input, got, tt.want)
		})
	}
}

func TestToTime(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		input  any
		wantOk bool
	}{
		{"time.Time", now, true},
		{"*time.Time", &now, true},
		{"nil pointer", (*time.Time)(nil), false},
		{"string", "2023-01-01", false},
		{"int", 123, false},
		{"nil", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, ok := toTime(tt.input)
			assertOkBool(t, "toTime", tt.input, ok, tt.wantOk)
		})
	}
}

func assertOkBool(t *testing.T, name string, input any, ok, wantOk bool) {
	if ok != wantOk {
		t.Errorf("%s(%v) ok = %v, want %v", name, input, ok, wantOk)
	}
}
