package output

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestRegisterTableDataRenderer_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func() {
			defer wg.Done()

			RegisterTableDataRenderer(
				Format("race-test-"+strconv.Itoa(i)),
				func(w io.Writer, data *TableData, opts RenderOptions) error { return nil },
			)
		}()

		go func() {
			defer wg.Done()

			_, _ = getTableDataRenderer(Format("race-test-" + strconv.Itoa(i)))
		}()
	}

	wg.Wait()
}

func TestRegisteredTableDataFormats(t *testing.T) {
	t.Parallel()

	formats := RegisteredTableDataFormats()

	// Root init() registers Markdown and Tree as TableData renderers.
	for _, exp := range []Format{FormatMarkdown, FormatTree} {
		if !slices.Contains(formats, exp) {
			t.Errorf("expected format %q to be registered, but it was not. Registered: %v", exp, formats)
		}
	}
}

func TestRegisteredAnyDataFormats(t *testing.T) {
	t.Parallel()

	formats := RegisteredAnyDataFormats()

	// Root does not register any any-data renderers in init(); this should
	// return an empty slice (or whatever sub-modules have registered).
	// Just verify the call is safe and returns no duplicates.
	seen := make(map[Format]bool, len(formats))
	for _, f := range formats {
		if seen[f] {
			t.Errorf("duplicate format in RegisteredAnyDataFormats: %q", f)
		}

		seen[f] = true
	}
}

func TestRenderAnyData_UnsupportedFormatError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := RenderAnyData("data", Format("unknown"), RenderOptions{Writer: &buf})
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
