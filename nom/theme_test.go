package nom

import (
	"context"
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestThemeDefault_MatchesGlobalColors(t *testing.T) {
	t.Parallel()

	if ThemeDefault.Colors.Running != Colors.Running {
		t.Errorf("ThemeDefault.Colors.Running = %v, want %v", ThemeDefault.Colors.Running, Colors.Running)
	}

	if ThemeDefault.Colors.Completed != Colors.Completed {
		t.Errorf("ThemeDefault.Colors.Completed = %v, want %v", ThemeDefault.Colors.Completed, Colors.Completed)
	}

	if ThemeDefault.Colors.Pending != Colors.Pending {
		t.Errorf("ThemeDefault.Colors.Pending = %v, want %v", ThemeDefault.Colors.Pending, Colors.Pending)
	}

	if ThemeDefault.Colors.Failed != Colors.Failed {
		t.Errorf("ThemeDefault.Colors.Failed = %v, want %v", ThemeDefault.Colors.Failed, Colors.Failed)
	}

	if ThemeDefault.Colors.Fallback != Colors.Fallback {
		t.Errorf("ThemeDefault.Colors.Fallback = %v, want %v", ThemeDefault.Colors.Fallback, Colors.Fallback)
	}

	if ThemeDefault.Colors.Phase != Colors.Phase {
		t.Errorf("ThemeDefault.Colors.Phase = %v, want %v", ThemeDefault.Colors.Phase, Colors.Phase)
	}
}

func TestTheme_StatusColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status ActivityStatus
		want   color.Color
	}{
		{ActivityStatusPending, ThemeDefault.Colors.Pending},
		{ActivityStatusRunning, ThemeDefault.Colors.Running},
		{ActivityStatusCompleted, ThemeDefault.Colors.Completed},
		{ActivityStatusFailed, ThemeDefault.Colors.Failed},
		{ActivityStatus(99), ThemeDefault.Colors.Fallback},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			t.Parallel()

			got := ThemeDefault.StatusColor(tt.status)
			if got != tt.want {
				t.Errorf("StatusColor(%v) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestTheme_ColorFor(t *testing.T) {
	t.Parallel()

	if got := ThemeDefault.ColorFor(ActivityKindTask, ActivityStatusRunning); got != ThemeDefault.Colors.Running {
		t.Errorf("ColorFor(Task, Running) = %v, want %v", got, ThemeDefault.Colors.Running)
	}

	if got := ThemeDefault.ColorFor(ActivityKindPhase, ActivityStatusPending); got != ThemeDefault.Colors.Phase {
		t.Errorf("ColorFor(Phase, Pending) = %v, want %v", got, ThemeDefault.Colors.Phase)
	}
}

func TestTheme_SymbolFor_Override(t *testing.T) {
	t.Parallel()

	customTheme := Theme{
		Colors: ThemeDefault.Colors,
		Symbols: map[ActivityStatus]Symbol{
			ActivityStatusRunning: "►",
		},
	}

	if got := customTheme.SymbolFor(ActivityStatusRunning, SymbolRunning); got != "►" {
		t.Errorf("SymbolFor(Running) = %q, want %q", got, "►")
	}

	if got := customTheme.SymbolFor(ActivityStatusPending, SymbolPending); got != SymbolPending {
		t.Errorf("SymbolFor(Pending) = %q, want %q", got, SymbolPending)
	}
}

func TestWithTheme_SetsSubscriberTheme(t *testing.T) {
	t.Parallel()

	custom := Theme{
		Colors: SemanticColors{
			Running: lipgloss.Color("#ff0000"),
		},
		Symbols: map[ActivityStatus]Symbol{
			ActivityStatusRunning: "R",
		},
	}

	ns := NewNOMSubscriber(WithTheme(custom))

	testhelpers.AssertEqual(t, "theme colors.Running", custom, ns.Theme().Colors.Running, custom.Colors.Running)
	testhelpers.AssertEqual(t, "theme symbol", custom, ns.Theme().Symbols[ActivityStatusRunning], Symbol("R"))
}

func TestSnapshotActivities_UsesThemeColor(t *testing.T) {
	t.Parallel()

	custom := Theme{
		Colors: SemanticColors{
			Running:   lipgloss.Color("#123456"),
			Completed: lipgloss.Color("#abcdef"),
			Pending:   lipgloss.Color("#111111"),
			Failed:    lipgloss.Color("#222222"),
			Fallback:  lipgloss.Color("#333333"),
			Phase:     lipgloss.Color("#444444"),
		},
		Symbols: map[ActivityStatus]Symbol{
			ActivityStatusRunning: "R",
		},
	}

	ns := NewNOMSubscriber(WithTheme(custom))
	ctx := context.Background()

	if err := ns.OnEvent(ctx, ActivityRegistered{ID: "build", Name: "Build"}); err != nil {
		t.Fatalf("ActivityRegistered error: %v", err)
	}

	if err := ns.OnEvent(ctx, ActivityStarted{ID: "build", Name: "Build"}); err != nil {
		t.Fatalf("ActivityStarted error: %v", err)
	}

	snap := ns.SnapshotActivities()[ActivityID("build")]
	if snap.Color != custom.Colors.Running {
		t.Errorf("snapshot color = %v, want %v", snap.Color, custom.Colors.Running)
	}

	if snap.Symbol != "R" {
		t.Errorf("snapshot symbol = %q, want %q", snap.Symbol, "R")
	}
}

func TestSnapshotActivities_UsesThemePhaseColor(t *testing.T) {
	t.Parallel()

	custom := Theme{
		Colors: SemanticColors{
			Running:   lipgloss.Color("#123456"),
			Completed: lipgloss.Color("#abcdef"),
			Pending:   lipgloss.Color("#111111"),
			Failed:    lipgloss.Color("#222222"),
			Fallback:  lipgloss.Color("#333333"),
			Phase:     lipgloss.Color("#444444"),
		},
		Symbols: map[ActivityStatus]Symbol{},
	}

	ns := NewNOMSubscriber(WithTheme(custom))
	ctx := context.Background()

	if err := ns.OnEvent(ctx, ActivityRegistered{ID: "phase", Name: "Phase", Kind: ActivityKindPhase}); err != nil {
		t.Fatalf("ActivityRegistered error: %v", err)
	}

	if err := ns.OnEvent(ctx, ActivityStarted{ID: "phase", Name: "Phase"}); err != nil {
		t.Fatalf("ActivityStarted error: %v", err)
	}

	snap := ns.SnapshotActivities()[ActivityID("phase")]
	if snap.Color != custom.Colors.Phase {
		t.Errorf("phase snapshot color = %v, want %v", snap.Color, custom.Colors.Phase)
	}
}

// TestThemePresets_NordAndMonochrome_NonNil verifies that the Nord and
// Monochrome theme presets have non-nil color values for all semantic slots.
// These themes had zero callers and zero tests before this.
func TestThemePresets_NordAndMonochrome_NonNil(t *testing.T) {
	t.Parallel()

	themes := []struct {
		name  string
		theme Theme
	}{
		{"ThemeNord", ThemeNord},
		{"ThemeMonochrome", ThemeMonochrome},
	}

	for _, tc := range themes {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			colors := []struct {
				label string
				c     color.Color
			}{
				{"Running", tc.theme.Colors.Running},
				{"Completed", tc.theme.Colors.Completed},
				{"Pending", tc.theme.Colors.Pending},
				{"Failed", tc.theme.Colors.Failed},
				{"Info", tc.theme.Colors.Fallback},
				{"Phase", tc.theme.Colors.Phase},
			}

			for _, c := range colors {
				if c.c == nil {
					t.Errorf("%s.Colors.%s is nil", tc.name, c.label)
				}
			}

			// Verify the theme can resolve a status color without panicking.
			got := tc.theme.StatusColor(ActivityStatusRunning)
			if got == nil {
				t.Errorf("%s.StatusColor(Running) returned nil", tc.name)
			}
		})
	}
}
