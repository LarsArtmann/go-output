package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// TestTUIRendersProgress verifies that the TUI's NOM tree rendering shows the
// sub-step progress message. This proves the TUI consumes ActivityProgress via
// its subscriber and renders it through the delegated RenderVisibleEntry path —
// the data and view layers are both complete, contrary to status-report item #6.
func TestTUIRendersProgress(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 20
	model.displayMode = DisplayModeNOM
	model.treeStartLine = 2

	ctx := context.Background()
	sub := model.nomSubscriber

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{ID: "wf", Name: "test"})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: "build", Name: "go-mod-tidy"})
	_ = sub.OnEvent(ctx, nom.ActivityProgress{
		ID:      "build",
		Name:    "go-mod-tidy",
		Message: "Tidying [2/26]",
	})

	model.dependencyTree = sub.GetDependencyTree()

	got := model.renderDependencyTree()
	if !strings.Contains(got, "Tidying [2/26]") {
		t.Errorf("TUI tree should render progress message:\n%s", got)
	}
}

// TestTUIRendersRetry verifies that the TUI's NOM tree rendering shows the
// retry suffix and optional reason via the delegated render path.
func TestTUIRendersRetry(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 20
	model.displayMode = DisplayModeNOM
	model.treeStartLine = 2

	ctx := context.Background()
	sub := model.nomSubscriber

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{ID: "wf", Name: "test"})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: "lint", Name: "golangci-lint"})
	_ = sub.OnEvent(ctx, nom.ActivityFailed{ID: "lint", Name: "golangci-lint"})
	_ = sub.OnEvent(ctx, nom.ActivityRetrying{
		ID: "lint", Name: "golangci-lint", Attempt: 1, Reason: "timeout",
	})

	model.dependencyTree = sub.GetDependencyTree()

	got := model.renderDependencyTree()
	if !strings.Contains(got, string(nom.SymbolRetrying)) {
		t.Errorf("TUI tree should render retry symbol:\n%s", got)
	}

	if !strings.Contains(got, "(timeout)") {
		t.Errorf("TUI tree should render retry reason:\n%s", got)
	}
}

// TestTUISummaryShowsRemaining verifies that the TUI NOM summary bar shows
// "~Xm left" when the subscriber has pending-activity estimates.
func TestTUISummaryShowsRemaining(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning

	ctx := context.Background()
	sub := model.nomSubscriber

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{ID: "wf", Name: "test"})
	_ = sub.OnEvent(ctx, nom.ActivityRegistered{ID: "a", Name: "a"})
	sub.SetEstimatedTime("a", 2*time.Minute)

	summary := model.renderNOMSummaryBar()
	if !strings.Contains(summary, "left") {
		t.Errorf("TUI summary should show '~Xm left' when estimates present:\n%s", summary)
	}
}
