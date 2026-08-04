package nom

import (
	"strings"
	"testing"
	"time"
)

func TestCollectVisibleNodes_RootsSortedByPriority(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("aaa-completed"), nil)
	dt.AddActivity(ActivityID("bbb-running"), nil)
	dt.AddActivity(ActivityID("ccc-pending"), nil)
	dt.AddActivity(ActivityID("ddd-failed"), nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("aaa-completed"), "AAA Completed", ActivityStatusCompleted, 5*time.Second)
	snaps.set(ActivityID("bbb-running"), "BBB Running", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("ccc-pending"), "CCC Pending", ActivityStatusPending, 0)
	snaps.set(ActivityID("ddd-failed"), "DDD Failed", ActivityStatusFailed, 2*time.Second)

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	visible := dt.collectVisibleNodes(snaps.snaps, 10)
	if len(visible) < 4 {
		t.Fatalf("expected 4 visible entries, got %d", len(visible))
	}

	wantOrder := []string{
		"ddd-failed",
		"bbb-running",
		"ccc-pending",
		"aaa-completed",
	}

	for i, want := range wantOrder {
		node := visible[i].Node
		if node == nil {
			t.Errorf("visible[%d].Node is nil", i)
			continue
		}

		if string(node.ID) != want {
			t.Errorf("visible[%d].ID = %q, want %q (priority order: failed > running > pending > completed)",
				i, node.ID, want)
		}
	}
}

func TestCollectVisibleNodes_RunningRootsVisibleWithManyCompleted(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	for i := range 20 {
		dt.AddActivity(ActivityID(string(rune('a'+i))), nil)
	}

	dt.AddActivity(ActivityID("z-running"), nil)

	snaps := newSnapshotBuilder()
	for i := range 20 {
		id := ActivityID(string(rune('a' + i)))
		snaps.set(id, string(id), ActivityStatusCompleted, 1*time.Second)
	}

	snaps.set(ActivityID("z-running"), "Z Running", ActivityStatusRunning, 500*time.Millisecond)

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	visible := dt.collectVisibleNodes(snaps.snaps, 5)

	var foundRunning bool

	for _, entry := range visible {
		if entry.Node != nil && string(entry.Node.ID) == "z-running" {
			foundRunning = true
			break
		}
	}

	if !foundRunning {
		t.Errorf("running root 'z-running' must be visible in viewport of 5 even with 20 completed roots; visible IDs: %v",
			visibleNodeIDs(visible))
	}
}

func TestWalkSubtree_PartialPhaseCollapseUnderPressure(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	phaseID := ActivityID("golangci-lint")
	dt.AddActivity(phaseID, nil)

	for i := range 10 {
		childID := ActivityID("golangci-lint [mod" + string(rune('0'+i)) + "]")
		dt.AddActivity(childID, []ActivityID{phaseID})
	}

	snaps := newSnapshotBuilder()
	snaps.setPhase(phaseID, "golangci-lint", ActivityStatusRunning, 5*time.Second)

	for i := range 10 {
		childID := ActivityID("golangci-lint [mod" + string(rune('0'+i)) + "]")
		status := ActivityStatusCompleted
		if i >= 5 {
			status = ActivityStatusRunning
		}

		snaps.set(childID, string(childID), status, 2*time.Second)
	}

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	visible := dt.collectVisibleNodes(snaps.snaps, 3)

	if len(visible) != 1 {
		t.Fatalf("expected 1 visible entry (collapsed phase), got %d: %+v", len(visible), visible)
	}

	entry := visible[0]
	if entry.PhaseCounts == nil {
		t.Fatal("expected PhaseCounts on collapsed phase entry")
	}

	pc := *entry.PhaseCounts
	if pc.Completed != 5 {
		t.Errorf("PhaseCounts.Completed = %d, want 5", pc.Completed)
	}

	if pc.Running != 5 {
		t.Errorf("PhaseCounts.Running = %d, want 5", pc.Running)
	}

	if pc.Total() != 10 {
		t.Errorf("PhaseCounts.Total() = %d, want 10", pc.Total())
	}

	if !pc.IsPartial() {
		t.Error("PhaseCounts.IsPartial() = false, want true (phase has running children)")
	}
}

