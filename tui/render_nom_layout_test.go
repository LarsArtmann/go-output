package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/larsartmann/go-output/nom"
)

// TestNOMStyle_TreeStartLineMatchesRenderedLayout pins the mouse-click
// coordinate mapping to the actual rendered layout. The click handler maps
// mouse.Y via m.treeStartLine; if the layout changes (title block, message
// section, section gaps) without nomTreeStartLine following, clicks silently
// select the wrong node — the exact bug fixed alongside this test. This test
// renders, splits the output into lines, and asserts the line at
// treeStartLine really is the first tree row, for both the no-message and
// with-message layouts.
func TestNOMStyle_TreeStartLineMatchesRenderedLayout(t *testing.T) {
	t.Parallel()

	newModel := func() *ProgressModel {
		model := newTestModel()
		model.width = 80
		model.height = 24
		model.displayMode = DisplayModeNOM

		tree := setupTestTree(model)
		_ = tree.AddActivity(nom.ActivityID("step-a"), nil)
		_ = tree.AddActivity(nom.ActivityID("step-b"), []nom.ActivityID{"step-a"})
		_ = tree.GetRootNodes()

		addTestActivity(model, "step-a", "Step A", nom.ActivityStatusPending)
		addTestActivity(model, "step-b", "Step B", nom.ActivityStatusPending)

		return model
	}

	assertFirstTreeRow := func(t *testing.T, model *ProgressModel, label string) {
		t.Helper()

		rendered := model.renderNOMStyle()
		lines := strings.Split(rendered, "\n")

		if model.treeStartLine < 0 || model.treeStartLine >= len(lines) {
			t.Fatalf("%s: treeStartLine %d out of range (%d lines)", label, model.treeStartLine, len(lines))
		}

		for i := 0; i < model.treeStartLine; i++ {
			if strings.Contains(ansi.Strip(lines[i]), "Step A") {
				t.Fatalf(
					"%s: first tree row renders at line %d but treeStartLine is %d:\n%s",
					label, i, model.treeStartLine, rendered,
				)
			}
		}

		first := ansi.Strip(lines[model.treeStartLine])
		if !strings.Contains(first, "Step A") {
			t.Errorf(
				"%s: line at treeStartLine %d is %q, want the first tree row containing step-a:\n%s",
				label, model.treeStartLine, first, rendered,
			)
		}
	}

	t.Run("no message", func(t *testing.T) {
		t.Parallel()
		assertFirstTreeRow(t, newModel(), "no message")
	})

	t.Run("with message", func(t *testing.T) {
		t.Parallel()

		model := newModel()
		model.currentMessage = "building..."
		assertFirstTreeRow(t, model, "with message")
	})
}
