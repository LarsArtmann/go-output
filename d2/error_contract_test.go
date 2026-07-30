package d2

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinels_Is_ThroughWrapping(t *testing.T) {
	t.Parallel()

	t.Run("ErrInvalidDirection from ParseDirection", func(t *testing.T) {
		t.Parallel()

		_, err := ParseDirection("bogus")
		if err == nil {
			t.Fatal("expected error from ParseDirection")
		}

		wrapped := fmt.Errorf("diagram config: %w", err)

		if !errors.Is(wrapped, ErrInvalidDirection) {
			t.Errorf("errors.Is(wrapped, ErrInvalidDirection) = false; err=%v", wrapped)
		}
	})

	t.Run("ErrInvalidNodeShape from ParseNodeShape", func(t *testing.T) {
		t.Parallel()

		_, err := ParseNodeShape("bogus")
		if err == nil {
			t.Fatal("expected error from ParseNodeShape")
		}

		wrapped := fmt.Errorf("node config: %w", err)

		if !errors.Is(wrapped, ErrInvalidNodeShape) {
			t.Errorf("errors.Is(wrapped, ErrInvalidNodeShape) = false; err=%v", wrapped)
		}
	})

	t.Run("ErrInvalidArrowType from ParseArrowType", func(t *testing.T) {
		t.Parallel()

		_, err := ParseArrowType("bogus")
		if err == nil {
			t.Fatal("expected error from ParseArrowType")
		}

		wrapped := fmt.Errorf("edge config: %w", err)

		if !errors.Is(wrapped, ErrInvalidArrowType) {
			t.Errorf("errors.Is(wrapped, ErrInvalidArrowType) = false; err=%v", wrapped)
		}
	})

	t.Run("ErrInvalidConstraint from ParseConstraint", func(t *testing.T) {
		t.Parallel()

		_, err := ParseConstraint("bogus")
		if err == nil {
			t.Fatal("expected error from ParseConstraint")
		}

		wrapped := fmt.Errorf("layout constraint: %w", err)

		if !errors.Is(wrapped, ErrInvalidConstraint) {
			t.Errorf("errors.Is(wrapped, ErrInvalidConstraint) = false; err=%v", wrapped)
		}
	})

	t.Run("ErrInvalidTextTransform from ParseTextTransform", func(t *testing.T) {
		t.Parallel()

		_, err := ParseTextTransform("bogus")
		if err == nil {
			t.Fatal("expected error from ParseTextTransform")
		}

		wrapped := fmt.Errorf("label style: %w", err)

		if !errors.Is(wrapped, ErrInvalidTextTransform) {
			t.Errorf("errors.Is(wrapped, ErrInvalidTextTransform) = false; err=%v", wrapped)
		}
	})

	t.Run("sentinels are distinct", func(t *testing.T) {
		t.Parallel()

		_, dirErr := ParseDirection("bogus")
		if errors.Is(dirErr, ErrInvalidNodeShape) {
			t.Errorf("direction error should not match ErrInvalidNodeShape")
		}

		_, shapeErr := ParseNodeShape("bogus")
		if errors.Is(shapeErr, ErrInvalidDirection) {
			t.Errorf("shape error should not match ErrInvalidDirection")
		}
	})
}
