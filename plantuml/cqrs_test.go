package plantuml

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

	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("a", "Alpha")).
		AddNode(*output.NewGraphNode("b", "Beta")).
		AddEdge(*output.NewGraphEdge("a", "b")).
		Build()

	var buf bytes.Buffer

	if err := Write(&buf, g); err != nil {
		t.Fatalf("Write returned unexpected error on success: %v", err)
	}

	if !strings.Contains(buf.String(), "@startuml") {
		t.Errorf("Write output missing diagram start, got %q", buf.String())
	}
}

// TestWrite_PropagatesWriterError verifies Write surfaces io.Writer failures.
func TestWrite_PropagatesWriterError(t *testing.T) {
	t.Parallel()

	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("a", "Alpha")).
		Build()

	err := Write(&testhelpers.ErrorWriter{}, g)
	if err == nil {
		t.Fatal("expected error from failing writer, got nil")
	}

	if !errors.Is(err, testhelpers.ErrWrite) {
		t.Errorf("error should wrap testhelpers.ErrWrite, got %v", err)
	}
}
