package nom

import (
	"context"
	"image/color"
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestStatusRegistry_CoreStatuses(t *testing.T) {
	t.Parallel()

	core := []struct {
		status   ActivityStatus
		name     string
		symbol   Symbol
		interest int
	}{
		{ActivityStatusPending, "pending", SymbolPending, 2},
		{ActivityStatusRunning, "running", SymbolRunning, 1},
		{ActivityStatusCompleted, "completed", SymbolCompleted, 3},
		{ActivityStatusFailed, "failed", SymbolFailed, 0},
	}

	for _, tc := range core {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			def, ok := LookupStatus(tc.status)
			if !ok {
				t.Fatalf("LookupStatus(%d) returned false for core status", tc.status)
			}

			testhelpers.AssertEqual(t, "Name", tc.status, def.Name, tc.name)
			testhelpers.AssertEqual(t, "Symbol", tc.status, def.Symbol, tc.symbol)
			testhelpers.AssertEqual(t, "Interest", tc.status, def.Interest, tc.interest)
		})
	}
}

func TestStatusRegistry_AllActivityStatuses_Core(t *testing.T) {
	t.Parallel()

	statuses := AllActivityStatuses()

	want := []ActivityStatus{
		ActivityStatusPending,
		ActivityStatusRunning,
		ActivityStatusCompleted,
		ActivityStatusFailed,
	}

	if len(statuses) < len(want) {
		t.Fatalf("AllActivityStatuses() returned %d values, want at least %d", len(statuses), len(want))
	}

	for i, s := range want {
		if statuses[i] != s {
			t.Errorf("AllActivityStatuses()[%d] = %v, want %v", i, statuses[i], s)
		}
	}
}

func TestStatusRegistry_AllRegisteredStatuses_Core(t *testing.T) {
	t.Parallel()

	defs := AllRegisteredStatuses()

	want := []string{"pending", "running", "completed", "failed"}
	if len(defs) < len(want) {
		t.Fatalf("AllRegisteredStatuses() returned %d values, want at least %d", len(defs), len(want))
	}

	for i, name := range want {
		if defs[i].Name != name {
			t.Errorf("AllRegisteredStatuses()[%d].Name = %q, want %q", i, defs[i].Name, name)
		}
	}
}

func TestStatusRegistry_RegisterStatus(t *testing.T) {
	t.Parallel()

	custom := RegisterStatus(
		"cached",
		"⚡",
		color.RGBA{R: 0x38, G: 0xbd, B: 0xf2, A: 0xff},
		2,
		output.NodeShapeCylinder,
		output.GraphStyle{
			Fill:      "#38bdf2",
			Stroke:    "#0ea5e9",
			FontColor: "#0f172a",
		},
	)

	if custom < ActivityStatusFailed+1 {
		t.Fatalf("RegisterStatus returned ID %d, want >= 4", custom)
	}

	def, ok := LookupStatus(custom)
	if !ok {
		t.Fatalf("LookupStatus(%d) returned false after registration", custom)
	}

	testhelpers.AssertEqual(t, "Name", custom, def.Name, "cached")
	testhelpers.AssertEqual(t, "Symbol", custom, def.Symbol, Symbol("⚡"))
	testhelpers.AssertEqual(t, "Interest", custom, def.Interest, 2)
}

func TestStatusRegistry_RegisterStatus_DeduplicatesByName(t *testing.T) {
	t.Parallel()

	id1 := RegisterStatus(
		"skipped-test",
		"⊘",
		Colors.Fallback,
		2,
		output.NodeShapeEllipse,
		output.GraphStyle{},
	)

	id2 := RegisterStatus(
		"skipped-test",
		"X",
		Colors.Failed,
		0,
		output.NodeShapeDiamond,
		output.GraphStyle{Fill: "#000"},
	)

	if id1 != id2 {
		t.Fatalf("RegisterStatus returned different IDs for same name: %d vs %d", id1, id2)
	}

	def, ok := LookupStatus(id1)
	if !ok {
		t.Fatalf("LookupStatus(%d) returned false", id1)
	}

	// First registration wins.
	testhelpers.AssertEqual(t, "Symbol", id1, def.Symbol, Symbol("⊘"))
}

