package output

import "slices"

// Direction represents a canonical layout direction for diagrams.
// It bridges D2 vocabulary ("down"/"right") and DOT vocabulary ("TB"/"LR")
// through a single canonical type.
type Direction string

const (
	DirectionDown  Direction = "down"
	DirectionUp    Direction = "up"
	DirectionLeft  Direction = "left"
	DirectionRight Direction = "right"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var AllDirections = []Direction{DirectionDown, DirectionUp, DirectionLeft, DirectionRight}

// IsValid reports whether d is a recognized Direction.
func (d Direction) IsValid() bool {
	return slices.Contains(AllDirections, d)
}

// ToD2Direction converts Direction to D2's direction string.
// D2 uses "" for default (down), "right", "left", "up".
func (d Direction) ToD2Direction() string {
	if d == DirectionDown {
		return "" // D2 default is down (empty string)
	}

	return string(d)
}

// ToRankDir converts Direction to DOT's rankdir string.
// DOT uses "TB" (top-to-bottom), "LR" (left-to-right), "BT", "RL".
func (d Direction) ToRankDir() string {
	switch d {
	case DirectionDown:
		return "TB"
	case DirectionUp:
		return "BT"
	case DirectionLeft:
		return "RL"
	case DirectionRight:
		return "LR"
	default:
		return "TB"
	}
}
