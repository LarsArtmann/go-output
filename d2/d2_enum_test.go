package d2

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

func TestParseNodeShape(t *testing.T) {
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
				ParseNodeShape,
				"ParseNodeShape",
				s,
				func(v NodeShape) string { return string(v) },
			)
		}
	})

	t.Run("invalid shape", func(t *testing.T) {
		t.Parallel()

		_, err := ParseNodeShape("invalid")
		if err == nil {
			t.Error("ParseNodeShape(invalid) should return error")
		}
	})
}

func TestD2NodeShapeValidation(t *testing.T) {
	t.Parallel()

	if !ShapeCircle.IsValid() {
		t.Error("ShapeCircle should be valid")
	}

	invalid := NodeShape("nonexistent")
	if invalid.IsValid() {
		t.Error("nonexistent shape should not be valid")
	}

	values := ShapeRectangle.AllowedValues()
	assertAllowedValueCount(t, values, 20)
}

func TestParseArrowType(t *testing.T) {
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
				ParseArrowType,
				"ParseArrowType",
				s,
				func(v ArrowType) string { return string(v) },
			)
		}
	})

	t.Run("invalid arrow", func(t *testing.T) {
		t.Parallel()

		_, err := ParseArrowType("invalid")
		if err == nil {
			t.Error("ParseArrowType(invalid) should return error")
		}
	})

	t.Run("empty arrow is valid", func(t *testing.T) {
		t.Parallel()

		a := ArrowType("")
		if !a.IsValid() {
			t.Error("empty arrow type should be valid (means no arrow)")
		}
	})

	t.Run("parse empty arrow returns none", func(t *testing.T) {
		t.Parallel()

		got, err := ParseArrowType("")
		if err != nil {
			t.Errorf("ParseArrowType(\"\") unexpected error: %v", err)
		}

		if got != ArrowNone {
			t.Errorf("ParseArrowType(\"\") = %q, want %q", got, ArrowNone)
		}
	})
}

func TestD2ArrowTypeValidation(t *testing.T) {
	t.Parallel()

	if !ArrowArrow.IsValid() {
		t.Error("ArrowArrow should be valid")
	}

	values := ArrowArrow.AllowedValues()
	assertAllowedValueCount(t, values, 12)
}

func TestParseDirection(t *testing.T) {
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
			got, err := ParseDirection(tc.input)
			if err != nil {
				t.Errorf("ParseDirection(%q) unexpected error: %v", tc.input, err)
			}

			if string(got) != tc.want {
				t.Errorf("ParseDirection(%q) = %q, want %q", tc.input, got, tc.want)
			}
		}
	})

	t.Run("invalid direction", func(t *testing.T) {
		t.Parallel()

		_, err := ParseDirection("diagonal")
		if err == nil {
			t.Error("ParseDirection(diagonal) should return error")
		}
	})
}

func TestD2DirectionValidation(t *testing.T) {
	t.Parallel()

	if !DirRight.IsValid() {
		t.Error("DirRight should be valid")
	}

	invalid := Direction("down-under")
	if invalid.IsValid() {
		t.Error("invalid direction should not be valid")
	}

	values := DirRight.AllowedValues()
	assertAllowedValueCount(t, values, 4)
}

func TestConstraint(t *testing.T) {
	t.Parallel()

	t.Run("all constraints", func(t *testing.T) {
		t.Parallel()

		all := AllConstraints()
		if len(all) != 3 {
			t.Errorf("AllConstraints() returned %d, want 3", len(all))
		}
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		if ConstraintPrimary.String() != "primary_key" {
			t.Errorf(
				"ConstraintPrimary.String() = %q, want %q",
				ConstraintPrimary.String(), "primary_key",
			)
		}
	})

	t.Run("allowed values", func(t *testing.T) {
		t.Parallel()

		values := ConstraintPrimary.AllowedValues()
		if len(values) != 3 {
			t.Errorf("AllowedValues() returned %d, want 3", len(values))
		}
	})

	t.Run("parse valid", func(t *testing.T) {
		t.Parallel()

		got, err := ParseConstraint("unique")
		if err != nil {
			t.Fatalf("ParseConstraint(unique) error: %v", err)
		}

		if got != ConstraintUnique {
			t.Errorf("ParseConstraint(unique) = %v, want %v", got, ConstraintUnique)
		}
	})

	t.Run("is valid", func(t *testing.T) {
		t.Parallel()

		if !ConstraintForeign.IsValid() {
			t.Error("ConstraintForeign should be valid")
		}

		invalid := Constraint("cascade")
		if invalid.IsValid() {
			t.Error("invalid constraint should not be valid")
		}
	})
}
