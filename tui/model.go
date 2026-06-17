package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
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
	if m.selectedNode == node.ActivityID {
		m.selectedNode = ""
	} else {
		m.selectedNode = node.ActivityID
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
	m.activities = m.nomSubscriber.GetActivities()

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
	running, completed, failed, _ := m.nomSubscriber.GetActivityCounts()
	if failed > 0 {
		m.workflowState = WorkflowStateErrored
	} else if running == 0 && completed > 0 {
		m.workflowState = WorkflowStateCompleted
	}
}
