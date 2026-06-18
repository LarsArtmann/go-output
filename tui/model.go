package tui

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/larsartmann/go-output/nom"
)

const chromeLinesAboveTree = 5

// ============================================================================
// BUBBLE TEA MODEL IMPLEMENTATION
// ============================================================================
// TickCmd creates a recurring tick command for the progress model.
func TickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init implements tea.Model interface.
func (m *ProgressModel) Init() tea.Cmd {
	return TickCmd()
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
	case ProgressUpdateMsg:
		return m.handleProgressUpdate(msg)
	case StepUpdateMsg:
		return m.handleStepUpdate(msg)
	case ErrorMsg:
		return m.handleError(msg)
	case StateTransitionMsg:
		return m.handleStateTransition(msg)
	case TickMsg:
		return m.handleTick(msg)
	case CancelMsg:
		return m, tea.Quit
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

	if m.dependencyTree == nil || len(m.visibleNodes) == 0 {
		return m, nil
	}

	// Map the physical Y coordinate to a visibleNodes index. This assumes
	// one terminal line per node, which holds because the tree is rendered
	// via RenderWithWidth(maxHeight, maxWidth) which truncates long lines
	// to prevent wrapping.
	treeLine := mouse.Y - m.treeStartLine - chromeLinesAboveTree + m.scrollOffset
	if treeLine < 0 || treeLine >= len(m.visibleNodes) {
		m.selectedNode = ""
		return m, nil
	}

	node := m.visibleNodes[treeLine]
	nodeID := nom.ActivityID(node.ID.Get())
	if m.selectedNode == nodeID {
		m.selectedNode = ""
	} else {
		m.selectedNode = nodeID
	}

	return m, nil
}

func (m *ProgressModel) handleProgressUpdate(msg ProgressUpdateMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.CanAcceptUpdates() {
		return m, nil
	}

	m.lastUpdate = time.Now()

	switch msg.Type {
	case ProgressUpdate:
		m.currentProgress = msg.Progress
		if msg.Progress >= 100.0 && m.workflowState.CanTransitionTo(WorkflowStateCompleted) {
			m.workflowState = WorkflowStateCompleted
		}
	case MessageUpdate:
		m.currentMessage = msg.Message
	case StepUpdate:
		// Step updates are handled in ReportStep method
	}

	return m, nil
}

func (m *ProgressModel) handleTick(msg TickMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.CanAcceptTicks() {
		return m, nil
	}

	m.lastUpdate = time.Time(msg)
	if m.nomSubscriber != nil && m.nomSubscriber.IsEnabled() {
		m.syncNOMSubscriber()
	}

	return m, TickCmd()
}

func (m *ProgressModel) syncNOMSubscriber() {
	m.nomSubscriber.UpdateRunningActivityElapsed()
	m.nomSubscriber.SyncActivityTimingToTree()

	m.dependencyTree = m.nomSubscriber.GetDependencyTree()
	if m.nomSubscriber.IsWorkflowRunning() {
		if m.workflowState == WorkflowStateIdle {
			m.workflowState = WorkflowStateRunning
			m.startTime = m.nomSubscriber.GetStartTime()
		}
	} else if m.workflowState == WorkflowStateRunning {
		m.updateWorkflowCompletionState()
	}
}

func (m *ProgressModel) updateWorkflowCompletionState() {
	counts := m.nomSubscriber.GetActivityCounts()
	if counts.Failed > 0 {
		m.workflowState = WorkflowStateErrored
	} else if counts.Running == 0 && counts.Completed > 0 {
		m.workflowState = WorkflowStateCompleted
	}
}

// handleStepUpdate processes a step-based progress update on the TUI goroutine.
// It creates a new step or updates an existing matching/active step.
func (m *ProgressModel) handleStepUpdate(msg StepUpdateMsg) (tea.Model, tea.Cmd) {
	if !m.workflowState.CanAcceptUpdates() {
		return m, nil
	}

	m.lastUpdate = time.Now()

	for i := range m.steps {
		if m.steps[i].Message == msg.Message || m.steps[i].IsActive() {
			m.steps[i].Current = msg.Current
			m.steps[i].Total = msg.Total
			m.steps[i].Message = msg.Message

			if msg.Current >= msg.Total && m.steps[i].CompletedAt == nil {
				now := time.Now()
				m.steps[i].CompletedAt = &now
			}

			return m, nil
		}
	}

	m.steps = append(m.steps, ProgressStep{
		Current:   msg.Current,
		Total:     msg.Total,
		Message:   msg.Message,
		StartTime: time.Now(),
	})

	return m, nil
}

// handleError transitions to Errored state and displays the error message.
func (m *ProgressModel) handleError(msg ErrorMsg) (tea.Model, tea.Cmd) {
	if m.workflowState.CanTransitionTo(WorkflowStateErrored) {
		m.workflowState = WorkflowStateErrored
		m.currentMessage = fmt.Sprintf("Error: %v", msg.Err)
	}

	return m, nil
}

// handleStateTransition applies a validated workflow state transition.
func (m *ProgressModel) handleStateTransition(msg StateTransitionMsg) (tea.Model, tea.Cmd) {
	if m.workflowState.CanTransitionTo(msg.NewState) {
		m.workflowState = msg.NewState
	}

	return m, nil
}
