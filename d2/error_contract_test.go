package d2

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// assertWrappedTypedError wraps a parse error and asserts errors.AsType extracts
// the expected typed error through fmt.Errorf wrapping. Extracted as a generic
// helper so the test function's cognitive complexity stays low.
func assertWrappedTypedError[T error](
	t *testing.T,
	name string,
	parseErr func() error,
	wrapCtx string,
	checkExtracted func(t *testing.T, extracted T),
) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		t.Parallel()

		err := parseErr()
		if err == nil {
			t.Fatal("expected error")
		}

		wrapped := fmt.Errorf("%s: %w", wrapCtx, err)

		extracted, ok := errors.AsType[T](wrapped)
		if !ok {
			t.Fatalf("errors.AsType failed; err=%v", wrapped)
		}

		checkExtracted(t, extracted)
	})
}

func TestTypedErrors_AsType_ThroughWrapping(t *testing.T) {
	t.Parallel()

	assertWrappedTypedError[*InvalidDirectionError](t,
		"InvalidDirectionError from ParseDirection",
		func() error { _, err := ParseDirection("bogus"); return err },
		"diagram config",
		func(t *testing.T, e *InvalidDirectionError) {
			if e.Value != "bogus" {
				t.Errorf("Value = %q, want %q", e.Value, "bogus")
			}

			if len(e.Allowed) != len(directionValues) {
				t.Errorf("Allowed length = %d, want %d", len(e.Allowed), len(directionValues))
			}
		},
	)

	assertWrappedTypedError[*InvalidNodeShapeError](t,
		"InvalidNodeShapeError from ParseNodeShape",
		func() error { _, err := ParseNodeShape("bogus"); return err },
		"node config",
		func(t *testing.T, e *InvalidNodeShapeError) {
			if e.Value != "bogus" {
				t.Errorf("Value = %q, want %q", e.Value, "bogus")
			}

			if len(e.Allowed) != len(nodeShapeValues) {
				t.Errorf("Allowed length = %d, want %d", len(e.Allowed), len(nodeShapeValues))
			}
		},
	)

	assertWrappedTypedError[*InvalidArrowTypeError](t,
		"InvalidArrowTypeError from ParseArrowType",
		func() error { _, err := ParseArrowType("bogus"); return err },
		"edge config",
		func(t *testing.T, e *InvalidArrowTypeError) {
			if e.Value != "bogus" {
				t.Errorf("Value = %q, want %q", e.Value, "bogus")
			}

			if len(e.Allowed) != len(arrowTypeValues) {
				t.Errorf("Allowed length = %d, want %d", len(e.Allowed), len(arrowTypeValues))
			}
		},
	)

	assertWrappedTypedError[*InvalidConstraintError](t,
		"InvalidConstraintError from ParseConstraint",
		func() error { _, err := ParseConstraint("bogus"); return err },
		"layout constraint",
		func(t *testing.T, e *InvalidConstraintError) {
			if e.Value != "bogus" {
				t.Errorf("Value = %q, want %q", e.Value, "bogus")
			}

			if len(e.Allowed) != len(allConstraints) {
				t.Errorf("Allowed length = %d, want %d", len(e.Allowed), len(allConstraints))
			}
		},
	)

	assertWrappedTypedError[*InvalidTextTransformError](t,
		"InvalidTextTransformError from ParseTextTransform",
		func() error { _, err := ParseTextTransform("bogus"); return err },
		"label style",
		func(t *testing.T, e *InvalidTextTransformError) {
			if e.Value != "bogus" {
				t.Errorf("Value = %q, want %q", e.Value, "bogus")
			}

			if len(e.Allowed) != len(textTransformValues) {
				t.Errorf("Allowed length = %d, want %d", len(e.Allowed), len(textTransformValues))
			}
		},
	)

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
