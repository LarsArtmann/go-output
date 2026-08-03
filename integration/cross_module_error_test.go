package integration

import (
	"errors"
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/graph"
)

// TestCrossModule_ErrorMatching verifies that errors.Is and errors.AsType
// work correctly across module boundaries. This is the integration-level
// proof that the three-tier error system (sentinels + typed structs + wrapping)
// is consistent end-to-end.

func TestCrossModule_RootSentinels_Is_AcrossWrapping(t *testing.T) {
	t.Parallel()

	t.Run("ErrColumnMismatch through RenderTable dispatch", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"A", "B"})
		data.AddRow([]string{"1", "2"})
		data.SetFooter([]string{"total"})

		err := output.RenderTable(data, output.FormatJSON, output.RenderOptions{})
		if err == nil {
			t.Fatal("expected error for column count mismatch")
		}

		if !errors.Is(err, output.ErrColumnMismatch) {
			t.Errorf("errors.Is(err, ErrColumnMismatch) = false; err=%v", err)
		}
	})

	t.Run("ErrNilRow through Validate", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"A"})
		data.AddRow([]string{"1"})
		data.Rows = append(data.Rows, nil)

		err := data.Validate()
		if err == nil {
			t.Fatal("expected error from Validate with nil row")
		}

		if !errors.Is(err, output.ErrNilRow) {
			t.Errorf("errors.Is(err, ErrNilRow) = false; err=%v", err)
		}
	})
}

func TestCrossModule_RootTypedErrors_AsType(t *testing.T) {
	t.Parallel()

	t.Run("InvalidFormatError from ParseFormat", func(t *testing.T) {
		t.Parallel()

		_, err := output.ParseFormat("bogus")
		if err == nil {
			t.Fatal("expected error from ParseFormat")
		}

		wrapped := fmt.Errorf("config: %w", err)

		extracted, ok := errors.AsType[*output.InvalidFormatError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*output.InvalidFormatError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})

	t.Run("InvalidShapeError from ParseShape", func(t *testing.T) {
		t.Parallel()

		_, err := output.ParseShape("bogus")
		if err == nil {
			t.Fatal("expected error from ParseShape")
		}

		wrapped := fmt.Errorf("data: %w", err)

		extracted, ok := errors.AsType[*output.InvalidShapeError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*output.InvalidShapeError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})

	t.Run("InvalidColorModeError from ParseColorMode", func(t *testing.T) {
		t.Parallel()

		_, err := output.ParseColorMode("bogus")
		if err == nil {
			t.Fatal("expected error from ParseColorMode")
		}

		wrapped := fmt.Errorf("theme: %w", err)

		extracted, ok := errors.AsType[*output.InvalidColorModeError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*output.InvalidColorModeError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})
}

func TestCrossModule_SubModuleTypedErrors_AsType(t *testing.T) {
	t.Parallel()

	t.Run("d2.InvalidDirectionError extracted from integration", func(t *testing.T) {
		t.Parallel()

		_, err := d2.ParseDirection("bogus")
		if err == nil {
			t.Fatal("expected error from d2.ParseDirection")
		}

		wrapped := fmt.Errorf("diagram setup: %w", err)

		extracted, ok := errors.AsType[*d2.InvalidDirectionError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*d2.InvalidDirectionError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) == 0 {
			t.Error("Allowed should not be empty")
		}
	})

	t.Run("d2.InvalidNodeShapeError extracted from integration", func(t *testing.T) {
		t.Parallel()

		_, err := d2.ParseNodeShape("bogus")
		if err == nil {
			t.Fatal("expected error from d2.ParseNodeShape")
		}

		wrapped := fmt.Errorf("node config: %w", err)

		extracted, ok := errors.AsType[*d2.InvalidNodeShapeError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*d2.InvalidNodeShapeError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})

	t.Run("graph.InvalidRankDirError extracted from integration", func(t *testing.T) {
		t.Parallel()

		_, err := graph.ParseRankDir("bogus")
		if err == nil {
			t.Fatal("expected error from graph.ParseRankDir")
		}

		wrapped := fmt.Errorf("graph layout: %w", err)

		extracted, ok := errors.AsType[*graph.InvalidRankDirError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*graph.InvalidRankDirError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) == 0 {
			t.Error("Allowed should not be empty")
		}
	})

	t.Run("graph.InvalidSplineStyleError extracted from integration", func(t *testing.T) {
		t.Parallel()

		_, err := graph.ParseSplineStyle("bogus")
		if err == nil {
			t.Fatal("expected error from graph.ParseSplineStyle")
		}

		wrapped := fmt.Errorf("edge style: %w", err)

		extracted, ok := errors.AsType[*graph.InvalidSplineStyleError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*graph.InvalidSplineStyleError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}
	})
}

func TestCrossModule_ErrorsAreDistinct(t *testing.T) {
	t.Parallel()

	_, d2Err := d2.ParseDirection("bogus")
	_, graphErr := graph.ParseRankDir("bogus")
	_, rootErr := output.ParseFormat("bogus")

	if _, ok := errors.AsType[*d2.InvalidDirectionError](graphErr); ok {
		t.Error("graph error should not match *d2.InvalidDirectionError")
	}

	if _, ok := errors.AsType[*graph.InvalidRankDirError](d2Err); ok {
		t.Error("d2 error should not match *graph.InvalidRankDirError")
	}

	if _, ok := errors.AsType[*d2.InvalidDirectionError](rootErr); ok {
		t.Error("root error should not match *d2.InvalidDirectionError")
	}

	if _, ok := errors.AsType[*output.InvalidFormatError](d2Err); ok {
		t.Error("d2 error should not match *output.InvalidFormatError")
	}
}

func TestCrossModule_DeepWrappingPreserveType(t *testing.T) {
	t.Parallel()

	_, baseErr := d2.ParseNodeShape("bogus")
	wrapped1 := fmt.Errorf("layer 1: %w", baseErr)
	wrapped2 := fmt.Errorf("layer 2: %w", wrapped1)
	wrapped3 := fmt.Errorf("layer 3: %w", wrapped2)

	extracted, ok := errors.AsType[*d2.InvalidNodeShapeError](wrapped3)
	if !ok {
		t.Fatalf("errors.AsType should find the typed error through 3 layers of wrapping; err=%v", wrapped3)
	}

	if extracted.Value != "bogus" {
		t.Errorf("Value = %q, want %q after deep wrapping", extracted.Value, "bogus")
	}
}
