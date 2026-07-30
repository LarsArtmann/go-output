package output

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestTypedError_InvalidShapeError_AsType(t *testing.T) {
	t.Parallel()

	_, err := ParseShape("bogus")
	if err == nil {
		t.Fatal("expected error from ParseShape")
	}

	wrapped := fmt.Errorf("parsing config: %w", err)

	extracted, ok := errors.AsType[*InvalidShapeError](wrapped)
	if !ok {
		t.Fatalf("errors.AsType[*InvalidShapeError] failed; err=%v", wrapped)
	}

	if extracted.Value != "bogus" {
		t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
	}

	if len(extracted.Allowed) != len(AllShapes) {
		t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllShapes))
	}
}

func TestTypedError_InvalidColorModeError_AsType(t *testing.T) {
	t.Parallel()

	_, err := ParseColorMode("bogus")
	if err == nil {
		t.Fatal("expected error from ParseColorMode")
	}

	wrapped := fmt.Errorf("loading theme: %w", err)

	extracted, ok := errors.AsType[*InvalidColorModeError](wrapped)
	if !ok {
		t.Fatalf("errors.AsType[*InvalidColorModeError] failed; err=%v", wrapped)
	}

	if extracted.Value != "bogus" {
		t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
	}

	if len(extracted.Allowed) != len(AllColorModes) {
		t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllColorModes))
	}
}

func TestTypedError_InvalidFormatError_AsType(t *testing.T) {
	t.Parallel()

	_, err := ParseFormat("bogus")
	if err == nil {
		t.Fatal("expected error from ParseFormat")
	}

	wrapped := fmt.Errorf("cli flag: %w", err)

	extracted, ok := errors.AsType[*InvalidFormatError](wrapped)
	if !ok {
		t.Fatalf("errors.AsType[*InvalidFormatError] failed; err=%v", wrapped)
	}

	if extracted.Value != "bogus" {
		t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
	}

	if len(extracted.Allowed) != len(AllFormats) {
		t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllFormats))
	}
}

func TestTypedError_InvalidLineStyleError_AsType(t *testing.T) {
	t.Parallel()

	_, err := ParseLineStyle("bogus")
	if err == nil {
		t.Fatal("expected error from ParseLineStyle")
	}

	wrapped := fmt.Errorf("graph config: %w", err)

	extracted, ok := errors.AsType[*InvalidLineStyleError](wrapped)
	if !ok {
		t.Fatalf("errors.AsType[*InvalidLineStyleError] failed; err=%v", wrapped)
	}

	if extracted.Value != "bogus" {
		t.Errorf("Value = %q, want %q", extracted.Value, "bogus")
	}

	if len(extracted.Allowed) != len(AllLineStyles) {
		t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllLineStyles))
	}
}

func TestTypedError_InvalidNodeShapeError_AsType(t *testing.T) {
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

	if len(extracted.Allowed) != len(AllNodeShapes) {
		t.Errorf("Allowed len = %d, want %d", len(extracted.Allowed), len(AllNodeShapes))
	}
}

func TestTypedError_UnsupportedFormatError_AsType(t *testing.T) {
	t.Parallel()

	data := NewTable([]string{"A"})

	err := RenderTable(data, Format("nonexistent"), RenderOptions{})
	if err == nil {
		t.Fatal("expected error from RenderTable with unsupported format")
	}

	wrapped := fmt.Errorf("render pipeline: %w", err)

	extracted, ok := errors.AsType[*UnsupportedFormatError](wrapped)
	if !ok {
		t.Fatalf("errors.AsType[*UnsupportedFormatError] failed; err=%v", wrapped)
	}

	if extracted.Format != Format("nonexistent") {
		t.Errorf("Format = %v, want %v", extracted.Format, Format("nonexistent"))
	}
}

func TestTypedErrors_ErrorMessages_IncludeAllowedValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			"InvalidShapeError",
			&InvalidShapeError{Value: "x", Allowed: AllShapes},
			"(allowed: table, tree, graph)",
		},
		{
			"InvalidColorModeError",
			&InvalidColorModeError{Value: "x", Allowed: AllColorModes},
			"(allowed: auto, always, never)",
		},
		{
			"InvalidFormatError",
			&InvalidFormatError{Value: "x", Allowed: AllFormats},
			"(allowed: ",
		},
		{
			"InvalidLineStyleError",
			&InvalidLineStyleError{Value: "x", Allowed: AllLineStyles},
			"(allowed: solid, dashed, dotted)",
		},
		{
			"InvalidNodeShapeError",
			&InvalidNodeShapeError{Value: "x", Allowed: AllNodeShapes},
			"(allowed: box, ellipse, diamond, circle, cylinder, hexagon, parallelogram)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			msg := tt.err.Error()
			if !strings.Contains(msg, tt.contains) {
				t.Errorf("Error() = %q, expected to contain %q", msg, tt.contains)
			}
		})
	}
}
