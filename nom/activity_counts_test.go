package nom

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

func TestActivityCounts_Total(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   int
	}{
		{"zero", ActivityCounts{}, 0},
		{"only running", ActivityCounts{Running: 5}, 5},
		{"mixed", ActivityCounts{Running: 2, Completed: 3, Failed: 1, Pending: 4}, 10},
		{"all completed", ActivityCounts{Completed: 100}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.Total()
			if got != tt.want {
				t.Errorf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestActivityCounts_Summary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   string
	}{
		{"zero", ActivityCounts{}, ""},
		{"only running", ActivityCounts{Running: 3}, "⏵3"},
		{"only completed", ActivityCounts{Completed: 2}, "✔2"},
		{"running + completed", ActivityCounts{Running: 1, Completed: 2}, "⏵1 ✔2"},
		{"all four", ActivityCounts{Running: 1, Completed: 2, Failed: 3, Pending: 4}, "⏵1 ✔2 ⚠3 ○4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.Summary()
			if got != tt.want {
				t.Errorf("Summary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActivityCounts_CompletionPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   int
	}{
		{"zero", ActivityCounts{}, 0},
		{"all pending", ActivityCounts{Pending: 10}, 0},
		{"half done", ActivityCounts{Completed: 5, Pending: 5}, 50},
		{"all completed", ActivityCounts{Completed: 10}, 100},
		{"all failed", ActivityCounts{Failed: 10}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.CompletionPercent()
			if got != tt.want {
				t.Errorf("CompletionPercent() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestActivityCounts_SummaryColored(t *testing.T) {
	// NOT parallel: temporarily mutates global lipgloss.Writer.Profile.
	oldProfile := lipgloss.Writer.Profile
	lipgloss.Writer.Profile = colorprofile.ANSI

	t.Cleanup(func() { lipgloss.Writer.Profile = oldProfile })

	t.Run("zero counts produce empty string", func(t *testing.T) {
		got := ActivityCounts{}.SummaryColored(Colors)
		if got != "" {
			t.Errorf("SummaryColored() = %q, want empty", got)
		}
	})

	t.Run("all four with correct ANSI colors", func(t *testing.T) {
		counts := ActivityCounts{Running: 1, Completed: 2, Failed: 3, Pending: 4}
		got := counts.SummaryColored(Colors)

		checks := []struct {
			name    string
			ansiSgr string
			visible string
		}{
			{"running yellow", "\x1b[93m", "⏵1"},
			{"completed green", "\x1b[92m", "✔2"},
			{"failed red", "\x1b[91m", "⚠3"},
			{"pending gray", "\x1b[90m", "○4"},
		}

		for _, c := range checks {
			if !strings.Contains(got, c.ansiSgr+c.visible) {
				t.Errorf("%s: expected %q in output, got %q", c.name, c.ansiSgr+c.visible, got)
			}
		}

		if !strings.Contains(got, "\x1b[m") {
			t.Errorf("output should contain reset code, got %q", got)
		}
	})

	t.Run("zero counts omitted", func(t *testing.T) {
		got := ActivityCounts{Running: 0, Completed: 5}.SummaryColored(Colors)
		if strings.Contains(got, "⏵") {
			t.Errorf("running should be omitted when zero, got %q", got)
		}

		if !strings.Contains(got, "✔5") {
			t.Errorf("completed should be present, got %q", got)
		}
	})
}

func TestActivityCounts_SummaryColored_CustomTheme(t *testing.T) {
	// NOT parallel: temporarily mutates global lipgloss.Writer.Profile.
	oldProfile := lipgloss.Writer.Profile
	lipgloss.Writer.Profile = colorprofile.TrueColor

	t.Cleanup(func() { lipgloss.Writer.Profile = oldProfile })

	counts := ActivityCounts{Running: 1, Completed: 2}
	custom := SemanticColors{
		Running:   lipgloss.Color("#ff0000"),
		Completed: lipgloss.Color("#00ff00"),
		Pending:   lipgloss.Color("#808080"),
		Failed:    lipgloss.Color("#ff00ff"),
		Fallback:  lipgloss.Color("#00ffff"),
		Phase:     lipgloss.Color("#800080"),
	}
	got := counts.SummaryColored(custom)

	if !strings.Contains(got, "\x1b[38;2;255;0;0m") {
		t.Errorf("running should use #ff0000 (TrueColor), got %q", got)
	}

	if !strings.Contains(got, "\x1b[38;2;0;255;0m") {
		t.Errorf("completed should use #00ff00 (TrueColor), got %q", got)
	}
}
