package graph

import (
	"github.com/larsartmann/go-output/enum"
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

//nolint:gochecknoglobals // Global variable used for value iteration.
var rankDirValues = []RankDir{
	RankDirTB,
	RankDirLR,
	RankDirBT,
	RankDirRL,
}

// InvalidRankDirError is returned when an invalid rank direction is provided.
type InvalidRankDirError struct {
	Value string
}

// Error returns a descriptive error message for the invalid rank direction.
func (e *InvalidRankDirError) Error() string {
	return "invalid rank direction: " + e.Value + " (allowed: TB, LR, BT, RL)"
}

// ParseRankDir converts a string to RankDir, returning an error if invalid.
func ParseRankDir(s string) (RankDir, error) {
	v, err := enum.Parse(rankDirValues, s, func(r RankDir) string { return string(r) })
	if err != nil {
		return "", &InvalidRankDirError{Value: s}
	}

	return v, nil
}

// String returns the string representation of the rank direction.
func (r RankDir) String() string {
	return string(r)
}

// AllowedValues returns all valid rank direction values.
func (RankDir) AllowedValues() []string {
	return enum.AllowedValues(rankDirValues)
}

// IsValid checks if the rank direction is valid.
func (r RankDir) IsValid() bool {
	return enum.Contains(rankDirValues, r)
}

// SplineStyle controls the edge routing style of a DOT graph.
type SplineStyle string

// SplineStyle constants define the valid edge routing styles for DOT graphs.
const (
	SplineOrtho     SplineStyle = "ortho"     // Orthogonal routing (default)
	SplineSpline    SplineStyle = "spline"    // Curved splines
	SplinePolyline  SplineStyle = "polyline"  // Straight line segments
	SplineLine      SplineStyle = "line"      // Straight lines
	SplineCurved    SplineStyle = "curved"    // Curved routing
	SplineNone      SplineStyle = "none"      // No edge routing
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var splineStyleValues = []SplineStyle{
	SplineOrtho,
	SplineSpline,
	SplinePolyline,
	SplineLine,
	SplineCurved,
	SplineNone,
}

// InvalidSplineStyleError is returned when an invalid spline style is provided.
type InvalidSplineStyleError struct {
	Value string
}

// Error returns a descriptive error message for the invalid spline style.
func (e *InvalidSplineStyleError) Error() string {
	return "invalid spline style: " + e.Value + " (allowed: ortho, spline, polyline, line, curved, none)"
}

// ParseSplineStyle converts a string to SplineStyle, returning an error if invalid.
func ParseSplineStyle(s string) (SplineStyle, error) {
	v, err := enum.Parse(splineStyleValues, s, func(s SplineStyle) string { return string(s) })
	if err != nil {
		return "", &InvalidSplineStyleError{Value: s}
	}

	return v, nil
}

// String returns the string representation of the spline style.
func (s SplineStyle) String() string {
	return string(s)
}

// AllowedValues returns all valid spline style values.
func (SplineStyle) AllowedValues() []string {
	return enum.AllowedValues(splineStyleValues)
}

// IsValid checks if the spline style is valid.
func (s SplineStyle) IsValid() bool {
	return enum.Contains(splineStyleValues, s)
}
