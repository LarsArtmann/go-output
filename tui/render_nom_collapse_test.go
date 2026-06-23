package tui

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/go-output/nom"
)

// TestProgressModel_RenderDependencyTree_CollapseMarkerNoPanic reproduces the
// production panic: when completed children are elided under terminal height
// pressure, the tree emits a synthetic collapse-marker line. The old TUI render
// path (VisibleNodesWithSnapshots + RenderNode) smuggled that marker through as
// a nil *ActivityNode and crashed dereferencing node.ID in RenderNode. The fixed
// path (VisibleEntriesWithSnapshots + RenderVisibleEntry) renders the marker and
// never dereferences a nil.
func TestProgressModel_RenderDependencyTree_CollapseMarkerNoPanic(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 12 // treeHeight = 12 - chromeLines(8) = 4 → elision kicks in
	model.displayMode = DisplayModeNOM
	model.treeStartLine = 2

	ctx := context.Background()
	sub := model.nomSubscriber

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{
		ID: nom.WorkflowID("wf-collapse"), Name: nom.WorkflowName("Collapse"),
	})
	_ = sub.OnEvent(ctx, nom.ActivityRegistered{ID: "root", Name: "Root"})

	// 6 children: first 4 completed (elided), last 2 running (kept).
	for i := range 6 {
		id := nom.ActivityID(fmt.Sprintf("c%d", i))
		name := nom.ActivityName(fmt.Sprintf("Child %d", i))
		_ = sub.OnEvent(ctx, nom.ActivityRegistered{ID: id, Name: name, Deps: []nom.ActivityID{"root"}})

		if i < 4 {
			_ = sub.OnEvent(ctx, nom.ActivityCompleted{ID: id, Name: name})
		} else {
			_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: id, Name: name, Deps: []nom.ActivityID{"root"}})
		}
	}

	model.dependencyTree = sub.GetDependencyTree()

	// Must not panic and must surface the collapse marker.
	got := model.renderDependencyTree()
	if got == "" {
		t.Fatal("renderDependencyTree() should produce output under height pressure")
	}

	if !strings.Contains(got, "⋯") || !strings.Contains(got, "4 completed") {
		t.Errorf("expected a '⋯ 4 completed' collapse marker line, got:\n%s", got)
	}
}

// TestProgressModel_ScrollRendering_ShowsCorrectWindow verifies that scrolling
// the dependency tree shows the correct entries for a given scrollOffset, not
// the entire tree. Regression for the O(n) render+clip → entry-level windowing
// optimization: the old code rendered the entire tree then string-split to
// extract the visible lines; the new code slices the entry list before rendering.
func TestProgressModel_ScrollRendering_ShowsCorrectWindow(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 18 // treeHeight = 18 - chromeLines(8) = 10
	model.displayMode = DisplayModeNOM
	model.treeStartLine = 2
	model.workflowState = workflowStateRunning

	tree := newTestTree(30) // root + 29 children = 30 entries
	model.dependencyTree = tree

	// No scroll: first 10 entries visible
	model.scrollOffset = 0

	got := model.renderDependencyTree()
	if got == "" {
		t.Fatal("scrollOffset=0 should produce output")
	}

	if len(model.visibleEntries) == 0 {
		t.Fatal("scrollOffset=0 should populate visibleEntries")
	}

	// All entries at offset 0 should be the first 10 nodes
	entryCount0 := len(model.visibleEntries)

	// Scroll to offset 5: entries 5..14 should be visible
	model.scrollOffset = 5
	_ = model.renderDependencyTree()

	if len(model.visibleEntries) == 0 {
		t.Fatal("scrollOffset=5 should populate visibleEntries")
	}

	// The first visible entry at offset 5 should differ from offset 0
	if model.visibleEntries[0].Node != nil &&
		model.visibleEntries[0].Node.ID == nom.ActivityID("root") {
		t.Error("scrollOffset=5 should not show root at first position")
	}

	// Scroll to bottom sentinel: should clamp to last page
	model.scrollOffset = scrollToBottomSentinel
	_ = model.renderDependencyTree()

	maxExpected := 30 - entryCount0
	if model.scrollOffset != maxExpected {
		t.Errorf("scrollToBottom should clamp to %d, got %d", maxExpected, model.scrollOffset)
	}
}
