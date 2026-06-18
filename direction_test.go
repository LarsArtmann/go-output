package output

import (
	"testing"
)

func TestDirection_ToD2Direction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dir  Direction
		want string
	}{
		{"down is empty (D2 default)", DirectionDown, ""},
		{"up", DirectionUp, "up"},
		{"left", DirectionLeft, "left"},
		{"right", DirectionRight, "right"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.dir.ToD2Direction()
			if got != tt.want {
				t.Errorf("ToD2Direction() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirection_ToRankDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dir  Direction
		want string
	}{
		{"down maps to TB", DirectionDown, "TB"},
		{"up maps to BT", DirectionUp, "BT"},
		{"left maps to RL", DirectionLeft, "RL"},
		{"right maps to LR", DirectionRight, "LR"},
		{"unknown defaults to TB", Direction("diagonal"), "TB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.dir.ToRankDir()
			if got != tt.want {
				t.Errorf("ToRankDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDirection_RoundTrip(t *testing.T) {
	t.Parallel()

	// Verify all 4 canonical directions produce distinct DOT rankdir values.
	seen := make(map[string]bool)

	for _, d := range []Direction{DirectionDown, DirectionUp, DirectionLeft, DirectionRight} {
		rd := d.ToRankDir()
		if seen[rd] {
			t.Errorf("Direction %s produced duplicate RankDir %q", d, rd)
		}

		seen[rd] = true
	}
}
