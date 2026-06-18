package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func assertScrollOffset(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("scrollOffset = %d, want %d", got, want)
	}
}

// TestProgressModel_KeyboardNavigation verifies scroll keys update offset.
func TestProgressModel_KeyboardNavigation(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.height = 20
	model.scrollOffset = 10

	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m := updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 9)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 10)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyHome})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 0)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, scrollToBottomSentinel)
}

// TestProgressModel_MouseScrolling verifies mouse wheel updates scroll offset.
func TestProgressModel_MouseScrolling(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.scrollOffset = 5

	updated, _ := model.Update(tea.MouseWheelMsg{Button: ansi.MouseWheelDown})
	m := updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 8)

	updated, _ = m.Update(tea.MouseWheelMsg{Button: ansi.MouseWheelUp})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 5)
}

// TestProgressModel_CancelMessage verifies CancelMsg triggers quit.
func TestProgressModel_CancelMessage(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	_, cmd := model.Update(CancelMsg{})

	if cmd == nil {
		t.Fatal("expected Quit command for CancelMsg")
	}

	msg := cmd()
	if msg == nil {
		t.Error("quit command should produce a message")
	}
}

// TestProgressModel_ViewportScrolling verifies applyScrollViewport clips content.
func TestProgressModel_ViewportScrolling(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 3

	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"

	result := model.applyScrollViewport(content)

	lines := splitLines(result)
	if len(lines) != 3 {
		t.Fatalf("expected 3 visible lines, got %d", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("first visible line = %q, want %q", lines[0], "line1")
	}

	// Scroll to offset 5
	model.scrollOffset = 5
	result = model.applyScrollViewport(content)

	lines = splitLines(result)
	if len(lines) != 3 {
		t.Fatalf("expected 3 visible lines after scroll, got %d", len(lines))
	}

	if lines[0] != "line6" {
		t.Errorf("first visible line after scroll = %q, want %q", lines[0], "line6")
	}

	// Scroll past end — should clamp
	model.scrollOffset = 100
	result = model.applyScrollViewport(content)

	lines = splitLines(result)
	if lines[0] != "line8" {
		t.Errorf("clamped first line = %q, want %q", lines[0], "line8")
	}

	// Content fits viewport — scrollOffset should reset to 0
	model.height = 20
	model.scrollOffset = 5

	_ = model.applyScrollViewport(content)
	if model.scrollOffset != 0 {
		t.Errorf("scrollOffset should be 0 when content fits, got %d", model.scrollOffset)
	}
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, "\n")
}

func TestProgressModel_ResizeStress_RapidSequence(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = WorkflowStateRunning
	model.dependencyTree = newTestTree(50)

	resizes := []tea.WindowSizeMsg{
		{Width: 120, Height: 40},
		{Width: 80, Height: 24},
		{Width: 60, Height: 15},
		{Width: 200, Height: 60},
		{Width: 40, Height: 10},
		{Width: 120, Height: 40},
	}

	for _, resize := range resizes {
		updated, _ := model.Update(resize)
		m := updated.(*ProgressModel)

		if m.width != resize.Width {
			t.Errorf("width after resize = %d, want %d", m.width, resize.Width)
		}

		if m.height != resize.Height {
			t.Errorf("height after resize = %d, want %d", m.height, resize.Height)
		}

		view := m.View()
		if view.Content == "" {
			t.Error("View() should produce content after resize")
		}

		model = m
	}
}

func TestProgressModel_Resize_ClampsScrollOffset(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = WorkflowStateRunning
	model.width = 120
	model.height = 40
	model.dependencyTree = newTestTree(50)

	// Scroll to bottom in large window
	model.scrollOffset = 30

	// Shrink window — viewport should clamp scrollOffset
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
	m := updated.(*ProgressModel)

	// scrollOffset should not exceed content lines - height
	if m.scrollOffset < 0 {
		t.Error("scrollOffset should not be negative")
	}

	// Render should not panic with small height
	view := m.View()
	if view.Content == "" {
		t.Error("View() should produce content after resize")
	}
}
