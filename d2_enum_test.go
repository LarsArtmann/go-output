package output

import (
	"testing"
)

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
			got, err := ParseD2NodeShape(s)
			if err != nil {
				t.Errorf("ParseD2NodeShape(%q) unexpected error: %v", s, err)
			}

			if string(got) != s {
				t.Errorf("ParseD2NodeShape(%q) = %q, want %q", s, got, s)
			}
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
	if len(values) != 20 {
		t.Errorf("AllowedValues() returned %d values, want 20", len(values))
	}
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
			got, err := ParseD2ArrowType(s)
			if err != nil {
				t.Errorf("ParseD2ArrowType(%q) unexpected error: %v", s, err)
			}

			if string(got) != s {
				t.Errorf("ParseD2ArrowType(%q) = %q, want %q", s, got, s)
			}
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
	if len(values) != 11 {
		t.Errorf("AllowedValues() returned %d values, want 11", len(values))
	}
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
	if len(values) != 4 {
		t.Errorf("AllowedValues() returned %d values, want 4", len(values))
	}
}
