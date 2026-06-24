package nom

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestRenderCompletion_Success(t *testing.T) {
	var buf bytes.Buffer

	sub := NewNOMStyleSubscriber()
	r := NewInlineRenderer(sub, &buf, 10)
	r.SetAppName("BuildFlow")

	r.RenderCompletion(CompletionResult{
		Success:     true,
		Elapsed:     12340 * time.Millisecond,
		TotalSteps:  42,
		FailedSteps: 0,
	})

	output := buf.String()
	if !strings.Contains(output, "✓") {
		t.Errorf("expected ✓ in success output, got: %s", output)
	}

	if !strings.Contains(output, "BuildFlow") {
		t.Errorf("expected app name in output, got: %s", output)
	}

	if !strings.Contains(output, "42 steps") {
		t.Errorf("expected step count in output, got: %s", output)
	}
}

func TestRenderCompletion_Failure(t *testing.T) {
	var buf bytes.Buffer

	sub := NewNOMStyleSubscriber()
	r := NewInlineRenderer(sub, &buf, 10)
	r.SetAppName("Workflow")

	r.RenderCompletion(CompletionResult{
		Success:     false,
		Elapsed:     45600 * time.Millisecond,
		TotalSteps:  42,
		FailedSteps: 3,
	})

	output := buf.String()
	if !strings.Contains(output, "✗") {
		t.Errorf("expected ✗ in failure output, got: %s", output)
	}

	if !strings.Contains(output, "3/42 steps failed") {
		t.Errorf("expected failure detail in output, got: %s", output)
	}
}

func TestRenderCompletion_RestoresCursor(t *testing.T) {
	var buf bytes.Buffer

	sub := NewNOMStyleSubscriber()
	r := NewInlineRenderer(sub, &buf, 10)

	// Simulate that cursor was hidden during rendering.
	r.renderMu.Lock()
	r.prevLines = 5
	r.renderMu.Unlock()

	r.RenderCompletion(CompletionResult{
		Success:    true,
		Elapsed:    1 * time.Second,
		TotalSteps: 1,
	})

	output := buf.String()

	// The ShowCursor escape code should be present since prevLines was > 0.
	if !strings.Contains(output, "\x1b[?25h") {
		t.Errorf("expected cursor to be restored (ShowCursor), got: %q", output)
	}
}
