package output

import (
	"testing"
)

// assertParseEqual verifies that parsing a valid value returns the expected result.
func assertParseEqual[T any](
	t *testing.T,
	fn func(string) (T, error),
	fnName, input string,
	toStr func(T) string,
) {
	t.Helper()

	got, err := fn(input)
	if err != nil {
		t.Errorf("%s(%q) unexpected error: %v", fnName, input, err)
	}

	if toStr(got) != input {
		t.Errorf("%s(%q) = %v, want %q", fnName, input, got, input)
	}
}

// assertAllowedValueCount verifies AllowedValues returns the expected number of values.
func assertAllowedValueCount(t *testing.T, values []string, want int) {
	t.Helper()

	if len(values) != want {
		t.Errorf("AllowedValues() returned %d values, want %d", len(values), want)
	}
}

func TestParseD2NodeShape(t *testing.T) {
	t.Parallel()

	t.Run("valid shapes", func(t *testing.T) {
		t.Parallel()

		shapes := []string{
			"rectangle", "circle", "diamond", "hexagon", "cloud",
			"cylinder", "person", "queue", "oval", "parallelogram",
			"triangle", "sql_table", "image", "code", "text", "class",
			"page", "step", "stored_data", "square",
		}

		for _, s := range shapes {
			assertParseEqual(
				t,
				ParseD2NodeShape,
				"ParseD2NodeShape",
				s,
				func(v D2NodeShape) string { return string(v) },
			)
		}
	})

	t.Run("invalid shape", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2NodeShape("invalid")
		if err == nil {
			t.Error("ParseD2NodeShape(invalid) should return error")
		}
	})
}

func TestD2NodeShapeValidation(t *testing.T) {
	t.Parallel()

	if !D2ShapeCircle.IsValid() {
		t.Error("D2ShapeCircle should be valid")
	}

	invalid := D2NodeShape("nonexistent")
	if invalid.IsValid() {
		t.Error("nonexistent shape should not be valid")
	}

	values := D2ShapeRectangle.AllowedValues()
	assertAllowedValueCount(t, values, 20)
}

func TestParseD2ArrowType(t *testing.T) {
	t.Parallel()

	t.Run("valid arrows", func(t *testing.T) {
		t.Parallel()

		arrows := []string{
			"arrow", "triangle", "diamond", "circle", "filled",
			"box", "cross", "cf-one", "cf-many",
			"cf-one-required", "cf-many-required",
		}

		for _, s := range arrows {
			assertParseEqual(
				t,
				ParseD2ArrowType,
				"ParseD2ArrowType",
				s,
				func(v D2ArrowType) string { return string(v) },
			)
		}
	})

	t.Run("invalid arrow", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2ArrowType("invalid")
		if err == nil {
			t.Error("ParseD2ArrowType(invalid) should return error")
		}
	})

	t.Run("empty arrow is valid", func(t *testing.T) {
		t.Parallel()

		a := D2ArrowType("")
		if !a.IsValid() {
			t.Error("empty arrow type should be valid (means no arrow)")
		}
	})
}

func TestD2ArrowTypeValidation(t *testing.T) {
	t.Parallel()

	if !D2ArrowArrow.IsValid() {
		t.Error("D2ArrowArrow should be valid")
	}

	values := D2ArrowArrow.AllowedValues()
	assertAllowedValueCount(t, values, 11)
}

func TestParseD2Direction(t *testing.T) {
	t.Parallel()

	t.Run("valid directions", func(t *testing.T) {
		t.Parallel()

		cases := []struct{ input, want string }{
			{"", ""},
			{"right", "right"},
			{"left", "left"},
			{"up", "up"},
		}

		for _, tc := range cases {
			got, err := ParseD2Direction(tc.input)
			if err != nil {
				t.Errorf("ParseD2Direction(%q) unexpected error: %v", tc.input, err)
			}

			if string(got) != tc.want {
				t.Errorf("ParseD2Direction(%q) = %q, want %q", tc.input, got, tc.want)
			}
		}
	})

	t.Run("invalid direction", func(t *testing.T) {
		t.Parallel()

		_, err := ParseD2Direction("diagonal")
		if err == nil {
			t.Error("ParseD2Direction(diagonal) should return error")
		}
	})
}

func TestD2DirectionValidation(t *testing.T) {
	t.Parallel()

	if !D2DirRight.IsValid() {
		t.Error("D2DirRight should be valid")
	}

	invalid := D2Direction("down-under")
	if invalid.IsValid() {
		t.Error("invalid direction should not be valid")
	}

	values := D2DirRight.AllowedValues()
	assertAllowedValueCount(t, values, 4)
}
