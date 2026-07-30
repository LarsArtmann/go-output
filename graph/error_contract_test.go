package graph

import (
	"errors"
	"fmt"
	"testing"
)

func TestTypedErrors_AsType_ThroughWrapping(t *testing.T) {
	t.Parallel()

	t.Run("InvalidRankDirError from ParseRankDir", func(t *testing.T) {
		t.Parallel()

		_, err := ParseRankDir("bogus")
		if err == nil {
			t.Fatal("expected error from ParseRankDir")
		}

		wrapped := fmt.Errorf("graph layout: %w", err)

		extracted, ok := errors.AsType[*InvalidRankDirError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidRankDirError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(AllRankDirs) {
			t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllRankDirs))
		}
	})

	t.Run("InvalidSplineStyleError from ParseSplineStyle", func(t *testing.T) {
		t.Parallel()

		_, err := ParseSplineStyle("bogus")
		if err == nil {
			t.Fatal("expected error from ParseSplineStyle")
		}

		wrapped := fmt.Errorf("edge routing: %w", err)

		extracted, ok := errors.AsType[*InvalidSplineStyleError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidSplineStyleError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(AllSplineStyles) {
			t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllSplineStyles))
		}
	})
}
