package tui

import (
	"fmt"
	"testing"

	"github.com/larsartmann/go-output/nom"
)

func newTestTree(nodeCount int) *nom.DependencyTree {
	tree := nom.NewDependencyTree()

	tree.AddActivity(nom.ActivityID("root"), nom.NewActivity("root", "Root"), nil)

	for i := 1; i < nodeCount; i++ {
		id := nom.ActivityID(fmt.Sprintf("step-%d", i))
		name := fmt.Sprintf("Step %d", i)
		tree.AddActivity(id, nom.NewActivity(string(id), name), []nom.ActivityID{"root"})
	}

	_ = tree.GetRootNodes()

	return tree
}

// assertSingleStep fails the test if the reporter doesn't have exactly one
// reported step. Used by tests that call ReportStep once and expect a single
// step to land in the model.
func assertSingleStep(t *testing.T, reporter *BubbleTeaProgressReporter) {
	t.Helper()

	if len(reporter.model.steps) != 1 {
		t.Fatalf("steps count = %d, want 1", len(reporter.model.steps))
	}
}

// assertStepCurrent fails the test if the first reported step's Current field
// does not match the expected value.
func assertStepCurrent(t *testing.T, reporter *BubbleTeaProgressReporter, want uint) {
	t.Helper()

	if got := reporter.model.steps[0].Current; got != want {
		t.Errorf("step current = %d, want %d", got, want)
	}
}
