package graph

import "testing"

func TestParseRankDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  RankDir
	}{
		{"TB", RankDirTB},
		{"LR", RankDirLR},
		{"BT", RankDirBT},
		{"RL", RankDirRL},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseRankDir(tt.input)
			if err != nil {
				t.Fatalf("ParseRankDir(%q) error = %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("ParseRankDir(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseRankDirInvalid(t *testing.T) {
	t.Parallel()

	_, err := ParseRankDir("DIAGONAL")
	if err == nil {
		t.Fatal("expected error for invalid rank direction")
	}
}

func TestRankDirIsValid(t *testing.T) {
	t.Parallel()

	valid := []RankDir{RankDirTB, RankDirLR, RankDirBT, RankDirRL}
	for _, v := range valid {
		if !v.IsValid() {
			t.Errorf("%q should be valid", v)
		}
	}

	if RankDir("XYZ").IsValid() {
		t.Error("XYZ should not be valid")
	}
}

func TestRankDirAllowedValues(t *testing.T) {
	t.Parallel()

	values := RankDirTB.AllowedValues()
	if len(values) != 4 {
		t.Errorf("expected 4 rank directions, got %d", len(values))
	}
}

func TestParseSplineStyle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  SplineStyle
	}{
		{"ortho", SplineOrtho},
		{"spline", SplineSpline},
		{"polyline", SplinePolyline},
		{"line", SplineLine},
		{"curved", SplineCurved},
		{"none", SplineNone},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseSplineStyle(tt.input)
			if err != nil {
				t.Fatalf("ParseSplineStyle(%q) error = %v", tt.input, err)
			}

			if got != tt.want {
				t.Errorf("ParseSplineStyle(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseSplineStyleInvalid(t *testing.T) {
	t.Parallel()

	_, err := ParseSplineStyle("zigzag")
	if err == nil {
		t.Fatal("expected error for invalid spline style")
	}
}

func TestSplineStyleIsValid(t *testing.T) {
	t.Parallel()

	valid := []SplineStyle{
		SplineOrtho, SplineSpline, SplinePolyline,
		SplineLine, SplineCurved, SplineNone,
	}
	for _, v := range valid {
		if !v.IsValid() {
			t.Errorf("%q should be valid", v)
		}
	}

	if SplineStyle("zigzag").IsValid() {
		t.Error("zigzag should not be valid")
	}
}
