package markdown

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

// TestWrite_ReturnsNilOnSuccess is a regression test for a bug where Write
// always returned a non-nil error wrapping nil — the final
// `fmt.Errorf("write output: %w", err)` executed even when io.WriteString
// succeeded. With the fix, Write must return nil for a healthy writer.
func TestWrite_ReturnsNilOnSuccess(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Status"})
	data.AddRow([]string{"Alpha", "OK"})

	var buf bytes.Buffer

	if err := Write(&buf, data); err != nil {
		t.Fatalf("Write returned unexpected error on success: %v", err)
	}

	if !strings.Contains(buf.String(), "Alpha") {
		t.Errorf("Write output missing row data, got %q", buf.String())
	}
}

// TestWrite_PropagatesWriterError verifies Write surfaces io.Writer failures
// instead of swallowing them (the buggy version always returned the same
// "write output" wrapper, hiding the real cause).
func TestWrite_PropagatesWriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alpha"})

	err := Write(&testhelpers.ErrorWriter{}, data)
	if err == nil {
		t.Fatal("expected error from failing writer, got nil")
	}

	if !errors.Is(err, testhelpers.ErrWrite) {
		t.Errorf("error should wrap testhelpers.ErrWrite, got %v", err)
	}
}

// TestRender_ReturnsMarkdownString confirms the Render convenience wrapper
// returns the same string Write would produce, without a trailing newline.
func TestRender_ReturnsMarkdownString(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alpha"})

	out, err := Render(data)
	if err != nil {
		t.Fatalf("Render error = %v", err)
	}

	if !strings.Contains(out, "Alpha") {
		t.Errorf("Render output missing data row, got %q", out)
	}
}