func TestWalkSubtree_PhaseNotCollapsedWhenViewportFits(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	phaseID := ActivityID("phase")
	dt.AddActivity(phaseID, nil)

	for i := range 5 {
		childID := ActivityID("child" + string(rune('0'+i)))
		dt.AddActivity(childID, []ActivityID{phaseID})
	}

	snaps := newSnapshotBuilder()
	snaps.setPhase(phaseID, "Phase", ActivityStatusRunning, 1*time.Second)

	for i := range 5 {
		childID := ActivityID("child" + string(rune('0'+i)))
		snaps.set(childID, string(childID), ActivityStatusRunning, 500*time.Millisecond)
	}

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	visible := dt.collectVisibleNodes(snaps.snaps, 50)

	if len(visible) != 6 {
		t.Fatalf("expected 6 visible entries (1 phase + 5 children, no collapse), got %d", len(visible))
	}

	if visible[0].PhaseCounts != nil {
		t.Error("phase should NOT be collapsed when viewport has room")
	}
}

func TestComputePartialPhaseCounts(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("p"), nil)
	dt.AddActivity(ActivityID("c1"), []ActivityID{"p"})
	dt.AddActivity(ActivityID("c2"), []ActivityID{"p"})
	dt.AddActivity(ActivityID("c3"), []ActivityID{"p"})
	dt.AddActivity(ActivityID("c4"), []ActivityID{"p"})

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("c1"), "C1", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("c2"), "C2", ActivityStatusFailed, 2*time.Second)
	snaps.set(ActivityID("c3"), "C3", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("c4"), "C4", ActivityStatusPending, 0)

	children := dt.GetNode(ActivityID("p")).Children
	pc := computePartialPhaseCounts(snaps.snaps, children)

	if pc.Completed != 1 {
		t.Errorf("Completed = %d, want 1", pc.Completed)
	}

	if pc.Failed != 1 {
		t.Errorf("Failed = %d, want 1", pc.Failed)
	}

	if pc.Running != 1 {
		t.Errorf("Running = %d, want 1", pc.Running)
	}

	if pc.Pending != 1 {
		t.Errorf("Pending = %d, want 1", pc.Pending)
	}

	if pc.Total() != 4 {
		t.Errorf("Total() = %d, want 4", pc.Total())
	}

	if !pc.IsPartial() {
		t.Error("IsPartial() = false, want true")
	}

	if pc.MaxElapsed != 3*time.Second {
		t.Errorf("MaxElapsed = %v, want 3s", pc.MaxElapsed)
	}
}

func TestComputePartialPhaseCounts_AllTerminal(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("p"), nil)
	dt.AddActivity(ActivityID("c1"), []ActivityID{"p"})
	dt.AddActivity(ActivityID("c2"), []ActivityID{"p"})

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("c1"), "C1", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("c2"), "C2", ActivityStatusFailed, 2*time.Second)

	children := dt.GetNode(ActivityID("p")).Children
	pc := computePartialPhaseCounts(snaps.snaps, children)

	if pc.IsPartial() {
		t.Error("IsPartial() = true, want false (all terminal)")
	}
}

func TestCountNonTerminalChildren(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("p"), nil)

	for _, id := range []string{"a", "b", "c", "d", "e"} {
		dt.AddActivity(ActivityID(id), []ActivityID{"p"})
	}

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "A", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("b"), "B", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("c"), "C", ActivityStatusFailed, 0)
	snaps.set(ActivityID("d"), "D", ActivityStatusRunning, 0)
	snaps.set(ActivityID("e"), "E", ActivityStatusPending, 0)

	children := dt.GetNode(ActivityID("p")).Children
	count := countNonTerminalChildren(snaps.snaps, children)

	if count != 2 {
		t.Errorf("count = %d, want 2 (1 running + 1 pending)", count)
	}
}

func TestPhaseCounts_Total(t *testing.T) {
	t.Parallel()

	pc := PhaseCounts{
		Completed: 3,
		Failed:    1,
		Running:   2,
		Pending:   4,
	}

	if pc.Total() != 10 {
		t.Errorf("Total() = %d, want 10", pc.Total())
	}
}

