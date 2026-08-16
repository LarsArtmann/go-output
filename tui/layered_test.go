package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

func TestProgressModel_LayeredMode_RendersLayers(t *testing.T) {
	t.Parallel()

	model, _ := setupLayeredTestModel(t, true)

	got := model.renderDependencyTree()

	if !strings.Contains(got, "Layer 0") {
		t.Errorf("expected Layer 0 in output:\n%s", got)
	}

	if !strings.Contains(got, "Layer 1") {
		t.Errorf("expected Layer 1 in output:\n%s", got)
	}

	if !strings.Contains(got, "Root") {
		t.Errorf("expected Root label in output:\n%s", got)
	}

	if !strings.Contains(got, "Child") {
		t.Errorf("expected Child label in output:\n%s", got)
	}
}

func TestProgressModel_LayeredMode_MouseClickSelectsNode(t *testing.T) {
	t.Parallel()

	model, _ := setupLayeredTestModel(t, true)

	// Render once so visibleEntries AND treeStartLine reflect the real layout.
	_ = model.renderNOMStyle()

	// Layer 0 header is tree-relative line 0; the first node row is line 1.
	clickY := model.treeStartLine + 1

	m := clickAt(model, clickY)

	if m.selectedNode != nom.ActivityID("root") {
		t.Errorf("selectedNode = %q, want %q", m.selectedNode, "root")
	}
}

func TestProgressModel_LayeredMode_ClickHeaderClearsSelection(t *testing.T) {
	t.Parallel()

	model, _ := setupLayeredTestModel(t, false)

	_ = model.renderNOMStyle()
	model.selectedNode = nom.ActivityID("root")

	clickY := model.treeStartLine + 0

	updatedModel, _ := model.Update(tea.MouseClickMsg{
		X: 5, Y: clickY, Button: tea.MouseLeft,
	})

	m := updatedModel.(*ProgressModel)

	if m.selectedNode != "" {
		t.Errorf("clicking a header should clear selection, got %q", m.selectedNode)
	}
}

func TestProgressModel_LayeredMode_SelectionHighlight(t *testing.T) {
	t.Parallel()

	model, _ := setupLayeredTestModel(t, true)

	model.selectedNode = nom.ActivityID("root")

	got := model.renderDependencyTree()

	if !strings.Contains(got, "Root") {
		t.Errorf("expected Root in rendered output:\n%s", got)
	}
}

// TestProgressModel_LayeredMode_InlineRendererDispatch uses the real render
// path that the inline renderer would take, proving RenderWithSnapshots already
// dispatches to layered mode without changes to the renderer itself.
func TestProgressModel_LayeredMode_InlineRendererDispatch(t *testing.T) {
	t.Parallel()

	sub := nom.NewNOMSubscriber(nom.WithCachePath(t.TempDir() + "/cache.csv"))
	sub.DependencyTree().SetRenderMode(nom.RenderModeLayered)

	ctx := t.Context()

	_ = sub.OnEvent(ctx, nom.ActivityRegistered{ID: "root", Name: "Root"})
	_ = sub.OnEvent(ctx, nom.ActivityRegistered{ID: "child", Name: "Child", Deps: []nom.ActivityID{"root"}})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: "root", Name: "Root"})

	buf := &strings.Builder{}
	renderer := nom.NewInlineRenderer(sub, buf, 20)
	renderer.SetNoColor(true)

	renderer.Draw()

	got := buf.String()

	if !strings.Contains(got, "Layer 0") {
		t.Errorf("expected Layer 0 in inline renderer output:\n%s", got)
	}

	if !strings.Contains(got, "Root") {
		t.Errorf("expected Root label in inline renderer output:\n%s", got)
	}
}
