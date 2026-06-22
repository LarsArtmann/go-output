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