func TestPhaseCounts_IsPartial(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pc   PhaseCounts
		want bool
	}{
		{
			name: "all terminal",
			pc:   PhaseCounts{Completed: 5, Failed: 1},
			want: false,
		},
		{
			name: "has running",
			pc:   PhaseCounts{Completed: 3, Running: 2},
			want: true,
		},
		{
			name: "has pending",
			pc:   PhaseCounts{Completed: 3, Pending: 1},
			want: true,
		},
		{
			name: "zero value",
			pc:   PhaseCounts{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.pc.IsPartial(); got != tt.want {
				t.Errorf("IsPartial() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatCollapsedPhaseLabel_PartialProgress(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Kind:  ActivityKindPhase,
		Label: "golangci-lint",
	}

	pc := PhaseCounts{
		Completed:  15,
		Running:    10,
		Pending:    5,
		MaxElapsed: 12 * time.Second,
	}

	display, c := formatCollapsedPhaseLabel(snap, pc, ThemeDefault)
	plain := stripANSI(display)

	if !strings.Contains(plain, "golangci-lint") {
		t.Errorf("display should contain phase label, got: %s", plain)
	}

	if !strings.Contains(plain, "15/30") {
		t.Errorf("display should show progress 15/30, got: %s", plain)
	}

	if !strings.Contains(plain, "12.0s") {
		t.Errorf("display should show max elapsed, got: %s", plain)
	}

	if c != ThemeDefault.Colors.Running {
		t.Error("color should be Running for partial phase")
	}
}

func TestFormatCollapsedPhaseLabel_AllTerminal(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Kind:  ActivityKindPhase,
		Label: "Code Formatting",
	}

	pc := PhaseCounts{
		Completed:  6,
		MaxElapsed: 4 * time.Second,
	}

	display, c := formatCollapsedPhaseLabel(snap, pc, ThemeDefault)
	plain := stripANSI(display)

	if !strings.Contains(plain, "6/6") {
		t.Errorf("display should show 6/6, got: %s", plain)
	}

	if c != ThemeDefault.Colors.Phase {
		t.Error("color should be Phase for fully terminal phase")
	}
}

func TestFormatCollapsedPhaseLabel_PartialWithFailures(t *testing.T) {
	t.Parallel()

	snap := ActivitySnapshot{
		Kind:  ActivityKindPhase,
		Label: "Tests",
	}

	pc := PhaseCounts{
		Completed:  3,
		Failed:     1,
		Running:    2,
		MaxElapsed: 5 * time.Second,
	}

	_, c := formatCollapsedPhaseLabel(snap, pc, ThemeDefault)

	if c != ThemeDefault.Colors.Failed {
		t.Error("color should be Failed when partial phase has failures")
	}
}

func TestSortRootsByPriority_RunningBeforeCompleted(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("alpha"), nil)
	dt.AddActivity(ActivityID("beta"), nil)
	dt.AddActivity(ActivityID("gamma"), nil)

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("alpha"), "Alpha", ActivityStatusCompleted, 10*time.Second)
	snaps.set(ActivityID("beta"), "Beta", ActivityStatusRunning, 5*time.Second)
	snaps.set(ActivityID("gamma"), "Gamma", ActivityStatusFailed, 2*time.Second)

	sorted := dt.sortRootsByPriority(dt.roots, snaps.snaps, nil)

	if len(sorted) != 3 {
		t.Fatalf("expected 3 sorted roots, got %d", len(sorted))
	}

	if string(sorted[0].ID) != "gamma" {
		t.Errorf("first root should be 'gamma' (failed), got %q", sorted[0].ID)
	}

	if string(sorted[1].ID) != "beta" {
		t.Errorf("second root should be 'beta' (running), got %q", sorted[1].ID)
	}

	if string(sorted[2].ID) != "alpha" {
		t.Errorf("third root should be 'alpha' (completed), got %q", sorted[2].ID)
	}
}

func TestSortRootsByPriority_SingleRoot(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("only"), nil)

	if err := dt.Build(); err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("only"), "Only", ActivityStatusRunning, 0)

	sorted := dt.sortRootsByPriority(dt.roots, snaps.snaps, nil)
	if len(sorted) != 1 {
		t.Fatalf("expected 1 root, got %d", len(sorted))
	}
}

func visibleNodeIDs(entries []VisibleEntry) []string {
	ids := make([]string, 0, len(entries))

	for _, e := range entries {
		if e.Node != nil {
			ids = append(ids, string(e.Node.ID))
		}
	}

	return ids
}
