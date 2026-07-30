package nom

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinels_Is_ThroughWrapping(t *testing.T) {
	t.Parallel()

	t.Run("ErrCycleDetected from Build", func(t *testing.T) {
		t.Parallel()

		tree := NewDependencyTree()
		_ = tree.AddActivity(ActivityID("a"), []ActivityID{"b"})
		_ = tree.AddActivity(ActivityID("b"), []ActivityID{"a"}) // creates cycle: a→b→a

		err := tree.Build()
		if err == nil {
			t.Fatal("expected cycle detection error")
		}

		wrapped := fmt.Errorf("building workflow: %w", err)

		if !errors.Is(wrapped, ErrCycleDetected) {
			t.Errorf("errors.Is(wrapped, ErrCycleDetected) = false; err=%v", wrapped)
		}
	})
}

func TestTypedErrors_AsType_ThroughWrapping(t *testing.T) {
	t.Parallel()

	t.Run("InvalidActivityStatusError from ParseActivityStatus", func(t *testing.T) {
		t.Parallel()

		_, err := ParseActivityStatus("bogus")
		if err == nil {
			t.Fatal("expected error from ParseActivityStatus")
		}

		wrapped := fmt.Errorf("event handler: %w", err)

		extracted, ok := errors.AsType[*InvalidActivityStatusError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidActivityStatusError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})

	t.Run("InvalidActivityKindError from ParseActivityKind", func(t *testing.T) {
		t.Parallel()

		_, err := ParseActivityKind("bogus")
		if err == nil {
			t.Fatal("expected error from ParseActivityKind")
		}

		wrapped := fmt.Errorf("activity setup: %w", err)

		extracted, ok := errors.AsType[*InvalidActivityKindError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidActivityKindError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})
}
