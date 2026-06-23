package tui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

const chromeLinesAboveTree = 5

// percentScale is the maximum progress value; progress is reported on a 0–100 scale,
// so 100.0 means fully complete. Used both to normalize the progress bar fill and
// to detect workflow completion.
const percentScale = 100.0

// ============================================================================
// BUBBLE TEA MODEL IMPLEMENTATION
// ============================================================================
// tickCmd creates a recurring tick command for the progress model.
func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init implements tea.Model interface.
func (m *ProgressModel) Init() tea.Cmd {
	return tickCmd()
}

// Update implements tea.Model interface with unified state management.
func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)
	case tea.MouseWheelMsg:
		return m.handleMouseWheel(msg)
	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)
	case progressUpdateMsg:
		return m.handleProgressUpdate(msg)
	case stepUpdateMsg:
		return m.handleStepUpdate(msg)
	case errorMsg:
		return m.handleError(msg)
	case stateTransitionMsg:
		return m.handleStateTransition(msg)
	case tickMsg:
		return m.handleTick(msg)
	}

	return m, nil
}

func (m *ProgressModel) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	return m, nil
}

func (m *ProgressModel) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		if m.cancelFunc != nil {
			m.cancelFunc()
		}

		return m, tea.Quit
	case "q":
		return m, tea.Quit
	case "up", "k":
		m.scrollUp(1)
	case "down", "j":
		m.scrollDown(1)
	case "pgup":
		m.scrollUp(m.height / 2)
	case "pgdown":
		m.scrollDown(m.height / 2)
	case "home", "g":
		m.scrollOffset = 0
	case "end", "G":
		m.scrollToBottom()
	case "?":
		m.showHelp = !m.showHelp
	}

	return m, nil
}

func (m *ProgressModel) scrollUp(lines int) {
	if m.scrollOffset > lines {
		m.scrollOffset -= lines
	} else {
		m.scrollOffset = 0
	}
}

func (m *ProgressModel) scrollDown(lines int) {
	m.scrollOffset += lines
}

// scrollToBottomSentinel is the offset value set by scrollToBottom; the
// render path clamps it to totalLines-m.height.
const scrollToBottomSentinel = 1 << 30

func (m *ProgressModel) scrollToBottom() {
	m.scrollOffset = scrollToBottomSentinel
}

func (m *ProgressModel) handleMouseWheel(msg tea.MouseWheelMsg) (tea.Model, tea.Cmd) {
	mouse := msg.Mouse()
	switch mouse.Button {
	case ansi.MouseWheelUp:
		m.scrollUp(3)
	case ansi.MouseWheelDown:
		m.scrollDown(3)
	}

	return m, nil
}

func (m *ProgressModel) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	mouse := msg.Mouse()
	if mouse.Button != tea.MouseLeft {
		return m, nil
	}

	if m.dependencyTree == nil || len(m.visibleEntries) == 0 {
		return m, nil
	}

	// Map the physical Y coordinate to a visibleEntries index. This assumes
	// one terminal line per entry (real node or collapse marker), which holds
	// because the tree is rendered via RenderVisibleEntry which truncates long
	// lines to prevent wrapping.
	treeLine := mouse.Y - m.treeStartLine - chromeLinesAboveTree
	if treeLine < 0 || treeLine >= len(m.visibleEntries) {
		m.selectedNode = ""
		return m, nil
	}

	entry := m.visibleEntries[treeLine]

	// Collapse-marker lines (entry.Node == nil) are not selectable: clicking
	// one clears any current selection instead of dereferencing a nil node.
	if entry.Node == nil {
		m.selectedNode = ""

		return m, nil
	}

	nodeID := entry.Node.ID
	if m.selectedNode == nodeID {
		m.selectedNode = ""
	} else {
		m.selectedNode = nodeID
	}

	return m, nil
}

func (m *ProgressModel) handleProgressUpdate(msg progressUpdateMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.canAcceptUpdates() {
		return m, nil
	}

	m.lastUpdate = time.Now()

	switch msg.Type {
	case progressUpdate:
		m.currentProgress = msg.Progress
		if msg.Progress >= percentScale && m.workflowState.canTransitionTo(workflowStateCompleted) {
			m.workflowState = workflowStateCompleted
		}
	case messageUpdate:
		m.currentMessage = msg.Message
	}

	return m, nil
}

func (m *ProgressModel) handleTick(msg tickMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.canAcceptTicks() {
		return m, nil
	}

	m.lastUpdate = time.Time(msg)
	if m.nomSubscriber != nil && m.nomSubscriber.IsEnabled() {
		m.syncNOMSubscriber()
	}

	return m, tickCmd()
}

func (m *ProgressModel) syncNOMSubscriber() {
	m.dependencyTree = m.nomSubscriber.GetDependencyTree()
	if m.nomSubscriber.IsWorkflowRunning() {
		if m.workflowState == workflowStateIdle {
			m.workflowState = workflowStateRunning
			m.startTime = m.nomSubscriber.GetStartTime()
		}
	} else if m.workflowState == workflowStateRunning {
		m.updateWorkflowCompletionState()
	}
}

func (m *ProgressModel) updateWorkflowCompletionState() {
	counts := m.nomSubscriber.GetActivityCounts()
	if counts.Failed > 0 {
		m.workflowState = workflowStateErrored
	} else if counts.Running == 0 && counts.Completed > 0 {
		m.workflowState = workflowStateCompleted
	}
}

// handleStepUpdate processes a step-based progress update on the TUI goroutine.
// It creates a new step or updates an existing step with a matching message.
func (m *ProgressModel) handleStepUpdate(msg stepUpdateMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.canAcceptUpdates() {
		return m, nil
	}

	m.lastUpdate = time.Now()

	for i := range m.steps {
		if m.steps[i].Message == msg.Message {
			m.steps[i].Current = msg.Current
			m.steps[i].Total = msg.Total
			m.steps[i].Message = msg.Message

			if msg.Total > 0 && msg.Current >= msg.Total && m.steps[i].CompletedAt == nil {
				now := time.Now()
				m.steps[i].CompletedAt = &now
			}

			return m, nil
		}
	}

	m.steps = append(m.steps, progressStep{
		Current:   msg.Current,
		Total:     msg.Total,
		Message:   msg.Message,
		StartTime: time.Now(),
	})

	return m, nil
}

// handleError transitions to Errored state and displays the error message.
func (m *ProgressModel) handleError(msg errorMsg) (tea.Model, tea.Cmd) {
	if m.workflowState.canTransitionTo(workflowStateErrored) {
		m.workflowState = workflowStateErrored
		m.currentMessage = fmt.Sprintf("Error: %v", msg.Err)
	}

	return m, nil
}

// handleStateTransition applies a validated workflow state transition.
func (m *ProgressModel) handleStateTransition(msg stateTransitionMsg) (tea.Model, tea.Cmd) {
	if m.workflowState.canTransitionTo(msg.NewState) {
		m.workflowState = msg.NewState
	}

	return m, nil
}