func TestStatusRegistry_CustomStatusMethods(t *testing.T) {
	t.Parallel()

	custom := RegisterStatus(
		"warn",
		"▲",
		Colors.Failed,
		1,
		output.NodeShapeDiamond,
		output.GraphStyle{Fill: "#f59e0b"},
	)

	testhelpers.AssertEqual(t, "String", custom, custom.String(), "warn")
	testhelpers.AssertEqual(t, "Symbol", custom, custom.GetSymbol(), Symbol("▲"))
	testhelpers.AssertEqual(t, "Interest", custom, custom.Interest(), 1)

	if custom.GetColor() == nil {
		t.Errorf("GetColor() returned nil for custom status")
	}

	if !custom.IsValid() {
		t.Errorf("IsValid() = false for registered custom status")
	}

	parsed, err := ParseActivityStatus("warn")
	if err != nil {
		t.Fatalf("ParseActivityStatus(%q) error: %v", "warn", err)
	}

	if parsed != custom {
		t.Errorf("ParseActivityStatus(%q) = %v, want %v", "warn", parsed, custom)
	}
}

func TestStatusRegistry_AllActivityStatuses_IncludesCustom(t *testing.T) {
	t.Parallel()

	custom := RegisterStatus(
		"unique-custom-"+t.Name(),
		"★",
		Colors.Fallback,
		2,
		output.NodeShapeBox,
		output.GraphStyle{},
	)

	after := AllActivityStatuses()

	if !slices.Contains(after, custom) {
		t.Errorf("AllActivityStatuses() does not include custom status %v", custom)
	}
}

func TestStatusRegistry_AllowedValues_IncludesCustom(t *testing.T) {
	t.Parallel()

	RegisterStatus(
		"allowed-custom-"+t.Name(),
		"✓",
		Colors.Completed,
		3,
		output.NodeShapeBox,
		output.GraphStyle{},
	)

	values := ActivityStatus(0).AllowedValues()

	if !slices.Contains(values, "allowed-custom-"+t.Name()) {
		t.Errorf("AllowedValues() does not include custom status")
	}
}

// TestRegisterStatus_RendersInTree proves the RegisterStatus() ghost system
// works end-to-end: register a custom "skipped" status, create a subscriber,
// set an activity to that status, take a snapshot, and verify the custom
// symbol appears in the rendered tree.
func TestRegisterStatus_RendersInTree(t *testing.T) {
	t.Parallel()

	skipped := RegisterStatus(
		"skipped-e2e-"+t.Name(),
		"⊘",
		Colors.Pending,
		0,
		output.NodeShapeBox,
		output.GraphStyle{},
	)

	sub := newTestSubscriber(t)
	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-skip"), "")
	registerActivity(sub, ctx, ActivityID("lint"), ActivityName("Lint"))

	// Manually set the activity to the custom status via SetActivityState
	// (GetActivity returns a copy, so we must replace the stored pointer).
	act := sub.GetActivity(ActivityID("lint"))
	if act == nil {
		t.Fatal("GetActivity returned nil")
	}

	act.Status = skipped
	act.Symbol = skipped.GetSymbol()
	act.Color = skipped.GetColor()
	sub.SetActivityState(ActivityID("lint"), act)

	snaps := sub.SnapshotActivities()
	rendered := sub.DependencyTree().RenderWithSnapshots(snaps, 20, 80)

	if !strings.Contains(rendered, "Lint") {
		t.Errorf("rendered tree should contain 'Lint', got:\n%s", rendered)
	}

	if !strings.Contains(rendered, "⊘") {
		t.Errorf("rendered tree should contain custom symbol '⊘', got:\n%s", rendered)
	}
}
