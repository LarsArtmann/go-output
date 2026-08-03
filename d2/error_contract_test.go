package d2

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestTypedErrors_AsType_ThroughWrapping(t *testing.T) {
	t.Parallel()

	t.Run("InvalidDirectionError from ParseDirection", func(t *testing.T) {
		t.Parallel()

		_, err := ParseDirection("bogus")
		if err == nil {
			t.Fatal("expected error from ParseDirection")
		}

		wrapped := fmt.Errorf("diagram config: %w", err)

		extracted, ok := errors.AsType[*InvalidDirectionError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidDirectionError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(directionValues) {
			t.Errorf("Allowed length = %d, want %d", len(extracted.Allowed), len(directionValues))
		}
	})

	t.Run("InvalidNodeShapeError from ParseNodeShape", func(t *testing.T) {
		t.Parallel()

		_, err := ParseNodeShape("bogus")
		if err == nil {
			t.Fatal("expected error from ParseNodeShape")
		}

		wrapped := fmt.Errorf("node config: %w", err)

		extracted, ok := errors.AsType[*InvalidNodeShapeError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidNodeShapeError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(nodeShapeValues) {
			t.Errorf("Allowed length = %d, want %d", len(extracted.Allowed), len(nodeShapeValues))
		}
	})

	t.Run("InvalidArrowTypeError from ParseArrowType", func(t *testing.T) {
		t.Parallel()

		_, err := ParseArrowType("bogus")
		if err == nil {
			t.Fatal("expected error from ParseArrowType")
		}

		wrapped := fmt.Errorf("edge config: %w", err)

		extracted, ok := errors.AsType[*InvalidArrowTypeError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidArrowTypeError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(arrowTypeValues) {
			t.Errorf("Allowed length = %d, want %d", len(extracted.Allowed), len(arrowTypeValues))
		}
	})

	t.Run("InvalidConstraintError from ParseConstraint", func(t *testing.T) {
		t.Parallel()

		_, err := ParseConstraint("bogus")
		if err == nil {
			t.Fatal("expected error from ParseConstraint")
		}

		wrapped := fmt.Errorf("layout constraint: %w", err)

		extracted, ok := errors.AsType[*InvalidConstraintError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidConstraintError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(allConstraints) {
			t.Errorf("Allowed length = %d, want %d", len(extracted.Allowed), len(allConstraints))
		}
	})

	t.Run("InvalidTextTransformError from ParseTextTransform", func(t *testing.T) {
		t.Parallel()

		_, err := ParseTextTransform("bogus")
		if err == nil {
			t.Fatal("expected error from ParseTextTransform")
		}

		wrapped := fmt.Errorf("label style: %w", err)

		extracted, ok := errors.AsType[*InvalidTextTransformError](wrapped)
		if !ok {
			t.Fatalf("errors.AsType[*InvalidTextTransformError] failed; err=%v", wrapped)
		}

		if extracted.Value != "bogus" {
			t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
		}

		if len(extracted.Allowed) != len(textTransformValues) {
			t.Errorf("Allowed length = %d, want %d", len(extracted.Allowed), len(textTransformValues))
		}
	})

	t.Run("typed errors are distinct", func(t *testing.T) {
		t.Parallel()

		_, dirErr := ParseDirection("bogus")
		if _, ok := errors.AsType[*InvalidNodeShapeError](dirErr); ok {
			t.Errorf("direction error should not match *InvalidNodeShapeError")
		}

		_, shapeErr := ParseNodeShape("bogus")
		if _, ok := errors.AsType[*InvalidDirectionError](shapeErr); ok {
			t.Errorf("shape error should not match *InvalidDirectionError")
		}
	})

	t.Run("error messages include allowed values", func(t *testing.T) {
		t.Parallel()

		err := &InvalidDirectionError{Value: "bogus", Allowed: directionValues}
		msg := err.Error()
		if !strings.Contains(msg, "(allowed:") {
			t.Errorf("error message should include allowed values list; got: %s", msg)
		}
		if !strings.Contains(msg, "bogus") {
			t.Errorf("error message should include the invalid value; got: %s", msg)
		}
	})
}
