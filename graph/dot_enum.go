package graph

import (
	"strings"

	output "github.com/larsartmann/go-output"
)

// RankDir controls the layout direction of a DOT graph.
type RankDir string

// RankDir constants define the valid layout directions for DOT graphs.
const (
	RankDirTB RankDir = "TB" // Top to bottom (default)
	RankDirLR RankDir = "LR" // Left to right
	RankDirBT RankDir = "BT" // Bottom to top
	RankDirRL RankDir = "RL" // Right to left
)

// AllRankDirs contains all valid rank direction values.
//
//nolint:gochecknoglobals // Global variable used for value iteration.
var AllRankDirs = []RankDir{
	RankDirTB,
	RankDirLR,
	RankDirBT,
	RankDirRL,
}

// InvalidRankDirError is returned when an invalid rank direction is provided.
type InvalidRankDirError struct {
	Value   string
	Allowed []RankDir
}

// Error returns a descriptive error message for the invalid rank direction.
func (e *InvalidRankDirError) Error() string {
	if len(e.Allowed) == 0 {
		return "invalid rank direction: " + e.Value
	}

	return "invalid rank direction: " + e.Value + " (allowed: " + strings.Join(output.EnumAllowedValues(e.Allowed), ", ") + ")"
}

// ParseRankDir converts a string to RankDir, returning an error if invalid.
func ParseRankDir(s string) (RankDir, error) {
	v, err := output.ParseEnum(AllRankDirs, s, func(r RankDir) string { return string(r) })
	if err != nil {
		return "", &InvalidRankDirError{Value: s, Allowed: AllRankDirs}
	}

	return v, nil
}

// String returns the string representation of the rank direction.
func (r RankDir) String() string {
	return string(r)
}

// AllowedValues returns all valid rank direction values.
func (RankDir) AllowedValues() []string {
	return output.EnumAllowedValues(AllRankDirs)
}

// IsValid checks if the rank direction is valid.
func (r RankDir) IsValid() bool {
	return output.ContainsEnum(AllRankDirs, r)
}

// SplineStyle controls the edge routing style of a DOT graph.
type SplineStyle string

// SplineStyle constants define the valid edge routing styles for DOT graphs.
const (
	SplineOrtho    SplineStyle = "ortho"    // Orthogonal routing (default)
	SplineSpline   SplineStyle = "spline"   // Curved splines
	SplinePolyline SplineStyle = "polyline" // Straight line segments
	SplineLine     SplineStyle = "line"     // Straight lines
	SplineCurved   SplineStyle = "curved"   // Curved routing
	SplineNone     SplineStyle = "none"     // No edge routing
)

// AllSplineStyles contains all valid spline style values.
//
//nolint:gochecknoglobals // Global variable used for value iteration.
var AllSplineStyles = []SplineStyle{
	SplineOrtho,
	SplineSpline,
	SplinePolyline,
	SplineLine,
	SplineCurved,
	SplineNone,
}

// InvalidSplineStyleError is returned when an invalid spline style is provided.
type InvalidSplineStyleError struct {
	Value   string
	Allowed []SplineStyle
}

// Error returns a descriptive error message for the invalid spline style.
func (e *InvalidSplineStyleError) Error() string {
	if len(e.Allowed) == 0 {
		return "invalid spline style: " + e.Value
	}

	return "invalid spline style: " + e.Value + " (allowed: " + strings.Join(output.EnumAllowedValues(e.Allowed), ", ") + ")"
}

// ParseSplineStyle converts a string to SplineStyle, returning an error if invalid.
func ParseSplineStyle(s string) (SplineStyle, error) {
	v, err := output.ParseEnum(AllSplineStyles, s, func(s SplineStyle) string { return string(s) })
	if err != nil {
		return "", &InvalidSplineStyleError{Value: s, Allowed: AllSplineStyles}
	}

	return v, nil
}

// String returns the string representation of the spline style.
func (s SplineStyle) String() string {
	return string(s)
}

// AllowedValues returns all valid spline style values.
func (SplineStyle) AllowedValues() []string {
	return output.EnumAllowedValues(AllSplineStyles)
}

// IsValid checks if the spline style is valid.
func (s SplineStyle) IsValid() bool {
	return output.ContainsEnum(AllSplineStyles, s)
}
