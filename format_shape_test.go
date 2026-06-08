package output

import (
	"strconv"
	"sync"
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestAllShapes(t *testing.T) {
	t.Parallel()

	want := []Shape{ShapeTable, ShapeTree, ShapeGraph}
	if len(AllShapes) != len(want) {
		t.Errorf("AllShapes length = %d, want %d", len(AllShapes), len(want))
	}

	for i, s := range AllShapes {
		if s != want[i] {
			t.Errorf("AllShapes[%d] = %v, want %v", i, s, want[i])
		}
	}
}

func TestParseShape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  Shape
		err   bool
	}{
		{"table", ShapeTable, false},
		{"tree", ShapeTree, false},
		{"graph", ShapeGraph, false},
		{"invalid", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseShape(tt.input)
			if tt.err {
				if err == nil {
					t.Errorf("ParseShape(%q) expected error, got nil", tt.input)
				}

				return
			}

			if err != nil {
				t.Errorf("ParseShape(%q) unexpected error: %v", tt.input, err)

				return
			}

			if got != tt.want {
				t.Errorf("ParseShape(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestShapeIsValid(t *testing.T) {
	t.Parallel()

	testhelpers.TestEnumIsValid(t, []Shape{
		ShapeTable,
		ShapeTree,
		ShapeGraph,
	}, []bool{true, true, true})
}

func TestShapeAllowedValues(t *testing.T) {
	t.Parallel()

	got := ShapeTable.AllowedValues()
	want := []string{"table", "tree", "graph"}

	if len(got) != len(want) {
		t.Errorf("Shape.AllowedValues() length = %d, want %d", len(got), len(want))

		return
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("Shape.AllowedValues()[%d] = %q, want %q", i, v, want[i])
		}
	}
}

func TestShapeString(t *testing.T) {
	t.Parallel()

	if ShapeTable.String() != "table" {
		t.Errorf("ShapeTable.String() = %q, want %q", ShapeTable.String(), "table")
	}

	if ShapeTree.String() != "tree" {
		t.Errorf("ShapeTree.String() = %q, want %q", ShapeTree.String(), "tree")
	}

	if ShapeGraph.String() != "graph" {
		t.Errorf("ShapeGraph.String() = %q, want %q", ShapeGraph.String(), "graph")
	}
}

func TestRegisterFormatShapes_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func() {
			defer wg.Done()

			RegisterFormatShapes(Format("race-test-"+strconv.Itoa(i)), ShapeTable)
		}()

		go func() {
			defer wg.Done()

			_, _ = getFormatShapes(Format("race-test-" + strconv.Itoa(i)))
		}()
	}

	wg.Wait()
}
