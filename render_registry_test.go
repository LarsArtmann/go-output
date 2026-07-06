package output

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestRegisterTableMarshaler_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func() {
			defer wg.Done()

			RegisterTableMarshaler(
				Format("race-test-"+strconv.Itoa(i)),
				func(w io.Writer, data *Table, opts RenderOptions) error { return nil },
			)
		}()

		go func() {
			defer wg.Done()

			_, _ = getTableMarshaler(Format("race-test-" + strconv.Itoa(i)))
		}()
	}

	wg.Wait()
}

func TestRegisteredTableMarshalFormats(t *testing.T) {
	t.Parallel()

	formats := RegisteredTableMarshalFormats()

	// Root in isolation does not register any Table renderers via init().
	// Sub-modules (markdown, tree, delimited, ...) self-register when imported.
	// Cross-module registration is tested in integration/.
	// Just verify the call is safe and returns no duplicates.
	seen := make(map[Format]bool, len(formats))
	for _, f := range formats {
		if seen[f] {
			t.Errorf("duplicate format in RegisteredTableMarshalFormats: %q", f)
		}

		seen[f] = true
	}
}

func TestRegisteredUnknownFormats(t *testing.T) {
	t.Parallel()

	formats := RegisteredUnknownFormats()

	// Root does not register any any-data renderers in init(); this should
	// return an empty slice (or whatever sub-modules have registered).
	// Just verify the call is safe and returns no duplicates.
	seen := make(map[Format]bool, len(formats))
	for _, f := range formats {
		if seen[f] {
			t.Errorf("duplicate format in RegisteredUnknownFormats: %q", f)
		}

		seen[f] = true
	}
}

func TestRenderUnknown_UnsupportedFormatError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := RenderUnknown("data", Format("unknown"), RenderOptions{Writer: &buf})
	if err == nil {
		t.Fatal("expected UnsupportedFormatError")
	}

	var unsupportedErr *UnsupportedFormatError
	if !errors.As(err, &unsupportedErr) {
		t.Errorf("expected UnsupportedFormatError, got %T: %v", err, err)
	}
}

func TestUnsupportedFormatErrorMessage(t *testing.T) {
	t.Parallel()

	err := &UnsupportedFormatError{Format: FormatJSON}

	msg := err.Error()
	if !strings.Contains(msg, "json") {
		t.Errorf("error message should contain format name, got: %q", msg)
	}
}
