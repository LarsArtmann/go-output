package nom

import (
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestCategory_TagPrefix_TreeMode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("build"), nil)
	_ = dt.AddActivity(ActivityID("test"), []ActivityID{"build"})
	dt.showCategory = true

	snaps := newSnapshotBuilder()
	snaps.setCategory(ActivityID("build"), "Build", ActivityStatusCompleted, ActivityCategory("build"))
	snaps.setCategory(ActivityID("test"), "Test", ActivityStatusRunning, ActivityCategory("test"))

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	if !strings.Contains(got, "[build]") {
		t.Errorf("expected [build] category tag in output:\n%s", got)
	}

	if !strings.Contains(got, "[test]") {
		t.Errorf("expected [test] category tag in output:\n%s", got)
	}
}

func TestCategory_TagPrefix_DisabledByDefault(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("build"), nil)

	snaps := newSnapshotBuilder()
	snaps.setCategory(ActivityID("build"), "Build", ActivityStatusCompleted, ActivityCategory("build"))

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	if strings.Contains(got, "[build]") {
		t.Errorf("category tag should not render when showCategory is off:\n%s", got)
	}
}

func TestCategory_TagPrefix_LayeredMode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("build"), nil)
	dt.SetRenderMode(RenderModeLayered)
	dt.showCategory = true

	snaps := newSnapshotBuilder()
	snaps.setCategory(ActivityID("build"), "Build", ActivityStatusRunning, ActivityCategory("compile"))

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	if !strings.Contains(got, "[compile]") {
		t.Errorf("expected [compile] category tag in layered output:\n%s", got)
	}
}

func TestCategory_SetViaEvent(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	ctx := t.Context()

	_ = sub.OnEvent(ctx, ActivityRegistered{
		ID:       ActivityID("build"),
		Name:     ActivityName("Build"),
		Category: ActivityCategory("build"),
	})

	act := sub.GetActivity(ActivityID("build"))
	if act == nil {
		t.Fatal("activity not found")
	}

	if act.Category != ActivityCategory("build") {
		t.Errorf("category = %q, want %q", act.Category, "build")
	}
}

func TestCategory_SetViaStarted(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	ctx := t.Context()

	_ = sub.OnEvent(ctx, ActivityStarted{
		ID:       ActivityID("deploy"),
		Name:     ActivityName("Deploy"),
		Category: ActivityCategory("deploy"),
	})

	act := sub.GetActivity(ActivityID("deploy"))
	if act == nil {
		t.Fatal("activity not found")
	}

	if act.Category != ActivityCategory("deploy") {
		t.Errorf("category = %q, want %q", act.Category, "deploy")
	}
}

func TestCategory_SetActivityCategory(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	ctx := t.Context()

	_ = sub.OnEvent(ctx, ActivityRegistered{ID: ActivityID("a"), Name: ActivityName("A")})

	sub.SetActivityCategory(ActivityID("a"), ActivityCategory("test"))

	act := sub.GetActivity(ActivityID("a"))
	if act.Category != ActivityCategory("test") {
		t.Errorf("category = %q, want %q", act.Category, "test")
	}
}

func TestCategory_SnapshotPropagation(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	ctx := t.Context()

	_ = sub.OnEvent(ctx, ActivityRegistered{
		ID:       ActivityID("build"),
		Name:     ActivityName("Build"),
		Category: ActivityCategory("ci"),
	})

	snaps := sub.SnapshotActivities()
	snap := snaps[ActivityID("build")]

	if snap.Category != ActivityCategory("ci") {
		t.Errorf("snapshot category = %q, want %q", snap.Category, "ci")
	}
}

func TestCategory_TagRendering_NoAnsiLeak(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	_ = dt.AddActivity(ActivityID("a"), nil)
	dt.showCategory = true

	snaps := newSnapshotBuilder()
	snaps.setCategory(ActivityID("a"), "Alpha", ActivityStatusRunning, ActivityCategory("x"))

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	// The [x] tag should be present in stripped output.
	stripped := lipgloss.NewStyle().Render(got)
	_ = stripped

	// Verify plain text after ANSI strip contains the tag.
	if !strings.Contains(stripANSI(got), "[x]") {
		t.Errorf("expected [x] tag in stripped output:\n%s", got)
	}
}

func TestCategory_ColorTint(t *testing.T) {
	t.Parallel()

	catColor := lipgloss.Color("#ff0000")

	theme := ThemeDefault
	theme.CategoryColors = map[ActivityCategory]color.Color{
		ActivityCategory("build"): catColor,
	}

	snap := ActivitySnapshot{
		Label:    "Compile",
		Status:   ActivityStatusRunning,
		Symbol:   SymbolRunning,
		Color:    Colors.Running,
		Category: ActivityCategory("build"),
		Kind:     ActivityKindTask,
	}

	_, c := formatActivityLabelWithOptions(snap, theme, labelOptions{ShowCategory: true})
	if c == nil {
		t.Fatal("expected non-nil color for category tint")
	}

	// Without ShowCategory, color should be the default status color.
	_, c2 := formatActivityLabelWithOptions(snap, theme, labelOptions{ShowCategory: false})
	if c2 == nil {
		t.Fatal("expected non-nil default color")
	}

	// The tinted color should differ from the default.
	if c == c2 {
		t.Error("category tint should override default status color")
	}
}

func TestCategory_ColorTint_NoCategoryColor(t *testing.T) {
	t.Parallel()

	theme := ThemeDefault
	// No CategoryColors defined at all.

	snap := ActivitySnapshot{
		Label:    "Compile",
		Status:   ActivityStatusRunning,
		Symbol:   SymbolRunning,
		Color:    Colors.Running,
		Category: ActivityCategory("build"),
		Kind:     ActivityKindTask,
	}

	_, c := formatActivityLabelWithOptions(snap, theme, labelOptions{ShowCategory: true})
	// Should fall back to the snapshot's default color.
	if c != snap.Color {
		t.Error("expected default status color when theme has no category color")
	}
}
